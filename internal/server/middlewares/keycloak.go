package middlewares

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"miracummapper/internal/config"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
)

func FetchKeycloakCerts(config *config.Config) (jwk.Set, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", config.Env.KeycloakHost, config.Env.KeycloakRealm)
	var keys struct {
		Keys []jwk.Key `json:"keys"`
	}

	for i := 0; i < config.File.DatabaseConfig.Retry; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("reading response body: %w", err)
			}
			err = json.Unmarshal(body, &keys)
			if err != nil {
				return nil, fmt.Errorf("unmarshalling response: %w", err)
			}
			keySet := jwk.NewSet()
			for _, key := range keys.Keys {
				keySet.Add(key)
			}
			return keySet, nil
		}
		log.Printf("Failed to fetch Keycloak certs. Retrying in %d seconds", config.File.DatabaseConfig.Sleep)
		if i != config.File.DatabaseConfig.Retry-1 {
			time.Sleep(time.Duration(config.File.DatabaseConfig.Sleep) * time.Second)
		}
	}
	return nil, fmt.Errorf("failed to fetch Keycloak certs after %d retries", config.File.DatabaseConfig.Retry)
}

// func NewFakeAuthenticator(config *config.Config) (*FakeAuthenticator, error) {
// 	keys, err := fetchKeycloakCerts(config)
// 	if err != nil {
// 		return nil, fmt.Errorf("fetching Keycloak certs: %w", err)
// 	}

// 	set := jwk.NewSet()
// 	for _, key := range keys.Keys {
// 		set.Add(key)
// 	}

// 	return &FakeAuthenticator{KeySet: set}, nil
// }

// type FakeAuthenticator struct {
// 	KeySet jwk.Set
// }

// var _ JWSValidator = (*FakeAuthenticator)(nil)

// func main() {
// 	// Example usage
// 	config := &config.Config{
// 		KeycloakHost:  "localhost:8081",
// 		KeycloakRealm: "master",
// 		File: config.FileConfig{
// 			DatabaseConfig: config.DatabaseConfig{
// 				Retry: 3,
// 				Sleep: 2,
// 			},
// 		},
// 	}

// 	authenticator, err := NewFakeAuthenticator(config)
// 	if err != nil {
// 		log.Fatalf("Failed to create authenticator: %v", err)
// 	}

// 	fmt.Printf("Authenticator created with keys: %v\n", authenticator.KeySet)
// }
