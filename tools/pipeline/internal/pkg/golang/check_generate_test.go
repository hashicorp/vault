// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

// generateFixture is a two-module workspace whose generate directives write
// files with sh, so the test needs no generator binaries. Module a writes
// out.txt; module b writes nothing unless a test adds a directive.
//
// Keep directives in interpreted strings with \n, as here. go generate scans
// every line of every Go file, including tests, so a raw string line starting
// with //go:generate would run whenever this module is generated.
func generateFixture() map[string]string {
	return map[string]string{
		"go.work":    "go 1.21\n\nuse (\n\t./a\n\t./b\n)\n",
		"a/go.mod":   "module example.com/a\n\ngo 1.21\n",
		"a/a.go":     "package a\n\n//go:generate sh -c \"echo generated > out.txt\"\n",
		"a/out.txt":  "generated\n",
		"b/go.mod":   "module example.com/b\n\ngo 1.21\n",
		"b/b.go":     "package b\n",
		".gitignore": "*.gen\n",
	}
}

// TestCheckGenerateReq_Run_UpToDate verifies that regenerating identical
// output passes.
func TestCheckGenerateReq_Run_UpToDate(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, generateFixture())
	gitCommitAll(t, root)

	res, err := (&CheckGenerateReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)
	require.True(t, res.Success(), res.ToTable())
	require.Empty(t, res.Changes)
	require.Empty(t, res.Diff)
	require.Contains(t, res.ToTable(), "Generated code is up to date in 2 modules.")
}

// TestCheckGenerateReq_Run_Changes verifies that modified, added, deleted and
// new git-ignored files are each reported, that tracked changes are diffed,
// and that changes stay in the working tree.
func TestCheckGenerateReq_Run_Changes(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	files := generateFixture()
	files["a/out.txt"] = "stale\n"
	files["b/b.go"] = "package b\n\n" +
		"//go:generate sh -c \"echo new > new.txt\"\n" +
		"//go:generate sh -c \"echo ignored > ignored.gen\"\n" +
		"//go:generate rm old.txt\n"
	files["b/old.txt"] = "old\n"
	writeFiles(t, root, files)
	gitCommitAll(t, root)

	res, err := (&CheckGenerateReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)
	require.False(t, res.Success())
	require.Equal(t, []*GeneratedChange{
		{Path: "a/out.txt", Change: GeneratedChangeModified},
		{Path: "b/ignored.gen", Change: GeneratedChangeAddedIgnored},
		{Path: "b/new.txt", Change: GeneratedChangeAdded},
		{Path: "b/old.txt", Change: GeneratedChangeDeleted},
	}, res.Changes)
	require.Contains(t, res.Diff, "-stale")
	require.Contains(t, res.Diff, "+generated")
	require.Contains(t, res.Diff, "-old")
	require.NotContains(t, res.Diff, "new.txt")
	require.Equal(t, "changed paths: 4, modules that failed to generate: 0 of 2", res.String())
	require.Contains(t, res.ToMarkdown(), "| added-ignored | b/ignored.gen |")
	require.Contains(t, res.ToMarkdown(), "````diff\n")

	got, err := os.ReadFile(filepath.Join(root, "a", "out.txt"))
	require.NoError(t, err)
	require.Equal(t, "generated\n", string(got))
}

// TestCheckGenerateReq_Run_DiffIgnoresGitConfig verifies that the reported
// diff is a plain patch even when git config would change how git diff
// renders it, e.g. a developer's prefix settings or an external diff tool.
func TestCheckGenerateReq_Run_DiffIgnoresGitConfig(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	files := generateFixture()
	files["a/out.txt"] = "stale\n"
	writeFiles(t, root, files)
	gitCommitAll(t, root)
	for _, args := range [][]string{
		{"config", "diff.noprefix", "true"},
		{"config", "diff.external", "echo external diff"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v: %s", args, out)
	}

	res, err := (&CheckGenerateReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)
	require.Contains(t, res.Diff, "--- a/a/out.txt\n+++ b/a/out.txt\n")
	require.NotContains(t, res.Diff, "external diff")
}

// TestTruncateDiff verifies that long diffs are capped with a marker, and that
// the cut never splits a multi-byte UTF-8 character.
func TestTruncateDiff(t *testing.T) {
	t.Parallel()

	require.Equal(t, "short", truncateDiff("short", 1024))

	// "é" is two bytes, so after a leading "a" byte 1024 is the middle of one.
	diff := "a" + strings.Repeat("é", 1024)
	got := truncateDiff(diff, 1024)
	require.True(t, utf8.ValidString(got))
	require.Equal(t, "a"+strings.Repeat("é", 511)+"\n... diff truncated at 1 KiB ...\n", got)
}

// TestCheckGenerateReq_Run_DirtyTree verifies that a file already modified
// before the run is still reported when generation changes it again, because
// detection compares contents rather than only git status.
func TestCheckGenerateReq_Run_DirtyTree(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, generateFixture())
	gitCommitAll(t, root)
	writeFiles(t, root, map[string]string{"a/out.txt": "local edit\n"})

	res, err := (&CheckGenerateReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)
	require.Equal(t, []*GeneratedChange{
		{Path: "a/out.txt", Change: GeneratedChangeModified},
	}, res.Changes)
}

// TestCheckGenerateReq_Run_ModuleError verifies that a failing directive is
// reported for its module and doesn't stop later modules from generating.
func TestCheckGenerateReq_Run_ModuleError(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	files := generateFixture()
	files["a/a.go"] = "package a\n\n//go:generate sh -c \"echo boom >&2; exit 3\"\n"
	files["b/b.go"] = "package b\n\n//go:generate sh -c \"echo new > new.txt\"\n"
	writeFiles(t, root, files)
	gitCommitAll(t, root)

	res, err := (&CheckGenerateReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)
	require.False(t, res.Success())
	require.Len(t, res.Modules, 2)
	require.Contains(t, res.Modules[0].Error, "boom")
	require.Empty(t, res.Modules[1].Error)
	require.Equal(t, []*GeneratedChange{{Path: "b/new.txt", Change: GeneratedChangeAdded}}, res.Changes)
	require.Contains(t, res.ToTable(), "error")
}

// TestCheckGenerateReq_Run_Tags verifies that --tags reaches go generate, so
// directives in files behind build constraints run only when asked.
func TestCheckGenerateReq_Run_Tags(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	files := generateFixture()
	files["b/tagged.go"] = "//go:build fixture\n\npackage b\n\n//go:generate sh -c \"echo tagged > tagged.txt\"\n"
	writeFiles(t, root, files)
	gitCommitAll(t, root)
	goWork := filepath.Join(root, "go.work")

	untagged, err := (&CheckGenerateReq{GoWork: goWork, Modules: []string{"b"}}).Run(t.Context())
	require.NoError(t, err)
	require.True(t, untagged.Success(), untagged.ToTable())
	require.Len(t, untagged.Modules, 1)

	tagged, err := (&CheckGenerateReq{GoWork: goWork, Tags: "fixture", Modules: []string{"b"}}).Run(t.Context())
	require.NoError(t, err)
	require.Equal(t, []*GeneratedChange{{Path: "b/tagged.txt", Change: GeneratedChangeAdded}}, tagged.Changes)
}

// TestCheckGenerateReq_Run_ToolsTag verifies that the tools tag is rejected,
// since it would run the directive that installs every tool.
func TestCheckGenerateReq_Run_ToolsTag(t *testing.T) {
	t.Parallel()

	_, err := (&CheckGenerateReq{Tags: "ent,tools"}).Run(t.Context())
	require.ErrorContains(t, err, "tools build tag")
}

// TestGenerateEnv verifies that tags are added to GOFLAGS, keeping existing
// flags, and that cross-compile settings are cleared.
func TestGenerateEnv(t *testing.T) {
	t.Parallel()

	require.Equal(t,
		[]string{"GOOS=", "GOARCH=", "CC=", "CC_FOR_TARGET="},
		generateEnv("-mod=readonly", nil),
	)
	require.Equal(t,
		[]string{"GOOS=", "GOARCH=", "CC=", "CC_FOR_TARGET=", "GOFLAGS=-mod=readonly -tags=ent,enterprise"},
		generateEnv("-mod=readonly", []string{"ent", "enterprise"}),
	)
	require.Equal(t,
		[]string{"GOOS=", "GOARCH=", "CC=", "CC_FOR_TARGET=", "GOFLAGS=-tags=ent"},
		generateEnv("", []string{"ent"}),
	)
}

// TestParsePorcelainZ verifies parsing of NUL-terminated porcelain v1 status,
// including paths with spaces, and that malformed entries are rejected.
func TestParsePorcelainZ(t *testing.T) {
	t.Parallel()

	codes, err := parsePorcelainZ([]byte(" M a/out.txt\x00?? b/new file.txt\x00!! bin/\x00D  old.txt\x00"))
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"a/out.txt":      " M",
		"b/new file.txt": "??",
		"bin/":           "!!",
		"old.txt":        "D ",
	}, codes)

	_, err = parsePorcelainZ([]byte("bad\x00"))
	require.ErrorContains(t, err, "unexpected git status entry")
}
