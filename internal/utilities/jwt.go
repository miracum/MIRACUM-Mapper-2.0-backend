package utilities

import (
	"fmt"

	"github.com/lestrrat-go/jwx/jwt"
)

func GetValueFromToken(token jwt.Token, key string, required bool) (string, error) {
	id, found := token.Get(key)
	if !found {
		if required {
			return "", fmt.Errorf("key not found")
		}
		return "", nil
	}
	idString, ok := id.(string)
	if !ok {
		if required {
			return "", fmt.Errorf("value is not a string")
		}
		return "", nil
	}
	return idString, nil
}
