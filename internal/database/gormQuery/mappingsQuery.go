package gormQuery

import (
	"fmt"
	"miracummapper/internal/database/models"

	"gorm.io/gorm"
)

// CreateMappingQuery implements database.Datastore.
func (gq *GormQuery) CreateMappingQuery(mapping *models.Mapping, checkFunc func(mapping *models.Mapping, project *models.Project) error) error {
	// start transaction, get codesystemrole ids for project, call check function, if no error, create mapping

	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		project := models.Project{
			Model: models.Model{
				ID: mapping.ProjectID,
			},
		}
		if err := tx.Preload("CodeSystemRoles", func(db *gorm.DB) *gorm.DB {
			return db.Order("Position ASC")
		}).Find(&project).Error; err != nil {
			return err
		}

		// Get all ConceptIDs from the Elements in the Mapping
		conceptIDs := make([]uint64, len(mapping.Elements))
		for i, element := range mapping.Elements {
			conceptIDs[i] = *element.ConceptID
		}

		// Load all Concepts with the ConceptIDs in one batch
		concepts := make([]models.Concept, len(conceptIDs))
		if err := tx.Find(&concepts, conceptIDs).Error; err != nil {
			return err
		}

		// Assign each Concept to its corresponding Element
		for i, concept := range concepts {
			mapping.Elements[i].Concept = concept
		}

		if err := checkFunc(mapping, &project); err != nil {
			return err
		}

		db := tx.Create(mapping)

		return db.Error
	})
	return err
}

// DeleteMappingQuery implements database.Datastore.
func (gq *GormQuery) DeleteMappingQuery(mapping *models.Mapping, mappingId int32) error {
	panic("unimplemented")
}

// GetAllMappingsQuery implements database.Datastore.
func (gq *GormQuery) GetAllMappingsQuery(mappings *[]models.Mapping, projectId int, pageSize int, offset int, sortBy string, sortOrder string) error {
	db := gq.Database.Where("project_id = ?", projectId).Preload("Elements.Concept.CodeSystem").Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize).Find(&mappings)
	return db.Error
}

// GetMappingQuery implements database.Datastore.
func (gq *GormQuery) GetMappingQuery(mapping *models.Mapping, projectId int, mappingId int32) error {
	db := gq.Database.Where("project_id = ?", projectId).Preload("Elements.Concept.CodeSystem", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC")
	}).First(mapping, mappingId)
	return db.Error
}

// UpdateMappingQuery implements database.Datastore.
func (gq *GormQuery) UpdateMappingQuery(mapping *models.Mapping, checkFunc func(oldMapping *models.Mapping, newMapping *models.Mapping) error) error {
	panic("unimplemented")
}
