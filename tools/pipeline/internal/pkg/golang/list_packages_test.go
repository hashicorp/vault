// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// listPackagesFixture is a two-module workspace where module a has a package
// that only exists with the fixture build tag.
var listPackagesFixture = map[string]string{
	"go.work":         "go 1.21\n\nuse (\n\t./a\n\t./b\n)\n",
	"a/go.mod":        "module example.com/a\n\ngo 1.21\n",
	"a/a.go":          "package a\n",
	"a/tagged/t.go":   "//go:build fixture\n\npackage tagged\n",
	"b/go.mod":        "module example.com/b\n\ngo 1.21\n",
	"b/b.go":          "package b\n",
	"b/internal/x.go": "package internal\n",
}

// TestListPackagesReq_Run verifies that each module's packages are listed
// from its own directory and that build tags change the package set.
func TestListPackagesReq_Run(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, listPackagesFixture)
	goWork := filepath.Join(root, "go.work")

	res, err := (&ListPackagesReq{GoWork: goWork}).Run(t.Context())
	require.NoError(t, err)
	require.Equal(t, []*ModulePackages{
		{Path: "example.com/a", Dir: "a", Packages: []string{"example.com/a"}},
		{Path: "example.com/b", Dir: "b", Packages: []string{"example.com/b", "example.com/b/internal"}},
	}, res.List.Modules)

	tagged, err := (&ListPackagesReq{GoWork: goWork, Tags: ",fixture"}).Run(t.Context())
	require.NoError(t, err)
	require.Equal(t, []string{"example.com/a", "example.com/a/tagged"}, tagged.List.Modules[0].Packages)

	require.Contains(t, res.ToMarkdown(), "| b | example.com/b/internal |")
}

// TestListPackagesReq_Run_ModuleFilter verifies that --module limits listing
// to the selected modules and rejects directories that aren't modules.
func TestListPackagesReq_Run_ModuleFilter(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, listPackagesFixture)
	goWork := filepath.Join(root, "go.work")

	res, err := (&ListPackagesReq{GoWork: goWork, Modules: []string{"b"}}).Run(t.Context())
	require.NoError(t, err)
	require.Len(t, res.List.Modules, 1)
	require.Equal(t, "b", res.List.Modules[0].Dir)

	_, err = (&ListPackagesReq{GoWork: goWork, Modules: []string{"c"}}).Run(t.Context())
	require.ErrorContains(t, err, "c is not a module")
}

// TestListPackagesReq_Run_EmptyModule verifies that a module without packages
// is reported with an empty list rather than dropped or null.
func TestListPackagesReq_Run_EmptyModule(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work":       "go 1.21\n\nuse ./a\n",
		"a/go.mod":      "module example.com/a\n\ngo 1.21\n",
		"a/tagged/t.go": "//go:build fixture\n\npackage tagged\n",
	})

	res, err := (&ListPackagesReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)

	b, err := res.ToJSON()
	require.NoError(t, err)
	require.JSONEq(t, `{"modules":[{"path":"example.com/a","dir":"a","packages":[]}]}`, string(b))
}

// TestListPackagesReq_Run_LoadError verifies that a module that fails to load
// fails the command, so its packages can't silently drop out of CI.
func TestListPackagesReq_Run_LoadError(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work":  "go 1.21\n\nuse ./a\n",
		"a/go.mod": "module example.com/a\n\ngo 1.21\n",
		"a/a.go":   "this is not go\n",
	})

	_, err := (&ListPackagesReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.ErrorContains(t, err, "listing packages in module a")
	require.ErrorContains(t, err, "expected 'package'")
}
