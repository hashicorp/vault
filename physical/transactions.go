package physical

import (
	multierror "github.com/hashicorp/go-multierror"
)

// TxnEntry is an operation that takes atomically as part of
// a transactional update. Only supported by Transactional backends.
type TxnEntry struct {
	Operation Operation
	Entry     *Entry
}

// Transactional is an optional interface for backends that
// support doing transactional updates of multiple keys. This is
// required for some features such as replication.
type Transactional interface {
	// The function to run a transaction
	Transaction([]*TxnEntry) error
}

type PseudoTransactional interface {
	// An internal function should do no locking or permit pool acquisition.
	// Depending on the backend and if it natively supports transactions, these
	// may simply chain to the normal backend functions.
	GetInternal(string) (*Entry, error)
	PutInternal(*Entry) error
	DeleteInternal(string) error
}

// Implements the transaction interface
func GenericTransactionHandler(t PseudoTransactional, txns []*TxnEntry) (retErr error) {
	rollbackStack := make([]*TxnEntry, 0, len(txns))
	var dirty bool

	// We walk the transactions in order; each successful operation goes into a
	// LIFO for rollback if we hit an error along the way
TxnWalk:
	for _, txn := range txns {
		switch txn.Operation {
		case DeleteOperation:
			entry, err := t.GetInternal(txn.Entry.Key)
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
					Key:      entry.Key,
					Value:    entry.Value,
					SealWrap: entry.SealWrap,
				},
			}
			err = t.DeleteInternal(txn.Entry.Key)
			if err != nil {
				retErr = multierror.Append(retErr, err)
				dirty = true
				break TxnWalk
			}
			rollbackStack = append([]*TxnEntry{rollbackEntry}, rollbackStack...)

		case PutOperation:
			entry, err := t.GetInternal(txn.Entry.Key)
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
						Key:      entry.Key,
						Value:    entry.Value,
						SealWrap: entry.SealWrap,
					},
				}
			}

			err = t.PutInternal(txn.Entry)
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
				err := t.DeleteInternal(txn.Entry.Key)
				if err != nil {
					retErr = multierror.Append(retErr, err)
				}
			case PutOperation:
				err := t.PutInternal(txn.Entry)
				if err != nil {
					retErr = multierror.Append(retErr, err)
				}
			}
		}
	}

	return
}
