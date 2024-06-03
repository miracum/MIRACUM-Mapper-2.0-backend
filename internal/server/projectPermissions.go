package server

import (
	"context"
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
	"miracummapper/internal/utilities"
)

// AddPermission implements api.StrictServerInterface.
func (s *Server) AddPermission(ctx context.Context, request api.AddPermissionRequestObject) (api.AddPermissionResponseObject, error) {
	panic("unimplemented")
}

// DeletePermission implements api.StrictServerInterface.
func (s *Server) DeletePermission(ctx context.Context, request api.DeletePermissionRequestObject) (api.DeletePermissionResponseObject, error) {
	panic("unimplemented")
}

// GetAllPermissions implements api.StrictServerInterface.
func (s *Server) GetAllPermissions(ctx context.Context, request api.GetAllPermissionsRequestObject) (api.GetAllPermissionsResponseObject, error) {

	projectId := request.ProjectId
	var permissions []models.ProjectPermission = []models.ProjectPermission{}

	if err := s.Database.GetProjectPermissionsQuery(&permissions, projectId); err != nil {
		return api.GetAllPermissions500JSONResponse{}, err
	}

	var apiPermissions []api.ProjectPermission = []api.ProjectPermission{}
	for _, permission := range permissions {
		apiPermissions = append(apiPermissions, transform.GormProjectPermissionToApiProjectPermission(permission))
	}

	return api.GetAllPermissions200JSONResponse(apiPermissions), nil
}

// GetPermission implements api.StrictServerInterface.
func (s *Server) GetPermission(ctx context.Context, request api.GetPermissionRequestObject) (api.GetPermissionResponseObject, error) {

	projectId := request.ProjectId
	userUid, err := utilities.ParseUUID(request.UserId)
	if err != nil {
		return api.GetPermission500JSONResponse{}, err
	}

	var permission models.ProjectPermission

	if err := s.Database.GetProjectPermissionQuery(&permission, projectId, userUid); err != nil {
		return api.GetPermission500JSONResponse{}, err
	}

	apiPermission := transform.GormProjectPermissionToApiProjectPermission(permission)

	return api.GetPermission200JSONResponse(apiPermission), nil
}

// UpdatePermission implements api.StrictServerInterface.
func (s *Server) UpdatePermission(ctx context.Context, request api.UpdatePermissionRequestObject) (api.UpdatePermissionResponseObject, error) {
	panic("unimplemented")
}
