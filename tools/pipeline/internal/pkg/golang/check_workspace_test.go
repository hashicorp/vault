// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
)

// checkWorkspaceFixture mirrors the vault layout: a root module, api, sdk and a
// nested api/auth/approle module.
func checkWorkspaceFixture() *Workspace {
	return &Workspace{
		GoWork: "/ws/go.work",
		Root:   "/ws",
		Modules: []*WorkspaceModule{
			{Path: "example.com/vault", Dir: ".", GoMod: "go.mod"},
			{Path: "example.com/vault/api", Dir: "api", GoMod: "api/go.mod"},
			{Path: "example.com/vault/api/auth/approle", Dir: "api/auth/approle", GoMod: "api/auth/approle/go.mod"},
			{Path: "example.com/vault/sdk", Dir: "sdk", GoMod: "sdk/go.mod"},
		},
	}
}

// TestFindReplaceViolations verifies each way a module can fail to build
// against the on-branch code of the workspace modules it requires, and that
// correct local replaces pass.
func TestFindReplaceViolations(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		moduleDir string
		goMod     string
		want      []*WorkspaceViolation
	}{
		"local replaces": {
			moduleDir: ".",
			goMod: `module example.com/vault
go 1.21
require (
	example.com/vault/api v1.0.0
	example.com/vault/sdk v0.1.0
	example.com/external v1.0.0
)
replace example.com/vault/api => ./api
replace example.com/vault/sdk => ./sdk/
`,
			want: []*WorkspaceViolation{},
		},
		"nested module replaces its parent": {
			moduleDir: "api/auth/approle",
			goMod: `module example.com/vault/api/auth/approle
go 1.21
require example.com/vault/api v1.0.0
replace example.com/vault/api => ../..
`,
			want: []*WorkspaceViolation{},
		},
		"missing replace": {
			moduleDir: "sdk",
			goMod: `module example.com/vault/sdk
go 1.21
require example.com/vault/api v1.0.0
`,
			want: []*WorkspaceViolation{{
				Kind:   WorkspaceViolationMissingLocalReplace,
				Module: "sdk",
				Detail: "requires example.com/vault/api v1.0.0 without a replace; use 'replace example.com/vault/api => ../api'",
			}},
		},
		"versioned old path": {
			moduleDir: "sdk",
			goMod: `module example.com/vault/sdk
go 1.21
require example.com/vault/api v1.0.0
replace example.com/vault/api v1.0.0 => ../api
`,
			want: []*WorkspaceViolation{{
				Kind:   WorkspaceViolationInvalidLocalReplace,
				Module: "sdk",
				Detail: "replace example.com/vault/api v1.0.0 => ../api only applies to version v1.0.0; use 'replace example.com/vault/api => ../api'",
			}},
		},
		"module replace": {
			moduleDir: "sdk",
			goMod: `module example.com/vault/sdk
go 1.21
require example.com/vault/api v1.0.0
replace example.com/vault/api => example.com/fork/api v1.1.0
`,
			want: []*WorkspaceViolation{{
				Kind:   WorkspaceViolationInvalidLocalReplace,
				Module: "sdk",
				Detail: "replace example.com/vault/api => example.com/fork/api v1.1.0 isn't a directory replace; use 'replace example.com/vault/api => ../api'",
			}},
		},
		"wrong directory": {
			moduleDir: "sdk",
			goMod: `module example.com/vault/sdk
go 1.21
require example.com/vault/api v1.0.0
replace example.com/vault/api => ../sdk
`,
			want: []*WorkspaceViolation{{
				Kind:   WorkspaceViolationInvalidLocalReplace,
				Module: "sdk",
				Detail: "replace example.com/vault/api => ../sdk resolves to sdk, not api; use 'replace example.com/vault/api => ../api'",
			}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ws := checkWorkspaceFixture()
			m, err := ws.ModuleByDir(test.moduleDir)
			require.NoError(t, err)
			modFile, err := modfile.Parse("go.mod", []byte(test.goMod), nil)
			require.NoError(t, err)

			require.Equal(t, test.want, findReplaceViolations(ws, m, modFile))
		})
	}
}

// TestFindPublishedReplaceViolations verifies that a published module may
// require workspace modules by version but may not replace them, since its
// importers ignore the replace. Replaces of modules outside the workspace are
// the module's own business and aren't reported.
func TestFindPublishedReplaceViolations(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		moduleDir string
		goMod     string
		want      []*WorkspaceViolation
	}{
		"requires a released version": {
			moduleDir: "sdk",
			goMod: `module example.com/vault/sdk
go 1.21
require example.com/vault/api v1.0.0
`,
			want: []*WorkspaceViolation{},
		},
		"replaces a module outside the workspace": {
			moduleDir: "sdk",
			goMod: `module example.com/vault/sdk
go 1.21
require example.com/external v1.0.0
replace example.com/external => example.com/fork v1.1.0
`,
			want: []*WorkspaceViolation{},
		},
		"local replace of a workspace module": {
			moduleDir: "api/auth/approle",
			goMod: `module example.com/vault/api/auth/approle
go 1.21
require example.com/vault/api v1.0.0
replace example.com/vault/api => ../..
`,
			want: []*WorkspaceViolation{{
				Kind:   WorkspaceViolationPublishedReplace,
				Module: "api/auth/approle",
				Detail: "published module replaces example.com/vault/api, which its importers ignore; remove 'replace example.com/vault/api => ../..' and require a released version of example.com/vault/api",
			}},
		},
		// Even a replace of a module that isn't required is reported, so the
		// published go.mod never carries a replace that tidy would leave behind.
		"versioned replace of an unrequired workspace module": {
			moduleDir: "sdk",
			goMod: `module example.com/vault/sdk
go 1.21
replace example.com/vault/api v1.0.0 => example.com/fork/api v1.1.0
`,
			want: []*WorkspaceViolation{{
				Kind:   WorkspaceViolationPublishedReplace,
				Module: "sdk",
				Detail: "published module replaces example.com/vault/api, which its importers ignore; remove 'replace example.com/vault/api v1.0.0 => example.com/fork/api v1.1.0' and require a released version of example.com/vault/api",
			}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ws := checkWorkspaceFixture()
			m, err := ws.ModuleByDir(test.moduleDir)
			require.NoError(t, err)
			modFile, err := modfile.Parse("go.mod", []byte(test.goMod), nil)
			require.NoError(t, err)

			require.Equal(t, test.want, findPublishedReplaceViolations(ws, m, modFile))
		})
	}
}

// TestPublishedModules verifies that published module patterns match module
// directories, including nested modules through globs, and that malformed or
// unmatched patterns are errors rather than silently exempting nothing.
func TestPublishedModules(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		patterns []string
		want     []string
		err      string
	}{
		"none": {
			want: []string{},
		},
		"directories and globs": {
			patterns: []string{"./sdk", "api", "api/auth/*/"},
			want:     []string{"api", "api/auth/approle", "sdk"},
		},
		"glob doesn't cross directories": {
			patterns: []string{"api/*"},
			err:      `published module pattern "api/*" matches no module in /ws/go.work`,
		},
		"no match": {
			patterns: []string{"sdk", "tools"},
			err:      `published module pattern "tools" matches no module in /ws/go.work`,
		},
		"malformed": {
			patterns: []string{"api/["},
			err:      `invalid published module pattern "api/[": syntax error in pattern`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			published, err := publishedModules(checkWorkspaceFixture(), test.patterns)
			if test.err != "" {
				require.EqualError(t, err, test.err)
				return
			}
			require.NoError(t, err)

			got := []string{}
			for m := range published {
				got = append(got, m.Dir)
			}
			require.ElementsMatch(t, test.want, got)
		})
	}
}

// TestFindUnregisteredModules verifies that committed go.mod files outside
// go.work are reported, except in directories the go command ignores.
func TestFindUnregisteredModules(t *testing.T) {
	t.Parallel()

	ws := checkWorkspaceFixture()
	got := findUnregisteredModules(ws, []string{
		"go.mod",
		"api/go.mod",
		"api/auth/approle/go.mod",
		"sdk/go.mod",
		"tools/new/go.mod",
		"sdk/helper/testdata/go.mod",
		"vendor/example.com/x/go.mod",
		".github/go.mod",
		"internal/_fixture/go.mod",
	})

	require.Equal(t, []*WorkspaceViolation{{
		Kind:   WorkspaceViolationUnregisteredModule,
		Module: "tools/new",
		Detail: "tools/new/go.mod is committed but go.work has no 'use ./tools/new'",
	}}, got)
}

// TestCheckWorkspaceReq_Run verifies the whole check against a git
// repository, including that only committed go.mod files are considered.
func TestCheckWorkspaceReq_Run(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work":    "go 1.21\n\nuse (\n\t.\n\t./api\n\t./sdk\n)\n",
		"go.mod":     "module example.com/vault\n\ngo 1.21\n\nrequire example.com/vault/api v1.0.0\n\nreplace example.com/vault/api => ./api\n",
		"api/go.mod": "module example.com/vault/api\n\ngo 1.21\n",
		"sdk/go.mod": "module example.com/vault/sdk\n\ngo 1.21\n\nrequire example.com/vault/api v1.0.0\n",
		// Committed but not in go.work.
		"tools/new/go.mod": "module example.com/vault/tools/new\n\ngo 1.21\n",
		// Ignored by the go command.
		"sdk/testdata/go.mod": "module example.com/fixture\n\ngo 1.21\n",
	})
	gitCommitAll(t, root)
	// Not committed, so not reported.
	writeFiles(t, root, map[string]string{
		"scratch/go.mod": "module example.com/scratch\n\ngo 1.21\n",
	})

	res, err := (&CheckWorkspaceReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)
	require.False(t, res.Success())
	require.Equal(t, "go.work and the go.mod files have 2 problems", res.String())
	require.Equal(t, []*WorkspaceViolation{
		{
			Kind:   WorkspaceViolationUnregisteredModule,
			Module: "tools/new",
			Detail: "tools/new/go.mod is committed but go.work has no 'use ./tools/new'",
		},
		{
			Kind:   WorkspaceViolationMissingLocalReplace,
			Module: "sdk",
			Detail: "requires example.com/vault/api v1.0.0 without a replace; use 'replace example.com/vault/api => ../api'",
		},
	}, res.Violations)

	// The output leads with the summary, and the table only has columns that
	// mean something for every row.
	require.Equal(t, `go.work and the go.mod files have 2 problems.

 MODULE     PROBLEM
 tools/new  tools/new/go.mod is committed but go.work has no 'use ./tools/new'
 sdk        requires example.com/vault/api v1.0.0 without a replace; use 'replace example.com/vault/api => ../api'`, res.ToTable())
	require.Equal(t, `# Go workspace check

go.work and the go.mod files have 2 problems.

| Module | Problem |
| --- | --- |
| tools/new | tools/new/go.mod is committed but go.work has no 'use ./tools/new' |
| sdk | requires example.com/vault/api v1.0.0 without a replace; use 'replace example.com/vault/api => ../api' |`, res.ToMarkdown())
}

// TestCheckWorkspaceReq_Run_Success verifies that a consistent workspace
// passes, including the root go.mod matched by the pathspec.
func TestCheckWorkspaceReq_Run_Success(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work":    "go 1.21\n\nuse (\n\t.\n\t./api\n)\n",
		"go.mod":     "module example.com/vault\n\ngo 1.21\n\nrequire example.com/vault/api v1.0.0\n\nreplace example.com/vault/api => ./api\n",
		"api/go.mod": "module example.com/vault/api\n\ngo 1.21\n",
	})
	gitCommitAll(t, root)

	res, err := (&CheckWorkspaceReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)
	require.True(t, res.Success(), res.ToTable())
	// Success is a single sentence, without an empty table.
	require.Equal(t, "go.work and the go.mod files of 2 modules agree.", res.ToTable())
	require.Equal(t, "# Go workspace check\n\ngo.work and the go.mod files of 2 modules agree.", res.ToMarkdown())

	b, err := res.ToJSON()
	require.NoError(t, err)
	require.Contains(t, string(b), `"violations":[]`)
}

// TestCheckWorkspaceReq_Run_Published verifies that published modules are
// checked for replaces instead of for missing local replaces, while the
// unpublished modules in the same workspace still need their local replaces.
func TestCheckWorkspaceReq_Run_Published(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work":    "go 1.21\n\nuse (\n\t.\n\t./api\n\t./api/auth/approle\n\t./sdk\n)\n",
		"go.mod":     "module example.com/vault\n\ngo 1.21\n\nrequire example.com/vault/sdk v1.0.0\n",
		"api/go.mod": "module example.com/vault/api\n\ngo 1.21\n",
		// Published and requires a release of api, which is what importers get.
		"sdk/go.mod": "module example.com/vault/sdk\n\ngo 1.21\n\nrequire example.com/vault/api v1.0.0\n",
		// Published but replaces api, which importers ignore.
		"api/auth/approle/go.mod": "module example.com/vault/api/auth/approle\n\ngo 1.21\n\nrequire example.com/vault/api v1.0.0\n\nreplace example.com/vault/api => ../..\n",
	})
	gitCommitAll(t, root)

	res, err := (&CheckWorkspaceReq{
		GoWork:    filepath.Join(root, "go.work"),
		Published: []string{"sdk", "api/auth/*"},
	}).Run(t.Context())
	require.NoError(t, err)
	require.Equal(t, []*WorkspaceViolation{
		{
			Kind:   WorkspaceViolationMissingLocalReplace,
			Module: ".",
			Detail: "requires example.com/vault/sdk v1.0.0 without a replace; use 'replace example.com/vault/sdk => ./sdk'",
		},
		{
			Kind:   WorkspaceViolationPublishedReplace,
			Module: "api/auth/approle",
			Detail: "published module replaces example.com/vault/api, which its importers ignore; remove 'replace example.com/vault/api => ../..' and require a released version of example.com/vault/api",
		},
	}, res.Violations)

	// A pattern that matches nothing fails the check rather than passing it.
	_, err = (&CheckWorkspaceReq{
		GoWork:    filepath.Join(root, "go.work"),
		Published: []string{"sdk", "api/auth"},
	}).Run(t.Context())
	require.ErrorContains(t, err, `published module pattern "api/auth" matches no module`)
}
