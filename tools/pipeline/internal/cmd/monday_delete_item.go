// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/monday"
	"github.com/spf13/cobra"
)

func newMondayDeleteItemCmd() *cobra.Command {
	req := &monday.DeleteItemReq{}

	cmd := &cobra.Command{
		Use:   "delete-item --board-id <id> --item-id <id> (--expect-name <name> | --dry-run)",
		Short: "Delete an item or subitem from a Monday.com board",
		Long: `Delete an item or subitem from a Monday.com board.

The item is selected by --board-id and --item-id. Names are not unique on a
board, so the name is never used to select. To delete a subitem, pass the
subitem board's id as --board-id.

This command is destructive and has no interactive prompt. Before deleting, it
reads the item and refuses to delete it if the item is not on --board-id or
its name does not equal --expect-name. --expect-name is required unless
--dry-run is set. Deleting an item also deletes its subitems.

--dry-run is read-only. It prints the item id, name, board id, group, parent
item (for a subitem) and subitem count. If --expect-name is also set, the dry
run checks it too.

Run --dry-run first, then copy the name into --expect-name:

  pipeline monday delete-item --board-id 123 --item-id 456 --dry-run
  pipeline monday delete-item --board-id 123 --item-id 456 --expect-name "Test item"`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := cmd.Context()

			result, err := req.Run(ctx, mondayCmdState.client)
			if err != nil {
				return fmt.Errorf("deleting Monday.com item: %w", err)
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
	cmd.Flags().StringVar(&req.ItemID, "item-id", "", "Monday.com id of the item or subitem to delete")
	cmd.Flags().StringVar(&req.ExpectName, "expect-name", "", "The exact name the item must have. Required unless --dry-run is set")
	cmd.Flags().BoolVar(&req.DryRun, "dry-run", false, "Print the item's details without deleting it")

	_ = cmd.MarkFlagRequired("board-id")
	_ = cmd.MarkFlagRequired("item-id")

	return cmd
}
