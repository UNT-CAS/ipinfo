package ipinfo

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "job_duration_seconds",
			Help:    "Jobs duration distribution",
			Buckets: []float64{.100, .200, .300, .400, .500, 1.000},
		},
		[]string{"status"},
	)
)

func PrometheusInit() {
	prometheus.MustRegister(duration)
	http.Handle("/metrics", promhttp.Handler())
}
