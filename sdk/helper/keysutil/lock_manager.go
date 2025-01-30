// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keysutil

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	shared                   = false
	exclusive                = true
	currentConvergentVersion = 3
)

var errNeedExclusiveLock = errors.New("an exclusive lock is needed for this operation")

// PolicyRequest holds values used when requesting a policy. Most values are
// only used during an upsert.
type PolicyRequest struct {
	// The storage to use
	Storage logical.Storage

	// The name of the policy
	Name string

	// The key type
	KeyType KeyType

	// The key size for variable key size algorithms
	KeySize int

	// Whether it should be derived
	Derived bool

	// Whether to enable convergent encryption
	Convergent bool

	// Whether to allow export
	Exportable bool

	// Whether to upsert
	Upsert bool

	// Whether to allow plaintext backup
	AllowPlaintextBackup bool

	// How frequently the key should automatically rotate
	AutoRotatePeriod time.Duration

	// AllowImportedKeyRotation indicates whether an imported key may be rotated by Vault
	AllowImportedKeyRotation bool

	// Indicates whether a private or public key is imported/upserted
	IsPrivateKey bool

	// The UUID of the managed key, if using one
	ManagedKeyUUID string

	// ParameterSet indicates the parameter set to use with ML-DSA and SLH-DSA keys
	ParameterSet string

	// HybridConfig contains the key types and parameters for hybrid keys
	HybridConfig HybridKeyConfig
}

type HybridKeyConfig struct {
	PQCKeyType KeyType
	ECKeyType  KeyType
}

type LockManager struct {
	useCache bool
	cache    Cache
	keyLocks []*locksutil.LockEntry
}

func NewLockManager(useCache bool, cacheSize int) (*LockManager, error) {
	// determine the type of cache to create
	var cache Cache
	switch {
	case !useCache:
	case cacheSize < 0:
		return nil, errors.New("cache size must be greater or equal to zero")
	case cacheSize == 0:
		cache = NewTransitSyncMap()
	case cacheSize > 0:
		newLRUCache, err := NewTransitLRU(cacheSize)
		if err != nil {
			return nil, errwrap.Wrapf("failed to create cache: {{err}}", err)
		}
		cache = newLRUCache
	}

	lm := &LockManager{
		useCache: useCache,
		cache:    cache,
		keyLocks: locksutil.CreateLocks(),
	}

	return lm, nil
}

func (lm *LockManager) GetCacheSize() int {
	if !lm.useCache {
		return 0
	}
	return lm.cache.Size()
}

func (lm *LockManager) GetUseCache() bool {
	return lm.useCache
}

func (lm *LockManager) InvalidatePolicy(name string) {
	if lm.useCache {
		lm.cache.Delete(name)
	}
}

func (lm *LockManager) InitCache(cacheSize int) error {
	if lm.useCache {
		switch {
		case cacheSize < 0:
			return errors.New("cache size must be greater or equal to zero")
		case cacheSize == 0:
			lm.cache = NewTransitSyncMap()
		case cacheSize > 0:
			newLRUCache, err := NewTransitLRU(cacheSize)
			if err != nil {
				return errwrap.Wrapf("failed to create cache: {{err}}", err)
			}
			lm.cache = newLRUCache
		}
	}
	return nil
}

// RestorePolicy acquires an exclusive lock on the policy name and restores the
// given policy along with the archive.
func (lm *LockManager) RestorePolicy(ctx context.Context, storage logical.Storage, name, backup string, force bool) error {
	backupBytes, err := base64.StdEncoding.DecodeString(backup)
	if err != nil {
		return err
	}

	var keyData KeyData
	err = jsonutil.DecodeJSON(backupBytes, &keyData)
	if err != nil {
		return err
	}

	// Set a different name if desired
	if name != "" {
		keyData.Policy.Name = name
	}

	name = keyData.Policy.Name

	// Grab the exclusive lock as we'll be modifying disk
	lock := locksutil.LockForKey(lm.keyLocks, name)
	lock.Lock()
	defer lock.Unlock()

	var ok bool
	var pRaw interface{}

	// If the policy is in cache and 'force' is not specified, error out. Anywhere
	// that would put it in the cache will also be protected by the mutex above,
	// so we don't need to re-check the cache later.
	if lm.useCache {
		pRaw, ok = lm.cache.Load(name)
		if ok && !force {
			return fmt.Errorf("key %q already exists", name)
		}
	}

	// Conditionally look up the policy from storage, depending on the use of
	// 'force' and if the policy was found in cache.
	//
	// - If was not found in cache and we are not using 'force', look for it in
	// storage. If found, error out.
	//
	// - If it was found in cache and we are using 'force', pRaw will not be nil
	// and we do not look the policy up from storage
	//
	// - If it was found in cache and we are not using 'force', we should have
	// returned above with error
	var p *Policy
	if pRaw == nil {
		p, err = lm.getPolicyFromStorage(ctx, storage, name)
		if err != nil {
			return err
		}
		if p != nil && !force {
			return fmt.Errorf("key %q already exists", name)
		}
	}

	// If both pRaw and p above are nil and 'force' is specified, we don't need to
	// grab policy locks as we have ensured it doesn't already exist, so there
	// will be no races as nothing else has this pointer. If 'force' was not used,
	// an error would have been returned by now if the policy already existed
	if pRaw != nil {
		p = pRaw.(*Policy)
	}
	if p != nil {
		p.l.Lock()
		defer p.l.Unlock()
	}

	// Restore the archived keys
	if keyData.ArchivedKeys != nil {
		err = keyData.Policy.storeArchive(ctx, storage, keyData.ArchivedKeys)
		if err != nil {
			return errwrap.Wrapf(fmt.Sprintf("failed to restore archived keys for key %q: {{err}}", name), err)
		}
	}

	// Mark that policy as a restored key
	keyData.Policy.RestoreInfo = &RestoreInfo{
		Time:    time.Now(),
		Version: keyData.Policy.LatestVersion,
	}

	// Restore the policy. This will also attempt to adjust the archive.
	err = keyData.Policy.Persist(ctx, storage)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("failed to restore the policy %q: {{err}}", name), err)
	}

	keyData.Policy.l = new(sync.RWMutex)

	// Update the cache to contain the restored policy
	if lm.useCache {
		lm.cache.Store(name, keyData.Policy)
	}
	return nil
}

func (lm *LockManager) BackupPolicy(ctx context.Context, storage logical.Storage, name string) (string, error) {
	var p *Policy
	var err error

	// Backup writes information about when the backup took place, so we get an
	// exclusive lock here
	lock := locksutil.LockForKey(lm.keyLocks, name)
	lock.Lock()
	defer lock.Unlock()

	var ok bool
	var pRaw interface{}

	if lm.useCache {
		pRaw, ok = lm.cache.Load(name)
	}
	if ok {
		p = pRaw.(*Policy)
		p.l.Lock()
		defer p.l.Unlock()
	} else {
		// If the policy doesn't exit in storage, error out
		p, err = lm.getPolicyFromStorage(ctx, storage, name)
		if err != nil {
			return "", err
		}
		if p == nil {
			return "", fmt.Errorf("key %q not found", name)
		}
	}

	if atomic.LoadUint32(&p.deleted) == 1 {
		return "", fmt.Errorf("key %q not found", name)
	}

	backup, err := p.Backup(ctx, storage)
	if err != nil {
		return "", err
	}

	return backup, nil
}

// When the function returns, if caching was disabled, the Policy's lock must
// be unlocked when the caller is done (and it should not be re-locked).
func (lm *LockManager) GetPolicy(ctx context.Context, req PolicyRequest, rand io.Reader) (retP *Policy, retUpserted bool, retErr error) {
	var p *Policy
	var err error
	var ok bool
	var pRaw interface{}

	// Check if it's in our cache. If so, return right away.
	if lm.useCache {
		pRaw, ok = lm.cache.Load(req.Name)
	}
	if ok {
		p = pRaw.(*Policy)
		if atomic.LoadUint32(&p.deleted) == 1 {
			return nil, false, nil
		}
		return p, false, nil
	}

	// We're not using the cache, or it wasn't found; get an exclusive lock.
	// This ensures that any other process writing the actual storage will be
	// finished before we load from storage.
	lock := locksutil.LockForKey(lm.keyLocks, req.Name)
	lock.Lock()

	// If we are using the cache, defer the lock unlock; otherwise we will
	// return from here with the lock still held.
	cleanup := func() {
		switch {
		// If using the cache we always unlock, the caller locks the policy
		// themselves
		case lm.useCache:
			lock.Unlock()

		// If not using the cache, if we aren't returning a policy the caller
		// doesn't have a lock, so we must unlock
		case retP == nil:
			lock.Unlock()
		}
	}

	// Check the cache again
	if lm.useCache {
		pRaw, ok = lm.cache.Load(req.Name)
	}
	if ok {
		p = pRaw.(*Policy)
		if atomic.LoadUint32(&p.deleted) == 1 {
			cleanup()
			return nil, false, nil
		}
		retP = p
		cleanup()
		return
	}

	// Load it from storage
	p, err = lm.getPolicyFromStorage(ctx, req.Storage, req.Name)
	if err != nil {
		cleanup()
		return nil, false, err
	}
	// We don't need to lock the policy as there would be no other holders of
	// the pointer

	if p == nil {
		// This is the only place we upsert a new policy, so if upsert is not
		// specified, or the lock type is wrong, unlock before returning
		if !req.Upsert {
			cleanup()
			return nil, false, nil
		}

		// We create the policy here, then at the end we do a LoadOrStore. If
		// it's been loaded since we last checked the cache, we return an error
		// to the user to let them know that their request can't be satisfied
		// because we don't know if the parameters match.

		switch req.KeyType {
		case KeyType_AES128_GCM96, KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305:
			if req.Convergent && !req.Derived {
				cleanup()
				return nil, false, fmt.Errorf("convergent encryption requires derivation to be enabled")
			}

		case KeyType_ECDSA_P256, KeyType_ECDSA_P384, KeyType_ECDSA_P521:
			if req.Derived || req.Convergent {
				cleanup()
				return nil, false, fmt.Errorf("key derivation and convergent encryption not supported for keys of type %v", req.KeyType)
			}

		case KeyType_ED25519:
			if req.Convergent {
				cleanup()
				return nil, false, fmt.Errorf("convergent encryption not supported for keys of type %v", req.KeyType)
			}

		case KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096:
			if req.Derived || req.Convergent {
				cleanup()
				return nil, false, fmt.Errorf("key derivation and convergent encryption not supported for keys of type %v", req.KeyType)
			}
		case KeyType_HMAC:
			if req.Derived || req.Convergent {
				cleanup()
				return nil, false, fmt.Errorf("key derivation and convergent encryption not supported for keys of type %v", req.KeyType)
			}

		case KeyType_MANAGED_KEY:
			if req.Derived || req.Convergent {
				cleanup()
				return nil, false, fmt.Errorf("key derivation and convergent encryption not supported for keys of type %v", req.KeyType)
			}

		case KeyType_AES128_CMAC, KeyType_AES256_CMAC:
			if req.Derived || req.Convergent {
				cleanup()
				return nil, false, fmt.Errorf("key derivation and convergent encryption not supported for keys of type %v", req.KeyType)
			}

		case KeyType_ML_DSA:
			if req.Derived || req.Convergent {
				cleanup()
				return nil, false, fmt.Errorf("key derivation and convergent encryption not supported for keys of type %v", req.KeyType)
			}

		case KeyType_HYBRID:
			if req.Derived || req.Convergent {
				cleanup()
				return nil, false, fmt.Errorf("key derivation and convergent encryption not supported for keys of type %v", req.KeyType)
			}

		default:
			cleanup()
			return nil, false, fmt.Errorf("unsupported key type %v", req.KeyType)
		}

		p = &Policy{
			l:                    new(sync.RWMutex),
			Name:                 req.Name,
			Type:                 req.KeyType,
			Derived:              req.Derived,
			Exportable:           req.Exportable,
			AllowPlaintextBackup: req.AllowPlaintextBackup,
			AutoRotatePeriod:     req.AutoRotatePeriod,
			KeySize:              req.KeySize,
			ParameterSet:         req.ParameterSet,
			HybridConfig:         req.HybridConfig,
		}

		if req.Derived {
			p.KDF = Kdf_hkdf_sha256
			if req.Convergent {
				p.ConvergentEncryption = true
				// As of version 3 we store the version within each key, so we
				// set to -1 to indicate that the value in the policy has no
				// meaning. We still, for backwards compatibility, fall back to
				// this value if the key doesn't have one, which means it will
				// only be -1 in the case where every key version is >= 3
				p.ConvergentVersion = -1
			}
		}

		// Performs the actual persist and does setup
		if p.Type == KeyType_MANAGED_KEY {
			err = p.RotateManagedKey(ctx, req.Storage, req.ManagedKeyUUID)
		} else {
			err = p.Rotate(ctx, req.Storage, rand)
		}
		if err != nil {
			cleanup()
			return nil, false, err
		}

		if lm.useCache {
			lm.cache.Store(req.Name, p)
		} else {
			p.l = &lock.RWMutex
			p.writeLocked = true
		}

		// We don't need to worry about upgrading since it will be a new policy
		retP = p
		retUpserted = true
		cleanup()
		return
	}

	if p.NeedsUpgrade() {
		if err := p.Upgrade(ctx, req.Storage, rand); err != nil {
			cleanup()
			return nil, false, err
		}
	}

	if lm.useCache {
		lm.cache.Store(req.Name, p)
	} else {
		p.l = &lock.RWMutex
		p.writeLocked = true
	}

	retP = p
	cleanup()
	return
}

func (lm *LockManager) ImportPolicy(ctx context.Context, req PolicyRequest, key []byte, rand io.Reader) error {
	var p *Policy
	var err error
	var ok bool
	var pRaw interface{}

	// Check if it's in our cache
	if lm.useCache {
		pRaw, ok = lm.cache.Load(req.Name)
	}
	if ok {
		p = pRaw.(*Policy)
		if atomic.LoadUint32(&p.deleted) == 1 {
			return nil
		}
	}

	// We're not using the cache, or it wasn't found; get an exclusive lock.
	// This ensures that any other process writing the actual storage will be
	// finished before we load from storage.
	lock := locksutil.LockForKey(lm.keyLocks, req.Name)
	lock.Lock()
	defer lock.Unlock()

	// Load it from storage
	p, err = lm.getPolicyFromStorage(ctx, req.Storage, req.Name)
	if err != nil {
		return err
	}

	if p == nil {
		p = &Policy{
			l:                        new(sync.RWMutex),
			Name:                     req.Name,
			Type:                     req.KeyType,
			Derived:                  req.Derived,
			Exportable:               req.Exportable,
			AllowPlaintextBackup:     req.AllowPlaintextBackup,
			AutoRotatePeriod:         req.AutoRotatePeriod,
			AllowImportedKeyRotation: req.AllowImportedKeyRotation,
			Imported:                 true,
		}
	}

	err = p.ImportPublicOrPrivate(ctx, req.Storage, key, req.IsPrivateKey, rand)
	if err != nil {
		return fmt.Errorf("error importing key: %s", err)
	}

	if lm.useCache {
		lm.cache.Store(req.Name, p)
	}

	return nil
}

func (lm *LockManager) DeletePolicy(ctx context.Context, storage logical.Storage, name string) error {
	var p *Policy
	var err error
	var ok bool
	var pRaw interface{}

	// We may be writing to disk, so grab an exclusive lock. This prevents bad
	// behavior when the cache is turned off. We also lock the shared policy
	// object to make sure no requests are in flight.
	lock := locksutil.LockForKey(lm.keyLocks, name)
	lock.Lock()
	defer lock.Unlock()

	if lm.useCache {
		pRaw, ok = lm.cache.Load(name)
	}
	if ok {
		p = pRaw.(*Policy)
		p.l.Lock()
		defer p.l.Unlock()
	}

	if p == nil {
		p, err = lm.getPolicyFromStorage(ctx, storage, name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("could not delete key; not found")
		}
	}

	if !p.DeletionAllowed {
		return fmt.Errorf("deletion is not allowed for this key")
	}

	atomic.StoreUint32(&p.deleted, 1)

	if lm.useCache {
		lm.cache.Delete(name)
	}

	err = storage.Delete(ctx, "policy/"+name)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error deleting key %q: {{err}}", name), err)
	}

	err = storage.Delete(ctx, "archive/"+name)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error deleting key %q archive: {{err}}", name), err)
	}

	return nil
}

func (lm *LockManager) getPolicyFromStorage(ctx context.Context, storage logical.Storage, name string) (*Policy, error) {
	return LoadPolicy(ctx, storage, "policy/"+name)
}
