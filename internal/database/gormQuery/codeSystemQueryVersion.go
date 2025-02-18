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
					if err := gq.Database.Save(&newerCodeSystemVersion).Error; err != nil { //Model(&models.CodeSystemVersion{}).Where("id = ?", newerCodeSystemVersion.ID).Update("version_id", newerCodeSystemVersion.VersionID).Error; err != nil {
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

func (gq *GormQuery) checkReleaseDateExists(codeSystemID uint32, releaseDate time.Time) (bool, error) {
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

func (gq *GormQuery) getNextVersionIdAndNewerCodeSystemVersions(codeSystemID uint32, releaseDate time.Time) (uint32, []models.CodeSystemVersion, error) {
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
		// get codeSystem so it can be returned in the api and then delete it
		if err := tx.First(&codeSystemVersion, codeSystemVersionId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found.", codeSystemVersionId))
			default:
				return err
			}
		}

		codeSystemRoles := []models.CodeSystemRole{}
		// if err := tx.Where(&models.CodeSystemRole{CodeSystemVersionID: uint32(codeSystemVersionId)}).Or(&models.CodeSystemRole{NextCodeSystemVersionID: uint32(codeSystemVersionId)}).Find(&codeSystemRoles).Error; err == nil {
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

func (gq *GormQuery) SetCodeSystemVersionImported(codeSystemVersionId int32, imported bool) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		var codeSystemVersion models.CodeSystemVersion
		if err := tx.First(&codeSystemVersion, codeSystemVersionId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found.", codeSystemVersionId))
			default:
				return err
			}
		}

		codeSystemVersion.Imported = imported
		if err := tx.Save(&codeSystemVersion).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// func (gq *GormQuery) CheckHasNoConceptsQuery(codeSystemId int32, codeSystemVersionId int32) error {
// 	var codeSystem models.CodeSystem
// 	var codeSystemVersion models.CodeSystemVersion
// 	var concept models.Concept
// 	err := gq.Database.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.First(&codeSystem, codeSystemId).Error; err != nil {
// 			switch {
// 			case errors.Is(err, gorm.ErrRecordNotFound):
// 				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
// 			default:
// 				return err
// 			}
// 		}
// 		if err := tx.Where("code_system_id = ?", codeSystemId).First(&codeSystemVersion, codeSystemVersionId).Error; err != nil {
// 			switch {
// 			case errors.Is(err, gorm.ErrRecordNotFound):
// 				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystemVersion with ID %d couldn't be found for CodeSystem with ID %d.", codeSystemVersionId, codeSystemId))
// 			default:
// 				return err
// 			}
// 		}

// 		query := tx.
// 			Preload("ValidFromVersion").
// 			Preload("ValidToVersion").
// 			Model(&models.Concept{}).
// 			Joins("JOIN code_system_versions AS valid_from_version ON valid_from_version.id = concepts.valid_from_version_id").
// 			Joins("JOIN code_system_versions AS valid_to_version ON valid_to_version.id = concepts.valid_to_version_id").
// 			Where("concepts.code_system_id = ?", codeSystemId).
// 			Where("valid_from_version.version_id <= ? AND valid_to_version.version_id >= ?", codeSystemVersion.VersionID, codeSystemVersion.VersionID).
// 			First(&concept)

// 		if err := query.Error; err != nil {
// 			switch {
// 			case errors.Is(err, gorm.ErrRecordNotFound):
// 				// Good case, no concepts found
// 				return nil
// 			default:
// 				return err
// 			}
// 		}
// 		// Bad case, concepts found
// 		return database.NewDBError(database.ClientError, fmt.Sprintf("CodeSystemVersion with ID %d has concepts.", codeSystemVersionId))
// 	})
// 	return err
// }

func (gq *GormQuery) GetImportedNeighborVersionIdsQuery(codeSystemId int32, codeSystemVersionId int32) (*uint32, *uint32, error) {
	var codeSystemVersion models.CodeSystemVersion
	if err := gq.GetCodeSystemVersionQuery(&codeSystemVersion, codeSystemId, codeSystemVersionId); err != nil {
		return nil, nil, err
	}

	var beforeVersionId *uint32
	beforeVersionId = nil
	var afterVersionId *uint32
	afterVersionId = nil

	var beforeCodeSystemVersions []models.CodeSystemVersion
	if err := gq.Database.Where("code_system_id = ? AND version_id < ?", codeSystemId, codeSystemVersion.VersionID).Order("version_id DESC").Find(&beforeCodeSystemVersions).Error; err != nil {
		return nil, nil, err
	}

	if len(beforeCodeSystemVersions) > 0 {
		for _, beforeCodeSystemVersion := range beforeCodeSystemVersions {
			if beforeCodeSystemVersion.Imported {
				beforeVersionId = &beforeCodeSystemVersion.VersionID
				break
			}
			// if err := gq.CheckHasNoConceptsQuery(codeSystemId, int32(beforeCodeSystemVersion.ID)); err != nil {
			// 	if errors.Is(err, database.ErrClientError) {
			// 		beforeVersionId = &beforeCodeSystemVersion.VersionID
			// 		break
			// 	} else {
			// 		return nil, nil, err
			// 	}
			// }
		}
	}

	var afterCodeSystemVersions []models.CodeSystemVersion
	if err := gq.Database.Where("code_system_id = ? AND version_id > ?", codeSystemId, codeSystemVersion.VersionID).Order("version_id ASC").Find(&afterCodeSystemVersions).Error; err != nil {
		return nil, nil, err
	}

	if len(afterCodeSystemVersions) > 0 {
		for _, afterCodeSystemVersion := range afterCodeSystemVersions {
			if afterCodeSystemVersion.Imported {
				afterVersionId = &afterCodeSystemVersion.VersionID
				break
			}
			// if err := gq.CheckHasNoConceptsQuery(codeSystemId, int32(afterCodeSystemVersion.ID)); err != nil {
			// 	if errors.Is(err, database.ErrClientError) {
			// 		afterVersionId = &afterCodeSystemVersion.VersionID
			// 	} else {
			// 		return nil, nil, err
			// 	}
			// }
		}
	}

	return beforeVersionId, afterVersionId, nil
}

// func (gq *GormQuery) GetFirstElementCodeSystemVersionQuery(codeSystem *models.CodeSystem, codeSystemId int32, concept *models.Concept) error {
// 	err := gq.Database.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.First(&codeSystem, codeSystemId).Error; err != nil {
// 			switch {
// 			case errors.Is(err, gorm.ErrRecordNotFound):
// 				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystem.ID))
// 			default:
// 				return err
// 			}
// 		}
// 		if err := tx.Where("code_system_id", codeSystemId).First(&concept).Error; err != nil {
// 			switch {
// 			case errors.Is(err, gorm.ErrRecordNotFound):
// 				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d has no elements.", codeSystem.ID))
// 			default:
// 				return err
// 			}
// 		}

// 		return nil
// 	})

// 	return err
// }

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
