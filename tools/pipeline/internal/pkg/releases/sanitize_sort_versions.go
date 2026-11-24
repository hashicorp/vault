// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/jedib0t/go-pretty/v6/table"
)

// versionPattern matches X.Y.Z with optional suffixes: -lts, -rcN, -ce, +ent, or other pre-release labels.
// Captured groups: (major)(minor)(patch)(suffix)(ent)
var versionPattern = regexp.MustCompile(
	`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-lts|-rc[1-9][0-9]*|-ce|-[^rc][a-zA-Z0-9]*)?(\+ent)?$`,
)

var rcSuffixPattern = regexp.MustCompile(`^-rc[0-9]+$`)

// SanitizeSortVersionsReq holds the request parameters for sanitizing and
// sorting a list of version strings.
type SanitizeSortVersionsReq struct {
	// Versions is the space-separated or slice of raw version strings to process.
	Versions []string
	// Descending controls sort direction. When true (the default) versions are
	// sorted highest-first. When false they are sorted lowest-first.
	Descending bool
	// WriteToGithubOutput controls whether the command layer writes to $GITHUB_OUTPUT.
	// It is not used by Run itself but is read by the command layer.
	WriteToGithubOutput bool
}

// SanitizeSortVersionsRes is the result of sanitizing and sorting versions.
type SanitizeSortVersionsRes struct {
	// SortedVersions is the deduplicated, normalized, sorted list of valid versions.
	SortedVersions []string `json:"sorted_versions"`
	// LTSVersions holds the base X.Y.Z for any input version carrying a -lts suffix.
	LTSVersions []string `json:"lts_versions"`
	// InvalidVersions holds any input strings that did not match the expected format.
	InvalidVersions []string `json:"invalid_versions"`
	// DuplicateVersions holds normalized versions that appeared more than once in
	// the input. Only the first occurrence is kept in SortedVersions.
	DuplicateVersions []string `json:"duplicate_versions"`
}

// Run sanitizes, normalizes, and sorts the requested versions.
func (r *SanitizeSortVersionsReq) Run(ctx context.Context) (*SanitizeSortVersionsRes, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	slog.Default().DebugContext(ctx, "starting sanitize-sort-versions", slog.Int("input_count", len(r.Versions)))

	res := &SanitizeSortVersionsRes{
		SortedVersions:    []string{},
		LTSVersions:       []string{},
		InvalidVersions:   []string{},
		DuplicateVersions: []string{},
	}

	seen := make(map[string]struct{})

	for _, v := range r.Versions {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}

		m := versionPattern.FindStringSubmatch(v)
		if m == nil {
			slog.Default().DebugContext(ctx, "invalid version", slog.String("version", v))
			res.InvalidVersions = append(res.InvalidVersions, v)
			continue
		}

		base := fmt.Sprintf("%s.%s.%s", m[1], m[2], m[3])
		suffix := m[4]
		entSuffix := m[5]

		// Collect LTS base versions separately.
		if suffix == "-lts" {
			res.LTSVersions = append(res.LTSVersions, base)
		}

		// Normalize: keep RC suffix and +ent if present, keep +ent alone, otherwise use bare X.Y.Z.
		var normalized string
		switch {
		case rcSuffixPattern.MatchString(suffix):
			normalized = v
		case entSuffix != "":
			normalized = base + entSuffix
		default:
			normalized = base
		}

		if _, alreadySeen := seen[normalized]; alreadySeen {
			slog.Default().WarnContext(ctx, "duplicate version encountered, skipping",
				slog.String("version", normalized),
			)
			res.DuplicateVersions = append(res.DuplicateVersions, normalized)
			continue
		}
		seen[normalized] = struct{}{}

		res.SortedVersions = append(res.SortedVersions, normalized)
	}

	// Sort by semver. Versions that fail semver parsing fall back to
	// lexicographic ordering at the end of the list.
	slices.SortFunc(res.SortedVersions, func(a, b string) int {
		sa, errA := semver.NewVersion(a)
		sb, errB := semver.NewVersion(b)

		var result int
		switch {
		case errA != nil && errB != nil:
			result = cmp.Compare(a, b) // both unparseable: lexicographic ascending
		case errA != nil:
			result = 1 // unparseable a sorts last
		case errB != nil:
			result = -1 // unparseable b sorts last
		default:
			result = sa.Compare(sb) // ascending
		}

		if r.Descending {
			return -result
		}
		return result
	})

	slog.Default().DebugContext(ctx, "completed sanitize-sort-versions",
		slog.Int("sorted_count", len(res.SortedVersions)),
		slog.Int("lts_count", len(res.LTSVersions)),
		slog.Int("invalid_count", len(res.InvalidVersions)),
	)

	return res, nil
}

// validate checks that the request is initialized and has input.
func (r *SanitizeSortVersionsReq) validate() error {
	if r == nil {
		return errors.New("uninitialized request")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *SanitizeSortVersionsRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling sanitize-sort-versions response to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *SanitizeSortVersionsRes) ToTable() string {
	if r == nil {
		return ""
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{"sorted_versions", "lts_versions", "invalid_versions", "duplicate_versions"})

	maxLen := max(max(max(len(r.SortedVersions), len(r.LTSVersions)), len(r.InvalidVersions)), len(r.DuplicateVersions))
	for i := range maxLen {
		row := table.Row{"", "", "", ""}
		if i < len(r.SortedVersions) {
			row[0] = r.SortedVersions[i]
		}
		if i < len(r.LTSVersions) {
			row[1] = r.LTSVersions[i]
		}
		if i < len(r.InvalidVersions) {
			row[2] = r.InvalidVersions[i]
		}
		if i < len(r.DuplicateVersions) {
			row[3] = r.DuplicateVersions[i]
		}
		t.AppendRow(row)
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t.Render()
}

// ToGithubOutput returns the outputs formatted for $GITHUB_OUTPUT as
// newline-delimited key=value pairs.
func (r *SanitizeSortVersionsRes) ToGithubOutput() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	out := strings.Builder{}
	out.WriteString(fmt.Sprintf("sorted_versions=%s\n", strings.Join(r.SortedVersions, " ")))
	out.WriteString(fmt.Sprintf("lts_versions=%s\n", strings.Join(r.LTSVersions, " ")))
	out.WriteString(fmt.Sprintf("invalid_versions=%s\n", strings.Join(r.InvalidVersions, " ")))
	out.WriteString(fmt.Sprintf("duplicate_versions=%s\n", strings.Join(r.DuplicateVersions, " ")))

	return []byte(out.String()), nil
}

// String returns a human-readable summary of the response.
func (r *SanitizeSortVersionsRes) String() string {
	if r == nil {
		return ""
	}

	w := strings.Builder{}
	w.WriteString(fmt.Sprintf("sorted_versions: %s\n", strings.Join(r.SortedVersions, " ")))
	w.WriteString(fmt.Sprintf("lts_versions: %s\n", strings.Join(r.LTSVersions, " ")))
	w.WriteString(fmt.Sprintf("invalid_versions: %s\n", strings.Join(r.InvalidVersions, " ")))
	w.WriteString(fmt.Sprintf("duplicate_versions: %s\n", strings.Join(r.DuplicateVersions, " ")))

	return w.String()
}
