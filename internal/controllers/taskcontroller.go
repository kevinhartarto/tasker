package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kevinhartarto/mytodolist/internal/database"
	"github.com/kevinhartarto/mytodolist/internal/models"
	"github.com/kevinhartarto/mytodolist/internal/utils"
)

type TaskController interface {
	CreateTask(fiber.Ctx) error

	CreateRemainder(fiber.Ctx) error

	GetTasks() error

	GetTasksByDay(fiber.Ctx) error

	GetTasksByFrequency(fiber.Ctx) error

	GetTaskByUuid(uuid.UUID, fiber.Ctx) error

	GetReminder(fiber.Ctx) error

	UpdateTask(fiber.Ctx) error

	UpdateRemainder(fiber.Ctx) error

	TaskFinished(fiber.Ctx) error
}

type taskController struct {
	db database.Service
}

var taskInstance *taskController

func NewTaskController(db database.Service) *taskController {

	if taskInstance != nil {
		return taskInstance
	}

	taskInstance = &taskController{
		db: db,
	}

	return taskInstance
}

func (tc *taskController) CreateTask(c *fiber.Ctx) error {
	var newTask models.Task

	if err := c.BodyParser(&newTask); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	newTask.TaskId = utils.GenerateNewUUID()

	if !utils.ValidateTask(newTask) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task name and id are required",
		})
	}

	result := tc.db.UseGorm().Create(&newTask)

	if result.Error != nil {
		return result.Error
	} else {
		message := fmt.Sprintf("Task %s (%v) created", newTask.Task, newTask.TaskId)
		return c.Status(fiber.StatusCreated).SendString(message)
	}
}

func (tc *taskController) GetTasks(c *fiber.Ctx) error {
	var tasks []models.Task

	result := tc.db.UseGorm().Where("NOT finished").Find(&tasks)

	if result.Error != nil {
		return result.Error
	} else {
		message, _ := json.Marshal(tasks)
		return c.Status(fiber.StatusOK).JSON(message)
	}
}

func (tc *taskController) GetTaskByUuid(uuid uuid.UUID, c *fiber.Ctx) error {
	var task models.Task
	result := tc.db.UseGorm().Where("NOT finished").First(&task, uuid)

	if result.Error != nil {
		return result.Error
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"task_id":     task.TaskId,
			"task":        task.Task,
			"description": task.Description,
		})
	}
}

func (tc *taskController) UpdateTask(c *fiber.Ctx) error {
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	result := tc.db.UseGorm().Save(&task)

	if result.Error != nil {
		return result.Error
	} else {
		message := fmt.Sprintf("Task %s (%v) updated", task.Task, task.TaskId)
		return c.Status(fiber.StatusCreated).SendString(message)
	}
}

func (tc *taskController) TaskFinished(c *fiber.Ctx) error {
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	result := tc.db.UseGorm().Model(&task).Update("finished", true)

	if result.Error != nil {
		return result.Error
	} else {
		message := fmt.Sprintf("Task %s (%v) finished", task.Task, task.TaskId)
		return c.Status(fiber.StatusCreated).SendString(message)
	}
}
