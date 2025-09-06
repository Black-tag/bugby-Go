package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/blacktag/bugby-Go/internal/metrics"
)


func MetricsMiddleware(m *metrics.Metrics) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            


			start := time.Now()


			recorder := &statusRecorder{ResponseWriter: w, status: 200}

            next.ServeHTTP(recorder, r)


			duration := time.Since(start).Seconds()

			m.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
			m.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)

			slog.Info("request handled",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", duration,
				


			)
            
        })
    }
}



type statusRecorder struct {
	http.ResponseWriter
	status int
}


func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
