package hud

import "strings"

// Trim will remove leading and trailing empty line
func Trim(s string) string {
	s = strings.TrimPrefix(s, "\n")
	s = strings.TrimSuffix(s, "\n")
	return s
}