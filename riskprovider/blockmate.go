package riskprovider

import (
	"context"
	"encoding/json"
	stdhttp "net/http"
	"time"

	"github.com/deividaspetraitis/wallet-screener/errors"
	"github.com/deividaspetraitis/wallet-screener/http"
	"github.com/deividaspetraitis/wallet-screener/slices"
	"github.com/deividaspetraitis/wallet-screener/token/jwt"
)

// Config represents Blockmate risk provider configuration.
type Config struct {
	APIKey string `mapstructure:"apikey"` // API key
}

// Blockmate is an implementation of walletscreener.WalletRiskScreeningProvider
type Blockmate struct {
	// apiKey is Blockmate API-Key used to authenticate and exchanged for JWT tokens.
	apiKey string

	// jwtToken is auth token used to authorise client when making API calls to blockmate.
	jwtToken string

	// client is an http.Client a library containing HTTP layer methods and helpers.
	*http.Client
}

// NewSumService constructs and returns new SumService instance.
func NewBlockMate(apiKey string, client *http.Client) (*Blockmate, error) {
	if len(apiKey) < 1 {
		return nil, errors.Newf("riskprovider: %s is not a valid API key token", apiKey) // API key is mandatory
	}

	return &Blockmate{
		apiKey: apiKey,
		Client: client,
	}, nil
}

// Request wraps http.Client Request method and makes sure that all requests are authorised by renewing JWT token.
// For more see http.Client docs.
func (c *Blockmate) Request(ctx context.Context, method, uri string, v []byte, options ...http.RequestOption) (*stdhttp.Response, error) {
	_, err := jwt.Parse(c.jwtToken)
	if err == nil {
		return c.authorizedRequest(ctx, method, uri, v, options...) // JWT is present and valid
	}

	c.jwtToken, err = c.AuthProject(ctx, c.apiKey)
	if err != nil {
		return nil, err // authorisation error
	}

	return c.authorizedRequest(ctx, method, uri, v, options...)
}

// authorizedRequest makes authorized requests to the API by adding authorization header with token.
func (c *Blockmate) authorizedRequest(ctx context.Context, method, uri string, v []byte, options ...http.RequestOption) (*stdhttp.Response, error) {
	return c.Client.Request(ctx, method, uri, v, append(options, http.WithBearerToken(c.jwtToken))...)
}

// jwtTokenResponse represents JWT token response from auth blockmate endpoint.
type jwtTokenResponse string

// UnmarshalHTTPResponse implements http.ResponseUnmarshaler
func (t *jwtTokenResponse) UnmarshalHTTPResponse(r *stdhttp.Response) error {
	type response struct {
		Token string `json:"token"`
	}

	var jwt response
	if err := json.NewDecoder(r.Body).Decode(&jwt); err != nil {
		return err
	}

	*t = jwtTokenResponse(jwt.Token)

	return nil
}

// AuthProject returns a JWT token for project.
func (c *Blockmate) AuthProject(ctx context.Context, apiKey string) (jwtToken string, err error) {
	var (
		opts     []http.RequestOption
		response jwtTokenResponse
	)

	opts = append(opts, http.WithHeader("X-API-KEY", c.apiKey))

	res, err := c.Client.Request(ctx, stdhttp.MethodGet, "auth", nil, opts...)
	if err == nil {
		defer res.Body.Close()
		err = http.UnmarshalResponse(res, &response)
	}

	jwtToken = string(response)

	return
}

// detailsCategory represents risk score category details.
type detailsCategory struct {
	Address      string `json:"address"`
	Name         string `json:"name"`
	CategoryName string `json:"category_name"`
	Risk         int    `json:"risk"`
}

// details represents risk score details.
type details struct {
	OwnCategories           []detailsCategory `json:"own_categories"`
	SourceOfFundsCategories []detailsCategory `json:"source_of_funds_categories"`
}

// getAddressRiskScoreDetails represents a response from score details blockmate endpoint.
type getAddressRiskScoreDetails struct {
	CaseID           string    `json:"case_id"`
	RequestDatetime  time.Time `json:"request_datetime"`
	ResponseDatetime time.Time `json:"response_datetime"`
	Chain            string    `json:"chain"`
	Address          string    `json:"address"`
	Name             string    `json:"name"`
	CategoryName     string    `json:"category_name"`
	Risk             int       `json:"risk"`
	Details          details   `json:"details"`
}

// UnmarshalHTTPResponse implements http.ResponseUnmarshaler
func (t *getAddressRiskScoreDetails) UnmarshalHTTPResponse(r *stdhttp.Response) error {
	return json.NewDecoder(r.Body).Decode(t)
}

// GetRiskCategories returns risk categories for given address on ethereum network.
func (c *Blockmate) GetRiskCategories(ctx context.Context, address string) (categories []string, err error) {
	var (
		response getAddressRiskScoreDetails
		opts     []http.RequestOption
	)

	opts = append(opts, http.WithQueryParam("address", address))
	opts = append(opts, http.WithQueryParam("chain", "eth"))

	res, err := c.Request(ctx, stdhttp.MethodGet, "risk/score/details", nil, opts...)
	if err == nil {
		defer res.Body.Close()
		err = http.UnmarshalResponse(res, &response)
	}

	categories = collectUniqueCategoryNames(response.Details.OwnCategories, response.Details.SourceOfFundsCategories)

	return
}

// collectUniqueCategoryNames parses detailsCategory and returns slice containing of unique CategoryName.
func collectUniqueCategoryNames(category ...[]detailsCategory) []string {
	var categories []string

	for _, detailsCategory := range category {
		for _, v := range detailsCategory {
			categories = append(categories, v.CategoryName)
		}
	}

	return slices.Unique(categories)
}
