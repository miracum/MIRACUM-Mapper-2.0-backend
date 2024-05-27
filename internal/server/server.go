package server

import (
	"context"
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/config"
	middlewares "miracummapper/internal/server/middlewares"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	strictgin "github.com/oapi-codegen/runtime/strictmiddleware/gin"
	"gorm.io/gorm"
)

type Server struct {
	Database *gorm.DB
	Config   *config.Config
}

func CreateServer(database *gorm.DB, config *config.Config) *Server {
	return &Server{Database: database, Config: config}
}

func CreateStrictMiddleware(v middlewares.JWSValidator) ([]api.StrictMiddlewareFunc, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: middlewares.NewAuthenticator(v),
			},
		})

	strictMiddlewareFuncs := make([]strictgin.StrictGinMiddlewareFunc, 0)
	for _, h := range []gin.HandlerFunc{validator} {
		strictMiddlewareFuncs = append(strictMiddlewareFuncs, func(f strictgin.StrictGinHandlerFunc, operationID string) strictgin.StrictGinHandlerFunc {
			return func(c *gin.Context, request interface{}) (response interface{}, err error) {
				h(c)
				// if c.IsAborted() {
				// 	return nil, nil
				// }
				return f(c, request)
			}
		})
	}

	return strictMiddlewareFuncs, nil
}

func CreateMiddleware(v middlewares.JWSValidator) ([]gin.HandlerFunc, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: middlewares.NewAuthenticator(v),
			},
		})

	return []gin.HandlerFunc{validator}, nil
}

// AddCodeSystem implements api.StrictServerInterface.
func (s *Server) AddCodeSystem(ctx context.Context, request api.AddCodeSystemRequestObject) (api.AddCodeSystemResponseObject, error) {
	panic("unimplemented")
}

// AddMapping implements api.StrictServerInterface.
func (s *Server) AddMapping(ctx context.Context, request api.AddMappingRequestObject) (api.AddMappingResponseObject, error) {
	panic("unimplemented")
}

// AddPermission implements api.StrictServerInterface.
func (s *Server) AddPermission(ctx context.Context, request api.AddPermissionRequestObject) (api.AddPermissionResponseObject, error) {
	panic("unimplemented")
}

// DeleteCodeSystem implements api.StrictServerInterface.
func (s *Server) DeleteCodeSystem(ctx context.Context, request api.DeleteCodeSystemRequestObject) (api.DeleteCodeSystemResponseObject, error) {
	panic("unimplemented")
}

// DeleteMapping implements api.StrictServerInterface.
func (s *Server) DeleteMapping(ctx context.Context, request api.DeleteMappingRequestObject) (api.DeleteMappingResponseObject, error) {
	panic("unimplemented")
}

// DeletePermission implements api.StrictServerInterface.
func (s *Server) DeletePermission(ctx context.Context, request api.DeletePermissionRequestObject) (api.DeletePermissionResponseObject, error) {
	panic("unimplemented")
}

// EditProject implements api.StrictServerInterface.
func (s *Server) EditProject(ctx context.Context, request api.EditProjectRequestObject) (api.EditProjectResponseObject, error) {
	panic("unimplemented")
}

// FindConceptByCode implements api.StrictServerInterface.
func (s *Server) FindConceptByCode(ctx context.Context, request api.FindConceptByCodeRequestObject) (api.FindConceptByCodeResponseObject, error) {
	panic("unimplemented")
}

// GetAllCodeSystemRoles implements api.StrictServerInterface.
func (s *Server) GetAllCodeSystemRoles(ctx context.Context, request api.GetAllCodeSystemRolesRequestObject) (api.GetAllCodeSystemRolesResponseObject, error) {
	panic("unimplemented")
}

// GetAllCodeSystems implements api.StrictServerInterface.
func (s *Server) GetAllCodeSystems(ctx context.Context, request api.GetAllCodeSystemsRequestObject) (api.GetAllCodeSystemsResponseObject, error) {
	panic("unimplemented")
}

// GetAllConcepts implements api.StrictServerInterface.
func (s *Server) GetAllConcepts(ctx context.Context, request api.GetAllConceptsRequestObject) (api.GetAllConceptsResponseObject, error) {
	panic("unimplemented")
}

// GetAllMappings implements api.StrictServerInterface.
func (s *Server) GetAllMappings(ctx context.Context, request api.GetAllMappingsRequestObject) (api.GetAllMappingsResponseObject, error) {
	panic("unimplemented")
}

// GetAllPermissions implements api.StrictServerInterface.
func (s *Server) GetAllPermissions(ctx context.Context, request api.GetAllPermissionsRequestObject) (api.GetAllPermissionsResponseObject, error) {
	panic("unimplemented")
}

// GetCodeSystem implements api.StrictServerInterface.
func (s *Server) GetCodeSystem(ctx context.Context, request api.GetCodeSystemRequestObject) (api.GetCodeSystemResponseObject, error) {
	panic("unimplemented")
}

// GetCodeSystemRole implements api.StrictServerInterface.
func (s *Server) GetCodeSystemRole(ctx context.Context, request api.GetCodeSystemRoleRequestObject) (api.GetCodeSystemRoleResponseObject, error) {
	panic("unimplemented")
}

// GetMappingByID implements api.StrictServerInterface.
func (s *Server) GetMappingByID(ctx context.Context, request api.GetMappingByIDRequestObject) (api.GetMappingByIDResponseObject, error) {
	panic("unimplemented")
}

// GetPermission implements api.StrictServerInterface.
func (s *Server) GetPermission(ctx context.Context, request api.GetPermissionRequestObject) (api.GetPermissionResponseObject, error) {
	panic("unimplemented")
}

// // Ping implements api.StrictServerInterface.
// func (s *StrictGormServer) Ping(ctx context.Context, request api.PingRequestObject) (api.PingResponseObject, error) {
// 	panic("unimplemented")
// }

// UpdateCodeSystem implements api.StrictServerInterface.
func (s *Server) UpdateCodeSystem(ctx context.Context, request api.UpdateCodeSystemRequestObject) (api.UpdateCodeSystemResponseObject, error) {
	panic("unimplemented")
}

// UpdateCodeSystemRole implements api.StrictServerInterface.
func (s *Server) UpdateCodeSystemRole(ctx context.Context, request api.UpdateCodeSystemRoleRequestObject) (api.UpdateCodeSystemRoleResponseObject, error) {
	panic("unimplemented")
}

// UpdateMapping implements api.StrictServerInterface.
func (s *Server) UpdateMapping(ctx context.Context, request api.UpdateMappingRequestObject) (api.UpdateMappingResponseObject, error) {
	panic("unimplemented")
}

// UpdatePermission implements api.StrictServerInterface.
func (s *Server) UpdatePermission(ctx context.Context, request api.UpdatePermissionRequestObject) (api.UpdatePermissionResponseObject, error) {
	panic("unimplemented")
}

var _ api.StrictServerInterface = (*Server)(nil)
