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

func CreateMiddleware(v middlewares.JWSValidator, config *config.Config) ([]gin.HandlerFunc, error) {
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

	allowedOrigins := config.File.CorsConfig.AllowedOrigins
	// warning if cors is set to allow all origins (contains "*")
	for _, origin := range allowedOrigins {
		if origin == "*" {
			fmt.Println("Warning: CORS is set to allow all origins. This is a security risk!")
			break
		}
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
