package framework

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// WALPrefix is the prefix within Storage where WAL entries will be written.
const WALPrefix = "wal/"

type WALEntry struct {
	ID        string      `json:"-"`
	Kind      string      `json:"type"`
	Data      interface{} `json:"data"`
	CreatedAt int64       `json:"created_at"`
}

// PutWAL writes some data to the WAL.
//
// The kind parameter is used by the framework to allow users to store
// multiple kinds of WAL data and to easily disambiguate what data they're
// expecting.
//
// Data within the WAL that is uncommitted (CommitWAL hasn't be called)
// will be given to the rollback callback when an rollback operation is
// received, allowing the backend to clean up some partial states.
//
// The data must be JSON encodable.
//
// This returns a unique ID that can be used to reference this WAL data.
// WAL data cannot be modified. You can only add to the WAL and commit existing
// WAL entries.
func PutWAL(ctx context.Context, s logical.Storage, kind string, data interface{}) (string, error) {
	value, err := json.Marshal(&WALEntry{
		Kind:      kind,
		Data:      data,
		CreatedAt: time.Now().UTC().Unix(),
	})
	if err != nil {
		return "", err
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}

	return id, s.Put(ctx, &logical.StorageEntry{
		Key:   WALPrefix + id,
		Value: value,
	})
}

// GetWAL reads a specific entry from the WAL. If the entry doesn't exist,
// then nil value is returned.
//
// The kind, value, and error are returned.
func GetWAL(ctx context.Context, s logical.Storage, id string) (*WALEntry, error) {
	entry, err := s.Get(ctx, WALPrefix+id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var raw WALEntry
	if err := jsonutil.DecodeJSON(entry.Value, &raw); err != nil {
		return nil, err
	}
	raw.ID = id

	return &raw, nil
}

// DeleteWAL commits the WAL entry with the given ID. Once committed,
// it is assumed that the operation was a success and doesn't need to
// be rolled back.
func DeleteWAL(ctx context.Context, s logical.Storage, id string) error {
	return s.Delete(ctx, WALPrefix+id)
}

// ListWAL lists all the entries in the WAL.
func ListWAL(ctx context.Context, s logical.Storage) ([]string, error) {
	keys, err := s.List(ctx, WALPrefix)
	if err != nil {
		return nil, err
	}

	for i, k := range keys {
		keys[i] = strings.TrimPrefix(k, WALPrefix)
	}

	return keys, nil
}
