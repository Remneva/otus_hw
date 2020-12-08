//nolint
package internalhttp

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func withTimeout(h http.HandlerFunc, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := context.WithTimeout(r.Context(), timeout)
		if err != nil {
			return
		}
		r = r.WithContext(ctx)
		h(w, r)
	}
}

func requestLoggerMiddleware(m *MyHandler, h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(startTime)
		m.l.Info("Request INFO",
			zap.Duration("Duration", duration),
			zap.String("Method", r.Method),
			zap.String("Host", r.Host),
			zap.String("Raw path URL", r.URL.RawPath),
			zap.String("User URL", r.URL.User.String()))
	})
}
