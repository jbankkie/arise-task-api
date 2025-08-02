package main

import (
	"Arise-test/configs"
	"Arise-test/internal/handler"
	"Arise-test/internal/model"
	"Arise-test/internal/repository"
	"Arise-test/internal/routes"
	"Arise-test/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Set Gin mode
	gin.SetMode(config.Server.GinMode)

	// Initialize database
	db, err := initDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database
	err = db.AutoMigrate(&model.User{}, &model.Task{}, &model.Category{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	// categoryRepo := repository.NewCategoryRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	taskHandler := handler.NewTaskHandler(taskService)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, userHandler, taskHandler)

	// Start server
	log.Printf("Starting server on port %s", config.Server.Port)
	if err := router.Run(":" + config.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDB(config *configs.Config) (*gorm.DB, error) {
	dsn := config.GetDatabaseDSN()

	log.Printf("Connecting to database with DSN: postgres://%s:***@%s:%s/%s",
		config.Database.User, config.Database.Host, config.Database.Port, config.Database.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
