package server

import (
	"context"
	"errors"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
	"miracummapper/internal/utilities"
)

var (
	// Define mappings from API parameters to database column names
	mappingSortColumns = map[api.GetAllMappingsParamsSortBy]string{
		api.GetAllMappingsParamsSortByComment:     "comment",
		api.GetAllMappingsParamsSortByCreated:     "created",
		api.GetAllMappingsParamsSortByEquivalence: "equivalence",
		api.GetAllMappingsParamsSortById:          "id",
		api.GetAllMappingsParamsSortByModified:    "modified",
		api.GetAllMappingsParamsSortByStatus:      "status",
	}

	// Define mappings from API parameters to sort orders
	mappingSortOrders = map[api.GetAllMappingsParamsSortOrder]string{
		api.GetAllMappingsParamsSortOrderAsc:  "ASC",
		api.GetAllMappingsParamsSortOrderDesc: "DESC",
	}
)

// GetAllMappings implements api.StrictServerInterface.
func (s *Server) GetAllMappings(ctx context.Context, request api.GetAllMappingsRequestObject) (api.GetAllMappingsResponseObject, error) {

	pageSize := *request.Params.PageSize
	offset := utilities.GetOffset(*request.Params.Page, pageSize)
	sortBy := mappingSortColumns[*request.Params.SortBy]
	sortOrder := mappingSortOrders[*request.Params.SortOrder]

	var mappings []models.Mapping = []models.Mapping{}

	if err := s.Database.GetAllMappingsQuery(&mappings, pageSize, offset, sortBy, sortOrder); err != nil {
		return api.GetAllMappings500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the mappings"}, err
	}

	var mapping []api.Mapping = []api.Mapping{}
	for _, m := range mappings {
		mapping = append(mapping, transform.GormMappingToApiMapping(m))
	}

	return api.GetAllMappings200JSONResponse(mapping), nil
}

// CreateMapping implements api.StrictServerInterface.
func (s *Server) CreateMapping(ctx context.Context, request api.CreateMappingRequestObject) (api.CreateMappingResponseObject, error) {

	project_id := request.ProjectId
	mappingDetails := request.Body

	if mappingDetails.Id != nil {
		return api.CreateMapping400JSONResponse{BadRequestErrorJSONResponse: "ID must not be provided"}, nil
	}

	mapping, err := transform.ApiMappingToGormMapping(*mappingDetails, project_id)
	if err != nil {
		return api.CreateMapping400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}

	if err := s.Database.CreateMappingQuery(&mapping); err != nil {
		switch {
		case errors.Is(err, database.ErrClientError):
			return api.CreateMapping400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.CreateMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to add the mapping"}, err
		}
	}
	id := int64(mapping.ID) // TODO this could result in a negative number
	mappingDetails.Id = &id
	return api.CreateMapping200JSONResponse(*mappingDetails), nil
}

// DeleteMapping implements api.StrictServerInterface.
func (s *Server) DeleteMapping(ctx context.Context, request api.DeleteMappingRequestObject) (api.DeleteMappingResponseObject, error) {
	panic("unimplemented")
}

// GetMapping implements api.StrictServerInterface.
func (s *Server) GetMapping(ctx context.Context, request api.GetMappingRequestObject) (api.GetMappingResponseObject, error) {
	panic("unimplemented")
}

// UpdateMapping implements api.StrictServerInterface.
func (s *Server) UpdateMapping(ctx context.Context, request api.UpdateMappingRequestObject) (api.UpdateMappingResponseObject, error) {
	panic("unimplemented")
}
