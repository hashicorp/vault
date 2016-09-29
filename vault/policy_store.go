package vault

import (
	"fmt"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/golang-lru"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
)

const (
	// policySubPath is the sub-path used for the policy store
	// view. This is nested under the system view.
	policySubPath = "policy/"

	// policyCacheSize is the number of policies that are kept cached
	policyCacheSize = 1024

	// responseWrappingPolicyName is the name of the fixed policy
	responseWrappingPolicyName = "response-wrapping"

	// responseWrappingPolicy is the policy that ensures cubbyhole response
	// wrapping can always succeed. Note that sys/wrapping/lookup isn't
	// contained here because using it would revoke the token anyways, so there
	// isn't much point.
	responseWrappingPolicy = `
path "cubbyhole/response" {
    capabilities = ["create", "read"]
}

path "sys/wrapping/unwrap" {
    capabilities = ["update"]
}
`

	// defaultPolicy is the "default" policy
	defaultPolicy = `
# Allow tokens to look up their own properties
path "auth/token/lookup-self" {
    capabilities = ["read"]
}

# Allow tokens to renew themselves
path "auth/token/renew-self" {
    capabilities = ["update"]
}

# Allow tokens to revoke themselves
path "auth/token/revoke-self" {
    capabilities = ["update"]
}

# Allow a token to look up its own capabilities on a path
path "sys/capabilities-self" {
    capabilities = ["update"]
}

# Allow a token to renew a lease via lease_id in the request body
path "sys/renew" {
    capabilities = ["update"]
}

# Allow a token to manage its own cubbyhole
path "cubbyhole/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
}

# Allow a token to list its cubbyhole (not covered by the splat above)
path "cubbyhole" {
    capabilities = ["list"]
}

# Allow a token to wrap arbitrary values in a response-wrapping token
path "sys/wrapping/wrap" {
    capabilities = ["update"]
}

# Allow a token to look up the creation time and TTL of a given
# response-wrapping token
path "sys/wrapping/lookup" {
    capabilities = ["update"]
}

# Allow a token to unwrap a response-wrapping token. This is a convenience to
# avoid client token swapping since this is also part of the response wrapping
# policy.
path "sys/wrapping/unwrap" {
    capabilities = ["update"]
}
`
)

var (
	immutablePolicies = []string{
		"root",
		responseWrappingPolicyName,
	}
	nonAssignablePolicies = []string{
		responseWrappingPolicyName,
	}
)

// PolicyStore is used to provide durable storage of policy, and to
// manage ACLs associated with them.
type PolicyStore struct {
	view *BarrierView
	lru  *lru.TwoQueueCache
}

// PolicyEntry is used to store a policy by name
type PolicyEntry struct {
	Version int
	Raw     string
}

// NewPolicyStore creates a new PolicyStore that is backed
// using a given view. It used used to durable store and manage named policy.
func NewPolicyStore(view *BarrierView, system logical.SystemView) *PolicyStore {
	p := &PolicyStore{
		view: view,
	}
	if !system.CachingDisabled() {
		cache, _ := lru.New2Q(policyCacheSize)
		p.lru = cache
	}

	return p
}

// setupPolicyStore is used to initialize the policy store
// when the vault is being unsealed.
func (c *Core) setupPolicyStore() error {
	// Create a sub-view
	view := c.systemBarrierView.SubView(policySubPath)

	// Create the policy store
	c.policyStore = NewPolicyStore(view, &dynamicSystemView{core: c})

	// Ensure that the default policy exists, and if not, create it
	policy, err := c.policyStore.GetPolicy("default")
	if err != nil {
		return errwrap.Wrapf("error fetching default policy from store: {{err}}", err)
	}
	if policy == nil {
		err := c.policyStore.createDefaultPolicy()
		if err != nil {
			return err
		}
	}

	// Ensure that the cubbyhole response wrapping policy exists
	policy, err = c.policyStore.GetPolicy(responseWrappingPolicyName)
	if err != nil {
		return errwrap.Wrapf("error fetching response-wrapping policy from store: {{err}}", err)
	}
	if policy == nil || policy.Raw != responseWrappingPolicy {
		err := c.policyStore.createResponseWrappingPolicy()
		if err != nil {
			return err
		}
	}

	return nil
}

// teardownPolicyStore is used to reverse setupPolicyStore
// when the vault is being sealed.
func (c *Core) teardownPolicyStore() error {
	c.policyStore = nil
	return nil
}

// SetPolicy is used to create or update the given policy
func (ps *PolicyStore) SetPolicy(p *Policy) error {
	defer metrics.MeasureSince([]string{"policy", "set_policy"}, time.Now())
	if p.Name == "" {
		return fmt.Errorf("policy name missing")
	}
	if strutil.StrListContains(immutablePolicies, p.Name) {
		return fmt.Errorf("cannot update %s policy", p.Name)
	}

	return ps.setPolicyInternal(p)
}

func (ps *PolicyStore) setPolicyInternal(p *Policy) error {
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

	if ps.lru != nil {
		// Update the LRU cache
		ps.lru.Add(p.Name, p)
	}
	return nil
}

// GetPolicy is used to fetch the named policy
func (ps *PolicyStore) GetPolicy(name string) (*Policy, error) {
	defer metrics.MeasureSince([]string{"policy", "get_policy"}, time.Now())
	if ps.lru != nil {
		// Check for cached policy
		if raw, ok := ps.lru.Get(name); ok {
			return raw.(*Policy), nil
		}
	}

	// Special case the root policy
	if name == "root" {
		p := &Policy{Name: "root"}
		if ps.lru != nil {
			ps.lru.Add(p.Name, p)
		}
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

	if ps.lru != nil {
		// Update the LRU cache
		ps.lru.Add(name, policy)
	}

	return policy, nil
}

// ListPolicies is used to list the available policies
func (ps *PolicyStore) ListPolicies() ([]string, error) {
	defer metrics.MeasureSince([]string{"policy", "list_policies"}, time.Now())
	// Scan the view, since the policy names are the same as the
	// key names.
	keys, err := CollectKeys(ps.view)

	for _, nonAssignable := range nonAssignablePolicies {
		deleteIndex := -1
		//Find indices of non-assignable policies in keys
		for index, key := range keys {
			if key == nonAssignable {
				// Delete collection outside the loop
				deleteIndex = index
				break
			}
		}
		// Remove non-assignable policies when found
		if deleteIndex != -1 {
			keys = append(keys[:deleteIndex], keys[deleteIndex+1:]...)
		}
	}

	return keys, err
}

// DeletePolicy is used to delete the named policy
func (ps *PolicyStore) DeletePolicy(name string) error {
	defer metrics.MeasureSince([]string{"policy", "delete_policy"}, time.Now())
	if strutil.StrListContains(immutablePolicies, name) {
		return fmt.Errorf("cannot delete %s policy", name)
	}
	if name == "default" {
		return fmt.Errorf("cannot delete default policy")
	}
	if err := ps.view.Delete(name); err != nil {
		return fmt.Errorf("failed to delete policy: %v", err)
	}

	if ps.lru != nil {
		// Clear the cache
		ps.lru.Remove(name)
	}
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

func (ps *PolicyStore) createDefaultPolicy() error {
	policy, err := Parse(defaultPolicy)
	if err != nil {
		return errwrap.Wrapf("error parsing default policy: {{err}}", err)
	}

	if policy == nil {
		return fmt.Errorf("parsing default policy resulted in nil policy")
	}

	policy.Name = "default"
	return ps.setPolicyInternal(policy)
}

func (ps *PolicyStore) createResponseWrappingPolicy() error {
	policy, err := Parse(responseWrappingPolicy)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error parsing %s policy: {{err}}", responseWrappingPolicyName), err)
	}

	if policy == nil {
		return fmt.Errorf("parsing %s policy resulted in nil policy", responseWrappingPolicyName)
	}

	policy.Name = responseWrappingPolicyName
	return ps.setPolicyInternal(policy)
}
