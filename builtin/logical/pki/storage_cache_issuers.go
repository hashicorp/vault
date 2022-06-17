package pki

import (
	"context"
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/vault/sdk/logical"
)

type issuerStorageCache struct {
	lock      sync.RWMutex
	invalid   bool
	idSet     map[issuerID]bool
	nameIDMap map[string]issuerID
	entries   *lru.Cache
	bundles   *lru.Cache
}

func InitIssuerStorageCache() *issuerStorageCache {
	var ret issuerStorageCache
	ret.invalid = true

	var err error

	ret.entries, err = lru.New(32)
	if err != nil {
		panic(err)
	}

	ret.bundles, err = lru.New(32)
	if err != nil {
		panic(err)
	}

	return &ret
}

func (c *issuerStorageCache) Invalidate(op func() error) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.invalid = true

	if op != nil {
		return op()
	}

	return nil
}

func (c *issuerStorageCache) reloadOnInvalidation(ctx context.Context, s logical.Storage) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.invalid {
		return nil
	}

	c.idSet = make(map[issuerID]bool)
	c.nameIDMap = make(map[string]issuerID)

	// Clear the LRU; this is necessary as some entries/bundles might've been deleted.
	c.entries.Purge()
	c.bundles.Purge()

	// List all issuers which exist.
	strList, err := s.List(ctx, issuerPrefix)
	if err != nil {
		return err
	}

	// Reset the issuer and name caches, populating the LRU.
	for _, issuerIdStr := range strList {
		issuerId := issuerID(issuerIdStr)

		// Fetch the specified issuer; it might've been deleted since the
		// list and thus returns an empty entry.
		rawEntry, err := s.Get(ctx, issuerPrefix+issuerIdStr)
		if err != nil {
			return err
		}
		if rawEntry == nil {
			continue
		}

		var entry issuerEntry
		if err := rawEntry.DecodeJSON(&entry); err != nil {
			return err
		}

		c.idSet[issuerId] = true
		if len(entry.Name) > 0 {
			c.nameIDMap[entry.Name] = issuerId
		}

		// Greedily add this entry to the LRU. We don't add bundles here.
		c.entries.Add(issuerId, &entry)
	}

	c.invalid = false
	return nil
}

func (c *issuerStorageCache) listIssuers(ctx context.Context, s logical.Storage) ([]issuerID, error) {
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
	result := make([]issuerID, 0, len(c.idSet))
	for entry := range c.idSet {
		result = append(result, entry)
	}

	return result, nil
}

func (c *issuerStorageCache) issuerWithID(ctx context.Context, s logical.Storage, id issuerID) (bool, error) {
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

func (c *issuerStorageCache) issuerWithName(ctx context.Context, s logical.Storage, name string) (issuerID, error) {
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
			return IssuerRefNotFound, err
		}

		// Now re-read-lock.
		c.lock.RLock()
		needUnlock = true
	}

	issuerId, ok := c.nameIDMap[name]
	if !ok {
		return IssuerRefNotFound, fmt.Errorf("unable to find PKI issuer for reference: %v", name)
	}

	return issuerId, nil
}

func (c *issuerStorageCache) fetchIssuerById(ctx context.Context, s logical.Storage, issuerId issuerID) (*issuerEntry, error) {
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
	if haveId, ok := c.idSet[issuerId]; !ok || !haveId {
		return nil, fmt.Errorf("pki issuer id %v does not exist", issuerId)
	}

	if entry, ok := c.entries.Get(issuerId); ok && entry != nil {
		e := entry.(*issuerEntry)
		return e, nil
	}

	// Otherwise, if it doesn't exist, fetch it and add it to the LRU. We
	// once again have to upgrade our read lock to a write lock.
	c.lock.RUnlock()
	needUnlock = false
	c.lock.Lock()
	defer c.lock.Unlock()

	rawEntry, err := s.Get(ctx, issuerPrefix+issuerId.String())
	if err != nil {
		return nil, err
	}
	if rawEntry == nil {
		return nil, fmt.Errorf("pki issuer id %s does not exist", issuerId)
	}

	var entry issuerEntry
	if err := rawEntry.DecodeJSON(&entry); err != nil {
		return nil, err
	}

	c.idSet[issuerId] = true
	if len(entry.Name) > 0 {
		c.nameIDMap[entry.Name] = issuerId
	}

	// Add this entry to the LRU.
	c.entries.Add(issuerId, &entry)
	copiedEntry := entry
	return &copiedEntry, nil
}
