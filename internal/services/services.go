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
