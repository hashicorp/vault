// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestAuditFile_NonExecutableModes verifies the actual file permissions,
// including the special mode that preserves an existing file's permissions.
func TestAuditFile_NonExecutableModes(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("Unix file permissions are not supported on Windows")
	}

	tests := []struct {
		name         string
		mode         string
		existingMode os.FileMode
		existing     bool
		padded       bool
		missingDir   bool
		checkReload  bool
		wantMode     os.FileMode
	}{
		{name: "configured executable mode", mode: "0755", wantMode: 0o644},
		{name: "execute-only mode", mode: "0111", wantMode: 0, checkReload: true},
		{name: "execute-only mode on existing file", mode: "0111", existing: true, existingMode: 0o755, wantMode: 0},
		{name: "write-only mode with missing directory", mode: "0300", missingDir: true, wantMode: 0o200},
		{name: "preserve non-executable mode", mode: "0000", existing: true, existingMode: 0o640, wantMode: 0o640},
		{name: "preserve mode with padded path", mode: "0000", existing: true, existingMode: 0o640, padded: true, wantMode: 0o640},
		{name: "clear existing executable mode", mode: "0000", existing: true, existingMode: 0o755, wantMode: 0o644},
		{name: "default mode", wantMode: 0o600},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			if tc.missingDir {
				dir = filepath.Join(dir, "logs")
			}
			file := filepath.Join(dir, "audit.log")
			if tc.existing {
				require.NoError(t, os.WriteFile(file, nil, 0o600))
				require.NoError(t, os.Chmod(file, tc.existingMode))
			}
			configuredPath := file
			if tc.padded {
				configuredPath = " " + file + " "
			}
			conf := &BackendConfig{
				Config:     map[string]string{"file_path": configuredPath, "mode": tc.mode},
				MountPath:  "audit/",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     hclog.NewNullLogger(),
			}

			backend, err := newFileBackend(conf, &noopHeaderFormatter{})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, backend.Close()) })

			info, err := os.Stat(file)
			require.NoError(t, err)
			require.Equal(t, tc.wantMode, info.Mode().Perm())
			if tc.checkReload {
				require.NoError(t, os.Chmod(file, 0o755))
				require.NoError(t, backend.Reload())
				info, err = os.Stat(file)
				require.NoError(t, err)
				require.Zero(t, info.Mode().Perm())
			}
		})
	}
}
