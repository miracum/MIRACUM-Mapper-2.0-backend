package server

import (
	"context"
	"errors"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
)

// GetAllCodeSystemRoles implements api.StrictServerInterface.
func (s *Server) GetAllCodeSystemRoles(ctx context.Context, request api.GetAllCodeSystemRolesRequestObject) (api.GetAllCodeSystemRolesResponseObject, error) {
	projectId := request.ProjectId
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
	codeSystemRole := request.Body

	db_codeSystemRole := transform.ApiUpdateCodeSystemRoleToGormCodeSystemRole(codeSystemRole, &projectId)

	if err := s.Database.UpdateCodeSystemRoleQuery(db_codeSystemRole, projectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.UpdateCodeSystemRole404JSONResponse(api.BadRequestErrorJSONResponse(err.Error())), nil
		// TODO
		case errors.Is(err, database.ErrClientError):
			return api.UpdateCodeSystemRole400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		// case errors.Is(err, database.???) error for trying to update status-/equivalenceRequired
		default:
			return api.UpdateCodeSystemRole500JSONResponse{}, err
		}
	}

	// TODO test if gorm returns full object after update and so everything is returned correctly
	return api.UpdateCodeSystemRole200JSONResponse(*transform.GormCodeSystemRoleToApiCodeSystemRole(db_codeSystemRole)), nil
}
