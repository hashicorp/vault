// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/monday"
	"github.com/spf13/cobra"
)

func newMondayGetItemCmd() *cobra.Command {
	req := &monday.GetItemReq{}

	cmd := &cobra.Command{
		Use:   "get-item --board-id <id> --item-id <id>",
		Short: "Get the details of an item or subitem on a Monday.com board",
		Long: `Get the details of an item or subitem on a Monday.com board.

This command is read-only. It prints the item id, name and board id, the group
(id and title), the state, the parent item id (for a subitem), each column
value (id, title, type and text), the subitem count and ids (for a parent
item), and the last-updated time. Use it to check an item's status before
calling update-item.

The item is selected by --board-id and --item-id. Monday.com does not scope
item lookups to a board, so the command exits non-zero if the item is not on
--board-id or does not exist. To read a subitem, pass the subitem board's id as
--board-id, the same as for update-item.

The table output shows each column's text. --format json also includes each
column's raw value.

  pipeline monday get-item --board-id 123 --item-id 456
  pipeline monday get-item --board-id 123 --item-id 456 --format json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := cmd.Context()

			result, err := req.Run(ctx, mondayCmdState.client)
			if err != nil {
				return fmt.Errorf("getting Monday.com item: %w", err)
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

			return nil
		},
	}
	cmd.Flags().StringVar(&req.BoardID, "board-id", "", "Monday.com id of the board the item is on. For a subitem, the subitem board's id")
	cmd.Flags().StringVar(&req.ItemID, "item-id", "", "Monday.com id of the item or subitem to read")

	_ = cmd.MarkFlagRequired("board-id")
	_ = cmd.MarkFlagRequired("item-id")

	return cmd
}
