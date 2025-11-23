package handlers

import (
    "net/http"
)

func (h *Handlers) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        h.respondWithError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    stats, err := h.prService.GetStats()
    if err != nil {
        h.respondWithError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
        return
    }

    h.respondWithJSON(w, map[string]interface{}{
        "stats": stats,
    }, http.StatusOK)
}