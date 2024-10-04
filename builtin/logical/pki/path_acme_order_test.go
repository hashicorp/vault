// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"net"
	"testing"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestACME_ValidateIdentifiersAgainstRole Verify the ACME order creation
// function verifies somewhat the identifiers that were provided have a
// decent chance of being allowed by the selected role.
func TestACME_ValidateIdentifiersAgainstRole(t *testing.T) {
	b, _ := CreateBackendWithStorage(t)

	tests := []struct {
		name        string
		role        *issuing.RoleEntry
		identifiers []*ACMEIdentifier
		expectErr   bool
	}{
		{
			name:        "verbatim-role-allows-dns-ip",
			role:        issuing.SignVerbatimRole(),
			identifiers: _buildACMEIdentifiers("test.com", "127.0.0.1"),
			expectErr:   false,
		},
		{
			name:        "default-role-does-not-allow-dns",
			role:        buildTestRole(t, nil),
			identifiers: _buildACMEIdentifiers("www.test.com"),
			expectErr:   true,
		},
		{
			name:        "default-role-allows-ip",
			role:        buildTestRole(t, nil),
			identifiers: _buildACMEIdentifiers("192.168.0.1"),
			expectErr:   false,
		},
		{
			name:        "disable-ip-sans-forbids-ip",
			role:        buildTestRole(t, map[string]interface{}{"allow_ip_sans": false}),
			identifiers: _buildACMEIdentifiers("192.168.0.1"),
			expectErr:   true,
		},
		{
			name: "role-no-wildcards-allowed-without",
			role: buildTestRole(t, map[string]interface{}{
				"allow_subdomains":            true,
				"allow_bare_domains":          true,
				"allowed_domains":             []string{"test.com"},
				"allow_wildcard_certificates": false,
			}),
			identifiers: _buildACMEIdentifiers("www.test.com", "test.com"),
			expectErr:   false,
		},
		{
			name: "role-no-wildcards-allowed-with-wildcard",
			role: buildTestRole(t, map[string]interface{}{
				"allow_subdomains":            true,
				"allowed_domains":             []string{"test.com"},
				"allow_wildcard_certificates": false,
			}),
			identifiers: _buildACMEIdentifiers("*.test.com"),
			expectErr:   true,
		},
		{
			name: "role-wildcards-allowed-with-wildcard",
			role: buildTestRole(t, map[string]interface{}{
				"allow_subdomains":            true,
				"allowed_domains":             []string{"test.com"},
				"allow_wildcard_certificates": true,
			}),
			identifiers: _buildACMEIdentifiers("*.test.com"),
			expectErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := b.validateIdentifiersAgainstRole(tt.role, tt.identifiers)

			if tt.expectErr {
				require.Error(t, err, "validateIdentifiersAgainstRole(%v, %v)", tt.role.ToResponseData(), tt.identifiers)
				// If we did return an error if should be classified as a ErrRejectedIdentifier
				require.ErrorIs(t, err, ErrRejectedIdentifier)
			} else {
				require.NoError(t, err, "validateIdentifiersAgainstRole(%v, %v)", tt.role.ToResponseData(), tt.identifiers)
			}
		})
	}
}

func _buildACMEIdentifiers(values ...string) []*ACMEIdentifier {
	var identifiers []*ACMEIdentifier

	for _, value := range values {
		identifiers = append(identifiers, _buildACMEIdentifier(value))
	}

	return identifiers
}

func _buildACMEIdentifier(val string) *ACMEIdentifier {
	ip := net.ParseIP(val)
	if ip == nil {
		identifier := &ACMEIdentifier{Type: "dns", Value: val, OriginalValue: val, IsWildcard: false}
		_, _, _ = identifier.MaybeParseWildcard()
		return identifier
	}

	return &ACMEIdentifier{Type: "ip", Value: val, OriginalValue: val, IsWildcard: false}
}

// Easily allow tests to create valid roles with proper defaults, since we don't have an easy
// way to generate roles with proper defaults, go through the createRole handler with the handlers
// field data so we pickup all the defaults specified there.
func buildTestRole(t *testing.T, config map[string]interface{}) *issuing.RoleEntry {
	b, s := CreateBackendWithStorage(t)

	path := pathRoles(b)
	fields := path.Fields
	if config == nil {
		config = map[string]interface{}{}
	}

	if _, exists := config["name"]; !exists {
		config["name"] = genUuid()
	}

	_, err := b.pathRoleCreate(ctx, &logical.Request{Storage: s}, &framework.FieldData{Raw: config, Schema: fields})
	require.NoError(t, err, "failed generating role with config %v", config)

	role, err := b.GetRole(ctx, s, config["name"].(string))
	require.NoError(t, err, "failed loading stored role")

	return role
}
