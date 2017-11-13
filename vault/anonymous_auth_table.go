package vault

import (
	"fmt"
	"strings"
	"sync"

	"github.com/armon/go-radix"
	"github.com/hashicorp/vault/logical"
)

// There are few things in here that might seem a bit counterintuitive, so
// a few notes on design choices.
// First: The use of a map to an empty interface. Right now, every mount point
// will display the exact same information anonymously. But, in the future, we
// might want to allow displaying more customized information per mount point.
// Map'ing to an empty structure allows us to add more to the data structure
// in the future if we want to, as opposed to serializing a list into the
// storage backend. Additionally, maps have improved lookup times compared to
// iterating through a slice -- we don't want to incur an O(M*N) lookup cost
// when iterating through every mount backend and every entry in the allowed list.
// Second, we're treating wildcard mounts differently from non-wildcard mounts.
// The reason is lookup semantics -- by storing wildcards in a radix tree, we
// can use optimized radix prefix-based lookups to see if a wildcard prefix
// matches a given mount point.
// They're stored separately (rather than including as a property in the
// authTableEntry struct) so that users can whitelist, e.g., both ldap and ldap*
// then un-whitelist ldap* and leave ldap whitelisted, and to have that play
// nicely with the data structures

const (
	anonAuthTableWildcardEntry = "allowed-wildcard-mounts"
	anonAuthTableEntry         = "allowed-mounts"
	anonAuthTableSubPath       = "anonymous-auth-table-config/"
)

type authTableEntry struct {
}

type AnonymousAuthTableConfig struct {
	AllowedWildcardMounts *radix.Tree
	AllowedMounts         map[string]*authTableEntry

	view         *BarrierView
	sync.RWMutex `json:"-"`
}

func (a *AnonymousAuthTableConfig) add(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	a.Lock()
	defer a.Unlock()

	if a.AllowedWildcardMounts == nil {
		a.AllowedWildcardMounts = radix.New()
	}
	if a.AllowedMounts == nil {
		a.AllowedMounts = make(map[string]*authTableEntry)
	}

	var entry *logical.StorageEntry
	var err error

	switch {
	case strings.HasSuffix(path, "*"):
		path = path[:len(path)-1]
		a.AllowedWildcardMounts.Insert(path, &authTableEntry{})
		entry, err = logical.StorageEntryJSON(anonAuthTableWildcardEntry, a.AllowedWildcardMounts.ToMap())
	default:
		path = sanitizeMountPath(path)
		a.AllowedMounts[path] = &authTableEntry{}
		entry, err = logical.StorageEntryJSON(anonAuthTableEntry, a.AllowedMounts)
	}

	if err != nil {
		return fmt.Errorf("failed to persist anonymous auth table entry: %v", err)
	}

	if err := a.view.Put(entry); err != nil {
		return fmt.Errorf("failed to persist anonymous auth table entry: %v", err)
	}

	return nil
}

func (a *AnonymousAuthTableConfig) remove(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	a.Lock()
	defer a.Unlock()

	var entry *logical.StorageEntry
	var err error

	switch {
	case strings.HasSuffix(path, "*"):
		path = path[:len(path)-1]
		_, _ = a.AllowedWildcardMounts.Delete(path)
		entry, err = logical.StorageEntryJSON(anonAuthTableWildcardEntry, a.AllowedWildcardMounts.ToMap())
	default:
		path = sanitizeMountPath(path)
		delete(a.AllowedMounts, path)
		entry, err = logical.StorageEntryJSON(anonAuthTableEntry, a.AllowedMounts)
	}

	if err != nil {
		return fmt.Errorf("failed to persist anonymous auth table config: %v", err)
	}

	if err := a.view.Put(entry); err != nil {
		return fmt.Errorf("failed to persist anonymous auth table config: %v", err)
	}

	return nil
}

func (a *AnonymousAuthTableConfig) read(path string) (*authTableEntry, error) {
	a.RLock()
	defer a.RUnlock()

	switch {
	case strings.HasSuffix(path, "*"):
		ent, ok := a.AllowedWildcardMounts.Get(path)
		if !ok {
			return nil, fmt.Errorf("didn't find entry for path %s", path)
		}
		return ent.(*authTableEntry), nil
	default:
		path = sanitizeMountPath(path)
		ent, ok := a.AllowedMounts[path]
		if !ok {
			return nil, fmt.Errorf("didn't find entry for path %s", path)
		}
		return ent, nil
	}
}

func (a *AnonymousAuthTableConfig) apply(entries []*MountEntry) (*logical.Response, error) {
	a.RLock()
	defer a.RUnlock()
	resp := &logical.Response{
		Data: make(map[string]interface{}),
	}
	for _, entry := range entries {
		path := entry.Path
		if _, ok := a.AllowedMounts[path]; ok {
			resp.Data[path] = map[string]interface{}{
				"type": entry.Type,
			}
			// Don't really need LongestPrefix here, just need something like AnyPrefix,
			// but LongestPrefix is probably good enough
		} else if _, _, ok := a.AllowedWildcardMounts.LongestPrefix(path); ok {
			resp.Data[path] = map[string]interface{}{
				"type": entry.Type,
			}
		}
	}

	return resp, nil
}

func (a *AnonymousAuthTableConfig) readAll() (map[string]*authTableEntry, error) {
	a.RLock()
	defer a.RUnlock()

	resp := make(map[string]*authTableEntry)
	for k, v := range a.AllowedMounts {
		resp[k] = v
	}
	for k, v := range a.AllowedWildcardMounts.ToMap() {
		resp[k+"*"] = v.(*authTableEntry)
	}

	return resp, nil
}

func (c *Core) setupAnonymousAuthTableConfig() error {
	view := c.systemBarrierView.SubView(anonAuthTableSubPath)
	out, err := view.Get(anonAuthTableEntry)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	allowedMounts := make(map[string]*authTableEntry)
	if out != nil {
		err = out.DecodeJSON(&allowedMounts)
		if err != nil {
			return err
		}
	}

	out, err = view.Get(anonAuthTableWildcardEntry)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}
	allowedWildcardMounts := make(map[string]interface{})
	if out != nil {
		err = out.DecodeJSON(&allowedWildcardMounts)
		if err != nil {
			return err
		}
	}

	c.anonAuthTable = &AnonymousAuthTableConfig{
		AllowedWildcardMounts: radix.NewFromMap(allowedWildcardMounts),
		AllowedMounts:         allowedMounts,
		view:                  view,
	}
	return nil
}
