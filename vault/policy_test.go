package vault

import (
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

	expect := []*PathPolicy{
		&PathPolicy{"", "deny"},
		&PathPolicy{"stage/", "sudo"},
		&PathPolicy{"prod/", "read"},
	}
	if !reflect.DeepEqual(p.Paths, expect) {
		t.Fatalf("bad: %#v", p)
	}
}

var rawPolicy = `
# Developer policy
name = "dev"

# Deny all paths by default
path "" {
	policy = "deny"
}

# Allow full access to staging
path "stage/" {
	policy = "sudo"
}

# Limited read privilege to production
path "prod/" {
	policy = "read"
}
`
