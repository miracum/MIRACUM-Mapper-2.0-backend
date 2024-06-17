package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
)

func GormCodeSystemToApiCodeSystem(codeSystem *models.CodeSystem) *api.CodeSystem {
	return &api.CodeSystem{
		Id:          int32(codeSystem.ID),
		Author:      codeSystem.Author,
		Description: codeSystem.Description,
		Name:        codeSystem.Name,
		Title:       codeSystem.Title,
		Uri:         codeSystem.Uri,
		Version:     codeSystem.Version,
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
		Version:     codeSystem.Version,
		Name:        codeSystem.Name,
		Title:       codeSystem.Title,
	}
}

func ApiCreateCodeSystemToGormCodeSystem(codeSystem *api.CreateCodeSystem) *models.CodeSystem {
	return &models.CodeSystem{
		Author:      codeSystem.Author,
		Description: codeSystem.Description,
		Uri:         codeSystem.Uri,
		Version:     codeSystem.Version,
		Name:        codeSystem.Name,
		Title:       codeSystem.Title,
	}
}
