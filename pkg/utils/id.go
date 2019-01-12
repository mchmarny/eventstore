package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

const (
	idPrefix = "id-"
)

// MakeID normalizes input string and returns encoded email
func MakeID(email string) string {

	// normalize
	s := strings.TrimSpace(email)
	s = strings.ToLower(s)

	// encode
	encoded := base64.StdEncoding.EncodeToString([]byte(s))

	return idPrefix + encoded

}

// ParseEmail parses email from the encoded id
func ParseEmail(id string) (email string, err error) {

	// check format
	if !strings.HasPrefix(id, idPrefix) {
		return "", fmt.Errorf("Invalid ID format: %s", id)
	}

	// trim
	id2 := strings.TrimPrefix(id, idPrefix)

	// decode
	decoded, err := base64.StdEncoding.DecodeString(id2)
	if err != nil {
		return "", err
	}

	return string(decoded), nil

}

// MakeUUID makes UUID string
func MakeUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Error while getting id: %v\n", err)
	}
	return id.String()
}
