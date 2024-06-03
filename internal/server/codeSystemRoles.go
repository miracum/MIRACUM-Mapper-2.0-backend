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

	var codeSystemRoles []models.CodeSystemRole = []models.CodeSystemRole{}

	if err := s.Database.GetAllCodeSystemRolesQuery(&codeSystemRoles, request.ProjectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetAllCodeSystemRoles404JSONResponse(err.Error()), nil
		default:
			return api.GetAllCodeSystemRoles500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the projects"}, err
		}
	}

	codeSystems := transform.GormCodeSystemRolesToApiCodeSystemRoles(codeSystemRoles)

	return api.GetAllCodeSystemRoles200JSONResponse(codeSystems), nil

}

// GetCodeSystemRole implements api.StrictServerInterface.
func (s *Server) GetCodeSystemRole(ctx context.Context, request api.GetCodeSystemRoleRequestObject) (api.GetCodeSystemRoleResponseObject, error) {
	var codeSystemRole models.CodeSystemRole = models.CodeSystemRole{}

	if err := s.Database.GetCodeSystemRoleQuery(&codeSystemRole, request.ProjectId, request.CodeSystemRoleId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetCodeSystemRole404JSONResponse(err.Error()), nil
		default:
			return api.GetCodeSystemRole500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the projects"}, err
		}
	}

	codeSystem := transform.GormCodeSystemRoleToApiCodeSystemRole(codeSystemRole)

	return api.GetCodeSystemRole200JSONResponse(codeSystem), nil
}

// UpdateCodeSystemRole implements api.StrictServerInterface.
func (s *Server) UpdateCodeSystemRole(ctx context.Context, request api.UpdateCodeSystemRoleRequestObject) (api.UpdateCodeSystemRoleResponseObject, error) {
	codeSystemRole := request.Body
	codeSystemRoleId := request.CodeSystemRoleId
	projectID := request.ProjectId

	if codeSystemRole.Id == nil {
		codeSystemRole.Id = &request.CodeSystemRoleId
	} else if *codeSystemRole.Id != request.CodeSystemRoleId {
		return api.UpdateCodeSystemRole400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(fmt.Sprintf("CodeSystemRole ID %d in URL does not match CodeSystemRole ID %d in body", codeSystemRoleId, *codeSystemRole.Id))}, nil
	}

	db_codeSystemRole := transform.ApiCodeSystemRoleToGormCodeSystemRole(*codeSystemRole)

	checkFunc := func(oldCodeSystemRole, newCodeSystemRole *models.CodeSystemRole) error {
		if oldCodeSystemRole.CodeSystemID != newCodeSystemRole.CodeSystemID {
			return database.NewDBError(database.NotFound, fmt.Sprintf("Specified SystemID %d does not match existing SystemID %d for CodeSystemRole", newCodeSystemRole.CodeSystem.ID, oldCodeSystemRole.CodeSystem.ID))
		} else if oldCodeSystemRole.CodeSystem.Name != newCodeSystemRole.CodeSystem.Name {
			return database.NewDBError(database.NotFound, fmt.Sprintf("Specified System Name %s does not match existing System Name %s for CodeSystemRole", newCodeSystemRole.CodeSystem.Name, oldCodeSystemRole.CodeSystem.Name))
		} else if oldCodeSystemRole.CodeSystem.Version != newCodeSystemRole.CodeSystem.Version {
			return database.NewDBError(database.NotFound, fmt.Sprintf("Specified System Version %s does not match existing System Version %s for CodeSystemRole", newCodeSystemRole.CodeSystem.Version, oldCodeSystemRole.CodeSystem.Version))
		}
		return nil
	}

	if err := s.Database.UpdateCodeSystemRoleQuery(&db_codeSystemRole, projectID, codeSystemRoleId, checkFunc); err != nil {
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

	return api.UpdateCodeSystemRole200JSONResponse(*codeSystemRole), nil
}
