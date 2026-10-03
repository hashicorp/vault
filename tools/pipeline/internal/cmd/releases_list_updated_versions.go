// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
)

var listReleaseUpdatedVersionsReq = &releases.ListUpdatedVersionsReq{}

func newReleasesListUpdatedVersionsCmd() *cobra.Command {
	updatedVersionsCmd := &cobra.Command{
		Use:   "updated-versions \"input_versions\"",
		Short: "String of input versions separated by spaces",
		Long:  "String of input versions separated by spaces",
		RunE:  runListUpdatedVersionsReq,
		Args:  cobra.MaximumNArgs(1), // s
	}

	return updatedVersionsCmd
}

func runListUpdatedVersionsReq(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := listReleaseUpdatedVersionsReq.Run(context.TODO(), args)
	if err != nil {
		return err
	}

	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}
