package usecases

import (
	"task_manager/domain"
	"task_manager/repositories"
)

func GetAllTasks() ([]domain.Task, error) { return repositories.GetAllTasks() }

func GetTaskByID(id string) (*domain.Task, error) { return repositories.GetTaskByID(id) }

func CreateTask(t domain.Task) (*domain.Task, error) { return repositories.AddTask(t) }

func UpdateTask(id string, t domain.Task) error { return repositories.UpdateTask(id, t) }

func DeleteTask(id string) error { return repositories.DeleteTask(id) }
