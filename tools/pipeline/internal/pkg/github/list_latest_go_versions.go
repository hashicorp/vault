// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

const goDevReleaseURL = "https://go.dev/dl/?mode=json"

// goDevRelease is the subset of the go.dev release JSON we need.
type goDevRelease struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

// ListLatestGoVersionsReq holds the parameters for fetching the latest Go versions.
type ListLatestGoVersionsReq struct {
	// HTTPClient is the HTTP client to use for go.dev requests. If nil,
	// http.DefaultClient is used.
	HTTPClient *http.Client
	// WriteToGithubOutput controls whether the command layer writes to $GITHUB_OUTPUT.
	WriteToGithubOutput bool
}

// ListLatestGoVersionsRes is the result of fetching the latest stable Go versions.
type ListLatestGoVersionsRes struct {
	// LatestStableVersion is the most recent stable Go release (e.g. "1.23.4").
	LatestStableVersion string `json:"latest_stable_version"`
	// PreviousStableVersion is the second most recent stable Go release.
	PreviousStableVersion string `json:"previous_stable_version"`
}

// Run fetches the latest stable Go versions from go.dev.
func (r *ListLatestGoVersionsReq) Run(ctx context.Context) (*ListLatestGoVersionsRes, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	slog.Default().DebugContext(ctx, "fetching latest Go versions", slog.String("url", goDevReleaseURL))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, goDevReleaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request to go.dev: %w", err)
	}

	httpClient := r.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching Go releases from go.dev: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("go.dev returned unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading go.dev response: %w", err)
	}

	var releases []goDevRelease
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, fmt.Errorf("parsing go.dev response: %w", err)
	}

	var stable []string
	for _, rel := range releases {
		if rel.Stable {
			stable = append(stable, strings.TrimPrefix(rel.Version, "go"))
		}
	}

	if len(stable) < 2 {
		return nil, fmt.Errorf("expected at least two stable Go releases from go.dev, got %d", len(stable))
	}

	res := &ListLatestGoVersionsRes{
		LatestStableVersion:   stable[0],
		PreviousStableVersion: stable[1],
	}

	slog.Default().DebugContext(ctx, "fetched latest Go versions",
		slog.String("latest", res.LatestStableVersion),
		slog.String("previous", res.PreviousStableVersion),
	)

	return res, nil
}

// validate checks that the request is initialized.
func (r *ListLatestGoVersionsReq) validate() error {
	if r == nil {
		return errors.New("uninitialized request")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *ListLatestGoVersionsRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling list-latest-go-versions response to JSON: %w", err)
	}

	return b, nil
}

// ToTable renders the response as a text table.
func (r *ListLatestGoVersionsRes) ToTable() string {
	if r == nil {
		return ""
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{"latest_stable_version", "previous_stable_version"})
	t.AppendRow(table.Row{r.LatestStableVersion, r.PreviousStableVersion})

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t.Render()
}

// ToGithubOutput returns the outputs keyed by the names the workflow expects.
// Each key must be written to $GITHUB_OUTPUT individually so that it is
// addressable as steps.<id>.outputs.<key>.
func (r *ListLatestGoVersionsRes) ToGithubOutput() (map[string]string, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	return map[string]string{
		"latest_stable_version":   r.LatestStableVersion,
		"previous_stable_version": r.PreviousStableVersion,
	}, nil
}

// String returns a human-readable summary.
func (r *ListLatestGoVersionsRes) String() string {
	if r == nil {
		return ""
	}

	return fmt.Sprintf("latest: %s, previous: %s\n", r.LatestStableVersion, r.PreviousStableVersion)
}
