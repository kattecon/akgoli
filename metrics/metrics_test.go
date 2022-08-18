package metrics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akshaal/akgoli/absos"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewMetrics(t *testing.T) {
	timeSvc := absos.NewFakeTimeSvc()
	timeSvc.Add(24 * 365 * 260 * time.Hour)
	m := NewMetrics(&MetricsParams{Prefix: "ff", AppVersion: "2.10.3"}, timeSvc)

	assert.Equal(t, "ff_xx", m.Prefixed("xx"))

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "my_test_counter321",
			Help: "my-test-desc",
		},
		[]string{"code", "method"},
	)
	m.MustRegister(counter)
	counter.WithLabelValues("404", "POST").Add(42)

	h := m.Handler()

	r := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()

	outB, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	out := string(outB)

	assert.Contains(t, out, "promhttp_metric_handler_errors_total")
	assert.Contains(t, out, "go_gc_duration_seconds")
	assert.Contains(t, out, "process_cpu_seconds_total")
	assert.Contains(t, out, "my_test_counter321{code=\"404\",method=\"POST\"} 42")
	assert.Contains(t, out, `ff_startup{version="2.10.3"} 1.403995421128655e+09`)
}

func TestNewMetricsWithoutDefaultCollectors(t *testing.T) {
	m := NewMetricsWithoutDefaultCollectors(&MetricsParams{Prefix: "mm", AppVersion: "2.10.3"})

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "my_test_counter321",
			Help: "my-test-desc",
		},
		[]string{"code", "method"},
	)
	m.MustRegister(counter)
	counter.WithLabelValues("404", "POST").Add(42)

	h := m.Handler()

	r := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()

	outB, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	out := string(outB)

	verify := func(name, s string) {
		t.Run(name, func(t *testing.T) {
			assert.Contains(t, out, "promhttp_metric_handler_errors_total")
			assert.NotContains(t, out, "go_build_info")
			assert.NotContains(t, out, "go_gc_cycles")
			assert.NotContains(t, out, "go_gc_duration_seconds")
			assert.NotContains(t, out, "process_cpu_seconds_total")
			assert.NotContains(t, out, "startup")
			assert.Contains(t, out, "my_test_counter321{code=\"404\",method=\"POST\"} 42")
		})
	}

	verify("via-handler", out)
	verify("via-dump", m.DumpAsTextForTest())
}

func TestJustDumpAsTextForTest(t *testing.T) {
	m := NewMetricsWithoutDefaultCollectors(&MetricsParams{Prefix: "mm", AppVersion: "2.10.3"})

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "my_test_counter321",
			Help: "my-test-desc",
		},
		[]string{"code", "method"},
	)
	m.MustRegister(counter)
	counter.WithLabelValues("404", "POST").Add(42)

	out := m.DumpAsTextForTest()

	assert.NotContains(t, out, "promhttp_metric_handler_errors_total")
	assert.Contains(t, out, "my_test_counter321{code=\"404\",method=\"POST\"} 42")
}
