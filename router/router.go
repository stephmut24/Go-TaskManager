package router

import (
	"net/http"
	"task_manager/controllers"
	"task_manager/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// auth routes
	router.POST("/register", controllers.Signup)
	router.POST("/login", controllers.Login)

	// Task routes
	taskRoutes := router.Group("/tasks")
	{
		// read endpoints: any authenticated user
		taskRoutes.GET("", middleware.AuthMiddleware(), controllers.GetTasks)
		taskRoutes.GET("/:id", middleware.AuthMiddleware(), controllers.GetTask)

		// write endpoints: admin only
		taskRoutes.POST("", middleware.AuthMiddleware(), middleware.AdminOnly(), controllers.AddTask)
		taskRoutes.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), controllers.UpdateTask)
		taskRoutes.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), controllers.DeleteTask)
	}

	// user management endpoints
	userRoutes := router.Group("/users")
	{
		// promote endpoint - admin only
		userRoutes.POST("/:id/promote", middleware.AuthMiddleware(), middleware.AdminOnly(), controllers.PromoteUser)
	}

	return router
}
