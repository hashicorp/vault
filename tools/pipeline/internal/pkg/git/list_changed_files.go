// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/jedib0t/go-pretty/v6/table"
	slogctx "github.com/veqryn/slog-context"
)

// ListChangedFilesReq holds the state and configuration for listing changed files
type ListChangedFilesReq struct {
	// DecodeRes is the result of decoding the pipeline configuration.
	DecodeRes *config.DecodeRes

	// Branch specifies the branch to compare against
	Branch string
	// Range specifies the commit range to compare (e.g., HEAD~5..HEAD)
	Range string
	// Commit specifies a specific commit SHA to analyze
	Commit string
	// GroupFiles requests that changed groups are added to each file

	GroupFiles bool
	// Write a specially formatted response to $GITHUB_OUTPUT
	WriteToGithubOutput bool
}

// ListChangedFilesRes represents the response from listing changed files
type ListChangedFilesRes struct {
	ChangedConfig *changed.Config    `json:"changed_config,omitempty"`
	Files         changed.Files      `json:"files,omitempty"`
	Groups        changed.FileGroups `json:"groups,omitempty"`
}

// ListChangedFilesGithubOutput is our GITHUB_OUTPUT type. It's a slimmed down
// type that only includes file names and groups.
type ListChangedFilesGithubOutput struct {
	Files  []string           `json:"files,omitempty"`
	Groups changed.FileGroups `json:"groups,omitempty"`
}

// Run executes the git list changed files operation
func (g *ListChangedFilesReq) Run(ctx context.Context, client *libgit.Client) (*ListChangedFilesRes, error) {
	slog.Default().DebugContext(ctx, "listing changed files from git")

	err := g.validate(ctx)
	if err != nil {
		return nil, err
	}

	execRes, err := g.getChangedFilesFromGit(ctx, client)
	if err != nil {
		return nil, err
	}

	res := &ListChangedFilesRes{}
	res.Files, err = g.parseChangedFiles(ctx, execRes.Stdout)
	if err != nil {
		return nil, err
	}

	if g.GroupFiles {
		// Store the changed config in the response
		res.ChangedConfig = g.DecodeRes.Config.ChangedFiles

		// Add group metadata to each file using the changed config file grouper.
		changed.GroupFiles(ctx, res.Files, res.ChangedConfig.FileGroups)

		// Get a set of all file groups from all changed files.
		res.Groups = changed.Groups(res.Files)
	}

	return res, nil
}

// ToJSON marshals the response to JSON.
func (r *ListChangedFilesRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling list changed files to JSON: %w", err)
	}

	return b, nil
}

// ToGithubOutput writes a simplified list of changed files to be used $GITHUB_OUTPUT
func (r *ListChangedFilesRes) ToGithubOutput() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	res := &ListChangedFilesGithubOutput{
		Groups: r.Groups,
	}
	if f := r.Files; f != nil {
		res.Files = f.Names()
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("marshaling list changed files GITHUB_OUTPUT to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
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
	t.AppendHeader(table.Row{"path", "groups"})
	for _, file := range r.Files {
		t.AppendRow(table.Row{file.Name(), file.Groups.String()})
	}
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t.Render()
}

// String returns a string representation of the response
func (r *ListChangedFilesRes) String() string {
	if r == nil || len(r.Files) == 0 {
		return "No changed files found"
	}

	w := strings.Builder{}
	for _, name := range r.Files.Names() {
		w.WriteString(name + "\n")
	}
	return w.String()
}

// validate checks that exactly one option is provided
func (g *ListChangedFilesReq) validate(ctx context.Context) error {
	if g == nil {
		return errors.New("uninitialized")
	}

	// Validate that exactly one option is provided
	optionsSet := 0
	if g.Branch != "" {
		optionsSet++
	}
	if g.Range != "" {
		optionsSet++
	}
	if g.Commit != "" {
		optionsSet++
	}

	if optionsSet == 0 {
		return errors.New("must specify one of: --branch, --range, or --commit")
	}

	if optionsSet > 1 {
		return errors.New("can only specify one of: --branch, --range, or --commit")
	}

	if g.GroupFiles {
		if err := g.DecodeRes.Validate(ctx); err != nil {
			return err
		}

		if g.DecodeRes.Config.ChangedFiles == nil {
			return errors.New("changed file grouping was enabled but no changed file grouping config was found in pipeline.hcl")
		}
	}

	return nil
}

// parseChangedFiles parses the raw client output into changed.Files
func (g *ListChangedFilesReq) parseChangedFiles(ctx context.Context, stdout []byte) (changed.Files, error) {
	slog.Default().DebugContext(ctx, "parsing changed files from git client output")
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	files := map[string]struct{}{}
	changedFiles := changed.Files{}

	for scanner.Scan() {
		files[scanner.Text()] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parsing changed files: %w", err)
	}
	for _, file := range slices.Sorted(maps.Keys(files)) {
		changedFiles = append(changedFiles, &changed.File{Filename: file})
	}

	return changedFiles, nil
}

// getChangedFilesFromGit gets the raw changed files output from the git client
func (g *ListChangedFilesReq) getChangedFilesFromGit(ctx context.Context, client *libgit.Client) (*libgit.ExecResponse, error) {
	if g == nil {
		return nil, errors.New("uninitialized")
	}

	switch {
	case g.Branch != "":
		ctx = slogctx.Append(ctx, slog.String("branch", g.Branch))
		slog.Default().DebugContext(ctx, "listing all changed files in branch")

		res, err := client.Log(ctx, &libgit.LogOpts{
			Target:     g.Branch,                                          // show all files for the branch
			Pretty:     libgit.LogPrettyFormatNone,                        // don't add extra formatting
			NameOnly:   true,                                              // list only the names of the files
			DiffFilter: []libgit.LogDiffFilter{libgit.LogDiffFilterAdded}, // only show added files
		})
		if err != nil {
			return res, fmt.Errorf("listing branch changed files: %s, %w", res.String(), err)
		}

		return res, nil
	case g.Range != "":
		ctx = slogctx.Append(ctx, slog.String("range", g.Range))
		slog.Default().DebugContext(ctx, "listing all changed files in target range")

		res, err := client.Log(ctx, &libgit.LogOpts{
			Target:   g.Range,                    // show all files for the range
			Pretty:   libgit.LogPrettyFormatNone, // don't add extra formatting
			NameOnly: true,                       // list only the names of the files
		})
		if err != nil {
			return res, fmt.Errorf("listing range changed files: %s, %w", res.String(), err)
		}

		return res, nil
	case g.Commit != "":
		ctx = slogctx.Append(ctx, slog.String("commit", g.Commit))
		slog.Default().DebugContext(ctx, "listing all changed files for commit")

		res, err := client.Show(ctx, &libgit.ShowOpts{
			Object:   g.Commit,                   // show all files for the range
			Pretty:   libgit.LogPrettyFormatNone, // don't add extra formatting
			NameOnly: true,                       // list only the names of the files
		})
		if err != nil {
			return res, fmt.Errorf("listing range changed files: %s, %w", res.String(), err)
		}

		return res, nil
	default:
		return nil, fmt.Errorf("listing range changed files: no supported target provided")
	}
}
