// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/stretchr/testify/require"
)

// testGoWorkspace writes a workspace with a root module and an api module,
// each with one package, and returns the go.work path.
func testGoWorkspace(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	for name, content := range map[string]string{
		"go.work":       "go 1.21\n\nuse (\n\t.\n\t./api\n)\n",
		"go.mod":        "module example.com/vault\n\ngo 1.21\n",
		"vault/core.go": "package vault\n",
		"api/go.mod":    "module example.com/vault/api\n\ngo 1.21\n",
		"api/client.go": "package api\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}

	return filepath.Join(root, "go.work")
}

// TestGoGroupPackagesCmd_Stdin verifies that '-' reads whitespace-separated
// import paths from stdin, which is how the test matrix pipes partitions in.
func TestGoGroupPackagesCmd_Stdin(t *testing.T) {
	goGroupPackagesReq = &golang.GroupPackagesReq{}
	rootCfg.format = "json"
	defer func() { rootCfg.format = "" }()

	var out bytes.Buffer
	cmd := newGoGroupPackagesCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetIn(strings.NewReader("example.com/vault/api\nexample.com/vault/vault  example.com/vault/api\n"))
	cmd.SetArgs([]string{"--go-work", testGoWorkspace(t), "-"})
	require.NoError(t, cmd.Execute(), out.String())

	require.JSONEq(t, `{"modules":[
		{"path":"example.com/vault","dir":".","packages":["example.com/vault/vault"]},
		{"path":"example.com/vault/api","dir":"api","packages":["example.com/vault/api"]}
	]}`, out.String())
}

// TestGoGroupPackagesCmd_Args verifies that arguments can mix paths relative
// to the current directory, import paths and '-', so a file in the tree can be
// passed directly rather than being read as a list of packages.
func TestGoGroupPackagesCmd_Args(t *testing.T) {
	goGroupPackagesReq = &golang.GroupPackagesReq{}
	rootCfg.format = "json"
	defer func() { rootCfg.format = "" }()

	goWork := testGoWorkspace(t)
	t.Chdir(filepath.Dir(goWork))

	var out bytes.Buffer
	cmd := newGoGroupPackagesCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetIn(strings.NewReader("example.com/vault/api\n"))
	cmd.SetArgs([]string{"api/client.go", "-", "vault/core.go"})
	require.NoError(t, cmd.Execute(), out.String())

	require.JSONEq(t, `{"modules":[
		{"path":"example.com/vault","dir":".","packages":["example.com/vault/vault"]},
		{"path":"example.com/vault/api","dir":"api","packages":["example.com/vault/api"]}
	]}`, out.String())
}

// TestGoGroupPackagesCmd_NoArgs verifies that the command requires input
// instead of silently waiting on stdin.
func TestGoGroupPackagesCmd_NoArgs(t *testing.T) {
	goGroupPackagesReq = &golang.GroupPackagesReq{}

	cmd := newGoGroupPackagesCmd()
	suppressOutput(t, cmd)
	cmd.SetArgs([]string{})
	require.ErrorContains(t, cmd.Execute(), "requires at least 1 arg")
}
