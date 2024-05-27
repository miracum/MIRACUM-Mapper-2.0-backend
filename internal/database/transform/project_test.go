package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGormProjectToAPIProjectDetails(t *testing.T) {
	// Given
	gormProject := models.Project{
		Model: models.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Description:         "Test Description",
		EquivalenceRequired: true,
		Name:                "Test Name",
		StatusRequired:      true,
		Version:             "1.0",
		CodeSystemRoles: []models.CodeSystemRole{
			{
				ID:           1,
				Type:         models.Source,
				Name:         "Test Role1",
				Position:     1,
				CodeSystemID: 1,
				CodeSystem: models.CodeSystem{
					Name:    "Test System",
					Version: "1.0",
				},
			},
		},
		Permissions: []models.ProjectPermission{
			{
				Role:   models.ProjectPermissionRole("admin"),
				UserID: uuid.New(),
				User: models.User{
					Id:       uuid.New(),
					UserName: "Test User",
				},
			},
		},
	}

	t.Run("normal case", func(t *testing.T) {
		// When
		apiProjectDetails := GormProjectToAPIProjectDetails(gormProject)

		// Then
		assert.Equal(t, int32(gormProject.ID), *apiProjectDetails.Id)
		assert.Equal(t, gormProject.Description, apiProjectDetails.Description)
		assert.Equal(t, gormProject.EquivalenceRequired, apiProjectDetails.EquivalenceRequired)
		assert.Equal(t, gormProject.Name, apiProjectDetails.Name)
		assert.Equal(t, gormProject.StatusRequired, apiProjectDetails.StatusRequired)
		assert.Equal(t, gormProject.Version, apiProjectDetails.Version)
		assert.Equal(t, gormProject.UpdatedAt.String(), *apiProjectDetails.Modified)
		// assert.Equal(t, gormProject.CreatedAt.String(), *apiProjectDetails.Created)

		// Add assertions for CodeSystemRoles
		assert.Equal(t, len(gormProject.CodeSystemRoles), len(apiProjectDetails.CodeSystemRoles))
		role := gormProject.CodeSystemRoles[0]
		assert.Equal(t, role.Name, apiProjectDetails.CodeSystemRoles[0].Name)
		assert.Equal(t, int32(role.ID), *apiProjectDetails.CodeSystemRoles[0].Id)
		assert.Equal(t, api.CodeSystemRoleType(role.Type), apiProjectDetails.CodeSystemRoles[0].Type)
		assert.Equal(t, role.CodeSystem.Name, *apiProjectDetails.CodeSystemRoles[0].System.Name)
		assert.Equal(t, role.CodeSystem.Version, *apiProjectDetails.CodeSystemRoles[0].System.Version)

		// Add assertions for Permissions
		assert.Equal(t, len(gormProject.Permissions), len(*apiProjectDetails.ProjectPermissions))
		permission := gormProject.Permissions[0]
		assert.Equal(t, api.ProjectPermissionRole(permission.Role), (*apiProjectDetails.ProjectPermissions)[0].Role)
		assert.Equal(t, permission.User.UserName, *(*apiProjectDetails.ProjectPermissions)[0].UserName)
		assert.Equal(t, permission.UserID.String(), (*apiProjectDetails.ProjectPermissions)[0].UserId)
	})

	t.Run("UpdatedAt is zero", func(t *testing.T) {
		// Given
		gormProject.UpdatedAt = time.Time{} // set UpdatedAt to zero value

		// When
		apiProjectDetails := GormProjectToAPIProjectDetails(gormProject)

		// Then
		assert.Equal(t, "", *apiProjectDetails.Modified) // Modified should be an empty string
	})
}

func TestApiProjectDetailsToGormProject(t *testing.T) {
	systemId := int32(1)

	// Given
	projectDetails := api.ProjectDetails{
		Name:                "Test Name",
		Description:         "Test Description",
		Version:             "1.0",
		EquivalenceRequired: true,
		StatusRequired:      true,
		CodeSystemRoles: []api.CodeSystemRole{
			{
				Name: "Test Role",
				Type: api.Source,
				System: struct {
					Id      *int32  `json:"id,omitempty"`
					Name    *string `json:"name,omitempty"`
					Version *string `json:"version,omitempty"`
				}{
					Id: &systemId,
				},
			},
			{
				Name: "Test Role2",
				Type: api.Source,
				System: struct {
					Id      *int32  `json:"id,omitempty"`
					Name    *string `json:"name,omitempty"`
					Version *string `json:"version,omitempty"`
				}{
					Id: &systemId,
				},
			},
		},
		ProjectPermissions: &[]api.ProjectPermission{
			{
				Role:   api.ProjectPermissionRole("admin"),
				UserId: uuid.New().String(),
			},
		},
	}

	t.Run("normal case", func(t *testing.T) {
		// When
		gormProject, err := ApiProjectDetailsToGormProject(projectDetails)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, projectDetails.Name, gormProject.Name)
		assert.Equal(t, projectDetails.Description, gormProject.Description)
		assert.Equal(t, projectDetails.EquivalenceRequired, gormProject.EquivalenceRequired)
		assert.Equal(t, projectDetails.StatusRequired, gormProject.StatusRequired)
		assert.Equal(t, projectDetails.Version, gormProject.Version)

		// Add assertions for CodeSystemRoles
		assert.Equal(t, projectDetails.CodeSystemRoles[0].Name, gormProject.CodeSystemRoles[0].Name)
		assert.Equal(t, models.CodeSystemRoleType(projectDetails.CodeSystemRoles[0].Type), gormProject.CodeSystemRoles[0].Type)
		assert.Equal(t, uint32(systemId), gormProject.CodeSystemRoles[0].CodeSystemID)
		assert.Equal(t, uint32(0), gormProject.CodeSystemRoles[0].Position)
		assert.Equal(t, uint32(1), gormProject.CodeSystemRoles[1].Position)

		// Add assertions for ProjectPermissions
		assert.Equal(t, models.ProjectPermissionRole((*projectDetails.ProjectPermissions)[0].Role), gormProject.Permissions[0].Role)
		assert.Equal(t, uuid.MustParse((*projectDetails.ProjectPermissions)[0].UserId), gormProject.Permissions[0].UserID)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Given
		invalidUUIDProjectDetails := projectDetails
		invalidUUIDProjectDetails.ProjectPermissions = &[]api.ProjectPermission{
			{
				Role:   api.ProjectPermissionRole("admin"),
				UserId: "invalid-uuid",
			},
		}

		// When
		_, err := ApiProjectDetailsToGormProject(invalidUUIDProjectDetails)

		// Then
		assert.Error(t, err)
	})
}

func TestGormProjectToAPIProject(t *testing.T) {
	// Given
	gormProject := models.Project{
		Model: models.Model{
			ID:        1,
			UpdatedAt: time.Now(),
		},
		Description:         "Test Description",
		EquivalenceRequired: true,
		Name:                "Test Name",
		StatusRequired:      true,
		Version:             "1.0",
	}

	t.Run("normal case", func(t *testing.T) {
		// When
		apiProject := GormProjectToAPIProject(gormProject)

		// Then
		assert.Equal(t, int32(gormProject.ID), *apiProject.Id)
		assert.Equal(t, gormProject.Description, apiProject.Description)
		assert.Equal(t, gormProject.EquivalenceRequired, apiProject.EquivalenceRequired)
		assert.Equal(t, gormProject.Name, apiProject.Name)
		assert.Equal(t, gormProject.StatusRequired, apiProject.StatusRequired)
		assert.Equal(t, gormProject.Version, apiProject.Version)
		assert.Equal(t, gormProject.UpdatedAt.String(), *apiProject.Modified)
	})

	t.Run("UpdatedAt is zero", func(t *testing.T) {
		// Given
		gormProject.UpdatedAt = time.Time{} // set UpdatedAt to zero value

		// When
		apiProject := GormProjectToAPIProject(gormProject)

		// Then
		assert.Equal(t, "", *apiProject.Modified) // Modified should be an empty string
	})
}
