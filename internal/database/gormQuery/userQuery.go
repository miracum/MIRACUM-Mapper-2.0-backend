package gormQuery

import (
	"errors"
	"fmt"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetAllUsersQuery implements database.Datastore.
func (gq *GormQuery) GetAllUsersQuery(users *[]models.User) error {
	db := gq.Database.Find(&users)
	return db.Error
}

// DeleteUserQuery implements database.Datastore.
func (gq *GormQuery) DeleteUserQuery(user *models.User, userId uuid.UUID) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		// get user so it can be returned in the api and then delete it
		if err := tx.First(&user, userId).Error; err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return database.NewDBError(database.NotFound, fmt.Sprintf("User with ID %d couldn't be found.", userId))
			default:
				return err
			}
		}

		db := tx.Delete(&user, userId)
		if db.Error != nil {
			return db.Error
		} else {
			if db.RowsAffected == 0 {
				return database.NewDBError(database.NotFound, fmt.Sprintf("User with ID %d couldn't be found.", userId))
			}
			return nil
		}
	})
	return err
}

func (gq *GormQuery) CreateOrUpdateUserQuery(user *models.User) error {
	err := gq.Database.Transaction(func(tx *gorm.DB) error {
		db := tx.Save(user)
		if db.Error != nil {
			return db.Error
		}

		return nil
	})
	return err
}
