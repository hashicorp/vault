package physical

import (
	"context"
	"errors"
	"strings"
)

var ErrRelativePath = errors.New("relative paths not supported")

// View represents a prefixed view of a physical backend
type View struct {
	backend Backend
	prefix  string
}

// Verify View satisfies the correct interfaces
var _ Backend = (*View)(nil)

// NewView takes an underlying physical backend and returns
// a view of it that can only operate with the given prefix.
func NewView(backend Backend, prefix string) *View {
	return &View{
		backend: backend,
		prefix:  prefix,
	}
}

// List the contents of the prefixed view
func (v *View) List(ctx context.Context, prefix string) ([]string, error) {
	if err := v.sanityCheck(prefix); err != nil {
		return nil, err
	}
	return v.backend.List(ctx, v.expandKey(prefix))
}

// Get the key of the prefixed view
func (v *View) Get(ctx context.Context, key string) (*Entry, error) {
	if err := v.sanityCheck(key); err != nil {
		return nil, err
	}
	entry, err := v.backend.Get(ctx, v.expandKey(key))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	entry.Key = v.truncateKey(entry.Key)

	return &Entry{
		Key:   entry.Key,
		Value: entry.Value,
	}, nil
}

// Put the entry into the prefix view
func (v *View) Put(ctx context.Context, entry *Entry) error {
	if err := v.sanityCheck(entry.Key); err != nil {
		return err
	}

	nested := &Entry{
		Key:   v.expandKey(entry.Key),
		Value: entry.Value,
	}
	return v.backend.Put(ctx, nested)
}

// Delete the entry from the prefix view
func (v *View) Delete(ctx context.Context, key string) error {
	if err := v.sanityCheck(key); err != nil {
		return err
	}
	return v.backend.Delete(ctx, v.expandKey(key))
}

// sanityCheck is used to perform a sanity check on a key
func (v *View) sanityCheck(key string) error {
	if strings.Contains(key, "..") {
		return ErrRelativePath
	}
	return nil
}

// expandKey is used to expand to the full key path with the prefix
func (v *View) expandKey(suffix string) string {
	return v.prefix + suffix
}

// truncateKey is used to remove the prefix of the key
func (v *View) truncateKey(full string) string {
	return strings.TrimPrefix(full, v.prefix)
}
