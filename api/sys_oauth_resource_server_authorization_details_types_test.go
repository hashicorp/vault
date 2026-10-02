// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSys_OAuthResourceServerAuthorizationDetailsTypes exercises API client
// methods for authorization_details discovery endpoints.
func TestSys_OAuthResourceServerAuthorizationDetailsTypes(t *testing.T) {
	t.Parallel()

	mockVaultServer := httptest.NewServer(http.HandlerFunc(mockOAuthResourceServerAuthorizationDetailsTypesHandler))
	defer mockVaultServer.Close()

	cfg := DefaultConfig()
	cfg.Address = mockVaultServer.URL
	client, err := NewClient(cfg)
	require.NoError(t, err)

	keys, err := client.Sys().ListOAuthResourceServerAuthorizationDetailsTypes()
	require.NoError(t, err)
	require.Equal(t, []string{"vault:path_access"}, keys)

	catalog, err := client.Sys().ReadOAuthResourceServerAuthorizationDetailsTypesCatalog()
	require.NoError(t, err)
	require.NotNil(t, catalog)
	require.Equal(t, "v1", catalog.APIVersion)
	require.Equal(t, "deny", catalog.DefaultUnknownTypeBehavior)
	require.Len(t, catalog.Types, 1)
	require.Equal(t, "vault:path_access", catalog.Types[0].Type)
	require.Equal(t, "2.0.0-ent", catalog.Types[0].IntroducedIn)

	typeDetail, err := client.Sys().ReadOAuthResourceServerAuthorizationDetailsType("vault:path_access")
	require.NoError(t, err)
	require.NotNil(t, typeDetail)
	require.Equal(t, "vault:path_access", typeDetail.Type)
	require.Equal(t, "exact_path_match_after_namespace_normalization", typeDetail.MatchingSemantics["path_match"])
	require.Equal(t, []string{"path", "capabilities"}, typeDetail.RequiredFields)
	require.Equal(t, "/v1/sys/oauth-resource-server/authorization-details-types/vault:path_access/schema", typeDetail.SchemaPath)

	schema, err := client.Sys().ReadOAuthResourceServerAuthorizationDetailsTypeSchema("vault:path_access")
	require.NoError(t, err)
	require.NotNil(t, schema)
	require.Equal(t, "vault:path_access", schema.Type)
	require.Equal(t, "2020-12", schema.JSONSchemaDraft)
	require.Equal(t, "object", schema.Schema["type"])
}

// mockOAuthResourceServerAuthorizationDetailsTypesHandler serves mocked
// authorization_details discovery responses for API client tests.
func mockOAuthResourceServerAuthorizationDetailsTypesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/v1/sys/oauth-resource-server/authorization-details-types" && r.URL.Query().Get("list") == "true" {
		_, _ = w.Write([]byte(oauthResourceServerAuthorizationDetailsTypesListResponse))
		return
	}
	if r.URL.Path == "/v1/sys/oauth-resource-server/authorization-details-types" {
		_, _ = w.Write([]byte(oauthResourceServerAuthorizationDetailsTypesCatalogResponse))
		return
	}
	if r.URL.Path == "/v1/sys/oauth-resource-server/authorization-details-types/vault:path_access" {
		_, _ = w.Write([]byte(oauthResourceServerAuthorizationDetailsTypeResponse))
		return
	}
	if r.URL.Path == "/v1/sys/oauth-resource-server/authorization-details-types/vault:path_access/schema" {
		_, _ = w.Write([]byte(oauthResourceServerAuthorizationDetailsTypeSchemaResponse))
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

const oauthResourceServerAuthorizationDetailsTypesListResponse = `{
  "request_id": "11111111-2222-3333-4444-555555555555",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "keys": [
      "vault:path_access"
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}`

const oauthResourceServerAuthorizationDetailsTypesCatalogResponse = `{
  "request_id": "11111111-2222-3333-4444-555555555556",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "api_version": "v1",
    "default_unknown_type_behavior": "deny",
    "types": [
      {
        "type": "vault:path_access",
        "status": "active",
        "schema_version": "1.0.0",
        "title": "Vault path access",
        "description": "Constrain token use to specific Vault paths and capabilities.",
        "docs_url": "https://developer.hashicorp.com/vault/docs/auth/jwt",
        "introduced_in": "2.0.0-ent",
        "schema_path": "/v1/sys/oauth-resource-server/authorization-details-types/vault:path_access/schema"
      }
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}`

const oauthResourceServerAuthorizationDetailsTypeResponse = `{
  "request_id": "11111111-2222-3333-4444-555555555557",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "type": "vault:path_access",
    "status": "active",
    "schema_version": "1.0.0",
    "title": "Vault path access",
    "description": "Constrain token use to specific Vault paths and capabilities.",
    "docs_url": "https://developer.hashicorp.com/vault/docs/auth/jwt",
    "introduced_in": "2.0.0-ent",
    "required_fields": [
      "path",
      "capabilities"
    ],
    "optional_fields": [
      "allowed_parameters",
      "denied_parameters",
      "required_parameters"
    ],
    "matching_semantics": {
      "path_match": "exact_path_match_after_namespace_normalization",
      "capability_match": "maps authorization_details.capabilities entries to logical operations"
    },
    "examples": [
      {
        "type": "vault:path_access",
        "path": "secret/data/app/config",
        "capabilities": [
          "read"
        ]
      }
    ],
    "schema_path": "/v1/sys/oauth-resource-server/authorization-details-types/vault:path_access/schema"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}`

const oauthResourceServerAuthorizationDetailsTypeSchemaResponse = `{
  "request_id": "11111111-2222-3333-4444-555555555558",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "type": "vault:path_access",
    "schema_version": "1.0.0",
    "json_schema_draft": "2020-12",
    "schema": {
      "type": "object",
      "required": [
        "type",
        "path",
        "capabilities"
      ]
    }
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}`
