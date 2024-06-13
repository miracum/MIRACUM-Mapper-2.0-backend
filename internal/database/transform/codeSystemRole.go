package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
)

func GormCodeSystemRolesToApiCodeSystemRoles(codeSystemRoles []models.CodeSystemRole) []api.CodeSystemRole {
	apiCodeSystemRoles := []api.CodeSystemRole{}
	for _, role := range codeSystemRoles {
		apiCodeSystemRoles = append(apiCodeSystemRoles, GormCodeSystemRoleToApiCodeSystemRole(role))
	}
	return apiCodeSystemRoles
}

func GormCodeSystemRoleToApiCodeSystemRole(codeSystemRole models.CodeSystemRole) api.CodeSystemRole {
	id := int32(codeSystemRole.ID)
	return api.CodeSystemRole{
		Id:   &id,
		Name: codeSystemRole.Name,
		System: struct {
			Id      int32   `json:"id"`
			Name    *string `json:"name,omitempty"`
			Version *string `json:"version,omitempty"`
		}{
			Id:      int32(codeSystemRole.CodeSystemID),
			Name:    &codeSystemRole.CodeSystem.Name,
			Version: &codeSystemRole.CodeSystem.Version,
		},
		Type: api.CodeSystemRoleType(codeSystemRole.Type),
	}
}

func ApiCodeSystemRoleToGormCodeSystemRole(codeSystemRole api.CodeSystemRole) models.CodeSystemRole {
	return models.CodeSystemRole{
		ID:           uint32(*codeSystemRole.Id),
		Name:         codeSystemRole.Name,
		CodeSystemID: uint32(codeSystemRole.System.Id),
		Type:         models.CodeSystemRoleType(codeSystemRole.Type),
		ProjectID:    uint32(codeSystemRole.System.Id),
		CodeSystem: models.CodeSystem{
			Model: models.Model{
				ID: uint32(codeSystemRole.System.Id),
			},
			Name:    *codeSystemRole.System.Name,
			Version: *codeSystemRole.System.Version,
		},
	}
}
