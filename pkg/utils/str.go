package utils

// Contains checks for presense of val in list
func Contains(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}
