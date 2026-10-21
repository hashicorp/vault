// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
)

var sanitizeSortVersionsReq = &releases.SanitizeSortVersionsReq{
	Descending: true, // default matches the original shell script behaviour
}

func newReleasesSanitizeSortVersionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sanitize-sort-versions --versions <v1 v2 ...>",
		Short: "Sanitize and sort a list of version strings",
		Long: `Sanitize, normalize, and sort a space-separated list of version strings.
Versions are sorted in descending (highest-first) order by default. Use --asc to
sort ascending instead. Versions carrying a -lts suffix are tracked separately.
Versions that do not match the expected format are reported as invalid.

Outputs sorted_versions, lts_versions, and invalid_versions. When --github-output
is set the results are written to $GITHUB_OUTPUT.`,
		RunE: runReleasesSanitizeSortVersionsCmd,
		Args: cobra.NoArgs,
	}

	var rawVersions string
	var asc bool
	cmd.PersistentFlags().StringVar(&rawVersions, "versions", "", "Space-separated list of version strings to process (e.g. '1.21.0 1.20.4-lts 1.19.10+ent')")
	cmd.PersistentFlags().BoolVar(&asc, "asc", false, "Sort versions ascending (lowest-first) instead of the default descending")
	cmd.PersistentFlags().BoolVar(&sanitizeSortVersionsReq.WriteToGithubOutput, "github-output", false, "Write outputs to $GITHUB_OUTPUT")

	err := cmd.MarkPersistentFlagRequired("versions")
	if err != nil {
		panic(err)
	}

	cmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		sanitizeSortVersionsReq.Versions = strings.Fields(rawVersions)
		sanitizeSortVersionsReq.Descending = !asc
		return nil
	}

	return cmd
}

func runReleasesSanitizeSortVersionsCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	res, err := sanitizeSortVersionsReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("sanitize-sort-versions: %w", err)
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

	if sanitizeSortVersionsReq.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}
		return writeToGithubOutput("sanitize-sort-versions", output)
	}

	return nil
}
