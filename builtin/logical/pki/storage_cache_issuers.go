package pki

import (
	"context"
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type issuerStorageCache struct {
	lock          sync.RWMutex
	invalid       bool
	idSet         map[issuerID]bool
	nameIDMap     map[string]issuerID
	entries       *lru.Cache
	bundles       *lru.Cache
	parsedBundles *lru.Cache
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

	ret.parsedBundles, err = lru.New(32)
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
	c.parsedBundles.Purge()

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

func (c *issuerStorageCache) fetchIssuerInfoById(ctx context.Context, s logical.Storage, b *backend, issuerId issuerID) (*issuerEntry, *certutil.CertBundle, *certutil.ParsedCertBundle, error) {
	// This method is more complex than the keys counterpart. Here, we handle
	// the logic for all three types of calls:
	//
	// 1. Issuer entry only,
	// 2. CertBundle as well,
	// 3. Upgrading all the way to a ParsedCertBundle.
	//
	// When present, we use the keyStorageCache to give us a version of the
	// bundle with the key pre-loaded and cache this, rather than potentially
	// storing a half-built bundle. This allows us to remove the private key
	// bits if they're not desired, rather than attempting to update the LRU
	// cache entry with new key material.
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
			return nil, nil, nil, err
		}

		// Now re-read-lock.
		c.lock.RLock()
		needUnlock = true
	}

	// Now we can service the request as expected.
	if haveId, ok := c.idSet[issuerId]; !ok || !haveId {
		return nil, nil, nil, fmt.Errorf("pki issuer id %v does not exist", issuerId)
	}

	entry, entryOk := c.entries.Get(issuerId)
	if b == nil && entryOk && entry != nil {
		// We can exit fast if we don't need the bundle. This only occurs when
		// the keyStorageCache value is nil.
		e := entry.(*issuerEntry)
		return e, nil, nil, nil
	}

	bundle, bundleOk := c.bundles.Get(issuerId)
	parsedBundle, parsedOk := c.parsedBundles.Get(issuerId)
	if bundleOk && bundle != nil && parsedOk && parsedBundle != nil {
		// If everything was loaded from cache, we're good.
		e := entry.(*issuerEntry)
		n := bundle.(*certutil.CertBundle)
		p := parsedBundle.(*certutil.ParsedCertBundle)
		return e, n, p, nil
	}

	if !entryOk || entry == nil {
		// Otherwise, if the entry doesn't exist, fetch it and add it to the
		// LRU. We once again have to upgrade our read lock to a write lock.
		c.lock.RUnlock()
		needUnlock = false
		c.lock.Lock()
		defer c.lock.Unlock()

		rawEntry, err := s.Get(ctx, issuerPrefix+issuerId.String())
		if err != nil {
			return nil, nil, nil, err
		}
		if rawEntry == nil {
			return nil, nil, nil, fmt.Errorf("pki issuer id %s does not exist", issuerId)
		}

		var storedEntry issuerEntry
		if err := rawEntry.DecodeJSON(&storedEntry); err != nil {
			return nil, nil, nil, err
		}

		c.idSet[issuerId] = true
		if len(storedEntry.Name) > 0 {
			c.nameIDMap[storedEntry.Name] = issuerId
		}

		// Add this entry to the LRU.
		c.entries.Add(issuerId, &storedEntry)
		copiedEntry := storedEntry
		entry = &copiedEntry
	}

	e := entry.(*issuerEntry)
	// If we don't have a key storage cache, don't bother building the cert
	// bundle entries.
	if b == nil {
		return e, nil, nil, nil
	}

	// Finally, we can finish building the cert bundle entries. Since we
	// either:
	//
	// 1. Missed a key entry and thus needed to invalidate our LRU copies,
	// 2. Missed either cert/parse copies,
	//
	// we always have to do this from scratch.
	var rawBundle certutil.CertBundle
	rawBundle.Certificate = e.Certificate
	rawBundle.CAChain = e.CAChain
	rawBundle.SerialNumber = e.SerialNumber
	if e.KeyID != keyID("") {
		keyEntry, err := b.keyCache.fetchKeyById(ctx, s, e.KeyID)
		if err != nil {
			return nil, nil, nil, err
		}

		rawBundle.PrivateKeyType = keyEntry.PrivateKeyType
		rawBundle.PrivateKey = keyEntry.PrivateKey
	}

	c.bundles.Add(issuerId, &rawBundle)
	copiedBundle := rawBundle
	bundle = &copiedBundle
	n := bundle.(*certutil.CertBundle)

	rawParsedBundle, err := parseCABundle(ctx, b, n)
	if err != nil {
		return nil, nil, nil, err
	}

	c.parsedBundles.Add(issuerId, rawParsedBundle)
	copiedParsedBundle := *rawParsedBundle
	parsedBundle = &copiedParsedBundle

	p := parsedBundle.(*certutil.ParsedCertBundle)
	return e, n, p, nil
}

func (c *issuerStorageCache) fetchIssuerById(ctx context.Context, s logical.Storage, issuerId issuerID) (*issuerEntry, error) {
	entry, _, _, err := c.fetchIssuerInfoById(ctx, s, nil, issuerId)
	return entry, err
}

func (c *issuerStorageCache) fetchCertBundleByIssuerId(ctx context.Context, s logical.Storage, b *backend, issuerId issuerID, loadKey bool) (*issuerEntry, *certutil.CertBundle, error) {
	entry, bundle, _, err := c.fetchIssuerInfoById(ctx, s, b, issuerId)

	if err != nil && !loadKey {
		bundle.PrivateKey = ""
	}

	return entry, bundle, err
}

func (c *issuerStorageCache) fetchParsedBundleByIssuerId(ctx context.Context, s logical.Storage, b *backend, issuerId issuerID, loadKey bool) (*issuerEntry, *certutil.CertBundle, *certutil.ParsedCertBundle, error) {
	return c.fetchIssuerInfoById(ctx, s, b, issuerId)
}
