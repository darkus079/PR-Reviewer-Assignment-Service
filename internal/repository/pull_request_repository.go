package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pr-reviewer-assignment-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPullRequestRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPullRequestRepository(db *pgxpool.Pool) *PostgresPullRequestRepository {
	return &PostgresPullRequestRepository{db: db}
}

func (r *PostgresPullRequestRepository) CreatePullRequest(ctx context.Context, pr *models.PullRequest) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	prQuery := `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	if pr.CreatedAt == nil {
		pr.CreatedAt = &now
	}

	_, err = tx.Exec(ctx, prQuery, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt)
	if err != nil {
		return err
	}

	if len(pr.AssignedReviewers) > 0 {
		err = r.setAssignedReviewersInTx(ctx, tx, pr.PullRequestID, pr.AssignedReviewers)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresPullRequestRepository) GetPullRequestByID(ctx context.Context, prID string) (*models.PullRequest, error) {
	query := `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at
		FROM pull_requests pr
		WHERE pr.pull_request_id = $1
	`

	var pr models.PullRequest
	var createdAt, mergedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, prID).Scan(
		&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &createdAt, &mergedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if createdAt.Valid {
		pr.CreatedAt = &createdAt.Time
	}
	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	reviewers, err := r.GetAssignedReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *PostgresPullRequestRepository) UpdatePullRequest(ctx context.Context, pr *models.PullRequest) error {
	query := `
		UPDATE pull_requests
		SET pull_request_name = $2, author_id = $3, status = $4, merged_at = $5
		WHERE pull_request_id = $1
	`

	result, err := r.db.Exec(ctx, query, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, pr.MergedAt)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresPullRequestRepository) DeletePullRequest(ctx context.Context, prID string) error {
	query := `DELETE FROM pull_requests WHERE pull_request_id = $1`

	result, err := r.db.Exec(ctx, query, prID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresPullRequestRepository) GetPullRequestsByReviewer(ctx context.Context, userID string) ([]*models.PullRequestShort, error) {
	query := `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
		FROM pull_requests pr
		JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id
		WHERE prr.user_id = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*models.PullRequestShort
	for rows.Next() {
		var pr models.PullRequestShort
		err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status)
		if err != nil {
			return nil, err
		}
		prs = append(prs, &pr)
	}

	return prs, rows.Err()
}

func (r *PostgresPullRequestRepository) MergePullRequest(ctx context.Context, prID string) error {
	query := `
		UPDATE pull_requests
		SET status = 'MERGED', merged_at = $2
		WHERE pull_request_id = $1 AND status = 'OPEN'
	`

	result, err := r.db.Exec(ctx, query, prID, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		exists, err := r.PullRequestExists(ctx, prID)
		if err != nil {
			return err
		}
		if !exists {
			return sql.ErrNoRows
		}
	}

	return nil
}

func (r *PostgresPullRequestRepository) PullRequestExists(ctx context.Context, prID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, prID).Scan(&exists)
	return exists, err
}

func (r *PostgresPullRequestRepository) GetAssignedReviewers(ctx context.Context, prID string) ([]string, error) {
	query := `
		SELECT user_id
		FROM pr_reviewers
		WHERE pull_request_id = $1
		ORDER BY assigned_at
	`

	rows, err := r.db.Query(ctx, query, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		reviewers = append(reviewers, userID)
	}

	return reviewers, rows.Err()
}

func (r *PostgresPullRequestRepository) SetAssignedReviewers(ctx context.Context, prID string, reviewers []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = r.setAssignedReviewersInTx(ctx, tx, prID, reviewers)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PostgresPullRequestRepository) setAssignedReviewersInTx(ctx context.Context, tx pgx.Tx, prID string, reviewers []string) error {
	deleteQuery := `DELETE FROM pr_reviewers WHERE pull_request_id = $1`
	_, err := tx.Exec(ctx, deleteQuery, prID)
	if err != nil {
		return err
	}

	if len(reviewers) > 0 {
		insertQuery := `
			INSERT INTO pr_reviewers (pull_request_id, user_id, assigned_at)
			VALUES ($1, $2, $3)
		`

		now := time.Now()
		for _, reviewerID := range reviewers {
			_, err = tx.Exec(ctx, insertQuery, prID, reviewerID, now)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
