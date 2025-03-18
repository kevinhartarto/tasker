package utils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kevinhartarto/tasker/internal/models"
)

func GenerateNewUUID() uuid.UUID {
	newUUID, err := uuid.NewUUID()

	if err != nil {
		fmt.Println("Failed to generate UUID")
	}

	return newUUID
}

func ParseUUID(value string) uuid.UUID {
	UUID, err := uuid.Parse(value)

	if err != nil {
		fmt.Println("Failed to parse UUID")
	}

	return UUID
}

func ValidateTask(task models.Task) bool {
	if task.Task == "" {
		return false
	}

	if task.TaskId == uuid.Nil {
		return false
	}

	return true
}

func ValidateReminder(reminder models.Reminder) bool {
	if reminder.Reminder == "" {
		return false
	}

	if reminder.ReminderId == uuid.Nil {
		return false
	}

	if reminder.TaskId == uuid.Nil {
		return false
	}

	if reminder.StartTime.IsZero() {
		return false
	}

	return true
}
