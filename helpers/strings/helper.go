package strings

import "strconv"

// FirstString get first string of string slice
func FirstString(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

// StringSliceContainsTrue return true when src []string contains any true-bool-string
func StringSliceContainsTrue(src []string) bool {
	for _, v := range src {
		b, err := strconv.ParseBool(v)
		if err == nil && b {
			return true
		}
	}
	return false
}
