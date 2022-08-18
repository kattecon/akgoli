package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	assert.True(t, ConstantTimeStringEquals("abc", "abc"))
	assert.True(t, ConstantTimeStringEquals("ac", "ac"))
	assert.True(t, ConstantTimeStringEquals("", ""))
	assert.False(t, ConstantTimeStringEquals("123", ""))
	assert.False(t, ConstantTimeStringEquals("ac", ""))
	assert.False(t, ConstantTimeStringEquals("bc", "abc"))
	assert.False(t, ConstantTimeStringEquals("abc", "bca"))
	assert.False(t, ConstantTimeStringEquals("Abc", "abc"))
	assert.False(t, ConstantTimeStringEquals("abc", "abc "))
	assert.False(t, ConstantTimeStringEquals(" abc", "abc"))
	assert.False(t, ConstantTimeStringEquals("abC", "abc"))
}
