// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/userpass"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestAudit_HMACFields verifies that all appropriate fields in audit
// request and response entries are HMACed properly. The fields in question are:
//   - request.headers.x-correlation-id
//   - request.data: all sub-fields
//   - respnse.auth.client_token
//   - response.auth.accessor
//   - response.data: all sub-fields
//   - response.wrap_info.token
//   - response.wrap_info.accessor
func TestAudit_HMACFields(t *testing.T) {
	const hmacPrefix = "hmac-sha256:"

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	tempDir := t.TempDir()
	logFile, err := os.CreateTemp(tempDir, "")
	require.NoError(t, err)
	devicePath := "file"
	deviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": logFile.Name(),
		},
	}

	_, err = client.Logical().Write("sys/config/auditing/request-headers/x-correlation-id", map[string]interface{}{
		"hmac": true,
	})
	require.NoError(t, err)

	// Request 1
	// Enable the audit device. A test probe request will audited along
	// with the associated creation response
	_, err = client.Logical().Write("sys/audit/"+devicePath, deviceData)
	require.NoError(t, err)

	// Request 2
	// Ensure the device has been created.
	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	require.Len(t, devices, 1)

	// Request 3
	// Enable the userpass auth method (this will be an audited action)
	err = client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	require.NoError(t, err)

	username := "jdoe"
	password := "abc123"

	// Request 4
	// Create a user with a password (another audited action)
	_, err = client.Logical().Write(fmt.Sprintf("auth/userpass/users/%s", username), map[string]interface{}{
		"password": password,
	})
	require.NoError(t, err)

	authInput, err := userpass.NewUserpassAuth(username, &userpass.Password{FromString: password})
	require.NoError(t, err)

	newClient, err := client.Clone()
	require.NoError(t, err)

	correlationID := "correlation-id-foo"
	newClient.AddHeader("x-correlation-id", correlationID)

	// Request 5
	authOutput, err := newClient.Auth().Login(context.Background(), authInput)
	require.NoError(t, err)

	// Request 6
	hashedPassword, err := client.Sys().AuditHash(devicePath, password)
	require.NoError(t, err)

	// Request 7
	hashedClientToken, err := client.Sys().AuditHash(devicePath, authOutput.Auth.ClientToken)
	require.NoError(t, err)

	// Request 8
	hashedAccessor, err := client.Sys().AuditHash(devicePath, authOutput.Auth.Accessor)
	require.NoError(t, err)

	// Request 9
	wrapResp, err := client.Logical().Write("sys/wrapping/wrap", map[string]interface{}{
		"foo": "bar",
	})
	require.NoError(t, err)

	// Request 10
	hashedBar, err := client.Sys().AuditHash(devicePath, "bar")
	require.NoError(t, err)

	// Request 11
	hashedWrapAccessor, err := client.Sys().AuditHash(devicePath, wrapResp.WrapInfo.Accessor)
	require.NoError(t, err)

	// Request 12
	hashedWrapToken, err := client.Sys().AuditHash(devicePath, wrapResp.WrapInfo.Token)
	require.NoError(t, err)

	// Request 13
	hashedCorrelationID, err := client.Sys().AuditHash(devicePath, correlationID)
	require.NoError(t, err)

	// Request 14
	// Disable the audit device. The request will be audited but not the response.
	_, err = client.Logical().Delete("sys/audit/" + devicePath)
	require.NoError(t, err)

	// Request 15
	// Ensure the device has been deleted. This will not be audited.
	devices, err = client.Sys().ListAudit()
	require.NoError(t, err)
	require.Len(t, devices, 0)

	entries := make([]map[string]interface{}, 0)
	scanner := bufio.NewScanner(logFile)

	for scanner.Scan() {
		entry := make(map[string]interface{})

		err := json.Unmarshal(scanner.Bytes(), &entry)
		require.NoError(t, err)

		entries = append(entries, entry)
	}

	// This count includes the initial test probe upon creation of the audit device
	require.Equal(t, 27, len(entries))

	loginReqEntry := entries[8]
	loginRespEntry := entries[9]

	loginRequestFromReq := loginReqEntry["request"].(map[string]interface{})
	loginRequestDataFromReq := loginRequestFromReq["data"].(map[string]interface{})
	loginHeadersFromReq := loginRequestFromReq["headers"].(map[string]interface{})

	loginRequestFromResp := loginRespEntry["request"].(map[string]interface{})
	loginRequestDataFromResp := loginRequestFromResp["data"].(map[string]interface{})
	loginHeadersFromResp := loginRequestFromResp["headers"].(map[string]interface{})

	loginAuth := loginRespEntry["auth"].(map[string]interface{})

	require.True(t, strings.HasPrefix(loginRequestDataFromReq["password"].(string), hmacPrefix))
	require.Equal(t, loginRequestDataFromReq["password"].(string), hashedPassword)

	require.True(t, strings.HasPrefix(loginRequestDataFromResp["password"].(string), hmacPrefix))
	require.Equal(t, loginRequestDataFromResp["password"].(string), hashedPassword)

	require.True(t, strings.HasPrefix(loginAuth["client_token"].(string), hmacPrefix))
	require.Equal(t, loginAuth["client_token"].(string), hashedClientToken)

	require.True(t, strings.HasPrefix(loginAuth["accessor"].(string), hmacPrefix))
	require.Equal(t, loginAuth["accessor"].(string), hashedAccessor)

	xCorrelationIDFromReq := loginHeadersFromReq["x-correlation-id"].([]interface{})
	require.Equal(t, len(xCorrelationIDFromReq), 1)
	require.True(t, strings.HasPrefix(xCorrelationIDFromReq[0].(string), hmacPrefix))
	require.Equal(t, xCorrelationIDFromReq[0].(string), hashedCorrelationID)

	xCorrelationIDFromResp := loginHeadersFromResp["x-correlation-id"].([]interface{})
	require.Equal(t, len(xCorrelationIDFromResp), 1)
	require.True(t, strings.HasPrefix(xCorrelationIDFromReq[0].(string), hmacPrefix))
	require.Equal(t, xCorrelationIDFromResp[0].(string), hashedCorrelationID)

	wrapReqEntry := entries[16]
	wrapRespEntry := entries[17]

	wrapRequestFromReq := wrapReqEntry["request"].(map[string]interface{})
	wrapRequestDataFromReq := wrapRequestFromReq["data"].(map[string]interface{})

	wrapRequestFromResp := wrapRespEntry["request"].(map[string]interface{})
	wrapRequestDataFromResp := wrapRequestFromResp["data"].(map[string]interface{})

	require.True(t, strings.HasPrefix(wrapRequestDataFromReq["foo"].(string), hmacPrefix))
	require.Equal(t, wrapRequestDataFromReq["foo"].(string), hashedBar)

	require.True(t, strings.HasPrefix(wrapRequestDataFromResp["foo"].(string), hmacPrefix))
	require.Equal(t, wrapRequestDataFromResp["foo"].(string), hashedBar)

	wrapResponseData := wrapRespEntry["response"].(map[string]interface{})
	wrapInfo := wrapResponseData["wrap_info"].(map[string]interface{})

	require.True(t, strings.HasPrefix(wrapInfo["accessor"].(string), hmacPrefix))
	require.Equal(t, wrapInfo["accessor"].(string), hashedWrapAccessor)

	require.True(t, strings.HasPrefix(wrapInfo["token"].(string), hmacPrefix))
	require.Equal(t, wrapInfo["token"].(string), hashedWrapToken)
}

// TestAudit_Headers validates that headers are audited correctly. This includes
// the default headers (x-correlation-id and user-agent) along with user-specified
// headers.
func TestAudit_Headers(t *testing.T) {
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	tempDir := t.TempDir()
	logFile, err := os.CreateTemp(tempDir, "")
	require.NoError(t, err)
	devicePath := "file"
	deviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": logFile.Name(),
		},
	}

	_, err = client.Logical().Write("sys/config/auditing/request-headers/x-some-header", map[string]interface{}{
		"hmac": false,
	})
	require.NoError(t, err)

	// User-Agent header is audited by default
	client.AddHeader("User-Agent", "foo-agent")

	// X-Some-Header has been added to audited headers manually
	client.AddHeader("X-Some-Header", "some-value")

	// X-Some-Other-Header will not be audited
	client.AddHeader("X-Some-Other-Header", "some-other-value")

	// Request 1
	// Enable the audit device. A test probe request will audited along
	// with the associated creation response
	_, err = client.Logical().Write("sys/audit/"+devicePath, deviceData)
	require.NoError(t, err)

	// Request 2
	// Ensure the device has been created.
	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	require.Len(t, devices, 1)

	// Request 3
	resp, err := client.Sys().SealStatus()
	require.NoError(t, err)
	require.NotEmpty(t, resp)

	expectedHeaders := map[string]interface{}{
		"user-agent":    []interface{}{"foo-agent"},
		"x-some-header": []interface{}{"some-value"},
	}

	entries := make([]map[string]interface{}, 0)
	scanner := bufio.NewScanner(logFile)

	for scanner.Scan() {
		entry := make(map[string]interface{})

		err := json.Unmarshal(scanner.Bytes(), &entry)
		require.NoError(t, err)

		request, ok := entry["request"].(map[string]interface{})
		require.True(t, ok)

		// test probe will not have headers set
		requestPath, ok := request["path"].(string)
		require.True(t, ok)

		if requestPath != "sys/audit/test" {
			headers, ok := request["headers"].(map[string]interface{})

			require.True(t, ok)
			require.Equal(t, expectedHeaders, headers)
		}

		entries = append(entries, entry)
	}

	// This count includes the initial test probe upon creation of the audit device
	require.Equal(t, 4, len(entries))
}
