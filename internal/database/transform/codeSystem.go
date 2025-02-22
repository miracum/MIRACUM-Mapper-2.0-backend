package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
)

func GormCodeSystemToApiCodeSystem(codeSystem *models.CodeSystem) *api.CodeSystem {
	return &api.CodeSystem{
		Id:          codeSystem.ID,
		Author:      codeSystem.Author,
		Description: codeSystem.Description,
		Name:        codeSystem.Name,
		Type:        api.CodeSystemType(codeSystem.Type),
		Title:       codeSystem.Title,
		Uri:         codeSystem.Uri,
	}
}

func GormCodeSystemToApiGetCodeSystem(codeSystem *models.CodeSystem) *api.GetCodeSystem {
	return &api.GetCodeSystem{
		Id:          codeSystem.ID,
		Author:      codeSystem.Author,
		Description: codeSystem.Description,
		Uri:         codeSystem.Uri,
		Name:        codeSystem.Name,
		Type:        api.GetCodeSystemType(codeSystem.Type),
		Title:       codeSystem.Title,
		Versions:    *GormCodeSystemVersionsToApiCodeSystemVersions(&codeSystem.CodeSystemVersions),
	}
}

func GormCodeSystemsToApiGetCodeSystems(codeSystems *[]models.CodeSystem) *[]api.GetCodeSystem {
	apiCodeSystems := []api.GetCodeSystem{}
	for _, codeSystem := range *codeSystems {
		apiCodeSystems = append(apiCodeSystems, *GormCodeSystemToApiGetCodeSystem(&codeSystem))
	}
	return &apiCodeSystems
}

func ApiCodeSystemToGormCodeSystem(codeSystem *api.CodeSystem) *models.CodeSystem {
	return &models.CodeSystem{
		Model: models.Model{
			ID: codeSystem.Id,
		},
		Author:      codeSystem.Author,
		Description: codeSystem.Description,
		Uri:         codeSystem.Uri,
		Name:        codeSystem.Name,
		Type:        models.CodeSystemType(codeSystem.Type),
		Title:       codeSystem.Title,
	}
}

func ApiCreateCodeSystemToGormCodeSystem(codeSystem *api.CreateCodeSystem) *models.CodeSystem {
	return &models.CodeSystem{
		Author:      codeSystem.Author,
		Description: codeSystem.Description,
		Uri:         codeSystem.Uri,
		Name:        codeSystem.Name,
		Type:        models.CodeSystemType(codeSystem.Type),
		Title:       codeSystem.Title,
	}
}
