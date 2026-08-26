// Package logging provides a Zap logger factory that connects log events to
// Prometheus metrics. NewLogger is the main entry point.
package logging

import (
	"github.com/kattecon/akgoli/metrics"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig controls the output format and verbosity of a logger created
// by NewLogger. The settings are read once during logger construction. Changing
// them afterward does not affect loggers that have already been built.
type LoggerConfig interface {
	// IsDebugLogging reports whether the logger should emit debug-level
	// messages. When false, the logger starts at info level.
	IsDebugLogging() bool

	// IsDevStyleLogging reports whether the logger should use Zap's
	// human-readable development format. When false, the logger uses Zap's
	// JSON production format.
	IsDevStyleLogging() bool
}

// SimpleLoggerConfigImpl is a mutable LoggerConfig with setter methods.
// Both fields default to false, which selects info-level JSON output.
type SimpleLoggerConfigImpl struct {
	debugLogging    bool
	devStyleLogging bool
}

// NewSimpleLoggerConfig creates a SimpleLoggerConfigImpl with both settings
// set to false (info-level JSON output).
func NewSimpleLoggerConfig() *SimpleLoggerConfigImpl {
	return &SimpleLoggerConfigImpl{}
}

func (c *SimpleLoggerConfigImpl) IsDebugLogging() bool {
	return c.debugLogging
}

// SetDebugLogging controls whether the logger emits debug-level messages.
func (c *SimpleLoggerConfigImpl) SetDebugLogging(v bool) {
	c.debugLogging = v
}

func (c *SimpleLoggerConfigImpl) IsDevStyleLogging() bool {
	return c.devStyleLogging
}

// SetDevStyleLogging controls whether the logger uses human-readable development format.
func (c *SimpleLoggerConfigImpl) SetDevStyleLogging(v bool) {
	c.devStyleLogging = v
}

// NewLogger creates a Zap logger and registers a Prometheus counter named
// <app-id>_log_events that counts log entries by level. Both cfg and m must
// be non-nil.
//
// The counter is registered with m.MustRegister, which panics if the same
// counter has already been registered in that Metrics instance. Call
// NewLogger only once per Metrics registry.
//
// When debug logging is enabled, the logger emits one "Logger initialized"
// entry on success and increments the debug counter. When debug logging is
// off, no entry is emitted and the debug counter stays at zero.
// Stack traces are disabled for all log levels. Both modes write to stderr.
// The function returns
// an error only when Zap's own configuration build fails.
func NewLogger(cfg LoggerConfig, m *metrics.Metrics) (*zap.Logger, error) {
	logEventsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: m.Prefixed("log_events"),
			Help: "Total number of log events logged.",
		},
		[]string{"level"},
	)
	m.MustRegister(logEventsCounter)

	// Pre-initialize every level series to zero. Prometheus and Grafana treat
	// a missing series differently from a zero-valued one, so PromQL rate()
	// and increase() return correct results from the start.
	logEventsCounter.WithLabelValues(zap.DebugLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.InfoLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.WarnLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.ErrorLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.DPanicLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.PanicLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.FatalLevel.String()).Add(0)

	var zapConfig zap.Config
	if cfg.IsDevStyleLogging() {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	// Stack traces add noise to structured logs in this library's use cases.
	zapConfig.DisableStacktrace = true

	if cfg.IsDebugLogging() {
		zapConfig.Level.SetLevel(zap.DebugLevel)
	} else {
		zapConfig.Level.SetLevel(zap.InfoLevel)
	}

	metricsHook := func(e zapcore.Entry) error {
		logEventsCounter.WithLabelValues(e.Level.String()).Inc()
		return nil
	}

	logger, err := zapConfig.Build(zap.Hooks(metricsHook))
	if err != nil {
		return nil, errors.Wrap(err, "unable to build zap logger")
	}

	logger.Debug("Logger initialized")

	return logger, nil
}
