// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/namespace"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

// TestGcpKmsDataProtectionCallCounts tests that we correctly store and track
// the GCP KMS data protection call counts by simulating billing operations.
func TestGcpKmsDataProtectionCallCounts(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		BillingConfig: billing.BillingConfig{
			MetricsUpdateCadence: 3 * time.Second,
		},
	}
	core, _, _, _ := vault.TestCoreUnsealedWithMetricsAndConfig(t, coreConfig)

	currentMonth := time.Now()
	ctx := namespace.RootContext(context.Background())

	// Get the consumption billing manager
	cbm := core.GetConsumptionBillingManager()
	require.NotNil(t, cbm)

	// Simulate GCP KMS operations by directly calling the billing manager
	// This tests the Vault-side tracking without needing the actual plugin

	// Simulate encrypt operation
	err := cbm.WriteBillingData(ctx, "gcpkms", map[string]interface{}{"count": uint64(1)})
	require.NoError(t, err)
	require.Equal(t, uint64(1), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Wait for storage update
	require.Eventually(t, func() bool {
		counts, err := core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
		return err == nil && counts == 1
	}, 5*time.Second, 100*time.Millisecond)
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Simulate decrypt operation
	err = cbm.WriteBillingData(ctx, "gcpkms", map[string]interface{}{"count": uint64(1)})
	require.NoError(t, err)
	require.Equal(t, uint64(1), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Wait for storage update
	require.Eventually(t, func() bool {
		counts, err := core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
		return err == nil && counts == 2
	}, 5*time.Second, 100*time.Millisecond)
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Simulate reencrypt operation
	err = cbm.WriteBillingData(ctx, "gcpkms", map[string]interface{}{"count": uint64(1)})
	require.NoError(t, err)
	require.Equal(t, uint64(1), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Wait for storage update
	require.Eventually(t, func() bool {
		counts, err := core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
		return err == nil && counts == 3
	}, 5*time.Second, 100*time.Millisecond)
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Simulate sign operation
	err = cbm.WriteBillingData(ctx, "gcpkms", map[string]interface{}{"count": uint64(1)})
	require.NoError(t, err)
	require.Equal(t, uint64(1), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Wait for storage update
	require.Eventually(t, func() bool {
		counts, err := core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
		return err == nil && counts == 4
	}, 5*time.Second, 100*time.Millisecond)
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Simulate verify operation
	err = cbm.WriteBillingData(ctx, "gcpkms", map[string]interface{}{"count": uint64(1)})
	require.NoError(t, err)
	require.Equal(t, uint64(1), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Wait for storage update
	require.Eventually(t, func() bool {
		counts, err := core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
		return err == nil && counts == 5
	}, 5*time.Second, 100*time.Millisecond)
	require.Equal(t, uint64(0), core.GetInMemoryGcpKmsDataProtectionCallCounts())

	// Run update again and make sure the value in storage is still 5
	counts, err := core.UpdateGcpKmsCallCounts(context.Background(), currentMonth)
	require.NoError(t, err)
	require.Equal(t, uint64(5), counts)

	// Verify the value in storage is still 5
	counts, err = core.GetStoredGcpKmsCallCounts(context.Background(), currentMonth)
	require.NoError(t, err)
	require.Equal(t, uint64(5), counts)
}

// TestOidcTokenBillingBothMethods tests OIDC token billing for both token creation methods:
// 1. Simple role-based tokens via identity/oidc/token/{role}
// 2. Provider-based tokens via the full authorization code flow (pathOIDCToken)
// This test runs on a single primary cluster and verifies that both methods correctly
// track duration-adjusted billing counts.
func TestOidcTokenBillingBothMethods(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
		BillingConfig: billing.BillingConfig{
			MetricsUpdateCadence: 5 * time.Second,
		},
	}
	clusterOpts := &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		DefaultHandlerProperties: vault.HandlerProperties{
			ListenerConfig: &configutil.Listener{},
		},
		NumCores: 1,
	}
	cluster := vault.NewTestCluster(t, coreConfig, clusterOpts)
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client
	ctx := context.Background()

	// Create a policy that allows reading OIDC tokens
	oidcPolicy := `path "identity/oidc/token/*" { capabilities = ["read"] }`
	_, err := client.Logical().Write("sys/policy/oidc-reader", map[string]interface{}{
		"policy": oidcPolicy,
	})
	require.NoError(t, err)

	// Enable userpass for entity creation
	err = client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	require.NoError(t, err)

	// Create a userpass user with the OIDC reader policy
	_, err = client.Logical().Write("auth/userpass/users/testuser", map[string]interface{}{
		"password": "testpass",
		"policies": "oidc-reader",
	})
	require.NoError(t, err)

	// Login to create entity
	loginResp, err := client.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
		"password": "testpass",
	})
	require.NoError(t, err)
	userToken := loginResp.Auth.ClientToken

	// METHOD 1: Configure simple role-based OIDC tokens (identity/oidc/token/{role})
	// Create OIDC key
	_, err = client.Logical().Write("identity/oidc/key/role-key", map[string]interface{}{})
	require.NoError(t, err)

	// Create OIDC role with 1-hour TTL
	_, err = client.Logical().Write("identity/oidc/role/test-role", map[string]interface{}{
		"key": "role-key",
		"ttl": "1h",
	})
	require.NoError(t, err)

	// Get the auto-generated client_id for the role
	secret, err := client.Logical().Read("identity/oidc/role/test-role")
	require.NoError(t, err)
	roleClientID := secret.Data["client_id"].(string)

	// Configure the key to allow this role's client_id
	_, err = client.Logical().Write("identity/oidc/key/role-key", map[string]interface{}{
		"allowed_client_ids": roleClientID,
	})
	require.NoError(t, err)

	// METHOD 2: Configure OIDC provider for authorization code flow (pathOIDCToken)
	// Create OIDC client with 2-hour ID token TTL and 1-hour access token TTL
	_, err = client.Logical().Write("identity/oidc/client/provider-client", map[string]interface{}{
		"redirect_uris":    []string{"https://localhost:8251/callback"},
		"assignments":      []string{"allow_all"},
		"id_token_ttl":     "2h",
		"access_token_ttl": "1h",
	})
	require.NoError(t, err)

	// Read the client to get client_id and client_secret
	clientResp, err := client.Logical().Read("identity/oidc/client/provider-client")
	require.NoError(t, err)
	providerClientID := clientResp.Data["client_id"].(string)
	providerClientSecret := clientResp.Data["client_secret"].(string)

	// Create OIDC provider
	_, err = client.Logical().Write("identity/oidc/provider/test-provider", map[string]interface{}{
		"allowed_client_ids": []string{providerClientID},
	})
	require.NoError(t, err)

	// Generate tokens using METHOD 1: role-based (identity/oidc/token/{role})
	// 2 tokens × 1 hour = 2 hours
	client.SetToken(userToken)
	for i := 0; i < 2; i++ {
		_, err := client.Logical().Read("identity/oidc/token/test-role")
		require.NoError(t, err)
	}

	// Generate tokens using METHOD 2: provider-based (authorization code flow)
	// 3 tokens × 2 hours (max of 2h ID token and 1h access token) = 6 hours
	client.SetToken(client.Token()) // Reset to root token
	for i := 0; i < 3; i++ {
		code := getAuthorizationCode(t, ctx, client, "test-provider", providerClientID, userToken)
		exchangeCodeForToken(t, ctx, client, "test-provider", code, providerClientID, providerClientSecret)
	}

	currentMonth := time.Now().UTC()

	// Total expected: 2 hours (role-based) + 6 hours (provider-based) = 8 hours
	expectedDurationAdjustedCount := vault.DurationAdjustedTokenCount(8 * time.Hour.Seconds())
	delta := 0.0001

	require.Eventually(t, func() bool {
		count, err := core.GetStoredOidcDurationAdjustedCount(ctx, currentMonth)
		if err != nil {
			return false
		}
		return count >= (expectedDurationAdjustedCount-delta) && count <= (expectedDurationAdjustedCount+delta)
	}, 10*time.Second, 500*time.Millisecond, "OIDC count not flushed to storage within timeout")

	// Verify exact value
	count, err := core.GetStoredOidcDurationAdjustedCount(ctx, currentMonth)
	require.NoError(t, err)
	require.InDelta(t, expectedDurationAdjustedCount, count, delta,
		"Expected 8 hours total: 2 hours from role-based tokens (2×1h) + 6 hours from provider tokens (3×2h)")
}

// exchangeCodeForToken is a test helper function to exchange authorization code for tokens via the OIDC provider token endpoint
func exchangeCodeForToken(t *testing.T, ctx context.Context, client *api.Client, providerName, code, clientID, clientSecret string) {
	// Prepare the token request with basic auth
	req := client.NewRequest("POST", "/v1/identity/oidc/provider/"+providerName+"/token")
	req.Headers.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
	req.BodyBytes = []byte(fmt.Sprintf(`{"code":"%s","grant_type":"authorization_code","redirect_uri":"https://localhost:8251/callback"}`, code))
	req.Headers.Set("Content-Type", "application/json")

	resp, err := client.RawRequestWithContext(ctx, req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)
}

// getAuthorizationCode is a test helper function to get authorization code from the OIDC provider authorize endpoint
func getAuthorizationCode(t *testing.T, ctx context.Context, client *api.Client, providerName, clientID, userToken string) string {
	// Save the original token
	originalToken := client.Token()

	// Use the user token (from userpass login) to authorize
	client.SetToken(userToken)

	// Use RawRequestWithContext to make the authorize request
	req := client.NewRequest("POST", "/v1/identity/oidc/provider/"+providerName+"/authorize")
	req.BodyBytes, _ = json.Marshal(map[string]interface{}{
		"client_id":     clientID,
		"scope":         "openid",
		"redirect_uri":  "https://localhost:8251/callback",
		"response_type": "code",
		"state":         "test-state",
		"nonce":         "test-nonce",
	})

	resp, err := client.RawRequestWithContext(ctx, req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Restore the original token
	client.SetToken(originalToken)

	// Parse the JSON response
	var authResult struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}
	err = json.NewDecoder(resp.Body).Decode(&authResult)
	require.NoError(t, err)
	require.NotEmpty(t, authResult.Code)

	return authResult.Code
}
