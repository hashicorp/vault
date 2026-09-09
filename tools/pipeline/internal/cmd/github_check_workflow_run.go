// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var checkWorkflowRunReq = &github.CheckWorkflowRunReq{}

func newGithubCheckWorkflowRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow-run [branches...]",
		Short: "Check that a recent workflow run has all required jobs and succeeded",
		Long: `Check that a recent workflow run has all required jobs and succeeded.

This command searches recent workflow runs to find one that meets all criteria:
  1. All required job patterns are matched
  2. All matched jobs succeeded (status=completed, conclusion=success/skipped)
  3. No other jobs failed (which could prevent required jobs from running)

Search Strategy:
  The command examines workflow runs from newest to oldest (up to --max-depth).
  For each run, it evaluates the jobs:

  • Complete & Successful: All patterns matched, all jobs succeeded → SUCCESS
  • Complete & Failed: All patterns matched, but some jobs failed → FAIL (stop searching)
  • Incomplete & Tainted: Some patterns missing AND other jobs failed → FAIL (stop searching)
  • Incomplete & Clean: Some patterns missing, no failures → CONTINUE to older run

  The "tainted" case is critical: if a run has job failures (even non-matched jobs),
  those failures likely prevented your required jobs from running. The check fails
  immediately rather than continuing to search older runs.

Branch Selection:
  Branches can be specified as positional arguments, or automatically determined from
  active versions in versions.hcl when no branches are provided. Use --only-lts or
  --only-ce to filter which active versions to check.

  Additional branch options:
  • --include-main: Include the 'main' branch in checks
  • --include-ce-prefix: Prefix for CE branches in enterprise repos (e.g., 'ce' results
    in 'ce/release/<version>' and 'ce/main')

Output Options:
  • --hide-successful-jobs: Hide successful jobs in text output (table/markdown) to
    reduce noise when many jobs pass

Examples:
  # Check specific branch for required jobs
  pipeline github check workflow-run main \
    --owner hashicorp \
    --repo vault-enterprise \
    --workflow build.yml \
    --pattern "^artifacts.*$" \
    --pattern "^test.*$"

  # Check multiple branches
  pipeline github check workflow-run main release/1.18.x \
    --workflow build.yml \
    --pattern "^build-.*-amd64$"

  # Check all active LTS branches (uses versions.hcl)
  pipeline github check workflow-run \
    --owner hashicorp \
    --repo vault-enterprise \
    --workflow build.yml \
    --only-lts \
    --pattern "^build-.*-amd64$"

  # Check main branch and all active versions
  pipeline github check workflow-run \
    --workflow build.yml \
    --include-main \
    --pattern "^build-.*-amd64$"

  # Check CE branches in enterprise repo
  pipeline github check workflow-run \
    --workflow build.yml \
    --only-ce \
    --include-ce-prefix ce \
    --pattern "^build-.*-amd64$"

  # Hide successful jobs to reduce output noise
  pipeline github check workflow-run main \
    --workflow build.yml \
    --pattern "^build-.*$" \
    --hide-successful-jobs

  # Check with custom search depth
  pipeline github check workflow-run main \
    --workflow ci.yml \
    --pattern "^test-.*$" \
    --max-depth 100

  # Output as JSON for automation
  pipeline github check workflow-run main \
    --workflow build.yml \
    --pattern "^build-.*$" \
    --format json

Exit Codes:
  0: All required jobs found and successful
  1: Error, required jobs missing/failed, or tainted run detected`,
		RunE: runGithubCheckWorkflowRunCmd,
		Args: cobra.ArbitraryArgs,
	}

	cmd.PersistentFlags().StringVarP(&checkWorkflowRunReq.Owner, "owner", "o", "hashicorp", "GitHub repository owner")
	cmd.PersistentFlags().StringVarP(&checkWorkflowRunReq.Repo, "repo", "r", "vault", "GitHub repository name")
	cmd.PersistentFlags().StringVarP(&checkWorkflowRunReq.Workflow, "workflow", "w", "", "Workflow file name (e.g., build.yml)")
	cmd.PersistentFlags().StringSliceVarP(&checkWorkflowRunReq.Patterns, "pattern", "p", []string{}, "Regex patterns for required job names (repeatable)")
	cmd.PersistentFlags().BoolVar(&checkWorkflowRunReq.OnlyLTS, "only-lts", false, "Only check LTS branches from active versions")
	cmd.PersistentFlags().BoolVar(&checkWorkflowRunReq.OnlyCE, "only-ce", false, "Only check CE active branches from active versions")
	cmd.PersistentFlags().BoolVar(&checkWorkflowRunReq.IncludeMain, "include-main", false, "Include 'main' branch in checks")
	cmd.PersistentFlags().StringVar(&checkWorkflowRunReq.CEBranchPrefix, "include-ce-prefix", "", "Prefix for CE branches in enterprise repos (e.g., 'ce' results in 'ce/release/<version>' and 'ce/main')")
	cmd.PersistentFlags().UintVar(&checkWorkflowRunReq.MaxDepth, "max-depth", 50, "Maximum depth to search workflow run history for matching run")
	cmd.PersistentFlags().BoolVar(&checkWorkflowRunReq.HideSuccessfulJobs, "hide-successful-jobs", false, "Hide successful jobs in text output (table/markdown)")
	cmd.PersistentFlags().BoolVar(&checkWorkflowRunReq.WriteToGithubOutput, "github-output", false, "Write result to $GITHUB_OUTPUT")
	cmd.PersistentFlags().StringVar(&checkWorkflowRunReq.GithubOutputKey, "github-output-key", "check-workflow-run", "The output key to use when writing to $GITHUB_OUTPUT")

	err := cmd.MarkPersistentFlagRequired("workflow")
	if err != nil {
		panic(err)
	}

	err = cmd.MarkPersistentFlagRequired("pattern")
	if err != nil {
		panic(err)
	}

	return cmd
}

func runGithubCheckWorkflowRunCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	// Set branches from positional arguments
	checkWorkflowRunReq.Branches = args

	// Pass versions configuration for active branch resolution
	checkWorkflowRunReq.Releases = rootCfg.versionsDecodeRes

	res, err := checkWorkflowRunReq.Run(cmd.Context(), githubCmdState.GithubV3)
	if err != nil {
		return fmt.Errorf("check workflow run: %w", err)
	}

	// Handle output formatting
	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "markdown":
		fmt.Println(res.ToMarkdown())
	default:
		fmt.Println(res.ToTable())
	}

	// Write to GitHub output if requested
	if checkWorkflowRunReq.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}
		return writeToGithubOutput(checkWorkflowRunReq.GithubOutputKey, output)
	}

	// Return error if check failed
	if !res.Success {
		return fmt.Errorf("workflow run check failed: %s", res.String())
	}

	return nil
}
