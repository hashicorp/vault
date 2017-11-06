package keysutil

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
)

const (
	shared    = false
	exclusive = true
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
}

type LockManager struct {
	// A lock for each named key
	locks map[string]*sync.RWMutex

	// A mutex for the map itself
	locksMutex sync.RWMutex

	// If caching is enabled, the map of name to in-memory policy cache
	cache map[string]*Policy

	// Used for global locking, and as the cache map mutex
	cacheMutex sync.RWMutex
}

func NewLockManager(cacheDisabled bool) *LockManager {
	lm := &LockManager{
		locks: map[string]*sync.RWMutex{},
	}
	if !cacheDisabled {
		lm.cache = map[string]*Policy{}
	}
	return lm
}

func (lm *LockManager) CacheActive() bool {
	return lm.cache != nil
}

func (lm *LockManager) InvalidatePolicy(name string) {
	// Check if it's in our cache. If so, return right away.
	if lm.CacheActive() {
		lm.cacheMutex.Lock()
		defer lm.cacheMutex.Unlock()
		delete(lm.cache, name)
	}
}

func (lm *LockManager) policyLock(name string, lockType bool) *sync.RWMutex {
	lm.locksMutex.RLock()
	lock := lm.locks[name]
	if lock != nil {
		// We want to give this up before locking the lock, but it's safe --
		// the only time we ever write to a value in this map is the first time
		// we access the value, so it won't be changing out from under us
		lm.locksMutex.RUnlock()
		if lockType == exclusive {
			lock.Lock()
		} else {
			lock.RLock()
		}
		return lock
	}

	lm.locksMutex.RUnlock()
	lm.locksMutex.Lock()

	// Don't defer the unlock call because if we get a valid lock below we want
	// to release the lock mutex right away to avoid the possibility of
	// deadlock by trying to grab the second lock

	// Check to make sure it hasn't been created since
	lock = lm.locks[name]
	if lock != nil {
		lm.locksMutex.Unlock()
		if lockType == exclusive {
			lock.Lock()
		} else {
			lock.RLock()
		}
		return lock
	}

	lock = &sync.RWMutex{}
	lm.locks[name] = lock
	lm.locksMutex.Unlock()
	if lockType == exclusive {
		lock.Lock()
	} else {
		lock.RLock()
	}

	return lock
}

func (lm *LockManager) UnlockPolicy(lock *sync.RWMutex, lockType bool) {
	if lockType == exclusive {
		lock.Unlock()
	} else {
		lock.RUnlock()
	}
}

// Get the policy with a read lock. If we get an error saying an exclusive lock
// is needed (for instance, for an upgrade/migration), give up the read lock,
// call again with an exclusive lock, then swap back out for a read lock.
func (lm *LockManager) GetPolicyShared(storage logical.Storage, name string) (*Policy, *sync.RWMutex, error) {
	p, lock, _, err := lm.getPolicyCommon(PolicyRequest{
		Storage: storage,
		Name:    name,
	}, shared)
	if err == nil ||
		(err != nil && err != errNeedExclusiveLock) {
		return p, lock, err
	}

	// Try again while asking for an exlusive lock
	p, lock, _, err = lm.getPolicyCommon(PolicyRequest{
		Storage: storage,
		Name:    name,
	}, exclusive)
	if err != nil || p == nil || lock == nil {
		return p, lock, err
	}

	lock.Unlock()

	p, lock, _, err = lm.getPolicyCommon(PolicyRequest{
		Storage: storage,
		Name:    name,
	}, shared)
	return p, lock, err
}

// Get the policy with an exclusive lock
func (lm *LockManager) GetPolicyExclusive(storage logical.Storage, name string) (*Policy, *sync.RWMutex, error) {
	p, lock, _, err := lm.getPolicyCommon(PolicyRequest{
		Storage: storage,
		Name:    name,
	}, exclusive)
	return p, lock, err
}

// Get the policy with a read lock; if it returns that an exclusive lock is
// needed, retry. If successful, call one more time to get a read lock and
// return the value.
func (lm *LockManager) GetPolicyUpsert(req PolicyRequest) (*Policy, *sync.RWMutex, bool, error) {
	req.Upsert = true

	p, lock, _, err := lm.getPolicyCommon(req, shared)
	if err == nil ||
		(err != nil && err != errNeedExclusiveLock) {
		return p, lock, false, err
	}

	// Try again while asking for an exlusive lock
	p, lock, upserted, err := lm.getPolicyCommon(req, exclusive)
	if err != nil || p == nil || lock == nil {
		return p, lock, upserted, err
	}
	lock.Unlock()

	req.Upsert = false
	// Now get a shared lock for the return, but preserve the value of upserted
	p, lock, _, err = lm.getPolicyCommon(req, shared)

	return p, lock, upserted, err
}

// When the function returns, a lock will be held on the policy if err == nil.
// It is the caller's responsibility to unlock.
func (lm *LockManager) getPolicyCommon(req PolicyRequest, lockType bool) (*Policy, *sync.RWMutex, bool, error) {
	lock := lm.policyLock(req.Name, lockType)

	var p *Policy
	var err error

	// Check if it's in our cache. If so, return right away.
	if lm.CacheActive() {
		lm.cacheMutex.RLock()
		p = lm.cache[req.Name]
		if p != nil {
			lm.cacheMutex.RUnlock()
			return p, lock, false, nil
		}
		lm.cacheMutex.RUnlock()
	}

	// Load it from storage
	p, err = lm.getStoredPolicy(req.Storage, req.Name)
	if err != nil {
		lm.UnlockPolicy(lock, lockType)
		return nil, nil, false, err
	}

	if p == nil {
		// This is the only place we upsert a new policy, so if upsert is not
		// specified, or the lock type is wrong, unlock before returning
		if !req.Upsert {
			lm.UnlockPolicy(lock, lockType)
			return nil, nil, false, nil
		}

		if lockType != exclusive {
			lm.UnlockPolicy(lock, lockType)
			return nil, nil, false, errNeedExclusiveLock
		}

		switch req.KeyType {
		case KeyType_AES256_GCM96:
			if req.Convergent && !req.Derived {
				lm.UnlockPolicy(lock, lockType)
				return nil, nil, false, fmt.Errorf("convergent encryption requires derivation to be enabled")
			}

		case KeyType_ECDSA_P256:
			if req.Derived || req.Convergent {
				lm.UnlockPolicy(lock, lockType)
				return nil, nil, false, fmt.Errorf("key derivation and convergent encryption not supported for keys of type %v", req.KeyType)
			}

		case KeyType_ED25519:
			if req.Convergent {
				lm.UnlockPolicy(lock, lockType)
				return nil, nil, false, fmt.Errorf("convergent encryption not supported for keys of type %v", req.KeyType)
			}

		case KeyType_RSA2048, KeyType_RSA4096:
			if req.Derived || req.Convergent {
				lm.UnlockPolicy(lock, lockType)
				return nil, nil, false, fmt.Errorf("key derivation and convergent encryption not supported for keys of type %v", req.KeyType)
			}

		default:
			lm.UnlockPolicy(lock, lockType)
			return nil, nil, false, fmt.Errorf("unsupported key type %v", req.KeyType)
		}

		p = &Policy{
			Name:       req.Name,
			Type:       req.KeyType,
			Derived:    req.Derived,
			Exportable: req.Exportable,
		}
		if req.Derived {
			p.KDF = Kdf_hkdf_sha256
			p.ConvergentEncryption = req.Convergent
			p.ConvergentVersion = 2
		}

		err = p.Rotate(req.Storage)
		if err != nil {
			lm.UnlockPolicy(lock, lockType)
			return nil, nil, false, err
		}

		if lm.CacheActive() {
			// Since we didn't have the policy in the cache, if there was no
			// error, write the value in.
			lm.cacheMutex.Lock()
			defer lm.cacheMutex.Unlock()
			// Make sure a policy didn't appear. If so, it will only be set if
			// there was no error, so assume it's good and return that
			exp := lm.cache[req.Name]
			if exp != nil {
				return exp, lock, false, nil
			}
			if err == nil {
				lm.cache[req.Name] = p
			}
		}

		// We don't need to worry about upgrading since it will be a new policy
		return p, lock, true, nil
	}

	if p.NeedsUpgrade() {
		if lockType == shared {
			lm.UnlockPolicy(lock, lockType)
			return nil, nil, false, errNeedExclusiveLock
		}

		err = p.Upgrade(req.Storage)
		if err != nil {
			lm.UnlockPolicy(lock, lockType)
			return nil, nil, false, err
		}
	}

	if lm.CacheActive() {
		// Since we didn't have the policy in the cache, if there was no
		// error, write the value in.
		lm.cacheMutex.Lock()
		defer lm.cacheMutex.Unlock()
		// Make sure a policy didn't appear. If so, it will only be set if
		// there was no error, so assume it's good and return that
		exp := lm.cache[req.Name]
		if exp != nil {
			return exp, lock, false, nil
		}
		if err == nil {
			lm.cache[req.Name] = p
		}
	}

	return p, lock, false, nil
}

func (lm *LockManager) DeletePolicy(storage logical.Storage, name string) error {
	lm.cacheMutex.Lock()
	lock := lm.policyLock(name, exclusive)
	defer lock.Unlock()
	defer lm.cacheMutex.Unlock()

	var p *Policy
	var err error

	if lm.CacheActive() {
		p = lm.cache[name]
	}
	if p == nil {
		p, err = lm.getStoredPolicy(storage, name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("could not delete policy; not found")
		}
	}

	if !p.DeletionAllowed {
		return fmt.Errorf("deletion is not allowed for this policy")
	}

	err = storage.Delete("policy/" + name)
	if err != nil {
		return fmt.Errorf("error deleting policy %s: %s", name, err)
	}

	err = storage.Delete("archive/" + name)
	if err != nil {
		return fmt.Errorf("error deleting archive %s: %s", name, err)
	}

	if lm.CacheActive() {
		delete(lm.cache, name)
	}

	return nil
}

func (lm *LockManager) getStoredPolicy(storage logical.Storage, name string) (*Policy, error) {
	// Check if the policy already exists
	raw, err := storage.Get("policy/" + name)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	// Decode the policy
	policy := &Policy{
		Keys: keyEntryMap{},
	}
	err = jsonutil.DecodeJSON(raw.Value, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}
