package server

import (
	"context"
	"errors"

	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
	"miracummapper/internal/utilities"
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

	return api.CreateCodeSystemVersion200JSONResponse(*transform.GormCodeSystemVersionToApiCodeSystemVersion(&db_codeSystemVersion, false)), nil
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

	return api.UpdateCodeSystemVersion200JSONResponse(*transform.GormCodeSystemVersionToApiCodeSystemVersion(&db_codeSystemVersion, false)), nil
}

// DeleteCodeSystemVersion implements api.StrictServerInterface.
func (s *Server) DeleteCodeSystemVersion(ctx context.Context, request api.DeleteCodeSystemVersionRequestObject) (api.DeleteCodeSystemVersionResponseObject, error) {
	codeSystemVersionId := request.CodesystemVersionId
	codeSystemId := request.CodesystemId
	var codeSystemVersion models.CodeSystemVersion

	if err := s.Database.DeleteCodeSystemVersionQuery(&codeSystemVersion, codeSystemId, codeSystemVersionId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.DeleteCodeSystemVersion404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.DeleteCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.DeleteCodeSystemVersion500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to delete the CodeSystemVersion"}, nil
		}
	}

	return api.DeleteCodeSystemVersion200JSONResponse(*transform.GormCodeSystemVersionToApiCodeSystemVersion(&codeSystemVersion, false)), nil
}

// GetImportStatus implements api.StrictServerInterface.
func (s *Server) GetImportStatus(ctx context.Context, request api.GetImportStatusRequestObject) (api.GetImportStatusResponseObject, error) {
	importStatus := utilities.GetImportStatus()
	var errorStringPointer *string
	if importStatus.Error != nil {
		errorString := importStatus.Error.Error()
		errorStringPointer = &errorString
	} else {
		errorStringPointer = nil
	}
	return api.GetImportStatus200JSONResponse(api.GetImportStatus200JSONResponse{Progress: importStatus.Progress, Running: importStatus.Running, Error: errorStringPointer}), nil
}
