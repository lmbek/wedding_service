package webserver

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for handler in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)

	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
)

// accessibleFromOutside is a global flag that controls whether security features are enabled
var accessibleFromOutside = false

func init() {
	// Register the Prometheus metrics
	prometheus.MustRegister(requestDuration, requestCount)
}

// statusRecorder is a custom implementation of http.ResponseWriter
// that records the status and captures the response body.
type statusRecorder struct {
	http.ResponseWriter
	status       int
	body         *bytes.Buffer
	originalBody []byte
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Capture the response body
func (r *statusRecorder) Write(p []byte) (n int, err error) {
	r.body.Write(p)                  // capture body content
	return r.ResponseWriter.Write(p) // write to original response
}

// Utility function to check if the request is from localhost or allowed IP (Docker network, etc.)
func isLocalRequest(r *http.Request) bool {
	// Check if the request is from localhost
	if r.RemoteAddr == "127.0.0.1" || r.RemoteAddr == "::1" {
		return true
	}

	// Check if the request is from a specific trusted IP range (e.g., Docker network)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}

	// Example: Allow Docker containers within a specific network range (adjust as needed)
	_, ipNet, _ := net.ParseCIDR("172.17.0.0/16") // Docker bridge network range
	return ipNet.Contains(net.ParseIP(ip))
}

// LoggingAndMetricsMiddleware logs requests and tracks metrics
func LoggingAndMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging and metrics for /metrics endpoint unless security is disabled
		if accessibleFromOutside && strings.HasPrefix(r.URL.Path, "/metrics") && !isLocalRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		// Skip WebSocket requests from logging and metrics tracking unless security is disabled
		if accessibleFromOutside && strings.HasPrefix(r.URL.Path, "/websocket") {
			next.ServeHTTP(w, r)
			return
		}

		// Capture request body (only once)
		var reqBody []byte
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
			// Restore the body for further use by the handler
			r.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		// Start measuring the time
		start := time.Now()
		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
			body:           bytes.NewBuffer([]byte{}),
		}

		// Proceed with request processing and track metrics
		next.ServeHTTP(rec, r)

		// Calculate request duration
		duration := time.Since(start).Seconds()
		statusStr := fmt.Sprintf("%d", rec.status)

		// Record Prometheus metrics
		requestDuration.WithLabelValues(r.Method, r.URL.Path, statusStr).Observe(duration)
		requestCount.WithLabelValues(r.Method, r.URL.Path, statusStr).Inc()

		// Log the expanded request data with bodies included
		log.Printf(
			"Method: %s | Path: %s | Status: %s | Duration: %.3fs | IP: %s | User-Agent: %s | Query Params: %s | Headers: %v | Request Body: %s | Response Body: %s | Content-Length: %d\n",
			r.Method,
			r.URL.Path,
			statusStr,
			duration,
			r.RemoteAddr,
			r.UserAgent(),
			r.URL.RawQuery,
			r.Header,
			string(reqBody),   // Request body (captured)
			rec.body.String(), // Response body (captured)
			r.ContentLength,
		)
	})
}
