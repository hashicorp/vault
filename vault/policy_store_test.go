// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"bytes"
	"context"
	"reflect"
	"sync"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func mockPolicyWithCore(t *testing.T, disableCache bool) (*Core, *PolicyStore) {
	conf := &CoreConfig{
		DisableCache: disableCache,
	}
	cluster := NewTestCluster(t, conf, nil)
	core := cluster.Cores[0].Core
	ps := core.policyStore

	return core, ps
}

func TestPolicyStore_Root(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Parallel()

		core, _, _ := TestCoreUnsealed(t)
		ps := core.policyStore
		testPolicyRoot(t, ps, namespace.RootNamespace, true)
	})
}

func testPolicyRoot(t *testing.T, ps *PolicyStore, ns *namespace.Namespace, expectFound bool) {
	// Get should return a special policy
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	p, err := ps.GetPolicy(ctx, "root", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Handle whether a root token is expected
	if expectFound {
		if p == nil {
			t.Fatalf("bad: %v", p)
		}
		if p.Name != "root" {
			t.Fatalf("bad: %v", p)
		}
	} else {
		if p != nil {
			t.Fatal("expected nil root policy")
		}
		// Create root policy for subsequent modification and deletion failure
		// tests
		p = &Policy{
			Name: "root",
		}
	}

	// Set should fail
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	err = ps.SetPolicy(ctx, p)
	if err.Error() != `cannot update "root" policy` {
		t.Fatalf("err: %v", err)
	}

	// Delete should fail
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	err = ps.DeletePolicy(ctx, "root", PolicyTypeACL)
	if err.Error() != `cannot delete "root" policy` {
		t.Fatalf("err: %v", err)
	}
}

func TestPolicyStore_CRUD(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		t.Run("cached", func(t *testing.T) {
			_, ps := mockPolicyWithCore(t, false)
			testPolicyStoreCRUD(t, ps, namespace.RootNamespace)
		})

		t.Run("no-cache", func(t *testing.T) {
			_, ps := mockPolicyWithCore(t, true)
			testPolicyStoreCRUD(t, ps, namespace.RootNamespace)
		})
	})
}

func testPolicyStoreCRUD(t *testing.T, ps *PolicyStore, ns *namespace.Namespace) {
	// Get should return nothing
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	p, err := ps.GetPolicy(ctx, "Dev", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if p != nil {
		t.Fatalf("bad: %v", p)
	}

	// Delete should be no-op
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	err = ps.DeletePolicy(ctx, "deV", PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// List should be blank
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	out, err := ps.ListPolicies(ctx, PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected := []string{defaultPolicyName, defaultCeilingPolicyName}
	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("bad: %v", out)
	}

	// Set should work
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	policy, _ := ParseACLPolicy(ns, aclPolicy, WithDenySlashInTemplatedPaths(ps.core.denySlashInTemplatedPolicyPaths))
	err = ps.SetPolicy(ctx, policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should work
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	p, err = ps.GetPolicy(ctx, "dEv", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(p, policy) {
		t.Fatalf("bad: %v", p)
	}

	// List should contain the two built-in assignable policies plus the new policy.
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	out, err = ps.ListPolicies(ctx, PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("bad: %v", out)
	}

	expected = []string{defaultPolicyName, defaultCeilingPolicyName, "dev"}
	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("expected: %v\ngot: %v", expected, out)
	}

	// Delete should be clear the entry
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	err = ps.DeletePolicy(ctx, "Dev", PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// List should contain one element
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	out, err = ps.ListPolicies(ctx, PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected = []string{defaultPolicyName, defaultCeilingPolicyName}
	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("bad: %v", out)
	}

	// Get should fail
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	p, err = ps.GetPolicy(ctx, "deV", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if p != nil {
		t.Fatalf("bad: %v", p)
	}
}

func TestPolicyStore_Predefined(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		_, ps := mockPolicyWithCore(t, false)
		testPolicyStorePredefined(t, ps, namespace.RootNamespace)
	})
}

// Test predefined policy handling
func testPolicyStorePredefined(t *testing.T, ps *PolicyStore, ns *namespace.Namespace) {
	// List should contain the built-in assignable ACL policies.
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	out, err := ps.ListPolicies(ctx, PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// This shouldn't contain response-wrapping since it's non-assignable.
	expected := []string{defaultPolicyName, defaultCeilingPolicyName}
	if !reflect.DeepEqual(expected, out) {
		t.Fatalf("bad: %v", out)
	}

	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	pDefaultCeiling, err := ps.GetPolicy(ctx, defaultCeilingPolicyName, PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if pDefaultCeiling == nil {
		t.Fatal("nil default ceiling policy")
	}
	if pDefaultCeiling.Raw != defaultCeilingPolicy {
		t.Fatalf("bad: expected\n%s\ngot\n%s\n", defaultCeilingPolicy, pDefaultCeiling.Raw)
	}
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	err = ps.DeletePolicy(ctx, pDefaultCeiling.Name, PolicyTypeACL)
	if err == nil {
		t.Fatalf("expected err deleting %s", pDefaultCeiling.Name)
	}

	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	updatedDefaultCeiling, err := ParseACLPolicy(ns, aclPolicy, WithDenySlashInTemplatedPaths(ps.core.denySlashInTemplatedPolicyPaths))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	updatedDefaultCeiling.Name = defaultCeilingPolicyName
	err = ps.SetPolicy(ctx, updatedDefaultCeiling)
	if err != nil {
		t.Fatalf("expected err to be nil updating %s: %v", updatedDefaultCeiling.Name, err)
	}
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	pDefaultCeiling, err = ps.GetPolicy(ctx, defaultCeilingPolicyName, PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if pDefaultCeiling == nil {
		t.Fatal("nil updated default ceiling policy")
	}
	if pDefaultCeiling.Raw != updatedDefaultCeiling.Raw {
		t.Fatalf("bad: expected\n%s\ngot\n%s\n", updatedDefaultCeiling.Raw, pDefaultCeiling.Raw)
	}

	// Response-wrapping policy checks
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	pCubby, err := ps.GetPolicy(ctx, "response-wrapping", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if pCubby == nil {
		t.Fatal("nil cubby policy")
	}
	if pCubby.Raw != responseWrappingPolicy {
		t.Fatalf("bad: expected\n%s\ngot\n%s\n", responseWrappingPolicy, pCubby.Raw)
	}
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	err = ps.SetPolicy(ctx, pCubby)
	if err == nil {
		t.Fatalf("expected err setting %s", pCubby.Name)
	}
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	err = ps.DeletePolicy(ctx, pCubby.Name, PolicyTypeACL)
	if err == nil {
		t.Fatalf("expected err deleting %s", pCubby.Name)
	}

	// Root policy checks, behavior depending on namespace
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	pRoot, err := ps.GetPolicy(ctx, "root", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if ns == namespace.RootNamespace {
		if pRoot == nil {
			t.Fatal("nil root policy")
		}
	} else {
		if pRoot != nil {
			t.Fatal("expected nil root policy")
		}
		pRoot = &Policy{
			Name: "root",
		}
	}
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	err = ps.SetPolicy(ctx, pRoot)
	if err == nil {
		t.Fatalf("expected err setting %s", pRoot.Name)
	}
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	err = ps.DeletePolicy(ctx, pRoot.Name, PolicyTypeACL)
	if err == nil {
		t.Fatalf("expected err deleting %s", pRoot.Name)
	}
}

func TestPolicyStore_ACL(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		_, ps := mockPolicyWithCore(t, false)
		testPolicyStoreACL(t, ps, namespace.RootNamespace)
	})
}

func testPolicyStoreACL(t *testing.T, ps *PolicyStore, ns *namespace.Namespace) {
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	policy, _ := ParseACLPolicy(ns, aclPolicy, WithDenySlashInTemplatedPaths(ps.core.denySlashInTemplatedPolicyPaths))
	err := ps.SetPolicy(ctx, policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	policy, _ = ParseACLPolicy(ns, aclPolicy2, WithDenySlashInTemplatedPaths(ps.core.denySlashInTemplatedPolicyPaths))
	err = ps.SetPolicy(ctx, policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	acl, err := ps.ACL(ctx, nil, map[string][]string{ns.ID: {"dev", "ops"}})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testLayeredACL(t, acl, ns)
}

func TestDefaultPolicy(t *testing.T) {
	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	policy, err := ParseACLPolicy(namespace.RootNamespace, defaultPolicy)
	if err != nil {
		t.Fatal(err)
	}
	acl, err := NewACL(ctx, []*Policy{policy})
	if err != nil {
		t.Fatal(err)
	}

	for name, tc := range map[string]struct {
		op            logical.Operation
		path          string
		expectAllowed bool
	}{
		"lookup self":            {logical.ReadOperation, "auth/token/lookup-self", true},
		"renew self":             {logical.UpdateOperation, "auth/token/renew-self", true},
		"revoke self":            {logical.UpdateOperation, "auth/token/revoke-self", true},
		"check own capabilities": {logical.UpdateOperation, "sys/capabilities-self", true},

		"read arbitrary path":     {logical.ReadOperation, "foo/bar", false},
		"login at arbitrary path": {logical.UpdateOperation, "auth/foo", false},
	} {
		t.Run(name, func(t *testing.T) {
			request := new(logical.Request)
			request.Operation = tc.op
			request.Path = tc.path

			result := acl.AllowOperation(ctx, request, false)
			if result.RootPrivs {
				t.Fatal("unexpected root")
			}
			if tc.expectAllowed != result.Allowed {
				t.Fatalf("Expected %v, got %v", tc.expectAllowed, result.Allowed)
			}
		})
	}
}

// TestPolicyStore_PoliciesByNamespaces tests the policiesByNamespaces function, which should return a slice of policy names for a given slice of namespaces.
func TestPolicyStore_PoliciesByNamespaces(t *testing.T) {
	_, ps := mockPolicyWithCore(t, false)

	ctxRoot := namespace.RootContext(context.Background())
	rootNs := namespace.RootNamespace

	parsedPolicy, _ := ParseACLPolicy(rootNs, aclPolicy, WithDenySlashInTemplatedPaths(ps.core.denySlashInTemplatedPolicyPaths))

	err := ps.SetPolicy(ctxRoot, parsedPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should work
	pResult, err := ps.GetPolicy(ctxRoot, "dev", PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(pResult, parsedPolicy) {
		t.Fatalf("bad: %v", pResult)
	}

	out, err := ps.policiesByNamespaces(ctxRoot, PolicyTypeACL, []*namespace.Namespace{rootNs})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expectedResult := []string{defaultPolicyName, defaultCeilingPolicyName, "dev"}
	if !reflect.DeepEqual(expectedResult, out) {
		t.Fatalf("expected: %v\ngot: %v", expectedResult, out)
	}
}

// TestPolicyStore_GetNonEGPPolicyType has five test cases:
//   - happy-acl and happy-rgp: we store a policy in the policy type map and
//     then look up its type successfully.
//   - not-in-map-acl and not-in-map-rgp: ensure that GetNonEGPPolicyType fails
//     returning a nil and an error when the policy doesn't exist in the map.
//   - unknown-policy-type: ensures that GetNonEGPPolicyType fails returning a nil
//     and an error when the policy type in the type map is a value that
//     does not map to a PolicyType.
func TestPolicyStore_GetNonEGPPolicyType(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		policyStoreKey       string
		policyStoreValue     any
		paramNamespace       string
		paramPolicyName      string
		paramPolicyType      PolicyType
		isErrorExpected      bool
		expectedErrorMessage string
	}{
		"happy-acl": {
			policyStoreKey:   "5:1AbcD:policy1",
			policyStoreValue: PolicyTypeACL,
			paramNamespace:   "1AbcD",
			paramPolicyName:  "policy1",
			paramPolicyType:  PolicyTypeACL,
		},
		"happy-rgp": {
			policyStoreKey:   "5:1AbcD:policy1",
			policyStoreValue: PolicyTypeRGP,
			paramNamespace:   "1AbcD",
			paramPolicyName:  "policy1",
			paramPolicyType:  PolicyTypeRGP,
		},
		"not-in-map-acl": {
			policyStoreKey:       "5:2WxyZ:policy2",
			policyStoreValue:     PolicyTypeACL,
			paramNamespace:       "1AbcD",
			paramPolicyName:      "policy1",
			isErrorExpected:      true,
			expectedErrorMessage: "policy does not exist in type map",
		},
		"not-in-map-rgp": {
			policyStoreKey:       "5:2WxyZ:policy2",
			policyStoreValue:     PolicyTypeRGP,
			paramNamespace:       "1AbcD",
			paramPolicyName:      "policy1",
			isErrorExpected:      true,
			expectedErrorMessage: "policy does not exist in type map",
		},
		"unknown-policy-type": {
			policyStoreKey:       "5:1AbcD:policy1",
			policyStoreValue:     7,
			paramNamespace:       "1AbcD",
			paramPolicyName:      "policy1",
			isErrorExpected:      true,
			expectedErrorMessage: "unknown policy type for: 5:1AbcD:policy1",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, ps := mockPolicyWithCore(t, false)
			ps.policyTypeMap.Store(tc.policyStoreKey, tc.policyStoreValue)
			got, err := ps.GetNonEGPPolicyType(tc.paramNamespace, tc.paramPolicyName)
			if tc.isErrorExpected {
				require.Error(t, err)
				require.Nil(t, got)
				require.EqualError(t, err, tc.expectedErrorMessage)

			}
			if !tc.isErrorExpected {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tc.paramPolicyType, *got)
			}
		})
	}
}

// TestPolicyStore_CacheKey_DoesNotNormalizeTraversalLikeNames verifies that
// policy cache keys preserve traversal-like policy names and are never path-normalized.
func TestPolicyStore_CacheKey_DoesNotNormalizeTraversalLikeNames(t *testing.T) {
	t.Parallel()

	key := policyCacheKey("root", "../child/admin")
	require.Equal(t, "4:root:../child/admin", key)
}

// TestPolicyStore_CacheKey_NoCollisionsForTraversalLikeSiblingAndRoot verifies
// sibling/root traversal-like names cannot collide with valid root keys.
func TestPolicyStore_CacheKey_NoCollisionsForTraversalLikeSiblingAndRoot(t *testing.T) {
	t.Parallel()

	rootKey := policyCacheKey("root", "admin")
	siblingAttemptKey := policyCacheKey("sibling", "../root/admin")
	require.NotEqual(t, rootKey, siblingAttemptKey)
}

// TestPolicyStore_GetNonEGPPolicyType_UsesSharedCacheKeyFormat verifies
// GetNonEGPPolicyType only uses the shared key format and does not match
// legacy path-joined keys.
func TestPolicyStore_GetNonEGPPolicyType_UsesSharedCacheKeyFormat(t *testing.T) {
	t.Parallel()

	ps := new(PolicyStore)

	ps.policyTypeMap.Store("root/dev", PolicyTypeACL)
	policyType, err := ps.GetNonEGPPolicyType("root", "dev")
	require.Error(t, err)
	require.Nil(t, policyType)

	key := policyCacheKey("root", "dev")
	ps.policyTypeMap.Store(key, PolicyTypeACL)
	policyType, err = ps.GetNonEGPPolicyType("root", "dev")
	require.NoError(t, err)
	require.NotNil(t, policyType)
	require.Equal(t, PolicyTypeACL, *policyType)
}

// TestPolicyStore_SetPolicyWithRequest_PolicyNameValidation verifies direct
// policy creation rejects invalid policy names before storage.
func TestPolicyStore_SetPolicyWithRequest_PolicyNameValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		policy string
	}{
		{name: "traversal-like", policy: "../team/read"},
		{name: "parent in middle", policy: "team/../read"},
		{name: "dot", policy: "."},
		{name: "dot segment in middle", policy: "team/./read"},
		{name: "empty canonical", policy: "   "},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := new(PolicyStore)

			p := &Policy{
				Name:      tc.policy,
				Raw:       aclPolicy,
				Type:      PolicyTypeACL,
				namespace: namespace.RootNamespace,
			}
			err := ps.SetPolicyWithRequest(context.Background(), p, nil)
			require.Error(t, err)
			require.Contains(t, err.Error(), "invalid policy name")
		})
	}
}

// TestPolicyStore_GetPolicy_DoesNotResolveInvalidNameFromCaches verifies that
// invalid policy names are treated as unresolved before token/EGP cache lookup.
func TestPolicyStore_GetPolicy_DoesNotResolveInvalidNameFromCaches(t *testing.T) {
	t.Parallel()

	cluster := NewTestCluster(t, nil, nil)
	ps := cluster.Cores[0].Core.policyStore
	ctx := namespace.RootContext(t.Context())

	invalidPolicyName := "sibling/../admin"
	invalidKey := policyCacheKey(namespace.RootNamespaceID, ps.sanitizeName(invalidPolicyName))
	ps.policyTypeMap.Store(invalidKey, PolicyTypeRGP)

	ps.tokenPoliciesLRU.Add(invalidKey, &Policy{
		Name:      invalidPolicyName,
		Type:      PolicyTypeRGP,
		Raw:       `path "secret/data/forbidden" { capabilities = ["read"] }`,
		namespace: namespace.RootNamespace,
	})
	policy, err := ps.GetPolicy(ctx, invalidPolicyName, PolicyTypeToken)
	require.NoError(t, err)
	require.Nil(t, policy)

	ps.egpLRU.Add(invalidKey, &Policy{
		Name:      invalidPolicyName,
		Type:      PolicyTypeEGP,
		Raw:       `main = "rule"`,
		namespace: namespace.RootNamespace,
	})
	policy, err = ps.GetPolicy(ctx, invalidPolicyName, PolicyTypeEGP)
	require.NoError(t, err)
	require.Nil(t, policy)

	policyType, err := ps.GetNonEGPPolicyType(namespace.RootNamespaceID, invalidPolicyName)
	require.ErrorIs(t, err, ErrPolicyNotExistInTypeMap)
	require.Nil(t, policyType)
}

// TestPolicyStore_ACL_SkipsLegacyInvalidPolicyNames verifies ACL construction
// ignores invalid legacy policy references even when they exist in storage,
// policyTypeMap, and cache.
func TestPolicyStore_ACL_SkipsLegacyInvalidPolicyNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		invalidPolicyName string
		primeCache        bool
	}{
		{
			name:              "cached sibling target",
			invalidPolicyName: "sibling/./admin",
			primeCache:        true,
		},
		{
			name:              "uncached root target",
			invalidPolicyName: "root/./admin",
			primeCache:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, ps := mockPolicyWithCore(t, false)
			ctx := namespace.RootContext(context.Background())

			validPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/allowed" {
	capabilities = ["read"]
}
`, WithDenySlashInTemplatedPaths(false))
			require.NoError(t, err)
			validPolicy.Name = "legacy-valid-acl"
			validPolicy.namespace = namespace.RootNamespace
			require.NoError(t, ps.SetPolicy(ctx, validPolicy))

			legacyInvalidPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/forbidden" {
	capabilities = ["read"]
}
`, WithDenySlashInTemplatedPaths(false))
			require.NoError(t, err)
			legacyInvalidPolicy.Name = tc.invalidPolicyName
			legacyInvalidPolicy.namespace = namespace.RootNamespace
			require.NoError(t, ps.setPolicyInternal(ctx, legacyInvalidPolicy, nil))

			invalidKey := policyCacheKey(namespace.RootNamespaceID, ps.sanitizeName(tc.invalidPolicyName))
			_, found := ps.policyTypeMap.Load(invalidKey)
			require.True(t, found)
			storedEntry, err := ps.getACLView(namespace.RootNamespace).Get(ctx, tc.invalidPolicyName)
			require.NoError(t, err)
			require.NotNil(t, storedEntry)

			if !tc.primeCache {
				ps.tokenPoliciesLRU.Remove(invalidKey)
			}

			acl, err := ps.ACL(ctx, nil, map[string][]string{
				namespace.RootNamespaceID: {"legacy-valid-acl", tc.invalidPolicyName},
			})
			require.NoError(t, err)

			allowed := acl.AllowOperation(ctx, &logical.Request{
				Path:      "secret/data/allowed",
				Operation: logical.ReadOperation,
			}, false)
			require.True(t, allowed.Allowed)

			forbidden := acl.AllowOperation(ctx, &logical.Request{
				Path:      "secret/data/forbidden",
				Operation: logical.ReadOperation,
			}, false)
			require.False(t, forbidden.Allowed)
		})
	}
}

// TestPolicyStore_DuplicateAttributes checks that the policyStore.ACL method rejects templated
// policies with duplicate attributes. The VAULT_ALLOW_PENDING_REMOVAL_DUPLICATE_HCL_ATTRIBUTES
// environment variable has been removed, so duplicate attributes now always fail.
func TestPolicyStore_DuplicateAttributes(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ps := core.policyStore
	dupAttrPolicy := aclPolicy + `
path "foo" {
	capabilities = ["list"]
	capabilities = ["read"]
}
`
	// ParseACLPolicy now rejects duplicate attributes, so construct the policy manually
	// to store the duplicate raw text and verify that re-parsing it fails.
	policy := &Policy{
		Name:      "dev",
		Type:      PolicyTypeACL,
		Templated: true,
		Raw:       dupAttrPolicy,
		namespace: namespace.RootNamespace,
	}
	ctx := namespace.RootContext(context.Background())
	err := ps.SetPolicy(ctx, policy)
	require.NoError(t, err)

	_, err = ps.ACL(ctx, nil, map[string][]string{namespace.RootNamespace.ID: {"dev", "ops"}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "error parsing templated policy \"dev\": failed to parse policy: The argument \"capabilities\" at 61:2 was already set. Each argument can only be defined once")
	ps.tokenPoliciesLRU.Purge()
	_, err = ps.GetPolicy(ctx, "dev", PolicyTypeACL)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse policy: failed to parse policy: The argument \"capabilities\" at 61:2 was already set. Each argument can only be defined once")
}

// TestPolicyStore_AllowedParametersWarning tests that a warning is logged when a policy containing
// allowed_parameters or denied_parameters is set in the policy store.
// TODO (DENIED_PARAMETERS_CHANGE): Remove this test after deprecation is done
func TestPolicyStore_AllowedParametersWarning(t *testing.T) {
	tests := []struct {
		name           string
		policyFragment string
		expectLog      bool
	}{
		{
			name: "allowed_parameters",
			policyFragment: `
path "foo" {
 allowed_parameters = {
  "param1" = ["val1", "val2"]
 }
 capabilities = ["read"]
}
`,
			expectLog: true,
		},
		{
			name: "denied_parameters",
			policyFragment: `
path "foo" {
 denied_parameters = {
  "param1" = ["val1", "val2"]
 }
 capabilities = ["read"]
}
`,
			expectLog: true,
		},
		{
			name: "no_parameters",
			policyFragment: `
path "foo" {
 capabilities = ["read"]
}
`,
			expectLog: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logOut := new(bytes.Buffer)
			conf := &CoreConfig{
				Logger: log.New(&log.LoggerOptions{
					Mutex:  &sync.Mutex{},
					Level:  log.Warn,
					Output: logOut,
				}),
			}
			core, _, _ := TestCoreUnsealedWithConfig(t, conf)
			ps := core.policyStore

			// First policy
			policy := aclPolicy + tc.policyFragment
			parsedPolicy, err := ParseACLPolicy(namespace.RootNamespace, policy, WithDenySlashInTemplatedPaths(core.denySlashInTemplatedPolicyPaths))
			require.NoError(t, err)

			ctx := namespace.RootContext(context.Background())
			err = ps.SetPolicy(ctx, parsedPolicy)
			require.NoError(t, err)

			if tc.expectLog {
				require.Contains(t, logOut.String(), "you're using 'allowed_parameters' or 'denied_parameters' in one or more policies")
			} else {
				require.NotContains(t, logOut.String(), "you're using 'allowed_parameters' or 'denied_parameters' in one or more policies")
			}

			// Reset log output and add a second policy
			logOut.Reset()
			err = ps.SetPolicy(ctx, parsedPolicy)
			require.NoError(t, err)

			// Ensure no additional log is generated for the second policy
			require.NotContains(t, logOut.String(), "you're using 'allowed_parameters' or 'denied_parameters' in one or more policies")
		})
	}
}
