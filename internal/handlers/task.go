package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/domain"
	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/service"
)

type TaskHandler struct {
    svc    *service.TaskService
    logger *slog.Logger
}

func NewTaskHandler(svc *service.TaskService, logger *slog.Logger) *TaskHandler {
    return &TaskHandler{svc: svc, logger: logger}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var task domain.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        h.logger.Error("Invalid request body", "error", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := h.svc.CreateTask(&task); err != nil {
        h.logger.Error("Failed to create task", "error", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        h.logger.Error("Invalid task ID", "error", err)
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    task, err := h.svc.GetTask(id)
    if err != nil {
        h.logger.Warn("Task not found", "id", id)
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
    tasks, err := h.svc.GetAllTasks()
    if err != nil {
        h.logger.Error("Failed to fetch tasks", "error", err)
        http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        h.logger.Error("Invalid task ID", "error", err)
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    if err := h.svc.DeleteTask(id); err != nil {
        h.logger.Warn("Task not found", "id", id)
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}