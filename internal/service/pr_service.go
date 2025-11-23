package service

import (
    "fmt"
    "math/rand"
    "time"
    "REVIEWER-SERVICE/internal/entities"
    "REVIEWER-SERVICE/internal/repository"
)

type PRService struct {
    prRepo   repository.PRRepository
    userRepo repository.UserRepository
    teamRepo repository.TeamRepository
}

func NewPRService(prRepo repository.PRRepository, userRepo repository.UserRepository, teamRepo repository.TeamRepository) *PRService {
    rand.Seed(time.Now().UnixNano())
    return &PRService{
        prRepo:   prRepo,
        userRepo: userRepo,
        teamRepo: teamRepo,
    }
}

func (s *PRService) CreatePR(prID, prName, authorID string) (*entities.PullRequest, error) {
    author, err := s.userRepo.GetUser(authorID)
    if err != nil {
        return nil, fmt.Errorf("author not found: %w", err)
    }

    teamMembers, err := s.userRepo.GetActiveUsersByTeam(author.TeamName)
    if err != nil {
        return nil, fmt.Errorf("failed to get team members: %w", err)
    }

    reviewers := s.selectRandomReviewers(teamMembers, authorID)

    pr := &entities.PullRequest{
        PullRequestID:    prID,
        PullRequestName:  prName,
        AuthorID:         authorID,
        Status:           entities.StatusOpen,
        AssignedReviewers: reviewers,
    }

    err = s.prRepo.CreatePR(pr)
    if err != nil {
        return nil, fmt.Errorf("failed to create PR: %w", err)
    }

    return pr, nil
}


func (s *PRService) selectRandomReviewers(teamMembers []*entities.User, authorID string) []string {

    var availableUsers []string
    for _, person := range teamMembers {
        if person.UserID != authorID {
            availableUsers = append(availableUsers, person.UserID)
        }
    }

    if len(availableUsers) == 0 {
        return []string{}
    }

    shuffled := make([]string, len(availableUsers))
    copy(shuffled, availableUsers)
    
    for i := range shuffled {
        j := rand.Intn(len(shuffled))
        shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
    }

    if len(shuffled) > 2 {
        return shuffled[:2]
    }
    return shuffled
}

func (s *PRService) MergePR(prID string) (*entities.PullRequest, error) {
    pr, err := s.prRepo.GetPR(prID)
    if err != nil {
        return nil, fmt.Errorf("PR not found: %w", err)
    }

    if pr.Status == entities.StatusMerged {
        return pr, nil
    }

    now := time.Now()
    pr.Status = entities.StatusMerged
    pr.MergedAt = &now

    err = s.prRepo.UpdatePR(pr)
    if err != nil {
        return nil, fmt.Errorf("failed to merge PR: %w", err)
    }

    return pr, nil
}

func (s *PRService) ReassignReviewer(prID, oldReviewerID string) (*entities.PullRequest, string, error) {
    pr, err := s.prRepo.GetPR(prID)
    if err != nil {
        return nil, "", fmt.Errorf("PR not found: %w", err)
    }

    if pr.Status == entities.StatusMerged {
        return nil, "", fmt.Errorf("cannot reassign on merged PR")
    }

    isAssigned := false
    for _, reviewer := range pr.AssignedReviewers {
        if reviewer == oldReviewerID {
            isAssigned = true
            break
        }
    }
    if !isAssigned {
        return nil, "", fmt.Errorf("reviewer %s is not assigned to this PR", oldReviewerID)
    }

    oldReviewer, err := s.userRepo.GetUser(oldReviewerID)
    if err != nil {
        return nil, "", fmt.Errorf("reviewer not found: %w", err)
    }

    teamMembers, err := s.userRepo.GetActiveUsersByTeam(oldReviewer.TeamName)
    if err != nil {
        return nil, "", fmt.Errorf("failed to get team members: %w", err)
    }

    currentReviewers := make(map[string]bool)
    for _, reviewer := range pr.AssignedReviewers {
        currentReviewers[reviewer] = true
    }

    var availableUsers []string
    for _, user := range teamMembers {
        if user.UserID != pr.AuthorID && 
           user.UserID != oldReviewerID && 
           !currentReviewers[user.UserID] {
            availableUsers = append(availableUsers, user.UserID)
        }
    }

    if len(availableUsers) == 0 {
        return nil, "", fmt.Errorf("no active replacement candidate in team")
    }

    newReviewerID := availableUsers[rand.Intn(len(availableUsers))]

    err = s.prRepo.ReplaceReviewer(prID, oldReviewerID, newReviewerID)
    if err != nil {
        return nil, "", fmt.Errorf("failed to replace reviewer: %w", err)
    }

    updatedPR, err := s.prRepo.GetPR(prID)
    if err != nil {
        return nil, "", fmt.Errorf("failed to get updated PR: %w", err)
    }

    return updatedPR, newReviewerID, nil
}

// func (s *PRService) ReassignReviewer(prID, oldReviewerID string) (*entities.PullRequest, string, error) {

//     pr, err := s.prRepo.GetPR(prID)
//     if err != nil {
//         return nil, "", fmt.Errorf("PR not found: %w", err)
//     }

//     if pr.Status == entities.StatusMerged {
//         return nil, "", fmt.Errorf("cannot reassign on merged PR")
//     }

//     isAssigned := false
//     for _, reviewer := range pr.AssignedReviewers {
//         if reviewer == oldReviewerID {
//             isAssigned = true
//             break
//         }
//     }
//     if !isAssigned {
//         return nil, "", fmt.Errorf("reviewer %s is not assigned to this PR", oldReviewerID)
//     }

//     oldReviewer, err := s.userRepo.GetUser(oldReviewerID)
//     if err != nil {
//         return nil, "", fmt.Errorf("reviewer not found: %w", err)
//     }

//     teamMembers, err := s.userRepo.GetActiveUsersByTeam(oldReviewer.TeamName)
//     if err != nil {
//         return nil, "", fmt.Errorf("failed to get team members: %w", err)
//     }

//     var availableUsers []string
//     for _, user := range teamMembers {
//         if user.UserID != oldReviewerID && user.UserID != pr.AuthorID {
//             availableUsers = append(availableUsers, user.UserID)
//         }
//     }

//     if len(availableUsers) == 0 {
//         return nil, "", fmt.Errorf("no active replacement candidate in team")
//     }

//     newReviewerID := availableUsers[rand.Intn(len(availableUsers))]

//     err = s.prRepo.ReplaceReviewer(prID, oldReviewerID, newReviewerID)
//     if err != nil {
//         return nil, "", fmt.Errorf("failed to replace reviewer: %w", err)
//     }

//     updatedPR, err := s.prRepo.GetPR(prID)
//     if err != nil {
//         return nil, "", fmt.Errorf("failed to get updated PR: %w", err)
//     }

//     return updatedPR, newReviewerID, nil
// }