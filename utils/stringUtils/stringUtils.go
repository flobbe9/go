package stringUtils

import "strings"

func StartsWith(str, subStr string) bool {
	// case: subStr cannot be a sub string of str
	if len(subStr) > len(str) {
		return false
	}

	return str[:len(subStr)] == subStr
}

func EndsWith(str, subStr string) bool {
	// case: subStr cannot be a sub string of str
	if len(subStr) > len(str) {
		return false
	}

	return str[len(str)-len(subStr):] == subStr
}

// Remove all leading and trailing forward slashes
func Unslash(str string) string {
	if len(str) == 0 {
		return str
	}

	for {
		if !StartsWith(str, "/") && !EndsWith(str, "/") {
			break
		}

		if StartsWith(str, "/") {
			str = str[1:]
		}

		if EndsWith(str, "/") {
			str = str[:len(str)-1]
		}
	}

	return str
}

// Indicates that [str] is empty or consists only of white space
func IsBlank(str string) bool {
	return len(strings.Trim(str, " ")) == 0
}