package service

import (
	"Arise-test/internal/model"
	"Arise-test/internal/repository"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *model.User) error
	GetUserByID(id uuid.UUID) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uuid.UUID) error
	ListUsers(limit, offset int) ([]model.User, error)
	ValidatePassword(hashedPassword, password string) bool
	HashPassword(password string) (string, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(user *model.User) error {
	// Check if user already exists
	if _, err := s.userRepo.GetByEmail(user.Email); err == nil {
		return errors.New("user with this email already exists")
	}

	if _, err := s.userRepo.GetByUsername(user.Username); err == nil {
		return errors.New("user with this username already exists")
	}

	// Hash password
	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepo.Create(user)
}

func (s *userService) GetUserByID(id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *userService) GetUserByUsername(username string) (*model.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *userService) UpdateUser(user *model.User) error {
	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}

func (s *userService) ListUsers(limit, offset int) ([]model.User, error) {
	return s.userRepo.List(limit, offset)
}

func (s *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *userService) ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
