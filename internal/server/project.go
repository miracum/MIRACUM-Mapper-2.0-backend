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
)

var (
	// Define mappings from API parameters to database column names
	projectSortColumns = map[api.GetAllProjectsParamsSortBy]string{
		api.GetAllProjectsParamsSortByDateCreated: "created",
		api.GetAllProjectsParamsSortById:          "id",
		api.GetAllProjectsParamsSortByName:        "name",
	}

	// Define mappings from API parameters to sort orders
	projectSortOrders = map[api.GetAllProjectsParamsSortOrder]string{
		api.GetAllProjectsParamsSortOrderAsc:  "ASC",
		api.GetAllProjectsParamsSortOrderDesc: "DESC",
	}
)

// GetAllProjects implements api.StrictServerInterface.
func (s *Server) GetAllProjects(ctx context.Context, request api.GetAllProjectsRequestObject) (api.GetAllProjectsResponseObject, error) {

	pageSize := *request.Params.PageSize
	offset := utilities.GetOffset(*request.Params.Page, pageSize)
	sortBy := projectSortColumns[*request.Params.SortBy]
	sortOrder := projectSortOrders[*request.Params.SortOrder]

	var projects []models.Project = []models.Project{}

	if err := s.Database.GetAllProjectsQuery(&projects, pageSize, offset, sortBy, sortOrder); err != nil {
		return api.GetAllProjects500JSONResponse{}, err
	}

	var apiProjects []api.Project = []api.Project{}
	for _, project := range projects {
		apiProjects = append(apiProjects, *transform.GormProjectToApiProject(&project))
	}

	return api.GetAllProjects200JSONResponse(apiProjects), nil
}

// CreateProject implements api.StrictServerInterface.
func (s *Server) CreateProject(ctx context.Context, request api.CreateProjectRequestObject) (api.CreateProjectResponseObject, error) {
	projectDetails := request.Body

	if len(projectDetails.CodeSystemRoles) == 0 {
		return api.CreateProject422JSONResponse("CodeSystemRoles must not be empty"), nil
	} else if len(projectDetails.ProjectPermissions) == 0 {
		return api.CreateProject422JSONResponse("Permissions must not be empty"), nil
	}

	// adminExists := false
	// for _, permission := range projectDetails.ProjectPermissions {
	// 	if permission.Role == "project_owner" {
	// 		adminExists = true
	// 		break
	// 	}
	// }
	// if !adminExists {
	// 	return api.CreateProject422JSONResponse("At least one admin permission is required"), nil
	// }

	project, err := transform.ApiCreateProjectDetailsToGormProject(projectDetails)
	if err != nil {
		return api.CreateProject400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}

	// Create the project in the database
	if err := s.Database.CreateProjectQuery(project); err != nil {
		switch {
		case errors.Is(err, database.ErrClientError):
			return api.CreateProject400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.CreateProject500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to create the project"}, nil
		}
	}
	// Create the project in the database
	// id := int32(project.Model.ID)
	// projectDetails.Id = &id
	// TODO check that project contains id etc
	return api.CreateProject200JSONResponse(transform.GormProjectToApiProjectDetails(project)), nil
}

func (s *Server) GetProject(ctx context.Context, request api.GetProjectRequestObject) (api.GetProjectResponseObject, error) {
	projectId := request.ProjectId
	var project models.Project

	if err := s.Database.GetProjectQuery(&project, projectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetProject404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.GetProject500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project"}, err
		}
	}

	projectDetails := transform.GormProjectToApiProjectDetails(&project)

	return api.GetProject200JSONResponse(projectDetails), nil
}

// UpdateProject implements api.StrictServerInterface.
func (s *Server) UpdateProject(ctx context.Context, request api.UpdateProjectRequestObject) (api.UpdateProjectResponseObject, error) {
	project := request.Body

	checkFunc := func(oldProject, newProject *models.Project) error {
		if oldProject.StatusRequired != newProject.StatusRequired || oldProject.EquivalenceRequired != newProject.EquivalenceRequired {
			return database.NewDBError(database.ClientError, "StatusRequired and EquivalenceRequired cannot be changed")
		}
		return nil
	}

	db_project := transform.ApiUpdateProjectToGormProject(project)
	if err := s.Database.UpdateProjectQuery(db_project, checkFunc); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.UpdateProject404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.UpdateProject400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.UpdateProject500JSONResponse{}, err
		}
	}

	return api.UpdateProject200JSONResponse(*transform.GormProjectToApiProject(db_project)), nil

}

// DeleteProject implements api.StrictServerInterface.
func (s *Server) DeleteProject(ctx context.Context, request api.DeleteProjectRequestObject) (api.DeleteProjectResponseObject, error) {

	projectId := request.ProjectId
	var project models.Project

	if err := s.Database.DeleteProjectQuery(&project, projectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.DeleteProject404JSONResponse(err.Error()), nil
		default:
			return api.DeleteProject500JSONResponse{InternalServerErrorJSONResponse: api.InternalServerErrorJSONResponse(database.InternalServerErrorMessage)}, nil
		}
	}

	return api.DeleteProject200JSONResponse(*transform.GormProjectToApiProject(&project)), nil
}
