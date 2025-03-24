package utils

import (
	"fmt"
	"os"
	"slices"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kevinhartarto/tasker/internal/models"
)

func GetEnvOrDefault(envName string, defaultValue string) any {
	if envValue := os.Getenv(envName); envValue != "" {
		return envValue
	}

	return defaultValue
}

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

	return ValidateReminderFrequencyAndNextReminder(reminder)
}

func ValidateReminderFrequencyAndNextReminder(reminder models.Reminder) bool {

	if reminder.Frequency == "" {
		return false
	}

	days := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
	switch reminder.Frequency {
	case "n":
	case "d", "w", "m", "y":
		if reminder.Interval == nil {
			return false
		}
	case "s":
		for _, day := range reminder.RepeatDays {
			if valid := slices.Contains(days, day); !valid {
				return false
			}
		}
	}

	return true
}
