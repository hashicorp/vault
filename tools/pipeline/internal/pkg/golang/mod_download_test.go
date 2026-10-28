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

// writeModuleProxy writes a module proxy to proxyDir, in the directory layout
// that GOPROXY=file://<proxyDir> serves, with a single module example.com/dep
// at v1.0.0.
func writeModuleProxy(t *testing.T, proxyDir string) {
	t.Helper()

	const goMod = "module example.com/dep\n\ngo 1.21\n"

	var zipped bytes.Buffer
	zw := zip.NewWriter(&zipped)
	for name, content := range map[string]string{
		"example.com/dep@v1.0.0/go.mod": goMod,
		"example.com/dep@v1.0.0/dep.go": "package dep\n",
	} {
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = w.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())

	writeFiles(t, proxyDir, map[string]string{
		"example.com/dep/@v/list":        "v1.0.0\n",
		"example.com/dep/@v/v1.0.0.info": `{"Version":"v1.0.0","Time":"2024-01-01T00:00:00Z"}`,
		"example.com/dep/@v/v1.0.0.mod":  goMod,
		"example.com/dep/@v/v1.0.0.zip":  zipped.String(),
	})
}

// TestRunGoModDownloadAll verifies that 'go mod download all' runs against the
// given go.work and records the checksums of the workspace's requirements in
// go.work.sum, even when the caller's environment sets GOWORK=off, as CI does
// for module-mode builds. The test serves its dependency from a file:// module
// proxy so it needs no network.
//
// It sets environment variables, so it can't run in parallel.
func TestRunGoModDownloadAll(t *testing.T) {
	root := t.TempDir()
	proxyDir := filepath.Join(root, "proxy")
	writeModuleProxy(t, proxyDir)

	workspaceDir := filepath.Join(root, "workspace")
	writeFiles(t, workspaceDir, map[string]string{
		"go.work": "go 1.21\n\nuse ./a\n",
		"a/go.mod": "module example.com/a\n\ngo 1.21\n\n" +
			"require example.com/dep v1.0.0\n",
		"a/a.go": "package a\n",
	})

	t.Setenv("GOPROXY", "file://"+filepath.ToSlash(proxyDir))
	t.Setenv("GOMODCACHE", filepath.Join(root, "modcache"))
	t.Setenv("GOFLAGS", "-modcacherw")
	t.Setenv("GOSUMDB", "off")
	t.Setenv("GOWORK", "off")

	goWork := filepath.Join(workspaceDir, "go.work")
	require.NoError(t, RunGoModDownloadAll(t.Context(), goWork))

	sum, err := os.ReadFile(filepath.Join(workspaceDir, "go.work.sum"))
	require.NoError(t, err)
	require.Contains(t, string(sum), "example.com/dep v1.0.0 h1:")
	require.Contains(t, string(sum), "example.com/dep v1.0.0/go.mod h1:")
}

// TestRunGoModDownloadAll_Error verifies that a failure carries the go
// command's stderr, which is what tells the user why go.work.sum wasn't
// updated, and that a missing go.work path is rejected before running go.
func TestRunGoModDownloadAll_Error(t *testing.T) {
	t.Parallel()

	require.ErrorContains(t, RunGoModDownloadAll(t.Context(), ""), "go.work path is required")

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work": "go 1.21\n\nuse ./missing\n",
	})

	err := RunGoModDownloadAll(t.Context(), filepath.Join(root, "go.work"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing")
}

// TestRunGoModDownload verifies that 'go mod download' runs against a single
// module, with GOWORK=off, and records the go.mod checksum of a requirement
// that the module's go.sum doesn't hold yet, which is what a workspace sync
// leaves behind when it raises a requirement. The test serves its dependency
// from a file:// module proxy so it needs no network.
//
// It sets environment variables, so it can't run in parallel.
func TestRunGoModDownload(t *testing.T) {
	root := t.TempDir()
	proxyDir := filepath.Join(root, "proxy")
	writeModuleProxy(t, proxyDir)

	moduleDir := filepath.Join(root, "module")
	writeFiles(t, moduleDir, map[string]string{
		"go.mod": "module example.com/a\n\ngo 1.21\n\nrequire example.com/dep v1.0.0\n",
		"a.go":   "package a\n",
	})

	t.Setenv("GOPROXY", "file://"+filepath.ToSlash(proxyDir))
	t.Setenv("GOMODCACHE", filepath.Join(root, "modcache"))
	t.Setenv("GOFLAGS", "-modcacherw")
	t.Setenv("GOSUMDB", "off")

	// A workspace that contains the module must not be used.
	writeFiles(t, root, map[string]string{"go.work": "go 1.21\n\nuse ./module\n"})
	t.Setenv("GOWORK", filepath.Join(root, "go.work"))

	require.NoError(t, RunGoModDownload(t.Context(), filepath.Join(moduleDir, "go.mod")))

	sum, err := os.ReadFile(filepath.Join(moduleDir, "go.sum"))
	require.NoError(t, err)
	require.Contains(t, string(sum), "example.com/dep v1.0.0/go.mod h1:")

	_, err = os.Stat(filepath.Join(root, "go.work.sum"))
	require.ErrorIs(t, err, os.ErrNotExist, "the workspace must not be used")
}

// TestRunGoModDownload_Error verifies that a failure carries the go command's
// stderr, so the user can tell why the module's go.sum wasn't updated, and that
// a missing path is rejected before running go.
func TestRunGoModDownload_Error(t *testing.T) {
	t.Parallel()

	require.ErrorContains(t, RunGoModDownload(t.Context(), ""), "go.mod path is required")

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.mod": "module example.com/a\n\ngo 1.21\n\nrequire example.com/dep v1.0.0\n\nreplace example.com/dep => ./missing\n",
	})

	err := RunGoModDownload(t.Context(), filepath.Join(root, "go.mod"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing")
}
