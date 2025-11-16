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

func TestStatisticServiceImpl_GetAssignmentsByUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPRRepo := mocks.NewMockPullRequestRepository(ctrl)
	mockTeamRepo := mocks.NewMockTeamRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	statSvc := NewStatisticService(mockPRRepo, mockTeamRepo, mockUserRepo)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		assignments := map[string]int{
			"user1": 5,
			"user2": 3,
		}

		user1 := &models.User{
			UserID:   "user1",
			Username: "User One",
			IsActive: true,
		}
		user2 := &models.User{
			UserID:   "user2",
			Username: "User Two",
			IsActive: true,
		}

		mockPRRepo.EXPECT().GetAssignmentsByUsers(ctx).Return(assignments, nil)
		mockUserRepo.EXPECT().GetUserByID(ctx, "user1").Return(user1, nil)
		mockUserRepo.EXPECT().GetUserByID(ctx, "user2").Return(user2, nil)

		stats, err := statSvc.GetAssignmentsByUsers(ctx)

		require.NoError(t, err)
		assert.Len(t, stats, 2)

		found := make(map[string]*models.UserAssignmentStats)
		for _, stat := range stats {
			found[stat.UserID] = stat
		}

		assert.Equal(t, "User One", found["user1"].Username)
		assert.Equal(t, 5, found["user1"].AssignmentCount)
		assert.Equal(t, "User Two", found["user2"].Username)
		assert.Equal(t, 3, found["user2"].AssignmentCount)
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockPRRepo.EXPECT().GetAssignmentsByUsers(ctx).Return(nil, expectedErr)

		stats, err := statSvc.GetAssignmentsByUsers(ctx)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "failed to get assignments by users")
	})

	t.Run("user lookup error", func(t *testing.T) {
		assignments := map[string]int{"user1": 1}
		expectedErr := errors.New("user lookup error")

		mockPRRepo.EXPECT().GetAssignmentsByUsers(ctx).Return(assignments, nil)
		mockUserRepo.EXPECT().GetUserByID(ctx, "user1").Return(nil, expectedErr)

		stats, err := statSvc.GetAssignmentsByUsers(ctx)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "failed to get user")
	})
}

func TestStatisticServiceImpl_GetPRCountByStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPRRepo := mocks.NewMockPullRequestRepository(ctrl)
	mockTeamRepo := mocks.NewMockTeamRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	statSvc := NewStatisticService(mockPRRepo, mockTeamRepo, mockUserRepo)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		statusCounts := map[string]int{
			"OPEN":   10,
			"MERGED": 25,
		}

		mockPRRepo.EXPECT().GetPRCountByStatus(ctx).Return(statusCounts, nil)

		stats, err := statSvc.GetPRCountByStatus(ctx)

		require.NoError(t, err)
		assert.Len(t, stats, 2)

		found := make(map[string]int)
		for _, stat := range stats {
			found[stat.Status] = stat.Count
		}

		assert.Equal(t, 10, found["OPEN"])
		assert.Equal(t, 25, found["MERGED"])
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockPRRepo.EXPECT().GetPRCountByStatus(ctx).Return(nil, expectedErr)

		stats, err := statSvc.GetPRCountByStatus(ctx)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "failed to get PR count by status")
	})
}

func TestStatisticServiceImpl_GetTeamStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPRRepo := mocks.NewMockPullRequestRepository(ctrl)
	mockTeamRepo := mocks.NewMockTeamRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	statSvc := NewStatisticService(mockPRRepo, mockTeamRepo, mockUserRepo)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		teams := []*models.Team{
			{
				TeamName:  "team1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				TeamName:  "team2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		fullTeam1 := &models.Team{
			TeamName: "team1",
			Members: []models.TeamMember{
				{UserID: "user1", IsActive: true},
				{UserID: "user2", IsActive: false},
			},
		}

		fullTeam2 := &models.Team{
			TeamName: "team2",
			Members: []models.TeamMember{
				{UserID: "user3", IsActive: true},
			},
		}

		mockTeamRepo.EXPECT().GetAllTeams(ctx).Return(teams, nil)
		mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, "team1").Return(fullTeam1, nil)
		mockPRRepo.EXPECT().GetTeamPRCount(ctx, "team1").Return(5, nil)
		mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, "team2").Return(fullTeam2, nil)
		mockPRRepo.EXPECT().GetTeamPRCount(ctx, "team2").Return(3, nil)

		stats, err := statSvc.GetTeamStatistics(ctx)

		require.NoError(t, err)
		assert.Len(t, stats, 2)

		found := make(map[string]*models.TeamStats)
		for _, stat := range stats {
			found[stat.TeamName] = stat
		}

		assert.Equal(t, 2, found["team1"].MemberCount)
		assert.Equal(t, 1, found["team1"].ActiveMemberCount)
		assert.Equal(t, 5, found["team1"].PRCount)

		assert.Equal(t, 1, found["team2"].MemberCount)
		assert.Equal(t, 1, found["team2"].ActiveMemberCount)
		assert.Equal(t, 3, found["team2"].PRCount)
	})

	t.Run("get all teams error", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockTeamRepo.EXPECT().GetAllTeams(ctx).Return(nil, expectedErr)

		stats, err := statSvc.GetTeamStatistics(ctx)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "failed to get all teams")
	})

	t.Run("get team with members error", func(t *testing.T) {
		teams := []*models.Team{{TeamName: "team1"}}
		expectedErr := errors.New("team error")

		mockTeamRepo.EXPECT().GetAllTeams(ctx).Return(teams, nil)
		mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, "team1").Return(nil, expectedErr)

		stats, err := statSvc.GetTeamStatistics(ctx)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "failed to get team with members")
	})
}
