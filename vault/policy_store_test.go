package vault

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func mockPolicyStore(t *testing.T) *PolicyStore {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	p := NewPolicyStore(view, logical.TestSystemView())
	return p
}

func mockPolicyStoreNoCache(t *testing.T) *PolicyStore {
	sysView := logical.TestSystemView()
	sysView.CachingDisabledVal = true
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	p := NewPolicyStore(view, sysView)
	return p
}

func TestPolicyStore_Root(t *testing.T) {
	ps := mockPolicyStore(t)

	// Get should return a special policy
	p, err := ps.GetPolicy("root")
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
	err = ps.SetPolicy(p)
	if err.Error() != "cannot update root policy" {
		t.Fatalf("err: %v", err)
	}

	// Delete should fail
	err = ps.DeletePolicy("root")
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
	p, err := ps.GetPolicy("dev")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if p != nil {
		t.Fatalf("bad: %v", p)
	}

	// Delete should be no-op
	err = ps.DeletePolicy("dev")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// List should be blank
	out, err := ps.ListPolicies()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("bad: %v", out)
	}

	// Set should work
	policy, _ := Parse(aclPolicy)
	err = ps.SetPolicy(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should work
	p, err = ps.GetPolicy("dev")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(p, policy) {
		t.Fatalf("bad: %v", p)
	}

	// List should be one element
	out, err = ps.ListPolicies()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 1 || out[0] != "dev" {
		t.Fatalf("bad: %v", out)
	}

	// Delete should be clear the entry
	err = ps.DeletePolicy("dev")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get should fail
	p, err = ps.GetPolicy("dev")
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
	err := core.setupPolicyStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// List should be two elements
	out, err := core.policyStore.ListPolicies()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// This shouldn't contain response-wrapping since it's non-assignable
	if len(out) != 1 || out[0] != "default" {
		t.Fatalf("bad: %v", out)
	}

	pCubby, err := core.policyStore.GetPolicy("response-wrapping")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if pCubby.Raw != responseWrappingPolicy {
		t.Fatalf("bad: expected\n%s\ngot\n%s\n", responseWrappingPolicy, pCubby.Raw)
	}
	pRoot, err := core.policyStore.GetPolicy("root")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = core.policyStore.SetPolicy(pCubby)
	if err == nil {
		t.Fatalf("expected err setting %s", pCubby.Name)
	}
	err = core.policyStore.SetPolicy(pRoot)
	if err == nil {
		t.Fatalf("expected err setting %s", pRoot.Name)
	}
	err = core.policyStore.DeletePolicy(pCubby.Name)
	if err == nil {
		t.Fatalf("expected err deleting %s", pCubby.Name)
	}
	err = core.policyStore.DeletePolicy(pRoot.Name)
	if err == nil {
		t.Fatalf("expected err deleting %s", pRoot.Name)
	}
}

func TestPolicyStore_ACL(t *testing.T) {
	ps := mockPolicyStore(t)

	policy, _ := Parse(aclPolicy)
	err := ps.SetPolicy(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	policy, _ = Parse(aclPolicy2)
	err = ps.SetPolicy(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err := ps.ACL("dev", "ops")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testLayeredACL(t, acl)
}

func TestPolicyStore_v1Upgrade(t *testing.T) {
	ps := mockPolicyStore(t)

	// Put a V1 record
	raw := `path "foo" { policy = "read" }`
	ps.view.Put(&logical.StorageEntry{Key: "old", Value: []byte(raw)})

	// Do a read
	p, err := ps.GetPolicy("old")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if p == nil || len(p.Paths) != 1 {
		t.Fatalf("bad policy: %#v", p)
	}

	// Check that glob is enabled
	if !p.Paths[0].Glob {
		t.Fatalf("should enable glob")
	}
}
