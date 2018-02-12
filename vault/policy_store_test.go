package vault

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/logical"
	log "github.com/mgutz/logxi/v1"
)

func mockPolicyStore(t *testing.T) *PolicyStore {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	p := NewPolicyStore(context.Background(), view, logical.TestSystemView(), logformat.NewVaultLogger(log.LevelTrace))
	return p
}

func mockPolicyStoreNoCache(t *testing.T) *PolicyStore {
	sysView := logical.TestSystemView()
	sysView.CachingDisabledVal = true
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	p := NewPolicyStore(context.Background(), view, sysView, logformat.NewVaultLogger(log.LevelTrace))
	return p
}

func TestPolicyStore_Root(t *testing.T) {
	ps := mockPolicyStore(t)

	// Get should return a special policy
	p, err := ps.GetPolicy(context.Background(), "root", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if p == nil {
		t.Fatalf("bad: %v", p)
	}
	if p.Name != "root" {
		t.Fatalf("bad: %v", p)
	}

	// Set should fail
	err = ps.SetPolicy(context.Background(), p)
	if err.Error() != "cannot update root policy" {
		t.Fatalf("err: %v", err)
	}

	// Delete should fail
	err = ps.DeletePolicy(context.Background(), "root", PolicyTypeACL)
	if err.Error() != "cannot delete root policy" {
		t.Fatalf("err: %v", err)
	}
}

func TestPolicyStore_CRUD(t *testing.T) {
	ps := mockPolicyStore(t)
	testPolicyStore_CRUD(t, ps)

	ps = mockPolicyStoreNoCache(t)
	testPolicyStore_CRUD(t, ps)
}

func testPolicyStore_CRUD(t *testing.T, ps *PolicyStore) {
	// Get should return nothing
	p, err := ps.GetPolicy(context.Background(), "Dev", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if p != nil {
		t.Fatalf("bad: %v", p)
	}

	// Delete should be no-op
	err = ps.DeletePolicy(context.Background(), "deV", PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// List should be blank
	out, err := ps.ListPolicies(context.Background(), PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("bad: %v", out)
	}

	// Set should work
	policy, _ := ParseACLPolicy(aclPolicy)
	err = ps.SetPolicy(context.Background(), policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should work
	p, err = ps.GetPolicy(context.Background(), "dEv", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(p, policy) {
		t.Fatalf("bad: %v", p)
	}

	// List should be one element
	out, err = ps.ListPolicies(context.Background(), PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 1 || out[0] != "dev" {
		t.Fatalf("bad: %v", out)
	}

	// Delete should be clear the entry
	err = ps.DeletePolicy(context.Background(), "Dev", PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should fail
	p, err = ps.GetPolicy(context.Background(), "deV", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if p != nil {
		t.Fatalf("bad: %v", p)
	}
}

// Test predefined policy handling
func TestPolicyStore_Predefined(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	// Ensure both default policies are created
	err := core.setupPolicyStore(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// List should be two elements
	out, err := core.policyStore.ListPolicies(context.Background(), PolicyTypeACL)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// This shouldn't contain response-wrapping since it's non-assignable
	if len(out) != 1 || out[0] != "default" {
		t.Fatalf("bad: %v", out)
	}

	pCubby, err := core.policyStore.GetPolicy(context.Background(), "response-wrapping", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if pCubby == nil {
		t.Fatal("nil cubby policy")
	}
	if pCubby.Raw != responseWrappingPolicy {
		t.Fatalf("bad: expected\n%s\ngot\n%s\n", responseWrappingPolicy, pCubby.Raw)
	}
	pRoot, err := core.policyStore.GetPolicy(context.Background(), "root", PolicyTypeToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if pRoot == nil {
		t.Fatal("nil root policy")
	}

	err = core.policyStore.SetPolicy(context.Background(), pCubby)
	if err == nil {
		t.Fatalf("expected err setting %s", pCubby.Name)
	}
	err = core.policyStore.SetPolicy(context.Background(), pRoot)
	if err == nil {
		t.Fatalf("expected err setting %s", pRoot.Name)
	}
	err = core.policyStore.DeletePolicy(context.Background(), pCubby.Name, PolicyTypeACL)
	if err == nil {
		t.Fatalf("expected err deleting %s", pCubby.Name)
	}
	err = core.policyStore.DeletePolicy(context.Background(), pRoot.Name, PolicyTypeACL)
	if err == nil {
		t.Fatalf("expected err deleting %s", pRoot.Name)
	}
}

func TestPolicyStore_ACL(t *testing.T) {
	ps := mockPolicyStore(t)

	policy, _ := ParseACLPolicy(aclPolicy)
	err := ps.SetPolicy(context.Background(), policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	policy, _ = ParseACLPolicy(aclPolicy2)
	err = ps.SetPolicy(context.Background(), policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err := ps.ACL(context.Background(), "dev", "ops")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testLayeredACL(t, acl)
}
