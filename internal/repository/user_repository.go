package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pr-reviewer-assignment-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (user_id, username, team_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET
			username = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active,
			updated_at = EXCLUDED.updated_at
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.Exec(ctx, query,
		user.UserID, user.Username, user.TeamName, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active, created_at, updated_at
		FROM users
		WHERE user_id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.UserID, &user.Username, &user.TeamName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $2, team_name = $3, is_active = $4, updated_at = $5
		WHERE user_id = $1
	`

	user.UpdatedAt = time.Now()

	result, err := r.db.Exec(ctx, query, user.UserID, user.Username, user.TeamName, user.IsActive, user.UpdatedAt)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE user_id = $1`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresUserRepository) GetUsersByTeam(ctx context.Context, teamName string) ([]*models.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active, created_at, updated_at
		FROM users
		WHERE team_name = $1
		ORDER BY username
	`

	rows, err := r.db.Query(ctx, query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, rows.Err()
}

func (r *PostgresUserRepository) GetActiveUsersByTeam(ctx context.Context, teamName string) ([]*models.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active, created_at, updated_at
		FROM users
		WHERE team_name = $1 AND is_active = true
		ORDER BY username
	`

	rows, err := r.db.Query(ctx, query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, rows.Err()
}

func (r *PostgresUserRepository) SetUserActiveStatus(ctx context.Context, userID string, isActive bool) error {
	query := `
		UPDATE users
		SET is_active = $2, updated_at = $3
		WHERE user_id = $1
	`

	result, err := r.db.Exec(ctx, query, userID, isActive, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresUserRepository) UserExists(ctx context.Context, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	return exists, err
}
