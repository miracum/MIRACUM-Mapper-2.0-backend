package server

import (
	"context"
	"errors"
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
)

// GetAllCodeSystemRoles implements api.StrictServerInterface.
func (s *Server) GetAllCodeSystemRoles(ctx context.Context, request api.GetAllCodeSystemRolesRequestObject) (api.GetAllCodeSystemRolesResponseObject, error) {
	projectId := request.ProjectId

	permissions, err := getUserPermissions(ctx, s, request.ProjectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.GetAllCodeSystemRoles404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.GetAllCodeSystemRoles500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(ProjectViewPermission, permissions) {
		return api.GetAllCodeSystemRoles403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to view code system roles in the project with ID %d", projectId))}, nil
	}

	var codeSystemRoles []models.CodeSystemRole = []models.CodeSystemRole{}
	if err := s.Database.GetAllCodeSystemRolesQuery(&codeSystemRoles, projectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetAllCodeSystemRoles404JSONResponse(err.Error()), nil
		default:
			return api.GetAllCodeSystemRoles500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the projects"}, err
		}
	}

	codeSystems := transform.GormCodeSystemRolesToApiCodeSystemRoles(&codeSystemRoles)

	return api.GetAllCodeSystemRoles200JSONResponse(*codeSystems), nil

}

// GetCodeSystemRole implements api.StrictServerInterface.
func (s *Server) GetCodeSystemRole(ctx context.Context, request api.GetCodeSystemRoleRequestObject) (api.GetCodeSystemRoleResponseObject, error) {
	projectId := request.ProjectId

	permissions, err := getUserPermissions(ctx, s, request.ProjectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.GetCodeSystemRole404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.GetCodeSystemRole500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(ProjectViewPermission, permissions) {
		return api.GetCodeSystemRole403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to view a code system role in the project with ID %d", projectId))}, nil
	}

	codeSystemRoleId := request.CodesystemRoleId
	var codeSystemRole models.CodeSystemRole = models.CodeSystemRole{}

	if err := s.Database.GetCodeSystemRoleQuery(&codeSystemRole, projectId, codeSystemRoleId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetCodeSystemRole404JSONResponse(err.Error()), nil
		default:
			return api.GetCodeSystemRole500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the projects"}, err
		}
	}

	codeSystem := transform.GormCodeSystemRoleToApiCodeSystemRole(&codeSystemRole)

	return api.GetCodeSystemRole200JSONResponse(*codeSystem), nil
}

// UpdateCodeSystemRole implements api.StrictServerInterface.
func (s *Server) UpdateCodeSystemRole(ctx context.Context, request api.UpdateCodeSystemRoleRequestObject) (api.UpdateCodeSystemRoleResponseObject, error) {
	projectId := request.ProjectId

	permissions, err := getUserPermissions(ctx, s, request.ProjectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.UpdateCodeSystemRole404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.UpdateCodeSystemRole500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(ProjectUpdatePermission, permissions) {
		return api.UpdateCodeSystemRole403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to update a code system role in the project with ID %d", projectId))}, nil
	}

	codeSystemRole := request.Body

	db_codeSystemRole := transform.ApiUpdateCodeSystemRoleToGormCodeSystemRole(codeSystemRole, &projectId)

	if err := s.Database.UpdateCodeSystemRoleQuery(db_codeSystemRole, projectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.UpdateCodeSystemRole404JSONResponse(api.BadRequestErrorJSONResponse(err.Error())), nil
		case errors.Is(err, database.ErrClientError):
			return api.UpdateCodeSystemRole400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		// case errors.Is(err, database.???) error for trying to update status-/equivalenceRequired
		default:
			return api.UpdateCodeSystemRole500JSONResponse{}, err
		}
	}

	return api.UpdateCodeSystemRole200JSONResponse(*transform.GormCodeSystemRoleToApiCodeSystemRole(db_codeSystemRole)), nil
}
