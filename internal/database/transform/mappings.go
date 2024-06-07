package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
)

// TODO return by reference. Currently, each time a copy is created --> unnecessary and slow
func GormMappingToApiMapping(mapping models.Mapping) api.Mapping {
	id := int64(mapping.ID)
	created := mapping.CreatedAt.String()
	modified := mapping.UpdatedAt.String()
	var api_mapping api.Mapping = api.Mapping{
		Comment:     mapping.Comment,
		Created:     &created,
		Equivalence: (*api.MappingEquivalence)(mapping.Equivalence),
		Id:          &id,
		Modified:    &modified,
		Status:      (*string)(mapping.Status),
	}
	var elements []api.Element = []api.Element{}
	for _, element := range mapping.Elements {
		elements = append(elements, GormElementToApiElement(element))
	}
	api_mapping.Elements = &elements
	return api_mapping
}

func GormElementToApiElement(element models.Element) api.Element {
	id := int32(element.CodeSystemRoleID)
	var api_element api.Element = api.Element{
		SystemId: &id,
		Concept: &api.Concept{
			Code:    &element.Concept.Code,
			Meaning: &element.Concept.Display,
		},
	}
	return api_element
}

func ApiMappingToGormMapping(mapping api.Mapping, project_id int32) (models.Mapping, error) {
	var elements []models.Element = []models.Element{}
	for _, element := range *mapping.Elements {
		elements = append(elements, ApiElementToGormElement(element))
	}

	db_mapping := models.Mapping{
		Comment:     mapping.Comment,
		Equivalence: (*models.Equivalence)(mapping.Equivalence),
		Status:      (*models.Status)(mapping.Status),
		Elements:    elements,
		ProjectID:   uint32(project_id),
	}
	return db_mapping, nil
}

func ApiElementToGormElement(element api.Element) models.Element {
	concept_id := uint32(*element.Concept.Id)
	return models.Element{
		CodeSystemRoleID: uint32(*element.SystemId),
		ConceptID:        &concept_id,
	}
}
