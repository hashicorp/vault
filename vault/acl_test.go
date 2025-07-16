// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestACL_NewACL(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		t.Parallel()
		testNewACL(t, namespace.RootNamespace)
	})
}

func testNewACL(t *testing.T, ns *namespace.Namespace) {
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	policy := []*Policy{{Name: "root"}}
	_, err := NewACL(ctx, policy)
	switch ns.ID {
	case namespace.RootNamespaceID:
		if err != nil {
			t.Fatal(err)
		}
	default:
		if err == nil {
			t.Fatal("expected an error")
		}
	}
}

func TestACL_MFAMethods(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		t.Parallel()
		testACLMFAMethods(t, namespace.RootNamespace)
	})
}

func testACLMFAMethods(t *testing.T, ns *namespace.Namespace) {
	mfaRules := `
path "secret/foo/*" {
	mfa_methods = ["mfa_method_1", "mfa_method_2", "mfa_method_3"]
}
path "secret/exact/path" {
	mfa_methods = ["mfa_method_4", "mfa_method_5"]
}
path "secret/split/definition" {
	mfa_methods = ["mfa_method_6", "mfa_method_7"]
}
path "secret/split/definition" {
	mfa_methods = ["mfa_method_7", "mfa_method_8", "mfa_method_9"]
}
	`

	policy, err := ParseACLPolicy(ns, mfaRules)
	if err != nil {
		t.Fatal(err)
	}

	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	acl, err := NewACL(ctx, []*Policy{policy})
	if err != nil {
		t.Fatal(err)
	}

	request := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/foo/testing/glob/pattern",
	}

	actual := acl.AllowOperation(ctx, request, false).MFAMethods
	expected := []string{"mfa_method_1", "mfa_method_2", "mfa_method_3"}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: MFA methods; expected: %#v\n actual: %#v\n", expected, actual)
	}

	request.Path = "secret/exact/path"
	actual = acl.AllowOperation(ctx, request, false).MFAMethods
	expected = []string{"mfa_method_4", "mfa_method_5"}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: MFA methods; expected: %#v\n actual: %#v\n", expected, actual)
	}

	request.Path = "secret/split/definition"
	actual = acl.AllowOperation(ctx, request, false).MFAMethods
	expected = []string{"mfa_method_6", "mfa_method_7", "mfa_method_8", "mfa_method_9"}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: MFA methods; expected: %#v\n actual: %#v\n", expected, actual)
	}
}

func TestACL_Capabilities(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		t.Parallel()
		policy := []*Policy{{Name: "root"}}
		ctx := namespace.RootContext(context.Background())
		acl, err := NewACL(ctx, policy)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		actual := acl.Capabilities(ctx, "any/path")
		expected := []string{"root"}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
		}
		testACLCapabilities(t, namespace.RootNamespace)
	})
}

func testACLCapabilities(t *testing.T, ns *namespace.Namespace) {
	// Create the root policy ACL
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	policy, err := ParseACLPolicy(ns, aclPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err := NewACL(ctx, []*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	actual := acl.Capabilities(ctx, "dev")
	expected := []string{"deny"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: path: %s\ngot\n%#v\nexpected\n%#v\n", "deny", actual, expected)
	}

	actual = acl.Capabilities(ctx, "dev/")
	expected = []string{"sudo", "read", "list", "update", "delete", "create"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: path: %s\ngot\n%#v\nexpected\n%#v\n", "dev/", actual, expected)
	}

	actual = acl.Capabilities(ctx, "stage/aws/test")
	expected = []string{"sudo", "read", "list", "update"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: path: %s\ngot\n%#v\nexpected\n%#v\n", "stage/aws/test", actual, expected)
	}
}

func TestACL_Root(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		t.Parallel()
		testACLRoot(t, namespace.RootNamespace)
	})
}

func testACLRoot(t *testing.T, ns *namespace.Namespace) {
	// Create the root policy ACL. Always create on root namespace regardless of
	// which namespace to ACL check on.
	policy := []*Policy{{Name: "root"}}
	acl, err := NewACL(namespace.RootContext(context.Background()), policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	request := new(logical.Request)
	request.Operation = logical.UpdateOperation
	request.Path = "sys/mount/foo"
	ctx := namespace.ContextWithNamespace(context.Background(), ns)

	authResults := acl.AllowOperation(ctx, request, false)
	if !authResults.RootPrivs {
		t.Fatalf("expected root")
	}
	if !authResults.Allowed {
		t.Fatalf("expected permissions")
	}
}

func TestACL_Single(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		t.Parallel()
		testACLSingle(t, namespace.RootNamespace)
	})
}

func testACLSingle(t *testing.T, ns *namespace.Namespace) {
	policy, err := ParseACLPolicy(ns, aclPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	acl, err := NewACL(ctx, []*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Type of operation is not important here as we only care about checking
	// sudo/root
	ctx = namespace.ContextWithNamespace(context.Background(), ns)
	request := new(logical.Request)
	request.Operation = logical.ReadOperation
	request.Path = "sys/mount/foo"

	authResults := acl.AllowOperation(ctx, request, false)
	if authResults.RootPrivs {
		t.Fatalf("unexpected root")
	}

	type tcase struct {
		op        logical.Operation
		path      string
		allowed   bool
		rootPrivs bool
	}
	tcases := []tcase{
		{logical.ReadOperation, "root", false, false},
		{logical.HelpOperation, "root", true, false},

		{logical.ReadOperation, "dev/foo", true, true},
		{logical.UpdateOperation, "dev/foo", true, true},

		{logical.DeleteOperation, "stage/foo", true, false},
		{logical.ListOperation, "stage/aws/foo", true, true},
		{logical.UpdateOperation, "stage/aws/foo", true, true},
		{logical.UpdateOperation, "stage/aws/policy/foo", true, true},

		{logical.DeleteOperation, "prod/foo", false, false},
		{logical.UpdateOperation, "prod/foo", false, false},
		{logical.ReadOperation, "prod/foo", true, false},
		{logical.ListOperation, "prod/foo", true, false},
		{logical.ReadOperation, "prod/aws/foo", false, false},

		{logical.ReadOperation, "foo/bar", true, true},
		{logical.ListOperation, "foo/bar", false, true},
		{logical.UpdateOperation, "foo/bar", false, true},
		{logical.CreateOperation, "foo/bar", true, true},

		{logical.ReadOperation, "baz/quux", true, false},
		{logical.CreateOperation, "baz/quux", true, false},
		{logical.PatchOperation, "baz/quux", true, false},
		{logical.ListOperation, "baz/quux", false, false},
		{logical.UpdateOperation, "baz/quux", false, false},

		// Path segment wildcards
		{logical.ReadOperation, "test/foo/bar/segment", false, false},
		{logical.ReadOperation, "test/foo/segment", true, false},
		{logical.ReadOperation, "test/bar/segment", true, false},
		{logical.ReadOperation, "test/segment/at/frond", false, false},
		{logical.ReadOperation, "test/segment/at/front", true, false},
		{logical.ReadOperation, "test/segment/at/end/foo", true, false},
		{logical.ReadOperation, "test/segment/at/end/foo/", false, false},
		{logical.ReadOperation, "test/segment/at/end/v2/foo/", true, false},
		{logical.ReadOperation, "test/segment/wildcard/at/foo/", true, false},
		{logical.ReadOperation, "test/segment/wildcard/at/end", true, false},
		{logical.ReadOperation, "test/segment/wildcard/at/end/", true, false},

		// Path segment wildcards vs glob
		{logical.ReadOperation, "1/2/3/4", false, false},
		{logical.ReadOperation, "1/2/3", true, false},
		{logical.UpdateOperation, "1/2/3", false, false},
		{logical.UpdateOperation, "1/2/3/4", true, false},
		{logical.CreateOperation, "1/2/3/4/5", true, false},
		{logical.RecoverOperation, "1/2/3/4", true, false},
	}

	for _, tc := range tcases {
		ctx := namespace.ContextWithNamespace(context.Background(), ns)
		request := new(logical.Request)
		request.Operation = tc.op
		request.Path = tc.path

		authResults := acl.AllowOperation(ctx, request, false)
		if authResults.Allowed != tc.allowed {
			t.Fatalf("bad: case %#v: %v, %v", tc, authResults.Allowed, authResults.RootPrivs)
		}
		if authResults.RootPrivs != tc.rootPrivs {
			t.Fatalf("bad: case %#v: %v, %v", tc, authResults.Allowed, authResults.RootPrivs)
		}
	}
}

func TestACL_Layered(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		t.Parallel()
		policy1, err := ParseACLPolicy(namespace.RootNamespace, aclPolicy)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		policy2, err := ParseACLPolicy(namespace.RootNamespace, aclPolicy2)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		acl, err := NewACL(namespace.RootContext(context.Background()), []*Policy{policy1, policy2})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		testLayeredACL(t, acl, namespace.RootNamespace)
	})
}

func testLayeredACL(t *testing.T, acl *ACL, ns *namespace.Namespace) {
	// Type of operation is not important here as we only care about checking
	// sudo/root
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	request := new(logical.Request)
	request.Operation = logical.ReadOperation
	request.Path = "sys/mount/foo"

	authResults := acl.AllowOperation(ctx, request, false)
	if authResults.RootPrivs {
		t.Fatalf("unexpected root")
	}

	type tcase struct {
		op        logical.Operation
		path      string
		allowed   bool
		rootPrivs bool
	}
	tcases := []tcase{
		{logical.ReadOperation, "root", false, false},
		{logical.HelpOperation, "root", true, false},

		{logical.ReadOperation, "dev/foo", true, true},
		{logical.UpdateOperation, "dev/foo", true, true},
		{logical.ReadOperation, "dev/hide/foo", false, false},
		{logical.UpdateOperation, "dev/hide/foo", false, false},

		{logical.DeleteOperation, "stage/foo", true, false},
		{logical.ListOperation, "stage/aws/foo", true, true},
		{logical.UpdateOperation, "stage/aws/foo", true, true},
		{logical.UpdateOperation, "stage/aws/policy/foo", false, false},

		{logical.DeleteOperation, "prod/foo", true, false},
		{logical.UpdateOperation, "prod/foo", true, false},
		{logical.ReadOperation, "prod/foo", true, false},
		{logical.ListOperation, "prod/foo", true, false},
		{logical.ReadOperation, "prod/aws/foo", false, false},

		{logical.ReadOperation, "sys/status", false, false},
		{logical.UpdateOperation, "sys/seal", true, true},

		{logical.ReadOperation, "foo/bar", false, false},
		{logical.ListOperation, "foo/bar", false, false},
		{logical.UpdateOperation, "foo/bar", false, false},
		{logical.CreateOperation, "foo/bar", false, false},

		{logical.ReadOperation, "baz/quux", false, false},
		{logical.ListOperation, "baz/quux", false, false},
		{logical.UpdateOperation, "baz/quux", false, false},
		{logical.CreateOperation, "baz/quux", false, false},
		{logical.PatchOperation, "baz/quux", false, false},
	}

	for _, tc := range tcases {
		ctx := namespace.ContextWithNamespace(context.Background(), ns)
		request := new(logical.Request)
		request.Operation = tc.op
		request.Path = tc.path

		authResults := acl.AllowOperation(ctx, request, false)
		if authResults.Allowed != tc.allowed {
			t.Fatalf("bad: case %#v: %v, %v", tc, authResults.Allowed, authResults.RootPrivs)
		}
		if authResults.RootPrivs != tc.rootPrivs {
			t.Fatalf("bad: case %#v: %v, %v", tc, authResults.Allowed, authResults.RootPrivs)
		}
	}
}

func TestACL_ParseMalformedPolicy(t *testing.T) {
	_, err := ParseACLPolicy(namespace.RootNamespace, `name{}`)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestACL_PolicyMerge(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		t.Parallel()
		testACLPolicyMerge(t, namespace.RootNamespace)
	})
}

func testACLPolicyMerge(t *testing.T, ns *namespace.Namespace) {
	policy, err := ParseACLPolicy(ns, mergingPolicies)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	acl, err := NewACL(ctx, []*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	type tcase struct {
		path           string
		minWrappingTTL *time.Duration
		maxWrappingTTL *time.Duration
		allowed        map[string][]interface{}
		denied         map[string][]interface{}
		required       []string
	}

	createDuration := func(seconds int) *time.Duration {
		ret := time.Duration(seconds) * time.Second
		return &ret
	}

	tcases := []tcase{
		{"foo/bar", nil, nil, nil, map[string][]interface{}{"zip": {}, "baz": {}}, []string{"baz"}},
		{"hello/universe", createDuration(50), createDuration(200), map[string][]interface{}{"foo": {}, "bar": {}}, nil, []string{"foo", "bar"}},
		{"allow/all", nil, nil, map[string][]interface{}{"*": {}, "test": {}, "test1": {"foo"}}, nil, nil},
		{"allow/all1", nil, nil, map[string][]interface{}{"*": {}, "test": {}, "test1": {"foo"}}, nil, nil},
		{"deny/all", nil, nil, nil, map[string][]interface{}{"*": {}, "test": {}}, nil},
		{"deny/all1", nil, nil, nil, map[string][]interface{}{"*": {}, "test": {}}, nil},
		{"value/merge", nil, nil, map[string][]interface{}{"test": {3, 4, 1, 2}}, map[string][]interface{}{"test": {3, 4, 1, 2}}, nil},
		{"value/empty", nil, nil, map[string][]interface{}{"empty": {}}, map[string][]interface{}{"empty": {}}, nil},
	}

	for _, tc := range tcases {
		policyPath := ns.Path + tc.path
		raw, ok := acl.exactRules.Get(policyPath)
		if !ok {
			t.Fatalf("Could not find acl entry for path %s", policyPath)
		}

		p := raw.(*ACLPermissions)
		if !reflect.DeepEqual(tc.allowed, p.AllowedParameters) {
			t.Fatalf("Allowed parameters did not match, Expected: %#v, Got: %#v", tc.allowed, p.AllowedParameters)
		}
		if !reflect.DeepEqual(tc.denied, p.DeniedParameters) {
			t.Fatalf("Denied parameters did not match, Expected: %#v, Got: %#v", tc.denied, p.DeniedParameters)
		}
		if !reflect.DeepEqual(tc.required, p.RequiredParameters) {
			t.Fatalf("Required parameters did not match, Expected: %#v, Got: %#v", tc.required, p.RequiredParameters)
		}
		if tc.minWrappingTTL != nil && *tc.minWrappingTTL != p.MinWrappingTTL {
			t.Fatalf("Min wrapping TTL did not match, Expected: %#v, Got: %#v", tc.minWrappingTTL, p.MinWrappingTTL)
		}
		if tc.minWrappingTTL != nil && *tc.maxWrappingTTL != p.MaxWrappingTTL {
			t.Fatalf("Max wrapping TTL did not match, Expected: %#v, Got: %#v", tc.maxWrappingTTL, p.MaxWrappingTTL)
		}
	}
}

func TestACL_AllowOperation(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		t.Parallel()
		testACLAllowOperation(t, namespace.RootNamespace)
	})
}

func testACLAllowOperation(t *testing.T, ns *namespace.Namespace) {
	policy, err := ParseACLPolicy(ns, permissionsPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	acl, err := NewACL(ctx, []*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	toperations := []logical.Operation{
		logical.UpdateOperation,
		logical.CreateOperation,
	}
	type tcase struct {
		path        string
		wrappingTTL *time.Duration
		parameters  []string
		allowed     bool
	}

	createDuration := func(seconds int) *time.Duration {
		ret := time.Duration(seconds) * time.Second
		return &ret
	}

	tcases := []tcase{
		{"dev/ops", nil, []string{"zip"}, true},
		{"foo/bar", nil, []string{"zap"}, false},
		{"foo/bar", nil, []string{"zip"}, false},
		{"foo/bar", createDuration(50), []string{"zip"}, false},
		{"foo/bar", createDuration(450), []string{"zip"}, false},
		{"foo/bar", createDuration(350), []string{"zip"}, true},
		{"foo/baz", nil, []string{"hello"}, false},
		{"foo/baz", createDuration(50), []string{"hello"}, false},
		{"foo/baz", createDuration(450), []string{"hello"}, true},
		{"foo/baz", nil, []string{"zap"}, false},
		{"broken/phone", nil, []string{"steve"}, false},
		{"working/phone", nil, []string{""}, false},
		{"working/phone", createDuration(450), []string{""}, false},
		{"working/phone", createDuration(350), []string{""}, true},
		{"hello/world", nil, []string{"one"}, false},
		{"tree/fort", nil, []string{"one"}, true},
		{"tree/fort", nil, []string{"foo"}, false},
		{"fruit/apple", nil, []string{"pear"}, false},
		{"fruit/apple", nil, []string{"one"}, false},
		{"cold/weather", nil, []string{"four"}, true},
		{"var/aws", nil, []string{"cold", "warm", "kitty"}, false},
		{"var/req", nil, []string{"cold", "warm", "kitty"}, false},
		{"var/req", nil, []string{"cold", "warm", "kitty", "foo"}, true},
	}

	for _, tc := range tcases {
		request := &logical.Request{
			Path: tc.path,
			Data: make(map[string]interface{}),
		}

		for _, parameter := range tc.parameters {
			request.Data[parameter] = ""
		}
		if tc.wrappingTTL != nil {
			request.WrapInfo = &logical.RequestWrapInfo{
				TTL: *tc.wrappingTTL,
			}
		}
		for _, op := range toperations {
			request.Operation = op
			ctx := namespace.ContextWithNamespace(context.Background(), ns)
			authResults := acl.AllowOperation(ctx, request, false)
			if authResults.Allowed != tc.allowed {
				t.Fatalf("bad: case %#v: %v", tc, authResults.Allowed)
			}
		}
	}
}

func TestACL_ValuePermissions(t *testing.T) {
	t.Run("root-ns", func(t *testing.T) {
		t.Parallel()
		testACLValuePermissions(t, namespace.RootNamespace)
	})
}

func testACLValuePermissions(t *testing.T, ns *namespace.Namespace) {
	policy, err := ParseACLPolicy(ns, valuePermissionsPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	acl, err := NewACL(ctx, []*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	toperations := []logical.Operation{
		logical.UpdateOperation,
		logical.CreateOperation,
	}
	type tcase struct {
		path       string
		parameters []string
		values     []interface{}
		allowed    bool
	}

	tcases := []tcase{
		{"dev/ops", []string{"allow"}, []interface{}{"good"}, true},
		{"dev/ops", []string{"allow"}, []interface{}{"bad"}, false},
		{"foo/bar", []string{"deny"}, []interface{}{"bad"}, false},
		{"foo/bar", []string{"deny"}, []interface{}{"bad glob"}, false},
		{"foo/bar", []string{"deny"}, []interface{}{"good"}, true},
		{"foo/bar", []string{"allow"}, []interface{}{"good"}, true},
		{"foo/bar", []string{"deny"}, []interface{}{nil}, true},
		{"foo/bar", []string{"allow"}, []interface{}{nil}, true},
		{"foo/baz", []string{"aLLow"}, []interface{}{"good"}, true},
		{"foo/baz", []string{"deny"}, []interface{}{"bad"}, false},
		{"foo/baz", []string{"deny"}, []interface{}{"good"}, false},
		{"foo/baz", []string{"allow", "deny"}, []interface{}{"good", "bad"}, false},
		{"foo/baz", []string{"deny", "allow"}, []interface{}{"good", "bad"}, false},
		{"foo/baz", []string{"deNy", "allow"}, []interface{}{"bad", "good"}, false},
		{"foo/baz", []string{"aLLow"}, []interface{}{"bad"}, false},
		{"foo/baz", []string{"Neither"}, []interface{}{"bad"}, false},
		{"foo/baz", []string{"allow"}, []interface{}{nil}, false},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"good"}, true},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"good1"}, true},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"good2"}, true},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"glob good2"}, false},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"glob good3"}, true},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"bad"}, false},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"bad"}, false},
		{"fizz/buzz", []string{"allow_multi", "allow"}, []interface{}{"good1", "good"}, true},
		{"fizz/buzz", []string{"deny_multi"}, []interface{}{"bad2"}, false},
		{"fizz/buzz", []string{"deny_multi", "allow_multi"}, []interface{}{"good", "good2"}, false},
		//	{"test/types", []string{"array"}, []interface{}{[1]string{"good"}}, true},
		{"test/types", []string{"map"}, []interface{}{map[string]interface{}{"good": "one"}}, true},
		{"test/types", []string{"map"}, []interface{}{map[string]interface{}{"bad": "one"}}, false},
		{"test/types", []string{"int"}, []interface{}{1}, true},
		{"test/types", []string{"int"}, []interface{}{3}, false},
		{"test/types", []string{"bool"}, []interface{}{false}, true},
		{"test/types", []string{"bool"}, []interface{}{true}, false},
		{"test/star", []string{"anything"}, []interface{}{true}, true},
		{"test/star", []string{"foo"}, []interface{}{true}, true},
		{"test/star", []string{"bar"}, []interface{}{false}, true},
		{"test/star", []string{"bar"}, []interface{}{true}, false},
	}

	for _, tc := range tcases {
		request := &logical.Request{
			Path: tc.path,
			Data: make(map[string]interface{}),
		}
		ctx := namespace.ContextWithNamespace(context.Background(), ns)

		for i, parameter := range tc.parameters {
			request.Data[parameter] = tc.values[i]
		}
		for _, op := range toperations {
			request.Operation = op
			authResults := acl.AllowOperation(ctx, request, false)
			if authResults.Allowed != tc.allowed {
				t.Fatalf("bad: case %#v: %v", tc, authResults.Allowed)
			}
		}
	}
}

func TestACL_SegmentWildcardPriority(t *testing.T) {
	ns := namespace.RootNamespace
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	type poltest struct {
		policy string
		path   string
	}

	// These test cases should each have a read rule and an update rule, where
	// the update rule wins out due to being more specific.
	poltests := []poltest{
		{
			// Verify edge conditions.  Here '*' is more specific both because
			// of first wildcard position (0 vs -1/infinity) and #wildcards.
			`
path "+/*" { capabilities = ["read"] }
path "*" { capabilities = ["update"] }
`,
			"foo/bar/bar/baz",
		},
		{
			// Verify edge conditions.  Here '+/*' is less specific because of
			// first wildcard position.
			`
path "+/*" { capabilities = ["read"] }
path "foo/+/*" { capabilities = ["update"] }
`,
			"foo/bar/bar/baz",
		},
		{
			// Verify that more wildcard segments is lower priority.
			`
path "foo/+/+/*" { capabilities = ["read"] }
path "foo/+/bar/baz" { capabilities = ["update"] }
`,
			"foo/bar/bar/baz",
		},
		{
			// Verify that more wildcard segments is lower priority.
			`
path "foo/+/+/baz" { capabilities = ["read"] }
path "foo/+/bar/baz" { capabilities = ["update"] }
`,
			"foo/bar/bar/baz",
		},
		{
			// Verify that first wildcard position is lower priority.
			// '(' is used here because it is lexicographically smaller than "+"
			`
path "foo/+/(ar/baz" { capabilities = ["read"] }
path "foo/(ar/+/baz" { capabilities = ["update"] }
`,
			"foo/(ar/(ar/baz",
		},

		{
			// Verify that a glob has lower priority, even if the prefix is the
			// same otherwise.
			`
path "foo/bar/+/baz*" { capabilities = ["read"] }
path "foo/bar/+/baz" { capabilities = ["update"] }
`,
			"foo/bar/bar/baz",
		},
		{
			// Verify that a shorter prefix has lower priority.
			`
path "foo/bar/+/b*" { capabilities = ["read"] }
path "foo/bar/+/ba*" { capabilities = ["update"] }
`,
			"foo/bar/bar/baz",
		},
	}

	for i, pt := range poltests {
		policy, err := ParseACLPolicy(ns, pt.policy)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		acl, err := NewACL(ctx, []*Policy{policy})
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		request := new(logical.Request)
		request.Path = pt.path

		request.Operation = logical.UpdateOperation
		authResults := acl.AllowOperation(ctx, request, false)
		if !authResults.Allowed {
			t.Fatalf("bad: case %d %#v: %v", i, pt, authResults.Allowed)
		}

		request.Operation = logical.ReadOperation
		authResults = acl.AllowOperation(ctx, request, false)
		if authResults.Allowed {
			t.Fatalf("bad: case %d %#v: %v", i, pt, authResults.Allowed)
		}
	}
}

func TestACL_SegmentWildcardPriority_BareMount(t *testing.T) {
	ns := namespace.RootNamespace
	ctx := namespace.ContextWithNamespace(context.Background(), ns)
	type poltest struct {
		policy    string
		mountpath string
		hasperms  bool
	}
	// These test cases should have one or more rules and a mount prefix.
	// hasperms should be true if there are non-deny perms that apply
	// to the mount prefix or something below it.
	poltests := []poltest{
		{
			`path "+" { capabilities = ["read"] }`,
			"foo/",
			true,
		},
		{
			`path "+/*" { capabilities = ["read"] }`,
			"foo/",
			true,
		},
		{
			`path "foo/+/+/*" { capabilities = ["read"] }`,
			"foo/",
			true,
		},
		{
			`path "foo/+/+/*" { capabilities = ["read"] }`,
			"foo/bar/",
			true,
		},
		{
			`path "foo/+/+/*" { capabilities = ["read"] }`,
			"foo/bar/bar/",
			true,
		},
		{
			`path "foo/+/+/*" { capabilities = ["read"] }`,
			"foo/bar/bar/baz/",
			true,
		},
		{
			`path "foo/+/+/baz" { capabilities = ["read"] }`,
			"foo/bar/bar/baz/",
			true,
		},
		{
			`path "foo/+/bar/baz" { capabilities = ["read"] }`,
			"foo/bar/bar/baz/",
			true,
		},
		{
			`path "foo/bar/+/baz*" { capabilities = ["read"] }`,
			"foo/bar/bar/baz/",
			true,
		},
		{
			`path "foo/bar/+/b*" { capabilities = ["read"] }`,
			"foo/bar/bar/baz/",
			true,
		},
		{
			`path "foo/+" { capabilities = ["read"] }`,
			"foo/",
			true,
		},
	}

	for i, pt := range poltests {
		policy, err := ParseACLPolicy(ns, pt.policy)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		acl, err := NewACL(ctx, []*Policy{policy})
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		hasperms := nil != acl.CheckAllowedFromNonExactPaths(pt.mountpath, true)
		if hasperms != pt.hasperms {
			t.Fatalf("bad: case %d: %#v", i, pt)
		}

	}
}

// NOTE: this test doesn't catch any races ATM
func TestACL_CreationRace(t *testing.T) {
	policy, err := ParseACLPolicy(namespace.RootNamespace, valuePermissionsPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	var wg sync.WaitGroup
	errs := make(chan error)
	stopTime := time.Now().Add(20 * time.Second)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				if time.Now().After(stopTime) {
					return
				}
				_, err := NewACL(namespace.RootContext(context.Background()), []*Policy{policy})
				if err != nil {
					errs <- fmt.Errorf("goroutine %d: %w", i, err)
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		t.Fatalf("err: %v", err)
	}
}

func TestACLGrantingPolicies(t *testing.T) {
	ns := namespace.RootNamespace
	policy, err := ParseACLPolicy(ns, grantingTestPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	merged, err := ParseACLPolicy(ns, grantingTestPolicyMerged)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	ctx := namespace.ContextWithNamespace(context.Background(), ns)

	type tcase struct {
		path     string
		op       logical.Operation
		policies []*Policy
		expected []logical.PolicyInfo
		allowed  bool
	}

	policyInfo := logical.PolicyInfo{
		Name:          "granting_policy",
		NamespaceId:   "root",
		NamespacePath: "",
		Type:          "acl",
	}
	mergedInfo := logical.PolicyInfo{
		Name:          "granting_policy_merged",
		NamespaceId:   "root",
		NamespacePath: "",
		Type:          "acl",
	}

	tcases := []tcase{
		{"kv/foo", logical.ReadOperation, []*Policy{policy}, []logical.PolicyInfo{policyInfo}, true},
		{"kv/foo", logical.UpdateOperation, []*Policy{policy}, []logical.PolicyInfo{policyInfo}, true},
		{"kv/bad", logical.ReadOperation, []*Policy{policy}, nil, false},
		{"kv/deny", logical.ReadOperation, []*Policy{policy}, nil, false},
		{"kv/path/foo", logical.ReadOperation, []*Policy{policy}, []logical.PolicyInfo{policyInfo}, true},
		{"kv/path/longer", logical.ReadOperation, []*Policy{policy}, []logical.PolicyInfo{policyInfo}, true},
		{"kv/foo", logical.ReadOperation, []*Policy{policy, merged}, []logical.PolicyInfo{policyInfo, mergedInfo}, true},
		{"kv/path/longer3", logical.ReadOperation, []*Policy{policy, merged}, []logical.PolicyInfo{mergedInfo}, true},
		{"kv/bar", logical.ReadOperation, []*Policy{policy, merged}, []logical.PolicyInfo{mergedInfo}, true},
		{"kv/deny", logical.ReadOperation, []*Policy{policy, merged}, nil, false},
		{"kv/path/longer", logical.UpdateOperation, []*Policy{policy, merged}, []logical.PolicyInfo{policyInfo}, true},
		{"kv/path/foo", logical.ReadOperation, []*Policy{policy, merged}, []logical.PolicyInfo{policyInfo, mergedInfo}, true},
		{"ns1/kv/foo", logical.RecoverOperation, []*Policy{policy, merged}, []logical.PolicyInfo{policyInfo}, true},
	}

	for _, tc := range tcases {
		request := &logical.Request{
			Path:      tc.path,
			Operation: tc.op,
		}

		acl, err := NewACL(ctx, tc.policies)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		authResults := acl.AllowOperation(ctx, request, false)
		if authResults.Allowed != tc.allowed {
			t.Fatalf("bad: case %#v: %v", tc, authResults.Allowed)
		}
		if !reflect.DeepEqual(authResults.GrantingPolicies, tc.expected) {
			t.Fatalf("bad: case %#v: got\n%#v\nexpected\n%#v\n", tc, authResults.GrantingPolicies, tc.expected)
		}
	}
}

var grantingTestPolicy = `
name = "granting_policy"
path "kv/foo" {
	capabilities = ["update", "read"]
}

path "kv/path/*" {
	capabilities = ["read"]
}

path "kv/path/longer" {
	capabilities = ["update", "read"]
}

path "kv/path/longer2" {
	capabilities = ["update"]
}

path "kv/deny" {
	capabilities = ["deny"]
}

path "ns1/kv/foo" {
	capabilities = ["update", "read", "recover"]
}
`

var grantingTestPolicyMerged = `
name = "granting_policy_merged"
path "kv/foo" {
	capabilities = ["update", "read"]
}

path "kv/bar" {
	capabilities = ["update", "read"]
}

path "kv/path/*" {
	capabilities = ["read"]
}

path "kv/path/longer" {
	capabilities = ["read"]
}

path "kv/path/longer3" {
	capabilities = ["read"]
}

path "kv/deny" {
	capabilities = ["update"]
}
`

var tokenCreationPolicy = `
name = "tokenCreation"
path "auth/token/create*" {
	capabilities = ["update", "create", "sudo"]
}
`

var aclPolicy = `
name = "DeV"
path "dev/*" {
	policy = "sudo"
}
path "stage/*" {
	policy = "write"
}
path "stage/aws/*" {
	policy = "read"
	capabilities = ["update", "sudo"]
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
path "foo/bar" {
	capabilities = ["read", "create", "sudo"]
}
path "baz/quux" {
	capabilities = ["read", "create", "patch"]
}
path "test/+/segment" {
	capabilities = ["read"]
}
path "+/segment/at/front" {
	capabilities = ["read"]
}
path "test/segment/at/end/+" {
	capabilities = ["read"]
}
path "test/segment/at/end/v2/+/" {
	capabilities = ["read"]
}
path "test/+/wildcard/+/*" {
	capabilities = ["read"]
}
path "test/+/wildcardglob/+/end*" {
	capabilities = ["read"]
}
path "1/2/*" {
	capabilities = ["create"]
}
path "1/2/+" {
	capabilities = ["read"]
}
path "1/2/+/+" {
	capabilities = ["update", "recover"]
}
`

var aclPolicy2 = `
name = "OpS"
path "dev/hide/*" {
	policy = "deny"
}
path "stage/aws/policy/*" {
	policy = "deny"
	# This should have no effect
	capabilities = ["read", "update", "sudo"]
}
path "prod/*" {
	policy = "write"
}
path "sys/seal" {
	policy = "sudo"
}
path "foo/bar" {
	capabilities = ["deny"]
}
path "baz/quux" {
	capabilities = ["deny"]
}
`

// test merging
var mergingPolicies = `
name = "ops"
path "foo/bar" {
	policy = "write"
	denied_parameters = {
		"baz" = []
	}
	required_parameters = ["baz"]
}
path "foo/bar" {
	policy = "write"
	denied_parameters = {
		"zip" = []
	}
}
path "hello/universe" {
	policy = "write"
	allowed_parameters = {
		"foo" = []
	}
	required_parameters = ["foo"]
	max_wrapping_ttl = 300
	min_wrapping_ttl = 100
}
path "hello/universe" {
	policy = "write"
	allowed_parameters = {
		"bar" = []
	}
	required_parameters = ["bar"]
	max_wrapping_ttl = 200
	min_wrapping_ttl = 50
}
path "allow/all" {
	policy = "write"
	allowed_parameters = {
		"test" = []
		"test1" = ["foo"]
	}
}
path "allow/all" {
	policy = "write"
	allowed_parameters = {
		"*" = []
	}
}
path "allow/all1" {
	policy = "write"
	allowed_parameters = {
		"*" = []
	}
}
path "allow/all1" {
	policy = "write"
	allowed_parameters = {
		"test" = []
		"test1" = ["foo"]
	}
}
path "deny/all" {
	policy = "write"
	denied_parameters = {
		"test" = []
	}
}
path "deny/all" {
	policy = "write"
	denied_parameters = {
		"*" = []
	}
}
path "deny/all1" {
	policy = "write"
	denied_parameters = {
		"*" = []
	}
}
path "deny/all1" {
	policy = "write"
	denied_parameters = {
		"test" = []
	}
}
path "value/merge" {
	policy = "write"
	allowed_parameters = {
		"test" = [1, 2]
	}
	denied_parameters = {
		"test" = [1, 2]
	}
}
path "value/merge" {
	policy = "write"
	allowed_parameters = {
		"test" = [3, 4]
	}
	denied_parameters = {
		"test" = [3, 4]
	}
}
path "value/empty" {
	policy = "write"
	allowed_parameters = {
		"empty" = []
	}
	denied_parameters = {
		"empty" = [1]
	}
}
path "value/empty" {
	policy = "write"
	allowed_parameters = {
		"empty" = [1]
	}
	denied_parameters = {
		"empty" = []
	}
}
`

// allow operation testing
var permissionsPolicy = `
name = "dev"
path "dev/*" {
	policy = "write"
	allowed_parameters = {
		"zip" = []
	}
}
path "foo/bar" {
	policy = "write"
	denied_parameters = {
		"zap" = []
	}
	min_wrapping_ttl = 300
	max_wrapping_ttl = 400
}
path "foo/baz" {
	policy = "write"
	allowed_parameters = {
		"hello" = []
	}
	denied_parameters = {
		"zap" = []
	}
	min_wrapping_ttl = 300
}
path "working/phone" {
	policy = "write"
	max_wrapping_ttl = 400
}
path "broken/phone" {
	policy = "write"
	allowed_parameters = {
	  "steve" = []
	}
	denied_parameters = {
	  "steve" = []
	}
}
path "hello/world" {
	policy = "write"
	allowed_parameters = {
		"*" = []
	}
	denied_parameters = {
		"*" = []
	}
}
path "tree/fort" {
	policy = "write"
	allowed_parameters = {
		"*" = []
	}
	denied_parameters = {
		"foo" = []
	}
}
path "fruit/apple" {
	policy = "write"
	allowed_parameters = {
		"pear" = []
	}
	denied_parameters = {
		"*" = []
	}
}
path "cold/weather" {
	policy = "write"
	allowed_parameters = {}
	denied_parameters = {}
}
path "var/aws" {
	policy = "write"
	allowed_parameters = {
		"*" = []
	}
	denied_parameters = {
		"soft" = []
		"warm" = []
		"kitty" = []
	}
}
path "var/req" {
	policy = "write"
	required_parameters = ["foo"]
}
`

// allow operation testing
var valuePermissionsPolicy = `
name = "op"
path "dev/*" {
	policy = "write"
	allowed_parameters = {
		"allow" = ["good"]
	}
}
path "foo/bar" {
	policy = "write"
	denied_parameters = {
		"deny" = ["bad*"]
	}
}
path "foo/baz" {
	policy = "write"
	allowed_parameters = {
		"ALLOW" = ["good"]
	}
	denied_parameters = {
		"dEny" = ["bad"]
	}
}
path "fizz/buzz" {
	policy = "write"
	allowed_parameters = {
		"allow_multi" = ["good", "good1", "good2", "*good3"]
		"allow" = ["good"]
	}
	denied_parameters = {
		"deny_multi" = ["bad", "bad1", "bad2"]
	}
}
path "test/types" {
	policy = "write"
	allowed_parameters = {
		"map" = [{"good" = "one"}]
		"int" = [1, 2]
		"bool" = [false]
	}
	denied_parameters = {
	}
}
path "test/star" {
	policy = "write"
	allowed_parameters = {
		"*" = []
		"foo" = []
		"bar" = [false]
	}
	denied_parameters = {
	}
}
`
