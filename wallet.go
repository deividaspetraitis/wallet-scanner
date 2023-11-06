package walletscreener

import (
	"context"

	"github.com/deividaspetraitis/wallet-screener/errors"
)

// HistoricalRiskCategory represents wallet historical risk category entity.
type HistoricalRiskCategory struct {
	Category string // Risk category
	Revision uint64 //  Revision of the category
}

// StoreWalletRiskCategoriesFunc stores risk categories for a given wallet address into database.
// This function is atomic, failure to store single category will result in failure storing the rest categories.
type StoreWalletRiskCategoriesFunc func(ctx context.Context, address string, categories []string) error

// ScreenWalletRiskCategories screens a wallet to fetch risk categories list for the given address from RiskProvider.
// Fetched categories will be stored into database for future reference.
func ScreenWalletRiskCategories(ctx context.Context, riskprovider WalletRiskScreeningProvider, storeRiskCategories StoreWalletRiskCategoriesFunc, address string) ([]string, error) {
	categories, err := riskprovider.GetRiskCategories(ctx, address)
	if err != nil {
		return nil, err
	}

	if err := storeRiskCategories(ctx, address, categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetWalletRiskCategories retrieves list of risk categories as slice of strings along revisions
// as slice of uint64 for given wallet address from the database.
// Length of categories and revisions always are the same and are indexed in same order.
type GetWalletRiskCategoriesFunc func(ctx context.Context, address string) ([]string, []uint64, error)

// GetWalletRiskCategoriesHistory retrieves history of risk categories for given wallet address.
func GetWalletRiskCategoriesHistory(ctx context.Context, getRiskCategories GetWalletRiskCategoriesFunc, address string) ([]*HistoricalRiskCategory, error) {
	categories, revision, err := getRiskCategories(ctx, address)
	if err != nil {
		return nil, errors.New("failed to fetch historical categories")
	}

	var result []*HistoricalRiskCategory
	for i, v := range categories {
		result = append(result, &HistoricalRiskCategory{
			Category: v,
			Revision: revision[i],
		})
	}

	return result, nil
}
