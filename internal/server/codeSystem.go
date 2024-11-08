package server

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
)

type CSVProcessingJob struct {
	CodeSystemID int32
	Reader       *csv.Reader
	Database     database.Datastore
}

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

// helper functions for import code system

func validateCSVHeader(reader *csv.Reader) (map[string]int, error) {
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading the CSV header: %v", err)
	}

	requiredColumns := map[string]int{"code": -1, "meaning": -1}
	for i, column := range header {
		if j, exists := requiredColumns[column]; j == -1 && exists {
			requiredColumns[column] = i
		} else if j != -1 {
			return nil, fmt.Errorf("error: column found multiple times in csv file: %s", column)
		}
	}
	for column, present := range requiredColumns {
		if present == -1 {
			return nil, fmt.Errorf("missing required column: %s", column)
		}
	}

	return requiredColumns, nil
}

func processCSVRows(job CSVProcessingJob, requiredColumns map[string]int) {
	codeIndex := requiredColumns["code"]
	meaningIndex := requiredColumns["meaning"]

	var concepts []models.Concept

	for {
		record, err := job.Reader.Read()
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

		concept := models.Concept{
			Code:         record[codeIndex],
			Display:      record[meaningIndex],
			CodeSystemID: uint32(job.CodeSystemID),
		}
		concepts = append(concepts, concept)

		// fmt.Printf("Code: %s, Meaning: %s\n", record[codeIndex], record[meaningIndex])
	}
	if err := job.Database.CreateConceptsQuery(&concepts); err != nil {
		log.Printf("Error inserting concepts into database: %v", err)
		return
	}

	log.Printf("CSV file processed successfully for CodeSystemID: %d", job.CodeSystemID)
}

// ImportCodeSystem implements api.StrictServerInterface
func (s *Server) ImportCodeSystem(ctx context.Context, request api.ImportCodeSystemRequestObject) (api.ImportCodeSystemResponseObject, error) {
	codeSystemId := request.CodesystemId
	var codeSystem models.CodeSystem
	var concept models.Concept

	if err := s.Database.GetFirstElementCodeSystemQuery(&codeSystem, codeSystemId, &concept); err != nil {
		if errors.Is(err, database.ErrNotFound) && codeSystem.ID == 0 {
			return api.ImportCodeSystem404JSONResponse(fmt.Sprintf("CodeSystem with ID %d couldn't be found.", request.CodesystemId)), nil
		} else if errors.Is(err, database.ErrNotFound) && codeSystem.ID != 0 && concept.ID == 0 {
		} else {
			return api.ImportCodeSystem500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystem"}, nil
		}
	}

	// if there already is a concept, exit with error
	if concept.ID != 0 {
		return api.ImportCodeSystem400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("CodeSystem already contains concepts")}, nil
	}

	// Read the entire CSV content into a buffer
	file := request.Body
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return api.ImportCodeSystem500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while reading the CSV file"}, nil
	}

	// Create a new reader for header validation
	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
	requiredColumns, err := validateCSVHeader(reader)
	if err != nil {
		return api.ImportCodeSystem400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}

	// Create another reader for row processing
	reader = csv.NewReader(bytes.NewReader(buf.Bytes()))
	go processCSVRows(CSVProcessingJob{CodeSystemID: codeSystemId, Reader: reader, Database: s.Database}, requiredColumns)

	return api.ImportCodeSystem202JSONResponse("CSV file is being processed"), nil
}
