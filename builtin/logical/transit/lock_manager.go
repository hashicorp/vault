package transit

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

// policyRequest holds values used when requesting a policy. Most values are
// only used during an upsert.
type policyRequest struct {
	// The storage to use
	storage logical.Storage

	// The name of the policy
	name string

	// The key type
	keyType KeyType

	// Whether it should be derived
	derived bool

	// Whether to enable convergent encryption
	convergent bool

	// Whether to upsert
	upsert bool
}

type lockManager struct {
	// A lock for each named key
	locks map[string]*sync.RWMutex

	// A mutex for the map itself
	locksMutex sync.RWMutex

	// If caching is enabled, the map of name to in-memory policy cache
	cache map[string]*policy

	// Used for global locking, and as the cache map mutex
	cacheMutex sync.RWMutex
}

func newLockManager(cacheDisabled bool) *lockManager {
	lm := &lockManager{
		locks: map[string]*sync.RWMutex{},
	}
	if !cacheDisabled {
		lm.cache = map[string]*policy{}
	}
	return lm
}

func (lm *lockManager) CacheActive() bool {
	return lm.cache != nil
}

func (lm *lockManager) policyLock(name string, lockType bool) *sync.RWMutex {
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

func (lm *lockManager) UnlockPolicy(lock *sync.RWMutex, lockType bool) {
	if lockType == exclusive {
		lock.Unlock()
	} else {
		lock.RUnlock()
	}
}

// Get the policy with a read lock. If we get an error saying an exclusive lock
// is needed (for instance, for an upgrade/migration), give up the read lock,
// call again with an exclusive lock, then swap back out for a read lock.
func (lm *lockManager) GetPolicyShared(storage logical.Storage, name string) (*policy, *sync.RWMutex, error) {
	p, lock, _, err := lm.getPolicyCommon(policyRequest{
		storage: storage,
		name:    name,
	}, shared)
	if err == nil ||
		(err != nil && err != errNeedExclusiveLock) {
		return p, lock, err
	}

	// Try again while asking for an exlusive lock
	p, lock, _, err = lm.getPolicyCommon(policyRequest{
		storage: storage,
		name:    name,
	}, exclusive)
	if err != nil || p == nil || lock == nil {
		return p, lock, err
	}

	lock.Unlock()

	p, lock, _, err = lm.getPolicyCommon(policyRequest{
		storage: storage,
		name:    name,
	}, shared)
	return p, lock, err
}

// Get the policy with an exclusive lock
func (lm *lockManager) GetPolicyExclusive(storage logical.Storage, name string) (*policy, *sync.RWMutex, error) {
	p, lock, _, err := lm.getPolicyCommon(policyRequest{
		storage: storage,
		name:    name,
	}, exclusive)
	return p, lock, err
}

// Get the policy with a read lock; if it returns that an exclusive lock is
// needed, retry. If successful, call one more time to get a read lock and
// return the value.
func (lm *lockManager) GetPolicyUpsert(req policyRequest) (*policy, *sync.RWMutex, bool, error) {
	req.upsert = true

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

	req.upsert = false
	// Now get a shared lock for the return, but preserve the value of upserted
	p, lock, _, err = lm.getPolicyCommon(req, shared)

	return p, lock, upserted, err
}

// When the function returns, a lock will be held on the policy if err == nil.
// It is the caller's responsibility to unlock.
func (lm *lockManager) getPolicyCommon(req policyRequest, lockType bool) (*policy, *sync.RWMutex, bool, error) {
	lock := lm.policyLock(req.name, lockType)

	var p *policy
	var err error

	// Check if it's in our cache. If so, return right away.
	if lm.CacheActive() {
		lm.cacheMutex.RLock()
		p = lm.cache[req.name]
		if p != nil {
			lm.cacheMutex.RUnlock()
			return p, lock, false, nil
		}
		lm.cacheMutex.RUnlock()
	}

	// Load it from storage
	p, err = lm.getStoredPolicy(req.storage, req.name)
	if err != nil {
		lm.UnlockPolicy(lock, lockType)
		return nil, nil, false, err
	}

	if p == nil {
		// This is the only place we upsert a new policy, so if upsert is not
		// specified, or the lock type is wrong, unlock before returning
		if !req.upsert {
			lm.UnlockPolicy(lock, lockType)
			return nil, nil, false, nil
		}

		if lockType != exclusive {
			lm.UnlockPolicy(lock, lockType)
			return nil, nil, false, errNeedExclusiveLock
		}

		switch req.keyType {
		case keyType_AES256_GCM96:
			if req.convergent && !req.derived {
				return nil, nil, false, fmt.Errorf("convergent encryption requires derivation to be enabled")
			}

		case keyType_ECDSA_P256:
			if req.derived || req.convergent {
				return nil, nil, false, fmt.Errorf("key derivation and convergent encryption not supported for keys of type %s", keyType_ECDSA_P256)
			}

		default:
			return nil, nil, false, fmt.Errorf("unsupported key type %v", req.keyType)
		}

		p = &policy{
			Name:    req.name,
			Type:    req.keyType,
			Derived: req.derived,
		}
		if req.derived {
			p.KDF = kdf_hkdf_sha256
			p.ConvergentEncryption = req.convergent
			p.ConvergentVersion = 2
		}

		err = p.rotate(req.storage)
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
			exp := lm.cache[req.name]
			if exp != nil {
				return exp, lock, false, nil
			}
			if err == nil {
				lm.cache[req.name] = p
			}
		}

		// We don't need to worry about upgrading since it will be a new policy
		return p, lock, true, nil
	}

	if p.needsUpgrade() {
		if lockType == shared {
			lm.UnlockPolicy(lock, lockType)
			return nil, nil, false, errNeedExclusiveLock
		}

		err = p.upgrade(req.storage)
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
		exp := lm.cache[req.name]
		if exp != nil {
			return exp, lock, false, nil
		}
		if err == nil {
			lm.cache[req.name] = p
		}
	}

	return p, lock, false, nil
}

func (lm *lockManager) DeletePolicy(storage logical.Storage, name string) error {
	lm.cacheMutex.Lock()
	lock := lm.policyLock(name, exclusive)
	defer lock.Unlock()
	defer lm.cacheMutex.Unlock()

	var p *policy
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

func (lm *lockManager) getStoredPolicy(storage logical.Storage, name string) (*policy, error) {
	// Check if the policy already exists
	raw, err := storage.Get("policy/" + name)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	// Decode the policy
	policy := &policy{
		Keys: keyEntryMap{},
	}
	err = jsonutil.DecodeJSON(raw.Value, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}
