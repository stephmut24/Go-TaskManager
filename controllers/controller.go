package controllers

import (
	"net/http"
	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

//signup

func Signup(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := data.Signup(user)
	if err != nil {
		if err == data.ErrUserAlreadyExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})

}

// login
func Login(ctx *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// authenticate and generate JWT token
	user, err := data.Login(loginData.Username, loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := GenerateJWT(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user": gin.H{
			"id":        user.ID.Hex(),
			"username":  user.UserName,
			"user_type": user.User_type,
		},
	})
}

func GetTasks(ctx *gin.Context) {
	tasks, err := data.GetAllTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible to get tasks",
		})
		return
	}
	if len(tasks) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Aucune tache disponible",
			"tasks":   []string{},
		})
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"tasks": tasks})
}

func GetTask(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := data.GetTaskByID(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, task)
}

func AddTask(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask, err := data.AddTask(task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create a task"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, newTask)
}

func UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedTask models.Task
	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := data.UpdateTask(id, updatedTask)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	err := data.DeleteTask(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée ou impossible à supprimer"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tâche supprimée avec succès"})
}

// Promote a user to admin
func PromoteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	err := data.PromoteUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found or could not be promoted"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}
