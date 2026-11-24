// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/monday"
	"github.com/spf13/cobra"
)

func newMondayArchiveBoardCmd() *cobra.Command {
	req := &monday.ArchiveBoardReq{}

	cmd := &cobra.Command{
		Use:   "archive-board --board-id <id> [--dry-run] [--hard-delete]",
		Short: "Archive or delete a Monday.com board",
		Long: `Archive or delete a Monday.com board, for example a throwaway board made by
create-board or duplicate-board.

The board is selected by --board-id. Names are not unique (the same name can
exist in two workspaces), so the name is never used to select.

This command is destructive and has no interactive prompt.

By default the board is archived, which can be undone in the Monday.com UI.
--hard-delete deletes the board instead. It is the only signal that selects
the irreversible path.

--dry-run is read-only. It prints the board id, name, state, workspace, item
count and last-updated time. Boards that are already archived or deleted are
still found, with state "archived" or "deleted".

Run --dry-run first, then archive it:

  pipeline monday archive-board --board-id 123 --dry-run
  pipeline monday archive-board --board-id 123`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := cmd.Context()

			result, err := req.Run(ctx, mondayCmdState.client)
			if err != nil {
				return fmt.Errorf("archiving Monday.com board: %w", err)
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
	cmd.Flags().StringVar(&req.BoardID, "board-id", "", "Monday.com id of the board to archive or delete")
	cmd.Flags().BoolVar(&req.DryRun, "dry-run", false, "Print the board's details without changing it")
	cmd.Flags().BoolVar(&req.HardDelete, "hard-delete", false, "Delete the board instead of archiving it")

	_ = cmd.MarkFlagRequired("board-id")

	return cmd
}
