// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// groupPackagesFixture is a workspace with a root module and two nested
// modules, mirroring vault, api and api/auth/approle. It also has the inputs
// that aren't packages of a workspace module: a non-Go file, a directory
// without Go files, an ignored testdata directory and a nested module that
// go.work doesn't use.
var groupPackagesFixture = map[string]string{
	"go.work":                     "go 1.21\n\nuse (\n\t.\n\t./api\n\t./api/auth/approle\n)\n",
	"go.mod":                      "module example.com/vault\n\ngo 1.21\n",
	"main.go":                     "package main\n",
	"README.md":                   "# vault\n",
	"command/command.go":          "package command\n",
	"docs/README.md":              "# docs\n",
	"vault/core.go":               "package vault\n",
	"vault/core_test.go":          "package vault\n",
	"vault/testdata/fixture.go":   "package testdata\n",
	"other/go.mod":                "module example.com/vault/other\n\ngo 1.21\n",
	"other/other.go":              "package other\n",
	"api/go.mod":                  "module example.com/vault/api\n\ngo 1.21\n",
	"api/client.go":               "package api\n",
	"api/auth/approle/go.mod":     "module example.com/vault/api/auth/approle\n\ngo 1.21\n",
	"api/auth/approle/approle.go": "package approle\n",
}

// TestGroupPackagesReq_Run verifies that import paths are grouped by the
// module that contains them, in go.work order, keeping input order and
// dropping duplicates. This is how the test matrix finds the module directory
// to test each partition's packages from.
func TestGroupPackagesReq_Run(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, groupPackagesFixture)

	res, err := (&GroupPackagesReq{
		GoWork: filepath.Join(root, "go.work"),
		Packages: []string{
			"example.com/vault/api/auth/approle",
			"example.com/vault/vault",
			"example.com/vault/api",
			"example.com/vault/command",
			"example.com/vault/vault",
		},
	}).Run(t.Context())
	require.NoError(t, err)
	require.Equal(t, []*ModulePackages{
		{Path: "example.com/vault", Dir: ".", Packages: []string{"example.com/vault/vault", "example.com/vault/command"}},
		{Path: "example.com/vault/api", Dir: "api", Packages: []string{"example.com/vault/api"}},
		{Path: "example.com/vault/api/auth/approle", Dir: "api/auth/approle", Packages: []string{"example.com/vault/api/auth/approle"}},
	}, res.List.Modules)

	require.Contains(t, res.ToMarkdown(), "| . | example.com/vault/vault<br/>example.com/vault/command |")
}

// TestGroupPackagesReq_Run_Paths verifies that directories and .go files,
// relative to Dir or absolute, resolve to the import path of their package. Two
// files in one package, or a file and its directory, yield the package once.
func TestGroupPackagesReq_Run_Paths(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, groupPackagesFixture)

	res, err := (&GroupPackagesReq{
		GoWork: filepath.Join(root, "go.work"),
		Dir:    root,
		Packages: []string{
			"vault/core.go",
			"api/auth/approle",
			"vault/core_test.go",
			"./vault",
			filepath.Join(root, "api", "client.go"),
			".",
			"example.com/vault/command",
		},
	}).Run(t.Context())
	require.NoError(t, err)
	require.Equal(t, []*ModulePackages{
		{Path: "example.com/vault", Dir: ".", Packages: []string{"example.com/vault/vault", "example.com/vault", "example.com/vault/command"}},
		{Path: "example.com/vault/api", Dir: "api", Packages: []string{"example.com/vault/api"}},
		{Path: "example.com/vault/api/auth/approle", Dir: "api/auth/approle", Packages: []string{"example.com/vault/api/auth/approle"}},
	}, res.List.Modules)
}

// TestGroupPackagesReq_Run_Invalid verifies that entries that aren't a package
// directory in a workspace module fail with an explanation instead of being
// grouped. Grouping them would hand go test a package that doesn't exist or
// that the module can't build.
func TestGroupPackagesReq_Run_Invalid(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		entry  string
		errHas string
	}{
		"import path outside the workspace": {
			entry:  "example.com/other",
			errHas: "example.com/other doesn't exist and isn't a package in any module in go.work",
		},
		"import path sharing a module path prefix": {
			entry:  "example.com/vaultx",
			errHas: "example.com/vaultx doesn't exist and isn't a package in any module in go.work",
		},
		"import path without a directory": {
			entry:  "example.com/vault/does/not/exist",
			errHas: "module . has no directory does/not/exist",
		},
		"relative path that doesn't exist": {
			entry:  "vault/nope.go",
			errHas: "vault/nope.go doesn't exist",
		},
		"non-Go file": {
			entry:  "README.md",
			errHas: "README.md is not a directory or a .go file",
		},
		"directory without Go files": {
			entry:  "docs",
			errHas: "module . has no Go files in docs",
		},
		"ignored directory": {
			entry:  "vault/testdata/fixture.go",
			errHas: "the go command ignores vault/testdata in module .",
		},
		"module that go.work doesn't use": {
			entry:  "other/other.go",
			errHas: "other is in a module that isn't in go.work",
		},
		"path outside the workspace": {
			entry:  "..",
			errHas: ".. is outside every module in go.work",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			writeFiles(t, root, groupPackagesFixture)

			_, err := (&GroupPackagesReq{
				GoWork:   filepath.Join(root, "go.work"),
				Dir:      root,
				Packages: []string{"example.com/vault/api", test.entry},
			}).Run(t.Context())
			require.ErrorContains(t, err, test.errHas)
			require.NotContains(t, err.Error(), "example.com/vault/api:")
		})
	}
}

// TestGroupPackagesReq_Run_Empty verifies that empty input produces an empty
// module list instead of an error.
func TestGroupPackagesReq_Run_Empty(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, groupPackagesFixture)

	res, err := (&GroupPackagesReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)

	b, err := res.ToJSON()
	require.NoError(t, err)
	require.JSONEq(t, `{"modules":[]}`, string(b))
}
