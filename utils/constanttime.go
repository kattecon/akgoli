package utils

import "crypto/subtle"

// Returns true iff the two string, x and y, have equal contents. Time depends on the length of the string.
// If length are different, then returns immediately.
func ConstantTimeStringEquals(x, y string) bool {
	if len(x) != len(y) {
		return false
	}

	var v byte

	for i := 0; i < len(x); i++ {
		v |= x[i] ^ y[i]
	}

	return subtle.ConstantTimeByteEq(v, 0) == 1
}
