package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"

	"gorm.io/gorm"
)

func CreateOrUpdateMapping(gq *GormQuery, mapping *models.Mapping, checkFunc func(mapping *models.Mapping, project *models.Project) ([]uint32, error), deleteMissingElements bool) error {
	// start transaction, get CodeSystemRole ids for project, call check function, if no error, create mapping

	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		project := models.Project{
			Model: models.Model{
				ID: mapping.ProjectID,
			},
		}
		if err := tx.Preload("CodeSystemRoles", func(db *gorm.DB) *gorm.DB {
			return db.Order("ID ASC")
		}).First(&project).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", project.ID))
			default:
				return err
			}
		}

		// Get all ConceptIDs from the Elements in the Mapping and find them in the database
		//  concepts := make([]models.Concept, len(mapping.Elements))
		for i, element := range mapping.Elements {
			var concept models.Concept
			if err := tx.First(&concept, element.ConceptID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return database.NewDBError(database.NotFound, fmt.Sprintf("Concept with ID %d couldn't be found.", *element.ConceptID))
				}
				return err
			}
			mapping.Elements[i].Concept = concept
		}

		unusedCodeSystemRoleIds, err := checkFunc(mapping, &project)
		if err != nil {
			return err
		}

		if len(unusedCodeSystemRoleIds) > 0 {
			if deleteMissingElements {
				// Delete all elements that are not in the provided list
				if err := tx.Where("mapping_id = ? AND code_system_role_id IN (?)", mapping.ID, unusedCodeSystemRoleIds).Delete(&models.Element{}).Error; err != nil {
					return err
				}
			} else {
				elements := []models.Element{}
				if err := tx.Where("mapping_id = ? AND code_system_role_id IN (?)", mapping.ID, unusedCodeSystemRoleIds).Preload("Concept").Find(&elements).Error; err != nil {
					return err
				}
				mapping.Elements = append(mapping.Elements, elements...)
			}
		}

		db := tx.Save(mapping)

		return db.Error
	})
	return err
}

// GetAllMappingsQuery implements database.Datastore.
func (gq *GormQuery) GetAllMappingsQuery(mappings *[]models.Mapping, projectId int, pageSize int, offset int, sortBy string, sortOrder string) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", projectId).Preload("Elements.Concept.CodeSystem").Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize).Find(&mappings).Error; err != nil {
			return err
		} else if len(*mappings) == 0 {
			if err := tx.First(&models.Project{}, projectId).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
				} else {
					*mappings = []models.Mapping{}
					return nil
				}
			}
		}
		return nil
	})
	return err

}

// CreateMappingQuery implements database.Datastore.
func (gq *GormQuery) CreateMappingQuery(mapping *models.Mapping, checkFunc func(mapping *models.Mapping, project *models.Project) ([]uint32, error)) error {
	return CreateOrUpdateMapping(gq, mapping, checkFunc, false)
}

// GetMappingQuery implements database.Datastore.
func (gq *GormQuery) GetMappingQuery(mapping *models.Mapping, projectId int, mappingId int64) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", projectId).Preload("Elements.Concept.CodeSystem").First(mapping, mappingId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				var project models.Project
				if err := tx.First(&project, projectId).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
					}
					return err
				} else {
					return database.NewDBError(database.NotFound, fmt.Sprintf("The mapping with id %d does not have a permission for the project with id %d.", mappingId, projectId))
				}
			default:
				return err
			}
		}
		return nil
	})
	return err
}

// UpdateMappingQuery implements database.Datastore.
func (gq *GormQuery) UpdateMappingQuery(mapping *models.Mapping, checkFunc func(mapping *models.Mapping, project *models.Project) ([]uint32, error), deleteMissingElements bool) error {
	// TODO it has to be checked if the mapping exists in the project. If there is another project which has the same CodeSystem Roles and the projectId of the other project is set in the update mapping url, the mapping would get moved to the other project
	return CreateOrUpdateMapping(gq, mapping, checkFunc, deleteMissingElements)
}

// DeleteMappingQuery implements database.Datastore.
func (gq *GormQuery) DeleteMappingQuery(mapping *models.Mapping) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", mapping.ProjectID).First(mapping, mapping.ID).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				// TODO This check to determine if the Project or the CodeSystemRole is not found is bad
				var project models.Project
				if err := tx.First(&project, mapping.ProjectID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", mapping.ProjectID))
					}
					return err
				} else {
					return database.NewDBError(database.NotFound, fmt.Sprintf("The Mapping with id %d does not exist in the project with id %d.", mapping.ID, mapping.ProjectID))
				}
			default:
				return err
			}
		}
		if err := gq.Database.Delete(mapping).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
