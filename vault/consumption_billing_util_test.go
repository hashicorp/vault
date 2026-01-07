// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	logicalAzure "github.com/hashicorp/vault-plugin-secrets-azure"
	logicalGcp "github.com/hashicorp/vault-plugin-secrets-gcp/plugin"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	logicalLDAP "github.com/hashicorp/vault-plugin-secrets-openldap"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	logicalAws "github.com/hashicorp/vault/builtin/logical/aws"
	logicalDatabase "github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

var roleLogicalBackends = map[string]logical.Factory{
	pluginconsts.SecretEngineAWS:      logicalAws.Factory,
	pluginconsts.SecretEngineAzure:    logicalAzure.Factory,
	pluginconsts.SecretEngineGCP:      logicalGcp.Factory,
	pluginconsts.SecretEngineKV:       logicalKv.Factory,
	pluginconsts.SecretEngineLDAP:     logicalLDAP.Factory,
	pluginconsts.SecretEngineDatabase: logicalDatabase.Factory,
	pluginconsts.SecretEngineOpenLDAP: logicalLDAP.Factory,
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
			require.Equal(t, tc.roleCounts.GCPStaticAccounts, retrievedCounts.GCPStaticAccounts)
			require.Equal(t, tc.roleCounts.GCPImpersonatedAccounts, retrievedCounts.GCPImpersonatedAccounts)
			require.Equal(t, tc.roleCounts.OpenLDAPDynamicRoles, retrievedCounts.OpenLDAPDynamicRoles)
			require.Equal(t, tc.roleCounts.OpenLDAPStaticRoles, retrievedCounts.OpenLDAPStaticRoles)
			require.Equal(t, tc.roleCounts.LDAPDynamicRoles, retrievedCounts.LDAPDynamicRoles)
			require.Equal(t, tc.roleCounts.LDAPStaticRoles, retrievedCounts.LDAPStaticRoles)
			require.Equal(t, tc.roleCounts.DatabaseDynamicRoles, retrievedCounts.DatabaseDynamicRoles)
			require.Equal(t, tc.roleCounts.DatabaseStaticRoles, retrievedCounts.DatabaseStaticRoles)
			require.Equal(t, tc.roleCounts.GCPRolesets, retrievedCounts.GCPRolesets)
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
	}

	core.mountsLock.RLock()
	defer core.mountsLock.RUnlock()
	for _, tc := range testCases {
		addRoleToStorage(t, core, tc.mount, tc.key, tc.numberOfKeys)
	}

	firstCounts := core.GetRoleCounts()
	require.Equal(t, &RoleCounts{
		AWSDynamicRoles:         5,
		AWSStaticRoles:          5,
		AzureDynamicRoles:       5,
		DatabaseDynamicRoles:    5,
		DatabaseStaticRoles:     5,
		GCPImpersonatedAccounts: 5,
		GCPRolesets:             5,
		GCPStaticAccounts:       5,
		LDAPDynamicRoles:        5,
		LDAPStaticRoles:         5,
		OpenLDAPDynamicRoles:    5,
		OpenLDAPStaticRoles:     5,
	}, firstCounts)

	counts, err := core.UpdateMaxRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	require.Equal(t, &RoleCounts{
		AWSDynamicRoles:         5,
		AWSStaticRoles:          5,
		AzureDynamicRoles:       5,
		DatabaseDynamicRoles:    5,
		DatabaseStaticRoles:     5,
		GCPImpersonatedAccounts: 5,
		GCPRolesets:             5,
		GCPStaticAccounts:       5,
		LDAPDynamicRoles:        5,
		LDAPStaticRoles:         5,
		OpenLDAPDynamicRoles:    5,
		OpenLDAPStaticRoles:     5,
	}, counts)

	counts, err = core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)
	// Verify that the max role counts are as expected
	require.Equal(t, &RoleCounts{
		AWSDynamicRoles:         5,
		AWSStaticRoles:          5,
		AzureDynamicRoles:       5,
		DatabaseDynamicRoles:    5,
		DatabaseStaticRoles:     5,
		GCPImpersonatedAccounts: 5,
		GCPRolesets:             5,
		GCPStaticAccounts:       5,
		LDAPDynamicRoles:        5,
		LDAPStaticRoles:         5,
		OpenLDAPDynamicRoles:    5,
		OpenLDAPStaticRoles:     5,
	}, counts)

	// Reduce the number of roles. The max counts should remain the same
	for _, tc := range testCases {
		deleteAllRolesFromStorage(t, core, tc.mount, tc.key)
		addRoleToStorage(t, core, tc.mount, tc.key, 2)
	}

	counts, err = core.UpdateMaxRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	require.Equal(t, &RoleCounts{
		AWSDynamicRoles:         5,
		AWSStaticRoles:          5,
		AzureDynamicRoles:       5,
		DatabaseDynamicRoles:    5,
		DatabaseStaticRoles:     5,
		GCPImpersonatedAccounts: 5,
		GCPRolesets:             5,
		GCPStaticAccounts:       5,
		LDAPDynamicRoles:        5,
		LDAPStaticRoles:         5,
		OpenLDAPDynamicRoles:    5,
		OpenLDAPStaticRoles:     5,
	}, counts)

	counts, err = core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	require.Equal(t, &RoleCounts{
		AWSDynamicRoles:         5,
		AWSStaticRoles:          5,
		AzureDynamicRoles:       5,
		DatabaseDynamicRoles:    5,
		DatabaseStaticRoles:     5,
		GCPImpersonatedAccounts: 5,
		GCPRolesets:             5,
		GCPStaticAccounts:       5,
		LDAPDynamicRoles:        5,
		LDAPStaticRoles:         5,
		OpenLDAPDynamicRoles:    5,
		OpenLDAPStaticRoles:     5,
	}, counts)

	// Increase the number of roles. The max counts should update
	for _, tc := range testCases {
		deleteAllRolesFromStorage(t, core, tc.mount, tc.key)
		addRoleToStorage(t, core, tc.mount, tc.key, 8)
	}

	counts, err = core.UpdateMaxRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	require.Equal(t, &RoleCounts{
		AWSDynamicRoles:         8,
		AWSStaticRoles:          8,
		AzureDynamicRoles:       8,
		DatabaseDynamicRoles:    8,
		DatabaseStaticRoles:     8,
		GCPImpersonatedAccounts: 8,
		GCPRolesets:             8,
		GCPStaticAccounts:       8,
		LDAPDynamicRoles:        8,
		LDAPStaticRoles:         8,
		OpenLDAPDynamicRoles:    8,
		OpenLDAPStaticRoles:     8,
	}, counts)

	counts, err = core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	require.Equal(t, &RoleCounts{
		AWSDynamicRoles:         8,
		AWSStaticRoles:          8,
		AzureDynamicRoles:       8,
		DatabaseDynamicRoles:    8,
		DatabaseStaticRoles:     8,
		GCPImpersonatedAccounts: 8,
		GCPRolesets:             8,
		GCPStaticAccounts:       8,
		LDAPDynamicRoles:        8,
		LDAPStaticRoles:         8,
		OpenLDAPDynamicRoles:    8,
		OpenLDAPStaticRoles:     8,
	}, counts)

	// Decrease the number of roles back to 5. The max counts should remain at 8
	for _, tc := range testCases {
		deleteAllRolesFromStorage(t, core, tc.mount, tc.key)
		addRoleToStorage(t, core, tc.mount, tc.key, 5)
	}

	counts, err = core.UpdateMaxRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	require.Equal(t, &RoleCounts{
		AWSDynamicRoles:         8,
		AWSStaticRoles:          8,
		AzureDynamicRoles:       8,
		DatabaseDynamicRoles:    8,
		DatabaseStaticRoles:     8,
		GCPImpersonatedAccounts: 8,
		GCPRolesets:             8,
		GCPStaticAccounts:       8,
		LDAPDynamicRoles:        8,
		LDAPStaticRoles:         8,
		OpenLDAPDynamicRoles:    8,
		OpenLDAPStaticRoles:     8,
	}, counts)

	counts, err = core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)
	require.Equal(t, &RoleCounts{
		AWSDynamicRoles:         8,
		AWSStaticRoles:          8,
		AzureDynamicRoles:       8,
		DatabaseDynamicRoles:    8,
		DatabaseStaticRoles:     8,
		GCPImpersonatedAccounts: 8,
		GCPRolesets:             8,
		GCPStaticAccounts:       8,
		LDAPDynamicRoles:        8,
		LDAPStaticRoles:         8,
		OpenLDAPDynamicRoles:    8,
		OpenLDAPStaticRoles:     8,
	}, counts)
}

func addRoleToStorage(t *testing.T, core *Core, mount string, key string, numberOfKeys int) {
	raw, ok := core.router.root.Get(mount + "/")
	if !ok {
		return
	}
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
