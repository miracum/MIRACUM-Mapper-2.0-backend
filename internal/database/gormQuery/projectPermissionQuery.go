package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (gq *GormQuery) CreateProjectPermissionQuery(projectPermission *models.ProjectPermission) error {
	db := gq.Database.Create(projectPermission)

	if db.Error != nil {
		// cast error to postgres error
		pgErr, ok := handlePgError(db.Error)
		if !ok {
			return db.Error
		}
		switch pgErr.Code {
		case "23503":
			switch pgErr.ConstraintName {
			case "fk_users_project_permissions":
				userID, err := extractIDFromErrorDetail(pgErr.Detail, "user_id")
				if err != nil {
					return db.Error
				}
				return database.NewDBError(database.ClientError, fmt.Sprintf("User with id %s specified in permissions does not exist", userID))

			case "fk_code_systems_code_system_roles":
				codeSystemID, err := extractIDFromErrorDetail(pgErr.Detail, "code_system_id")
				if err != nil {
					return db.Error
				}
				return database.NewDBError(database.ClientError, fmt.Sprintf("Code System with id %s specified in code system roles does not exist", codeSystemID))
			default:
				return db.Error
			}
		default:
			return db.Error
		}
	} else {
		return nil
	}
}

func (gq *GormQuery) GetProjectPermissionsQuery(projectPermissions *[]models.ProjectPermission, projectId int32) error {
	db := gq.Database.Where("project_id = ?", projectId).Preload("User").Find(&projectPermissions)

	if db.Error != nil {
		// pgErr, ok := handlePgError(db.Error)
		// if !ok {
		// 	return db.Error
		// }
		switch {
		case errors.Is(db.Error, gorm.ErrRecordNotFound):
			return database.NewDBError(database.NotFound, fmt.Sprintf("ProjectPermission with Project ID %d couldn't be found.", projectId))
		default:
			return db.Error
		}
	}
	if len(*projectPermissions) == 0 {
		return database.NewDBError(database.NotFound, fmt.Sprintf("ProjectPermission with Project ID %d couldn't be found.", projectId))
	}
	return nil
}

func (gq *GormQuery) GetProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32) error {
	db := gq.Database.Where("project_id = ? AND user_id = ?", projectId, projectPermission.UserID).Preload("User").First(projectPermission)

	if db.Error != nil {
		// pgErr, ok := handlePgError(db.Error)
		// if !ok {
		// 	return db.Error
		// }
		switch {
		case errors.Is(db.Error, gorm.ErrRecordNotFound):
			// TODO This check to determine if the Project or the CodeSystemRole is not found is bad
			var project models.Project
			if err := gq.Database.First(&project, projectId).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
				}
				return err
			}
			return database.NewDBError(database.NotFound, fmt.Sprintf("ProjectPermission with User ID %d couldn't be found in Project with ID %d.", projectPermission.UserID, projectId))
		default:
			return db.Error
		}
	}
	return nil
}

func (gq *GormQuery) DeleteProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32, userId uuid.UUID) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ? AND user_id = ?", projectId, userId).First(projectPermission).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				// TODO This check to determine if the Project or the CodeSystemRole is not found is bad
				var project models.Project
				if err := tx.First(&project, projectId).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
					}
					return err
				}
				return database.NewDBError(database.NotFound, fmt.Sprintf("The user with id %s does not have a permission for the project with id %d", userId, projectId))
			default:
				return err
			}
		}

		db := tx.Delete(projectPermission)
		if db.Error != nil {
			return db.Error
		} else {
			if db.RowsAffected == 0 {
				return database.NewDBError(database.NotFound, fmt.Sprintf("The user with id %s does not have a permission for the project with id %d", projectPermission.UserID, projectId))
			}
			return nil
		}
	})
	return err
}

func (gq *GormQuery) UpdateProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		oldProjectPermission := models.ProjectPermission{}

		if err := tx.Where("project_id = ? AND user_id = ?", projectId, projectPermission.UserID).First(&oldProjectPermission).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				// TODO This check to determine if the Project or the CodeSystemRole is not found is bad
				var project models.Project
				if err := tx.First(&project, projectId).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
					}
					return err
				}
				return database.NewDBError(database.NotFound, fmt.Sprintf("The user with id %s does not have a permission for the project with id %d", projectPermission.UserID, projectId))
			default:
				return err
			}
		}

		if err := tx.Save(projectPermission).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
