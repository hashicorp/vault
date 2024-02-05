// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestAuditFilteringOnDifferentFields validates that the audit device 'filter'
// option works as expected for the fields we allow filtering on. We create
// three audit devices, each with a different filter, and make some auditable
// requests, then we ensure that correct entries were written to the respective
// log files. The mount_type and namespace filters are tested in other tests in
// this package.
func TestAuditFilteringOnDifferentFields(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Create audit devices.
	tempDir := t.TempDir()
	mountPointFilterLogFile, err := os.CreateTemp(tempDir, "")
	mountPointFilterDevicePath := "mountpoint"
	mountPointFilterDeviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": mountPointFilterLogFile.Name(),
			"filter":    "mount_point == secret/",
		},
	}
	_, err = client.Logical().Write("sys/audit/"+mountPointFilterDevicePath, mountPointFilterDeviceData)
	require.NoError(t, err)

	operationFilterLogFile, err := os.CreateTemp(tempDir, "")
	operationFilterPath := "operation"
	operationFilterData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": operationFilterLogFile.Name(),
			"filter":    "operation == create",
		},
	}
	_, err = client.Logical().Write("sys/audit/"+operationFilterPath, operationFilterData)
	require.NoError(t, err)

	pathFilterLogFile, err := os.CreateTemp(tempDir, "")
	pathFilterDevicePath := "path"
	pathFilterDeviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": pathFilterLogFile.Name(),
			"filter":    "path == secret/foo",
		},
	}
	_, err = client.Logical().Write("sys/audit/"+pathFilterDevicePath, pathFilterDeviceData)
	require.NoError(t, err)

	// Ensure the devices have been created.
	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	_, ok := devices[mountPointFilterDevicePath+"/"]
	require.True(t, ok)
	_, ok = devices[operationFilterPath+"/"]
	require.True(t, ok)
	_, ok = devices[pathFilterDevicePath+"/"]
	require.True(t, ok)

	// A write to KV should produce an audit entry that is written to all the
	// audit devices.
	data := map[string]any{
		"foo": "bar",
	}
	err = client.KVv1("secret/").Put(context.Background(), "foo", data)
	require.NoError(t, err)

	// Disable the audit devices.
	err = client.Sys().DisableAudit(mountPointFilterDevicePath)
	require.NoError(t, err)
	err = client.Sys().DisableAudit(operationFilterPath)
	require.NoError(t, err)
	err = client.Sys().DisableAudit(pathFilterDevicePath)
	require.NoError(t, err)
	// Ensure the devices are no longer there.
	devices, err = client.Sys().ListAudit()
	require.NoError(t, err)
	require.Len(t, devices, 0)

	// Validate that only the entries matching the filters were written to each log file.
	entries := checkAuditEntries(t, mountPointFilterLogFile, "mount_point", "secret/")
	require.Equal(t, 2, entries)
	entries = checkAuditEntries(t, operationFilterLogFile, "operation", "create")
	require.Equal(t, 2, entries)
	entries = checkAuditEntries(t, pathFilterLogFile, "path", "secret/foo")
	require.Equal(t, 2, entries)
}

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
	client := cluster.Cores[0].Client

	// Create audit devices.
	tempDir := t.TempDir()
	filteredLogFile, err := os.CreateTemp(tempDir, "")
	filteredDevicePath := "filtered"
	filteredDeviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": filteredLogFile.Name(),
			"filter":    "mount_type == kv",
		},
	}
	_, err = client.Logical().Write("sys/audit/"+filteredDevicePath, filteredDeviceData)
	require.NoError(t, err)

	filteredLogFile2, err := os.CreateTemp(tempDir, "")
	filteredDevicePath2 := "filtered2"
	filteredDeviceData2 := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": filteredLogFile2.Name(),
			"filter":    "mount_type == kv",
		},
	}
	_, err = client.Logical().Write("sys/audit/"+filteredDevicePath2, filteredDeviceData2)
	require.NoError(t, err)

	nonFilteredLogFile, err := os.CreateTemp(tempDir, "")
	nonFilteredDevicePath := "nonfiltered"
	nonFilteredDeviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": nonFilteredLogFile.Name(),
		},
	}
	_, err = client.Logical().Write("sys/audit/"+nonFilteredDevicePath, nonFilteredDeviceData)
	require.NoError(t, err)

	// Ensure the devices have been created.
	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	_, ok := devices[filteredDevicePath+"/"]
	require.True(t, ok)
	_, ok = devices[filteredDevicePath2+"/"]
	require.True(t, ok)
	_, ok = devices[nonFilteredDevicePath+"/"]
	require.True(t, ok)

	// Ensure the non-filtered log file is not empty.
	nonFilteredLogSize := getFileSize(t, nonFilteredLogFile.Name())
	require.Positive(t, nonFilteredLogSize)

	// A write to KV should produce an audit entry that is written to the
	// filtered devices and the non-filtered device.
	data := map[string]any{
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
		numberOfEntries := checkAuditEntries(t, f, "mount_type", "kv")
		require.Equal(t, 2, numberOfEntries)
	}

	// Disable the audit devices.
	err = client.Sys().DisableAudit(filteredDevicePath)
	require.NoError(t, err)
	err = client.Sys().DisableAudit(filteredDevicePath2)
	require.NoError(t, err)
	err = client.Sys().DisableAudit(nonFilteredDevicePath)
	require.NoError(t, err)
	// Ensure the devices are no longer there.
	devices, err = client.Sys().ListAudit()
	require.NoError(t, err)
	require.Len(t, devices, 0)
}

// TestAuditFilteringFallbackDevice validates that the audit device 'fallback'
// option works as expected. We create two audit devices, one with 'fallback'
// enabled and one with a filter, and make some auditable requests, then we
// ensure that correct entries were written to the respective log files.
func TestAuditFilteringFallbackDevice(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	tempDir := t.TempDir()
	fallbackLogFile, err := os.CreateTemp(tempDir, "")
	fallbackDevicePath := "fallback"
	fallbackDeviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": fallbackLogFile.Name(),
			"fallback":  "true",
		},
	}
	_, err = client.Logical().Write("sys/audit/"+fallbackDevicePath, fallbackDeviceData)
	require.NoError(t, err)

	filteredLogFile, err := os.CreateTemp(tempDir, "")
	filteredDevicePath := "filtered"
	filteredDeviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": filteredLogFile.Name(),
			"filter":    "mount_type == kv",
		},
	}
	_, err = client.Logical().Write("sys/audit/"+filteredDevicePath, filteredDeviceData)
	require.NoError(t, err)

	// Ensure the devices have been created.
	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	_, ok := devices[fallbackDevicePath+"/"]
	require.True(t, ok)
	_, ok = devices[filteredDevicePath+"/"]
	require.True(t, ok)

	// A write to KV should produce an audit entry that is written to the
	// filtered device.
	data := map[string]any{
		"foo": "bar",
	}
	err = client.KVv1("secret/").Put(context.Background(), "foo", data)
	require.NoError(t, err)

	// Disable the audit devices.
	err = client.Sys().DisableAudit(fallbackDevicePath)
	require.NoError(t, err)
	err = client.Sys().DisableAudit(filteredDevicePath)
	require.NoError(t, err)
	// Ensure the devices are no longer there.
	devices, err = client.Sys().ListAudit()
	require.Len(t, devices, 0)

	// Validate that only the entries matching the filter were written to the filtered log file.
	numberOfEntries := checkAuditEntries(t, filteredLogFile, "mount_type", "kv")
	require.Equal(t, 2, numberOfEntries)

	// Validate that only the entries NOT matching the filter were written to the fallback log file.
	numberOfEntries = 0
	scanner := bufio.NewScanner(fallbackLogFile)
	var auditRecord map[string]any
	for scanner.Scan() {
		auditRequest := map[string]any{}
		err := json.Unmarshal(scanner.Bytes(), &auditRecord)
		require.NoError(t, err)
		req, ok := auditRecord["request"]
		require.True(t, ok, "failed to parse request data from audit log entry")

		auditRequest = req.(map[string]any)
		require.NotEqual(t, "kv", auditRequest["mount_type"])
		numberOfEntries += 1
	}
	// the fallback device will catch all non-kv related entries such as login, etc. there should be 7 in total.
	require.Equal(t, 7, numberOfEntries)
}

// TestAuditFilteringFilterForUnsupportedField validates that the audit device
// 'filter' option fails when the filter expression selector references an
// unsupported field and that the error prevents an audit device from being
// created.
func TestAuditFilteringFilterForUnsupportedField(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	tempDir := t.TempDir()
	filteredLogFile, err := os.CreateTemp(tempDir, "")
	filteredDevicePath := "filtered"
	filteredDeviceData := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": filteredLogFile.Name(),
			"filter":    "auth == foo", // 'auth' is not one of the fields we allow filtering on
		},
	}
	_, err = client.Logical().Write("sys/audit/"+filteredDevicePath, filteredDeviceData)
	require.Error(t, err)
	require.ErrorContains(t, err, "audit.NewEntryFilter: filter references an unsupported field: auth == foo")

	// Ensure the device has not been created.
	devices, err := client.Sys().ListAudit()
	require.NoError(t, err)
	require.Len(t, devices, 0)

	// Now we do the same test but with the 'skip_test' option set to true.
	filteredDeviceDataSkipTest := map[string]any{
		"type":        "file",
		"description": "",
		"local":       false,
		"options": map[string]any{
			"file_path": filteredLogFile.Name(),
			"filter":    "auth == foo", // 'auth' is not one of the fields we allow filtering on
			"skip_test": true,
		},
	}
	_, err = client.Logical().Write("sys/audit/"+filteredDevicePath, filteredDeviceDataSkipTest)
	require.Error(t, err)
	require.ErrorContains(t, err, "audit.NewEntryFilter: filter references an unsupported field: auth == foo")

	// Ensure the device has not been created.
	devices, err = client.Sys().ListAudit()
	require.NoError(t, err)
	require.Len(t, devices, 0)
}

// getFileSize returns the size of the given file in bytes.
func getFileSize(t *testing.T, filePath string) int64 {
	t.Helper()
	fi, err := os.Stat(filePath)
	require.NoError(t, err)
	return fi.Size()
}

// checkAuditEntries parses the audit log file and asserts that the given key
// has the expected value for each entry. It returns the number of entries that
// were parsed.
func checkAuditEntries(t *testing.T, logFile *os.File, key string, expectedValue any) int {
	t.Helper()
	counter := 0
	scanner := bufio.NewScanner(logFile)
	var auditRecord map[string]any
	for scanner.Scan() {
		auditRequest := map[string]any{}
		err := json.Unmarshal(scanner.Bytes(), &auditRecord)
		require.NoError(t, err)
		req, ok := auditRecord["request"]
		require.True(t, ok, "failed to parse request data from audit log entry")
		auditRequest = req.(map[string]any)
		require.Equal(t, expectedValue, auditRequest[key])
		counter += 1
	}
	return counter
}
