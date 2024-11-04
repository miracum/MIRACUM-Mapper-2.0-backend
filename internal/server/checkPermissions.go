package server

import (
	"context"
	"errors"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"slices"

	"github.com/google/uuid"
)

type ErrorType int

// this defined which roles are allowed to do which actions on a project.
var (
	ProjectViewPermission          *[]models.ProjectPermissionRole = &[]models.ProjectPermissionRole{models.ReviewerRole, models.ProjectOwnerRole, models.EditorRole}
	ProjectUpdatePermission        *[]models.ProjectPermissionRole = &[]models.ProjectPermissionRole{models.ProjectOwnerRole}
	MappingViewPermission          *[]models.ProjectPermissionRole = &[]models.ProjectPermissionRole{models.ReviewerRole, models.ProjectOwnerRole, models.EditorRole}
	MappingUpdatePermission        *[]models.ProjectPermissionRole = &[]models.ProjectPermissionRole{models.ProjectOwnerRole, models.EditorRole}
	MappingDeletePermission        *[]models.ProjectPermissionRole = &[]models.ProjectPermissionRole{models.ProjectOwnerRole, models.EditorRole}
	MappingCreatePermission        *[]models.ProjectPermissionRole = &[]models.ProjectPermissionRole{models.ProjectOwnerRole, models.EditorRole}
	MappingUpdateCommentPermission *[]models.ProjectPermissionRole = &[]models.ProjectPermissionRole{models.ProjectOwnerRole, models.EditorRole}
)

// an admin has all permissions
var adminRoles = &[]models.ProjectPermissionRole{models.EditorRole, models.ProjectOwnerRole, models.ReviewerRole}

type PermissionType int

const (
	ProjectPermission PermissionType = iota
	MappingPermission PermissionType = iota
)

// return the user id or nil if the user has admin rights
func getUserToCheckPermission(ctx context.Context) (*uuid.UUID, error) {
	if !IsAdminFromContext(ctx) {
		userId, err := GetUserIdFromContext(ctx)
		return &userId, err
	}
	return nil, nil
}

func getUserPermissions(ctx context.Context, server *Server, projectId int32) (*[]models.ProjectPermissionRole, error) {
	userId, err := getUserToCheckPermission(ctx)
	if err != nil {
		return nil, err
	}
	if userId != nil {
		var permission models.ProjectPermission
		if err := server.Database.GetProjectPermissionQuery(&permission, projectId, *userId); err != nil {
			switch {
			case errors.Is(err, database.ErrProjectNotFound):
				return nil, err
			case errors.Is(err, database.ErrNotFound):
				return &[]models.ProjectPermissionRole{}, nil
			default:
				return nil, err
			}
		}
		return &[]models.ProjectPermissionRole{permission.Role}, nil // currently a user can only have one role. Can be extended in the future
	}
	return adminRoles, nil
}

// this function takes a list of needed permissions and check if the actual permissions contain at least one of the needed permissions
func checkUserHasPermissions(neededPermissionRoles *[]models.ProjectPermissionRole, actualPermissions *[]models.ProjectPermissionRole) bool {
	for _, actualPermission := range *actualPermissions {
		if slices.Contains(*neededPermissionRoles, actualPermission) {
			return true
		}
	}
	return false
}
