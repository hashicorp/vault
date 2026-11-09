// Copyright IBM Corp. 2016, 2026
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
	// DisallowedGroups specifies the file groups that must not have changed.
	// They're added to the disallowed groups of the sync policy.
	DisallowedGroups []string
	// SyncPolicy is the name of a sync policy in pipeline.hcl. We add its
	// disallowed groups to ours and won't fail on the files it exempts.
	SyncPolicy string
}

// CheckChangedFilesRes represents the response from checking changed files
type CheckChangedFilesRes struct {
	// Inputs
	ChangedConfig *changed.Config          `json:"changed_config,omitempty"`
	ChangedFiles  changed.Files            `json:"changed_files,omitempty"`
	ChangedGroups changed.FileGroups       `json:"changed_groups,omitempty"`
	CheckedGroups changed.FileGroups       `json:"checked_groups,omitempty"`
	SyncPolicy    *config.SyncPolicyConfig `json:"sync_policy,omitempty"`
	// Outputs
	MatchedFiles  changed.Files      `json:"matched_files,omitempty"`
	MatchedGroups changed.FileGroups `json:"matched_groups,omitempty"`
	ExemptFiles   changed.Files      `json:"exempt_files,omitempty"`
}

// Run executes the git check changed files operation
func (g *CheckChangedFilesReq) Run(ctx context.Context, client *libgit.Client) (*CheckChangedFilesRes, error) {
	slog.Default().DebugContext(ctx, "checking changed files from git for disallowed groups")

	if err := g.validate(ctx); err != nil {
		return nil, err
	}

	syncPolicy, disallowedGroups, err := g.disallowedGroups()
	if err != nil {
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

	ctx = slogctx.Append(
		ctx,
		slog.String("disallowed-groups", disallowedGroups.String()),
		slog.String("sync-policy", g.SyncPolicy),
	)

	res := &CheckChangedFilesRes{
		ChangedConfig: listRes.ChangedConfig,
		ChangedFiles:  listRes.Files,
		ChangedGroups: listRes.Groups,
		CheckedGroups: disallowedGroups,
		SyncPolicy:    syncPolicy,
	}

	exemptFiles := map[string]struct{}{}
	if syncPolicy != nil {
		for _, file := range syncPolicy.ExemptFromDisallowedGroups {
			exemptFiles[file] = struct{}{}
		}
	}

	slog.Default().DebugContext(ctx, "checking changed files for disallowed groups")
	for _, file := range listRes.Files {
		matchedGroups := file.Groups.Intersection(disallowedGroups)
		if len(matchedGroups) < 1 {
			continue
		}

		if _, ok := exemptFiles[file.Name()]; ok {
			res.ExemptFiles = append(res.ExemptFiles, file)
			continue
		}

		res.MatchedFiles = append(res.MatchedFiles, file)
		res.MatchedGroups = res.MatchedGroups.Add(matchedGroups...)
	}

	if len(res.ExemptFiles) > 0 {
		slog.Default().DebugContext(slogctx.Append(
			ctx,
			slog.String("exempt-files", strings.Join(res.ExemptFiles.Names(), ", ")),
		), "found files that the sync policy exempts from disallowed groups")
	}

	if len(res.MatchedFiles) > 0 {
		ctx = slogctx.Append(
			ctx,
			slog.String("matched-groups", res.MatchedGroups.String()),
		)
		slog.Default().DebugContext(ctx, "found files matching disallowed groups")
	}

	return res, nil
}

// disallowedGroups returns the sync policy, if we've been given one, and all of
// the groups that we disallow. That's our disallowed groups plus those of the
// sync policy.
func (g *CheckChangedFilesReq) disallowedGroups() (*config.SyncPolicyConfig, changed.FileGroups, error) {
	groups := changed.FileGroups{}
	for _, group := range g.DisallowedGroups {
		groups = groups.Add(changed.FileGroup(group))
	}

	if g.SyncPolicy == "" {
		return nil, groups, nil
	}

	policy, err := g.DecodeRes.Config.SyncPolicy(g.SyncPolicy)
	if err != nil {
		return nil, nil, err
	}

	for _, group := range policy.DisallowedGroups {
		groups = groups.Add(changed.FileGroup(group))
	}

	return policy, groups, nil
}

func (g *CheckChangedFilesReq) validate(ctx context.Context) error {
	if g == nil {
		return errors.New("uninitialized")
	}

	if len(g.DisallowedGroups) < 1 && g.SyncPolicy == "" {
		return errors.New("no disallowed groups or sync policy have been configured")
	}

	if err := g.DecodeRes.Validate(ctx); err != nil {
		return err
	}

	if g.DecodeRes.Config.ChangedFiles == nil {
		return errors.New("no changed file grouping config was found in pipeline.hcl")
	}

	_, disallowedGroups, err := g.disallowedGroups()
	if err != nil {
		return err
	}

	// A sync policy doesn't have to disallow any groups, but if it doesn't and
	// we weren't given any either then we've got nothing to check.
	if len(disallowedGroups) < 1 {
		return errors.New("no disallowed groups have been configured")
	}

	// Make sure every disallowed group exists. A misspelled group would never
	// match anything and would quietly pass the check.
	configuredGroups := changed.FileGroups{}
	for _, group := range g.DecodeRes.Config.ChangedFiles.Groups {
		configuredGroups = configuredGroups.Add(changed.FileGroup(group.Name))
	}
	for _, group := range disallowedGroups {
		if _, ok := configuredGroups.In(group); !ok {
			return fmt.Errorf("disallowed group %s is not a changed file group in pipeline.hcl", group)
		}
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
	ExemptFiles   []string           `json:"exempt_files,omitempty"`
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
	if f := r.ExemptFiles; f != nil {
		output.ExemptFiles = f.Names()
	}

	b, err := json.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("marshaling check changed files GITHUB_OUTPUT to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table. We include exempt files so
// that we can always see what the sync policy let through.
func (r *CheckChangedFilesRes) ToTable() string {
	if r == nil || (len(r.MatchedFiles) < 1 && len(r.ExemptFiles) < 1) {
		return ""
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"path", "groups", "disallowed groups", "exempt"})
	for _, file := range r.MatchedFiles {
		t.AppendRow(table.Row{
			file.Name(),
			file.Groups.String(),
			r.CheckedGroups.Intersection(file.Groups).String(),
			"",
		})
	}
	for _, file := range r.ExemptFiles {
		t.AppendRow(table.Row{
			file.Name(),
			file.Groups.String(),
			r.CheckedGroups.Intersection(file.Groups).String(),
			"yes",
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
