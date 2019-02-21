package vault

import (
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
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
path "test/mfa" {
	capabilities = ["create", "sudo"]
	mfa_methods = ["my_totp", "my_totp2"]
}
path "test/+/segment" {
	capabilities = ["create", "sudo"]
}
path "test/segment/at/end/+" {
	capabilities = ["create", "sudo"]
}
path "test/segment/at/end/v2/+/" {
	capabilities = ["create", "sudo"]
}
path "test/+/wildcard/+/*" {
	capabilities = ["create", "sudo"]
}
path "test/+/wildcard/+/end*" {
	capabilities = ["create", "sudo"]
}
`)

func TestPolicy_Parse(t *testing.T) {
	p, err := ParseACLPolicy(namespace.RootNamespace, rawPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if p.Name != "dev" {
		t.Fatalf("bad name: %q", p.Name)
	}

	expect := []*PathRules{
		{
			Path:   "",
			Policy: "deny",
			Capabilities: []string{
				"deny",
			},
			Permissions: &ACLPermissions{CapabilitiesBitmap: DenyCapabilityInt},
			IsPrefix:    true,
		},
		{
			Path:   "stage/",
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
			IsPrefix: true,
		},
		{
			Path:   "prod/version",
			Policy: "read",
			Capabilities: []string{
				"read",
				"list",
			},
			Permissions: &ACLPermissions{CapabilitiesBitmap: (ReadCapabilityInt | ListCapabilityInt)},
		},
		{
			Path:   "foo/bar",
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
		},
		{
			Path: "foo/bar",
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
		},
		{
			Path: "foo/bar",
			Capabilities: []string{
				"create",
				"sudo",
			},
			AllowedParametersHCL: map[string][]interface{}{"zip": {}, "zap": {}},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
				AllowedParameters:  map[string][]interface{}{"zip": {}, "zap": {}},
			},
		},
		{
			Path: "baz/bar",
			Capabilities: []string{
				"create",
				"sudo",
			},
			DeniedParametersHCL: map[string][]interface{}{"zip": []interface{}{}, "zap": []interface{}{}},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
				DeniedParameters:   map[string][]interface{}{"zip": []interface{}{}, "zap": []interface{}{}},
			},
		},
		{
			Path: "biz/bar",
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
		},
		{
			Path:   "test/types",
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
			IsPrefix: false,
		},
		{
			Path: "test/req",
			Capabilities: []string{
				"create",
				"sudo",
			},
			RequiredParametersHCL: []string{"foo"},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
				RequiredParameters: []string{"foo"},
			},
		},
		{
			Path: "test/mfa",
			Capabilities: []string{
				"create",
				"sudo",
			},
			MFAMethodsHCL: []string{
				"my_totp",
				"my_totp2",
			},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
				MFAMethods: []string{
					"my_totp",
					"my_totp2",
				},
			},
		},
		{
			Path: "test/+/segment",
			Capabilities: []string{
				"create",
				"sudo",
			},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
			},
			HasSegmentWildcards: true,
		},
		{
			Path: "test/segment/at/end/+",
			Capabilities: []string{
				"create",
				"sudo",
			},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
			},
			HasSegmentWildcards: true,
		},
		{
			Path: "test/segment/at/end/v2/+/",
			Capabilities: []string{
				"create",
				"sudo",
			},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
			},
			HasSegmentWildcards: true,
		},
		{
			Path: "test/+/wildcard/+/*",
			Capabilities: []string{
				"create",
				"sudo",
			},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
			},
			HasSegmentWildcards: true,
		},
		{
			Path: "test/+/wildcard/+/end*",
			Capabilities: []string{
				"create",
				"sudo",
			},
			Permissions: &ACLPermissions{
				CapabilitiesBitmap: (CreateCapabilityInt | SudoCapabilityInt),
			},
			HasSegmentWildcards: true,
		},
	}

	if diff := deep.Equal(p.Paths, expect); diff != nil {
		t.Error(diff)
	}
}

func TestPolicy_ParseBadRoot(t *testing.T) {
	_, err := ParseACLPolicy(namespace.RootNamespace, strings.TrimSpace(`
name = "test"
bad  = "foo"
nope = "yes"
`))
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), `invalid key "bad" on line 2`) {
		t.Errorf("bad error: %q", err)
	}

	if !strings.Contains(err.Error(), `invalid key "nope" on line 3`) {
		t.Errorf("bad error: %q", err)
	}
}

func TestPolicy_ParseBadPath(t *testing.T) {
	// The wrong spelling is intended here
	_, err := ParseACLPolicy(namespace.RootNamespace, strings.TrimSpace(`
path "/" {
	capabilities = ["read"]
	capabilites  = ["read"]
}
`))
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), `invalid key "capabilites" on line 3`) {
		t.Errorf("bad error: %s", err)
	}
}

func TestPolicy_ParseBadPolicy(t *testing.T) {
	_, err := ParseACLPolicy(namespace.RootNamespace, strings.TrimSpace(`
path "/" {
	policy = "banana"
}
`))
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), `path "/": invalid policy "banana"`) {
		t.Errorf("bad error: %s", err)
	}
}

func TestPolicy_ParseBadWrapping(t *testing.T) {
	_, err := ParseACLPolicy(namespace.RootNamespace, strings.TrimSpace(`
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
	_, err := ParseACLPolicy(namespace.RootNamespace, strings.TrimSpace(`
path "/" {
	capabilities = ["read", "banana"]
}
`))
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), `path "/": invalid capability "banana"`) {
		t.Errorf("bad error: %s", err)
	}
}
