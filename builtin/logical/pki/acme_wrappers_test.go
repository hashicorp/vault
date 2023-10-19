// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestACMEIssuerRoleLoading validates the role and issuer loading logic within the base
// ACME wrapper is correct.
func TestACMEIssuerRoleLoading(t *testing.T) {
	b, s := CreateBackendWithStorage(t)

	_, err := CBWrite(b, s, "config/cluster", map[string]interface{}{
		"path":     "http://localhost:8200/v1/pki",
		"aia_path": "http://localhost:8200/cdn/pki",
	})
	require.NoError(t, err)

	_, err = CBWrite(b, s, "config/acme", map[string]interface{}{
		"enabled": true,
	})
	require.NoError(t, err)

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "myvault1.com",
		"issuer_name": "issuer-1",
		"key_type":    "ec",
	})
	require.NoError(t, err, "failed creating issuer issuer-1")

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "myvault2.com",
		"issuer_name": "issuer-2",
		"key_type":    "ec",
	})
	require.NoError(t, err, "failed creating issuer issuer-2")

	_, err = CBWrite(b, s, "roles/role-bad-issuer", map[string]interface{}{
		issuerRefParam: "non-existant",
		"no_store":     "false",
	})
	require.NoError(t, err, "failed creating role role-bad-issuer")

	_, err = CBWrite(b, s, "roles/role-no-store-enabled", map[string]interface{}{
		issuerRefParam: "issuer-2",
		"no_store":     "true",
	})
	require.NoError(t, err, "failed creating role role-no-store-enabled")

	_, err = CBWrite(b, s, "roles/role-issuer-2", map[string]interface{}{
		issuerRefParam: "issuer-2",
		"no_store":     "false",
	})
	require.NoError(t, err, "failed creating role role-issuer-2")

	tc := []struct {
		name               string
		roleName           string
		issuerName         string
		expectedIssuerName string
		expectErr          bool
	}{
		{name: "pass-default-use-default", roleName: "", issuerName: "", expectedIssuerName: "issuer-1", expectErr: false},
		{name: "pass-role-issuer-2", roleName: "role-issuer-2", issuerName: "", expectedIssuerName: "issuer-2", expectErr: false},
		{name: "pass-issuer-1-no-role", roleName: "", issuerName: "issuer-1", expectedIssuerName: "issuer-1", expectErr: false},
		{name: "fail-role-has-bad-issuer", roleName: "role-bad-issuer", issuerName: "", expectedIssuerName: "", expectErr: true},
		{name: "fail-role-no-store-enabled", roleName: "role-no-store-enabled", issuerName: "", expectedIssuerName: "", expectErr: true},
		{name: "fail-role-no-store-enabled", roleName: "role-no-store-enabled", issuerName: "", expectedIssuerName: "", expectErr: true},
		{name: "fail-role-does-not-exist", roleName: "non-existant", issuerName: "", expectedIssuerName: "", expectErr: true},
		{name: "fail-issuer-does-not-exist", roleName: "", issuerName: "non-existant", expectedIssuerName: "", expectErr: true},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			f := b.acmeWrapper(acmeWrapperOpts{}, func(acmeCtx *acmeContext, r *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
				if tt.roleName != acmeCtx.role.Name {
					return nil, fmt.Errorf("expected role %s but got %s", tt.roleName, acmeCtx.role.Name)
				}

				if tt.expectedIssuerName != acmeCtx.issuer.Name {
					return nil, fmt.Errorf("expected issuer %s but got %s", tt.expectedIssuerName, acmeCtx.issuer.Name)
				}

				return nil, nil
			})

			var acmePath string
			fieldRaw := map[string]interface{}{}
			if tt.issuerName != "" {
				fieldRaw[issuerRefParam] = tt.issuerName
				acmePath = "issuer/" + tt.issuerName + "/"
			}
			if tt.roleName != "" {
				fieldRaw["role"] = tt.roleName
				acmePath = acmePath + "roles/" + tt.roleName + "/"
			}

			acmePath = strings.TrimLeft(acmePath+"/acme/directory", "/")

			resp, err := f(context.Background(), &logical.Request{Path: acmePath, Storage: s}, &framework.FieldData{
				Raw:    fieldRaw,
				Schema: getCsrSignVerbatimSchemaFields(),
			})
			require.NoError(t, err, "all errors should be re-encoded")

			if tt.expectErr {
				require.NotEqual(t, 200, resp.Data[logical.HTTPStatusCode])
				require.Equal(t, ErrorContentType, resp.Data[logical.HTTPContentType])
			} else {
				if resp != nil {
					t.Fatalf("expected no error got %s", string(resp.Data[logical.HTTPRawBody].([]uint8)))
				}
			}
		})
	}
}
