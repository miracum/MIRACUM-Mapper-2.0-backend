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
	offset := GetOffset(*request.Params.Page, pageSize)
	sortBy := projectSortColumns[*request.Params.SortBy]
	sortOrder := projectSortOrders[*request.Params.SortOrder]

	var projects []models.Project = []models.Project{}

	// get user id and roles which are needed to see a project
	userId, err := getUserToCheckPermission(ctx)
	if err != nil {
		return api.GetAllProjects500JSONResponse{InternalServerErrorJSONResponse: "Can't determine the userId from the request"}, err
	}
	roles := ProjectViewPermission

	if err := s.Database.GetAllProjectsQuery(&projects, userId, roles, pageSize, offset, sortBy, sortOrder); err != nil {
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

	// var uuid *uuid.UUID = nil
	// var roles *[]models.ProjectPermissionRole = nil

	// // check if user is admin. If not, get the user id from the context (Setting uuid to nil will return all projects)
	// if !IsAdminFromContext(ctx) {
	// 	userId, err := GetUserIdFromContext(ctx)
	// 	if err != nil {
	// 		return api.GetProject500JSONResponse{InternalServerErrorJSONResponse: "Can't determine the userId from the request"}, err
	// 	}
	// 	uuid = &userId
	// 	roles = &[]models.ProjectPermissionRole{models.ReviewerRole, models.ProjectOwnerRole, models.EditorRole}
	// }

	permissions, err := getUserPermissions(ctx, s, projectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.GetProject404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.GetProject500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(ProjectViewPermission, permissions) {
		return api.GetProject403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to view project with ID %d", projectId))}, nil
	}

	if err := s.Database.GetProjectQuery(&project, projectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetProject404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.GetProject500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project"}, nil
		}
	}

	// Check if the user has the required permissions
	// if uuid != nil {
	// 	var permission models.ProjectPermission
	// 	if err := s.Database.GetProjectPermissionQuery(&permission, projectId, *uuid); err != nil {
	// 		switch {
	// 		case errors.Is(err, database.ErrNotFound):
	// 			return api.GetProject403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to view project with ID %d", projectId))}, nil
	// 		default:
	// 			return api.GetProject500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
	// 		}
	// 	}
	// 	// check if role is in roles
	// 	if slices.Contains(*roles, permission.Role) {
	// 		return api.GetProject403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to view project with ID %d", projectId))}, nil
	// 	}
	// }

	projectDetails := transform.GormProjectToApiProjectDetails(&project)

	return api.GetProject200JSONResponse(projectDetails), nil
}

// UpdateProject implements api.StrictServerInterface.
func (s *Server) UpdateProject(ctx context.Context, request api.UpdateProjectRequestObject) (api.UpdateProjectResponseObject, error) {
	project := request.Body

	permissions, err := getUserPermissions(ctx, s, project.Id)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.UpdateProject404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", project.Id)), nil
		default:
			return api.UpdateProject500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(ProjectUpdatePermission, permissions) {
		return api.UpdateProject403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to edit project with ID %d", project.Id))}, nil
	}

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
