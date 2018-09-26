package vault

import (
	"context"
	"reflect"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
)

func mockPolicyStore(t *testing.T) *PolicyStore {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	p, err := NewPolicyStore(context.Background(), nil, view, logical.TestSystemView(), logging.NewVaultLogger(log.Trace))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func mockPolicyStoreNoCache(t *testing.T) *PolicyStore {
	sysView := logical.TestSystemView()
	sysView.CachingDisabledVal = true
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	p, err := NewPolicyStore(context.Background(), nil, view, sysView, logging.NewVaultLogger(log.Trace))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

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
	acl, err := ps.ACL(ctx, nil, map[string][]string{ns.ID: []string{"dev", "ops"}})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testLayeredACL(t, acl, ns)
}
