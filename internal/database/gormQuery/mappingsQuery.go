package gormQuery

import (
	"fmt"
	"miracummapper/internal/database/models"

	"gorm.io/gorm"
)

// CreateMappingQuery implements database.Datastore.
func (gq *GormQuery) CreateMappingQuery(mapping *models.Mapping) error {
	db := gq.Database.Create(mapping)
	if db.Error != nil {
		return db.Error
	}
	panic("unimplemented")
}

// DeleteMappingQuery implements database.Datastore.
func (gq *GormQuery) DeleteMappingQuery(mapping *models.Mapping, mappingId int32) error {
	panic("unimplemented")
}

// GetAllMappingsQuery implements database.Datastore.
func (gq *GormQuery) GetAllMappingsQuery(mappings *[]models.Mapping, pageSize int, offset int, sortBy string, sortOrder string) error {
	db := gq.Database.Preload("Elements.Concepts.CodeSystem", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC")
	}).Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize).Find(&mappings)
	return db.Error
}

// GetMappingQuery implements database.Datastore.
func (gq *GormQuery) GetMappingQuery(mapping *models.Mapping, mappingId int32) error {
	panic("unimplemented")
}

// UpdateMappingQuery implements database.Datastore.
func (gq *GormQuery) UpdateMappingQuery(mapping *models.Mapping, checkFunc func(oldMapping *models.Mapping, newMapping *models.Mapping) error) error {
	panic("unimplemented")
}
