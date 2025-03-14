package utils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kevinhartarto/mytodolist/internal/models"
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
