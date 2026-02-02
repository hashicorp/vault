// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package changed

import (
	"context"
	"testing"

	"github.com/google/go-github/v81/github"
	"github.com/stretchr/testify/require"
)

func TestFileGroupDefaultCheckers(t *testing.T) {
	t.Parallel()

	for filename, groups := range map[string]FileGroups{
		".build/entrypoint.sh":                           {FileGroupPipeline},
		".github/actions/changed-files/actions.yml":      {FileGroupGithub, FileGroupPipeline},
		".github/workflows/build.yml":                    {FileGroupGithub, FileGroupPipeline},
		".github/workflows/build-artifacts-ce.yml":       {FileGroupCommunity, FileGroupGithub, FileGroupPipeline},
		".github/workflows/build-artifacts-ent.yml":      {FileGroupEnterprise, FileGroupGithub, FileGroupPipeline},
		".github/workflows/backport-ce-ent.yml":          {FileGroupCommunity, FileGroupGithub, FileGroupPipeline},
		".github/scripts/pr_comment.sh":                  {FileGroupGithub, FileGroupPipeline},
		".github/CODEOWNERS":                             {FileGroupGithub},
		".go-version":                                    {FileGroupGoToolchain},
		".release/ibm-pao/eboms/5900-BJ8.essentials.csv": {FileGroupEnterprise, FileGroupPipeline},
		".release/docker/ubi-docker-entrypoint.sh":       {FileGroupPipeline},
		"audit/backend_ce.go":                            {FileGroupGoApp, FileGroupCommunity},
		"audit/backend_config_ent.go":                    {FileGroupGoApp, FileGroupEnterprise},
		"builtin/logical/transit/something_ent.go":       {FileGroupGoApp, FileGroupEnterprise},
		"buf.yml":                                                                {FileGroupProto},
		"changelog/1726.txt":                                                     {FileGroupChangelog},
		"changelog/_1726.txt":                                                    {FileGroupChangelog},
		"command/server/config.go":                                               {FileGroupGoApp},
		"command/operator_raft_autopilot_state.go":                               {FileGroupGoApp, FileGroupAutopilot},
		"command/agent_ent_test.go":                                              {FileGroupGoApp, FileGroupEnterprise},
		"enos/enos-samples-ce-build.hcl":                                         {FileGroupCommunity, FileGroupEnos},
		"enos/enos-samples-ent-build.hcl":                                        {FileGroupEnos, FileGroupEnterprise},
		"enos/enos-scenario-smoke.hcl":                                           {FileGroupEnos},
		"enos/enos-scenario-autopilot-ent.hcl":                                   {FileGroupEnos, FileGroupEnterprise},
		"enos/modules/softhsm_create_vault_keys/scripts/create-keys.sh":          {FileGroupEnos},
		"enos/modules/softhsm_create_vault_keys/scripts/get-keys.sh":             {FileGroupEnos},
		"enos/modules/softhsm_distribute_vault_keys/main.tf":                     {FileGroupEnos},
		"enos/modules/softhsm_distribute_vault_keys/scripts/distribute-token.sh": {FileGroupEnos},
		"enos/modules/softhsm_init/main.tf":                                      {FileGroupEnos},
		"enos/modules/softhsm_init/scripts/init-softhsm.sh":                      {FileGroupEnos},
		"enos/modules/softhsm_install/main.tf":                                   {FileGroupEnos},
		"enos/modules/softhsm_install/scripts/find-shared-object.sh":             {FileGroupEnos},
		"enos/modules/verify_secrets_engines/scripts/identity-verify-entity.sh":  {FileGroupEnos},
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
		"scripts/docker/docker-entrypoint.sh":               {FileGroupPipeline},
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
