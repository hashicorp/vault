package vault

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func mockPolicyWithCore(t *testing.T, disableCache bool) (*Core, *PolicyStore) {
	conf := &CoreConfig{
		DisableCache: disableCache,
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, conf)
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
	if len(out) != 1 {
		t.Fatalf("bad: %v", out)
	}

	// Set should work
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	policy, _ := ParseACLPolicy(ns, aclPolicy)
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

	// List should contain two elements
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	out, err = ps.ListPolicies(ctx, PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("bad: %v", out)
	}

	expected := []string{"default", "dev"}
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
	if len(out) != 1 || out[0] != "default" {
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
	// List should be two elements
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	out, err := ps.ListPolicies(ctx, PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// This shouldn't contain response-wrapping since it's non-assignable
	if len(out) != 1 || out[0] != "default" {
		t.Fatalf("bad: %v", out)
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
	policy, _ := ParseACLPolicy(ns, aclPolicy)
	err := ps.SetPolicy(ctx, policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	policy, _ = ParseACLPolicy(ns, aclPolicy2)
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

	parsedPolicy, _ := ParseACLPolicy(rootNs, aclPolicy)

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

	expectedResult := []string{"default", "dev"}
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
			policyStoreKey:   "1AbcD/policy1",
			policyStoreValue: PolicyTypeACL,
			paramNamespace:   "1AbcD",
			paramPolicyName:  "policy1",
			paramPolicyType:  PolicyTypeACL,
		},
		"happy-rgp": {
			policyStoreKey:   "1AbcD/policy1",
			policyStoreValue: PolicyTypeRGP,
			paramNamespace:   "1AbcD",
			paramPolicyName:  "policy1",
			paramPolicyType:  PolicyTypeRGP,
		},
		"not-in-map-acl": {
			policyStoreKey:       "2WxyZ/policy2",
			policyStoreValue:     PolicyTypeACL,
			paramNamespace:       "1AbcD",
			paramPolicyName:      "policy1",
			isErrorExpected:      true,
			expectedErrorMessage: "policy does not exist in type map",
		},
		"not-in-map-rgp": {
			policyStoreKey:       "2WxyZ/policy2",
			policyStoreValue:     PolicyTypeRGP,
			paramNamespace:       "1AbcD",
			paramPolicyName:      "policy1",
			isErrorExpected:      true,
			expectedErrorMessage: "policy does not exist in type map",
		},
		"unknown-policy-type": {
			policyStoreKey:       "1AbcD/policy1",
			policyStoreValue:     7,
			paramNamespace:       "1AbcD",
			paramPolicyName:      "policy1",
			isErrorExpected:      true,
			expectedErrorMessage: "unknown policy type for: 1AbcD/policy1",
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
