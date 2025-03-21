package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kevinhartarto/tasker/internal/database"
	"github.com/kevinhartarto/tasker/internal/server"
)

var ctx = context.Background

func main() {

	// Ports
	apiPort := os.Getenv("PORT_API")
	if apiPort == "" {
		apiPort = "3030"
	}

	// Instances
	db := database.StartDB()
	app := server.NewHandler(db)

	apiAddr := fmt.Sprintf(":%v", apiPort)
	fmt.Printf("Server listening on http://localhost%s\n", apiAddr)

	log.Fatal(app.Listen(apiAddr), db.Close())
}
