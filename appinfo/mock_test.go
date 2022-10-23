package appinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	i := Mock()
	assert.Equal(t, "mock", i.AppIdName())
	assert.Equal(t, "1.2.3", i.AppVersion())
	assert.Equal(t, "100.500", i.GoVersion())
}
