package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pr-reviewer-assignment-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTeamRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTeamRepository(db *pgxpool.Pool) *PostgresTeamRepository {
	return &PostgresTeamRepository{db: db}
}

func (r *PostgresTeamRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	query := `
		INSERT INTO teams (team_name, created_at, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (team_name) DO UPDATE SET
			updated_at = EXCLUDED.updated_at
	`

	now := time.Now()
	team.CreatedAt = now
	team.UpdatedAt = now

	_, err := r.db.Exec(ctx, query, team.TeamName, team.CreatedAt, team.UpdatedAt)
	return err
}

func (r *PostgresTeamRepository) GetTeamByName(ctx context.Context, teamName string) (*models.Team, error) {
	query := `
		SELECT team_name, created_at, updated_at
		FROM teams
		WHERE team_name = $1
	`

	var team models.Team
	err := r.db.QueryRow(ctx, query, teamName).Scan(&team.TeamName, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &team, nil
}

func (r *PostgresTeamRepository) UpdateTeam(ctx context.Context, team *models.Team) error {
	query := `
		UPDATE teams
		SET updated_at = $2
		WHERE team_name = $1
	`

	team.UpdatedAt = time.Now()

	result, err := r.db.Exec(ctx, query, team.TeamName, team.UpdatedAt)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresTeamRepository) DeleteTeam(ctx context.Context, teamName string) error {
	query := `DELETE FROM teams WHERE team_name = $1`

	result, err := r.db.Exec(ctx, query, teamName)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresTeamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, teamName).Scan(&exists)
	return exists, err
}

func (r *PostgresTeamRepository) GetTeamWithMembers(ctx context.Context, teamName string) (*models.Team, error) {
	team, err := r.GetTeamByName(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	membersQuery := `
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY username
	`

	rows, err := r.db.Query(ctx, membersQuery, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMember
	for rows.Next() {
		var member models.TeamMember
		err := rows.Scan(&member.UserID, &member.Username, &member.IsActive)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	team.Members = members
	return team, nil
}
