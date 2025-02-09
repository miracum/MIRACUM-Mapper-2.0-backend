package server

import (
	//"bytes"
	"context"
	//"encoding/csv"
	"errors"
	//"fmt"
	//"io"
	//"log"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
)

// type CSVProcessingJob struct {
// 	CodeSystemID int32
// 	Reader       *csv.Reader
// 	Database     database.Datastore
// }

// CreateCodeSystemVersion implements api.StrictServerInterface.
func (s *Server) CreateCodeSystemVersion(ctx context.Context, request api.CreateCodeSystemVersionRequestObject) (api.CreateCodeSystemVersionResponseObject, error) {
	codeSystemId := request.CodesystemId
	codeSystemVersion := request.Body

	db_codeSystemVersion := *transform.ApiBaseCodeSystemVersionToGormCodeSystemVersion(codeSystemVersion, codeSystemId)
	if err := s.Database.CreateCodeSystemVersionQuery(&db_codeSystemVersion); err != nil {
		return api.CreateCodeSystemVersion500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to create the CodeSystemVersion"}, nil
	}

	return api.CreateCodeSystemVersion200JSONResponse(*transform.GormCodeSystemVersionToApiCodeSystemVersion(&db_codeSystemVersion)), nil
}

// UpdateCodeSystemVersion implements api.StrictServerInterface.
func (s *Server) UpdateCodeSystemVersion(ctx context.Context, request api.UpdateCodeSystemVersionRequestObject) (api.UpdateCodeSystemVersionResponseObject, error) {
	codeSystemId := request.CodesystemId
	codeSystemVersion := request.Body

	db_codeSystemVersion := *transform.ApiCodeSystemVersionToGormCodeSystemVersion(codeSystemVersion, codeSystemId)
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
	panic("unimplemented")
}

// helper functions for import code system

// func validateCSVHeader(reader *csv.Reader) (map[string]int, error) {
// 	header, err := reader.Read()
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading the CSV header: %v", err)
// 	}

// 	requiredColumns := map[string]int{"code": -1, "meaning": -1}
// 	for i, column := range header {
// 		if j, exists := requiredColumns[column]; j == -1 && exists {
// 			requiredColumns[column] = i
// 		} else if j != -1 {
// 			return nil, fmt.Errorf("error: column found multiple times in csv file: %s", column)
// 		}
// 	}
// 	for column, present := range requiredColumns {
// 		if present == -1 {
// 			return nil, fmt.Errorf("missing required column: %s", column)
// 		}
// 	}

// 	return requiredColumns, nil
// }

// func processCSVRows(job CSVProcessingJob, requiredColumns map[string]int) {
// 	codeIndex := requiredColumns["code"]
// 	meaningIndex := requiredColumns["meaning"]

// 	var concepts []models.Concept

// 	for {
// 		record, err := job.Reader.Read()
// 		if err != nil {
// 			if err == csv.ErrFieldCount {
// 				log.Printf("CSV file has inconsistent number of fields")
// 				return
// 			}
// 			if err == io.EOF {
// 				break
// 			}
// 			log.Printf("Error reading CSV file: %v", err)
// 			return
// 		}

// 		concept := models.Concept{
// 			Code:         record[codeIndex],
// 			Display:      record[meaningIndex],
// 			CodeSystemID: uint32(job.CodeSystemID),
// 		}
// 		concepts = append(concepts, concept)

// 		// fmt.Printf("Code: %s, Meaning: %s\n", record[codeIndex], record[meaningIndex])
// 	}
// 	if err := job.Database.CreateConceptsQuery(&concepts); err != nil {
// 		log.Printf("Error inserting concepts into database: %v", err)
// 		return
// 	}

// 	log.Printf("CSV file processed successfully for CodeSystemID: %d", job.CodeSystemID)
// }

// // ImportCodeSystem implements api.StrictServerInterface
// func (s *Server) ImportCodeSystem(ctx context.Context, request api.ImportCodeSystemRequestObject) (api.ImportCodeSystemResponseObject, error) {
// 	codeSystemId := request.CodesystemId
// 	var codeSystem models.CodeSystem
// 	var concept models.Concept

// 	if err := s.Database.GetFirstElementCodeSystemQuery(&codeSystem, codeSystemId, &concept); err != nil {
// 		if errors.Is(err, database.ErrNotFound) && codeSystem.ID == 0 {
// 			return api.ImportCodeSystem404JSONResponse(fmt.Sprintf("CodeSystem with ID %d couldn't be found.", request.CodesystemId)), nil
// 		} else if errors.Is(err, database.ErrNotFound) && codeSystem.ID != 0 && concept.ID == 0 {
// 		} else {
// 			return api.ImportCodeSystem500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystem"}, nil
// 		}
// 	}

// 	// if there already is a concept, exit with error
// 	if concept.ID != 0 {
// 		return api.ImportCodeSystem400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("CodeSystem already contains concepts")}, nil
// 	}

// 	// Read the entire CSV content into a buffer
// 	file := request.Body
// 	buf := new(bytes.Buffer)
// 	if _, err := buf.ReadFrom(file); err != nil {
// 		return api.ImportCodeSystem500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while reading the CSV file"}, nil
// 	}

// 	// Create a new reader for header validation
// 	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
// 	requiredColumns, err := validateCSVHeader(reader)
// 	if err != nil {
// 		return api.ImportCodeSystem400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
// 	}

// 	// Create another reader for row processing
// 	reader = csv.NewReader(bytes.NewReader(buf.Bytes()))
// 	go processCSVRows(CSVProcessingJob{CodeSystemID: codeSystemId, Reader: reader, Database: s.Database}, requiredColumns)

// 	return api.ImportCodeSystem202JSONResponse("CSV file is being processed"), nil
// }
