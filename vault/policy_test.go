package vault

import (
	"reflect"
	"testing"
)

func TestPolicy_TakesPrecedence(t *testing.T) {
	type tcase struct {
		a, b       string
		precedence bool
	}
	tests := []tcase{
		tcase{PathPolicyDeny, PathPolicyDeny, true},
		tcase{PathPolicyDeny, PathPolicyRead, true},
		tcase{PathPolicyDeny, PathPolicyWrite, true},
		tcase{PathPolicyDeny, PathPolicySudo, true},

		tcase{PathPolicyRead, PathPolicyDeny, false},
		tcase{PathPolicyRead, PathPolicyRead, false},
		tcase{PathPolicyRead, PathPolicyWrite, false},
		tcase{PathPolicyRead, PathPolicySudo, false},

		tcase{PathPolicyWrite, PathPolicyDeny, false},
		tcase{PathPolicyWrite, PathPolicyRead, true},
		tcase{PathPolicyWrite, PathPolicyWrite, false},
		tcase{PathPolicyWrite, PathPolicySudo, false},

		tcase{PathPolicySudo, PathPolicyDeny, false},
		tcase{PathPolicySudo, PathPolicyRead, true},
		tcase{PathPolicySudo, PathPolicyWrite, true},
		tcase{PathPolicySudo, PathPolicySudo, false},
	}
	for idx, test := range tests {
		a := &PathPolicy{Policy: test.a}
		b := &PathPolicy{Policy: test.b}
		if out := a.TakesPrecedence(b); out != test.precedence {
			t.Fatalf("bad: idx %d expect: %v out: %v",
				idx, test.precedence, out)
		}
	}
}

func TestPolicy_Parse(t *testing.T) {
	p, err := Parse(rawPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if p.Name != "dev" {
		t.Fatalf("bad: %#v", p)
	}

	expect := []*PathPolicy{
		&PathPolicy{"", "deny", true},
		&PathPolicy{"stage/", "sudo", true},
		&PathPolicy{"prod/version", "read", false},
	}
	if !reflect.DeepEqual(p.Paths, expect) {
		t.Fatalf("bad: %#v", p)
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
`
