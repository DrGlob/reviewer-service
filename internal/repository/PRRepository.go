package repository

import (
	"fmt"
	"REVIEWER-SERVICE/internal/entities"
)

type prRepository struct {
	db *PostgresDB
}

func NewPRRepository(db *PostgresDB) *prRepository {
    return &prRepository{db: db}
}

func (r *prRepository) CreatePR(pr *entities.PullRequest) error {
	tx, err := r.db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status) 
		VALUES ($1, $2, $3, $4)`,
		pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	for _, reviewerID := range pr.AssignedReviewers {
		_, err = tx.Exec(`
			INSERT INTO pr_reviewers (pr_id, reviewer_id) 
			VALUES ($1, $2)`,
			pr.PullRequestID, reviewerID,
		)
		if err != nil {
			return fmt.Errorf("failed to assign reviewer %s: %w", reviewerID, err)
		}
	}

	return tx.Commit()
}

func (r *prRepository) GetPR(prID string) (*entities.PullRequest, error) {
	var pr entities.PullRequest
	
	err := r.db.DB.QueryRow(`
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pull_requests 
		WHERE pull_request_id = $1`,
		prID,
	).Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
	
	if err != nil {
		return nil, fmt.Errorf("PR '%s' not found: %w", prID, err)
	}

	reviewers, err := r.GetPRReviewers(prID)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *prRepository) UpdatePR(pr *entities.PullRequest) error {
	_, err := r.db.DB.Exec(`
		UPDATE pull_requests 
		SET pull_request_name = $1, author_id = $2, status = $3, merged_at = $4
		WHERE pull_request_id = $5`,
		pr.PullRequestName, pr.AuthorID, pr.Status, pr.MergedAt, pr.PullRequestID,
	)
	if err != nil {
		return fmt.Errorf("failed to update PR: %w", err)
	}
	return nil
}

func (r *prRepository) AssignReviewers(prID string, reviewerIDs []string) error {
	tx, err := r.db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM pr_reviewers WHERE pr_id = $1", prID)
	if err != nil {
		return fmt.Errorf("failed to clear old reviewers: %w", err)
	}

	for _, reviewerID := range reviewerIDs {
		_, err = tx.Exec(`
			INSERT INTO pr_reviewers (pr_id, reviewer_id) 
			VALUES ($1, $2)`,
			prID, reviewerID,
		)
		if err != nil {
			return fmt.Errorf("failed to assign reviewer %s: %w", reviewerID, err)
		}
	}

	return tx.Commit()
}

func (r *prRepository) ReplaceReviewer(prID, oldReviewerID, newReviewerID string) error {
    tx, err := r.db.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()


    var exists bool
    err = tx.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM pr_reviewers WHERE pr_id = $1 AND reviewer_id = $2)`,
        prID, oldReviewerID,
    ).Scan(&exists)
    
    if err != nil {
        return fmt.Errorf("failed to check reviewer assignment: %w", err)
    }
    if !exists {
        return fmt.Errorf("reviewer %s is not assigned to this PR", oldReviewerID)
    }

    var newReviewerExists bool
    err = tx.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM pr_reviewers WHERE pr_id = $1 AND reviewer_id = $2)`,
        prID, newReviewerID,
    ).Scan(&newReviewerExists)
    
    if err != nil {
        return fmt.Errorf("failed to check new reviewer: %w", err)
    }
    if newReviewerExists {
        return fmt.Errorf("reviewer %s is already assigned to this PR", newReviewerID)
    }

    _, err = tx.Exec(`
        UPDATE pr_reviewers 
        SET reviewer_id = $1 
        WHERE pr_id = $2 AND reviewer_id = $3`,
        newReviewerID, prID, oldReviewerID,
    )
    if err != nil {
        return fmt.Errorf("failed to replace reviewer: %w", err)
    }

    return tx.Commit()
}

func (r *prRepository) GetPRsByReviewer(userID string) ([]*entities.PullRequestShort, error) {
	rows, err := r.db.DB.Query(`
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
		FROM pull_requests pr
		JOIN pr_reviewers rev ON pr.pull_request_id = rev.pr_id
		WHERE rev.reviewer_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get PRs by reviewer: %w", err)
	}
	defer rows.Close()

	var prs []*entities.PullRequestShort
	for rows.Next() {
		var pr entities.PullRequestShort
		err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan PR: %w", err)
		}
		prs = append(prs, &pr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return prs, nil
}

func (r *prRepository) GetPRReviewers(prID string) ([]string, error) {
	rows, err := r.db.DB.Query(`
		SELECT reviewer_id 
		FROM pr_reviewers 
		WHERE pr_id = $1`,
		prID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR reviewers: %w", err)
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var reviewer string
		err := rows.Scan(&reviewer)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reviewer: %w", err)
		}
		reviewers = append(reviewers, reviewer)
	}

	return reviewers, nil
}

func (r *prRepository) IsPRMerged(prID string) (bool, error) {
	var status string
	err := r.db.DB.QueryRow(`
		SELECT status FROM pull_requests WHERE pull_request_id = $1`,
		prID,
	).Scan(&status)
	
	if err != nil {
		return false, fmt.Errorf("failed to check PR status: %w", err)
	}

	return status == "MERGED", nil
}

// GetStats возвращает статистику назначений
func (r *prRepository) GetStats() (*Stats, error) {
    stats := &Stats{}

    // Статистика по пользователям (кто сколько PR отревьюил)
    rows, err := r.db.DB.Query(`
        SELECT reviewer_id, COUNT(*) as review_count 
        FROM pr_reviewers 
        GROUP BY reviewer_id 
        ORDER BY review_count DESC
    `)
    if err != nil {
        return nil, fmt.Errorf("failed to get user stats: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var userStat UserStat
        err := rows.Scan(&userStat.UserID, &userStat.ReviewCount)
        if err != nil {
            return nil, fmt.Errorf("failed to scan user stat: %w", err)
        }
        stats.UserStats = append(stats.UserStats, userStat)
    }

    // Статистика по PR (сколько ревьюеров у каждого PR)
    rows, err = r.db.DB.Query(`
        SELECT pr_id, COUNT(*) as reviewer_count 
        FROM pr_reviewers 
        GROUP BY pr_id 
        ORDER BY reviewer_count DESC
    `)
    if err != nil {
        return nil, fmt.Errorf("failed to get PR stats: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var prStat PRStat
        err := rows.Scan(&prStat.PRID, &prStat.ReviewerCount)
        if err != nil {
            return nil, fmt.Errorf("failed to scan PR stat: %w", err)
        }
        stats.PRStats = append(stats.PRStats, prStat)
    }

    // Общая статистика
    err = r.db.DB.QueryRow(`
        SELECT 
            COUNT(DISTINCT pr_id) as total_prs,
            COUNT(*) as total_review_assignments,
            COUNT(DISTINCT reviewer_id) as unique_reviewers
        FROM pr_reviewers
    `).Scan(&stats.TotalPRs, &stats.TotalReviewAssignments, &stats.UniqueReviewers)
    
    if err != nil {
        return nil, fmt.Errorf("failed to get overall stats: %w", err)
    }

    return stats, nil
}