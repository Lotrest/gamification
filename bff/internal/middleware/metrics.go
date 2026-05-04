package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	requestDuration *prometheus.HistogramVec
	requestsTotal   *prometheus.CounterVec
}

func NewMetrics(registry *prometheus.Registry) *Metrics {
	metrics := &Metrics{
		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "cdek_bff_request_duration_seconds",
			Help:    "HTTP request duration for the BFF service.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "route", "status"}),
		requestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "cdek_bff_requests_total",
			Help: "HTTP requests handled by the BFF service.",
		}, []string{"method", "route", "status"}),
	}

	registry.MustRegister(metrics.requestDuration, metrics.requestsTotal)
	return metrics
}

func (m *Metrics) Middleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		started := time.Now()
		err := ctx.Next()
		status := strconv.Itoa(ctx.Response().StatusCode())
		route := ctx.Route().Path

		m.requestsTotal.WithLabelValues(ctx.Method(), route, status).Inc()
		m.requestDuration.WithLabelValues(ctx.Method(), route, status).Observe(time.Since(started).Seconds())

		return err
	}
}
