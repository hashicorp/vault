// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"slices"

	libgitclient "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/jedib0t/go-pretty/v6/table"
)

// ListActiveVersionsReq is a request to list the active branch versions from the
// .release/metadata file
type ListActiveVersionsReq struct {
	// VersionsDecodeRes is the result of decoding .release/versions.hcl. If we
	// auto-loaded it during our command initialization then we can return the
	// contents.
	VersionsDecodeRes *DecodeRes
	// Write the active versions to $GITHUB_OUTPUT
	WriteToGithubOutput bool
}

// ListActiveVersionsRes are the active versions and associated metadata for the repo
type ListActiveVersionsRes struct {
	VersionsConfig *VersionsConfig `json:"versions_config,omitempty"`
}

// ListActiveVersionsGithubOutput is our GITHUB_OUTPUT type. While ListActiveVersionsReq is designed to match the schema of the releases source file, this type
// is designed for maximal utility in Github Actions workflows and their associated built-in functions.
type ListActiveVersionsGithubOutput struct {
	VersionsConfig   *VersionsConfig `json:"versions_config,omitempty"`
	Versions         []string        `json:"versions,omitempty"`
	CEActiveVersions []string        `json:"ce_active_versions,omitempty"`
	LTSVersions      []string        `json:"lts_versions,omitempty"`
}

// Run runs the dynamic configuration request
func (l *ListActiveVersionsReq) Run(ctx context.Context, git *libgitclient.Client) (*ListActiveVersionsRes, error) {
	if l == nil {
		return nil, fmt.Errorf("list active versions request is uninitialized")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	slog.Default().DebugContext(ctx, "running list active versions request")
	if err := l.VersionsDecodeRes.Validate(ctx); err != nil {
		return nil, err
	}

	return &ListActiveVersionsRes{VersionsConfig: l.VersionsDecodeRes.Config}, nil
}

// ToJSON marshals the response to JSON.
func (l *ListActiveVersionsRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(l)
	if err != nil {
		return nil, fmt.Errorf("marshaling list changed files to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (l *ListActiveVersionsRes) ToTable() string {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"version", "ce active", "lts"})
	for _, version := range slices.Sorted(maps.Keys(l.VersionsConfig.ActiveVersion.Versions)) {
		values := l.VersionsConfig.ActiveVersion.Versions[version]
		t.AppendRow(table.Row{version, values.CEActive, values.LTS})
	}
	return t.Render()
}

// ToGithubOutput writes a JSON encoded versions of ListActiveVersionsRes to
// $GITHUB_OUTPUT. We use an intermediate type to structure the data in a more
// suitable fashion to make usage within Github Actions easier.
func (r ListActiveVersionsRes) ToGithubOutput() ([]byte, error) {
	res := &ListActiveVersionsGithubOutput{
		VersionsConfig:   r.VersionsConfig,
		Versions:         slices.Sorted(maps.Keys(r.VersionsConfig.ActiveVersion.Versions)),
		CEActiveVersions: []string{},
		LTSVersions:      []string{},
	}

	for version, cfg := range r.VersionsConfig.ActiveVersion.Versions {
		if cfg.CEActive {
			res.CEActiveVersions = append(res.CEActiveVersions, version)
		}
		if cfg.LTS {
			res.LTSVersions = append(res.LTSVersions, version)
		}
	}
	slices.Sort(res.CEActiveVersions)
	slices.Sort(res.LTSVersions)

	b, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("marshaling list active versions GITHUB_OUTPUT to JSON: %w", err)
	}

	return b, nil
}
