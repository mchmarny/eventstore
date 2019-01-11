package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
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

	// decode
	decoded, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return "", err
	}

	e := string(decoded)

	// check format
	if !strings.HasPrefix(e, idPrefix) {
		return "", fmt.Errorf("Invalid ID: %s", e)
	}

	return strings.TrimPrefix(e, idPrefix), nil

}
