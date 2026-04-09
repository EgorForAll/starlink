package rest

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			l := logger.With().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Logger()

			r = r.WithContext(l.WithContext(r.Context()))

			next.ServeHTTP(w, r)

			log.Ctx(r.Context()).Info().
				Dur("duration", time.Since(start)).
				Msg("request completed")
		})
	}
}