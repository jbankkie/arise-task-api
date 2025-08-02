package repository

import (
	"Arise-test/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *model.Task) error
	GetByID(id uuid.UUID) (*model.Task, error)
	GetByUserID(userID uuid.UUID, limit, offset int) ([]model.Task, error)
	GetByStatus(userID uuid.UUID, status model.TaskStatus, limit, offset int) ([]model.Task, error)
	GetByCategory(categoryID uuid.UUID, limit, offset int) ([]model.Task, error)
	Update(task *model.Task) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]model.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) GetByID(id uuid.UUID) (*model.Task, error) {
	var task model.Task
	err := r.db.Preload("User").Preload("Category").First(&task, "id = ?", id).Error
	return &task, err
}

func (r *taskRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Preload("Category").Where("user_id = ?", userID).
		Limit(limit).Offset(offset).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetByStatus(userID uuid.UUID, status model.TaskStatus, limit, offset int) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Preload("Category").Where("user_id = ? AND status = ?", userID, status).
		Limit(limit).Offset(offset).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetByCategory(categoryID uuid.UUID, limit, offset int) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Preload("User").Where("category_id = ?", categoryID).
		Limit(limit).Offset(offset).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) Update(task *model.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Task{}, "id = ?", id).Error
}

func (r *taskRepository) List(limit, offset int) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Preload("User").Preload("Category").
		Limit(limit).Offset(offset).Find(&tasks).Error
	return tasks, err
}
