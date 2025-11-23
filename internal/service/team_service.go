package service

import (
    "fmt"
    "REVIEWER-SERVICE/internal/entities"
    "REVIEWER-SERVICE/internal/repository"
)

type TeamService struct {
    teamRepo repository.TeamRepository
    userRepo repository.UserRepository
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository) *TeamService {
    return &TeamService{
        teamRepo: teamRepo,
        userRepo: userRepo,
    }
}

func (s *TeamService) CreateTeam(team *entities.Team) error {
    exists, err := s.teamRepo.TeamExists(team.TeamName)
    if err != nil {
        return fmt.Errorf("failed to check team existence: %w", err)
    }
    if exists {
        return fmt.Errorf("team '%s' already exists", team.TeamName)
    }

    err = s.teamRepo.CreateTeam(team)
    if err != nil {
        return fmt.Errorf("failed to create team: %w", err)
    }

    return nil
}

func (s *TeamService) GetTeam(teamName string) (*entities.Team, error) {
    team, err := s.teamRepo.GetTeam(teamName)
    if err != nil {
        return nil, fmt.Errorf("failed to get team: %w", err)
    }
    return team, nil
}