package transit

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/logical"
)

const (
	shared    = false
	exclusive = true
)

type lockManager struct {
	// A lock for each named key
	locks map[string]*sync.RWMutex

	// A mutex for the map itself
	locksMutex sync.RWMutex

	// If caching is enabled, the map of name to in-memory policy cache
	cache map[string]*Policy

	// Used for global locking, and as the cache map mutex
	globalMutex sync.RWMutex
}

func newLockManager(cacheDisabled bool) *lockManager {
	lm := &lockManager{
		locks: map[string]*sync.RWMutex{},
	}
	if !cacheDisabled {
		lm.cache = map[string]*Policy{}
	}
	return lm
}

func (lm *lockManager) CacheActive() bool {
	return lm.cache != nil
}

func (lm *lockManager) LockAll(name string) {
	lm.globalMutex.Lock()
	lm.LockPolicy(name, exclusive)
}

func (lm *lockManager) UnlockAll(name string) {
	lm.UnlockPolicy(name, exclusive)
	lm.globalMutex.Unlock()
}

func (lm *lockManager) LockPolicy(name string, writeLock bool) {
	lm.locksMutex.RLock()
	lock := lm.locks[name]
	if lock != nil {
		// We want to give this up before locking the lock, but it's safe --
		// the only time we ever write to a value in this map is the first time
		// we access the value, so it won't be changing out from under us
		lm.locksMutex.RUnlock()
		if writeLock {
			lock.Lock()
		} else {
			lock.RLock()
		}
		return
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
		if writeLock {
			lock.Lock()
		} else {
			lock.RLock()
		}
		return
	}

	lock = &sync.RWMutex{}
	lm.locks[name] = lock
	lm.locksMutex.Unlock()
	if writeLock {
		lock.Lock()
	} else {
		lock.RLock()
	}
}

func (lm *lockManager) UnlockPolicy(name string, writeLock bool) {
	lm.locksMutex.RLock()
	lock := lm.locks[name]
	lm.locksMutex.RUnlock()

	if writeLock {
		lock.Unlock()
	} else {
		lock.RUnlock()
	}
}

func (lm *lockManager) GetPolicy(storage logical.Storage, name string) (*Policy, bool, error) {
	p, lt, _, err := lm.getPolicyCommon(storage, name, false, false)
	return p, lt, err
}

func (lm *lockManager) GetPolicyUpsert(storage logical.Storage, name string, derived bool) (*Policy, bool, bool, error) {
	return lm.getPolicyCommon(storage, name, true, derived)
}

// When the function returns, a lock will be held on the policy if err == nil.
// The type of lock will be indicated by the return value. It is the caller's
// responsibility to unlock.
func (lm *lockManager) getPolicyCommon(storage logical.Storage, name string, upsert, derived bool) (p *Policy, lockType bool, upserted bool, err error) {
	// If we are using a cache, lock it now to avoid having to do really
	// complicated lock juggling as we call various functions. We'll also defer
	// the store into the cache.
	lockType = shared
	lm.LockPolicy(name, shared)

	if lm.CacheActive() {
		lm.globalMutex.RLock()
		p = lm.cache[name]
		if p != nil {
			lm.globalMutex.RUnlock()
			return
		}
		lm.globalMutex.RUnlock()

		// When we return, since we didn't have the policy in the cache, if
		// there was no error, write the value in.
		defer func() {
			lm.globalMutex.Lock()
			defer lm.globalMutex.Unlock()
			// Make sure a policy didn't appear. If so, it will only be set if
			// there was no error, so now just clear the error and return that
			// policy.
			exp := lm.cache[name]
			if exp != nil {
				upserted = false
				err = nil
				p = exp
				return
			}

			if err == nil {
				lm.cache[name] = p
			}
		}()
	}

	p, err = lm.getStoredPolicy(storage, name)
	if err != nil {
		lm.UnlockPolicy(name, shared)
		return
	}

	if p == nil {
		if !upsert {
			lm.UnlockPolicy(name, shared)
			return
		}

		// Get an exlusive lock; on success, check again to ensure that no
		// policy exists. Note that if we are using a cache we will already be
		// serializing this entire code path and it's currently the only one
		// that generates policies, so we don't need to check the cache here;
		// simply checking the disk again is sufficient.
		lm.UnlockPolicy(name, shared)
		lockType = exclusive
		lm.LockPolicy(name, exclusive)

		p, err = lm.getStoredPolicy(storage, name)
		if err != nil {
			defer lm.UnlockPolicy(name, exclusive)
			return
		}
		if p != nil {
			return
		}

		upserted = true

		p = &Policy{
			Name:       name,
			CipherMode: "aes-gcm",
			Derived:    derived,
		}
		if derived {
			p.KDFMode = kdfMode
		}

		err = p.rotate(storage)
		if err != nil {
			defer lm.UnlockPolicy(name, exclusive)
			p = nil
		}

		// We don't need to worry about upgrading since it will be a new policy
		return
	}

	if p.needsUpgrade() {
		lm.UnlockPolicy(name, shared)
		lockType = exclusive
		lm.LockPolicy(name, exclusive)

		// Reload the policy with the write lock to ensure we still need the upgrade
		p, err = lm.getStoredPolicy(storage, name)
		if err != nil {
			defer lm.UnlockPolicy(name, exclusive)
			return
		}
		if p == nil {
			defer lm.UnlockPolicy(name, exclusive)
			err = fmt.Errorf("error reloading policy for upgrade")
			return
		}

		if !p.needsUpgrade() {
			// Already happened, return the newly loaded policy
			return
		}

		err = p.upgrade(storage)
		if err != nil {
			defer lm.UnlockPolicy(name, exclusive)
		}
	}

	return
}

func (lm *lockManager) DeletePolicy(storage logical.Storage, name string) error {
	lm.LockAll(name)
	defer lm.UnlockAll(name)

	var p *Policy
	var err error

	if lm.CacheActive() {
		p = lm.cache[name]
		if p == nil {
			return fmt.Errorf("could not delete policy; not found")
		}
	} else {
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

// When this function returns it's the responsibility of the caller to call
// UnlockPolicy if err is nil and policy is not nil
func (lm *lockManager) RefreshPolicy(storage logical.Storage, name string) (p *Policy, err error) {
	lm.LockPolicy(name, exclusive)

	if lm.CacheActive() {
		p = lm.cache[name]
		if p != nil {
			return
		}
		err = fmt.Errorf("could not refresh policy; not found")
		defer lm.UnlockPolicy(name, exclusive)
		return
	}

	p, err = lm.getStoredPolicy(storage, name)
	if err != nil {
		defer lm.UnlockPolicy(name, exclusive)
		return
	}

	if p == nil {
		err = fmt.Errorf("could not refresh policy; not found")
		defer lm.UnlockPolicy(name, exclusive)
	}

	if p.needsUpgrade() {
		err = p.upgrade(storage)
		if err != nil {
			defer lm.UnlockPolicy(name, exclusive)
		}
	}

	return
}

func (lm *lockManager) getStoredPolicy(storage logical.Storage, name string) (*Policy, error) {
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
		Keys: KeyEntryMap{},
	}
	err = json.Unmarshal(raw.Value, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}
