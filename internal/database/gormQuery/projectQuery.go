package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"

	"gorm.io/gorm"
)

func (gq *GormQuery) GetAllProjectsQuery(projects *[]models.Project, pageSize int, offset int, sortBy string, sortOrder string) error {
	db := gq.Database.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize).Find(&projects)
	return db.Error
}

// CreateProject implements database.Datastore.
func (gq *GormQuery) CreateProjectQuery(project *models.Project) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&project).Error
		if err != nil {
			// cast error to postgres error
			pgErr, ok := handlePgError(err)
			if !ok {
				return err
			}
			switch pgErr.Code {
			case "23503":
				switch pgErr.ConstraintName {
				case "fk_users_project_permissions":
					userID, err := extractIDFromErrorDetail(pgErr.Detail, "user_id")
					if err != nil {
						return err
					}
					return database.NewDBError(database.ClientError, fmt.Sprintf("User with id %s specified in permissions does not exist", userID))

				case "fk_code_systems_code_system_roles":
					codeSystemID, err := extractIDFromErrorDetail(pgErr.Detail, "code_system_id")
					if err != nil {
						return err
					}
					return database.NewDBError(database.ClientError, fmt.Sprintf("Code System with id %s specified in code system roles does not exist", codeSystemID))
				default:
					return err
				}
			default:
				return err
			}
		} else {
			if err = tx.Preload("CodeSystemRoles").Preload("CodeSystemRoles.CodeSystem").Preload("Permissions.User").First(&project, project.ID).Error; err != nil {
				return nil
			}
			return err
		}
	})
	return err
}

func (gq *GormQuery) GetProjectQuery(project *models.Project, projectId int32) error {
	db := gq.Database.Preload("CodeSystemRoles", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC") // TODO maybe sort at the end is also possible
	}).Preload("CodeSystemRoles.CodeSystem").Preload("Permissions.User").First(&project, projectId) // TODO & is not necessary here
	if db.Error != nil {
		pgErr, ok := handlePgError(db.Error)
		if db.Error.Error() == "record not found" {
			return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
		}
		if !ok {
			return db.Error
		}
		switch {
		case errors.Is(pgErr, gorm.ErrRecordNotFound):
			return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
		default:
			return pgErr
		}
	} else {
		return nil
	}
}

// UpdateProject implements database.Datastore.
func (gq *GormQuery) UpdateProjectQuery(project *models.Project, checkFunc func(oldProject, newProject *models.Project) error) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		project_old := models.Project{}

		if err := tx.First(&project_old, project.ID).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", project.ID))
			default:
				return err
			}
		}

		// Call the check function
		if err := checkFunc(&project_old, project); err != nil {
			return err
		}

		// // TODO avoid checking these fields in this function, move to endpoint logic
		// if project_old.StatusRequired != project.StatusRequired || project_old.EquivalenceRequired != project.EquivalenceRequired {
		// 	return database.NewDBError(database.ClientError, "StatusRequired and EquivalenceRequired cannot be updated")
		// }

		// won't create new record because tx.First already checked that a project with that ID exists
		if err := tx.Save(&project).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteProject implements database.Datastore.
func (gq *GormQuery) DeleteProjectQuery(project *models.Project, projectId int32) error {

	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		// get project so it can be returned in the api and then delete it
		if err := tx.First(&project, projectId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
			default:
				return err
			}
		}

		db := tx.Delete(&project, projectId)
		if db.Error != nil {
			// switch {
			// case errors.Is(db.Error, gorm.ErrRecordNotFound):
			// 	return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
			// default:
			return db.Error
			// }
		} else {
			if db.RowsAffected == 0 {
				return database.NewDBError(database.NotFound, fmt.Sprintf("Project with ID %d couldn't be found.", projectId))
			}
			return nil
		}
	})
	return err
}
