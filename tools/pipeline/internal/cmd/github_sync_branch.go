// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"errors"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/spf13/cobra"
)

var syncGithubBranchReq = github.SyncBranchReq{}

func newSyncGithubBranchCmd() *cobra.Command {
	syncBranchCmd := &cobra.Command{
		Use:   "branch",
		Short: "Sync a branch between two repositories",
		Long:  "Sync a branch between two repositories by merging the --from branch into the --to branch and pushing the result on success",
		RunE:  runSyncGithubBranchCmd,
	}

	syncBranchCmd.PersistentFlags().StringVar(&syncGithubBranchReq.FromOrigin, "from-origin", "from", "The origin name to use for the from-branch")
	syncBranchCmd.PersistentFlags().StringVar(&syncGithubBranchReq.FromOwner, "from-owner", "hashicorp", "The Github organization hosting the from branch")
	syncBranchCmd.PersistentFlags().StringVar(&syncGithubBranchReq.FromRepo, "from-repo", "vault-enterprise", "The Github repository to sync from")
	syncBranchCmd.PersistentFlags().StringVar(&syncGithubBranchReq.FromBranch, "from-branch", "", "The name of the branch we want to sync from")
	syncBranchCmd.PersistentFlags().StringVar(&syncGithubBranchReq.ToOrigin, "to-origin", "to", "The origin name to use for the to-branch")
	syncBranchCmd.PersistentFlags().StringVar(&syncGithubBranchReq.ToOwner, "to-owner", "hashicorp", "The Github organization hosting the to branch")
	syncBranchCmd.PersistentFlags().StringVar(&syncGithubBranchReq.ToRepo, "to-repo", "vault", "The Github repository to sync to")
	syncBranchCmd.PersistentFlags().StringVar(&syncGithubBranchReq.ToBranch, "to-branch", "", "The name of the branch we want to sync to")
	syncBranchCmd.PersistentFlags().StringVarP(&syncGithubBranchReq.RepoDir, "repo-dir", "d", "", "The path to the vault repository dir. If not set a temporary directory will be used")
	syncBranchCmd.PersistentFlags().StringSliceVarP(&syncGithubBranchReq.DisallowedGroups, "disallowed-groups", "g", nil, "Enable changed file group and disallow if any files match the given groups")

	err := syncBranchCmd.MarkPersistentFlagRequired("from-branch")
	if err != nil {
		panic(err)
	}
	err = syncBranchCmd.MarkPersistentFlagRequired("to-branch")
	if err != nil {
		panic(err)
	}

	return syncBranchCmd
}

func runSyncGithubBranchCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	syncGithubBranchReq.DecodeRes = rootCfg.configDecodeRes
	res, err := syncGithubBranchReq.Run(cmd.Context(), githubCmdState.GithubV3, rootCfg.git)

	switch rootCfg.format {
	case "json":
		b, err1 := res.ToJSON()
		if err1 != nil {
			return errors.Join(err, err1)
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.ToTable(err).Render())
	}

	return err
}
