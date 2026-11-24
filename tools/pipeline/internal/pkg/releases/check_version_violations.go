// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// violationVersionPattern matches the normalized version format produced by
// sanitize-sort-versions: X.Y.Z or X.Y.Z-rcN (no +ent, no -lts).
var violationVersionPattern = regexp.MustCompile(
	`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-rc([1-9][0-9]*))?$`,
)

// CheckVersionViolationsReq holds the parameters for checking a sorted,
// normalized list of versions for rule violations.
type CheckVersionViolationsReq struct {
	// Versions is the sorted, normalized list to check (output of sanitize-sort-versions).
	Versions []string
	// WriteToGithubOutput controls whether the command layer writes to $GITHUB_OUTPUT.
	WriteToGithubOutput bool
}

// CheckVersionViolationsRes is the result of the violations check.
type CheckVersionViolationsRes struct {
	// HasViolations is true when at least one rule was violated.
	HasViolations bool `json:"has_violations"`
	// InvalidVersionFormat holds versions that did not match X.Y.Z or X.Y.Z-rcN.
	InvalidVersionFormat []string `json:"invalid_version_format_violations"`
	// MajorMinorDuplicates holds major.minor keys that appeared more than once.
	MajorMinorDuplicates []string `json:"major_minor_with_duplicates_violations"`
	// PatchZeroHits holds all versions with patch level 0 (there may only be one).
	PatchZeroHits []string `json:"patch_zero_hits"`
	// GaOrRcNotTopVersion holds GA or RC versions that were not the highest in the list.
	GaOrRcNotTopVersion []string `json:"ga_or_rc_not_top_version_violations"`
	// RCVersions holds all RC versions found (more than one is a violation).
	RCVersions []string `json:"rc_version_violations"`
	// GAVersions holds all GA (patch=0, no RC) versions found (more than one is a violation).
	GAVersions []string `json:"ga_version_violations"`
	// GaAndRcMixedInSet holds the merged GA+RC list when both appear together.
	GaAndRcMixedInSet []string `json:"ga_and_rc_mixed_in_set_violations"`
	// RCWithInvalidPatchLevel holds RC versions whose patch level is not 0.
	RCWithInvalidPatchLevel []string `json:"release_candidate_with_invalid_patch_level_violations"`
}

// Run checks the provided sorted, normalized versions for the seven violation rules.
func (r *CheckVersionViolationsReq) Run(ctx context.Context) (*CheckVersionViolationsRes, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	slog.Default().DebugContext(ctx, "starting check-version-violations",
		slog.Int("version_count", len(r.Versions)),
	)

	res := &CheckVersionViolationsRes{
		InvalidVersionFormat:    []string{},
		MajorMinorDuplicates:    []string{},
		PatchZeroHits:           []string{},
		GaOrRcNotTopVersion:     []string{},
		RCVersions:              []string{},
		GAVersions:              []string{},
		GaAndRcMixedInSet:       []string{},
		RCWithInvalidPatchLevel: []string{},
	}

	majorMinorMap := make(map[string]int)

	for _, version := range r.Versions {
		m := violationVersionPattern.FindStringSubmatch(version)
		if m == nil {
			slog.Default().DebugContext(ctx, "invalid version format", slog.String("version", version))
			res.InvalidVersionFormat = append(res.InvalidVersionFormat, version)
			continue
		}

		patch, _ := strconv.Atoi(m[3])
		rc := m[4] // e.g. "-rc1", or "" for GA/release

		majorMinorKey := fmt.Sprintf("%s.%s", m[1], m[2])

		// Rule 1: track major.minor to detect duplicates.
		majorMinorMap[majorMinorKey]++

		// Rule 2: patch level 0 tracking.
		if patch == 0 {
			res.PatchZeroHits = append(res.PatchZeroHits, version)
		}

		// Rule 3: GA (patch=0) or RC must be the highest version (first in the
		// descending-sorted list).
		if (patch == 0 || rc != "") && version != r.Versions[0] {
			res.GaOrRcNotTopVersion = append(res.GaOrRcNotTopVersion, version)
		}

		// Rule 4 & 7: classify into RC, GA, or plain release.
		if rc != "" {
			res.RCVersions = append(res.RCVersions, version)
		} else if patch == 0 {
			res.GAVersions = append(res.GAVersions, version)
		}

		// Rule 5: RC must have patch level 0.
		if rc != "" && patch != 0 {
			res.RCWithInvalidPatchLevel = append(res.RCWithInvalidPatchLevel, version)
		}
	}

	// Rule 1 (continued): collect major.minor keys with more than one version.
	for key, count := range majorMinorMap {
		if count > 1 {
			res.MajorMinorDuplicates = append(res.MajorMinorDuplicates, key)
		}
	}

	// Rule 4: GA and RC cannot coexist.
	if len(res.GAVersions) > 0 && len(res.RCVersions) > 0 {
		res.GaAndRcMixedInSet = append(res.GAVersions, res.RCVersions...)
	}

	res.HasViolations = len(res.InvalidVersionFormat) > 0 ||
		len(res.PatchZeroHits) > 1 ||
		len(res.GaOrRcNotTopVersion) > 0 ||
		len(res.MajorMinorDuplicates) > 0 ||
		len(res.RCVersions) > 1 ||
		len(res.GAVersions) > 1 ||
		len(res.GaAndRcMixedInSet) > 0 ||
		len(res.RCWithInvalidPatchLevel) > 0

	slog.Default().DebugContext(ctx, "completed check-version-violations",
		slog.Bool("has_violations", res.HasViolations),
	)

	return res, nil
}

// validate checks that the request is initialized.
func (r *CheckVersionViolationsReq) validate() error {
	if r == nil {
		return errors.New("uninitialized request")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *CheckVersionViolationsRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling check-version-violations response to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CheckVersionViolationsRes) ToTable() string {
	if r == nil {
		return ""
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{"violation", "versions"})
	appendViolation := func(name string, items []string) {
		if len(items) > 0 {
			t.AppendRow(table.Row{name, strings.Join(items, ", ")})
		}
	}
	appendViolation("invalid_version_format", r.InvalidVersionFormat)
	appendViolation("major_minor_duplicates", r.MajorMinorDuplicates)
	appendViolation("patch_zero_hits", r.PatchZeroHits)
	appendViolation("ga_or_rc_not_top_version", r.GaOrRcNotTopVersion)
	appendViolation("multiple_rc_versions", r.RCVersions)
	appendViolation("multiple_ga_versions", r.GAVersions)
	appendViolation("ga_and_rc_mixed", r.GaAndRcMixedInSet)
	appendViolation("rc_invalid_patch_level", r.RCWithInvalidPatchLevel)

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t.Render()
}

// ToGithubOutput returns the outputs formatted for $GITHUB_OUTPUT, matching
// the key names the compile-version-rule-violations workflow step expects.
func (r *CheckVersionViolationsRes) ToGithubOutput() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	join := func(items []string) string { return strings.Join(items, ",") }

	out := strings.Builder{}
	out.WriteString(fmt.Sprintf("has_violations=%t\n", r.HasViolations))
	out.WriteString(fmt.Sprintf("invalid_version_format_violations=%s\n", join(r.InvalidVersionFormat)))
	out.WriteString(fmt.Sprintf("patch_zero_hits=%s\n", join(r.PatchZeroHits)))
	out.WriteString(fmt.Sprintf("ga_or_rc_not_top_version_violations=%s\n", join(r.GaOrRcNotTopVersion)))
	out.WriteString(fmt.Sprintf("major_minor_with_duplicates_violations=%s\n", join(r.MajorMinorDuplicates)))
	out.WriteString(fmt.Sprintf("ga_and_rc_mixed_in_set_violations=%s\n", join(r.GaAndRcMixedInSet)))
	out.WriteString(fmt.Sprintf("release_candidate_with_invalid_patch_level_violations=%s\n", join(r.RCWithInvalidPatchLevel)))
	out.WriteString(fmt.Sprintf("rc_version_violations=%s\n", join(r.RCVersions)))
	out.WriteString(fmt.Sprintf("ga_version_violations=%s\n", join(r.GAVersions)))

	return []byte(out.String()), nil
}

// String returns a human-readable summary.
func (r *CheckVersionViolationsRes) String() string {
	if r == nil {
		return ""
	}

	if !r.HasViolations {
		return "No version violations found.\n"
	}

	w := strings.Builder{}
	w.WriteString("Version violations found:\n")
	report := func(label string, items []string) {
		if len(items) > 0 {
			w.WriteString(fmt.Sprintf("  %s: %s\n", label, strings.Join(items, ", ")))
		}
	}
	report("invalid_version_format", r.InvalidVersionFormat)
	report("major_minor_duplicates", r.MajorMinorDuplicates)
	report("patch_zero_hits (>1 not allowed)", r.PatchZeroHits)
	report("ga_or_rc_not_top_version", r.GaOrRcNotTopVersion)
	report("multiple_rc_versions", r.RCVersions)
	report("multiple_ga_versions", r.GAVersions)
	report("ga_and_rc_mixed_in_set", r.GaAndRcMixedInSet)
	report("rc_with_invalid_patch_level", r.RCWithInvalidPatchLevel)

	return w.String()
}
