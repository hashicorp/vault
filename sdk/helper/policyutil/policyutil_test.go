// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package policyutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestValidatePolicyName verifies policy-name canonicalization and segment-aware
// traversal validation.
func TestValidatePolicyName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "accepts valid name",
			input: " Team-Read ",
			want:  "team-read",
		},
		{
			name:  "accepts embedded dot dot substring",
			input: "policy..backup",
			want:  "policy..backup",
		},
		{
			name:  "accepts embedded triple dots",
			input: "version...read",
			want:  "version...read",
		},
		{
			name:    "rejects empty canonical name",
			input:   " \t ",
			wantErr: true,
		},
		{
			name:  "accepts slash-separated name",
			input: "team/read",
			want:  "team/read",
		},
		{
			name:  "accepts backslash-containing name",
			input: `team\read`,
			want:  `team\read`,
		},
		{
			name:  "accepts dotted suffix",
			input: "read.v2",
			want:  "read.v2",
		},
		{
			name:  "accepts hidden-style segment",
			input: ".hidden",
			want:  ".hidden",
		},
		{
			name:  "accepts hidden subsegment",
			input: "team/.hidden",
			want:  "team/.hidden",
		},
		{
			name:  "accepts prefixed dotdot subsegment",
			input: "team/..hidden",
			want:  "team/..hidden",
		},
		{
			name:  "accepts suffixed dotdot subsegment",
			input: "team/hidden..",
			want:  "team/hidden..",
		},
		{
			name:    "rejects parent segment name",
			input:   "..",
			wantErr: true,
		},
		{
			name:    "rejects parent segment at beginning",
			input:   "../admin",
			wantErr: true,
		},
		{
			name:    "rejects parent segment in middle",
			input:   "team/../admin",
			wantErr: true,
		},
		{
			name:    "rejects parent segment at end",
			input:   "team/..",
			wantErr: true,
		},
		{
			name:    "rejects multi-level parent traversal",
			input:   "../root/admin",
			wantErr: true,
		},
		{
			name:    "rejects exact dot segment name",
			input:   ".",
			wantErr: true,
		},
		{
			name:    "rejects dot segment at beginning",
			input:   "./admin",
			wantErr: true,
		},
		{
			name:    "rejects dot segment in middle",
			input:   "team/./admin",
			wantErr: true,
		},
		{
			name:    "rejects dot segment at end",
			input:   "team/.",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ValidatePolicyName(tc.input)
			if tc.wantErr {
				require.Error(t, err, "expected validation error for %q", tc.input)
				return
			}

			require.NoError(t, err, "expected no validation error for %q", tc.input)
			require.Equal(t, tc.want, got, "expected canonical policy name")
		})
	}
}

// TestSanitizePolicies verifies policy sanitization preserves valid values and
// default-policy handling semantics.
func TestSanitizePolicies(t *testing.T) {
	t.Parallel()

	expected := []string{"foo", "bar"}
	actual := SanitizePolicies([]string{"foo", "bar"}, false)
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)

	// If 'default' is already added, do not remove it.
	expected = []string{"foo", "bar", "default"}
	actual = SanitizePolicies([]string{"foo", "bar", "default"}, false)
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)
}

// TestParsePolicies verifies parsing normalizes policy lists and root-policy
// precedence semantics.
func TestParsePolicies(t *testing.T) {
	t.Parallel()

	expected := []string{"foo", "bar", "default"}
	actual := ParsePolicies("foo,bar")
	// add default if not present.
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)

	// do not add default more than once.
	actual = ParsePolicies("foo,bar,default")
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)

	// handle spaces and tabs.
	actual = ParsePolicies(" foo ,	bar	,   default")
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)

	// ignore all others if root is present.
	expected = []string{"root"}
	actual = ParsePolicies("foo,bar,root")
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)

	// with spaces and tabs.
	expected = []string{"root"}
	actual = ParsePolicies("foo ,bar, root		")
	require.True(t, EquivalentPolicies(expected, actual), "bad: expected:%s\ngot:%s\n", expected, actual)
}

// TestEquivalentPolicies verifies semantic policy-list equivalence, including
// case-insensitivity and implicit default handling.
func TestEquivalentPolicies(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		A        []string
		B        []string
		Expected bool
	}{
		"nil": {
			A:        nil,
			B:        nil,
			Expected: true,
		},
		"empty": {
			A:        []string{"foo", "bar"},
			B:        []string{},
			Expected: false,
		},
		"missing": {
			A:        []string{"foo", "bar"},
			B:        []string{"foo"},
			Expected: false,
		},
		"equal": {
			A:        []string{"bar", "foo"},
			B:        []string{"bar", "foo"},
			Expected: true,
		},
		"default": {
			A:        []string{"bar", "foo"},
			B:        []string{"foo", "default", "bar"},
			Expected: true,
		},
		"case-insensitive": {
			A:        []string{"test"},
			B:        []string{"Test"},
			Expected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.Expected, EquivalentPolicies(tc.A, tc.B))
		})
	}
}

// TestCanonicalizePolicyName verifies the single canonicalization function used
// by both ACL parameter enforcement and the downstream consumers that store
// policy names. Policy names are matched case-insensitively and ignoring
// surrounding whitespace, so representations that differ only in case or
// surrounding whitespace must collapse to the same canonical name. If they did
// not, a caller could submit a representation that evades a denied_parameters
// ACL check but still resolves to the denied policy when Vault stores it.
func TestCanonicalizePolicyName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    string
		expected string
	}{
		"already canonical":            {"super-admin", "super-admin"},
		"mixed case":                   {"Super-Admin", "super-admin"},
		"upper case":                   {"SUPER-ADMIN", "super-admin"},
		"leading space":                {" super-admin", "super-admin"},
		"trailing space":               {"super-admin ", "super-admin"},
		"surrounding spaces":           {"  super-admin  ", "super-admin"},
		"surrounding tabs":             {"\tsuper-admin\t", "super-admin"},
		"surrounding newlines":         {"\nsuper-admin\n", "super-admin"},
		"carriage return and newline":  {"\r\nsuper-admin\r\n", "super-admin"},
		"vertical tab and form feed":   {"\vsuper-admin\f", "super-admin"},
		"unicode no-break space":       {"\u00a0super-admin\u00a0", "super-admin"},
		"unicode line separator":       {"\u2028super-admin\u2029", "super-admin"},
		"unicode ideographic space":    {"\u3000super-admin\u3000", "super-admin"},
		"whitespace and case combined": {"  Super-Admin\t", "super-admin"},
		"empty string":                 {"", ""},
		"only whitespace":              {"   \t\n", ""},
		// Interior whitespace is meaningful: it denotes a different policy name
		// and must not be collapsed away.
		"interior whitespace preserved": {" super admin ", "super admin"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, CanonicalizePolicyName(tc.input))
			// Canonicalization must be idempotent, otherwise ACL evaluation and
			// downstream storage could still disagree after one pass.
			require.Equal(t, tc.expected, CanonicalizePolicyName(tc.expected))
		})
	}
}

// TestCanonicalizePolicyName_MatchesSanitizePolicies verifies that
// CanonicalizePolicyName agrees with SanitizePolicies, which is the function
// downstream consumers use. ACL enforcement relies on these producing identical
// results; any divergence reintroduces the canonicalization bypass.
func TestCanonicalizePolicyName_MatchesSanitizePolicies(t *testing.T) {
	t.Parallel()

	inputs := []string{
		"super-admin", "Super-Admin", "SUPER-ADMIN", " super-admin ",
		"\tsuper-admin\n", "\u00a0Super-Admin\u2028", "  read-only  ",
	}

	for _, input := range inputs {
		sanitized := SanitizePolicies([]string{input}, DoNotAddDefaultPolicy)
		require.Len(t, sanitized, 1)
		require.Equal(t, sanitized[0], CanonicalizePolicyName(input),
			"CanonicalizePolicyName diverged from SanitizePolicies for %q", input)
	}
}
