package server

import (
	"context"
	"miracummapper/internal/api"
)

func (s *Server) Ping(ctx context.Context, request api.PingRequestObject) (api.PingResponseObject, error) {

	message := "pong"
	return api.Ping200JSONResponse{
		Message: &message,
	}, nil

}
