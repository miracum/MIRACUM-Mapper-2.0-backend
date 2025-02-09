package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"strings"

	"gorm.io/gorm"
)

// func (gq *GormQuery) GetAllCodeSystemsQuery(codeSystems *[]models.CodeSystem) error {
// 	db := gq.Database.Find(&codeSystems)
// 	return db.Error
// }

func (gq *GormQuery) CreateCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion) error {
	return gq.Database.Create(&codeSystemVersion).Error
}

// func (gq *GormQuery) GetCodeSystemQuery(codeSystem *models.CodeSystem, codeSystemId int32) error {
// 	db := gq.Database.First(&codeSystem, codeSystemId)
// 	if db.Error != nil {
// 		switch {
// 		case errors.Is(db.Error, gorm.ErrRecordNotFound):
// 			return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
// 		default:
// 			return db.Error
// 		}
// 	} else {
// 		return nil
// 	}
// }

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

// func (gq *GormQuery) GetFirstElementCodeSystemQuery(codeSystem *models.CodeSystem, codeSystemId int32, concept *models.Concept) error {
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
