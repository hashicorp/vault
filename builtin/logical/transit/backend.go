// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
)

const (
	operationPrefixTransit = "transit"

	// Minimum cache size for transit backend
	minCacheSize = 10
)

var ErrKeyTypeEntOnly = "key type %s is only available in enterprise versions of Vault"

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
			b.pathRotate(),
			b.pathRewrap(),
			b.pathWrappingKey(),
			b.pathImport(),
			b.pathImportVersion(),
			b.pathKeys(),
			b.pathListKeys(),
			b.pathBYOKExportKeys(),
			b.pathExportKeys(),
			b.pathKeysConfig(),
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
			b.pathConfigKeys(),
			b.pathCreateCsr(),
			b.pathImportCertChain(),
		},

		Secrets:        []*framework.Secret{},
		Invalidate:     b.invalidate,
		BackendType:    logical.TypeLogical,
		PeriodicFunc:   b.periodicFunc,
		InitializeFunc: b.initialize,
		Clean:          b.cleanup,
	}

	b.backendUUID = conf.BackendUUID
	b.initializeRotationQueue()

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
	b.setupEnt()

	return &b, nil
}

type backend struct {
	*framework.Backend
	entBackend

	lm *keysutil.LockManager
	// billingDataCounts tracks successful data protection operations
	// for this backend instance. It's intended for test assertions and avoids
	// cross-test/package contamination from global counters.
	billingDataCounts billing.DataProtectionCallCounts
	// Lock to make changes to any of the backend's cache configuration.
	configMutex          sync.RWMutex
	cacheSizeChanged     bool
	checkAutoRotateAfter time.Time
	autoRotateOnce       sync.Once
	rotationQueue        *rotationQueue
	backendUUID          string
}

type keyRotationEntry struct {
	rotateAt time.Time
	keyPath  string
}

func (kre keyRotationEntry) isZero() bool {
	return kre.rotateAt.IsZero() && kre.keyPath == ""
}

// rotationQueue stores information about which keys need to be rotated at what time.  It's a priority queue, which
// implements hash.Interface.
type rotationQueue []keyRotationEntry

var _ heap.Interface = &rotationQueue{}

func (rq rotationQueue) Len() int { return len(rq) }
func (rq rotationQueue) Less(i, j int) bool {
	if i < 0 || j < 0 || i >= rq.Len() || j >= rq.Len() { // If out of bounds, don't switch
		return false
	}
	return rq[i].rotateAt.Before(rq[j].rotateAt)
}

func (rq rotationQueue) Swap(i, j int) {
	if i < 0 || j < 0 || i >= len(rq) || j >= len(rq) {
		return
	}
	rq[i], rq[j] = rq[j], rq[i]
}

func (rq *rotationQueue) Push(kre any) {
	if kre.(keyRotationEntry).isZero() {
		return
	}
	*rq = append(*rq, kre.(keyRotationEntry))
}

func (rq *rotationQueue) Pop() any {
	if len(*rq) == 0 {
		return keyRotationEntry{}
	}
	old := *rq
	n := len(old)
	item := (*rq)[n-1]
	*rq = old[0 : n-1]
	return item
}

func (rq rotationQueue) Peek() keyRotationEntry {
	if len(rq) == 0 {
		return keyRotationEntry{}
	} else {
		return rq[0]
	}
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

// incrementBillingCounts atomically increments the transit billing data counts
func (b *backend) incrementBillingCounts(ctx context.Context, count uint64) error {
	// If we are a test, we need to increment this testing structure to verify the counts are correct.
	if b.billingDataCounts.Transit != nil {
		b.billingDataCounts.Transit.Add(count)
	}

	// Write billling data
	return b.ConsumptionBillingManager.WriteBillingData(ctx, "transit", map[string]interface{}{
		"count": count,
	})
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
			b.configMutex.RUnlock()
			return nil, false, err
		}
		if currentCacheSize != storedCacheSize {
			err = b.lm.InitCache(storedCacheSize)
			if err != nil {
				b.configMutex.RUnlock()
				return nil, false, err
			}
		}
		// Release the read lock and acquire the write lock
		b.configMutex.RUnlock()
		b.configMutex.Lock()
		defer b.configMutex.Unlock()
		b.cacheSizeChanged = false
	} else {
		b.configMutex.RUnlock()
	}
	p, _, err := b.lm.GetPolicy(ctx, polReq, rand)
	if err != nil {
		return p, false, err
	}

	if p != nil && p.Type.IsEnterpriseOnly() && !constants.IsEnterprise {
		p.Unlock()
		return nil, false, fmt.Errorf(ErrKeyTypeEntOnly, p.Type)
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

	b.invalidateEnt(ctx, key)
}

// periodicFunc is a central collection of functions that run on an interval.
// Anything that should be called regularly can be placed within this method.
func (b *backend) periodicFunc(ctx context.Context, req *logical.Request) error {
	// These operations ensure the auto-rotate only happens once simultaneously.
	var err error
	didAutoRotate := false
	autoRotateOnceFn := func() {
		err = b.autoRotateKeys(ctx, req)
		didAutoRotate = true
	}
	b.autoRotateOnce.Do(autoRotateOnceFn)
	if didAutoRotate {
		b.autoRotateOnce = sync.Once{}
	}

	if err != nil {
		return err
	}

	return b.periodicFuncEnt(ctx, req)
}

// autoRotateKeys retrieves all transit keys and rotates those which have an
// auto rotate period defined which has passed. This operation only happens
// on primary nodes and performance secondary nodes which have a local mount.
func (b *backend) autoRotateKeys(ctx context.Context, req *logical.Request) error {
	// Early exit if not a primary or performance secondary with a local mount.
	if b.System().ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary)) {
		return nil
	}

	// Collect errors in a multierror to ensure a single failure doesn't prevent
	// all keys from being rotated.
	var errs *multierror.Error

	// If we haven't initialized yet, exit out
	if b.rotationQueue == nil {
		errs = multierror.Append(errs, errors.New("auto-rotation can not run because backend is not initialized"))
		return errs.ErrorOrNil()
	}

	// Between once-per-hour rotations, check to see if any keys that need rotating are in the heap
	if b.rotationQueue.Len() == 0 {
		// No keys are scheduled for rotation this hour
	} else {
		for !b.rotationQueue.Peek().isZero() && b.rotationQueue.Peek().rotateAt.Before(time.Now()) {
			kre := heap.Pop(b.rotationQueue).(keyRotationEntry)
			_, err := b.rotateByPath(ctx, req, kre.keyPath)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}

	// Only check for autorotation once an hour to avoid unnecessarily iterating
	// over all keys too frequently.
	if time.Now().Before(b.checkAutoRotateAfter) {
		return nil
	}
	b.checkAutoRotateAfter = time.Now().Add(1 * time.Hour)

	// Retrieve all keys and loop over them to check if they need to be rotated.
	keys, err := req.Storage.List(ctx, "policy/")
	if err != nil {
		return err
	}

	for _, key := range keys {
		kre, err := b.rotateByPath(ctx, req, key)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		if !kre.isZero() {
			heap.Push(b.rotationQueue, kre)
		}
	}

	return errs.ErrorOrNil()
}

func (b *backend) rotateByPath(ctx context.Context, req *logical.Request, keyPath string) (keyRotationEntry, error) {
	p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage:     req.Storage,
		Name:        keyPath,
		WriteLocked: true,
	}, b.GetRandomReader())
	if err != nil {
		return keyRotationEntry{}, err
	}

	// If the policy is nil, move onto the next one.
	if p == nil {
		return keyRotationEntry{}, nil
	}

	// rotateIfRequired properly acquires/releases the lock on p
	kre, err := b.rotateIfRequired(ctx, req, keyPath, p)
	if err != nil {
		return kre, err
	}

	p.Unlock()
	return kre, err
}

// rotateIfRequired rotates a key if it is due for autorotation.  If it isn't due for autorotation, but will be due
// soon (within an hour), it returns a keyRotationEntry.
func (b *backend) rotateIfRequired(ctx context.Context, req *logical.Request, key string, p *keysutil.Policy) (kre keyRotationEntry, err error) {
	// If the key is imported, it can only be rotated from within Vault if allowed.
	if p.Imported && !p.AllowImportedKeyRotation {
		return keyRotationEntry{}, nil
	}

	// If the policy's automatic rotation period is 0, it should not
	// automatically rotate.
	if p.AutoRotatePeriod == 0 {
		return keyRotationEntry{}, nil
	}

	// We can't auto-rotate managed keys
	if p.Type == keysutil.KeyType_MANAGED_KEY {
		return keyRotationEntry{}, nil
	}

	// Retrieve the latest version of the policy and determine if it is time to rotate.
	latestKey := p.Keys[strconv.Itoa(p.LatestVersion)]
	autoRotateAt := latestKey.CreationTime.Add(p.AutoRotatePeriod)
	if time.Now().After(autoRotateAt) {
		if b.Logger().IsDebug() {
			b.Logger().Debug("automatically rotating key", "key", key)
		}
		return keyRotationEntry{}, p.Rotate(ctx, req.Storage, b.GetRandomReader())
	} else {
		// Check if it will be time to rotate the key within the next hour
		if time.Now().Add(1 * time.Hour).After(autoRotateAt) {
			kre = keyRotationEntry{
				rotateAt: autoRotateAt,
				keyPath:  key,
			}
			return kre, nil
		}
	}

	return keyRotationEntry{}, nil
}

func (b *backend) initialize(ctx context.Context, request *logical.InitializationRequest) error {
	return b.initializeEnt(ctx, request)
}

func (b *backend) initializeRotationQueue() {
	rotationQueue := &rotationQueue{}
	heap.Init(rotationQueue)
	b.rotationQueue = rotationQueue
}

func (b *backend) cleanup(ctx context.Context) {
	b.cleanupEnt(ctx)
}
