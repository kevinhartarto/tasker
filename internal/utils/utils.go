package utils

import "github.com/google/uuid"

func GenerateNewUUID() (uuid.UUID, error) {
	return uuid.NewUUID()
}
