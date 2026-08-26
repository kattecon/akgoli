package utils

import (
	"strings"
	"unicode/utf8"
)

// MaskAll replaces every rune in s with an asterisk. The returned string has
// the same rune count as the input, so the caller should be aware that the
// length of the original value is preserved and visible.
// An empty string returns an empty string.
// See also Masked for use as a zap.Field in structured logs.
func MaskAll(s string) string {
	return strings.Repeat("*", utf8.RuneCountInString(s))
}
