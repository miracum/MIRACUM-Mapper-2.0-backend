package server

import (
	"context"
	"errors"
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
	"sort"
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
		api.Code:    "code",
		api.Meaning: "display",
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

var (
	// Define mappings from API parameters to database column names
	newConceptsSortColumns = map[api.GetAllNewConceptsParamsSortBy]string{
		api.GetAllNewConceptsParamsSortByCode:    "code",
		api.GetAllNewConceptsParamsSortByMeaning: "display",
	}

	// Define mappings from API parameters to sort orders
	newConceptsSortOrders = map[api.GetAllNewConceptsParamsSortOrder]string{
		api.GetAllNewConceptsParamsSortOrderAsc:  "ASC",
		api.GetAllNewConceptsParamsSortOrderDesc: "DESC",
	}
)

type VersionWithConcepts struct {
	Version  models.CodeSystemVersion
	Concepts map[string]models.Concept
}

type VersionsWithConcepts map[int32]VersionWithConcepts

// GetAllNewConcepts implements api.StrictServerInterface.
func (s *Server) GetAllNewConcepts(ctx context.Context, request api.GetAllNewConceptsRequestObject) (api.GetAllNewConceptsResponseObject, error) {
	sortBy := newConceptsSortColumns[*request.Params.SortBy]
	sortOrder := newConceptsSortOrders[*request.Params.SortOrder]

	var codeSystemId int32 = request.CodesystemId

	var codeSystem models.CodeSystem
	if err := s.Database.GetCodeSystemQuery(&codeSystem, codeSystemId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetAllNewConcepts404JSONResponse(err.Error()), nil
		default:
			return api.GetAllNewConcepts500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystem"}, nil
		}
	}

	sort.Slice(codeSystem.CodeSystemVersions, func(i, j int) bool {
		return codeSystem.CodeSystemVersions[i].ReleaseDate.Before(codeSystem.CodeSystemVersions[j].ReleaseDate)
	})
	orderedVersions := codeSystem.CodeSystemVersions

	concepts, err := s.Database.GetAllConceptsNewByVersionQuery(codeSystemId, orderedVersions, sortBy, sortOrder)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetAllNewConcepts404JSONResponse(err.Error()), nil
		default:
			return api.GetAllNewConcepts500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the Concepts"}, err
		}
	}

	versions := make(map[int32]VersionWithConcepts)
	for count, version := range orderedVersions {
		versionWithConcepts := VersionWithConcepts{
			Version:  version,
			Concepts: make(map[string]models.Concept),
		}
		for _, concept := range (*concepts)[version.ID] {
			isInOlderVersions := false
			for _, olderVersion := range orderedVersions[:count] {
				if _, exists := versions[olderVersion.ID].Concepts[concept.Code]; exists {
					isInOlderVersions = true
					break
				}
			}
			if !isInOlderVersions {
				versionWithConcepts.Concepts[concept.Code] = concept
			}
		}
		versions[version.ID] = versionWithConcepts
	}

	apiVersionsWithConcepts := []api.CodeSystemVersionWithConcepts{}
	for _, versionWithConcepts := range versions {
		apiVersionWithConcepts := transform.GormVersionToApiVersionWithConcepts(&versionWithConcepts.Version)
		if apiVersionWithConcepts == nil {
			continue
		}
		for _, concept := range versionWithConcepts.Concepts {
			apiVersionWithConcepts.Concepts = append(apiVersionWithConcepts.Concepts, *transform.GormConceptToApiConcept(&concept))
		}
		apiVersionsWithConcepts = append(apiVersionsWithConcepts, *apiVersionWithConcepts)
	}

	return api.GetAllNewConcepts200JSONResponse(apiVersionsWithConcepts), nil
}
