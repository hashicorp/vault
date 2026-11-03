// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/monday"
	"github.com/spf13/cobra"
)

func newMondayCreateSubitemCmd() *cobra.Command {
	req := &monday.CreateSubitemReq{}
	var githubOutput bool

	cmd := &cobra.Command{
		Use:   "create-subitem --parent-item-id <id> --item-name <name>",
		Short: "Create a new sub-item under an existing Monday.com item",
		Long:  "Create a new sub-item under an existing Monday.com item",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := cmd.Context()

			result, err := req.Run(ctx, mondayCmdState.client)
			if err != nil {
				return fmt.Errorf("creating Monday.com subitem: %w", err)
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
				if err := writeToGithubOutput("subitem-id", []byte(result.SubitemID)); err != nil {
					return err
				}
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&req.ParentItemID, "parent-item-id", "", "Monday.com id for the item that receives the new sub-item")
	cmd.Flags().StringVar(&req.ItemName, "item-name", "", "The name of the sub-item to create")
	cmd.Flags().BoolVar(&githubOutput, "github-output", false, "Write the new sub-item's subitem-id to $GITHUB_OUTPUT")

	_ = cmd.MarkFlagRequired("parent-item-id")
	_ = cmd.MarkFlagRequired("item-name")

	return cmd
}
