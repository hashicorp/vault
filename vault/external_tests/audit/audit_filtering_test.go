// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestAuditFilteringMultipleDevices validates that the audit device 'filter'
// option works as expected and multiple audit devices with the same filter all
// write the relevant entries to the logs. We create two audit devices that
// filter out all events that are not for the KV mount type and one without
// filters, make some auditable requests that both match and do not match the
// filters, and ensure there are audit log entries for the former but not the
// latter.
func TestAuditFilteringMultipleDevices(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client, err := cluster.Cores[0].Client.Clone()
	require.NoError(t, err)
	client.SetToken(cluster.RootToken)

	// Create audit devices.
	tempDir := t.TempDir()
	filteredLogFile, err := os.CreateTemp(tempDir, "")
	filteredDevicePath := "filtered"
	filteredDeviceData := map[string]interface{}{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]interface{}{
			"file_path": filteredLogFile.Name(),
			"filter":    "mount_type == kv",
		},
	}
	_, err = client.Logical().Write("sys/audit/"+filteredDevicePath, filteredDeviceData)
	require.NoError(t, err)

	filteredLogFile2, err := os.CreateTemp(tempDir, "")
	filteredDevicePath2 := "filtered2"
	filteredDeviceData2 := map[string]interface{}{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]interface{}{
			"file_path": filteredLogFile2.Name(),
			"filter":    "mount_type == kv",
		},
	}
	_, err = client.Logical().Write("sys/audit/"+filteredDevicePath2, filteredDeviceData2)
	require.NoError(t, err)

	nonFilteredLogFile, err := os.CreateTemp(tempDir, "")
	nonFilteredDevicePath := "nonfiltered"
	nonFilteredDeviceData := map[string]interface{}{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]interface{}{
			"file_path": nonFilteredLogFile.Name(),
		},
	}
	_, err = client.Logical().Write("sys/audit/"+nonFilteredDevicePath, nonFilteredDeviceData)
	require.NoError(t, err)

	// Ensure the non-filtered log file is not empty.
	nonFilteredLogSize := getFileSize(t, nonFilteredLogFile.Name())
	require.Positive(t, nonFilteredLogSize)

	// A write to KV should produce an audit entry that is written to the
	// filtered devices and the non-filtered device.
	data := map[string]interface{}{
		"foo": "bar",
	}
	err = client.KVv1("secret/").Put(context.Background(), "foo", data)
	require.NoError(t, err)
	// Ensure the non-filtered log file was written to.
	oldNonFilteredLogSize := nonFilteredLogSize
	nonFilteredLogSize = getFileSize(t, nonFilteredLogFile.Name())
	require.Greater(t, nonFilteredLogSize, oldNonFilteredLogSize)

	// Parse the filtered logs and verify that the filters only allowed entries
	// with mount_type value being 'kv'. While we're at it, ensure that the
	// numbers of entries are correct in both files.
	filteredLogFiles := []*os.File{filteredLogFile, filteredLogFile2}
	for _, f := range filteredLogFiles {
		counter := 0
		decoder := json.NewDecoder(f)
		var auditRecord map[string]interface{}
		for decoder.Decode(&auditRecord) == nil {
			auditRequest := map[string]interface{}{}
			if req, ok := auditRecord["request"]; ok {
				auditRequest = req.(map[string]interface{})
			} else {
				t.Fatal("failed to parse request data from audit log entry")
			}

			require.Equal(t, "kv", auditRequest["mount_type"])
			counter += 1
		}
		require.Equal(t, 2, counter)
	}

	// Disable the audit devices.
	err = client.Sys().DisableAudit(filteredDevicePath)
	require.NoError(t, err)
	err = client.Sys().DisableAudit(filteredDevicePath2)
	require.NoError(t, err)
	err = client.Sys().DisableAudit(nonFilteredDevicePath)
	require.NoError(t, err)
	// Ensure the devices are no longer there.
	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	_, ok := devices[filteredDevicePath]
	require.False(t, ok)
	_, ok = devices[filteredDevicePath2]
	require.False(t, ok)
	_, ok = devices[nonFilteredDevicePath]
	require.False(t, ok)
}

func getFileSize(t *testing.T, filePath string) int64 {
	t.Helper()
	fi, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}
	size := fi.Size()
	return size
}
