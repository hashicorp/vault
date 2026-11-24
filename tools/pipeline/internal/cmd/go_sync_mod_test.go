// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/stretchr/testify/require"
)

// testGoMods writes a source go.mod (in-memory only) and a dest go.mod (on disk)
// into a temp dir and returns their paths. The two modules differ so there is
// always a delta.
func testGoMods(t *testing.T) (srcPath, dstPath string) {
	t.Helper()
	dir := t.TempDir()

	srcPath = filepath.Join(dir, "src.mod")
	require.NoError(t, os.WriteFile(srcPath, []byte(`module github.com/example/src

go 1.21

require github.com/foo/bar v1.2.0
`), 0o644))

	dstPath = filepath.Join(dir, "dst.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(`module github.com/example/dst

go 1.21

require github.com/foo/bar v1.1.0
`), 0o644))

	return srcPath, dstPath
}

// TestGoSyncModCmd_TooFewArgs verifies that passing fewer than two positional
// arguments returns an error without executing the sync logic.
func TestGoSyncModCmd_TooFewArgs(t *testing.T) {
	goSyncModReq = &golang.SyncModReq{
		A:        &golang.ModSource{},
		B:        &golang.ModSource{},
		DiffOpts: golang.DefaultSyncDiffOpts(),
	}

	cmd := newGoSyncModCmd()
	suppressOutput(t, cmd)
	cmd.SetArgs([]string{"only-one-file.mod"})
	err := cmd.Execute()
	require.Error(t, err, "one argument should fail argument validation")
	require.Contains(t, err.Error(), "two local file paths are required")
}

// TestGoSyncModCmd_TooManyArgs verifies that passing more than two positional
// arguments returns an error describing the unexpected argument count.
func TestGoSyncModCmd_TooManyArgs(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(f, []byte("module m\ngo 1.21\n"), 0o644))

	goSyncModReq = &golang.SyncModReq{
		A:        &golang.ModSource{},
		B:        &golang.ModSource{},
		DiffOpts: golang.DefaultSyncDiffOpts(),
	}

	cmd := newGoSyncModCmd()
	suppressOutput(t, cmd)
	cmd.SetArgs([]string{f, f, f})
	err := cmd.Execute()
	require.Error(t, err, "three arguments should fail argument validation")
	require.Contains(t, err.Error(), "expected two path arguments")
}

// TestGoSyncModCmd_DefaultTableOutput verifies that the command writes a table
// to stdout by default when changes are present between the two go.mod files.
func TestGoSyncModCmd_DefaultTableOutput(t *testing.T) {
	srcPath, dstPath := testGoMods(t)

	goSyncModReq = &golang.SyncModReq{
		A:        &golang.ModSource{},
		B:        &golang.ModSource{},
		DiffOpts: golang.DefaultSyncDiffOpts(),
	}

	var buf bytes.Buffer
	cmd := newGoSyncModCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--dry-run", srcPath, dstPath})
	require.NoError(t, cmd.Execute())
	output := buf.String()

	require.NotEmpty(t, output, "expected table output for pending changes")
	require.Contains(t, output, "github.com/foo/bar")
	// Default output must not be valid JSON.
	var jsonCheck map[string]any
	require.Error(t, json.Unmarshal([]byte(output), &jsonCheck), "default output should be table, not JSON")
}

// TestGoSyncModCmd_FormatJSON verifies that --format json produces valid JSON
// output whose changes field contains the synced module entry.
func TestGoSyncModCmd_FormatJSON(t *testing.T) {
	srcPath, dstPath := testGoMods(t)

	goSyncModReq = &golang.SyncModReq{
		A:        &golang.ModSource{},
		B:        &golang.ModSource{},
		DiffOpts: golang.DefaultSyncDiffOpts(),
	}
	rootCfg.format = "json"
	defer func() { rootCfg.format = "" }()

	var buf bytes.Buffer
	cmd := newGoSyncModCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--dry-run", srcPath, dstPath})
	require.NoError(t, cmd.Execute())
	output := buf.String()

	require.NotEmpty(t, output)
	var parsed map[string]any
	require.NoError(t, json.Unmarshal([]byte(output), &parsed), "output should be valid JSON")
	require.Contains(t, parsed, "changes", "JSON output must include a changes key")
}

// TestGoSyncModCmd_FormatMarkdown verifies that --format markdown produces a
// pipe-delimited markdown table containing the changed module path.
func TestGoSyncModCmd_FormatMarkdown(t *testing.T) {
	srcPath, dstPath := testGoMods(t)

	goSyncModReq = &golang.SyncModReq{
		A:        &golang.ModSource{},
		B:        &golang.ModSource{},
		DiffOpts: golang.DefaultSyncDiffOpts(),
	}
	rootCfg.format = "markdown"
	defer func() { rootCfg.format = "" }()

	var buf bytes.Buffer
	cmd := newGoSyncModCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--dry-run", srcPath, dstPath})
	require.NoError(t, cmd.Execute())
	output := buf.String()

	require.NotEmpty(t, output, "markdown output should not be empty when changes are present")
	require.Contains(t, output, "|", "markdown table must contain pipe delimiters")
	require.Contains(t, output, "github.com/foo/bar", "markdown table must contain the changed module")
}
