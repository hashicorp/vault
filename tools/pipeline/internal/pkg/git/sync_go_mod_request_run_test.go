// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"archive/zip"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/stretchr/testify/require"
)

// setUpGoWorkspaceRepo creates a repository with two branches that have a Go
// workspace, and returns its directory with the dest branch ce/main checked out.
// The modules are:
//
//   - root requires dep and replaces sdk with the local sdk module. Its dep
//     requirement differs between the branches.
//   - sdk requires dep v1.1.0 on main but a dep version that doesn't exist on
//     ce/main. The root module is synced before sdk and can't be tidied until sdk
//     is synced, which shows that we sync every module before tidying any of them.
//   - tool requires dep v1.0.0 on main and v1.1.0 on ce/main. Its downgrade must
//     survive even though the workspace selects dep v1.1.0.
//   - leaf has no dependencies, so tidying it leaves it without a go.sum, which
//     we can't stage.
//   - extra is in the go.work of main but not that of ce/main, which shows that
//     the modules to sync come from the source rather than the dest. It isn't
//     tidied because it isn't in the workspace of ce/main.
//
// It serves dep from a file:// module proxy so we need no network, and sets
// environment variables, so a test that uses it can't run in parallel. The
// repository is made with the git client, with testGitConfig.
func setUpGoWorkspaceRepo(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	repoDir := filepath.Join(root, "repo")

	write := func(name, content string) {
		t.Helper()

		require.NoError(t, os.MkdirAll(filepath.Dir(name), 0o755))
		require.NoError(t, os.WriteFile(name, []byte(content), 0o644))
	}
	repo := libgit.NewClient(libgit.WithDir(repoDir), libgit.WithConfig(testGitConfig))

	// A module proxy with dep v1.0.0 and v1.1.0.
	proxyDir := filepath.Join(root, "proxy")
	for _, version := range []string{"v1.0.0", "v1.1.0"} {
		var zipped bytes.Buffer
		zw := zip.NewWriter(&zipped)
		for name, content := range map[string]string{
			"example.com/dep@" + version + "/go.mod": "module example.com/dep\n\ngo 1.21\n",
			"example.com/dep@" + version + "/dep.go": "package dep\n",
		} {
			w, err := zw.Create(name)
			require.NoError(t, err)
			_, err = w.Write([]byte(content))
			require.NoError(t, err)
		}
		require.NoError(t, zw.Close())

		versionDir := filepath.Join(proxyDir, "example.com/dep/@v")
		write(filepath.Join(versionDir, version+".info"), `{"Version":"`+version+`","Time":"2024-01-01T00:00:00Z"}`)
		write(filepath.Join(versionDir, version+".mod"), "module example.com/dep\n\ngo 1.21\n")
		write(filepath.Join(versionDir, version+".zip"), zipped.String())
	}
	write(filepath.Join(proxyDir, "example.com/dep/@v/list"), "v1.0.0\nv1.1.0\n")

	t.Setenv("GOPROXY", "file://"+filepath.ToSlash(proxyDir))
	t.Setenv("GOMODCACHE", filepath.Join(root, "modcache"))
	t.Setenv("GOFLAGS", "-modcacherw")
	t.Setenv("GOSUMDB", "off")
	t.Setenv("GOWORK", "off")

	// The dest branch comes first, and then main, the source, is made from it.
	require.NoError(t, os.MkdirAll(repoDir, 0o755))
	gitRes, err := repo.Init(t.Context(), &libgit.InitOpts{Quiet: true, InitialBranch: "ce/main"})
	require.NoError(t, err, gitRes.String())
	write(filepath.Join(repoDir, "go.work"), "go 1.21\n\nuse (\n\t.\n\t./leaf\n\t./sdk\n\t./tool\n)\n")
	write(filepath.Join(repoDir, "go.mod"), "module example.com/root\n\ngo 1.21\n\n"+
		"require (\n\texample.com/dep v1.0.0\n\texample.com/root/sdk v0.0.0\n)\n\n"+
		"replace example.com/root/sdk => ./sdk\n")
	write(filepath.Join(repoDir, "root.go"), "package root\n\nimport (\n\t_ \"example.com/dep\"\n\t_ \"example.com/root/sdk\"\n)\n")
	write(filepath.Join(repoDir, "sdk/go.mod"), "module example.com/root/sdk\n\ngo 1.21\n\nrequire example.com/dep v1.9.9\n")
	write(filepath.Join(repoDir, "sdk/sdk.go"), "package sdk\n\nimport _ \"example.com/dep\"\n")
	write(filepath.Join(repoDir, "tool/go.mod"), "module example.com/root/tool\n\ngo 1.21\n\nrequire example.com/dep v1.1.0\n")
	write(filepath.Join(repoDir, "tool/tool.go"), "package tool\n\nimport _ \"example.com/dep\"\n")
	write(filepath.Join(repoDir, "leaf/go.mod"), "module example.com/root/leaf\n\ngo 1.21\n")
	write(filepath.Join(repoDir, "leaf/leaf.go"), "package leaf\n")
	write(filepath.Join(repoDir, "extra/go.mod"), "module example.com/root/extra\n\ngo 1.21\n")
	write(filepath.Join(repoDir, "extra/extra.go"), "package extra\n")
	gitRes, err = repo.Add(t.Context(), &libgit.AddOpts{All: true})
	require.NoError(t, err, gitRes.String())
	gitRes, err = repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, Message: "ce/main"})
	require.NoError(t, err, gitRes.String())

	gitRes, err = repo.Checkout(t.Context(), &libgit.CheckoutOpts{Quiet: true, NewBranch: "main"})
	require.NoError(t, err, gitRes.String())
	write(filepath.Join(repoDir, "go.work"), "go 1.21\n\nuse (\n\t.\n\t./extra\n\t./leaf\n\t./sdk\n\t./tool\n)\n")
	write(filepath.Join(repoDir, "go.mod"), "module example.com/root\n\ngo 1.21\n\n"+
		"require (\n\texample.com/dep v1.1.0\n\texample.com/root/sdk v0.0.0\n)\n\n"+
		"replace example.com/root/sdk => ./sdk\n")
	write(filepath.Join(repoDir, "sdk/go.mod"), "module example.com/root/sdk\n\ngo 1.21\n\nrequire example.com/dep v1.1.0\n")
	write(filepath.Join(repoDir, "tool/go.mod"), "module example.com/root/tool\n\ngo 1.21\n\nrequire example.com/dep v1.0.0\n")
	gitRes, err = repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, All: true, Message: "main"})
	require.NoError(t, err, gitRes.String())

	gitRes, err = repo.Checkout(t.Context(), &libgit.CheckoutOpts{Quiet: true, Branch: "ce/main"})
	require.NoError(t, err, gitRes.String())

	// Resolve symlinks, as the root git reports has them resolved.
	resolved, err := filepath.EvalSymlinks(repoDir)
	require.NoError(t, err)

	return resolved
}

// TestSyncGoModReq_Run_workspaceModules verifies that, with workspace modules,
// we sync every module in the go.work of the source branch into a new branch of
// the dest, tidy every module in the go.work of the dest without raising their
// requirements to workspace-selected versions, and commit the go.mod and go.sum
// files that changed. We run from a subdirectory of the repository, as the
// modules of a go.work are relative to the root and
// nothing we read, write, or stage may depend on where we run. See
// setUpGoWorkspaceRepo for what the modules show.
//
// It changes the working directory and uses a fixture that sets environment
// variables, so it can't run in parallel.
func TestSyncGoModReq_Run_workspaceModules(t *testing.T) {
	repoDir := setUpGoWorkspaceRepo(t)
	t.Chdir(filepath.Join(repoDir, "sdk"))

	// The client we check the repository with runs from the root. The one we
	// give to Run doesn't, as it runs from where the command does.
	repo := libgit.NewClient(libgit.WithDir(repoDir), libgit.WithConfig(testGitConfig))
	mainRes, err := repo.RevParse(t.Context(), &libgit.RevParseOpts{Args: []string{"main"}})
	require.NoError(t, err, mainRes.String())
	destRes, err := repo.RevParse(t.Context(), &libgit.RevParseOpts{Args: []string{"ce/main"}})
	require.NoError(t, err, destRes.String())

	req := &SyncGoModReq{
		SourceBranch:     "main",
		DestBranch:       "ce/main",
		NewBranch:        "sync-test",
		WorkspaceModules: true,
		DiffOpts:         golang.DefaultSyncDiffOpts(),
		Commit:           true,
	}

	res, err := req.Run(t.Context(), libgit.NewClient(libgit.WithConfig(testGitConfig)))
	require.NoError(t, err)
	require.Equal(t, "sync-test", res.NewBranch)
	require.NotEmpty(t, res.CommitSHA)
	require.Nil(t, res.GoWorkSum)

	// Every module in the go.work of main was synced. The root and sdk modules
	// changed, as did tool. We tidied every module in the go.work of ce/main,
	// including leaf though it didn't change, but not extra.
	require.Len(t, res.Results, 5)
	for i, want := range []struct {
		path    string
		changes int
		tidy    bool
	}{
		{"go.mod", 1, true},
		{"extra/go.mod", 0, false},
		{"leaf/go.mod", 0, true},
		{"sdk/go.mod", 1, true},
		{"tool/go.mod", 1, true},
	} {
		require.Equal(t, want.path, res.Results[i].Path)
		require.Len(t, res.Results[i].Changes, want.changes, want.path)
		require.Equal(t, want.tidy, res.Results[i].TidyRan, want.path)
		require.Empty(t, res.Results[i].TidyError, want.path)
	}

	// Tool keeps its own requirement even though the workspace selects v1.1.0.
	diffRes, err := repo.Diff(t.Context(), &libgit.DiffOpts{
		NameOnly: true,
		Commits:  []string{"ce/main", "sync-test"},
	})
	require.NoError(t, err, diffRes.String())
	require.ElementsMatch(t, []string{
		"go.mod", "go.sum", "sdk/go.mod", "sdk/go.sum", "tool/go.mod", "tool/go.sum",
	}, strings.Fields(string(diffRes.Stdout)))
	for _, mod := range []string{"go.mod", "sdk/go.mod"} {
		showRes, err := repo.Show(t.Context(), &libgit.ShowOpts{Object: "sync-test:" + mod})
		require.NoError(t, err, showRes.String())
		require.Contains(t, string(showRes.Stdout), "example.com/dep v1.1.0", mod)
		require.NotContains(t, string(showRes.Stdout), "example.com/dep v1.0.0", mod)
	}
	showRes, err := repo.Show(t.Context(), &libgit.ShowOpts{Object: "sync-test:tool/go.mod"})
	require.NoError(t, err, showRes.String())
	require.Contains(t, string(showRes.Stdout), "example.com/dep v1.0.0")

	// We never touch the branches we sync, and we're left on the new one with
	// everything committed.
	afterRes, err := repo.RevParse(t.Context(), &libgit.RevParseOpts{Args: []string{"main"}})
	require.NoError(t, err, afterRes.String())
	require.Equal(t, string(mainRes.Stdout), string(afterRes.Stdout))
	afterRes, err = repo.RevParse(t.Context(), &libgit.RevParseOpts{Args: []string{"ce/main"}})
	require.NoError(t, err, afterRes.String())
	require.Equal(t, string(destRes.Stdout), string(afterRes.Stdout))
	headRes, err := repo.RevParse(t.Context(), &libgit.RevParseOpts{AbbrevRef: "loose", Args: []string{"HEAD"}})
	require.NoError(t, err, headRes.String())
	require.Equal(t, "sync-test", strings.TrimSpace(string(headRes.Stdout)))
	statusRes, err := repo.Status(t.Context(), &libgit.StatusOpts{Porcelain: true})
	require.NoError(t, err, statusRes.String())
	require.Empty(t, string(statusRes.Stdout))

	// Tidying each module again leaves its requirements unchanged.
	goCmd := func(dir string, env []string, args ...string) {
		t.Helper()

		cmd := exec.Command("go", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), env...)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "go %v in %s: %s", args, dir, out)
	}
	for _, mod := range []string{".", "sdk", "tool", "leaf"} {
		goCmd(filepath.Join(repoDir, mod), []string{"GOWORK=off"}, "mod", "tidy")
	}
	statusRes, err = repo.Status(t.Context(), &libgit.StatusOpts{Porcelain: true})
	require.NoError(t, err, statusRes.String())
	require.Empty(t, string(statusRes.Stdout))
}

// TestSyncGoModReq_Run_paths verifies that without workspace modules we sync,
// tidy, and commit only the go.mod paths we're given, relative to the working
// directory, and don't touch the workspace. The root module is listed before sdk
// and can't be tidied until sdk is synced, which shows that we sync every module
// before tidying any of them. See setUpGoWorkspaceRepo for the modules.
//
// It changes the working directory and uses a fixture that sets environment
// variables, so it can't run in parallel.
func TestSyncGoModReq_Run_paths(t *testing.T) {
	repoDir := setUpGoWorkspaceRepo(t)
	t.Chdir(repoDir)

	repo := libgit.NewClient(libgit.WithDir(repoDir), libgit.WithConfig(testGitConfig))

	req := &SyncGoModReq{
		SourceBranch: "main",
		DestBranch:   "ce/main",
		NewBranch:    "sync-test",
		Paths:        []string{"go.mod", "sdk/go.mod"},
		DiffOpts:     golang.DefaultSyncDiffOpts(),
		Commit:       true,
	}
	require.NoError(t, req.validate())

	res, err := req.Run(t.Context(), libgit.NewClient(libgit.WithConfig(testGitConfig)))
	require.NoError(t, err)
	require.NotEmpty(t, res.CommitSHA)
	require.Nil(t, res.GoWorkSum)

	// Only the modules we were given, which both changed and were tidied.
	require.Len(t, res.Results, 2)
	for i, path := range []string{"go.mod", "sdk/go.mod"} {
		require.Equal(t, path, res.Results[i].Path)
		require.Len(t, res.Results[i].Changes, 1, path)
		require.True(t, res.Results[i].TidyRan, path)
		require.Empty(t, res.Results[i].TidyError, path)
	}

	// Nothing else changed, including tool, which wasn't selected for syncing.
	diffRes, err := repo.Diff(t.Context(), &libgit.DiffOpts{
		NameOnly: true,
		Commits:  []string{"ce/main", "sync-test"},
	})
	require.NoError(t, err, diffRes.String())
	require.ElementsMatch(t, []string{"go.mod", "go.sum", "sdk/go.mod", "sdk/go.sum"}, strings.Fields(string(diffRes.Stdout)))
	showRes, err := repo.Show(t.Context(), &libgit.ShowOpts{Object: "sync-test:tool/go.mod"})
	require.NoError(t, err, showRes.String())
	require.Contains(t, string(showRes.Stdout), "example.com/dep v1.1.0")
	statusRes, err := repo.Status(t.Context(), &libgit.StatusOpts{Porcelain: true})
	require.NoError(t, err, statusRes.String())
	require.Empty(t, string(statusRes.Stdout))
}

// TestSyncGoModReq_Run_workspaceModulesNoGoWork verifies that when the source
// branch has no go.work, workspace mode fails rather than falling back to
// syncing only the root go.mod, which would look like a success while leaving
// every other module unsynced, and that it does so before it creates the new
// branch so that nothing is left behind.
func TestSyncGoModReq_Run_workspaceModulesNoGoWork(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	repo := libgit.NewClient(libgit.WithDir(dir), libgit.WithConfig(testGitConfig))

	gitRes, err := repo.Init(t.Context(), &libgit.InitOpts{Quiet: true, InitialBranch: "main"})
	require.NoError(t, err, gitRes.String())
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/root\n\ngo 1.21\n"), 0o644))
	gitRes, err = repo.Add(t.Context(), &libgit.AddOpts{All: true})
	require.NoError(t, err, gitRes.String())
	gitRes, err = repo.Commit(t.Context(), &libgit.CommitOpts{Quiet: true, Message: "main"})
	require.NoError(t, err, gitRes.String())
	gitRes, err = repo.Branch(t.Context(), &libgit.BranchOpts{BranchName: "dest"})
	require.NoError(t, err, gitRes.String())

	req := &SyncGoModReq{
		SourceBranch:     "main",
		DestBranch:       "dest",
		WorkspaceModules: true,
		DiffOpts:         golang.DefaultSyncDiffOpts(),
		Commit:           true,
	}

	res, err := req.Run(t.Context(), repo)
	require.ErrorContains(t, err, "git show main:go.work")
	require.Nil(t, res)

	branchRes, err := repo.Branch(t.Context(), &libgit.BranchOpts{Format: "%(refname:short)"})
	require.NoError(t, err, branchRes.String())
	require.ElementsMatch(t, []string{"main", "dest"}, strings.Fields(string(branchRes.Stdout)))
	headRes, err := repo.RevParse(t.Context(), &libgit.RevParseOpts{AbbrevRef: "loose", Args: []string{"HEAD"}})
	require.NoError(t, err, headRes.String())
	require.Equal(t, "main", strings.TrimSpace(string(headRes.Stdout)))
}

// TestSyncGoModReq_Run_workspaceModulesTrailingSpace verifies that a trailing
// space in the checkout name cannot redirect writes into a sibling directory.
// The fixture sets environment variables and the test changes directory, so it runs sequentially.
func TestSyncGoModReq_Run_workspaceModulesTrailingSpace(t *testing.T) {
	repoDir := setUpGoWorkspaceRepo(t)
	checkout := repoDir + " "
	require.NoError(t, os.Rename(repoDir, checkout))
	require.NoError(t, os.Mkdir(repoDir, 0o755))
	const siblingMod = "module example.com/untouched\n\ngo 1.21\n"
	require.NoError(t, os.WriteFile(filepath.Join(repoDir, "go.mod"), []byte(siblingMod), 0o644))
	t.Chdir(checkout)

	req := &SyncGoModReq{
		SourceBranch: "main", DestBranch: "ce/main", NewBranch: "sync-test",
		WorkspaceModules: true, DiffOpts: golang.DefaultSyncDiffOpts(), Commit: true,
	}
	res, syncErr := req.Run(t.Context(), libgit.NewClient(libgit.WithConfig(testGitConfig)))
	require.NoError(t, syncErr)
	require.NotEmpty(t, res.CommitSHA)
	sibling, readErr := os.ReadFile(filepath.Join(repoDir, "go.mod"))
	require.NoError(t, readErr)
	require.Equal(t, siblingMod, string(sibling))
	repo := libgit.NewClient(libgit.WithDir(checkout), libgit.WithConfig(testGitConfig))
	statusRes, statusErr := repo.Status(t.Context(), &libgit.StatusOpts{Porcelain: true})
	require.NoError(t, statusErr, statusRes.String())
	require.Empty(t, statusRes.Stdout)
}
