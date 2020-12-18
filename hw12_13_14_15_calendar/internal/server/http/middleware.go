package internalhttp

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func requestLoggerMiddleware(m *MyHandler, h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(startTime)
		m.app.Log.Info("Request INFO",
			zap.Duration("Duration", duration),
			zap.String("Method", r.Method),
			zap.String("Host", r.Host),
			zap.String("Raw path URL", r.URL.RawPath))
	})
}
