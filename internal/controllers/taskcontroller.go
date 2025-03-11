package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kevinhartarto/mytodolist/internal/database"
	"github.com/kevinhartarto/mytodolist/internal/models"
	"github.com/kevinhartarto/mytodolist/internal/utils"
)

type TaskController interface {

	// Get all task groups
	// return all active task groups
	GetAllTaskGroups(c *fiber.Ctx) error

	// Get all tasks
	// return all active tasks
	GetAllTasks(c *fiber.Ctx) error

	// Get task group
	// return a task group by UUID
	GetTaskGroup(c *fiber.Ctx) error

	// Get task
	// return a task by UUID
	GetTask(c *fiber.Ctx) error

	// Get all task in a task group
	// return all tasks within a task group
	GetAllTasksByTaskGroup(c *fiber.Ctx) error

	// Update a task group
	// return status of update
	UpdateTaskGroup(c *fiber.Ctx) error

	// Update a task
	// return status of update
	UpdateTask(c *fiber.Ctx) error

	// Create a task
	// return new task id
	CreateTask(c *fiber.Ctx) error

	// Create a task group
	// return new task group id
	CreateTaskGroup(c *fiber.Ctx) error
}

var (
	taskInstance *taskController

	taskGroups []models.TaskGroup
	tasks      []models.Task

	taskGroup models.TaskGroup
	task      models.Task

	affectedRows int64
)

type taskController struct {
	db database.Service
}

func NewTaskController(db database.Service) *taskController {

	if taskInstance != nil {
		return taskInstance
	}

	taskInstance = &taskController{
		db: db,
	}

	return taskInstance
}

func (tc *taskController) GetAllTaskGroups(c *fiber.Ctx) error {
	if err := tc.db.UseGorm().Where("deprecated is false").Find(&taskGroups).Error; err != nil {
		fmt.Printf("WARNING - all task groups, %v", err)
		return c.SendString("Record not found")
	}

	if len(taskGroups) == 0 {
		return c.SendString("No task group found! \n please create new task group")
	} else {
		result, _ := json.Marshal(taskGroups)
		return c.SendString(string(result))
	}
}

func (tc *taskController) GetAllTasks(c *fiber.Ctx) error {

	if err := tc.db.UseGorm().Where("deprecated is false").Find(&tasks).Error; err != nil {
		fmt.Printf("WARNING - all tasks, %v", err)
		return c.SendString("Record not found")
	}

	if len(tasks) == 0 {
		return c.SendString("No task found! \n please create new task")
	} else {
		result, _ := json.Marshal(tasks)
		return c.SendString(string(result))
	}
}

func (tc *taskController) GetTaskGroup(c *fiber.Ctx) error {
	if err := c.BodyParser(&taskGroup); err != nil {
		return err
	}

	if err := tc.db.UseGorm().First(&taskGroup).Error; err != nil {
		fmt.Printf("WARNING - task group, %v", err)
		return c.SendString("Record not found")
	}

	result, _ := json.Marshal(&taskGroup)
	return c.SendString(string(result))
}

func (tc *taskController) GetTask(c *fiber.Ctx) error {
	if err := c.BodyParser(&task); err != nil {
		return err
	}

	if err := tc.db.UseGorm().First(&task).Error; err != nil {
		fmt.Printf("WARNING - task, %v", err)
		return c.SendString("Record not found")
	}

	result, _ := json.Marshal(&task)
	return c.SendString(string(result))
}

func (tc *taskController) GetAllTasksByTaskGroup(c *fiber.Ctx) error {
	if err := c.BodyParser(&taskGroup); err != nil {
		return err
	}

	if err := tc.db.UseGorm().Table("task").
		Select("task.*").
		Joins("JOIN task_group_task ON task_group_task.task_id = task.task_id").
		Where("task_group_task.task_group_id = ?", taskGroup.TaskGroupId).
		Find(&task).Error; err != nil {
		fmt.Printf("WARNING - all tasks by task group, %v", err)
		return c.SendString("Record not found")
	}

	result, _ := json.Marshal(&task)
	return c.SendString(string(result))
}

func (tc *taskController) UpdateTaskGroup(c *fiber.Ctx) error {
	var updateTaskGroup struct {
		taskGroup   models.TaskGroup `json:"update_task_group"`
		updateType  string           `json:"update_type"`
		updateValue bool             `json:"update_value"`
	}

	if err := c.BodyParser(&updateTaskGroup); err != nil {
		return c.SendString("Record not found")
	}

	// Type of Update (update, deprecated)
	switch updateTaskGroup.updateType {
	case "update":
		affectedRows = tc.db.UseGorm().Save(&updateTaskGroup.taskGroup).RowsAffected
	case "deprecated":
		affectedRows = tc.db.UseGorm().Model(&updateTaskGroup.taskGroup).
			Update("deprecated", updateTaskGroup.updateValue).RowsAffected
	}

	// This is not a batch updates
	// We're expect only 1 affected row
	if affectedRows == 1 {
		result, _ := json.Marshal(&updateTaskGroup.taskGroup)
		return c.SendString(string(result))
	}
	return c.SendStatus(fiber.StatusBadRequest)
}

func (tc *taskController) UpdateTask(c *fiber.Ctx) error {
	var updateTask struct {
		task        models.Task `json:"update_task"`
		updateType  string      `json:"update_type"`
		updateValue bool        `json:"update_value"`
	}

	if err := c.BodyParser(&updateTask); err != nil {
		return err
	}

	// Type of Update (update, deprecated)
	switch updateTask.updateType {
	case "update":
		affectedRows = tc.db.UseGorm().Save(&updateTask.task).RowsAffected
	case "deprecated":
		affectedRows = tc.db.UseGorm().Model(&updateTask.task).
			Update("deprecated", updateTask.updateValue).RowsAffected
	}

	// This is not a batch updates
	// We're expect only 1 affected row
	if affectedRows == 1 {
		result, _ := json.Marshal(&updateTask.task)
		return c.SendString(string(result))
	}
	return c.SendStatus(fiber.StatusBadRequest)
}

func (tc *taskController) CreateTask(c *fiber.Ctx) error {
	if err := c.BodyParser(&task); err != nil {
		return err
	}

	task.TaskId, _ = utils.GenerateNewUUID()

	if err := tc.db.UseGorm().Create(&task).Error; err != nil {
		return err
	}

	response := fmt.Sprintf("Task (%v) created", task.TaskId)
	return c.SendString(response)
}

func (tc *taskController) CreateTaskGroup(c *fiber.Ctx) error {
	if err := c.BodyParser(&taskGroup); err != nil {
		return err
	}
	taskGroup.TaskGroupId, _ = utils.GenerateNewUUID()

	if err := tc.db.UseGorm().Create(&taskGroup).Error; err != nil {
		return err
	}

	response := fmt.Sprintf("Task group (%v) created", taskGroup.TaskGroupId)
	return c.SendString(response)
}
