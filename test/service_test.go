package test

import (
	"Arise-test/internal/model"
	"Arise-test/internal/repository"
	"Arise-test/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_CreateUser(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	err := userService.CreateUser(user)

	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)

	// Verify password is hashed
	assert.NotEqual(t, "password123", user.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	assert.NoError(t, err)
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// Create first user
	user1 := &model.User{
		Username:  "user1",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user1)
	require.NoError(t, err)

	// Try to create second user with same email
	user2 := &model.User{
		Username:  "user2",
		Email:     "test@example.com", // Same email
		Password:  "password456",
		FirstName: "Another",
		LastName:  "User",
	}
	err = userService.CreateUser(user2)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user with this email already exists")
}

func TestUserService_GetUserByEmail(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	// Get user by email
	foundUser, err := userService.GetUserByEmail("test@example.com")

	require.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, "test@example.com", foundUser.Email)
}

func TestUserService_GetUserByID(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	// Get user by ID
	foundUser, err := userService.GetUserByID(user.ID)

	require.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, "testuser", foundUser.Username)
}

func TestUserService_UpdateUser(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	// Update user
	user.FirstName = "Updated"
	user.LastName = "Name"
	err = userService.UpdateUser(user)

	require.NoError(t, err)
	assert.Equal(t, "Updated", user.FirstName)
	assert.Equal(t, "Name", user.LastName)
}

func TestTaskService_CreateTask(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
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

	err = taskService.CreateTask(task)

	require.NoError(t, err)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, model.TaskStatusPending, task.Status)
	assert.Equal(t, user.ID, task.UserID)
}

func TestTaskService_GetTaskByID(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	// Create test task
	task := &model.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      model.TaskStatusPending,
		Priority:    model.TaskPriorityMedium,
		UserID:      user.ID,
	}
	err = taskService.CreateTask(task)
	require.NoError(t, err)

	// Get task by ID
	foundTask, err := taskService.GetTaskByID(task.ID)

	require.NoError(t, err)
	assert.Equal(t, task.ID, foundTask.ID)
	assert.Equal(t, "Test Task", foundTask.Title)
}

func TestTaskService_GetTasksByUserID(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	// Create multiple tasks
	for i := 0; i < 3; i++ {
		task := &model.Task{
			Title:    "Test Task",
			Status:   model.TaskStatusPending,
			Priority: model.TaskPriorityMedium,
			UserID:   user.ID,
		}
		err := taskService.CreateTask(task)
		require.NoError(t, err)
	}

	// Get user tasks
	tasks, err := taskService.GetTasksByUserID(user.ID, 10, 0)

	require.NoError(t, err)
	assert.Len(t, tasks, 3)
	for _, task := range tasks {
		assert.Equal(t, user.ID, task.UserID)
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	// Create test task
	task := &model.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      model.TaskStatusPending,
		Priority:    model.TaskPriorityMedium,
		UserID:      user.ID,
	}
	err = taskService.CreateTask(task)
	require.NoError(t, err)

	// Update task
	task.Title = "Updated Task"
	task.Status = model.TaskStatusCompleted
	err = taskService.UpdateTask(task)

	require.NoError(t, err)
	assert.Equal(t, "Updated Task", task.Title)
	assert.Equal(t, model.TaskStatusCompleted, task.Status)
}

func TestTaskService_DeleteTask(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	// Create test user
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	// Create test task
	task := &model.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      model.TaskStatusPending,
		Priority:    model.TaskPriorityMedium,
		UserID:      user.ID,
	}
	err = taskService.CreateTask(task)
	require.NoError(t, err)

	// Delete task
	err = taskService.DeleteTask(task.ID)
	require.NoError(t, err)

	// Verify task is deleted
	foundTask, err := taskService.GetTaskByID(task.ID)
	assert.Error(t, err)
	assert.Nil(t, foundTask)
}
