// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestSupplementalAuditData validates if a plugin populates the logical.Response SupplementalAuditRequestData and
// SupplementalAuditResponseData that we populate the audit entry appropriately also applying the mount's
// HMAC keys to the appropriate request/response fields
func TestSupplementalAuditData(t *testing.T) {
	t.Parallel()

	testHandlerWithAuditOnly := func(ctx context.Context, l *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"secret": "my-fancy-secret",
			},
			SupplementalAuditRequestData: map[string]any{
				"foo": "bar",
				"baz": "qux",
			},
			SupplementalAuditResponseData: map[string]any{
				"foo":  "bar",
				"baz":  "qux",
				"quux": "corge",
			},
		}, nil
	}

	testHandlerNoAudit := func(ctx context.Context, l *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"secret": "my-fancy-secret",
			},
		}, nil
	}

	operationsWithAudit := map[logical.Operation]framework.OperationHandler{
		logical.UpdateOperation: &framework.PathOperation{Callback: testHandlerWithAuditOnly},
	}

	operationsNoAudit := map[logical.Operation]framework.OperationHandler{
		logical.UpdateOperation: &framework.PathOperation{Callback: testHandlerNoAudit},
	}

	conf := &vault.CoreConfig{
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		LogicalBackends: map[string]logical.Factory{
			"audittest": func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
				b := new(framework.Backend)
				b.BackendType = logical.TypeLogical
				b.Paths = []*framework.Path{
					{Pattern: "with-audit-fields", Operations: operationsWithAudit},
					{Pattern: "no-audit-fields", Operations: operationsNoAudit},
				}
				err := b.Setup(ctx, config)
				return b, err
			},
		},
	}

	cluster := minimal.NewTestSoloCluster(t, conf)
	client := cluster.Cores[0].Client

	auditLog := filepath.Join(t.TempDir(), "audit.log")
	devicePath := "file"
	deviceData := map[string]any{
		"type": "file",
		"options": map[string]any{
			"file_path": auditLog,
		},
	}
	_, err := client.Logical().Write("sys/audit/"+devicePath, deviceData)
	require.NoError(t, err)

	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	require.Len(t, devices, 1)

	err = client.Sys().Mount("audittest", &api.MountInput{
		Type: "audittest",
		Config: api.MountConfigInput{
			AuditNonHMACRequestKeys:  []string{"foo"},
			AuditNonHMACResponseKeys: []string{"baz", "secret"},
		},
	})
	require.NoError(t, err)

	// Call our API with audit fields
	resp, err := client.Logical().Write("audittest/with-audit-fields", map[string]interface{}{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	// Make sure we only have 1 element and it's secret within the Data field
	require.Len(t, resp.Data, 1)
	require.Contains(t, resp.Data, "secret")

	// Call our API with no audit fields
	resp, err = client.Logical().Write("audittest/no-audit-fields", map[string]interface{}{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	// Make sure we only have 1 element and it's secret within the Data field
	require.Len(t, resp.Data, 1)
	require.Contains(t, resp.Data, "secret")

	entries := make([]map[string]interface{}, 0)
	auditFile, err := os.OpenFile(auditLog, os.O_RDONLY, 0o644)
	require.NoError(t, err, "failed to open audit log")
	scanner := bufio.NewScanner(auditFile)

	// Collect the two entries we really care about
	for scanner.Scan() {
		entry := make(map[string]interface{})

		err := json.Unmarshal(scanner.Bytes(), &entry)
		require.NoError(t, err)

		if isResponseEntryForPath(entry, "audittest/with-audit-fields") {
			entries = append(entries, entry)
		}

		if isResponseEntryForPath(entry, "audittest/no-audit-fields") {
			entries = append(entries, entry)
		}
	}

	// We expect to have 2 entries, the first with audit_only_fields set, the other doesn't
	require.Equal(t, 2, len(entries))

	{
		// Make sure the request object within has audit_only_fields, it should contain a different
		// set of keys than the response audit_only_fields, and the values should have been hmac'd
		// based on the mount's AuditNonHMACRequestKeys value which is ["foo"]
		entryWithAuditFields := entries[0]

		entryRequest := castToMap(t, entryWithAuditFields["request"])
		require.Contains(t, entryRequest, "supplemental_audit_data")
		entryRequestAuditOnlyFields := castToStringMap(t, entryRequest["supplemental_audit_data"])
		require.Contains(t, entryRequestAuditOnlyFields, "foo")
		require.Contains(t, entryRequestAuditOnlyFields, "baz")
		requireHmaced(t, entryRequestAuditOnlyFields["baz"])
		require.Equal(t, "bar", entryRequestAuditOnlyFields["foo"])
		require.Len(t, entryRequestAuditOnlyFields, 2)
	}

	{
		// Make sure the audit response data field only contains the secret field and not any of our audit only fields
		entryWithAuditFields := entries[0]
		entryResponse := castToMap(t, entryWithAuditFields["response"])
		entryResponseData := castToStringMap(t, entryResponse["data"])
		require.Contains(t, entryResponseData, "secret")
		require.Len(t, entryResponseData, 1)
		require.Equal(t, "my-fancy-secret", entryResponseData["secret"])
	}

	{
		// Make sure the audit response audit only fields contains the three fields we set, and we properly
		// applied the AuditNonHMACResponseKeys to those keys, see mount config above, but we expect keys
		// ["baz", "secret"] to be cleared
		entryWithAuditFields := entries[0]
		entryResponse := castToMap(t, entryWithAuditFields["response"])
		entryResponseAuditOnly := castToStringMap(t, entryResponse["supplemental_audit_data"])
		require.Contains(t, entryResponseAuditOnly, "foo")
		require.Contains(t, entryResponseAuditOnly, "baz")
		require.Contains(t, entryResponseAuditOnly, "quux")
		requireHmaced(t, entryResponseAuditOnly["foo"])
		requireHmaced(t, entryResponseAuditOnly["quux"])
		require.Equal(t, "qux", entryResponseAuditOnly["baz"])
		require.Len(t, entryResponseAuditOnly, 3)
	}

	{
		// Now validate the audit entry with no additional audit entries doesn't have the new fields in the audit entry
		entryNoAudit := entries[1]
		entryRequest := castToMap(t, entryNoAudit["request"])
		require.NotContains(t, entryRequest, "supplemental_audit_data")

		entryResponseNoAudit := castToMap(t, entryNoAudit["response"])
		require.NotContains(t, entryResponseNoAudit, "supplemental_audit_data")

		// We still should see the secret not hmac'd in the response data
		entryResponseData := castToStringMap(t, entryResponseNoAudit["data"])
		require.Contains(t, entryResponseData, "secret")
		require.Len(t, entryResponseData, 1)
		require.Equal(t, "my-fancy-secret", entryResponseData["secret"])
	}
}

func requireHmaced(t testing.TB, val string) {
	t.Helper()

	parts := strings.Split(val, ":")
	require.Len(t, parts, 2, "splitting hmac value %q should have 2 parts found %d", val, len(parts))
	require.Equal(t, "hmac-sha256", parts[0])
	require.Equal(t, 64, len(parts[1]), "expected hmac'd field %q to have a length of 64 characters", val)
}

func castToStringMap(t testing.TB, val interface{}) map[string]string {
	valMap, ok := val.(map[string]interface{})
	if !ok {
		t.Fatalf("value is not a map was: %T", val)
	}

	stringMap := make(map[string]string)
	for k, v := range valMap {
		if s, ok := v.(string); ok {
			stringMap[k] = s
		} else {
			t.Fatalf("value for key %q is not a string was: %T", k, v)
		}
	}

	return stringMap
}

func castToMap(t testing.TB, val interface{}) map[string]interface{} {
	t.Helper()

	if valMap, ok := val.(map[string]interface{}); ok {
		return valMap
	}

	t.Fatalf("Value is not a map was %T", val)
	return nil
}

func isResponseEntryForPath(entry map[string]interface{}, desiredPath string) bool {
	if typeRaw, ok := entry["type"]; !ok || typeRaw != "response" {
		return false
	}
	requestRaw, ok := entry["request"]
	if !ok || requestRaw == nil {
		return false
	}
	request := requestRaw.(map[string]interface{})
	if pathRaw, ok := request["path"]; !ok || pathRaw != desiredPath {
		return false
	}
	return true
}
