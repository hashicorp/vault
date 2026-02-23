// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	"github.com/spf13/cobra"
)

var validateConfigReq = &config.ValidateReq{}

func newConfigValidateCmd() *cobra.Command {
	validateCmd := &cobra.Command{
		Use:   "validate [--pipeline-config .release/pipeline.hcl]",
		Short: "Validate pipeline configuration",
		Long:  "Validate the pipeline.hcl configuration file for syntax and semantic errors",
		RunE:  runConfigValidateCmd,
		Args:  cobra.NoArgs,
	}

	return validateCmd
}

func runConfigValidateCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	validateConfigReq.DecodeRes = rootCfg.configDecodeRes
	res, err := validateConfigReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	switch rootCfg.format {
	case "json":
		b, err := res.ToJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		fmt.Println(res.String())
	}

	return nil
}
