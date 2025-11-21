package main

import (
	"task_manager/config"
	"task_manager/data"
	"task_manager/router"
)

func main() {

	config.LoadEnv()
	config.ConnectDB()
	data.InitTaskCollection() 
	
	
	r := router.SetupRouter()

	port := config.GetEnv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
