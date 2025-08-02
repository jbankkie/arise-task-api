# Task Manager API

A modern RESTful API for task management built with Go, Gin web framework, GORM ORM, and PostgreSQL database. Supports both Docker containerized deployment and local development with WSL.

## 🚀 Features

- **User Management**: Complete CRUD operations for users with authentication
- **Task Management**: Create, update, delete, and list tasks with status tracking
- **Category System**: Organize tasks by categories
- **Database**: PostgreSQL with GORM ORM and auto-migration
- **RESTful API**: Clean API design following REST conventions
- **Docker Support**: Full containerization with Docker Compose
- **WSL Support**: Local development using Windows Subsystem for Linux
- **Health Monitoring**: Built-in health check endpoints

## 🛠 Technologies Used

- **Go 1.23+**: Modern Go with latest features
- **Gin Framework**: High-performance HTTP web framework
- **GORM**: Feature-rich ORM library for Go
- **PostgreSQL 15**: Reliable relational database
- **Docker & Docker Compose**: Container orchestration
- **UUID**: Unique identifiers for all entities
- **WSL2**: Windows Subsystem for Linux for development

## 📁 Project Structure

```
arise-task-api/
│
├── cmd/                  # Application entry point
│   └── main.go          # Main server file
│
├── configs/              # Configuration package
│   └── config.go        # Centralized configuration loader
├── .env                 # Environment variables (root)
│
├── internal/             # Private application code
│   ├── handler/          # HTTP request handlers (controllers)
│   │   ├── user_handler.go
│   │   └── task_handler.go
│   ├── model/            # Database models
│   │   └── models.go
│   ├── repository/       # Data access layer
│   │   ├── user_repository.go
│   │   ├── task_repository.go
│   │   └── category_repository.go
│   ├── service/          # Business logic layer
│   │   ├── user_service.go
│   │   └── task_service.go
│   └── routes/           # API routing
│       └── routes.go
│
├── test/                 # Test files
├── docker-compose.yml    # Docker Compose configuration
├── Dockerfile            # Docker image configuration
├── init.sql              # Database initialization
├── go.mod / go.sum       # Go module files
└── README.md             # This file
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.23+**
- **Docker Desktop** with WSL2 integration
- **Windows Subsystem for Linux (WSL2)** with Ubuntu

### Option 1: Docker Compose (Recommended for Production)

1. **Clone the repository**
   ```bash
   git clone https://github.com/jbankkie/arise-task-api.git
   cd arise-task-api
   ```

2. **Start all services**
   ```bash
   # Start PostgreSQL, API, and pgAdmin
   docker compose up -d
   
   # View logs
   docker compose logs -f app
   ```

3. **Access the services**
   - **API**: http://localhost:8080
   - **Health Check**: http://localhost:8080/health
   - **pgAdmin**: http://localhost:5050
     - Email: `admin@taskmanager.com`
     - Password: `admin`

### Option 2: Local Development (WSL + Go)

Perfect for active development with hot reload and debugging.

1. **Setup WSL and Docker Integration**
   ```bash
   # Install Ubuntu WSL
   wsl --install Ubuntu
   
   # Enable Docker Desktop WSL integration:
   # Docker Desktop → Settings → Resources → WSL Integration
   # Enable "Ubuntu" distribution
   ```

2. **Start PostgreSQL in WSL**
   ```bash
   # Run PostgreSQL container in Ubuntu WSL
   wsl -d Ubuntu docker run --name postgres-wsl \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=password \
     -e POSTGRES_DB=taskmanager \
     -e POSTGRES_HOST_AUTH_METHOD=trust \
     -p 5432:5432 -d postgres:15-alpine
   ```

3. **Get WSL IP address and update configuration**
   ```bash
   # Get WSL IP
   wsl -d Ubuntu hostname -I
   # Example output: 172.19.18.36
   
   # Update .env file with WSL IP
   # DB_HOST=172.19.18.36
   ```

4. **Run Go application locally**
   ```bash
   # Install dependencies
   go mod tidy
   
   # Run the application
   go run cmd/main.go
   ```

5. **Development commands**
   ```bash
   # Stop PostgreSQL
   wsl -d Ubuntu docker stop postgres-wsl
   wsl -d Ubuntu docker rm postgres-wsl
   
   # Restart PostgreSQL
   wsl -d Ubuntu docker start postgres-wsl
   ```

## 🔧 Configuration

### Environment Variables

Create or update `.env` file in the root directory:

```env
# Server Configuration
PORT=8080
GIN_MODE=debug

# Database Configuration
# For Docker Compose: use 'postgres'
# For WSL Development: use WSL IP address (e.g., 172.19.18.36)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=taskmanager

# Security
JWT_SECRET=your-development-secret-key
```

### WSL IP Address Configuration

The WSL IP address changes on each Windows restart. Update your `.env` file:

```bash
# Get current WSL IP
wsl -d Ubuntu hostname -I

# Update .env file
DB_HOST=<WSL_IP_ADDRESS>
```

## 📚 API Documentation

### Health Check
```http
GET /health
```

### User Endpoints
```http
POST   /api/v1/users          # Create new user
GET    /api/v1/users/:id      # Get user by ID
PUT    /api/v1/users/:id      # Update user
DELETE /api/v1/users/:id      # Delete user
GET    /api/v1/users          # List all users (with pagination)
```

### Task Endpoints
```http
POST   /api/v1/tasks          # Create new task
GET    /api/v1/tasks/:id      # Get task by ID
PUT    /api/v1/tasks/:id      # Update task
DELETE /api/v1/tasks/:id      # Delete task
GET    /api/v1/tasks          # Get user's tasks (with filtering & pagination)
```

## 💡 API Usage Examples

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
    "status": "pending",
    "due_date": "2025-12-31T23:59:59Z"
  }'
```

### Get Tasks with Filtering
```bash
# Get all tasks
curl http://localhost:8080/api/v1/tasks

# Get tasks by status
curl "http://localhost:8080/api/v1/tasks?status=pending"

# Get tasks with pagination
curl "http://localhost:8080/api/v1/tasks?limit=5&offset=0"

# Get high priority tasks
curl "http://localhost:8080/api/v1/tasks?priority=high"
```

### PowerShell Examples (Windows)
```powershell
# Health Check
Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get

# Create User
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users" -Method Post -ContentType "application/json" -Body '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Get Users
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users" -Method Get
```

## 🐳 Docker Commands

### Development with Docker Compose
```bash
# Start all services
docker compose up -d

# Start only PostgreSQL
docker compose up postgres -d

# Start only the API
docker compose up app -d

# View logs
docker compose logs -f app
docker compose logs -f postgres

# Stop services
docker compose down

# Remove volumes (⚠️ deletes all data)
docker compose down -v

# Rebuild and restart
docker compose up --build -d
```

### WSL Development Commands
```bash
# Start PostgreSQL in WSL
wsl -d Ubuntu docker run --name postgres-wsl \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=taskmanager \
  -e POSTGRES_HOST_AUTH_METHOD=trust \
  -p 5432:5432 -d postgres:15-alpine

# Check WSL IP
wsl -d Ubuntu hostname -I

# Stop PostgreSQL
wsl -d Ubuntu docker stop postgres-wsl
wsl -d Ubuntu docker rm postgres-wsl

# Check PostgreSQL status
wsl -d Ubuntu docker ps | grep postgres
```

## 🗄️ Database Schema

### Users Table
```sql
- id           UUID PRIMARY KEY
- username     VARCHAR UNIQUE NOT NULL
- email        VARCHAR UNIQUE NOT NULL  
- password     VARCHAR NOT NULL (hashed)
- first_name   VARCHAR
- last_name    VARCHAR
- created_at   TIMESTAMP
- updated_at   TIMESTAMP
- deleted_at   TIMESTAMP (soft delete)
```

### Tasks Table
```sql
- id           UUID PRIMARY KEY
- title        VARCHAR NOT NULL
- description  TEXT
- status       ENUM (pending, in_progress, completed, cancelled)
- priority     ENUM (low, medium, high, urgent)
- due_date     TIMESTAMP
- user_id      UUID FOREIGN KEY
- category_id  UUID FOREIGN KEY
- created_at   TIMESTAMP
- updated_at   TIMESTAMP
- deleted_at   TIMESTAMP (soft delete)
```

### Categories Table
```sql
- id           UUID PRIMARY KEY
- name         VARCHAR NOT NULL
- description  TEXT
- color        VARCHAR
- user_id      UUID FOREIGN KEY
- created_at   TIMESTAMP
- updated_at   TIMESTAMP
- deleted_at   TIMESTAMP (soft delete)
```

## 🛠️ Development

### Running Tests
```bash
go test ./test/...
go test -v ./test/...
```

### Building the Application
```bash
# Build for current platform
go build -o bin/app cmd/main.go

# Build for Linux (Docker)
GOOS=linux GOARCH=amd64 go build -o bin/app-linux cmd/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o bin/app.exe cmd/main.go

# Note: bin/ directory is ignored in .gitignore
```

### Hot Reload Development
```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with hot reload (create .air.toml config first)
air
```

### Creating Air Configuration
```bash
# Generate default air config
air init

# Or create custom .air.toml
```

## 🔧 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check if PostgreSQL is running
   docker ps | grep postgres
   
   # Check WSL IP
   wsl -d Ubuntu hostname -I
   
   # Update .env with correct DB_HOST
   ```

2. **Port Already in Use**
   ```bash
   # Find process using port 8080
   netstat -ano | findstr :8080
   
   # Kill process (replace PID)
   taskkill /PID <PID> /F
   ```

3. **WSL Docker Issues**
   ```bash
   # Restart Docker Desktop
   # Enable WSL integration in Docker Desktop settings
   # Settings → Resources → WSL Integration → Enable Ubuntu
   ```

### Logs and Debugging
```bash
# Application logs
docker compose logs -f app

# Database logs  
docker compose logs -f postgres

# WSL PostgreSQL logs
wsl -d Ubuntu docker logs postgres-wsl
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)