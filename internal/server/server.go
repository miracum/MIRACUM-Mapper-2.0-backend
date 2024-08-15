package server

import (
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/config"
	"miracummapper/internal/database"
	"miracummapper/internal/database/gormQuery"
	middlewares "miracummapper/internal/server/middlewares"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	strictgin "github.com/oapi-codegen/runtime/strictmiddleware/gin"
	"gorm.io/gorm"
)

type Server struct {
	Database database.Datastore
	Config   *config.Config
}

func CreateServerWithGormDB(database *gorm.DB, config *config.Config) *Server {
	return &Server{Database: &gormQuery.GormQuery{Database: database}, Config: config}
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

	cors := middlewares.SetupCORS()

	return []gin.HandlerFunc{cors, validator}, nil
}

var _ api.StrictServerInterface = (*Server)(nil)
