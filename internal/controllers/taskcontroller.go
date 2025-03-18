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

	GetTasks() error

	GetFinishedTasks() error

	GetTasksByDay(fiber.Ctx) error

	GetTasksByFrequency(fiber.Ctx) error

	GetTaskByUuid(uuid.UUID, fiber.Ctx) error

	UpdateTask(fiber.Ctx) error

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
		var response []map[string]interface{}
		for _, task := range tasks {
			response = append(response, map[string]interface{}{
				"task_id":     task.TaskId,
				"task":        task.Task,
				"description": task.Description,
			})
		}

		if response == nil {
			return c.Status(fiber.StatusOK).SendString("Tasks not found")
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func (tc *taskController) GetFinishedTasks(c *fiber.Ctx) error {
	var tasks []models.Task

	result := tc.db.UseGorm().Where("finished").Find(&tasks)

	if result.Error != nil {
		return result.Error
	} else {
		var response []map[string]interface{}
		for _, task := range tasks {
			response = append(response, map[string]interface{}{
				"task_id":     task.TaskId,
				"task":        task.Task,
				"description": task.Description,
			})
		}

		if response == nil {
			return c.Status(fiber.StatusOK).SendString("Tasks not found")
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func (tc *taskController) GetTaskByUuid(uuid uuid.UUID, c *fiber.Ctx) error {
	var task models.Task
	result := tc.db.UseGorm().First(&task, uuid)

	if result.Error != nil {
		return result.Error
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"task_id":     task.TaskId,
			"task":        task.Task,
			"description": task.Description,
			"finished":    task.Finished,
		})
	}
}

func (tc *taskController) GetTasksByDay(c *fiber.Ctx) error {
	type Result struct {
		taskId      uuid.UUID
		task        string
		description string
		finished    bool
	}

	var data map[string]interface{}
	var result Result

	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	tc.db.UseGorm().Model(&models.Task{}).
		Select("task.task_id, task.task, task.description, task.finished").
		Joins("join reminder using (task_id)").Where("reminder.repeat_days in ?", data).Scan(&result)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"task_id":     result.taskId,
		"task":        result.task,
		"description": result.description,
		"finished":    result.finished,
	})
}

func (tc *taskController) GetTasksByFrequency(c *fiber.Ctx) error {
	type Result struct {
		taskId      uuid.UUID
		task        string
		description string
		finished    bool
	}

	var data map[string]interface{}
	var result Result

	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	tc.db.UseGorm().Model(&models.Task{}).
		Select("task.task_id, task.task, task.description, task.finished").
		Joins("join reminder using (task_id)").Where("reminder.frequency in ?", data).Scan(&result)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"task_id":     result.taskId,
		"task":        result.task,
		"description": result.description,
		"finished":    result.finished,
	})
}

func (tc *taskController) UpdateTask(c *fiber.Ctx) error {
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	result := tc.db.UseGorm().Where("task_id = ?", task.TaskId).Save(&task)

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

	result := tc.db.UseGorm().Model(&task).Where("task_id = ?", task.TaskId).Update("finished", true)

	if result.Error != nil {
		return result.Error
	} else {
		message := fmt.Sprintf("Task %s (%v) finished", task.Task, task.TaskId)
		return c.Status(fiber.StatusCreated).SendString(message)
	}
}
