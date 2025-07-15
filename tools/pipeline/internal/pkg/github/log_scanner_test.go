// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogScannerScan_groups(t *testing.T) {
	for name, test := range map[string]struct {
		onlySteps []string
		expected  []*LogEntry
	}{
		"all": {
			expected: []*LogEntry{
				{
					StepName: "Operating System",
					SetupLog: []byte(`2024-12-10T00:12:39.4701675Z Ubuntu
2024-12-10T00:12:39.4702661Z 24.04.1
2024-12-10T00:12:39.4784102Z LTS`),
				},

				{
					StepName: "Runner Image",
					SetupLog: []byte(`2024-12-10T00:12:39.4788307Z Image: ubuntu-24.04
2024-12-10T00:12:39.4789692Z Version: 20241201.1.0
2024-12-10T00:12:39.4792093Z Included Software: https://github.com/actions/runner-images/blob/ubuntu24/20241201.1/images/ubuntu/Ubuntu2404-Readme.md
2024-12-10T00:12:39.4793925Z Image Release: https://github.com/actions/runner-images/releases/tag/ubuntu24%2F20241201.1`),
				},
				{
					StepName: "Runner Image Provisioner",
					SetupLog: []byte(`2024-12-10T00:12:39.4796788Z 2.0.385.1`),
				},
				{
					StepName: "GITHUB_TOKEN Permissions",
					SetupLog: []byte(`2024-12-10T00:12:39.4801616Z Contents: read
2024-12-10T00:12:39.4802428Z Metadata: read`),
					BodyLog: []byte(`2024-12-10T00:12:39.4806392Z Secret source: Actions
2024-12-10T00:12:39.4807540Z Prepare workflow directory
2024-12-10T00:12:39.5129473Z Prepare all required actions
2024-12-10T00:12:39.5167065Z Getting action download info
2024-12-10T00:12:39.9037777Z Download action repository 'actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683' (SHA:11bd71901bbe5b1630ceea73d27597364c9af683)
2024-12-10T00:12:40.0290605Z Download action repository 'hashicorp/vault-action@d1720f055e0635fd932a1d2a48f87a666a57906c' (SHA:d1720f055e0635fd932a1d2a48f87a666a57906c)
2024-12-10T00:12:40.3303861Z Download action repository 'hashicorp/setup-terraform@v3' (SHA:b9cd54a3c349d3f38e8881555d616ced269862dd)
2024-12-10T00:12:40.6598264Z Download action repository 'aws-actions/configure-aws-credentials@e3dd6a429d7300a6a4c196c26e071d42e0343502' (SHA:e3dd6a429d7300a6a4c196c26e071d42e0343502)
2024-12-10T00:12:40.9874063Z Download action repository 'hashicorp/action-setup-enos@v1' (SHA:b9fa53484a1e8fdcc7b02a118bcf01d65b9414c9)
2024-12-10T00:12:41.3078672Z Download action repository 'actions/download-artifact@fa0a91b85d4f404e444e00e005971372dc801d16' (SHA:fa0a91b85d4f404e444e00e005971372dc801d16)
2024-12-10T00:12:41.6838043Z Download action repository 'actions/upload-artifact@b4b15b8c7c6ac21ea08fcf65892d2ee8f75cf882' (SHA:b4b15b8c7c6ac21ea08fcf65892d2ee8f75cf882)
2024-12-10T00:12:41.7759930Z Download action repository 'hashicorp/actions-slack-status@v2.0.1' (SHA:1a3f63b30bd476aee1f3bd6f9d8f2aacc4f14d81)
2024-12-10T00:12:42.1628380Z Getting action download info
2024-12-10T00:12:42.3588968Z Download action repository 'actions/github-script@60a0d83039c74a4aee543508d2ffcb1c3799cdea' (SHA:60a0d83039c74a4aee543508d2ffcb1c3799cdea)
2024-12-10T00:12:42.6849336Z Download action repository 'slackapi/slack-github-action@70cd7be8e40a46e8b0eced40b0de447bdb42f68e' (SHA:70cd7be8e40a46e8b0eced40b0de447bdb42f68e)
2024-12-10T00:12:43.0076240Z Getting action download info
2024-12-10T00:12:43.1908900Z Getting action download info
2024-12-10T00:12:43.4507571Z Getting action download info
2024-12-10T00:12:43.6489981Z Uses: hashicorp/vault/.github/workflows/test-run-enos-scenario-matrix.yml@refs/heads/main (59489a88821ff0431b1ec1b22208a92ecccce183)`),
				},
				{
					StepName: "Inputs",
					SetupLog: []byte(`2024-12-10T00:12:43.6492960Z   build-artifact-name: vault_1.18.3-1_amd64.deb
2024-12-10T00:12:43.6493938Z   sample-max: 2
2024-12-10T00:12:43.6494344Z   sample-name: release_ce_linux_amd64_deb
2024-12-10T00:12:43.6494796Z   runs-on: "ubuntu-latest"
2024-12-10T00:12:43.6495195Z   ssh-key-name: vault-ci-ssh-key
2024-12-10T00:12:43.6495583Z   vault-edition: ce
2024-12-10T00:12:43.6495989Z   vault-revision: 2767f8ee6214d03498b32d776173e1f336281bc5
2024-12-10T00:12:43.6496467Z   vault-version: 1.18.3`),
					BodyLog: []byte(`2024-12-10T00:12:43.6498105Z Complete job name: Test vault_1.18.3-1_amd64.deb / run proxy arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir`),
				},
				{
					StepName: "Run actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683",
					SetupLog: []byte(`2024-12-10T00:12:43.7201109Z with:
2024-12-10T00:12:43.7201732Z   ref: 2767f8ee6214d03498b32d776173e1f336281bc5
2024-12-10T00:12:43.7202218Z   repository: hashicorp/vault
2024-12-10T00:12:43.7202777Z   token: ***
2024-12-10T00:12:43.7203128Z   ssh-strict: true
2024-12-10T00:12:43.7203468Z   ssh-user: git
2024-12-10T00:12:43.7203808Z   persist-credentials: true
2024-12-10T00:12:43.7204183Z   clean: true
2024-12-10T00:12:43.7204535Z   sparse-checkout-cone-mode: true
2024-12-10T00:12:43.7204940Z   fetch-depth: 1
2024-12-10T00:12:43.7205270Z   fetch-tags: false
2024-12-10T00:12:43.7205608Z   show-progress: true
2024-12-10T00:12:43.7205950Z   lfs: false
2024-12-10T00:12:43.7206267Z   submodules: false
2024-12-10T00:12:43.7206604Z   set-safe-directory: true`),
					BodyLog: []byte(`2024-12-10T00:12:43.9369921Z Syncing repository: hashicorp/vault`),
				},
				{
					StepName: "Getting Git version info",
					SetupLog: []byte(`2024-12-10T00:12:43.9372499Z Working directory is '/home/runner/work/vault/vault'
2024-12-10T00:12:43.9373351Z [command]/usr/bin/git version
2024-12-10T00:12:43.9487550Z git version 2.47.1`),
					BodyLog: []byte(`2024-12-10T00:12:43.9533901Z Temporarily overriding HOME='/home/runner/work/_temp/115b40a3-18c7-458a-9579-a62a03b06fec' before making global git config changes
2024-12-10T00:12:43.9535633Z Adding repository directory to the temporary git global config as a safe directory
2024-12-10T00:12:43.9546389Z [command]/usr/bin/git config --global --add safe.directory /home/runner/work/vault/vault
2024-12-10T00:12:43.9583806Z Deleting the contents of '/home/runner/work/vault/vault'`),
				},
				{
					StepName: "Initializing the repository",
					SetupLog: []byte(`2024-12-10T00:12:43.9592246Z [command]/usr/bin/git init /home/runner/work/vault/vault
2024-12-10T00:12:43.9681668Z hint: Using 'master' as the name for the initial branch. This default branch name
2024-12-10T00:12:43.9682921Z hint: is subject to change. To configure the initial branch name to use in all
2024-12-10T00:12:43.9684045Z hint: of your new repositories, which will suppress this warning, call:
2024-12-10T00:12:43.9684886Z hint:
2024-12-10T00:12:43.9685567Z hint: 	git config --global init.defaultBranch <name>
2024-12-10T00:12:43.9686287Z hint:
2024-12-10T00:12:43.9686783Z hint: Names commonly chosen instead of 'master' are 'main', 'trunk' and
2024-12-10T00:12:43.9687500Z hint: 'development'. The just-created branch can be renamed via this command:
2024-12-10T00:12:43.9688022Z hint:
2024-12-10T00:12:43.9688363Z hint: 	git branch -m <name>
2024-12-10T00:12:43.9692962Z Initialized empty Git repository in /home/runner/work/vault/vault/.git/
2024-12-10T00:12:43.9704240Z [command]/usr/bin/git remote add origin https://github.com/hashicorp/vault`),
				},
				{
					StepName: "Disabling automatic garbage collection",
					SetupLog: []byte(`2024-12-10T00:12:43.9750441Z [command]/usr/bin/git config --local gc.auto 0`),
				},
				{
					StepName: "Setting up auth",
					SetupLog: []byte(`2024-12-10T00:12:43.9786405Z [command]/usr/bin/git config --local --name-only --get-regexp core\.sshCommand
2024-12-10T00:12:43.9816121Z [command]/usr/bin/git submodule foreach --recursive sh -c "git config --local --name-only --get-regexp 'core\.sshCommand' && git config --local --unset-all 'core.sshCommand' || :"
2024-12-10T00:12:44.0274103Z [command]/usr/bin/git config --local --name-only --get-regexp http\.https\:\/\/github\.com\/\.extraheader
2024-12-10T00:12:44.0302372Z [command]/usr/bin/git submodule foreach --recursive sh -c "git config --local --name-only --get-regexp 'http\.https\:\/\/github\.com\/\.extraheader' && git config --local --unset-all 'http.https://github.com/.extraheader' || :"
2024-12-10T00:12:44.0522906Z [command]/usr/bin/git config --local http.https://github.com/.extraheader AUTHORIZATION: basic ***`),
				},
				{
					StepName: "Fetching the repository",
					SetupLog: []byte(`2024-12-10T00:12:44.0572949Z [command]/usr/bin/git -c protocol.version=2 fetch --no-tags --prune --no-recurse-submodules --depth=1 origin 2767f8ee6214d03498b32d776173e1f336281bc5
2024-12-10T00:12:45.5127876Z From https://github.com/hashicorp/vault
2024-12-10T00:12:45.5130722Z  * branch            2767f8ee6214d03498b32d776173e1f336281bc5 -> FETCH_HEAD`),
				},
				{
					StepName: "Determining the checkout info",
					BodyLog: []byte(`2024-12-10T00:12:45.5169128Z [command]/usr/bin/git sparse-checkout disable
2024-12-10T00:12:45.5214350Z [command]/usr/bin/git config --local --unset-all extensions.worktreeConfig`),
				},
				{
					StepName: "Checking out the ref",
					SetupLog: []byte(`2024-12-10T00:12:45.5246868Z [command]/usr/bin/git checkout --progress --force 2767f8ee6214d03498b32d776173e1f336281bc5
2024-12-10T00:12:46.1483364Z Note: switching to '2767f8ee6214d03498b32d776173e1f336281bc5'.
2024-12-10T00:12:46.1484000Z
2024-12-10T00:12:46.1484523Z You are in 'detached HEAD' state. You can look around, make experimental
2024-12-10T00:12:46.1485553Z changes and commit them, and you can discard any commits you make in this
2024-12-10T00:12:46.1486616Z state without impacting any branches by switching back to a branch.
2024-12-10T00:12:46.1487236Z
2024-12-10T00:12:46.1487698Z If you want to create a new branch to retain commits you create, you may
2024-12-10T00:12:46.1488691Z do so (now or later) by using -c with the switch command. Example:
2024-12-10T00:12:46.1489249Z
2024-12-10T00:12:46.1489526Z   git switch -c <new-branch-name>
2024-12-10T00:12:46.1489940Z
2024-12-10T00:12:46.1490209Z Or undo this operation with:
2024-12-10T00:12:46.1490632Z
2024-12-10T00:12:46.1490894Z   git switch -
2024-12-10T00:12:46.1491230Z
2024-12-10T00:12:46.1492016Z Turn off this advice by setting config variable advice.detachedHead to false
2024-12-10T00:12:46.1492672Z
2024-12-10T00:12:46.1493261Z HEAD is now at 2767f8e Update vault-plugin-secrets-openldap to v0.14.4 (#29131) (#29133)`),
					BodyLog: []byte(`2024-12-10T00:12:46.1556639Z [command]/usr/bin/git log -1 --format=%H
2024-12-10T00:12:46.1578665Z 2767f8ee6214d03498b32d776173e1f336281bc5`),
				},
				{
					StepName: "Run enos scenario launch --timeout 45m0s --chdir ./enos proxy arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir",
					SetupLog: []byte(`2024-12-10T00:12:50.2289674Z [36;1menos scenario launch --timeout 45m0s --chdir ./enos proxy arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir[0m
2024-12-10T00:12:50.2314107Z shell: /usr/bin/bash -e {0}
2024-12-10T00:12:50.2314496Z env:
2024-12-10T00:12:50.2314951Z   GITHUB_TOKEN: ***
2024-12-10T00:12:50.2315358Z   ENOS_DEBUG_DATA_ROOT_DIR: /tmp/enos-debug-data
2024-12-10T00:12:50.2331925Z   DYNAMIC_CONFIG_KEY: 1.18.3-2024-12-49
2024-12-10T00:12:50.2332387Z   DYNAMIC_CONFIG_PATH: enos/enos-dynamic-config.hcl`),
					BodyLog: []byte(`2024-12-10T00:12:50.7389391Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] generate: running
2024-12-10T00:12:51.2323596Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] generate: success!
2024-12-10T00:12:51.2325463Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] init: running
2024-12-10T00:12:56.2324410Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] init: running
2024-12-10T00:12:57.8598120Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] init: success!
2024-12-10T00:12:57.8602036Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] validate: running
2024-12-10T00:13:00.2734861Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] validate: success!
2024-12-10T00:13:00.2737223Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] plan: running
2024-12-10T00:13:04.6929951Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] plan: success!
2024-12-10T00:13:04.6933016Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] apply: running
2024-12-10T00:13:09.6937929Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] apply: running
2024-12-10T00:16:27.4736392Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] apply: failed!
2024-12-10T00:16:27.4738430Z [31m╷[0m
2024-12-10T00:16:27.4739095Z [31m│[0m [1m[31mError: [0m[1mexit status 1
2024-12-10T00:16:27.4739835Z [31m│[0m
2024-12-10T00:16:27.4740433Z [31m│[0m Error: Execution Error
2024-12-10T00:16:27.4741075Z [31m│[0m
2024-12-10T00:16:27.4742287Z [31m│[0m   with module.verify_secrets_engines_create.enos_remote_exec.kv_put_secret_test["0"],
2024-12-10T00:16:27.4744486Z [31m│[0m   on ../../modules/verify_secrets_engines/modules/create/kv.tf line 106, in resource "enos_remote_exec" "kv_put_secret_test":
2024-12-10T00:16:27.4746078Z [31m│[0m  106: resource "enos_remote_exec" "kv_put_secret_test" {
2024-12-10T00:16:27.4746974Z [31m│[0m
2024-12-10T00:16:27.4747719Z [31m│[0m failed to execute commands due to: running script:
2024-12-10T00:16:27.4749074Z [31m│[0m [/home/runner/work/vault/vault/enos/modules/verify_secrets_engines/scripts/kv-put.sh]
2024-12-10T00:16:27.4750274Z [31m│[0m failed, due to: 1 error occurred:
2024-12-10T00:16:27.4751566Z [31m│[0m 	* executing script: Process exited with status 2: Error writing data to
2024-12-10T00:16:27.4752733Z [31m│[0m secret/data/smoke-0: Error making API request.
2024-12-10T00:16:27.4753591Z [31m│[0m
2024-12-10T00:16:27.4754454Z [31m│[0m URL: PUT http://127.0.0.1:8200/v1/secret/data/smoke-0
2024-12-10T00:16:27.4755467Z [31m│[0m Code: 403. Errors:
2024-12-10T00:16:27.4756215Z [31m│[0m
2024-12-10T00:16:27.4756934Z [31m│[0m * 1 error occurred:
2024-12-10T00:16:27.4757726Z [31m│[0m 	* permission denied
2024-12-10T00:16:27.4951718Z Scenario: proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] failed!
2024-12-10T00:16:27.4952646Z   Init: success!
2024-12-10T00:16:27.4952980Z   Validate: success!
2024-12-10T00:16:27.4953315Z   Plan: success!`),
					ErrorLog: []byte(`2024-12-10T00:16:27.4962854Z Process completed with exit code 1.`),
				},
				{
					StepName: "Run enos scenario launch --timeout 45m0s --chdir ./enos proxy arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir",
					SetupLog: []byte(`2024-12-10T00:16:27.5028497Z [36;1menos scenario launch --timeout 45m0s --chdir ./enos proxy arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir[0m
2024-12-10T00:16:27.5056211Z shell: /usr/bin/bash -e {0}
2024-12-10T00:16:27.5056606Z env:
2024-12-10T00:16:27.5057108Z   GITHUB_TOKEN: ***
2024-12-10T00:16:27.5074106Z   DYNAMIC_CONFIG_KEY: 1.18.3-2024-12-49
2024-12-10T00:16:27.5074569Z   DYNAMIC_CONFIG_PATH: enos/enos-dynamic-config.hcl`),
					BodyLog: []byte(`2024-12-10T00:16:28.0207793Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] generate: running
2024-12-10T00:16:28.5089608Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] generate: success!
2024-12-10T00:16:28.5093852Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] init: running
2024-12-10T00:16:29.8095857Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] init: success!
2024-12-10T00:16:29.8097587Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] validate: running
2024-12-10T00:16:32.1198366Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] validate: success!
2024-12-10T00:16:32.1199975Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] plan: running
2024-12-10T00:16:37.0846706Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] plan: success!
2024-12-10T00:16:37.0848336Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] apply: running
2024-12-10T00:16:42.0856135Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] apply: running
2024-12-10T00:16:46.3363724Z proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] apply: failed!
2024-12-10T00:16:46.3365518Z [31m╷[0m
2024-12-10T00:16:46.3366178Z [31m│[0m [1m[31mError: [0m[1mexit status 1
2024-12-10T00:16:46.3366866Z [31m│[0m
2024-12-10T00:16:46.3367513Z [31m│[0m Error: Execution Error
2024-12-10T00:16:46.3368230Z [31m│[0m
2024-12-10T00:16:46.3369253Z [31m│[0m   with module.verify_secrets_engines_create.enos_remote_exec.kv_put_secret_test["1"],
2024-12-10T00:16:46.3370933Z [31m│[0m   on ../../modules/verify_secrets_engines/modules/create/kv.tf line 106, in resource "enos_remote_exec" "kv_put_secret_test":
2024-12-10T00:16:46.3372804Z [31m│[0m  106: resource "enos_remote_exec" "kv_put_secret_test" {
2024-12-10T00:16:46.3373675Z [31m│[0m
2024-12-10T00:16:46.3374416Z [31m│[0m failed to execute commands due to: running script:
2024-12-10T00:16:46.3375692Z [31m│[0m [/home/runner/work/vault/vault/enos/modules/verify_secrets_engines/scripts/kv-put.sh]
2024-12-10T00:16:46.3376812Z [31m│[0m failed, due to: 1 error occurred:
2024-12-10T00:16:46.3377846Z [31m│[0m 	* executing script: Process exited with status 2: Error writing data to
2024-12-10T00:16:46.3379243Z [31m│[0m secret/data/smoke-1: Error making API request.
2024-12-10T00:16:46.3380011Z [31m│[0m
2024-12-10T00:16:46.3380747Z [31m│[0m URL: PUT http://127.0.0.1:8200/v1/secret/data/smoke-1
2024-12-10T00:16:46.3381808Z [31m│[0m Code: 403. Errors:
2024-12-10T00:16:46.3382441Z [31m│[0m
2024-12-10T00:16:46.3383036Z [31m│[0m * 1 error occurred:
2024-12-10T00:16:46.3383733Z [31m│[0m 	* permission denied
2024-12-10T00:16:46.3384376Z [31m│[0m
2024-12-10T00:16:46.3384902Z [31m│[0m
2024-12-10T00:16:46.3385425Z [31m│[0m
2024-12-10T00:16:46.3385951Z [31m│[0m
2024-12-10T00:16:46.3386514Z [31m│[0m
2024-12-10T00:16:46.3387094Z [31m│[0m output:
2024-12-10T00:16:46.3387970Z [31m│[0m Error writing data to secret/data/smoke-1: Error making API request.
2024-12-10T00:16:46.3388879Z [31m│[0m
2024-12-10T00:16:46.3389512Z [31m│[0m URL: PUT http://127.0.0.1:8200/v1/secret/data/smoke-1
2024-12-10T00:16:46.3390043Z [31m│[0m Code: 403. Errors:
2024-12-10T00:16:46.3390423Z [31m│[0m
2024-12-10T00:16:46.3390775Z [31m│[0m * 1 error occurred:
2024-12-10T00:16:46.3391187Z [31m│[0m 	* permission denied
2024-12-10T00:16:46.3391955Z [31m│[0m
2024-12-10T00:16:46.3392544Z [31m│[0m
2024-12-10T00:16:46.3393117Z [31m│[0m
2024-12-10T00:16:46.3393728Z [31m│[0m SSH Transport Config:
2024-12-10T00:16:46.3394448Z [31m│[0m             user : ubuntu
2024-12-10T00:16:46.3395054Z [31m│[0m             host : 23.22.214.188
2024-12-10T00:16:46.3395515Z [31m│[0m      private_key : null
2024-12-10T00:16:46.3396144Z [31m│[0m private_key_path : /home/runner/work/vault/vault/enos/support/private_key.pem
2024-12-10T00:16:46.3396966Z [31m│[0m       passphrase : null
2024-12-10T00:16:46.3397419Z [31m│[0m  passphrase_path : null
2024-12-10T00:16:46.3397813Z [31m│[0m
2024-12-10T00:16:46.3398141Z [31m│[0m
2024-12-10T00:16:46.3398495Z [31m│[0m Application Logs:
2024-12-10T00:16:46.3398993Z [31m│[0m   vault: /tmp/enos-debug-data/vault_23.22.214.188.log
2024-12-10T00:16:46.3399481Z [31m│[0m
2024-12-10T00:16:46.3571781Z Scenario: proxy [arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir] failed!
2024-12-10T00:16:46.3572720Z   Init: success!
2024-12-10T00:16:46.3573082Z   Validate: success!
2024-12-10T00:16:46.3573432Z   Plan: success!`),
					ErrorLog: []byte(`2024-12-10T00:16:46.3575376Z Process completed with exit code 1.`),
				},
				{
					StepName: "Run hashicorp/actions-slack-status@v2.0.1",
					SetupLog: []byte(`2024-12-10T00:18:21.7000653Z with:
2024-12-10T00:18:21.7002008Z   failure-message: enos scenario launch proxy arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir failed.
Triggering event: "repository_dispatch"
Actor: "crt-orchestrator[bot]"
2024-12-10T00:18:21.7003270Z   status: failure
2024-12-10T00:18:21.7003900Z   slack-webhook-url: ***
2024-12-10T00:18:21.7004273Z env:
2024-12-10T00:18:21.7004715Z   GITHUB_TOKEN: ***
2024-12-10T00:18:21.7021520Z   DYNAMIC_CONFIG_KEY: 1.18.3-2024-12-49
2024-12-10T00:18:21.7021989Z   DYNAMIC_CONFIG_PATH: enos/enos-dynamic-config.hcl`),
				},
				{
					StepName: "Run actions/github-script@60a0d83039c74a4aee543508d2ffcb1c3799cdea",
					SetupLog: []byte(`2024-12-10T00:18:21.7064074Z with:
2024-12-10T00:18:21.7064579Z   script: return require(process.env.GITHUB_ACTION_PATH + '/main.js')({context, core})

2024-12-10T00:18:21.7065286Z   github-token: ***
2024-12-10T00:18:21.7065816Z   debug: false
2024-12-10T00:18:21.7066175Z   user-agent: actions/github-script
2024-12-10T00:18:21.7066603Z   result-encoding: json
2024-12-10T00:18:21.7066955Z   retries: 0
2024-12-10T00:18:21.7067326Z   retry-exempt-status-codes: 400,401,403,404,422
2024-12-10T00:18:21.7067754Z env:
2024-12-10T00:18:21.7068163Z   GITHUB_TOKEN: ***
2024-12-10T00:18:21.7085096Z   DYNAMIC_CONFIG_KEY: 1.18.3-2024-12-49
2024-12-10T00:18:21.7085559Z   DYNAMIC_CONFIG_PATH: enos/enos-dynamic-config.hcl
2024-12-10T00:18:21.7086012Z   INPUT_STATUS: failure
2024-12-10T00:18:21.7086382Z   INPUT_SUCCESS-MESSAGE:
2024-12-10T00:18:21.7087631Z   INPUT_FAILURE-MESSAGE: enos scenario launch proxy arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir failed.
Triggering event: "repository_dispatch"
Actor: "crt-orchestrator[bot]"
2024-12-10T00:18:21.7088915Z   INPUT_CANCELLED-MESSAGE:
2024-12-10T00:18:21.7089296Z   INPUT_SKIPPED-MESSAGE:`),
					BodyLog: []byte(`2024-12-10T00:18:22.6819662Z Post job cleanup.
2024-12-10T00:18:22.6892801Z Post job cleanup.
2024-12-10T00:18:22.8076185Z Post job cleanup.
2024-12-10T00:18:22.9022071Z [command]/usr/bin/git version
2024-12-10T00:18:22.9062573Z git version 2.47.1
2024-12-10T00:18:22.9107683Z Temporarily overriding HOME='/home/runner/work/_temp/dbe6a6bb-4627-4ebf-86bc-d7e567dc0df4' before making global git config changes
2024-12-10T00:18:22.9109480Z Adding repository directory to the temporary git global config as a safe directory
2024-12-10T00:18:22.9122431Z [command]/usr/bin/git config --global --add safe.directory /home/runner/work/vault/vault
2024-12-10T00:18:22.9159445Z [command]/usr/bin/git config --local --name-only --get-regexp core\.sshCommand
2024-12-10T00:18:22.9193356Z [command]/usr/bin/git submodule foreach --recursive sh -c "git config --local --name-only --get-regexp 'core\.sshCommand' && git config --local --unset-all 'core.sshCommand' || :"
2024-12-10T00:18:22.9439897Z [command]/usr/bin/git config --local --name-only --get-regexp http\.https\:\/\/github\.com\/\.extraheader
2024-12-10T00:18:22.9461015Z http.https://github.com/.extraheader
2024-12-10T00:18:22.9474504Z [command]/usr/bin/git config --local --unset-all http.https://github.com/.extraheader
2024-12-10T00:18:22.9505151Z [command]/usr/bin/git submodule foreach --recursive sh -c "git config --local --name-only --get-regexp 'http\.https\:\/\/github\.com\/\.extraheader' && git config --local --unset-all 'http.https://github.com/.extraheader' || :"
2024-12-10T00:18:22.9846381Z Cleaning up orphan processes`),
				},
			},
		},
		"only_steps": {
			onlySteps: []string{
				"Operating System",
				"Runner Image Provisioner",
				"Disabling automatic garbage collection",
				"Setting up auth",
				"Run hashicorp/actions-slack-status@v2.0.1",
			},
			expected: []*LogEntry{
				{
					StepName: "Operating System",
					SetupLog: []byte(`2024-12-10T00:12:39.4701675Z Ubuntu
2024-12-10T00:12:39.4702661Z 24.04.1
2024-12-10T00:12:39.4784102Z LTS`),
				},
				{
					StepName: "Runner Image Provisioner",
					SetupLog: []byte(`2024-12-10T00:12:39.4796788Z 2.0.385.1`),
				},
				{
					StepName: "Disabling automatic garbage collection",
					SetupLog: []byte(`2024-12-10T00:12:43.9750441Z [command]/usr/bin/git config --local gc.auto 0`),
				},
				{
					StepName: "Setting up auth",
					SetupLog: []byte(`2024-12-10T00:12:43.9786405Z [command]/usr/bin/git config --local --name-only --get-regexp core\.sshCommand
2024-12-10T00:12:43.9816121Z [command]/usr/bin/git submodule foreach --recursive sh -c "git config --local --name-only --get-regexp 'core\.sshCommand' && git config --local --unset-all 'core.sshCommand' || :"
2024-12-10T00:12:44.0274103Z [command]/usr/bin/git config --local --name-only --get-regexp http\.https\:\/\/github\.com\/\.extraheader
2024-12-10T00:12:44.0302372Z [command]/usr/bin/git submodule foreach --recursive sh -c "git config --local --name-only --get-regexp 'http\.https\:\/\/github\.com\/\.extraheader' && git config --local --unset-all 'http.https://github.com/.extraheader' || :"
2024-12-10T00:12:44.0522906Z [command]/usr/bin/git config --local http.https://github.com/.extraheader AUTHORIZATION: basic ***`),
				},
				{
					StepName: "Run hashicorp/actions-slack-status@v2.0.1",
					SetupLog: []byte(`2024-12-10T00:18:21.7000653Z with:
2024-12-10T00:18:21.7002008Z   failure-message: enos scenario launch proxy arch:amd64 artifact_source:artifactory artifact_type:package backend:raft config_mode:file consul_edition:ent consul_version:1.14.11 distro:ubuntu edition:ce ip_version:4 seal:shamir failed.
Triggering event: "repository_dispatch"
Actor: "crt-orchestrator[bot]"
2024-12-10T00:18:21.7003270Z   status: failure
2024-12-10T00:18:21.7003900Z   slack-webhook-url: ***
2024-12-10T00:18:21.7004273Z env:
2024-12-10T00:18:21.7004715Z   GITHUB_TOKEN: ***
2024-12-10T00:18:21.7021520Z   DYNAMIC_CONFIG_KEY: 1.18.3-2024-12-49
2024-12-10T00:18:21.7021989Z   DYNAMIC_CONFIG_PATH: enos/enos-dynamic-config.hcl`),
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			f, err := os.Open(filepath.Join("./testfixtures/actions.log"))
			require.NoError(t, err)
			scanner := NewLogScaner(WithLogScannerOnlySteps(test.onlySteps))
			res, err := scanner.Scan(f)
			require.NoError(t, err)
			require.Len(t, res, len(test.expected))
			require.Equal(t, test.expected, res)
		})
	}
}
