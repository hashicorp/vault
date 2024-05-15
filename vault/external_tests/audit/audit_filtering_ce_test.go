// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package audit

import (
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestAuditFilteringInCE ensures that the audit device 'filter'
// option is only supported in the enterprise edition of the product.
func TestAuditFilteringInCE(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Attempt to create an audit device with filtering enabled.
	mountPointFilterDevicePath := "mountpoint"
	mountPointFilterDeviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": "/tmp/audit.log",
			"filter":    "mount_point == secret/",
		},
	}
	_, err := client.Logical().Write("sys/audit/"+mountPointFilterDevicePath, mountPointFilterDeviceData)
	require.Error(t, err)
	require.ErrorContains(t, err, "enterprise-only options supplied")

	// Ensure the device has not been created.
	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	require.Len(t, devices, 0)
}

// TestAuditFilteringFallbackDeviceInCE validates that the audit device
// 'fallback' option is only available in the enterprise edition of the product.
func TestAuditFilteringFallbackDeviceInCE(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	fallbackDevicePath := "fallback"
	fallbackDeviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": "/tmp/audit.log",
			"fallback":  "true",
		},
	}
	_, err := client.Logical().Write("sys/audit/"+fallbackDevicePath, fallbackDeviceData)
	require.Error(t, err)
	require.ErrorContains(t, err, "enterprise-only options supplied")

	// Ensure the device has not been created.
	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	require.Len(t, devices, 0)
}
