// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"errors"
	"fmt"

	githubpkg "github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
)

var githubSyncGoModReq = &githubpkg.SyncGoModReq{
	DiffOpts:    golang.DefaultSyncDiffOpts(),
	Commit:      true,
	PullRequest: true,
}

func newGithubSyncGoModCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-mod [flags]",
		Short: "sync go.mod files from a remote A branch into a new branch forked from B, then push and open a PR",
		Long: `Synchronize go.mod files from a remote A-side branch into a new intermediate
branch forked from a B-side branch. The intermediate branch is pushed to the
TO remote and a pull request is opened against the B-side base branch.

We clone (or reuse) the B-side repository, add the A-side as a separate remote,
read each go.mod file via git show, apply our standard Go module diff to the
working-tree copy, run go mod tidy for any file with changes, commit, push, and
open a PR.

Examples:
  # Sync the main into ce/main branch. Don't strictly require equal replace or
  # require directives so that only shared modules are synchronized. We also
  # explicitly exlude plugins who have enterprise counterparts that are
  # different.
  pipeline github sync go-mod \
    --a-branch main \
    --b-branch ce/main \
    --strict-replace=false \
    --strict-require=false \
    -p go.mod \
    -p api/go.mod \
    -p api/auth/approle/go.mod \
    -p api/auth/aws/go.mod \
    -p api/auth/azure/go.mod \
    -p api/auth/cert/go.mod \
    -p api/auth/gcp/go.mod \
    -p api/auth/kubernetes/go.mod \
    -p api/auth/ldap/go.mod \
    -p api/auth/userpass/go.mod \
    -p tools/pipeline/go.mod \
    -p sdk/go.mod \
    --exclude-require github.com/hashicorp/vault-plugin-auth-azure \
    --exclude-require github.com/hashicorp/vault-plugin-secrets-azure \
    --exclude-require github.com/hashicorp/vault-plugin-secrets-openldap

  # Cross-repo sync (hashicorp/vault-enterprise → hashicorp/vault).
  pipeline github sync go-mod \
    --a-owner hashicorp --a-repo vault-enterprise --a-branch main \
    --b-owner hashicorp --b-repo vault --b-branch main

  # With a ticket ID (branch becomes VAULT-1234-main-into-ce-main).
  pipeline github sync go-mod \
    --a-branch main --b-branch ce/main \
    --ticket VAULT-1234

  # Use an existing local clone to skip re-cloning.
  pipeline github sync go-mod \
    --a-branch main --b-branch ce/main \
    --repo-dir /path/to/vault-enterprise

  # Push only — no PR.
  pipeline github sync go-mod \
    --a-branch main --b-branch ce/main \
    --pull-request=false
`,
		RunE: runGithubSyncGoModCmd,
		Args: cobra.NoArgs,
	}

	// A-side flags
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.AHost, "a-host", "github.com", "GitHub host for the A (source) remote")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.AOwner, "a-owner", "hashicorp", "GitHub owner for the A (source) remote")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.ARepo, "a-repo", "vault-enterprise", "GitHub repo for the A (source) remote")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.ABranch, "a-branch", "", "source branch to read go.mod files from (required)")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.AOrigin, "a-origin", "", "remote name for A in the local clone (auto-resolved when empty)")

	// B-side flags
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.BHost, "b-host", "github.com", "GitHub host for the B (dest) remote")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.BOwner, "b-owner", "hashicorp", "GitHub owner for the B (dest) remote")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.BRepo, "b-repo", "vault-enterprise", "GitHub repo for the B (dest) remote")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.BBranch, "b-branch", "", "dest branch to fork the intermediate branch from (required)")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.BOrigin, "b-origin", "", "remote name for B in the local clone (auto-resolved when empty)")

	// TO-side flags (all default to the corresponding B value)
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.ToHost, "to-host", "", "push target GitHub host (defaults to --b-host)")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.ToOwner, "to-owner", "", "push target GitHub owner (defaults to --b-owner)")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.ToRepo, "to-repo", "", "push target GitHub repo (defaults to --b-repo)")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.ToBranch, "to-branch", "", "intermediate branch name (auto-generated when empty)")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.ToOrigin, "to-origin", "", "push target remote name in the local clone (defaults to --b-origin)")

	// Repo isolation
	cmd.PersistentFlags().StringVarP(&githubSyncGoModReq.RepoDir, "repo-dir", "d", "", "path to an existing local clone. A temp dir is created and cleaned up when empty.")

	// Path flag
	cmd.PersistentFlags().StringArrayVarP(&githubSyncGoModReq.Paths, "path", "p", []string{"go.mod"}, "repeatable. Path to a go.mod file to sync.")

	// Commit / ticket flags
	cmd.PersistentFlags().BoolVar(&githubSyncGoModReq.Commit, "commit", true, "commit the synced changes to the intermediate branch. Only useful with --repo-dir when false.")
	cmd.PersistentFlags().StringVarP(&githubSyncGoModReq.CommitMessage, "message", "m", "", "commit subject override (default: auto-generated from branch names)")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.Ticket, "ticket", "", "bare ticket ID prepended to the branch name and commit subject, e.g. VAULT-1234")

	// PR flags
	cmd.PersistentFlags().BoolVar(&githubSyncGoModReq.PullRequest, "pull-request", true, "open a pull request after pushing the intermediate branch")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.PRTitle, "pr-title", "", "PR title override (default: auto-generated from branch names)")
	cmd.PersistentFlags().StringVar(&githubSyncGoModReq.PRBase, "pr-base", "", "PR base branch (defaults to --b-branch)")
	cmd.PersistentFlags().StringArrayVar(&githubSyncGoModReq.PRReviewers, "pr-reviewer", nil, "repeatable. GitHub login to request a review from.")
	cmd.PersistentFlags().StringArrayVar(&githubSyncGoModReq.PRAssignees, "pr-assignee", nil, "repeatable. GitHub login to assign the PR to.")

	// Sync option flags
	cmd.PersistentFlags().BoolVar(&githubSyncGoModReq.DiffOpts.ParseLax, "lax", false, "parse go.mod files in lax mode, ignoring unknown directives")
	cmd.PersistentFlags().BoolVar(&githubSyncGoModReq.DiffOpts.Require, "require", true, "sync require directives")
	cmd.PersistentFlags().BoolVar(&githubSyncGoModReq.DiffOpts.Replace, "replace", true, "sync replace directives")
	cmd.PersistentFlags().BoolVar(&githubSyncGoModReq.DiffOpts.Go, "go", false, "sync the 'go' directive")
	cmd.PersistentFlags().BoolVar(&githubSyncGoModReq.DiffOpts.Toolchain, "toolchain", false, "sync the 'toolchain' directive")
	cmd.PersistentFlags().BoolVar(&githubSyncGoModReq.DiffOpts.StrictDiffRequire, "strict-require", true, "add requires present in source but absent in dest. When false, only update shared requires.")
	cmd.PersistentFlags().BoolVar(&githubSyncGoModReq.DiffOpts.StrictDiffReplace, "strict-replace", true, "add replaces present in source but absent in dest. When false, only update shared replaces.")
	cmd.PersistentFlags().StringArrayVar(&githubSyncGoModReq.DiffOpts.ExcludeRequire, "exclude-require", nil, "repeatable glob/literal pattern. Skips matching module paths in require directives.")
	cmd.PersistentFlags().StringArrayVar(&githubSyncGoModReq.DiffOpts.ExcludeReplace, "exclude-replace", nil, "repeatable glob/literal pattern. Skips matching module paths in replace directives.")

	// Required flags
	_ = cmd.MarkPersistentFlagRequired("a-branch")
	_ = cmd.MarkPersistentFlagRequired("b-branch")

	return cmd
}

func runGithubSyncGoModCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	res, runErr := githubSyncGoModReq.Run(cmd.Context(), githubCmdState.GithubV3, rootCfg.git)
	printErr := printGithubSyncGoModResult(cmd, res, runErr)

	return errors.Join(runErr, printErr)
}

// printGithubSyncGoModResult renders the full sync result in the requested output
// format. Metadata (A/B/branch/commit/PR) is printed first, followed by the
// per-path results table.
func printGithubSyncGoModResult(cmd *cobra.Command, res *githubpkg.SyncGoModRes, err error) error {
	out := cmd.OutOrStdout()

	switch rootCfg.format {
	case "json":
		b, jsonErr := res.ToJSON()
		if jsonErr != nil {
			return jsonErr
		}
		fmt.Fprintln(out, string(b))
	case "markdown":
		fmt.Fprintln(out, res.ToMarkdown(err))
	default:
		fmt.Fprintln(out, res.ToString(err))
	}

	return nil
}
