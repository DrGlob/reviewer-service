package main

import (
    "log"
    "net/http"
    "REVIEWER-SERVICE/internal/handlers"
    "REVIEWER-SERVICE/internal/repository"
    "REVIEWER-SERVICE/internal/service"
)

func main() {
    db, err := repository.NewPostgresDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    teamRepo := repository.NewTeamRepository(db)
    userRepo := repository.NewUserRepository(db)
    prRepo := repository.NewPRRepository(db)

    teamService := service.NewTeamService(teamRepo, userRepo)
    userService := service.NewUserService(userRepo, prRepo)
    prService := service.NewPRService(prRepo, userRepo, teamRepo)

    handlers := handlers.NewHandlers(teamService, userService, prService)

    http.HandleFunc("/team/add", handlers.CreateTeamHandler)
    http.HandleFunc("/team/get", handlers.GetTeamHandler)
    http.HandleFunc("/users/setIsActive", handlers.SetUserActiveHandler)
    http.HandleFunc("/users/getReview", handlers.GetUserReviewPRsHandler)
    http.HandleFunc("/pullRequest/create", handlers.CreatePRHandler)
    http.HandleFunc("/pullRequest/merge", handlers.MergePRHandler)
    http.HandleFunc("/pullRequest/reassign", handlers.ReassignReviewerHandler)
    http.HandleFunc("/stats", handlers.GetStatsHandler)

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Server failed:", err)
    }
}