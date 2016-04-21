package transit

import (
	"sync"

	"github.com/hashicorp/vault/logical"
)

// policyCache implements CRUD operations with a simple locking cache of
// policies
type cachingPolicyCRUD struct {
	sync.RWMutex
	cache map[string]lockingPolicy
}

func newCachingPolicyCRUD() *cachingPolicyCRUD {
	return &cachingPolicyCRUD{
		cache: map[string]lockingPolicy{},
	}
}

func (p *cachingPolicyCRUD) getPolicy(storage logical.Storage, name string) (lockingPolicy, error) {
	// We don't defer this since we may need to give it up and get a write lock
	p.RLock()

	// First, see if we're in the cache -- if so, return that
	if p.cache[name] != nil {
		defer p.RUnlock()
		return p.cache[name], nil
	}

	// If we didn't find anything, we'll need to write into the cache, plus possibly
	// persist the entry, so lock the cache
	p.RUnlock()
	p.Lock()
	defer p.Unlock()

	// Check one more time to ensure that another process did not write during
	// our lock switcheroo.
	if p.cache[name] != nil {
		return p.cache[name], nil
	}

	return p.refreshPolicy(storage, name)
}

func (p *cachingPolicyCRUD) refreshPolicy(storage logical.Storage, name string) (lockingPolicy, error) {
	// Check once more to ensure it hasn't been added to the cache since the lock was acquired
	if p.cache[name] != nil {
		return p.cache[name], nil
	}

	// Note that we don't need to create the locking entry until the end,
	// because the policy wasn't in the cache so we don't know about it, and we
	// hold the cache lock so nothing else can be writing it in right now
	policy, err := fetchPolicyFromStorage(storage, name)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, nil
	}

	lp := &mutexLockingPolicy{
		policy: policy,
		mutex:  &sync.RWMutex{},
	}
	p.cache[name] = lp

	return lp, nil
}

// generatePolicy is used to create a new named policy with a randomly
// generated key. The caller should hold the write lock prior to calling this.
func (p *cachingPolicyCRUD) generatePolicy(storage logical.Storage, name string, derived bool) (lockingPolicy, error) {
	policy, err := generatePolicyCommon(p, storage, name, derived)
	if err != nil {
		return nil, err
	}

	// Now we need to check again in the cache to ensure the policy wasn't
	// created since we ran generatePolicy and then got the lock. A policy
	// being created holds a write lock until it's done (starting from this
	// point), so it'll be in the cache at this point.
	if lp := p.cache[name]; lp != nil {
		return lp, nil
	}

	lp := &mutexLockingPolicy{
		policy: policy,
		mutex:  &sync.RWMutex{},
	}
	p.cache[name] = lp

	// Return the policy
	return lp, nil
}

// deletePolicy deletes a policy
func (p *cachingPolicyCRUD) deletePolicy(storage logical.Storage, lp lockingPolicy, name string) error {
	err := deletePolicyCommon(p, lp, storage, name)
	if err != nil {
		return err
	}

	delete(p.cache, name)

	return nil
}
