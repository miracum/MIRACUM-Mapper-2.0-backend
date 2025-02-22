package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
)

// TODO return by reference. Currently, each time a copy is created --> unnecessary and slow
func GormMappingToApiMapping(mapping models.Mapping) api.Mapping {
	id := mapping.ID
	var modified string
	if !mapping.UpdatedAt.IsZero() {
		modified = mapping.UpdatedAt.String()
	} else {
		modified = ""
	}

	var created string
	if !mapping.CreatedAt.IsZero() {
		created = mapping.CreatedAt.String()
	} else {
		created = ""
	}

	var api_mapping api.Mapping = api.Mapping{
		Comment:     mapping.Comment,
		Created:     created,
		Equivalence: (*api.MappingEquivalence)(mapping.Equivalence),
		Id:          id,
		Modified:    modified,
		Status:      (*api.MappingStatus)(mapping.Status),
	}
	var elements []api.FullElement = []api.FullElement{}
	for _, element := range mapping.Elements {
		elements = append(elements, GormElementToApiFullElement(&element))
	}
	api_mapping.Elements = elements
	return api_mapping
}

func GormElementToApiFullElement(element *models.Element) api.FullElement {
	codeSystemRole := element.CodeSystemRoleID
	var apiFullElement api.FullElement = api.FullElement{
		CodeSystemRole: &codeSystemRole,
		Concept: &api.Concept{
			Id:      *element.ConceptID,
			Code:    element.Concept.Code,
			Meaning: element.Concept.Display,
		},
		NextConcept: &api.Concept{
			Id:      *element.NextConceptID,
			Code:    element.NextConcept.Code,
			Meaning: element.NextConcept.Display,
		},
	}
	return apiFullElement
}

func ApiCreateMappingToGormMapping(mapping api.CreateMapping, projectId int32) models.Mapping {
	var elements []models.Element = []models.Element{}
	for _, element := range *mapping.Elements {
		conceptId := *element.Concept
		codeSystemRoleId := *element.CodeSystemRole
		elements = append(elements, models.Element{
			ConceptID:        &conceptId,
			CodeSystemRoleID: codeSystemRoleId,
		})
	}

	dbMapping := models.Mapping{
		Comment:     mapping.Comment,
		Equivalence: (*models.Equivalence)(mapping.Equivalence),
		Status:      (*models.MappingStatus)(mapping.Status),
		Elements:    elements,
		ProjectID:   projectId,
	}
	return dbMapping
}

func ApiUpdateMappingToGormMapping(mapping api.UpdateMapping, projectId int32) models.Mapping {
	var elements []models.Element = []models.Element{}
	for _, element := range *mapping.Elements {
		conceptId := *element.Concept
		codeSystemRoleId := *element.CodeSystemRole
		elements = append(elements, models.Element{
			ConceptID:        &conceptId,
			CodeSystemRoleID: codeSystemRoleId,
			MappingID:        mapping.Id,
		})
	}

	dbMapping := models.Mapping{
		Model: models.Model{
			ID: mapping.Id,
		},
		Comment:     mapping.Comment,
		Equivalence: (*models.Equivalence)(mapping.Equivalence),
		Status:      (*models.MappingStatus)(mapping.Status),
		Elements:    elements,
		ProjectID:   projectId,
	}

	return dbMapping

}
