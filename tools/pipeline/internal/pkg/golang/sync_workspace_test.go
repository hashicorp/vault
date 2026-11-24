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

// TestGoWorkSumSyncResult_Status verifies that the state of go.work.sum is
// described in one word, and that an error takes precedence over an update.
func TestGoWorkSumSyncResult_Status(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		res  *GoWorkSumSyncResult
		want string
	}{
		"unchanged":              {res: &GoWorkSumSyncResult{}, want: "unchanged"},
		"updated":                {res: &GoWorkSumSyncResult{Updated: true}, want: "updated"},
		"failed":                 {res: &GoWorkSumSyncResult{Error: "download failed"}, want: "failed"},
		"failed after an update": {res: &GoWorkSumSyncResult{Updated: true, Error: "download failed"}, want: "failed"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.want, test.res.Status())
		})
	}
}

// TestWorkspaceGoModPaths verifies that the modules of the go.work in a
// directory come back in go.work order, that a directory without a go.work has
// none, and that a go.work we can't use is an error rather than an empty
// workspace.
func TestWorkspaceGoModPaths(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		goWork  string
		missing bool // the directory has no go.work
		want    []string
		errHas  string
	}{
		"modules in go.work order": {
			goWork: "go 1.21\n\nuse (\n\t.\n\t./sdk\n\t./api/auth/approle\n)\n",
			want:   []string{"go.mod", "sdk/go.mod", "api/auth/approle/go.mod"},
		},
		"no go.work": {missing: true},
		"malformed go.work": {
			goWork: "use (\n",
			errHas: "listing workspace modules",
		},
		"go.work using a directory outside the repository": {
			goWork: "go 1.21\n\nuse ../other\n",
			errHas: "not inside",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			if !test.missing {
				require.NoError(t, os.WriteFile(filepath.Join(dir, "go.work"), []byte(test.goWork), 0o644))
			}

			paths, err := WorkspaceGoModPaths(dir)
			if test.errHas != "" {
				require.ErrorContains(t, err, test.errHas)
				require.Nil(t, paths)

				return
			}
			require.NoError(t, err)
			require.Equal(t, test.want, paths)
		})
	}
}

// TestSyncWorkspace_noGoWork verifies that a directory without a go.work has
// nothing to sync: no result, nothing to stage, and no error.
func TestSyncWorkspace_noGoWork(t *testing.T) {
	t.Parallel()

	res, stage, err := SyncWorkspace(t.Context(), t.TempDir(), []string{"go.mod"})
	require.NoError(t, err)
	require.Nil(t, res)
	require.Empty(t, stage)
}

// TestSyncWorkspace_failure verifies that when a go command fails we report why
// in the result, for the PR body, and return an error, and that we don't offer
// go.work.sum for staging.
func TestSyncWorkspace_failure(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(dir, "go.work"), []byte("go 1.21\n\nuse ./missing\n"), 0o644,
	))

	res, stage, err := SyncWorkspace(t.Context(), dir, []string{"missing/go.mod"})
	require.ErrorContains(t, err, "syncing the workspace")
	require.NotNil(t, res)
	require.False(t, res.Updated)
	require.Contains(t, res.Error, "missing")
	require.Empty(t, stage)
}

// TestSyncWorkspace verifies that go.work.sum is brought up to date from the
// requirements of the workspace's modules, that it's only offered for staging
// when it changed, and that a second sync with nothing new to record reports
// no change. The test serves its dependency from a file:// module
// proxy so it needs no network. It sets environment variables, so it can't run
// in parallel.
func TestSyncWorkspace(t *testing.T) {
	root := t.TempDir()

	const depMod = "module example.com/dep\n\ngo 1.21\n"
	var depZip bytes.Buffer
	zw := zip.NewWriter(&depZip)
	for name, content := range map[string]string{
		"example.com/dep@v1.0.0/go.mod": depMod,
		"example.com/dep@v1.0.0/dep.go": "package dep\n",
	} {
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = w.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())

	proxyDir := filepath.Join(root, "proxy")
	workspaceDir := filepath.Join(root, "workspace")
	for name, content := range map[string]string{
		filepath.Join(proxyDir, "example.com/dep/@v/list"):        "v1.0.0\n",
		filepath.Join(proxyDir, "example.com/dep/@v/v1.0.0.info"): `{"Version":"v1.0.0","Time":"2024-01-01T00:00:00Z"}`,
		filepath.Join(proxyDir, "example.com/dep/@v/v1.0.0.mod"):  depMod,
		filepath.Join(proxyDir, "example.com/dep/@v/v1.0.0.zip"):  depZip.String(),
		filepath.Join(workspaceDir, "go.work"):                    "go 1.21\n\nuse ./a\n",
		filepath.Join(workspaceDir, "a/go.mod"):                   "module example.com/a\n\ngo 1.21\n\nrequire example.com/dep v1.0.0\n",
		filepath.Join(workspaceDir, "a/a.go"):                     "package a\n",
	} {
		require.NoError(t, os.MkdirAll(filepath.Dir(name), 0o755))
		require.NoError(t, os.WriteFile(name, []byte(content), 0o644))
	}

	t.Setenv("GOPROXY", "file://"+filepath.ToSlash(proxyDir))
	t.Setenv("GOMODCACHE", filepath.Join(root, "modcache"))
	t.Setenv("GOFLAGS", "-modcacherw")
	t.Setenv("GOSUMDB", "off")
	t.Setenv("GOWORK", "off")

	res, stage, err := SyncWorkspace(t.Context(), workspaceDir, []string{"a/go.mod"})
	require.NoError(t, err)
	require.Equal(t, &GoWorkSumSyncResult{Updated: true}, res)
	require.Equal(t, []string{filepath.Join(workspaceDir, "go.work.sum")}, stage)

	sum, err := os.ReadFile(filepath.Join(workspaceDir, "go.work.sum"))
	require.NoError(t, err)
	require.Contains(t, string(sum), "example.com/dep v1.0.0 h1:")

	res, stage, err = SyncWorkspace(t.Context(), workspaceDir, []string{"a/go.mod"})
	require.NoError(t, err)
	require.Equal(t, &GoWorkSumSyncResult{}, res)
	require.Empty(t, stage)
}
