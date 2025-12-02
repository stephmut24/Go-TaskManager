package main

import (
	"fmt"
	"os"

	"task_manager/config"
	"task_manager/delivery/routers"
	"task_manager/repositories"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	// initialize repositories after DB connection
	repositories.InitUserCollection()
	repositories.InitTaskCollection()

	r := routers.SetupRouter()

	port := config.GetEnv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Starting server on port", port)
	if err := r.Run(
		":" + port,
	); err != nil {
		fmt.Fprintln(os.Stderr, "server error:", err)
	}
}
