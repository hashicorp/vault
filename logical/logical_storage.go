package logical

import (
	"context"
	"fmt"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/file"
	"github.com/hashicorp/vault/physical/inmem"
)

type LogicalType string

const (
	LogicalTypeInmem LogicalType = "inmem"
	LogicalTypeFile  LogicalType = "file"
)

type LogicalStorage struct {
	logicalType LogicalType
	underlying  physical.Backend
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

func (s *LogicalStorage) init() {
}

func NewLogicalStorage(logicalType LogicalType, config map[string]string, logger log.Logger) (*LogicalStorage, error) {
	s := &LogicalStorage{}
	var err error
	switch logicalType {
	case LogicalTypeInmem:
		s.underlying, err = inmem.NewInmem(nil, nil)
		if err != nil {
			return nil, err
		}
	case LogicalTypeFile:
		s.underlying, err = file.NewFileBackend(config, logger)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported logical type %q", logicalType)
	}

	return s, nil
}
