package internalhttp

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func requestLoggerMiddleware(m *MyHandler, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(startTime)
		m.app.Log.Info("Request INFO",
			zap.Duration("Duration", duration),
			zap.String("Method", r.Method),
			zap.String("Host", r.Host),
			zap.String("Raw path URL", r.URL.RawPath))
	}
}

var headers = map[string]string{
	"Content-Type": "application/json; charset=utf-8",
	"test":         "test",
}

func headerSetter(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		for k, v := range headers {
			rw.Header().Set(k, v)
		}
		fn(rw, req)
	}
}
