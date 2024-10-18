package server

import (
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/config"
	"miracummapper/internal/database"
	"miracummapper/internal/database/gormQuery"
	middlewares "miracummapper/internal/server/middlewares"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	"gorm.io/gorm"
)

type Server struct {
	Database      database.Datastore
	Config        *config.Config
	Authenticator *middlewares.Authenticator
}

func CreateServer(database *gorm.DB, config *config.Config, authenticator *middlewares.Authenticator) *Server {
	return &Server{Database: &gormQuery.GormQuery{Database: database}, Config: config, Authenticator: authenticator}
}

// func CreateStrictMiddleware(v middlewares.JWSValidator) ([]api.StrictMiddlewareFunc, error) {
// 	spec, err := api.GetSwagger()
// 	if err != nil {
// 		return nil, fmt.Errorf("loading spec: %w", err)
// 	}

// 	validator := middleware.OapiRequestValidatorWithOptions(spec,
// 		&middleware.Options{
// 			Options: openapi3filter.Options{
// 				AuthenticationFunc: middlewares.NewAuthenticator(v),
// 			},
// 		})

// 	strictMiddlewareFuncs := make([]strictgin.StrictGinMiddlewareFunc, 0)
// 	for _, handler := range []gin.HandlerFunc{validator} {
// 		strictMiddlewareFuncs = append(strictMiddlewareFuncs, func(strictHandler strictgin.StrictGinHandlerFunc, operationID string) strictgin.StrictGinHandlerFunc {
// 			return func(c *gin.Context, request interface{}) (response interface{}, err error) {
// 				handler(c)
// 				// if c.IsAborted() {
// 				// 	return nil, nil
// 				// }
// 				return strictHandler(c, request)
// 			}
// 		})
// 	}

// 	return strictMiddlewareFuncs, nil
// }

func CreateMiddleware(v middlewares.JWSValidator) ([]gin.HandlerFunc, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: middlewares.NewAuthenticate(v),
			},
			ErrorHandler: CustomErrorHandler,
		})

	// TODO: make this configurable through config file
	allowedOrigins := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:8080",
		"http://localhost:80",
		"http://localhost",
		"http://127.0.0.1",
		"http://localhost:443",
		"http://localhost:18512",
	}

	cors := middlewares.SetupCORS(allowedOrigins)

	return []gin.HandlerFunc{cors, validator}, nil
}

func CustomErrorHandler(c *gin.Context, message string, statusCode int) {
	if statusCode == http.StatusUnauthorized || strings.Contains(message, middlewares.ErrorTokenExpiredApi.Reason) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": message})
		return
	}
	c.JSON(statusCode, gin.H{"error": message})
}

var _ api.StrictServerInterface = (*Server)(nil)
