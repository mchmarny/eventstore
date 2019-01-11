package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

const (
	idPrefix = "id-"
)

// MakeID normalizes input string and returns hashed value
func MakeID(val string) string {

	// normalize
	s := strings.TrimSpace(val)
	s = strings.ToLower(s)

	// sign
	hasher := md5.New()
	hasher.Write([]byte(s))

	return idPrefix + hex.EncodeToString(hasher.Sum(nil))

}
