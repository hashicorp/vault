// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package changed

import (
	"context"
	"testing"

	"github.com/google/go-github/v68/github"
	"github.com/stretchr/testify/require"
)

func TestFileGroupCheckerDir(t *testing.T) {
	t.Parallel()

	for filename, groups := range map[string]FileGroups{
		".github/actions/changed-files/actions.yml": {FileGroupPipeline},
		".github/workflows/build.yml":               {FileGroupPipeline},
		"changelog/16455.txt":                       {FileGroupChangelog},
		"enos/Makefile":                             {FileGroupEnos},
		"tools/pipeline/main.go":                    {FileGroupTools},
		"ui/lib/ldap/index.js":                      {FileGroupWebUI},
		"website/content/api-docs/index.mdx":        {FileGroupDocs},
	} {
		t.Run(filename, func(t *testing.T) {
			file := &File{File: &github.CommitFile{Filename: &filename}}
			require.Equal(t, groups, FileGroupCheckerDir(context.Background(), file))
		})
	}
}

func TestFileGroupCheckerFileName(t *testing.T) {
	t.Parallel()

	for filename, groups := range map[string]FileGroups{
		"buf.yml":      {FileGroupProto},
		"CHANGELOG.md": {FileGroupChangelog},
		"CODEOWNERS":   {FileGroupPipeline},
		"go.mod":       {FileGroupGoApp, FileGroupGoModules},
		"go.sum":       {FileGroupGoApp, FileGroupGoModules},
	} {
		t.Run(filename, func(t *testing.T) {
			file := &File{File: &github.CommitFile{Filename: &filename}}
			require.Equal(t, groups, FileGroupCheckerFileName(context.Background(), file))
		})
	}
}

func TestFileGroupCheckerFileGo(t *testing.T) {
	t.Parallel()

	for filename, groups := range map[string]FileGroups{
		"vault/acl.go":                             {FileGroupGoApp},
		"command/server/config.go":                 {FileGroupGoApp},
		"tools/pipeline/main.go":                   {FileGroupGoApp, FileGroupPipeline},
		"command/operator_raft_autopilot_state.go": {FileGroupGoApp, FileGroupAutopilot},
		"physical/raft/raft_autopilot.go":          {FileGroupGoApp, FileGroupAutopilot},
		"http/util_stubs_oss.go":                   {FileGroupGoApp, FileGroupCommunity},
		"audit/backend_ce.go":                      {FileGroupGoApp, FileGroupCommunity},
		"vault/activity_log_util_ent.go":           {FileGroupGoApp, FileGroupEnterprise},
	} {
		t.Run(filename, func(t *testing.T) {
			file := &File{File: &github.CommitFile{Filename: &filename}}
			require.Equal(t, groups, FileGroupCheckerFileGo(context.Background(), file))
		})
	}
}

func TestFileGroupCheckerProto(t *testing.T) {
	t.Parallel()

	for filename, groups := range map[string]FileGroups{
		"physical/raft/types.proto":       {FileGroupProto},
		"helper/identity/mfa/types.proto": {FileGroupProto},
	} {
		t.Run(filename, func(t *testing.T) {
			file := &File{File: &github.CommitFile{Filename: &filename}}
			require.Equal(t, groups, FileGroupCheckerFileProto(context.Background(), file))
		})
	}
}
