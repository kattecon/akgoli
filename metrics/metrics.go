// Package metrics provides a private Prometheus registry with application-aware
// metric naming. All metric names are prefixed with the application's ID name
// from appinfo via the Prefixed method.
package metrics

import (
	"bytes"
	"net/http"

	"github.com/kattecon/akgoli/absos"
	"github.com/kattecon/akgoli/appinfo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
)

// Metrics wraps a private Prometheus registry. Using a private registry
// instead of the global default avoids collisions with other libraries and
// gives full control over which collectors are registered. Create instances
// with NewMetrics or NewMetricsWithoutDefaultCollectors; the zero value is
// not usable.
type Metrics struct {
	reg     *prometheus.Registry
	appInfo appinfo.AppInfo
}

// NewMetricsWithoutDefaultCollectors creates a Metrics with an empty registry.
// No Go runtime or process collectors are registered. Use this in tests to
// avoid process-collector noise. The first call to Handler on the returned
// Metrics still registers promhttp handler-instrumentation metrics.
func NewMetricsWithoutDefaultCollectors(appInfo appinfo.AppInfo) *Metrics {
	return &Metrics{
		reg:     prometheus.NewRegistry(),
		appInfo: appInfo,
	}
}

// NewMetrics creates a Metrics with Go runtime collectors, process collectors,
// and a startup gauge named <app-id>_startup. The startup gauge records the
// time of creation as fractional Unix seconds (derived from UnixNano),
// labeled with the application version from appInfo. Both appInfo and
// timeSvc must be non-nil.
func NewMetrics(appInfo appinfo.AppInfo, timeSvc absos.TimeSvc) *Metrics {
	m := NewMetricsWithoutDefaultCollectors(appInfo)

	m.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	startupGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: m.Prefixed("startup"),
			Help: "Startup time.",
		},
		[]string{"version"},
	)
	m.MustRegister(startupGauge)
	// Store startup time as fractional Unix seconds. Prometheus convention
	// uses seconds as the base unit for timestamps.
	startupGauge.WithLabelValues(appInfo.AppVersion()).Set(float64(timeSvc.Now().UnixNano()) / 1e9)

	return m
}

// Prefixed returns name prefixed with the application ID and an underscore.
// For example, if AppIdName returns "myapp", Prefixed("requests") returns
// "myapp_requests". No validation is performed on name.
func (m *Metrics) Prefixed(name string) string {
	return m.appInfo.AppIdName() + "_" + name
}

// MustRegister registers the given collectors with the private registry.
// It panics if any collector is invalid or has already been registered.
func (m *Metrics) MustRegister(cs ...prometheus.Collector) {
	m.reg.MustRegister(cs...)
}

// Handler returns an http.Handler that serves the collected metrics in
// Prometheus text exposition format. The handler also registers its own
// promhttp_metric_handler_* instrumentation metrics.
//
// Response compression is disabled to keep scrape output deterministic.
// At most 10 concurrent scrape requests are served. Requests beyond that
// limit receive HTTP 503. Gather errors produce HTTP 500 and increment the
// promhttp_metric_handler_errors_total counter.
func (m *Metrics) Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		m.reg,
		promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{
			Registry: m.reg,
			// Keep scrape output byte-identical across requests to simplify
			// debugging and test assertions.
			DisableCompression: true,
			// Limit concurrent scrapes to protect memory under load.
			MaxRequestsInFlight: 10,
		}),
	)
}

// DumpAsTextForTest gathers all registered metrics and returns them in
// Prometheus text exposition format. This method is intended for test
// assertions. It returns an empty string when no metrics are registered.
// Metric families and samples are sorted in the registry's canonical order.
// It panics if gathering or encoding fails.
func (m *Metrics) DumpAsTextForTest() string {
	mfs, err := m.reg.Gather()
	if err != nil {
		panic(err)
	}

	b := &bytes.Buffer{}

	for _, mf := range mfs {
		if _, err := expfmt.MetricFamilyToText(b, mf); err != nil {
			panic(err)
		}
	}

	return b.String()
}
