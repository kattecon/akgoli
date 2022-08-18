package testutils

import (
	"bytes"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BufferingLogger struct {
	Logger *zap.Logger
	Buffer *bytes.Buffer
}

func (l BufferingLogger) JsonNoDoubleQuotes() string {
	return strings.ReplaceAll(l.Buffer.String(), "\"", "'")
}

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
