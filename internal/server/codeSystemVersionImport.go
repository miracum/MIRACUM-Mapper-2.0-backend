package server

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"path/filepath"
	"strings"
)

// This file contains functions for processing CSV files for importing generic and LOINC code systems.

// Structs for processing CSV files

type csvIndexMain struct {
	code        int
	display     []int
	description int
	status      int
}

type csvIndexReplaceBy struct {
	code        int
	mapTo       int
	equivalence int
	comment     int
}

// Helper functions for conversion

func convertConceptStatus(status string) models.ConceptStatus {
	switch strings.ToLower(status) {
	case "active":
		return models.ActiveConcept
	case "trial":
		return models.Trial
	case "deprecated":
		return models.Deprecated
	case "discouraged":
		return models.Discouraged
	default:
		return models.ActiveConcept
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

// Function to check the file type based on the extension and content type
func checkFileType(file *multipart.FileHeader, expected string) error {
	filename := file.Filename
	fileExtension := filepath.Ext(filename)
	if fileExtension == fmt.Sprintf(".%s", expected) {
		return nil
	}

	mimeType := file.Header.Get("Content-Type")
	expectedMimeType := ""
	if expected == "csv" {
		expectedMimeType = "text/csv"
	} else if expected == "txt" {
		expectedMimeType = "text/plain"
	}
	if mimeType == expectedMimeType {
		return nil
	}

	if mimeType == "" && fileExtension == "" {
		return fmt.Errorf("no content type or filename with extension provided")
	} else if mimeType == "" {
		return fmt.Errorf("unsupported file extension: %s", fileExtension)
	} else if fileExtension == "" {
		return fmt.Errorf("unsupported content type: %s", mimeType)
	} else {
		return fmt.Errorf("unsupported content type: %s and file extension: %s", mimeType, fileExtension)
	}
}

// Functions for CSV header validation and column retrieval

func getCSVColumnsMain(codeSystemType models.CodeSystemType) ([]string, []string) {
	switch codeSystemType {
	case models.GENERIC:
		return []string{"code", "display", "status"}, []string{"description"}
	case models.LOINC:
		return []string{"LOINC_NUM", "SHORTNAME", "LONG_COMMON_NAME", "STATUS", "DefinitionDescription"}, []string{}
	default:
		return nil, nil
	}
}

func getCSVColumnsReplaceBy(codeSystemType models.CodeSystemType) ([]string, []string) {
	switch codeSystemType {
	case models.GENERIC:
		return []string{"code", "map_to"}, []string{"equivalence", "comment"}
	case models.LOINC:
		return []string{"LOINC", "MAP_TO", "COMMENT"}, []string{}
	default:
		return nil, nil
	}
}

func getCSVIndexMain(codeSystemType models.CodeSystemType, columnsIndex map[string]int) csvIndexMain {
	switch codeSystemType {
	case models.GENERIC:
		return csvIndexMain{
			code:        columnsIndex["code"],
			display:     []int{columnsIndex["display"]},
			description: columnsIndex["description"],
			status:      columnsIndex["status"],
		}
	case models.LOINC:
		return csvIndexMain{
			code:        columnsIndex["LOINC_NUM"],
			display:     []int{columnsIndex["LONG_COMMON_NAME"], columnsIndex["SHORTNAME"]},
			description: columnsIndex["DefinitionDescription"],
			status:      columnsIndex["STATUS"],
		}
	default:
		return csvIndexMain{}
	}
}

func getCSVIndexReplaceBy(codeSystemType models.CodeSystemType, columnsIndex map[string]int) csvIndexReplaceBy {
	switch codeSystemType {
	case models.GENERIC:
		return csvIndexReplaceBy{
			code:        columnsIndex["code"],
			mapTo:       columnsIndex["map_to"],
			equivalence: columnsIndex["equivalence"],
			comment:     columnsIndex["comment"],
		}
	case models.LOINC:
		return csvIndexReplaceBy{
			code:        columnsIndex["LOINC"],
			mapTo:       columnsIndex["MAP_TO"],
			equivalence: -1,
			comment:     columnsIndex["COMMENT"],
		}
	default:
		return csvIndexReplaceBy{}
	}
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

// Function to process CSV rows of the files

func processCSVRowsMain(reader *csv.Reader, csvIndex csvIndexMain) (*[]database.ConceptImport, int32, error) {
	var concepts []database.ConceptImport

	for {
		record, err := reader.Read()
		if err != nil {
			if err == csv.ErrFieldCount {
				return nil, 400, fmt.Errorf("Main CSV file has inconsistent number of fields")
			}
			if err == io.EOF {
				break
			}
			return nil, 500, fmt.Errorf("An Error occurred while reading the Main CSV file: %v", err)
		}

		var description *string
		if csvIndex.description != -1 && record[csvIndex.description] != "" {
			description = &record[csvIndex.description]
		}

		conceptImport := database.ConceptImport{
			Code:        record[csvIndex.code],
			Display:     getDisplayName(csvIndex.display, record),
			Description: description,
			Status:      convertConceptStatus(record[csvIndex.status]),
		}

		concepts = append(concepts, conceptImport)
	}

	return &concepts, 200, nil
}

func processCSVRowsReplaceBy(reader *csv.Reader, csvIndex csvIndexReplaceBy, codeSystemId int32) (*[]models.ConceptReplaceBy, int32, error) {
	var replaceByConcepts []models.ConceptReplaceBy

	for {
		record, err := reader.Read()
		if err != nil {
			if err == csv.ErrFieldCount {
				return nil, 400, fmt.Errorf("Replace by CSV file has inconsistent number of fields")
			}
			if err == io.EOF {
				break
			}
			return nil, 500, fmt.Errorf("An Error occurred while reading the Replace by CSV file: %v", err)
		}

		code := record[csvIndex.code]
		mapTo := record[csvIndex.mapTo]
		if code == "" || mapTo == "" || code == mapTo {
			continue // Skip invalid entries
		}

		var equivalence *models.ConceptReplaceByEquivalence = nil
		if csvIndex.equivalence != -1 && record[csvIndex.equivalence] != "" {
			equivalenceValue := models.ConceptReplaceByEquivalence(record[csvIndex.equivalence])
			equivalence = &equivalenceValue
		}

		var comment *string
		if csvIndex.comment != -1 && record[csvIndex.comment] != "" {
			comment = &record[csvIndex.comment]
		}
		replaceByConcept := models.ConceptReplaceBy{
			Code:         record[csvIndex.code],
			MapTo:        record[csvIndex.mapTo],
			CodeSystemID: codeSystemId,
			Equivalence:  equivalence,
			Comment:      comment,
		}

		replaceByConcepts = append(replaceByConcepts, replaceByConcept)
	}

	return &replaceByConcepts, 200, nil
}
