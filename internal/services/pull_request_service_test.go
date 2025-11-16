package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pr-reviewer-assignment-service/internal/mocks"
	"pr-reviewer-assignment-service/internal/models"
)

func TestPullRequestServiceImpl_selectRandomReviewers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPRRepo := mocks.NewMockPullRequestRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTeamRepo := mocks.NewMockTeamRepository(ctrl)
	mockUserSvc := mocks.NewMockUserRepository(ctrl)

	prSvc := &PullRequestServiceImpl{
		prRepo:   mockPRRepo,
		userRepo: mockUserRepo,
		teamRepo: mockTeamRepo,
		userSvc:  &UserServiceImpl{userRepo: mockUserSvc},
		randGen:  NewPullRequestService(mockPRRepo, mockUserRepo, mockTeamRepo, &UserServiceImpl{userRepo: mockUserSvc}).randGen,
	}

	candidates := []*models.User{
		{UserID: "user1", Username: "User 1"},
		{UserID: "user2", Username: "User 2"},
		{UserID: "user3", Username: "User 3"},
		{UserID: "user4", Username: "User 4"},
	}

	t.Run("select fewer than available", func(t *testing.T) {
		selected := prSvc.selectRandomReviewers(candidates, 2)

		assert.Len(t, selected, 2)
		userIDs := make(map[string]bool)
		for _, user := range selected {
			assert.False(t, userIDs[user.UserID], "User %s selected multiple times", user.UserID)
			userIDs[user.UserID] = true
		}
	})

	t.Run("select all available", func(t *testing.T) {
		selected := prSvc.selectRandomReviewers(candidates, 4)

		assert.Len(t, selected, 4)
		userIDs := make(map[string]bool)
		for _, user := range selected {
			userIDs[user.UserID] = true
		}
		assert.Len(t, userIDs, 4)
	})

	t.Run("select more than available", func(t *testing.T) {
		selected := prSvc.selectRandomReviewers(candidates, 10)

		assert.Len(t, selected, 4)
	})

	t.Run("empty candidates", func(t *testing.T) {
		selected := prSvc.selectRandomReviewers([]*models.User{}, 2)

		assert.Len(t, selected, 0)
	})

	t.Run("select zero", func(t *testing.T) {
		selected := prSvc.selectRandomReviewers(candidates, 0)

		assert.Len(t, selected, 0)
	})
}

func TestPullRequestServiceImpl_GetUserPullRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPRRepo := mocks.NewMockPullRequestRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTeamRepo := mocks.NewMockTeamRepository(ctrl)
	mockUserSvc := mocks.NewMockUserRepository(ctrl)

	prSvc := NewPullRequestService(mockPRRepo, mockUserRepo, mockTeamRepo, &UserServiceImpl{userRepo: mockUserSvc})

	ctx := context.Background()
	userID := "test-user"

	t.Run("success", func(t *testing.T) {
		expectedPRs := []*models.PullRequestShort{
			{PullRequestID: "pr1", PullRequestName: "PR 1", AuthorID: "author1", Status: "OPEN"},
			{PullRequestID: "pr2", PullRequestName: "PR 2", AuthorID: "author2", Status: "MERGED"},
		}

		mockUserSvc.EXPECT().UserExists(ctx, userID).Return(true, nil)
		mockPRRepo.EXPECT().GetPullRequestsByReviewer(ctx, userID).Return(expectedPRs, nil)

		prs, err := prSvc.GetUserPullRequests(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, expectedPRs, prs)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserSvc.EXPECT().UserExists(ctx, userID).Return(false, nil)

		prs, err := prSvc.GetUserPullRequests(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, prs)
		assert.Contains(t, err.Error(), "invalid user")
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockUserSvc.EXPECT().UserExists(ctx, userID).Return(true, nil)
		mockPRRepo.EXPECT().GetPullRequestsByReviewer(ctx, userID).Return(nil, expectedErr)

		prs, err := prSvc.GetUserPullRequests(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, prs)
		assert.Contains(t, err.Error(), "failed to get pull requests for user")
	})
}
