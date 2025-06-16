package server

import (
	"context"
	"errors"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/utilities"
)

// ImportCodeSystemVersionIcd implements api.StrictServerInterface.
func (s *Server) ImportCodeSystemVersionIcd(ctx context.Context, request api.ImportCodeSystemVersionIcdRequestObject) (api.ImportCodeSystemVersionIcdResponseObject, error) {
	codeSystemId := request.CodesystemId
	codeSystemVersionId := request.CodesystemVersionId

	// Check if the CodeSystem exists and is of type ICD_10_GM
	var codeSystem models.CodeSystem
	if err := s.Database.GetCodeSystemQuery(&codeSystem, codeSystemId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.ImportCodeSystemVersionIcd404JSONResponse(err.Error()), nil
		default:
			return api.ImportCodeSystemVersionIcd500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystem"}, nil
		}
	}

	if codeSystem.Type != models.ICD_10_GM {
		return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("This endpoint can only be used for ICD_10_GM CodeSystems. Type is: " + string(codeSystem.Type))}, nil
	}

	// Check if the CodeSystemVersion exists and is not already imported
	var codeSystemVersion models.CodeSystemVersion
	if err := s.Database.GetCodeSystemVersionQuery(&codeSystemVersion, codeSystemId, codeSystemVersionId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.ImportCodeSystemVersionIcd404JSONResponse(err.Error()), nil
		default:
			return api.ImportCodeSystemVersionIcd500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the CodeSystemVersion"}, nil
		}
	}

	if codeSystemVersion.Imported {
		return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("CodeSystemVersion is already imported")}, nil
	}

	return processFilesIcd(request.Body, codeSystemId, codeSystemVersionId, s.Database)
}

func processFilesIcd(body *api.ImportCodeSystemVersionIcdJSONRequestBody, codeSystemId int32, codeSystemVersionId int32, db database.Datastore) (api.ImportCodeSystemVersionIcdResponseObject, error) {

	// Validate the request body
	if body == nil {
		return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("No body provided")}, nil
	}

	codeSystemFile, ok := (*body)["codeSystem"].(map[string]any)
	if !ok || codeSystemFile == nil {
		return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("CodeSystem file is missing in the request body")}, nil
	}
	if resourceType, ok := codeSystemFile["resourceType"].(string); !ok || resourceType != "CodeSystem" {
		return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Incorrect FHIR resource type in codeSystem, expected 'CodeSystem'")}, nil
	}

	conceptMap, ok := (*body)["conceptMap"].(map[string]any)
	if !ok || conceptMap == nil {
		return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("ConceptMap file is missing in the request body")}, nil
	}
	if resourceType, ok := conceptMap["resourceType"].(string); !ok || resourceType != "ConceptMap" {
		return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Incorrect FHIR resource type in conceptMap, expected 'ConceptMap'")}, nil
	}

	// Retrieve concepts
	var concepts []database.ConceptImport

	if conceptsList, ok := codeSystemFile["concept"].([]any); ok {
		for _, concept := range conceptsList {
			if conceptMap, ok := concept.(map[string]any); ok {
				code, ok := conceptMap["code"].(string)
				if !ok || code == "" {
					return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Code is missing in CodeSystem file")}, nil
				}
				display, ok := conceptMap["display"].(string)
				if !ok || display == "" {
					return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Display is missing in CodeSystem file")}, nil
				}
				newConcept := database.ConceptImport{
					Code:    code,
					Display: display,
					Status:  models.ActiveConcept,
				}
				concepts = append(concepts, newConcept)
			} else {
				return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Invalid concept format in CodeSystem file")}, nil
			}
		}
	} else {
		return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Invalid concept format in CodeSystem file")}, nil
	}

	// Retrieve replaceBy items
	var replaceByItems []models.ConceptReplaceBy

	if groupList, ok := conceptMap["group"].([]any); ok {
		for _, group := range groupList {
			groupMap, ok := group.(map[string]any)
			if !ok {
				return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Invalid group format in ConceptMap file")}, nil
			}

			elements, ok := groupMap["element"].([]any)
			if !ok {
				return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Element is missing in ConceptMap file")}, nil
			}

			for _, element := range elements {
				elementMap, ok := element.(map[string]any)
				if !ok {
					return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Invalid element format in ConceptMap file")}, nil
				}

				code, ok := elementMap["code"].(string)
				if !ok || code == "" {
					continue // Skip this element if code is missing
				}

				targetList, ok := elementMap["target"].([]any)
				if !ok || len(targetList) == 0 {
					continue // Skip this element if no targets are present
				}

				for _, target := range targetList {
					targetMap, ok := target.(map[string]any)
					if !ok {
						return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Invalid target format in ConceptMap file")}, nil
					}

					mapTo, ok := targetMap["code"].(string)
					if !ok || mapTo == "" {
						continue // Skip this target if mapTo code is missing
					}

					if code == mapTo {
						continue // Skip if code and mapTo are the same
					}

					equivalenceStr, ok := targetMap["equivalence"].(string)
					if !ok || equivalenceStr == "" {
						return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("Equivalence is missing in ConceptMap file")}, nil
					}
					equivalenceValue := models.ConceptReplaceByEquivalence(equivalenceStr)
					equivalence := &equivalenceValue

					replaceByItem := models.ConceptReplaceBy{
						Code:         code,
						MapTo:        mapTo,
						CodeSystemID: codeSystemId,
						Equivalence:  equivalence,
					}
					replaceByItems = append(replaceByItems, replaceByItem)
				}
			}
		}
	}

	// Start the import process
	if !utilities.TryImporting() {
		return api.ImportCodeSystemVersionIcd400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("An import is already in progress. Please wait until it is finished.")}, nil
	}
	go db.ImportConcepts(codeSystemId, codeSystemVersionId, &concepts, &replaceByItems)
	return api.ImportCodeSystemVersionIcd202JSONResponse("Files processed successfully. Starting to create and update concepts in the background."), nil
}
