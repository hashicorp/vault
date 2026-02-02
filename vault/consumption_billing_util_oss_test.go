// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// verifyExpectedRoleCounts verifies that the actual role counts match expected values.
// In OSS, AzureStaticRoles should be 0 since they're only supported in Enterprise.
func verifyExpectedRoleCounts(t *testing.T, actual *RoleCounts, baseCount int) {
	expected := &RoleCounts{
		AWSDynamicRoles:            baseCount,
		AWSStaticRoles:             baseCount,
		AzureDynamicRoles:          baseCount,
		AzureStaticRoles:           0, // OSS: Azure Static roles not supported
		DatabaseDynamicRoles:       baseCount,
		DatabaseStaticRoles:        baseCount,
		GCPImpersonatedAccounts:    baseCount,
		GCPRolesets:                baseCount,
		GCPStaticAccounts:          baseCount,
		LDAPDynamicRoles:           baseCount,
		LDAPStaticRoles:            baseCount,
		OpenLDAPDynamicRoles:       baseCount,
		OpenLDAPStaticRoles:        baseCount,
		AlicloudDynamicRoles:       baseCount,
		RabbitMQDynamicRoles:       baseCount,
		ConsulDynamicRoles:         baseCount,
		NomadDynamicRoles:          baseCount,
		KubernetesDynamicRoles:     baseCount,
		MongoDBAtlasDynamicRoles:   baseCount,
		TerraformCloudDynamicRoles: baseCount,
	}
	require.Equal(t, expected, actual)
}

// testCMACOperations is a no-op in OSS since CMAC is an Enterprise-only feature.
// Returns the current count unchanged.
func testCMACOperations(t *testing.T, core *Core, ctx context.Context, root string, currentCount int64) int64 {
	// CMAC is not supported in OSS, so we don't perform any operations
	return currentCount
}
