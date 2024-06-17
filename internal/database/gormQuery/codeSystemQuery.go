package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"

	"gorm.io/gorm"
)

func (gq *GormQuery) GetAllCodeSystemsQuery(codeSystems *[]models.CodeSystem) error {
	db := gq.Database.Find(&codeSystems)
	return db.Error
}

func (gq *GormQuery) CreateCodeSystemQuery(codeSystem *models.CodeSystem) error {
	return gq.Database.Create(&codeSystem).Error
}

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

func (gq *GormQuery) UpdateCodeSystemQuery(codeSystem *models.CodeSystem) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		// TODO check fields that are not allowed to change
		if err := tx.First(&models.CodeSystem{}, codeSystem.ID).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystem.ID))
			default:
				return err
			}
		}

		if err := tx.Save(&codeSystem).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (gq *GormQuery) DeleteCodeSystemQuery(codeSystem *models.CodeSystem, codeSystemId int32) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		// get codeSystem so it can be returned in the api and then delete it
		if err := tx.First(&codeSystem, codeSystemId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
			default:
				return err
			}
		}

		db := tx.Delete(&codeSystem, codeSystemId)
		if db.Error != nil {
			return db.Error
		} else {
			if db.RowsAffected == 0 {
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
			}
			return nil
		}
	})
	return err
}
