package handlers

import (
    "encoding/json"
    "net/http"
    "REVIEWER-SERVICE/internal/entities"
)

func (h *Handlers) CreateTeamHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        h.respondWithError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var team entities.Team
    if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
        h.respondWithError(w, "INVALID_INPUT", "Invalid JSON", http.StatusBadRequest)
        return
    }

    err := h.teamService.CreateTeam(&team)
    if err != nil {
        if err.Error() == "team already exists" {
            h.respondWithError(w, "TEAM_EXISTS", "team_name already exists", http.StatusBadRequest)
        } else {
            h.respondWithError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
        }
        return
    }

    h.respondWithJSON(w, map[string]interface{}{
        "team": team,
    }, http.StatusCreated)
}

func (h *Handlers) GetTeamHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        h.respondWithError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    teamName := r.URL.Query().Get("team_name")
    if teamName == "" {
        h.respondWithError(w, "INVALID_INPUT", "team_name parameter is required", http.StatusBadRequest)
        return
    }

    team, err := h.teamService.GetTeam(teamName)
    if err != nil {
        if err.Error() == "team not found" {
            h.respondWithError(w, "NOT_FOUND", "resource not found", http.StatusNotFound)
        } else {
            h.respondWithError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
        }
        return
    }

    h.respondWithJSON(w, team, http.StatusOK)
}