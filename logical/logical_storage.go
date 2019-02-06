package logical

import (
	"context"

	"github.com/hashicorp/vault/physical"
)

type LogicalStorage struct {
	underlying physical.Backend
}

func (s *LogicalStorage) Get(ctx context.Context, key string) (*StorageEntry, error) {
	entry, err := s.underlying.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	return &StorageEntry{
		Key:      entry.Key,
		Value:    entry.Value,
		SealWrap: entry.SealWrap,
	}, nil
}

func (s *LogicalStorage) Put(ctx context.Context, entry *StorageEntry) error {
	return s.underlying.Put(ctx, &physical.Entry{
		Key:      entry.Key,
		Value:    entry.Value,
		SealWrap: entry.SealWrap,
	})
}

func (s *LogicalStorage) Delete(ctx context.Context, key string) error {
	return s.underlying.Delete(ctx, key)
}

func (s *LogicalStorage) List(ctx context.Context, prefix string) ([]string, error) {
	return s.underlying.List(ctx, prefix)
}

func (s *LogicalStorage) Underlying() physical.Backend {
	return s.underlying
}

func NewLogicalStorage(underlying physical.Backend) *LogicalStorage {
	return &LogicalStorage{
		underlying: underlying,
	}
}
