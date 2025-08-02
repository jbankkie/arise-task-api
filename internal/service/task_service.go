package service

import (
	"Arise-test/internal/model"
	"Arise-test/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

type TaskService interface {
	CreateTask(task *model.Task) error
	GetTaskByID(id uuid.UUID) (*model.Task, error)
	GetTasksByUserID(userID uuid.UUID, limit, offset int) ([]model.Task, error)
	GetTasksByStatus(userID uuid.UUID, status model.TaskStatus, limit, offset int) ([]model.Task, error)
	GetTasksByCategory(categoryID uuid.UUID, limit, offset int) ([]model.Task, error)
	UpdateTask(task *model.Task) error
	UpdateTaskStatus(id uuid.UUID, status model.TaskStatus) error
	DeleteTask(id uuid.UUID) error
	ListTasks(limit, offset int) ([]model.Task, error)
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
	}
}

func (s *taskService) CreateTask(task *model.Task) error {
	if task.Title == "" {
		return errors.New("task title is required")
	}

	if task.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}

	return s.taskRepo.Create(task)
}

func (s *taskService) GetTaskByID(id uuid.UUID) (*model.Task, error) {
	return s.taskRepo.GetByID(id)
}

func (s *taskService) GetTasksByUserID(userID uuid.UUID, limit, offset int) ([]model.Task, error) {
	return s.taskRepo.GetByUserID(userID, limit, offset)
}

func (s *taskService) GetTasksByStatus(userID uuid.UUID, status model.TaskStatus, limit, offset int) ([]model.Task, error) {
	return s.taskRepo.GetByStatus(userID, status, limit, offset)
}

func (s *taskService) GetTasksByCategory(categoryID uuid.UUID, limit, offset int) ([]model.Task, error) {
	return s.taskRepo.GetByCategory(categoryID, limit, offset)
}

func (s *taskService) UpdateTask(task *model.Task) error {
	task.UpdatedAt = time.Now()
	return s.taskRepo.Update(task)
}

func (s *taskService) UpdateTaskStatus(id uuid.UUID, status model.TaskStatus) error {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return err
	}

	task.Status = status
	task.UpdatedAt = time.Now()
	return s.taskRepo.Update(task)
}

func (s *taskService) DeleteTask(id uuid.UUID) error {
	return s.taskRepo.Delete(id)
}

func (s *taskService) ListTasks(limit, offset int) ([]model.Task, error) {
	return s.taskRepo.List(limit, offset)
}
