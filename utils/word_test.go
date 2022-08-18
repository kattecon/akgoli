package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAsciiWord(t *testing.T) {
	assert.True(t, IsAsciiWord("hello"))
	assert.True(t, IsAsciiWord("h"))
	assert.True(t, IsAsciiWord("Hello"))
	assert.True(t, IsAsciiWord("helL"))
	assert.False(t, IsAsciiWord("нет"))
	assert.False(t, IsAsciiWord("ø"))
	assert.False(t, IsAsciiWord("123"))
	assert.False(t, IsAsciiWord(""))
	assert.False(t, IsAsciiWord("eee2ee"))
	assert.False(t, IsAsciiWord("eee2"))
	assert.False(t, IsAsciiWord("1"))
	assert.False(t, IsAsciiWord("d2"))
	assert.False(t, IsAsciiWord("2d"))
	assert.False(t, IsAsciiWord(" "))
	assert.False(t, IsAsciiWord("x/"))
}

func TestIsAsciiWordWithDigits(t *testing.T) {
	assert.True(t, IsAsciiWordWithDigits("hello"))
	assert.True(t, IsAsciiWordWithDigits("h"))
	assert.True(t, IsAsciiWordWithDigits("Hello"))
	assert.True(t, IsAsciiWordWithDigits("helL"))
	assert.False(t, IsAsciiWordWithDigits("нет"))
	assert.False(t, IsAsciiWordWithDigits("ø"))
	assert.False(t, IsAsciiWordWithDigits("123"))
	assert.False(t, IsAsciiWordWithDigits(""))
	assert.True(t, IsAsciiWordWithDigits("eee2ee"))
	assert.True(t, IsAsciiWordWithDigits("eee2"))
	assert.False(t, IsAsciiWordWithDigits("1"))
	assert.True(t, IsAsciiWordWithDigits("d2"))
	assert.False(t, IsAsciiWordWithDigits("2d"))
	assert.False(t, IsAsciiWordWithDigits(" "))
	assert.False(t, IsAsciiWordWithDigits("x/"))
}
