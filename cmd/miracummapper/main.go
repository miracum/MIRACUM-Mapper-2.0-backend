package main

import (
	"log"
	"miracummapper/internal/api"
	"miracummapper/internal/config"
	"miracummapper/internal/database"
	"miracummapper/internal/server/middlewares"

	// "miracummapper/internal/routes"
	"miracummapper/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Printf("Preparing Miracum Mapper Backend...")

	config := config.NewConfig()

	// Db := database.NewDBConnection(config)

	// database.Migrate(db)
	db := database.NewGormConnection(config)

	// r := routes.SetupRouter()

	// // run on port from config
	// r.Run(":" + config.Env.Port)

	r := gin.Default()

	fa, err := middlewares.NewFakeAuthenticator()
	if err != nil {
		log.Fatalln("error creating authenticator:", err)
	}

	// Create middleware for validating tokens.
	// mw, err := server.CreateStrictMiddleware(fa)
	mw, err := server.CreateMiddleware(fa)
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}
	r.Use(mw...)
	// r.Use(middleware.Logger())

	// svr := server.CreateServer(db, config)
	// svr := server.CreateStrictServer(Db, config)
	svr := server.CreateServer(db, config)

	strictHandler := api.NewStrictHandler(svr, nil)

	api.RegisterHandlers(r, strictHandler)

	// r.Use(mw...)

	// api.RegisterHandlers(r, svr)

	normalJWS, err := fa.CreateJWSWithClaims([]string{"normal"})
	if err != nil {
		log.Fatalln("error creating normal JWS:", err)
	}

	adminJWS, err := fa.CreateJWSWithClaims([]string{"admin"})
	if err != nil {
		log.Fatalln("error creating admin JWS:", err)
	}

	log.Println("Normal token", string(normalJWS))
	log.Println("Admin token", string(adminJWS))

	r.Run(":" + config.Env.Port)
}
