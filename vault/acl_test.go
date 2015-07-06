package vault

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestACL_Root(t *testing.T) {
	// Create the root policy ACL
	policy := []*Policy{&Policy{Name: "root"}}
	acl, err := NewACL(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !acl.RootPrivilege("sys/mount/foo") {
		t.Fatalf("expected root")
	}
	if !acl.AllowOperation(logical.WriteOperation, "sys/mount/foo") {
		t.Fatalf("expected permission")
	}
}

func TestACL_Single(t *testing.T) {
	policy, err := Parse(aclPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err := NewACL([]*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if acl.RootPrivilege("sys/mount/foo") {
		t.Fatalf("unexpected root")
	}

	type tcase struct {
		op     logical.Operation
		path   string
		expect bool
	}
	tcases := []tcase{
		{logical.ReadOperation, "root", false},
		{logical.HelpOperation, "root", true},

		{logical.ReadOperation, "dev/foo", true},
		{logical.WriteOperation, "dev/foo", true},

		{logical.DeleteOperation, "stage/foo", true},
		{logical.WriteOperation, "stage/aws/foo", false},
		{logical.WriteOperation, "stage/aws/policy/foo", true},

		{logical.DeleteOperation, "prod/foo", false},
		{logical.WriteOperation, "prod/foo", false},
		{logical.ReadOperation, "prod/foo", true},
		{logical.ListOperation, "prod/foo", true},
		{logical.ReadOperation, "prod/aws/foo", false},
	}

	for _, tc := range tcases {
		out := acl.AllowOperation(tc.op, tc.path)
		if out != tc.expect {
			t.Fatalf("bad: case %#v: %v", tc, out)
		}
	}
}

func TestACL_Layered(t *testing.T) {
	policy1, err := Parse(aclPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	policy2, err := Parse(aclPolicy2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err := NewACL([]*Policy{policy1, policy2})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testLayeredACL(t, acl)
}

func testLayeredACL(t *testing.T, acl *ACL) {
	if acl.RootPrivilege("sys/mount/foo") {
		t.Fatalf("unexpected root")
	}

	type tcase struct {
		op     logical.Operation
		path   string
		expect bool
	}
	tcases := []tcase{
		{logical.ReadOperation, "root", false},
		{logical.HelpOperation, "root", true},

		{logical.ReadOperation, "dev/hide/foo", false},
		{logical.WriteOperation, "dev/hide/foo", false},

		{logical.DeleteOperation, "stage/foo", true},
		{logical.WriteOperation, "stage/aws/foo", false},
		{logical.WriteOperation, "stage/aws/policy/foo", false},

		{logical.DeleteOperation, "prod/foo", true},
		{logical.WriteOperation, "prod/foo", true},
		{logical.ReadOperation, "prod/foo", true},
		{logical.ListOperation, "prod/foo", true},
		{logical.ReadOperation, "prod/aws/foo", false},

		{logical.ReadOperation, "sys/status", false},
		{logical.WriteOperation, "sys/seal", true},
	}

	for _, tc := range tcases {
		out := acl.AllowOperation(tc.op, tc.path)
		if out != tc.expect {
			t.Fatalf("bad: case %#v: %v", tc, out)
		}
	}
}

var aclPolicy = `
name = "dev"
path "dev/*" {
	policy = "sudo"
}
path "stage/*" {
	policy = "write"
}
path "stage/aws/*" {
	policy = "read"
}
path "stage/aws/policy/*" {
	policy = "sudo"
}
path "prod/*" {
	policy = "read"
}
path "prod/aws/*" {
	policy = "deny"
}
path "sys/*" {
	policy = "deny"
}
`

var aclPolicy2 = `
name = "ops"
path "dev/hide/*" {
	policy = "deny"
}
path "stage/aws/policy/*" {
	policy = "deny"
}
path "prod/*" {
	policy = "write"
}
path "sys/seal" {
	policy = "write"
}
`
