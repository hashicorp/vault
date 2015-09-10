package vault

import (
	"fmt"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/golang-lru"
	"github.com/hashicorp/vault/logical"
)

const (
	// policySubPath is the sub-path used for the policy store
	// view. This is nested under the system view.
	policySubPath = "policy/"

	// policyCacheSize is the number of policies that are kept cached
	policyCacheSize = 1024
)

// PolicyStore is used to provide durable storage of policy, and to
// manage ACLs associated with them.
type PolicyStore struct {
	view *BarrierView
	lru  *lru.Cache
}

// PolicyEntry is used to store a policy by name
type PolicyEntry struct {
	Version int
	Raw     string
}

// NewPolicyStore creates a new PolicyStore that is backed
// using a given view. It used used to durable store and manage named policy.
func NewPolicyStore(view *BarrierView) *PolicyStore {
	cache, _ := lru.New(policyCacheSize)
	p := &PolicyStore{
		view: view,
		lru:  cache,
	}
	return p
}

// setupPolicyStore is used to initialize the policy store
// when the vault is being unsealed.
func (c *Core) setupPolicyStore() error {
	// Create a sub-view
	view := c.systemBarrierView.SubView(policySubPath)

	// Create the policy store
	c.policy = NewPolicyStore(view)
	return nil
}

// teardownPolicyStore is used to reverse setupPolicyStore
// when the vault is being sealed.
func (c *Core) teardownPolicyStore() error {
	c.policy = nil
	return nil
}

// SetPolicy is used to create or update the given policy
func (ps *PolicyStore) SetPolicy(p *Policy) error {
	defer metrics.MeasureSince([]string{"policy", "set_policy"}, time.Now())
	if p.Name == "root" {
		return fmt.Errorf("cannot update root policy")
	}
	if p.Name == "" {
		return fmt.Errorf("policy name missing")
	}

	// Create the entry
	entry, err := logical.StorageEntryJSON(p.Name, &PolicyEntry{
		Version: 2,
		Raw:     p.Raw,
	})
	if err != nil {
		return fmt.Errorf("failed to create entry: %v", err)
	}
	if err := ps.view.Put(entry); err != nil {
		return fmt.Errorf("failed to persist policy: %v", err)
	}

	// Update the LRU cache
	ps.lru.Add(p.Name, p)
	return nil
}

// GetPolicy is used to fetch the named policy
func (ps *PolicyStore) GetPolicy(name string) (*Policy, error) {
	defer metrics.MeasureSince([]string{"policy", "get_policy"}, time.Now())
	// Check for cached policy
	if raw, ok := ps.lru.Get(name); ok {
		return raw.(*Policy), nil
	}

	// Special case the root policy
	if name == "root" {
		p := &Policy{Name: "root"}
		ps.lru.Add(p.Name, p)
		return p, nil
	}

	// Load the policy in
	out, err := ps.view.Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy: %v", err)
	}
	if out == nil {
		return nil, nil
	}

	// In Vault 0.1.X we stored the raw policy, but in
	// Vault 0.2 we switch to the PolicyEntry
	policyEntry := new(PolicyEntry)
	var policy *Policy
	if err := out.DecodeJSON(policyEntry); err == nil {
		// Parse normally
		p, err := Parse(policyEntry.Raw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse policy: %v", err)
		}
		p.Name = name
		policy = p

	} else {
		// On error, attempt to use V1 parsing
		p, err := Parse(string(out.Value))
		if err != nil {
			return nil, fmt.Errorf("failed to parse policy: %v", err)
		}
		p.Name = name

		// V1 used implicit glob, we need to do a fix-up
		for _, pp := range p.Paths {
			pp.Glob = true
		}
		policy = p
	}

	// Update the LRU cache
	ps.lru.Add(name, policy)
	return policy, nil
}

// ListPolicies is used to list the available policies
func (ps *PolicyStore) ListPolicies() ([]string, error) {
	defer metrics.MeasureSince([]string{"policy", "list_policies"}, time.Now())
	// Scan the view, since the policy names are the same as the
	// key names.
	return CollectKeys(ps.view)
}

// DeletePolicy is used to delete the named policy
func (ps *PolicyStore) DeletePolicy(name string) error {
	defer metrics.MeasureSince([]string{"policy", "delete_policy"}, time.Now())
	if name == "root" {
		return fmt.Errorf("cannot delete root policy")
	}
	if err := ps.view.Delete(name); err != nil {
		return fmt.Errorf("failed to delete policy: %v", err)
	}

	// Clear the cache
	ps.lru.Remove(name)
	return nil
}

// ACL is used to return an ACL which is built using the
// named policies.
func (ps *PolicyStore) ACL(names ...string) (*ACL, error) {
	// Fetch the policies
	var policy []*Policy
	for _, name := range names {
		p, err := ps.GetPolicy(name)
		if err != nil {
			return nil, fmt.Errorf("failed to get policy '%s': %v", name, err)
		}
		policy = append(policy, p)
	}

	// Construct the ACL
	acl, err := NewACL(policy)
	if err != nil {
		return nil, fmt.Errorf("failed to construct ACL: %v", err)
	}
	return acl, nil
}
