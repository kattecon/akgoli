package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstErr(t *testing.T) {
	err := ConstError("xxx")
	assert.Equal(t, "xxx", err.Error())

	err = ConstError("yyy")
	assert.Equal(t, "yyy", err.Error())
}
