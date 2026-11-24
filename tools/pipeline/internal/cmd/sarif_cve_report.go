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

var cveReportReq = &sarif.CVEReportReq{}

func newSarifCVEReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cve-report [sarif-file]",
		Short: "Create a CVE report from a SARIF file",
		Args:  cobra.ExactArgs(1),
		RunE:  runSarifCVEReportCmd,
		Long: `Create CVE report from a SARIF file

This command reads a SARIF (Static Analysis Results Interchange Format) file,
identifies findings linked to CVE (Common Vulnerabilities and Exposures) entries,
and generates a report in the specified format. This is intended to be used
on the results of the DAST/Zap scans of products so the CSV format that we
write adheres to a specific schema.

Examples:
  # Create table report to stdout
  pipeline sarif cve-report sarif.json

  # Create CSV report to file
  pipeline sarif cve-report sarif.json --format csv --out report.csv

  # Create JSON report
  pipeline sarif cve-report sarif.json --format json`,
	}

	cmd.Flags().StringVarP(&cveReportReq.OutputPath, "out", "o", "", "Output file path (default: stdout)")

	return cmd
}

func runSarifCVEReportCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	cveReportReq.SarifPath = args[0]

	res, err := cveReportReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("generating CVE report: %w", err)
	}

	// Format output based on rootCfg.format
	var output []byte
	switch rootCfg.format {
	case "json":
		output, err = res.ToJSON()
	case "csv":
		output, err = res.ToCSV()
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
