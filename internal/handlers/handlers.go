package handlers

import (
    "encoding/json"
    "net/http"
    "REVIEWER-SERVICE/internal/service"
)

type Handlers struct {
    teamService *service.TeamService
    userService *service.UserService
    prService   *service.PRService
}

func NewHandlers(teamService *service.TeamService, userService *service.UserService, prService *service.PRService) *Handlers {
    return &Handlers{
        teamService: teamService,
        userService: userService,
        prService:   prService,
    }
}

func (h *Handlers) respondWithError(w http.ResponseWriter, code string, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error": map[string]string{
            "code":    code,
            "message": message,
        },
    })
}

func (h *Handlers) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}