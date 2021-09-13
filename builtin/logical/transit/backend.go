package transit

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// Minimum cache size for transit backend
const minCacheSize = 10

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend(ctx, conf)
	if err != nil {
		return nil, err
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend(ctx context.Context, conf *logical.BackendConfig) (*backend, error) {
	var b backend
	b.Backend = &framework.Backend{
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"archive/",
				"policy/",
			},
		},

		Paths: []*framework.Path{
			// Rotate/Config needs to come before Keys
			// as the handler is greedy
			b.pathConfig(),
			b.pathRotate(),
			b.pathRewrap(),
			b.pathKeys(),
			b.pathListKeys(),
			b.pathExportKeys(),
			b.pathEncrypt(),
			b.pathDecrypt(),
			b.pathDatakey(),
			b.pathRandom(),
			b.pathHash(),
			b.pathHMAC(),
			b.pathSign(),
			b.pathVerify(),
			b.pathBackup(),
			b.pathRestore(),
			b.pathTrim(),
			b.pathCacheConfig(),
		},

		Secrets:     []*framework.Secret{},
		Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}

	// determine cacheSize to use. Defaults to 0 which means unlimited
	cacheSize := 0
	useCache := !conf.System.CachingDisabled()
	if useCache {
		var err error
		cacheSize, err = GetCacheSizeFromStorage(ctx, conf.StorageView)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving cache size from storage: %w", err)
		}

		if cacheSize != 0 && cacheSize < minCacheSize {
			b.Logger().Warn("size %d is less than minimum %d. Cache size is set to %d", cacheSize, minCacheSize, minCacheSize)
			cacheSize = minCacheSize
		}
	}

	var err error
	b.lm, err = keysutil.NewLockManager(useCache, cacheSize)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

type backend struct {
	*framework.Backend
	lm *keysutil.LockManager
	// Lock to make changes to any of the backend's cache configuration.
	configMutex sync.RWMutex
	cacheSizeChanged bool
}

func GetCacheSizeFromStorage(ctx context.Context, s logical.Storage) (int, error) {
	size := 0
	entry, err := s.Get(ctx, "config/cache")
	if err != nil {
		return 0, err
	}
	if entry != nil {
		var storedCache configCache
		if err := entry.DecodeJSON(&storedCache); err != nil {
			return 0, err
		}
		size = storedCache.Size
	}
	return size, nil
}

// Update cache size and get policy
func (b *backend) GetPolicy(ctx context.Context, polReq keysutil.PolicyRequest, rand io.Reader) (retP *keysutil.Policy, retUpserted bool, retErr error) {
	// Acquire read lock to read cacheSizeChanged
	b.configMutex.RLock()
	if b.lm.GetUseCache() && b.cacheSizeChanged {
		var err error
		currentCacheSize := b.lm.GetCacheSize()
		storedCacheSize, err := GetCacheSizeFromStorage(ctx, polReq.Storage)
		if err != nil {
			return nil, false, err
		}
		if currentCacheSize != storedCacheSize {
			err = b.lm.InitCache(storedCacheSize)
			if err != nil {
				return nil, false, err
			}
		}
		// Release the read lock and acquire the write lock
		b.configMutex.RUnlock()
		b.configMutex.Lock()
		defer b.configMutex.Unlock()
		b.cacheSizeChanged = false
	}
	p, _, err := b.lm.GetPolicy(ctx, polReq, rand)
	if err != nil {
		return p, false, err
	}
	return p, true, nil
}

func (b *backend) invalidate(ctx context.Context, key string) {
	if b.Logger().IsDebug() {
		b.Logger().Debug("invalidating key", "key", key)
	}
	switch {
	case strings.HasPrefix(key, "policy/"):
		name := strings.TrimPrefix(key, "policy/")
		b.lm.InvalidatePolicy(name)
	case strings.HasPrefix(key, "cache-config/"):
		// Acquire the lock to set the flag to indicate that cache size needs to be refreshed from storage
		b.configMutex.Lock()
		defer b.configMutex.Unlock()
		b.cacheSizeChanged = true
	}
}
