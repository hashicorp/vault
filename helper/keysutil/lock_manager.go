package keysutil

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/logical"
)

const (
	shared                   = false
	exclusive                = true
	currentConvergentVersion = 3
)

var (
	errNeedExclusiveLock = errors.New("an exclusive lock is needed for this operation")
)

// PolicyRequest holds values used when requesting a policy. Most values are
// only used during an upsert.
type PolicyRequest struct {
	// The storage to use
	Storage logical.Storage

	// The name of the policy
	Name string

	// The key type
	KeyType KeyType

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
}

type LockManager struct {
	useCache bool
	// If caching is enabled, the map of name to in-memory policy cache
	cache sync.Map

	keyLocks []*locksutil.LockEntry
}

func NewLockManager(cacheDisabled bool) *LockManager {
	lm := &LockManager{
		useCache: !cacheDisabled,
		keyLocks: locksutil.CreateLocks(),
	}
	return lm
}

func (lm *LockManager) CacheActive() bool {
	return lm.useCache
}

func (lm *LockManager) InvalidatePolicy(name string) {
	lm.cache.Delete(name)
}

// RestorePolicy acquires an exclusive lock on the policy name and restores the
// given policy along with the archive.
func (lm *LockManager) RestorePolicy(ctx context.Context, storage logical.Storage, name, backup string) error {
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

	// If the policy is in cache, error out. Anywhere that would put it in the
	// cache will also be protected by the mutex above, so we don't need to
	// re-check the cache later.
	_, ok := lm.cache.Load(name)
	if ok {
		return fmt.Errorf(fmt.Sprintf("key %q already exists", name))
	}

	// If the policy exists in storage, error out
	p, err := lm.getPolicyFromStorage(ctx, storage, name)
	if err != nil {
		return err
	}
	if p != nil {
		return fmt.Errorf(fmt.Sprintf("key %q already exists", name))
	}

	// We don't need to grab policy locks as we have ensured it doesn't already
	// exist, so there will be no races as nothing else has this pointer.

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
	lm.cache.Store(name, keyData.Policy)

	return nil
}

func (lm *LockManager) BackupPolicy(ctx context.Context, storage logical.Storage, name string) (string, error) {
	var p *Policy
	var err error

	// Backup writes information about when the bacup took place, so we get an
	// exclusive lock here
	lock := locksutil.LockForKey(lm.keyLocks, name)
	lock.Lock()
	defer lock.Unlock()

	pRaw, ok := lm.cache.Load(name)
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
			return "", fmt.Errorf(fmt.Sprintf("key %q not found", name))
		}
	}

	if atomic.LoadUint32(&p.deleted) == 1 {
		return "", fmt.Errorf(fmt.Sprintf("key %q not found", name))
	}

	backup, err := p.Backup(ctx, storage)
	if err != nil {
		return "", err
	}

	return backup, nil
}

// When the function returns, if caching was disabled, the Policy's lock must
// be unlocked when the caller is done (and it should not be re-locked).
func (lm *LockManager) GetPolicy(ctx context.Context, req PolicyRequest) (retP *Policy, retUpserted bool, retErr error) {
	var p *Policy
	var err error

	// Check if it's in our cache. If so, return right away.
	pRaw, ok := lm.cache.Load(req.Name)
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
	pRaw, ok = lm.cache.Load(req.Name)
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
		case KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305:
			if req.Convergent && !req.Derived {
				cleanup()
				return nil, false, fmt.Errorf("convergent encryption requires derivation to be enabled")
			}

		case KeyType_ECDSA_P256:
			if req.Derived || req.Convergent {
				cleanup()
				return nil, false, fmt.Errorf("key derivation and convergent encryption not supported for keys of type %v", req.KeyType)
			}

		case KeyType_ED25519:
			if req.Convergent {
				cleanup()
				return nil, false, fmt.Errorf("convergent encryption not supported for keys of type %v", req.KeyType)
			}

		case KeyType_RSA2048, KeyType_RSA4096:
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
		err = p.Rotate(ctx, req.Storage)
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
		if err := p.Upgrade(ctx, req.Storage); err != nil {
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

func (lm *LockManager) DeletePolicy(ctx context.Context, storage logical.Storage, name string) error {
	var p *Policy
	var err error

	// We may be writing to disk, so grab an exclusive lock. This prevents bad
	// behavior when the cache is turned off. We also lock the shared policy
	// object to make sure no requests are in flight.
	lock := locksutil.LockForKey(lm.keyLocks, name)
	lock.Lock()
	defer lock.Unlock()

	pRaw, ok := lm.cache.Load(name)
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

	lm.cache.Delete(name)

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
