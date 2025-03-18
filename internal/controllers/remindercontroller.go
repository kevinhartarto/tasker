package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kevinhartarto/tasker/internal/database"
	"github.com/kevinhartarto/tasker/internal/models"
	"github.com/kevinhartarto/tasker/internal/utils"
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

	newReminder.ReminderId = utils.GenerateNewUUID()
	if newReminder.StartTime.IsZero() {
		newReminder.StartTime = time.Now()
	}

	if !utils.ValidateReminder(newReminder) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Reminder name, id, start time and task id are required",
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

func (rc *reminderController) UpdateRemainder(c *fiber.Ctx) error {
	var reminder models.Reminder
	var data map[string]interface{}

	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	rc.db.UseGorm().Model(&reminder).Where("reminder_id = ?", data["reminder_id"]).Updates(data)
	result := rc.db.UseGorm().Where("reminder_id = ?", data["reminder_id"]).First(&reminder)

	if result.Error != nil {
		return result.Error
	} else {
		message := fmt.Sprintf("Reminder %s (%v) for task (%v) updated",
			reminder.Reminder, reminder.ReminderId, reminder.TaskId)
		return c.Status(fiber.StatusCreated).SendString(message)
	}
}
