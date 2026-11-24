// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	"github.com/jedib0t/go-pretty/v6/table"
)

// ListChangedFilesReq is a request to list workflows runs. The fields represent
// various criteria we can use to filter.
type ListChangedFilesReq struct {
	// DecodeRes is the result of decoding the pipeline configuration.
	DecodeRes *config.DecodeRes
	// Our target information
	Owner      string
	Repo       string
	PullNumber int
	CommitSHA  string
	// Options
	GroupFiles          bool
	WriteToGithubOutput bool
	// DisallowedGroups specifies the file groups that we'll check changed files
	// for. They're added to the disallowed groups of the sync policy.
	DisallowedGroups changed.FileGroups
	// SyncPolicy is the name of a sync policy in pipeline.hcl. We add its
	// disallowed groups to ours and report any files it exempts separately.
	SyncPolicy string
}

// ListChangedFilesRes is a list workflows response.
type ListChangedFilesRes struct {
	ChangedConfig *changed.Config    `json:"changed_config,omitempty"`
	Files         changed.Files      `json:"files,omitempty"`
	Groups        changed.FileGroups `json:"groups,omitempty"`
	// DisallowedFiles are changed files that are in a disallowed group and
	// aren't exempt. DisallowedGroups are the groups that they matched.
	DisallowedFiles  changed.Files      `json:"disallowed_files,omitempty"`
	DisallowedGroups changed.FileGroups `json:"disallowed_groups,omitempty"`
	// ExemptFiles are changed files that are in a disallowed group but are
	// exempt in the sync policy.
	ExemptFiles changed.Files `json:"exempt_files,omitempty"`
}

// ListChangedFilesGithubOutput is out GITHUB_OUTPUT type. It's a slimmed down
// type that only include file names and groups. Our workflows rely on these
// JSON keys in conditionals so don't change them without updating those too.
type ListChangedFilesGithubOutput struct {
	Files            []string           `json:"files,omitempty"`
	Groups           changed.FileGroups `json:"groups,omitempty"`
	DisallowedFiles  []string           `json:"disallowed_files,omitempty"`
	DisallowedGroups changed.FileGroups `json:"disallowed_groups,omitempty"`
}

// Run runs the request to gather all instances of the workflow that match
// our filter criteria.
func (r *ListChangedFilesReq) Run(ctx context.Context, client *gh.Client) (*ListChangedFilesRes, error) {
	var err error
	res := &ListChangedFilesRes{}

	if err = r.validate(ctx); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	if r.CommitSHA != "" {
		files, err := r.getCommitFiles(ctx, client)
		if err != nil {
			return nil, err
		}
		res.Files = append(res.Files, files...)
	}

	if r.PullNumber != 0 {
		files, err := r.getPullFiles(ctx, client)
		if err != nil {
			return nil, err
		}
		res.Files = append(res.Files, files...)
	}

	if r.GroupFiles {
		// Store the changed config in the response
		res.ChangedConfig = r.DecodeRes.Config.ChangedFiles

		// Add group metadata to each file using the changed config file grouper.
		changed.GroupFiles(ctx, res.Files, res.ChangedConfig.FileGroups)

		// Get a set of all file groups from all changed files.
		res.Groups = changed.Groups(res.Files)
	}

	if r.checkDisallowed() {
		policy, groups, err := r.disallowedGroups()
		if err != nil {
			return nil, err
		}

		var exempt []string
		if policy != nil {
			exempt = policy.ExemptFromDisallowedGroups
		}
		res.DisallowedFiles, res.DisallowedGroups, res.ExemptFiles = res.Files.Disallowed(groups, exempt)
	}

	return res, nil
}

// checkDisallowed returns whether or not we need to check changed files for
// disallowed groups.
func (r *ListChangedFilesReq) checkDisallowed() bool {
	return r.SyncPolicy != "" || len(r.DisallowedGroups) > 0
}

// disallowedGroups returns our sync policy, if we've been given one, and all of
// the groups that we disallow.
func (r *ListChangedFilesReq) disallowedGroups() (*config.SyncPolicy, changed.FileGroups, error) {
	var policy *config.SyncPolicy
	if r.SyncPolicy != "" {
		var err error
		policy, err = r.DecodeRes.Config.SyncPolicy(r.SyncPolicy)
		if err != nil {
			return nil, nil, err
		}
	}

	groups, err := r.DecodeRes.Config.DisallowedGroups(
		config.WithDisallowedGroupsGroups(r.DisallowedGroups),
		config.WithDisallowedGroupsSyncPolicies(policy),
	)
	if err != nil {
		return nil, nil, err
	}

	return policy, groups, nil
}

// validate ensures that we've been given the minimum filter arguments necessary to complete a
// request. It is always recommended that additional fitlers be given to reduce the response size
// and not exhaust API limits.
func (r *ListChangedFilesReq) validate(ctx context.Context) error {
	if r == nil {
		return errors.New("failed to initialize request")
	}

	if r.Owner == "" {
		return errors.New("no github organization has been provided")
	}

	if r.Repo == "" {
		return errors.New("no github repository has been provided")
	}

	if r.PullNumber == 0 && r.CommitSHA == "" {
		return errors.New("no pull request number or commit SHA has been provided")
	}

	if r.GroupFiles {
		if err := r.DecodeRes.Validate(ctx); err != nil {
			return err
		}

		if r.DecodeRes.Config.ChangedFiles == nil {
			return errors.New("changed file grouping was enabled but no changed file grouping config was found in pipeline.hcl")
		}
	}

	if r.checkDisallowed() {
		if !r.GroupFiles {
			return errors.New("changed file grouping must be enabled to check for disallowed groups")
		}

		if _, _, err := r.disallowedGroups(); err != nil {
			return err
		}
	}

	return nil
}

// getCommitFiles attempts to locate the workflow associated with our workflow name.
func (r *ListChangedFilesReq) getCommitFiles(ctx context.Context, client *gh.Client) (changed.Files, error) {
	opts := &gh.ListOptions{PerPage: PerPageMax}
	files := changed.Files{}
	for {
		commit, res, err := client.Repositories.GetCommit(ctx, r.Owner, r.Repo, r.CommitSHA, opts)
		if err != nil {
			return nil, err
		}

		for _, f := range commit.Files {
			files = append(files, &changed.File{GithubCommitFile: f})
		}

		if res.NextPage == 0 {
			return files, nil
		}

		opts.Page = res.NextPage
	}
}

// getPullFiles attempts to locate the workflow associated with our workflow name.
func (r *ListChangedFilesReq) getPullFiles(ctx context.Context, client *gh.Client) (changed.Files, error) {
	opts := &gh.ListOptions{PerPage: PerPageMax}
	files := changed.Files{}
	for {
		fl, res, err := client.PullRequests.ListFiles(ctx, r.Owner, r.Repo, r.PullNumber, opts)
		if err != nil {
			return nil, err
		}

		for _, f := range fl {
			files = append(files, &changed.File{GithubCommitFile: f})
		}

		if res.NextPage == 0 {
			return files, nil
		}

		opts.Page = res.NextPage
	}
}

// ToJSON marshals the response to JSON.
func (r *ListChangedFilesRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling list changed files to JSON: %w", err)
	}

	return b, nil
}

// ToGithubOutput writes a simplified list of changed files to be used $GITHUB_OUTPUT
func (r ListChangedFilesRes) ToGithubOutput() ([]byte, error) {
	res := &ListChangedFilesGithubOutput{
		Groups:           r.Groups,
		DisallowedGroups: r.DisallowedGroups,
	}
	if f := r.Files; f != nil {
		res.Files = f.Names()
	}
	if f := r.DisallowedFiles; f != nil {
		res.DisallowedFiles = f.Names()
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("marshaling list changed files GITHUB_OUTPUT to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table. If we've found any disallowed
// or exempt files we'll add a column for them.
func (r *ListChangedFilesRes) ToTable(groups bool) string {
	if !groups {
		w := strings.Builder{}
		for _, name := range r.Files.Names() {
			w.WriteString(name + "\n")
		}

		return w.String()
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	disallowed := map[string]string{}
	for _, file := range r.DisallowedFiles {
		disallowed[file.Name()] = "yes"
	}
	for _, file := range r.ExemptFiles {
		disallowed[file.Name()] = "exempt"
	}

	header := table.Row{"path", "groups"}
	if len(disallowed) > 0 {
		header = append(header, "disallowed")
	}
	t.AppendHeader(header)
	for _, file := range r.Files {
		row := table.Row{file.Name(), file.Groups.String()}
		if len(disallowed) > 0 {
			row = append(row, disallowed[file.Name()])
		}
		t.AppendRow(row)
	}
	return t.Render()
}
