package routers

import (
	"net/http"
	dcontrollers "task_manager/delivery/controllers"
	infra "task_manager/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// auth routes
	router.POST("/register", dcontrollers.Signup)
	router.POST("/login", dcontrollers.Login)

	// Task routes
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("", infra.AuthMiddleware(), dcontrollers.GetTasks)
		taskRoutes.GET("/:id", infra.AuthMiddleware(), dcontrollers.GetTask)

		taskRoutes.POST("", infra.AuthMiddleware(), infra.AdminOnly(), dcontrollers.AddTask)
		taskRoutes.PUT("/:id", infra.AuthMiddleware(), infra.AdminOnly(), dcontrollers.UpdateTask)
		taskRoutes.DELETE("/:id", infra.AuthMiddleware(), infra.AdminOnly(), dcontrollers.DeleteTask)
	}

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/:id/promote", infra.AuthMiddleware(), infra.AdminOnly(), dcontrollers.PromoteUser)
	}

	return router
}
