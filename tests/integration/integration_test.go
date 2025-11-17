package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"pr-reviewer-assignment-service/internal/models"
	"pr-reviewer-assignment-service/internal/repository"
	"pr-reviewer-assignment-service/internal/services"
)

var (
	dbContainer testcontainers.Container
	dbPool      *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5)),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	dbContainer = pgContainer

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to create connection pool: %v", err)
	}

	dbPool = pool

	if err := runMigrations(connStr); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	code := m.Run()

	if err := pgContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	dbPool.Close()

	os.Exit(code)
}

func runMigrations(connStr string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	defer db.Close()

	config := migratepg.Config{}
	driver, err := migratepg.WithInstance(db, &config)
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func setupTestData(t *testing.T) {
	ctx := context.Background()

	_, err := dbPool.Exec(ctx, "DELETE FROM pr_reviewers")
	require.NoError(t, err)

	_, err = dbPool.Exec(ctx, "DELETE FROM pull_requests")
	require.NoError(t, err)

	_, err = dbPool.Exec(ctx, "DELETE FROM users")
	require.NoError(t, err)

	_, err = dbPool.Exec(ctx, "DELETE FROM teams")
	require.NoError(t, err)
}

func TestIntegration_FullWorkflow(t *testing.T) {
	ctx := context.Background()
	setupTestData(t)

	userRepo := repository.NewPostgresUserRepository(dbPool)
	teamRepo := repository.NewPostgresTeamRepository(dbPool)
	prRepo := repository.NewPostgresPullRequestRepository(dbPool)

	userSvc := services.NewUserService(userRepo)
	teamSvc := services.NewTeamService(teamRepo, userRepo)
	prSvc := services.NewPullRequestService(prRepo, userRepo, teamRepo, userSvc)

	teamMembers := []models.TeamMember{
		{UserID: "user1", Username: "User One", IsActive: true},
		{UserID: "user2", Username: "User Two", IsActive: true},
		{UserID: "user3", Username: "User Three", IsActive: true},
	}

	team, err := teamSvc.CreateTeamWithMembers(ctx, "test-team", teamMembers)
	require.NoError(t, err)
	assert.Equal(t, "test-team", team.TeamName)
	assert.Len(t, team.Members, 3)

	pr := &models.PullRequest{
		PullRequestID:   "pr-001",
		PullRequestName: "Test PR",
		AuthorID:        "user1",
		Status:          models.PRStatusOpen,
	}

	createdPR, err := prSvc.CreatePullRequest(ctx, pr)
	require.NoError(t, err)
	assert.Equal(t, "pr-001", createdPR.PullRequestID)
	assert.Len(t, createdPR.AssignedReviewers, 2)

	reviewers, err := prRepo.GetAssignedReviewers(ctx, "pr-001")
	require.NoError(t, err)
	assert.Len(t, reviewers, 2)
	assert.NotContains(t, reviewers, "user1")

	err = prSvc.MergePullRequest(ctx, "pr-001")
	require.NoError(t, err)

	mergedPR, err := prRepo.GetPullRequestByID(ctx, "pr-001")
	require.NoError(t, err)
	assert.Equal(t, models.PRStatusMerged, mergedPR.Status)
	assert.NotNil(t, mergedPR.MergedAt)

	userPRs, err := prSvc.GetUserPullRequests(ctx, reviewers[0])
	require.NoError(t, err)
	assert.Len(t, userPRs, 1)
	assert.Equal(t, "pr-001", userPRs[0].PullRequestID)
}
