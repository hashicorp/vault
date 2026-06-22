package ydb

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/physical"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type ydbTxWrapper struct {
	tx    query.TxActor
	table string
}

func (w *ydbTxWrapper) GetInternal(ctx context.Context, key string) (*physical.Entry, error) {
	stmt := fmt.Sprintf("SELECT key, value FROM %s WHERE key = $key", w.table)
	params := ydb.ParamsBuilder().Param("$key").Text(key).Build()

	res, err := w.tx.Query(ctx, stmt, query.WithParameters(params))
	if err != nil {
		return nil, err
	}
	defer res.Close(ctx)

	for rs, rerr := range res.ResultSets(ctx) {
		if rerr != nil {
			return nil, rerr
		}
		for row, rerr := range rs.Rows(ctx) {
			if rerr != nil {
				return nil, rerr
			}
			var k string
			var v []byte
			if err := row.Scan(&k, &v); err != nil {
				return nil, err
			}
			return &physical.Entry{Key: k, Value: v}, nil
		}
	}
	return nil, nil
}

func (w *ydbTxWrapper) PutInternal(ctx context.Context, entry *physical.Entry) error {
	stmt := fmt.Sprintf("UPSERT INTO %s (key, value) VALUES ($key, $value)", w.table)
	params := ydb.ParamsBuilder().
		Param("$key").Text(entry.Key).
		Param("$value").Bytes(entry.Value).Build()
	return w.tx.Exec(ctx, stmt, query.WithParameters(params))
}

func (w *ydbTxWrapper) DeleteInternal(ctx context.Context, key string) error {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE key = $key", w.table)
	params := ydb.ParamsBuilder().Param("$key").Text(key).Build()
	return w.tx.Exec(ctx, stmt, query.WithParameters(params))
}

func (y *YDBBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	if len(txns) == 0 {
		return nil
	}
	return y.db.Query().DoTx(ctx, func(ctx context.Context, tx query.TxActor) error {
		w := &ydbTxWrapper{tx: tx, table: y.table}
		return physical.GenericTransactionHandler(ctx, w, txns)
	})
}

func (y *YDBBackend) TransactionLimits() (int, int) {
	// These defaults are intentionally conservative.
	// Actual YDB limits are yet to be validated.
	return y.transactionMaxEntries, y.transactionMaxSize
}
