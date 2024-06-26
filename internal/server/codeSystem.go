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

// GetAllCodeSystems implements api.StrictServerInterface.
func (s *Server) GetAllCodeSystems(ctx context.Context, request api.GetAllCodeSystemsRequestObject) (api.GetAllCodeSystemsResponseObject, error) {
	var codeSystems []models.CodeSystem

	if err := s.Database.GetAllCodeSystemsQuery(&codeSystems); err != nil {
		// switch {
		// case errors.Is(err, database.ErrNotFound):
		// 	return api.GetAllCodeSystems404JSONResponse("No CodeSystems found"), nil
		// default:
		// 	return api.GetAllCodeSystems500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystems"}, err
		// }
		return api.GetAllCodeSystems500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystems"}, nil
	}

	var apiCodeSystems []api.CodeSystem = []api.CodeSystem{}
	for _, codeSystem := range codeSystems {
		apiCodeSystems = append(apiCodeSystems, *transform.GormCodeSystemToApiCodeSystem(&codeSystem))
	}

	return api.GetAllCodeSystems200JSONResponse(apiCodeSystems), nil
}

// CreateCodeSystem implements api.StrictServerInterface.
func (s *Server) CreateCodeSystem(ctx context.Context, request api.CreateCodeSystemRequestObject) (api.CreateCodeSystemResponseObject, error) {
	codeSystem := request.Body

	db_codeSystem := *transform.ApiCreateCodeSystemToGormCodeSystem(codeSystem)
	if err := s.Database.CreateCodeSystemQuery(&db_codeSystem); err != nil {
		return api.CreateCodeSystem500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to create the CodeSystem"}, nil
	}

	return api.CreateCodeSystem200JSONResponse(*transform.GormCodeSystemToApiCodeSystem(&db_codeSystem)), nil
	// TODO return 201 return api.CreateCodeSystem201JSONResponse(*transform.GormCodeSystemToApiCodeSystem(&db_codeSystem)), nil
}

// GetCodeSystem implements api.StrictServerInterface.
func (s *Server) GetCodeSystem(ctx context.Context, request api.GetCodeSystemRequestObject) (api.GetCodeSystemResponseObject, error) {
	codeSystemId := request.CodesystemId
	var codeSystem models.CodeSystem

	if err := s.Database.GetCodeSystemQuery(&codeSystem, codeSystemId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetCodeSystem404JSONResponse(fmt.Sprintf("CodeSystem with ID %d couldn't be found.", request.CodesystemId)), nil
		default:
			return api.GetCodeSystem500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystem"}, nil
		}
	}

	return api.GetCodeSystem200JSONResponse(*transform.GormCodeSystemToApiCodeSystem(&codeSystem)), nil

}

// UpdateCodeSystem implements api.StrictServerInterface.
func (s *Server) UpdateCodeSystem(ctx context.Context, request api.UpdateCodeSystemRequestObject) (api.UpdateCodeSystemResponseObject, error) {
	codeSystem := request.Body

	db_codeSystem := *transform.ApiCodeSystemToGormCodeSystem(codeSystem)
	if err := s.Database.UpdateCodeSystemQuery(&db_codeSystem); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.UpdateCodeSystem404JSONResponse(err.Error()), nil
		default:
			return api.UpdateCodeSystem500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to update the CodeSystem"}, nil
		}
	}

	return api.UpdateCodeSystem200JSONResponse(*transform.GormCodeSystemToApiCodeSystem(&db_codeSystem)), nil

}

// DeleteCodeSystem implements api.StrictServerInterface.
func (s *Server) DeleteCodeSystem(ctx context.Context, request api.DeleteCodeSystemRequestObject) (api.DeleteCodeSystemResponseObject, error) {
	codeSystemId := request.CodesystemId
	var codeSystem models.CodeSystem

	if err := s.Database.DeleteCodeSystemQuery(&codeSystem, codeSystemId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.DeleteCodeSystem404JSONResponse(err.Error()), nil
		default:
			return api.DeleteCodeSystem500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to delete the CodeSystem"}, nil
			// TODO or: return api.DeleteCodeSystem500JSONResponse{InternalServerErrorJSONResponse: err.Error()}, nil
			// TODO or: return api.DeleteCodeSystem500JSONResponse{InternalServerErrorJSONResponse: database.InternalServerErrorMessage}, nil
		}
	}

	return api.DeleteCodeSystem200JSONResponse(*transform.GormCodeSystemToApiCodeSystem(&codeSystem)), nil
}
