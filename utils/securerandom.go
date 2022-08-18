package utils

import "crypto/rand"

var secRndLetters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenSecureRandomId(n int) string {
	b := make([]byte, n)
	nr, err := rand.Read(b)
	if nr != n || err != nil {
		panic(err)
	}

	l := byte(len(secRndLetters))
	for i, c := range b {
		b[i] = secRndLetters[c%l]
	}

	return string(b)
}
