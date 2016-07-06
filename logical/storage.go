package logical

import (
	"fmt"

	"github.com/hashicorp/vault/helper/jsonutil"
)

// Storage is the way that logical backends are able read/write data.
type Storage interface {
	List(prefix string) ([]string, error)
	Get(string) (*StorageEntry, error)
	Put(*StorageEntry) error
	Delete(string) error
}

// StorageEntry is the entry for an item in a Storage implementation.
type StorageEntry struct {
	Key   string
	Value []byte
}

// DecodeJSON decodes the 'Value' present in StorageEntry.
func (e *StorageEntry) DecodeJSON(out interface{}) error {
	return jsonutil.DecodeJSON(e.Value, out)
}

// StorageEntryJSON creates a StorageEntry with a JSON-encoded value.
func StorageEntryJSON(k string, v interface{}) (*StorageEntry, error) {
	encodedBytes, err := jsonutil.EncodeJSON(v)
	if err != nil {
		return nil, fmt.Errorf("failed to encode storage entry: %v", err)
	}

	return &StorageEntry{
		Key:   k,
		Value: encodedBytes,
	}, nil
}
