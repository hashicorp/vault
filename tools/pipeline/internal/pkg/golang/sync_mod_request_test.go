// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSyncModReq_Require_UpdatesAndAdds verifies that when strict mode is on,
// requires present in source but absent in dest are added, and requires present
// in both with different versions are updated to the source version.
func TestSyncModReq_Require_UpdatesAndAdds(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require (
	github.com/foo/bar v1.2.0
	github.com/foo/baz v1.2.0
	github.com/foo/new v1.3.0
)
`
	const dstMod = `module github.com/example/dst

go 1.21

require (
	github.com/foo/bar v1.1.0
	github.com/foo/baz v1.2.0
)
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Changes, 2, "expected update to bar and addition of new")

	// Find the update.
	var updated, added *SyncChange
	for _, c := range res.Changes {
		switch c.Module {
		case "github.com/foo/bar":
			updated = c
		case "github.com/foo/new":
			added = c
		}
	}
	require.NotNil(t, updated)
	require.Equal(t, SyncActionUpdated, updated.Action)
	require.Equal(t, "v1.1.0", updated.OldVersion)
	require.Equal(t, "v1.2.0", updated.NewVersion)

	require.NotNil(t, added)
	require.Equal(t, SyncActionAdded, added.Action)
	require.Equal(t, "v1.3.0", added.NewVersion)

	// Verify file was written.
	written, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	require.Contains(t, string(written), "v1.2.0")
	require.Contains(t, string(written), "github.com/foo/new")
}

// TestSyncModReq_Require_StrictFalse verifies that when strict mode is off,
// requires present only in source are not added to dest, and only shared
// modules with different versions are updated.
func TestSyncModReq_Require_StrictFalse(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require (
	github.com/foo/bar v1.2.0
	github.com/foo/new v1.3.0
)
`
	const dstMod = `module github.com/example/dst

go 1.21

require (
	github.com/foo/bar v1.1.0
)
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	opts.StrictDiffRequire = false
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.Len(t, res.Changes, 1, "should only update bar, not add new")
	require.Equal(t, "github.com/foo/bar", res.Changes[0].Module)
	require.Equal(t, SyncActionUpdated, res.Changes[0].Action)

	// foo/new should not appear in the dest file.
	written, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	require.NotContains(t, string(written), "github.com/foo/new")
}

// TestSyncModReq_Replace_UpdatesAndAdds verifies that replace directives with
// different new-path/version are updated, and new replaces are added when strict.
func TestSyncModReq_Replace_UpdatesAndAdds(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require github.com/foo/bar v1.0.0

replace (
	github.com/foo/bar => github.com/foo/bar v1.2.0
	github.com/foo/new => github.com/foo/new v2.0.0
)
`
	const dstMod = `module github.com/example/dst

go 1.21

require github.com/foo/bar v1.0.0

replace github.com/foo/bar => github.com/foo/bar v1.1.0
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.Len(t, res.Changes, 2)

	var updated, added *SyncChange
	for _, c := range res.Changes {
		if c.Module == "github.com/foo/bar" {
			updated = c
		}
		if c.Module == "github.com/foo/new" {
			added = c
		}
	}
	require.NotNil(t, updated)
	require.Equal(t, SyncActionUpdated, updated.Action)
	require.Equal(t, "v1.1.0", updated.OldVersion)
	require.Equal(t, "v1.2.0", updated.NewVersion)

	require.NotNil(t, added)
	require.Equal(t, SyncActionAdded, added.Action)
}

// TestSyncModReq_Replace_StrictFalse verifies that replace directives present
// only in source are not added to dest when strict mode is off.
func TestSyncModReq_Replace_StrictFalse(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require github.com/foo/bar v1.0.0

replace (
	github.com/foo/bar => github.com/foo/bar v1.2.0
	github.com/foo/new => github.com/foo/new v2.0.0
)
`
	const dstMod = `module github.com/example/dst

go 1.21

require github.com/foo/bar v1.0.0

replace github.com/foo/bar => github.com/foo/bar v1.1.0
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	opts.StrictDiffReplace = false
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.Len(t, res.Changes, 1)
	require.Equal(t, "github.com/foo/bar", res.Changes[0].Module)
	require.Equal(t, SyncActionUpdated, res.Changes[0].Action)

	written, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	require.NotContains(t, string(written), "github.com/foo/new")
}

// TestSyncModReq_ExcludeRequire_GlobPattern verifies that requires matching a
// glob pattern are skipped during sync.
//
// Note: matchesExcludePattern uses path.Match, where '*' does not cross '/'.
// To exclude a module and its subpaths, provide one pattern per path component
// (e.g. "github.com/hashicorp/vault-enterprise" and
// "github.com/hashicorp/vault-enterprise/api"), or use a single pattern scoped
// to the exact path level (e.g. "github.com/hashicorp/vault-enterprise/a*" to
// match only direct children whose name starts with 'a').
func TestSyncModReq_ExcludeRequire_GlobPattern(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require (
	github.com/hashicorp/vault-enterprise v1.2.0
	github.com/hashicorp/vault-enterprise/api v1.2.0
	github.com/foo/bar v1.2.0
)
`
	const dstMod = `module github.com/example/dst

go 1.21

require (
	github.com/hashicorp/vault-enterprise v1.0.0
	github.com/hashicorp/vault-enterprise/api v1.0.0
	github.com/foo/bar v1.0.0
)
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	// Two patterns: one for the root module, one for the /api subpath.
	// path.Match '*' does not cross '/', so each path component needs its own
	// pattern or literal entry.
	opts.ExcludeRequire = []string{
		"github.com/hashicorp/vault-enterprise",
		"github.com/hashicorp/vault-enterprise/api",
	}
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)

	// Only foo/bar should be changed; vault-enterprise matches should be excluded.
	require.Len(t, res.Changes, 1)
	require.Equal(t, "github.com/foo/bar", res.Changes[0].Module)

	written, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	// vault-enterprise should stay at v1.0.0.
	require.NotContains(t, string(written), "vault-enterprise v1.2.0")
}

// TestSyncModReq_ExcludeRequire_LiteralPattern verifies that a literal
// (non-glob) exclude pattern correctly skips an exact module path.
func TestSyncModReq_ExcludeRequire_LiteralPattern(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require (
	github.com/foo/bar v1.2.0
	github.com/foo/baz v1.2.0
)
`
	const dstMod = `module github.com/example/dst

go 1.21

require (
	github.com/foo/bar v1.0.0
	github.com/foo/baz v1.0.0
)
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	opts.ExcludeRequire = []string{"github.com/foo/bar"}
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.Len(t, res.Changes, 1)
	require.Equal(t, "github.com/foo/baz", res.Changes[0].Module)
}

// TestSyncModReq_ExcludeReplace_MatchesOldPath verifies that a replace
// directive is skipped when its old path matches an exclude pattern.
func TestSyncModReq_ExcludeReplace_MatchesOldPath(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require github.com/foo/bar v1.0.0

replace (
	github.com/foo/bar => github.com/foo/bar v1.2.0
	github.com/foo/baz => github.com/foo/baz v2.0.0
)
`
	const dstMod = `module github.com/example/dst

go 1.21

require github.com/foo/bar v1.0.0

replace (
	github.com/foo/bar => github.com/foo/bar v1.1.0
	github.com/foo/baz => github.com/foo/baz v1.0.0
)
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	opts.Require = false
	opts.ExcludeReplace = []string{"github.com/foo/bar"}
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)

	// Only baz should be updated; bar's replace is excluded.
	require.Len(t, res.Changes, 1)
	require.Equal(t, "github.com/foo/baz", res.Changes[0].Module)

	written, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	require.Contains(t, string(written), "github.com/foo/bar v1.1.0")
}

// TestSyncModReq_ExcludeReplace_MatchesNewPath verifies that a replace
// directive is skipped when its new path matches an exclude pattern.
func TestSyncModReq_ExcludeReplace_MatchesNewPath(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require github.com/foo/bar v1.0.0

replace github.com/foo/bar => github.com/internal/bar v1.2.0
`
	const dstMod = `module github.com/example/dst

go 1.21

require github.com/foo/bar v1.0.0

replace github.com/foo/bar => github.com/internal/bar v1.1.0
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	opts.Require = false
	// Exclude by new path pattern.
	opts.ExcludeReplace = []string{"github.com/internal/*"}
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.Empty(t, res.Changes, "replace with excluded new path should be skipped")

	// File should not have been updated.
	written, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	require.Contains(t, string(written), "github.com/internal/bar v1.1.0")
}

// TestSyncModReq_RequireFlagFalse verifies that setting Require=false entirely
// skips require sync even when diffs exist.
func TestSyncModReq_RequireFlagFalse(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require github.com/foo/bar v1.2.0
`
	const dstMod = `module github.com/example/dst

go 1.21

require github.com/foo/bar v1.0.0
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	origContent := []byte(dstMod)
	require.NoError(t, os.WriteFile(dstPath, origContent, 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	opts.Require = false
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.Empty(t, res.Changes)

	// File should be unchanged.
	written, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	require.Equal(t, string(origContent), string(written))
}

// TestSyncModReq_ReplaceFlagFalse verifies that setting Replace=false entirely
// skips replace sync even when diffs exist.
func TestSyncModReq_ReplaceFlagFalse(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require github.com/foo/bar v1.0.0

replace github.com/foo/bar => github.com/foo/bar v1.2.0
`
	const dstMod = `module github.com/example/dst

go 1.21

require github.com/foo/bar v1.0.0

replace github.com/foo/bar => github.com/foo/bar v1.1.0
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	origContent := []byte(dstMod)
	require.NoError(t, os.WriteFile(dstPath, origContent, 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	opts.Require = false
	opts.Replace = false
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.Empty(t, res.Changes)
}

// TestSyncModReq_GoDirective verifies that the go directive is only synced
// when the Go option is explicitly enabled (it defaults false).
func TestSyncModReq_GoDirective(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.22
`
	const dstMod = `module github.com/example/dst

go 1.21
`

	t.Run("go flag false (default)", func(t *testing.T) {
		t.Parallel()

		src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}
		dir := t.TempDir()
		dstPath := filepath.Join(dir, "go.mod")
		require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
		dst := &ModSource{}
		require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

		opts := DefaultSyncDiffOpts()
		opts.Require = false
		opts.Replace = false
		req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

		res, err := req.Run(context.Background())
		require.NoError(t, err)
		require.Empty(t, res.Changes, "go directive should not be synced by default")

		written, err := os.ReadFile(dstPath)
		require.NoError(t, err)
		require.Contains(t, string(written), "go 1.21")
	})

	t.Run("go flag true", func(t *testing.T) {
		t.Parallel()

		src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}
		dir := t.TempDir()
		dstPath := filepath.Join(dir, "go.mod")
		require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
		dst := &ModSource{}
		require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

		opts := DefaultSyncDiffOpts()
		opts.Require = false
		opts.Replace = false
		opts.Go = true
		req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

		res, err := req.Run(context.Background())
		require.NoError(t, err)
		require.Len(t, res.Changes, 1)
		require.Equal(t, DirectiveGo, res.Changes[0].Directive)
		require.Equal(t, "1.21", res.Changes[0].OldVersion)
		require.Equal(t, "1.22", res.Changes[0].NewVersion)
	})
}

// TestSyncModReq_ToolchainDirective verifies that the toolchain directive is
// only synced when the Toolchain option is explicitly enabled (it defaults false).
func TestSyncModReq_ToolchainDirective(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

toolchain go1.22.0
`
	const dstMod = `module github.com/example/dst

go 1.21

toolchain go1.21.0
`

	t.Run("toolchain flag false (default)", func(t *testing.T) {
		t.Parallel()

		src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}
		dir := t.TempDir()
		dstPath := filepath.Join(dir, "go.mod")
		require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
		dst := &ModSource{}
		require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

		opts := DefaultSyncDiffOpts()
		opts.Require = false
		opts.Replace = false
		req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

		res, err := req.Run(context.Background())
		require.NoError(t, err)
		require.Empty(t, res.Changes, "toolchain directive should not be synced by default")

		written, err := os.ReadFile(dstPath)
		require.NoError(t, err)
		require.Contains(t, string(written), "toolchain go1.21.0")
	})

	t.Run("toolchain flag true", func(t *testing.T) {
		t.Parallel()

		src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}
		dir := t.TempDir()
		dstPath := filepath.Join(dir, "go.mod")
		require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
		dst := &ModSource{}
		require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

		opts := DefaultSyncDiffOpts()
		opts.Require = false
		opts.Replace = false
		opts.Toolchain = true
		req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

		res, err := req.Run(context.Background())
		require.NoError(t, err)
		require.Len(t, res.Changes, 1)
		require.Equal(t, DirectiveToolchain, res.Changes[0].Directive)
		require.Equal(t, "go1.21.0", res.Changes[0].OldVersion)
		require.Equal(t, "go1.22.0", res.Changes[0].NewVersion)
	})
}

// TestSyncModReq_WritesFile verifies that when changes are present, the dest
// go.mod file on disk is updated with the new content.
func TestSyncModReq_WritesFile(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require github.com/foo/bar v1.2.0
`
	const dstMod = `module github.com/example/dst

go 1.21

require github.com/foo/bar v1.0.0
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	req := &SyncModReq{A: src, B: dst, DiffOpts: DefaultSyncDiffOpts()}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, res.Changes)

	written, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	require.Contains(t, string(written), "v1.2.0", "dest file should have updated version")
	require.NotContains(t, string(written), "v1.0.0", "dest file should not have old version")
}

// TestSyncModReq_NoChanges_NoWrite verifies that when source and dest have
// identical content, no file write occurs and the result has no changes.
func TestSyncModReq_NoChanges_NoWrite(t *testing.T) {
	t.Parallel()

	const mod = `module github.com/example/mod

go 1.21

require github.com/foo/bar v1.0.0
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(mod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(mod), 0o644))

	// Record the file modification time before running.
	infoBefore, err := os.Stat(dstPath)
	require.NoError(t, err)

	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))
	req := &SyncModReq{A: src, B: dst, DiffOpts: DefaultSyncDiffOpts()}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.Empty(t, res.Changes)
	require.Empty(t, res.Error)

	// Mod time should not have changed.
	infoAfter, err := os.Stat(dstPath)
	require.NoError(t, err)
	require.Equal(t, infoBefore.ModTime(), infoAfter.ModTime(), "file should not have been written")
}

// TestSyncModReq_NilSource verifies that a nil source (A) returns an error
// immediately without panicking.
func TestSyncModReq_NilSource(t *testing.T) {
	t.Parallel()

	req := &SyncModReq{
		A:        nil,
		B:        &ModSource{Name: "dst/go.mod"},
		DiffOpts: DefaultSyncDiffOpts(),
	}

	res, err := req.Run(context.Background())
	require.Error(t, err)
	require.Nil(t, res)
}

// TestSyncModReq_NilDest verifies that a nil dest (B) returns an error
// immediately without panicking.
func TestSyncModReq_NilDest(t *testing.T) {
	t.Parallel()

	req := &SyncModReq{
		A:        &ModSource{Name: "src/go.mod", Data: []byte("module m\ngo 1.21\n")},
		B:        nil,
		DiffOpts: DefaultSyncDiffOpts(),
	}

	res, err := req.Run(context.Background())
	require.Error(t, err)
	require.Nil(t, res)
}

// TestSyncModReq_ParseError verifies that invalid source go.mod content returns
// an error whose message contains "diffing".
func TestSyncModReq_ParseError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte("module github.com/example/dst\ngo 1.21\n"), 0o644))
	bSrc := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, bSrc))

	req := &SyncModReq{
		A:        &ModSource{Name: "src/go.mod", Data: []byte("invalid go mod content !@#")},
		B:        bSrc,
		DiffOpts: DefaultSyncDiffOpts(),
	}

	res, err := req.Run(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "diffing")
	// res may be non-nil with Error set, or nil — both are acceptable.
	_ = res
}

// TestSyncModReq_GoDirective_NilSource verifies that when the source file has
// no go directive, syncing the go directive (--go=true) is a no-op and the
// destination file is not modified.
func TestSyncModReq_GoDirective_NilSource(t *testing.T) {
	t.Parallel()

	// Source has no go directive at all.
	const srcMod = "module github.com/example/src\n"
	const dstMod = "module github.com/example/dst\n\ngo 1.21\n"

	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	opts.Require = false
	opts.Replace = false
	opts.Go = true
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.Empty(t, res.Changes, "no go directive in source means no sync")

	written, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	require.Contains(t, string(written), "go 1.21", "dest go directive should be unchanged")
}

// TestSyncModReq_ToolchainDirective_NilSource verifies that when the source
// file has no toolchain directive, syncing the toolchain (--toolchain=true) is
// a no-op and the destination file is not modified.
func TestSyncModReq_ToolchainDirective_NilSource(t *testing.T) {
	t.Parallel()

	// Source has no toolchain directive; dest has one.
	const srcMod = "module github.com/example/src\n\ngo 1.21\n"
	const dstMod = "module github.com/example/dst\n\ngo 1.21\n\ntoolchain go1.21.0\n"

	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	opts := DefaultSyncDiffOpts()
	opts.Require = false
	opts.Replace = false
	opts.Toolchain = true
	req := &SyncModReq{A: src, B: dst, DiffOpts: opts}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.Empty(t, res.Changes, "no toolchain directive in source means no sync")

	written, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	require.Contains(t, string(written), "toolchain go1.21.0", "dest toolchain should be unchanged")
}

// TestSyncModReq_DryRun verifies that when DryRun is true, changes are
// computed and returned in res.Changes, but the destination file on disk is
// not modified.
func TestSyncModReq_DryRun(t *testing.T) {
	t.Parallel()

	const srcMod = `module github.com/example/src

go 1.21

require (
	github.com/foo/bar v1.2.0
	github.com/foo/new v1.3.0
)
`
	const dstMod = `module github.com/example/dst

go 1.21

require github.com/foo/bar v1.0.0
`
	src := &ModSource{Name: "source/go.mod", Data: []byte(srcMod)}

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(dstPath, []byte(dstMod), 0o644))
	dst := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(dstPath, dst))

	req := &SyncModReq{A: src, B: dst, DiffOpts: DefaultSyncDiffOpts(), DryRun: true}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.NotNil(t, res)

	// Changes must be populated even in dry-run mode.
	require.NotEmpty(t, res.Changes, "dry run should still compute changes")

	var updated, added *SyncChange
	for _, c := range res.Changes {
		switch c.Module {
		case "github.com/foo/bar":
			updated = c
		case "github.com/foo/new":
			added = c
		}
	}
	require.NotNil(t, updated, "bar update should be present in changes")
	require.Equal(t, SyncActionUpdated, updated.Action)
	require.Equal(t, "v1.0.0", updated.OldVersion)
	require.Equal(t, "v1.2.0", updated.NewVersion)

	require.NotNil(t, added, "new module addition should be present in changes")
	require.Equal(t, SyncActionAdded, added.Action)

	// The file on disk must be unchanged.
	written, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	require.Equal(t, dstMod, string(written), "dry run must not write to disk")
}

// TestSyncModRes_ToTable_WithChanges verifies that ToTable produces a non-empty
// table containing the directive, action, and module columns when changes are present.
func TestSyncModRes_ToTable_WithChanges(t *testing.T) {
	t.Parallel()

	res := &SyncModRes{
		Path: "b/go.mod",
		Changes: []*SyncChange{
			{
				Directive:  DirectiveRequire,
				Module:     "github.com/foo/bar",
				OldVersion: "v1.0.0",
				NewVersion: "v1.2.0",
				Action:     SyncActionUpdated,
			},
			{
				Directive:  DirectiveRequire,
				Module:     "github.com/foo/new",
				NewVersion: "v1.3.0",
				Action:     SyncActionAdded,
			},
		},
	}

	out := res.ToTable()
	require.NotEmpty(t, out, "ToTable should produce output when changes are present")
	require.Contains(t, out, "github.com/foo/bar", "table must contain updated module path")
	require.Contains(t, out, "github.com/foo/new", "table must contain added module path")
	require.Contains(t, out, string(SyncActionUpdated), "table must contain the updated action")
	require.Contains(t, out, string(SyncActionAdded), "table must contain the added action")
	require.Contains(t, out, "v1.0.0", "table must contain old version")
	require.Contains(t, out, "v1.2.0", "table must contain new version")
}

// TestSyncModRes_ToTable_NoChanges verifies that ToTable produces a summary
// row with a zero change count and no error when the result has no changes.
func TestSyncModRes_ToTable_NoChanges(t *testing.T) {
	t.Parallel()

	res := &SyncModRes{Path: "b/go.mod"}

	out := res.ToTable()
	require.NotEmpty(t, out, "ToTable should produce a summary row even with no changes")
	require.Contains(t, out, "0", "change count must be zero in the no-changes summary row")
}

// TestSyncModRes_ToTable_NilRes verifies that calling ToTable on a nil pointer
// returns an empty string without panicking.
func TestSyncModRes_ToTable_NilRes(t *testing.T) {
	t.Parallel()

	var res *SyncModRes
	require.Empty(t, res.ToTable(), "nil receiver must return an empty string")
}

// TestSyncModRes_ToMarkdown_WithChanges verifies that ToMarkdown produces a
// pipe-delimited markdown table containing all expected columns and values.
func TestSyncModRes_ToMarkdown_WithChanges(t *testing.T) {
	t.Parallel()

	res := &SyncModRes{
		Path: "b/go.mod",
		Changes: []*SyncChange{
			{
				Directive:  DirectiveReplace,
				Module:     "github.com/foo/bar",
				OldVersion: "v1.1.0",
				NewVersion: "v1.2.0",
				Action:     SyncActionUpdated,
			},
		},
	}

	out := res.ToMarkdown()
	require.NotEmpty(t, out, "ToMarkdown should produce output when changes are present")
	require.Contains(t, out, "|", "markdown table must use pipe delimiters")
	require.Contains(t, out, "github.com/foo/bar", "markdown table must contain the module path")
	require.Contains(t, out, string(DirectiveReplace), "markdown table must contain the directive")
	require.Contains(t, out, string(SyncActionUpdated), "markdown table must contain the action")
	require.Contains(t, out, "v1.1.0", "markdown table must contain old version")
	require.Contains(t, out, "v1.2.0", "markdown table must contain new version")
}

// TestSyncModRes_ToMarkdown_NoChanges verifies that ToMarkdown produces a
// summary row with a zero change count when there are no changes.
func TestSyncModRes_ToMarkdown_NoChanges(t *testing.T) {
	t.Parallel()

	res := &SyncModRes{Path: "b/go.mod"}

	out := res.ToMarkdown()
	require.NotEmpty(t, out, "ToMarkdown should produce a summary row even with no changes")
	require.Contains(t, out, "|", "markdown output must still use pipe delimiters")
	require.Contains(t, out, "0", "change count must be zero in the no-changes summary row")
}

// TestSyncModRes_ToMarkdown_NilRes verifies that calling ToMarkdown on a nil
// pointer returns an empty string without panicking.
func TestSyncModRes_ToMarkdown_NilRes(t *testing.T) {
	t.Parallel()

	var res *SyncModRes
	require.Empty(t, res.ToMarkdown(), "nil receiver must return an empty string")
}

// TestSyncModRes_ToJSON_WithChanges verifies that ToJSON marshals a result with
// changes to valid JSON containing the path, directive, module, and action fields.
func TestSyncModRes_ToJSON_WithChanges(t *testing.T) {
	t.Parallel()

	res := &SyncModRes{
		Path: "b/go.mod",
		Changes: []*SyncChange{
			{
				Directive:  DirectiveRequire,
				Module:     "github.com/foo/bar",
				OldVersion: "v1.0.0",
				NewVersion: "v1.2.0",
				Action:     SyncActionUpdated,
			},
		},
	}

	b, err := res.ToJSON()
	require.NoError(t, err, "ToJSON must not error on a valid result")
	require.NotEmpty(t, b, "ToJSON must produce non-empty output")

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed), "ToJSON output must be valid JSON")

	require.Equal(t, "b/go.mod", parsed["path"], "JSON must include the dest path")

	changes, ok := parsed["changes"].([]any)
	require.True(t, ok, "JSON must contain a changes array")
	require.Len(t, changes, 1, "changes array should have exactly one entry")

	change := changes[0].(map[string]any)
	require.Equal(t, "require", change["directive"])
	require.Equal(t, "github.com/foo/bar", change["module"])
	require.Equal(t, "updated", change["action"])
	require.Equal(t, "v1.2.0", change["new_version"])
	require.Equal(t, "v1.0.0", change["old_version"])
}

// TestSyncModRes_ToJSON_NoChanges verifies that ToJSON omits the changes key
// when the result has no changes (omitempty on the Changes field).
func TestSyncModRes_ToJSON_NoChanges(t *testing.T) {
	t.Parallel()

	res := &SyncModRes{Path: "b/go.mod"}

	b, err := res.ToJSON()
	require.NoError(t, err, "ToJSON must not error when there are no changes")
	require.NotEmpty(t, b)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed))
	require.Equal(t, "b/go.mod", parsed["path"])
	// Changes is omitempty; when empty the key should be absent.
	_, hasChanges := parsed["changes"]
	require.False(t, hasChanges, "changes key must be absent when there are no changes (omitempty)")
}

// TestSyncModRes_ToJSON_NilRes verifies that calling ToJSON on a nil pointer
// returns an error without panicking.
func TestSyncModRes_ToJSON_NilRes(t *testing.T) {
	t.Parallel()

	var res *SyncModRes
	b, err := res.ToJSON()
	require.Error(t, err, "nil receiver must return an error")
	require.Nil(t, b, "nil receiver must return nil bytes")
}

// TestSetUpGoModSourceFromPath_ValidFile verifies that a valid file is read
// into source.Data and source.Name is set to the provided path.
func TestSetUpGoModSourceFromPath_ValidFile(t *testing.T) {
	t.Parallel()

	const content = "module github.com/example/mod\n\ngo 1.21\n"
	dir := t.TempDir()
	p := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(p, []byte(content), 0o644))

	src := &ModSource{}
	require.NoError(t, SetUpGoModSourceFromPath(p, src))

	require.Equal(t, p, src.Name, "Name must be set to the provided path")
	require.Equal(t, []byte(content), src.Data, "Data must contain the file contents")
}

// TestSetUpGoModSourceFromPath_NonExistentFile verifies that a missing file
// returns an error.
func TestSetUpGoModSourceFromPath_NonExistentFile(t *testing.T) {
	t.Parallel()

	src := &ModSource{}
	err := SetUpGoModSourceFromPath("/nonexistent/path/go.mod", src)
	require.Error(t, err, "missing file must return an error")
}

// TestSetUpGoModSourceFromPath_NilSource verifies that passing nil source
// returns an error immediately.
func TestSetUpGoModSourceFromPath_NilSource(t *testing.T) {
	t.Parallel()

	err := SetUpGoModSourceFromPath("go.mod", nil)
	require.Error(t, err, "nil source must return an error")
	require.Contains(t, err.Error(), "mod source")
}
