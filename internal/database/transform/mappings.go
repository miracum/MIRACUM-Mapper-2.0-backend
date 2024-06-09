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
	codeSystemRole := int32(element.CodeSystemRoleID)
	var apiFullElement api.FullElement = api.FullElement{
		CodeSystemRole: &codeSystemRole,
		Concept: &api.Concept{
			Id:      int64(*element.ConceptID),
			Code:    element.Concept.Code,
			Meaning: element.Concept.Display,
		},
	}
	return apiFullElement
}

func ApiCreateMappingToGormMapping(mapping api.CreateMapping, projectId int32) models.Mapping {
	var elements []models.Element = []models.Element{}
	for _, element := range *mapping.Elements {
		conceptId := uint64(*element.Concept)
		codeSystemRoleId := uint32(*element.CodeSystemRole)
		elements = append(elements, models.Element{
			ConceptID:        &conceptId,
			CodeSystemRoleID: codeSystemRoleId,
		})
	}

	db_mapping := models.Mapping{
		Comment:     mapping.Comment,
		Equivalence: (*models.Equivalence)(mapping.Equivalence),
		Status:      (*models.Status)(mapping.Status),
		Elements:    elements,
		ProjectID:   uint32(projectId),
	}
	return db_mapping
}

// func ApiConceptsToGormElement(concept api.Concept) models.Element {
// 	concept_id := uint32(*concept.Id)
// 	return models.Element{
// 		CodeSystemRoleID: uint32(*element.SystemId),
// 		ConceptID:        &concept_id,
// 	}
// }
