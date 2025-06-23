package server

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"mime/multipart"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/utilities"
)

// ImportCodeSystemVersionSnomed implements api.StrictServerInterface.
func (s *Server) ImportCodeSystemVersionSnomed(ctx context.Context, request api.ImportCodeSystemVersionSnomedRequestObject) (api.ImportCodeSystemVersionSnomedResponseObject, error) {
	codeSystemId := request.CodesystemId
	codeSystemVersionId := request.CodesystemVersionId

	// Check if the CodeSystem exists and is of type SNOMED
	var codeSystem models.CodeSystem
	if err := s.Database.GetCodeSystemQuery(&codeSystem, codeSystemId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.ImportCodeSystemVersionSnomed404JSONResponse(err.Error()), nil
		default:
			return api.ImportCodeSystemVersionSnomed500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystem"}, nil
		}
	}

	if codeSystem.Type != models.SNOMED_CT {
		return api.ImportCodeSystemVersionSnomed400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("This endpoint can only be used for SNOMED CT CodeSystems. Type is: " + string(codeSystem.Type))}, nil
	}

	// Check if the CodeSystemVersion exists and is not already imported
	var codeSystemVersion models.CodeSystemVersion
	if err := s.Database.GetCodeSystemVersionQuery(&codeSystemVersion, codeSystemId, codeSystemVersionId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.ImportCodeSystemVersionSnomed404JSONResponse(err.Error()), nil
		default:
			return api.ImportCodeSystemVersionSnomed500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystemVersion"}, nil
		}
	}

	if codeSystemVersion.Imported {
		return api.ImportCodeSystemVersionSnomed400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("CodeSystemVersion is already imported")}, nil
	}

	// Read the multipart form data
	form, err := request.Body.ReadForm(1024 * 1024 * 1024 * 10) // 10GB limit
	if err != nil {
		return api.ImportCodeSystemVersionSnomed500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to read the body"}, nil
	}

	if len(form.File["concept"]) == 0 {
		return api.ImportCodeSystemVersionSnomed400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("No concept file provided")}, nil
	}
	conceptFile := form.File["concept"][0]

	if len(form.File["description"]) == 0 {
		return api.ImportCodeSystemVersionSnomed400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("No description file provided")}, nil
	}
	descriptionFile := form.File["description"][0]

	return processFilesSnomed(conceptFile, descriptionFile, codeSystemId, codeSystemVersionId, s.Database)
}

func processFilesSnomed(concept *multipart.FileHeader, description *multipart.FileHeader, codeSystemId int32, codeSystemVersionId int32, db database.Datastore) (api.ImportCodeSystemVersionSnomedResponseObject, error) {
	// Check file types
	err := checkFileType(concept, "txt")
	if err != nil {
		return api.ImportCodeSystemVersionSnomed400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}
	err = checkFileType(description, "txt")
	if err != nil {
		return api.ImportCodeSystemVersionSnomed400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}

	// Validate concept file header
	conceptFile, err := concept.Open()
	if err != nil {
		return api.ImportCodeSystemVersionSnomed500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to open the concept file"}, nil
	}
	defer conceptFile.Close()

	conceptFileReader := csv.NewReader(conceptFile)
	conceptFileReader.Comma = '\t' // SNOMED files are tab-separated
	requiredColumns, optionalColumns := getCSVColumnsSnomedConcepts()
	columnsIndex, err := validateCSVHeader(conceptFileReader, requiredColumns, optionalColumns)
	if err != nil {
		return api.ImportCodeSystemVersionSnomed400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}
	csvIndexConcepts := getCSVIndexSnomedConcepts(columnsIndex)

	// Validate description file header
	descriptionFile, err := description.Open()
	if err != nil {
		return api.ImportCodeSystemVersionSnomed500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to open the description file"}, nil
	}
	defer descriptionFile.Close()

	descriptionFileReader := csv.NewReader(descriptionFile)
	descriptionFileReader.Comma = '\t' // SNOMED files are tab-separated
	requiredColumnsDescription, optionalColumnsDescription := getCSVColumnsSnomedDescriptions()
	columnsIndexDescription, err := validateCSVHeader(descriptionFileReader, requiredColumnsDescription, optionalColumnsDescription)
	if err != nil {
		return api.ImportCodeSystemVersionSnomed400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}
	csvIndexDescriptions := getCSVIndexSnomedDescriptions(columnsIndexDescription)

	// Read files
	descriptionFile2, err := description.Open()
	if err != nil {
		return api.ImportCodeSystemVersionSnomed500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to open the description file"}, nil
	}
	defer descriptionFile2.Close()
	descriptionFileReader2 := bufio.NewScanner(descriptionFile2)

	concepts, status, err := processCSVRowsSnomed(conceptFileReader, descriptionFileReader2, csvIndexConcepts, csvIndexDescriptions)
	if err != nil {
		switch status {
		case 400:
			return api.ImportCodeSystemVersionSnomed400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.ImportCodeSystemVersionSnomed500JSONResponse{InternalServerErrorJSONResponse: api.InternalServerErrorJSONResponse(err.Error())}, nil
		}
	}

	// Start the import process
	if !utilities.TryImporting() {
		return api.ImportCodeSystemVersionSnomed400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("An import is already in progress. Please wait until it is finished.")}, nil
	}
	go db.ImportConcepts(codeSystemId, codeSystemVersionId, concepts, nil)
	return api.ImportCodeSystemVersionSnomed202JSONResponse("Files processed successfully. Starting to create and update concepts in the background."), nil
}
