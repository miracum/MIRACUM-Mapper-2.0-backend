package middlewares

// import (
// 	"fmt"
// 	"miracummapper/internal/api"

// 	"github.com/getkin/kin-openapi/openapi3filter"
// 	"github.com/gin-gonic/gin"
// 	middleware "github.com/oapi-codegen/gin-middleware"
// )

// func CreateMiddleware(v JWSValidator) ([]gin.HandlerFunc, error) {
// 	spec, err := api.GetSwagger()
// 	if err != nil {
// 		return nil, fmt.Errorf("loading spec: %w", err)
// 	}

// 	validator := middleware.OapiRequestValidatorWithOptions(spec,
// 		&middleware.Options{
// 			Options: openapi3filter.Options{
// 				AuthenticationFunc: NewAuthenticator(v),
// 			},
// 		})

// 	return []gin.HandlerFunc{validator}, nil
// }
