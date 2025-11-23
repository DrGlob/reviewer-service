package main

import (
    "log"
    "net/http"
    "REVIEWER-SERVICE/internal/handlers"
    "REVIEWER-SERVICE/internal/repository"
    "REVIEWER-SERVICE/internal/service"
)

func main() {
    // 1. Инициализируем БД
    db, err := repository.NewPostgresDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // 2. Инициализируем репозитории
    teamRepo := repository.NewTeamRepository(db)
    userRepo := repository.NewUserRepository(db)
    prRepo := repository.NewPRRepository(db)

    // 3. Инициализируем сервисы
    teamService := service.NewTeamService(teamRepo, userRepo)
    userService := service.NewUserService(userRepo, prRepo)
    prService := service.NewPRService(prRepo, userRepo, teamRepo)

    // 4. Инициализируем обработчики
    handlers := handlers.NewHandlers(teamService, userService, prService)

    // 5. Настраиваем роутинг
    http.HandleFunc("/team/add", handlers.CreateTeamHandler)
    http.HandleFunc("/team/get", handlers.GetTeamHandler)
    http.HandleFunc("/users/setIsActive", handlers.SetUserActiveHandler)
    http.HandleFunc("/users/getReview", handlers.GetUserReviewPRsHandler)
    http.HandleFunc("/pullRequest/create", handlers.CreatePRHandler)
    http.HandleFunc("/pullRequest/merge", handlers.MergePRHandler)
    http.HandleFunc("/pullRequest/reassign", handlers.ReassignReviewerHandler)

    // 6. Health check
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // 7. Запускаем сервер
    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Server failed:", err)
    }
}