package stringUtils

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
