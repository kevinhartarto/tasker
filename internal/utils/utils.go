package utils

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kevinhartarto/tasker/internal/logger"
	"github.com/kevinhartarto/tasker/internal/models"
)

var log = logger.GetLogger()

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

	if reminder.Frequency == "" {
		return false
	} else {
		switch reminder.Frequency {
		// Same day reminder
		case "n":
			if reminder.RepeatSameday {
				// need to have interval in minutes
				// or must not more than 1 day (1440 minutes)
				if reminder.Interval != nil ||
					reminder.IntervalInMinutes == nil ||
					*reminder.IntervalInMinutes >= 1440 {
					return false
				}
				// check next reminder
				return ValidateNextReminder(reminder)
			}
		case "d", "w", "m", "y":
			if !validateFrequenctTypeRepeats(reminder) {
				return false
			}
			// check next reminder
			return ValidateNextReminder(reminder)
		case "s":
			days := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
			valid := false
			for _, day := range reminder.RepeatDays {
				valid = slices.Contains(days, day)
			}

			if !validateFrequenctTypeRepeats(reminder) || !valid {
				return false
			}
			// check next reminder
			return ValidateNextReminder(reminder)
		default:
			break
		}
	}
	return true
}

func validateFrequenctTypeRepeats(reminder models.Reminder) bool {
	if reminder.Interval == nil ||
		reminder.RepeatSameday ||
		reminder.IntervalInMinutes != nil ||
		reminder.RepeatUntil.IsZero() ||
		reminder.NextReminder.IsZero() {
		return false
	}
	return true
}

func ValidateNextReminder(reminder models.Reminder) bool {
	expectedDate := reminder.StartTime

	switch reminder.Frequency {
	case "n":
		expectedDate = expectedDate.Add(time.Minute * time.Duration(*reminder.IntervalInMinutes))
	case "d", "w", "m", "y":
		expectedDate = expectedDate.Add(time.Duration(*reminder.Interval))
	case "s":
		// get next reminder week day
		nextDateWeekday := strings.ToLower(reminder.NextReminder.Weekday().String())

		// compare week day with repeat days
		if !slices.Contains(reminder.RepeatDays, nextDateWeekday[0:3]) &&
			(reminder.NextReminder.After(*reminder.RepeatUntil)) {
			return false
		}

		return true
	}

	return (reminder.NextReminder.Equal(expectedDate)) && (reminder.NextReminder.Before(*reminder.RepeatUntil))
}

func SendDesktopNotification(level string, title string, msgBody string) {
	switch strings.ToLower(level) {
	case "notify":
		if err := beeep.Notify(title, msgBody, "resources/assets/normal_notification.jpg"); err != nil {
			log.Info("Failed to send notification (" + title + ")")
		}
	case "alert":
		if err := beeep.Alert(title, msgBody, "resources/assets/normal_notification.jpg"); err != nil {
			log.Info("Failed to send alert (" + title + ")")
		}
	default:
		log.Info("Unknown dekstop notification level, please use either 'notify' or 'alert' !")
	}
}
