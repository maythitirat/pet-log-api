package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

// Logger is a middleware that logs the start and end of each request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap response writer to capture status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		
		// Process request
		next.ServeHTTP(ww, r)
		
		// Log after request is complete
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Int("status", ww.Status()).
			Int("bytes", ww.BytesWritten()).
			Dur("latency", time.Since(start)).
			Str("request_id", middleware.GetReqID(r.Context())).
			Msg("request completed")
	})
}
