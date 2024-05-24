package server

import (
	"context"
	"miracummapper/internal/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindPets implements all the handlers in the ServerInterface
func (s *StrictGormServer) Ping(ctx context.Context, request api.PingRequestObject) (api.PingResponseObject, error) {

	message := "pong"
	return api.Ping200JSONResponse{
		Message: &message,
	}, nil

}

// FindPets implements all the handlers in the ServerInterface
func (s *StrictServer) Ping(ctx context.Context, request api.PingRequestObject) (api.PingResponseObject, error) {

	message := "pong"
	return api.Ping200JSONResponse{
		Message: &message,
	}, nil

}

// Ping implements codegen.ServerInterface.
func (s *Server) Ping(c *gin.Context) {
	// message := "pong"
	// return codegen.Ping200JSONResponse{
	// 	Message: &message,
	// }, nil

	resp := "pong"

	c.JSON(http.StatusOK, resp)
}
