package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"

	"gorm.io/gorm"
)

/*
	IMPORTANT: Most of the functions in this file do not check if the CodeSystem, CodeSystemVersion or Concept exists in the database.
	This is done for better performance (to avoid unnecessary queries).
	Therefore, the calling Function should check if the CodeSystem, CodeSystemVersion or Concept exists in the database.
*/

func (gq *GormQuery) SetCodeSystemVersionImported(codeSystemVersionId int32, imported bool) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		var codeSystemVersion models.CodeSystemVersion
		if err := tx.First(&codeSystemVersion, codeSystemVersionId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found.", codeSystemVersionId))
			default:
				return err
			}
		}

		codeSystemVersion.Imported = imported
		if err := tx.Save(&codeSystemVersion).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (gq *GormQuery) GetImportedNeighborVersionIds(codeSystemId int32, codeSystemVersionId int32) (uint32, *uint32, *uint32, error) {
	var codeSystemVersion models.CodeSystemVersion
	if err := gq.GetCodeSystemVersionQuery(&codeSystemVersion, codeSystemId, codeSystemVersionId); err != nil {
		return 0, nil, nil, err
	}

	versionId := codeSystemVersion.VersionID

	var beforeVersionId *uint32
	beforeVersionId = nil
	var afterVersionId *uint32
	afterVersionId = nil

	var beforeCodeSystemVersions []models.CodeSystemVersion
	if err := gq.Database.Where("code_system_id = ? AND version_id < ?", codeSystemId, codeSystemVersion.VersionID).Order("version_id DESC").Find(&beforeCodeSystemVersions).Error; err != nil {
		return versionId, nil, nil, err
	}

	if len(beforeCodeSystemVersions) > 0 {
		for _, beforeCodeSystemVersion := range beforeCodeSystemVersions {
			if beforeCodeSystemVersion.Imported {
				beforeVersionId = &beforeCodeSystemVersion.VersionID
				break
			}
		}
	}

	var afterCodeSystemVersions []models.CodeSystemVersion
	if err := gq.Database.Where("code_system_id = ? AND version_id > ?", codeSystemId, codeSystemVersion.VersionID).Order("version_id ASC").Find(&afterCodeSystemVersions).Error; err != nil {
		return versionId, nil, nil, err
	}

	if len(afterCodeSystemVersions) > 0 {
		for _, afterCodeSystemVersion := range afterCodeSystemVersions {
			if afterCodeSystemVersion.Imported {
				afterVersionId = &afterCodeSystemVersion.VersionID
				break
			}
		}
	}

	return versionId, beforeVersionId, afterVersionId, nil
}

func (gq GormQuery) CreateConceptQuery(concept *models.Concept) error {
	if err := gq.Database.Create(concept).Error; err != nil {
		return err
	}
	return nil
}

func (gq GormQuery) UpdateConceptQuery(concept *models.Concept) error {
	// concept.DisplaySearchVector should not be updated
	if err := gq.Database.Model(&models.Concept{}).Where("id = ?", concept.ID).Updates(map[string]interface{}{
		"code":                  concept.Code,
		"display":               concept.Display,
		"code_system_id":        concept.CodeSystemID,
		"description":           concept.Description,
		"status":                concept.Status,
		"valid_from_version_id": concept.ValidFromVersionID,
		"valid_to_version_id":   concept.ValidToVersionID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (gq GormQuery) getConceptsByCode(code string, codeSystemId int32) ([]models.Concept, error) {
	var concepts []models.Concept
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		query := tx.
			Preload("ValidFromVersion").
			Preload("ValidToVersion").
			Model(&models.Concept{}).
			Where("code_system_id = ?", codeSystemId).
			Where("code = ?", code)

		if err := query.Find(&concepts).Error; err != nil {
			return err
		} else if len(concepts) == 0 {
			return database.NewDBError(database.NotFound, fmt.Sprintf("Concept with code %s couldn't be found in CodeSystem with ID %d.", code, codeSystemId))
		}
		return nil
	})
	return concepts, err
}

func (gq GormQuery) GetNeighborConceptsQuery(code string, codeSystemId int32, versionId uint32, beforeVersionId *uint32, afterVersionId *uint32) (database.NeighborConcepts, error) {
	if concepts, err := gq.getConceptsByCode(code, codeSystemId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return database.NeighborConcepts{NeighborType: database.NeighborConceptsTypeNone}, nil
		default:
			return database.NeighborConcepts{}, err
		}
	} else {
		var beforeConcepts []models.Concept
		var afterConcepts []models.Concept
		var surroundingConcepts []models.Concept

		for _, concept := range concepts {
			if beforeVersionId != nil && concept.ValidToVersion.VersionID == *beforeVersionId {
				beforeConcepts = append(beforeConcepts, concept)
			} else if afterVersionId != nil && concept.ValidFromVersion.VersionID == *afterVersionId {
				afterConcepts = append(afterConcepts, concept)
			} else if concept.ValidFromVersion.VersionID < versionId && concept.ValidToVersion.VersionID > versionId {
				surroundingConcepts = append(surroundingConcepts, concept)
			}
		}

		if (len(beforeConcepts) > 1 || len(afterConcepts) > 1 || len(surroundingConcepts) > 1) || ((len(beforeConcepts) == 1 || len(afterConcepts) == 1) && len(surroundingConcepts) == 1) {
			return database.NeighborConcepts{}, database.NewDBError(database.InternalServerError, fmt.Sprintf("Invalid stored concepts found while getting Neighbors for Concept with code %s in CodeSystem with ID %d.", code, codeSystemId))
		} else if len(surroundingConcepts) == 1 {
			return database.NeighborConcepts{SurroundingConcept: &surroundingConcepts[0], NeighborType: database.NeighborConceptsTypeSurrounding}, nil
		} else if len(beforeConcepts) == 1 && len(afterConcepts) == 1 {
			return database.NeighborConcepts{BeforeConcept: &beforeConcepts[0], AfterConcept: &afterConcepts[0], NeighborType: database.NeighborConceptsTypeBeforeAndAfter}, nil
		} else if len(beforeConcepts) == 1 {
			return database.NeighborConcepts{BeforeConcept: &beforeConcepts[0], NeighborType: database.NeighborConceptsTypeBefore}, nil
		} else if len(afterConcepts) == 1 {
			return database.NeighborConcepts{AfterConcept: &afterConcepts[0], NeighborType: database.NeighborConceptsTypeAfter}, nil
		} else {
			return database.NeighborConcepts{NeighborType: database.NeighborConceptsTypeNone}, nil
		}
	}
}
