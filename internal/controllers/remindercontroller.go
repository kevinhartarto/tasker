package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kevinhartarto/mytodolist/internal/database"
	"github.com/kevinhartarto/mytodolist/internal/models"
	"github.com/kevinhartarto/mytodolist/internal/utils"
)

type ReminderController interface {
	CreateRemainder(fiber.Ctx) error

	GetAllReminders(fiber.Ctx) error

	GetReminderByUuid(uuid.UUID, fiber.Ctx) error

	GetReminderByTaskUuid(uuid.UUID, fiber.Ctx) error

	UpdateRemainder(fiber.Ctx) error
}

type reminderController struct {
	db database.Service
}

var reminderInstance *reminderController

func NewReminderController(db database.Service) *reminderController {
	if reminderInstance != nil {
		return reminderInstance
	}

	reminderInstance = &reminderController{
		db: db,
	}

	return reminderInstance
}

func (rc *reminderController) CreateRemainder(c *fiber.Ctx) error {
	var newReminder models.Reminder

	if err := c.BodyParser(&newReminder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	if newReminder.StartTime.IsZero() {
		newReminder.StartTime = time.Now()
	}

	if !utils.ValidateReminder(newReminder) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task name and id are required",
		})
	}

	result := rc.db.UseGorm().Create(&newReminder)

	if result.Error != nil {
		return result.Error
	} else {
		message := fmt.Sprintf("Reminder %s (%v) for task (%v) created",
			newReminder.Reminder, newReminder.ReminderId, newReminder.TaskId)
		return c.Status(fiber.StatusCreated).SendString(message)
	}
}

func (rc *reminderController) GetAllReminders(c *fiber.Ctx) error {
	var reminders []models.Reminder
	result := rc.db.UseGorm().Find(&reminders)

	if result.Error != nil {
		return result.Error
	} else {
		var response []map[string]interface{}
		for _, reminder := range reminders {
			response = append(response, map[string]interface{}{
				"reminder_id":         reminder.ReminderId,
				"reminder":            reminder.Reminder,
				"task_id":             reminder.TaskId,
				"description":         reminder.Description,
				"start_time":          reminder.StartTime,
				"frequency":           reminder.Frequency,
				"repeat_days":         reminder.RepeatDays,
				"repeat_sameday":      reminder.RepeatSameday,
				"repeat_until":        reminder.RepeatUntil,
				"interval":            reminder.Interval,
				"interval_in_minutes": reminder.IntervalInMinutes,
				"updated_at":          reminder.UpdatedAt,
			})
		}

		if response == nil {
			return c.Status(fiber.StatusOK).SendString("Reminders not found")
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func (rc *reminderController) GetReminderByUuid(uuid uuid.UUID, c *fiber.Ctx) error {
	var reminder models.Reminder
	result := rc.db.UseGorm().First(&reminder, uuid)

	if result.Error != nil {
		return result.Error
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"reminder_id":         reminder.ReminderId,
			"reminder":            reminder.Reminder,
			"task_id":             reminder.TaskId,
			"description":         reminder.Description,
			"start_time":          reminder.StartTime,
			"frequency":           reminder.Frequency,
			"repeat_days":         reminder.RepeatDays,
			"repeat_sameday":      reminder.RepeatSameday,
			"repeat_until":        reminder.RepeatUntil,
			"interval":            reminder.Interval,
			"interval_in_minutes": reminder.IntervalInMinutes,
			"updated_at":          reminder.UpdatedAt,
		})
	}
}

func (rc *reminderController) GetReminderByTaskUuid(uuid uuid.UUID, c *fiber.Ctx) error {
	var reminder models.Reminder
	result := rc.db.UseGorm().Where("task_id = ?", uuid).First(&reminder)

	if result.Error != nil {
		return result.Error
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"reminder_id":         reminder.ReminderId,
			"reminder":            reminder.Reminder,
			"task_id":             reminder.TaskId,
			"description":         reminder.Description,
			"start_time":          reminder.StartTime,
			"frequency":           reminder.Frequency,
			"repeat_days":         reminder.RepeatDays,
			"repeat_sameday":      reminder.RepeatSameday,
			"repeat_until":        reminder.RepeatUntil,
			"interval":            reminder.Interval,
			"interval_in_minutes": reminder.IntervalInMinutes,
			"updated_at":          reminder.UpdatedAt,
		})
	}
}

func (rc *reminderController) UpdateTask(c *fiber.Ctx) error {
	var reminder models.Reminder

	if err := c.BodyParser(&reminder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	result := rc.db.UseGorm().Where("reminder_id = ? and task_id = ?",
		reminder.ReminderId, reminder.TaskId).Save(&reminder)

	if result.Error != nil {
		return result.Error
	} else {
		message := fmt.Sprintf("Reminder %s (%v) for task (%v) updated",
			reminder.Reminder, reminder.ReminderId, reminder.TaskId)
		return c.Status(fiber.StatusCreated).SendString(message)
	}
}
