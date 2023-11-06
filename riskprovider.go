package walletscreener

import "context"

// WalletRiskScreeningProvider represents wallet risk screening provider.
type WalletRiskScreeningProvider interface {
	// GetRiskCategories returns a list of risk categories for the given address.
	GetRiskCategories(ctx context.Context, address string) ([]string, error)
}
