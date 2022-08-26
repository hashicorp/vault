package logical

import (
	"context"
	"errors"
	"strings"
)

type StorageView struct {
	storage Storage
	prefix  string
}

var ErrRelativePath = errors.New("relative paths not supported")

func NewStorageView(storage Storage, prefix string) *StorageView {
	return &StorageView{
		storage: storage,
		prefix:  prefix,
	}
}

// logical.Storage impl.
func (s *StorageView) List(ctx context.Context, prefix string) ([]string, error) {
	if err := s.SanityCheck(prefix); err != nil {
		return nil, err
	}
	return s.storage.List(ctx, s.ExpandKey(prefix))
}

// logical.Storage impl.
func (s *StorageView) Get(ctx context.Context, key string) (*StorageEntry, error) {
	if err := s.SanityCheck(key); err != nil {
		return nil, err
	}
	entry, err := s.storage.Get(ctx, s.ExpandKey(key))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	entry.Key = s.TruncateKey(entry.Key)

	return &StorageEntry{
		Key:      entry.Key,
		Value:    entry.Value,
		SealWrap: entry.SealWrap,
	}, nil
}

// logical.Storage impl.
func (s *StorageView) Put(ctx context.Context, entry *StorageEntry) error {
	if entry == nil {
		return errors.New("cannot write nil entry")
	}

	if err := s.SanityCheck(entry.Key); err != nil {
		return err
	}

	expandedKey := s.ExpandKey(entry.Key)

	nested := &StorageEntry{
		Key:      expandedKey,
		Value:    entry.Value,
		SealWrap: entry.SealWrap,
	}

	return s.storage.Put(ctx, nested)
}

// logical.Storage impl.
func (s *StorageView) Delete(ctx context.Context, key string) error {
	if err := s.SanityCheck(key); err != nil {
		return err
	}

	expandedKey := s.ExpandKey(key)

	return s.storage.Delete(ctx, expandedKey)
}

func (s *StorageView) Prefix() string {
	return s.prefix
}

// SubView constructs a nested sub-view using the given prefix
func (s *StorageView) SubView(prefix string) *StorageView {
	sub := s.ExpandKey(prefix)
	return &StorageView{storage: s.storage, prefix: sub}
}

// SanityCheck is used to perform a sanity check on a key
func (s *StorageView) SanityCheck(key string) error {
	if strings.Contains(key, "..") {
		return ErrRelativePath
	}
	return nil
}

// ExpandKey is used to expand to the full key path with the prefix
func (s *StorageView) ExpandKey(suffix string) string {
	return s.prefix + suffix
}

// TruncateKey is used to remove the prefix of the key
func (s *StorageView) TruncateKey(full string) string {
	return strings.TrimPrefix(full, s.prefix)
}
