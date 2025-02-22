package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"strings"

	"gorm.io/gorm"
)

func (gq *GormQuery) GetAllConceptsQuery(concepts *[]models.Concept, codeSystemId int32, pageSize int, offset int, sortBy string, sortOrder string, meaning string, code string) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&models.CodeSystem{}, codeSystemId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
			default:
				return err
			}
		}

		query := tx.
			Model(&models.Concept{}).
			Where("code_system_id = ?", codeSystemId)

		// Add code condition if code is not empty
		if code != "" {
			query = query.Where("LOWER(code) LIKE LOWER(?)", code+"%")
		}

		// Add meaning condition if meaning is not empty
		if meaning != "" {
			formattedMeaning := strings.Join(strings.Fields(meaning), ":* & ") + ":*" // Adjust for partial matches
			query = query.Where("display_search_vector @@ to_tsquery('english', ?)", formattedMeaning)
			// Tests for similarity searches. These were very slow and therefore not used in the final implementation. I a search should be implemented which is not part of a autocomplete but is ok to take e.g a second to complete, similarity searches with the pg_trgm extension could be used here.
			// query = query.Where("display_search_vector @@ to_tsquery(?) OR similarity(display, ?) > 0.3", formattedMeaning, meaning)
			// query = query.Where("similarity(display, ?) > 0.8", meaning)
			// query = query.Select("*, similarity(display, ?) > set_limit(0.99) AS s", meaning).Order("s DESC")
		}

		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize)

		return query.Find(&concepts).Error
	})
	return err
}

func (gq *GormQuery) GetAllConceptsByVersionQuery(concepts *[]models.Concept, codeSystemId int32, codeSystemVersionId int32, pageSize int, offset int, sortBy string, sortOrder string, meaning string, code string) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&models.CodeSystem{}, codeSystemId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
			default:
				return err
			}
		}

		var codeSystemVersion models.CodeSystemVersion
		if err := tx.Where("code_system_id = ?", codeSystemId).First(&codeSystemVersion, codeSystemVersionId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found for CodeSystem with ID %d.", codeSystemVersionId, codeSystemId))
			default:
				return err
			}
		}
		versionId := codeSystemVersion.VersionID

		query := tx.
			Preload("ValidFromVersion").
			Preload("ValidToVersion").
			Model(&models.Concept{}).
			Joins("JOIN code_system_versions AS valid_from_version ON valid_from_version.id = concepts.valid_from_version_id").
			Joins("JOIN code_system_versions AS valid_to_version ON valid_to_version.id = concepts.valid_to_version_id").
			Where("concepts.code_system_id = ?", codeSystemId).
			Where("valid_from_version.version_id <= ? AND valid_to_version.version_id >= ?", versionId, versionId)

		// Add code condition if code is not empty
		if code != "" {
			query = query.Where("LOWER(code) LIKE LOWER(?)", code+"%")
		}

		// Add meaning condition if meaning is not empty
		if meaning != "" {
			formattedMeaning := strings.Join(strings.Fields(meaning), ":* & ") + ":*" // Adjust for partial matches
			query = query.Where("display_search_vector @@ to_tsquery('english', ?)", formattedMeaning)
		}

		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize)

		return query.Find(&concepts).Error
	})
	return err
}
