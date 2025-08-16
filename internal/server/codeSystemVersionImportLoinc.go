package server

import (
	"context"
	"encoding/csv"
	"errors"
	"mime/multipart"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/utilities"
)

// ImportCodeSystemVersionLoinc implements api.StrictServerInterface.
func (s *Server) ImportCodeSystemVersionLoinc(ctx context.Context, request api.ImportCodeSystemVersionLoincRequestObject) (api.ImportCodeSystemVersionLoincResponseObject, error) {
	codeSystemId := request.CodesystemId
	codeSystemVersionId := request.CodesystemVersionId

	// Check if the CodeSystem exists and is of type LOINC
	var codeSystem models.CodeSystem
	if err := s.Database.GetCodeSystemQuery(&codeSystem, codeSystemId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.ImportCodeSystemVersionLoinc404JSONResponse(err.Error()), nil
		default:
			return api.ImportCodeSystemVersionLoinc500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystem"}, nil
		}
	}

	if codeSystem.Type != models.LOINC {
		return api.ImportCodeSystemVersionLoinc400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("This endpoint can only be used for LOINC CodeSystems. Type is: " + string(codeSystem.Type))}, nil
	}

	// Check if the CodeSystemVersion exists and is not already imported
	var codeSystemVersion models.CodeSystemVersion
	if err := s.Database.GetCodeSystemVersionQuery(&codeSystemVersion, codeSystemId, codeSystemVersionId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.ImportCodeSystemVersionLoinc404JSONResponse(err.Error()), nil
		default:
			return api.ImportCodeSystemVersionLoinc500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystemVersion"}, nil
		}
	}

	if codeSystemVersion.Imported {
		return api.ImportCodeSystemVersionLoinc400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("CodeSystemVersion is already imported")}, nil
	}

	// Read the multipart form data
	form, err := request.Body.ReadForm(1024 * 1024 * 1024) // 1GB limit
	if err != nil {
		return api.ImportCodeSystemVersionLoinc500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to read the body"}, nil
	}

	if len(form.File["loinc"]) == 0 {
		return api.ImportCodeSystemVersionLoinc400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("No loinc file provided")}, nil
	}
	mainFile := form.File["loinc"][0]

	if len(form.File["map_to"]) == 0 {
		return api.ImportCodeSystemVersionLoinc400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("No rmap to file provided")}, nil
	}
	replaceByFile := form.File["map_to"][0]

	return processFilesLoinc(mainFile, replaceByFile, codeSystemId, codeSystemVersionId, s.Database)
}

func processFilesLoinc(main *multipart.FileHeader, replaceBy *multipart.FileHeader, codeSystemId int32, codeSystemVersionId int32, db database.Datastore) (api.ImportCodeSystemVersionLoincResponseObject, error) {
	// Check file types
	err := checkFileType(main, "csv")
	if err != nil {
		return api.ImportCodeSystemVersionLoinc400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}
	err = checkFileType(replaceBy, "csv")
	if err != nil {
		return api.ImportCodeSystemVersionLoinc400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}

	// Validate main file header
	mainFile, err := main.Open()
	if err != nil {
		return api.ImportCodeSystemVersionLoinc500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to open the main file"}, nil
	}
	defer mainFile.Close()

	mainFileReader := csv.NewReader(mainFile)
	requiredColumns, optionalColumns := getCSVColumnsMain(models.LOINC)
	columnsIndex, err := validateCSVHeader(mainFileReader, requiredColumns, optionalColumns)
	if err != nil {
		return api.ImportCodeSystemVersionLoinc400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}
	csvIndexMain := getCSVIndexMain(models.LOINC, columnsIndex)

	// Read main file
	concepts, status, err := processCSVRowsMain(mainFileReader, csvIndexMain)
	if err != nil {
		switch status {
		case 400:
			return api.ImportCodeSystemVersionLoinc400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.ImportCodeSystemVersionLoinc500JSONResponse{InternalServerErrorJSONResponse: api.InternalServerErrorJSONResponse(err.Error())}, nil
		}
	}

	// Validate replace by file header
	replaceByFile, err := replaceBy.Open()
	if err != nil {
		return api.ImportCodeSystemVersionLoinc500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to open the replace by file"}, nil
	}
	defer replaceByFile.Close()

	replaceByFileReader := csv.NewReader(replaceByFile)
	requiredColumnsReplaceBy, optionalColumnsReplaceBy := getCSVColumnsReplaceBy(models.LOINC)
	columnsIndexReplaceBy, err := validateCSVHeader(replaceByFileReader, requiredColumnsReplaceBy, optionalColumnsReplaceBy)
	if err != nil {
		return api.ImportCodeSystemVersionLoinc400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}
	csvIndexReplaceBy := getCSVIndexReplaceBy(models.LOINC, columnsIndexReplaceBy)

	// Read replace by file
	replaceByItems, status, err := processCSVRowsReplaceBy(replaceByFileReader, csvIndexReplaceBy, codeSystemId)
	if err != nil {
		switch status {
		case 400:
			return api.ImportCodeSystemVersionLoinc400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.ImportCodeSystemVersionLoinc500JSONResponse{InternalServerErrorJSONResponse: api.InternalServerErrorJSONResponse(err.Error())}, nil
		}
	}

	// Start the import process
	if !utilities.TryImporting() {
		return api.ImportCodeSystemVersionLoinc400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("An import is already in progress. Please wait until it is finished.")}, nil
	}
	go db.ImportConcepts(codeSystemId, codeSystemVersionId, concepts, replaceByItems)
	return api.ImportCodeSystemVersionLoinc202JSONResponse("Files processed successfully. Starting to create and update concepts in the background."), nil
}
