package testutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureStderrNoDoubleQuotes(t *testing.T) {
	old := os.Stderr

	s := CaptureStderrNoDoubleQuotes(func() {
		fmt.Fprint(os.Stderr, "hello")
	})

	assert.Equal(t, old, os.Stderr)
	assert.Equal(t, "hello", s)
}

func TestCaptureStdoutNoDoubleQuotes(t *testing.T) {
	old := os.Stdout

	s := CaptureStdoutNoDoubleQuotes(func() {
		fmt.Fprint(os.Stdout, "world")
	})

	assert.Equal(t, old, os.Stdout)
	assert.Equal(t, "world", s)
}

func TestCapturePanicValue(t *testing.T) {
	assert.Equal(t, nil, CapturePanicValue(func() {}))
	assert.Equal(t, "xxx", CapturePanicValue(func() { panic("xxx") }))
}
