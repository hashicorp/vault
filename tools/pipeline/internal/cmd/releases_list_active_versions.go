// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
)

var listReleaseActiveVersionsReq = &releases.ListActiveVersionsReq{}

func newReleasesListActiveVersionsCmd() *cobra.Command {
	activeVersionsCmd := &cobra.Command{
		Use:   "active-versions [--versions-config .release/versions.hcl]",
		Short: "List the active versions from .release/versions.hcl",
		Long: `List the active versions from .release/versions.hcl

This command reads the .release/versions.hcl file and outputs information about
active Vault versions, including branch names, version numbers, and metadata for
both Community Edition (CE) and Enterprise Edition (ENT).

OUTPUT MODES:

1. Default (human-readable table):
   Displays versions in a formatted table with columns for version, CE active status,
   LTS status, enterprise branch, and CE branch (if applicable).

2. --format json (machine-readable):
   Outputs a JSON structure containing:
   - versions_config: Full configuration from .release/versions.hcl
   - versions: Array of version objects with metadata and branch names
   Each version object includes: version, ce_active, lts, enterprise_branch, and
   ce_branch (if ce_active is true).

3. --github-output (GitHub Actions integration):
   Writes to $GITHUB_OUTPUT with key "active-versions=" followed by JSON containing:
   - All fields from --format json
   - Additional arrays: versions, ce_active_versions, lts_versions, active_branches,
     ce_active_branches, lts_active_branches, all_active_branches
   - Three matrix arrays optimized for GitHub Actions workflows:
     * active_versions_matrix: ENT-only versions
     * ce_active_versions_matrix: CE-only versions
     * all_active_versions_matrix: Both CE and ENT versions
   Each matrix entry contains: branch, version, edition ("ce" or "ent"), and lts (boolean).

KEY DIFFERENCES between --format json and --github-output:
  --format json: Simpler structure focused on version metadata, outputs to stdout
	--github-output: Comprehensive structure with pre-computed matrix arrays and branch
  lists, writes to $GITHUB_OUTPUT file with "active-versions=" key prefix for direct
  use in GitHub Actions workflows`,
		RunE: runListActiveVersionsReq,
		Args: cobra.NoArgs,
		Example: `  # View as human-readable table (default)
	$ pipeline releases list active-versions
  Output:
  VERSION  CE ACTIVE  LTS    ENTERPRISE BRANCH   CE BRANCH
  1.19.x   false      true   release/1.19.x+ent
  1.20.x   false      false  release/1.20.x+ent
  1.21.x   false      false  release/1.21.x+ent
  2.0.x    true       true   release/2.0.x+ent   release/2.0.x

  # View as JSON (simpler structure, outputs to stdout, machine readable)
  $ pipeline releases list active-versions --format json --include-ce-prefix ce
  Output:
  {
    "versions_config": {
      "schema": 1,
      "active_versions": {
        "versions": {
          "1.19.x": {"ce_active": false, "lts": true},
          "1.20.x": {"ce_active": false, "lts": false},
          "1.21.x": {"ce_active": false, "lts": false},
          "2.0.x": {"ce_active": true, "lts": true}
        }
      }
    },
    "versions": [
      {"version": "1.19.x", "ce_active": false, "lts": true, "enterprise_branch": "release/1.19.x+ent"},
      {"version": "1.20.x", "ce_active": false, "lts": false, "enterprise_branch": "release/1.20.x+ent"},
      {"version": "1.21.x", "ce_active": false, "lts": false, "enterprise_branch": "release/1.21.x+ent"},
      {"version": "2.0.x", "ce_active": true, "lts": true, "enterprise_branch": "release/2.0.x+ent", "ce_branch": "ce/release/2.0.x"}
    ]
  }

  # Write to GitHub Actions output (comprehensive structure with matrix arrays)
  $ pipeline releases list active-versions --github-output --include-main --include-ce-prefix ce
  Writes to $GITHUB_OUTPUT with key prefix "active-versions=":
  active-versions={"versions_config":{...},"versions":["1.19.x","1.20.x","1.21.x","2.0.x"],"ce_active_versions":["2.0.x"],"lts_versions":["1.19.x","2.0.x"],"active_branches":["main","release/1.19.x+ent","release/1.20.x+ent","release/1.21.x+ent","release/2.0.x+ent"],"ce_active_branches":["ce/main","ce/release/2.0.x"],"lts_active_branches":["release/1.19.x+ent","release/2.0.x+ent"],"all_active_branches":["ce/main","ce/release/2.0.x","main","release/1.19.x+ent","release/1.20.x+ent","release/1.21.x+ent","release/2.0.x+ent"],"active_versions_matrix":[{"branch":"main","version":"main","edition":"ent","lts":false},{"branch":"release/1.19.x+ent","version":"1.19.x","edition":"ent","lts":true},{"branch":"release/1.20.x+ent","version":"1.20.x","edition":"ent","lts":false},{"branch":"release/1.21.x+ent","version":"1.21.x","edition":"ent","lts":false},{"branch":"release/2.0.x+ent","version":"2.0.x","edition":"ent","lts":true}],"ce_active_versions_matrix":[{"branch":"ce/main","version":"main","edition":"ce","lts":false},{"branch":"ce/release/2.0.x","version":"2.0.x","edition":"ce","lts":true}],"all_active_versions_matrix":[{"branch":"ce/main","version":"main","edition":"ce","lts":false},{"branch":"ce/release/2.0.x","version":"2.0.x","edition":"ce","lts":true},{"branch":"main","version":"main","edition":"ent","lts":false},{"branch":"release/1.19.x+ent","version":"1.19.x","edition":"ent","lts":true},{"branch":"release/1.20.x+ent","version":"1.20.x","edition":"ent","lts":false},{"branch":"release/1.21.x+ent","version":"1.21.x","edition":"ent","lts":false},{"branch":"release/2.0.x+ent","version":"2.0.x","edition":"ent","lts":true}]}

  GitHub Actions Workflow Examples:

  Example 1: ENT-only builds (5 branches: main + 4 release branches)
  strategy:
    matrix:
      include: ${{ fromJSON(needs.setup.outputs.active-versions).active_versions_matrix }}

  Example 2: CE-only builds (2 branches: ce/main + ce/release/2.0.x)
  strategy:
    matrix:
      include: ${{ fromJSON(needs.setup.outputs.active-versions).ce_active_versions_matrix }}

  Example 3: Both CE and ENT builds (7 branches total)
  strategy:
    matrix:
      include: ${{ fromJSON(needs.setup.outputs.active-versions).all_active_versions_matrix }}

  Using matrix fields in workflow steps:
  - uses: actions/checkout@v4
    with:
      ref: ${{ matrix.branch }}  (e.g., "main", "ce/main", "release/1.19.x+ent")
  - name: Build
    env:
      EDITION: ${{ matrix.edition }}  ("ce" or "ent")
      VERSION: ${{ matrix.version }}  (e.g., "main", "1.19.x", "2.0.x")
      IS_LTS: ${{ matrix.lts }}       (true or false)
  - name: Conditional logic based on edition
    run: |
      if [ "${{ matrix.edition }}" = "ce" ]; then
        echo "Building Community Edition"
      else
        echo "Building Enterprise Edition"
      fi`,
	}

	activeVersionsCmd.PersistentFlags().BoolVar(&listReleaseActiveVersionsReq.WriteToGithubOutput, "github-output", false, "Write 'active-versions' to $GITHUB_OUTPUT for use in GitHub Actions workflows")
	activeVersionsCmd.PersistentFlags().BoolVar(&listReleaseActiveVersionsReq.IncludeMain, "include-main", false, "Include 'main' branch in output (both ENT and CE if --include-ce-prefix is set)")
	activeVersionsCmd.PersistentFlags().StringVar(&listReleaseActiveVersionsReq.CEPrefix, "include-ce-prefix", "", "Prefix to add to CE branches (e.g., 'ce' results in 'ce/release/<version>' and 'ce/main')")

	return activeVersionsCmd
}

func runListActiveVersionsReq(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	listReleaseActiveVersionsReq.VersionsDecodeRes = rootCfg.versionsDecodeRes
	res, err := listReleaseActiveVersionsReq.Run(cmd.Context(), rootCfg.git)
	if err != nil {
		return err
	}

	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON(listReleaseActiveVersionsReq.CEPrefix)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.ToTable(listReleaseActiveVersionsReq.CEPrefix))
	}

	if listReleaseActiveVersionsReq.WriteToGithubOutput {
		output, err := res.ToGithubOutput(listReleaseActiveVersionsReq.IncludeMain, listReleaseActiveVersionsReq.CEPrefix)
		if err != nil {
			return err
		}

		return writeToGithubOutput("active-versions", output)
	}

	return err
}
