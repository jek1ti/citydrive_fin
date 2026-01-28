package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jekiti/citydrive/processing/internal/handlers"
)

func NewServer() *gin.Engine {
	router := gin.Default()
	healthHandler := handlers.NewHealthHandler()

	router.GET("/health/liveness", healthHandler.Liveness)
	router.GET("/health/readiness", healthHandler.Readiness)
	
	return router
}
