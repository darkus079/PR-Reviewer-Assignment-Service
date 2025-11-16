package services

import (
	"context"

	"pr-reviewer-assignment-service/internal/models"
)

type UserService interface {
	SetUserActiveStatus(ctx context.Context, userID string, isActive bool) error
	ValidateUserExists(ctx context.Context, userID string) error
	GetUserWithTeam(ctx context.Context, userID string) (*models.User, error)
}

type TeamService interface {
	CreateTeamWithMembers(ctx context.Context, teamName string, members []models.TeamMember) (*models.Team, error)
	GetTeamWithMembers(ctx context.Context, teamName string) (*models.Team, error)
}

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string) error
	ReassignReviewer(ctx context.Context, prID string, oldReviewerID string) error
	GetUserPullRequests(ctx context.Context, userID string) ([]*models.PullRequestShort, error)
}

type StatisticService interface {
	GetAssignmentsByUsers(ctx context.Context) ([]*models.UserAssignmentStats, error)
	GetPRCountByStatus(ctx context.Context) ([]*models.PRStatusStats, error)
	GetTeamStatistics(ctx context.Context) ([]*models.TeamStats, error)
}
