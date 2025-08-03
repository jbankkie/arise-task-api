package test

import (
	"Arise-test/internal/handler"
	"Arise-test/internal/model"
	"Arise-test/internal/repository"
	"Arise-test/internal/service"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestUserHandler_CreateUser(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := setupTestRouter()
	router.POST("/users", userHandler.CreateUser)

	user := map[string]interface{}{
		"username":   "testuser",
		"email":      "test@example.com",
		"password":   "password123",
		"first_name": "Test",
		"last_name":  "User",
	}

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	userData, ok := response["user"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "testuser", userData["username"])
	assert.Equal(t, "test@example.com", userData["email"])
}

func TestUserHandler_CreateUser_InvalidData(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := setupTestRouter()
	router.POST("/users", userHandler.CreateUser)

	// Missing required fields
	user := map[string]interface{}{
		"username": "testuser",
		// missing email and password
	}

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestUserHandler_GetUser(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

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

	router := setupTestRouter()
	router.GET("/users/:id", userHandler.GetUser)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", user.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	userData, ok := response["user"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "testuser", userData["username"])
	assert.Equal(t, "test@example.com", userData["email"])
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := setupTestRouter()
	router.GET("/users/:id", userHandler.GetUser)

	// Use a non-existent UUID
	req, _ := http.NewRequest("GET", "/users/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestTaskHandler_CreateTask(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	router := setupTestRouter()
	router.POST("/tasks", func(c *gin.Context) {
		// Set userID in context (simulate auth middleware)
		c.Set("userID", user.ID)
		taskHandler.CreateTask(c)
	})

	task := map[string]interface{}{
		"title":       "Test Task",
		"description": "This is a test task",
		"status":      "pending",
		"priority":    "medium",
		"user_id":     user.ID.String(),
	}

	jsonData, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	taskData, ok := response["task"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Test Task", taskData["title"])
	assert.Equal(t, "pending", taskData["status"])
}

func TestTaskHandler_GetTask(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	// Create test user first
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

	router := setupTestRouter()
	router.GET("/tasks/:id", taskHandler.GetTask)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/tasks/%s", task.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	taskData, ok := response["task"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Test Task", taskData["title"])
	assert.Equal(t, "pending", taskData["status"])
}

func TestTaskHandler_GetUserTasks(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	// Create test user first
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
			Title:    fmt.Sprintf("Task %d", i+1),
			Status:   model.TaskStatusPending,
			Priority: model.TaskPriorityMedium,
			UserID:   user.ID,
		}
		err := taskService.CreateTask(task)
		require.NoError(t, err)
	}

	router := setupTestRouter()
	router.GET("/tasks", func(c *gin.Context) {
		// Set userID in context (simulate auth middleware)
		c.Set("userID", user.ID)
		taskHandler.GetUserTasks(c)
	})

	req, _ := http.NewRequest("GET", fmt.Sprintf("/tasks?user_id=%s", user.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	tasks, ok := response["tasks"].([]interface{})
	require.True(t, ok)
	assert.Len(t, tasks, 3)
}

func TestCategoryHandler_CreateCategory(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	router := setupTestRouter()
	router.POST("/categories", categoryHandler.CreateCategory)

	category := map[string]interface{}{
		"name":        "Work Projects",
		"description": "All work-related project tasks",
		"color":       "#2196F3",
	}

	jsonData, _ := json.Marshal(category)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/categories?user_id=%s", user.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	categoryData, ok := response["category"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Work Projects", categoryData["name"])
	assert.Equal(t, "All work-related project tasks", categoryData["description"])
	assert.Equal(t, "#2196F3", categoryData["color"])
}

func TestCategoryHandler_CreateCategory_InvalidData(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	router := setupTestRouter()
	router.POST("/categories", categoryHandler.CreateCategory)

	// Missing required fields
	category := map[string]interface{}{
		"description": "Missing name field",
		// missing name
	}

	jsonData, _ := json.Marshal(category)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/categories?user_id=%s", user.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestCategoryHandler_GetCategory(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	// Create test category
	category := &model.Category{
		Name:        "Personal Goals",
		Description: "Personal development and goals",
		Color:       "#9C27B0",
		UserID:      user.ID,
	}
	err = categoryService.CreateCategory(category)
	require.NoError(t, err)

	router := setupTestRouter()
	router.GET("/categories/:id", categoryHandler.GetCategory)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/categories/%s", category.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	categoryData, ok := response["category"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Personal Goals", categoryData["name"])
	assert.Equal(t, "Personal development and goals", categoryData["description"])
	assert.Equal(t, "#9C27B0", categoryData["color"])
}

func TestCategoryHandler_GetCategory_NotFound(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	router := setupTestRouter()
	router.GET("/categories/:id", categoryHandler.GetCategory)

	// Use a non-existent UUID
	req, _ := http.NewRequest("GET", "/categories/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestCategoryHandler_GetUserCategories(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	// Create multiple categories
	categoryNames := []string{"Work", "Health", "Finance"}
	for _, name := range categoryNames {
		category := &model.Category{
			Name:   name,
			UserID: user.ID,
		}
		err := categoryService.CreateCategory(category)
		require.NoError(t, err)
	}

	router := setupTestRouter()
	router.GET("/categories", categoryHandler.GetUserCategories)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/categories?user_id=%s", user.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	categories, ok := response["categories"].([]interface{})
	require.True(t, ok)
	assert.Len(t, categories, 3)
}

func TestCategoryHandler_UpdateCategory(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	// Create test category
	category := &model.Category{
		Name:        "Original Name",
		Description: "Original description",
		Color:       "#000000",
		UserID:      user.ID,
	}
	err = categoryService.CreateCategory(category)
	require.NoError(t, err)

	router := setupTestRouter()
	router.PUT("/categories/:id", categoryHandler.UpdateCategory)

	updateData := map[string]interface{}{
		"name":        "Updated Name",
		"description": "Updated description",
		"color":       "#FF5722",
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/categories/%s", category.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	categoryData, ok := response["category"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Updated Name", categoryData["name"])
	assert.Equal(t, "Updated description", categoryData["description"])
	assert.Equal(t, "#FF5722", categoryData["color"])
}

func TestCategoryHandler_DeleteCategory(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("Test database not available")
	}
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Create test user first
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	err := userService.CreateUser(user)
	require.NoError(t, err)

	// Create test category
	category := &model.Category{
		Name:        "To Delete",
		Description: "This category will be deleted",
		UserID:      user.ID,
	}
	err = categoryService.CreateCategory(category)
	require.NoError(t, err)

	router := setupTestRouter()
	router.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/categories/%s", category.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "message")
	assert.Equal(t, "category deleted successfully", response["message"])
}
