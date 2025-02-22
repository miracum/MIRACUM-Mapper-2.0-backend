package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"time"

	"gorm.io/gorm"
)

func (gq *GormQuery) GetCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion, codeSystemId int32, codeSystemVersionId int32) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&models.CodeSystem{}, codeSystemVersion.CodeSystemID).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemVersion.CodeSystemID))
			default:
				return err
			}
		}
		if err := tx.Where("code_system_id = ?", codeSystemId).First(&codeSystemVersion, codeSystemVersionId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found for CodeSystem with ID %d.", codeSystemVersionId, codeSystemId))
			default:
				return err
			}
		} else {
			return nil
		}
	})
	return err
}

func (gq *GormQuery) CreateCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&models.CodeSystem{}, codeSystemVersion.CodeSystemID).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemVersion.CodeSystemID))
			default:
				return err
			}
		}

		if exists, err := checkReleaseDateExists(tx, codeSystemVersion.CodeSystemID, codeSystemVersion.ReleaseDate); err != nil {
			return err
		} else if exists {
			return database.NewDBError(database.ClientError, fmt.Sprintf("CodeSystemVersion with release date %s already exists for CodeSystem with ID %d.", codeSystemVersion.ReleaseDate, codeSystemVersion.CodeSystemID))
		}

		if nextVersionId, newerCodeSystemVersions, err := getNextVersionIdAndNewerCodeSystemVersions(tx, codeSystemVersion.CodeSystemID, codeSystemVersion.ReleaseDate); err != nil {
			return err
		} else {
			if len(newerCodeSystemVersions) > 0 {
				for _, newerCodeSystemVersion := range newerCodeSystemVersions {
					newerCodeSystemVersion.VersionID++
					if err := tx.Save(&newerCodeSystemVersion).Error; err != nil {
						return err
					}
				}
			}
			codeSystemVersion.VersionID = nextVersionId
		}

		return tx.Create(&codeSystemVersion).Error
	})
	return err
}

func checkReleaseDateExists(db *gorm.DB, codeSystemID int32, releaseDate time.Time) (bool, error) {
	if err := db.First(&models.CodeSystemVersion{}, "code_system_id = ? AND release_date = ?", codeSystemID, releaseDate).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func getNextVersionIdAndNewerCodeSystemVersions(db *gorm.DB, codeSystemID int32, releaseDate time.Time) (int32, []models.CodeSystemVersion, error) {
	var newerCodeSystemVersions []models.CodeSystemVersion
	if err := db.Order("release_date DESC").Find(&newerCodeSystemVersions, "code_system_id = ? AND release_date > ?", codeSystemID, releaseDate).Error; err != nil {
		return 0, nil, err
	}
	if len(newerCodeSystemVersions) > 0 {
		return newerCodeSystemVersions[len(newerCodeSystemVersions)-1].VersionID, newerCodeSystemVersions, nil
	} else {
		var lastCodeSystemVersion models.CodeSystemVersion
		if err := db.Order("version_id DESC").Limit(1).Find(&lastCodeSystemVersion, "code_system_id = ?", codeSystemID).Error; err != nil {
			return 0, nil, err
		} else {
			return lastCodeSystemVersion.VersionID + 1, nil, nil
		}
	}
}

func (gq *GormQuery) UpdateCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&models.CodeSystem{}, codeSystemVersion.CodeSystemID).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemVersion.CodeSystemID))
			default:
				return err
			}
		}

		var existingCodeSystemVersion models.CodeSystemVersion
		if err := tx.Where("code_system_id = ?", codeSystemVersion.CodeSystemID).First(&existingCodeSystemVersion, codeSystemVersion.ID).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found for CodeSystem with ID %d.", codeSystemVersion.ID, codeSystemVersion.CodeSystemID))
			default:
				return err
			}
		}

		codeSystemVersion.VersionID = existingCodeSystemVersion.VersionID
		codeSystemVersion.ReleaseDate = existingCodeSystemVersion.ReleaseDate
		codeSystemVersion.Imported = existingCodeSystemVersion.Imported
		return tx.Save(&codeSystemVersion).Error
	})
	return err
}

func (gq *GormQuery) DeleteCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion, codeSystemVersionId int32) error {
	return database.NewDBError(database.ClientError, "CodeSystemVersion cannot be deleted at the moment.")
	// err := gq.Database.Transaction(func(tx *gorm.DB) error {
	// 	if err := tx.First(&codeSystemVersion, codeSystemVersionId).Error; err != nil {
	// 		switch {
	// 		case errors.Is(err, gorm.ErrRecordNotFound):
	// 			return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found.", codeSystemVersionId))
	// 		default:
	// 			return err
	// 		}
	// 	}

	// 	codeSystemRoles := []models.CodeSystemRole{}
	// 	if err := tx.Find(&codeSystemRoles, "code_system_version_id = ? OR next_code_system_version_id = ?", codeSystemVersionId, codeSystemVersionId).Error; err == nil {
	// 		if len(codeSystemRoles) > 0 {
	// 			projectIds := []string{}
	// 			for _, role := range codeSystemRoles {
	// 				projectIds = append(projectIds, fmt.Sprintf("Id: %d", role.ProjectID))
	// 			}
	// 			return database.NewDBError(database.ClientError, fmt.Sprintf("CodeSystemVersion cannot be deleted if it is in use in these projects: %s", strings.Join(projectIds, ", ")))
	// 		}
	// 	}

	// 	db := tx.Delete(&codeSystemVersion, codeSystemVersionId)
	// 	if db.Error != nil {
	// 		return db.Error
	// 	} else {
	// 		if db.RowsAffected == 0 {
	// 			return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found.", codeSystemVersionId))
	// 		}
	// 		return nil
	// 	}
	// })
	// return err
}

// func (gq *GormQuery) CreateConceptsQuery(concepts *[]models.Concept) error {
// 	if len(*concepts) == 0 {
// 		return nil
// 	}

// 	err := gq.Database.Transaction(func(tx *gorm.DB) error {
// 		batchSize := 100
// 		totalConcepts := len(*concepts)

// 		// Create batch of 100 concepts at a time in the database for better performance
// 		for i := 0; i < totalConcepts; i += batchSize {
// 			end := i + batchSize
// 			if end > totalConcepts {
// 				end = totalConcepts
// 			}

// 			batch := (*concepts)[i:end]

// 			// The log level is set to silent as the batch create can create a huge amount of logs slowing down the create process significantly for huge code systems
// 			if err := tx.Create(&batch).Error; err != nil { // this results in a extremely huge log(in debug mode), consider using this: .Session(&gorm.Session{Logger: tx.Logger.LogMode(logger.Silent)})
// 				return err
// 			}
// 		}
// 		return nil
// 	})

// 	return err
// }
