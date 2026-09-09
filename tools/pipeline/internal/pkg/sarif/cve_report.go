// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package sarif

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	gosarif "github.com/owenrumney/go-sarif/v3/pkg/report/v22/sarif"
)

// CVEReportReq is the request to generate a CVE report
type CVEReportReq struct {
	SarifPath  string // Path to SARIF file
	OutputPath string // Output file path (empty for stdout)
}

// CVEReportRes is the response for generating a CVE report
type CVEReportRes struct {
	OutputPath string // Output file path
	CVEs       []*CVE // Found CVEs in the report
}

// CVE represents a CVE-linked finding in a SARIF report
type CVE struct {
	Vulnerability         string `json:"vulnerability"`
	PackageLibraryName    string `json:"packageLibraryName"`
	PackageLibraryVersion string `json:"packageLibraryVersion"`
	RuleID                string `json:"ruleId,omitempty"`
	Message               string `json:"message,omitempty"`
	URI                   string `json:"uri,omitempty"`
}

type ruleInfo struct {
	name   string
	hasCVE bool
}

// Run executes the CVE report generation
func (r *CVEReportReq) Run(ctx context.Context) (*CVEReportRes, error) {
	report, err := gosarif.Open(r.SarifPath)
	if err != nil {
		return nil, fmt.Errorf("opening SARIF file: %w", err)
	}

	return &CVEReportRes{
		OutputPath: r.OutputPath,
		CVEs:       sarifCVEs(report),
	}, nil
}

// ToJSON marshals the response to JSON
func (r *CVEReportRes) ToJSON() ([]byte, error) {
	b, err := json.MarshalIndent(r.CVEs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling CVE report to JSON: %w", err)
	}
	return b, nil
}

// ToCSV marshals the response to CSV format
func (r *CVEReportRes) ToCSV() ([]byte, error) {
	if len(r.CVEs) == 0 {
		return []byte(""), nil
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"vulnerability", "packageLibraryName", "packageLibraryVersion"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("writing CSV header: %w", err)
	}

	for _, f := range r.CVEs {
		row := []string{f.Vulnerability, f.PackageLibraryName, f.PackageLibraryVersion}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("writing CSV entry: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("flushing CSV writer: %w", err)
	}

	return buf.Bytes(), nil
}

// ToTable marshals the response to a text table
func (r *CVEReportRes) ToTable(err error) (table.Writer, error) {
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

	if r == nil || len(r.CVEs) == 0 {
		t.AppendHeader(table.Row{"message"})
		t.AppendRow(table.Row{"No CVE-linked findings in report."})

		return t, err
	}

	t.AppendHeader(table.Row{"VULNERABILITY", "PACKAGE", "VERSION"})
	for _, f := range r.CVEs {
		pkg := f.PackageLibraryName
		if pkg == "" {
			pkg = "-"
		}

		ver := f.PackageLibraryVersion
		if ver == "" {
			ver = "-"
		}

		t.AppendRow(table.Row{f.Vulnerability, pkg, ver})
	}

	return t, nil
}

func sarifCVEs(report *gosarif.Report) []*CVE {
	// Build rule map with CVE status
	ruleMap := make(map[string]*ruleInfo)

	for _, run := range report.Runs {
		if run.Tool.Driver == nil || run.Tool.Driver.Rules == nil {
			continue
		}

		for _, rule := range run.Tool.Driver.Rules {
			hasCVE := false

			// Check if rule has CVE relationship
			if rule.Relationships != nil {
				for _, rel := range rule.Relationships {
					if rel.Target != nil &&
						rel.Target.ToolComponent != nil &&
						rel.Target.ToolComponent.Name != nil &&
						strings.ToUpper(*rel.Target.ToolComponent.Name) == "CVE" {
						hasCVE = true
						break
					}
				}
			}

			// Determine rule name (priority: Name > ShortDescription > ID)
			name := ""
			if rule.Name != nil {
				name = *rule.Name
			} else if rule.ShortDescription != nil && rule.ShortDescription.Text != nil {
				name = *rule.ShortDescription.Text
			} else if rule.ID != nil {
				name = *rule.ID
			}

			if rule.ID != nil {
				ruleMap[*rule.ID] = &ruleInfo{
					name:   name,
					hasCVE: hasCVE,
				}
			}
		}
	}

	cves := []*CVE{}
	seen := make(map[string]bool)

	for _, run := range report.Runs {
		if run.Results == nil {
			continue
		}

		for _, result := range run.Results {
			if result.RuleID == nil {
				continue
			}

			ruleInfo, ok := ruleMap[*result.RuleID]
			if !ok || !ruleInfo.hasCVE {
				continue
			}

			vulnName := ruleInfo.name
			// Deduplicate by vulnerability name
			if seen[vulnName] {
				continue
			}
			seen[vulnName] = true

			cve := &CVE{
				Vulnerability:         vulnName,
				PackageLibraryName:    "", // For now we leave this blank but it's there for the CSV interface
				PackageLibraryVersion: "", // For now we leave this blank but it's there for the CSV interface
				RuleID:                *result.RuleID,
			}

			// Extract message
			if result.Message.Text != nil {
				cve.Message = *result.Message.Text
			}

			// Extract URI if available
			if result.Locations != nil {
				loc := result.Locations[0]
				if loc.PhysicalLocation != nil &&
					loc.PhysicalLocation.ArtifactLocation != nil &&
					loc.PhysicalLocation.ArtifactLocation.URI != nil {
					cve.URI = *loc.PhysicalLocation.ArtifactLocation.URI
				}
			}

			cves = append(cves, cve)
		}
	}

	return cves
}
