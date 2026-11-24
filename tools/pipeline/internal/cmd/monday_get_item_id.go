// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/monday"
	"github.com/spf13/cobra"
)

func newMondayGetItemIDCmd() *cobra.Command {
	req := &monday.GetItemIDReq{}
	var githubOutput bool

	cmd := &cobra.Command{
		Use:   "get-item-id (--board-id <id> | --parent-id <id>) --item-name <name>",
		Short: "Get the item id from the monday board by item-name, use parent-id instead of board-id when getting subitem",
		Long:  "Get the item id from the monday board by item-name, use parent-id instead of board-id when getting subitem",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := cmd.Context()

			result, err := req.Run(ctx, mondayCmdState.client)
			if err != nil {
				return fmt.Errorf("getting Monday item id: %w", err)
			}

			switch rootCfg.format {
			case "json":
				b, err := result.ToJSON()
				if err != nil {
					return err
				}
				fmt.Println(string(b))
			default:
				fmt.Println(result.ToTable().Render())
			}

			if githubOutput {
				if err := writeToGithubOutput("item-id", []byte(result.ItemID)); err != nil {
					return err
				}
				if err := writeToGithubOutput("board-id", []byte(result.BoardID)); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&req.BoardID, "board-id", "", "Monday.com id for the board that contains the item.")
	cmd.Flags().StringVar(&req.ItemName, "item-name", "", "Name of the Monday.com item for which we want the id.")
	cmd.Flags().BoolVar(&githubOutput, "github-output", false, "Write board item's item-id and board-id to $GITHUB_OUTPUT")
	cmd.Flags().StringVar(&req.ParentID, "parent-id", "", "Monday.com id of the parent item whose subitems are searched. Cannot be used with --board-id flag.")
	_ = cmd.MarkFlagRequired("item-name")
	cmd.MarkFlagsOneRequired("board-id", "parent-id")
	cmd.MarkFlagsMutuallyExclusive("board-id", "parent-id")

	return cmd
}
