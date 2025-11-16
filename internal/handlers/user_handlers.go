package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetUserActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type GetUserReviewRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

func (h *Handler) SetUserActive(c *gin.Context) {
	var req SetUserActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()}})
		return
	}

	err := h.userService.SetUserActiveStatus(c.Request.Context(), req.UserID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "USER_UPDATE_FAILED", "message": err.Error()}})
		return
	}

	user, err := h.userService.GetUserWithTeam(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "USER_GET_FAILED", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserReviews(c *gin.Context) {
	var req GetUserReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()}})
		return
	}

	err := h.userService.ValidateUserExists(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "USER_NOT_FOUND", "message": err.Error()}})
		return
	}

	prs, err := h.prService.GetUserPullRequests(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "PR_GET_FAILED", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pull_requests": prs})
}
