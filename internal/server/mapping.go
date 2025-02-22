package server

import (
	"context"
	"errors"
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
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
		api.Asc:  "ASC",
		api.Desc: "DESC",
	}
)

// callback function to check if the mapping is valid
var checkFunc = func(mapping *models.Mapping, project *models.Project) ([]int32, error) {
	var errorMessages []string

	if !project.StatusRequired && mapping.Status != nil {
		errorMessages = append(errorMessages, "The project does not allow to set a status.")
	}
	if !project.EquivalenceRequired && mapping.Equivalence != nil {
		errorMessages = append(errorMessages, "The project does not allow to set an equivalence.")
	}

	codeSystemRoleIDs := make(map[int32]bool)
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
					errorMessages = append(errorMessages, fmt.Sprintf("The CodeSystemRole %d has the CodeSystem %d which does not match the CodeSystem %d of the Concept %d", role.ID, role.CodeSystemID, element.Concept.CodeSystemID, *element.ConceptID))
				}
				break
			}
			// since the roles are ordered in ascending order, break can be performed here because all remaining roles will have an ID greater than the codeSystemRoleID of the element
			if element.CodeSystemRoleID < role.ID {
				break
			}
		}

		if !isValid {
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid mapping: CodeSystemRoleID %d is not valid", element.CodeSystemRoleID))
		}
	}

	if len(errorMessages) > 0 {
		return nil, database.NewDBError(database.ClientError, strings.Join(errorMessages, "; "))
	}

	unusedCodeSystemRoleIDs := make([]int32, 0)
	for _, role := range project.CodeSystemRoles {
		if _, exists := codeSystemRoleIDs[role.ID]; !exists {
			unusedCodeSystemRoleIDs = append(unusedCodeSystemRoleIDs, role.ID)
		}
	}

	return unusedCodeSystemRoleIDs, nil
}

// GetAllMappings implements api.StrictServerInterface.
func (s *Server) GetAllMappings(ctx context.Context, request api.GetAllMappingsRequestObject) (api.GetAllMappingsResponseObject, error) {

	pageSize := *request.Params.PageSize
	offset := GetOffset(*request.Params.Page, pageSize)
	sortBy := mappingSortColumns[*request.Params.SortBy]
	sortOrder := mappingSortOrders[*request.Params.SortOrder]

	projectId := request.ProjectId

	permissions, err := getUserPermissions(ctx, s, request.ProjectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.GetAllMappings404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.GetAllMappings500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(MappingViewPermission, permissions) {
		return api.GetAllMappings403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to view mappings in the project with ID %d", projectId))}, nil
	}

	var mappings []models.Mapping = []models.Mapping{}

	if err := s.Database.GetAllMappingsQuery(&mappings, projectId, pageSize, offset, sortBy, sortOrder); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetAllMappings404JSONResponse(err.Error()), nil
		default:
			return api.GetAllMappings500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the mappings"}, nil
		}
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

	permissions, err := getUserPermissions(ctx, s, request.ProjectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.CreateMapping404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.CreateMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(MappingCreatePermission, permissions) {
		return api.CreateMapping403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to create mappings in the project with ID %d", projectId))}, nil
	}

	mapping := transform.ApiCreateMappingToGormMapping(*createMapping, projectId)

	if err := s.Database.CreateMappingQuery(&mapping, checkFunc); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.CreateMapping404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.CreateMapping422JSONResponse(err.Error()), nil
		default:
			return api.CreateMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to add the mapping"}, nil
		}
	}

	return api.CreateMapping200JSONResponse(transform.GormMappingToApiMapping(mapping)), nil
}

// UpdateMapping implements api.StrictServerInterface. TODO: reviewer can edit comment, status and equivalence
func (s *Server) UpdateMapping(ctx context.Context, request api.UpdateMappingRequestObject) (api.UpdateMappingResponseObject, error) {
	projectId := request.ProjectId
	updateMapping := *request.Body

	permissions, err := getUserPermissions(ctx, s, request.ProjectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.UpdateMapping404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.UpdateMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(MappingUpdatePermission, permissions) {
		return api.UpdateMapping403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to update mappings in the project with ID %d", projectId))}, nil
	}

	dbMapping := transform.ApiUpdateMappingToGormMapping(updateMapping, projectId)

	if err := s.Database.UpdateMappingQuery(&dbMapping, checkFunc, true); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.UpdateMapping404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.UpdateMapping422JSONResponse(err.Error()), nil
		default:
			return api.UpdateMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to update the mapping"}, nil
		}
	}

	return api.UpdateMapping200JSONResponse(transform.GormMappingToApiMapping(dbMapping)), nil
}

// PatchMapping implements api.StrictServerInterface. TODO: reviewer can edit comment, status and equivalence
func (s *Server) PatchMapping(ctx context.Context, request api.PatchMappingRequestObject) (api.PatchMappingResponseObject, error) {
	projectId := request.ProjectId
	updateMapping := *request.Body

	permissions, err := getUserPermissions(ctx, s, request.ProjectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.PatchMapping404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.PatchMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(MappingUpdatePermission, permissions) {
		return api.PatchMapping403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to update mappings in the project with ID %d", projectId))}, nil
	}

	dbMapping := transform.ApiUpdateMappingToGormMapping(updateMapping, projectId)

	if err := s.Database.UpdateMappingQuery(&dbMapping, checkFunc, false); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.PatchMapping404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.PatchMapping422JSONResponse(err.Error()), nil
		default:
			return api.PatchMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to update the mapping"}, nil
		}
	}

	return api.PatchMapping200JSONResponse(transform.GormMappingToApiMapping(dbMapping)), nil
}

// GetMapping implements api.StrictServerInterface.
func (s *Server) GetMapping(ctx context.Context, request api.GetMappingRequestObject) (api.GetMappingResponseObject, error) {

	projectId := request.ProjectId
	mappingId := request.MappingId

	permissions, err := getUserPermissions(ctx, s, request.ProjectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.GetMapping404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.GetMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(MappingViewPermission, permissions) {
		return api.GetMapping403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to get a mapping in the project with ID %d", projectId))}, nil
	}

	mapping := models.Mapping{}

	if err := s.Database.GetMappingQuery(&mapping, projectId, mappingId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetMapping404JSONResponse(err.Error()), nil
		default:
			return api.GetMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the mapping"}, nil
		}
	}

	return api.GetMapping200JSONResponse(transform.GormMappingToApiMapping(mapping)), nil
}

// DeleteMapping implements api.StrictServerInterface.
func (s *Server) DeleteMapping(ctx context.Context, request api.DeleteMappingRequestObject) (api.DeleteMappingResponseObject, error) {

	projectId := request.ProjectId
	mappingId := request.MappingId

	permissions, err := getUserPermissions(ctx, s, request.ProjectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.DeleteMapping404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.DeleteMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(MappingDeletePermission, permissions) {
		return api.DeleteMapping403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to delete a mapping in the project with ID %d", projectId))}, nil
	}

	mapping := models.Mapping{
		Model: models.Model{
			ID: mappingId,
		},
		ProjectID: projectId,
	}

	if err := s.Database.DeleteMappingQuery(&mapping); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.DeleteMapping404JSONResponse(err.Error()), nil
		default:
			return api.DeleteMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to delete the mapping"}, nil
		}
	}

	return api.DeleteMapping200JSONResponse(transform.GormMappingToApiMapping(mapping)), nil
}
