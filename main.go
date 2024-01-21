package main

import (
	"Ecom/config"
	"Ecom/connections"
	_ "Ecom/docs"
	"Ecom/routes"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// @title Ecom Services API
// @version 1.0
// @description An E-commerce API service in Go using Gin Framework
// @host localhost:8080
// @BasePath /
func main() {

	err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = connections.InitializeDB()
	if err != nil {
		log.Fatal(err)
	}

	router := routes.SetupRouter()

	fmt.Println("Starting Router...")
	router.Run(":8080")

	defer connections.DB.Close()
}
