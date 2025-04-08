package server

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/utilities"
	"path/filepath"
)

// Structs for processing JSON files

type CodeSystemJSON struct {
	ResourceType string        `json:"resourceType"`
	Concepts     []ConpectJSON `json:"concept"`
}

type ConpectJSON struct {
	Code       string  `json:"code"`
	Display    *string `json:"display"`
	Definition *string `json:"definition"`
}

// Structs for processing CSV files

type csvIndex struct {
	code        int
	display     []int
	description int
	status      int
}

// CodeSystemType specific functions

func getConceptStatus(codeSystemType models.CodeSystemType, status string) models.ConceptStatus {
	switch codeSystemType {
	case models.GENERIC, models.LOINC:
		switch status {
		case "ACTIVE":
			return models.ActiveConcept
		case "TRIAL":
			return models.Trial
		case "DEPRECATED":
			return models.Deprecated
		case "DISCOURAGED":
			return models.Discouraged
		default:
			return models.ActiveConcept
		}
	default:
		return models.ActiveConcept
	}
}

func getCSVColumns(codeSystemType models.CodeSystemType) ([]string, []string) {
	switch codeSystemType {
	case models.GENERIC:
		return []string{"code", "display", "status"}, []string{"description"}
	case models.LOINC:
		return []string{"LOINC_NUM", "SHORTNAME", "LONG_COMMON_NAME", "STATUS", "DefinitionDescription"}, []string{}
	default:
		return nil, nil
	}
}

func getCSVIndex(codeSystemType models.CodeSystemType, columnsIndex map[string]int) csvIndex {
	switch codeSystemType {
	case models.GENERIC:
		return csvIndex{
			code:        columnsIndex["code"],
			display:     []int{columnsIndex["display"]},
			description: columnsIndex["description"],
			status:      columnsIndex["status"],
		}
	case models.LOINC:
		return csvIndex{
			code:        columnsIndex["LOINC_NUM"],
			display:     []int{columnsIndex["LONG_COMMON_NAME"], columnsIndex["SHORTNAME"]},
			description: columnsIndex["DefinitionDescription"],
			status:      columnsIndex["STATUS"],
		}
	default:
		return csvIndex{}
	}
}

// main function for processing the file

func processFile(file *multipart.Part, codeSystemId int32, codeSystemVersionId int32, codeSystemType models.CodeSystemType, db database.Datastore) (api.ImportCodeSystemVersionResponseObject, error) {
	fileType, err := getFileType(file)
	if err != nil {
		return api.ImportCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}

	if fileType == "csv" {
		switch codeSystemType {
		case models.GENERIC, models.LOINC:
			return processCSVFile(file, codeSystemId, codeSystemVersionId, codeSystemType, db)
		default:
			return api.ImportCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("The import of a csv file is currently not supported for this code system type")}, nil
		}
	} else if fileType == "json" {
		return api.ImportCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Please use the import-json endpoint for uploading a json file")}, nil
		// switch codeSystemType {
		// case models.GENERIC, models.ICD_10_GM:
		// 	return processJSONFile(file, codeSystemId, codeSystemVersionId, codeSystemType, db)
		// default:
		// 	return api.ImportCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("The import of a json file is currently not supported for this code system type")}, nil
		// }
	} else {
		return api.ImportCodeSystemVersion500JSONResponse{InternalServerErrorJSONResponse: api.InternalServerErrorJSONResponse("An error occurred while trying to determine the file type")}, nil
	}
}

// helper functions

func getFileType(file *multipart.Part) (string, error) {
	filename := file.FileName()
	fileExtension := filepath.Ext(filename)
	if fileExtension == ".csv" {
		return "csv", nil
	} else if fileExtension == ".json" {
		return "json", nil
	}

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "text/csv" {
		return "csv", nil
	} else if mimeType == "application/json" {
		return "json", nil
	}

	if mimeType == "" && fileExtension == "" {
		return "", fmt.Errorf("no content type or filename with extension provided")
	} else if mimeType == "" {
		return "", fmt.Errorf("unsupported file extension: %s", fileExtension)
	} else if fileExtension == "" {
		return "", fmt.Errorf("unsupported content type: %s", mimeType)
	} else {
		return "", fmt.Errorf("unsupported content type: %s and file extension: %s", mimeType, fileExtension)
	}
}

func getDisplayName(displayIndex []int, record []string) string {
	displayName := ""
	for _, index := range displayIndex {
		if name := record[index]; name != "" {
			if displayName == "" {
				displayName = name
			} else {
				displayName = fmt.Sprintf("%s | %s", displayName, record[index])
			}
		}
	}
	return displayName
}

// functions for processing CSV files

func processCSVFile(file *multipart.Part, codeSystemId int32, codeSystemVersionId int32, codeSystemType models.CodeSystemType, db database.Datastore) (api.ImportCodeSystemVersionResponseObject, error) {
	reader := csv.NewReader(file)
	requiredColumns, optionalColumns := getCSVColumns(codeSystemType)
	columnsIndex, err := validateCSVHeader(reader, requiredColumns, optionalColumns)
	if err != nil {
		return api.ImportCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}
	csvIndex := getCSVIndex(codeSystemType, columnsIndex)

	return processCSVRows(reader, codeSystemId, codeSystemVersionId, codeSystemType, csvIndex, db)
}

func validateCSVHeader(reader *csv.Reader, requiredColumns []string, optionalColumns []string) (map[string]int, error) {
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading the CSV header: %v", err)
	}

	columnsIndex := make(map[string]int)
	for _, column := range requiredColumns {
		columnsIndex[column] = -1
	}
	for _, column := range optionalColumns {
		columnsIndex[column] = -1
	}

	for i, column := range header {
		if j, exists := columnsIndex[column]; j == -1 && exists {
			columnsIndex[column] = i
		} else if j != -1 && exists {
			return nil, fmt.Errorf("error: column found multiple times in csv file: %s", column)
		}
	}

	for column, present := range columnsIndex {
		if present == -1 {
			for _, requiredColumn := range requiredColumns {
				if column == requiredColumn {
					return nil, fmt.Errorf("missing required column: %s", column)
				}
			}
		}
	}

	return columnsIndex, nil
}

func processCSVRows(reader *csv.Reader, codeSystemId int32, codeSystemVersionId int32, codeSystemType models.CodeSystemType, csvIndex csvIndex, db database.Datastore) (api.ImportCodeSystemVersionResponseObject, error) {
	var concepts []database.ConceptImport

	for {
		record, err := reader.Read()
		if err != nil {
			if err == csv.ErrFieldCount {
				return api.ImportCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("CSV file has inconsistent number of fields")}, nil
			}
			if err == io.EOF {
				break
			}
			return api.ImportCodeSystemVersion500JSONResponse{InternalServerErrorJSONResponse: api.InternalServerErrorJSONResponse(fmt.Sprintf("An Error occurred while reading the CSV file: %v", err))}, nil
		}

		var description *string
		if csvIndex.description != -1 && record[csvIndex.description] != "" {
			description = &record[csvIndex.description]
		}

		conceptImport := database.ConceptImport{
			Code:        record[csvIndex.code],
			Display:     getDisplayName(csvIndex.display, record),
			Description: description,
			Status:      getConceptStatus(codeSystemType, record[csvIndex.status]),
		}

		concepts = append(concepts, conceptImport)
	}

	if !utilities.TryImporting() {
		return api.ImportCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("An import is already in progress. Please wait until it is finished.")}, nil
	}
	go db.CreateConcepts(codeSystemId, codeSystemVersionId, &concepts)
	return api.ImportCodeSystemVersion202JSONResponse("CSV file processed successfully. Starting to create and update concepts in the background."), nil
}

// functions for processing JSON files

func processJSONFile(file *api.ImportCodeSystemVersionJsonJSONRequestBody, codeSystemId int32, codeSystemVersionId int32, codeSystemType models.CodeSystemType, db database.Datastore) (api.ImportCodeSystemVersionJsonResponseObject, error) {
	switch codeSystemType {
	case models.ICD_10_GM:
		break
	default:
		return api.ImportCodeSystemVersionJson400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("The import of a json file is currently not supported for this code system type")}, nil
	}

	if file == nil {
		return api.ImportCodeSystemVersionJson400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("No body provided")}, nil
	}

	if resourceType, ok := (*file)["resourceType"].(string); !ok || resourceType != "CodeSystem" {
		return api.ImportCodeSystemVersionJson400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Incorrect FHIR resource type, expected 'CodeSystem'")}, nil
	}

	var concepts []database.ConceptImport

	if conceptsList, ok := (*file)["concept"].([]any); ok {
		for _, concept := range conceptsList {
			if conceptMap, ok := concept.(map[string]any); ok {
				code, ok := conceptMap["code"].(string)
				if !ok || code == "" {
					return api.ImportCodeSystemVersionJson400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Code is missing in JSON file")}, nil
				}
				display, ok := conceptMap["display"].(string)
				if !ok || display == "" {
					return api.ImportCodeSystemVersionJson400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Display is missing in JSON file")}, nil
				}
				newConcept := database.ConceptImport{
					Code:    code,
					Display: display,
					Status:  models.ActiveConcept,
				}
				concepts = append(concepts, newConcept)
			} else {
				return api.ImportCodeSystemVersionJson400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Invalid concept format in JSON file")}, nil
			}
		}
	} else {
		return api.ImportCodeSystemVersionJson400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Invalid concept format in JSON file")}, nil
	}

	if !utilities.TryImporting() {
		return api.ImportCodeSystemVersionJson400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("An import is already in progress. Please wait until it is finished.")}, nil
	}
	go db.CreateConcepts(codeSystemId, codeSystemVersionId, &concepts)
	return api.ImportCodeSystemVersionJson202JSONResponse("JSON file processed successfully. Starting to create and update concepts in the background."), nil
}
