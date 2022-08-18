package logging

import (
	"github.com/akshaal/akgoli/metrics"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig interface {
	// Whether to output debugging log messages or not.
	IsDebugLogging() bool

	// Outputs json formatted log messes if true. Otherwise outputs human readable ones.
	IsDevStyleLogging() bool
}

type SimpleLoggerConfigImpl struct {
	debugLogging    bool
	devStyleLogging bool
}

// This one is mainly for unit tests... usually there is a config instance based upon environment variables, some other logic..
func NewSimpleLoggerConfig() *SimpleLoggerConfigImpl {
	return &SimpleLoggerConfigImpl{}
}

func (c *SimpleLoggerConfigImpl) IsDebugLogging() bool {
	return c.debugLogging
}

func (c *SimpleLoggerConfigImpl) SetDebugLogging(v bool) {
	c.debugLogging = v
}

func (c *SimpleLoggerConfigImpl) IsDevStyleLogging() bool {
	return c.devStyleLogging
}

func (c *SimpleLoggerConfigImpl) SetDevStyleLogging(v bool) {
	c.devStyleLogging = v
}

func NewLogger(cfg LoggerConfig, m *metrics.Metrics) (*zap.Logger, error) {
	// Counter for log events
	logEventsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: m.Prefixed("log_events"),
			Help: "Total number of log events logged.",
		},
		[]string{"level"},
	)
	m.MustRegister(logEventsCounter)

	// Set initial counter to zero for each log level
	// (grafana/prometheus might otherwise have problems doing math and distinguishing between missing/zero values...)
	logEventsCounter.WithLabelValues(zap.DebugLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.InfoLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.WarnLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.ErrorLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.DPanicLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.PanicLevel.String()).Add(0)
	logEventsCounter.WithLabelValues(zap.FatalLevel.String()).Add(0)

	// Logger itself...

	var zapConfig zap.Config
	if cfg.IsDevStyleLogging() {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

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
