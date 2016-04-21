package transit

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/logical"
)

type lockingPolicy interface {
	Lock()
	RLock()
	Unlock()
	RUnlock()
	Policy() *Policy
	SetPolicy(*Policy)
}

type policyCRUD interface {
	// getPolicy returns a lockingPolicy. It performs its own locking according
	// to implementation.
	getPolicy(storage logical.Storage, name string) (lockingPolicy, error)

	// refreshPolicy returns a lockingPolicy. It does not perform its own
	// locking; a write lock must be held before calling.
	refreshPolicy(storage logical.Storage, name string) (lockingPolicy, error)

	// generatePolicy generates and returns a lockingPolicy. A write lock must
	// be held before calling.
	generatePolicy(storage logical.Storage, name string, derived bool) (lockingPolicy, error)

	// deletePolicy deletes a lockingPolicy. A write lock must be held on both
	// the CRUD implementation and the lockingPolicy before calling.
	deletePolicy(storage logical.Storage, lp lockingPolicy, name string) error

	// These are generally satisfied by embedded mutexes in the implementing struct
	Lock()
	RLock()
	Unlock()
	RUnlock()
}

// The mutex is kept separate from the struct since we may set it to its own
// mutex (if the object is shared) or a shared mutext (if the object isn't
// shared and only the locking is)
type mutexLockingPolicy struct {
	mutex  *sync.RWMutex
	policy *Policy
}

func (m *mutexLockingPolicy) Lock() {
	m.mutex.Lock()
}

func (m *mutexLockingPolicy) RLock() {
	m.mutex.RLock()
}

func (m *mutexLockingPolicy) Unlock() {
	m.mutex.Unlock()
}

func (m *mutexLockingPolicy) RUnlock() {
	m.mutex.RUnlock()
}

func (m *mutexLockingPolicy) Policy() *Policy {
	return m.policy
}

func (m *mutexLockingPolicy) SetPolicy(p *Policy) {
	m.policy = p
}

// fetchPolicyFromStorage fetches the policy from backend storage. The caller
// should hold the write lock when calling this, to handle upgrades.
func fetchPolicyFromStorage(storage logical.Storage, name string) (*Policy, error) {
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

	persistNeeded := false
	// Ensure we've moved from Key -> Keys
	if policy.Key != nil && len(policy.Key) > 0 {
		policy.migrateKeyToKeysMap()
		persistNeeded = true
	}

	// With archiving, past assumptions about the length of the keys map are no longer valid
	if policy.LatestVersion == 0 && len(policy.Keys) != 0 {
		policy.LatestVersion = len(policy.Keys)
		persistNeeded = true
	}

	// We disallow setting the version to 0, since they start at 1 since moving
	// to rotate-able keys, so update if it's set to 0
	if policy.MinDecryptionVersion == 0 {
		policy.MinDecryptionVersion = 1
		persistNeeded = true
	}

	// On first load after an upgrade, copy keys to the archive
	if policy.ArchiveVersion == 0 {
		persistNeeded = true
	}

	if persistNeeded {
		err = policy.Persist(storage)
		if err != nil {
			return nil, err
		}
	}

	return policy, nil
}

// generatePolicyCommon is used to create a new named policy with a randomly
// generated key. The caller should have a write lock prior to calling this.
func generatePolicyCommon(p policyCRUD, storage logical.Storage, name string, derived bool) (*Policy, error) {
	// Make sure this doesn't exist in case it was created before we got the write lock
	policy, err := fetchPolicyFromStorage(storage, name)
	if err != nil {
		return nil, err
	}
	if policy != nil {
		return policy, nil
	}

	// Create the policy object
	policy = &Policy{
		Name:       name,
		CipherMode: "aes-gcm",
		Derived:    derived,
	}
	if derived {
		policy.KDFMode = kdfMode
	}

	err = policy.rotate(storage)
	if err != nil {
		return nil, err
	}

	return policy, err
}

// deletePolicyCommon deletes a policy. The caller should hold the write lock
// for both the policy and lockingPolicy prior to calling this.
func deletePolicyCommon(p policyCRUD, lp lockingPolicy, storage logical.Storage, name string) error {
	if lp.Policy() == nil {
		// This got deleted before we grabbed the lock
		return fmt.Errorf("policy already deleted")
	}

	// Verify this hasn't changed
	if !lp.Policy().DeletionAllowed {
		return fmt.Errorf("deletion not allowed for policy %s", name)
	}

	err := storage.Delete("policy/" + name)
	if err != nil {
		return fmt.Errorf("error deleting policy %s: %s", name, err)
	}

	err = storage.Delete("archive/" + name)
	if err != nil {
		return fmt.Errorf("error deleting archive %s: %s", name, err)
	}

	lp.SetPolicy(nil)

	return nil
}
