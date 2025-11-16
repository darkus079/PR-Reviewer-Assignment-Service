package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pr-reviewer-assignment-service/internal/mocks"
	"pr-reviewer-assignment-service/internal/models"
)

func TestUserServiceImpl_SetUserActiveStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := NewUserService(mockRepo)

	ctx := context.Background()
	userID := "test-user"

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().UserExists(ctx, userID).Return(true, nil)
		mockRepo.EXPECT().SetUserActiveStatus(ctx, userID, true).Return(nil)

		err := userSvc.SetUserActiveStatus(ctx, userID, true)

		assert.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().UserExists(ctx, userID).Return(false, nil)

		err := userSvc.SetUserActiveStatus(ctx, userID, true)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("repository error on exists check", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockRepo.EXPECT().UserExists(ctx, userID).Return(false, expectedErr)

		err := userSvc.SetUserActiveStatus(ctx, userID, true)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check user existence")
	})

	t.Run("repository error on update", func(t *testing.T) {
		expectedErr := errors.New("update error")
		mockRepo.EXPECT().UserExists(ctx, userID).Return(true, nil)
		mockRepo.EXPECT().SetUserActiveStatus(ctx, userID, true).Return(expectedErr)

		err := userSvc.SetUserActiveStatus(ctx, userID, true)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update user active status")
	})
}

func TestUserServiceImpl_ValidateUserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := NewUserService(mockRepo)

	ctx := context.Background()
	userID := "test-user"

	t.Run("user exists", func(t *testing.T) {
		mockRepo.EXPECT().UserExists(ctx, userID).Return(true, nil)

		err := userSvc.ValidateUserExists(ctx, userID)

		assert.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().UserExists(ctx, userID).Return(false, nil)

		err := userSvc.ValidateUserExists(ctx, userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockRepo.EXPECT().UserExists(ctx, userID).Return(false, expectedErr)

		err := userSvc.ValidateUserExists(ctx, userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check user existence")
	})
}

func TestUserServiceImpl_GetUserWithTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := NewUserService(mockRepo)

	ctx := context.Background()
	userID := "test-user"

	t.Run("success", func(t *testing.T) {
		expectedUser := &models.User{
			UserID:    userID,
			Username:  "testuser",
			TeamName:  "test-team",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(expectedUser, nil)

		user, err := userSvc.GetUserWithTeam(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, nil)

		user, err := userSvc.GetUserWithTeam(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, expectedErr)

		user, err := userSvc.GetUserWithTeam(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to get user")
	})
}
