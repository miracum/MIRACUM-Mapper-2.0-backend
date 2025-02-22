package server

import (
	"context"
	"errors"

	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
)

// CreateCodeSystemVersion implements api.StrictServerInterface.
func (s *Server) CreateCodeSystemVersion(ctx context.Context, request api.CreateCodeSystemVersionRequestObject) (api.CreateCodeSystemVersionResponseObject, error) {
	codeSystemId := request.CodesystemId
	codeSystemVersion := request.Body

	db_codeSystemVersion := *transform.ApiBaseCodeSystemVersionToGormCodeSystemVersion(codeSystemVersion, codeSystemId)
	if err := s.Database.CreateCodeSystemVersionQuery(&db_codeSystemVersion); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.CreateCodeSystemVersion404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.CreateCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.CreateCodeSystemVersion500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to create the CodeSystemVersion"}, nil
		}
	}

	return api.CreateCodeSystemVersion200JSONResponse(*transform.GormCodeSystemVersionToApiCodeSystemVersion(&db_codeSystemVersion)), nil
}

// UpdateCodeSystemVersion implements api.StrictServerInterface.
func (s *Server) UpdateCodeSystemVersion(ctx context.Context, request api.UpdateCodeSystemVersionRequestObject) (api.UpdateCodeSystemVersionResponseObject, error) {
	codeSystemId := request.CodesystemId
	codeSystemVersion := request.Body

	db_codeSystemVersion := *transform.ApiUpdateCodeSystemVersionToGormCodeSystemVersion(codeSystemVersion, codeSystemId)
	if err := s.Database.UpdateCodeSystemVersionQuery(&db_codeSystemVersion); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.UpdateCodeSystemVersion404JSONResponse(err.Error()), nil
		default:
			return api.UpdateCodeSystemVersion500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to update the CodeSystemVersion"}, nil
		}
	}

	return api.UpdateCodeSystemVersion200JSONResponse(*transform.GormCodeSystemVersionToApiCodeSystemVersion(&db_codeSystemVersion)), nil
}

// DeleteCodeSystemVersion implements api.StrictServerInterface.
func (s *Server) DeleteCodeSystemVersion(ctx context.Context, request api.DeleteCodeSystemVersionRequestObject) (api.DeleteCodeSystemVersionResponseObject, error) {
	codeSystemVersionId := request.CodesystemVersionId
	var codeSystemVersion models.CodeSystemVersion

	if err := s.Database.DeleteCodeSystemVersionQuery(&codeSystemVersion, codeSystemVersionId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.DeleteCodeSystemVersion404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.DeleteCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.DeleteCodeSystemVersion500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to delete the CodeSystemVersion"}, nil
		}
	}

	return api.DeleteCodeSystemVersion200JSONResponse(*transform.GormCodeSystemVersionToApiCodeSystemVersion(&codeSystemVersion)), nil
}

// ImportCodeSystemVersion implements api.StrictServerInterface.
func (s *Server) ImportCodeSystemVersion(ctx context.Context, request api.ImportCodeSystemVersionRequestObject) (api.ImportCodeSystemVersionResponseObject, error) {
	codeSystemId := request.CodesystemId
	codeSystemVersionId := request.CodesystemVersionId

	var codeSystem models.CodeSystem
	if err := s.Database.GetCodeSystemQuery(&codeSystem, codeSystemId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.ImportCodeSystemVersion404JSONResponse(err.Error()), nil
		default:
			return api.ImportCodeSystemVersion500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystem"}, nil
		}
	}

	var codeSystemVersion models.CodeSystemVersion
	if err := s.Database.GetCodeSystemVersionQuery(&codeSystemVersion, codeSystemId, codeSystemVersionId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.ImportCodeSystemVersion404JSONResponse(err.Error()), nil
		default:
			return api.ImportCodeSystemVersion500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystemVersion"}, nil
		}
	}

	if codeSystemVersion.Imported {
		return api.ImportCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("CodeSystemVersion is already imported")}, nil
	}

	file, err := request.Body.NextPart()
	if err != nil {
		return api.ImportCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Error reading the file")}, nil
	}
	codeSystemType := codeSystem.Type
	return processFile(file, codeSystemId, codeSystemVersionId, codeSystemType, s.Database)
}
