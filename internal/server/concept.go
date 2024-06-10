package server

import (
	"context"
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
	"miracummapper/internal/utilities"
)

var (
	// Define mappings from API parameters to database column names
	conceptSortColumns = map[api.GetAllConceptsParamsSortBy]string{
		api.Code:    "code",
		api.Meaning: "display",
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
	offset := utilities.GetOffset(*request.Params.Page, pageSize)
	sortBy := conceptSortColumns[*request.Params.SortBy]
	sortOrder := conceptSortOrders[*request.Params.SortOrder]

	var meaning, code string
	if request.Params.MeaningSearch != nil {
		meaning = *request.Params.MeaningSearch
	}
	if request.Params.CodeSearch != nil {
		code = *request.Params.CodeSearch
	}

	var concepts []models.Concept = []models.Concept{}

	if err := s.Database.GetAllConceptsQuery(&concepts, pageSize, offset, sortBy, sortOrder, meaning, code); err != nil {
		return api.GetAllConcepts500JSONResponse{}, err
	}

	var apiConcepts []api.Concept = []api.Concept{}
	for _, concept := range concepts {
		apiConcepts = append(apiConcepts, *transform.GormConceptToApiConcept(&concept))
	}

	return api.GetAllConcepts200JSONResponse(apiConcepts), nil
}
