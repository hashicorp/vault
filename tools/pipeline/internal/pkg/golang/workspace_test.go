// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// writeFiles writes files, keyed by slash-separated path relative to root,
// creating parent directories as needed. Fixture modules are written at test
// time because a committed go.mod outside testdata would become a nested
// module of tools/pipeline.
func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()

	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}
}

// gitCommitAll initializes a git repository in dir if needed and commits every
// file in it. It ignores the user's global and system git config so commit
// signing and hooks can't interfere.
func gitCommitAll(t *testing.T, dir string) {
	t.Helper()

	env := append(os.Environ(),
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_SYSTEM=/dev/null",
		"GIT_AUTHOR_NAME=test",
		"GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test",
		"GIT_COMMITTER_EMAIL=test@example.com",
	)
	for _, args := range [][]string{
		{"init", "--quiet"},
		{"add", "--all"},
		{"commit", "--quiet", "--allow-empty", "--message", "fixture"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = env
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v: %s", args, out)
	}
}

// TestFindGoWork verifies that go.work is found by walking up from a nested
// directory, so commands work when run from inside any module.
func TestFindGoWork(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work":               "go 1.21\n\nuse ./a\n",
		"a/go.mod":              "module example.com/a\n\ngo 1.21\n",
		"a/internal/pkg/doc.go": "package pkg\n",
	})

	got, err := FindGoWork(filepath.Join(root, "a", "internal", "pkg"))
	require.NoError(t, err)
	require.Equal(t, filepath.Join(root, "go.work"), got)
}

// TestFindGoWork_NotFound verifies that a directory tree without go.work
// returns an error instead of an empty path.
func TestFindGoWork_NotFound(t *testing.T) {
	t.Parallel()

	_, err := FindGoWork(t.TempDir())
	require.ErrorContains(t, err, "no go.work file found")
}

// TestLoadWorkspace verifies that modules are loaded in go.work order with
// their module paths, normalized directories and go.mod paths.
func TestLoadWorkspace(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work":                 "go 1.21\n\nuse (\n\t.\n\t./api\n\t./api/auth/approle/\n)\n",
		"go.mod":                  "module example.com/root\n\ngo 1.21\n",
		"api/go.mod":              "module example.com/root/api\n\ngo 1.21\n",
		"api/auth/approle/go.mod": "module example.com/root/api/auth/approle\n\ngo 1.21\n",
	})

	ws, err := LoadWorkspace(t.Context(), filepath.Join(root, "go.work"))
	require.NoError(t, err)
	require.Equal(t, filepath.Join(root, "go.work"), ws.GoWork)
	require.Equal(t, root, ws.Root)
	require.Equal(t, []*WorkspaceModule{
		{Path: "example.com/root", Dir: ".", GoMod: "go.mod"},
		{Path: "example.com/root/api", Dir: "api", GoMod: "api/go.mod"},
		{Path: "example.com/root/api/auth/approle", Dir: "api/auth/approle", GoMod: "api/auth/approle/go.mod"},
	}, ws.Modules)
	require.Equal(t, filepath.Join(root, "api", "auth", "approle"), ws.AbsDir(ws.Modules[2]))
}

// TestLoadWorkspace_Errors verifies that broken workspaces fail to load rather
// than silently dropping modules from CI.
func TestLoadWorkspace_Errors(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		files map[string]string
		err   string
	}{
		"use without go.mod": {
			files: map[string]string{
				"go.work": "go 1.21\n\nuse ./a\n",
			},
			err: "go.work uses ./a, which has no go.mod",
		},
		"go.mod without module directive": {
			files: map[string]string{
				"go.work":  "go 1.21\n\nuse ./a\n",
				"a/go.mod": "go 1.21\n",
			},
			err: "has no module directive",
		},
		"duplicate module path": {
			files: map[string]string{
				"go.work":  "go 1.21\n\nuse (\n\t./a\n\t./b\n)\n",
				"a/go.mod": "module example.com/same\n\ngo 1.21\n",
				"b/go.mod": "module example.com/same\n\ngo 1.21\n",
			},
			err: "module example.com/same is used twice in go.work: a and b",
		},
		"unparsable go.work": {
			files: map[string]string{
				"go.work": "go 1.21\n\nnot-a-directive\n",
			},
			err: "parsing go.work",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			writeFiles(t, root, test.files)
			_, err := LoadWorkspace(t.Context(), filepath.Join(root, "go.work"))
			require.ErrorContains(t, err, test.err)
		})
	}
}

// TestWorkspace_ModuleForPackage verifies that packages map to the module with
// the longest matching path, so nested modules aren't tested from their
// parent's directory.
func TestWorkspace_ModuleForPackage(t *testing.T) {
	t.Parallel()

	ws := &Workspace{Modules: []*WorkspaceModule{
		{Path: "github.com/hashicorp/vault", Dir: "."},
		{Path: "github.com/hashicorp/vault/api", Dir: "api"},
		{Path: "github.com/hashicorp/vault/api/auth/approle", Dir: "api/auth/approle"},
		{Path: "github.com/hashicorp/vault/vault/hcp_link/proto", Dir: "vault/hcp_link/proto"},
	}}

	for importPath, wantDir := range map[string]string{
		"github.com/hashicorp/vault":                             ".",
		"github.com/hashicorp/vault/vault":                       ".",
		"github.com/hashicorp/vault/vault/hcp_link":              ".",
		"github.com/hashicorp/vault/vault/hcp_link/proto/node":   "vault/hcp_link/proto",
		"github.com/hashicorp/vault/api":                         "api",
		"github.com/hashicorp/vault/api/auth":                    "api",
		"github.com/hashicorp/vault/api/auth/approle":            "api/auth/approle",
		"github.com/hashicorp/vault/api/auth/approlex":           "api",
		"github.com/hashicorp/vault/apiextra":                    ".",
		"github.com/hashicorp/vault/api/auth/approle/internal/x": "api/auth/approle",
	} {
		t.Run(importPath, func(t *testing.T) {
			t.Parallel()

			m := ws.ModuleForPackage(importPath)
			require.NotNil(t, m)
			require.Equal(t, wantDir, m.Dir)
		})
	}

	require.Nil(t, ws.ModuleForPackage("github.com/hashicorp/vaultx"))
	require.Nil(t, ws.ModuleForPackage("example.com/other"))
}

// TestWorkspace_SelectModules verifies that --module selections are
// normalized, returned in go.work order, and that unknown directories fail.
func TestWorkspace_SelectModules(t *testing.T) {
	t.Parallel()

	ws := &Workspace{GoWork: "/ws/go.work", Root: "/ws", Modules: []*WorkspaceModule{
		{Path: "example.com/root", Dir: "."},
		{Path: "example.com/root/api", Dir: "api"},
		{Path: "example.com/root/sdk", Dir: "sdk"},
	}}

	all, err := ws.SelectModules(nil)
	require.NoError(t, err)
	require.Equal(t, ws.Modules, all)

	selected, err := ws.SelectModules([]string{"./sdk/", "./", "sdk"})
	require.NoError(t, err)
	require.Equal(t, []*WorkspaceModule{ws.Modules[0], ws.Modules[2]}, selected)

	_, err = ws.SelectModules([]string{"nope"})
	require.ErrorContains(t, err, "nope is not a module in /ws/go.work")
}

// TestNormalizeTags verifies that tag lists from workflow inputs, which may
// have leading commas or spaces, become a clean, ordered, de-duplicated list.
func TestNormalizeTags(t *testing.T) {
	t.Parallel()

	for input, want := range map[string][]string{
		"":                        {},
		",deadlock":               {"deadlock"},
		"ent,enterprise,deadlock": {"ent", "enterprise", "deadlock"},
		" ent  enterprise\t":      {"ent", "enterprise"},
		"ent,ent, enterprise,ent": {"ent", "enterprise"},
	} {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, want, NormalizeTags(input))
		})
	}
}
