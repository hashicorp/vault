// Copyright IBM Corp. 2016, 2026
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
We sync every go.mod before we tidy any of them, because tidying a module depends
on the go.mod of the local modules it replaces. We then run 'go mod tidy' for any
file with changes and commit the result.

The go.mod files to sync come from --path, relative to the working directory, or
from every module in the go.work at the root of <source-branch> with
--workspace-modules. In workspace mode we tidy every module in the go.work of
<dest-branch> with GOWORK=off. The go.mod and go.sum files that change are committed.
We work from the root of the repository wherever you run the command from. The
go commands use your environment, so set GOPRIVATE, GONOSUMDB, and GOFLAGS as you
would for go mod tidy.

A new forked intermediate branch from <dest-branch> is created with the changes
and we never touch dest directly.

Examples:
# Sync go.mod from main into ce/main (creates a main-into-ce-main branch).
  pipeline git sync go-mod main ce/main

# With a ticket ID (creates VAULT-1234-main-into-ce-main).
  pipeline git sync go-mod main ce/main --ticket VAULT-1234

# Leave the changes in the working tree without committing them.
  pipeline git sync go-mod main ce/main --commit=false

# Sync only specific modules, relative to the working directory.
  pipeline git sync go-mod main ce/main -p go.mod -p api/go.mod -p sdk/go.mod

# Sync every module in the go.work of main, with the modules that only exist
# in enterprise excluded. This represents something much closer to what we
# would actually do if we triggered the sync from local branches.
  GOPRIVATE=github.com/hashicorp pipeline git sync go-mod main ce/main \
    --workspace-modules \
    --strict-replace=false \
    --strict-require=false \
    --exclude-require github.com/hashicorp/vault-plugin-auth-azure \
    --exclude-require github.com/hashicorp/vault-plugin-secrets-azure \
    --exclude-require github.com/hashicorp/vault-plugin-secrets-openldap \
    --ticket VAULT-49777
`,
		RunE: runGitSyncGoModCmd,
		Args: cobra.ExactArgs(2),
	}

	// Path flags
	cmd.PersistentFlags().StringArrayVarP(&gitSyncGoModReq.Paths, "path", "p", nil, "repeatable. Path to a go.mod file to sync, relative to the working directory. Defaults to go.mod unless --workspace-modules is set.")
	cmd.PersistentFlags().BoolVar(&gitSyncGoModReq.WorkspaceModules, "workspace-modules", false, "sync the go.mod of every module in the go.work at the root of the source branch instead of --path, then tidy every module in the go.work of the dest branch with GOWORK=off and commit the go.mod and go.sum files that change.")

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
