// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package trace

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestStartDebugTrace tests the debug trace functionality creating real
// files and traces.
func TestStartDebugTrace(t *testing.T) {
	t.Run("error_on_non_existent_dir", func(t *testing.T) {
		_, _, err := StartDebugTrace("non-existent-dir", "filePrefix")
		require.Error(t, err)
		require.Contains(t, err.Error(), "does not exist")
	})

	t.Run("error_on_non_dir", func(t *testing.T) {
		f, err := os.CreateTemp("", "")
		require.NoError(t, err)
		require.NoError(t, f.Close())
		_, _, err = StartDebugTrace(f.Name(), "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "is not a directory")
	})

	t.Run("error_on_failed_to_create_trace_file", func(t *testing.T) {
		noWriteFolder := filepath.Join(os.TempDir(), "no-write-permissions")
		// create folder without write permission
		err := os.Mkdir(noWriteFolder, 0o000)
		t.Cleanup(func() {
			os.RemoveAll(noWriteFolder)
		})
		require.NoError(t, err)
		_, _, err = StartDebugTrace(noWriteFolder, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create trace file")
	})

	t.Run("error_trying_to_start_second_concurrent_trace", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "")
		require.NoError(t, err)
		t.Cleanup(func() {
			os.RemoveAll(dir)
		})
		_, stop, err := StartDebugTrace(dir, "filePrefix")
		require.NoError(t, err)
		_, stopNil, err := StartDebugTrace(dir, "filePrefix")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to start trace")
		require.NoError(t, stop())
		require.Nil(t, stopNil)
	})

	t.Run("error_when_stating_tmp_dir_with_restricted_permissions", func(t *testing.T) {
		// this test relies on setting TMPDIR so skip it if we're not on a Unix system
		if runtime.GOOS == "windows" {
			t.Skip("skipping test on Windows")
		}

		tmpMissingPermissions := filepath.Join(t.TempDir(), "missing_permissions")
		err := os.Mkdir(tmpMissingPermissions, 0o000)
		require.NoError(t, err)
		t.Setenv("TMPDIR", tmpMissingPermissions)
		_, _, err = StartDebugTrace("", "filePrefix")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to stat trace directory")
	})

	t.Run("successful_trace_generates_non_empty_file", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "")
		require.NoError(t, err)
		t.Cleanup(func() {
			os.RemoveAll(dir)
		})
		file, stop, err := StartDebugTrace(dir, "filePrefix")
		require.NoError(t, err)
		require.NoError(t, stop())
		f, err := os.Stat(file)
		require.NoError(t, err)
		require.Greater(t, f.Size(), int64(0))
	})

	t.Run("successful_creation_of_tmp_dir", func(t *testing.T) {
		os.RemoveAll(filepath.Join(os.TempDir(), "vault-traces"))
		file, stop, err := StartDebugTrace("", "filePrefix")
		require.NoError(t, err)
		require.NoError(t, stop())
		require.Contains(t, file, filepath.Join(os.TempDir(), "vault-traces", "filePrefix"))
		f, err := os.Stat(file)
		require.NoError(t, err)
		require.Greater(t, f.Size(), int64(0))
	})

	t.Run("successful_trace_with_existing_tmp_dir", func(t *testing.T) {
		os.Mkdir(filepath.Join(os.TempDir(), "vault-traces"), 0o700)
		file, stop, err := StartDebugTrace("", "filePrefix")
		require.NoError(t, err)
		require.NoError(t, stop())
		require.Contains(t, file, filepath.Join(os.TempDir(), "vault-traces", "filePrefix"))
		f, err := os.Stat(file)
		require.NoError(t, err)
		require.Greater(t, f.Size(), int64(0))
	})
}
