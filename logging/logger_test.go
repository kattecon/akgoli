package logging

import (
	"testing"

	"github.com/akshaal/akgoli/metrics"
	"github.com/akshaal/akgoli/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewLoggerDefault(t *testing.T) {
	m := metrics.NewMetricsWithoutDefaultCollectors(&metrics.MetricsParams{Prefix: "xx"})
	out := testutils.CaptureStderrNoDoubleQuotes(func() {
		cfg := NewSimpleLoggerConfig()
		logger, err := NewLogger(cfg, m)
		logger.Debug("debug", zap.Error(err))
		logger.Info("info", zap.Error(err))
		logger.Error("error", zap.Error(err))
	})

	assert.Contains(t, out, "'msg':'error'")
	assert.Contains(t, out, "'msg':'info'")
	assert.NotContains(t, out, "'msg':'debug'")
	assert.NotContains(t, out, "akgoli", "must not contain stacktrace for errors")
	assert.Contains(t, m.DumpAsTextForTest(), `xx_log_events{level="error"} 1`)
	assert.Contains(t, m.DumpAsTextForTest(), `xx_log_events{level="info"} 1`)
	assert.Contains(t, m.DumpAsTextForTest(), `xx_log_events{level="debug"} 0`)
	assert.Contains(t, m.DumpAsTextForTest(), `xx_log_events{level="fatal"} 0`)
}

func TestNewLoggerDevStyle(t *testing.T) {
	m := metrics.NewMetricsWithoutDefaultCollectors(&metrics.MetricsParams{Prefix: "yy"})
	out := testutils.CaptureStderrNoDoubleQuotes(func() {
		cfg := NewSimpleLoggerConfig()
		cfg.SetDevStyleLogging(true)
		logger, err := NewLogger(cfg, m)
		logger.Debug("debug", zap.Error(err))
		logger.Info("info", zap.Error(err))
		logger.Error("error", zap.Error(err))
	})

	assert.Contains(t, out, "error")
	assert.Contains(t, out, "\tinfo\n")
	assert.NotContains(t, out, "debug")
	assert.NotContains(t, out, "akgoli", "must not contain stacktrace for errors")
	assert.Contains(t, m.DumpAsTextForTest(), `yy_log_events{level="error"} 1`)
	assert.Contains(t, m.DumpAsTextForTest(), `yy_log_events{level="info"} 1`)
	assert.Contains(t, m.DumpAsTextForTest(), `yy_log_events{level="debug"} 0`)
	assert.Contains(t, m.DumpAsTextForTest(), `yy_log_events{level="fatal"} 0`)
}

func TestNewLoggerDebugLogging(t *testing.T) {
	m := metrics.NewMetricsWithoutDefaultCollectors(&metrics.MetricsParams{Prefix: "zz"})
	out := testutils.CaptureStderrNoDoubleQuotes(func() {
		cfg := NewSimpleLoggerConfig()
		cfg.SetDebugLogging(true)
		logger, err := NewLogger(cfg, m)
		logger.Debug("debug", zap.Error(err))
		logger.Info("info", zap.Error(err))
		logger.Error("error", zap.Error(err))
	})

	assert.Contains(t, out, "'msg':'error'")
	assert.Contains(t, out, "'msg':'info'")
	assert.Contains(t, out, "'msg':'debug'")
	assert.NotContains(t, out, "akgoli", "must not contain stacktrace for errors")
	assert.Contains(t, m.DumpAsTextForTest(), `zz_log_events{level="error"} 1`)
	assert.Contains(t, m.DumpAsTextForTest(), `zz_log_events{level="info"} 1`)

	// '2' because there is also "Logged initialized message"
	assert.Contains(t, m.DumpAsTextForTest(), `zz_log_events{level="debug"} 2`)

	assert.Contains(t, m.DumpAsTextForTest(), `zz_log_events{level="fatal"} 0`)
}
