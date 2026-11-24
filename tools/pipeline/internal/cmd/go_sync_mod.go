// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"errors"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
)

var goSyncModReq = &golang.SyncModReq{
	A:        &golang.ModSource{},
	B:        &golang.ModSource{},
	DiffOpts: golang.DefaultSyncDiffOpts(),
}

func newGoSyncModCmd() *cobra.Command {
	goSyncModCmd := &cobra.Command{
		Use:   "mod </path/a/go.mod> </path/b/go.mod> [flags]",
		Short: "sync a destination go.mod with a source go.mod",
		Long: `Synchronize a destination go.mod file with a source go.mod file.

We read both files from disk, compute what differs between them, and apply the
source (A) directives to the destination (B) in-place. No git operations, no
tidy — just the file mutation. Tidy is the caller's responsibility.

Examples:
  # Sync every directive from A into B.
  pipeline go sync mod a.mod b.mod

  # Only update requires/replaces that already exist in both files (don't add new ones).
  pipeline go sync mod --strict-replace=false --strict-require=false a.mod b.mod

  # Exclude specific modules from the sync.
  pipeline go sync mod \
    --strict-replace=false \
    --strict-require=false \
    --exclude-require github.com/hashicorp/vault-plugin-auth-azure \
    --exclude-require github.com/hashicorp/vault-plugin-secrets-azure \
    --exclude-require github.com/hashicorp/vault-plugin-secrets-openldap \
    a.mod b.mod

  # The real-world usage — syncing ent into ce with the azure/openldap modules
  # excluded (they live only in enterprise).
  pipeline go sync mod /tmp/ent.mod /tmp/ce.mod \
    --strict-replace=false \
    --strict-require=false \
    --exclude-require github.com/hashicorp/vault-plugin-auth-azure \
    --exclude-require github.com/hashicorp/vault-plugin-secrets-azure \
    --exclude-require github.com/hashicorp/vault-plugin-secrets-openldap \
    --log debug
`,
		RunE: runGoSyncModCmd,
		Args: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			switch len(args) {
			case 2:
				if err := golang.SetUpGoModSourceFromPath(args[0], goSyncModReq.A); err != nil {
					return err
				}
				return golang.SetUpGoModSourceFromPath(args[1], goSyncModReq.B)
			case 0, 1:
				return errors.New("two local file paths are required")
			default:
				return fmt.Errorf("expected two path arguments, got %d", len(args))
			}
		},
	}

	goSyncModCmd.PersistentFlags().BoolVar(&goSyncModReq.DiffOpts.ParseLax, "lax", false, "parse go.mod files in lax mode, ignoring unknown directives")
	goSyncModCmd.PersistentFlags().BoolVar(&goSyncModReq.DiffOpts.Require, "require", true, "sync require directives")
	goSyncModCmd.PersistentFlags().BoolVar(&goSyncModReq.DiffOpts.Replace, "replace", true, "sync replace directives")
	goSyncModCmd.PersistentFlags().BoolVar(&goSyncModReq.DiffOpts.Go, "go", false, "sync the 'go' directive")
	goSyncModCmd.PersistentFlags().BoolVar(&goSyncModReq.DiffOpts.Toolchain, "toolchain", false, "sync the 'toolchain' directive")
	goSyncModCmd.PersistentFlags().BoolVar(&goSyncModReq.DiffOpts.StrictDiffRequire, "strict-require", true, "add requires present in source but absent in dest. When false, only update shared requires.")
	goSyncModCmd.PersistentFlags().BoolVar(&goSyncModReq.DiffOpts.StrictDiffReplace, "strict-replace", true, "add replaces present in source but absent in dest. When false, only update shared replaces.")
	goSyncModCmd.PersistentFlags().StringArrayVar(&goSyncModReq.DiffOpts.ExcludeRequire, "exclude-require", nil, "repeatable glob/literal pattern. Skips matching module paths in require directives.")
	goSyncModCmd.PersistentFlags().StringArrayVar(&goSyncModReq.DiffOpts.ExcludeReplace, "exclude-replace", nil, "repeatable glob/literal pattern. Skips matching module paths in replace directives (old or new path).")
	goSyncModCmd.PersistentFlags().BoolVarP(&goSyncModReq.DryRun, "dry-run", "n", false, "simulate the sync without writing changes to disk")

	return goSyncModCmd
}

func runGoSyncModCmd(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	res, err := goSyncModReq.Run(cmd.Context())
	if err != nil {
		return err
	}

	return printGoSyncModResult(cmd, res)
}

// printGoSyncModResult writes the sync result to cmd.OutOrStdout() in the format
// requested via rootCfg.format.
func printGoSyncModResult(cmd *cobra.Command, res *golang.SyncModRes) error {
	out := cmd.OutOrStdout()

	switch rootCfg.format {
	case "json":
		b, jsonErr := res.ToJSON()
		if jsonErr != nil {
			return jsonErr
		}
		fmt.Fprintln(out, string(b))
	case "markdown":
		fmt.Fprintln(out, res.ToMarkdown())
	default:
		if text := res.ToTable(); text != "" {
			fmt.Fprintln(out, text)
		}
	}

	return nil
}
