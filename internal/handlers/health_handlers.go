package handlers

import (
	"context"
	"net/http"

	"pr-reviewer-assignment-service/internal/repository"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	userRepo repository.UserRepository
}

func NewHealthHandler(userRepo repository.UserRepository) *HealthHandler {
	return &HealthHandler{
		userRepo: userRepo,
	}
}

func (h *HealthHandler) Health(c *gin.Context) {
	_, err := h.userRepo.UserExists(context.Background(), "health-check")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"database": gin.H{
				"status":  "error",
				"message": "Database connection check failed",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"database": gin.H{
			"status": "ok",
		},
	})
}
