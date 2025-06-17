package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"

	"gorm.io/gorm"
)

// this query converts the gorm object to an api object
func (gq *GormQuery) GetConceptReplaceBies(code string, codeSystemId int32, codeSystemVersionId int32) (*[]api.ConceptReplaceBy, error) {
	apiConceptReplaceBies := []api.ConceptReplaceBy{}

	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&models.CodeSystem{}, codeSystemId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
			default:
				return err
			}
		}

		var conceptReplaceBies []models.ConceptReplaceBy = []models.ConceptReplaceBy{}
		if err := tx.Where("code_system_id = ? AND code = ?", codeSystemId, code).Find(&conceptReplaceBies).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return nil
			default:
				return err
			}
		}

		for _, conceptReplaceBy := range conceptReplaceBies {
			var concept = models.Concept{}
			if err := gq.GetConceptQuery(&concept, conceptReplaceBy.MapTo, conceptReplaceBy.CodeSystemID, codeSystemVersionId); err != nil {
				switch {
				case errors.Is(err, gorm.ErrRecordNotFound):
					continue
				default:
					return err
				}
			}
			apiConceptReplaceBies = append(apiConceptReplaceBies, api.ConceptReplaceBy{
				Code:         conceptReplaceBy.Code,
				MapTo:        *transform.GormConceptToApiConcept(&concept),
				CodeSystemId: conceptReplaceBy.CodeSystemID,
				Equivalence:  (*api.ConceptReplaceByEquivalence)(conceptReplaceBy.Equivalence),
				Comment:      conceptReplaceBy.Comment,
			})
		}

		return nil
	})
	return &apiConceptReplaceBies, err
}
