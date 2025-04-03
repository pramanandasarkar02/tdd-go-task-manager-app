package service

import (
	"errors"

	"testing"

	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/domain"
	"github.com/stretchr/testify/assert"
)

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

func TestCreateTask(t *testing.T) {
    repo := &mockTaskRepository{tasks: make(map[int]*domain.Task)}
    svc := NewTaskService(repo)

    task := &domain.Task{ID: 1, Title: "Write TDD example"}
    err := svc.CreateTask(task)
    assert.NoError(t, err)

    savedTask, err := repo.FindByID(1)
    assert.NoError(t, err)
    assert.Equal(t, "Write TDD example", savedTask.Title)

    invalidTask := &domain.Task{ID: 2, Title: ""}
    err = svc.CreateTask(invalidTask)
    assert.Error(t, err)
    assert.Equal(t, "title cannot be empty", err.Error())
}

func TestGetTask(t *testing.T) {
    repo := &mockTaskRepository{tasks: make(map[int]*domain.Task)}
    svc := NewTaskService(repo)

    task := &domain.Task{ID: 1, Title: "Test Task"}
    repo.Save(task)

    foundTask, err := svc.GetTask(1)
    assert.NoError(t, err)
    assert.Equal(t, "Test Task", foundTask.Title)

    _, err = svc.GetTask(999)
    assert.Error(t, err)
    assert.Equal(t, "task not found", err.Error())
}

func TestGetAllTasks(t *testing.T) {
    repo := &mockTaskRepository{tasks: make(map[int]*domain.Task)}
    svc := NewTaskService(repo)

    repo.Save(&domain.Task{ID: 1, Title: "Task 1"})
    repo.Save(&domain.Task{ID: 2, Title: "Task 2"})

    tasks, err := svc.GetAllTasks()
    assert.NoError(t, err)
    assert.Len(t, tasks, 2)
}

func TestDeleteTask(t *testing.T) {
    repo := &mockTaskRepository{tasks: make(map[int]*domain.Task)}
    svc := NewTaskService(repo)

    repo.Save(&domain.Task{ID: 1, Title: "Task 1"})
    err := svc.DeleteTask(1)
    assert.NoError(t, err)

    _, err = repo.FindByID(1)
    assert.Error(t, err)

    err = svc.DeleteTask(999)
    assert.Error(t, err)
    assert.Equal(t, "task not found", err.Error())
}