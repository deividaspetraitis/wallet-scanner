package riskprovider

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/deividaspetraitis/wallet-screener/http"

	"github.com/google/go-cmp/cmp"
)

func newGetAddressRiskScoreDetails(t *testing.T) *getAddressRiskScoreDetails {
	t.Helper()
	return &getAddressRiskScoreDetails{
		CaseID:           "e8f0db90-5a31-44b0-930d-e83a4d573947",
		RequestDatetime:  time.Date(2023, 10, 4, 15, 18, 21, 0, time.UTC),
		ResponseDatetime: time.Date(2023, 10, 4, 15, 18, 21, 0, time.UTC),
		Chain:            "eth",
		Address:          "0xe9e9afac38e64728f1afbb2b65dec7be7c704c05",
		Name:             "unknown",
		CategoryName:     "Banned",
		Risk:             100,
		Details: details{
			OwnCategories: []detailsCategory{
				{
					Address:      "0xe9e9afac38e64728f1afbb2b65dec7be7c704c05",
					Name:         "unknown",
					CategoryName: "Banned",
					Risk:         100,
				},
			},
			SourceOfFundsCategories: []detailsCategory{
				{
					Address:      "0xe9e9afac38e64728f1afbb2b65dec7be7c704c05",
					Name:         "unknown",
					CategoryName: "Banned",
					Risk:         100,
				},
			},
		},
	}
}

var getAddressRiskScoreDetailsPayload = []byte(`{
  "case_id": "e8f0db90-5a31-44b0-930d-e83a4d573947",
  "request_datetime": "2023-10-04T15:18:21Z",
  "response_datetime": "2023-10-04T15:18:21Z",
  "chain": "eth",
  "address": "0xe9e9afac38e64728f1afbb2b65dec7be7c704c05",
  "name": "unknown",
  "category_name": "Banned",
  "risk": 100,
  "details": {
    "own_categories": [
      {
        "address": "0xe9e9afac38e64728f1afbb2b65dec7be7c704c05",
        "name": "unknown",
        "category_name": "Banned",
        "risk": 100
      }
    ],
    "source_of_funds_categories": [
      {
        "address": "0xe9e9afac38e64728f1afbb2b65dec7be7c704c05",
        "name": "unknown",
        "category_name": "Banned",
        "risk": 100
      }
    ]
  }
}`)

func TestUnmarshalHTTPResponse(t *testing.T) {
	t.Run("getAddressRiskScoreDetails", func(t *testing.T) {
		var response getAddressRiskScoreDetails
		w := httptest.NewRecorder()

		if _, err := w.WriteString(string(getAddressRiskScoreDetailsPayload)); err != nil {
			t.Errorf("got %v, want %v", err, nil)
		}

		if err := http.UnmarshalResponse(w.Result(), &response); err != nil {
			t.Errorf("got %v, want %v", err, nil)
		}

		if expected := newGetAddressRiskScoreDetails(t); !cmp.Equal(&response, expected) {
			t.Errorf("got %v, want %v", response, expected)
		}
	})

	// TODO: test token
}
