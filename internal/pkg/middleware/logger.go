package middleware

import (
	"net/http"
	"time"

	"github.com/zsbahtiar/lionparcel-test/internal/pkg/logger"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode >= 400 {
		rw.body = append(rw.body, b...)
	}
	return rw.ResponseWriter.Write(b)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get request details
		path := r.URL.Path
		method := r.Method
		ip := r.RemoteAddr
		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			ip = forwardedFor
		}

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		latency := time.Since(start)

		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", ip),
			zap.Int("status", rw.statusCode),
			zap.Duration("latency", latency),
		}

		if rw.statusCode >= 400 {
			if len(rw.body) > 0 {
				fields = append(fields, zap.String("response", string(rw.body)))
			}
			logger.Error("Request failed", fields...)
		} else {
			logger.Info("Request completed", fields...)
		}
	})
}
