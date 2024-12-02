package middlewares

////////////////////////////////////////////////////////////////
// CODE COPIED FROM OFFICIAL DOCUMENTATION AND MODIFIED
// see: https://github.com/oapi-codegen/oapi-codegen/tree/main/examples/authenticated-api
////////////////////////////////////////////////////////////////

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/lestrrat-go/jwx/jwt"
	middleware "github.com/oapi-codegen/gin-middleware"
)

// JWSValidator is used to validate JWS payloads and return a JWT if they're
// valid
type JWSValidator interface {
	ValidateJWS(jws string) (jwt.Token, error)
}

const JWTClaimsContextKey = "jwt_claims"

var (
	ErrNoAuthHeader      = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("authorization header is malformed")
	ErrClaimsInvalid     = errors.New("provided claims do not match expected scopes")
	ErrorTokenExpired    = errors.New("token Expired")
)

var ErrorTokenExpiredApi = &AuthenticationError{
	Reason:   "token expired",
	Response: &ResponseError{StatusCode: http.StatusUnauthorized},
}

type AuthenticationError struct {
	Reason   string
	Response *ResponseError
}

func (e *AuthenticationError) Error() string {
	return e.Reason
}

type ResponseError struct {
	StatusCode int
}

// GetJWSFromRequest extracts a JWS string from an Authorization: Bearer <jws> header
func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	// Check for the Authorization header.
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	// We expect a header value of the form "Bearer <token>", with 1 space after
	// Bearer, per spec.
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

func NewAuthenticate(v JWSValidator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(v, ctx, input)
	}
}

// Authenticate uses the specified validator to ensure a JWT is valid, then makes
// sure that the claims provided by the JWT match the scopes as required in the API.
func Authenticate(v JWSValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	// Our security scheme is named BearerAuth, ensure this is the case
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	// Now, we need to get the JWS from the request, to match the request expectations
	// against request contents.
	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jws: %w", err)
	}

	// if the JWS is valid, we have a JWT, which will contain a bunch of claims.
	token, err := v.ValidateJWS(jws)
	if err != nil {
		if errors.Is(err, ErrorTokenExpired) {
			return ErrorTokenExpiredApi
		}
		return fmt.Errorf("validating JWS: %w", err)
	}

	// We've got a valid token now, and we can look into its claims to see whether
	// they match. Every single scope must be present in the claims.
	err = CheckTokenClaims(input.Scopes, token)

	if err != nil {
		return fmt.Errorf("token claims don't match: %w", err)
	}

	// Set the property on the echo context so the handler is able to
	// access the claims data we generate in here.
	eCtx := middleware.GetGinContext(ctx)
	eCtx.Set(JWTClaimsContextKey, token)

	return nil
}

// GetClaimsFromToken returns a list of roles from the token for the specified client.
func GetClaimsFromToken(t jwt.Token) ([]string, error) {
	resourceAccess, found := t.Get("resource_access")
	if !found {
		// If the resource_access claim isn't found, return an empty list of roles.
		return make([]string, 0), nil
	}

	resourceAccessMap, ok := resourceAccess.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("resource_access claim is unexpected type")
	}

	clientAccess, found := resourceAccessMap[ClientID]
	if !found {
		// If the client-specific access isn't found, return an empty list of roles.
		return make([]string, 0), nil
	}

	clientAccessMap, ok := clientAccess.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("client access claim is unexpected type")
	}

	roles, found := clientAccessMap["roles"]
	if !found {
		// If the roles aren't found, return an empty list of roles.
		return make([]string, 0), nil
	}

	rolesList, ok := roles.([]interface{})
	if !ok {
		return nil, fmt.Errorf("roles claim is unexpected type")
	}

	claims := make([]string, len(rolesList))
	for i, role := range rolesList {
		claims[i], ok = role.(string)
		if !ok {
			return nil, fmt.Errorf("roles[%d] is not a string", i)
		}
	}

	return claims, nil
}

func CheckTokenClaims(expectedClaims []string, t jwt.Token) error {
	claims, err := GetClaimsFromToken(t)
	if err != nil {
		return fmt.Errorf("getting claims from token: %w", err)
	}
	// Put the claims into a map, for quick access.
	claimsMap := make(map[string]bool, len(claims))
	for _, c := range claims {
		claimsMap[c] = true
	}

	for _, e := range expectedClaims {
		if claimsMap[e] {
			return nil
		}
	}
	return ErrClaimsInvalid
}
