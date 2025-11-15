package services

import (
	"context"
	"errors"
	"fmt"

	"pr-reviewer-assignment-service/internal/models"
	"pr-reviewer-assignment-service/internal/repository"
)

type TeamServiceImpl struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository) *TeamServiceImpl {
	return &TeamServiceImpl{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *TeamServiceImpl) CreateTeamWithMembers(ctx context.Context, teamName string, members []models.TeamMember) (*models.Team, error) {
	exists, err := s.teamRepo.TeamExists(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}
	if exists {
		return nil, errors.New("team already exists")
	}

	team := &models.Team{
		TeamName: teamName,
		Members:  members,
	}

	err = s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	for _, member := range members {
		user := &models.User{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: teamName,
			IsActive: member.IsActive,
		}

		err = s.userRepo.CreateUser(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create/update user %s: %w", member.UserID, err)
		}
	}

	return team, nil
}

func (s *TeamServiceImpl) GetTeamWithMembers(ctx context.Context, teamName string) (*models.Team, error) {
	team, err := s.teamRepo.GetTeamWithMembers(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team with members: %w", err)
	}
	if team == nil {
		return nil, errors.New("team not found")
	}

	return team, nil
}
