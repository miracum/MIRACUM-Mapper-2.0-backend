package server

import (
	"context"
	"miracummapper/internal/api"
)

// FindConceptByCode implements api.StrictServerInterface.
func (s *Server) FindConceptByCode(ctx context.Context, request api.FindConceptByCodeRequestObject) (api.FindConceptByCodeResponseObject, error) {
	panic("unimplemented")
}

// GetAllConcepts implements api.StrictServerInterface.
func (s *Server) GetAllConcepts(ctx context.Context, request api.GetAllConceptsRequestObject) (api.GetAllConceptsResponseObject, error) {
	panic("unimplemented")
}
