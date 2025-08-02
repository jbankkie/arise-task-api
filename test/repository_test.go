package test

import (
	"Arise-test/internal/model"
	"Arise-test/internal/repository"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&model.User{}, &model.Task{}, &model.Category{})

	return db
}

func TestUserRepository(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)

	// Test Create User
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}

	err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test Get User by Email
	retrievedUser, err := userRepo.GetByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	if retrievedUser.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", retrievedUser.Username)
	}

	// Test Get User by ID
	userByID, err := userRepo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if userByID.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", userByID.Email)
	}
}

func TestTaskRepository(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Create a test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	userRepo.Create(user)

	// Test Create Task
	task := &model.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      model.TaskStatusPending,
		Priority:    model.TaskPriorityMedium,
		UserID:      user.ID,
	}

	err := taskRepo.Create(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Test Get Task by ID
	retrievedTask, err := taskRepo.GetByID(task.ID)
	if err != nil {
		t.Fatalf("Failed to get task by ID: %v", err)
	}

	if retrievedTask.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", retrievedTask.Title)
	}

	// Test Get Tasks by User ID
	tasks, err := taskRepo.GetByUserID(user.ID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get tasks by user ID: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}
