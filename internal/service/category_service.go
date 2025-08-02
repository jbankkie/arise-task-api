package service

import (
	"Arise-test/internal/model"
	"Arise-test/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateCategory(category *model.Category) error
	GetCategoryByID(id uuid.UUID) (*model.Category, error)
	GetCategoriesByUserID(userID uuid.UUID) ([]model.Category, error)
	UpdateCategory(category *model.Category) error
	DeleteCategory(id uuid.UUID) error
	ListCategories(limit, offset int) ([]model.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(category *model.Category) error {
	if category.Name == "" {
		return errors.New("category name is required")
	}

	if category.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}

	return s.categoryRepo.Create(category)
}

func (s *categoryService) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *categoryService) GetCategoriesByUserID(userID uuid.UUID) ([]model.Category, error) {
	return s.categoryRepo.GetByUserID(userID)
}

func (s *categoryService) UpdateCategory(category *model.Category) error {
	return s.categoryRepo.Update(category)
}

func (s *categoryService) DeleteCategory(id uuid.UUID) error {
	return s.categoryRepo.Delete(id)
}

func (s *categoryService) ListCategories(limit, offset int) ([]model.Category, error) {
	return s.categoryRepo.List(limit, offset)
}
