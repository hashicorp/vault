package vault

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func mockPolicyStore(t *testing.T) *PolicyStore {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	p := NewPolicyStore(view)
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
	ps.view.Put(&logical.StorageEntry{"old", []byte(raw)})

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
