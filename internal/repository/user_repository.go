package repository

import (
    "fmt"
    "REVIEWER-SERVICE/internal/entities"
)

type userRepository struct {
    db *PostgresDB
}

func NewUserRepository(db *PostgresDB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *entities.User) error {
    tx, err := r.db.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec(
        "INSERT INTO users (user_id, username, team_name, is_active) VALUES ($1, $2, $3, $4)",
        user.UserID, user.Username, user.TeamName, user.IsActive,
    )
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    return tx.Commit()
}


func (r *userRepository) GetUser(userID string) (*entities.User, error) {
    var user entities.User
    
    err := r.db.DB.QueryRow(`
        SELECT user_id, username, team_name, is_active 
        FROM users 
        WHERE user_id = $1`,
        userID,
    ).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
    
    if err != nil {
        return nil, fmt.Errorf("user '%s' not found: %w", userID, err)
    }

    return &user, nil
}

func (r *userRepository) UpdateUser(user *entities.User) error {
    tx, err := r.db.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec(
        "UPDATE users SET username = $1, team_name = $2, is_active = $3 WHERE user_id = $4",
        user.Username, user.TeamName, user.IsActive, user.UserID,
    )
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }

    return tx.Commit()
}

func (r *userRepository) GetActiveUsersByTeam(teamName string) ([]*entities.User, error) {

    rows, err := r.db.DB.Query(`
        SELECT user_id, username, team_name, is_active 
        FROM users 
        WHERE team_name = $1 AND is_active = true`,
        teamName,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get active users: %w", err)
    }
    defer rows.Close()

    var users []*entities.User

    for rows.Next() {
        var user entities.User
        err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
        if err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, &user)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error during rows iteration: %w", err)
    }

    return users, nil
}

