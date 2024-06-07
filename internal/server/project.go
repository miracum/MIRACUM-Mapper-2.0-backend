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
	projectSortColumns = map[api.GetProjectsParamsSortBy]string{
		api.GetProjectsParamsSortByDateCreated: "created",
		api.GetProjectsParamsSortById:          "id",
		api.GetProjectsParamsSortByName:        "name",
	}

	// Define mappings from API parameters to sort orders
	projectSortOrders = map[api.GetProjectsParamsSortOrder]string{
		api.GetProjectsParamsSortOrderAsc:  "ASC",
		api.GetProjectsParamsSortOrderDesc: "DESC",
	}
)

func (s *Server) GetProject(ctx context.Context, request api.GetProjectRequestObject) (api.GetProjectResponseObject, error) {
	var project models.Project

	if err := s.Database.GetProjectQuery(&project, request.ProjectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetProject404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", request.ProjectId)), nil
		default:
			return api.GetProject500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the projects"}, err
		}
	}

	projectDetails := transform.GormProjectToApiProjectDetails(project)

	return api.GetProject200JSONResponse(projectDetails), nil
}

// AddProject implements api.StrictServerInterface.
func (s *Server) AddProject(ctx context.Context, request api.AddProjectRequestObject) (api.AddProjectResponseObject, error) {
	projectDetails := request.Body

	// if len(projectDetails.CodeSystemRoles) == 0 {
	// 	return api.AddProject422JSONResponse("CodeSystemRoles are required"), nil
	// }

	if projectDetails.Id != nil {
		return api.AddProject400JSONResponse{BadRequestErrorJSONResponse: "ID must not be provided"}, nil
	}

	project, err := transform.ApiProjectDetailsToGormProject(*projectDetails)
	if err != nil {
		return api.AddProject400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		// switch {
		// case errors.Is(err, transform.ErrInvalidUUID):
		// 	return api.AddProject400JSONResponse{BadRequestErrorJSONResponse: "Invalid uuid provided"}, nil
		// default:
		// 	return api.AddProject500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to create the project"}, nil
		// }
	}

	// Create the project in the database
	if err := s.Database.CreateProjectQuery(project); err != nil {
		switch {
		case errors.Is(err, database.ErrClientError):
			return api.AddProject400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.AddProject500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to create the project"}, nil
		}
	}
	// Create the project in the database
	id := int32(project.Model.ID)
	projectDetails.Id = &id
	return api.AddProject200JSONResponse(*projectDetails), nil
}

// GetProjects implements api.StrictServerInterface.
func (s *Server) GetProjects(ctx context.Context, request api.GetProjectsRequestObject) (api.GetProjectsResponseObject, error) {

	pageSize := *request.Params.PageSize
	offset := utilities.GetOffset(*request.Params.Page, pageSize)
	sortBy := projectSortColumns[*request.Params.SortBy]
	sortOrder := projectSortOrders[*request.Params.SortOrder]

	var projects []models.Project = []models.Project{}

	if err := s.Database.GetProjectsQuery(&projects, pageSize, offset, sortBy, sortOrder); err != nil {
		return api.GetProjects500JSONResponse{}, err
	}

	var apiProjects []api.Project = []api.Project{}
	for _, project := range projects {
		apiProjects = append(apiProjects, transform.GormProjectToApiProject(project))
	}

	return api.GetProjects200JSONResponse(apiProjects), nil
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

	api_project := transform.GormProjectToApiProject(project)
	return api.DeleteProject200JSONResponse(api_project), nil
}

// EditProject implements api.StrictServerInterface.
func (s *Server) EditProject(ctx context.Context, request api.EditProjectRequestObject) (api.EditProjectResponseObject, error) {
	project := request.Body
	projectId := request.ProjectId

	if project.Id == nil {
		project.Id = &projectId
	} else {
		if *project.Id != projectId {
			return api.EditProject400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(fmt.Sprintf("Project ID %d in URL does not match project ID %d in body", projectId, *project.Id))}, nil
		}
	}

	checkFunc := func(oldProject, newProject *models.Project) error {
		if oldProject.StatusRequired != newProject.StatusRequired || oldProject.EquivalenceRequired != newProject.EquivalenceRequired {
			return database.NewDBError(database.ClientError, "StatusRequired and EquivalenceRequired cannot be changed")
		}
		return nil
	}

	db_project := transform.ApiProjectToGormProject(*project)
	if err := s.Database.UpdateProjectQuery(&db_project, checkFunc); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.EditProject404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		// TODO
		case errors.Is(err, database.ErrClientError):
			return api.EditProject400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		// case errors.Is(err, database.???) error for trying to update status-/equivalenceRequired
		default:
			return api.EditProject500JSONResponse{}, err
		}
	}

	return api.EditProject200JSONResponse(*project), nil

}
