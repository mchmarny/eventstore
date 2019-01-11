package utils

import (
	"log"
	"strings"
	"testing"
)

func TestID(t *testing.T) {

	testEmail := "Test@Chmarny.com"

	id1 := MakeID(testEmail)
	log.Printf("ID1: %s", id1)

	id2 := MakeID(strings.ToLower(testEmail))
	log.Printf("ID2: %s", id2)

	if id1 != id2 {
		t.Errorf("Failed to generate case insensitive ID")
	}

	email, err := ParseEmail(id2)
	if err != nil {
		t.Errorf("Error parsing email: %v", err)
	}

	if email != id2 {
		t.Errorf("Failed to generate case insensitive ID")
	}

}
