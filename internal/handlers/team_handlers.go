package handlers

import (
	"net/http"

	"pr-reviewer-assignment-service/internal/models"

	"github.com/gin-gonic/gin"
)

type CreateTeamRequest struct {
	TeamName string              `json:"team_name" binding:"required"`
	Members  []models.TeamMember `json:"members" binding:"required,dive"`
}

type GetTeamRequest struct {
	TeamName string `json:"team_name" binding:"required"`
}

func (h *Handler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()}})
		return
	}

	team, err := h.teamService.CreateTeamWithMembers(c.Request.Context(), req.TeamName, req.Members)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "TEAM_CREATE_FAILED", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *Handler) GetTeam(c *gin.Context) {
	var req GetTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()}})
		return
	}

	team, err := h.teamService.GetTeamWithMembers(c.Request.Context(), req.TeamName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "TEAM_GET_FAILED", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, team)
}
