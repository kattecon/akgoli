package utils

import "crypto/subtle"

// ConstantTimeStringEquals reports whether the two strings x and y have equal
// contents. When the strings have the same length, execution time depends only
// on the length and not on the content. When the lengths differ, the function
// returns false immediately, which leaks the length difference through timing.
func ConstantTimeStringEquals(x, y string) bool {
	if len(x) != len(y) {
		return false
	}

	// Manual XOR loop instead of subtle.ConstantTimeCompare to avoid the
	// string-to-[]byte allocation that the stdlib function would require.
	//
	// XOR each byte pair and OR the results into a single accumulator. If any
	// byte differs, a nonzero bit propagates into v. The loop must visit every
	// byte without an early exit so that the running time reveals nothing about
	// where the first difference occurs. Replacing this loop with == or adding
	// a break would reintroduce a timing side channel.
	var v byte

	for i := 0; i < len(x); i++ {
		v |= x[i] ^ y[i]
	}

	return subtle.ConstantTimeByteEq(v, 0) == 1
}
