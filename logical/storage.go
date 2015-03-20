package logical

import (
	"bytes"
	"encoding/json"
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

func (e *StorageEntry) DecodeJSON(out interface{}) error {
	return json.Unmarshal(e.Value, out)
}

// StorageEntryJSON creates a StorageEntry with a JSON-encoded value.
func StorageEntryJSON(k string, v interface{}) (*StorageEntry, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return &StorageEntry{
		Key:   k,
		Value: buf.Bytes(),
	}, nil
}
