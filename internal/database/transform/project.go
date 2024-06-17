package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
	"miracummapper/internal/utilities"
)

func GormProjectToApiProjectDetails(project *models.Project) api.ProjectDetails {
	var modified string
	if !project.UpdatedAt.IsZero() {
		modified = project.UpdatedAt.String()
	} else {
		modified = ""
	}
	var projectDetails api.ProjectDetails = api.ProjectDetails{
		Description:         project.Description,
		EquivalenceRequired: project.EquivalenceRequired,
		Id:                  int32(project.ID),
		Modified:            modified,
		Name:                project.Name,
		StatusRequired:      project.StatusRequired,
		Version:             project.Version,
	}

	projectDetails.CodeSystemRoles = *GormCodeSystemRolesToApiCodeSystemRoles(&project.CodeSystemRoles)

	projectDetails.ProjectPermissions = GormProjectPermissionsToApiProjectPermissions(&project.Permissions)

	return projectDetails
}

func ApiCreateProjectDetailsToGormProject(projectDetails *api.CreateProjectDetails) (*models.Project, error) {
	project := models.Project{
		Name:                projectDetails.Name,
		Description:         projectDetails.Description,
		Version:             projectDetails.Version,
		EquivalenceRequired: projectDetails.EquivalenceRequired,
		StatusRequired:      projectDetails.StatusRequired,
	}

	project.CodeSystemRoles = *ApiCreateCodeSystemRolesToGormCodeSystemRoles(&projectDetails.CodeSystemRoles)

	// Append the ProjectPermissions
	for _, permission := range *projectDetails.ProjectPermissions {
		userID, err := utilities.ParseUUID(permission.UserId)
		if err != nil {
			return nil, err
		}
		project.Permissions = append(project.Permissions, models.ProjectPermission{
			Role:   models.ProjectPermissionRole(permission.Role),
			UserID: userID,
		})
	}
	return &project, nil
}

func GormProjectToApiProject(project *models.Project) *api.Project {
	var modified string
	if !project.UpdatedAt.IsZero() {
		modified = project.UpdatedAt.String()
	} else {
		modified = ""
	}
	return &api.Project{
		Description:         project.Description,
		EquivalenceRequired: project.EquivalenceRequired,
		Id:                  int32(project.ID),
		Modified:            modified,
		Name:                project.Name,
		StatusRequired:      project.StatusRequired,
		Version:             project.Version,
	}
}

func ApiUpdateProjectToGormProject(project *api.UpdateProject) *models.Project {
	return &models.Project{
		Model: models.Model{
			ID: uint32(project.Id),
		},
		Name:                project.Name,
		Description:         project.Description,
		Version:             project.Version,
		EquivalenceRequired: project.EquivalenceRequired,
		StatusRequired:      project.StatusRequired,
	}
}
