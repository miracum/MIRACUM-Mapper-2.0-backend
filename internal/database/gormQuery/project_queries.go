package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"

	"gorm.io/gorm"
)

func (gq *GormQuery) GetProjectQuery(project *models.Project, projectId int32) error {
	db := gq.Database.Preload("CodeSystemRoles", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC")
	}).Preload("CodeSystemRoles.CodeSystem").Preload("Permissions.User").First(&project, projectId)
	if db.Error != nil {
		switch {
		case errors.Is(db.Error, gorm.ErrRecordNotFound):
			return database.ErrRecordNotFound
		default:
			return db.Error
		}
	} else {
		return nil
	}
}

// AddProject implements database.Datastore.
func (gq *GormQuery) AddProjectQuery(project *models.Project) error {
	db := gq.Database.Create(&project)
	// if strings.Contains(err.Error(), "foreign key") && strings.Contains(err.Error(), "user_id") {
	// 	// Extract the user ID from the error message
	// 	regex := regexp.MustCompile(`\((.*?)\)`)
	// 	matches := regex.FindStringSubmatch(err.Error())
	// 	if len(matches) > 1 {
	// 		userID := matches[1]
	// 		return api.AddProject422JSONResponse(fmt.Sprintf("User with ID %s does not exist", userID)), nil
	// 	}
	// }
	// if(db.Error != nil) {
	// 	switch {
	// 		case errors.Is(db.Error, gorm.ErrForeignKeyViolated):
	// 			if(strings.Contains(err.Error(), "fk_users_project_permissions")){

	// 			}

	// 	}
	// }
	return db.Error
}

// DeleteProject implements database.Datastore.
func (gq *GormQuery) DeleteProjectQuery(project *models.Project, projectId int32) error {
	db := gq.Database.Delete(&project, projectId)
	switch {
	case errors.Is(db.Error, gorm.ErrRecordNotFound):
		return database.ErrRecordNotFound
	default:
		return db.Error
	}
}

// GetProjectsQuery implements database.Datastore.
func (gq *GormQuery) GetProjectsQuery(projects *[]models.Project, pageSize int, offset int, sortBy string, sortOrder string) error {
	db := gq.Database.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize).Find(&projects)
	return db.Error
}

// UpdateProject implements database.Datastore.
func (gq *GormQuery) UpdateProjectQuery(project *models.Project) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		project_old := models.Project{}

		if err := tx.First(&project_old, project.ID).Error; err != nil {
			return err
			// TODO insert correct error: check if project exists (record not found)
		}

		// TODO avoid checking these fields in this function, move to endpoint logic
		if project_old.StatusRequired != project.StatusRequired || project_old.EquivalenceRequired != project.EquivalenceRequired {
			// TODO change error
			return errors.New("StatusRequired and EquivalenceRequired cannot be updated")
		}

		// won't create new record because tx.First already checked that a project with that ID exists
		if err := tx.Save(&project).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
