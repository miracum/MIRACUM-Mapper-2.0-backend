package gormQuery

import (
	"miracummapper/internal/database"

	"gorm.io/gorm"
)

type GormQuery struct {
	Database *gorm.DB
}

// Ensure Store implements the Datastore interface
var _ database.Datastore = &GormQuery{}
