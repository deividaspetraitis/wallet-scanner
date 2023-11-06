package api

import (
	"testing"

	"github.com/deividaspetraitis/wallet-screener/validator"
)

func TestScreenWalletRiskCategoriesRequest(t *testing.T) {
	var testcases = []struct {
		Address string
		Error   error
	}{
		{"", ErrAddressNotValid},
		{"0x4E9ce36E442e55EcD9025B9a6E0D88485d628A67", nil},
	}

	for _, v := range testcases {
		req := ScreenWalletRiskCategoriesRequest{
			Address: v.Address,
		}

		err := validator.Validate(&req)
		if err != v.Error {
			t.Errorf("got %v, want %v", err, v.Error)
		}
	}

}
