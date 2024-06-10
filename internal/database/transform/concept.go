package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
)

func GormConceptToApiConcept(gormConcept *models.Concept) *api.Concept {
	return &api.Concept{
		Id:      int64(gormConcept.ID),
		Meaning: gormConcept.Display,
		Code:    gormConcept.Code,
	}
}
