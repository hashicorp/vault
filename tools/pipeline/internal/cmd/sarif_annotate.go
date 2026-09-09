// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/sarif"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var annotateReq = &sarif.AnnotateReq{}

func newSarifAnnotateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "annotate [sarif-file] [key:value]...",
		Short: "Annotate SARIF file with metadata",
		Long: `Annotate SARIF run-level properties with user defined metadata.

This command adds metadata to a SARIF file's run-level properties using
key:value pairs. The metadata applies to each run in the report.

Key:value pairs support dot notation for nested properties:
  - Flat: product:vault-enterprise version:1.21.2
  - Nested: product.name:vault-enterprise product.version:1.21.2

The annotation is added to each run's properties bag per SARIF 2.1.0
specification (§3.8 Property Bags, §3.14 run object).

Examples:
  # View summary table (default)
  pipeline sarif annotate sarif.json product:vault-enterprise version:1.21.2

  # Update file in-place
  pipeline sarif annotate sarif.json product:vault version:1.21.2 --write

  # Write to new file
  pipeline sarif annotate sarif.json product:vault version:1.21.2 --out annotated.json

  # Output annotated JSON to stdout
  pipeline sarif annotate sarif.json product:vault version:1.21.2 --format json

  # Nested properties with dot notation
  pipeline sarif annotate sarif.json product.name:vault product.version:1.21.2 environment:prod`,
		Args: cobra.MinimumNArgs(2),
		RunE: runSarifAnnotateCmd,
	}

	cmd.Flags().BoolVarP(&annotateReq.Write, "write", "w", false, "Update file in-place")
	cmd.Flags().StringVarP(&annotateReq.OutputPath, "out", "o", "", "Output file path")

	cmd.MarkFlagsMutuallyExclusive("write", "out")

	return cmd
}

func runSarifAnnotateCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	annotateReq.SarifPath = args[0]
	annotateReq.KeyValuePairs = args[1:]

	res, err := annotateReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("annotating SARIF file: %w", err)
	}

	// Handle writing to our output files
	var output []byte
	if res.OutputPath != "" || annotateReq.Write {
		output, err = res.ToJSON()
		if err != nil {
			return err
		}

		writeFile := func(path string) error {
			mode := fs.FileMode(0o644)
			apath, err := filepath.Abs(path)
			if err == nil {
				stat, err := os.Stat(apath)
				if err == nil {
					mode = stat.Mode()
				}
			}

			return os.WriteFile(path, output, mode)
		}

		if res.OutputPath != "" {
			err = writeFile(res.OutputPath)
			if err != nil {
				return fmt.Errorf("writing to output path: %w", err)
			}
		}

		if annotateReq.Write {
			err = writeFile(annotateReq.SarifPath)
			if err != nil {
				return fmt.Errorf("writing to origin sarif: %w", err)
			}
		}
	}

	switch rootCfg.format {
	case "json":
		if len(output) > 0 {
			// We already parsed our output as JSON for our outpath or in-place annotation
			break
		}
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
		return err
	}
	_, err = os.Stdout.Write(output)
	if err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
