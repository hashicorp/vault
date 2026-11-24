// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
)

// verifyBinariesCustomFetcher allows unit tests to inject a mock fetcher.
var verifyBinariesCustomFetcher releases.BinaryVersionsFetcher

func newReleasesVerifyBinariesCmd() *cobra.Command {
	var (
		versionsString          string
		shouldWriteGithubOutput bool
		failOnMissing           bool
	)

	cmd := &cobra.Command{
		Use:   "binaries --versions <v1 v2 ...>",
		Short: "Verify production release binaries on releases.hashicorp.com",
		Long: `Verify production release binaries on releases.hashicorp.com.

Checks that the Vault release binaries for given versions exist on releases.hashicorp.com
with all required distributions, variants, and architectures.

OUTPUT MODES:
1. Default (human-readable table):
   Displays verification results in a formatted table with details of valid/missing items.
2. --format json:
   Outputs full JSON report with all verified metadata, missing variants, missing OS, etc.
3. --github-output:
   Writes output formatted for GitHub Actions to $GITHUB_OUTPUT.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			req := releases.NewVerifyBinariesReq()
			req.VersionsString = versionsString
			if rootCfg.versionsDecodeRes != nil {
				req.VersionsConfig = rootCfg.versionsDecodeRes.Config
			}
			if verifyBinariesCustomFetcher != nil {
				req.BinaryFetcher = verifyBinariesCustomFetcher
			}

			res, err := req.Run(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to verify binaries: %w", err)
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

			if shouldWriteGithubOutput {
				output, err := res.ToGithubOutput()
				if err != nil {
					return err
				}
				if err := writeToGithubOutput("release_binary_verification", output); err != nil {
					return err
				}
			}

			if failOnMissing && res.HasMissing {
				return fmt.Errorf("binary verification failed: missing binaries, variants, or OS detected")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&versionsString, "versions", "", "Space-separated list of vault versions to verify (e.g. '2.2.1 1.21.13 1.19.24')")
	cmd.Flags().BoolVar(&shouldWriteGithubOutput, "github-output", false, "Write output to $GITHUB_OUTPUT")
	cmd.Flags().BoolVar(&failOnMissing, "fail-on-missing", false, "Exit with non-zero error if missing versions, variants, or OS are detected")

	_ = cmd.MarkFlagRequired("versions")

	return cmd
}
