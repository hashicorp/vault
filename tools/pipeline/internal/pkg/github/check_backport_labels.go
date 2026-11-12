// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"regexp"
	"slices"
	"strings"

	libgithub "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	backportLabelPrefix = "backport/"
	// defaultBackportLabelColor is used when no backport label exists to copy from.
	defaultBackportLabelColor = "70F5B5"
)

// backportVersionRegex matches release line versions like "1.19.x". Only labels
// that match are managed so that labels like "backport/skip" are never touched.
var backportVersionRegex = regexp.MustCompile(`^\d+\.\d+\.x$`)

// BackportLabelsAPI is the subset of the GitHub issues API used to manage labels.
type BackportLabelsAPI interface {
	ListLabels(ctx context.Context, owner, repo string, opts *libgithub.ListOptions) ([]*libgithub.Label, *libgithub.Response, error)
	CreateLabel(ctx context.Context, owner, repo string, label *libgithub.Label) (*libgithub.Label, *libgithub.Response, error)
	DeleteLabel(ctx context.Context, owner, repo, name string) (*libgithub.Response, error)
}

// CheckBackportLabelsReq is a request to verify that the backport labels in a
// repository match the active release branches.
type CheckBackportLabelsReq struct {
	Owner string
	Repo  string
	// VersionsDecodeRes is the decoded .release/versions.hcl. Its active versions
	// are the source of truth for the expected labels.
	VersionsDecodeRes *releases.DecodeRes
	// CreateMissing creates missing labels. When false the check is report-only.
	CreateMissing bool
	// DeleteStale deletes stale labels. When false the check is report-only.
	DeleteStale bool
	// WriteToGithubOutput controls whether the command layer writes to $GITHUB_OUTPUT.
	WriteToGithubOutput bool
}

// CheckBackportLabelsRes is the result of the backport label check.
type CheckBackportLabelsRes struct {
	// InSync is true when no missing or stale labels remain after any requested changes.
	InSync bool `json:"in_sync"`
	// Expected are the labels that should exist.
	Expected []string `json:"expected"`
	// Existing are the managed backport labels that were found.
	Existing []string `json:"existing"`
	// Missing are expected labels that do not exist, excluding any that were created.
	Missing []string `json:"missing"`
	// Stale are existing labels with no active release branch, excluding any that were deleted.
	Stale []string `json:"stale"`
	// Created are the labels created by this run.
	Created []string `json:"created"`
	// Deleted are the labels deleted by this run.
	Deleted []string `json:"deleted"`
}

// BackportLabelForVersion returns the backport label for a release line.
// Example: "1.19.x" -> "backport/1.19.x"
func BackportLabelForVersion(version string) string {
	return backportLabelPrefix + version
}

// Run compares the repository's backport labels to the active versions and
// optionally creates missing and deletes stale labels.
func (r *CheckBackportLabelsReq) Run(ctx context.Context, api BackportLabelsAPI) (*CheckBackportLabelsRes, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	expected, err := r.expectedLabels(ctx)
	if err != nil {
		return nil, err
	}

	// An empty active version list means a bad versions.hcl, and deleting on it
	// would remove every backport label.
	if r.DeleteStale && len(expected) == 0 {
		return nil, errors.New("refusing to delete stale labels: versions.hcl has no active versions")
	}

	labels, err := listBackportLabels(ctx, api, r.Owner, r.Repo)
	if err != nil {
		return nil, err
	}

	existing := map[string]*libgithub.Label{}
	for _, l := range labels {
		existing[l.GetName()] = l
	}

	res := &CheckBackportLabelsRes{
		Expected: slices.Sorted(maps.Keys(expected)),
		Existing: slices.Sorted(maps.Keys(existing)),
		Missing:  []string{},
		Stale:    []string{},
		Created:  []string{},
		Deleted:  []string{},
	}

	for _, name := range res.Expected {
		if _, ok := existing[name]; !ok {
			res.Missing = append(res.Missing, name)
		}
	}
	for _, name := range res.Existing {
		if _, ok := expected[name]; !ok {
			res.Stale = append(res.Stale, name)
		}
	}

	if r.CreateMissing && len(res.Missing) > 0 {
		template := backportLabelTemplate(labels)
		for _, name := range res.Missing {
			slog.Default().InfoContext(ctx, "creating backport label", slog.String("label", name))
			_, _, err := api.CreateLabel(ctx, r.Owner, r.Repo, &libgithub.Label{
				Name:        libgithub.Ptr(name),
				Color:       libgithub.Ptr(template.GetColor()),
				Description: libgithub.Ptr(template.GetDescription()),
			})
			if err != nil {
				return res, fmt.Errorf("creating label %q: %w", name, err)
			}
			res.Created = append(res.Created, name)
		}
		res.Missing = []string{}
	}

	if r.DeleteStale && len(res.Stale) > 0 {
		for _, name := range res.Stale {
			slog.Default().InfoContext(ctx, "deleting stale backport label", slog.String("label", name))
			if _, err := api.DeleteLabel(ctx, r.Owner, r.Repo, name); err != nil {
				return res, fmt.Errorf("deleting label %q: %w", name, err)
			}
			res.Deleted = append(res.Deleted, name)
		}
		res.Stale = []string{}
	}

	res.InSync = len(res.Missing) == 0 && len(res.Stale) == 0

	return res, nil
}

// expectedLabels returns the set of backport labels that should exist. The
// main branch has no backport label.
func (r *CheckBackportLabelsReq) expectedLabels(ctx context.Context) (map[string]struct{}, error) {
	if err := r.VersionsDecodeRes.Validate(ctx); err != nil {
		return nil, err
	}

	expected := map[string]struct{}{}
	for version := range r.VersionsDecodeRes.Config.ActiveVersion.Versions {
		expected[BackportLabelForVersion(version)] = struct{}{}
	}

	return expected, nil
}

// listBackportLabels returns every managed backport label in the repository.
func listBackportLabels(ctx context.Context, api BackportLabelsAPI, owner, repo string) ([]*libgithub.Label, error) {
	var found []*libgithub.Label
	opts := &libgithub.ListOptions{PerPage: 100}

	for {
		labels, resp, err := api.ListLabels(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("listing labels for %s/%s: %w", owner, repo, err)
		}

		for _, l := range labels {
			name := l.GetName()
			if strings.HasPrefix(name, backportLabelPrefix) &&
				backportVersionRegex.MatchString(strings.TrimPrefix(name, backportLabelPrefix)) {
				found = append(found, l)
			}
		}

		if resp == nil || resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return found, nil
}

// backportLabelTemplate returns the label whose color and description are
// copied to new labels. It falls back to the default color when none exist.
func backportLabelTemplate(labels []*libgithub.Label) *libgithub.Label {
	if len(labels) > 0 {
		return labels[0]
	}

	return &libgithub.Label{Color: libgithub.Ptr(defaultBackportLabelColor)}
}

func (r *CheckBackportLabelsReq) validate() error {
	if r == nil {
		return errors.New("uninitialized request")
	}

	if r.Owner == "" {
		return errors.New("owner is required")
	}

	if r.Repo == "" {
		return errors.New("repo is required")
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *CheckBackportLabelsRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling check-backport-labels response to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CheckBackportLabelsRes) ToTable() string {
	if r == nil {
		return ""
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{"label", "status"})
	for _, group := range []struct {
		status string
		labels []string
	}{
		{"created", r.Created},
		{"deleted", r.Deleted},
		{"missing", r.Missing},
		{"stale", r.Stale},
	} {
		for _, name := range group.labels {
			t.AppendRow(table.Row{name, group.status})
		}
	}
	for _, name := range r.Expected {
		if !slices.Contains(r.Missing, name) && !slices.Contains(r.Created, name) {
			t.AppendRow(table.Row{name, "ok"})
		}
	}

	t.SuppressTrailingSpaces()

	return t.Render()
}

// ToGithubOutput returns the result formatted for $GITHUB_OUTPUT.
func (r *CheckBackportLabelsRes) ToGithubOutput() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	w := strings.Builder{}
	w.WriteString(fmt.Sprintf("in_sync=%t\n", r.InSync))
	w.WriteString(fmt.Sprintf("missing=%s\n", strings.Join(r.Missing, " ")))
	w.WriteString(fmt.Sprintf("stale=%s\n", strings.Join(r.Stale, " ")))
	w.WriteString(fmt.Sprintf("created=%s\n", strings.Join(r.Created, " ")))
	w.WriteString(fmt.Sprintf("deleted=%s\n", strings.Join(r.Deleted, " ")))

	return []byte(w.String()), nil
}

// String returns a human-readable summary.
func (r *CheckBackportLabelsRes) String() string {
	if r == nil {
		return ""
	}

	return fmt.Sprintf("missing: [%s], stale: [%s]", strings.Join(r.Missing, " "), strings.Join(r.Stale, " "))
}
