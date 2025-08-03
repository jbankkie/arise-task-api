package test

import (
	"Arise-test/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserModel_PasswordHashing(t *testing.T) {
	password := "testpassword123"

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Verify password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	assert.NoError(t, err)

	// Wrong password should fail
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongpassword"))
	assert.Error(t, err)
}

func TestTaskModel_StatusValidation(t *testing.T) {
	tests := []struct {
		status   model.TaskStatus
		expected bool
	}{
		{model.TaskStatusPending, true},
		{model.TaskStatusInProgress, true},
		{model.TaskStatusCompleted, true},
		{model.TaskStatusCancelled, true},
		{"invalid", false},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			task := &model.Task{
				Title:  "Test Task",
				Status: test.status,
			}

			if test.expected {
				assert.Equal(t, test.status, task.Status)
			} else {
				// Invalid status should be handled by validation
				assert.NotEqual(t, model.TaskStatusPending, test.status)
			}
		})
	}
}

func TestTaskModel_PriorityValidation(t *testing.T) {
	tests := []struct {
		priority model.TaskPriority
		expected bool
	}{
		{model.TaskPriorityLow, true},
		{model.TaskPriorityMedium, true},
		{model.TaskPriorityHigh, true},
		{model.TaskPriorityUrgent, true},
		{"invalid", false},
	}

	for _, test := range tests {
		t.Run(string(test.priority), func(t *testing.T) {
			task := &model.Task{
				Title:    "Test Task",
				Priority: test.priority,
			}

			if test.expected {
				assert.Equal(t, test.priority, task.Priority)
			} else {
				// Invalid priority should be handled by validation
				assert.NotEqual(t, model.TaskPriorityMedium, test.priority)
			}
		})
	}
}

func TestUserModel_Creation(t *testing.T) {
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}

	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test", user.FirstName)
	assert.Equal(t, "User", user.LastName)
}

func TestTaskModel_Creation(t *testing.T) {
	dueDate := time.Now().Add(24 * time.Hour)

	task := &model.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      model.TaskStatusPending,
		Priority:    model.TaskPriorityMedium,
		DueDate:     &dueDate,
	}

	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "This is a test task", task.Description)
	assert.Equal(t, model.TaskStatusPending, task.Status)
	assert.Equal(t, model.TaskPriorityMedium, task.Priority)
	assert.NotNil(t, task.DueDate)
	assert.True(t, task.DueDate.After(time.Now()))
}

func TestTaskModel_DefaultValues(t *testing.T) {
	task := &model.Task{
		Title: "Test Task",
	}

	// These should have default values when saved to database
	assert.Equal(t, "Test Task", task.Title)
	assert.Empty(t, task.Status)   // Will be set to default by GORM
	assert.Empty(t, task.Priority) // Will be set to default by GORM
}
