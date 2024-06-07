package server

import (
	"context"
	"errors"
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
	"miracummapper/internal/utilities"
)

// AddPermission implements api.StrictServerInterface.
func (s *Server) AddPermission(ctx context.Context, request api.AddPermissionRequestObject) (api.AddPermissionResponseObject, error) {

	projectId := request.ProjectId
	permission := *request.Body

	dbPermission, err := transform.ApiProjectPermissionToGormProjectPermission(permission, projectId)
	if err != nil {
		return api.AddPermission422JSONResponse(err.Error()), nil
	}

	if err := s.Database.CreateProjectPermissionQuery(&dbPermission); err != nil {
		switch {
		case errors.Is(err, database.ErrClientError):
			return api.AddPermission422JSONResponse(err.Error()), err
		default:
			return api.AddPermission500JSONResponse{}, err
		}
	}

	return api.AddPermission200JSONResponse(permission), nil
	//TODO return api.AddPermission201JSONResponse{}, nil
}

// DeletePermission implements api.StrictServerInterface.
func (s *Server) DeletePermission(ctx context.Context, request api.DeletePermissionRequestObject) (api.DeletePermissionResponseObject, error) {

	projectId := request.ProjectId
	userUuid, err := utilities.ParseUUID(request.UserId)
	if err != nil {
		// TODO return 422 instead of 400
		return api.DeletePermission400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(fmt.Sprintf("Invalid User ID: %s", err.Error()))}, err
	}

	var permission models.ProjectPermission

	if err := s.Database.DeleteProjectPermissionQuery(&permission, projectId, userUuid); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.DeletePermission404JSONResponse(err.Error()), err
		default:
			return api.DeletePermission500JSONResponse{}, err
		}
	}

	api_permission := transform.GormProjectPermissionToApiProjectPermission(permission)
	return api.DeletePermission200JSONResponse(api_permission), nil
}

// GetAllPermissions implements api.StrictServerInterface.
func (s *Server) GetAllPermissions(ctx context.Context, request api.GetAllPermissionsRequestObject) (api.GetAllPermissionsResponseObject, error) {

	projectId := request.ProjectId
	var permissions []models.ProjectPermission = []models.ProjectPermission{}

	if err := s.Database.GetProjectPermissionsQuery(&permissions, projectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetAllPermissions404JSONResponse(err.Error()), err
		default:
			return api.GetAllPermissions500JSONResponse{}, err
		}
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
	userUuid, err := utilities.ParseUUID(request.UserId)
	if err != nil {
		// TODO return 422 instead of 400
		return api.GetPermission400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(fmt.Sprintf("Invalid User ID: %s", err.Error()))}, err
	}

	var permission models.ProjectPermission

	if err := s.Database.GetProjectPermissionQuery(&permission, projectId, userUuid); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetPermission404JSONResponse(err.Error()), err
		default:
			// TODO returns empty body?
			return api.GetPermission500JSONResponse{}, err
		}
	}

	apiPermission := transform.GormProjectPermissionToApiProjectPermission(permission)

	return api.GetPermission200JSONResponse(apiPermission), nil
}

// UpdatePermission implements api.StrictServerInterface.
func (s *Server) UpdatePermission(ctx context.Context, request api.UpdatePermissionRequestObject) (api.UpdatePermissionResponseObject, error) {
	projectId := request.ProjectId
	permission := *request.Body

	dbPermission, err := transform.ApiProjectPermissionToGormProjectPermission(permission, projectId)
	if err != nil {
		return api.UpdatePermission422JSONResponse(err.Error()), nil
	}

	if err := s.Database.UpdateProjectPermissionQuery(&dbPermission, projectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.UpdatePermission404JSONResponse(err.Error()), err
		default:
			return api.UpdatePermission500JSONResponse{}, err
		}
	}

	return api.UpdatePermission200JSONResponse(permission), nil
}
