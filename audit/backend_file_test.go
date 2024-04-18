// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
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

	backendConfig := &BackendConfig{
		Config: map[string]string{
			"path": file,
			"mode": modeStr,
		},
		MountPath:  "foo/bar",
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Logger:     hclog.NewNullLogger(),
	}
	_, err = newFileBackend(backendConfig, &NoopHeaderFormatter{})
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

	backendConfig := &BackendConfig{
		Config: map[string]string{
			"path": f.Name(),
		},
		MountPath:  "foo/bar",
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Logger:     hclog.NewNullLogger(),
	}

	_, err = newFileBackend(backendConfig, &NoopHeaderFormatter{})
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

	backendConfig := &BackendConfig{
		Config: map[string]string{
			"path": f.Name(),
			"mode": "0000",
		},
		MountPath:  "foo/bar",
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Logger:     hclog.NewNullLogger(),
	}

	_, err = newFileBackend(backendConfig, &NoopHeaderFormatter{})
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

	backendConfig := &BackendConfig{
		Config: map[string]string{
			"file_path": file,
			"mode":      modeStr,
		},
		MountPath:  "foo/bar",
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Logger:     hclog.NewNullLogger(),
	}

	_, err = newFileBackend(backendConfig, &NoopHeaderFormatter{})
	require.NoError(t, err)

	info, err := os.Stat(file)
	require.NoError(t, err)
	require.Equalf(t, os.FileMode(mode), info.Mode(), "File mode does not match.")
}

// TestFileBackend_newFileBackend ensures that we can correctly configure the sink
// node on the Backend, and any incorrect parameters result in the relevant errors.
func TestFileBackend_newFileBackend(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		mountPath      string
		filePath       string
		mode           string
		format         string
		wantErr        bool
		expectedErrMsg string
		expectedName   string
	}{
		"name-empty": {
			mountPath:      "",
			format:         "json",
			wantErr:        true,
			expectedErrMsg: "mount path cannot be empty: invalid configuration",
		},
		"name-whitespace": {
			mountPath:      "   ",
			format:         "json",
			wantErr:        true,
			expectedErrMsg: "mount path cannot be empty: invalid configuration",
		},
		"filePath-empty": {
			mountPath:      "foo",
			filePath:       "",
			format:         "json",
			wantErr:        true,
			expectedErrMsg: "file path is required: invalid configuration",
		},
		"filePath-whitespace": {
			mountPath:      "foo",
			filePath:       "   ",
			format:         "json",
			wantErr:        true,
			expectedErrMsg: "file path is required: invalid configuration",
		},
		"filePath-stdout-lower": {
			mountPath:    "foo",
			expectedName: "stdout",
			filePath:     "stdout",
			format:       "json",
		},
		"filePath-stdout-upper": {
			mountPath:    "foo",
			expectedName: "stdout",
			filePath:     "STDOUT",
			format:       "json",
		},
		"filePath-stdout-mixed": {
			mountPath:    "foo",
			expectedName: "stdout",
			filePath:     "StdOut",
			format:       "json",
		},
		"filePath-discard-lower": {
			mountPath:    "foo",
			expectedName: "discard",
			filePath:     "discard",
			format:       "json",
		},
		"filePath-discard-upper": {
			mountPath:    "foo",
			expectedName: "discard",
			filePath:     "DISCARD",
			format:       "json",
		},
		"filePath-discard-mixed": {
			mountPath:    "foo",
			expectedName: "discard",
			filePath:     "DisCArd",
			format:       "json",
		},
		"format-empty": {
			mountPath:      "foo",
			filePath:       "/tmp/",
			format:         "",
			wantErr:        true,
			expectedErrMsg: "unsupported \"format\": invalid configuration",
		},
		"format-whitespace": {
			mountPath:      "foo",
			filePath:       "/tmp/",
			format:         "   ",
			wantErr:        true,
			expectedErrMsg: "unsupported \"format\": invalid configuration",
		},
		"filePath-weird-with-mode-zero": {
			mountPath:      "foo",
			filePath:       "/tmp/qwerty",
			format:         "json",
			mode:           "0",
			wantErr:        true,
			expectedErrMsg: "file sink creation failed for path \"/tmp/qwerty\": unable to determine existing file mode: stat /tmp/qwerty: no such file or directory",
		},
		"happy": {
			mountPath:    "foo",
			filePath:     "/tmp/log",
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

			cfg := &BackendConfig{
				SaltView:   &logical.InmemStorage{},
				SaltConfig: &salt.Config{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"file_path": tc.filePath,
					"mode":      tc.mode,
					"format":    tc.format,
				},
				MountPath: tc.mountPath,
			}
			b, err := newFileBackend(cfg, &NoopHeaderFormatter{})

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrMsg)
				require.Nil(t, b)
			} else {
				require.NoError(t, err)
				require.Len(t, b.nodeIDList, 2) // Expect formatter + the sink
				require.Len(t, b.nodeMap, 2)
				id := b.nodeIDList[1]
				node := b.nodeMap[id]
				require.Equal(t, eventlogger.NodeTypeSink, node.Type())
				mc, ok := node.(*event.MetricsCounter)
				require.True(t, ok)
				require.Equal(t, tc.expectedName, mc.Name)
			}
		})
	}
}
