// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/jedib0t/go-pretty/v6/table"
	slogctx "github.com/veqryn/slog-context"
)

// CheckChangedFilesReq holds the state and configuration for checking changed files
type CheckChangedFilesReq struct {
	// DecodeRes is the result of decoding the pipeline configuration.
	DecodeRes *config.DecodeRes

	// Branch specifies the branch to compare against
	Branch string
	// Range specifies the commit range to compare (e.g., HEAD~5..HEAD)
	Range string
	// Commit specifies a specific commit SHA to analyze
	Commit string

	// Write a specially formatted response to $GITHUB_OUTPUT
	WriteToGithubOutput bool
	// DisallowedGroups specifies the file groups that must not have changed
	DisallowedGroups []string
}

// CheckChangedFilesRes represents the response from checking changed files
type CheckChangedFilesRes struct {
	// Inputs
	ChangedConfig *changed.Config    `json:"changed_config,omitempty"`
	ChangedFiles  changed.Files      `json:"changed_files,omitempty"`
	ChangedGroups changed.FileGroups `json:"changed_groups,omitempty"`
	CheckedGroups changed.FileGroups `json:"checked_groups,omitempty"`
	// Outputs
	MatchedFiles  changed.Files      `json:"matched_files,omitempty"`
	MatchedGroups changed.FileGroups `json:"matched_groups,omitempty"`
}

// Run executes the git check changed files operation
func (g *CheckChangedFilesReq) Run(ctx context.Context, client *libgit.Client) (*CheckChangedFilesRes, error) {
	slog.Default().DebugContext(ctx, "checking changed files from git for disallowed groups")

	if err := g.validate(ctx); err != nil {
		return nil, err
	}

	listReq := &ListChangedFilesReq{
		DecodeRes:  g.DecodeRes,
		Branch:     g.Branch,
		Range:      g.Range,
		Commit:     g.Commit,
		GroupFiles: true,
	}

	listRes, err := listReq.Run(ctx, client)
	if err != nil {
		return nil, err
	}

	disallowdGroups := changed.FileGroups{}
	for _, g := range g.DisallowedGroups {
		disallowdGroups = disallowdGroups.Add(changed.FileGroup(g))
	}

	ctx = slogctx.Append(ctx,
		slog.String("disallowed-groups", disallowdGroups.String()),
	)

	res := &CheckChangedFilesRes{
		ChangedConfig: listRes.ChangedConfig,
		ChangedFiles:  listRes.Files,
		ChangedGroups: listRes.Groups,
		CheckedGroups: disallowdGroups,
		MatchedGroups: disallowdGroups.Intersection(listRes.Groups),
	}

	slog.Default().DebugContext(ctx, "checking changed files for disallowed groups")
	matchedFiles := changed.Files{}
	for _, file := range listRes.Files {
		if i := file.Groups.Intersection(disallowdGroups); len(i) > 0 {
			matchedFiles = append(matchedFiles, file)
		}
	}
	if len(matchedFiles) > 0 {
		res.MatchedFiles = matchedFiles
		ctx = slogctx.Append(ctx,
			slog.String("matched-groups", res.MatchedGroups.String()),
		)
		slog.Default().DebugContext(ctx, "found files matching disallowed groups")
	}

	return res, nil
}

func (g *CheckChangedFilesReq) validate(ctx context.Context) error {
	if g == nil {
		return errors.New("uninitialized")
	}

	if len(g.DisallowedGroups) < 1 {
		return fmt.Errorf("no disallowed groups have been configured")
	}

	if err := g.DecodeRes.Validate(ctx); err != nil {
		return err
	}

	if g.DecodeRes.Config.ChangedFiles == nil {
		return errors.New("no changed file grouping config was found in pipeline.hcl")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *CheckChangedFilesRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling check changed files to JSON: %w", err)
	}

	return b, nil
}

// CheckChangedFilesGithubOutput is our GITHUB_OUTPUT type for check command
type CheckChangedFilesGithubOutput struct {
	ChangedFiles  []string           `json:"changed_files,omitempty"`
	ChangedGroups changed.FileGroups `json:"changed_groups,omitempty"`
	CheckedGroups changed.FileGroups `json:"checked_groups,omitempty"`
	MatchedGroups changed.FileGroups `json:"matched_groups,omitempty"`
	MatchedFiles  []string           `json:"matched_files,omitempty"`
}

// ToGithubOutput writes a simplified check result to be used in $GITHUB_OUTPUT
func (r *CheckChangedFilesRes) ToGithubOutput() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	output := &CheckChangedFilesGithubOutput{
		ChangedGroups: r.ChangedGroups,
		CheckedGroups: r.CheckedGroups,
		MatchedGroups: r.MatchedGroups,
	}
	if f := r.ChangedFiles; f != nil {
		output.ChangedFiles = f.Names()
	}
	if f := r.MatchedFiles; f != nil {
		output.MatchedFiles = f.Names()
	}

	b, err := json.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("marshaling check changed files GITHUB_OUTPUT to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CheckChangedFilesRes) ToTable() string {
	if r == nil || len(r.MatchedGroups) < 1 {
		return ""
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"path", "groups", "disallowed groups"})
	for _, file := range r.MatchedFiles {
		t.AppendRow(table.Row{
			file.Name(),
			file.Groups.String(),
			r.CheckedGroups.Intersection(file.Groups),
		})
	}
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t.Render()
}

// String returns a string representation of the response
func (r *CheckChangedFilesRes) String() string {
	if r == nil || len(r.ChangedFiles) == 0 {
		return "No changed files found"
	}

	w := strings.Builder{}
	for _, name := range r.ChangedFiles.Names() {
		w.WriteString(name + "\n")
	}
	return w.String()
}
