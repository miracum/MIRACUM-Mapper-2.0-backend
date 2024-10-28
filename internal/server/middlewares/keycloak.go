package middlewares

import (
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

	for i := 0; i < config.File.DatabaseConfig.Retry; i++ {
		log.Printf("Fetching Keycloak certs from %s", url)
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("reading response body: %w", err)
			}

			set, err := jwk.Parse(body)
			if err != nil {
				fmt.Println("Error parsing JWK Set:", err)
				return nil, err
			}
			return set, nil
			// return keySet, nil
		}
		log.Printf("Failed to fetch Keycloak certs. Retrying in %d seconds", config.File.DatabaseConfig.Sleep)
		log.Printf("Error: %v", err)
		if i != config.File.DatabaseConfig.Retry-1 {
			time.Sleep(time.Duration(config.File.DatabaseConfig.Sleep) * time.Second)
		}
	}
	return nil, fmt.Errorf("failed to fetch Keycloak certs after %d retries", config.File.DatabaseConfig.Retry)
}