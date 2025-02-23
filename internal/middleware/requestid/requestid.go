package requestid

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/RinatZaynet/CouchFilmCritic/internal/cookie/tracking"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/random"
	"github.com/RinatZaynet/CouchFilmCritic/internal/jwtutill"
)

const reqIDLen = 10

type ctxKeyRequestID int

const requestIDKey ctxKeyRequestID = 0

func New(log *slog.Logger, manager *jwtutill.Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/requestID"),
		)

		log.Info("requestID middleware enabled")

		handler := func(w http.ResponseWriter, r *http.Request) {
			reqID := tracking.GetRequestID(r)

			if reqID == "" {
				tracking.CreateRequestIDCookie(w, random.RandomString(reqIDLen))
			}

			ctx := r.Context()

			ctx = context.WithValue(ctx, requestIDKey, reqID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(handler)
	}
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {

		return reqID
	}

	return ""
}
