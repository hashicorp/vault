// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/monday"
	"github.com/spf13/cobra"
)

func newMondayDuplicateBoardCmd() *cobra.Command {
	req := &monday.DuplicateBoardReq{}
	var githubOutput bool

	cmd := &cobra.Command{
		Use:   "duplicate-board --template-id <id> --board-name <name> [--duplicate-type <type>]",
		Short: "Create a new Monday.com board by duplicating a template board",
		Long:  "Create a new Monday.com board by duplicating a template board",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			ctx := cmd.Context()

			result, err := req.Run(ctx, mondayCmdState.client)
			if err != nil {
				return fmt.Errorf("creating Monday.com board: %w", err)
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
				if err := writeToGithubOutput("board-id", []byte(result.BoardID)); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&req.TemplateID, "template-id", "", "The board ID of the Monday.com template to duplicate")
	cmd.Flags().StringVar(&req.BoardName, "board-name", "", "The name to assign to the new board")
	cmd.Flags().StringVar((*string)(&req.DuplicateType), "duplicate-type", string(monday.DuplicateBoardWithStructure), "The duplication rule to apply (duplicate_board_with_structure, duplicate_board_with_pulses, duplicate_board_with_pulses_and_updates)")
	cmd.Flags().BoolVar(&githubOutput, "github-output", false, "Write the new board-id to $GITHUB_OUTPUT")

	_ = cmd.MarkFlagRequired("template-id")
	_ = cmd.MarkFlagRequired("board-name")

	return cmd
}
