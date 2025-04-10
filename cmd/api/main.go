package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kevinhartarto/tasker/internal/database"
	"github.com/kevinhartarto/tasker/internal/logger"
	"github.com/kevinhartarto/tasker/internal/server"
	"github.com/kevinhartarto/tasker/internal/utils"
)

var log = logger.GetLogger()

func main() {
	// Instances
	gorm := database.Start()
	log.Info("Tasker connection with database established.")

	redis := server.StartRedis()
	app := server.NewHandler(gorm, *redis)

	apiPort := utils.GetEnvOrDefault("PORT_API", "3030")
	apiAddr := fmt.Sprintf(":%v", apiPort)

	log.Info("Running Tasker server", "Port", apiAddr)
	utils.SendDesktopNotification("notify", "Tasker Running", "Tasker is now running")
	if appErr := app.Listen(apiAddr); appErr != nil {
		log.Error("Failed to start Tasker, exiting...")
		closeApp(app, gorm)
	}
}

func closeApp(app *fiber.App, gorm database.Database) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit // Wait for termination signal

	log.Info("Shutting down tasker...")

	// Gracefully shut down Fiber
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Info("Fiber shutdown error", "message: ", err)
	}

	// Closing database connection
	gorm.Close()
	log.Info("Tasker shutdown complete.")
}
