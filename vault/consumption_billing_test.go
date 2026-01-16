// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/stretchr/testify/require"
)

var secretEngineBackends = map[string]struct {
	mount string
	key   string
}{
	"AWS Dynamic Roles": {
		mount: pluginconsts.SecretEngineAWS,
		key:   "role/",
	},
	"AWS Static Roles": {
		mount: pluginconsts.SecretEngineAWS,
		key:   "static-roles/",
	},
	"Azure Dynamic Roles": {
		mount: pluginconsts.SecretEngineAzure,
		key:   "roles/",
	},
	"Database Dynamic Roles": {
		mount: pluginconsts.SecretEngineDatabase,
		key:   "role/",
	},
	"Database Static Roles": {
		mount: pluginconsts.SecretEngineDatabase,
		key:   "static-role/",
	},
	"GCP Impersonated Accounts": {
		mount: pluginconsts.SecretEngineGCP,
		key:   "impersonated-account/",
	},
	"GCP Rolesets": {
		mount: pluginconsts.SecretEngineGCP,
		key:   "roleset/",
	},
	"GCP Static Accounts": {
		mount: pluginconsts.SecretEngineGCP,
		key:   "static-account/",
	},
	"LDAP Dynamic Roles": {
		mount: pluginconsts.SecretEngineLDAP,
		key:   "role/",
	},
	"LDAP Static Roles": {
		mount: pluginconsts.SecretEngineLDAP,
		key:   "static-role/",
	},
	"OpenLDAP Dynamic Roles": {
		mount: pluginconsts.SecretEngineOpenLDAP,
		key:   "role/",
	},
	"OpenLDAP Static Roles": {
		mount: pluginconsts.SecretEngineOpenLDAP,
		key:   "static-role/",
	},
}

// TestConsumptionBillingMetricsWorker tests that we correctly update the consumption metrics at
// regular intervals
func TestConsumptionBillingMetricsWorker(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: roleLogicalBackends,
		BillingConfig: billing.BillingConfig{
			MetricsUpdateCadence: 5 * time.Second,
		},
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

	for _, tc := range secretEngineBackends {
		addRoleToStorage(t, core, tc.mount, tc.key, 5)
	}
	timer := time.NewTimer(5 * time.Second)
	expectedCounts := RoleCounts{
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
	}

	_ = <-timer.C
	// Check that the billing metrics have been updated
	counts, err := core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	require.Equal(t, *counts, expectedCounts)

	for _, tc := range secretEngineBackends {
		deleteAllRolesFromStorage(t, core, tc.mount, tc.key)
		addRoleToStorage(t, core, tc.mount, tc.key, 3)
	}

	timer.Reset(5 * time.Second)

	_ = <-timer.C
	// Check that the billing metrics have been updated
	counts, err = core.GetStoredHWMRoleCounts(context.Background(), billing.ReplicatedPrefix, time.Now())
	require.NoError(t, err)

	require.Equal(t, *counts, expectedCounts)
}
