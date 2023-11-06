package api

import (
	"encoding/json"
	"net/http"

	"github.com/deividaspetraitis/wallet-screener"
	"github.com/deividaspetraitis/wallet-screener/errors"
	"github.com/deividaspetraitis/wallet-screener/log"

	"github.com/gorilla/mux"
)

// API errors
var (
	ErrAddressNotValid = errors.New("given address is not valid wallet address")
)

// Ethereum hexadecimal address is derived from the last 20 bytes
// of the public key controlling the account with 0x appended in front.
// e.g., 0x71C7656EC7ab88b098defB751B7401B5f6d8976F
const ethWalletAddressLength = 42

// ScreenWalletRiskCategoriesRequest represents HTTP request for screening a wallet for risk categories.
type ScreenWalletRiskCategoriesRequest struct {
	Address string
}

// Validate parses request fields and returns whether they contain valid data.
// Validate implements validator.Validator.
// TODO: implement more sophisticated rule
func (r *ScreenWalletRiskCategoriesRequest) Validate() error {
	if len(r.Address) < ethWalletAddressLength {
		return ErrAddressNotValid
	}
	return nil
}

// UnmarshalHTTP implements http.RequestUnmarshaler.
func (r *ScreenWalletRiskCategoriesRequest) UnmarshalHTTPRequest(req *http.Request) error {
	*r = ScreenWalletRiskCategoriesRequest{
		Address: mux.Vars(req)["address"],
	}
	log.Println("address", req.URL)
	return r.Validate()
}

// ScreenWalletRiskCategoriesResponse represents a response for ScreenWalletRiskCategoriesRequest.
type ScreenWalletRiskCategoriesResponse struct {
	Categories []string `json:"categories"`
}

// MarshalHTTP implements http.Marshaler.
func (r *ScreenWalletRiskCategoriesResponse) MarshalHTTP(w http.ResponseWriter) error {
	if r.Categories == nil {
		r.Categories = []string{}
	}

	return json.NewEncoder(w).Encode(r)
}

// GetWalletRiskCategoriesHistory represents HTTP request for retrieving historical risk categories for a wallet.
type GetWalletRiskCategoriesHistoryRequest struct {
	Address string
}

// Validate parses request fields and returns whether they contain valid data.
// Validate implements validator.Validator.
// TODO: implement more sophisticated rule
func (r *GetWalletRiskCategoriesHistoryRequest) Validate() error {
	if len(r.Address) < ethWalletAddressLength {
		return ErrAddressNotValid
	}
	return nil
}

// UnmarshalHTTP implements http.RequestUnmarshaler.
func (r *GetWalletRiskCategoriesHistoryRequest) UnmarshalHTTPRequest(req *http.Request) error {
	*r = GetWalletRiskCategoriesHistoryRequest{
		Address: mux.Vars(req)["address"],
	}
	return r.Validate()
}

// HistoricalRiskCategory represents wallet historical risk category entity.
type HistoricalRiskCategory struct {
	Category string `json:"category"`
	Revision uint64 `json:"revision"`
}

// NewGetWalletRiskCategoriesHistoryRespone constructs a new response for GetWalletRiskCategoriesHistoryRequest.
func NewGetWalletRiskCategoriesHistoryRespone(categories []*walletscreener.HistoricalRiskCategory) *GetWalletRiskCategoriesHistoryRespone {
	return &GetWalletRiskCategoriesHistoryRespone{
		input: categories,
	}
}

// GetWalletRiskCategoriesHistoryRespone represents a response for GetWalletRiskCategoriesHistoryRequest.
type GetWalletRiskCategoriesHistoryRespone struct {
	input []*walletscreener.HistoricalRiskCategory // state

	Categories []*HistoricalRiskCategory `json:"categories"`
}

// MarshalHTTP implements http.Marshaler.
func (r *GetWalletRiskCategoriesHistoryRespone) MarshalHTTP(w http.ResponseWriter) error {
	for _, v := range r.input {
		r.Categories = append(r.Categories, &HistoricalRiskCategory{
			Category: v.Category,
			Revision: v.Revision,
		})
	}

	if r.Categories == nil {
		r.Categories = []*HistoricalRiskCategory{}
	}

	return json.NewEncoder(w).Encode(r)
}
