package service

import (
	"errors"

	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/domain"
	"github.com/pramanandasarkar02/tdd-go-task-manager-app/internal/repository"
)

type TaskService struct {
    repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task *domain.Task) error {
    if task.Title == "" {
        return errors.New("title cannot be empty")
    }
    return s.repo.Save(task)
}

func (s *TaskService) GetTask(id int) (*domain.Task, error) {
    return s.repo.FindByID(id)
}

func (s *TaskService) GetAllTasks() ([]domain.Task, error) {
    return s.repo.FindAll()
}

func (s *TaskService) DeleteTask(id int) error {
    return s.repo.Delete(id)
}