// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/require"
)

// TestCheckGoModDiffReq_validate_workspaceModules verifies that workspace
// modules can't be combined with explicit go.mod paths, and that a go.work
// path is only accepted together with workspace modules, since it would
// otherwise be silently ignored.
func TestCheckGoModDiffReq_validate_workspaceModules(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		paths            []string
		workspaceModules bool
		goWork           string
		errHas           string
	}{
		"paths": {
			paths: []string{"go.mod", "sdk/go.mod"},
		},
		"workspace modules": {
			workspaceModules: true,
		},
		"workspace modules with go.work": {
			workspaceModules: true,
			goWork:           "go.work",
		},
		"workspace modules with paths": {
			paths:            []string{"go.mod"},
			workspaceModules: true,
			errHas:           "cannot be combined",
		},
		"go.work without workspace modules": {
			goWork: "go.work",
			errHas: "requires workspace modules",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &CheckGoModDiffReq{
				AOwner:           "hashicorp",
				ARepo:            "vault-enterprise",
				ABranch:          "main",
				BOwner:           "hashicorp",
				BRepo:            "vault",
				BBranch:          "ce/main",
				Paths:            test.paths,
				WorkspaceModules: test.workspaceModules,
				GoWork:           test.goWork,
			}

			err := req.validate()
			if test.errHas == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, test.errHas)
		})
	}
}

// TestCheckGoModDiffReq_Run_workspaceModulesListFailure verifies that Run
// fails before cloning anything when it can't list the workspace modules. CI
// relies on this: falling back to the root go.mod would hide drift in every
// other module. The nil git client would panic if Run tried to clone.
func TestCheckGoModDiffReq_Run_workspaceModulesListFailure(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		files  map[string]string
		errHas string
	}{
		"missing go.work": {
			errHas: "reading go.work",
		},
		"go.work without modules": {
			files:  map[string]string{"go.work": "go 1.21\n"},
			errHas: "has no modules",
		},
		"go.work using a directory without go.mod": {
			files:  map[string]string{"go.work": "go 1.21\n\nuse ./sdk\n"},
			errHas: "which has no go.mod",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			for name, content := range test.files {
				require.NoError(t, os.WriteFile(filepath.Join(root, name), []byte(content), 0o644))
			}

			req := &CheckGoModDiffReq{
				AOwner:           "hashicorp",
				ARepo:            "vault-enterprise",
				ABranch:          "main",
				BOwner:           "hashicorp",
				BRepo:            "vault",
				BBranch:          "ce/main",
				WorkspaceModules: true,
				GoWork:           filepath.Join(root, "go.work"),
			}

			res, err := req.Run(t.Context(), nil, nil)
			require.ErrorContains(t, err, test.errHas)
			require.Nil(t, res)
		})
	}
}

// TestWorkspaceGoModPaths verifies that workspace modules become go.mod paths
// relative to the go.work directory, in go.work order, which are the paths
// go-mod-diff checks out on each branch.
func TestWorkspaceGoModPaths(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	for name, content := range map[string]string{
		"go.work":                 "go 1.21\n\nuse (\n\t.\n\t./sdk\n\t./api/auth/approle\n)\n",
		"go.mod":                  "module example.com/root\n\ngo 1.21\n",
		"sdk/go.mod":              "module example.com/root/sdk\n\ngo 1.21\n",
		"api/auth/approle/go.mod": "module example.com/root/api/auth/approle\n\ngo 1.21\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}

	paths, err := workspaceGoModPaths(t.Context(), filepath.Join(root, "go.work"))
	require.NoError(t, err)
	require.Equal(t, []string{"go.mod", "sdk/go.mod", "api/auth/approle/go.mod"}, paths)
}

// TestCheckGoModDiffRes_errors verifies that a path that couldn't be compared,
// such as a go.mod missing on one branch, is shown in the table and fails the
// check alongside the differences of every other path, rather than looking
// like a path without differences.
func TestCheckGoModDiffRes_errors(t *testing.T) {
	t.Parallel()

	diff := func(directive golang.Directive) *golang.Diff {
		return &golang.Diff{
			Directive: directive,
			Diff: &difflib.UnifiedDiff{
				A:        []string{"a\n"},
				B:        []string{"b\n"},
				FromFile: "a",
				ToFile:   "b",
			},
		}
	}
	missing := errors.New("reading sdk/go.mod on B branch: no such file")
	res := &CheckGoModDiffRes{Diffs: []*CheckGoModDiff{
		{Path: "go.mod", ModDiff: golang.ModDiff{diff(golang.DirectiveGo)}},
		{Path: "sdk/go.mod", Err: missing, Error: missing.Error()},
		{Path: "api/go.mod", ModDiff: golang.ModDiff{diff(golang.DirectiveRequire), diff(golang.DirectiveReplace)}},
		{Path: "version/go.mod"},
	}}

	err := res.Err()
	require.ErrorContains(t, err, "3 differences were found")
	require.ErrorIs(t, err, missing)

	tbl, err := res.ToTable(nil)
	require.NoError(t, err)
	require.Contains(t, tbl.RenderMarkdown(), "| sdk/go.mod | error | reading sdk/go.mod on B branch: no such file |")

	// The error message survives JSON; the error value isn't marshaled.
	b, err := res.ToJSON()
	require.NoError(t, err)
	require.Contains(t, string(b), `"error":"reading sdk/go.mod on B branch: no such file"`)
	require.NotContains(t, string(b), `"err"`)

	require.NoError(t, (&CheckGoModDiffRes{Diffs: []*CheckGoModDiff{{Path: "go.mod"}}}).Err())
}
