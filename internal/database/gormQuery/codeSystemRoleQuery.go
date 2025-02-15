package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"

	"gorm.io/gorm"
)

func (gq *GormQuery) GetAllCodeSystemRolesQuery(codeSystemRoles *[]models.CodeSystemRole, projectId int32) error {
	db := gq.Database.Preload("CodeSystem").Preload("CodeSystemVersion").Preload("NextCodeSystemVersion").Where("project_id = ?", projectId).Find(&codeSystemRoles)
	if db.Error != nil {
		pgErr, ok := handlePgError(db.Error)
		if !ok {
			return db.Error
		}
		switch {
		case errors.Is(pgErr, gorm.ErrRecordNotFound):
			return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
		default:
			return pgErr
		}
	} else if db.RowsAffected == 0 {
		return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
	}
	if len(*codeSystemRoles) == 0 {
		return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d has no CodeSystemRoles.", projectId))
	}
	return nil
}

func (gq *GormQuery) GetCodeSystemRoleQuery(codeSystemRole *models.CodeSystemRole, projectId int32, codeSystemRoleId int32) error {
	db := gq.Database.Preload("CodeSystem").Preload("CodeSystemVersion").Preload("NextCodeSystemVersion").
		Where("project_id = ?", projectId).
		First(&codeSystemRole, codeSystemRoleId)
	if db.Error != nil {
		switch {
		case errors.Is(db.Error, gorm.ErrRecordNotFound):
			var project models.Project
			if err := gq.Database.First(&project, projectId).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
				}
				return err
			}
			return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemRole with ID %d couldn't be found in Project with ID %d.", codeSystemRoleId, projectId))
		default:
			return db.Error
		}
	}
	return nil
}

func (gq *GormQuery) UpdateCodeSystemRoleQuery(codeSystemRole *models.CodeSystemRole, projectId int32) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		oldCodeSystemRole := models.CodeSystemRole{}

		if err := tx.Preload("CodeSystem").Preload("CodeSystemVersion").Preload("NextCodeSystemVersion").Where("project_id = ?", projectId).First(&oldCodeSystemRole, codeSystemRole.ID).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				var project models.Project
				if err := tx.First(&project, projectId).Error; err != nil {
					switch {
					case errors.Is(err, gorm.ErrRecordNotFound):
						return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
					default:
						return err
					}
				} else {
					return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemRole with ID %d couldn't be found.", codeSystemRole.ID))
				}
			default:
				return err
			}
		}

		codeSystemRole.CodeSystemID = oldCodeSystemRole.CodeSystemID
		codeSystemRole.Position = oldCodeSystemRole.Position

		if err := tx.Save(codeSystemRole).Error; err != nil {
			return err
		}
		codeSystemRole.CodeSystem = oldCodeSystemRole.CodeSystem
		return nil
	})
	return err

}
