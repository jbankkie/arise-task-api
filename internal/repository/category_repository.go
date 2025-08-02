package repository

import (
	"Arise-test/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	GetByID(id uuid.UUID) (*model.Category, error)
	GetByUserID(userID uuid.UUID) ([]model.Category, error)
	Update(category *model.Category) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]model.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id uuid.UUID) (*model.Category, error) {
	var category model.Category
	err := r.db.Preload("Tasks").First(&category, "id = ?", id).Error
	return &category, err
}

func (r *categoryRepository) GetByUserID(userID uuid.UUID) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("user_id = ?", userID).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Category{}, "id = ?", id).Error
}

func (r *categoryRepository) List(limit, offset int) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Limit(limit).Offset(offset).Find(&categories).Error
	return categories, err
}
