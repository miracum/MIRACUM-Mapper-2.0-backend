package server

import (
	"context"
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/config"
	"miracummapper/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StrictGormServer struct {
	Database *gorm.DB
	Config   *config.Config
}

func CreateStrictGormServer(database *gorm.DB, config *config.Config) *StrictGormServer {
	return &StrictGormServer{Database: database, Config: config}
}

// AddCodeSystem implements api.StrictServerInterface.
func (s *StrictGormServer) AddCodeSystem(ctx context.Context, request api.AddCodeSystemRequestObject) (api.AddCodeSystemResponseObject, error) {
	panic("unimplemented")
}

// AddMapping implements api.StrictServerInterface.
func (s *StrictGormServer) AddMapping(ctx context.Context, request api.AddMappingRequestObject) (api.AddMappingResponseObject, error) {
	panic("unimplemented")
}

// AddPermission implements api.StrictServerInterface.
func (s *StrictGormServer) AddPermission(ctx context.Context, request api.AddPermissionRequestObject) (api.AddPermissionResponseObject, error) {
	panic("unimplemented")
}

// AddProject implements api.StrictServerInterface.
func (s *StrictGormServer) AddProject(ctx context.Context, request api.AddProjectRequestObject) (api.AddProjectResponseObject, error) {
	projectDetails := request.Body

	// Validate the project details, must contain at least one code system role
	if len(projectDetails.CodeSystemRoles) == 0 {
		return api.AddProject422JSONResponse("CodeSystemRoles are required"), nil
	}

	// Create a new project
	project := models.Project{
		Name:                projectDetails.Name,
		Description:         projectDetails.Description,
		Version:             projectDetails.Version,
		EquivalenceRequired: projectDetails.EquivalenceRequired,
		StatusRequired:      projectDetails.StatusRequired,
	}

	// for i, role := range projectDetails.CodeSystemRoles {
	// 	project.CodeSystemRoles = append(project.CodeSystemRoles, models.CodeSystemRole{
	// 		Name:         role.Name,
	// 		Type:         models.CodeSystemRoleType(role.Type),
	// 		Position:     uint32(i),
	// 		CodeSystemID: uint32(*role.System.Id),
	// 	})
	// }

	for _, permission := range *projectDetails.ProjectPermissions {
		userID, err := uuid.Parse(permission.UserId)
		if err != nil {
			return api.AddProject500JSONResponse{InternalServerErrorJSONResponse: "Invalid uuid provided"}, err
		}
		project.Permissions = append(project.Permissions, models.ProjectPermission{
			Role:   models.ProjectPermissionRole(permission.Role),
			UserID: userID,
		})
	}

	// Create the project along with its associations
	// s.Database.Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "project_id"}},
	// 	DoUpdates: clause.AssignmentColumns([]string{"role", "user_id"}),
	// }).Create(&project)

	s.Database.Clauses(clause.OnConflict{DoNothing: true}).Create(&project)

	// Return the ID of the newly created project
	return api.AddProject200JSONResponse(*projectDetails), nil
}

// DeleteCodeSystem implements api.StrictServerInterface.
func (s *StrictGormServer) DeleteCodeSystem(ctx context.Context, request api.DeleteCodeSystemRequestObject) (api.DeleteCodeSystemResponseObject, error) {
	panic("unimplemented")
}

// DeleteMapping implements api.StrictServerInterface.
func (s *StrictGormServer) DeleteMapping(ctx context.Context, request api.DeleteMappingRequestObject) (api.DeleteMappingResponseObject, error) {
	panic("unimplemented")
}

// DeletePermission implements api.StrictServerInterface.
func (s *StrictGormServer) DeletePermission(ctx context.Context, request api.DeletePermissionRequestObject) (api.DeletePermissionResponseObject, error) {
	panic("unimplemented")
}

// DeleteProject implements api.StrictServerInterface.
func (s *StrictGormServer) DeleteProject(ctx context.Context, request api.DeleteProjectRequestObject) (api.DeleteProjectResponseObject, error) {
	panic("unimplemented")
}

// EditProject implements api.StrictServerInterface.
func (s *StrictGormServer) EditProject(ctx context.Context, request api.EditProjectRequestObject) (api.EditProjectResponseObject, error) {
	panic("unimplemented")
}

// FindConceptByCode implements api.StrictServerInterface.
func (s *StrictGormServer) FindConceptByCode(ctx context.Context, request api.FindConceptByCodeRequestObject) (api.FindConceptByCodeResponseObject, error) {
	panic("unimplemented")
}

// GetAllCodeSystemRoles implements api.StrictServerInterface.
func (s *StrictGormServer) GetAllCodeSystemRoles(ctx context.Context, request api.GetAllCodeSystemRolesRequestObject) (api.GetAllCodeSystemRolesResponseObject, error) {
	panic("unimplemented")
}

// GetAllCodeSystems implements api.StrictServerInterface.
func (s *StrictGormServer) GetAllCodeSystems(ctx context.Context, request api.GetAllCodeSystemsRequestObject) (api.GetAllCodeSystemsResponseObject, error) {
	panic("unimplemented")
}

// GetAllConcepts implements api.StrictServerInterface.
func (s *StrictGormServer) GetAllConcepts(ctx context.Context, request api.GetAllConceptsRequestObject) (api.GetAllConceptsResponseObject, error) {
	panic("unimplemented")
}

// GetAllMappings implements api.StrictServerInterface.
func (s *StrictGormServer) GetAllMappings(ctx context.Context, request api.GetAllMappingsRequestObject) (api.GetAllMappingsResponseObject, error) {
	panic("unimplemented")
}

// GetAllPermissions implements api.StrictServerInterface.
func (s *StrictGormServer) GetAllPermissions(ctx context.Context, request api.GetAllPermissionsRequestObject) (api.GetAllPermissionsResponseObject, error) {
	panic("unimplemented")
}

// GetCodeSystem implements api.StrictServerInterface.
func (s *StrictGormServer) GetCodeSystem(ctx context.Context, request api.GetCodeSystemRequestObject) (api.GetCodeSystemResponseObject, error) {
	panic("unimplemented")
}

// GetCodeSystemRole implements api.StrictServerInterface.
func (s *StrictGormServer) GetCodeSystemRole(ctx context.Context, request api.GetCodeSystemRoleRequestObject) (api.GetCodeSystemRoleResponseObject, error) {
	panic("unimplemented")
}

// GetMappingByID implements api.StrictServerInterface.
func (s *StrictGormServer) GetMappingByID(ctx context.Context, request api.GetMappingByIDRequestObject) (api.GetMappingByIDResponseObject, error) {
	panic("unimplemented")
}

// GetPermission implements api.StrictServerInterface.
func (s *StrictGormServer) GetPermission(ctx context.Context, request api.GetPermissionRequestObject) (api.GetPermissionResponseObject, error) {
	panic("unimplemented")
}

// GetProjects implements api.StrictServerInterface.
func (s *StrictGormServer) GetProjects(ctx context.Context, request api.GetProjectsRequestObject) (api.GetProjectsResponseObject, error) {

	pageSize := *request.Params.PageSize
	offset := (*request.Params.Page - 1) * pageSize
	sortBy := *request.Params.SortBy
	switch sortBy {
	case "dateCreated":
		sortBy = "created"
	}
	sortOrder := *request.Params.SortOrder
	switch sortOrder {
	case "asc":
		sortOrder = "ASC"
	case "desc":
		sortOrder = "DESC"
	}

	var projects []models.Project = []models.Project{}
	err := s.Database.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize).Find(&projects).Error
	if err != nil {
		return nil, err
	}

	// Convert projects to api.Projects
	var apiProjects []api.Project = []api.Project{}
	for _, project := range projects {
		apiProjects = append(apiProjects, convertToAPIProject(project))
	}

	return api.GetProjects200JSONResponse(apiProjects), nil
}

func convertToAPIProject(project models.Project) api.Project {
	id := int32(project.ID)
	var modified string
	if !project.UpdatedAt.IsZero() {
		modified = project.UpdatedAt.String()
	} else {
		modified = ""
	}
	return api.Project{
		Description:         project.Description,
		EquivalenceRequired: project.EquivalenceRequired,
		Id:                  &id,
		Modified:            &modified, // Assign the string pointer
		Name:                project.Name,
		StatusRequired:      project.StatusRequired,
		Version:             project.Version,
	}
}

// // Ping implements api.StrictServerInterface.
// func (s *StrictGormServer) Ping(ctx context.Context, request api.PingRequestObject) (api.PingResponseObject, error) {
// 	panic("unimplemented")
// }

// UpdateCodeSystem implements api.StrictServerInterface.
func (s *StrictGormServer) UpdateCodeSystem(ctx context.Context, request api.UpdateCodeSystemRequestObject) (api.UpdateCodeSystemResponseObject, error) {
	panic("unimplemented")
}

// UpdateCodeSystemRole implements api.StrictServerInterface.
func (s *StrictGormServer) UpdateCodeSystemRole(ctx context.Context, request api.UpdateCodeSystemRoleRequestObject) (api.UpdateCodeSystemRoleResponseObject, error) {
	panic("unimplemented")
}

// UpdateMapping implements api.StrictServerInterface.
func (s *StrictGormServer) UpdateMapping(ctx context.Context, request api.UpdateMappingRequestObject) (api.UpdateMappingResponseObject, error) {
	panic("unimplemented")
}

// UpdatePermission implements api.StrictServerInterface.
func (s *StrictGormServer) UpdatePermission(ctx context.Context, request api.UpdatePermissionRequestObject) (api.UpdatePermissionResponseObject, error) {
	panic("unimplemented")
}

var _ api.StrictServerInterface = (*StrictGormServer)(nil)
