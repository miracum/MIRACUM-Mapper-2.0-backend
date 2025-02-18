package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"strings"

	"gorm.io/gorm"
)

func (gq *GormQuery) GetAllConceptsQuery(concepts *[]models.Concept, codeSystemId int32, pageSize int, offset int, sortBy string, sortOrder string, meaning string, code string) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		query := tx.
			Preload("ValidFromVersion").
			Preload("ValidToVersion").
			Model(&models.Concept{}).
			Where("code_system_id = ?", codeSystemId)

		// Add code condition if code is not empty
		if code != "" {
			query = query.Where("LOWER(code) LIKE LOWER(?)", code+"%")
		}

		// Add meaning condition if meaning is not empty
		if meaning != "" {
			formattedMeaning := strings.Join(strings.Fields(meaning), ":* & ") + ":*" // Adjust for partial matches
			query = query.Where("display_search_vector @@ to_tsquery('english', ?)", formattedMeaning)
			// Tests for similarity searches. These were very slow and therefore not used in the final implementation. I a search should be implemented which is not part of a autocomplete but is ok to take e.g a second to complete, similarity searches with the pg_trgm extension could be used here.

			// query = query.Where("display_search_vector @@ to_tsquery(?) OR similarity(display, ?) > 0.3", formattedMeaning, meaning)
			// query = query.Where("similarity(display, ?) > 0.8", meaning)
			// query = query.Select("*, similarity(display, ?) > set_limit(0.99) AS s", meaning).Order("s DESC")
		}

		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize)

		if err := query.Find(&concepts).Error; err != nil {
			return err
		} else if len(*concepts) == 0 {
			var codesystem models.CodeSystem
			if err := tx.First(&codesystem, codeSystemId).Error; err != nil {
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
			} else {
				*concepts = []models.Concept{}
				return nil
			}
		}
		return nil
	})
	return err
}

func (gq *GormQuery) GetAllConceptsByVersionQuery(concepts *[]models.Concept, codeSystemId int32, codeSystemVersionId int32, pageSize int, offset int, sortBy string, sortOrder string, meaning string, code string) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		query := tx.
			Preload("ValidFromVersion").
			Preload("ValidToVersion").
			Model(&models.Concept{}).
			Joins("JOIN code_system_versions AS valid_from_version ON valid_from_version.id = concepts.valid_from_version_id").
			Joins("JOIN code_system_versions AS valid_to_version ON valid_to_version.id = concepts.valid_to_version_id").
			Where("concepts.code_system_id = ?", codeSystemId).
			Where("valid_from_version.version_id <= ? AND valid_to_version.version_id >= ?", codeSystemVersionId, codeSystemVersionId)

		// Add code condition if code is not empty
		if code != "" {
			query = query.Where("LOWER(code) LIKE LOWER(?)", code+"%")
		}

		// Add meaning condition if meaning is not empty
		if meaning != "" {
			formattedMeaning := strings.Join(strings.Fields(meaning), ":* & ") + ":*" // Adjust for partial matches
			query = query.Where("display_search_vector @@ to_tsquery('english', ?)", formattedMeaning)
			// Tests for similarity searches. These were very slow and therefore not used in the final implementation. I a search should be implemented which is not part of a autocomplete but is ok to take e.g a second to complete, similarity searches with the pg_trgm extension could be used here.

			// query = query.Where("display_search_vector @@ to_tsquery(?) OR similarity(display, ?) > 0.3", formattedMeaning, meaning)
			// query = query.Where("similarity(display, ?) > 0.8", meaning)
			// query = query.Select("*, similarity(display, ?) > set_limit(0.99) AS s", meaning).Order("s DESC")
		}

		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize)

		if err := query.Find(&concepts).Error; err != nil {
			return err
		} else if len(*concepts) == 0 {
			var codesystem models.CodeSystem
			if err := tx.First(&codesystem, codeSystemId).Error; err != nil {
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
			} else if err := tx.First(&models.CodeSystemVersion{}, codeSystemVersionId).Error; err != nil {
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found.", codeSystemVersionId))
			} else {
				*concepts = []models.Concept{}
				return nil
			}
		}
		return nil
	})
	return err
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

func (gq GormQuery) GetNeighborConceptsQuery(code string, codeSystemId int32, codeSystemVersionId int32) (database.NeighborConcepts, error) {
	var versionId uint32
	var codeSystemVersion models.CodeSystemVersion
	if err := gq.GetCodeSystemVersionQuery(&codeSystemVersion, codeSystemId, codeSystemVersionId); err != nil {
		return database.NeighborConcepts{}, err
	} else {
		versionId = codeSystemVersion.VersionID
	}

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

		beforeVersionId, afterVersionId, err := gq.GetImportedNeighborVersionIdsQuery(codeSystemId, codeSystemVersionId)
		if err != nil {
			return database.NeighborConcepts{}, err
		}

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
