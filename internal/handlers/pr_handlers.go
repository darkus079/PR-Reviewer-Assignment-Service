package handlers

import (
	"net/http"

	"pr-reviewer-assignment-service/internal/models"

	"github.com/gin-gonic/gin"
)

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

type ReassignPRRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldReviewerID string `json:"old_reviewer_id" binding:"required"`
}

func (h *Handler) CreatePullRequest(c *gin.Context) {
	var req CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()}})
		return
	}

	pr := &models.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
		Status:          models.PRStatusOpen,
	}

	createdPR, err := h.prService.CreatePullRequest(c.Request.Context(), pr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "PR_CREATE_FAILED", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, createdPR)
}

func (h *Handler) MergePullRequest(c *gin.Context) {
	var req MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()}})
		return
	}

	err := h.prService.MergePullRequest(c.Request.Context(), req.PullRequestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "PR_MERGE_FAILED", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pull request merged successfully"})
}

func (h *Handler) ReassignReviewer(c *gin.Context) {
	var req ReassignPRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()}})
		return
	}

	err := h.prService.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldReviewerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "PR_REASSIGN_FAILED", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reviewer reassigned successfully"})
}
