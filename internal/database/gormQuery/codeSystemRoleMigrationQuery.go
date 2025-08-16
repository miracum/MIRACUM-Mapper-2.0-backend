package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"time"

	"gorm.io/gorm"
)

func (gq *GormQuery) StartMigrationQuery(projectId int32, codeSystemRoleId int32, versionId int32) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		var project models.Project
		if err := tx.Preload("CodeSystemRoles").Preload("CodeSystemRoles.CodeSystemVersion").Preload("CodeSystemRoles.NextCodeSystemVersion").First(&project, projectId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
			default:
				return err
			}
		}

		var codeSystemRole *models.CodeSystemRole
		for _, role := range project.CodeSystemRoles {
			if role.ID == codeSystemRoleId {
				codeSystemRole = &role
			}
			if role.NextCodeSystemVersionID != nil {
				return database.NewDBError(database.ClientError, "There is already a migration in progress for this project.")
			}
		}
		if codeSystemRole == nil {
			return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemRole with ID %d couldn't be found in Project with ID %d.", codeSystemRoleId, projectId))
		}

		var codeSystemVersion *models.CodeSystemVersion
		if err := tx.Where("code_system_id = ?", codeSystemRole.CodeSystemID).First(&codeSystemVersion, versionId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found.", versionId))
			default:
				return err
			}
		}

		if codeSystemVersion.VersionID <= codeSystemRole.CodeSystemVersion.VersionID {
			return database.NewDBError(database.ClientError, fmt.Sprintf("CodeSystemVersion with ID %d is not a newer version than CodeSystemVersion with ID %d.", versionId, codeSystemRole.CodeSystemVersion.ID))
		}

		codeSystemRole.NextCodeSystemVersionID = &codeSystemVersion.ID

		if err := tx.Save(&codeSystemRole).Error; err != nil {
			return err
		}

		return nil
	})
	return err
}

func (gq *GormQuery) CancelMigrationQuery(projectId int32) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		var project models.Project
		if err := tx.Preload("CodeSystemRoles").Preload("CodeSystemRoles.CodeSystemVersion").Preload("CodeSystemRoles.NextCodeSystemVersion").First(&project, projectId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
			default:
				return err
			}
		}

		var codeSystemRole *models.CodeSystemRole
		for _, role := range project.CodeSystemRoles {
			if role.NextCodeSystemVersionID != nil {
				codeSystemRole = &role
			}
		}
		if codeSystemRole == nil {
			return database.NewDBError(database.ClientError, "There is no migration in progress for this project.")
		}

		codeSystemRole.NextCodeSystemVersionID = nil
		codeSystemRole.NextCodeSystemVersion = models.CodeSystemVersion{}

		if err := tx.Save(&codeSystemRole).Error; err != nil {
			return err
		}

		var mappings []models.Mapping
		if err := tx.Where("project_id = ?", projectId).Preload("Elements.NextConcept.CodeSystem").Find(&mappings).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("Mappings for Project with ID %d couldn't be found.", projectId))
			default:
				return err
			}
		}

		for _, mapping := range mappings {
			for _, element := range mapping.Elements {
				if element.CodeSystemRoleID == codeSystemRole.ID {
					if element.NextConceptID != nil {
						if err := tx.Model(models.Element{}).Where("mapping_id = ? AND code_system_role_id = ?", mapping.ID, codeSystemRole.ID).Updates(map[string]any{
							"next_concept_id": nil,
						}).Error; err != nil {
							return err
						}
					}
					break
				}
			}
		}

		return nil
	})
	return err
}

func (gq *GormQuery) GetMigrationCodeSystemRoleQuery(projectId int32) (*models.CodeSystemRole, error) {
	var codeSystemRole *models.CodeSystemRole

	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		var project models.Project
		if err := tx.Preload("CodeSystemRoles").Preload("CodeSystemRoles.CodeSystem").Preload("CodeSystemRoles.CodeSystemVersion").Preload("CodeSystemRoles.NextCodeSystemVersion").First(&project, projectId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
			default:
				return err
			}
		}

		for _, role := range project.CodeSystemRoles {
			if role.NextCodeSystemVersionID != nil {
				codeSystemRole = &role
			}
		}
		if codeSystemRole == nil {
			return database.NewDBError(database.ClientError, "There is no migration in progress for this project.")
		}
		return nil
	})
	return codeSystemRole, err
}

func (gq *GormQuery) GetMigrationValidToVersionIdsQuery(codeSystemId int32, nextCodeSystemVersionId int32) ([]int32, error) {
	var validToVersionIds []int32

	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		var codeSystem models.CodeSystem
		if err := gq.Database.Preload("CodeSystemVersions").First(&codeSystem, codeSystemId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
			default:
				return err
			}
		}
		var nextCodeSystemVersionVersionId int32 = 0
		for _, version := range codeSystem.CodeSystemVersions {
			if version.ID == nextCodeSystemVersionId {
				nextCodeSystemVersionVersionId = version.VersionID
			}
		}
		if nextCodeSystemVersionVersionId == 0 {
			return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found for CodeSystem with ID %d.", nextCodeSystemVersionId, codeSystemId))
		}
		for _, version := range codeSystem.CodeSystemVersions {
			if version.VersionID >= nextCodeSystemVersionVersionId {
				validToVersionIds = append(validToVersionIds, version.ID)
			}
		}
		return nil
	})
	return validToVersionIds, err
}

func (gq *GormQuery) GetMigrationValidFromVersionIdsQuery(codeSystemId int32, codeSystemVersionId int32) ([]int32, error) {
	var validFromVersionIds []int32

	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		var codeSystem models.CodeSystem
		if err := gq.Database.Preload("CodeSystemVersions").First(&codeSystem, codeSystemId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
			default:
				return err
			}
		}
		var codeSystemVersionVersionId int32 = 0
		for _, version := range codeSystem.CodeSystemVersions {
			if version.ID == codeSystemVersionId {
				codeSystemVersionVersionId = version.VersionID
			}
		}
		if codeSystemVersionVersionId == 0 {
			return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found for CodeSystem with ID %d.", codeSystemVersionId, codeSystemId))
		}
		for _, version := range codeSystem.CodeSystemVersions {
			if version.VersionID <= codeSystemVersionVersionId {
				validFromVersionIds = append(validFromVersionIds, version.ID)
			}
		}
		return nil
	})
	return validFromVersionIds, err
}

func (gq *GormQuery) FinishMigrationQuery(codeSystemRole *models.CodeSystemRole, mappings *[]models.Mapping, statusRequired bool) error {
	codeSystemRoleId := codeSystemRole.ID

	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		codeSystemRole.CodeSystemVersionID = *codeSystemRole.NextCodeSystemVersionID
		codeSystemRole.CodeSystemVersion = codeSystemRole.NextCodeSystemVersion
		codeSystemRole.NextCodeSystemVersionID = nil
		codeSystemRole.NextCodeSystemVersion = models.CodeSystemVersion{}
		if err := tx.Save(&codeSystemRole).Error; err != nil {
			return err
		}
		for _, mapping := range *mappings {
			commentString := ""
			for _, element := range mapping.Elements {
				if element.CodeSystemRoleID == codeSystemRoleId {
					if element.NextConceptID != nil {
						commentString = getMigrationCommentString(codeSystemRole, &element.Concept, &element.NextConcept, element.Concept.Code != element.NextConcept.Code)

						element.ConceptID = element.NextConceptID
						element.Concept = element.NextConcept
						element.NextConceptID = nil
						element.NextConcept = models.Concept{}
						if err := tx.Save(&element).Error; err != nil {
							return err
						}
					}
					break
				}
			}
			if commentString != "" {
				if mapping.Comment == nil || *mapping.Comment == "" {
					mapping.Comment = &commentString
				} else {
					newComment := commentString + " " + *mapping.Comment
					mapping.Comment = &newComment
				}
				if statusRequired {
					newStatus := models.Migrated
					mapping.Status = &newStatus
				}
				if err := tx.Save(&mapping).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func getMigrationCommentString(codeSystemRole *models.CodeSystemRole, oldConcept *models.Concept, newConcept *models.Concept, conceptChanged bool) string {
	year, month, day := time.Now().Date()
	comment := fmt.Sprintf("Migrated %s %s ", codeSystemRole.Type, codeSystemRole.CodeSystem.Name)
	if codeSystemRole.Name != "" {
		comment += fmt.Sprintf("(%s) ", codeSystemRole.Name)
	}
	comment += fmt.Sprintf("on %d-%02d-%02d: ", year, month, day)

	if conceptChanged {
		comment += fmt.Sprintf("Concept [%s] (%s) was replaced by [%s] (%s). ", oldConcept.Code, oldConcept.Display, newConcept.Code, newConcept.Display)
	} else {
		if oldConcept.Status != newConcept.Status {
			comment += fmt.Sprintf("Status has changed from [%s] to [%s]. ", oldConcept.Status, newConcept.Status)
		}
		if oldConcept.Display != newConcept.Display {
			comment += fmt.Sprintf("Display has changed from [%s] to [%s]. ", oldConcept.Display, newConcept.Display)
		}
		if oldConcept.Description != nil && newConcept.Description != nil {
			if *oldConcept.Description != *newConcept.Description {
				comment += fmt.Sprintf("Description has changed from '[%s]' to '[%s]'. ", *oldConcept.Description, *newConcept.Description)
			}
		} else if oldConcept.Description == nil && newConcept.Description != nil {
			comment += fmt.Sprintf("Description was added: [%s]. ", *newConcept.Description)
		} else if oldConcept.Description != nil && newConcept.Description == nil {
			comment += fmt.Sprintf("Description was removed: [%s]. ", *oldConcept.Description)
		}
	}
	return fmt.Sprintf("--- %s ---", comment)
}

func (gq *GormQuery) MigrationElementSetNextConceptQuery(mappingId int32, codeSystemRoleId int32, nextConceptId *int32) error {
	if err := gq.Database.Model(&models.Element{}).Where("mapping_id = ? AND code_system_role_id = ?", mappingId, codeSystemRoleId).Updates(map[string]any{
		"next_concept_id": nextConceptId,
	}).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return database.NewDBError(database.NotFound, fmt.Sprintf("Element with mapping ID %d and code system role ID %d couldn't be found.", mappingId, codeSystemRoleId))
		default:
			return err
		}
	}
	return nil
}
