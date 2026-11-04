// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// ToReport formats the verification result as a human-readable text report.
func (r *VerifyBinariesRes) ToReport() string {
	if r == nil {
		return ""
	}

	var validVersionsOutput strings.Builder
	validVersionsOutput.WriteString("Versions:\n  ")
	for _, v := range r.ValidVersions {
		fmt.Fprintf(&validVersionsOutput, "✓ %s, ", v)
	}
	validVersionsOutput.WriteString("\n")

	if !r.HasMissing {
		var out strings.Builder
		out.WriteString("*Available Versions:*\n")
		out.WriteString(validVersionsOutput.String())
		return out.String()
	}

	var validVariantsOutput strings.Builder
	var validOSOutput strings.Builder

	for _, v := range r.ValidVersions {
		detail := r.Details[v]
		if detail == nil {
			continue
		}

		fmt.Fprintf(&validVariantsOutput, "Version: %s\n  ", v)
		for _, varName := range detail.ValidVariants {
			fmt.Fprintf(&validVariantsOutput, "✓ %s, ", varName)
		}
		validVariantsOutput.WriteString("\n")

		fmt.Fprintf(&validOSOutput, "Version: %s", v)
		var sortedVariants []string
		for varName := range detail.VariantsOS {
			sortedVariants = append(sortedVariants, varName)
		}
		slices.Sort(sortedVariants)

		for _, varName := range sortedVariants {
			osDetail := detail.VariantsOS[varName]
			fmt.Fprintf(&validOSOutput, "\n  Variant: %s\n  ", varName)
			if varName == detail.BaseVersion || varName == detail.BaseVersion+"+ent" {
				for _, osName := range osDetail.ValidOS {
					fmt.Fprintf(&validOSOutput, "  ✓ %s, ", osName)
				}
			} else {
				if slices.Contains(osDetail.ValidOS, "linux") {
					validOSOutput.WriteString("  ✓ linux")
				}
			}
		}
		validOSOutput.WriteString("\n")
	}

	var out strings.Builder
	out.WriteString("*Available Versions:*\n")
	out.WriteString(validVersionsOutput.String())

	out.WriteString("\n*Available Variants:*\n")
	out.WriteString(validVariantsOutput.String())

	out.WriteString("\n*Available OS:*\n")
	out.WriteString(validOSOutput.String())

	if len(r.MissingVersions) > 0 {
		out.WriteString("\n*Missing Release Binary Versions:*\n")
		for _, v := range r.MissingVersions {
			fmt.Fprintf(&out, "Missing Versions: %s\n", v)
		}
	}

	if len(r.MissingVariants) > 0 {
		out.WriteString("\n*Missing Release Binary Variants:*\n")
		for _, mv := range r.MissingVariants {
			fmt.Fprintf(&out, "%s\n", mv)
		}
	}

	if len(r.MissingOS) > 0 || len(r.InvalidOS) > 0 {
		out.WriteString("\n*Missing Release Binary OS:*\n")
		for _, mo := range r.MissingOS {
			fmt.Fprintf(&out, "%s\n", mo)
		}
		for _, io := range r.InvalidOS {
			fmt.Fprintf(&out, "%s\n", io)
		}
	}

	return out.String()
}

// ToTable renders a formatted table showing the verification summary.
func (r *VerifyBinariesRes) ToTable() string {
	if r == nil {
		return ""
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Version", "Status", "Variants Count", "Valid Variants", "Missing Variants", "Missing/Invalid OS"})

	versionsToRender := r.AllVersions
	if len(versionsToRender) == 0 {
		for v := range r.Details {
			versionsToRender = append(versionsToRender, v)
		}
		slices.Sort(versionsToRender)
	}

	for _, v := range versionsToRender {
		detail, ok := r.Details[v]
		if !ok || detail.Status == "missing" {
			t.AppendRow(table.Row{v, "MISSING", "—", "—", "All", "All"})
			continue
		}

		status := "VALID"
		missingVars := "—"
		if len(detail.MissingVariants) > 0 {
			status = "INCOMPLETE"
			missingVars = strings.Join(detail.MissingVariants, ", ")
		}

		// Sort variant names for deterministic table output
		var sortedVariants []string
		for varName := range detail.VariantsOS {
			sortedVariants = append(sortedVariants, varName)
		}
		slices.Sort(sortedVariants)

		var osIssues []string
		for _, varName := range sortedVariants {
			osDetail := detail.VariantsOS[varName]
			if len(osDetail.MissingOS) > 0 {
				osIssues = append(osIssues, fmt.Sprintf("%s missing [%s]", osDetail.Variant, strings.Join(osDetail.MissingOS, ",")))
			}
			if len(osDetail.InvalidOS) > 0 {
				osIssues = append(osIssues, fmt.Sprintf("%s invalid [%s]", osDetail.Variant, strings.Join(osDetail.InvalidOS, ",")))
			}
		}
		osStr := "—"
		if len(osIssues) > 0 {
			status = "INCOMPLETE"
			osStr = strings.Join(osIssues, "; ")
		}

		t.AppendRow(table.Row{
			v,
			status,
			fmt.Sprintf("%d", len(detail.ValidVariants)),
			strings.Join(detail.ValidVariants, ", "),
			missingVars,
			osStr,
		})
	}

	var sb strings.Builder
	sb.WriteString(t.Render())
	sb.WriteString("\n")

	if r.HasMissing {
		sb.WriteString("\nOmissions & Violations:\n")
		for _, v := range r.MissingVersions {
			fmt.Fprintf(&sb, "  • Missing Version: %s\n", v)
		}
		for _, mv := range r.MissingVariants {
			fmt.Fprintf(&sb, "  • %s\n", mv)
		}
		for _, mo := range r.MissingOS {
			fmt.Fprintf(&sb, "  • %s\n", mo)
		}
		for _, io := range r.InvalidOS {
			fmt.Fprintf(&sb, "  • %s\n", io)
		}
	}

	return sb.String()
}

// ToJSON returns pretty-printed JSON of the verification results.
func (r *VerifyBinariesRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}
	return json.MarshalIndent(r, "", "  ")
}

// ToGithubOutput formats the output for $GITHUB_OUTPUT.
func (r *VerifyBinariesRes) ToGithubOutput() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}
	return []byte(r.ToReport()), nil
}
