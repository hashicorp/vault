// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestValidateGoModPath verifies that module writes stay within their directory
// and cannot follow symlinked directories, go.mod files, or go.sum files.
func TestValidateGoModPath(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		path   string
		link   string
		target string
		errHas string
	}{
		"regular module without go.sum": {path: "sdk/go.mod"},
		"root module":                   {path: "go.mod"},
		"missing module":                {path: "missing/go.mod", errHas: "reading"},
		"outside directory":             {path: "../outside/go.mod", errHas: "not inside"},
		"symlinked module directory":    {path: "linked/go.mod", link: "linked", target: ".", errHas: "symlink"},
		"symlinked go.mod":              {path: "sdk/go.mod", link: "sdk/go.mod", target: "go.mod", errHas: "symlink"},
		"symlinked go.sum":              {path: "sdk/go.mod", link: "sdk/go.sum", target: "go.sum", errHas: "symlink"},
		"dangling go.sum symlink":       {path: "sdk/go.mod", link: "sdk/go.sum", target: "missing", errHas: "symlink"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			dir, outside := filepath.Join(root, "repo"), filepath.Join(root, "outside")
			writeFiles(t, dir, map[string]string{
				"go.mod":     "module example.com/root\n\ngo 1.21\n",
				"sdk/go.mod": "module example.com/sdk\n\ngo 1.21\n",
			})
			writeFiles(t, outside, map[string]string{"go.mod": "module example.com/outside\n", "go.sum": ""})
			if test.link != "" {
				link := filepath.Join(dir, test.link)
				if test.link == "sdk/go.mod" {
					require.NoError(t, os.Remove(link))
				}
				require.NoError(t, os.Symlink(filepath.Join(outside, test.target), link))
			}

			err := ValidateGoModPath(t.Context(), dir, test.path)
			if test.errHas != "" {
				require.ErrorContains(t, err, test.errHas)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
