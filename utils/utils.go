package utils

import (
	"fmt"

	"github.com/google/uuid"
)

// Generate a unique manga ID (you can implement your own logic here)
func GenerateUUID() (u uuid.UUID) {
	u, err := uuid.NewUUID()
	if err != nil {
		fmt.Printf("Error generating UUID: %v\n", err)
		return
	}

	return
}
