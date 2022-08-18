package utils

import "strings"

func MaskAll(s string) string {
	return strings.Repeat("*", len(s))
}
