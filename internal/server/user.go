package server

import (
	"context"
	"errors"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
	"miracummapper/internal/utilities"

	"github.com/lestrrat-go/jwx/jwt"
)

// DeleteUser implements api.StrictServerInterface.
func (s *Server) DeleteUser(ctx context.Context, request api.DeleteUserRequestObject) (api.DeleteUserResponseObject, error) {
	userIdString := request.UserId

	var user models.User
	userId, err := utilities.ParseUUID(userIdString)
	if err != nil {
		return api.DeleteUser400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}

	if err := s.Database.DeleteUserQuery(&user, userId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.DeleteUser404JSONResponse(err.Error()), nil
		default:
			return api.DeleteUser500JSONResponse{InternalServerErrorJSONResponse: api.InternalServerErrorJSONResponse(database.InternalServerErrorMessage)}, nil
		}
	}

	return api.DeleteUser200JSONResponse(*transform.GormUserToApiUser(&user)), nil
}

// GetAllUsers implements api.StrictServerInterface.
func (s *Server) GetAllUsers(ctx context.Context, request api.GetAllUsersRequestObject) (api.GetAllUsersResponseObject, error) {
	var users []models.User = []models.User{}

	if err := s.Database.GetAllUsersQuery(&users); err != nil {
		return api.GetAllUsers500JSONResponse{}, err
	}

	var apiUsers []api.User = []api.User{}

	for _, user := range users {
		apiUsers = append(apiUsers, *transform.GormUserToApiUser(&user))
	}

	return api.GetAllUsers200JSONResponse(apiUsers), nil
}

// Login implements api.StrictServerInterface.
func (s *Server) Login(ctx context.Context, request api.LoginRequestObject) (api.LoginResponseObject, error) {
	// Get jwt
	jwt := ctx.Value("jwt_claims").(jwt.Token)
	if jwt == nil {
		return api.Login400JSONResponse{}, nil
	}

	// get required values from jwt for user
	id := jwt.Subject()

	userName, _ := utilities.GetValueFromToken(jwt, "preferred_username", false)

	fullName, _ := utilities.GetValueFromToken(jwt, "name", false)

	email, _ := utilities.GetValueFromToken(jwt, "email", false)

	// create user object
	apiUser := api.User{
		Id:       id,
		Username: userName,
		Email:    &email,
		Fullname: &fullName,
	}

	gormUser, err := transform.ApiUserToGormUser(&apiUser)
	if err != nil {
		return api.Login400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
	}
	if err := s.Database.CreateOrUpdateUserQuery(gormUser); err != nil {
		return api.Login500JSONResponse{InternalServerErrorJSONResponse: api.InternalServerErrorJSONResponse(err.Error())}, nil
	}

	return api.Login200JSONResponse(apiUser), nil

}
