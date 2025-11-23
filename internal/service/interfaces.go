package service

import "REVIEWER-SERVICE/internal/entities"

type TeamServiceInterface interface {
    CreateTeam(team *entities.Team) error
    GetTeam(teamName string) (*entities.Team, error)
}

type UserServiceInterface interface {
    SetUserActive(userID string, isActive bool) (*entities.User, error)
    GetUserReviewPRs(userID string) ([]*entities.PullRequestShort, error)
}

type PRServiceInterface interface {
    CreatePR(prID, prName, authorID string) (*entities.PullRequest, error)
    MergePR(prID string) (*entities.PullRequest, error)
    ReassignReviewer(prID, oldReviewerID string) (*entities.PullRequest, string, error)
}