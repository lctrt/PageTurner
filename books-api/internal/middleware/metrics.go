package middleware

import (
	"net/http"
	"strconv"
	"time"

	"books/internal/metrics"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := chiMiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		route := chi.RouteContext(r.Context()).RoutePattern()
		if route == "" {
			route = "unknown"
		}

		metrics.HTTPRequestsTotal.
			WithLabelValues(r.Method, route, strconv.Itoa(ww.Status())).Inc()
		metrics.HTTPRequestDuration.
			WithLabelValues(r.Method, route).
			Observe(time.Since(start).Seconds())
	})
}
