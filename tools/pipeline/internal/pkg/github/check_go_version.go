// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	gh "github.com/google/go-github/v83/github"
	"github.com/jedib0t/go-pretty/v6/table"
)

// BranchGoVersionStatus is the status of a single branch's Go version.
type BranchGoVersionStatus struct {
	Branch           string `json:"branch"`
	GoVersion        string `json:"go_version"`
	Status           string `json:"status"`            // "up-to-date", "out-of-date", "file-missing"
	AcceptedVersions string `json:"accepted_versions"` // e.g. "1.23.4 / 1.22.9"
	// LatestGo is the latest stable Go version. Downstream release workflows
	// compare it against go_version, so it must stay in the JSON rows.
	LatestGo string `json:"latest_go"`
}

// CheckGoVersionReq holds the parameters for checking .go-version across a set
// of branches.
type CheckGoVersionReq struct {
	Owner               string
	Repo                string
	Branches            []string
	HTTPClient          *http.Client // if nil, http.DefaultClient is used
	WriteToGithubOutput bool
}

// CheckGoVersionRes is the result of checking Go versions across branches.
type CheckGoVersionRes struct {
	// Statuses holds a per-branch result.
	Statuses []*BranchGoVersionStatus `json:"statuses"`
	// MismatchCount is the number of branches that are out-of-date.
	MismatchCount int `json:"version_mismatch_count"`
	// FileMissingCount is the number of branches where .go-version was not found.
	FileMissingCount int `json:"file_missing_count"`
}

// Run fetches the latest stable Go versions from go.dev, then checks each
// branch's .go-version file against the latest (N) and previous (N-1) versions.
// A branch is "up-to-date" if it uses either the latest or the previous stable release.
func (r *CheckGoVersionReq) Run(ctx context.Context, client *gh.Client) (*CheckGoVersionRes, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	goVersions, err := (&ListLatestGoVersionsReq{HTTPClient: r.HTTPClient}).Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching latest Go versions: %w", err)
	}

	latest := goVersions.LatestStableVersion
	previous := goVersions.PreviousStableVersion
	acceptedVersions := fmt.Sprintf("%s / %s", latest, previous)

	slog.Default().DebugContext(ctx, "starting check go-version",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.Int("branch_count", len(r.Branches)),
		slog.String("latest", latest),
		slog.String("previous", previous),
	)

	res := &CheckGoVersionRes{
		Statuses: []*BranchGoVersionStatus{},
	}

	for _, branch := range r.Branches {
		status, err := r.checkBranch(ctx, client, branch, latest, previous, acceptedVersions)
		if err != nil {
			return nil, err
		}
		res.Statuses = append(res.Statuses, status)

		switch status.Status {
		case "out-of-date":
			res.MismatchCount++
		case "file-missing":
			res.FileMissingCount++
		}
	}

	slog.Default().DebugContext(ctx, "completed check go-version",
		slog.Int("mismatch_count", res.MismatchCount),
		slog.Int("file_missing_count", res.FileMissingCount),
	)

	return res, nil
}

// checkBranch fetches .go-version for a single branch and returns its status.
func (r *CheckGoVersionReq) checkBranch(ctx context.Context, client *gh.Client, branch, latest, previous, acceptedVersions string) (*BranchGoVersionStatus, error) {
	slog.Default().DebugContext(ctx, "checking branch go version",
		slog.String("branch", branch),
	)

	fileContent, _, _, err := client.Repositories.GetContents(
		ctx,
		r.Owner,
		r.Repo,
		".go-version",
		&gh.RepositoryContentGetOptions{Ref: branch},
	)

	status := &BranchGoVersionStatus{
		Branch:           branch,
		AcceptedVersions: acceptedVersions,
		LatestGo:         latest,
	}

	if err != nil {
		var ghErr *gh.ErrorResponse
		if errors.As(err, &ghErr) && ghErr.Response != nil && ghErr.Response.StatusCode == http.StatusNotFound {
			slog.Default().WarnContext(ctx, ".go-version file not found", slog.String("branch", branch))
			status.GoVersion = "N/A"
			status.Status = "file-missing"
			return status, nil
		}
		return nil, fmt.Errorf("fetching .go-version for branch %s: %w", branch, err)
	}

	if fileContent == nil {
		slog.Default().WarnContext(ctx, ".go-version file content was nil", slog.String("branch", branch))
		status.GoVersion = "N/A"
		status.Status = "file-missing"
		return status, nil
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, fmt.Errorf("decoding .go-version for branch %s: %w", branch, err)
	}

	branchVersion := strings.TrimSpace(content)
	if branchVersion == "" {
		slog.Default().WarnContext(ctx, ".go-version file is empty", slog.String("branch", branch))
		status.GoVersion = "N/A"
		status.Status = "file-missing"
		return status, nil
	}

	status.GoVersion = branchVersion
	if branchVersion == latest || branchVersion == previous {
		status.Status = "up-to-date"
	} else {
		status.Status = "out-of-date"
	}

	slog.Default().DebugContext(ctx, "branch go version checked",
		slog.String("branch", branch),
		slog.String("version", branchVersion),
		slog.String("status", status.Status),
	)

	return status, nil
}

// validate checks that the request has the required fields.
func (r *CheckGoVersionReq) validate() error {
	if r == nil {
		return errors.New("uninitialized request")
	}

	if r.Owner == "" {
		return errors.New("owner is required")
	}

	if r.Repo == "" {
		return errors.New("repo is required")
	}

	if len(r.Branches) == 0 {
		return errors.New("at least one branch is required")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *CheckGoVersionRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling check-go-version response to JSON: %w", err)
	}

	return b, nil
}

// ToTable renders the response as a text table.
func (r *CheckGoVersionRes) ToTable() string {
	if r == nil {
		return ""
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{"branch", "go_version", "status", "accepted_versions"})
	for _, s := range r.Statuses {
		t.AppendRow(table.Row{s.Branch, s.GoVersion, s.Status, s.AcceptedVersions})
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t.Render()
}

// ToGithubOutput returns the outputs keyed by the names the workflow expects.
// Each key must be written to $GITHUB_OUTPUT individually so that it is
// addressable as steps.<id>.outputs.<key>.
func (r *CheckGoVersionRes) ToGithubOutput() (map[string]string, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	statusJSON, err := json.Marshal(r.Statuses)
	if err != nil {
		return nil, fmt.Errorf("marshaling branch go version statuses: %w", err)
	}

	return map[string]string{
		"vault_branch_go_statuses": string(statusJSON),
		"version_mismatch_count":   strconv.Itoa(r.MismatchCount),
		"file_missing_count":       strconv.Itoa(r.FileMissingCount),
	}, nil
}

// String returns a human-readable summary.
func (r *CheckGoVersionRes) String() string {
	if r == nil {
		return ""
	}

	w := strings.Builder{}
	for _, s := range r.Statuses {
		w.WriteString(fmt.Sprintf("[%s] %s — %s (accepted: %s)\n", s.Status, s.Branch, s.GoVersion, s.AcceptedVersions))
	}
	w.WriteString(fmt.Sprintf("mismatches: %d, file-missing: %d\n", r.MismatchCount, r.FileMissingCount))

	return w.String()
}
