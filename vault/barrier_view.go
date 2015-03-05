package vault

import "strings"

// BarrierView is used to wrap a barrier and ensure all access is automatically
// prefixed. This means that nothing outside of the given prefix can be
// accessed through the view, which is an additional layer of security when
// interacting with the security barrier. Conceptually this is like a
// "chroot" into the barrier.
type BarrierView struct {
	barrier SecurityBarrier
	prefix  string
}

// NewBarrierView takes an underlying security barrier and returns
// a view of it that can only operate with the given prefix.
func NewBarrierView(barrier SecurityBarrier, prefix string) *BarrierView {
	return &BarrierView{
		barrier: barrier,
		prefix:  prefix,
	}
}

// Put is used to insert or update an entry
func (v *BarrierView) Put(entry *Entry) error {
	nested := &Entry{
		Key:   v.expandKey(entry.Key),
		Value: entry.Value,
	}
	return v.barrier.Put(nested)
}

// Get is used to fetch an entry
func (v *BarrierView) Get(key string) (*Entry, error) {
	entry, err := v.barrier.Get(v.expandKey(key))
	if err != nil {
		return nil, err
	}
	if entry != nil {
		entry.Key = v.truncateKey(entry.Key)
	}
	return entry, nil
}

// Delete is used to permanently delete an entry
func (v *BarrierView) Delete(key string) error {
	return v.barrier.Delete(v.expandKey(key))
}

// List is used ot list all the keys under a given
// prefix, up to the next prefix.
func (v *BarrierView) List(prefix string) ([]string, error) {
	return v.barrier.List(v.expandKey(prefix))
}

// expandKey is used to expand to the full key path with the prefix
func (v *BarrierView) expandKey(suffix string) string {
	return v.prefix + suffix
}

// truncateKey is used to remove the prefix of the key
func (v *BarrierView) truncateKey(full string) string {
	return strings.TrimPrefix(full, v.prefix)
}
