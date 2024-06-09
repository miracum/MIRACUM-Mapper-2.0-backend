package server

import (
	"context"
	"errors"
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
	"miracummapper/internal/utilities"
	"strings"
)

var (
	// Define mappings from API parameters to database column names
	mappingSortColumns = map[api.GetAllMappingsParamsSortBy]string{
		api.GetAllMappingsParamsSortByComment:     "comment",
		api.GetAllMappingsParamsSortByCreated:     "created",
		api.GetAllMappingsParamsSortByEquivalence: "equivalence",
		api.GetAllMappingsParamsSortById:          "ID",
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

	projectId := int(request.ProjectId)
	var mappings []models.Mapping = []models.Mapping{}

	if err := s.Database.GetAllMappingsQuery(&mappings, projectId, pageSize, offset, sortBy, sortOrder); err != nil {
		return api.GetAllMappings500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the mappings"}, nil
	}

	var mapping []api.Mapping = []api.Mapping{}
	for _, m := range mappings {
		mapping = append(mapping, transform.GormMappingToApiMapping(m))
	}

	return api.GetAllMappings200JSONResponse(mapping), nil
}

// CreateMapping implements api.StrictServerInterface.
func (s *Server) CreateMapping(ctx context.Context, request api.CreateMappingRequestObject) (api.CreateMappingResponseObject, error) {

	projectId := request.ProjectId
	createMapping := request.Body

	mapping := transform.ApiCreateMappingToGormMapping(*createMapping, projectId)

	checkFunc := func(mapping *models.Mapping, project *models.Project) error {
		var errorMessages []string

		if !project.StatusRequired && mapping.Status != nil {
			errorMessages = append(errorMessages, "The project does not allow to set a status")
		}
		if !project.EquivalenceRequired && mapping.Equivalence != nil {
			errorMessages = append(errorMessages, "The project does not allow to set an equivalence")
		}

		codeSystemRoleIDs := make(map[uint32]bool)
		for _, element := range mapping.Elements {
			isValid := false
			for _, role := range project.CodeSystemRoles {
				if element.CodeSystemRoleID == role.ID {
					if _, exists := codeSystemRoleIDs[role.ID]; exists {
						errorMessages = append(errorMessages, fmt.Sprintf("Duplicate CodeSystemRoleID %d", role.ID))
					} else {
						codeSystemRoleIDs[role.ID] = true
						isValid = true
					}
					if role.CodeSystemID != element.Concept.CodeSystemID {
						errorMessages = append(errorMessages, fmt.Sprintf("The CodeSystemRole %d has the CodeSystem %d which does not match the CodeSystem %d of the Concept %d", role.ID, role.CodeSystemID, element.Concept.CodeSystemID, element.ConceptID))
					}
					break
				}
				// since the roles are ordered in ascending order, we can break here because all remaining roles will have an ID greater than the codeSystemRoleID of the element
				if element.CodeSystemRoleID < role.ID {
					break
				}
			}

			if !isValid || element.CodeSystemRoleID > project.CodeSystemRoles[len(project.CodeSystemRoles)-1].ID {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid mapping: CodeSystemRoleID %d is not valid", element.CodeSystemRoleID))
			}
		}

		if len(errorMessages) > 0 {
			return database.NewDBError(database.ClientError, strings.Join(errorMessages, "; "))
		}

		return nil
	}

	if err := s.Database.CreateMappingQuery(&mapping, checkFunc); err != nil {
		switch {
		case errors.Is(err, database.ErrClientError):
			return api.CreateMapping400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.CreateMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to add the mapping"}, err
		}
	}
	id := uint64(mapping.ID) // TODO this could result in a negative number
	mapping.ModelBigId.ID = id

	return api.CreateMapping200JSONResponse(transform.GormMappingToApiMapping(mapping)), nil
}

// DeleteMapping implements api.StrictServerInterface.
func (s *Server) DeleteMapping(ctx context.Context, request api.DeleteMappingRequestObject) (api.DeleteMappingResponseObject, error) {
	panic("unimplemented")
}

// GetMapping implements api.StrictServerInterface.
func (s *Server) GetMapping(ctx context.Context, request api.GetMappingRequestObject) (api.GetMappingResponseObject, error) {

	projectId := int(request.ProjectId)
	mappingId := request.MappingId

	mapping := models.Mapping{}

	if err := s.Database.GetMappingQuery(&mapping, projectId, mappingId); err != nil {
		// TODO error handling, 404 etc
		return api.GetMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to add the mapping"}, nil
	}

	return api.GetMapping200JSONResponse(transform.GormMappingToApiMapping(mapping)), nil
}

// UpdateMapping implements api.StrictServerInterface.
func (s *Server) UpdateMapping(ctx context.Context, request api.UpdateMappingRequestObject) (api.UpdateMappingResponseObject, error) {
	panic("unimplemented")
}
