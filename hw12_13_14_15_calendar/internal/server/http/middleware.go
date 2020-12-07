package internalhttp

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
	})
}

func withTimeout(h http.HandlerFunc, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(r.Context(), timeout)
		r = r.WithContext(ctx)

		h(w, r)
	}
}

func requestLogger(m *MyHandler, h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Now().Sub(startTime)
		m.l.Info("Request INFO",
			zap.Duration("Duration", duration),
			zap.String("Method", r.Method),
			zap.String("Host", r.Host),
			zap.String("Raw path URL", r.URL.RawPath),
			zap.String("Scheme URL", r.URL.Scheme),
			zap.String("User URL", r.URL.User.String()))

	})
}
