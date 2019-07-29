package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ilyakaznacheev/secret/database"
	"github.com/ilyakaznacheev/secret/handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	var conf RedisConfig
	cleanenv.ReadEnv(&conf)

	metrics := map[string]*metricSet{
		"secret_post": newMetricSet("secret_post"),
		"secret_get":  newMetricSet("secret_get"),
	}

	db, err := database.NewRedisDB(conf.Host+":"+conf.Port, conf.Password, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	h := handler.NewSecretHandler(db)

	router := gin.Default()

	v1 := router.Group("/v1")
	v1.POST("/secret", metricsMiddleware(h.PostSecret, metrics["secret_post"]))
	v1.GET("/secret/:hash", metricsMiddleware(h.GetSecret, metrics["secret_get"]))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Run service
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func metricsMiddleware(hf gin.HandlerFunc, ms *metricSet) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		hf(c)
		respTimeMSec := float64(time.Now().Sub(start).Nanoseconds()) / 1000.0

		ms.requestCounter.Inc()
		ms.responceTimeGauge.Set(respTimeMSec)
		ms.responceTimeQuantile.Observe(respTimeMSec)
	}
}

type metricSet struct {
	endpointName         string
	requestCounter       prometheus.Counter
	responceTimeGauge    prometheus.Gauge
	responceTimeQuantile prometheus.Summary
}

// newMetricSet creates a metric for an endpoint
func newMetricSet(endpoint string) *metricSet {
	cl := map[string]string{
		"ip":       getLocalIP(),
		"endpoint": endpoint,
	}

	ms := metricSet{
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

// RedisConfig ia redis-related configuration
type RedisConfig struct {
	Port     string `env:"REDIS_PORT" env-default:"5050"`
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Password string `env:"REDIS_PASSWORD"`
}

// Config is an application configuration structure
type Config struct {
	Redis RedisConfig
}
