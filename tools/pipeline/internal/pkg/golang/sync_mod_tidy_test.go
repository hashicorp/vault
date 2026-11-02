// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTidySyncedModules verifies the part of a sync that follows syncing the
// directives: which modules we tidy, which of their files we offer for staging,
// without syncing the workspace. The modules have no dependencies so that
// tidying them needs no network, which also means they have no go.sum.
func TestTidySyncedModules(t *testing.T) {
	t.Parallel()

	files := func() map[string]string {
		return map[string]string{
			"go.work":    "go 1.21\n\nuse (\n\t./a\n\t./b\n)\n",
			"a/go.mod":   "module example.com/a\n\ngo 1.21\n",
			"a/a.go":     "package a\n",
			"b/go.mod":   "module example.com/b\n\ngo 1.21\n",
			"b/b.go":     "package b\n",
			"c/go.mod":   "module example.com/c\n\ngo 1.21\n",
			"c/c.go":     "package c\n",
			"bad/go.mod": "this is not a go.mod\n",
		}
	}

	changed := func(p string) *GoModPathSyncResult {
		return &GoModPathSyncResult{SyncModRes: &SyncModRes{Path: p, Changes: []*SyncChange{{}}}}
	}
	unchanged := func(p string) *GoModPathSyncResult {
		return &GoModPathSyncResult{SyncModRes: &SyncModRes{Path: p}}
	}
	tidied := func(results []*GoModPathSyncResult) map[string]bool {
		ran := map[string]bool{}
		for _, pr := range results {
			ran[pr.Path] = pr.TidyRan
		}

		return ran
	}

	t.Run("nothing changed leaves the working tree alone", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeFiles(t, dir, files())

		// Not even in workspace mode.
		res, err := TidySyncedModules(t.Context(), dir, []*GoModPathSyncResult{unchanged("a/go.mod")}, true)
		require.NoError(t, err)
		require.Equal(t, map[string]bool{"a/go.mod": false}, tidied(res.Results))
		require.Nil(t, res.GoWorkSum)
		require.Empty(t, res.Stage)
	})

	t.Run("a changed module is tidied and its go.mod staged", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeFiles(t, dir, files())

		res, err := TidySyncedModules(
			t.Context(), dir, []*GoModPathSyncResult{changed("a/go.mod"), unchanged("c/go.mod")}, false,
		)
		require.NoError(t, err)

		// Only the module that changed, and not the others, nor the workspace. The
		// module has no go.sum, which we'd fail to stage, so we don't offer it.
		require.Equal(t, map[string]bool{"a/go.mod": true, "c/go.mod": false}, tidied(res.Results))
		require.Nil(t, res.GoWorkSum)
		require.Equal(t, []string{filepath.Join(dir, "a/go.mod")}, res.Stage)
	})

	t.Run("workspace mode tidies every module without syncing the workspace", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeFiles(t, dir, files())

		res, err := TidySyncedModules(t.Context(), dir, []*GoModPathSyncResult{changed("a/go.mod")}, true)
		require.NoError(t, err)

		// b didn't change and we hadn't synced it, but it's in the workspace so we
		// tidied it and added its result. c isn't in the workspace.
		require.Len(t, res.Results, 2)
		require.Equal(t, map[string]bool{"a/go.mod": true, "b/go.mod": true}, tidied(res.Results))
		require.Nil(t, res.GoWorkSum)
		require.Equal(t, []string{filepath.Join(dir, "a/go.mod"), filepath.Join(dir, "b/go.mod")}, res.Stage)
	})

	t.Run("a module that fails to tidy is reported and not staged", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeFiles(t, dir, files())

		res, err := TidySyncedModules(
			t.Context(), dir, []*GoModPathSyncResult{changed("bad/go.mod"), changed("a/go.mod")}, false,
		)
		require.ErrorContains(t, err, "go mod tidy in")

		// We still tidy the module that can be tidied.
		bad := res.Results[0]
		require.True(t, bad.TidyRan)
		require.Error(t, bad.TidyErr)
		require.NotEmpty(t, bad.TidyError)
		require.True(t, res.Results[1].TidyRan)
		require.Empty(t, res.Results[1].TidyError)
		require.Equal(t, []string{filepath.Join(dir, "a/go.mod")}, res.Stage)
	})

	t.Run("a go.work that we can't read is reported but the sync carries on", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeFiles(t, dir, files())
		require.NoError(t, os.WriteFile(filepath.Join(dir, "go.work"), []byte("use (\n"), 0o644))

		res, err := TidySyncedModules(t.Context(), dir, []*GoModPathSyncResult{changed("a/go.mod")}, true)
		require.ErrorContains(t, err, "listing workspace modules")
		require.NotNil(t, res.GoWorkSum)
		require.Contains(t, res.GoWorkSum.Error, "listing workspace modules")

		// We don't know the workspace, so we tidy only the module that changed.
		require.Equal(t, map[string]bool{"a/go.mod": true}, tidied(res.Results))
		require.Equal(t, []string{filepath.Join(dir, "a/go.mod")}, res.Stage)
	})
}

// TestTidySyncedModules_snapshotFailure verifies that a failed initial snapshot
// neither reports modules as tidied nor offers their untidied files for staging.
func TestTidySyncedModules_snapshotFailure(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	const aMod = "module example.com/a\n\ngo 1.21\n\nrequire example.com/unused v1.0.0\n\nreplace example.com/unused => ./unused\n"
	writeFiles(t, dir, map[string]string{
		"go.work":        "go 1.21\n\nuse (\n ./a\n ./b\n)\n",
		"a/go.mod":       aMod,
		"a/a.go":         "package a\n",
		"b/go.mod":       "module example.com/b\n\ngo 1.21\n",
		"b/go.sum/.keep": "",
	})
	res, err := TidySyncedModules(t.Context(), dir, []*GoModPathSyncResult{
		{SyncModRes: &SyncModRes{Path: "a/go.mod", Changes: []*SyncChange{{}}}},
	}, true)
	require.ErrorContains(t, err, "reading")
	require.Empty(t, res.Stage)
	require.Len(t, res.Results, 2)
	for _, result := range res.Results {
		require.False(t, result.TidyRan, result.Path)
		require.Contains(t, result.TidyError, "reading", result.Path)
	}
	got, readErr := os.ReadFile(filepath.Join(dir, "a/go.mod"))
	require.NoError(t, readErr)
	require.Equal(t, aMod, string(got))
}

// TestTidySyncedModules_invalidWorkspaceModule verifies that destination-only
// entries cannot redirect tidy outside the checkout or to an enclosing module.
func TestTidySyncedModules_invalidWorkspaceModule(t *testing.T) {
	t.Parallel()

	for name, symlink := range map[string]bool{"missing go.mod": false, "external symlink": true} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir, outside := t.TempDir(), t.TempDir()
			const untouchedMod = "module example.com/untouched\n\ngo 1.21\n\nrequire example.com/unused v1.0.0\n\nreplace example.com/unused => ./unused\n"
			writeFiles(t, dir, map[string]string{
				"go.work":  "go 1.21\n\nuse (\n ./a\n ./other\n)\n",
				"go.mod":   untouchedMod,
				"root.go":  "package root\n",
				"a/go.mod": "module example.com/a\n\ngo 1.21\n",
				"a/a.go":   "package a\n",
			})
			writeFiles(t, outside, map[string]string{"go.mod": untouchedMod, "other.go": "package other\n"})
			if symlink {
				require.NoError(t, os.Symlink(outside, filepath.Join(dir, "other")))
			} else {
				require.NoError(t, os.Mkdir(filepath.Join(dir, "other"), 0o755))
			}

			res, err := TidySyncedModules(t.Context(), dir, []*GoModPathSyncResult{
				{SyncModRes: &SyncModRes{Path: "a/go.mod", Changes: []*SyncChange{{}}}},
			}, true)
			require.Error(t, err)
			require.Equal(t, []string{filepath.Join(dir, "a/go.mod")}, res.Stage)
			require.Len(t, res.Results, 2)
			require.False(t, res.Results[1].TidyRan)
			require.NotEmpty(t, res.Results[1].TidyError)
			for _, root := range []string{dir, outside} {
				got, readErr := os.ReadFile(filepath.Join(root, "go.mod"))
				require.NoError(t, readErr)
				require.Equal(t, untouchedMod, string(got), root)
			}
		})
	}
}

// TestModulesToTidy verifies which modules we tidy after syncing. Nothing is
// tidied when the sync changed nothing, so the working tree is left alone.
// Otherwise a changed module is tidied, and in workspace mode so is every module
// in the workspace, in go.work order, because a module that didn't change can
// need tidying when a module it replaces did. A module that failed to sync isn't
// tidied, and a changed module that isn't in the workspace is tidied last.
func TestModulesToTidy(t *testing.T) {
	t.Parallel()

	changed := func(p string) *GoModPathSyncResult {
		return &GoModPathSyncResult{SyncModRes: &SyncModRes{
			Path: p, Changes: []*SyncChange{{}},
		}}
	}
	unchanged := func(p string) *GoModPathSyncResult {
		return &GoModPathSyncResult{SyncModRes: &SyncModRes{Path: p}}
	}
	failed := func(p string) *GoModPathSyncResult {
		return &GoModPathSyncResult{SyncModRes: &SyncModRes{Path: p, Error: "boom"}}
	}

	for name, test := range map[string]struct {
		results   []*GoModPathSyncResult
		workspace []string
		want      []string
	}{
		"nothing changed": {
			results:   []*GoModPathSyncResult{unchanged("go.mod"), unchanged("sdk/go.mod")},
			workspace: []string{"go.mod", "sdk/go.mod"},
		},
		"no results": {},
		"changed modules without a workspace": {
			results: []*GoModPathSyncResult{changed("go.mod"), unchanged("api/go.mod"), changed("sdk/go.mod")},
			want:    []string{"go.mod", "sdk/go.mod"},
		},
		"every workspace module in go.work order": {
			results: []*GoModPathSyncResult{
				changed("go.mod"), unchanged("api/go.mod"), changed("sdk/go.mod"), unchanged("version/go.mod"),
			},
			workspace: []string{"sdk/go.mod", "go.mod", "version/go.mod", "api/go.mod"},
			want:      []string{"sdk/go.mod", "go.mod", "version/go.mod", "api/go.mod"},
		},
		"workspace module that wasn't synced": {
			results:   []*GoModPathSyncResult{changed("go.mod")},
			workspace: []string{"go.mod", "only-in-b/go.mod"},
			want:      []string{"go.mod", "only-in-b/go.mod"},
		},
		"changed module outside the workspace is last": {
			results:   []*GoModPathSyncResult{changed("outside/go.mod"), changed("go.mod")},
			workspace: []string{"go.mod", "sdk/go.mod"},
			want:      []string{"go.mod", "sdk/go.mod", "outside/go.mod"},
		},
		"failed module isn't tidied": {
			results:   []*GoModPathSyncResult{changed("go.mod"), failed("sdk/go.mod")},
			workspace: []string{"go.mod", "sdk/go.mod"},
			want:      []string{"go.mod"},
		},
		"nil results are skipped": {
			results: []*GoModPathSyncResult{nil, {}, changed("go.mod")},
			want:    []string{"go.mod"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.want, modulesToTidy(test.results, test.workspace))
		})
	}
}

// TestResultFor verifies that we return the result of a module we synced, and
// add a result for a module we didn't, which is how a module that's only in the
// go.work of the branch we synced into shows up in the results.
func TestResultFor(t *testing.T) {
	t.Parallel()

	synced := &GoModPathSyncResult{SyncModRes: &SyncModRes{Path: "go.mod"}}
	results := []*GoModPathSyncResult{nil, {}, synced}

	got, results := resultFor(results, "go.mod")
	require.Same(t, synced, got)
	require.Len(t, results, 3)

	added, results := resultFor(results, "only-in-b/go.mod")
	require.Equal(t, "only-in-b/go.mod", added.Path)
	require.Len(t, results, 4)

	got, results = resultFor(results, "only-in-b/go.mod")
	require.Same(t, added, got)
	require.Len(t, results, 4)
}

// TestExistingFiles verifies that we only offer files that exist for staging,
// because git add fails on a path that doesn't, as with the go.sum of a module
// that has no dependencies.
func TestExistingFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	goMod := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(goMod, []byte("module example.com/a\n"), 0o644))

	require.Equal(t, []string{goMod}, existingFiles([]string{goMod, filepath.Join(dir, "go.sum")}))
	require.Empty(t, existingFiles(nil))
}
