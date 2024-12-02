package server

import (
	"miracummapper/internal/config"

	"github.com/gin-gonic/gin"
)

func CreateGin(config *config.Config) *gin.Engine {
	if !config.File.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()

	engine.Use(gin.Recovery())

	if config.File.Debug {
		engine.Use(gin.Logger())
	}

	return engine
}
