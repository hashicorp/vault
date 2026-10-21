// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

// entSuffixRepos is the set of repository names whose tags carry the +ent suffix.
var entSuffixRepos = map[string]struct{}{
	"vault-enterprise": {},
}

var checkVersionTagsReq = &github.CheckVersionTagsReq{}

func newGithubCheckVersionTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version-tags --owner <org> --repo <repo> --versions <v1 v2 ...>",
		Short: "Validate that versions satisfy tag sequencing rules",
		Long: `Fetches all repository tags from GitHub and validates each version against
four tag sequencing rules:

  Rule 0: The exact tag must not already exist.
  Rule 1: A GA version (patch=0) requires at least one RC tag for that major.minor.
  Rule 2: An rcN version (N>1) requires the rc(N-1) tag to exist.
  Rule 3: An rc1 version must be the first tag for that major.minor.
  Rule 4: A minor release (patch>0) requires the previous patch-level tag to exist.

For repositories in the enterprise set (e.g. vault-enterprise) the +ent tag
suffix is applied automatically unless --tag-suffix is provided explicitly.
Use --tag-prefix and --tag-suffix to override the tag format for any project.
When --github-output is set the results are written to $GITHUB_OUTPUT.`,
		RunE: runGithubCheckVersionTagsCmd,
		Args: cobra.NoArgs,
	}

	var rawVersions string
	cmd.PersistentFlags().StringVarP(&checkVersionTagsReq.Owner, "owner", "o", "hashicorp", "GitHub organization")
	cmd.PersistentFlags().StringVarP(&checkVersionTagsReq.Repo, "repo", "r", "", "GitHub repository")
	cmd.PersistentFlags().StringVar(&rawVersions, "versions", "", "Space-separated list of versions to validate (e.g. '1.21.0-rc1 1.20.4 1.19.10')")
	cmd.PersistentFlags().StringVar(&checkVersionTagsReq.TagPrefix, "tag-prefix", "v", "Tag prefix (e.g. 'v')")
	cmd.PersistentFlags().StringVar(&checkVersionTagsReq.TagSuffix, "tag-suffix", "", "Tag suffix applied to every tag (e.g. '+ent'). Derived automatically for known enterprise repos when omitted.")
	cmd.PersistentFlags().BoolVar(&checkVersionTagsReq.WriteToGithubOutput, "github-output", false, "Write outputs to $GITHUB_OUTPUT")

	for _, flag := range []string{"repo", "versions"} {
		if err := cmd.MarkPersistentFlagRequired(flag); err != nil {
			panic(err)
		}
	}

	cmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		checkVersionTagsReq.Versions = strings.Fields(rawVersions)

		// Apply the +ent suffix automatically for known enterprise repos unless
		// the caller has already set --tag-suffix explicitly.
		if !c.Flags().Changed("tag-suffix") {
			if _, isEnt := entSuffixRepos[checkVersionTagsReq.Repo]; isEnt {
				checkVersionTagsReq.TagSuffix = "+ent"
			}
		}

		return nil
	}

	return cmd
}

func runGithubCheckVersionTagsCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	res, err := checkVersionTagsReq.Run(cmd.Context(), githubCmdState.GithubV3)
	if err != nil {
		return fmt.Errorf("check-version-tags: %w", err)
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

	if checkVersionTagsReq.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}
		return writeToGithubOutput("version-tags", output)
	}

	if !res.Passed {
		return fmt.Errorf("version tag validation failed:\n%s", res.String())
	}

	return nil
}
