package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTakeRunes(t *testing.T) {
	assert.Equal(t, "", TakeRunes("", 0))
	assert.Equal(t, "", TakeRunes("", 1))
	assert.Equal(t, "", TakeRunes("", 2))
	assert.Equal(t, "", TakeRunes("a", 0))
	assert.Equal(t, "a", TakeRunes("a", 1))
	assert.Equal(t, "a", TakeRunes("a", 2))
	assert.Equal(t, "", TakeRunes("ab", 0))
	assert.Equal(t, "a", TakeRunes("ab", 1))
	assert.Equal(t, "ab", TakeRunes("ab", 2))
	assert.Equal(t, "ab", TakeRunes("ab", 3))
	assert.Equal(t, "", TakeRunes("абв1", 0))
	assert.Equal(t, "а", TakeRunes("абв1", 1))
	assert.Equal(t, "аб", TakeRunes("абв1", 2))
	assert.Equal(t, "абв", TakeRunes("абв1", 3))
	assert.Equal(t, "абв1", TakeRunes("абв1", 4))
	assert.Equal(t, "абв1", TakeRunes("абв1", 5))
}

func TestDropRunes(t *testing.T) {
	assert.Equal(t, "", DropRunes("", 0))
	assert.Equal(t, "", DropRunes("", 1))
	assert.Equal(t, "", DropRunes("", 2))
	assert.Equal(t, "a", DropRunes("a", 0))
	assert.Equal(t, "", DropRunes("a", 1))
	assert.Equal(t, "", DropRunes("a", 2))
	assert.Equal(t, "ab", DropRunes("ab", 0))
	assert.Equal(t, "b", DropRunes("ab", 1))
	assert.Equal(t, "", DropRunes("ab", 2))
	assert.Equal(t, "", DropRunes("ab", 3))
	assert.Equal(t, "абв1", DropRunes("абв1", 0))
	assert.Equal(t, "бв1", DropRunes("абв1", 1))
	assert.Equal(t, "в1", DropRunes("абв1", 2))
	assert.Equal(t, "1", DropRunes("абв1", 3))
	assert.Equal(t, "", DropRunes("абв1", 4))
	assert.Equal(t, "", DropRunes("абв1", 5))
}
