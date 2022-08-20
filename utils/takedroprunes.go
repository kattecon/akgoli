package utils

func indexOfRunePos(s string, pos int) int {
	for idx := range s {
		if pos == 0 {
			return idx
		}
		pos -= 1
	}
	return len(s)
}

func TakeRunes(s string, n int) string {
	return s[:indexOfRunePos(s, n)]
}

func DropRunes(s string, n int) string {
	return s[indexOfRunePos(s, n):]
}
