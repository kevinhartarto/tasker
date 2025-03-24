package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kevinhartarto/tasker/internal/database"
	"github.com/kevinhartarto/tasker/internal/models"
	"github.com/kevinhartarto/tasker/internal/utils"
)

type TaskController interface {

	// Create a task
	// return task name and uuid
	CreateTask(fiber.Ctx) error

	// Query all unfinished tasks
	// return an array of unfinished tasks
	GetTasks() error

	// Query all finished tasks
	// return an array of finished tasks
	GetFinishedTasks() error

	// Query all tasks group by day
	// return an array of tasks group by day
	GetTasksByDay(fiber.Ctx) error

	// Query all tasks group by frequency
	// return an array of tasks group by frequency
	GetTasksByFrequency(fiber.Ctx) error

	// Query a task by uuid
	// return task details
	GetTaskByUuid(uuid.UUID, fiber.Ctx) error

	// Update a task
	UpdateTask(fiber.Ctx) error

	// Change task status to finished
	TaskFinished(fiber.Ctx) error
}

type taskController struct {
	db database.Database
}

var taskInstance *taskController

func NewTaskController(db database.Database) *taskController {

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

	result := tc.db.Gorm().Create(&newTask)

	if result.Error != nil {
		return result.Error
	} else {
		message := fmt.Sprintf("Task %s (%v) created", newTask.Task, newTask.TaskId)
		return c.Status(fiber.StatusCreated).SendString(message)
	}
}

func (tc *taskController) GetTasks(c *fiber.Ctx) error {
	var tasks []models.Task

	result := tc.db.Gorm().Where("NOT finished").Find(&tasks)

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

	result := tc.db.Gorm().Where("finished").Find(&tasks)

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
	result := tc.db.Gorm().First(&task, uuid)

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

	tc.db.Gorm().Model(&models.Task{}).
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

	tc.db.Gorm().Model(&models.Task{}).
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

	result := tc.db.Gorm().Where("task_id = ?", task.TaskId).Save(&task)

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

	result := tc.db.Gorm().Model(&task).Where("task_id = ?", task.TaskId).Update("finished", true)

	if result.Error != nil {
		return result.Error
	} else {
		message := fmt.Sprintf("Task %s (%v) finished", task.Task, task.TaskId)
		return c.Status(fiber.StatusCreated).SendString(message)
	}
}
