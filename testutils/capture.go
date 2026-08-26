package testutils

import (
	"bytes"
	"io"
	"os"
	"strings"
)

// CapturePanicValue calls f and returns the value passed to panic. If f does
// not panic, the returned value is nil.
func CapturePanicValue(f func()) (recovered any) {
	defer func() {
		recovered = recover()
	}()

	f()

	return
}

// CaptureStderrNoDoubleQuotes replaces os.Stderr with a pipe, calls f, reads
// everything f wrote to stderr, and returns the output with every double-quote
// character replaced by a single quote.
//
// The replacement is process-wide. Tests that use this function must not run
// concurrently with other code that reads or writes os.Stderr. If f writes
// more data than the OS pipe buffer can hold, f will block because the pipe
// is drained only after f returns. If f panics, os.Stderr is restored but no
// captured output is returned.
func CaptureStderrNoDoubleQuotes(f func()) string {
	r, w, _ := os.Pipe() // Pipe failure is catastrophic (fd exhaustion); nil w panics on next write.

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

// CaptureStdoutNoDoubleQuotes replaces os.Stdout with a pipe, calls f, reads
// everything f wrote to stdout, and returns the output with every double-quote
// character replaced by a single quote.
//
// The replacement is process-wide. Tests that use this function must not run
// concurrently with other code that reads or writes os.Stdout. If f writes
// more data than the OS pipe buffer can hold, f will block because the pipe
// is drained only after f returns. If f panics, os.Stdout is restored but no
// captured output is returned.
func CaptureStdoutNoDoubleQuotes(f func()) string {
	r, w, _ := os.Pipe() // Pipe failure is catastrophic (fd exhaustion); nil w panics on next write.

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
