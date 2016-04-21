package transit

import (
	"sync"

	"github.com/hashicorp/vault/logical"
)

// Directly implements CRUD operations without caching, mapped to the backend,
// but implements shared locking to ensure that we can't overwrite data on the
// backend from multiple operators
type simplePolicyCRUD struct {
	sync.RWMutex
	locks         map[string]*sync.RWMutex
	locksMapMutex sync.RWMutex
}

func newSimplePolicyCRUD() *simplePolicyCRUD {
	return &simplePolicyCRUD{
		locks: map[string]*sync.RWMutex{},
	}
}

func (p *simplePolicyCRUD) ensureLockExists(name string) {
	p.locksMapMutex.RLock()

	if p.locks[name] == nil {
		p.locksMapMutex.RUnlock()
		p.locksMapMutex.Lock()
		// Make sure nothing has appeared since we switched the lock type
		if p.locks[name] == nil {
			p.locks[name] = &sync.RWMutex{}
		}
		p.locksMapMutex.Unlock()
		return
	}

	p.locksMapMutex.RUnlock()
}

// See general comments on the interface method
func (p *simplePolicyCRUD) getPolicy(storage logical.Storage, name string) (lockingPolicy, error) {
	// Use a write lock since fetching the policy can cause a need for upgrade persistence
	p.Lock()
	defer p.Unlock()

	return p.refreshPolicy(storage, name)
}

// See general comments on the interface method
func (p *simplePolicyCRUD) refreshPolicy(storage logical.Storage, name string) (lockingPolicy, error) {
	p.ensureLockExists(name)

	policy, err := fetchPolicyFromStorage(storage, name)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, nil
	}

	lp := &mutexLockingPolicy{
		policy: policy,
		mutex:  p.locks[name],
	}

	return lp, nil
}

// See general comments on the interface method
func (p *simplePolicyCRUD) generatePolicy(storage logical.Storage, name string, derived bool) (lockingPolicy, error) {
	p.ensureLockExists(name)

	policy, err := generatePolicyCommon(p, storage, name, derived)
	if err != nil {
		return nil, err
	}

	lp := &mutexLockingPolicy{
		policy: policy,
		mutex:  p.locks[name],
	}

	return lp, nil
}

// See general comments on the interface method
func (p *simplePolicyCRUD) deletePolicy(storage logical.Storage, lp lockingPolicy, name string) error {
	return deletePolicyCommon(p, lp, storage, name)
}
