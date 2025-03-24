package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/kevinhartarto/tasker/internal/database"
	"github.com/kevinhartarto/tasker/internal/logger"
	"github.com/kevinhartarto/tasker/internal/server"
	"github.com/kevinhartarto/tasker/internal/utils"
)

func main() {
	mainGroup := slog.Group(
		"environment",
		"stage", utils.GetEnvOrDefault("STAGE", "DEV"),
		"service", "tasker",
	)
	log := logger.GetLogger().With(mainGroup)

	// Instances
	gorm := database.Start()
	log.Info("Tasker connection with database established.")

	app := server.NewHandler(gorm)

	apiPort := utils.GetEnvOrDefault("PORT_API", "3030")
	apiAddr := fmt.Sprintf(":%v", apiPort)
	appErr := app.Listen(apiAddr)
	if appErr != nil {
		log.Error("Failed to start tasker, exiting...")
		closeApp(app, gorm)
	}

	log.Info(fmt.Sprintf("Server listening on http://localhost%s", apiAddr))
	log.Info("Tasker started and running.")
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
