package utils

// indexOfRunePos returns the byte offset of the boundary after the first pos
// runes in s. If pos exceeds the rune count, it returns len(s).
func indexOfRunePos(s string, pos int) int {
	for idx := range s {
		if pos == 0 {
			return idx
		}
		pos -= 1
	}
	return len(s)
}

// TakeRunes returns the first n runes of s as a string. If n is greater than
// the number of runes in s, it returns s unchanged. This function operates on
// rune boundaries, so it never splits a multi-byte UTF-8 character. The
// result shares storage with s, so retaining a short result keeps the full
// input string's allocation alive.
func TakeRunes(s string, n int) string {
	return s[:indexOfRunePos(s, n)]
}

// DropRunes returns s with the first n runes removed. If n is greater than
// the number of runes in s, it returns an empty string. This function
// operates on rune boundaries, so it never splits a multi-byte UTF-8
// character. The result shares storage with s.
func DropRunes(s string, n int) string {
	return s[indexOfRunePos(s, n):]
}
