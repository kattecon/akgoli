package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskAll(t *testing.T) {
	assert.Equal(t, "", MaskAll(""))
	assert.Equal(t, "*", MaskAll("x"))
	assert.Equal(t, "**", MaskAll("xy"))
	assert.Equal(t, "*******", MaskAll("xyz-abc"))
	assert.Equal(t, "****", MaskAll("1абв"))
}
