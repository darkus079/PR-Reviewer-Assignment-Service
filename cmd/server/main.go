package main

import (
	"log"

	"pr-reviewer-assignment-service/internal/config"
	"pr-reviewer-assignment-service/internal/database"
	"pr-reviewer-assignment-service/internal/handlers"
	"pr-reviewer-assignment-service/internal/middleware"
	"pr-reviewer-assignment-service/internal/repository"
	"pr-reviewer-assignment-service/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewPostgresUserRepository(db.Pool)
	teamRepo := repository.NewPostgresTeamRepository(db.Pool)
	prRepo := repository.NewPostgresPullRequestRepository(db.Pool)

	userSvc := services.NewUserService(userRepo)
	teamSvc := services.NewTeamService(teamRepo, userRepo)
	prSvc := services.NewPullRequestService(prRepo, userRepo, teamRepo, userSvc)
	statSvc := services.NewStatisticService(prRepo, teamRepo, userRepo)

	handler := handlers.NewHandler(teamSvc, userSvc, prSvc, statSvc)
	healthHandler := handlers.NewHealthHandler(userRepo)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware())

	r.GET("/health", healthHandler.Health)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware("admin-token", "user-token"))
	{
		team := api.Group("/team")
		{
			team.POST("/add", handler.CreateTeam)
			team.GET("/get", handler.GetTeam)
		}

		user := api.Group("/users")
		{
			user.POST("/setIsActive", middleware.AdminOnlyMiddleware(), handler.SetUserActive)
			user.GET("/getReview", handler.GetUserReviews)
		}

		pr := api.Group("/pullRequest")
		{
			pr.POST("/create", handler.CreatePullRequest)
			pr.POST("/merge", handler.MergePullRequest)
			pr.POST("/reassign", handler.ReassignReviewer)
		}
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(r.Run(":" + cfg.Server.Port))
}
