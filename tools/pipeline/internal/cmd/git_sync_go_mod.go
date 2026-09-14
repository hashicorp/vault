// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"errors"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
)

var gitSyncGoModReq = &git.SyncGoModReq{
	DiffOpts: golang.DefaultSyncDiffOpts(),
	Commit:   true,
}

func newGitSyncGoModCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-mod <source-branch> <dest-branch> [flags]",
		Short: "sync go.mod files from a source branch into a dest branch",
		Long: `Synchronize go.mod files between two local branches.

We read the source go.mod files via git, compute what differs from the
working-tree copies on <dest-branch>, and apply the source directives in-place.
We then run 'go mod tidy' for any file with changes and commit the result.

A new forked intermediate branch from <dest-branch> is created with the changes
and we never touch dest directly.

Examples:
# Sync go.mod from main into ce/main (creates a main-into-ce-main branch).
  pipeline git sync go-mod main ce/main

# With a ticket ID (creates VAULT-1234-main-into-ce-main).
  pipeline git sync go-mod main ce/main --ticket VAULT-1234

# Stage changes without committing.
  pipeline git sync go-mod main ce/main --commit=false

# Sync multiple paths with excluded modules and a ticket. This represents
# something much closer to what we would actually do if we triggered the
# sync from local branches.
  GOPRIVATE=github.com/hashicorp pipeline git sync go-mod main ce/main \
    --strict-replace=false \
    --strict-require=false \
    -p ../../go.mod \
    -p ../../api/go.mod \
    -p ../../api/auth/approle/go.mod \
    -p ../../api/auth/aws/go.mod \
    -p ../../api/auth/azure/go.mod \
    -p ../../api/auth/cert/go.mod \
    -p ../../api/auth/gcp/go.mod \
    -p ../../api/auth/kubernetes/go.mod \
    -p ../../api/auth/ldap/go.mod \
    -p ../../api/auth/userpass/go.mod \
    -p go.mod \
    -p ../../sdk/go.mod \
    --exclude-require github.com/hashicorp/vault-plugin-auth-azure \
    --exclude-require github.com/hashicorp/vault-plugin-secrets-azure \
    --exclude-require github.com/hashicorp/vault-plugin-secrets-openldap \
    --ticket VAULT-49777
`,
		RunE: runGitSyncGoModCmd,
		Args: cobra.ExactArgs(2),
	}

	// Path flag
	cmd.PersistentFlags().StringArrayVarP(&gitSyncGoModReq.Paths, "path", "p", []string{"go.mod"}, "repeatable. Path to a go.mod file to sync.")

	// Branch / commit flags
	cmd.PersistentFlags().StringVar(&gitSyncGoModReq.NewBranch, "branch", "", "intermediate branch name (one is created when omitted)")
	cmd.PersistentFlags().BoolVar(&gitSyncGoModReq.Commit, "commit", true, "commit the synced changes to the intermediate branch")
	cmd.PersistentFlags().StringVarP(&gitSyncGoModReq.CommitMessage, "message", "m", "", "commit subject override")
	cmd.PersistentFlags().StringVar(&gitSyncGoModReq.Ticket, "ticket", "", "ticket ID to include in the message, e.g. VAULT-1234")

	// Sync option flags
	cmd.PersistentFlags().BoolVar(&gitSyncGoModReq.DiffOpts.ParseLax, "lax", false, "parse go.mod files in lax mode, ignoring unknown directives")
	cmd.PersistentFlags().BoolVar(&gitSyncGoModReq.DiffOpts.Require, "require", true, "sync require directives")
	cmd.PersistentFlags().BoolVar(&gitSyncGoModReq.DiffOpts.Replace, "replace", true, "sync replace directives")
	cmd.PersistentFlags().BoolVar(&gitSyncGoModReq.DiffOpts.Go, "go", false, "sync the 'go' directive")
	cmd.PersistentFlags().BoolVar(&gitSyncGoModReq.DiffOpts.Toolchain, "toolchain", false, "sync the 'toolchain' directive")
	cmd.PersistentFlags().BoolVar(&gitSyncGoModReq.DiffOpts.StrictDiffRequire, "strict-require", true, "add requires present in source but absent in dest. When false, only update shared requires.")
	cmd.PersistentFlags().BoolVar(&gitSyncGoModReq.DiffOpts.StrictDiffReplace, "strict-replace", true, "add replaces present in source but absent in dest. When false, only update shared replaces.")
	cmd.PersistentFlags().StringArrayVar(&gitSyncGoModReq.DiffOpts.ExcludeRequire, "exclude-require", nil, "repeatable glob/literal pattern. Skips matching module paths in require directives.")
	cmd.PersistentFlags().StringArrayVar(&gitSyncGoModReq.DiffOpts.ExcludeReplace, "exclude-replace", nil, "repeatable glob/literal pattern. Skips matching module paths in replace directives.")

	return cmd
}

func runGitSyncGoModCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	gitSyncGoModReq.SourceBranch = args[0]
	gitSyncGoModReq.DestBranch = args[1]

	res, err := gitSyncGoModReq.Run(cmd.Context(), rootCfg.git)
	if err != nil {
		// Still print what we have on partial failure so the caller can see which
		// paths succeeded and which didn't.
		if res != nil {
			err = errors.Join(err, printSyncGoModResult(cmd, res))
		}
		return fmt.Errorf("git sync go-mod: %w", err)
	}

	return printSyncGoModResult(cmd, res)
}

// printSyncGoModResult emits the branch/commit header followed by the per-path
// table in the requested output format.
func printSyncGoModResult(cmd *cobra.Command, res *git.SyncGoModRes) error {
	out := cmd.OutOrStdout()

	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(b))
	case "markdown":
		t := res.ToTable()
		t.SetOutputMirror(out)
		t.RenderMarkdown()
	default:
		t := res.ToTable()
		t.SetOutputMirror(out)
		t.Render()
	}

	return nil
}
