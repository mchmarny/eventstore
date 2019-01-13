package utils

import (
	"log"
	"os"
	"strconv"
	"strings"
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

// EnvVarAsBool returns true if set to one of the trueStrings values
// false else. EnvVarAsBool is case insensitive
func EnvVarAsBool(key string, fallbackValue bool) bool {
	val := strings.ToLower(MustGetEnv(key, strconv.FormatBool(fallbackValue)))
	return Contains(trueStrings, val)
}
