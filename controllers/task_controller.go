package controllers

import (
	"task_manager/models"
	"task_manager/data"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetTasks(ctx *gin.Context) {
	tasks := data.GetAllTasks()
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

func AddTask(ctx *gin.Context){}

func UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedTask models.Task
	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := data.UpdateTask(id, updatedTask); err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return 
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task Updated"})
	
}
func DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := data.DeleteTask(id); err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task removed"})
}