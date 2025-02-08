package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    // Register the metrics
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}

type Response struct {
    Message string `json:"message"`
}

func instrumentHandler(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create a custom response writer to capture the status code
        rw := &responseWriter{w, http.StatusOK}
        
        // Call the original handler
        handler(rw, r)
        
        // Record metrics
        duration := time.Since(start).Seconds()
        httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
        httpRequestsTotal.WithLabelValues(
            r.Method,
            r.URL.Path,
            strconv.Itoa(rw.statusCode),
        ).Inc()
    }
}

// Custom response writer to capture status code
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    response := Response{
        Message: "Hello, World!",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func main() {
    // API routes
    http.HandleFunc("/api", instrumentHandler(apiHandler))
    http.HandleFunc("/health", instrumentHandler(healthHandler))
    
    // Prometheus metrics endpoint
    http.Handle("/metrics", promhttp.Handler())
    
    // Start server
    log.Printf("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
