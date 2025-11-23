package repository

import "REVIEWER-SERVICE/internal/entities"

type Stats struct {
    TotalPRs               int        `json:"total_prs"`
    TotalReviewAssignments int        `json:"total_review_assignments"`
    UniqueReviewers        int        `json:"unique_reviewers"`
    UserStats              []UserStat `json:"user_stats"`
    PRStats                []PRStat   `json:"pr_stats"`
}

type UserStat struct {
    UserID      string `json:"user_id"`
    ReviewCount int    `json:"review_count"`
}

type PRStat struct {
    PRID          string `json:"pr_id"`
    ReviewerCount int    `json:"reviewer_count"`
}


type TeamRepository interface {
    CreateTeam(team *entities.Team) error
    GetTeam(teamName string) (*entities.Team, error)
    TeamExists(teamName string) (bool, error)
}

type UserRepository interface {
    CreateUser(user *entities.User) error
    GetUser(userID string) (*entities.User, error)
    UpdateUser(user *entities.User) error
    GetActiveUsersByTeam(teamName string) ([]*entities.User, error)
}

type PRRepository interface {
    CreatePR(pr *entities.PullRequest) error
    GetPR(prID string) (*entities.PullRequest, error)
    UpdatePR(pr *entities.PullRequest) error
    AssignReviewers(prID string, reviewerIDs []string) error
    ReplaceReviewer(prID, oldReviewerID, newReviewerID string) error
    GetPRsByReviewer(userID string) ([]*entities.PullRequestShort, error)
    GetStats() (*Stats, error)
}