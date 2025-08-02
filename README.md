# Task Manager API

A RESTful API for task management built with Go, Gin, GORM, and supports both PostgreSQL and MongoDB.

## Project Structure

```
your-app/
│
├── cmd/                  # main.go และ entry point
│   └── main.go
│
├── configs/              # env, config.yaml
│   ├── .env
│   └── config.yaml
│
├── internal/
│   ├── handler/          # Gin handler (controller)
│   │   ├── user_handler.go
│   │   └── task_handler.go
│   ├── model/            # Struct ของ DB
│   │   └── models.go
│   ├── repository/       # DB access layer (GORM/Mongo)
│   │   ├── user_repository.go
│   │   ├── task_repository.go
│   │   └── category_repository.go
│   ├── service/          # Logic ของระบบ
│   │   ├── user_service.go
│   │   └── task_service.go
│   └── routes/           # Routing setup
│       └── routes.go
│
├── test/                 # Unit tests
│   └── repository_test.go
│
├── Dockerfile
├── docker-compose.yml
├── init.sql              # PostgreSQL initialization
├── go.mod / go.sum
└── README.md
```

## Features

- **User Management**: Create, update, delete, and list users
- **Task Management**: CRUD operations for tasks with status and priority
- **Category Management**: Organize tasks by categories
- **Database Support**: PostgreSQL with GORM ORM
- **RESTful API**: Clean API design with proper HTTP methods
- **Docker Support**: Easy deployment with Docker Compose
- **Database Migration**: Auto migration with GORM

## Technologies Used

- **Go 1.18+**
- **Gin Web Framework**: Fast HTTP web framework
- **GORM**: Object-relational mapping library
- **PostgreSQL**: Relational database with GORM
- **Docker & Docker Compose**: Containerization
- **UUID**: For unique identifiers

## Quick Start

### Prerequisites

- Go 1.18 or higher
- Docker and Docker Compose

### 1. Clone the repository

```bash
git clone <repository-url>
cd Arise-test
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Start with Docker Compose

```bash
# Start PostgreSQL and App services
docker-compose up postgres app -d

# View logs
docker-compose logs -f app
```

### 4. Access the services

- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **pgAdmin** (PostgreSQL UI): http://localhost:5050
  - Email: admin@taskmanager.com
  - Password: admin

### 5. Manual setup (without Docker)

```bash
# 1. Setup PostgreSQL
createdb taskmanager

# 2. Update configs/.env with your database credentials

# 3. Run the application
go run cmd/main.go
```

## API Endpoints

### Health Check
```
GET /health
```

### Users
```
POST   /api/v1/users          # Create user
GET    /api/v1/users/:id      # Get user by ID
PUT    /api/v1/users/:id      # Update user
DELETE /api/v1/users/:id      # Delete user
GET    /api/v1/users          # List users (with pagination)
```

### Tasks
```
POST   /api/v1/tasks          # Create task
GET    /api/v1/tasks/:id      # Get task by ID
PUT    /api/v1/tasks/:id      # Update task
DELETE /api/v1/tasks/:id      # Delete task
GET    /api/v1/tasks          # Get user's tasks (with pagination & filtering)
```

## API Examples

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Create Task
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete Project",
    "description": "Finish the task management API",
    "priority": "high",
    "due_date": "2024-12-31T23:59:59Z"
  }'
```

### Get Tasks with filtering
```bash
# Get all tasks
curl http://localhost:8080/api/v1/tasks

# Get tasks by status
curl http://localhost:8080/api/v1/tasks?status=pending

# Get tasks with pagination
curl http://localhost:8080/api/v1/tasks?limit=5&offset=0
```

## Database Configuration

### PostgreSQL
- Default connection: `localhost:5432`
- Database: `taskmanager`
- User/Password: `postgres/password`

## Environment Variables

Create `configs/.env` file:

```env
PORT=8080
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=taskmanager

# JWT Secret
JWT_SECRET=your-secret-key
```

## Development

### Run tests
```bash
go test ./test/...
```

### Build application
```bash
go build -o bin/app cmd/main.go
```

### Run with hot reload (install air first)
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

## Docker Commands

```bash
# Build and start all services
docker-compose up --build

# Start only PostgreSQL service
docker-compose up postgres

# View logs
docker-compose logs -f app

# Stop all services
docker-compose down

# Remove volumes (caution: deletes data)
docker-compose down -v
```

## Database Schema

### Users Table
- `id` (UUID, Primary Key)
- `username` (String, Unique)
- `email` (String, Unique)
- `password` (String, Hashed)
- `first_name` (String)
- `last_name` (String)
- `created_at`, `updated_at`, `deleted_at`

### Tasks Table
- `id` (UUID, Primary Key)
- `title` (String, Required)
- `description` (Text)
- `status` (Enum: pending, in_progress, completed, cancelled)
- `priority` (Enum: low, medium, high, urgent)
- `due_date` (Timestamp, Optional)
- `user_id` (UUID, Foreign Key)
- `category_id` (UUID, Foreign Key, Optional)
- `created_at`, `updated_at`, `deleted_at`

### Categories Table
- `id` (UUID, Primary Key)
- `name` (String, Required)
- `description` (Text)
- `color` (String)
- `user_id` (UUID, Foreign Key)
- `created_at`, `updated_at`, `deleted_at`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
├── internal
│   └── app
│       └── app.go      # Main application logic
├── pkg
│   └── utils.go        # Utility functions
├── go.mod              # Module dependencies
├── go.sum              # Module checksums
└── README.md           # Project documentation
```

## Getting Started

### Prerequisites
- Go 1.16 or later
- A working Go environment

### Installation
1. Clone the repository:
   ```
   git clone <repository-url>
   cd Arise-test
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

### Running the Application
To run the application, execute the following command:
```
go run cmd/main.go
```

### Usage
Once the application is running, you can access it at `http://localhost:8080` (or the port specified in your application).

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.