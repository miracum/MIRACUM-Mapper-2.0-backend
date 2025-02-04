package middlewares

////////////////////////////////////////////////////////////////
// CODE COPIED FROM OFFICIAL DOCUMENTATION AND MODIFIED
// see: https://github.com/oapi-codegen/oapi-codegen/tree/main/examples/authenticated-api
////////////////////////////////////////////////////////////////

import (
	"fmt"
	"miracummapper/internal/config"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

const HardCodedAudience = "account"
const PermissionsClaim = "resource_access"

var ClientID string
var Issuer string

type Authenticator struct {
	KeySet jwk.Set
}

var _ JWSValidator = (*Authenticator)(nil)

// NewAuthenticator creates an authenticator example which uses a hard coded
func NewAuthenticator(keySet jwk.Set, config *config.Config) (*Authenticator, error) {
	ClientID = config.Env.KeycloakClientId
	// String with Hostname and realm
	Issuer = config.Env.KeycloakUrl + "/realms/" + config.Env.KeycloakRealm
	return &Authenticator{KeySet: keySet}, nil

	// Could be used to create own access tokens without using keycloak
	// normalJWS, err := fa.CreateJWSWithClaims([]string{"normal"})
	// if err != nil {
	// 	log.Fatalln("error creating normal JWS:", err)
	// }

	// adminJWS, err := fa.CreateJWSWithClaims([]string{"admin"})
	// if err != nil {
	// 	log.Fatalln("error creating admin JWS:", err)
	// }

	// log.Println("Normal token", string(normalJWS))
	// log.Println("Admin token", string(adminJWS))
}

// ValidateJWS ensures that the critical JWT claims needed to ensure that we
// trust the JWT are present and with the correct values.
func (f *Authenticator) ValidateJWS(jwsString string) (jwt.Token, error) {

	token, err := jwt.Parse([]byte(jwsString), jwt.WithKeySet(f.KeySet),
		jwt.WithAudience(HardCodedAudience), jwt.WithIssuer(Issuer))
	if err != nil {
		return nil, fmt.Errorf("parsing JWT: %w", err)
	}

	// Check if the token is expired
	if err := jwt.Validate(token, jwt.WithClock(jwt.ClockFunc(time.Now))); err != nil {
		// check error for message exp not satisfied
		if err.Error() == "exp not satisfied" {
			return nil, ErrorTokenExpired
		}
		return nil, fmt.Errorf("validating JWT: %w", err)
	}

	return token, nil
}

// This code can be used to create own JWTs to authenticate against the api. Currently, only ones provided by Keycloak are accepted

// SignToken takes a JWT and signs it with our private key, returning a JWS.
// func (f *Authenticator) SignToken(t jwt.Token) ([]byte, error) {
// 	hdr := jws.NewHeaders()
// 	if err := hdr.Set(jws.AlgorithmKey, jwa.ES256); err != nil {
// 		return nil, fmt.Errorf("setting algorithm: %w", err)
// 	}
// 	if err := hdr.Set(jws.TypeKey, "JWT"); err != nil {
// 		return nil, fmt.Errorf("setting type: %w", err)
// 	}
// 	if err := hdr.Set(jws.KeyIDKey, KeyID); err != nil {
// 		return nil, fmt.Errorf("setting Key ID: %w", err)
// 	}
// 	return jwt.Sign(t, jwa.ES256, f.PrivateKey, jwt.WithHeaders(hdr))
// }

// CreateJWSWithClaims is a helper function to create JWT's with the specified
// claims.
// func (f *Authenticator) CreateJWSWithClaims(claims []string) ([]byte, error) {
// 	t := jwt.New()
// 	err := t.Set(jwt.IssuerKey, FakeIssuer)
// 	if err != nil {
// 		return nil, fmt.Errorf("setting issuer: %w", err)
// 	}
// 	err = t.Set(jwt.AudienceKey, FakeAudience)
// 	if err != nil {
// 		return nil, fmt.Errorf("setting audience: %w", err)
// 	}
// 	err = t.Set(PermissionsClaim, claims)
// 	if err != nil {
// 		return nil, fmt.Errorf("setting permissions: %w", err)
// 	}
// 	return f.SignToken(t)
// }
