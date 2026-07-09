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
	// Include 'main' branch in output
	IncludeMain bool
	// Prefix to add to CE branches (e.g., 'ce' for 'ce/release/<version>')
	CEPrefix string
}

// ListActiveVersionsRes are the active versions and associated metadata for the repo
type ListActiveVersionsRes struct {
	VersionsConfig *VersionsConfig `json:"versions_config,omitempty"`
}

// ActiveVersionMatrixEntry represents a single active version with metadata.
// The idea is to render this as JSON to make building matrices for Github
// Actions easy.
type ActiveVersionMatrixEntry struct {
	Branch  string `json:"branch"`  // Full branch name (e.g., "main", "release/1.19.x+ent", "ce/release/2.0.x")
	Version string `json:"version"` // Version number (e.g., "1.19.x", "2.0.x", "main")
	Edition string `json:"edition"` // "ce" or "ent", derived solely based on ce_active. No other editions metadata (ent.hsm etc.) is available
	LTS     bool   `json:"lts"`     // Whether this is an LTS version
}

// ListActiveVersionsGithubOutput is our GITHUB_OUTPUT type optimized for GitHub Actions workflows.
type ListActiveVersionsGithubOutput struct {
	VersionsConfig          *VersionsConfig             `json:"versions_config,omitempty"`            // Full config from .release/versions.hcl
	Versions                []string                    `json:"versions,omitempty"`                   // e.g., ["1.19.x", "1.20.x", "1.21.x", "2.0.x"]
	CEActiveVersions        []string                    `json:"ce_active_versions,omitempty"`         // e.g., ["2.0.x"] (versions with ce_active: true)
	LTSVersions             []string                    `json:"lts_versions,omitempty"`               // e.g., ["1.19.x"] (versions with lts: true)
	ActiveBranches          []string                    `json:"active_branches,omitempty"`            // e.g., ["release/1.19.x+ent", "release/1.20.x+ent"]
	CEActiveBranches        []string                    `json:"ce_active_branches,omitempty"`         // e.g., ["ce/release/2.0.x"] (with --include-ce-prefix ce)
	LTSActiveBranches       []string                    `json:"lts_active_branches,omitempty"`        // e.g., ["release/1.19.x+ent"]
	AllActiveBranches       []string                    `json:"all_active_branches,omitempty"`        // ENT + CE branches combined
	ActiveVersionsMatrix    []*ActiveVersionMatrixEntry `json:"active_versions_matrix,omitempty"`     // ENT only: [{branch:"main",version:"main",edition:"ent",lts:false}, ...]
	CEActiveVersionsMatrix  []*ActiveVersionMatrixEntry `json:"ce_active_versions_matrix,omitempty"`  // CE only: [{branch:"ce/main",version:"main",edition:"ce",lts:false}, ...]
	AllActiveVersionsMatrix []*ActiveVersionMatrixEntry `json:"all_active_versions_matrix,omitempty"` // Both: ENT + CE entries combined
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

// ActiveVersionWithMetadata represents a version with its associated branch names
type ActiveVersionWithMetadata struct {
	Version          string `json:"version"`
	CEActive         bool   `json:"ce_active"`
	LTS              bool   `json:"lts"`
	EnterpriseBranch string `json:"enterprise_branch"`
	CEBranch         string `json:"ce_branch,omitempty"`
}

// ListActiveVersionsJSONOutput is the JSON output structure
type ListActiveVersionsJSONOutput struct {
	VersionsConfig *VersionsConfig             `json:"versions_config,omitempty"`
	Versions       []ActiveVersionWithMetadata `json:"versions,omitempty"`
}

// ToJSON marshals the response to JSON.
func (l *ListActiveVersionsRes) ToJSON(cePrefix string) ([]byte, error) {
	output := &ListActiveVersionsJSONOutput{
		VersionsConfig: l.VersionsConfig,
		Versions:       []ActiveVersionWithMetadata{},
	}

	for _, version := range slices.Sorted(maps.Keys(l.VersionsConfig.ActiveVersion.Versions)) {
		cfg := l.VersionsConfig.ActiveVersion.Versions[version]
		vwb := ActiveVersionWithMetadata{
			Version:          version,
			CEActive:         cfg.CEActive,
			LTS:              cfg.LTS,
			EnterpriseBranch: enterpriseReleaseBranchForVersion(version),
		}

		if cfg.CEActive {
			vwb.CEBranch = ceReleaseBranchForVersion(version, cePrefix)
		}

		output.Versions = append(output.Versions, vwb)
	}

	b, err := json.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("marshaling list active versions to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (l *ListActiveVersionsRes) ToTable(cePrefix string) string {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"version", "ce active", "lts", "enterprise branch", "ce branch"})
	for _, version := range slices.Sorted(maps.Keys(l.VersionsConfig.ActiveVersion.Versions)) {
		values := l.VersionsConfig.ActiveVersion.Versions[version]
		entBranch := enterpriseReleaseBranchForVersion(version)
		ceBranch := ""

		// If CE active, show CE branch
		if values.CEActive {
			ceBranch = ceReleaseBranchForVersion(version, cePrefix)
		}

		t.AppendRow(table.Row{version, values.CEActive, values.LTS, entBranch, ceBranch})
	}
	return t.Render()
}

// ToGithubOutput writes a JSON encoded versions of ListActiveVersionsRes to
// $GITHUB_OUTPUT. We use an intermediate type to structure the data in a more
// suitable fashion to make usage within Github Actions easier.
func (l ListActiveVersionsRes) ToGithubOutput(includeMain bool, cePrefix string) ([]byte, error) {
	res := &ListActiveVersionsGithubOutput{
		VersionsConfig:          l.VersionsConfig,
		Versions:                slices.Sorted(maps.Keys(l.VersionsConfig.ActiveVersion.Versions)),
		CEActiveVersions:        []string{},
		LTSVersions:             []string{},
		ActiveBranches:          []string{},
		CEActiveBranches:        []string{},
		LTSActiveBranches:       []string{},
		AllActiveBranches:       []string{},
		ActiveVersionsMatrix:    []*ActiveVersionMatrixEntry{}, // ENT only (for backward compatibility)
		CEActiveVersionsMatrix:  []*ActiveVersionMatrixEntry{}, // CE only
		AllActiveVersionsMatrix: []*ActiveVersionMatrixEntry{}, // Both CE and ENT
	}

	// Generate branch names from versions
	for version, cfg := range l.VersionsConfig.ActiveVersion.Versions {
		// Enterprise branch (all versions)
		entBranch := enterpriseReleaseBranchForVersion(version)
		res.ActiveBranches = append(res.ActiveBranches, entBranch)
		res.AllActiveBranches = append(res.AllActiveBranches, entBranch)

		// CE branch (only if ce_active: true)
		if cfg.CEActive {
			res.CEActiveVersions = append(res.CEActiveVersions, version)
			ceBranch := ceReleaseBranchForVersion(version, cePrefix)
			res.CEActiveBranches = append(res.CEActiveBranches, ceBranch)
			res.AllActiveBranches = append(res.AllActiveBranches, ceBranch)
		}

		// LTS branch (only if lts: true)
		if cfg.LTS {
			res.LTSVersions = append(res.LTSVersions, version)
			ltsBranch := enterpriseReleaseBranchForVersion(version)
			res.LTSActiveBranches = append(res.LTSActiveBranches, ltsBranch)
		}

		// NEW: Generate active version matrix entries (always)
		// Enterprise entry (always included)
		entEntry := &ActiveVersionMatrixEntry{
			Branch:  entBranch,
			Version: version,
			Edition: "ent",
			LTS:     cfg.LTS,
		}
		res.ActiveVersionsMatrix = append(res.ActiveVersionsMatrix, entEntry)
		res.AllActiveVersionsMatrix = append(res.AllActiveVersionsMatrix, entEntry)

		// CE entry (only if ce_active: true)
		if cfg.CEActive {
			ceEntry := &ActiveVersionMatrixEntry{
				Branch:  ceReleaseBranchForVersion(version, cePrefix),
				Version: version,
				Edition: "ce",
				LTS:     cfg.LTS,
			}
			res.CEActiveVersionsMatrix = append(res.CEActiveVersionsMatrix, ceEntry)
			res.AllActiveVersionsMatrix = append(res.AllActiveVersionsMatrix, ceEntry)
		}
	}

	// Handle --include-main flag
	if includeMain {
		res.ActiveBranches = append(res.ActiveBranches, "main")

		ceMain := "main"
		if cePrefix != "" {
			ceMain = fmt.Sprintf("%s/main", cePrefix)
		}
		res.CEActiveBranches = append(res.CEActiveBranches, ceMain)

		// Add both to all_active_branches
		res.AllActiveBranches = append(res.AllActiveBranches, "main")
		if cePrefix != "" {
			res.AllActiveBranches = append(res.AllActiveBranches, ceMain)
		}

		// Enterprise main
		entMainEntry := &ActiveVersionMatrixEntry{
			Branch:  "main",
			Version: "main",
			Edition: "ent",
			LTS:     false,
		}
		res.ActiveVersionsMatrix = append(res.ActiveVersionsMatrix, entMainEntry)
		res.AllActiveVersionsMatrix = append(res.AllActiveVersionsMatrix, entMainEntry)

		// CE main (if prefix specified)
		if cePrefix != "" {
			ceMainEntry := &ActiveVersionMatrixEntry{
				Branch:  ceMain,
				Version: "main",
				Edition: "ce",
				LTS:     false,
			}
			res.CEActiveVersionsMatrix = append(res.CEActiveVersionsMatrix, ceMainEntry)
			res.AllActiveVersionsMatrix = append(res.AllActiveVersionsMatrix, ceMainEntry)
		}
	}

	// Sort all slices
	slices.Sort(res.CEActiveVersions)
	slices.Sort(res.LTSVersions)
	slices.Sort(res.ActiveBranches)
	slices.Sort(res.CEActiveBranches)
	slices.Sort(res.LTSActiveBranches)
	slices.Sort(res.AllActiveBranches)

	// Sort active version matrices by branch name
	byBranch := func(a, b *ActiveVersionMatrixEntry) int {
		return slices.Compare([]byte(a.Branch), []byte(b.Branch))
	}
	slices.SortFunc(res.ActiveVersionsMatrix, byBranch)
	slices.SortFunc(res.CEActiveVersionsMatrix, byBranch)
	slices.SortFunc(res.AllActiveVersionsMatrix, byBranch)

	b, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("marshaling list active versions GITHUB_OUTPUT to JSON: %w", err)
	}

	return b, nil
}

// enterpriseReleaseBranchForVersion returns the enterprise release branch name for a version
func enterpriseReleaseBranchForVersion(version string) string {
	return fmt.Sprintf("release/%s+ent", version)
}

// ceReleaseBranchForVersion returns the CE release branch name for a version with optional prefix
func ceReleaseBranchForVersion(version string, prefix string) string {
	branch := fmt.Sprintf("release/%s", version)
	if prefix != "" {
		return fmt.Sprintf("%s/%s", prefix, branch)
	}
	return branch
}
