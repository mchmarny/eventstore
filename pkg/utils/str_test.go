package utils

import (
	"testing"
)

func TestContainsSuccess(t *testing.T) {
	if !Contains([]string{"test1", "test2", "test3"}, "test2") {
		t.Errorf("Contains failed to find a valid item in array")
	}
}

func TestContainsFailure(t *testing.T) {
	if Contains([]string{"test1", "test2", "test3"}, "test4") {
		t.Errorf("Contains should have found an invalid item in array")
	}
}
