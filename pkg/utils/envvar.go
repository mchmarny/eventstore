package utils

import (
	"log"
	"os"
)

var (
	trueStrings = []string{
		"true", "1", "yes", "tak",
	}
)

// MustGetEnv gets sets value or sets it to default when not set
func MustGetEnv(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	if fallbackValue == "" {
		log.Fatalf("Required env var (%s) not set", key)
	}

	return fallbackValue
}
