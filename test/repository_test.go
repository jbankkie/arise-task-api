package test

import (
	"Arise-test/internal/model"
	"Arise-test/internal/repository"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB() *gorm.DB {
	// Use PostgreSQL test database (require Docker running)
	dsn := "postgres://postgres:password@localhost:5433/taskmanager_test?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		// If PostgreSQL is not available, skip tests
		log.Println("PostgreSQL test database not available, skipping tests...")
		log.Println("Run 'docker run --name postgres-test -e POSTGRES_PASSWORD=password -e POSTGRES_DB=taskmanager_test -p 5433:5432 -d postgres:15-alpine' to set up test DB")
		return nil
	}

	// Clean up any existing data
	db.Exec("DROP TABLE IF EXISTS tasks CASCADE")
	db.Exec("DROP TABLE IF EXISTS users CASCADE")
	db.Exec("DROP TABLE IF EXISTS categories CASCADE")

	// Migrate the schema
	err = db.AutoMigrate(&model.User{}, &model.Task{}, &model.Category{})
	if err != nil {
		log.Fatal("Failed to migrate test database:", err)
	}

	return db
}

func TestMain(m *testing.M) {
	// Setup
	log.Println("Setting up test environment...")

	// Run tests
	code := m.Run()

	// Teardown
	log.Println("Cleaning up test environment...")

	os.Exit(code)
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)

	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}

	err := userRepo.Create(user)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.False(t, user.CreatedAt.IsZero())
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Test GetByEmail
	foundUser, err := userRepo.GetByEmail("test@example.com")

	require.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, "test@example.com", foundUser.Email)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)

	foundUser, err := userRepo.GetByEmail("nonexistent@example.com")

	assert.Error(t, err)
	assert.Nil(t, foundUser)
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Test GetByID
	foundUser, err := userRepo.GetByID(user.ID)

	require.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, "testuser", foundUser.Username)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Update user
	user.FirstName = "Updated"
	user.LastName = "Name"
	err = userRepo.Update(user)

	require.NoError(t, err)
	assert.Equal(t, "Updated", user.FirstName)
	assert.Equal(t, "Name", user.LastName)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Delete user
	err = userRepo.Delete(user.ID)
	require.NoError(t, err)

	// Verify user is deleted (soft delete)
	foundUser, err := userRepo.GetByID(user.ID)
	assert.Error(t, err) // Should return error for soft deleted record
	assert.Nil(t, foundUser)
}

func TestTaskRepository_Create(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create test task
	dueDate := time.Now().Add(24 * time.Hour)
	task := &model.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      model.TaskStatusPending,
		Priority:    model.TaskPriorityMedium,
		DueDate:     &dueDate,
		UserID:      user.ID,
	}

	err = taskRepo.Create(task)

	require.NoError(t, err)
	assert.NotNil(t, task)
	assert.NotEqual(t, uuid.Nil, task.ID)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, model.TaskStatusPending, task.Status)
	assert.Equal(t, user.ID, task.UserID)
}

func TestTaskRepository_GetByID(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create test task
	task := &model.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      model.TaskStatusPending,
		Priority:    model.TaskPriorityMedium,
		UserID:      user.ID,
	}
	err = taskRepo.Create(task)
	require.NoError(t, err)

	// Test GetByID
	foundTask, err := taskRepo.GetByID(task.ID)

	require.NoError(t, err)
	assert.Equal(t, task.ID, foundTask.ID)
	assert.Equal(t, "Test Task", foundTask.Title)
}

func TestTaskRepository_GetByUserID(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create multiple tasks
	for i := 0; i < 3; i++ {
		task := &model.Task{
			Title:    fmt.Sprintf("Task %d", i+1),
			Status:   model.TaskStatusPending,
			Priority: model.TaskPriorityMedium,
			UserID:   user.ID,
		}
		err := taskRepo.Create(task)
		require.NoError(t, err)
	}

	// Get tasks by user ID
	tasks, err := taskRepo.GetByUserID(user.ID, 10, 0)

	require.NoError(t, err)
	assert.Len(t, tasks, 3)
	for _, task := range tasks {
		assert.Equal(t, user.ID, task.UserID)
	}
}

func TestTaskRepository_Update(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create test task
	task := &model.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      model.TaskStatusPending,
		Priority:    model.TaskPriorityMedium,
		UserID:      user.ID,
	}
	err = taskRepo.Create(task)
	require.NoError(t, err)

	// Update task
	task.Title = "Updated Task"
	task.Status = model.TaskStatusCompleted
	err = taskRepo.Update(task)

	require.NoError(t, err)
	assert.Equal(t, "Updated Task", task.Title)
	assert.Equal(t, model.TaskStatusCompleted, task.Status)
}

func TestTaskRepository_Delete(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create test task
	task := &model.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      model.TaskStatusPending,
		Priority:    model.TaskPriorityMedium,
		UserID:      user.ID,
	}
	err = taskRepo.Create(task)
	require.NoError(t, err)

	// Delete task
	err = taskRepo.Delete(task.ID)
	require.NoError(t, err)

	// Verify task is deleted (soft delete)
	foundTask, err := taskRepo.GetByID(task.ID)
	assert.Error(t, err) // Should return error for soft deleted record
	assert.Nil(t, foundTask)
}

func TestCategoryRepository_Create(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create test category
	category := &model.Category{
		Name:        "Work",
		Description: "Work related tasks",
		Color:       "#FF5722",
		UserID:      user.ID,
	}

	err = categoryRepo.Create(category)

	require.NoError(t, err)
	assert.NotNil(t, category)
	assert.NotEqual(t, uuid.Nil, category.ID)
	assert.Equal(t, "Work", category.Name)
	assert.Equal(t, "Work related tasks", category.Description)
	assert.Equal(t, "#FF5722", category.Color)
	assert.Equal(t, user.ID, category.UserID)
}

func TestCategoryRepository_GetByID(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create test category
	category := &model.Category{
		Name:        "Personal",
		Description: "Personal tasks",
		Color:       "#4CAF50",
		UserID:      user.ID,
	}
	err = categoryRepo.Create(category)
	require.NoError(t, err)

	// Test GetByID
	foundCategory, err := categoryRepo.GetByID(category.ID)

	require.NoError(t, err)
	assert.Equal(t, category.ID, foundCategory.ID)
	assert.Equal(t, "Personal", foundCategory.Name)
	assert.Equal(t, "#4CAF50", foundCategory.Color)
}

func TestCategoryRepository_GetByUserID(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create multiple categories
	categoryNames := []string{"Work", "Personal", "Shopping"}
	for _, name := range categoryNames {
		category := &model.Category{
			Name:   name,
			UserID: user.ID,
		}
		err := categoryRepo.Create(category)
		require.NoError(t, err)
	}

	// Get categories by user ID
	categories, err := categoryRepo.GetByUserID(user.ID)

	require.NoError(t, err)
	assert.Len(t, categories, 3)
	for _, category := range categories {
		assert.Equal(t, user.ID, category.UserID)
		assert.Contains(t, categoryNames, category.Name)
	}
}

func TestCategoryRepository_Update(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create test category
	category := &model.Category{
		Name:        "Old Name",
		Description: "Old description",
		Color:       "#000000",
		UserID:      user.ID,
	}
	err = categoryRepo.Create(category)
	require.NoError(t, err)

	// Update category
	category.Name = "Updated Name"
	category.Description = "Updated description"
	category.Color = "#FFFFFF"
	err = categoryRepo.Update(category)

	require.NoError(t, err)
	assert.Equal(t, "Updated Name", category.Name)
	assert.Equal(t, "Updated description", category.Description)
	assert.Equal(t, "#FFFFFF", category.Color)
}

func TestCategoryRepository_Delete(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)

	// Create test category
	category := &model.Category{
		Name:        "To Delete",
		Description: "This category will be deleted",
		UserID:      user.ID,
	}
	err = categoryRepo.Create(category)
	require.NoError(t, err)

	// Delete category
	err = categoryRepo.Delete(category.ID)
	require.NoError(t, err)

	// Verify category is deleted (soft delete)
	foundCategory, err := categoryRepo.GetByID(category.ID)
	assert.Error(t, err) // Should return error for soft deleted record
	assert.Nil(t, foundCategory)
}

func TestCategoryRepository_List(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Create test users
	user1 := &model.User{
		Username:  "user1",
		Email:     "user1@example.com",
		Password:  "hashedpassword",
		FirstName: "User",
		LastName:  "One",
	}
	err := userRepo.Create(user1)
	require.NoError(t, err)

	user2 := &model.User{
		Username:  "user2",
		Email:     "user2@example.com",
		Password:  "hashedpassword",
		FirstName: "User",
		LastName:  "Two",
	}
	err = userRepo.Create(user2)
	require.NoError(t, err)

	// Create categories for both users
	categories := []struct {
		name   string
		userID uuid.UUID
	}{
		{"Work", user1.ID},
		{"Personal", user1.ID},
		{"Shopping", user2.ID},
		{"Health", user2.ID},
		{"Travel", user1.ID},
	}

	for _, cat := range categories {
		category := &model.Category{
			Name:   cat.name,
			UserID: cat.userID,
		}
		err := categoryRepo.Create(category)
		require.NoError(t, err)
	}

	// Test List with pagination
	allCategories, err := categoryRepo.List(10, 0)

	require.NoError(t, err)
	assert.Len(t, allCategories, 5)

	// Test List with limit
	limitedCategories, err := categoryRepo.List(3, 0)

	require.NoError(t, err)
	assert.Len(t, limitedCategories, 3)

	// Test List with offset
	offsetCategories, err := categoryRepo.List(10, 2)

	require.NoError(t, err)
	assert.Len(t, offsetCategories, 3) // 5 total - 2 offset = 3
}
