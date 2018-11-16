package physical

import (
	"context"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	log "github.com/hashicorp/go-hclog"
)

var ErrNonUTF8 = errors.New("key contains invalid UTF-8 characters")
var ErrNonPrintable = errors.New("key contains invalid UTF-8 characters")

// ErrorInjector is used to add errors into underlying physical requests
type StorageEncoding struct {
	Backend
}

// TransactionalErrorInjector is the transactional version of the error
// injector
type TransactionalStorageEncoding struct {
	*StorageEncoding
	Transactional
}

// Verify ErrorInjector satisfies the correct interfaces
var _ Backend = (*StorageEncoding)(nil)
var _ Transactional = (*TransactionalStorageEncoding)(nil)

// NewErrorInjector returns a wrapped physical backend to inject error
func NewStorageEncoding(b Backend, logger log.Logger) *StorageEncoding {
	return &StorageEncoding{
		Backend: b,
	}
}

// NewTransactionalErrorInjector creates a new transactional ErrorInjector
func NewTransactionalStorageEncoding(b Backend, errorPercent int, logger log.Logger) *TransactionalStorageEncoding {
	return &TransactionalStorageEncoding{
		StorageEncoding: NewStorageEncoding(b, logger),
		Transactional:   b.(Transactional),
	}
}

func (e *StorageEncoding) containsNonPrintableChars(key string) bool {
	idx := strings.IndexFunc(key, func(c rune) bool {
		return !unicode.IsPrint(c)
	})

	return idx != -1
}

func (e *StorageEncoding) Put(ctx context.Context, entry *Entry) error {
	if !utf8.ValidString(entry.Key) {
		return ErrNonUTF8
	}

	if e.containsNonPrintableChars(entry.Key) {
		return ErrNonPrintable
	}

	return e.Backend.Put(ctx, entry)
}

func (e *StorageEncoding) Delete(ctx context.Context, key string) error {
	if !utf8.ValidString(key) {
		return ErrNonUTF8
	}

	if e.containsNonPrintableChars(key) {
		return ErrNonPrintable
	}

	return e.Backend.Delete(ctx, key)
}

func (e *TransactionalStorageEncoding) Transaction(ctx context.Context, txns []*TxnEntry) error {
	for _, txn := range txns {
		if !utf8.ValidString(txn.Entry.Key) {
			return ErrNonUTF8
		}

		if e.containsNonPrintableChars(txn.Entry.Key) {
			return ErrNonPrintable
		}

	}

	return e.Transactional.Transaction(ctx, txns)
}
