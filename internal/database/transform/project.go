package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
	"miracummapper/internal/utilities"
)

func GormProjectToApiProjectDetails(project models.Project) api.ProjectDetails {
	id := int32(project.ID)
	// modified := project.UpdatedAt.String()
	var modified string
	if !project.UpdatedAt.IsZero() {
		modified = project.UpdatedAt.String()
	} else {
		modified = ""
	}
	var projectDetails api.ProjectDetails = api.ProjectDetails{
		Description:         project.Description,
		EquivalenceRequired: project.EquivalenceRequired,
		Id:                  &id,
		Modified:            &modified,
		Name:                project.Name,
		StatusRequired:      project.StatusRequired,
		Version:             project.Version,
	}

	projectDetails.CodeSystemRoles = GormCodeSystemRolesToApiCodeSystemRoles(project.CodeSystemRoles)

	// Map Permissions
	var permissions []api.ProjectPermission
	for _, perm := range project.Permissions {
		permissions = append(permissions, api.ProjectPermission{
			Role:     api.ProjectPermissionRole(perm.Role),
			UserId:   perm.UserID.String(),
			UserName: &perm.User.UserName, // Assuming User is preloaded
		})
	}
	projectDetails.ProjectPermissions = &permissions

	return projectDetails
}

func ApiProjectDetailsToGormProject(projectDetails api.ProjectDetails) (*models.Project, error) {
	project := models.Project{
		Name:                projectDetails.Name,
		Description:         projectDetails.Description,
		Version:             projectDetails.Version,
		EquivalenceRequired: projectDetails.EquivalenceRequired,
		StatusRequired:      projectDetails.StatusRequired,
	}

	// Append the CodeSystemRoles
	for i, role := range projectDetails.CodeSystemRoles {
		project.CodeSystemRoles = append(project.CodeSystemRoles, models.CodeSystemRole{
			Name:         role.Name,
			Type:         models.CodeSystemRoleType(role.Type),
			Position:     uint32(i),
			CodeSystemID: uint32(role.System.Id),
		})
	}

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

func GormProjectToApiProject(project models.Project) api.Project {
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
		Modified:            &modified,
		Name:                project.Name,
		StatusRequired:      project.StatusRequired,
		Version:             project.Version,
	}
}

func ApiProjectToGormProject(project api.Project) models.Project {
	return models.Project{
		Model: models.Model{
			ID: uint32(*project.Id),
		},
		Name:                project.Name,
		Description:         project.Description,
		Version:             project.Version,
		EquivalenceRequired: project.EquivalenceRequired,
		StatusRequired:      project.StatusRequired,
	}
}
