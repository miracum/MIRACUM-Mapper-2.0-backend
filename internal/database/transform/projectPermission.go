package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
)

func GormProjectPermissionToApiProjectPermission(projectPermission models.ProjectPermission) api.ProjectPermission {
	var permission api.ProjectPermission
	permission = api.ProjectPermission{
		Role:     api.ProjectPermissionRole(projectPermission.Role),
		UserId:   projectPermission.UserID.String(),
		UserName: &projectPermission.User.UserName, // possibly nil
	}

	return permission
}
