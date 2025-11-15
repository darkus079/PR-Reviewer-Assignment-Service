package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pr-reviewer-assignment-service/internal/models"
	"pr-reviewer-assignment-service/internal/repository"
)

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (s *UserServiceImpl) SetUserActiveStatus(ctx context.Context, userID string, isActive bool) error {
	exists, err := s.userRepo.UserExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return errors.New("user not found")
	}

	err = s.userRepo.SetUserActiveStatus(ctx, userID, isActive)
	if err != nil {
		return fmt.Errorf("failed to update user active status: %w", err)
	}

	return nil
}

func (s *UserServiceImpl) ValidateUserExists(ctx context.Context, userID string) error {
	exists, err := s.userRepo.UserExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return errors.New("user not found")
	}
	return nil
}

func (s *UserServiceImpl) GetUserWithTeam(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
