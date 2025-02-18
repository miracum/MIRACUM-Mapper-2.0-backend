package server

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
)

// Structs for processing CSV files

type ConceptImport struct {
	Code        string
	Display     string
	Description *string
	Status      models.ConceptStatus
}

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
		return api.ImportCodeSystemVersion500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while reading the CSV file"}, nil
	}
	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))

	// Validate the header
	requiredColumns, optionalColumns := getImportColumns(codeSystemType)
	columnsIndex, err := validateCSVHeader(reader, requiredColumns, optionalColumns)
	if err != nil {
		return api.ImportCodeSystemVersion400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}

	// Process the rows
	go processCSVRows(reader, codeSystemId, codeSystemVersionId, codeSystemType, columnsIndex, db)
	return api.ImportCodeSystemVersion202JSONResponse("CSV file is being processed"), nil
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

func processCSVRows(reader *csv.Reader, codeSystemId int32, codeSystemVersionId int32, codeSystemType models.CodeSystemType, columnsIndex map[string]int, db database.Datastore) {
	csvIndex := getCsvIndex(codeSystemType, columnsIndex)

	var concepts []ConceptImport

	for {
		record, err := reader.Read()
		if err != nil {
			if err == csv.ErrFieldCount {
				log.Printf("CSV file has inconsistent number of fields")
				return
			}
			if err == io.EOF {
				break
			}
			log.Printf("Error reading CSV file: %v", err)
			return
		}

		var description *string
		if csvIndex.description != -1 && record[csvIndex.description] != "" {
			description = &record[csvIndex.description]
		}

		conceptImport := ConceptImport{
			Code:        record[csvIndex.code],
			Display:     getDisplayName(csvIndex.display, record),
			Description: description,
			Status:      getConceptStatus(codeSystemType, record[csvIndex.status]),
		}

		concepts = append(concepts, conceptImport)
	}

	log.Printf("CSV file processed successfully for CodeSystemID: %d", codeSystemId)
	createConcepts(codeSystemId, codeSystemVersionId, &concepts, db)

	if err := db.SetCodeSystemVersionImported(codeSystemVersionId, true); err != nil {
		log.Printf("Error setting CodeSystemVersion as imported: %v", err)
	}
}

func createConcepts(codeSystemId int32, codeSystemVersionId int32, concepts *[]ConceptImport, db database.Datastore) {
	for _, concept := range *concepts {
		neighborConcepts, err := db.GetNeighborConceptsQuery(concept.Code, codeSystemId, codeSystemVersionId)
		if err != nil {
			log.Printf("Error getting neighbor concepts: %v", err)
			return
		}
		switch neighborConcepts.NeighborType {
		case database.NeighborConceptsTypeNone:
			createNewConcept(codeSystemId, codeSystemVersionId, &concept, db)
		case database.NeighborConceptsTypeBefore:
			beforeConcept := neighborConcepts.BeforeConcept
			if !conceptsAreEqual(&concept, beforeConcept) {
				createNewConcept(codeSystemId, codeSystemVersionId, &concept, db)
			} else {
				beforeConcept.ValidToVersionID = uint32(codeSystemVersionId)
				if err := db.UpdateConceptQuery(beforeConcept); err != nil {
					log.Printf("Error updating concept: %v", err)
				}
			}
		case database.NeighborConceptsTypeAfter:
			afterConcept := neighborConcepts.AfterConcept
			if !conceptsAreEqual(&concept, afterConcept) {
				createNewConcept(codeSystemId, codeSystemVersionId, &concept, db)
			} else {
				afterConcept.ValidFromVersionID = uint32(codeSystemVersionId)
				if err := db.UpdateConceptQuery(afterConcept); err != nil {
					log.Printf("Error updating concept: %v", err)
				}
			}
		case database.NeighborConceptsTypeBeforeAndAfter:
			beforeConcept := neighborConcepts.BeforeConcept
			afterConcept := neighborConcepts.AfterConcept
			if !conceptsAreEqual(&concept, beforeConcept) && !conceptsAreEqual(&concept, afterConcept) {
				createNewConcept(codeSystemId, codeSystemVersionId, &concept, db)
			} else if conceptsAreEqual(&concept, beforeConcept) && !conceptsAreEqual(&concept, afterConcept) {
				beforeConcept.ValidToVersionID = uint32(codeSystemVersionId)
				if err := db.UpdateConceptQuery(beforeConcept); err != nil {
					log.Printf("Error updating concept: %v", err)
				}
			} else if !conceptsAreEqual(&concept, beforeConcept) && conceptsAreEqual(&concept, afterConcept) {
				afterConcept.ValidFromVersionID = uint32(codeSystemVersionId)
				if err := db.UpdateConceptQuery(afterConcept); err != nil {
					log.Printf("Error updating concept: %v", err)
				}
			} else {
				log.Printf("Error: Concept is before and after")
			}
		case database.NeighborConceptsTypeSurrounding:
			surroundingConcept := neighborConcepts.SurroundingConcept
			if !conceptsAreEqual(&concept, surroundingConcept) {
				// TODO
			} else {
				// Do nothing
			}
		}
	}
}

// Helper functions

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

func createNewConcept(codeSystemId int32, codeSystemVersionId int32, concept *ConceptImport, db database.Datastore) {
	newConcept := models.Concept{
		Code:               concept.Code,
		Display:            concept.Display,
		Description:        concept.Description,
		Status:             concept.Status,
		CodeSystemID:       uint32(codeSystemId),
		ValidFromVersionID: uint32(codeSystemVersionId),
		ValidToVersionID:   uint32(codeSystemVersionId),
	}
	if err := db.CreateConceptQuery(&newConcept); err != nil {
		log.Printf("Error creating concept: %v", err)
	}
}

func conceptsAreEqual(conceptImport *ConceptImport, conceptDB *models.Concept) bool {
	var descriptionsAreEqual bool
	if conceptImport.Description == nil || conceptDB.Description == nil {
		descriptionsAreEqual = conceptImport.Description == conceptDB.Description
	} else {
		descriptionsAreEqual = *conceptImport.Description == *conceptDB.Description
	}
	return conceptImport.Display == conceptDB.Display && descriptionsAreEqual && conceptImport.Status == conceptDB.Status
}
