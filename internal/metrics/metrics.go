package metrics

import (
	
	"log/slog"
	"net/http"

	
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)




type Metrics struct {
	RequestsTotal *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	DBQueries *prometheus.CounterVec
	BugsGauge prometheus.Gauge

}








func NewMetrics(reg prometheus.Registerer) *Metrics {
	logger := slog.Default().With(
		"metrics handler", "NewMetricaHandler",
	)
	logger.Info("entered NewMetrics")
    m := &Metrics{
        RequestsTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: "bugby",
                Name:      "http_requests_total",
                Help:      "Total number of HTTP requests",
            },
            []string{"method", "path", "status"},
        ),
        RequestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Namespace: "bugby",
                Name:      "http_request_duration_seconds",
                Help:      "Duration of HTTP requests",
                Buckets:   prometheus.DefBuckets,
            },
            []string{"method", "path"},
        ),
        DBQueries: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: "bugby",
                Name:      "db_queries_total",
                Help:      "Total number of DB queries",
            },
            []string{"query", "status"},
        ),
        BugsGauge: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Namespace: "bugby",
                Name:      "bugs_total",
                Help:      "Current number of bugs in system",
            },
        ),
    }

    reg.MustRegister(m.RequestsTotal, m.RequestDuration, m.DBQueries, m.BugsGauge)
    return m
}

func MetricsHandler() http.Handler {
    return promhttp.Handler()
}
