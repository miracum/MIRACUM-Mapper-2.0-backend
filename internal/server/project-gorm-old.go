package server

import (
	"context"
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"

	"github.com/google/uuid"
)

type project_permission_join_user struct {
	Pp_role   string
	User_id   uuid.UUID
	User_name string
}

type code_system_role_join_code_system struct {
	Csr_id             uint32
	Csr_name           string
	Csr_position       uint32
	Csr_code_system_id uint32
	System_id          uint32
	System_name        string
	System_Version     string
	Csr_type           string
}

// GetProject implements api.StrictServerInterface.
func (s *StrictGormServer) GetProjectOld(ctx context.Context, request api.GetProjectRequestObject) (api.GetProjectResponseObject, error) {
	var project models.Project
	var code_system_roles []models.CodeSystemRole
	var project_permissions []models.ProjectPermission
	var csr_join_cs []code_system_role_join_code_system
	var pp_join_user []project_permission_join_user

	s.Database.First(&project, request.ProjectId)
	s.Database.Model(&code_system_roles).Select("code_system_roles.id, code_system_roles.name, code_system_roles.position, code_system_roles.code_system_id, code_systems.id, code_systems.name, code_systems.version, code_system_roles.type").Joins("JOIN code_systems ON code_system_roles.code_system_id = code_systems.id").Where("code_system_roles.project_id = ?", request.ProjectId).Scan(&csr_join_cs)

	s.Database.Model(&project_permissions).Select("project_permissions.role, users.id, users.user_name").Joins("JOIN users ON project_permissions.user_id = users.id").Where("project_permissions.project_id = ?", request.ProjectId).Scan(&pp_join_user)

	id := int32(project.ID)
	modified := project.UpdatedAt.String()
	var pd api.ProjectDetails = api.ProjectDetails{
		Description:         project.Description,
		EquivalenceRequired: project.EquivalenceRequired,
		Id:                  &id,
		Modified:            &modified,
		Name:                project.Name,
		StatusRequired:      project.StatusRequired,
		Version:             project.Version,
	}

	var csrs []api.CodeSystemRole
	for _, v := range csr_join_cs {
		id := int32(v.Csr_id)
		system_id := int32(v.System_id)
		csr := api.CodeSystemRole{
			Id:       &id,
			Name:     v.Csr_name,
			Position: int32(v.Csr_position),
			System: struct {
				Id      *int32  `json:"id,omitempty"`
				Name    *string `json:"name,omitempty"`
				Version *string `json:"version,omitempty"`
			}{Id: &system_id,
				Name:    &v.System_name,
				Version: &v.System_Version},
			Type: api.CodeSystemRoleType(v.Csr_type),
		}
		csrs = append(csrs, csr)
	}

	var pps []api.ProjectPermission
	for _, v := range pp_join_user {
		pp := api.ProjectPermission{
			Role:     api.ProjectPermissionRole(v.Pp_role),
			UserId:   v.User_id.String(),
			UserName: &v.User_name,
		}
		pps = append(pps, pp)
	}

	pd.CodeSystemRoles = csrs
	pd.ProjectPermissions = &pps

	return api.GetProject200JSONResponse(pd), nil
}
