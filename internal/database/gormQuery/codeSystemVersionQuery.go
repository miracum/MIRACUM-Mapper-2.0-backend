package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

func (gq *GormQuery) CreateCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := gq.Database.First(&models.CodeSystem{}, codeSystemVersion.CodeSystemID).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemVersion.CodeSystemID))
			default:
				return err
			}
		}

		if exists, err := gq.checkReleaseDateExists(codeSystemVersion.CodeSystemID, codeSystemVersion.ReleaseDate); err != nil {
			return err
		} else if exists {
			return database.NewDBError(database.ClientError, fmt.Sprintf("CodeSystemVersion with release date %s already exists for CodeSystem with ID %d.", codeSystemVersion.ReleaseDate, codeSystemVersion.CodeSystemID))
		}

		if nextVersionId, newerCodeSystemVersions, err := gq.getNextVersionIdAndNewerCodeSystemVersions(codeSystemVersion.CodeSystemID, codeSystemVersion.ReleaseDate); err != nil {
			return err
		} else {
			if len(newerCodeSystemVersions) > 0 {
				for _, newerCodeSystemVersion := range newerCodeSystemVersions {
					newerCodeSystemVersion.VersionID++
					if err := gq.Database.Save(&newerCodeSystemVersion).Error; err != nil {
						return err
					}
				}
			}
			codeSystemVersion.VersionID = nextVersionId
		}

		if err := gq.Database.Create(&codeSystemVersion).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (gq *GormQuery) checkReleaseDateExists(codeSystemID int32, releaseDate time.Time) (bool, error) {
	var codeSystemVersion models.CodeSystemVersion
	query := gq.Database.First(&codeSystemVersion, "code_system_id = ? AND release_date = ?", codeSystemID, releaseDate)
	if err := query.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (gq *GormQuery) getNextVersionIdAndNewerCodeSystemVersions(codeSystemID int32, releaseDate time.Time) (int32, []models.CodeSystemVersion, error) {
	var newerCodeSystemVersions []models.CodeSystemVersion
	query := gq.Database.Order("release_date DESC").Find(&newerCodeSystemVersions, "code_system_id = ? AND release_date > ?", codeSystemID, releaseDate)
	if err := query.Error; err != nil {
		return 0, nil, err
	}
	if len(newerCodeSystemVersions) > 0 {
		return newerCodeSystemVersions[len(newerCodeSystemVersions)-1].VersionID, newerCodeSystemVersions, nil
	} else {
		var lastCodeSystemVersion models.CodeSystemVersion
		if err := gq.Database.Order("version_id DESC").Limit(1).Find(&lastCodeSystemVersion, "code_system_id = ?", codeSystemID).Error; err != nil {
			return 0, nil, err
		} else {
			return lastCodeSystemVersion.VersionID + 1, nil, nil
		}
	}
}

func (gq *GormQuery) GetCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion, codeSystemId int32, codeSystemVersionId int32) error {
	var codeSystem models.CodeSystem
	if err := gq.GetCodeSystemQuery(&codeSystem, codeSystemId); err != nil {
		return err
	}
	db := gq.Database.First(&codeSystemVersion, codeSystemVersionId)
	if db.Error != nil {
		switch {
		case errors.Is(db.Error, gorm.ErrRecordNotFound):
			return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found.", codeSystemVersionId))
		default:
			return db.Error
		}
	} else {
		return nil
	}
}

func (gq *GormQuery) UpdateCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&models.CodeSystemVersion{}, codeSystemVersion.ID).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found.", codeSystemVersion.ID))
			default:
				return err
			}
		}

		if err := tx.Save(&codeSystemVersion).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (gq *GormQuery) DeleteCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion, codeSystemVersionId int32) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&codeSystemVersion, codeSystemVersionId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found.", codeSystemVersionId))
			default:
				return err
			}
		}

		codeSystemRoles := []models.CodeSystemRole{}
		if err := tx.Find(&codeSystemRoles, "code_system_version_id = ? OR next_code_system_version_id = ?", codeSystemVersionId, codeSystemVersionId).Error; err == nil {
			if len(codeSystemRoles) > 0 {
				projectIds := []string{}
				for _, role := range codeSystemRoles {
					projectIds = append(projectIds, fmt.Sprintf("Id: %d", role.ProjectID))
				}
				return database.NewDBError(database.ClientError, fmt.Sprintf("CodeSystemVersion cannot be deleted if it is in use in these projects: %s", strings.Join(projectIds, ", ")))
			}
		}

		db := tx.Delete(&codeSystemVersion, codeSystemVersionId)
		if db.Error != nil {
			return db.Error
		} else {
			if db.RowsAffected == 0 {
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found.", codeSystemVersionId))
			}
			return nil
		}
	})
	return err
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
