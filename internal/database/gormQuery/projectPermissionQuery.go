package gormQuery

import (
	"miracummapper/internal/database/models"

	"github.com/google/uuid"
)

func (gq *GormQuery) GetProjectPermissionsQuery(projectPermissions *[]models.ProjectPermission, projectId int32) error {
	db := gq.Database.Where("project_id = ?", projectId).Preload("User").Find(&projectPermissions)

	if db.Error != nil {
		switch {
		default:
			return db.Error
		}
	} else {
		return nil
	}
}

func (gq *GormQuery) GetProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32, userId uuid.UUID) error {
	db := gq.Database.Where("project_id = ? AND user_id = ?", projectId, userId).Preload("User").First(projectPermission)

	if db.Error != nil {
		switch {
		default:
			return db.Error
		}
	} else {
		return nil
	}

}
