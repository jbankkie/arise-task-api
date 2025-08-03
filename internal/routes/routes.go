package routes

import (
	"Arise-test/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	taskHandler *handler.TaskHandler,
	categoryHandler *handler.CategoryHandler,
) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("/", userHandler.ListUsers)
		}

		// Task routes
		tasks := v1.Group("/tasks")
		{
			tasks.POST("/", taskHandler.CreateTask)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
			tasks.GET("/", taskHandler.GetUserTasks)
		}

		// Category routes
		categories := v1.Group("/categories")
		{
			categories.POST("/", categoryHandler.CreateCategory)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
			categories.GET("/", categoryHandler.GetUserCategories)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Task Manager API is running",
		})
	})
}
