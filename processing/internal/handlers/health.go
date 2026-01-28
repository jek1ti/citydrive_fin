package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jekiti/citydrive/processing/internal/domain"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Liveness(c *gin.Context) {
	var response domain.HealthResponse
	response.Status = "alive"
	response.Timestamp = time.Now().Unix()
	c.JSON(200, response)
}

func (h *HealthHandler) Readiness(c *gin.Context) {
	var response domain.HealthResponse
	response.Status = "ready"
	response.Timestamp = time.Now().Unix()
	c.JSON(200, response)
}
