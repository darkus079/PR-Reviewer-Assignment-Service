package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"pr-reviewer-assignment-service/internal/models"
	"pr-reviewer-assignment-service/internal/repository"
)

type PullRequestServiceImpl struct {
	prRepo   repository.PullRequestRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
	userSvc  UserService
	randGen  *rand.Rand
}

func NewPullRequestService(
	prRepo repository.PullRequestRepository,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
	userSvc UserService,
) *PullRequestServiceImpl {
	return &PullRequestServiceImpl{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
		userSvc:  userSvc,
		randGen:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *PullRequestServiceImpl) CreatePullRequest(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error) {
	err := s.userSvc.ValidateUserExists(ctx, pr.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("invalid author: %w", err)
	}

	author, err := s.userSvc.GetUserWithTeam(ctx, pr.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get author team: %w", err)
	}

	if author.TeamName == "" {
		return nil, errors.New("author is not assigned to any team")
	}

	activeMembers, err := s.userRepo.GetActiveUsersByTeam(ctx, author.TeamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	var candidates []*models.User
	for _, member := range activeMembers {
		if member.UserID != pr.AuthorID {
			candidates = append(candidates, member)
		}
	}

	selectedReviewers := s.selectRandomReviewers(candidates, 2)
	pr.AssignedReviewers = make([]string, len(selectedReviewers))
	for i, reviewer := range selectedReviewers {
		pr.AssignedReviewers[i] = reviewer.UserID
	}

	err = s.prRepo.CreatePullRequest(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	return pr, nil
}

func (s *PullRequestServiceImpl) MergePullRequest(ctx context.Context, prID string) error {
	exists, err := s.prRepo.PullRequestExists(ctx, prID)
	if err != nil {
		return fmt.Errorf("failed to check PR existence: %w", err)
	}
	if !exists {
		return errors.New("pull request not found")
	}

	err = s.prRepo.MergePullRequest(ctx, prID)
	if err != nil {
		return fmt.Errorf("failed to merge pull request: %w", err)
	}

	return nil
}

func (s *PullRequestServiceImpl) ReassignReviewer(ctx context.Context, prID string, oldReviewerID string) error {
	pr, err := s.prRepo.GetPullRequestByID(ctx, prID)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}
	if pr == nil {
		return errors.New("pull request not found")
	}

	if pr.Status == models.PRStatusMerged {
		return errors.New("cannot reassign reviewer for merged pull request")
	}

	isAssigned := false
	for _, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldReviewerID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return errors.New("reviewer is not assigned to this pull request")
	}

	oldReviewer, err := s.userSvc.GetUserWithTeam(ctx, oldReviewerID)
	if err != nil {
		return fmt.Errorf("failed to get reviewer team: %w", err)
	}

	if oldReviewer.TeamName == "" {
		return errors.New("reviewer is not assigned to any team")
	}

	activeMembers, err := s.userRepo.GetActiveUsersByTeam(ctx, oldReviewer.TeamName)
	if err != nil {
		return fmt.Errorf("failed to get team members: %w", err)
	}

	var candidates []*models.User
	for _, member := range activeMembers {
		if member.UserID != pr.AuthorID && member.UserID != oldReviewerID {
			candidates = append(candidates, member)
		}
	}

	if len(candidates) == 0 {
		return errors.New("no candidate reviewers available")
	}

	newReviewer := candidates[s.randGen.Intn(len(candidates))]

	for i, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldReviewerID {
			pr.AssignedReviewers[i] = newReviewer.UserID
			break
		}
	}

	err = s.prRepo.SetAssignedReviewers(ctx, prID, pr.AssignedReviewers)
	if err != nil {
		return fmt.Errorf("failed to update reviewers: %w", err)
	}

	return nil
}

func (s *PullRequestServiceImpl) selectRandomReviewers(candidates []*models.User, count int) []*models.User {
	if len(candidates) <= count {
		return candidates
	}

	available := make([]*models.User, len(candidates))
	copy(available, candidates)

	selected := make([]*models.User, 0, count)

	for i := 0; i < count && len(available) > 0; i++ {
		randomIndex := s.randGen.Intn(len(available))

		selected = append(selected, available[randomIndex])

		available = append(available[:randomIndex], available[randomIndex+1:]...)
	}

	return selected
}
