package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB() (*PostgresDB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnvInt("DB_PORT", 5432)
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "service")
	dbname := getEnv("DB_NAME", "reviewer")

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	postgresDB := &PostgresDB{DB: db}

	if err := postgresDB.InitSchema(); err != nil {
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL and initialized schema")
	return postgresDB, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (p *PostgresDB) Close() error {
	return p.DB.Close()
}

func (p* PostgresDB) InitSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS teams (
            team_name VARCHAR(100) PRIMARY KEY,
            created_at TIMESTAMP DEFAULT NOW()
        )`,

        `CREATE TABLE IF NOT EXISTS users (
            user_id VARCHAR(100) PRIMARY KEY,
            username VARCHAR(100) NOT NULL,
            team_name VARCHAR(100) REFERENCES teams(team_name) ON DELETE CASCADE,
            is_active BOOLEAN DEFAULT TRUE,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW()
        )`,

        `CREATE TABLE IF NOT EXISTS pull_requests (
            pull_request_id VARCHAR(100) PRIMARY KEY,
            pull_request_name VARCHAR(200) NOT NULL,
            author_id VARCHAR(100) REFERENCES users(user_id),
            status VARCHAR(20) DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
            created_at TIMESTAMP DEFAULT NOW(),
            merged_at TIMESTAMP NULL
        )`,

        `CREATE TABLE IF NOT EXISTS pr_reviewers (
            pr_id VARCHAR(100) REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
            reviewer_id VARCHAR(100) REFERENCES users(user_id),
            assigned_at TIMESTAMP DEFAULT NOW(),
            PRIMARY KEY (pr_id, reviewer_id)
        )`,

        `CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_name, is_active)`,
        `CREATE INDEX IF NOT EXISTS idx_prs_status ON pull_requests(status)`,
        `CREATE INDEX IF NOT EXISTS idx_prs_author ON pull_requests(author_id)`,
        `CREATE INDEX IF NOT EXISTS idx_reviewers_pr ON pr_reviewers(pr_id)`,
        `CREATE INDEX IF NOT EXISTS idx_reviewers_user ON pr_reviewers(reviewer_id)`,
	}

	for _, query := range queries {
        if _, err := p.DB.Exec(query); err != nil {
            return fmt.Errorf("failed to execute query '%s': %w", query, err)
        }
    }
    
    log.Println("Database schema initialized successfully")
    return nil
}

