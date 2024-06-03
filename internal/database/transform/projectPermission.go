package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
	"miracummapper/internal/utilities"
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

func ApiProjectPermissionToGormProjectPermission(projectPermission api.ProjectPermission, projectId int32) (models.ProjectPermission, error) {
	var permission models.ProjectPermission

	userUuid, err := utilities.ParseUUID(projectPermission.UserId)
	if err != nil {
		return permission, err
	}
	permission = models.ProjectPermission{
		Role:      models.ProjectPermissionRole(projectPermission.Role),
		UserID:    userUuid,
		ProjectID: uint32(projectId),
	}

	return permission, nil
}
