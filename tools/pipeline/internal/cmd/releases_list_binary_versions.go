// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var listBinaryVersionsReq = &releases.ListBinaryVersionsReq{}

func newReleasesListBinaryVersionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "binary-versions <versions>",
		Short: "List available binary release variants for the given version labels",
		Long:  "List available binary release variants for the given version labels",
		RunE:  runListBinaryVersions,
		Args:  cobra.MinimumNArgs(1), // Require at least the versions argument
	}

	// Allows the user to control whether table or JSON is printed to stdout.
	cmd.PersistentFlags().StringVar(&rootCfg.format, "format", "table", `Output format: table|json`)
	return cmd
}

func runListBinaryVersions(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	listBinaryVersionsReq.VersionsString = args[0]

	// Optional second argument: boolean controlling whether JSON is saved to a file.
	saveToFile := false
	if len(args) > 1 {
		v, err := strconv.ParseBool(args[1])
		if err != nil {
			return fmt.Errorf("second argument must be true or false: %w", err)
		}
		saveToFile = v
	}

	// Executes the backend query to retrieve binary version metadata.
	res, err := listBinaryVersionsReq.Run(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to list binary versions: %w", err)
	}

	// Pre-encode JSON once since we may print it or save it depending on flags/args.
	jsonBytes, err := res.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	// If the user requested saving the response, write the JSON output to a local file.
	// This occurs regardless of the --format flag.
	if saveToFile {
		fileName := "binary-versions-output.json"
		if err := os.WriteFile(fileName, jsonBytes, 0o644); err != nil {
			return fmt.Errorf("failed to write JSON to file: %w", err)
		}
		fmt.Printf("Saved JSON output to %s\n", fileName)
	}

	switch rootCfg.format {
	case "json":
		fmt.Println(string(jsonBytes))
	default:
		printBinaryTable(res)
	}

	return nil
}

func printBinaryTable(res *releases.ListBinaryVersionsRes) {
	// Pretty-print table writer for terminal output.
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Version", "Status", "Variants Count", "Variants", "OS"})

	// Iterate through all versions requested, including ones the backend marked missing.
	for _, label := range res.AllVersions {
		entry, ok := res.ValidVersions[label]

		// Default values for versions not found.
		status := "MISSING"
		count := "—"
		variantsStr := ""
		osStr := ""

		if ok {
			// Version exists — fill in detailed information.
			status = strings.ToUpper(entry.Status)
			count = fmt.Sprintf("%d", len(entry.Variants))

			// Collect variant labels and all OS values.
			var variantNames []string
			osSet := make(map[string]struct{}) // Deduplicate OS entries across variants.

			for _, v := range entry.Variants {
				variantNames = append(variantNames, v.Variant)
				for _, osName := range v.OS {
					osSet[osName] = struct{}{}
				}
			}

			// Join variant names like: "enterprise, oss, fips"
			variantsStr = strings.Join(variantNames, ", ")

			// Convert deduped OS set into sorted list for stable output.
			var osList []string
			for osName := range osSet {
				osList = append(osList, osName)
			}
			slices.Sort(osList)
			osStr = strings.Join(osList, ", ")
		}

		// Add row representing this version label.
		t.AppendRow(table.Row{label, status, count, variantsStr, osStr})
	}

	t.Render()
}
