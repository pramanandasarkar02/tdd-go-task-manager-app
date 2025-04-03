package repository

import (
	"errors"

	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/domain"
)

type TaskRepository interface {
    Save(task *domain.Task) error
    FindByID(id int) (*domain.Task, error)
    FindAll() ([]domain.Task, error)
    Delete(id int) error
}

type InMemoryTaskRepository struct {
    tasks map[int]*domain.Task
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
    return &InMemoryTaskRepository{tasks: make(map[int]*domain.Task)}
}

func (r *InMemoryTaskRepository) Save(task *domain.Task) error {
    if task.ID <= 0 {
        return errors.New("task ID must be positive")
    }
    r.tasks[task.ID] = task
    return nil
}

func (r *InMemoryTaskRepository) FindByID(id int) (*domain.Task, error) {
    task, ok := r.tasks[id]
    if !ok {
        return nil, errors.New("task not found")
    }
    return task, nil
}

func (r *InMemoryTaskRepository) FindAll() ([]domain.Task, error) {
    tasks := make([]domain.Task, 0, len(r.tasks))
    for _, task := range r.tasks {
        tasks = append(tasks, *task)
    }
    return tasks, nil
}

func (r *InMemoryTaskRepository) Delete(id int) error {
    if _, ok := r.tasks[id]; !ok {
        return errors.New("task not found")
    }
    delete(r.tasks, id)
    return nil
}