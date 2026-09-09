// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package sarif

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	gosarif "github.com/owenrumney/go-sarif/v3/pkg/report/v22/sarif"
)

// AnnotateReq is the request to annotate a SARIF file
type AnnotateReq struct {
	SarifPath     string   // Path to SARIF file
	KeyValuePairs []string // Key:value pairs to add as properties
	Write         bool     // Update file in-place
	OutputPath    string   // Output file path (empty for stdout)
}

// AnnotateRes is the response for SARIF annotation
type AnnotateRes struct {
	SarifPath  string          // Original SARIF file path
	OutputPath string          // Output file path (if different)
	Updated    bool            // Whether file was updated
	Properties map[string]any  // Properties that were added
	Report     *gosarif.Report // Annotated SARIF report
}

// Run executes the SARIF annotation
func (r *AnnotateReq) Run(ctx context.Context) (*AnnotateRes, error) {
	// Validate inputs
	if len(r.KeyValuePairs) == 0 {
		return nil, fmt.Errorf("at least one metadata key:value pair is required")
	}

	properties, err := parsePropertiesFromKeyValuePairs(r.KeyValuePairs)
	if err != nil {
		return nil, fmt.Errorf("parsing key:value pairs: %w", err)
	}

	report, err := gosarif.Open(r.SarifPath)
	if err != nil {
		return nil, fmt.Errorf("opening SARIF file: %w", err)
	}

	err = annotateSarif(report, properties)
	if err != nil {
		return nil, fmt.Errorf("annotating SARIF: %w", err)
	}

	outputPath := r.OutputPath
	if r.Write {
		outputPath = r.SarifPath
	}

	return &AnnotateRes{
		SarifPath:  r.SarifPath,
		OutputPath: outputPath,
		Updated:    r.Write || r.OutputPath != "",
		Properties: properties,
		Report:     report,
	}, nil
}

// ToJSON marshals the annotated SARIF report to JSON
func (r *AnnotateRes) ToJSON() ([]byte, error) {
	b, err := json.MarshalIndent(r.Report, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling SARIF report to JSON: %w", err)
	}

	return append(b, '\n'), nil
}

// ToTable marshals the response to a text table
func (r *AnnotateRes) ToTable(err error) (table.Writer, error) {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	defer t.SuppressEmptyColumns()
	defer t.SuppressTrailingSpaces()

	if err != nil {
		t.AppendHeader(table.Row{"error"})
		t.AppendRow(table.Row{err.Error()})
		return t, err
	}

	if r == nil {
		t.AppendHeader(table.Row{"message"})
		t.AppendRow(table.Row{"No annotation performed."})
		return t, nil
	}

	t.AppendHeader(table.Row{"KEY", "VALUE"})

	// Add driver information if available
	if r.Report != nil && len(r.Report.Runs) > 0 && r.Report.Runs[0].Tool != nil && r.Report.Runs[0].Tool.Driver != nil {
		driver := r.Report.Runs[0].Tool.Driver
		if driver.Name != nil && *driver.Name != "" {
			t.AppendRow(table.Row{"Tool", *driver.Name})
		}
		if driver.Version != nil && *driver.Version != "" {
			t.AppendRow(table.Row{"Tool Version", *driver.Version})
		}
		if driver.InformationURI != nil && *driver.InformationURI != "" {
			t.AppendRow(table.Row{"Tool URI", *driver.InformationURI})
		}
		t.AppendSeparator()
	}

	// Add each property as a row
	for key, value := range r.Properties {
		valueStr := formatPropertyValue(value)
		t.AppendRow(table.Row{key, valueStr})
	}

	// Add summary row
	t.AppendSeparator()
	if r.Report != nil {
		t.AppendRow(table.Row{"Runs annotated", fmt.Sprintf("%d", len(r.Report.Runs))})
	}

	if r.Updated {
		if r.OutputPath == r.SarifPath {
			t.AppendRow(table.Row{"Status", "Updated in-place"})
		} else {
			t.AppendRow(table.Row{"Output", r.OutputPath})
		}
	}

	return t, nil
}

// parsePropertiesFromKeyValuePairs parses command line key:value pairs into
// a property map. While it would be safe to assume kye:value pairs are always
// string:string, the SARIF library expects map[string]any so we use that type
// everywhere for props.
func parsePropertiesFromKeyValuePairs(args []string) (map[string]any, error) {
	properties := make(map[string]any)

	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid key:value format: %s", arg)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" || value == "" {
			return nil, fmt.Errorf("key and value cannot be empty: %s", arg)
		}

		// Handle dot notation for nested keys
		if strings.Contains(key, ".") {
			if err := setNestedProperty(properties, key, value); err != nil {
				return nil, err
			}
		} else {
			properties[key] = value
		}
	}

	return properties, nil
}

// setNestedProperty creates nested map structure from dot notation.
func setNestedProperty(props map[string]any, key, value string) error {
	keys := strings.Split(key, ".")

	for i := 0; i < len(keys)-1; i++ {
		k := keys[i]
		if existing, exists := props[k]; exists {
			// If exists, it must be a map
			if m, ok := existing.(map[string]any); ok {
				props = m
			} else {
				return fmt.Errorf("cannot create nested property %s: %s is not a map", key, k)
			}
		} else {
			// Create new nested map
			newMap := make(map[string]any)
			props[k] = newMap
			props = newMap
		}
	}

	props[keys[len(keys)-1]] = value

	return nil
}

// formatPropertyValue formats a property value for display in table
func formatPropertyValue(value any) string {
	switch v := value.(type) {
	case map[string]any:
		var parts []string
		for key, val := range v {
			parts = append(parts, fmt.Sprintf("%s: %s", key, formatPropertyValue(val)))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// annotateSarif adds properties to all runs in the SARIF report
func annotateSarif(report *gosarif.Report, properties map[string]any) error {
	if report == nil {
		return fmt.Errorf("report is nil")
	}

	if len(report.Runs) == 0 {
		return fmt.Errorf("report has no runs")
	}

	for _, run := range report.Runs {
		if run.Properties == nil {
			run.Properties = &gosarif.PropertyBag{
				Properties: make(map[string]any),
			}
		}
		if run.Properties.Properties == nil {
			run.Properties.Properties = make(map[string]any)
		}

		maps.Copy(run.Properties.Properties, properties)
	}

	return nil
}
