// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/monday"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/spf13/cobra"
)

var mondayCmdState = &struct {
	client *graphql.Client
}{}

func newMondayCmd() *cobra.Command {
	mondayCmd := &cobra.Command{
		Use:   "monday",
		Short: "Monday.com board interactions",
		Long:  "Monday.com board interactions",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			token, ok := os.LookupEnv("MONDAY_TOKEN")

			if !ok || token == "" {
				return fmt.Errorf("MONDAY_TOKEN env variable isn't set")
			}

			mondayCmdState.client = monday.NewClient(token)
			return nil
		},
	}

	mondayCmd.AddCommand(newMondayCreateBoardCmd())
	mondayCmd.AddCommand(newMondayDuplicateBoardCmd())
	mondayCmd.AddCommand(newMondayUpdateItemCmd())
	mondayCmd.AddCommand(newMondayGetItemIDCmd())

	return mondayCmd
}
