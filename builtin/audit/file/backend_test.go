// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package file

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestAuditFile_fileModeNew verifies that the backend Factory correctly sets
// the file mode when the mode argument is set.
func TestAuditFile_fileModeNew(t *testing.T) {
	t.Parallel()

	modeStr := "0777"
	mode, err := strconv.ParseUint(modeStr, 8, 32)
	require.NoError(t, err)

	file := filepath.Join(t.TempDir(), "auditTest.txt")

	backendConfig := &audit.BackendConfig{
		Config: map[string]string{
			"path": file,
			"mode": modeStr,
		},
		MountPath:  "foo/bar",
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Logger:     hclog.NewNullLogger(),
	}
	_, err = Factory(context.Background(), backendConfig, nil)
	require.NoError(t, err)

	info, err := os.Stat(file)
	require.NoErrorf(t, err, "cannot retrieve file mode from `Stat`")
	require.Equalf(t, os.FileMode(mode), info.Mode(), "File mode does not match.")
}

// TestAuditFile_fileModeExisting verifies that the backend Factory correctly sets
// the mode on an existing file.
func TestAuditFile_fileModeExisting(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "auditTest.log")
	require.NoErrorf(t, err, "Failure to create test file.")

	err = os.Chmod(f.Name(), 0o777)
	require.NoErrorf(t, err, "Failure to chmod temp file for testing.")

	err = f.Close()
	require.NoErrorf(t, err, "Failure to close temp file for test.")

	backendConfig := &audit.BackendConfig{
		Config: map[string]string{
			"path": f.Name(),
		},
		MountPath:  "foo/bar",
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Logger:     hclog.NewNullLogger(),
	}

	_, err = Factory(context.Background(), backendConfig, nil)
	require.NoError(t, err)

	info, err := os.Stat(f.Name())
	require.NoErrorf(t, err, "cannot retrieve file mode from `Stat`")
	require.Equalf(t, os.FileMode(0o600), info.Mode(), "File mode does not match.")
}

// TestAuditFile_fileMode0000 verifies that setting the audit file mode to
// "0000" prevents Vault from modifying the permissions of the file.
func TestAuditFile_fileMode0000(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "auditTest.log")
	require.NoErrorf(t, err, "Failure to create test file.")

	err = os.Chmod(f.Name(), 0o777)
	require.NoErrorf(t, err, "Failure to chmod temp file for testing.")

	err = f.Close()
	require.NoErrorf(t, err, "Failure to close temp file for test.")

	backendConfig := &audit.BackendConfig{
		Config: map[string]string{
			"path": f.Name(),
			"mode": "0000",
		},
		MountPath:  "foo/bar",
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Logger:     hclog.NewNullLogger(),
	}

	_, err = Factory(context.Background(), backendConfig, nil)
	require.NoError(t, err)

	info, err := os.Stat(f.Name())
	require.NoErrorf(t, err, "cannot retrieve file mode from `Stat`. The error is %v", err)
	require.Equalf(t, os.FileMode(0o777), info.Mode(), "File mode does not match.")
}

// TestAuditFile_EventLogger_fileModeNew verifies that the Factory function
// correctly sets the file mode when the useEventLogger argument is set to
// true.
func TestAuditFile_EventLogger_fileModeNew(t *testing.T) {
	modeStr := "0777"
	mode, err := strconv.ParseUint(modeStr, 8, 32)
	require.NoError(t, err)

	file := filepath.Join(t.TempDir(), "auditTest.txt")

	backendConfig := &audit.BackendConfig{
		Config: map[string]string{
			"path": file,
			"mode": modeStr,
		},
		MountPath:  "foo/bar",
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Logger:     hclog.NewNullLogger(),
	}

	_, err = Factory(context.Background(), backendConfig, nil)
	require.NoError(t, err)

	info, err := os.Stat(file)
	require.NoErrorf(t, err, "Cannot retrieve file mode from `Stat`")
	require.Equalf(t, os.FileMode(mode), info.Mode(), "File mode does not match.")
}

// TestBackend_formatterConfig ensures that all the configuration values are parsed correctly.
func TestBackend_formatterConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config          map[string]string
		want            audit.FormatterConfig
		wantErr         bool
		expectedMessage string
	}{
		"happy-path-json": {
			config: map[string]string{
				"format":               audit.JSONFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "true",
			},
			want: audit.FormatterConfig{
				Raw:                true,
				HMACAccessor:       true,
				ElideListResponses: true,
				RequiredFormat:     "json",
			}, wantErr: false,
		},
		"happy-path-jsonx": {
			config: map[string]string{
				"format":               audit.JSONxFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "true",
			},
			want: audit.FormatterConfig{
				Raw:                true,
				HMACAccessor:       true,
				ElideListResponses: true,
				RequiredFormat:     "jsonx",
			},
			wantErr: false,
		},
		"invalid-format": {
			config: map[string]string{
				"format":               " squiggly ",
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "true",
			},
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "audit.NewFormatterConfig: error applying options: audit.(format).validate: 'squiggly' is not a valid format: invalid parameter",
		},
		"invalid-hmac-accessor": {
			config: map[string]string{
				"format":        audit.JSONFormat.String(),
				"hmac_accessor": "maybe",
			},
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "file.formatterConfig: unable to parse 'hmac_accessor': strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
		"invalid-log-raw": {
			config: map[string]string{
				"format":        audit.JSONFormat.String(),
				"hmac_accessor": "true",
				"log_raw":       "maybe",
			},
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "file.formatterConfig: unable to parse 'log_raw': strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
		"invalid-elide-bool": {
			config: map[string]string{
				"format":               audit.JSONFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "maybe",
			},
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "file.formatterConfig: unable to parse 'elide_list_responses': strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
	}
	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := formatterConfig(tc.config)
			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedMessage)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.want, got)
		})
	}
}

// TestBackend_configureFilterNode ensures that configureFilterNode handles various
// filter values as expected. Empty (including whitespace) strings should return
// no error but skip configuration of the node.
func TestBackend_configureFilterNode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		filter           string
		shouldSkipNode   bool
		wantErr          bool
		expectedErrorMsg string
	}{
		"happy": {
			filter: "operation == update",
		},
		"empty": {
			filter:         "",
			shouldSkipNode: true,
		},
		"spacey": {
			filter:         "    ",
			shouldSkipNode: true,
		},
		"bad": {
			filter:           "___qwerty",
			wantErr:          true,
			expectedErrorMsg: "file.(Backend).configureFilterNode: error creating filter node: audit.NewEntryFilter: cannot create new audit filter",
		},
		"unsupported-field": {
			filter:           "foo == bar",
			wantErr:          true,
			expectedErrorMsg: "filter references an unsupported field: foo == bar",
		},
	}
	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			b := &Backend{
				nodeIDList: []eventlogger.NodeID{},
				nodeMap:    map[eventlogger.NodeID]eventlogger.Node{},
			}

			err := b.configureFilterNode(tc.filter)

			switch {
			case tc.wantErr:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedErrorMsg)
				require.Len(t, b.nodeIDList, 0)
				require.Len(t, b.nodeMap, 0)
			case tc.shouldSkipNode:
				require.NoError(t, err)
				require.Len(t, b.nodeIDList, 0)
				require.Len(t, b.nodeMap, 0)
			default:
				require.NoError(t, err)
				require.Len(t, b.nodeIDList, 1)
				require.Len(t, b.nodeMap, 1)
				id := b.nodeIDList[0]
				node := b.nodeMap[id]
				require.Equal(t, eventlogger.NodeTypeFilter, node.Type())
			}
		})
	}
}

// TestBackend_configureFormatterNode ensures that configureFormatterNode
// populates the nodeIDList and nodeMap on Backend when given valid formatConfig.
func TestBackend_configureFormatterNode(t *testing.T) {
	t.Parallel()

	b := &Backend{
		nodeIDList: []eventlogger.NodeID{},
		nodeMap:    map[eventlogger.NodeID]eventlogger.Node{},
	}

	formatConfig, err := audit.NewFormatterConfig()
	require.NoError(t, err)

	err = b.configureFormatterNode("juan", formatConfig, hclog.NewNullLogger())

	require.NoError(t, err)
	require.Len(t, b.nodeIDList, 1)
	require.Len(t, b.nodeMap, 1)
	id := b.nodeIDList[0]
	node := b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeFormatter, node.Type())
}

// TestBackend_configureSinkNode ensures that we can correctly configure the sink
// node on the Backend, and any incorrect parameters result in the relevant errors.
func TestBackend_configureSinkNode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name           string
		filePath       string
		mode           string
		format         string
		wantErr        bool
		expectedErrMsg string
		expectedName   string
	}{
		"name-empty": {
			name:           "",
			wantErr:        true,
			expectedErrMsg: "file.(Backend).configureSinkNode: name is required: invalid parameter",
		},
		"name-whitespace": {
			name:           "   ",
			wantErr:        true,
			expectedErrMsg: "file.(Backend).configureSinkNode: name is required: invalid parameter",
		},
		"filePath-empty": {
			name:           "foo",
			filePath:       "",
			wantErr:        true,
			expectedErrMsg: "file.(Backend).configureSinkNode: file path is required: invalid parameter",
		},
		"filePath-whitespace": {
			name:           "foo",
			filePath:       "   ",
			wantErr:        true,
			expectedErrMsg: "file.(Backend).configureSinkNode: file path is required: invalid parameter",
		},
		"filePath-stdout-lower": {
			name:         "foo",
			expectedName: "stdout",
			filePath:     "stdout",
			format:       "json",
		},
		"filePath-stdout-upper": {
			name:         "foo",
			expectedName: "stdout",
			filePath:     "STDOUT",
			format:       "json",
		},
		"filePath-stdout-mixed": {
			name:         "foo",
			expectedName: "stdout",
			filePath:     "StdOut",
			format:       "json",
		},
		"filePath-discard-lower": {
			name:         "foo",
			expectedName: "discard",
			filePath:     "discard",
			format:       "json",
		},
		"filePath-discard-upper": {
			name:         "foo",
			expectedName: "discard",
			filePath:     "DISCARD",
			format:       "json",
		},
		"filePath-discard-mixed": {
			name:         "foo",
			expectedName: "discard",
			filePath:     "DisCArd",
			format:       "json",
		},
		"format-empty": {
			name:           "foo",
			filePath:       "/tmp/",
			format:         "",
			wantErr:        true,
			expectedErrMsg: "file.(Backend).configureSinkNode: format is required: invalid parameter",
		},
		"format-whitespace": {
			name:           "foo",
			filePath:       "/tmp/",
			format:         "   ",
			wantErr:        true,
			expectedErrMsg: "file.(Backend).configureSinkNode: format is required: invalid parameter",
		},
		"filePath-weird-with-mode-zero": {
			name:           "foo",
			filePath:       "/tmp/qwerty",
			format:         "json",
			mode:           "0",
			wantErr:        true,
			expectedErrMsg: "file.(Backend).configureSinkNode: file sink creation failed for path \"/tmp/qwerty\": event.NewFileSink: unable to determine existing file mode: stat /tmp/qwerty: no such file or directory",
		},
		"happy": {
			name:         "foo",
			filePath:     "/tmp/audit.log",
			mode:         "",
			format:       "json",
			wantErr:      false,
			expectedName: "foo",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			b := &Backend{
				nodeIDList: []eventlogger.NodeID{},
				nodeMap:    map[eventlogger.NodeID]eventlogger.Node{},
			}

			err := b.configureSinkNode(tc.name, tc.filePath, tc.mode, tc.format)

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrMsg)
				require.Len(t, b.nodeIDList, 0)
				require.Len(t, b.nodeMap, 0)
			} else {
				require.NoError(t, err)
				require.Len(t, b.nodeIDList, 1)
				require.Len(t, b.nodeMap, 1)
				id := b.nodeIDList[0]
				node := b.nodeMap[id]
				require.Equal(t, eventlogger.NodeTypeSink, node.Type())
				mc, ok := node.(*event.MetricsCounter)
				require.True(t, ok)
				require.Equal(t, tc.expectedName, mc.Name)
			}
		})
	}
}

// TestBackend_configureFilterFormatterSink ensures that configuring all three
// types of nodes on a Backend works as expected, i.e. we have all three nodes
// at the end and nothing gets overwritten. The order of calls influences the
// slice of IDs on the Backend.
func TestBackend_configureFilterFormatterSink(t *testing.T) {
	t.Parallel()

	b := &Backend{
		nodeIDList: []eventlogger.NodeID{},
		nodeMap:    map[eventlogger.NodeID]eventlogger.Node{},
	}

	formatConfig, err := audit.NewFormatterConfig()
	require.NoError(t, err)

	err = b.configureFilterNode("path == bar")
	require.NoError(t, err)

	err = b.configureFormatterNode("juan", formatConfig, hclog.NewNullLogger())
	require.NoError(t, err)

	err = b.configureSinkNode("foo", "/tmp/foo", "0777", "json")
	require.NoError(t, err)

	require.Len(t, b.nodeIDList, 3)
	require.Len(t, b.nodeMap, 3)

	id := b.nodeIDList[0]
	node := b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeFilter, node.Type())

	id = b.nodeIDList[1]
	node = b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeFormatter, node.Type())

	id = b.nodeIDList[2]
	node = b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeSink, node.Type())
}

// TestBackend_Factory_Conf is used to ensure that any configuration which is
// supplied, is validated and tested.
func TestBackend_Factory_Conf(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := map[string]struct {
		backendConfig        *audit.BackendConfig
		isErrorExpected      bool
		expectedErrorMessage string
	}{
		"nil-salt-config": {
			backendConfig: &audit.BackendConfig{
				SaltConfig: nil,
			},
			isErrorExpected:      true,
			expectedErrorMessage: "file.Factory: nil salt config",
		},
		"nil-salt-view": {
			backendConfig: &audit.BackendConfig{
				SaltConfig: &salt.Config{},
			},
			isErrorExpected:      true,
			expectedErrorMessage: "file.Factory: nil salt view",
		},
		"nil-logger": {
			backendConfig: &audit.BackendConfig{
				MountPath:  "discard",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     nil,
			},
			isErrorExpected:      true,
			expectedErrorMessage: "file.Factory: nil logger",
		},
		"fallback-device-with-filter": {
			backendConfig: &audit.BackendConfig{
				MountPath:  "discard",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"fallback":  "true",
					"file_path": discard,
					"filter":    "mount_type == kv",
				},
			},
			isErrorExpected:      true,
			expectedErrorMessage: "file.Factory: cannot configure a fallback device with a filter: invalid parameter",
		},
		"non-fallback-device-with-filter": {
			backendConfig: &audit.BackendConfig{
				MountPath:  "discard",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"fallback":  "false",
					"file_path": discard,
					"filter":    "mount_type == kv",
				},
			},
			isErrorExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			be, err := Factory(ctx, tc.backendConfig, nil)

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrorMessage)
			default:
				require.NoError(t, err)
				require.NotNil(t, be)
			}
		})
	}
}

// TestBackend_IsFallback ensures that the 'fallback' config setting is parsed
// and set correctly, then exposed via the interface method IsFallback().
func TestBackend_IsFallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := map[string]struct {
		backendConfig      *audit.BackendConfig
		isFallbackExpected bool
	}{
		"fallback": {
			backendConfig: &audit.BackendConfig{
				MountPath:  "discard",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"fallback":  "true",
					"file_path": discard,
				},
			},
			isFallbackExpected: true,
		},
		"no-fallback": {
			backendConfig: &audit.BackendConfig{
				MountPath:  "discard",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"fallback":  "false",
					"file_path": discard,
				},
			},
			isFallbackExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			be, err := Factory(ctx, tc.backendConfig, nil)
			require.NoError(t, err)
			require.NotNil(t, be)
			require.Equal(t, tc.isFallbackExpected, be.IsFallback())
		})
	}
}
