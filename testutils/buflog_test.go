package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestDebugLevel(t *testing.T) {
	buflog := NewBufferingLogger(zapcore.DebugLevel)
	buflog.Logger.Info("x", zap.Int("a", 123))
	buflog.Logger.Debug("y", zap.Int("b", 333))
	assert.Equal(t, "{'level':'info','msg':'x','a':123}\n{'level':'debug','msg':'y','b':333}\n", buflog.JsonNoDoubleQuotes())
}

func TestInfoLevel(t *testing.T) {
	buflog := NewBufferingLogger(zapcore.InfoLevel)
	buflog.Logger.Info("x", zap.Int("a", 123))
	buflog.Logger.Debug("y", zap.Int("b", 333))
	assert.Equal(t, "{'level':'info','msg':'x','a':123}\n", buflog.JsonNoDoubleQuotes())
}
