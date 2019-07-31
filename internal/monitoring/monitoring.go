package monitoring

import (
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsMiddleware adds Prometheus monitoring to the endpoint handler
func MetricsMiddleware(hf gin.HandlerFunc, endpoint string) gin.HandlerFunc {
	ms := NewMetricSet(endpoint)

	return func(c *gin.Context) {
		start := time.Now()
		hf(c)
		// nanosec to millisec
		respTimeMSec := float64(time.Now().Sub(start).Nanoseconds()) / 1000000.0

		ms.requestCounter.Inc()
		ms.responceTimeGauge.Set(respTimeMSec)
		ms.responceTimeQuantile.Observe(respTimeMSec)
	}
}

// MetricSet is a set of Prometheus metrics of endpoint
type MetricSet struct {
	endpointName         string
	requestCounter       prometheus.Counter
	responceTimeGauge    prometheus.Gauge
	responceTimeQuantile prometheus.Summary
}

// NewMetricSet creates a metric for an endpoint
func NewMetricSet(endpoint string) *MetricSet {
	cl := map[string]string{
		"ip":       getLocalIP(),
		"endpoint": endpoint,
	}

	ms := MetricSet{
		endpointName: endpoint,

		requestCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace:   "secret",
				Name:        "request_number",
				Help:        "Number of API endpoint requests",
				ConstLabels: cl,
			},
		),

		responceTimeGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "secret",
				Name:        "request_processing_time_ms",
				Help:        "API endpoint response time in milliseconds",
				ConstLabels: cl,
			},
		),

		responceTimeQuantile: promauto.NewSummary(
			prometheus.SummaryOpts{
				Namespace:   "secret",
				Name:        "request_processing_time_summary_ms",
				Help:        "API endpoint response time quantile",
				ConstLabels: cl,
				Objectives: map[float64]float64{
					0.5:  0.005,
					0.95: 0.0005,
					0.99: 0.0001,
				},
			},
		),
	}

	return &ms
}

func getLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "-"
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		// handle err
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				return v.IP.String()
			case *net.IPAddr:
				return v.IP.String()
			}
		}
	}
	return "-"
}
