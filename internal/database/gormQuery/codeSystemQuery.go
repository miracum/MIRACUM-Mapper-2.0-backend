package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"strings"

	"gorm.io/gorm"
)

func (gq *GormQuery) GetAllCodeSystemsQuery(codeSystems *[]models.CodeSystem) error {
	return gq.Database.Preload("CodeSystemVersions").Find(&codeSystems).Error
}

func (gq *GormQuery) CreateCodeSystemQuery(codeSystem *models.CodeSystem) error {
	return gq.Database.Create(&codeSystem).Error
}

func (gq *GormQuery) GetCodeSystemQuery(codeSystem *models.CodeSystem, codeSystemId int32) error {
	if err := gq.Database.Preload("CodeSystemVersions").First(&codeSystem, codeSystemId).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
		default:
			return err
		}
	} else {
		return nil
	}
}

func (gq *GormQuery) UpdateCodeSystemQuery(codeSystem *models.CodeSystem) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&models.CodeSystem{}, codeSystem.ID).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystem.ID))
			default:
				return err
			}
		}

		return tx.Save(&codeSystem).Error
	})
	return err
}

func (gq *GormQuery) DeleteCodeSystemQuery(codeSystem *models.CodeSystem, codeSystemId int32) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		// get codeSystem so it can be returned in the api and then delete it
		if err := tx.Preload("CodeSystemRoles").First(&codeSystem, codeSystemId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemId))
			default:
				return err
			}
		}

		// check if codeSystem is used in any projects
		if len(codeSystem.CodeSystemRoles) > 0 {
			projectIds := []string{}
			for _, role := range codeSystem.CodeSystemRoles {
				projectIds = append(projectIds, fmt.Sprintf("Id: %d", role.ProjectID))
			}
			return database.NewDBError(database.ClientError, fmt.Sprintf("CodeSystem cannot be deleted if it is in use in these projects: %s", strings.Join(projectIds, ", ")))
		}

		return tx.Delete(&codeSystem, codeSystemId).Error
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
