package testutils

import (
	"bytes"
	"io"
	"os"
	"strings"
)

func CapturePanicValue(f func()) (recovered interface{}) {
	defer func() {
		recovered = recover()
	}()

	f()

	return
}

func CaptureStderrNoDoubleQuotes(f func()) string {
	r, w, _ := os.Pipe()

	old := os.Stderr
	os.Stderr = w

	defer func() {
		os.Stderr = old
	}()

	f()

	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return strings.ReplaceAll(buf.String(), "\"", "'")
}

func CaptureStdoutNoDoubleQuotes(f func()) string {
	r, w, _ := os.Pipe()

	old := os.Stdout
	os.Stdout = w

	defer func() {
		os.Stdout = old
	}()

	f()

	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return strings.ReplaceAll(buf.String(), "\"", "'")
}
