package services

import (
	"context"
	"fmt"

	"pr-reviewer-assignment-service/internal/models"
	"pr-reviewer-assignment-service/internal/repository"
)

type StatisticServiceImpl struct {
	prRepo   repository.PullRequestRepository
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
}

func NewStatisticService(
	prRepo repository.PullRequestRepository,
	teamRepo repository.TeamRepository,
	userRepo repository.UserRepository,
) *StatisticServiceImpl {
	return &StatisticServiceImpl{
		prRepo:   prRepo,
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *StatisticServiceImpl) GetAssignmentsByUsers(ctx context.Context) ([]*models.UserAssignmentStats, error) {
	assignments, err := s.prRepo.GetAssignmentsByUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments by users: %w", err)
	}

	var stats []*models.UserAssignmentStats
	for userID, count := range assignments {
		user, err := s.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user %s: %w", userID, err)
		}
		if user != nil {
			stats = append(stats, &models.UserAssignmentStats{
				UserID:          userID,
				Username:        user.Username,
				AssignmentCount: count,
			})
		}
	}

	return stats, nil
}

func (s *StatisticServiceImpl) GetPRCountByStatus(ctx context.Context) ([]*models.PRStatusStats, error) {
	statusCounts, err := s.prRepo.GetPRCountByStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR count by status: %w", err)
	}

	var stats []*models.PRStatusStats
	for status, count := range statusCounts {
		stats = append(stats, &models.PRStatusStats{
			Status: status,
			Count:  count,
		})
	}

	return stats, nil
}

func (s *StatisticServiceImpl) GetTeamStatistics(ctx context.Context) ([]*models.TeamStats, error) {
	teams, err := s.teamRepo.GetAllTeams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all teams: %w", err)
	}

	var stats []*models.TeamStats
	for _, team := range teams {
		fullTeam, err := s.teamRepo.GetTeamWithMembers(ctx, team.TeamName)
		if err != nil {
			return nil, fmt.Errorf("failed to get team with members %s: %w", team.TeamName, err)
		}

		activeCount := 0
		for _, member := range fullTeam.Members {
			if member.IsActive {
				activeCount++
			}
		}

		prCount, err := s.prRepo.GetTeamPRCount(ctx, team.TeamName)
		if err != nil {
			return nil, fmt.Errorf("failed to get PR count for team %s: %w", team.TeamName, err)
		}

		stats = append(stats, &models.TeamStats{
			TeamName:          team.TeamName,
			MemberCount:       len(fullTeam.Members),
			ActiveMemberCount: activeCount,
			PRCount:           prCount,
		})
	}

	return stats, nil
}
