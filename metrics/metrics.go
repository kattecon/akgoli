package metrics

import (
	"bytes"
	"net/http"

	"github.com/akshaal/akgoli/absos"
	"github.com/akshaal/akgoli/appinfo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
)

type Metrics struct {
	reg     *prometheus.Registry
	appInfo appinfo.AppInfo
}

func NewMetricsWithoutDefaultCollectors(appInfo appinfo.AppInfo) *Metrics {
	return &Metrics{
		reg:     prometheus.NewRegistry(),
		appInfo: appInfo,
	}
}

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
	startupGauge.WithLabelValues(appInfo.AppVersion()).Set(float64(timeSvc.Now().UnixNano()) / 1e9)

	return m
}

func (m *Metrics) Prefixed(name string) string {
	return m.appInfo.AppIdName() + "_" + name
}

func (m *Metrics) MustRegister(cs ...prometheus.Collector) {
	m.reg.MustRegister(cs...)
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		m.reg,
		promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{
			Registry:            m.reg,
			DisableCompression:  true,
			MaxRequestsInFlight: 10,
		}),
	)
}

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
