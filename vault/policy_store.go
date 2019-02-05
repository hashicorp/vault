package vault

import (
	"context"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
)

const (
	// policySubPath is the sub-path used for the policy store view. This is
	// nested under the system view. policyRGPSubPath/policyEGPSubPath are
	// similar but for RGPs/EGPs.
	policyACLSubPath = "policy/"
	policyRGPSubPath = "policy-rgp/"
	policyEGPSubPath = "policy-egp/"

	// policyCacheSize is the number of policies that are kept cached
	policyCacheSize = 1024

	// defaultPolicyName is the name of the default policy
	defaultPolicyName = "default"

	// responseWrappingPolicyName is the name of the fixed policy
	responseWrappingPolicyName = "response-wrapping"

	// controlGroupPolicyName is the name of the fixed policy for control group
	// tokens
	controlGroupPolicyName = "control-group"

	// responseWrappingPolicy is the policy that ensures cubbyhole response
	// wrapping can always succeed.
	responseWrappingPolicy = `
path "cubbyhole/response" {
    capabilities = ["create", "read"]
}

path "sys/wrapping/unwrap" {
    capabilities = ["update"]
}
`
	// controlGroupPolicy is the policy that ensures control group requests can
	// commit themselves
	controlGroupPolicy = `
path "cubbyhole/control-group" {
    capabilities = ["update", "create", "read"]
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

# Allow a token to look up its own entity by id or name
path "identity/entity/id/{{identity.entity.id}}" {
  capabilities = ["read"]
}
path "identity/entity/name/{{identity.entity.name}}" {
  capabilities = ["read"]
}


# Allow a token to look up its resultant ACL from all policies. This is useful
# for UIs. It is an internal path because the format may change at any time
# based on how the internal ACL features and capabilities change.
path "sys/internal/ui/resultant-acl" {
    capabilities = ["read"]
}

# Allow a token to renew a lease via lease_id in the request body; old path for
# old clients, new path for newer
path "sys/renew" {
    capabilities = ["update"]
}
path "sys/leases/renew" {
    capabilities = ["update"]
}

# Allow looking up lease properties. This requires knowing the lease ID ahead
# of time and does not divulge any sensitive information.
path "sys/leases/lookup" {
    capabilities = ["update"]
}

# Allow a token to manage its own cubbyhole
path "cubbyhole/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
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

# Allow general purpose tools
path "sys/tools/hash" {
    capabilities = ["update"]
}
path "sys/tools/hash/*" {
    capabilities = ["update"]
}
path "sys/tools/random" {
    capabilities = ["update"]
}
path "sys/tools/random/*" {
    capabilities = ["update"]
}

# Allow checking the status of a Control Group request if the user has the
# accessor
path "sys/control-group/request" {
    capabilities = ["update"]
}
`
)

var (
	immutablePolicies = []string{
		"root",
		responseWrappingPolicyName,
		controlGroupPolicyName,
	}
	nonAssignablePolicies = []string{
		responseWrappingPolicyName,
		controlGroupPolicyName,
	}
)

// PolicyStore is used to provide durable storage of policy, and to
// manage ACLs associated with them.
type PolicyStore struct {
	entPolicyStore

	core    *Core
	aclView *BarrierView
	rgpView *BarrierView
	egpView *BarrierView

	tokenPoliciesLRU *lru.TwoQueueCache
	egpLRU           *lru.TwoQueueCache

	// This is used to ensure that writes to the store (acl/rgp) or to the egp
	// path tree don't happen concurrently. We are okay reading stale data so
	// long as there aren't concurrent writes.
	modifyLock *sync.RWMutex

	// Stores whether a token policy is ACL or RGP
	policyTypeMap sync.Map

	// logger is the server logger copied over from core
	logger log.Logger
}

// PolicyEntry is used to store a policy by name
type PolicyEntry struct {
	sentinelPolicy

	Version   int
	Raw       string
	Templated bool
	Type      PolicyType
}

// NewPolicyStore creates a new PolicyStore that is backed
// using a given view. It used used to durable store and manage named policy.
func NewPolicyStore(ctx context.Context, core *Core, baseView *BarrierView, system logical.SystemView, logger log.Logger) (*PolicyStore, error) {
	ps := &PolicyStore{
		aclView:    baseView.SubView(policyACLSubPath),
		rgpView:    baseView.SubView(policyRGPSubPath),
		egpView:    baseView.SubView(policyEGPSubPath),
		modifyLock: new(sync.RWMutex),
		logger:     logger,
		core:       core,
	}

	ps.extraInit()

	if !system.CachingDisabled() {
		cache, _ := lru.New2Q(policyCacheSize)
		ps.tokenPoliciesLRU = cache
		cache, _ = lru.New2Q(policyCacheSize)
		ps.egpLRU = cache
	}

	aclView := ps.getACLView(namespace.RootNamespace)
	keys, err := logical.CollectKeys(namespace.RootContext(ctx), aclView)
	if err != nil {
		ps.logger.Error("error collecting acl policy keys", "error", err)
		return nil, err
	}
	for _, key := range keys {
		index := ps.cacheKey(namespace.RootNamespace, ps.sanitizeName(key))
		ps.policyTypeMap.Store(index, PolicyTypeACL)
	}

	if err := ps.loadNamespacePolicies(ctx, core); err != nil {
		return nil, err
	}

	// Special-case root; doesn't exist on disk but does need to be found
	ps.policyTypeMap.Store(ps.cacheKey(namespace.RootNamespace, "root"), PolicyTypeACL)
	return ps, nil
}

// setupPolicyStore is used to initialize the policy store
// when the vault is being unsealed.
func (c *Core) setupPolicyStore(ctx context.Context) error {
	// Create the policy store
	var err error
	sysView := &dynamicSystemView{core: c}
	psLogger := c.baseLogger.Named("policy")
	c.AddLogger(psLogger)
	c.policyStore, err = NewPolicyStore(ctx, c, c.systemBarrierView, sysView, psLogger)
	if err != nil {
		return err
	}

	if c.ReplicationState().HasState(consts.ReplicationPerformanceSecondary) {
		// Policies will sync from the primary
		return nil
	}

	// Ensure that the default policy exists, and if not, create it
	if err := c.policyStore.loadACLPolicy(ctx, defaultPolicyName, defaultPolicy); err != nil {
		return err
	}
	// Ensure that the response wrapping policy exists
	if err := c.policyStore.loadACLPolicy(ctx, responseWrappingPolicyName, responseWrappingPolicy); err != nil {
		return err
	}
	// Ensure that the control group policy exists
	if err := c.policyStore.loadACLPolicy(ctx, controlGroupPolicyName, controlGroupPolicy); err != nil {
		return err
	}

	return nil
}

// teardownPolicyStore is used to reverse setupPolicyStore
// when the vault is being sealed.
func (c *Core) teardownPolicyStore() error {
	c.policyStore = nil
	return nil
}

func (ps *PolicyStore) invalidate(ctx context.Context, name string, policyType PolicyType) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		ps.logger.Error("unable to invalidate key, no namespace info passed", "key", name)
		return
	}

	// This may come with a prefixed "/" due to joining the file path
	saneName := strings.TrimPrefix(name, "/")
	index := ps.cacheKey(ns, saneName)

	ps.modifyLock.Lock()
	defer ps.modifyLock.Unlock()

	// We don't lock before removing from the LRU here because the worst that
	// can happen is we load again if something since added it
	switch policyType {
	case PolicyTypeACL, PolicyTypeRGP:
		if ps.tokenPoliciesLRU != nil {
			ps.tokenPoliciesLRU.Remove(index)
		}

	case PolicyTypeEGP:
		if ps.egpLRU != nil {
			ps.egpLRU.Remove(index)
		}

	default:
		// Can't do anything
		return
	}

	// Force a reload
	out, err := ps.switchedGetPolicy(ctx, name, policyType, false)
	if err != nil {
		ps.logger.Error("error fetching policy after invalidation", "name", saneName)
	}

	// If true, the invalidation was actually a delete, so we may need to
	// perform further deletion tasks. We skip the physical deletion just in
	// case another process has re-written the policy; instead next time Get is
	// called the values will be loaded back in.
	if out == nil {
		ps.switchedDeletePolicy(ctx, name, policyType, false)
	}

	return
}

// SetPolicy is used to create or update the given policy
func (ps *PolicyStore) SetPolicy(ctx context.Context, p *Policy) error {
	defer metrics.MeasureSince([]string{"policy", "set_policy"}, time.Now())
	if p == nil {
		return fmt.Errorf("nil policy passed in for storage")
	}
	if p.Name == "" {
		return fmt.Errorf("policy name missing")
	}
	// Policies are normalized to lower-case
	p.Name = ps.sanitizeName(p.Name)
	if strutil.StrListContains(immutablePolicies, p.Name) {
		return fmt.Errorf("cannot update %q policy", p.Name)
	}

	return ps.setPolicyInternal(ctx, p)
}

func (ps *PolicyStore) setPolicyInternal(ctx context.Context, p *Policy) error {
	ps.modifyLock.Lock()
	defer ps.modifyLock.Unlock()

	// Get the appropriate view based on policy type and namespace
	view := ps.getBarrierView(p.namespace, p.Type)
	if view == nil {
		return fmt.Errorf("unable to get the barrier subview for policy type %q", p.Type)
	}

	if err := ps.parseEGPPaths(p); err != nil {
		return err
	}

	// Create the entry
	entry, err := logical.StorageEntryJSON(p.Name, &PolicyEntry{
		Version:        2,
		Raw:            p.Raw,
		Type:           p.Type,
		Templated:      p.Templated,
		sentinelPolicy: p.sentinelPolicy,
	})
	if err != nil {
		return errwrap.Wrapf("failed to create entry: {{err}}", err)
	}

	// Construct the cache key
	index := ps.cacheKey(p.namespace, p.Name)

	switch p.Type {
	case PolicyTypeACL:
		rgpView := ps.getRGPView(p.namespace)
		rgp, err := rgpView.Get(ctx, entry.Key)
		if err != nil {
			return errwrap.Wrapf("failed looking up conflicting policy: {{err}}", err)
		}
		if rgp != nil {
			return fmt.Errorf("cannot reuse policy names between ACLs and RGPs")
		}

		if err := view.Put(ctx, entry); err != nil {
			return errwrap.Wrapf("failed to persist policy: {{err}}", err)
		}

		ps.policyTypeMap.Store(index, PolicyTypeACL)

		if ps.tokenPoliciesLRU != nil {
			ps.tokenPoliciesLRU.Add(index, p)
		}

	case PolicyTypeRGP:
		aclView := ps.getACLView(p.namespace)
		acl, err := aclView.Get(ctx, entry.Key)
		if err != nil {
			return errwrap.Wrapf("failed looking up conflicting policy: {{err}}", err)
		}
		if acl != nil {
			return fmt.Errorf("cannot reuse policy names between ACLs and RGPs")
		}

		if err := ps.handleSentinelPolicy(ctx, p, view, entry); err != nil {
			return err
		}

		ps.policyTypeMap.Store(index, PolicyTypeRGP)

		// We load here after successfully loading into Sentinel so that on
		// error we will try loading again on the next get
		if ps.tokenPoliciesLRU != nil {
			ps.tokenPoliciesLRU.Add(index, p)
		}

	case PolicyTypeEGP:
		if err := ps.handleSentinelPolicy(ctx, p, view, entry); err != nil {
			return err
		}

		// We load here after successfully loading into Sentinel so that on
		// error we will try loading again on the next get
		if ps.egpLRU != nil {
			ps.egpLRU.Add(index, p)
		}

	default:
		return fmt.Errorf("unknown policy type, cannot set")
	}

	return nil
}

// GetPolicy is used to fetch the named policy
func (ps *PolicyStore) GetPolicy(ctx context.Context, name string, policyType PolicyType) (*Policy, error) {
	return ps.switchedGetPolicy(ctx, name, policyType, true)
}

func (ps *PolicyStore) switchedGetPolicy(ctx context.Context, name string, policyType PolicyType, grabLock bool) (*Policy, error) {
	defer metrics.MeasureSince([]string{"policy", "get_policy"}, time.Now())

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	// Policies are normalized to lower-case
	name = ps.sanitizeName(name)
	index := ps.cacheKey(ns, name)

	var cache *lru.TwoQueueCache
	var view *BarrierView

	switch policyType {
	case PolicyTypeACL:
		cache = ps.tokenPoliciesLRU
		view = ps.getACLView(ns)
	case PolicyTypeRGP:
		cache = ps.tokenPoliciesLRU
		view = ps.getRGPView(ns)
	case PolicyTypeEGP:
		cache = ps.egpLRU
		view = ps.getEGPView(ns)
	case PolicyTypeToken:
		cache = ps.tokenPoliciesLRU
		val, ok := ps.policyTypeMap.Load(index)
		if !ok {
			// Doesn't exist
			return nil, nil
		}
		policyType = val.(PolicyType)
		switch policyType {
		case PolicyTypeACL:
			view = ps.getACLView(ns)
		case PolicyTypeRGP:
			view = ps.getRGPView(ns)
		default:
			return nil, fmt.Errorf("invalid type of policy in type map: %q", policyType)
		}
	}

	if cache != nil {
		// Check for cached policy
		if raw, ok := cache.Get(index); ok {
			return raw.(*Policy), nil
		}
	}

	// Special case the root policy
	if policyType == PolicyTypeACL && name == "root" && ns.ID == namespace.RootNamespaceID {
		p := &Policy{
			Name:      "root",
			namespace: namespace.RootNamespace,
		}
		if cache != nil {
			cache.Add(index, p)
		}
		return p, nil
	}

	if grabLock {
		ps.modifyLock.Lock()
		defer ps.modifyLock.Unlock()
	}

	// See if anything has added it since we got the lock
	if cache != nil {
		if raw, ok := cache.Get(index); ok {
			return raw.(*Policy), nil
		}
	}

	// Nil-check on the view before proceeding to retrive from storage
	if view == nil {
		return nil, fmt.Errorf("unable to get the barrier subview for policy type %q", policyType)
	}

	out, err := view.Get(ctx, name)
	if err != nil {
		return nil, errwrap.Wrapf("failed to read policy: {{err}}", err)
	}

	if out == nil {
		return nil, nil
	}

	policyEntry := new(PolicyEntry)
	policy := new(Policy)
	err = out.DecodeJSON(policyEntry)
	if err != nil {
		return nil, errwrap.Wrapf("failed to parse policy: {{err}}", err)
	}

	// Set these up here so that they're available for loading into
	// Sentinel
	policy.Name = name
	policy.Raw = policyEntry.Raw
	policy.Type = policyEntry.Type
	policy.Templated = policyEntry.Templated
	policy.sentinelPolicy = policyEntry.sentinelPolicy
	policy.namespace = ns
	switch policyEntry.Type {
	case PolicyTypeACL:
		// Parse normally
		p, err := ParseACLPolicy(ns, policyEntry.Raw)
		if err != nil {
			return nil, errwrap.Wrapf("failed to parse policy: {{err}}", err)
		}
		policy.Paths = p.Paths

		// Reset this in case they set the name in the policy itself
		policy.Name = name

		ps.policyTypeMap.Store(index, PolicyTypeACL)

	case PolicyTypeRGP:
		if err := ps.handleSentinelPolicy(ctx, policy, nil, nil); err != nil {
			return nil, err
		}

		ps.policyTypeMap.Store(index, PolicyTypeRGP)

	case PolicyTypeEGP:
		if err := ps.handleSentinelPolicy(ctx, policy, nil, nil); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown policy type %q", policyEntry.Type.String())
	}

	if cache != nil {
		// Update the LRU cache
		cache.Add(index, policy)
	}

	return policy, nil
}

// ListPolicies is used to list the available policies
func (ps *PolicyStore) ListPolicies(ctx context.Context, policyType PolicyType) ([]string, error) {
	defer metrics.MeasureSince([]string{"policy", "list_policies"}, time.Now())

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ns == nil {
		return nil, namespace.ErrNoNamespace
	}

	// Get the appropriate view based on policy type and namespace
	view := ps.getBarrierView(ns, policyType)
	if view == nil {
		return []string{}, fmt.Errorf("unable to get the barrier subview for policy type %q", policyType)
	}

	// Scan the view, since the policy names are the same as the
	// key names.
	var keys []string
	switch policyType {
	case PolicyTypeACL:
		keys, err = logical.CollectKeys(ctx, view)
	case PolicyTypeRGP:
		return logical.CollectKeys(ctx, view)
	case PolicyTypeEGP:
		return logical.CollectKeys(ctx, view)
	default:
		return nil, fmt.Errorf("unknown policy type %q", policyType)
	}

	// We only have non-assignable ACL policies at the moment
	for _, nonAssignable := range nonAssignablePolicies {
		deleteIndex := -1
		// Find indices of non-assignable policies in keys
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
func (ps *PolicyStore) DeletePolicy(ctx context.Context, name string, policyType PolicyType) error {
	return ps.switchedDeletePolicy(ctx, name, policyType, true)
}

func (ps *PolicyStore) switchedDeletePolicy(ctx context.Context, name string, policyType PolicyType, physicalDeletion bool) error {
	defer metrics.MeasureSince([]string{"policy", "delete_policy"}, time.Now())

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	// If not set, the call comes from invalidation, where we'll already have
	// grabbed the lock
	if physicalDeletion {
		ps.modifyLock.Lock()
		defer ps.modifyLock.Unlock()
	}

	// Policies are normalized to lower-case
	name = ps.sanitizeName(name)
	index := ps.cacheKey(ns, name)

	view := ps.getBarrierView(ns, policyType)
	if view == nil {
		return fmt.Errorf("unable to get the barrier subview for policy type %q", policyType)
	}

	switch policyType {
	case PolicyTypeACL:
		if strutil.StrListContains(immutablePolicies, name) {
			return fmt.Errorf("cannot delete %q policy", name)
		}
		if name == "default" {
			return fmt.Errorf("cannot delete default policy")
		}

		if physicalDeletion {
			err := view.Delete(ctx, name)
			if err != nil {
				return errwrap.Wrapf("failed to delete policy: {{err}}", err)
			}
		}

		if ps.tokenPoliciesLRU != nil {
			// Clear the cache
			ps.tokenPoliciesLRU.Remove(index)
		}

		ps.policyTypeMap.Delete(index)

	case PolicyTypeRGP:
		if physicalDeletion {
			err := view.Delete(ctx, name)
			if err != nil {
				return errwrap.Wrapf("failed to delete policy: {{err}}", err)
			}
		}

		if ps.tokenPoliciesLRU != nil {
			// Clear the cache
			ps.tokenPoliciesLRU.Remove(index)
		}

		ps.policyTypeMap.Delete(index)

		defer ps.core.invalidateSentinelPolicy(policyType, index)

	case PolicyTypeEGP:
		if physicalDeletion {
			err := view.Delete(ctx, name)
			if err != nil {
				return errwrap.Wrapf("failed to delete policy: {{err}}", err)
			}
		}

		if ps.egpLRU != nil {
			// Clear the cache
			ps.egpLRU.Remove(index)
		}

		defer ps.core.invalidateSentinelPolicy(policyType, index)

		ps.invalidateEGPTreePath(index)
	}

	return nil
}

type TemplateError struct {
	Err error
}

func (t *TemplateError) WrappedErrors() []error {
	return []error{t.Err}
}

func (t *TemplateError) Error() string {
	return t.Err.Error()
}

// ACL is used to return an ACL which is built using the
// named policies.
func (ps *PolicyStore) ACL(ctx context.Context, entity *identity.Entity, policyNames map[string][]string) (*ACL, error) {
	var policies []*Policy
	// Fetch the policies
	for nsID, nsPolicyNames := range policyNames {
		policyNS, err := NamespaceByID(ctx, nsID, ps.core)
		if err != nil {
			return nil, err
		}
		if policyNS == nil {
			return nil, namespace.ErrNoNamespace
		}
		policyCtx := namespace.ContextWithNamespace(ctx, policyNS)
		for _, nsPolicyName := range nsPolicyNames {
			p, err := ps.GetPolicy(policyCtx, nsPolicyName, PolicyTypeToken)
			if err != nil {
				return nil, errwrap.Wrapf("failed to get policy: {{err}}", err)
			}
			if p != nil {
				policies = append(policies, p)
			}
		}
	}

	var fetchedGroups bool
	var groups []*identity.Group
	for i, policy := range policies {
		if policy.Type == PolicyTypeACL && policy.Templated {
			if !fetchedGroups {
				fetchedGroups = true
				if entity != nil {
					directGroups, inheritedGroups, err := ps.core.identityStore.groupsByEntityID(entity.ID)
					if err != nil {
						return nil, errwrap.Wrapf("failed to fetch group memberships: {{err}}", err)
					}
					groups = append(directGroups, inheritedGroups...)
				}
			}
			p, err := parseACLPolicyWithTemplating(policy.namespace, policy.Raw, true, entity, groups)
			if err != nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("error parsing templated policy %q: {{err}}", policy.Name), err)
			}
			p.Name = policy.Name
			policies[i] = p
		}
	}

	// Construct the ACL
	acl, err := NewACL(ctx, policies)
	if err != nil {
		return nil, errwrap.Wrapf("failed to construct ACL: {{err}}", err)
	}

	return acl, nil
}

// loadACLPolicy is used to load default ACL policies. The default policies will
// be loaded to all namespaces.
func (ps *PolicyStore) loadACLPolicy(ctx context.Context, policyName, policyText string) error {
	return ps.loadACLPolicyNamespaces(ctx, policyName, policyText)
}

// loadACLPolicyInternal is used to load default ACL policies in a specific
// namespace.
func (ps *PolicyStore) loadACLPolicyInternal(ctx context.Context, policyName, policyText string) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	// Check if the policy already exists
	policy, err := ps.GetPolicy(ctx, policyName, PolicyTypeACL)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error fetching %s policy from store: {{err}}", policyName), err)
	}
	if policy != nil {
		if !strutil.StrListContains(immutablePolicies, policyName) || policyText == policy.Raw {
			return nil
		}
	}

	policy, err = ParseACLPolicy(ns, policyText)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error parsing %s policy: {{err}}", policyName), err)
	}

	if policy == nil {
		return fmt.Errorf("parsing %q policy resulted in nil policy", policyName)
	}

	policy.Name = policyName
	policy.Type = PolicyTypeACL
	return ps.setPolicyInternal(ctx, policy)
}

func (ps *PolicyStore) sanitizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func (ps *PolicyStore) cacheKey(ns *namespace.Namespace, name string) string {
	return path.Join(ns.ID, name)
}
