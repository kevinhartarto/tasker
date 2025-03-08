package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kevinhartarto/mytodolist/internal/controllers"
	"github.com/kevinhartarto/mytodolist/internal/database"
)

func NewHandler(db database.Service) *fiber.App {

	app := fiber.New()
	app.Use(healthcheck.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// API groups
	api := app.Group("/api")

	// Version 1
	v1 := api.Group("/v1")

	// Tasks
	list := controllers.NewTaskController(db)
	day := controllers.NewDayController(db)
	listAPI := v1.Group("/list")
	listAPI.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	listAPI.Get("/", func(c *fiber.Ctx) error {
		return list.GetAllTaskGroups(c)
	})
	listAPI.Get("/tasks", func(c *fiber.Ctx) error {
		return list.GetAllTasks(c)
	})
	listAPI.Get("/group", func(c *fiber.Ctx) error {
		return list.GetTaskGroup(c)
	})
	listAPI.Get("/task", func(c *fiber.Ctx) error {
		return list.GetTask(c)
	})
	listAPI.Get("/group/tasks", func(c *fiber.Ctx) error {
		return list.GetAllTasksByTaskGroup(c)
	})
	listAPI.Put("/group/update", func(c *fiber.Ctx) error {
		return list.UpdateTaskGroup(c)
	})
	listAPI.Put("/task/update", func(c *fiber.Ctx) error {
		return list.UpdateTask(c)
	})
	listAPI.Get("/days", func(c *fiber.Ctx) error {
		return day.GetAllDays(c)
	})
	listAPI.Get("/day", func(c *fiber.Ctx) error {
		return day.GetDay(c)
	})

	return app
}
