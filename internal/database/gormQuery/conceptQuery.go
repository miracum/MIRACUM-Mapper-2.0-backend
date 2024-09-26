package gormQuery

import (
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"strings"

	"gorm.io/gorm"
)

func (gq *GormQuery) GetAllConceptsQuery(concepts *[]models.Concept, codeSystemId int32, pageSize int, offset int, sortBy string, sortOrder string, meaning string, code string) error {
	// return gq.Database.Transaction(func(tx *gorm.DB) error {
	// 	if err := tx.Where("display ILIKE ? AND code ILIKE ?", "%"+meaning+"%", code+"%").Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize).Find(&concepts).Error; err != nil {
	// 		return err
	// 	} else if len(*concepts) == 0 {
	// 		*concepts = []models.Concept{}
	// 		return nil
	// 	}
	// 	return nil
	// })

	// Start building the query
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		query := tx.
			Model(&models.Concept{}).
			Where("code_system_id = ?", codeSystemId)

		// Add code condition if code is not empty
		if code != "" {
			query = query.Where("LOWER(code) ILIKE ?", strings.ToLower(code)+"%")
		}

		// Add meaning condition if meaning is not empty
		if meaning != "" {
			formattedMeaning := strings.Join(strings.Fields(meaning), ":* & ") + ":*" // Adjust for partial matches
			// query = query.Where("display_search_vector @@ to_tsquery(?) OR similarity(display, ?) > 0.3", formattedMeaning, meaning)
			query = query.Where("display_search_vector @@ to_tsquery(?)", formattedMeaning)
			// query = query.Where("similarity(display, ?) > 0.8", meaning)
			// query = query.Select("*, similarity(display, ?) > set_limit(0.99) AS s", meaning).Order("s DESC")
		}

		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize)

		// Execute the query
		if err := query.Find(&concepts).Error; err != nil {
			return err
		} else if len(*concepts) == 0 {
			var codesystem models.CodeSystem
			if err := tx.First(&codesystem, codeSystemId).Error; err != nil {
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
			} else {
				*concepts = []models.Concept{}
				return nil
			}
		}
		return nil
	})
	return err
}
