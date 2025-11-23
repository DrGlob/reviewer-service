package handlers

import (
    "encoding/json"
    "net/http"
)


func (h *Handlers) SetUserActiveHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        h.respondWithError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        UserID   string `json:"user_id"`
        IsActive bool   `json:"is_active"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        h.respondWithError(w, "INVALID_INPUT", "Invalid JSON", http.StatusBadRequest)
        return
    }

    user, err := h.userService.SetUserActive(request.UserID, request.IsActive)
    if err != nil {
        if err.Error() == "user not found" {
            h.respondWithError(w, "NOT_FOUND", "resource not found", http.StatusNotFound)
        } else {
            h.respondWithError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
        }
        return
    }

    h.respondWithJSON(w, map[string]interface{}{
        "user": user,
    }, http.StatusOK)
}

func (h *Handlers) GetUserReviewPRsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        h.respondWithError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    userID := r.URL.Query().Get("user_id")
    if userID == "" {
        h.respondWithError(w, "INVALID_INPUT", "user_id parameter is required", http.StatusBadRequest)
        return
    }

    prs, err := h.userService.GetUserReviewPRs(userID)
    if err != nil {
        h.respondWithError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
        return
    }

    h.respondWithJSON(w, map[string]interface{}{
        "user_id":       userID,
        "pull_requests": prs,
    }, http.StatusOK)
}