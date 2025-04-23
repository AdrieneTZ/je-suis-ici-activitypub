package middlewares

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type responseWriter struct {
	w      http.ResponseWriter
	status int
}

// Logger returns logger middleware
func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// wrap ResponseWriter to get status code
			wrapResWriter := &responseWriter{w: w, status: http.StatusOK}

			// handle request
			next.ServeHTTP(wrapResWriter, r)

			// count request duration
			duration := time.Since(start)

			// log request in key-value pair format
			logger.Info("HTTP Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.Int("status", wrapResWriter.status),
				zap.Duration("duration", duration),
			)
		})
	}
}

func (rw *responseWriter) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.w.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.w.WriteHeader(statusCode)
}
