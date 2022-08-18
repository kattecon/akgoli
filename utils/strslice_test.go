package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveStrFromSliceInPlace(t *testing.T) {
	slice1 := []string{"a", "x", "x", "1", "b"}
	slice2 := RemoveStrFromSliceInPlace(slice1, "x")
	assert.Equal(t, []string{"a", "1", "b", "1", "b"}, slice1)
	assert.Equal(t, []string{"a", "1", "b"}, slice2)

	slice3 := RemoveStrFromSliceInPlace(slice2, "a")
	assert.Equal(t, []string{"1", "b"}, slice3)

	slice4 := RemoveStrFromSliceInPlace(slice3, "1")
	assert.Equal(t, []string{"b"}, slice4)

	slice5 := RemoveStrFromSliceInPlace(slice4, "b")
	assert.Equal(t, []string{}, slice5)

	slice6 := RemoveStrFromSliceInPlace(slice5, "z")
	assert.Equal(t, []string{}, slice6)
}

func TestTrimAllInPlace(t *testing.T) {
	slice := []string{"   ", " x "}

	assert.Equal(t, []string{"", "x"}, TrimAllInPlace(slice))
	assert.Equal(t, []string{"", "x"}, slice)
}
