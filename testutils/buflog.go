package testutils

import (
	"bytes"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// BufferingLogger holds a Zap logger that writes JSON log entries into an
// in-memory buffer. It is intended for tests that need to assert on log output.
// Create instances with NewBufferingLogger; the zero value is not usable.
//
// The underlying buffer is not protected by a mutex, so the logger and buffer
// must be used from a single goroutine at a time. Do not log concurrently or
// read Buffer / call JsonNoDoubleQuotes while logging is in progress.
type BufferingLogger struct {
	Logger *zap.Logger
	Buffer *bytes.Buffer
}

// JsonNoDoubleQuotes returns the buffered log output with every double-quote
// character replaced by a single quote. The result is not valid JSON. It is
// an assertion-friendly string that avoids double-quote escaping in Go test
// literals. Do not pass the result to a JSON parser.
func (l BufferingLogger) JsonNoDoubleQuotes() string {
	return strings.ReplaceAll(l.Buffer.String(), "\"", "'")
}

// NewBufferingLogger creates a BufferingLogger that captures log entries at
// the given level or above. The logger produces JSON-encoded output with
// caller and timestamp fields omitted for stable test comparisons.
//
// The returned logger is not safe for concurrent use. See BufferingLogger.
func NewBufferingLogger(level zapcore.Level) BufferingLogger {
	b := &bytes.Buffer{}

	config := zap.NewProductionEncoderConfig()
	config.CallerKey = zapcore.OmitKey
	config.TimeKey = zapcore.OmitKey

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(b),
		level,
	)

	logger := zap.New(core)

	return BufferingLogger{logger, b}
}
