package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"

	"github.com/oapi-codegen/runtime/types"
)

func GormConceptToApiConcept(gormConcept *models.Concept) *api.Concept {
	if gormConcept == nil {
		return nil
	}
	if gormConcept.ID == 0 {
		return nil
	}
	return &api.Concept{
		Id:          gormConcept.ID,
		Meaning:     gormConcept.Display,
		Code:        gormConcept.Code,
		Description: gormConcept.Description,
		Status:      api.ConceptStatus(gormConcept.Status),
	}
}

func GormVersionToApiVersionWithConcepts(gormVersion *models.CodeSystemVersion) *api.CodeSystemVersionWithConcepts {
	if gormVersion == nil {
		return nil
	}

	return &api.CodeSystemVersionWithConcepts{
		Id:          gormVersion.ID,
		VersionName: gormVersion.VersionName,
		ReleaseDate: types.Date{Time: gormVersion.ReleaseDate},
		Imported:    gormVersion.Imported,
		ProjectUses: []string{},
		Concepts:    []api.Concept{},
	}
}
