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

// GetAllPermissions implements api.StrictServerInterface.
func (s *Server) GetAllPermissions(ctx context.Context, request api.GetAllPermissionsRequestObject) (api.GetAllPermissionsResponseObject, error) {
	projectId := request.ProjectId
	var permissions []models.ProjectPermission = []models.ProjectPermission{}

	if err := s.Database.GetAllProjectPermissionsQuery(&permissions, projectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetAllPermissions404JSONResponse(err.Error()), nil
		default:
			return api.GetAllPermissions500JSONResponse{}, nil
		}
	}

	apiPermissions := transform.GormProjectPermissionsToApiProjectPermissions(&permissions)

	return api.GetAllPermissions200JSONResponse(*apiPermissions), nil
}

// CreatePermission implements api.StrictServerInterface.
func (s *Server) CreatePermission(ctx context.Context, request api.CreatePermissionRequestObject) (api.CreatePermissionResponseObject, error) {
	projectId := request.ProjectId
	permission := request.Body

	dbPermission, err := transform.ApiSendProjectPermissionToGormProjectPermission(permission, projectId)
	if err != nil {
		return api.CreatePermission422JSONResponse(err.Error()), nil
	}

	if err := s.Database.CreateProjectPermissionQuery(dbPermission); err != nil {
		switch {
		case errors.Is(err, database.ErrClientError):
			return api.CreatePermission422JSONResponse(err.Error()), nil
		default:
			return api.CreatePermission500JSONResponse{}, nil
		}
	}

	// TODO test if dbPermission contains all the information
	return api.CreatePermission200JSONResponse(*transform.GormProjectPermissionToApiProjectPermission(dbPermission)), nil
	// TODO return api.CreatePermission201JSONResponse{}, nil
}

// UpdatePermission implements api.StrictServerInterface.
func (s *Server) UpdatePermission(ctx context.Context, request api.UpdatePermissionRequestObject) (api.UpdatePermissionResponseObject, error) {
	projectId := request.ProjectId
	permission := request.Body

	dbPermission, err := transform.ApiSendProjectPermissionToGormProjectPermission(permission, projectId)
	if err != nil {
		return api.UpdatePermission422JSONResponse(err.Error()), nil
	}

	if err := s.Database.UpdateProjectPermissionQuery(dbPermission); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.UpdatePermission404JSONResponse(err.Error()), nil
		default:
			return api.UpdatePermission500JSONResponse{}, nil
		}
	}

	// TODO test that dbPermission contains all the information
	return api.UpdatePermission200JSONResponse(*transform.GormProjectPermissionToApiProjectPermission(dbPermission)), nil
}

// GetPermission implements api.StrictServerInterface.
func (s *Server) GetPermission(ctx context.Context, request api.GetPermissionRequestObject) (api.GetPermissionResponseObject, error) {
	projectId := request.ProjectId
	userUuid, err := utilities.ParseUUID(request.UserId)
	if err != nil {
		// TODO return 422 instead of 400
		return api.GetPermission400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(fmt.Sprintf("Invalid User ID: %s", err.Error()))}, nil
	}

	var permission models.ProjectPermission

	if err := s.Database.GetProjectPermissionQuery(&permission, projectId, userUuid); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetPermission404JSONResponse(err.Error()), nil
		default:
			// TODO returns empty body?
			return api.GetPermission500JSONResponse{}, nil
		}
	}

	return api.GetPermission200JSONResponse(*transform.GormProjectPermissionToApiProjectPermission(&permission)), nil
}

// DeletePermission implements api.StrictServerInterface.
func (s *Server) DeletePermission(ctx context.Context, request api.DeletePermissionRequestObject) (api.DeletePermissionResponseObject, error) {
	projectId := request.ProjectId
	userUuid, err := utilities.ParseUUID(request.UserId)
	if err != nil {
		// TODO return 422 instead of 400
		return api.DeletePermission400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(fmt.Sprintf("Invalid User ID: %s", err.Error()))}, nil
	}

	var permission models.ProjectPermission

	if err := s.Database.DeleteProjectPermissionQuery(&permission, projectId, userUuid); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.DeletePermission404JSONResponse(err.Error()), nil
		default:
			return api.DeletePermission500JSONResponse{}, nil
		}
	}

	return api.DeletePermission200JSONResponse(*transform.GormProjectPermissionToApiProjectPermission(&permission)), nil
}
