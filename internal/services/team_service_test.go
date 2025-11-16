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

func TestTeamServiceImpl_CreateTeamWithMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTeamRepo := mocks.NewMockTeamRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	teamSvc := NewTeamService(mockTeamRepo, mockUserRepo)

	ctx := context.Background()
	teamName := "test-team"
	members := []models.TeamMember{
		{UserID: "user1", Username: "User One", IsActive: true},
		{UserID: "user2", Username: "User Two", IsActive: false},
	}

	t.Run("success", func(t *testing.T) {
		expectedTeam := &models.Team{
			TeamName: teamName,
			Members:  members,
		}

		mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, nil)
		mockTeamRepo.EXPECT().CreateTeam(ctx, gomock.Any()).Return(nil)
		mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil).Times(2)

		team, err := teamSvc.CreateTeamWithMembers(ctx, teamName, members)

		require.NoError(t, err)
		assert.Equal(t, expectedTeam.TeamName, team.TeamName)
		assert.Equal(t, expectedTeam.Members, team.Members)
	})

	t.Run("team already exists", func(t *testing.T) {
		mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(true, nil)

		team, err := teamSvc.CreateTeamWithMembers(ctx, teamName, members)

		assert.Error(t, err)
		assert.Nil(t, team)
		assert.Contains(t, err.Error(), "team already exists")
	})

	t.Run("team exists check error", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, expectedErr)

		team, err := teamSvc.CreateTeamWithMembers(ctx, teamName, members)

		assert.Error(t, err)
		assert.Nil(t, team)
		assert.Contains(t, err.Error(), "failed to check team existence")
	})

	t.Run("create team error", func(t *testing.T) {
		expectedErr := errors.New("create error")
		mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, nil)
		mockTeamRepo.EXPECT().CreateTeam(ctx, gomock.Any()).Return(expectedErr)

		team, err := teamSvc.CreateTeamWithMembers(ctx, teamName, members)

		assert.Error(t, err)
		assert.Nil(t, team)
		assert.Contains(t, err.Error(), "failed to create team")
	})

	t.Run("create user error", func(t *testing.T) {
		expectedErr := errors.New("user create error")
		mockTeamRepo.EXPECT().TeamExists(ctx, teamName).Return(false, nil)
		mockTeamRepo.EXPECT().CreateTeam(ctx, gomock.Any()).Return(nil)
		mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(expectedErr)

		team, err := teamSvc.CreateTeamWithMembers(ctx, teamName, members)

		assert.Error(t, err)
		assert.Nil(t, team)
		assert.Contains(t, err.Error(), "failed to create/update user")
	})
}

func TestTeamServiceImpl_GetTeamWithMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTeamRepo := mocks.NewMockTeamRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	teamSvc := NewTeamService(mockTeamRepo, mockUserRepo)

	ctx := context.Background()
	teamName := "test-team"

	t.Run("success", func(t *testing.T) {
		expectedTeam := &models.Team{
			TeamName:  teamName,
			Members:   []models.TeamMember{{UserID: "user1", Username: "User One", IsActive: true}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(expectedTeam, nil)

		team, err := teamSvc.GetTeamWithMembers(ctx, teamName)

		require.NoError(t, err)
		assert.Equal(t, expectedTeam, team)
	})

	t.Run("team not found", func(t *testing.T) {
		mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(nil, nil)

		team, err := teamSvc.GetTeamWithMembers(ctx, teamName)

		assert.Error(t, err)
		assert.Nil(t, team)
		assert.Contains(t, err.Error(), "team not found")
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockTeamRepo.EXPECT().GetTeamWithMembers(ctx, teamName).Return(nil, expectedErr)

		team, err := teamSvc.GetTeamWithMembers(ctx, teamName)

		assert.Error(t, err)
		assert.Nil(t, team)
		assert.Contains(t, err.Error(), "failed to get team with members")
	})
}
