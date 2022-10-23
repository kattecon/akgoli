package appinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	i := Get()
	assert.Equal(t, "test", i.AppIdName())
	assert.Equal(t, "3.2.1", i.AppVersion())
	assert.NotEmpty(t, i.GoVersion())
}

func withSavedValues(f func()) {
	origVersion := version
	origGoVersion := goVersion
	origIdName := idName

	defer func() {
		version = origVersion
		goVersion = origGoVersion
		idName = origIdName
	}()

	f()
}

func Test2(t *testing.T) {
	i := Get()

	withSavedValues(func() {
		version = ""
		goVersion = ""
		idName = ""

		assert.Equal(t, "unknown", i.AppIdName())
		assert.Equal(t, "unknown", i.AppVersion())
		assert.Equal(t, "unknown", i.GoVersion())
	})
}

func Test3(t *testing.T) {
	i := Get()

	withSavedValues(func() {
		version = "a"
		goVersion = "b"
		idName = "c"

		assert.Equal(t, "c", i.AppIdName())
		assert.Equal(t, "a", i.AppVersion())
		assert.Equal(t, "b", i.GoVersion())
	})
}
