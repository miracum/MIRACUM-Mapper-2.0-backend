package gormQuery

import (
	"errors"
	"fmt"
	"log"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/utilities"

	"gorm.io/gorm"
)

/*
	IMPORTANT: Most of the functions in this file do not check if the CodeSystem, CodeSystemVersion or Concept exists in the database.
	This is done for better performance (to avoid unnecessary queries).
	Therefore, the calling Function should check if the CodeSystem, CodeSystemVersion or Concept exists in the database.
*/

type NeighborConceptsType int

const (
	NeighborConceptsTypeBefore NeighborConceptsType = iota
	NeighborConceptsTypeAfter
	NeighborConceptsTypeBeforeAndAfter
	NeighborConceptsTypeSurrounding
	NeighborConceptsTypeNone
)

type NeighborConcepts struct {
	BeforeConcept      *models.Concept
	AfterConcept       *models.Concept
	SurroundingConcept *models.Concept
	NeighborType       NeighborConceptsType
}

func (gq GormQuery) ImportConcepts(codeSystemId int32, codeSystemVersionId int32, concepts *[]database.ConceptImport, replaceBies *[]models.ConceptReplaceBy) {
	defer utilities.DoneImporting()

	numConcepts := len(*concepts)
	numReplaceBies := 0
	if replaceBies != nil {
		numReplaceBies = len(*replaceBies)
	}
	totalNum := numConcepts + numReplaceBies

	utilities.SetImportBegin()

	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&models.CodeSystem{}, codeSystemId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return fmt.Errorf("CodeSystem with ID %d couldn't be found", codeSystemVersionId)
			default:
				return err
			}
		}
		if err := tx.Where("code_system_id = ?", codeSystemId).First(&models.CodeSystemVersion{}, codeSystemVersionId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return fmt.Errorf("CodeSystemVersion with ID %d couldn't be found for CodeSystem with ID %d", codeSystemVersionId, codeSystemId)
			default:
				return err
			}
		}

		versionId, beforeVersionId, afterVersionId, err := getImportedNeighborVersionIds(tx, codeSystemId, codeSystemVersionId)
		if err != nil {
			return fmt.Errorf("error getting neighbor version IDs: %v", err)
		}

		for i, concept := range *concepts {
			utilities.SetImportProgress(i * 100 / totalNum)

			neighborConcepts, err := getNeighborConcepts(tx, concept.Code, codeSystemId, versionId, beforeVersionId, afterVersionId)
			if err != nil {
				return fmt.Errorf("error getting neighbor concepts: %v", err)
			}
			switch neighborConcepts.NeighborType {
			case NeighborConceptsTypeNone:
				if err := createNewConcept(tx, codeSystemId, codeSystemVersionId, &concept); err != nil {
					return err
				}
			case NeighborConceptsTypeBefore:
				beforeConcept := neighborConcepts.BeforeConcept
				if !conceptsAreEqual(&concept, beforeConcept) {
					if err := createNewConcept(tx, codeSystemId, codeSystemVersionId, &concept); err != nil {
						return err
					}
				} else {
					beforeConcept.ValidToVersionID = codeSystemVersionId
					if err := UpdateConceptQuery(tx, beforeConcept); err != nil {
						return err
					}
				}
			case NeighborConceptsTypeAfter:
				afterConcept := neighborConcepts.AfterConcept
				if !conceptsAreEqual(&concept, afterConcept) {
					if err := createNewConcept(tx, codeSystemId, codeSystemVersionId, &concept); err != nil {
						return err
					}
				} else {
					afterConcept.ValidFromVersionID = codeSystemVersionId
					if err := UpdateConceptQuery(tx, afterConcept); err != nil {
						return err
					}
				}
			case NeighborConceptsTypeBeforeAndAfter:
				beforeConcept := neighborConcepts.BeforeConcept
				afterConcept := neighborConcepts.AfterConcept
				if !conceptsAreEqual(&concept, beforeConcept) && !conceptsAreEqual(&concept, afterConcept) {
					if err := createNewConcept(tx, codeSystemId, codeSystemVersionId, &concept); err != nil {
						return err
					}
				} else if conceptsAreEqual(&concept, beforeConcept) && !conceptsAreEqual(&concept, afterConcept) {
					beforeConcept.ValidToVersionID = codeSystemVersionId
					if err := UpdateConceptQuery(tx, beforeConcept); err != nil {
						return err
					}
				} else if !conceptsAreEqual(&concept, beforeConcept) && conceptsAreEqual(&concept, afterConcept) {
					afterConcept.ValidFromVersionID = codeSystemVersionId
					if err := UpdateConceptQuery(tx, afterConcept); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("error: Concept is before and after and both are equal (indicates invalid data)")
				}
			case NeighborConceptsTypeSurrounding:
				surroundingConcept := neighborConcepts.SurroundingConcept
				if !conceptsAreEqual(&concept, surroundingConcept) {
					return fmt.Errorf("error: Surrounding concept is not equal to new concept (indicates invalid data)")
				} else {
					// Do nothing
				}
			}
		}

		if replaceBies != nil {
			for i, replaceBy := range *replaceBies {
				utilities.SetImportProgress((i + numConcepts) * 100 / totalNum)

				var existingReplaceBy models.ConceptReplaceBy
				if err := tx.Where("code_system_id = ? AND code = ? AND map_to = ?", replaceBy.CodeSystemID, replaceBy.Code, replaceBy.MapTo).First(&existingReplaceBy).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						// If it doesn't exist, create it
						if err := tx.Create(&replaceBy).Error; err != nil {
							return fmt.Errorf("error creating ConceptReplaceBy: %v", err)
						}
					} else {
						return fmt.Errorf("error checking existing ConceptReplaceBy: %v", err)
					}
				} else {
					// If it exists, update it
					existingReplaceBy.Equivalence = replaceBy.Equivalence
					existingReplaceBy.Comment = replaceBy.Comment
					if err := tx.Save(&existingReplaceBy).Error; err != nil {
						return fmt.Errorf("error updating ConceptReplaceBy: %v", err)
					}
				}
			}
		}

		if err := setCodeSystemVersionImported(tx, codeSystemVersionId, true); err != nil {
			return fmt.Errorf("error setting CodeSystemVersion as imported: %v", err)
		}

		return nil
	})

	if err != nil {
		utilities.SetImportDone(err)
		log.Printf("Error importing concepts: %v", err)
	} else {
		utilities.SetImportDone(nil)
		log.Print("Successfully imported concepts")
	}
}

func getImportedNeighborVersionIds(db *gorm.DB, codeSystemId int32, codeSystemVersionId int32) (int32, *int32, *int32, error) {
	var codeSystemVersion models.CodeSystemVersion
	if err := db.Where("code_system_id = ?", codeSystemId).First(&codeSystemVersion, codeSystemVersionId).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return 0, nil, nil, fmt.Errorf("CodeSystemVersion with ID %d couldn't be found for CodeSystem with ID %d", codeSystemVersionId, codeSystemId)
		default:
			return 0, nil, nil, err
		}
	}

	versionId := codeSystemVersion.VersionID

	var beforeVersionId *int32
	beforeVersionId = nil
	var afterVersionId *int32
	afterVersionId = nil

	var beforeCodeSystemVersions []models.CodeSystemVersion
	if err := db.Where("code_system_id = ? AND version_id < ?", codeSystemId, versionId).Order("version_id DESC").Find(&beforeCodeSystemVersions).Error; err != nil {
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
	if err := db.Where("code_system_id = ? AND version_id > ?", codeSystemId, versionId).Order("version_id ASC").Find(&afterCodeSystemVersions).Error; err != nil {
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

func getNeighborConcepts(db *gorm.DB, code string, codeSystemId int32, versionId int32, beforeVersionId *int32, afterVersionId *int32) (NeighborConcepts, error) {
	var concepts []models.Concept
	if err := db.Preload("ValidFromVersion").Preload("ValidToVersion").Model(&models.Concept{}).Where("code_system_id = ?", codeSystemId).Where("code = ?", code).Find(&concepts).Error; err != nil {
		return NeighborConcepts{}, err
	}
	if len(concepts) == 0 {
		return NeighborConcepts{NeighborType: NeighborConceptsTypeNone}, nil
	}

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
		return NeighborConcepts{}, fmt.Errorf("invalid stored concepts found while getting Neighbors for Concept with code %s in CodeSystem with ID %d", code, codeSystemId)
	} else if len(surroundingConcepts) == 1 {
		return NeighborConcepts{SurroundingConcept: &surroundingConcepts[0], NeighborType: NeighborConceptsTypeSurrounding}, nil
	} else if len(beforeConcepts) == 1 && len(afterConcepts) == 1 {
		return NeighborConcepts{BeforeConcept: &beforeConcepts[0], AfterConcept: &afterConcepts[0], NeighborType: NeighborConceptsTypeBeforeAndAfter}, nil
	} else if len(beforeConcepts) == 1 {
		return NeighborConcepts{BeforeConcept: &beforeConcepts[0], NeighborType: NeighborConceptsTypeBefore}, nil
	} else if len(afterConcepts) == 1 {
		return NeighborConcepts{AfterConcept: &afterConcepts[0], NeighborType: NeighborConceptsTypeAfter}, nil
	} else {
		return NeighborConcepts{NeighborType: NeighborConceptsTypeNone}, nil
	}

}

func createNewConcept(db *gorm.DB, codeSystemId int32, codeSystemVersionId int32, concept *database.ConceptImport) error {
	newConcept := models.Concept{
		Code:               concept.Code,
		Display:            concept.Display,
		Description:        concept.Description,
		Status:             concept.Status,
		CodeSystemID:       codeSystemId,
		ValidFromVersionID: codeSystemVersionId,
		ValidToVersionID:   codeSystemVersionId,
	}
	if err := db.Create(&newConcept).Error; err != nil {
		return fmt.Errorf("error creating concept: %v", err)
	}
	return nil
}

func conceptsAreEqual(conceptImport *database.ConceptImport, conceptDB *models.Concept) bool {
	var descriptionsAreEqual bool
	if conceptImport.Description == nil || conceptDB.Description == nil {
		descriptionsAreEqual = conceptImport.Description == conceptDB.Description
	} else {
		descriptionsAreEqual = *conceptImport.Description == *conceptDB.Description
	}
	return conceptImport.Display == conceptDB.Display && descriptionsAreEqual && conceptImport.Status == conceptDB.Status
}

func setCodeSystemVersionImported(db *gorm.DB, codeSystemVersionId int32, imported bool) error {
	var codeSystemVersion models.CodeSystemVersion
	if err := db.First(&codeSystemVersion, codeSystemVersionId).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return fmt.Errorf("CodeSystemVersion with ID %d couldn't be found", codeSystemVersionId)
		default:
			return err
		}
	}

	codeSystemVersion.Imported = imported
	if err := db.Save(&codeSystemVersion).Error; err != nil {
		return err
	}
	return nil
}
