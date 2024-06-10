package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"

	"gorm.io/gorm"
)

// CreateCodeSystemQuery implements database.Datastore.
func (gq *GormQuery) GetAllCodeSystemsQuery(codeSystems *[]models.CodeSystem) error {
	db := gq.Database.Find(&codeSystems)
	return db.Error
}

// GetCodeSystemQuery implements database.Datastore.
func (gq *GormQuery) GetCodeSystemQuery(codeSystem *models.CodeSystem, codeSystemId int32) error {
	db := gq.Database.First(&codeSystem, codeSystemId)
	if db.Error != nil {
		switch {
		case errors.Is(db.Error, gorm.ErrRecordNotFound):
			return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
		default:
			return db.Error
		}
	} else {
		return nil
	}
}
