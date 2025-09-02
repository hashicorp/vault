// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package physical

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

// TxnEntry is an operation that takes atomically as part of
// a transactional update. Only supported by Transactional backends.
type TxnEntry struct {
	Operation Operation
	Entry     *Entry
}

func (t *TxnEntry) String() string {
	return fmt.Sprintf("Operation: %s. Entry: %s", t.Operation, t.Entry)
}

// Transactional is an optional interface for backends that
// support doing transactional updates of multiple keys. This is
// required for some features such as replication.
type Transactional interface {
	// The function to run a transaction
	Transaction(context.Context, []*TxnEntry) error
}

type TransactionalBackend interface {
	Backend
	Transactional
}

// TransactionalLimits SHOULD be implemented by all TransactionalBackend
// implementations. It is separate for backwards compatibility reasons since
// this in a public SDK module. If a TransactionalBackend does not implement
// this, the historic default limits of 63 entries and 128kb (based on Consul's
// limits) are used by replication internals when encoding batches of
// transactions.
type TransactionalLimits interface {
	TransactionalBackend

	// TransactionLimits must return the limits of how large each transaction may
	// be. The limits returned indicate how many individual operation entries are
	// supported in total and an overall size limit on the contents of each
	// transaction if applicable. Vault will deduct any meta-operations it needs
	// to add from the maxEntries given. maxSize will be compared against the sum
	// of the key and value sizes for all operations in a transaction. The backend
	// should provide a reasonable margin of safety for any overhead it may have
	// while encoding, for example Consul's encoded transaction in JSON must fit
	// in the configured max transaction size so it must leave adequate room for
	// JSON encoding overhead on top of the raw key and value sizes.
	//
	// If zero is returned for either value, the replication internals will use
	// historic reasonable defaults. This allows middleware implementations such
	// as cache layers to either pass through to the underlying backend if it
	// implements this interface, or to return zeros to indicate that the
	// implementer should apply whatever defaults it would use if the middleware
	// were not present.
	TransactionLimits() (maxEntries int, maxSize int)
}

type PseudoTransactional interface {
	// An internal function should do no locking or permit pool acquisition.
	// Depending on the backend and if it natively supports transactions, these
	// may simply chain to the normal backend functions.
	GetInternal(context.Context, string) (*Entry, error)
	PutInternal(context.Context, *Entry) error
	DeleteInternal(context.Context, string) error
}

// Implements the transaction interface
func GenericTransactionHandler(ctx context.Context, t PseudoTransactional, txns []*TxnEntry) (retErr error) {
	rollbackStack := make([]*TxnEntry, 0, len(txns))
	var dirty bool

	// Update all of our GET transaction entries, so we can populate existing values back at the wal layer.
	for _, txn := range txns {
		if txn.Operation == GetOperation {
			entry, err := t.GetInternal(ctx, txn.Entry.Key)
			if err != nil {
				return err
			}
			if entry != nil {
				txn.Entry.Value = entry.Value
			}
		}
	}

	// We walk the transactions in order; each successful operation goes into a
	// LIFO for rollback if we hit an error along the way
TxnWalk:
	for _, txn := range txns {
		switch txn.Operation {
		case DeleteOperation:
			entry, err := t.GetInternal(ctx, txn.Entry.Key)
			if err != nil {
				retErr = multierror.Append(retErr, err)
				dirty = true
				break TxnWalk
			}
			if entry == nil {
				// Nothing to delete or roll back
				continue
			}
			rollbackEntry := &TxnEntry{
				Operation: PutOperation,
				Entry: &Entry{
					Key:   entry.Key,
					Value: entry.Value,
				},
			}
			err = t.DeleteInternal(ctx, txn.Entry.Key)
			if err != nil {
				retErr = multierror.Append(retErr, err)
				dirty = true
				break TxnWalk
			}
			rollbackStack = append([]*TxnEntry{rollbackEntry}, rollbackStack...)

		case PutOperation:
			entry, err := t.GetInternal(ctx, txn.Entry.Key)
			if err != nil {
				retErr = multierror.Append(retErr, err)
				dirty = true
				break TxnWalk
			}

			// Nothing existed so in fact rolling back requires a delete
			var rollbackEntry *TxnEntry
			if entry == nil {
				rollbackEntry = &TxnEntry{
					Operation: DeleteOperation,
					Entry: &Entry{
						Key: txn.Entry.Key,
					},
				}
			} else {
				rollbackEntry = &TxnEntry{
					Operation: PutOperation,
					Entry: &Entry{
						Key:   entry.Key,
						Value: entry.Value,
					},
				}
			}

			err = t.PutInternal(ctx, txn.Entry)
			if err != nil {
				retErr = multierror.Append(retErr, err)
				dirty = true
				break TxnWalk
			}
			rollbackStack = append([]*TxnEntry{rollbackEntry}, rollbackStack...)
		}
	}

	// Need to roll back because we hit an error along the way
	if dirty {
		// While traversing this, if we get an error, we continue anyways in
		// best-effort fashion
		for _, txn := range rollbackStack {
			switch txn.Operation {
			case DeleteOperation:
				err := t.DeleteInternal(ctx, txn.Entry.Key)
				if err != nil {
					retErr = multierror.Append(retErr, err)
				}
			case PutOperation:
				err := t.PutInternal(ctx, txn.Entry)
				if err != nil {
					retErr = multierror.Append(retErr, err)
				}
			}
		}
	}

	return
}
