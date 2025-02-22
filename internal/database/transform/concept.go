package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
)

func GormConceptToApiConcept(gormConcept *models.Concept) *api.Concept {
	return &api.Concept{
		Id:          gormConcept.ID,
		Meaning:     gormConcept.Display,
		Code:        gormConcept.Code,
		Description: gormConcept.Description,
		Status:      api.ConceptStatus(gormConcept.Status),
	}
}
