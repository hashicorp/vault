package vault

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var rawPolicy = strings.TrimSpace(`
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
# Also tests stripping of leading slash
path "/foo/bar" {
	policy = "read"
}

# Add MFA method required for "sudo" on foobar
path "foo/bar" {
	capabilities = ["sudo"]
	mfa_methods = ["mfa1"]
}

# Add capabilities for creation and sudo to foobar
# This will be separate; they are combined when compiled into an ACL
path "foo/bar" {
	capabilities = ["create", "sudo"]
}

# Add a second method to detect merging
path "foo/bar" {
	capabilities = ["sudo"]
	mfa_methods = ["mfa2"]
}
`)

func TestPolicy_Parse(t *testing.T) {
	p, err := Parse(rawPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if p.Name != "dev" {
		t.Fatalf("bad name: %q", p.Name)
	}

	expect := []*PathCapabilities{
		&PathCapabilities{"", "deny",
			[]string{
				"deny",
			}, DenyCapabilityInt, []string(nil), true},
		&PathCapabilities{"stage/", "sudo",
			[]string{
				"create",
				"read",
				"update",
				"delete",
				"list",
				"sudo",
			}, CreateCapabilityInt | ReadCapabilityInt | UpdateCapabilityInt |
				DeleteCapabilityInt | ListCapabilityInt | SudoCapabilityInt, []string(nil), true},
		&PathCapabilities{"prod/version", "read",
			[]string{
				"read",
				"list",
			}, ReadCapabilityInt | ListCapabilityInt, []string(nil), false},
		&PathCapabilities{"foo/bar", "read",
			[]string{
				"read",
				"list",
			}, ReadCapabilityInt | ListCapabilityInt, []string(nil), false},
		&PathCapabilities{"foo/bar", "",
			[]string{
				"sudo",
			}, SudoCapabilityInt, []string{"mfa1"}, false},
		&PathCapabilities{"foo/bar", "",
			[]string{
				"create",
				"sudo",
			}, CreateCapabilityInt | SudoCapabilityInt, []string(nil), false},
		&PathCapabilities{"foo/bar", "",
			[]string{
				"sudo",
			}, SudoCapabilityInt, []string{"mfa2"}, false},
	}
	if !reflect.DeepEqual(p.Paths, expect) {
		var foundPaths, expectPaths string
		for _, path := range p.Paths {
			foundPaths = fmt.Sprintf("%s\n%#v", foundPaths, *path)
		}
		for _, path := range expect {
			expectPaths = fmt.Sprintf("%s\n%#v", expectPaths, *path)
		}
		t.Errorf("expected \n\n%s\n\n to be \n\n%s\n\n", foundPaths, expectPaths)
	}
}

func TestPolicy_ParseBadRoot(t *testing.T) {
	_, err := Parse(strings.TrimSpace(`
name = "test"
bad  = "foo"
nope = "yes"
`))
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), "invalid key 'bad' on line 2") {
		t.Errorf("bad error: %q", err)
	}

	if !strings.Contains(err.Error(), "invalid key 'nope' on line 3") {
		t.Errorf("bad error: %q", err)
	}
}

func TestPolicy_ParseBadPath(t *testing.T) {
	_, err := Parse(strings.TrimSpace(`
path "/" {
	capabilities = ["read"]
	capabilites  = ["read"]
}
`))
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), "invalid key 'capabilites' on line 3") {
		t.Errorf("bad error: %s", err)
	}
}

func TestPolicy_ParseBadPolicy(t *testing.T) {
	_, err := Parse(strings.TrimSpace(`
path "/" {
	policy = "banana"
}
`))
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), `path "/": invalid policy 'banana'`) {
		t.Errorf("bad error: %s", err)
	}
}

func TestPolicy_ParseBadCapabilities(t *testing.T) {
	_, err := Parse(strings.TrimSpace(`
path "/" {
	capabilities = ["read", "banana"]
}
`))
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), `path "/": invalid capability 'banana'`) {
		t.Errorf("bad error: %s", err)
	}
}
