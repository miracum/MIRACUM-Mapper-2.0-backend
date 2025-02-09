package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
)

// TODO löschen
func GormCodeSystemToApiCodeSystem(codeSystem *models.CodeSystem) *api.CodeSystem {
	return &api.CodeSystem{
		Id:          int32(codeSystem.ID),
		Author:      codeSystem.Author,
		Description: codeSystem.Description,
		Name:        codeSystem.Name,
		Title:       codeSystem.Title,
		Uri:         codeSystem.Uri,
		//Versions:    codeSystem.CodeSystemVersions,
		//Version:     "TODO", // TODO codeSystem.Version,
	}
}

func GormCodeSystemToApiGetCodeSystem(codeSystem *models.CodeSystem) *api.GetCodeSystem {
	return &api.GetCodeSystem{
		Id:          int32(codeSystem.ID),
		Author:      codeSystem.Author,
		Description: codeSystem.Description,
		Uri:         codeSystem.Uri,
		Name:        codeSystem.Name,
		Title:       codeSystem.Title,
		Versions:    *GormCodeSystemVersionsToApiCodeSystemVersions(&codeSystem.CodeSystemVersions),
	}
}

func ApiCodeSystemToGormCodeSystem(codeSystem *api.CodeSystem) *models.CodeSystem {
	return &models.CodeSystem{
		Model: models.Model{
			ID: uint32(codeSystem.Id),
		},
		Author:      codeSystem.Author,
		Description: codeSystem.Description,
		Uri:         codeSystem.Uri,
		Name:        codeSystem.Name,
		Title:       codeSystem.Title,
	}
}

func ApiCreateCodeSystemToGormCodeSystem(codeSystem *api.CreateCodeSystem) *models.CodeSystem {
	return &models.CodeSystem{
		Author:      codeSystem.Author,
		Description: codeSystem.Description,
		Uri:         codeSystem.Uri,
		Name:        codeSystem.Name,
		Title:       codeSystem.Title,
	}
}
