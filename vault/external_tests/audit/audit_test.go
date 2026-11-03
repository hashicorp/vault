// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/userpass"
	vaultaudit "github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func restartCoreWithAuditPrefixing(t *testing.T, cluster *vault.TestCluster, allow bool) {
	t.Helper()

	cluster.StopCore(t, 0)
	cluster.Cores[0].CoreConfig.AllowAuditLogPrefixing = allow
	cluster.StartCore(t, 0, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
		DefaultHandlerProperties: vault.HandlerProperties{
			ListenerConfig: &configutil.Listener{},
		},
	}, false)
	cluster.UnsealCores(t)
}

// TestAuditPrefixConfiguration_File verifies that file audit devices reject
// prefixes when disabled
func TestAuditPrefixConfiguration_File(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                   string
		allowAuditLogPrefixing bool
		prefix                 string
		errContains            string
	}{
		{
			name:                   "prefix forbidden but included",
			allowAuditLogPrefixing: false,
			prefix:                 "MY-TEST-PREFIX",
			errContains:            "audit prefixing is not enabled",
		},
		{
			name:                   "prefix allowed but contains newline",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\n",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains carriage return",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\r",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains null byte",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\x00",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains Unicode line separator",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\u2028",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains Unicode paragraph separator",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\u2029",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed and not included",
			allowAuditLogPrefixing: true,
		},
		{
			name:                   "prefix forbidden and not included",
			allowAuditLogPrefixing: false,
		},
		{
			name:                   "prefix allowed and included",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// start cluster with audit log prefixing enabled or disabled
			cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
				AllowAuditLogPrefixing: tc.allowAuditLogPrefixing,
			})
			if tc.allowAuditLogPrefixing {
				restartCoreWithAuditPrefixing(t, cluster, true)
			}
			client := cluster.Cores[0].Client

			auditLogPath := filepath.Join(t.TempDir(), "audit.log")
			require.NoError(t, os.WriteFile(auditLogPath, nil, 0o600))

			options := map[string]any{
				"file_path": auditLogPath,
			}
			if tc.prefix != "" {
				options["prefix"] = tc.prefix
			}

			_, err := client.Logical().Write("sys/audit/test", map[string]any{
				"type":    vaultaudit.TypeFile,
				"options": options,
			})

			if tc.errContains != "" {
				// validate request fails and no probe writes sneak through
				require.ErrorContains(t, err, tc.errContains)
				devices, listErr := client.Sys().ListAudit()
				require.NoError(t, listErr)
				require.NotContains(t, devices, "test/")
				auditFile, readErr := os.ReadFile(auditLogPath)
				require.NoError(t, readErr)
				require.Empty(t, auditFile)
				return
			}

			// validate request succeeds and file contains prefix
			require.NoError(t, err)
			auditFile, readErr := os.ReadFile(auditLogPath)
			require.NoError(t, readErr)
			require.NotEmpty(t, auditFile)
			if tc.prefix != "" {
				require.True(t, strings.HasPrefix(string(auditFile), tc.prefix))
			}
		})
	}
}

// TestAuditPrefixConfiguration_PersistedPrefix verifies that a persisted audit
// device remains available without its prefix when prefixing is later disabled.
func TestAuditPrefixConfiguration_PersistedPrefix(t *testing.T) {
	t.Parallel()

	// start vault with prefixing allowed
	config := &vault.CoreConfig{AllowAuditLogPrefixing: true}
	opts := &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	}
	cluster := vault.NewTestCluster(t, config, opts)
	restartCoreWithAuditPrefixing(t, cluster, true)
	client := cluster.Cores[0].Client

	auditLogPath := filepath.Join(t.TempDir(), "audit.log")
	prefix := "MY-TEST-PREFIX"
	require.NoError(t, client.Sys().EnableAuditWithOptions("test", &api.EnableAuditOptions{
		Type: vaultaudit.TypeFile,
		Options: map[string]string{
			"file_path":             auditLogPath,
			vaultaudit.OptionPrefix: prefix,
		},
	}))

	// confirm audit log entries contain MY-TEST-PREFIX
	auditFile, err := os.ReadFile(auditLogPath)
	require.NoError(t, err)
	require.NotEmpty(t, auditFile)
	require.True(t, strings.HasPrefix(string(auditFile), prefix))

	// stop vault, delete audit log contents
	cluster.StopCore(t, 0)
	require.NoError(t, os.Truncate(auditLogPath, 0))

	// start vault with prefixing prohibited
	cluster.Cores[0].CoreConfig.AllowAuditLogPrefixing = false
	cluster.StartCore(t, 0, opts, false)
	cluster.UnsealCores(t)

	client = cluster.Cores[0].Client
	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	require.Contains(t, devices, "test/")

	auditFile, err = os.ReadFile(auditLogPath)
	require.NoError(t, err)
	require.NotEmpty(t, auditFile)
	require.False(t, strings.HasPrefix(string(auditFile), prefix))
}

// TestAuditPrefixConfiguration_PersistedInvalidControlPrefix verifies that a
// persisted audit device with a control-character prefix is still loaded, with
// the invalid prefix ignored at setup time.
func TestAuditPrefixConfiguration_PersistedInvalidControlPrefix(t *testing.T) {
	t.Parallel()

	config := &vault.CoreConfig{
		AllowAuditLogPrefixing: true,
		EnableRaw:              true,
	}
	opts := &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	}
	cluster := vault.NewTestCluster(t, config, opts)
	restartCoreWithAuditPrefixing(t, cluster, true)
	client := cluster.Cores[0].Client

	auditLogPath := filepath.Join(t.TempDir(), "audit.log")
	require.NoError(t, client.Sys().EnableAuditWithOptions("test", &api.EnableAuditOptions{
		Type: vaultaudit.TypeFile,
		Options: map[string]string{
			"file_path":             auditLogPath,
			vaultaudit.OptionPrefix: "MY-TEST-PREFIX",
		},
	}))

	const invalidPrefix = "MY-TEST-PREFIX\n"
	rawAuditTable, err := client.Logical().Read("sys/raw/core/audit")
	require.NoError(t, err)
	require.NotNil(t, rawAuditTable)

	rawValue, ok := rawAuditTable.Data["value"].(string)
	require.True(t, ok)

	var auditTable vault.MountTable
	require.NoError(t, json.Unmarshal([]byte(rawValue), &auditTable))

	var updated bool
	for _, entry := range auditTable.Entries {
		if entry.Path != "test/" {
			continue
		}
		entry.Options[vaultaudit.OptionPrefix] = invalidPrefix
		updated = true
	}
	require.True(t, updated)

	updatedAuditTable, err := json.Marshal(auditTable)
	require.NoError(t, err)
	_, err = client.Logical().Write("sys/raw/core/audit", map[string]any{
		"value": string(updatedAuditTable),
	})
	require.NoError(t, err)

	cluster.StopCore(t, 0)
	require.NoError(t, os.Truncate(auditLogPath, 0))

	cluster.StartCore(t, 0, opts, false)
	cluster.UnsealCores(t)
	client = cluster.Cores[0].Client

	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	require.Contains(t, devices, "test/")

	_, err = client.Logical().Read("auth/token/lookup-self")
	require.NoError(t, err)

	auditFile, err := os.ReadFile(auditLogPath)
	require.NoError(t, err)
	require.NotEmpty(t, auditFile)
	require.False(t, strings.HasPrefix(string(auditFile), "MY-TEST-PREFIX"))
}

// TestAuditPrefixConfiguration_Socket verifies that socket audit devices reject
// prefixes when disabled
func TestAuditPrefixConfiguration_Socket(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                   string
		allowAuditLogPrefixing bool
		prefix                 string
		errContains            string
	}{
		{
			name:                   "prefix forbidden but included",
			allowAuditLogPrefixing: false,
			prefix:                 "MY-TEST-PREFIX",
			errContains:            "audit prefixing is not enabled",
		},
		{
			name:                   "prefix allowed but contains newline",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\n",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains carriage return",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\r",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains null byte",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\x00",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains Unicode line separator",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\u2028",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains Unicode paragraph separator",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\u2029",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed and not included",
			allowAuditLogPrefixing: true,
		},
		{
			name:                   "prefix forbidden and not included",
			allowAuditLogPrefixing: false,
		},
		{
			name:                   "prefix allowed and included",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX",
		},
	}

	socketTimeout := 2 * time.Second
	if !testhelpers.IsLocalOrRegressionTests() {
		socketTimeout = 5 * time.Second
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// configure vault cluster with audit prefixing enabled or disabled
			cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
				AllowAuditLogPrefixing: tc.allowAuditLogPrefixing,
			})
			if tc.allowAuditLogPrefixing {
				restartCoreWithAuditPrefixing(t, cluster, true)
			}
			client := cluster.Cores[0].Client

			// start tcp listener
			socketListener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1")})
			require.NoError(t, err)
			require.NoError(t, socketListener.SetDeadline(time.Now().Add(socketTimeout)))
			t.Cleanup(func() {
				require.NoError(t, socketListener.Close())
			})

			socketCh := make(chan net.Conn, 1)
			socketErrCh := make(chan error, 1)
			// separate go routine that waits for vault to make connection with tcp listener
			go func() {
				conn, acceptErr := socketListener.Accept()
				if acceptErr != nil {
					socketErrCh <- acceptErr
					return
				}
				socketCh <- conn
			}()

			options := map[string]any{
				"address":     socketListener.Addr().String(),
				"socket_type": "tcp",
			}
			if tc.prefix != "" {
				options["prefix"] = tc.prefix
			}

			_, err = client.Logical().Write("sys/audit/test", map[string]any{
				"type":    vaultaudit.TypeSocket,
				"options": options,
			})

			if tc.errContains != "" {
				// validate request fails and socket receives no probes
				require.ErrorContains(t, err, tc.errContains)
				devices, listErr := client.Sys().ListAudit()
				require.NoError(t, listErr)
				require.NotContains(t, devices, "test/")
				select {
				case conn := <-socketCh:
					require.NoError(t, conn.Close())
					require.FailNow(t, "rejected socket audit device received a test probe")
				case acceptErr := <-socketErrCh:
					// a timeout confirms that vault did not connect to the audit device
					var netErr net.Error
					require.ErrorAs(t, acceptErr, &netErr)
					require.True(t, netErr.Timeout())
				}
				return
			}

			// validate request succeeds and socket includes prefix
			require.NoError(t, err)
			select {
			case conn := <-socketCh:
				t.Cleanup(func() {
					require.NoError(t, conn.Close())
				})
				// add timeout so that ReadFull doesn't wait indefinitely
				require.NoError(t, conn.SetReadDeadline(time.Now().Add(socketTimeout)))
				if tc.prefix != "" {
					socketPrefix := make([]byte, len(tc.prefix))
					_, readErr := io.ReadFull(conn, socketPrefix)
					require.NoError(t, readErr)
					require.Equal(t, tc.prefix, string(socketPrefix))
					return
				}

				auditEntry := make([]byte, 1)
				_, readErr := io.ReadFull(conn, auditEntry)
				require.NoError(t, readErr)
			case acceptErr := <-socketErrCh:
				require.NoError(t, acceptErr)
			}
		})
	}
}

// TestAuditPrefixConfiguration_Syslog verifies that syslog audit devices reject
// prefixes when disabled
func TestAuditPrefixConfiguration_Syslog(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                   string
		allowAuditLogPrefixing bool
		prefix                 string
		errContains            string
	}{
		{
			name:                   "prefix forbidden but included",
			allowAuditLogPrefixing: false,
			prefix:                 "MY-TEST-PREFIX",
			errContains:            "audit prefixing is not enabled",
		},
		{
			name:                   "prefix allowed but contains newline",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\n",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains carriage return",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\r",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains null byte",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\x00",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains Unicode line separator",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\u2028",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed but contains Unicode paragraph separator",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX\u2029",
			errContains:            "audit prefix may not contain control characters",
		},
		{
			name:                   "prefix allowed and not included",
			allowAuditLogPrefixing: true,
		},
		{
			name:                   "prefix forbidden and not included",
			allowAuditLogPrefixing: false,
		},
		{
			name:                   "prefix allowed and included",
			allowAuditLogPrefixing: true,
			prefix:                 "MY-TEST-PREFIX",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// configure vault cluster with audit prefixing allowed or disallowed
			cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
				AllowAuditLogPrefixing: tc.allowAuditLogPrefixing,
			})
			if tc.allowAuditLogPrefixing {
				restartCoreWithAuditPrefixing(t, cluster, true)
			}
			client := cluster.Cores[0].Client

			options := map[string]any{
				"facility": "LOCAL0",
			}
			if tc.prefix != "" {
				options["prefix"] = tc.prefix
			}

			_, err := client.Logical().Write("sys/audit/test", map[string]any{
				"type":    vaultaudit.TypeSyslog,
				"options": options,
			})

			if tc.errContains != "" {
				require.ErrorContains(t, err, tc.errContains)
				devices, listErr := client.Sys().ListAudit()
				require.NoError(t, listErr)
				require.NotContains(t, devices, "test/")
				return
			}
			require.NoError(t, err)
		})
	}
}

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

type testAuditStartupPlugin struct {
	*framework.Backend
}

func testAuditStartupBackend(t testing.TB, panic *atomic.Bool) logical.Factory {
	return func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		be := &testAuditStartupPlugin{Backend: &framework.Backend{BackendType: logical.TypeLogical}}
		if err := be.Setup(ctx, config); err != nil {
			return nil, err
		}
		extendedSys, ok := be.System().(logical.ExtendedSystemView)
		if !ok {
			return nil, fmt.Errorf("expected ExtendedSystemView, got %T", be.System())
		}

		wait := make(chan struct{})
		// run in a goroutine to more closely mimic the behavior of a plugin where this causes a
		// problem.
		go func() {
			defer close(wait)
			defer func() {
				if r := recover(); r != nil {
					be.Logger().Error("panic during setup", "error", r)
					panic.Store(true)
				}
			}()
			// try to audit a request before the auditBroker has been created
			if err := extendedSys.Auditor().AuditRequest(ctx, &logical.LogInput{}); err != nil && !errors.Is(err, consts.ErrSealed) {
				be.Logger().Error("error auditing request", "error", err)
			}
		}()
		<-wait
		return be, nil
	}
}

// TestAudit_BeforePostUnseal verifies that an audit request can be made before unseal without causing a panic. The test
// mounts a backend that will attempt to audit a request in a goroutine during its setup. The seals and unseals the
// cluster to force the backend to be re-created and perform an audit before the audit broker has been created.
func TestAudit_BeforePostUnseal(t *testing.T) {
	didPanic := new(atomic.Bool)
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"test": testAuditStartupBackend(t, didPanic),
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})

	testhelpers.WaitForActiveNode(t, cluster)
	err := cluster.Cores[0].Client.Sys().Mount("test", &api.MountInput{
		Type: "test",
	})
	require.NoError(t, err)
	// Seal and unseal to trigger a re-creation of all the mounts
	cluster.Cores[0].Seal(t)
	cluster.UnsealCores(t)
	testhelpers.WaitForActiveNode(t, cluster)
	require.False(t, didPanic.Load())
}

// TestAudit_PluginDirectorySecurityCheck_WithPathOrFilePath validates that the
// audit security check still rejects a real audit target inside the plugin
// directory and that Windows-style drive-root paths passed via either 'path' or
// 'file_path' do not trigger false plugin-directory rejections
func TestAudit_PluginDirectorySecurityCheck_WithPathOrFilePath(t *testing.T) {
	pluginDir := corehelpers.MakeTestPluginDir(t)

	tests := []struct {
		name        string
		pluginDir   string
		auditPath   string
		wantErr     bool
		errContains string
	}{
		{
			name:        "inside plugin dir is rejected",
			pluginDir:   pluginDir,
			auditPath:   pluginDir + "/audit.log",
			wantErr:     true,
			errContains: "audit file target may not be in the plugin directory",
		},
		{
			name:      "windows: outside plugin dir is allowed",
			pluginDir: `C:\Vault\plugins`,
			auditPath: `C:\Vault\logs\audit.log`,
		},
		{
			name:      "windows: drive root is allowed",
			pluginDir: `C:\Vault\plugins`,
			auditPath: `C:\audit.log`,
		},
		{
			name:      "windows: alternate drive outside plugin dir is allowed",
			pluginDir: `D:\Vault\plugins`,
			auditPath: `D:\Vault\logs\audit.log`,
		},
		{
			name:      "windows: alternate drive root is allowed",
			pluginDir: `D:\Vault\plugins`,
			auditPath: `D:\audit.log`,
		},
		{
			name:      "windows: cross-volume audit path is allowed",
			pluginDir: `D:\Vault\plugins`,
			auditPath: `E:\audit.log`,
		},
	}

	for _, optionKey := range []string{"path", "file_path"} {
		for _, tc := range tests {
			t.Run(fmt.Sprintf("%s/%s", optionKey, tc.name), func(t *testing.T) {
				cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
					PluginDirectory: tc.pluginDir,
				})
				client := cluster.Cores[0].Client

				deviceData := map[string]any{
					"type":        "file",
					"description": "test audit device",
					"local":       false,
					"options": map[string]any{
						optionKey: tc.auditPath,
					},
				}

				_, err := client.Logical().Write("sys/audit/test", deviceData)
				if tc.wantErr {
					require.Error(t, err)
					require.ErrorContains(t, err, tc.errContains)
					return
				}

				if err != nil {
					require.NotContains(t, err.Error(), "audit file target may not be in the plugin directory")
				}
			})
		}
	}
}
