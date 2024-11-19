// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package trace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartDebugTrace(t *testing.T) {
	t.Run("error_on_empty_dir", func(t *testing.T) {
		_, _, err := StartDebugTrace("", "filePrefix")
		require.Error(t, err)
		require.Contains(t, err.Error(), "trace directory is required")
	})

	t.Run("error_on_non_existent_dir", func(t *testing.T) {
		_, _, err := StartDebugTrace("non-existent-dir", "filePrefix")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to stat trace directory")
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

	t.Run("successful_trace_generates_non_empty_file", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "")
		require.NoError(t, err)
		t.Cleanup(func() {
			os.RemoveAll(dir)
		})
		file, stop, err := StartDebugTrace(dir, "filePrefix")
		require.NoError(t, err)
		stop()
		f, err := os.Stat(file)
		require.NoError(t, err)
		require.True(t, f.Size() > 0)
	})
}
