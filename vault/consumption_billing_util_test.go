// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	logicalAlicloud "github.com/hashicorp/vault-plugin-secrets-alicloud"
	logicalAzure "github.com/hashicorp/vault-plugin-secrets-azure"
	logicalGcp "github.com/hashicorp/vault-plugin-secrets-gcp/plugin"
	logicalKubernetes "github.com/hashicorp/vault-plugin-secrets-kubernetes"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	logicalMongoDBAtlas "github.com/hashicorp/vault-plugin-secrets-mongodbatlas"
	logicalLDAP "github.com/hashicorp/vault-plugin-secrets-openldap"
	logicalTerraform "github.com/hashicorp/vault-plugin-secrets-terraform"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	logicalAws "github.com/hashicorp/vault/builtin/logical/aws"
	logicalConsul "github.com/hashicorp/vault/builtin/logical/consul"
	logicalDatabase "github.com/hashicorp/vault/builtin/logical/database"
	logicalNomad "github.com/hashicorp/vault/builtin/logical/nomad"
	logicalRabbitMQ "github.com/hashicorp/vault/builtin/logical/rabbitmq"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

var roleLogicalBackends = map[string]logical.Factory{
	pluginconsts.SecretEngineAWS:          logicalAws.Factory,
	pluginconsts.SecretEngineAzure:        logicalAzure.Factory,
	pluginconsts.SecretEngineGCP:          logicalGcp.Factory,
	pluginconsts.SecretEngineKV:           logicalKv.Factory,
	pluginconsts.SecretEngineLDAP:         logicalLDAP.Factory,
	pluginconsts.SecretEngineDatabase:     logicalDatabase.Factory,
	pluginconsts.SecretEngineOpenLDAP:     logicalLDAP.Factory,
	pluginconsts.SecretEngineAlicloud:     logicalAlicloud.Factory,
	pluginconsts.SecretEngineRabbitMQ:     logicalRabbitMQ.Factory,
	pluginconsts.SecretEngineConsul:       logicalConsul.Factory,
	pluginconsts.SecretEngineNomad:        logicalNomad.Factory,
	pluginconsts.SecretEngineKubernetes:   logicalKubernetes.Factory,
	pluginconsts.SecretEngineMongoDBAtlas: logicalMongoDBAtlas.Factory,
	pluginconsts.SecretEngineTerraform:    logicalTerraform.Factory,
}

// TestStoreAndGetMaxRoleCounts verifies that we can store and retrieve the HWM role counts correctly
func TestStoreAndGetMaxRoleCounts(t *testing.T) {
	coreConfig := &CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			pluginconsts.AuthTypeUserpass: userpass.Factory,
		},
	}
	core, _, _ := TestCoreUnsealedWithConfig(t, coreConfig)

	testCases := []struct {
		description     string
		localPathPrefix string
		monthOffset     int
		roleCounts      *RoleCounts
	}{
		{
			description:     "Local storage, current month",
			localPathPrefix: billing.LocalPrefix,
			monthOffset:     0,
			roleCounts: &RoleCounts{
				AWSDynamicRoles:  10,
				AWSStaticRoles:   15,
				LDAPDynamicRoles: 8,
			},
		},
		{
			description:     "Replicated storage, previous month",
			localPathPrefix: billing.ReplicatedPrefix,
			monthOffset:     -1,
			roleCounts: &RoleCounts{
				DatabaseDynamicRoles: 5,
				DatabaseStaticRoles:  7,
				GCPRolesets:          3,
			},
		},
		{
			description:     "Replicated storage, current month",
			localPathPrefix: billing.ReplicatedPrefix,
			monthOffset:     0,
			roleCounts: &RoleCounts{
				AWSDynamicRoles:  12,
				AWSStaticRoles:   18,
				LDAPDynamicRoles: 6,
				GCPRolesets:      4,
			},
		},
		{
			description:     "Local storage, previous month with 4 role counts",
			localPathPrefix: billing.LocalPrefix,
			monthOffset:     -1,
			roleCounts: &RoleCounts{
				AWSDynamicRoles:  8,
				AWSStaticRoles:   10,
				LDAPDynamicRoles: 5,
				GCPRolesets:      2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			month := time.Now().AddDate(0, tc.monthOffset, 0)

			err := core.storeMaxRoleCountsLocked(context.Background(), tc.roleCounts, tc.localPathPrefix, month)
			require.NoError(t, err)

			retrievedCounts, err := core.GetStoredHWMRoleCounts(context.Background(), tc.localPathPrefix, month)
			require.NoError(t, err)

			require.Equal(t, tc.roleCounts.AWSDynamicRoles, retrievedCounts.AWSDynamicRoles)
			require.Equal(t, tc.roleCounts.AWSStaticRoles, retrievedCounts.AWSStaticRoles)
			require.Equal(t, tc.roleCounts.AzureDynamicRoles, retrievedCounts.AzureDynamicRoles)
			require.Equal(t, tc.roleCounts.AzureStaticRoles, retrievedCounts.AzureStaticRoles)
			require.Equal(t, tc.roleCounts.GCPStaticAccounts, retrievedCounts.GCPStaticAccounts)
			require.Equal(t, tc.roleCounts.GCPImpersonatedAccounts, retrievedCounts.GCPImpersonatedAccounts)
			require.Equal(t, tc.roleCounts.OpenLDAPDynamicRoles, retrievedCounts.OpenLDAPDynamicRoles)
			require.Equal(t, tc.roleCounts.OpenLDAPStaticRoles, retrievedCounts.OpenLDAPStaticRoles)
			require.Equal(t, tc.roleCounts.LDAPDynamicRoles, retrievedCounts.LDAPDynamicRoles)
			require.Equal(t, tc.roleCounts.LDAPStaticRoles, retrievedCounts.LDAPStaticRoles)
			require.Equal(t, tc.roleCounts.DatabaseDynamicRoles, retrievedCounts.DatabaseDynamicRoles)
			require.Equal(t, tc.roleCounts.DatabaseStaticRoles, retrievedCounts.DatabaseStaticRoles)
			require.Equal(t, tc.roleCounts.GCPRolesets, retrievedCounts.GCPRolesets)
			require.Equal(t, tc.roleCounts.AlicloudDynamicRoles, retrievedCounts.AlicloudDynamicRoles)
			require.Equal(t, tc.roleCounts.RabbitMQDynamicRoles, retrievedCounts.RabbitMQDynamicRoles)
			require.Equal(t, tc.roleCounts.ConsulDynamicRoles, retrievedCounts.ConsulDynamicRoles)
			require.Equal(t, tc.roleCounts.NomadDynamicRoles, retrievedCounts.NomadDynamicRoles)
			require.Equal(t, tc.roleCounts.KubernetesDynamicRoles, retrievedCounts.KubernetesDynamicRoles)
			require.Equal(t, tc.roleCounts.MongoDBAtlasDynamicRoles, retrievedCounts.MongoDBAtlasDynamicRoles)
			require.Equal(t, tc.roleCounts.TerraformCloudDynamicRoles, retrievedCounts.TerraformCloudDynamicRoles)
		})
	}
}

// TestHWMRoleCounts tests that we correctly store and track the HWM role counts
func TestHWMRoleCounts(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
	}
	core, _, root := TestCoreUnsealedWithConfig(t, coreConfig)

	for mount := range roleLogicalBackends {
		req := logical.TestRequest(t, logical.CreateOperation, fmt.Sprintf("sys/mounts/%v", mount))
		req.Data["type"] = mount
		req.ClientToken = root
		ctx := namespace.RootContext(context.Background())

		resp, err := core.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.Nil(t, resp.Error())
	}

	testCases := map[string]struct {
		mount        string
		key          string
		numberOfKeys int
	}{
		"AWS Dynamic Roles": {
			mount:        pluginconsts.SecretEngineAWS,
			key:          "role/",
			numberOfKeys: 5,
		},
		"AWS Static Roles": {
			mount:        pluginconsts.SecretEngineAWS,
			key:          "static-roles/",
			numberOfKeys: 5,
		},
		"Azure Dynamic Roles": {
			mount:        pluginconsts.SecretEngineAzure,
			key:          "roles/",
			numberOfKeys: 5,
		},
		"Azure Static Roles": {
			mount:        pluginconsts.SecretEngineAzure,
			key:          "static-roles/",
			numberOfKeys: 5,
		},
		"Database Dynamic Roles": {
			mount:        pluginconsts.SecretEngineDatabase,
			key:          "role/",
			numberOfKeys: 5,
		},
		"Database Static Roles": {
			mount:        pluginconsts.SecretEngineDatabase,
			key:          "static-role/",
			numberOfKeys: 5,
		},
		"GCP Impersonated Accounts": {
			mount:        pluginconsts.SecretEngineGCP,
			key:          "impersonated-account/",
			numberOfKeys: 5,
		},
		"GCP Rolesets": {
			mount:        pluginconsts.SecretEngineGCP,
			key:          "roleset/",
			numberOfKeys: 5,
		},
		"GCP Static Accounts": {
			mount:        pluginconsts.SecretEngineGCP,
			key:          "static-account/",
			numberOfKeys: 5,
		},
		"LDAP Dynamic Roles": {
			mount:        pluginconsts.SecretEngineLDAP,
			key:          "role/",
			numberOfKeys: 5,
		},
		"LDAP Static Roles": {
			mount:        pluginconsts.SecretEngineLDAP,
			key:          "static-role/",
			numberOfKeys: 5,
		},
		"OpenLDAP Dynamic Roles": {
			mount:        pluginconsts.SecretEngineOpenLDAP,
			key:          "role/",
			numberOfKeys: 5,
		},
		"OpenLDAP Static Roles": {
			mount:        pluginconsts.SecretEngineOpenLDAP,
			key:          "static-role/",
			numberOfKeys: 5,
		},
		"Alicloud Dynamic Roles": {
			mount:        pluginconsts.SecretEngineAlicloud,
			key:          "role/",
			numberOfKeys: 5,
		},
		"RabbitMQ Dynamic Roles": {
			mount:        pluginconsts.SecretEngineRabbitMQ,
			key:          "role/",
			numberOfKeys: 5,
		},
		"Consul Dynamic Roles": {
			mount:        pluginconsts.SecretEngineConsul,
			key:          "policy/",
			numberOfKeys: 5,
		},
		"Nomad Dynamic Roles": {
			mount:        pluginconsts.SecretEngineNomad,
			key:          "role/",
			numberOfKeys: 5,
		},
		"Kubernetes Dynamic Roles": {
			mount:        pluginconsts.SecretEngineKubernetes,
			key:          "roles/",
			numberOfKeys: 5,
		},
		"MongoDB Atlas Dynamic Roles": {
			mount:        pluginconsts.SecretEngineMongoDBAtlas,
			key:          "roles/",
			numberOfKeys: 5,
		},
		"Terraform Cloud Dynamic Roles": {
			mount:        pluginconsts.SecretEngineTerraform,
			key:          "role/",
			numberOfKeys: 5,
		},
	}

	// Sleep to prevent race conditions during the role initialization
	time.Sleep(1 * time.Second)

	core.mountsLock.RLock()
	defer core.mountsLock.RUnlock()
	for _, tc := range testCases {
		addRoleToStorage(t, core, tc.mount, tc.key, tc.numberOfKeys)
	}

	firstCounts := core.GetRoleCounts()
	verifyExpectedRoleCounts(t, firstCounts, 5)

	counts, err := core.UpdateMaxRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	verifyExpectedRoleCounts(t, counts, 5)

	counts, err = core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)
	// Verify that the max role counts are as expected
	verifyExpectedRoleCounts(t, counts, 5)

	// Reduce the number of roles. The max counts should remain the same
	for _, tc := range testCases {
		deleteAllRolesFromStorage(t, core, tc.mount, tc.key)
		addRoleToStorage(t, core, tc.mount, tc.key, 2)
	}

	counts, err = core.UpdateMaxRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	verifyExpectedRoleCounts(t, counts, 5)

	counts, err = core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	verifyExpectedRoleCounts(t, counts, 5)

	// Increase the number of roles. The max counts should update
	for _, tc := range testCases {
		deleteAllRolesFromStorage(t, core, tc.mount, tc.key)
		addRoleToStorage(t, core, tc.mount, tc.key, 8)
	}

	counts, err = core.UpdateMaxRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	verifyExpectedRoleCounts(t, counts, 8)

	counts, err = core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	verifyExpectedRoleCounts(t, counts, 8)

	// Decrease the number of roles back to 5. The max counts should remain at 8
	for _, tc := range testCases {
		deleteAllRolesFromStorage(t, core, tc.mount, tc.key)
		addRoleToStorage(t, core, tc.mount, tc.key, 5)
	}

	counts, err = core.UpdateMaxRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	verifyExpectedRoleCounts(t, counts, 8)

	counts, err = core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)
	verifyExpectedRoleCounts(t, counts, 8)
}

// TestHWMKvSecretsCounts tests that we correctly store and track the HWM kv counts
// for both kv-v1 and kv-v2 mounts.
func TestHWMKvSecretsCounts(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
		BillingConfig: billing.BillingConfig{
			MetricsUpdateCadence: 3 * time.Second,
		},
	}
	core, _, root := TestCoreUnsealedWithConfig(t, coreConfig)

	// Add 1 kv-v1 mount and 1 kv-v2 mount in the root namespace
	for _, mount := range []string{"kv-v1", "kv-v2"} {
		req := logical.TestRequest(t, logical.CreateOperation, fmt.Sprintf("sys/mounts/%v", mount))
		req.Data["type"] = mount
		req.ClientToken = root
		ctx := namespace.RootContext(context.Background())

		_, err := core.HandleRequest(ctx, req)
		require.NoError(t, err)
	}

	// Add two secrets to each mount
	for _, mount := range []string{"kv-v1", "kv-v2"} {
		for i := 0; i < 2; i++ {
			secretName := fmt.Sprintf("secret-%d", i)
			addKvSecretToStorage(t, namespace.RootContext(context.Background()), core, mount, root, secretName, mount)
		}
	}

	// Verify that the max kv counts are as expected
	timer := time.NewTimer(3 * time.Second)
	_ = <-timer.C
	counts, err := core.GetStoredHWMKvCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)
	require.Equal(t, 4, counts)

	// Add one more secret to the kv-v1 mount
	addKvSecretToStorage(t, namespace.RootContext(context.Background()), core, "kv-v1", root, "secret-3", "kv-v1")

	// Wait for the metrics update
	timer = time.NewTimer(3 * time.Second)
	_ = <-timer.C

	// Verify that the max kv counts are updated
	counts, err = core.GetStoredHWMKvCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)
	require.Equal(t, 5, counts)

	// Now delete one secret from the kv-v2 mount
	deleteKvSecretFromStorage(t, namespace.RootContext(context.Background()), core, "kv-v2", root, "secret-1", "kv-v2")

	// Wait for any metrics updates to complete
	timer = time.NewTimer(3 * time.Second)
	_ = <-timer.C

	// Verify that the max kv counts are still the same
	counts, err = core.GetStoredHWMKvCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)
	require.Equal(t, 5, counts)
}

// TestDataProtectionCallCounts tests that we correctly store and track the data protection call counts
func TestDataProtectionCallCounts(t *testing.T) {
	t.Parallel()
	coreConfig := &CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
		BillingConfig: billing.BillingConfig{
			MetricsUpdateCadence: 3 * time.Second,
		},
	}

	core, _, root := TestCoreUnsealedWithConfig(t, coreConfig)

	// Mount transit backend
	req := logical.TestRequest(t, logical.CreateOperation, "sys/mounts/transit")
	req.Data["type"] = "transit"
	req.ClientToken = root
	ctx := namespace.RootContext(context.Background())
	_, err := core.HandleRequest(ctx, req)
	require.NoError(t, err)

	// Reset the transit counters
	billing.CurrentDataProtectionCallCounts.Transit = 0

	// Create an encryption key
	req = logical.TestRequest(t, logical.CreateOperation, "transit/keys/foo")
	req.Data["type"] = "aes256-gcm96"
	req.ClientToken = root
	_, err = core.HandleRequest(ctx, req)
	require.NoError(t, err)

	// Perform encryption on the key
	req = logical.TestRequest(t, logical.UpdateOperation, "transit/encrypt/foo")
	req.Data["plaintext"] = "dGhlIHF1aWNrIGJyb3duIGZveA=="
	req.ClientToken = root
	resp, err := core.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)

	// Verify that the transit counter is incremented (replicated mount by default)
	require.Equal(t, int64(1), billing.CurrentDataProtectionCallCounts.Transit)

	// Get the ciphertext from the encryption response
	ciphertext, ok := resp.Data["ciphertext"].(string)
	require.True(t, ok)
	require.NotEmpty(t, ciphertext)

	// Now perform decryption using the ciphertext
	req = logical.TestRequest(t, logical.UpdateOperation, "transit/decrypt/foo")
	req.Data["ciphertext"] = ciphertext
	req.ClientToken = root
	_, err = core.HandleRequest(ctx, req)
	require.NoError(t, err)

	// Verify that the transit counter is incremented
	require.Equal(t, int64(2), billing.CurrentDataProtectionCallCounts.Transit)

	// Test rewrap operation
	req = logical.TestRequest(t, logical.UpdateOperation, "transit/rewrap/foo")
	req.Data["ciphertext"] = ciphertext
	req.ClientToken = root
	resp, err = core.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)

	// Verify that the transit counter is incremented
	require.Equal(t, int64(3), billing.CurrentDataProtectionCallCounts.Transit)

	// Get the new ciphertext from rewrap
	newCiphertext, ok := resp.Data["ciphertext"].(string)
	require.True(t, ok)
	require.NotEmpty(t, newCiphertext)

	// Test datakey generation
	req = logical.TestRequest(t, logical.UpdateOperation, "transit/datakey/plaintext/foo")
	req.ClientToken = root
	resp, err = core.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)

	// Verify that the transit counter is incremented
	require.Equal(t, int64(4), billing.CurrentDataProtectionCallCounts.Transit)

	// Test HMAC generation
	req = logical.TestRequest(t, logical.UpdateOperation, "transit/hmac/foo")
	req.Data["input"] = "dGhlIHF1aWNrIGJyb3duIGZveA=="
	req.ClientToken = root
	resp, err = core.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)

	// Verify that the transit counter is incremented
	require.Equal(t, int64(5), billing.CurrentDataProtectionCallCounts.Transit)

	// Get the HMAC value
	hmacValue, ok := resp.Data["hmac"].(string)
	require.True(t, ok)
	require.NotEmpty(t, hmacValue)

	// Test HMAC verification
	req = logical.TestRequest(t, logical.UpdateOperation, "transit/verify/foo")
	req.Data["input"] = "dGhlIHF1aWNrIGJyb3duIGZveA=="
	req.Data["hmac"] = hmacValue
	req.ClientToken = root
	resp, err = core.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)

	// Verify that the transit counter is incremented
	require.Equal(t, int64(6), billing.CurrentDataProtectionCallCounts.Transit)

	// Verify the HMAC is valid
	hmacValid, ok := resp.Data["valid"].(bool)
	require.True(t, ok)
	require.True(t, hmacValid)

	// Create a signing key for sign/verify operations
	req = logical.TestRequest(t, logical.CreateOperation, "transit/keys/signing-key")
	req.Data["type"] = "ecdsa-p256"
	req.ClientToken = root
	_, err = core.HandleRequest(ctx, req)
	require.NoError(t, err)

	// Test sign operation
	req = logical.TestRequest(t, logical.UpdateOperation, "transit/sign/signing-key")
	req.Data["input"] = "dGhlIHF1aWNrIGJyb3duIGZveA=="
	req.ClientToken = root
	resp, err = core.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)

	// Verify that the transit counter is incremented
	require.Equal(t, int64(7), billing.CurrentDataProtectionCallCounts.Transit)

	// Get the signature
	signature, ok := resp.Data["signature"].(string)
	require.True(t, ok)
	require.NotEmpty(t, signature)

	// Test verify operation
	req = logical.TestRequest(t, logical.UpdateOperation, "transit/verify/signing-key")
	req.Data["input"] = "dGhlIHF1aWNrIGJyb3duIGZveA=="
	req.Data["signature"] = signature
	req.ClientToken = root
	resp, err = core.HandleRequest(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)

	// Verify that the transit counter is incremented
	require.Equal(t, int64(8), billing.CurrentDataProtectionCallCounts.Transit)

	// Verify the signature is valid
	signatureValid, ok := resp.Data["valid"].(bool)
	require.True(t, ok)
	require.True(t, signatureValid)

	// Test CMAC operations (ENT only - will be no-op in OSS)
	currentCount := billing.CurrentDataProtectionCallCounts.Transit
	currentCount = testCMACOperations(t, core, ctx, root, currentCount)

	// Verify that the transit counter matches expected count
	require.Equal(t, currentCount, billing.CurrentDataProtectionCallCounts.Transit)

	// Now test persisting the summed counts - store and retrieve counts
	// First, update the data protection call counts (this will sum current counter with stored value)
	summedCounts, err := core.UpdateDataProtectionCallCounts(ctx, time.Now())
	require.NoError(t, err)
	require.NotNil(t, summedCounts)
	require.Equal(t, currentCount, summedCounts.Transit)

	// Verify the counter was reset after update
	require.Equal(t, int64(0), billing.CurrentDataProtectionCallCounts.Transit, "Counter should be reset after update")

	// Retrieve the stored counts
	storedCounts, err := core.GetStoredDataProtectionCallCounts(ctx, time.Now())
	require.NoError(t, err)
	require.NotNil(t, storedCounts)
	require.Equal(t, currentCount, storedCounts.Transit)

	// Perform more operations to increase the counter
	req = logical.TestRequest(t, logical.UpdateOperation, "transit/encrypt/foo")
	req.Data["plaintext"] = "dGhlIHF1aWNrIGJyb3duIGZveA=="
	req.ClientToken = root
	_, err = core.HandleRequest(ctx, req)
	require.NoError(t, err)

	// Counter should now be 1 (reset + 1 operation)
	require.Equal(t, int64(1), billing.CurrentDataProtectionCallCounts.Transit)

	// Update counts again - should sum the new count (1) with the stored count (currentCount)
	summedCounts, err = core.UpdateDataProtectionCallCounts(ctx, time.Now())
	require.NoError(t, err)
	expectedSum := currentCount + 1
	require.Equal(t, expectedSum, summedCounts.Transit, "Count should be sum of stored and current")

	// Verify the counter was reset after update
	require.Equal(t, int64(0), billing.CurrentDataProtectionCallCounts.Transit, "Counter should be reset after update")

	// Verify stored counts are now the sum
	storedCounts, err = core.GetStoredDataProtectionCallCounts(ctx, time.Now())
	require.NoError(t, err)
	require.Equal(t, expectedSum, storedCounts.Transit)

	// Add more operations without manually resetting
	for i := 0; i < 3; i++ {
		req = logical.TestRequest(t, logical.UpdateOperation, "transit/encrypt/foo")
		req.Data["plaintext"] = "dGhlIHF1aWNrIGJyb3duIGZveA=="
		req.ClientToken = root
		_, err = core.HandleRequest(ctx, req)
		require.NoError(t, err)
	}

	// Counter should be 3
	require.Equal(t, int64(3), billing.CurrentDataProtectionCallCounts.Transit)

	// Update counts - should sum 3 with the previous stored sum
	summedCounts, err = core.UpdateDataProtectionCallCounts(ctx, time.Now())
	require.NoError(t, err)
	expectedSum = expectedSum + 3
	require.Equal(t, expectedSum, summedCounts.Transit, "Count should continue to sum")

	// Verify the counter was reset after update
	require.Equal(t, int64(0), billing.CurrentDataProtectionCallCounts.Transit, "Counter should be reset after update")

	// Verify stored counts
	storedCounts, err = core.GetStoredDataProtectionCallCounts(ctx, time.Now())
	require.NoError(t, err)
	require.Equal(t, expectedSum, storedCounts.Transit)

	// Update again without any new operations
	// This verifies we don't double-count
	summedCounts, err = core.UpdateDataProtectionCallCounts(ctx, time.Now())
	require.NoError(t, err)
	require.Equal(t, expectedSum, summedCounts.Transit, "Count should remain the same when no new operations occurred")

	// Verify stored counts haven't changed
	storedCounts, err = core.GetStoredDataProtectionCallCounts(ctx, time.Now())
	require.NoError(t, err)
	require.Equal(t, expectedSum, storedCounts.Transit, "Stored count should remain the same")

	// Verify counter is still at 0
	require.Equal(t, int64(0), billing.CurrentDataProtectionCallCounts.Transit, "Counter should still be 0")
}

func addRoleToStorage(t *testing.T, core *Core, mount string, key string, numberOfKeys int) {
	raw, ok := core.router.root.Get(mount + "/")
	if !ok {
		return
	}
	require.NotNil(t, raw)
	re := raw.(*routeEntry)
	storageView := re.storageView

	// Write to storage simulating adding a static, dynamic role, impersonated account or roleset
	// This bypasses the API to add data
	for i := 0; i < numberOfKeys; i++ {
		roleKey := fmt.Sprintf("%srole-%d", key, i)
		// Create a role with a unique key
		err := storageView.Put(context.Background(), &logical.StorageEntry{
			Key:   roleKey,
			Value: []byte("foo"),
		})
		require.NoError(t, err)
	}
	// Verify that the role is stored
	list, err := storageView.List(context.Background(), key)
	require.NoError(t, err)
	require.Len(t, list, numberOfKeys)
}

func deleteAllRolesFromStorage(t *testing.T, core *Core, mount string, key string) {
	raw, ok := core.router.root.Get(mount + "/")
	if !ok {
		return
	}
	re := raw.(*routeEntry)
	storageView := re.storageView

	// List all roles
	list, err := storageView.List(context.Background(), key)
	require.NoError(t, err)

	// Delete each role
	for _, role := range list {
		err := storageView.Delete(context.Background(), fmt.Sprintf("%s%s", key, role))
		require.NoError(t, err)
	}

	// Verify that all roles are deleted
	list, err = storageView.List(context.Background(), key)
	require.NoError(t, err)
	require.Len(t, list, 0)
}

func addKvSecretToStorage(t *testing.T, ctx context.Context, core *Core, mount string, token string, secretName string, kvVersion string) {
	var req *logical.Request
	switch kvVersion {
	case "kv-v2":
		// KV v2 expects writes to /data/<path> with a nested "data" payload
		req = logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("%v/data/%s", mount, secretName))
		req.Data["data"] = map[string]interface{}{
			"foo": "bar",
		}
	case "kv-v1":
		// KV v1 expects writes directly to /<path> with a flat payload
		req = logical.TestRequest(t, logical.UpdateOperation, fmt.Sprintf("%v/%s", mount, secretName))
		req.Data["foo"] = "bar"
	default:
		t.Fatalf("invalid kv version: %s", kvVersion)
	}
	req.ClientToken = token
	_, err := core.HandleRequest(ctx, req)
	require.NoError(t, err)
}

func deleteKvSecretFromStorage(t *testing.T, ctx context.Context, core *Core, mount string, token string, secretName string, kvVersion string) {
	var req *logical.Request
	switch kvVersion {
	case "kv-v2":
		req = logical.TestRequest(t, logical.DeleteOperation, fmt.Sprintf("%v/data/%s", mount, secretName))
	case "kv-v1":
		req = logical.TestRequest(t, logical.DeleteOperation, fmt.Sprintf("%v/%s", mount, secretName))
	default:
		t.Fatalf("invalid kv version: %s", kvVersion)
	}
	req.ClientToken = token
	_, err := core.HandleRequest(ctx, req)
	require.NoError(t, err)
}
