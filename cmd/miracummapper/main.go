package main

import (
	"log"
	"miracummapper/internal/api"
	"miracummapper/internal/config"
	"miracummapper/internal/database"
	"miracummapper/internal/server/middlewares"

	"miracummapper/internal/server"
)

func main() {
	log.Printf("Preparing Miracum Mapper Backend...")

	config := config.NewConfig()

	db := database.NewGormConnection(config)

	keySet, err := middlewares.FetchKeycloakCerts(config)
	if err != nil {
		log.Fatalf("Failed to fetch Keycloak certs: %v", err)
	}

	r := server.CreateGin(config)

	auth, err := middlewares.NewAuthenticator(keySet, config) // keySet
	if err != nil {
		log.Fatalln("error creating authenticator:", err)
	}

	// Create middleware for validating tokens.
	mw, err := server.CreateMiddleware(auth, config)
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}
	r.Use(mw...)

	svr := server.CreateServer(db, config, auth)

	strictHandler := api.NewStrictHandler(svr, nil)

	api.RegisterHandlers(r, strictHandler)

	r.Run(":" + config.Env.Port)
}
