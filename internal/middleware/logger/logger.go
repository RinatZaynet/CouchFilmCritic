package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/RinatZaynet/CouchFilmCritic/internal/middleware/requestid"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		handler := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", requestid.GetRequestID(r.Context())),
			)

			t1 := time.Now()

			defer func() {
				entry.Info("request completed",
					slog.String("duration", time.Since(t1).String()),
				)
			}()
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(handler)
	}
}
