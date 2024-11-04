package server

import (
	"context"
	"miracummapper/internal/server/middlewares"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
)

func GetOffset(page int, pageSize int) int {
	return (page - 1) * pageSize
}

func GetJWTFromContext(ctx context.Context) (jwt.Token, error) {
	jwt := ctx.Value("jwt_claims").(jwt.Token)
	if jwt == nil {
		return jwt, nil
	}
	return jwt, nil
}

func GetUserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	jwt, err := GetJWTFromContext(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}
	id := jwt.Subject()
	uuid, err := uuid.Parse(id)
	if err != nil {
		return uuid, err
	}
	return uuid, nil
}

func IsAdminFromContext(ctx context.Context) bool {
	jwt, err := GetJWTFromContext(ctx)
	if err != nil {
		return false
	}
	err = middlewares.CheckTokenClaims([]string{"abc"}, jwt) // config.KeycloakAdminScope
	return err == nil
}
