package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/deividaspetraitis/wallet-screener/log"
)

func TestRequestRate(t *testing.T) {
	var testcases = []struct {
		limit     int
		threshold time.Duration

		requests   int
		statusCode int
	}{
		// should pass: requests <= limit
		{
			limit:     10,
			threshold: time.Second,

			requests:   10,
			statusCode: http.StatusOK,
		},
		// should fail: requests > limit
		{
			limit:     100,
			threshold: time.Second,

			requests:   101,
			statusCode: http.StatusTooManyRequests,
		},
	}

	for _, tt := range testcases {
		// last request result
		var w *httptest.ResponseRecorder

		// construct test target
		middleware := RequestRate(tt.limit, tt.threshold, log.Default(), nil)

		for i := 0; i < tt.requests; i++ {
			req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
			w = httptest.NewRecorder()

			middleware.ServeHTTP(w, req)
		}

		if statusCode := w.Result().StatusCode; statusCode != tt.statusCode {
			t.Errorf("HTTP status got %v, want %v", statusCode, tt.statusCode)
		}
	}
}
