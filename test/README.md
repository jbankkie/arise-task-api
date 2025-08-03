# Unit Tests Documentation

## Overview

This project includes comprehensive unit tests covering all layers of the Clean Architecture:
- **Model layer**: Data structures, validation, and business rules
- **Repository layer**: Database operations and data access
- **Service layer**: Business logic and data processing
- **Handler layer**: HTTP endpoints and request/response handling

All tests are designed to ensure code quality, maintainability, and reliability.

## Test Structure

```
test/
├── model_test.go       # Model validation and creation tests
├── repository_test.go  # Database operation tests (requires PostgreSQL)
├── service_test.go     # Business logic tests (requires PostgreSQL)
└── handler_test.go     # HTTP handler tests (requires PostgreSQL)
```

## Running Tests

### Model Tests (No Database Required)
```bash
# Run all model tests
go test ./test/ -run "TestUserModel|TestTaskModel" -v

# Run specific model test
go test ./test/ -run TestUserModel_PasswordHashing -v
```

### Repository Tests (Requires PostgreSQL)
```bash
# Setup test database first
docker run --name postgres-test \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=taskmanager_test \
  -p 5432:5432 -d postgres:15-alpine

# Run repository tests
go test ./test/ -run TestUserRepository -v
go test ./test/ -run TestTaskRepository -v
```

### Service Tests (Requires PostgreSQL)
```bash
# Run service tests
go test ./test/ -run TestUserService -v
go test ./test/ -run TestTaskService -v
```

### Handler Tests (Requires PostgreSQL)
```bash
# Run handler tests
go test ./test/ -run TestUserHandler -v
go test ./test/ -run TestTaskHandler -v
```

### All Tests
```bash
# Run all available tests
go test ./test/ -v

# Run with coverage
go test ./test/ -v -cover

# Generate coverage report
go test ./test/ -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Test Categories

### 1. Model Tests ✅ (No Dependencies)
**Status**: Fully implemented and working
- **TestUserModel_PasswordHashing**: Password hashing and verification
- **TestUserModel_Creation**: User model creation and field validation
- **TestTaskModel_StatusValidation**: Task status enum validation (pending, in_progress, completed, cancelled)
- **TestTaskModel_PriorityValidation**: Task priority enum validation (low, medium, high, urgent)
- **TestTaskModel_Creation**: Task model creation with all fields
- **TestTaskModel_DefaultValues**: Default value handling

### 2. Repository Tests ⚠️ (Requires PostgreSQL)
**Status**: Implemented but requires PostgreSQL test database
- **TestUserRepository_Create**: User creation in database
- **TestUserRepository_GetByEmail**: Retrieve user by email
- **TestUserRepository_GetByID**: Retrieve user by ID
- **TestUserRepository_Update**: Update user information
- **TestUserRepository_Delete**: Soft delete user
- **TestTaskRepository_Create**: Task creation in database
- **TestTaskRepository_GetByID**: Retrieve task by ID
- **TestTaskRepository_GetByUserID**: Get tasks for specific user
- **TestTaskRepository_Update**: Update task information
- **TestTaskRepository_Delete**: Soft delete task

### 3. Service Tests ⚠️ (Requires PostgreSQL)
**Status**: Implemented but requires PostgreSQL test database
- **TestUserService_CreateUser**: User creation with password hashing
- **TestUserService_CreateUser_DuplicateEmail**: Duplicate email validation
- **TestUserService_GetUserByEmail**: Service layer user retrieval
- **TestUserService_GetUserByID**: Service layer user by ID
- **TestUserService_UpdateUser**: Service layer user updates
- **TestTaskService_CreateTask**: Task creation through service
- **TestTaskService_GetTaskByID**: Service layer task retrieval
- **TestTaskService_GetTasksByUserID**: Get user's tasks via service
- **TestTaskService_UpdateTask**: Service layer task updates
- **TestTaskService_DeleteTask**: Service layer task deletion
- **TestCategoryService_CreateCategory**: Category creation through service
- **TestCategoryService_GetCategoryByID**: Service layer category retrieval
- **TestCategoryService_GetCategoriesByUserID**: Get user's categories via service
- **TestCategoryService_UpdateCategory**: Service layer category updates
- **TestCategoryService_DeleteCategory**: Service layer category deletion

### 4. Handler Tests ⚠️ (Requires PostgreSQL)
**Status**: Implemented but requires PostgreSQL test database
- **TestUserHandler_CreateUser**: HTTP POST /users endpoint
- **TestUserHandler_CreateUser_InvalidData**: Invalid input handling
- **TestUserHandler_GetUser**: HTTP GET /users/:id endpoint
- **TestUserHandler_GetUser_NotFound**: 404 error handling
- **TestTaskHandler_CreateTask**: HTTP POST /tasks endpoint
- **TestTaskHandler_GetTask**: HTTP GET /tasks/:id endpoint
- **TestTaskHandler_GetUserTasks**: HTTP GET /tasks with user filter
- **TestCategoryHandler_CreateCategory**: HTTP POST /categories endpoint
- **TestCategoryHandler_GetCategory**: HTTP GET /categories/:id endpoint
- **TestCategoryHandler_GetUserCategories**: HTTP GET /categories with user filter
- **TestCategoryHandler_UpdateCategory**: HTTP PUT /categories/:id endpoint
- **TestCategoryHandler_DeleteCategory**: HTTP DELETE /categories/:id endpoint

## Test Database Setup

For tests that require PostgreSQL:

```bash
# Start test database
docker run --name postgres-test \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=taskmanager_test \
  -p 5432:5432 -d postgres:15-alpine

# Stop test database
docker stop postgres-test
docker rm postgres-test
```

## Test Coverage Goals

- **Model Layer**: ✅ 100% coverage
- **Repository Layer**: 🎯 90%+ coverage
- **Service Layer**: 🎯 85%+ coverage  
- **Handler Layer**: 🎯 80%+ coverage

## Continuous Integration

Tests can be integrated into CI/CD pipelines:

```yaml
# GitHub Actions example
- name: Run Tests
  run: |
    docker run -d --name postgres-test \
      -e POSTGRES_PASSWORD=password \
      -e POSTGRES_DB=taskmanager_test \
      -p 5432:5432 postgres:15-alpine
    
    sleep 10  # Wait for PostgreSQL to start
    go test ./test/ -v -cover
```

## Notes

- Model tests run without external dependencies
- Repository/Service/Handler tests require PostgreSQL test database
- Tests use isolated transactions for data consistency
- Soft delete testing ensures data integrity
- Password hashing validation ensures security
- HTTP handler tests validate complete request/response cycle
