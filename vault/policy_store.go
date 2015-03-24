package vault

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
)

const (
	// policySubPath is the sub-path used for the policy store
	// view. This is nested under the system view.
	policySubPath = "policy/"
)

// PolicyStore is used to provide durable storage of policy, and to
// manage ACLs associated with them.
type PolicyStore struct {
	view *BarrierView
}

// NewPolicyStore creates a new PolicyStore that is backed
// using a given view. It used used to durable store and manage named policy.
func NewPolicyStore(view *BarrierView) *PolicyStore {
	p := &PolicyStore{
		view: view,
	}
	return p
}

// setupPolicyStore is used to initialize the policy store
// when the vault is being unsealed.
func (c *Core) setupPolicyStore() error {
	// Create a sub-view
	view := c.systemView.SubView(policySubPath)

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
	if p.Name == "root" {
		return fmt.Errorf("cannot update root policy")
	}
	if p.Name == "" {
		return fmt.Errorf("policy name missing")
	}

	entry := &logical.StorageEntry{
		Key:   p.Name,
		Value: []byte(p.Raw),
	}
	if err := ps.view.Put(entry); err != nil {
		return fmt.Errorf("failed to persist policy: %v", err)
	}
	return nil
}

// GetPolicy is used to fetch the named policy
func (ps *PolicyStore) GetPolicy(name string) (*Policy, error) {
	// TODO: Cache policy

	// Special case the root policy
	if name == "root" {
		p := &Policy{Name: "root"}
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

	// Parse into a policy object
	p, err := Parse(string(out.Value))
	if err != nil {
		return nil, fmt.Errorf("failed to parse policy: %v", err)
	}
	return p, nil
}

// ListPolicies is used to list the available policies
func (ps *PolicyStore) ListPolicies() ([]string, error) {
	// Scan the view, since the policy names are the same as the
	// key names.
	return CollectKeys(ps.view)
}

// DeletePolicy is used to delete the named policy
func (ps *PolicyStore) DeletePolicy(name string) error {
	if name == "root" {
		return fmt.Errorf("cannot delete root policy")
	}
	if err := ps.view.Delete(name); err != nil {
		return fmt.Errorf("failed to delete policy: %v", err)
	}
	return nil
}

// ACL is used to return an ACL which is built using the
// named policies.
func (ps *PolicyStore) ACL(names ...string) (*ACL, error) {
	// TODO: Cache ACLs
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
