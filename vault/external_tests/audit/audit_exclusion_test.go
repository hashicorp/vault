// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"testing"

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestAudit_Exclusion_ByVaultVersion ensures that the audit device 'exclude'
// option is only supported in the enterprise edition of the product.
func TestAudit_Exclusion_ByVaultVersion(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Attempt to create an audit device with exclusion enabled.
	mountPointFilterDevicePath := "mountpoint"
	mountPointFilterDeviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": "discard",
			"exclude":   "[ { \"fields\": [ \"/response/data\" ] } ]",
		},
	}

	_, err := client.Logical().Write("sys/audit/"+mountPointFilterDevicePath, mountPointFilterDeviceData)
	if constants.IsEnterprise {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
		require.ErrorContains(t, err, "enterprise-only options supplied")
	}

	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	if constants.IsEnterprise {
		require.Len(t, devices, 1)
	} else {
		// Ensure the device has not been created.
		require.Len(t, devices, 0)
	}
}
