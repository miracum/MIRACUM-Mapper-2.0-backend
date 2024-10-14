package middlewares

import (
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// PrivateKey is an ECDSA private key which was generated with the following
// command:
//
//	openssl ecparam -name prime256v1 -genkey -noout -out ecprivatekey.pem
//
// We are using a hard coded key here in this example, but in real applications,
// you would never do this. Your JWT signing key must never be in your application,
// only the public key.
// const PrivateKey = `-----BEGIN EC PRIVATE KEY-----
// MHcCAQEEIN2dALnjdcZaIZg4QuA6Dw+kxiSW502kJfmBN3priIhPoAoGCCqGSM49
// AwEHoUQDQgAE4pPyvrB9ghqkT1Llk0A42lixkugFd/TBdOp6wf69O9Nndnp4+HcR
// s9SlG/8hjB2Hz42v4p3haKWv3uS1C6ahCQ==
// -----END EC PRIVATE KEY-----`

const HardCodedKeyID = `xah9Ht7EMFI0WfaLRIdJsVLLH2BzRdHT2qzowq8PkH4`
const HardCodedIssuer = "http://localhost:8081/realms/master"
const HardCodedAudience = "account"
const PermissionsClaim = "perm"

type Authenticator struct {
	// PrivateKey *ecdsa.PrivateKey
	KeySet jwk.Set
}

var _ JWSValidator = (*Authenticator)(nil)

// NewAuthenticator creates an authenticator example which uses a hard coded
// ECDSA key to validate JWT's that it has signed itself.
func NewFakeAuthenticator(keySet jwk.Set) (*Authenticator, error) {
	// privKey, err := ecdsafile.LoadEcdsaPrivateKey([]byte(PrivateKey))
	// if err != nil {
	// 	return nil, fmt.Errorf("loading PEM private key: %w", err)
	// }

	// set := jwk.NewSet()
	// pubKey := jwk.NewECDSAPublicKey()

	// err = pubKey.FromRaw(&privKey.PublicKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("parsing jwk key: %w", err)
	// }

	// err = pubKey.Set(jwk.AlgorithmKey, jwa.ES256)
	// if err != nil {
	// 	return nil, fmt.Errorf("setting key algorithm: %w", err)
	// }

	// err = pubKey.Set(jwk.KeyIDKey, KeyID)
	// if err != nil {
	// 	return nil, fmt.Errorf("setting key ID: %w", err)
	// }

	// set.Add(pubKey)

	// return &FakeAuthenticator{PrivateKey: privKey, KeySet: set}, nil
	return &Authenticator{KeySet: keySet}, nil
}

// ValidateJWS ensures that the critical JWT claims needed to ensure that we
// trust the JWT are present and with the correct values.
func (f *Authenticator) ValidateJWS(jwsString string) (jwt.Token, error) {
	return jwt.Parse([]byte(jwsString), jwt.WithKeySet(f.KeySet),
		jwt.WithAudience(HardCodedAudience), jwt.WithIssuer(HardCodedIssuer))
}

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
