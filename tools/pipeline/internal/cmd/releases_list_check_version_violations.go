// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
)

var checkVersionViolationsReq = &releases.CheckVersionViolationsReq{}

func newReleasesCheckVersionViolationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-version-violations --versions <v1 v2 ...>",
		Short: "Check a sorted list of versions for rule violations",
		Long: `Checks a sorted, normalized list of versions for the following violations:

  Rule 1: Only one version per major.minor series is allowed.
  Rule 2: Only one version with patch level 0 is allowed.
  Rule 3: A GA (patch=0) or RC version must be the highest version in the list.
  Rule 4: GA and RC versions cannot coexist in the same release set.
  Rule 5: RC versions must have patch level 0 (e.g. 1.2.0-rc1, not 1.2.1-rc1).
  Rule 6: Version format must be valid (X.Y.Z or X.Y.Z-rcN, no leading zeros).
  Rule 7: Only one GA version is allowed per release set.

Input should be the sorted_versions output from 'pipeline releases list sanitize-sort-versions'.
When --github-output is set the results are written to $GITHUB_OUTPUT using the
same output key names as the check-for-version-violations-ent.sh script.`,
		RunE: runReleasesCheckVersionViolationsCmd,
		Args: cobra.NoArgs,
	}

	var rawVersions string
	cmd.PersistentFlags().StringVar(&rawVersions, "versions", "", "Space-separated sorted list of versions to check (e.g. '1.21.0-rc1 1.20.4 1.19.10')")
	cmd.PersistentFlags().BoolVar(&checkVersionViolationsReq.WriteToGithubOutput, "github-output", false, "Write outputs to $GITHUB_OUTPUT")

	if err := cmd.MarkPersistentFlagRequired("versions"); err != nil {
		panic(err)
	}

	cmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		checkVersionViolationsReq.Versions = strings.Fields(rawVersions)
		return nil
	}

	return cmd
}

func runReleasesCheckVersionViolationsCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	res, err := checkVersionViolationsReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("check-version-violations: %w", err)
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

	if checkVersionViolationsReq.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}
		return writeToGithubOutput("check-version-violations", output)
	}

	if res.HasViolations {
		return fmt.Errorf("version violations detected:\n%s", res.String())
	}

	return nil
}
