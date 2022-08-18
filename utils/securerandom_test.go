package utils

import (
	"testing"

	"github.com/agnivade/levenshtein"
	"github.com/stretchr/testify/assert"
)

func TestGenSecureRandomId(t *testing.T) {
	for i := 0; i < 1000; i++ {
		for _, l := range []int{80, 100, 200} {
			s1 := GenSecureRandomId(l)
			s2 := GenSecureRandomId(l)
			assert.NotEqual(t, s1, s2)
			assert.Equal(t, l, len(s1))
			assert.Equal(t, l, len(s2))

			dist := float64(levenshtein.ComputeDistance(s1, s2)) / float64(l)
			assert.Greater(t, dist, 0.8)
		}
	}
}
