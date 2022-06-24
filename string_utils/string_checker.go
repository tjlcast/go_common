package string_utils

import "strings"

func EmptyString(str string) bool {
	str = strings.Trim(str, " ")
	str = strings.Trim(str, "\"")
	if len(str) == 0 {
		return true
	}
	return false
}