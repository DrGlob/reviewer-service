package repository

import "REVIEWER-SERVICE/internal/entities"

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
}