package main

import (
	"Ecom/config"
	"Ecom/connections"
	"Ecom/routes"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

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
