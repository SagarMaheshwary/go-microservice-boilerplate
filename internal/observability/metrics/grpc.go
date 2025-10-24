package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	GRPCRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests received, labeled by method and status.",
		},
		[]string{"method", "status"},
	)

	GRPCRequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Histogram of gRPC request latencies (seconds).",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

type GRPCMetrics struct{}

func (GRPCMetrics) Register(r *prometheus.Registry) {
	r.MustRegister(GRPCRequestCounter, GRPCRequestLatency)
}
