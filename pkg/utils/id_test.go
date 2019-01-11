package utils

import (
	"log"
	"strings"
	"testing"
)

func TestMakeID(t *testing.T) {

	testStr := "Test@Chmarny.com"

	id1 := MakeID(testStr)
	log.Printf("ID1: %s", id1)

	id2 := MakeID(strings.ToLower(testStr))
	log.Printf("ID2: %s", id2)

	if id1 != id2 {
		t.Errorf("Failed to generate case insensitive ID")
	}

}
