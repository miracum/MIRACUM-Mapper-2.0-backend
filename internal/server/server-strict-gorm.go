package server

import (
	"context"
	"miracummapper/internal/api"
	"miracummapper/internal/config"

	"gorm.io/gorm"
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
