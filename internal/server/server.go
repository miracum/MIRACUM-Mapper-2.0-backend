package server

import (
	"database/sql"
	"fmt"
	"miracummapper/internal/api"

	"miracummapper/internal/config"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"

	middlewares "miracummapper/internal/server/middlewares"

	middleware "github.com/oapi-codegen/gin-middleware"
)

type Server struct {
	Database *sql.DB
	Config   *config.Config
}

// Server implements the api.ServerInterface
var _ api.ServerInterface = (*Server)(nil)

func CreateServer(database *sql.DB, config *config.Config) *Server {
	return &Server{Database: database, Config: config}
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

// AddCodeSystem implements codegen.ServerInterface.
func (s *Server) AddCodeSystem(c *gin.Context) {
	panic("unimplemented")
}

// AddMapping implements codegen.ServerInterface.
func (s *Server) AddMapping(c *gin.Context, projectId int32) {
	panic("unimplemented")
}

// AddPermission implements codegen.ServerInterface.
func (s *Server) AddPermission(c *gin.Context, projectId int32, userId string) {
	panic("unimplemented")
}

// AddProject implements codegen.ServerInterface.
func (s *Server) AddProject(c *gin.Context) {
	panic("unimplemented")
}

// DeleteCodeSystem implements codegen.ServerInterface.
func (s *Server) DeleteCodeSystem(c *gin.Context, codeSystemId int32) {
	panic("unimplemented")
}

// DeleteMapping implements codegen.ServerInterface.
func (s *Server) DeleteMapping(c *gin.Context, projectId int32, mappingId int64) {
	panic("unimplemented")
}

// DeletePermission implements codegen.ServerInterface.
func (s *Server) DeletePermission(c *gin.Context, projectId int32, userId string) {
	panic("unimplemented")
}

// DeleteProject implements codegen.ServerInterface.
func (s *Server) DeleteProject(c *gin.Context, projectId int32) {
	panic("unimplemented")
}

// EditProject implements codegen.ServerInterface.
func (s *Server) EditProject(c *gin.Context, projectId int32) {
	panic("unimplemented")
}

// FindConceptByCode implements codegen.ServerInterface.
func (s *Server) FindConceptByCode(c *gin.Context, codeSystemId int32, params api.FindConceptByCodeParams) {
	panic("unimplemented")
}

// GetAllCodeSystemRoles implements codegen.ServerInterface.
func (s *Server) GetAllCodeSystemRoles(c *gin.Context, projectId int32) {
	panic("unimplemented")
}

// GetAllCodeSystems implements codegen.ServerInterface.
func (s *Server) GetAllCodeSystems(c *gin.Context) {
	panic("unimplemented")
}

// GetAllConcepts implements codegen.ServerInterface.
func (s *Server) GetAllConcepts(c *gin.Context, codeSystemId int32, params api.GetAllConceptsParams) {
	panic("unimplemented")
}

// GetAllMappings implements codegen.ServerInterface.
func (s *Server) GetAllMappings(c *gin.Context, projectId int32, params api.GetAllMappingsParams) {
	panic("unimplemented")
}

// GetAllPermissions implements codegen.ServerInterface.
func (s *Server) GetAllPermissions(c *gin.Context, projectId int32) {
	panic("unimplemented")
}

// GetCodeSystem implements codegen.ServerInterface.
func (s *Server) GetCodeSystem(c *gin.Context, codeSystemId int32) {
	panic("unimplemented")
}

// GetCodeSystemRole implements codegen.ServerInterface.
func (s *Server) GetCodeSystemRole(c *gin.Context, projectId int32, codeSystemRoleId int32) {
	panic("unimplemented")
}

// GetMappingByID implements codegen.ServerInterface.
func (s *Server) GetMappingByID(c *gin.Context, projectId int32, mappingId int64) {
	panic("unimplemented")
}

// GetPermission implements codegen.ServerInterface.
func (s *Server) GetPermission(c *gin.Context, projectId int32, userId string) {
	panic("unimplemented")
}

// GetProject implements codegen.ServerInterface.
func (s *Server) GetProject(c *gin.Context, projectId int32) {
	panic("unimplemented")
}

// UpdateCodeSystem implements codegen.ServerInterface.
func (s *Server) UpdateCodeSystem(c *gin.Context, codeSystemId int32) {
	panic("unimplemented")
}

// UpdateCodeSystemRole implements codegen.ServerInterface.
func (s *Server) UpdateCodeSystemRole(c *gin.Context, projectId int32, codeSystemRoleId int32) {
	panic("unimplemented")
}

// UpdateMapping implements codegen.ServerInterface.
func (s *Server) UpdateMapping(c *gin.Context, projectId int32, mappingId int64) {
	panic("unimplemented")
}

// UpdatePermission implements codegen.ServerInterface.
func (s *Server) UpdatePermission(c *gin.Context, projectId int32, userId string) {
	panic("unimplemented")
}
