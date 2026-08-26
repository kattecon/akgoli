package utils

import "crypto/rand"

var secRndLetters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// GenSecureRandomId generates a random alphanumeric string of length n using
// crypto/rand as the random source. The output alphabet is a-z, A-Z, 0-9
// (62 characters).
//
// Each random byte is mapped to the alphabet with modulo 62. Because 256 is
// not evenly divisible by 62, the first 8 alphabet characters appear with
// probability 5/256 and the remaining 54 with probability 4/256. This bias is
// small but present. Callers that need a perfectly uniform distribution should
// use a different generator.
//
// Panics if n is negative or if crypto/rand fails. Passing n == 0 returns
// an empty string.
func GenSecureRandomId(n int) string {
	b := make([]byte, n)
	nr, err := rand.Read(b)
	if nr != n || err != nil {
		panic(err)
	}

	l := byte(len(secRndLetters))
	// Modulo mapping: each byte selects one of 62 alphabet characters. The
	// resulting distribution is slightly non-uniform (see doc comment above).
	for i, c := range b {
		b[i] = secRndLetters[c%l]
	}

	return string(b)
}
