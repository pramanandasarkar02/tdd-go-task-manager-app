package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gorilla/mux"
	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/domain"
	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/service"
	"github.com/stretchr/testify/assert"
)

// Define the mockTaskRepository directly in this file
type mockTaskRepository struct {
    tasks map[int]*domain.Task
}

func (m *mockTaskRepository) Save(task *domain.Task) error {
    if task.ID <= 0 {
        return errors.New("task ID must be positive")
    }
    if task.Title == "" {
        return errors.New("title cannot be empty")
    }
    m.tasks[task.ID] = task
    return nil
}

func (m *mockTaskRepository) FindByID(id int) (*domain.Task, error) {
    task, ok := m.tasks[id]
    if !ok {
        return nil, errors.New("task not found")
    }
    return task, nil
}

func (m *mockTaskRepository) FindAll() ([]domain.Task, error) {
    tasks := make([]domain.Task, 0, len(m.tasks))
    for _, task := range m.tasks {
        tasks = append(tasks, *task)
    }
    return tasks, nil
}

func (m *mockTaskRepository) Delete(id int) error {
    if _, ok := m.tasks[id]; !ok {
        return errors.New("task not found")
    }
    delete(m.tasks, id)
    return nil
}

func TestCreateTaskHandler(t *testing.T) {
    repo := &mockTaskRepository{tasks: make(map[int]*domain.Task)}
    svc := service.NewTaskService(repo)
    handler := NewTaskHandler(svc, slog.Default())

    task := domain.Task{ID: 1, Title: "Test Task"}
    body, _ := json.Marshal(task)
    req, _ := http.NewRequest("POST", "/tasks", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    handler.CreateTask(rr, req)

    assert.Equal(t, http.StatusCreated, rr.Code)
    var createdTask domain.Task
    json.NewDecoder(rr.Body).Decode(&createdTask)
    assert.Equal(t, "Test Task", createdTask.Title)
}

func TestGetTaskHandler(t *testing.T) {
    repo := &mockTaskRepository{tasks: make(map[int]*domain.Task)}
    svc := service.NewTaskService(repo)
    handler := NewTaskHandler(svc, slog.Default())

    repo.Save(&domain.Task{ID: 1, Title: "Test Task"})
    req, _ := http.NewRequest("GET", "/tasks/1", nil)
    req = mux.SetURLVars(req, map[string]string{"id": "1"})

    rr := httptest.NewRecorder()
    handler.GetTask(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    var task domain.Task
    json.NewDecoder(rr.Body).Decode(&task)
    assert.Equal(t, "Test Task", task.Title)
}

func TestGetAllTasksHandler(t *testing.T) {
    repo := &mockTaskRepository{tasks: make(map[int]*domain.Task)}
    svc := service.NewTaskService(repo)
    handler := NewTaskHandler(svc, slog.Default())

    repo.Save(&domain.Task{ID: 1, Title: "Task 1"})
    repo.Save(&domain.Task{ID: 2, Title: "Task 2"})
    req, _ := http.NewRequest("GET", "/tasks", nil)

    rr := httptest.NewRecorder()
    handler.GetAllTasks(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    var tasks []domain.Task
    json.NewDecoder(rr.Body).Decode(&tasks)
    assert.Len(t, tasks, 2)
}

func TestDeleteTaskHandler(t *testing.T) {
    repo := &mockTaskRepository{tasks: make(map[int]*domain.Task)}
    svc := service.NewTaskService(repo)
    handler := NewTaskHandler(svc, slog.Default())

    repo.Save(&domain.Task{ID: 1, Title: "Task 1"})
    req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
    req = mux.SetURLVars(req, map[string]string{"id": "1"})

    rr := httptest.NewRecorder()
    handler.DeleteTask(rr, req)

    assert.Equal(t, http.StatusNoContent, rr.Code)
    _, err := repo.FindByID(1)
    assert.Error(t, err)
}