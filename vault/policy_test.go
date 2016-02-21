package vault

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPolicy_Parse(t *testing.T) {
	p, err := Parse(rawPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if p.Name != "dev" {
		t.Fatalf("bad: %#v", p)
	}

	expect := []*PathCapabilities{
		&PathCapabilities{"", "deny",
			[]string{
				"deny",
			}, DenyCapabilityInt, true},
		&PathCapabilities{"stage/", "sudo",
			[]string{
				"create",
				"read",
				"update",
				"delete",
				"list",
				"sudo",
			}, CreateCapabilityInt | ReadCapabilityInt | UpdateCapabilityInt |
				DeleteCapabilityInt | ListCapabilityInt | SudoCapabilityInt, true},
		&PathCapabilities{"prod/version", "read",
			[]string{
				"read",
				"list",
			}, ReadCapabilityInt | ListCapabilityInt, false},
		&PathCapabilities{"foo/bar", "read",
			[]string{
				"read",
				"list",
			}, ReadCapabilityInt | ListCapabilityInt, false},
		&PathCapabilities{"foo/bar", "",
			[]string{
				"create",
				"sudo",
			}, CreateCapabilityInt | SudoCapabilityInt, false},
	}
	if !reflect.DeepEqual(p.Paths, expect) {
		ret := fmt.Sprintf("bad:\nexpected:\n")
		for _, v := range expect {
			ret = fmt.Sprintf("%s\n%#v", ret, *v)
		}
		ret = fmt.Sprintf("%s\n\ngot:\n", ret)
		for _, v := range p.Paths {
			ret = fmt.Sprintf("%s\n%#v", ret, *v)
		}
		t.Fatalf("%s\n", ret)
	}
}

var rawPolicy = `
# Developer policy
name = "dev"

# Deny all paths by default
path "*" {
	policy = "deny"
}

# Allow full access to staging
path "stage/*" {
	policy = "sudo"
}

# Limited read privilege to production
path "prod/version" {
	policy = "read"
}

# Read access to foobar
path "foo/bar" {
	policy = "read"
}

# Add capabilities for creation and sudo to foobar
# This will be separate; they are combined when compiled into an ACL
path "foo/bar" {
	capabilities = ["create", "sudo"]
}
`
