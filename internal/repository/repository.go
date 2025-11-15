package repository

import (
	"context"

	"pr-reviewer-assignment-service/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID string) error

	GetUsersByTeam(ctx context.Context, teamName string) ([]*models.User, error)
	GetActiveUsersByTeam(ctx context.Context, teamName string) ([]*models.User, error)
	SetUserActiveStatus(ctx context.Context, userID string, isActive bool) error
	UserExists(ctx context.Context, userID string) (bool, error)
}

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamByName(ctx context.Context, teamName string) (*models.Team, error)
	UpdateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, teamName string) error

	TeamExists(ctx context.Context, teamName string) (bool, error)
	GetTeamWithMembers(ctx context.Context, teamName string) (*models.Team, error)

	// Методы для статистики
	GetAllTeams(ctx context.Context) ([]*models.Team, error)
}

type PullRequestRepository interface {
	CreatePullRequest(ctx context.Context, pr *models.PullRequest) error
	GetPullRequestByID(ctx context.Context, prID string) (*models.PullRequest, error)
	UpdatePullRequest(ctx context.Context, pr *models.PullRequest) error
	DeletePullRequest(ctx context.Context, prID string) error

	GetPullRequestsByReviewer(ctx context.Context, userID string) ([]*models.PullRequestShort, error)
	MergePullRequest(ctx context.Context, prID string) error
	PullRequestExists(ctx context.Context, prID string) (bool, error)
	GetAssignedReviewers(ctx context.Context, prID string) ([]string, error)
	SetAssignedReviewers(ctx context.Context, prID string, reviewers []string) error

	// Методы для статистики
	GetPRCountByStatus(ctx context.Context) (map[string]int, error)
	GetAssignmentsByUsers(ctx context.Context) (map[string]int, error)
	GetTeamPRCount(ctx context.Context, teamName string) (int, error)
}
