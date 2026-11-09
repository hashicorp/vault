// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/monday"
	"github.com/spf13/cobra"
)

func newMondayCreateItemCmd() *cobra.Command {
	req := &monday.CreateItemReq{}
	var githubOutput bool

	cmd := &cobra.Command{
		Use:   "create-item --board-id <id> --item-name <name> [--group-id <id>] [--status <label>]",
		Short: "Create a new item on a Monday.com board",
		Long:  "Create a new top-level item on a Monday.com board",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			ctx := cmd.Context()

			result, err := req.Run(ctx, mondayCmdState.client)
			if err != nil {
				return fmt.Errorf("creating Monday.com item: %w", err)
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
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&req.BoardID, "board-id", "", "Monday.com id for the board that receives the new item")
	cmd.Flags().StringVar(&req.ItemName, "item-name", "", "The name of the item to create")
	cmd.Flags().StringVar(&req.GroupID, "group-id", "", "Id of the group that receives the item. Defaults to the board's top group")
	cmd.Flags().StringVar(&req.Status, "status", "", "Initial status label of the item, e.g. \"Not started\". Case-sensitive. When empty, the status stays blank")
	cmd.Flags().BoolVar(&githubOutput, "github-output", false, "Write the new item's item-id to $GITHUB_OUTPUT")

	_ = cmd.MarkFlagRequired("board-id")
	_ = cmd.MarkFlagRequired("item-name")

	return cmd
}
