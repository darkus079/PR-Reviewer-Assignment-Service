package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"pr-reviewer-assignment-service/internal/config"
	"pr-reviewer-assignment-service/internal/database"
	"pr-reviewer-assignment-service/internal/handlers"
	"pr-reviewer-assignment-service/internal/middleware"
	"pr-reviewer-assignment-service/internal/repository"
	"pr-reviewer-assignment-service/internal/services"
)

var (
	e2eDBContainer testcontainers.Container
	e2eDBPool      *pgxpool.Pool
	testServer     *httptest.Server
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:15-alpine",
		postgres.WithDatabase("e2edb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(2*time.Minute)),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	e2eDBContainer = pgContainer

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to resolve postgres host: %v", err)
	}

	mappedPort, err := pgContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("failed to resolve postgres port: %v", err)
	}

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:         "8080",
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Database: config.DatabaseConfig{
			Host:     host,
			Port:     mappedPort.Port(),
			User:     "postgres",
			Password: "postgres",
			DBName:   "e2edb",
			SSLMode:  "disable",
		},
	}

	pool, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("failed to create connection pool: %v", err)
	}

	e2eDBPool = pool.Pool

	if err := runE2EMigrations(connStr); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	router := setupTestRouter(pool.Pool)
	testServer = httptest.NewServer(router)

	code := m.Run()

	testServer.Close()

	if err := pgContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	e2eDBPool.Close()
	os.Exit(code)
}

func runE2EMigrations(connStr string) error {
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

func setupTestRouter(dbPool *pgxpool.Pool) *gin.Engine {
	userRepo := repository.NewPostgresUserRepository(dbPool)
	teamRepo := repository.NewPostgresTeamRepository(dbPool)
	prRepo := repository.NewPostgresPullRequestRepository(dbPool)

	userSvc := services.NewUserService(userRepo)
	teamSvc := services.NewTeamService(teamRepo, userRepo)
	prSvc := services.NewPullRequestService(prRepo, userRepo, teamRepo, userSvc)
	statSvc := services.NewStatisticService(prRepo, teamRepo, userRepo)

	handler := handlers.NewHandler(teamSvc, userSvc, prSvc, statSvc)
	healthHandler := handlers.NewHealthHandler(userRepo)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

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

	return r
}

func TestE2E_FullAPIWorkflow(t *testing.T) {
	setupE2ETestData(t)

	teamRequest := map[string]interface{}{
		"team_name": "e2e-team",
		"members": []map[string]interface{}{
			{"user_id": "e2e-user1", "username": "E2E User 1", "is_active": true},
			{"user_id": "e2e-user2", "username": "E2E User 2", "is_active": true},
			{"user_id": "e2e-user3", "username": "E2E User 3", "is_active": true},
		},
	}

	body, _ := json.Marshal(teamRequest)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var teamResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&teamResp)
	resp.Body.Close()

	assert.Equal(t, "e2e-team", teamResp["team_name"])

	prRequest := map[string]interface{}{
		"pull_request_id":   "e2e-pr-001",
		"pull_request_name": "E2E Test PR",
		"author_id":         "e2e-user1",
	}

	body, _ = json.Marshal(prRequest)
	req, _ = http.NewRequest("POST", testServer.URL+"/api/pullRequest/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer user-token")

	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var prResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&prResp)
	resp.Body.Close()

	assert.Equal(t, "e2e-pr-001", prResp["pull_request_id"])
	assert.Contains(t, prResp, "assigned_reviewers")

	teamGetRequest := map[string]interface{}{
		"team_name": "e2e-team",
	}

	body, _ = json.Marshal(teamGetRequest)
	req, _ = http.NewRequest("GET", testServer.URL+"/api/team/get", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer user-token")

	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var teamGetResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&teamGetResp)
	resp.Body.Close()

	assert.Equal(t, "e2e-team", teamGetResp["team_name"])
	assert.Len(t, teamGetResp["members"], 3)

	mergeRequest := map[string]interface{}{
		"pull_request_id": "e2e-pr-001",
	}

	body, _ = json.Marshal(mergeRequest)
	req, _ = http.NewRequest("POST", testServer.URL+"/api/pullRequest/merge", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer user-token")

	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var mergeResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&mergeResp)
	resp.Body.Close()

	assert.Contains(t, mergeResp, "message")
}

func TestE2E_HealthCheck(t *testing.T) {
	req, _ := http.NewRequest("GET", testServer.URL+"/health", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var healthResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&healthResp)
	resp.Body.Close()

	assert.Equal(t, "ok", healthResp["status"])
	assert.Contains(t, healthResp, "database")
}

func TestE2E_Authentication(t *testing.T) {
	req, _ := http.NewRequest("POST", testServer.URL+"/api/team/add", nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	req, _ = http.NewRequest("POST", testServer.URL+"/api/team/add", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func setupE2ETestData(t *testing.T) {
	ctx := context.Background()

	tables := []string{"pr_reviewers", "pull_requests", "users", "teams"}
	for _, table := range tables {
		_, err := e2eDBPool.Exec(ctx, "DELETE FROM "+table)
		require.NoError(t, err)
	}
}
