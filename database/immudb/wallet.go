package immudb

import (
	"context"

	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/deividaspetraitis/wallet-screener/errors"
)

// StoreWalletCategories implements StoreWalletRiskCategoriesFunc.
func StoreWalletRiskCategories(ctx context.Context, db immudb.ImmuClient, address string, categories []string) error {
	var kvs []*schema.KeyValue
	for _, v := range categories {
		kvs = append(kvs, &schema.KeyValue{
			Key:   []byte(address),
			Value: []byte(v),
		})
	}

	_, err := db.SetAll(ctx, &schema.SetRequest{
		KVs: kvs,
	})
	if err != nil {
		return errors.Wrap(err, "failed to store address scores")
	}

	return nil
}

// GetWalletRiskCategories implements GetWalletRiskCategoriesFunc.
func GetWalletRiskCategories(ctx context.Context, db immudb.ImmuClient, address string) ([]string, []uint64, error) {
	entries, err := db.History(ctx, &schema.HistoryRequest{
		Key: []byte(address),
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to retrieve category history for address %s", address)
	}

	var (
		categories []string
		revisions  []uint64
	)
	for _, v := range entries.GetEntries() {
		categories = append(categories, string(v.GetValue()))
		revisions = append(revisions, v.GetRevision())
	}

	return categories, revisions, nil
}
