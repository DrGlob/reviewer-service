package service

import (
    "fmt"
    "REVIEWER-SERVICE/internal/entities"
    "REVIEWER-SERVICE/internal/repository"
)

type UserService struct {
    userRepo repository.UserRepository
    prRepo   repository.PRRepository
}

func NewUserService(userRepo repository.UserRepository, prRepo repository.PRRepository) *UserService {
    return &UserService{
        userRepo: userRepo,
        prRepo:   prRepo,
    }
}

func (s *UserService) SetUserActive(userID string, isActive bool) (*entities.User, error) {
    user, err := s.userRepo.GetUser(userID)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    user.IsActive = isActive
    err = s.userRepo.UpdateUser(user)
    if err != nil {
        return nil, fmt.Errorf("failed to update user: %w", err)
    }

    return user, nil
}

func (s *UserService) GetUserReviewPRs(userID string) ([]*entities.PullRequestShort, error) {
    prs, err := s.prRepo.GetPRsByReviewer(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user review PRs: %w", err)
    }
    return prs, nil
}