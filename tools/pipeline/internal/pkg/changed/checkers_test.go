// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package changed

import (
	"context"
	"testing"

	"github.com/google/go-github/v68/github"
	"github.com/stretchr/testify/require"
)

func TestFileGroupDefaultCheckers(t *testing.T) {
	t.Parallel()

	for filename, groups := range map[string]FileGroups{
		".build/entrypoint.sh":                      {FileGroupPipeline},
		".github/actions/changed-files/actions.yml": {FileGroupPipeline},
		".github/workflows/build.yml":               {FileGroupPipeline},
		".github/workflows/build-artifacts-ce.yml":  {FileGroupCommunity, FileGroupPipeline},
		".github/workflows/build-artifacts-ent.yml": {FileGroupEnterprise, FileGroupPipeline},
		".github/workflows/backport-ce-ent.yml":     {FileGroupCommunity, FileGroupPipeline},
		".go-version":                               {FileGroupGoToolchain},
		"audit/backend_ce.go":                       {FileGroupGoApp, FileGroupCommunity},
		"audit/backend_config_ent.go":               {FileGroupGoApp, FileGroupEnterprise},
		"builtin/logical/transit/something_ent.go":  {FileGroupGoApp, FileGroupEnterprise},
		"buf.yml":                                           {FileGroupProto},
		"changelog/1726.txt":                                {FileGroupChangelog},
		"changelog/_1726.txt":                               {FileGroupChangelog, FileGroupEnterprise},
		"command/server/config.go":                          {FileGroupGoApp},
		"command/operator_raft_autopilot_state.go":          {FileGroupGoApp, FileGroupAutopilot},
		"command/agent_ent_test.go":                         {FileGroupGoApp, FileGroupEnterprise},
		"enos/enos-samples-ce-build.hcl":                    {FileGroupCommunity, FileGroupEnos},
		"enos/enos-samples-ent-build.hcl":                   {FileGroupEnos, FileGroupEnterprise},
		"enos/enos-scenario-smoke.hcl":                      {FileGroupEnos},
		"enos/enos-scenario-autopilot-ent.hcl":              {FileGroupEnos, FileGroupEnterprise},
		"go.mod":                                            {FileGroupGoApp, FileGroupGoToolchain},
		"go.sum":                                            {FileGroupGoApp, FileGroupGoToolchain},
		"helper/identity/mfa/types.proto":                   {FileGroupProto},
		"http/util_stubs_oss.go":                            {FileGroupGoApp, FileGroupCommunity},
		"physical/raft/raft_autopilot.go":                   {FileGroupGoApp, FileGroupAutopilot},
		"physical/raft/types.proto":                         {FileGroupProto},
		"scripts/ci-helper.sh":                              {FileGroupPipeline},
		"scripts/cross/Dockerfile-ent":                      {FileGroupEnterprise, FileGroupPipeline},
		"scripts/cross/Dockerfile-ent-hsm":                  {FileGroupEnterprise, FileGroupPipeline},
		"scripts/dev/hsm/README.md":                         {FileGroupEnterprise, FileGroupPipeline},
		"scripts/dist-ent.sh":                               {FileGroupEnterprise, FileGroupPipeline},
		"scripts/testing/test-vault-license.sh":             {FileGroupEnterprise, FileGroupPipeline},
		"scripts/testing/upgrade/README.md":                 {FileGroupEnterprise, FileGroupPipeline},
		"sdk/database/dbplugin/v5/proto/database_ent.pb.go": {FileGroupGoApp, FileGroupEnterprise},
		"sdk/database/dbplugin/v5/proto/database_ent.proto": {FileGroupEnterprise, FileGroupProto},
		"specs/merkle-tree/spec.md":                         {FileGroupEnterprise},
		"tools/pipeline/main.go":                            {FileGroupPipeline},
		"ui/lib/ldap/index.js":                              {FileGroupWebUI},
		"vault/acl.go":                                      {FileGroupGoApp},
		"vault/activity_log_util_ent.go":                    {FileGroupGoApp, FileGroupEnterprise},
		"vault/identity_store_ent_test.go":                  {FileGroupGoApp, FileGroupEnterprise},
		"vault_ent/go.mod":                                  {FileGroupGoApp, FileGroupEnterprise, FileGroupGoToolchain},
		"vault_ent/go.sum":                                  {FileGroupGoApp, FileGroupEnterprise, FileGroupGoToolchain},
		"vault_ent/requires_ent.go":                         {FileGroupGoApp, FileGroupEnterprise},
		"website/content/api-docs/index.mdx":                {FileGroupDocs},
		"CHANGELOG.md":                                      {FileGroupChangelog},
		"CODEOWNERS":                                        {FileGroupPipeline},
		"Dockerfile":                                        {FileGroupPipeline},
		"Makefile":                                          {FileGroupPipeline},
		"README.md":                                         {FileGroupDocs},
	} {
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			file := &File{File: &github.CommitFile{Filename: &filename}}
			Group(context.Background(), file, DefaultFileGroupCheckers...)
			require.Equal(t, groups, file.Groups)
		})
	}
}
