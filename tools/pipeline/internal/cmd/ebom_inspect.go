// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/ebom"
	"github.com/spf13/cobra"
)

var inspectEBOMReq = &ebom.InspectEBOMReq{}

func newEBOMInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect <xlsx-file>",
		Short: "Inspect eAssembly sheets and inferred slugs in an IBM EBOM .xlsx file",
		Long: `Inspect an IBM EBOM .xlsx files.

Inspect IBM EBOMS saved as .xlsx. If the EBOM is saved as .xls it will need
to be converted to .xlsx via Excel or LibreOffice first.

Reports each eAssembly sheet with its 0-based index, PID, eAssembly
description, and inferred slug. You can use this before converting a
multi-page EBOM to verify the slugs. Then use the --slugs arguments
for covert to override any that need it.

Examples:
  pipeline ebom inspect ./5900-BJF_1.19.10.xlsx
  pipeline ebom inspect ./file.xlsx --pid 5900-BJ8 --format json`,
		Args: cobra.ExactArgs(1),
		RunE: runEBOMInspectCmd,
	}

	cmd.Flags().StringVar(&inspectEBOMReq.PIDOverride, "pid", "",
		"Override the PID shown in output")

	return cmd
}

func runEBOMInspectCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	inspectEBOMReq.XLSXPath = args[0]

	res, err := inspectEBOMReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("inspect EBOM: %w", err)
	}

	switch rootCfg.format {
	case "json":
		b, jsonErr := res.ToJSON()
		if jsonErr != nil {
			return jsonErr
		}
		fmt.Println(string(b))
	case "markdown":
		fmt.Println(res.ToMarkdown())
	default:
		fmt.Println(res.ToTable())
	}
	return nil
}
