package middleware

import (
	"net/http"
	"time"

	"github.com/deividaspetraitis/wallet-screener/log"

	"golang.org/x/time/rate"
)

// limiter is a global limiter declaration used across multiple requests.
var limiter *rate.Limiter

// RequestRate verifies that request rate is not higher per given threshold.
// Once given rate is exhausted further HTTP requests will be terminated with HTTP 429.
func RequestRate(limit int, threshold time.Duration, logger log.Logger, h http.Handler) http.Handler {
	if limiter == nil {
		limiter = rate.NewLimiter(rate.Every(threshold), limit)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		if h != nil {
			h.ServeHTTP(w, r)
		}
	})
}
