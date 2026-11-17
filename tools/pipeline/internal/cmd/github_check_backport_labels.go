// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var (
	checkBackportLabelsReq           = &github.CheckBackportLabelsReq{}
	checkBackportLabelsFailOnMissing bool
)

func newGithubCheckBackportLabelsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backport-labels",
		Short: "Check that backport labels match the active release branches",
		Long: `Compares the 'backport/<version>' labels in the repository to the active
versions in .release/versions.hcl. Every active version (e.g. 1.19.x for
release/1.19.x+ent) must have a backport label, and labels for versions that
are no longer active are reported as stale. The main branch has no label.

By default the command is report-only and exits zero even when labels are
missing or stale. Use --fail-on-missing to exit non-zero when drift remains,
and --create-missing and --delete-stale to fix the drift. New labels copy the color and description of an existing backport label.`,
		Example: `  # Report drift
  pipeline github check backport-labels

  # Report drift and fail if any remains
  pipeline github check backport-labels --fail-on-missing

  # Create missing and delete stale labels
  pipeline github check backport-labels --create-missing --delete-stale`,
		RunE: runGithubCheckBackportLabelsCmd,
		Args: cobra.NoArgs,
	}

	cmd.PersistentFlags().BoolVar(&checkBackportLabelsReq.CreateMissing, "create-missing", false, "Create missing backport labels")
	cmd.PersistentFlags().BoolVar(&checkBackportLabelsReq.DeleteStale, "delete-stale", false, "Delete stale backport labels")
	cmd.PersistentFlags().BoolVar(&checkBackportLabelsFailOnMissing, "fail-on-missing", false, "Exit with non-zero error if missing or stale backport labels remain")
	cmd.PersistentFlags().BoolVar(&checkBackportLabelsReq.WriteToGithubOutput, "github-output", false, "Write outputs to $GITHUB_OUTPUT")

	return cmd
}

func runGithubCheckBackportLabelsCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	checkBackportLabelsReq.Owner = githubCmdState.Owner
	checkBackportLabelsReq.Repo = githubCmdState.Repo
	checkBackportLabelsReq.VersionsDecodeRes = rootCfg.versionsDecodeRes

	res, err := checkBackportLabelsReq.Run(cmd.Context(), githubCmdState.GithubV3.Issues)
	if err != nil {
		return fmt.Errorf("check-backport-labels: %w", err)
	}

	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.ToTable())
	}

	if checkBackportLabelsReq.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}
		if err := writeToGithubOutput("backport-labels", output); err != nil {
			return err
		}
	}

	if checkBackportLabelsFailOnMissing && !res.InSync {
		return fmt.Errorf("backport labels are out of sync: %s", res.String())
	}

	return nil
}
