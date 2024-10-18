package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
	"miracummapper/internal/utilities"
)

func GormUserToApiUser(user *models.User) *api.User {
	apiUser := api.User{
		Id:       user.Id.String(),
		Username: user.UserName,
		Email:    &user.Email,
		Fullname: &user.FullName,
	}

	return &apiUser
}

func ApiUserToGormUser(user *api.User) (*models.User, error) {
	userID, err := utilities.ParseUUID(user.Id)
	if err != nil {
		return nil, err
	}
	gormUser := models.User{
		Id:       userID,
		UserName: user.Username,
		FullName: *user.Fullname,
		Email:    *user.Email,
	}

	return &gormUser, nil
}
