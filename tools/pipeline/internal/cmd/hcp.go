// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/hcp"
	"github.com/spf13/cobra"
)

var hcpCmdState = &struct {
	client *hcp.Client
}{}

func newHCPCmd() *cobra.Command {
	env := ""
	hcpCmd := &cobra.Command{
		Use:   "hcp",
		Short: "HCP commands",
		Long:  "HCP commands",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			hcpCmdState.client = hcp.NewClient(
				hcp.WithEnvironment(hcp.Environment(env)),
				hcp.WithLoadAuthFromEnv(),
			)
			if _, set := os.LookupEnv("HCP_PASSWORD"); !set {
				fmt.Println("\x1b[1;33;49mWARNING\x1b[0m: HCP_PASSWORD has not been set. You probably want to set it and HCP_USERNAME in order to authenticate with the image service")
			}
		},
	}
	hcpCmd.AddCommand(newHCPShowCmd())
	hcpCmd.AddCommand(newHCPWaitCmd())

	hcpCmd.PersistentFlags().StringVarP(&env, "environment", "e", "prod", "The HCP environment to use. E.g. dev, int, prod")

	return hcpCmd
}
