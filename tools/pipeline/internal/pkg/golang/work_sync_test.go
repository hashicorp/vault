// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRunGoWorkSync verifies that 'go work sync' raises a workspace module's
// requirement to the version the workspace selects, even when the caller's
// environment sets GOWORK=off, as CI does for module-mode builds. Module a
// requires example.com/dep v1.0.0 and module b requires v1.1.0, so the
// workspace selects v1.1.0 and module a is behind. The test serves its
// dependency from a file:// module proxy so it needs no network.
//
// It sets environment variables, so it can't run in parallel.
func TestRunGoWorkSync(t *testing.T) {
	root := t.TempDir()
	proxyDir := filepath.Join(root, "proxy")
	for _, version := range []string{"v1.0.0", "v1.1.0"} {
		const goMod = "module example.com/dep\n\ngo 1.21\n"

		var zipped bytes.Buffer
		zw := zip.NewWriter(&zipped)
		for name, content := range map[string]string{
			"example.com/dep@" + version + "/go.mod": goMod,
			"example.com/dep@" + version + "/dep.go": "package dep\n",
		} {
			w, err := zw.Create(name)
			require.NoError(t, err)
			_, err = w.Write([]byte(content))
			require.NoError(t, err)
		}
		require.NoError(t, zw.Close())

		writeFiles(t, proxyDir, map[string]string{
			"example.com/dep/@v/" + version + ".info": `{"Version":"` + version + `","Time":"2024-01-01T00:00:00Z"}`,
			"example.com/dep/@v/" + version + ".mod":  goMod,
			"example.com/dep/@v/" + version + ".zip":  zipped.String(),
		})
	}
	writeFiles(t, proxyDir, map[string]string{"example.com/dep/@v/list": "v1.0.0\nv1.1.0\n"})

	workspaceDir := filepath.Join(root, "workspace")
	writeFiles(t, workspaceDir, map[string]string{
		"go.work":  "go 1.21\n\nuse (\n\t./a\n\t./b\n)\n",
		"a/go.mod": "module example.com/a\n\ngo 1.21\n\nrequire example.com/dep v1.0.0\n",
		"a/a.go":   "package a\n\nimport _ \"example.com/dep\"\n",
		"b/go.mod": "module example.com/b\n\ngo 1.21\n\nrequire example.com/dep v1.1.0\n",
		"b/b.go":   "package b\n\nimport _ \"example.com/dep\"\n",
	})

	t.Setenv("GOPROXY", "file://"+filepath.ToSlash(proxyDir))
	t.Setenv("GOMODCACHE", filepath.Join(root, "modcache"))
	t.Setenv("GOFLAGS", "-modcacherw")
	t.Setenv("GOSUMDB", "off")
	t.Setenv("GOWORK", "off")

	require.NoError(t, RunGoWorkSync(t.Context(), filepath.Join(workspaceDir, "go.work")))

	a, err := os.ReadFile(filepath.Join(workspaceDir, "a/go.mod"))
	require.NoError(t, err)
	require.Contains(t, string(a), "example.com/dep v1.1.0")
	require.NotContains(t, string(a), "example.com/dep v1.0.0")
}

// TestRunGoWorkSync_Error verifies that a failure carries the go command's
// stderr, so the user can tell why the workspace wasn't synced, and that a
// missing go.work path is rejected before running go.
func TestRunGoWorkSync_Error(t *testing.T) {
	t.Parallel()

	require.ErrorContains(t, RunGoWorkSync(t.Context(), ""), "go.work path is required")

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work": "go 1.21\n\nuse ./missing\n",
	})

	err := RunGoWorkSync(t.Context(), filepath.Join(root, "go.work"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing")
}
