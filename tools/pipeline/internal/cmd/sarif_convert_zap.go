// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/sarif"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var convertZapReq = &sarif.ConvertZapReq{}

func newSarifConvertZapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-zap [sarif-file]",
		Short: "Convert SARIF to ZAP format",
		Args:  cobra.ExactArgs(1),
		RunE:  runSarifConvertZapCmd,
		Long: `Convert SARIF file to ZAP format

This command reads a SARIF JSON file and converts it to ZAP (Zed Attack Proxy)
XML-like JSON format. The conversion groups findings by site and alert type,
creating a hierarchical structure suitable for security scanning tools.

Examples:
  # View as a human readable table (default)
  pipeline sarif convert-zap sarif.json

  # Convert to JSON (stdout)
  pipeline sarif convert-zap sarif.json --format json

  # Convert to JSON file
  pipeline sarif convert-zap sarif.json --format json --out zap-report.json

  # View as markdown table
  pipeline sarif convert-zap sarif.json --format markdown`,
	}

	cmd.Flags().StringVarP(&convertZapReq.OutputPath, "out", "o", "", "Output file path (default: stdout)")

	return cmd
}

func runSarifConvertZapCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	convertZapReq.SarifPath = args[0]

	res, err := convertZapReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("converting to ZAP format: %w", err)
	}

	var output []byte
	switch rootCfg.format {
	case "json":
		output, err = res.ToJSON()
	default:
		var t table.Writer
		t, err = res.ToTable(nil)
		if err == nil {
			if rootCfg.format == "markdown" {
				output = []byte(t.RenderMarkdown())
			} else {
				output = []byte(t.Render())
			}
		}
	}

	if err != nil {
		return fmt.Errorf("formatting output: %w", err)
	}

	// Write output
	if res.OutputPath != "" {
		err = os.WriteFile(res.OutputPath, output, 0o644)
	} else {
		_, err = os.Stdout.Write(output)
	}

	if err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
