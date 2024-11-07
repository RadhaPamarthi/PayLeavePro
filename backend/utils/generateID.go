// utils/generateID.go
package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GenerateEmployeeID creates an employee ID with the format: first 2 letters of first and last name + 5 random digits
func GenerateEmployeeID(firstName string, lastName string) string {
	if len(firstName) < 2 || len(lastName) < 2 {
		return ""
	}

	// Take first 2 letters from first and last name
	prefix := strings.ToUpper(firstName[:2] + lastName[:2])

	// Seed and generate 5 random digits
	rand.Seed(time.Now().UnixNano())
	randomNumbers := rand.Intn(100000)

	// Format the ID
	id := fmt.Sprintf("%s%05d", prefix, randomNumbers)
	return id
}
