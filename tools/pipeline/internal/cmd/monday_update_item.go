// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/monday"
	"github.com/spf13/cobra"
)

func newMondayUpdateItemCmd() *cobra.Command {
	req := &monday.UpdateItemReq{}
	var githubOutput bool

	cmd := &cobra.Command{
		Use:   "update-item --board-id <id> --item-id <id> --column-id <id> --value <value>",
		Short: "Change the value of a column on a Monday.com board",
		Long:  "Change the value of a column on a Monday.com board",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := cmd.Context()

			result, err := req.Run(ctx, mondayCmdState.client)
			if err != nil {
				return fmt.Errorf("updating Monday board column: %w", err)
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
	cmd.Flags().StringVar(&req.BoardID, "board-id", "", "Monday.com id for the board we're modifying")
	cmd.Flags().StringVar(&req.ItemID, "item-id", "", "Item id for the item within the board we're updating")
	cmd.Flags().StringVar(&req.ColumnID, "column-id", "", "The column within the item upon which we're making the update")
	cmd.Flags().StringVar(&req.Value, "value", "", "The actual value that the column will take on after update")

	_ = cmd.MarkFlagRequired("board-id")
	_ = cmd.MarkFlagRequired("item-id")
	_ = cmd.MarkFlagRequired("column-id")
	_ = cmd.MarkFlagRequired("value")
	cmd.Flags().BoolVar(&githubOutput, "github-output", false, "Write updated item's item-id to $GITHUB_OUTPUT")

	return cmd
}
