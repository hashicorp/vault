// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	libgithub "github.com/google/go-github/v83/github"
	"github.com/jedib0t/go-pretty/v6/table"
)

// CheckVersionTagsReq holds the parameters for validating that a set of
// versions satisfy the tag sequencing rules for a given repository.
type CheckVersionTagsReq struct {
	Owner     string
	Repo      string
	Versions  []string
	TagPrefix string // e.g. "v"
	TagSuffix string // e.g. "+ent"
	// WriteToGithubOutput controls whether the command layer writes to $GITHUB_OUTPUT.
	// It is not used by Run itself but is read by the command layer.
	WriteToGithubOutput bool
}

// CheckVersionTagsRes is the result of the version-tag check.
type CheckVersionTagsRes struct {
	// Passed is true when every version passed all tag rules.
	Passed bool `json:"passed"`
	// Results holds per-version outcomes.
	Results []*VersionTagCheckResult `json:"results"`
}

// VersionTagCheckResult holds the outcome for a single version.
type VersionTagCheckResult struct {
	Version string `json:"version"`
	Passed  bool   `json:"passed"`
	Rule    string `json:"rule"`
	Detail  string `json:"detail"`
}

// Run fetches all tags from GitHub and validates each version against the tag
// sequencing rules.
func (r *CheckVersionTagsReq) Run(ctx context.Context, client *libgithub.Client) (*CheckVersionTagsRes, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	slog.Default().DebugContext(ctx, "starting check-version-tags",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.Int("version_count", len(r.Versions)),
	)

	tags, err := r.fetchTags(ctx, client)
	if err != nil {
		return nil, err
	}

	tagSet := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		tagSet[t] = struct{}{}
	}

	res := &CheckVersionTagsRes{Passed: true}

	for _, version := range r.Versions {
		result := r.checkVersion(version, tagSet)
		res.Results = append(res.Results, result)
		if !result.Passed {
			res.Passed = false
		}
	}

	slog.Default().DebugContext(ctx, "completed check-version-tags", slog.Bool("passed", res.Passed))

	return res, nil
}

// fetchTags retrieves all repository tags from GitHub matching the configured
// prefix and suffix, paging through all results.
func (r *CheckVersionTagsReq) fetchTags(ctx context.Context, client *libgithub.Client) ([]string, error) {
	slog.Default().DebugContext(ctx, "fetching repository tags",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
	)

	var allTags []string
	opts := &libgithub.ListOptions{PerPage: 100}

	for {
		tags, resp, err := client.Repositories.ListTags(ctx, r.Owner, r.Repo, opts)
		if err != nil {
			return nil, fmt.Errorf("listing tags for %s/%s: %w", r.Owner, r.Repo, err)
		}

		for _, t := range tags {
			name := t.GetName()
			if strings.HasPrefix(name, r.TagPrefix) && strings.HasSuffix(name, r.TagSuffix) {
				allTags = append(allTags, name)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	slog.Default().DebugContext(ctx, "fetched tags", slog.Int("count", len(allTags)))

	return allTags, nil
}

// checkVersion applies the four tag sequencing rules to a single version.
func (r *CheckVersionTagsReq) checkVersion(version string, tagSet map[string]struct{}) *VersionTagCheckResult {
	result := &VersionTagCheckResult{Version: version, Passed: true}

	major, minor, patch, rc, ok := parseVersion(version)
	if !ok {
		result.Passed = false
		result.Rule = "parse"
		result.Detail = fmt.Sprintf("version %q could not be parsed", version)
		return result
	}

	currentTag := fmt.Sprintf("%s%d.%d.%d%s%s", r.TagPrefix, major, minor, patch, rc, r.TagSuffix)

	// Rule 0: the exact tag must not already exist.
	if _, exists := tagSet[currentTag]; exists {
		result.Passed = false
		result.Rule = "rule-0"
		result.Detail = fmt.Sprintf("tag %q already exists", currentTag)
		return result
	}

	switch {
	case rc == "" && patch == 0:
		// Rule 1: GA (patch=0) requires at least one RC tag for the same major.minor.
		rcPattern := regexp.MustCompile(
			fmt.Sprintf(`^%s%d\.%d\.0-rc[0-9]+%s$`,
				regexp.QuoteMeta(r.TagPrefix), major, minor, regexp.QuoteMeta(r.TagSuffix)),
		)
		if !matchesAny(rcPattern, tagSet) {
			result.Passed = false
			result.Rule = "rule-1"
			result.Detail = fmt.Sprintf("no RC tag found for %d.%d.0", major, minor)
		}

	case rc == "-rc1":
		// Rule 3: rc1 must be the first tag for that major.minor — no prior tags allowed.
		priorPattern := regexp.MustCompile(
			fmt.Sprintf(`^%s%d\.%d\..*%s$`,
				regexp.QuoteMeta(r.TagPrefix), major, minor, regexp.QuoteMeta(r.TagSuffix)),
		)
		if matchesAny(priorPattern, tagSet) {
			result.Passed = false
			result.Rule = "rule-3"
			result.Detail = fmt.Sprintf("prior tags exist for %d.%d, rc1 must be first", major, minor)
		}

	case strings.HasPrefix(rc, "-rc") && rc != "-rc1":
		// Rule 2: rcN (N>1) requires rc(N-1) to exist.
		n, _ := strconv.Atoi(strings.TrimPrefix(rc, "-rc"))
		prevTag := fmt.Sprintf("%s%d.%d.0-rc%d%s", r.TagPrefix, major, minor, n-1, r.TagSuffix)
		if _, exists := tagSet[prevTag]; !exists {
			result.Passed = false
			result.Rule = "rule-2"
			result.Detail = fmt.Sprintf("prerequisite tag %q not found", prevTag)
		}

	default:
		// Rule 4: minor release (patch>0, no RC) requires patch-1 to exist.
		prevTag := fmt.Sprintf("%s%d.%d.%d%s", r.TagPrefix, major, minor, patch-1, r.TagSuffix)
		if _, exists := tagSet[prevTag]; !exists {
			result.Passed = false
			result.Rule = "rule-4"
			result.Detail = fmt.Sprintf("prerequisite tag %q not found", prevTag)
		}
	}

	return result
}

// parseVersion breaks "X.Y.Z" or "X.Y.Z-rcN" into its numeric components.
func parseVersion(v string) (major, minor, patch int, rc string, ok bool) {
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(-rc[0-9]+)?$`)
	m := re.FindStringSubmatch(v)
	if m == nil {
		return 0, 0, 0, "", false
	}

	major, _ = strconv.Atoi(m[1])
	minor, _ = strconv.Atoi(m[2])
	patch, _ = strconv.Atoi(m[3])
	rc = m[4]

	return major, minor, patch, rc, true
}

// matchesAny reports whether any key in tagSet matches pattern.
func matchesAny(pattern *regexp.Regexp, tagSet map[string]struct{}) bool {
	for tag := range tagSet {
		if pattern.MatchString(tag) {
			return true
		}
	}

	return false
}

// validate checks that the request has the required fields.
func (r *CheckVersionTagsReq) validate() error {
	if r == nil {
		return errors.New("uninitialized request")
	}

	if r.Owner == "" {
		return errors.New("owner is required")
	}

	if r.Repo == "" {
		return errors.New("repo is required")
	}

	if len(r.Versions) == 0 {
		return errors.New("at least one version is required")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *CheckVersionTagsRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling check-version-tags response to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CheckVersionTagsRes) ToTable() string {
	if r == nil {
		return ""
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{"version", "passed", "rule", "detail"})
	for _, result := range r.Results {
		t.AppendRow(table.Row{result.Version, result.Passed, result.Rule, result.Detail})
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t.Render()
}

// ToGithubOutput returns the failure map formatted for $GITHUB_OUTPUT using
// the same key=value multiline format as the original shell script.
func (r *CheckVersionTagsRes) ToGithubOutput() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	w := strings.Builder{}
	w.WriteString(fmt.Sprintf("passed=%t\n", r.Passed))

	if !r.Passed {
		w.WriteString("version_failure_map<<EOF\n")
		for _, result := range r.Results {
			if !result.Passed {
				w.WriteString(fmt.Sprintf("%s:%s\n", result.Version, result.Detail))
			}
		}
		w.WriteString("EOF\n")
	}

	return []byte(w.String()), nil
}

// String returns a human-readable summary.
func (r *CheckVersionTagsRes) String() string {
	if r == nil {
		return ""
	}

	w := strings.Builder{}
	for _, result := range r.Results {
		status := "PASS"
		if !result.Passed {
			status = "FAIL"
		}
		w.WriteString(fmt.Sprintf("[%s] %s", status, result.Version))
		if result.Detail != "" {
			w.WriteString(fmt.Sprintf(" — %s (%s)", result.Detail, result.Rule))
		}
		w.WriteString("\n")
	}

	return w.String()
}
