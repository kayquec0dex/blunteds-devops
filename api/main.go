package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total de requests recebidos",
		},
		[]string{"method", "route", "status"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Duração das requests em segundos",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	requestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "api_requests_in_flight",
			Help: "Requests sendo processadas no momento",
		},
	)
)

func metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}

		requestsInFlight.Inc()
		defer requestsInFlight.Dec()

		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		requestsTotal.WithLabelValues(c.Request.Method, route, status).Inc()
		requestDuration.WithLabelValues(c.Request.Method, route).Observe(duration)
	}
}

func main() {
	log.Println("[INFO] Iniciando blunteds-devops-api...")

	r := gin.Default()
	r.Use(metricsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		log.Printf("[INFO] health check requested from %s", c.ClientIP())
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/", func(c *gin.Context) {
		log.Printf("[INFO] root requested from %s", c.ClientIP())
		c.JSON(http.StatusOK, gin.H{
			"message": "Bem-vindo à DevOps API! v2",
		})
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	log.Println("[INFO] Server listening on :8080")
	r.Run(":8080")
}