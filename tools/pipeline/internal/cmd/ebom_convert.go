// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/ebom"
	"github.com/spf13/cobra"
)

var convertEBOMToCSVReq = &ebom.ConvertEBOMToCSVReq{}

func newEBOMConvertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert <xlsx-file>",
		Short: "Convert an IBM EBOM .xlsx file to CSV",
		Long: `Convert an IBM EBOM .xlsx to per-sheet CSV files.

Convert IBM EBOMS saved as .xlsx into CSV. If the EBOM is saved as .xls it
will need to be converted to .xlsx via Excel or LibreOffice first.

CSV is written per eAssembly sheet. For single-sheet EBOMS spreadsheets, the
output is <product-id>.csv. For multi-sheet EBOMS it is <product-id>.<slug>.csv
where the slug is inferred from the eAssembly description.

Use --slugs to override inferred slugs for specific sheets by 0-based index:

Examples:
  pipeline ebom convert ./5900-BJF_1.19.10.xlsx
  pipeline ebom convert ./file.xlsx --pid 5900-BJ8
  pipeline ebom convert ./file.xlsx --slugs 0:essentials --slugs 1:standard
  pipeline ebom convert ./file.xlsx --out-dir /tmp/eboms`,
		Args: cobra.ExactArgs(1),
		RunE: runEBOMConvertCmd,
	}

	cmd.Flags().StringVarP(&convertEBOMToCSVReq.OutDir, "out-dir", "o", "",
		"Output directory (default: <git-root>/.release/ibm-pao/eboms, or '.' if not in a git repo)")
	cmd.Flags().StringVar(&convertEBOMToCSVReq.PIDOverride, "pid", "",
		"Override the PID in the output filename")
	cmd.Flags().StringArrayVar(&convertEBOMToCSVReq.Slugs, "slugs", nil,
		"Per-sheet slug overrides in N:slug format (0-based index, repeatable)")

	return cmd
}

func runEBOMConvertCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	convertEBOMToCSVReq.XLSXPath = args[0]
	convertEBOMToCSVReq.RepoRoot = rootCfg.repoRoot

	res, err := convertEBOMToCSVReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("convert EBOM to CSV: %w", err)
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
