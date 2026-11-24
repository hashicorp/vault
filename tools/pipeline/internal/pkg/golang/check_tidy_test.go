// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCheckTidyReq_Run verifies that a module with an unused requirement is
// reported as untidy with the diff tidy would apply, while a tidy module
// passes. The dependency is a local replace so the test runs offline.
func TestCheckTidyReq_Run(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work":  "go 1.21\n\nuse (\n\t./a\n\t./b\n)\n",
		"a/go.mod": "module example.com/a\n\ngo 1.21\n\nrequire example.com/b v0.0.0\n\nreplace example.com/b => ../b\n",
		"a/a.go":   "package a\n",
		"b/go.mod": "module example.com/b\n\ngo 1.21\n",
		"b/b.go":   "package b\n",
	})

	res, err := (&CheckTidyReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)
	require.False(t, res.Success())
	require.Equal(t, "1 of 2 modules not tidy, 0 failed", res.String())
	require.Len(t, res.Modules, 2)

	require.Equal(t, "a", res.Modules[0].Dir)
	require.False(t, res.Modules[0].Tidy)
	require.Contains(t, res.Modules[0].Diff, "-require example.com/b v0.0.0")
	require.Empty(t, res.Modules[0].Error)

	require.Equal(t, "b", res.Modules[1].Dir)
	require.True(t, res.Modules[1].Tidy)

	// Markdown keeps the diff out of the table so its lines stay intact.
	md := res.ToMarkdown()
	require.Contains(t, md, "| a | not tidy |  |")
	require.Contains(t, md, "`a`:\n\n```diff\n")
	require.Contains(t, md, "\n-require example.com/b v0.0.0\n")
	require.NotContains(t, md, "<br/>")
	require.Contains(t, res.ToTable(), "-require example.com/b v0.0.0")

	tidyOnly, err := (&CheckTidyReq{GoWork: filepath.Join(root, "go.work"), Modules: []string{"b"}}).Run(t.Context())
	require.NoError(t, err)
	require.True(t, tidyOnly.Success())
}

// TestCheckTidyReq_Run_Error verifies that a module tidy can't process is
// reported as an error rather than as tidy.
func TestCheckTidyReq_Run_Error(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work":  "go 1.21\n\nuse ./a\n",
		"a/go.mod": "module example.com/a\n\ngo 1.21\n\nrequire example.com/missing v0.0.0\n\nreplace example.com/missing => ../missing\n",
		"a/a.go":   "package a\n\nimport _ \"example.com/missing\"\n",
	})

	res, err := (&CheckTidyReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)
	require.False(t, res.Success())
	require.False(t, res.Modules[0].Tidy)
	require.Empty(t, res.Modules[0].Diff)
	require.Contains(t, res.Modules[0].Error, "missing")
	require.Equal(t, "0 of 1 modules not tidy, 1 failed", res.String())
}
