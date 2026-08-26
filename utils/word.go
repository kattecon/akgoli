package utils

// IsAsciiWord reports whether s contains only ASCII letters (a-z, A-Z).
// It returns false for an empty string, digits, Unicode letters, spaces,
// and any other non-ASCII-letter byte.
func IsAsciiWord(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			return false
		}
	}

	return true
}

// IsAsciiWordWithDigits reports whether s starts with an ASCII letter and
// contains only ASCII letters and digits after the first character.
// It returns false for an empty string, a string starting with a digit,
// Unicode letters, and any other non-ASCII-alphanumeric byte.
func IsAsciiWordWithDigits(s string) bool {
	if s == "" {
		return false
	}

	// Letters pass at any position. Digits pass after the first character.
	// Everything else is rejected.
	for i, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			if !(i != 0 && c >= '0' && c <= '9') {
				return false
			}
		}
	}

	return true
}
