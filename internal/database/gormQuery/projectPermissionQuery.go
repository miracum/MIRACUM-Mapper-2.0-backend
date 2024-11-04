package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (gq *GormQuery) GetAllProjectPermissionsQuery(projectPermissions *[]models.ProjectPermission, projectId int32) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", projectId).Preload("User").Find(&projectPermissions).Error; err != nil {
			return err
		} else if len(*projectPermissions) == 0 {
			var project models.Project
			if err := tx.First(&project, projectId).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
				}
				return err
			} else {
				*projectPermissions = []models.ProjectPermission{}
				return nil
			}
		}
		return nil
	})
	return err
}

func (gq *GormQuery) CreateProjectPermissionQuery(projectPermission *models.ProjectPermission) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(projectPermission).Error; err != nil {
			pgErr, ok := handlePgError(err)
			if !ok {
				return err
			}
			switch pgErr.Code {
			case "23503":
				switch pgErr.ConstraintName {
				case "fk_users_project_permissions":
					userID, err := extractIDFromErrorDetail(pgErr.Detail, "user_id")
					if err != nil {
						return err
					}
					return database.NewDBError(database.ClientError, fmt.Sprintf("User with id %s specified in permissions does not exist", userID))
				// TODO add error code (testen)
				case "fk_project_project_permissions":
					projectID, err := extractIDFromErrorDetail(pgErr.Detail, "project_id")
					if err != nil {
						return err
					}
					return database.NewDBError(database.ClientError, fmt.Sprintf("Project with id %s specified in permissions does not exist", projectID))
				default:
					return err
				}
			case "23505":
				if pgErr.ConstraintName == "project_permissions_pkey" {
					return database.NewDBError(database.ClientError, "The user already has a permission for this project. If you want to change it, please update the existing permission.")
				} else {
					return err
				}
			default:
				return err
			}
		} else {
			if err = tx.First(&projectPermission.User, projectPermission.UserID).Error; err != nil {
				return err
			}
			return nil
		}
	})
	return err
}

func (gq *GormQuery) GetProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32, userId uuid.UUID) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ? AND user_id = ?", projectId, userId).Preload("User").First(projectPermission).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				// TODO This check to determine if the Project or the CodeSystemRole is not found is bad
				var project models.Project
				if err := tx.First(&project, projectId).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return database.NewDBErrorWithTable(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId), database.ProjectTable)
					}
					return err
				} else {
					return database.NewDBError(database.NotFound, fmt.Sprintf("The user with id %s does not have a permission for the project with id %d or doesn't exist.", userId, projectId))
				}
			default:
				return err
			}
		}
		return nil
	})
	return err
}

func (gq *GormQuery) UpdateProjectPermissionQuery(projectPermission *models.ProjectPermission) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		oldProjectPermission := models.ProjectPermission{}
		projectId := projectPermission.ProjectID
		userId := projectPermission.UserID

		if err := tx.Where("project_id = ? AND user_id = ?", projectId, userId).Preload("User").First(&oldProjectPermission).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				// TODO This check to determine if the Project or the User is not found is bad
				var project models.Project
				if err := tx.First(&project, projectId).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
					}
					return err
				}
				return database.NewDBError(database.NotFound, fmt.Sprintf("The user with id %s does not have a permission for the project with id %d or doesn't exist", userId, projectId))
			default:
				return err
			}
		}

		if err := tx.Save(projectPermission).Error; err != nil {
			return err
		}
		projectPermission.User = oldProjectPermission.User
		return nil
	})
	return err
}

func (gq *GormQuery) DeleteProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32, userId uuid.UUID) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ? AND user_id = ?", projectId, userId).Preload("User").First(projectPermission).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				// TODO This check to determine if the Project or the CodeSystemRole is not found is bad
				var project models.Project
				if err := tx.First(&project, projectId).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
					}
					return err
				} else {
					return database.NewDBError(database.NotFound, fmt.Sprintf("The user with id %s does not have a permission for the project with id %d or doesn't exist.", userId, projectId))
				}
			default:
				return err
			}
		}

		if err := tx.Delete(projectPermission).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
