package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/handlers"
	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/logger"
	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/repository"
	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/service"
)

func main() {
    log := logger.NewLogger()
    repo := repository.NewInMemoryTaskRepository()
    svc := service.NewTaskService(repo)
    handler := handlers.NewTaskHandler(svc, log)

    r := mux.NewRouter()

    // API routes
    r.HandleFunc("/tasks", handler.CreateTask).Methods("POST")
    r.HandleFunc("/tasks", handler.GetAllTasks).Methods("GET")
    r.HandleFunc("/tasks/{id}", handler.GetTask).Methods("GET")
    r.HandleFunc("/tasks/{id}", handler.DeleteTask).Methods("DELETE")

    // Serve static files and UI
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "static/index.html")
    })
    log.Info("Starting server on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Error("Server failed", "error", err)
    }
}