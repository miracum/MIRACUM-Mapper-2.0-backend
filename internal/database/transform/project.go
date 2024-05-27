package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"

	"github.com/google/uuid"
)

func GormProjectToAPIProjectDetails(project models.Project) api.ProjectDetails {
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

	// Map CodeSystemRoles
	for _, role := range project.CodeSystemRoles {
		id := int32(role.ID)
		system_id := int32(role.CodeSystemID)
		projectDetails.CodeSystemRoles = append(projectDetails.CodeSystemRoles, api.CodeSystemRole{
			Id:       &id,
			Name:     role.Name,
			Position: int32(role.Position),
			System: struct {
				Id      *int32  `json:"id,omitempty"`
				Name    *string `json:"name,omitempty"`
				Version *string `json:"version,omitempty"`
			}{
				Id:      &system_id,
				Name:    &role.CodeSystem.Name,
				Version: &role.CodeSystem.Version,
			},
			Type: api.CodeSystemRoleType(role.Type),
		})
	}

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
			CodeSystemID: uint32(*role.System.Id),
		})
	}

	// Append the ProjectPermissions
	for _, permission := range *projectDetails.ProjectPermissions {
		userID, err := uuid.Parse(permission.UserId)
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

func GormProjectToAPIProject(project models.Project) api.Project {
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
