package utils

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

func IsAsciiWordWithDigits(s string) bool {
	if s == "" {
		return false
	}

	for i, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			if !(i != 0 && c >= '0' && c <= '9') {
				return false
			}
		}
	}

	return true
}
