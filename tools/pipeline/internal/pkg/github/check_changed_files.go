// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	"github.com/jedib0t/go-pretty/v6/table"
	slogctx "github.com/veqryn/slog-context"
)

// CheckChangedFilesReq is a request to check the files that changed in a pull
// request or commit for any that are in disallowed groups.
type CheckChangedFilesReq struct {
	// DecodeRes is the result of decoding the pipeline configuration.
	DecodeRes *config.DecodeRes
	// Our target information
	Owner      string
	Repo       string
	PullNumber int
	CommitSHA  string
	// DisallowedGroups specifies the file groups that must not have changed.
	// They're added to the disallowed groups of the sync policy.
	DisallowedGroups changed.FileGroups
	// SyncPolicy is the name of a sync policy in pipeline.hcl. We add its
	// disallowed groups to ours and won't fail on the files it exempts.
	SyncPolicy string
}

// CheckChangedFilesRes is the response from checking changed files.
type CheckChangedFilesRes struct {
	// Inputs
	ChangedConfig *changed.Config    `json:"changed_config,omitempty"`
	ChangedFiles  changed.Files      `json:"changed_files,omitempty"`
	ChangedGroups changed.FileGroups `json:"changed_groups,omitempty"`
	CheckedGroups changed.FileGroups `json:"checked_groups,omitempty"`
	SyncPolicy    *config.SyncPolicy `json:"sync_policy,omitempty"`
	// Outputs
	MatchedFiles  changed.Files      `json:"matched_files,omitempty"`
	MatchedGroups changed.FileGroups `json:"matched_groups,omitempty"`
	ExemptFiles   changed.Files      `json:"exempt_files,omitempty"`
}

// Run lists the changed files of the pull request or commit and checks them
// for any that are in disallowed groups.
func (r *CheckChangedFilesReq) Run(ctx context.Context, client *gh.Client) (*CheckChangedFilesRes, error) {
	if err := r.validate(ctx); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	listReq := &ListChangedFilesReq{
		DecodeRes:        r.DecodeRes,
		Owner:            r.Owner,
		Repo:             r.Repo,
		PullNumber:       r.PullNumber,
		CommitSHA:        r.CommitSHA,
		GroupFiles:       true,
		DisallowedGroups: r.DisallowedGroups,
		SyncPolicy:       r.SyncPolicy,
	}

	syncPolicy, disallowedGroups, err := listReq.disallowedGroups()
	if err != nil {
		return nil, err
	}

	ctx = slogctx.Append(
		ctx,
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.Int("pr", r.PullNumber),
		slog.String("commit", r.CommitSHA),
		slog.String("disallowed-groups", disallowedGroups.String()),
		slog.String("sync-policy", r.SyncPolicy),
	)

	slog.Default().DebugContext(ctx, "listing changed files from github")
	listRes, err := listReq.Run(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("listing changed files: %w", err)
	}

	res := &CheckChangedFilesRes{
		ChangedConfig: listRes.ChangedConfig,
		ChangedFiles:  listRes.Files,
		ChangedGroups: listRes.Groups,
		CheckedGroups: disallowedGroups,
		SyncPolicy:    syncPolicy,
		MatchedFiles:  listRes.DisallowedFiles,
		MatchedGroups: listRes.DisallowedGroups,
		ExemptFiles:   listRes.ExemptFiles,
	}

	if len(res.ExemptFiles) > 0 {
		slog.Default().DebugContext(slogctx.Append(
			ctx,
			slog.String("exempt-files", strings.Join(res.ExemptFiles.Names(), ", ")),
		), "found files that the sync policy exempts from disallowed groups")
	}

	if len(res.MatchedFiles) > 0 {
		slog.Default().DebugContext(slogctx.Append(
			ctx,
			slog.String("matched-groups", res.MatchedGroups.String()),
		), "found files matching disallowed groups")
	}

	return res, nil
}

func (r *CheckChangedFilesReq) validate(ctx context.Context) error {
	if r == nil {
		return errors.New("uninitialized")
	}

	if len(r.DisallowedGroups) < 1 && r.SyncPolicy == "" {
		return errors.New("no disallowed groups or sync policy have been configured")
	}

	return r.DecodeRes.Validate(ctx)
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
