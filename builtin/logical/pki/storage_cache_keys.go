package pki

import (
	"context"
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/vault/sdk/logical"
)

type keyStorageCache struct {
	lock      sync.RWMutex
	invalid   bool
	idSet     map[keyID]bool
	nameIDMap map[string]keyID
	entries   *lru.Cache
}

func InitKeyStorageCache() *keyStorageCache {
	var ret keyStorageCache
	ret.invalid = true

	var err error
	ret.entries, err = lru.New(32)
	if err != nil {
		panic(err)
	}

	return &ret
}

func (c *keyStorageCache) Invalidate(op func() error) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.invalid = true

	if op != nil {
		return op()
	}

	return nil
}

func (c *keyStorageCache) reloadOnInvalidation(ctx context.Context, s logical.Storage) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.invalid {
		return nil
	}

	c.idSet = make(map[keyID]bool)
	c.nameIDMap = make(map[string]keyID)

	// Clear the LRU; this is necessary as some entries might've been deleted.
	c.entries.Purge()

	// List all keys which exist.
	strList, err := s.List(ctx, keyPrefix)
	if err != nil {
		return err
	}

	// Reset the key and name caches, populating the LRU.
	for _, keyIdStr := range strList {
		keyId := keyID(keyIdStr)

		// Fetch the specified key; it might've been deleted since the
		// list and thus returns an empty entry.
		rawEntry, err := s.Get(ctx, keyPrefix+keyIdStr)
		if err != nil {
			return err
		}
		if rawEntry == nil {
			continue
		}

		var entry keyEntry
		if err := rawEntry.DecodeJSON(&entry); err != nil {
			return err
		}

		c.idSet[keyId] = true
		if len(entry.Name) > 0 {
			c.nameIDMap[entry.Name] = keyId
		}

		// Greedily add this entry to the LRU.
		c.entries.Add(keyId, &entry)
	}

	c.invalid = false
	return nil
}

func (c *keyStorageCache) listKeys(ctx context.Context, s logical.Storage) ([]keyID, error) {
	needUnlock := true
	c.lock.RLock()
	defer func() {
		if needUnlock {
			c.lock.RUnlock()
		}
	}()

	if c.invalid {
		// Release our read lock so we can race to grab a write lock.
		c.lock.RUnlock()
		needUnlock = false

		if err := c.reloadOnInvalidation(ctx, s); err != nil {
			return nil, err
		}

		// Now re-read-lock.
		c.lock.RLock()
		needUnlock = true
	}

	// Now we can service the request as expected.
	result := make([]keyID, 0, len(c.idSet))
	for entry := range c.idSet {
		result = append(result, entry)
	}

	return result, nil
}

func (c *keyStorageCache) keyWithID(ctx context.Context, s logical.Storage, id keyID) (bool, error) {
	needUnlock := true
	c.lock.RLock()
	defer func() {
		if needUnlock {
			c.lock.RUnlock()
		}
	}()

	if c.invalid {
		// Release our read lock so we can race to grab a write lock.
		c.lock.RUnlock()
		needUnlock = false

		if err := c.reloadOnInvalidation(ctx, s); err != nil {
			return false, err
		}

		// Now re-read-lock.
		c.lock.RLock()
		needUnlock = true
	}

	present, ok := c.idSet[id]
	return ok && present, nil
}

func (c *keyStorageCache) keyWithName(ctx context.Context, s logical.Storage, name string) (keyID, error) {
	needUnlock := true
	c.lock.RLock()
	defer func() {
		if needUnlock {
			c.lock.RUnlock()
		}
	}()

	if c.invalid {
		// Release our read lock so we can race to grab a write lock.
		c.lock.RUnlock()
		needUnlock = false

		if err := c.reloadOnInvalidation(ctx, s); err != nil {
			return KeyRefNotFound, err
		}

		// Now re-read-lock.
		c.lock.RLock()
		needUnlock = true
	}

	keyId, ok := c.nameIDMap[name]
	if !ok || len(keyId) == 0 {
		return KeyRefNotFound, fmt.Errorf("unable to find PKI key for reference: %v", name)
	}

	return keyId, nil
}

func (c *keyStorageCache) fetchKeyById(ctx context.Context, s logical.Storage, keyId keyID) (*keyEntry, error) {
	needUnlock := true

	c.lock.RLock()
	defer func() {
		if needUnlock {
			c.lock.RUnlock()
		}
	}()

	if c.invalid {
		// Release our read lock so we can race to grab a write lock.
		c.lock.RUnlock()
		needUnlock = false

		if err := c.reloadOnInvalidation(ctx, s); err != nil {
			return nil, err
		}

		// Now re-read-lock.
		c.lock.RLock()
		needUnlock = true
	}

	// Now we can service the request as expected.
	if haveId, ok := c.idSet[keyId]; !ok || !haveId {
		return nil, fmt.Errorf("pki key id %v does not exist", keyId)
	}

	if entry, ok := c.entries.Get(keyId); ok && entry != nil {
		e := entry.(*keyEntry)
		return e, nil
	}

	// Otherwise, if it doesn't exist, fetch it and add it to the LRU. We
	// once again have to upgrade our read lock to a write lock.
	c.lock.RUnlock()
	needUnlock = false
	c.lock.Lock()
	defer c.lock.Unlock()

	rawEntry, err := s.Get(ctx, keyPrefix+keyId.String())
	if err != nil {
		return nil, err
	}
	if rawEntry == nil {
		return nil, fmt.Errorf("pki key id %s does not exist", keyId)
	}

	var entry keyEntry
	if err := rawEntry.DecodeJSON(&entry); err != nil {
		return nil, err
	}

	c.idSet[keyId] = true
	if len(entry.Name) > 0 {
		c.nameIDMap[entry.Name] = keyId
	}

	// Add this entry to the LRU.
	c.entries.Add(keyId, &entry)
	copiedEntry := entry
	return &copiedEntry, nil
}
