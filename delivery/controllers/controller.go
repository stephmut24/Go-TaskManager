package controllers

import (
	"net/http"
	"task_manager/domain"
	"task_manager/usecases"

	"github.com/gin-gonic/gin"
)

func Signup(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := usecases.RegisterUser(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(ctx *gin.Context) {
	var payload struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, user, err := usecases.LoginUser(payload.Username, payload.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user": gin.H{
			"id":        user.ID.Hex(),
			"username":  user.Username,
			"user_type": user.UserType,
		},
	})
}

func GetTasks(ctx *gin.Context) {
	tasks, err := usecases.GetAllTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible to get tasks"})
		return
	}
	if len(tasks) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Aucune tache disponible", "tasks": []string{}})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"tasks": tasks})
}

func GetTask(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := usecases.GetTaskByID(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, task)
}

func AddTask(ctx *gin.Context) {
	var task domain.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTask, err := usecases.CreateTask(task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create a task"})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, newTask)
}

func UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var updated domain.Task
	if err := ctx.ShouldBindJSON(&updated); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := usecases.UpdateTask(id, updated); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := usecases.DeleteTask(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée ou impossible à supprimer"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Tâche supprimée avec succès"})
}

func PromoteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := usecases.PromoteUser(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found or could not be promoted"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}
