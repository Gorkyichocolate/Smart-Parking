// pkg/metrics/metrics.go
package metrics

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	// HTTP метрики
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	RequestErrors   *prometheus.CounterVec

	// gRPC метрики
	GrpcRequestsTotal   *prometheus.CounterVec
	GrpcRequestDuration *prometheus.HistogramVec
	GrpcRequestErrors   *prometheus.CounterVec

	// Database метрики
	DBConnections   prometheus.Gauge
	DBQueryDuration *prometheus.HistogramVec
	DBQueryErrors   *prometheus.CounterVec

	// RabbitMQ метрики
	RabbitMQMessagesSent     prometheus.Counter
	RabbitMQMessagesReceived prometheus.Counter
	RabbitMQPublishErrors    prometheus.Counter
	RabbitMQConsumeErrors    prometheus.Counter

	// Service specific
	PaymentCreatedTotal prometheus.Counter
	PaymentFailedTotal  prometheus.Counter
	EmailSentTotal      prometheus.Counter
	EmailFailedTotal    prometheus.Counter
}

var metrics *Metrics

func InitMetrics(serviceName string) *Metrics {
	metrics = &Metrics{
		// HTTP метрики
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_http_requests_total", serviceName),
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    fmt.Sprintf("%s_http_request_duration_seconds", serviceName),
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		RequestErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_http_request_errors_total", serviceName),
				Help: "Total number of HTTP request errors",
			},
			[]string{"method", "endpoint"},
		),

		// gRPC метрики
		GrpcRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_grpc_requests_total", serviceName),
				Help: "Total number of gRPC requests",
			},
			[]string{"method", "status"},
		),
		GrpcRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    fmt.Sprintf("%s_grpc_request_duration_seconds", serviceName),
				Help:    "gRPC request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method"},
		),
		GrpcRequestErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_grpc_request_errors_total", serviceName),
				Help: "Total number of gRPC request errors",
			},
			[]string{"method"},
		),

		// Database метрики
		DBConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: fmt.Sprintf("%s_db_connections", serviceName),
				Help: "Number of active database connections",
			},
		),
		DBQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    fmt.Sprintf("%s_db_query_duration_seconds", serviceName),
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"query"},
		),
		DBQueryErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_db_query_errors_total", serviceName),
				Help: "Total number of database query errors",
			},
			[]string{"query"},
		),

		// RabbitMQ метрики
		RabbitMQMessagesSent: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_rabbitmq_messages_sent_total", serviceName),
				Help: "Total number of RabbitMQ messages sent",
			},
		),
		RabbitMQMessagesReceived: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_rabbitmq_messages_received_total", serviceName),
				Help: "Total number of RabbitMQ messages received",
			},
		),
		RabbitMQPublishErrors: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_rabbitmq_publish_errors_total", serviceName),
				Help: "Total number of RabbitMQ publish errors",
			},
		),
		RabbitMQConsumeErrors: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_rabbitmq_consume_errors_total", serviceName),
				Help: "Total number of RabbitMQ consume errors",
			},
		),

		// Service specific
		PaymentCreatedTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_payment_created_total", serviceName),
				Help: "Total number of payments created",
			},
		),
		PaymentFailedTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_payment_failed_total", serviceName),
				Help: "Total number of payments failed",
			},
		),
		EmailSentTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_email_sent_total", serviceName),
				Help: "Total number of emails sent",
			},
		),
		EmailFailedTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_email_failed_total", serviceName),
				Help: "Total number of emails failed",
			},
		),
	}

	return metrics
}

func GetMetrics() *Metrics {
	return metrics
}

// StartMetricsServer starts Prometheus metrics HTTP server
func StartMetricsServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()
	log.Printf("Metrics server started on port %s", port)
}
