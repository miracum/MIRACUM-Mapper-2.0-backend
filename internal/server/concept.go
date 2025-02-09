package server

import (
	"context"
	"errors"
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
)

var (
	// Define mappings from API parameters to database column names
	conceptSortColumns = map[api.GetAllConceptsParamsSortBy]string{
		api.GetAllConceptsParamsSortByCode:    "code",
		api.GetAllConceptsParamsSortByMeaning: "display",
	}

	// Define mappings from API parameters to sort orders
	conceptSortOrders = map[api.GetAllConceptsParamsSortOrder]string{
		api.GetAllConceptsParamsSortOrderAsc:  "ASC",
		api.GetAllConceptsParamsSortOrderDesc: "DESC",
	}
)

// GetAllConcepts implements api.StrictServerInterface.
func (s *Server) GetAllConcepts(ctx context.Context, request api.GetAllConceptsRequestObject) (api.GetAllConceptsResponseObject, error) {
	pageSize := *request.Params.PageSize
	offset := GetOffset(*request.Params.Page, pageSize)
	sortBy := conceptSortColumns[*request.Params.SortBy]
	sortOrder := conceptSortOrders[*request.Params.SortOrder]

	var meaning, code string
	if request.Params.MeaningSearch != nil {
		meaning = *request.Params.MeaningSearch
	}
	if request.Params.CodeSearch != nil {
		code = *request.Params.CodeSearch
	}

	var codeSystemId int32 = request.CodesystemId
	var concepts []models.Concept = []models.Concept{}

	if err := s.Database.GetAllConceptsQuery(&concepts, codeSystemId, pageSize, offset, sortBy, sortOrder, meaning, code); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetAllConcepts404JSONResponse(fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId)), nil
		default:
			return api.GetAllConcepts500JSONResponse{}, err
		}

	}

	var apiConcepts []api.Concept = []api.Concept{}
	for _, concept := range concepts {
		apiConcepts = append(apiConcepts, *transform.GormConceptToApiConcept(&concept))
	}

	return api.GetAllConcepts200JSONResponse(apiConcepts), nil
}

var (
	// Define mappings from API parameters to database column names
	conceptByVersionSortColumns = map[api.GetAllConceptsByVersionParamsSortBy]string{
		api.GetAllConceptsByVersionParamsSortByCode:    "code",
		api.GetAllConceptsByVersionParamsSortByMeaning: "display",
	}

	// Define mappings from API parameters to sort orders
	conceptByVersionSortOrders = map[api.GetAllConceptsByVersionParamsSortOrder]string{
		api.GetAllConceptsByVersionParamsSortOrderAsc:  "ASC",
		api.GetAllConceptsByVersionParamsSortOrderDesc: "DESC",
	}
)

// GetAllConceptsByVersion implements api.StrictServerInterface.
func (s *Server) GetAllConceptsByVersion(ctx context.Context, request api.GetAllConceptsByVersionRequestObject) (api.GetAllConceptsByVersionResponseObject, error) {
	pageSize := *request.Params.PageSize
	offset := GetOffset(*request.Params.Page, pageSize)
	sortBy := conceptByVersionSortColumns[*request.Params.SortBy]
	sortOrder := conceptByVersionSortOrders[*request.Params.SortOrder]

	var meaning, code string
	if request.Params.MeaningSearch != nil {
		meaning = *request.Params.MeaningSearch
	}
	if request.Params.CodeSearch != nil {
		code = *request.Params.CodeSearch
	}

	var codeSystemId int32 = request.CodesystemId
	var codeSystemVersionId int32 = request.CodesystemVersionId
	var concepts []models.Concept = []models.Concept{}

	if err := s.Database.GetAllConceptsByVersionQuery(&concepts, codeSystemId, codeSystemVersionId, pageSize, offset, sortBy, sortOrder, meaning, code); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetAllConceptsByVersion404JSONResponse(err.Error()), nil
		default:
			return api.GetAllConceptsByVersion500JSONResponse{}, err
		}

	}

	var apiConcepts []api.Concept = []api.Concept{}
	for _, concept := range concepts {
		apiConcepts = append(apiConcepts, *transform.GormConceptToApiConcept(&concept))
	}

	return api.GetAllConceptsByVersion200JSONResponse(apiConcepts), nil
}
