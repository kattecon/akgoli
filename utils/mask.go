package utils

import (
	"strings"
	"unicode/utf8"
)

func MaskAll(s string) string {
	return strings.Repeat("*", utf8.RuneCountInString(s))
}
