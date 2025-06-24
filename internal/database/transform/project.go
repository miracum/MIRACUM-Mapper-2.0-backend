package transform

import (
	"fmt"
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

	var created string
	if !project.CreatedAt.IsZero() {
		created = project.CreatedAt.String()
	} else {
		created = ""
	}

	var projectDetails api.ProjectDetails = api.ProjectDetails{
		Description:         project.Description,
		EquivalenceRequired: project.EquivalenceRequired,
		Id:                  project.ID,
		Modified:            modified,
		Created:             created,
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
	userPermissionsMap := make(map[string]bool)

	for _, permission := range projectDetails.ProjectPermissions {
		userID, err := utilities.ParseUUID(permission.UserId)
		if err != nil {
			return nil, err
		}

		if _, exists := userPermissionsMap[permission.UserId]; exists {
			return nil, fmt.Errorf("duplicate permission found for user ID: %s", permission.UserId)
		}

		userPermissionsMap[permission.UserId] = true

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

	var created string
	if !project.CreatedAt.IsZero() {
		created = project.CreatedAt.String()
	} else {
		created = ""
	}

	return &api.Project{
		Description:         project.Description,
		EquivalenceRequired: project.EquivalenceRequired,
		Id:                  project.ID,
		Modified:            modified,
		Name:                project.Name,
		StatusRequired:      project.StatusRequired,
		Version:             project.Version,
		Created:             created,
	}
}

func ApiUpdateProjectToGormProject(project *api.UpdateProject) *models.Project {
	return &models.Project{
		Model: models.Model{
			ID: project.Id,
		},
		Name:                project.Name,
		Description:         project.Description,
		Version:             project.Version,
		EquivalenceRequired: project.EquivalenceRequired,
		StatusRequired:      project.StatusRequired,
	}
}

func GormMigrationOptionsToApiMigrationOptions(project *models.Project, newerCodeSystemVersions *map[int32][]models.CodeSystemVersion) *api.MigrationOptions {
	var modified string
	if !project.UpdatedAt.IsZero() {
		modified = project.UpdatedAt.String()
	} else {
		modified = ""
	}

	var created string
	if !project.CreatedAt.IsZero() {
		created = project.CreatedAt.String()
	} else {
		created = ""
	}

	var migrationOptions api.MigrationOptions = api.MigrationOptions{
		Description:         project.Description,
		EquivalenceRequired: project.EquivalenceRequired,
		Id:                  project.ID,
		Modified:            modified,
		Created:             created,
		Name:                project.Name,
		StatusRequired:      project.StatusRequired,
		Version:             project.Version,
	}

	migrationOptions.CodeSystemRoles = *GormCodeSystemRolesMigrationToApiCodeSystemRolesMigration(&project.CodeSystemRoles, newerCodeSystemVersions)

	return &migrationOptions
}

func GormCodeSystemRolesMigrationToApiCodeSystemRolesMigration(codeSystemRoles *[]models.CodeSystemRole, newerCodeSystemVersions *map[int32][]models.CodeSystemVersion) *[]api.CodeSystemRoleMigration {
	apiCodeSystemRoles := []api.CodeSystemRoleMigration{}
	for _, role := range *codeSystemRoles {
		var newerVersions = (*newerCodeSystemVersions)[role.ID]
		apiCodeSystemRole := api.CodeSystemRoleMigration{
			Id:   role.ID,
			Name: role.Name,
			System: struct {
				Id            int32                    `json:"id"`
				Name          string                   `json:"name"`
				NewerVersions *[]api.CodeSystemVersion `json:"newer_versions,omitempty"`
				NextVersion   *api.CodeSystemVersion   `json:"next_version,omitempty"`
				Version       api.CodeSystemVersion    `json:"version"`
			}{
				Id:            role.CodeSystemID,
				Name:          role.CodeSystem.Name,
				NewerVersions: GormCodeSystemVersionsToApiCodeSystemVersions(&newerVersions, false),
				NextVersion:   GormCodeSystemVersionToApiCodeSystemVersion(&role.NextCodeSystemVersion, false),
				Version:       *GormCodeSystemVersionToApiCodeSystemVersion(&role.CodeSystemVersion, false),
			},
			Type: api.CodeSystemRoleMigrationType(role.Type),
		}
		apiCodeSystemRoles = append(apiCodeSystemRoles, apiCodeSystemRole)
	}
	return &apiCodeSystemRoles
}
