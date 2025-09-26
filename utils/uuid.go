package utils

import (
	"github.com/google/uuid"
)

func GenerateUUIDv7() string {
	id, _ := uuid.NewV7()
	return id.String()
}
