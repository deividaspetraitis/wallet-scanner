package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/deividaspetraitis/wallet-screener/errors"

	"github.com/gorilla/mux"
)

func TestGetRiskCategories(t *testing.T) {
	var testcases = []struct {
		address           string
		getRiskCategories getRiskCategoriesFunc

		response   string
		statusCode int
	}{
		// not a valid address
		{
			address: "abc",
			getRiskCategories: func(ctx context.Context, address string) ([]string, error) {
				return nil, nil
			},
			statusCode: http.StatusBadRequest,
		},
		// empty categories list
		{
			address: "0x4E9ce36E442e55EcD9025B9a6E0D88485d628A67",
			getRiskCategories: func(ctx context.Context, address string) ([]string, error) {
				return nil, nil
			},
			response:   `{"categories":[]}`,
			statusCode: http.StatusOK,
		},
		// non-empty categories list
		{
			address: "0x4E9ce36E442e55EcD9025B9a6E0D88485d628A67",
			getRiskCategories: func(ctx context.Context, address string) ([]string, error) {
				return []string{"category1", "category2"}, nil
			},
			response:   `{"categories":["category1","category2"]}`,
			statusCode: http.StatusOK,
		},
		// service error
		{
			address: "0x4E9ce36E442e55EcD9025B9a6E0D88485d628A67",
			getRiskCategories: func(ctx context.Context, address string) ([]string, error) {
				return nil, errors.New("test getRiskCategories errors")
			},
			response:   "",
			statusCode: http.StatusInternalServerError,
		},
	}

	for i, tt := range testcases {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost/wallet/%s/categories", tt.address), nil)
		w := httptest.NewRecorder()

		// To add the vars to the context we need to create a router through which we can pass the request.
		// TODO: tests should be not aware of routing mechanism.
		router := mux.NewRouter()
		router.HandleFunc("/wallet/{address}/categories", GetRiskCategories(tt.getRiskCategories))

		router.ServeHTTP(w, req)

		if statusCode := w.Result().StatusCode; statusCode != tt.statusCode {
			t.Errorf("#%d HTTP status got %v, want %v", i, statusCode, tt.statusCode)
		}

		// we do apply TrimSpace to clean up response coming from HTTP protocol
		if response := strings.TrimSpace(w.Body.String()); response != tt.response {
			t.Errorf("#%d HTTP status got %v, want %s", i, response, tt.response)
		}
	}
}
