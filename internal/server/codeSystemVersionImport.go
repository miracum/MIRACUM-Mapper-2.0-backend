package server

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/utilities"
)

// Structs for processing CSV files

type csvIndex struct {
	code        int
	display     []int
	description int
	status      int
}

// CodeSystemType specific functions

func getImportColumns(codeSystemType models.CodeSystemType) ([]string, []string) {
	switch codeSystemType {
	case models.GENERIC:
		return []string{"code", "display", "status"}, []string{"description"}
	case models.LOINC:
		return []string{"LOINC_NUM", "SHORTNAME", "LONG_COMMON_NAME", "STATUS", "DefinitionDescription"}, []string{}
	default:
		return nil, nil
	}
}

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

func getCsvIndex(codeSystemType models.CodeSystemType, columnsIndex map[string]int) csvIndex {
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

// main functions for processing CSV files

func processFile(file io.Reader, codeSystemId int32, codeSystemVersionId int32, codeSystemType models.CodeSystemType, db database.Datastore) (api.ImportCodeSystemVersionResponseObject, error) {
	// Read the entire CSV content into a buffer and create a new reader
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return api.ImportCodeSystemVersion500JSONResponse{InternalServerErrorJSONResponse: api.InternalServerErrorJSONResponse(fmt.Sprintf("An Error occurred while reading the CSV file: %v", err))}, nil
	}
	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))

	// Validate the header
	requiredColumns, optionalColumns := getImportColumns(codeSystemType)
	columnsIndex, err := validateCSVHeader(reader, requiredColumns, optionalColumns)
	if err != nil {
		return api.ImportCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}

	// Process the rows
	return processCSVRows(reader, codeSystemId, codeSystemVersionId, codeSystemType, columnsIndex, db)
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

func processCSVRows(reader *csv.Reader, codeSystemId int32, codeSystemVersionId int32, codeSystemType models.CodeSystemType, columnsIndex map[string]int, db database.Datastore) (api.ImportCodeSystemVersionResponseObject, error) {
	csvIndex := getCsvIndex(codeSystemType, columnsIndex)

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
