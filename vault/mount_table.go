package vault

import (
	"context"
	"reflect"
	"sort"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// MountTable is used to represent the internal mount table
type MountTable struct {
	Type    string        `json:"type"`
	Entries []*MountEntry `json:"entries"`

	// For SPv2 backend this will have a reference to the
	// StoragePacker.  In the transition MountTable points back to
	// Core to use Core's methods.
	core   *Core      `-`
	logger log.Logger `-`
}

// Now that mount tables aren't simply a data type, tests that check
// equality can't just use reflect.DeepEqual
func matchingMountTables(table1 *MountTable, table2 *MountTable) bool {
	if table1.Type != table2.Type {
		return false
	}
	if !reflect.DeepEqual(table1.Entries, table2.Entries) {
		return false
	}
	return true
}

// shallowClone returns a copy of the mount table that
// keeps the MountEntry locations, so as not to invalidate
// other locations holding pointers. Care needs to be taken
// if modifying entries rather than modifying the table itself
func (t *MountTable) shallowClone() *MountTable {
	mt := &MountTable{
		Type:    t.Type,
		Entries: make([]*MountEntry, len(t.Entries)),
		core:    t.core,
		logger:  t.logger,
	}
	for i, e := range t.Entries {
		mt.Entries[i] = e
	}
	return mt
}

//
// Low-level operations on the entries
//

// setTaint is used to set the taint on given entry Accepts either the mount
// entry's path or namespace + path, i.e. <ns-path>/secret/ or <ns-path>/token/
func (t *MountTable) setTaint(ctx context.Context, path string, value bool) (*MountEntry, error) {
	n := len(t.Entries)
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		if entry := t.Entries[i]; entry.Path == path && entry.Namespace().ID == ns.ID {
			t.Entries[i].Tainted = value
			return t.Entries[i], nil
		}
	}
	return nil, nil
}

// remove is used to remove a given path entry; returns the entry that was
// removed
func (t *MountTable) remove(ctx context.Context, path string) (*MountEntry, error) {
	n := len(t.Entries)
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		if entry := t.Entries[i]; entry.Path == path && entry.Namespace().ID == ns.ID {
			t.Entries[i], t.Entries[n-1] = t.Entries[n-1], nil
			t.Entries = t.Entries[:n-1]
			return entry, nil
		}
	}
	return nil, nil
}

func (t *MountTable) find(ctx context.Context, path string) (*MountEntry, error) {
	n := len(t.Entries)
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		if entry := t.Entries[i]; entry.Path == path && entry.Namespace().ID == ns.ID {
			return entry, nil
		}
	}
	return nil, nil
}

// sortEntriesByPath sorts the entries in the table by path and returns the
// table; this is useful for tests
func (t *MountTable) sortEntriesByPath() *MountTable {
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Path < t.Entries[j].Path
	})
	return t
}

// sortEntriesByPath sorts the entries in the table by path and returns the
// table; this is useful for tests
func (t *MountTable) sortEntriesByPathDepth() *MountTable {
	sort.Slice(t.Entries, func(i, j int) bool {
		return len(strings.Split(t.Entries[i].Namespace().Path+t.Entries[i].Path, "/")) < len(strings.Split(t.Entries[j].Namespace().Path+t.Entries[j].Path, "/"))
	})
	return t
}

//
// High-level operations that persist state
//

// Return a replacement MountTable to use
// (TODO: this return value will go away after the conversion to SPv2.)
func (t *MountTable) addMountEntry(ctx context.Context, entry *MountEntry, updateStorage bool) (*MountTable, error) {
	newTable := t.shallowClone()
	newTable.Entries = append(newTable.Entries, entry)
	if updateStorage {
		if err := t.core.persistMounts(ctx, newTable, &entry.Local); err != nil {
			t.logger.Error("failed to update mount table", "error", err)
			if err == logical.ErrReadOnly && t.core.perfStandby {
				return nil, err
			}
			return nil, logical.CodedError(500, "failed to update mount table")
		}
	}
	return newTable, nil
}

// removeMountEntry is used to remove an entry from the mount table
// Return a replacement MountTable to use
// (TODO: this return value will go away after the conversion to SPv2.)
func (t *MountTable) removeMountEntry(ctx context.Context, path string, updateStorage bool) (*MountTable, error) {
	// Remove the entry from the mount table
	newTable := t.shallowClone()
	entry, err := newTable.remove(ctx, path)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		t.logger.Error("nil entry found removing entry in mounts table", "path", path)
		return nil, logical.CodedError(500, "failed to remove entry in mounts table")
	}

	// When unmounting all entries the JSON code will load back up from storage
	// as a nil slice, which kills tests...just set it nil explicitly
	if len(newTable.Entries) == 0 {
		newTable.Entries = nil
	}

	if updateStorage {
		// Update the mount table
		if err := t.core.persistMounts(ctx, newTable, &entry.Local); err != nil {
			t.logger.Error("failed to remove entry from mounts table", "error", err)
			return nil, logical.CodedError(500, "failed to remove entry from mounts table")
		}
	}
	return newTable, nil
}

// taintMountEntry is used to mark an entry in the mount table as tainted
// Performed in-place already.
func (t *MountTable) taintMountEntry(ctx context.Context, path string, updateStorage bool) error {
	// As modifying the taint of an entry affects shallow clones,
	// we simply use the original
	entry, err := t.setTaint(ctx, path, true)
	if err != nil {
		return err
	}
	if entry == nil {
		t.logger.Error("nil entry found tainting entry in mounts table", "path", path)
		return logical.CodedError(500, "failed to taint entry in mounts table")
	}

	if updateStorage {
		// Update the mount table
		if err := t.core.persistMounts(ctx, t, &entry.Local); err != nil {
			if err == logical.ErrReadOnly && t.core.perfStandby {
				return err
			}

			t.logger.Error("failed to taint entry in mounts table", "error", err)
			return logical.CodedError(500, "failed to taint entry in mounts table")
		}
	}

	return nil
}
