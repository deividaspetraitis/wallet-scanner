package http

import (
	"context"
	"net/http"

	"github.com/deividaspetraitis/wallet-screener"
	"github.com/deividaspetraitis/wallet-screener/log"
	"github.com/deividaspetraitis/wallet-screener/pkg/api/v1"
)

// getRiskCategoriesFunc decouples actual check implementation and allows easily test HTTP handler.
type getRiskCategoriesFunc func(ctx context.Context, address string) ([]string, error)

// GetRiskCategories responds with risk categories list for given address.
func GetRiskCategories(getRiskCategories getRiskCategoriesFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// It's always json.
		w.Header().Set("Content-Type", "application/json")

		var request api.ScreenWalletRiskCategoriesRequest
		if err := UnmarshalRequest(r, &request); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"handler": "wallet",
				"method":  "GetRiskCategories",
			}).Println("unable to unmarshal request data")

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		categories, err := getRiskCategories(r.Context(), request.Address)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"handler": "wallet",
				"method":  "GetRiskCategories",
			}).Println("encountered an error retrieving risk categories")

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := api.ScreenWalletRiskCategoriesResponse{
			Categories: categories,
		}

		w.WriteHeader(http.StatusOK)
		if err := Marshal(w, &response); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"handler": "wallet",
				"method":  "GetRiskCategories",
			}).Println("unable to marshal response data")

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// historyFunc decouples actual check implementation and allows easily test HTTP handler.
type getRiskCategoriesHistoryFunc func(ctx context.Context, address string) ([]*walletscreener.HistoricalRiskCategory, error)

// History responds with risk categories history list for given address.
func GetRiskCategoriesHistory(getRiskCategoriesHistory getRiskCategoriesHistoryFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// It's always json.
		w.Header().Set("Content-Type", "application/json")

		var request api.GetWalletRiskCategoriesHistoryRequest
		if err := UnmarshalRequest(r, &request); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"handler": "wallet",
				"method":  "GetRiskCategoriesHistory",
			}).Println("unable to unmarshal request data")

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		categories, err := getRiskCategoriesHistory(r.Context(), request.Address)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"handler": "wallet",
				"method":  "GetRiskCategoriesHistory",
			}).Println("encountered an error retrieving risk categories history")

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := api.NewGetWalletRiskCategoriesHistoryRespone(categories)

		w.WriteHeader(http.StatusOK)
		if err := Marshal(w, response); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"handler": "wallet",
				"method":  "GetRiskCategoriesHistory",
			}).Println("unable to marshal response data")

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
