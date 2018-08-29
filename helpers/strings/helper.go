package strings

// FirstString get first string of string slice
func FirstString(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}
