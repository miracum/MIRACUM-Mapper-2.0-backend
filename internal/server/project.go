package server

import (
	"context"
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
	"miracummapper/internal/utilities"
)

var (
	// Define mappings from API parameters to database column names
	sortColumns = map[api.GetProjectsParamsSortBy]string{
		"dateCreated": "created",
		"id":          "id",
		"name":        "name",
	}

	// Define mappings from API parameters to sort orders
	sortOrders = map[api.GetProjectsParamsSortOrder]string{
		"asc":  "ASC",
		"desc": "DESC",
	}
)

func (s *Server) GetProject(ctx context.Context, request api.GetProjectRequestObject) (api.GetProjectResponseObject, error) {
	var project models.Project

	if err := s.Database.Preload("CodeSystemRoles.CodeSystem").Preload("Permissions.User").First(&project, request.ProjectId).Error; err != nil {
		// Handle error
		if err.Error() == "record not found" {
			return api.GetProject404Response{}, nil
		}
		return api.GetProject500JSONResponse{}, err
	}

	projectDetails := transform.GormProjectToAPIProjectDetails(project)

	return api.GetProject200JSONResponse(projectDetails), nil
}

// AddProject implements api.StrictServerInterface.
func (s *Server) AddProject(ctx context.Context, request api.AddProjectRequestObject) (api.AddProjectResponseObject, error) {
	projectDetails := request.Body

	if len(projectDetails.CodeSystemRoles) == 0 {
		return api.AddProject422JSONResponse("CodeSystemRoles are required"), nil
	}

	project, err := transform.ApiProjectDetailsToGormProject(*projectDetails)
	if err != nil {
		return api.AddProject500JSONResponse{InternalServerErrorJSONResponse: "Invalid uuid provided"}, err
	}

	// Create the project in the database
	if err := s.Database.Create(&project).Error; err != nil {
		return api.AddProject500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to create the project"}, err
	}
	return api.AddProject200JSONResponse(*projectDetails), nil
}

// GetProjects implements api.StrictServerInterface.
func (s *Server) GetProjects(ctx context.Context, request api.GetProjectsRequestObject) (api.GetProjectsResponseObject, error) {

	pageSize := *request.Params.PageSize
	offset := utilities.GetOffset(*request.Params.Page, pageSize)
	sortBy := sortColumns[*request.Params.SortBy]
	sortOrder := sortOrders[*request.Params.SortOrder]

	var projects []models.Project = []models.Project{}

	if err := s.Database.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize).Find(&projects).Error; err != nil {
		return nil, err
	}

	// Convert projects to api.Projects
	var apiProjects []api.Project = []api.Project{}
	for _, project := range projects {
		apiProjects = append(apiProjects, transform.GormProjectToAPIProject(project))
	}

	return api.GetProjects200JSONResponse(apiProjects), nil
}

// DeleteProject implements api.StrictServerInterface.
func (s *Server) DeleteProject(ctx context.Context, request api.DeleteProjectRequestObject) (api.DeleteProjectResponseObject, error) {
	// NINA IST HIER AM WERK; PFOTEN WEG

	project_id := request.ProjectId
	var project models.Project

	if err := s.Database.First(&project, project_id).Error; err != nil {
		if err.Error() == "record not found" {
			return api.DeleteProject404Response{}, nil
		}
		return api.DeleteProject500JSONResponse{}, err
	}

	s.Database.Delete(&models.Project{}, project_id)

	api_project := transform.GormProjectToAPIProject(project)
	return api.DeleteProject200JSONResponse(api_project), nil
}
