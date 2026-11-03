// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"os"
	"path/filepath"
	"testing"

	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/stretchr/testify/require"
)

// testGitConfig is the git config our test repositories run with. It gives
// commits an identity, and stops the global and system config of whoever runs
// the tests, such as commit signing or a hooks directory, from changing what
// the tests do. It's set with the git client's WithConfig, which has higher
// precedence than any config file.
var testGitConfig = map[string]string{
	"user.name":      "test",
	"user.email":     "test@example.com",
	"commit.gpgsign": "false",
	"core.hooksPath": "/dev/null",
}

// TestSyncGoModFromRef verifies that we read the source go.mod from the ref and
// the dest go.mod from the directory we're given, and not from the working
// directory of the process, which is where git would resolve a path of its own
// relative to; that the result has the path we were given rather than the path
// we wrote to, as the path is how a result is matched to a module of a go.work;
// and that a go.mod that's missing from either side is reported in the result,
// as well as returned, so that the sync of the other modules can carry on.
func TestSyncGoModFromRef(t *testing.T) {
	t.Parallel()

	const (
		sourceMod = "module example.com/root\n\ngo 1.21\n\nrequire example.com/dep v1.1.0\n"
		destMod   = "module example.com/root\n\ngo 1.21\n\nrequire example.com/dep v1.0.0\n"
	)

	for name, test := range map[string]struct {
		p           string
		destFile    string // the go.mod in the working tree
		wantChanges int
		wantDest    string
		errHas      string
	}{
		"root go.mod": {
			p:           "go.mod",
			destFile:    "go.mod",
			wantChanges: 1,
			wantDest:    "example.com/dep v1.1.0",
		},
		"go.mod of a module in a directory": {
			p:           "sdk/go.mod",
			destFile:    "sdk/go.mod",
			wantChanges: 1,
			wantDest:    "example.com/dep v1.1.0",
		},
		"missing from the ref": {
			p:        "missing/go.mod",
			destFile: "go.mod",
			errHas:   "git show main:missing/go.mod",
		},
		"missing from the working tree": {
			p:        "sdk/go.mod",
			destFile: "go.mod", // there's no sdk/go.mod on the branch we're on
			errHas:   "reading",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			repo := libgit.NewClient(libgit.WithDir(dir), libgit.WithConfig(testGitConfig))
			write := func(name, content string) {
				t.Helper()

				require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644))
			}

			// main has the go.mod we sync from, in the root and in sdk. We then
			// check out a branch that has older requirements, in the root only
			// unless the test wants an sdk too.
			gitRes, err := repo.Init(t.Context(), &libgit.InitOpts{Quiet: true, InitialBranch: "main"})
			require.NoError(t, err, gitRes.String())
			write("go.mod", sourceMod)
			write("sdk/go.mod", sourceMod)
			gitRes, err = repo.Add(t.Context(), &libgit.AddOpts{All: true})
			require.NoError(t, err, gitRes.String())
			gitRes, err = repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, Message: "main"})
			require.NoError(t, err, gitRes.String())
			gitRes, err = repo.Checkout(t.Context(), &libgit.CheckoutOpts{Quiet: true, NewBranch: "dest"})
			require.NoError(t, err, gitRes.String())
			write("go.mod", destMod)
			if test.destFile == "sdk/go.mod" {
				write("sdk/go.mod", destMod)
			} else {
				require.NoError(t, os.RemoveAll(filepath.Join(dir, "sdk")))
			}
			gitRes, err = repo.Add(t.Context(), &libgit.AddOpts{All: true})
			require.NoError(t, err, gitRes.String())
			gitRes, err = repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, Message: "dest"})
			require.NoError(t, err, gitRes.String())

			res, err := SyncGoModFromRef(
				t.Context(), repo, "main", dir, test.p, golang.DefaultSyncDiffOpts(),
			)

			// Whatever the outcome, the result is the module we were given.
			require.NotNil(t, res)
			require.Equal(t, test.p, res.Path)

			if test.errHas != "" {
				require.ErrorContains(t, err, test.errHas)
				require.Contains(t, res.Error, test.errHas)

				return
			}
			require.NoError(t, err)
			require.Empty(t, res.Error)
			require.Len(t, res.Changes, test.wantChanges)

			got, err := os.ReadFile(filepath.Join(dir, test.destFile))
			require.NoError(t, err)
			require.Contains(t, string(got), test.wantDest)
			require.NotContains(t, string(got), "example.com/dep v1.0.0")
		})
	}
}

// TestWorkspaceGoModPathsFromRef verifies that the modules to sync come from
// the go.work committed at the root of the given ref, in go.work order, and
// not from the branch that's checked out or the working tree. Falling back to
// the root go.mod when go.work can't be read would leave every other module
// unsynced while looking like a success, so a missing, empty, or malformed
// go.work is an error.
func TestWorkspaceGoModPathsFromRef(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		goWork string // the go.work committed on the main branch, if any
		want   []string
		errHas string
	}{
		"modules in go.work order": {
			goWork: "go 1.21\n\nuse (\n\t.\n\t./sdk\n\t./api/auth/approle\n)\n",
			want:   []string{"go.mod", "sdk/go.mod", "api/auth/approle/go.mod"},
		},
		"missing go.work": {
			errHas: "git show main:go.work",
		},
		"go.work without modules": {
			goWork: "go 1.21\n",
			errHas: "has no modules",
		},
		"go.work using a directory outside the repository": {
			goWork: "go 1.21\n\nuse ../other\n",
			errHas: "not inside",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			repo := libgit.NewClient(libgit.WithDir(dir), libgit.WithConfig(testGitConfig))

			// Commit the go.work to main, then check out a branch with a
			// different go.work and an uncommitted change to it. Neither may
			// affect what we read from main.
			gitRes, err := repo.Init(t.Context(), &libgit.InitOpts{Quiet: true, InitialBranch: "main"})
			require.NoError(t, err, gitRes.String())
			require.NoError(t, os.WriteFile(filepath.Join(dir, "README"), []byte("fixture\n"), 0o644))
			if test.goWork != "" {
				require.NoError(t, os.WriteFile(filepath.Join(dir, "go.work"), []byte(test.goWork), 0o644))
			}
			gitRes, err = repo.Add(t.Context(), &libgit.AddOpts{All: true})
			require.NoError(t, err, gitRes.String())
			gitRes, err = repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, Message: "main"})
			require.NoError(t, err, gitRes.String())
			gitRes, err = repo.Checkout(t.Context(), &libgit.CheckoutOpts{Quiet: true, NewBranch: "other"})
			require.NoError(t, err, gitRes.String())
			require.NoError(t, os.WriteFile(filepath.Join(dir, "go.work"), []byte("go 1.21\n\nuse ./other\n"), 0o644))
			gitRes, err = repo.Add(t.Context(), &libgit.AddOpts{All: true})
			require.NoError(t, err, gitRes.String())
			gitRes, err = repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, Message: "other"})
			require.NoError(t, err, gitRes.String())
			require.NoError(t, os.WriteFile(filepath.Join(dir, "go.work"), []byte("go 1.21\n\nuse ./dirty\n"), 0o644))

			paths, err := WorkspaceGoModPathsFromRef(t.Context(), repo, "main")
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

// TestSyncGoModFromRef_symlink verifies that syncing directives cannot overwrite
// an external module through a destination directory or go.mod symlink.
func TestSyncGoModFromRef_symlink(t *testing.T) {
	t.Parallel()

	for name, directory := range map[string]bool{"module directory": true, "go.mod": false} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir, outside := t.TempDir(), t.TempDir()
			repo := libgit.NewClient(libgit.WithDir(dir), libgit.WithConfig(testGitConfig))
			gitRes, initErr := repo.Init(t.Context(), &libgit.InitOpts{Quiet: true, InitialBranch: "main"})
			require.NoError(t, initErr, gitRes.String())
			require.NoError(t, os.Mkdir(filepath.Join(dir, "sdk"), 0o755))
			require.NoError(t, os.WriteFile(filepath.Join(dir, "sdk/go.mod"),
				[]byte("module example.com/sdk\n\ngo 1.21\n\nrequire example.com/dep v1.1.0\n"), 0o644))
			gitRes, addErr := repo.Add(t.Context(), &libgit.AddOpts{All: true})
			require.NoError(t, addErr, gitRes.String())
			gitRes, commitErr := repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, Message: "source"})
			require.NoError(t, commitErr, gitRes.String())

			const externalMod = "module example.com/sdk\n\ngo 1.21\n\nrequire example.com/dep v1.0.0\n"
			require.NoError(t, os.WriteFile(filepath.Join(outside, "go.mod"), []byte(externalMod), 0o644))
			require.NoError(t, os.Remove(filepath.Join(dir, "sdk/go.mod")))
			if directory {
				require.NoError(t, os.Remove(filepath.Join(dir, "sdk")))
				require.NoError(t, os.Symlink(outside, filepath.Join(dir, "sdk")))
			} else {
				require.NoError(t, os.Symlink(filepath.Join(outside, "go.mod"), filepath.Join(dir, "sdk/go.mod")))
			}

			res, syncErr := SyncGoModFromRef(t.Context(), repo, "main", dir, "sdk/go.mod", golang.DefaultSyncDiffOpts())
			require.ErrorContains(t, syncErr, "symlink")
			require.Contains(t, res.Error, "symlink")
			got, readErr := os.ReadFile(filepath.Join(outside, "go.mod"))
			require.NoError(t, readErr)
			require.Equal(t, externalMod, string(got))
		})
	}
}
