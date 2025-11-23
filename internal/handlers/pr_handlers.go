package handlers

import (
    "encoding/json"
    "net/http"
)

func (h *Handlers) CreatePRHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        h.respondWithError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        PullRequestID   string `json:"pull_request_id"`
        PullRequestName string `json:"pull_request_name"`
        AuthorID        string `json:"author_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        h.respondWithError(w, "INVALID_INPUT", "Invalid JSON", http.StatusBadRequest)
        return
    }

    pr, err := h.prService.CreatePR(request.PullRequestID, request.PullRequestName, request.AuthorID)
    if err != nil {
        switch err.Error() {
        case "author not found", "team not found":
            h.respondWithError(w, "NOT_FOUND", "resource not found", http.StatusNotFound)
        case "PR already exists":
            h.respondWithError(w, "PR_EXISTS", "PR id already exists", http.StatusConflict)
        default:
            h.respondWithError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
        }
        return
    }

    h.respondWithJSON(w, map[string]interface{}{
        "pr": pr,
    }, http.StatusCreated)
}

func (h *Handlers) MergePRHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        h.respondWithError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        PullRequestID string `json:"pull_request_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        h.respondWithError(w, "INVALID_INPUT", "Invalid JSON", http.StatusBadRequest)
        return
    }

    pr, err := h.prService.MergePR(request.PullRequestID)
    if err != nil {
        if err.Error() == "PR not found" {
            h.respondWithError(w, "NOT_FOUND", "resource not found", http.StatusNotFound)
        } else {
            h.respondWithError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
        }
        return
    }

    h.respondWithJSON(w, map[string]interface{}{
        "pr": pr,
    }, http.StatusOK)
}

func (h *Handlers) ReassignReviewerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        h.respondWithError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        PullRequestID string `json:"pull_request_id"`
        OldUserID     string `json:"old_user_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        h.respondWithError(w, "INVALID_INPUT", "Invalid JSON", http.StatusBadRequest)
        return
    }

    pr, newReviewerID, err := h.prService.ReassignReviewer(request.PullRequestID, request.OldUserID)
    if err != nil {
        switch err.Error() {
        case "PR not found", "reviewer not found":
            h.respondWithError(w, "NOT_FOUND", "resource not found", http.StatusNotFound)
        case "cannot reassign on merged PR":
            h.respondWithError(w, "PR_MERGED", "cannot reassign on merged PR", http.StatusConflict)
        case "reviewer is not assigned to this PR":
            h.respondWithError(w, "NOT_ASSIGNED", "reviewer is not assigned to this PR", http.StatusConflict)
        case "no active replacement candidate in team":
            h.respondWithError(w, "NO_CANDIDATE", "no active replacement candidate in team", http.StatusConflict)
        default:
            h.respondWithError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
        }
        return
    }

    h.respondWithJSON(w, map[string]interface{}{
        "pr":          pr,
        "replaced_by": newReviewerID,
    }, http.StatusOK)
}