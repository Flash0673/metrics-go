package middleware

import (
	"net/http"
	"time"

	"github.com/Flash0673/metrics-go/pkg/logging/logger"
	"go.uber.org/zap"
)

type responseData struct {
	statusCode int
	size       int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write переопределяем метод Write
func (lr *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lr.Write(b)
	lr.responseData.size += size
	return size, err
}

// WriteHeader переопределяем метод WriteHeader
func (lr *loggingResponseWriter) WriteHeader(statusCode int) {
	lr.WriteHeader(statusCode)
	lr.responseData.statusCode = statusCode
}

// RequestResponseLogging .
func RequestResponseLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		method := r.Method
		startTime := time.Now()

		rd := &responseData{}
		loggingRW := &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   rd,
		}

		h.ServeHTTP(loggingRW, r)

		logger.Logger.Info(
			"Request/Response data:",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.Duration("handling time", time.Since(startTime)),
			zap.Int("code", rd.statusCode),
			zap.Int("size", rd.size),
		)
	})
}
