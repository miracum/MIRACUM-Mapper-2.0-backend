package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
	"miracummapper/internal/utilities"
)

func GormProjectPermissionsToApiProjectPermissions(projectPermissions *[]models.ProjectPermission) *[]api.ProjectPermission {
	var permissions []api.ProjectPermission
	for _, permission := range *projectPermissions {
		permissions = append(permissions, *GormProjectPermissionToApiProjectPermission(&permission))
	}
	return &permissions
}

func GormProjectPermissionToApiProjectPermission(projectPermission *models.ProjectPermission) *api.ProjectPermission {
	return &api.ProjectPermission{
		Role:     api.ProjectPermissionRole(projectPermission.Role),
		UserId:   projectPermission.UserID.String(),
		UserName: projectPermission.User.UserName,
	}
}

func convertToGormProjectPermission(userId string, role models.ProjectPermissionRole, projectId int32) (*models.ProjectPermission, error) {
	var permission models.ProjectPermission

	userUuid, err := utilities.ParseUUID(userId)
	if err != nil {
		return &permission, err
	}
	permission = models.ProjectPermission{
		Role:      role,
		UserID:    userUuid,
		ProjectID: uint32(projectId),
	}

	return &permission, nil
}

func ApiProjectPermissionToGormProjectPermission(projectPermission *api.ProjectPermission, projectId int32) (*models.ProjectPermission, error) {
	return convertToGormProjectPermission(projectPermission.UserId, models.ProjectPermissionRole(projectPermission.Role), projectId)
}

func ApiSendProjectPermissionToGormProjectPermission(projectPermission *api.SendProjectPermission, projectId int32) (*models.ProjectPermission, error) {
	return convertToGormProjectPermission(projectPermission.UserId, models.ProjectPermissionRole(projectPermission.Role), projectId)
}
