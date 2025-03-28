package server

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/kevinhartarto/tasker/internal/controllers"
	"github.com/kevinhartarto/tasker/internal/database"
	"github.com/kevinhartarto/tasker/internal/utils"
)

func getCorsConfig() cors.Config {

	localConfig := cors.Config{
		AllowOrigins: "*",
	}

	return localConfig
}

func NewHandler(database database.Database) *fiber.App {

	app := fiber.New()
	app.Use(healthcheck.New())
	app.Use(cors.New(getCorsConfig()))

	logFile, err := os.OpenFile("api.log", os.O_RDWR|os.O_SYNC|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	app.Use(logger.New(logger.Config{
		Format:   "[${time}-${pid}] (${ip}) ${status} - ${method} ${path} | ${body}​\n",
		TimeZone: "Local",
		Output:   logFile,
	}))

	// API groups
	api := app.Group("/api")

	// Version 1
	v1 := api.Group("/v1")
	v1.Get("/metrics", monitor.New(monitor.Config{
		Title: "Tasker Metrics Page",
	}))

	// Tasker APIs
	list := controllers.NewTaskController(database)
	listAPI := v1.Group("/list")

	listAPI.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// Tasks
	listAPI.Get("/tasks", func(c *fiber.Ctx) error {
		return list.GetTasks(c)
	})
	listAPI.Get("/tasks/finished", func(c *fiber.Ctx) error {
		return list.GetFinishedTasks(c)
	})
	listAPI.Get("/task/:uuid", func(c *fiber.Ctx) error {
		uuid := utils.ParseUUID(c.Params("uuid"))
		return list.GetTaskByUuid(uuid, c)
	})
	listAPI.Post("/task", func(c *fiber.Ctx) error {
		return list.CreateTask(c)
	})
	listAPI.Put("/task", func(c *fiber.Ctx) error {
		return list.UpdateTask(c)
	})
	listAPI.Delete("/task", func(c *fiber.Ctx) error {
		return list.TaskFinished(c)
	})

	// Reminders
	reminder := controllers.NewReminderController(database)
	listAPI.Get("/reminders", func(c *fiber.Ctx) error {
		return reminder.GetAllReminders(c)
	})
	listAPI.Get("/reminder/:uuid", func(c *fiber.Ctx) error {
		uuid := utils.ParseUUID(c.Params("uuid"))
		return reminder.GetReminderByUuid(uuid, c)
	})
	listAPI.Get("/reminder/task/:uuid", func(c *fiber.Ctx) error {
		uuid := utils.ParseUUID(c.Params("uuid"))
		return reminder.GetReminderByTaskUuid(uuid, c)
	})
	listAPI.Post("/reminder", func(c *fiber.Ctx) error {
		return reminder.CreateRemainder(c)
	})
	listAPI.Put("/reminder", func(c *fiber.Ctx) error {
		return reminder.UpdateRemainder(c)
	})

	return app
}
