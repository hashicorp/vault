// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/stretchr/testify/require"
)

// TestGoCheckTidyCmd_Failure verifies that an untidy module fails the command
// with a summary of what isn't tidy. Fixes are left to the caller, since how
// a repository tidies its modules isn't the tool's business.
func TestGoCheckTidyCmd_Failure(t *testing.T) {
	goCheckTidyReq = &golang.CheckTidyReq{}

	root := t.TempDir()
	for name, content := range map[string]string{
		"go.work":  "go 1.21\n\nuse (\n\t./a\n\t./b\n)\n",
		"a/go.mod": "module example.com/a\n\ngo 1.21\n\nrequire example.com/b v0.0.0\n\nreplace example.com/b => ../b\n",
		"b/go.mod": "module example.com/b\n\ngo 1.21\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}

	cmd := newGoCheckTidyCmd()
	suppressOutput(t, cmd)
	cmd.SetArgs([]string{"--go-work", filepath.Join(root, "go.work")})
	err := cmd.Execute()
	require.EqualError(t, err, "check failed: 1 of 2 modules not tidy, 0 failed")
}
