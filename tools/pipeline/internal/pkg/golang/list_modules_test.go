// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestListModulesReq_Run verifies that the JSON output exposes each module's
// go.mod path, which workflows use instead of hardcoded lists.
func TestListModulesReq_Run(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"go.work":    "go 1.21\n\nuse (\n\t.\n\t./sdk\n)\n",
		"go.mod":     "module example.com/root\n\ngo 1.21\n",
		"sdk/go.mod": "module example.com/root/sdk\n\ngo 1.21\n",
	})

	res, err := (&ListModulesReq{GoWork: filepath.Join(root, "go.work")}).Run(t.Context())
	require.NoError(t, err)

	b, err := res.ToJSON()
	require.NoError(t, err)
	var decoded struct {
		Modules []struct {
			GoMod string `json:"go_mod"`
		} `json:"modules"`
	}
	require.NoError(t, json.Unmarshal(b, &decoded))
	require.Len(t, decoded.Modules, 2)
	require.Equal(t, "go.mod", decoded.Modules[0].GoMod)
	require.Equal(t, "sdk/go.mod", decoded.Modules[1].GoMod)

	require.Contains(t, res.ToTable(), "example.com/root/sdk")
	require.Contains(t, res.ToMarkdown(), "| example.com/root/sdk | sdk |")
}
