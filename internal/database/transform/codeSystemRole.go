package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
)

func GormCodeSystemRolesToApiCodeSystemRoles(codeSystemRoles *[]models.CodeSystemRole) *[]api.CodeSystemRole {
	apiCodeSystemRoles := []api.CodeSystemRole{}
	for _, role := range *codeSystemRoles {
		apiCodeSystemRoles = append(apiCodeSystemRoles, *GormCodeSystemRoleToApiCodeSystemRole(&role))
	}
	return &apiCodeSystemRoles
}

func GormCodeSystemRoleToApiCodeSystemRole(codeSystemRole *models.CodeSystemRole) *api.CodeSystemRole {
	return &api.CodeSystemRole{
		Id:   int32(codeSystemRole.ID),
		Name: codeSystemRole.Name,
		System: struct {
			Id          int32   `json:"id"`
			Name        string  `json:"name"`
			NextVersion *string `json:"nextVersion,omitempty"`
			Version     string  `json:"version"`
		}{
			Id:          int32(codeSystemRole.CodeSystemID),
			Name:        codeSystemRole.CodeSystem.Name,
			NextVersion: &codeSystemRole.NextCodeSystemVersion.VersionName, // TODO check if this is correct
			Version:     codeSystemRole.CodeSystemVersion.VersionName,
		},
		Type: api.CodeSystemRoleType(codeSystemRole.Type),
	}
}

func ApiUpdateCodeSystemRoleToGormCodeSystemRole(codeSystemRole *api.UpdateCodeSystemRole, projectId *api.ProjectId) *models.CodeSystemRole {
	return &models.CodeSystemRole{
		ID:        uint32(codeSystemRole.Id),
		Type:      models.CodeSystemRoleType(codeSystemRole.Type),
		Name:      codeSystemRole.Name,
		ProjectID: uint32(*projectId),
	}
}

func ApiCreateCodeSystemRoleToGormCodeSystemRole(codeSystemRole *api.CreateCodeSystemRole) *models.CodeSystemRole {
	return &models.CodeSystemRole{
		Type:                models.CodeSystemRoleType(codeSystemRole.Type),
		Name:                codeSystemRole.Name,
		CodeSystemID:        uint32(codeSystemRole.System),
		CodeSystemVersionID: uint32(codeSystemRole.Version),
	}
}

func ApiCreateCodeSystemRolesToGormCodeSystemRoles(codeSystemRoles *[]api.CreateCodeSystemRole) *[]models.CodeSystemRole {
	gormCodeSystemRoles := []models.CodeSystemRole{}
	for i, role := range *codeSystemRoles {
		gormRole := ApiCreateCodeSystemRoleToGormCodeSystemRole(&role)
		gormRole.Position = uint32(i)
		gormCodeSystemRoles = append(gormCodeSystemRoles, *gormRole)
	}
	return &gormCodeSystemRoles
}
