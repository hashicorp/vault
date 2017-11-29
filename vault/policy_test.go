package vault

import (
	"reflect"
	"strings"
	"testing"
	"time"
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
# Also tests stripping of leading slash and parsing of min/max as string and
# integer
path "/foo/bar" {
	policy = "read"
	min_wrapping_ttl = 300
	max_wrapping_ttl = "1h"
}

# Add capabilities for creation and sudo to foobar
# This will be separate; they are combined when compiled into an ACL
# Also tests reverse string/int handling to the above
path "foo/bar" {
	capabilities = ["create", "sudo"]
	min_wrapping_ttl = "300s"
	max_wrapping_ttl = 3600
}

# Check that only allowed_parameters are being added to foobar
path "foo/bar" {
	capabilities = ["create", "sudo"]
	allowed_parameters = {
	  "zip" = []
	  "zap" = []
	}
}

# Check that only denied_parameters are being added to bazbar
path "baz/bar" {
	capabilities = ["create", "sudo"]
	denied_parameters = {
	  "zip" = []
	  "zap" = []
	}
}

# Check that both allowed and denied parameters are being added to bizbar
path "biz/bar" {
	capabilities = ["create", "sudo"]
	allowed_parameters = {
	  "zim" = []
	  "zam" = []
	}
	denied_parameters = {
	  "zip" = []
	  "zap" = []
	}
}
path "test/types" {
	capabilities = ["create", "sudo"]
	allowed_parameters = {
		"map" = [{"good" = "one"}]
		"int" = [1, 2]
	}
	denied_parameters = {
		"string" = ["test"]
		"bool" = [false]
	}
}
path "test/req" {
	capabilities = ["create", "sudo"]
	required_parameters = ["foo"]
}
`)

func TestPolicy_Parse(t *testing.T) {
	p, err := ParseACLPolicy(rawPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if p.Name != "dev" {
		t.Fatalf("bad name: %q", p.Name)
	}

	expect := []*PathRules{
		&PathRules{
			Prefix: "",
			Policy: "deny",
			Capabilities: []string{
				"deny",
			},
			Permissions: &ACLPermissions{CapabilitiesBitmap: DenyCapabilityInt},
			Glob:        true,
		},
		&PathRules{
			Prefix: "stage/",
			Policy: "sudo",
			Capabilities: []string{
				"create",
				"read",
				"update",
				"delete",
				"list",
				"sudo",
			},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | ReadCapabilityInt | UpdateCapabilityInt | DeleteCapabilityInt | ListCapabilityInt | SudoCapabilityInt),
			},
			Glob: true,
		},
		&PathRules{
			Prefix: "prod/version",
			Policy: "read",
			Capabilities: []string{
				"read",
				"list",
			},
			Permissions: &ACLPermissions{CapabilitiesBitmap: (ReadCapabilityInt | ListCapabilityInt)},
			Glob:        false,
		},
		&PathRules{
			Prefix: "foo/bar",
			Policy: "read",
			Capabilities: []string{
				"read",
				"list",
			},
			MinWrappingTTLHCL: 300,
			MaxWrappingTTLHCL: "1h",
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (ReadCapabilityInt | ListCapabilityInt),
				MinWrappingTTL:     300 * time.Second,
				MaxWrappingTTL:     3600 * time.Second,
			},
			Glob: false,
		},
		&PathRules{
			Prefix: "foo/bar",
			Policy: "",
			Capabilities: []string{
				"create",
				"sudo",
			},
			MinWrappingTTLHCL: "300s",
			MaxWrappingTTLHCL: 3600,
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
				MinWrappingTTL:     300 * time.Second,
				MaxWrappingTTL:     3600 * time.Second,
			},
			Glob: false,
		},
		&PathRules{
			Prefix: "foo/bar",
			Policy: "",
			Capabilities: []string{
				"create",
				"sudo",
			},
			AllowedParametersHCL: map[string][]interface{}{"zip": {}, "zap": {}},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
				AllowedParameters:  map[string][]interface{}{"zip": {}, "zap": {}},
			},
			Glob: false,
		},
		&PathRules{
			Prefix: "baz/bar",
			Policy: "",
			Capabilities: []string{
				"create",
				"sudo",
			},
			DeniedParametersHCL: map[string][]interface{}{"zip": []interface{}{}, "zap": []interface{}{}},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
				DeniedParameters:   map[string][]interface{}{"zip": []interface{}{}, "zap": []interface{}{}},
			},
			Glob: false,
		},
		&PathRules{
			Prefix: "biz/bar",
			Policy: "",
			Capabilities: []string{
				"create",
				"sudo",
			},
			AllowedParametersHCL: map[string][]interface{}{"zim": {}, "zam": {}},
			DeniedParametersHCL:  map[string][]interface{}{"zip": {}, "zap": {}},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
				AllowedParameters:  map[string][]interface{}{"zim": {}, "zam": {}},
				DeniedParameters:   map[string][]interface{}{"zip": {}, "zap": {}},
			},
			Glob: false,
		},
		&PathRules{
			Prefix: "test/types",
			Policy: "",
			Capabilities: []string{
				"create",
				"sudo",
			},
			AllowedParametersHCL: map[string][]interface{}{"map": []interface{}{map[string]interface{}{"good": "one"}}, "int": []interface{}{1, 2}},
			DeniedParametersHCL:  map[string][]interface{}{"string": []interface{}{"test"}, "bool": []interface{}{false}},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
				AllowedParameters:  map[string][]interface{}{"map": []interface{}{map[string]interface{}{"good": "one"}}, "int": []interface{}{1, 2}},
				DeniedParameters:   map[string][]interface{}{"string": []interface{}{"test"}, "bool": []interface{}{false}},
			},
			Glob: false,
		},
		&PathRules{
			Prefix: "test/req",
			Policy: "",
			Capabilities: []string{
				"create",
				"sudo",
			},
			RequiredParametersHCL: []string{"foo"},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
				RequiredParameters: []string{"foo"},
			},
			Glob: false,
		},
	}
	if !reflect.DeepEqual(p.Paths, expect) {
		t.Errorf("expected \n\n%#v\n\n to be \n\n%#v\n\n", p.Paths, expect)
	}
}

func TestPolicy_ParseBadRoot(t *testing.T) {
	_, err := ParseACLPolicy(strings.TrimSpace(`
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
	_, err := ParseACLPolicy(strings.TrimSpace(`
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
	_, err := ParseACLPolicy(strings.TrimSpace(`
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

func TestPolicy_ParseBadWrapping(t *testing.T) {
	_, err := ParseACLPolicy(strings.TrimSpace(`
path "/" {
	policy = "read"
	min_wrapping_ttl = 400
	max_wrapping_ttl = 200
}
`))
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), `max_wrapping_ttl cannot be less than min_wrapping_ttl`) {
		t.Errorf("bad error: %s", err)
	}
}

func TestPolicy_ParseBadCapabilities(t *testing.T) {
	_, err := ParseACLPolicy(strings.TrimSpace(`
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
