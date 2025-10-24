package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/config"
)

// MetricCollector allows each module to self-register its metrics.
type MetricCollector interface {
	Register(r *prometheus.Registry)
}

type MetricsService interface {
	Register(collectors ...MetricCollector)
	RegisterDefault()
	Handler() http.Handler
}

type metricsService struct {
	Config   *config.Metrics
	Registry *prometheus.Registry
}

func NewMetricsService(cfg *config.Metrics, collectors ...MetricCollector) MetricsService {
	registry := prometheus.NewRegistry()

	m := &metricsService{
		Registry: registry,
		Config:   cfg,
	}

	if cfg.EnableDefaultMetrics {
		m.RegisterDefault()
	}
	m.Register(collectors...)

	return m
}

func (m *metricsService) Register(collectors ...MetricCollector) {
	for _, c := range collectors {
		c.Register(m.Registry)
	}
}

func (m *metricsService) RegisterDefault() {
	m.Registry.MustRegister(
		promcollectors.NewGoCollector(),
		promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
	)
}

func (m *metricsService) Handler() http.Handler {
	return promhttp.HandlerFor(
		m.Registry,
		promhttp.HandlerOpts{},
	)
}
