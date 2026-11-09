// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package exactdeny

import (
	"testing"

	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestACL_DenyAppliesToCaseVariants verifies that an explicit deny also applies
// to case variants of its path, but only when it is at least as specific as the
// rule that matched the request as written. A less specific deny must never
// override a more specific allow, and LIST requests with a trailing slash must
// be evaluated the same way as the equivalent READ.
func TestACL_DenyAppliesToCaseVariants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		policy  string
		op      logical.Operation
		path    string
		allowed bool
	}{
		{
			name: "glob allow, exact deny, case variant",
			policy: `
path "secret/*"   { capabilities = ["read"] }
path "secret/foo" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/FOO",
			allowed: false,
		},
		{
			name: "glob allow, exact deny, canonical path",
			policy: `
path "secret/*"   { capabilities = ["read"] }
path "secret/foo" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/foo",
			allowed: false,
		},
		{
			name: "glob allow, exact deny, unrelated path",
			policy: `
path "secret/*"   { capabilities = ["read"] }
path "secret/foo" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/bar",
			allowed: true,
		},
		{
			name: "segment wildcard allow, exact deny, case variant",
			policy: `
path "secret/+/data" { capabilities = ["read"] }
path "secret/foo/data" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/FOO/data",
			allowed: false,
		},
		{
			name: "exact allow, less specific glob deny, case variant of allow",
			policy: `
path "secret/Foo" { capabilities = ["read"] }
path "secret/*"   { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/Foo",
			allowed: true,
		},
		{
			name: "exact allow, equally specific exact deny differing by case",
			policy: `
path "secret/Foo" { capabilities = ["read"] }
path "secret/foo" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/Foo",
			allowed: false,
		},
		{
			name: "exact list allow, equally specific exact deny, trailing slash",
			policy: `
path "secret/Foo" { capabilities = ["list"] }
path "secret/foo" { capabilities = ["deny"] }`,
			op:      logical.ListOperation,
			path:    "secret/Foo/",
			allowed: false,
		},
		{
			name: "glob list allow, exact deny, trailing slash",
			policy: `
path "secret/*"   { capabilities = ["list"] }
path "secret/foo" { capabilities = ["deny"] }`,
			op:      logical.ListOperation,
			path:    "secret/FOO/",
			allowed: false,
		},
		{
			name: "glob list allow, prefix deny, trailing slash",
			policy: `
path "secret/*"      { capabilities = ["list"] }
path "secret/foo/*"  { capabilities = ["deny"] }`,
			op:      logical.ListOperation,
			path:    "secret/FOO/",
			allowed: false,
		},
		{
			name: "more specific list allow beats less specific deny, trailing slash",
			policy: `
path "secret/FOO/*" { capabilities = ["list"] }
path "secret/*"     { capabilities = ["deny"] }`,
			op:      logical.ListOperation,
			path:    "secret/FOO/",
			allowed: true,
		},
		{
			name: "no deny rules",
			policy: `
path "secret/*" { capabilities = ["read"] }`,
			op:      logical.ReadOperation,
			path:    "secret/FOO",
			allowed: true,
		},
		{
			name: "deny rule with uppercase letters applies to lowercase request",
			policy: `
path "secret/*"   { capabilities = ["read"] }
path "secret/Foo" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/foo",
			allowed: false,
		},
		{
			name: "glob allow, prefix deny, nested case variant",
			policy: `
path "secret/*"     { capabilities = ["read"] }
path "secret/foo/*" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/FOO/bar",
			allowed: false,
		},
		{
			name: "glob allow, prefix deny, variant in deeper segments",
			policy: `
path "secret/*"     { capabilities = ["read"] }
path "secret/foo/*" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/Foo/BAR/Baz",
			allowed: false,
		},
		{
			name: "glob allow, prefix deny, path above the denied prefix stays allowed",
			policy: `
path "secret/*"     { capabilities = ["read"] }
path "secret/foo/*" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/FOO",
			allowed: true,
		},
		{
			name: "glob allow, prefix deny, sibling with shared stem stays allowed",
			policy: `
path "secret/*"     { capabilities = ["read"] }
path "secret/foo/*" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/FOOBAR/baz",
			allowed: true,
		},
		{
			name: "glob allow, mid-segment glob deny, case variant",
			policy: `
path "secret/*"    { capabilities = ["read"] }
path "secret/foo*" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/FOOBAR",
			allowed: false,
		},
		{
			name: "glob allow, mid-segment glob deny, non-matching stem",
			policy: `
path "secret/*"    { capabilities = ["read"] }
path "secret/foo*" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/FO",
			allowed: true,
		},
		{
			name: "exact allow beats less specific prefix deny, case variant of allow",
			policy: `
path "secret/*"         { capabilities = ["read"] }
path "secret/FOO/bar"   { capabilities = ["read"] }
path "secret/foo/*"     { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/FOO/bar",
			allowed: true,
		},
		{
			name: "exact allow and equally specific exact deny differing by case, allow listed in lowercase",
			policy: `
path "secret/foo" { capabilities = ["read"] }
path "secret/FOO" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/foo",
			allowed: false,
		},
		{
			name: "mixed-case allow prefix, lowercase exact deny, uppercase request",
			policy: `
path "secret/Foo/*"   { capabilities = ["read"] }
path "secret/foo/bar" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/Foo/BAR",
			allowed: false,
		},
		{
			name: "segment wildcard allow, segment wildcard deny, case variant",
			policy: `
path "secret/*"       { capabilities = ["read"] }
path "secret/+/data"  { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/Anything/DATA",
			allowed: false,
		},
		{
			name: "segment wildcard deny does not match other suffix",
			policy: `
path "secret/*"       { capabilities = ["read"] }
path "secret/+/data"  { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/Anything/OTHER",
			allowed: true,
		},
		{
			name: "segment wildcard deny with uppercase in rule applies to lowercase request",
			policy: `
path "secret/*"       { capabilities = ["read"] }
path "secret/+/Data"  { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/x/data",
			allowed: false,
		},
		{
			name: "segment wildcard allow, less specific segment wildcard deny",
			policy: `
path "secret/+/data/foo" { capabilities = ["read"] }
path "secret/+/+/Foo"    { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/X/data/foo",
			allowed: true,
		},
		{
			name: "segment wildcard allow, more specific exact deny",
			policy: `
path "secret/+/data/+"       { capabilities = ["read"] }
path "secret/team/data/foo"  { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/TEAM/Data/FOO",
			allowed: false,
		},
		{
			name: "deny blocks update",
			policy: `
path "secret/*"   { capabilities = ["read", "update"] }
path "secret/foo" { capabilities = ["deny"] }`,
			op:      logical.UpdateOperation,
			path:    "secret/FOO",
			allowed: false,
		},
		{
			name: "deny blocks delete",
			policy: `
path "secret/*"   { capabilities = ["read", "delete"] }
path "secret/foo" { capabilities = ["deny"] }`,
			op:      logical.DeleteOperation,
			path:    "secret/Foo",
			allowed: false,
		},
		{
			name: "deny blocks create",
			policy: `
path "secret/*"   { capabilities = ["create"] }
path "secret/foo" { capabilities = ["deny"] }`,
			op:      logical.CreateOperation,
			path:    "secret/fOO",
			allowed: false,
		},
		{
			name: "deny blocks patch",
			policy: `
path "secret/*"   { capabilities = ["patch"] }
path "secret/foo" { capabilities = ["deny"] }`,
			op:      logical.PatchOperation,
			path:    "secret/FOO",
			allowed: false,
		},
		{
			name: "non-deny rule on the lowercase path does not block case variants",
			policy: `
path "secret/*"   { capabilities = ["read"] }
path "secret/foo" { capabilities = ["list"] }`,
			op:      logical.ReadOperation,
			path:    "secret/FOO",
			allowed: true,
		},
		{
			name: "deny with a leading slash in the request",
			policy: `
path "secret/*"   { capabilities = ["read"] }
path "secret/foo" { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "/secret/FOO",
			allowed: false,
		},
		{
			name: "non-ASCII case variant of a denied path",
			policy: `
path "secret/*"       { capabilities = ["read"] }
path "secret/über"    { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/ÜBER",
			allowed: false,
		},
		{
			// Lowercasing İ (2 bytes) yields i + combining dot (3 bytes), so the
			// wildcard offsets of both rules must be recomputed after folding for
			// the equal-specificity comparison to stay correct.
			name: "non-ASCII allow is not overridden by less specific deny despite length change",
			policy: `
path "secret/İ/+/data" { capabilities = ["read"] }
path "secret/İ/+/+"    { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/İ/x/data",
			allowed: true,
		},
		{
			name: "non-ASCII deny with length-changing lowercase applies to variant",
			policy: `
path "secret/*"          { capabilities = ["read"] }
path "secret/İstanbul"   { capabilities = ["deny"] }`,
			op:      logical.ReadOperation,
			path:    "secret/İSTANBUL",
			allowed: false,
		},
		{
			name: "exact allow nested in prefix deny, canonical path stays allowed",
			policy: `
path "secret/*"         { capabilities = ["read"] }
path "secret/team/*"    { capabilities = ["deny"] }
path "secret/team/ok"   { capabilities = ["read"] }`,
			op:      logical.ReadOperation,
			path:    "secret/team/ok",
			allowed: true,
		},
		{
			// The request as written only matches the broad glob, so the nested
			// exact allow on the lowercase path does not protect the variant.
			name: "exact allow nested in prefix deny, case variant is denied",
			policy: `
path "secret/*"         { capabilities = ["read"] }
path "secret/team/*"    { capabilities = ["deny"] }
path "secret/team/ok"   { capabilities = ["read"] }`,
			op:      logical.ReadOperation,
			path:    "secret/TEAM/ok",
			allowed: false,
		},
		{
			name: "exact allow nested in prefix deny, case variant of sibling is denied",
			policy: `
path "secret/*"         { capabilities = ["read"] }
path "secret/team/*"    { capabilities = ["deny"] }
path "secret/team/ok"   { capabilities = ["read"] }`,
			op:      logical.ReadOperation,
			path:    "secret/TEAM/other",
			allowed: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := namespace.ContextWithNamespace(t.Context(), namespace.RootNamespace)
			policy, err := vault.ParseACLPolicy(namespace.RootNamespace, tc.policy)
			require.NoError(t, err)
			acl, err := vault.NewACL(ctx, []*vault.Policy{policy})
			require.NoError(t, err)

			result := acl.AllowOperation(ctx, &logical.Request{Operation: tc.op, Path: tc.path}, false)
			require.Equal(t, tc.allowed, result.Allowed)
		})
	}
}

// TestACL_DenyAppliesToCaseVariants_MultiplePolicies verifies that case-folded
// deny handling works when the allow and the deny come from different
// policies, and that the outcome does not depend on policy order.
func TestACL_DenyAppliesToCaseVariants_MultiplePolicies(t *testing.T) {
	t.Parallel()

	const (
		allowGlob     = `path "secret/*" { capabilities = ["read"] }`
		denyExact     = `path "secret/foo" { capabilities = ["deny"] }`
		allowExact    = `path "secret/Foo" { capabilities = ["read"] }`
		denyGlob      = `path "secret/*" { capabilities = ["deny"] }`
		allowExactLow = `path "secret/foo" { capabilities = ["read"] }`
	)

	tests := []struct {
		name     string
		policies []string
		path     string
		allowed  bool
	}{
		{
			name:     "deny in a separate policy applies to case variant",
			policies: []string{allowGlob, denyExact},
			path:     "secret/FOO",
			allowed:  false,
		},
		{
			name:     "deny in a separate policy applies regardless of policy order",
			policies: []string{denyExact, allowGlob},
			path:     "secret/FOO",
			allowed:  false,
		},
		{
			name:     "more specific allow in another policy beats less specific deny",
			policies: []string{allowExact, denyGlob},
			path:     "secret/Foo",
			allowed:  true,
		},
		{
			name:     "more specific allow beats less specific deny regardless of policy order",
			policies: []string{denyGlob, allowExact},
			path:     "secret/Foo",
			allowed:  true,
		},
		{
			name:     "same lowercase path allowed in one policy and denied in another",
			policies: []string{allowGlob, allowExactLow, denyExact},
			path:     "secret/Foo",
			allowed:  false,
		},
		{
			name:     "unrelated policy without deny rules does not affect the result",
			policies: []string{allowGlob, denyExact, `path "other/*" { capabilities = ["read"] }`},
			path:     "secret/bar",
			allowed:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := namespace.ContextWithNamespace(t.Context(), namespace.RootNamespace)
			var policies []*vault.Policy
			for _, p := range tc.policies {
				policy, err := vault.ParseACLPolicy(namespace.RootNamespace, p)
				require.NoError(t, err)
				policies = append(policies, policy)
			}
			acl, err := vault.NewACL(ctx, policies)
			require.NoError(t, err)

			result := acl.AllowOperation(ctx, &logical.Request{Operation: logical.ReadOperation, Path: tc.path}, false)
			require.Equal(t, tc.allowed, result.Allowed)
		})
	}
}

// TestACL_DenyAppliesToCaseVariants_Capabilities verifies that the
// capabilities-check path (used by sys/capabilities*) reports "deny" for case
// variants of a denied path, so it agrees with request enforcement.
func TestACL_DenyAppliesToCaseVariants_Capabilities(t *testing.T) {
	t.Parallel()

	ctx := namespace.ContextWithNamespace(t.Context(), namespace.RootNamespace)
	policy, err := vault.ParseACLPolicy(namespace.RootNamespace, `
path "secret/*"      { capabilities = ["read", "list"] }
path "secret/foo"    { capabilities = ["deny"] }
path "secret/team/*" { capabilities = ["deny"] }`)
	require.NoError(t, err)
	acl, err := vault.NewACL(ctx, []*vault.Policy{policy})
	require.NoError(t, err)

	tests := []struct {
		path string
		want []string
	}{
		{path: "secret/foo", want: []string{"deny"}},
		{path: "secret/FOO", want: []string{"deny"}},
		{path: "secret/FoO", want: []string{"deny"}},
		{path: "secret/TEAM/x", want: []string{"deny"}},
		{path: "secret/bar", want: []string{"read", "list"}},
		{path: "secret/BAR", want: []string{"read", "list"}},
	}
	for _, tc := range tests {
		require.Equal(t, tc.want, acl.Capabilities(ctx, tc.path), "path %q", tc.path)
	}
}
