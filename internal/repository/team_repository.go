package repository

import (
    "fmt"
    "REVIEWER-SERVICE/internal/entities"
)

type teamRepository struct {
    db *PostgresDB
}

func NewTeamRepository(db *PostgresDB) TeamRepository {
    return &teamRepository{db: db}
}

func (r *teamRepository) CreateTeam(team *entities.Team) error {
    tx, err := r.db.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec("INSERT INTO teams (team_name) VALUES ($1)", team.TeamName)
    if err != nil {
        return fmt.Errorf("failed to create team: %w", err)
    }

    for _, member := range team.Members {
        _, err = tx.Exec(
            "INSERT INTO users (user_id, username, team_name, is_active) VALUES ($1, $2, $3, $4)",
            member.UserID, member.Username, team.TeamName, member.IsActive,
        )
        if err != nil {
            return fmt.Errorf("failed to create user %s: %w", member.UserID, err)
        }
    }

    return tx.Commit()
}

func (r *teamRepository) GetTeam(teamName string) (*entities.Team, error) {
    var team entities.Team
    team.TeamName = teamName
    
    var exists bool
    err := r.db.DB.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`,
        teamName,
    ).Scan(&exists)
    
    if err != nil {
        return nil, fmt.Errorf("failed to check team existence: %w", err)
    }
	if !exists {
        return nil, fmt.Errorf("team '%s' doesn't exist", teamName)
    }

    rows, err := r.db.DB.Query(`
        SELECT user_id, username, is_active 
        FROM users 
        WHERE team_name = $1`,
        teamName,
    )

    if err != nil {
        return nil, fmt.Errorf("failed to get team members: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var member entities.TeamMember
        err := rows.Scan(&member.UserID, &member.Username, &member.IsActive)
        if err != nil {
            return nil, fmt.Errorf("failed to scan team member: %w", err)
        }
        team.Members = append(team.Members, member)
    }

    return &team, nil
}

func (r *teamRepository) TeamExists(teamName string) (bool, error) {
	var exists bool
	err := r.db.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`,
		teamName,
	).Scan(&exists)
	
	if err != nil {
		return false, fmt.Errorf("failed to check team existence: %w", err)
	}
	return exists, nil
}