package handlers

import (
	"pr-reviewer-assignment-service/internal/services"
)

type Handler struct {
	teamService      services.TeamService
	userService      services.UserService
	prService        services.PullRequestService
	statisticService services.StatisticService
}

func NewHandler(
	teamService services.TeamService,
	userService services.UserService,
	prService services.PullRequestService,
	statisticService services.StatisticService,
) *Handler {
	return &Handler{
		teamService:      teamService,
		userService:      userService,
		prService:        prService,
		statisticService: statisticService,
	}
}
