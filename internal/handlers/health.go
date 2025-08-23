package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *HTTPHandler) HealthCheck(ctx *gin.Context) {
	services := make(map[string]string)

	// Check Redis connectivity
	if h.rateLimiter != nil {
		if err := h.rateLimiter.Ping(); err != nil {
			services["redis"] = "unhealthy"
		} else {
			services["redis"] = "healthy"
		}
	} else {
		services["redis"] = "not configured"
	}

	response := healthResponse{
		Status:    "healthy",
		Services:  services,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, response)
}
