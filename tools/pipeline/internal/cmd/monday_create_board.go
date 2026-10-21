// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/monday"
	"github.com/spf13/cobra"
)

func newMondayCreateBoardCmd() *cobra.Command {
	req := &monday.CreateBoardReq{}
	var githubOutput bool

	cmd := &cobra.Command{
		Use:   "create-board --board-name <name> [--board-kind <kind>]",
		Short: "Create a new Monday.com board",
		Long:  "Create a new Monday.com board",
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

	cmd.Flags().StringVar(&req.BoardName, "board-name", "", "The name to assign to the new board")
	cmd.Flags().StringVar((*string)(&req.BoardKind), "board-kind", string(monday.BoardKindPublic), "The visibility of the new board (public, private, share)")
	cmd.Flags().BoolVar(&githubOutput, "github-output", false, "Write the new board-id to $GITHUB_OUTPUT")

	_ = cmd.MarkFlagRequired("board-name")

	return cmd
}
