// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	gh "github.com/google/go-github/v83/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/jedib0t/go-pretty/v6/table"
	slogctx "github.com/veqryn/slog-context"
)

// ListOpenPRsReq is a request to list open pull requests across one or more branches.
type ListOpenPRsReq struct {
	Owner    string
	Repo     string
	Branches []string

	// Releases is the decoded versions.hcl configuration. When Branches is
	// empty, active release branches are derived from this config.
	Releases *releases.DecodeRes

	// IncludeDrafts controls whether draft PRs are included. Defaults to false,
	// matching the behavior of the open-pr-checker workflow.
	IncludeDrafts bool

	// Output opts
	WriteToGithubOutput bool
	GithubOutputKey     string
}

// OpenPullRequest holds the fields we surface for each open PR.
type OpenPullRequest struct {
	Number    int      `json:"number"`
	URL       string   `json:"url"`
	Author    string   `json:"author"`
	Milestone string   `json:"milestone,omitempty"`
	Backports []string `json:"backports,omitempty"`
}

// ListOpenPRsRes is the response from listing open pull requests.
// BranchPRs maps each requested branch to its list of open PRs.
type ListOpenPRsRes struct {
	Owner     string                        `json:"owner"`
	Repo      string                        `json:"repo"`
	BranchPRs map[string][]*OpenPullRequest `json:"branch_prs"`
}

// ListOpenPRsGithubOutput is the slimmed-down shape written to
// $GITHUB_OUTPUT — a map of branch name to raw PR list, matching the
// open_prs_map output of the existing workflow.
type ListOpenPRsGithubOutput struct {
	BranchPRs map[string][]*OpenPullRequest `json:"branch_prs"`
}

// Run fetches open pull requests for each requested branch.
func (r *ListOpenPRsReq) Run(
	ctx context.Context,
	client *gh.Client,
) (*ListOpenPRsRes, error) {
	if err := r.validate(ctx); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	branches, err := r.resolveBranches()
	if err != nil {
		return nil, fmt.Errorf("resolving branches: %w", err)
	}

	slog.Default().DebugContext(
		ctx, "listing open pull requests",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.Any("branches", branches),
		slog.Bool("include-drafts", r.IncludeDrafts),
	)

	res := &ListOpenPRsRes{
		Owner:     r.Owner,
		Repo:      r.Repo,
		BranchPRs: make(map[string][]*OpenPullRequest, len(branches)),
	}

	for _, branch := range branches {
		prs, err := listOpenPRs(ctx, client, r.Owner, r.Repo, branch, r.IncludeDrafts)
		if err != nil {
			return nil, fmt.Errorf("listing open pull requests for branch %q: %w", branch, err)
		}
		res.BranchPRs[branch] = prs
	}

	return res, nil
}

// validate ensures the request is well-formed before any API calls are made.
func (r *ListOpenPRsReq) validate(ctx context.Context) error {
	if r == nil {
		return errors.New("uninitialized request")
	}

	if r.Owner == "" {
		return errors.New("no github organization has been provided")
	}

	if r.Repo == "" {
		return errors.New("no github repository has been provided")
	}

	// Branches can be empty when Releases is provided for automatic resolution.
	if len(r.Branches) == 0 && r.Releases == nil {
		return errors.New("either branches or a versions.hcl releases configuration is required")
	}

	if r.Releases != nil {
		if err := r.Releases.Validate(ctx); err != nil {
			return fmt.Errorf("invalid releases configuration: %w", err)
		}
	}

	return nil
}

// resolveBranches returns the branches to query. When explicit branches were
// provided they are used as-is; otherwise all active release branches are
// derived from the versions.hcl configuration.
func (r *ListOpenPRsReq) resolveBranches() ([]string, error) {
	if len(r.Branches) > 0 {
		return r.Branches, nil
	}

	if r.Releases == nil || r.Releases.Config == nil || r.Releases.Config.ActiveVersion == nil {
		return nil, errors.New("no branches provided and no releases configuration available")
	}

	isEntRepo := releases.IsEnterpriseRepo(r.Repo)

	var branches []string
	for version := range r.Releases.Config.ActiveVersion.Versions {
		var branch string
		if isEntRepo {
			branch = releases.EnterpriseReleaseBranchForVersion(version)
		} else {
			branch = releases.CEReleaseBranchForVersion(version, "")
		}
		branches = append(branches, branch)
	}

	return branches, nil
}

// ToJSON marshals the response to JSON.
func (r *ListOpenPRsRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling list open pull requests to JSON: %w", err)
	}

	return b, nil
}

// ToTable renders the response as a plain-text table.
func (r *ListOpenPRsRes) ToTable() string {
	if r == nil {
		return ""
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	t.AppendHeader(table.Row{"branch", "pr", "author", "milestone", "backports"})

	for _, branch := range r.sortedBranches() {
		prs := r.BranchPRs[branch]
		if len(prs) == 0 {
			t.AppendRow(table.Row{branch, "no open PRs", "", "", ""})
			continue
		}
		for _, pr := range prs {
			t.AppendRow(table.Row{
				branch,
				pr.URL,
				pr.Author,
				pr.Milestone,
				strings.Join(pr.Backports, " "),
			})
		}
	}

	return t.Render()
}

// ToMarkdown renders the response as a Markdown table.
func (r *ListOpenPRsRes) ToMarkdown() string {
	if r == nil {
		return ""
	}

	t := table.NewWriter()
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	t.AppendHeader(table.Row{"branch", "pr", "author", "milestone", "backports"})

	for _, branch := range r.sortedBranches() {
		prs := r.BranchPRs[branch]
		if len(prs) == 0 {
			t.AppendRow(table.Row{branch, "no open PRs", "", "", ""})
			continue
		}
		for _, pr := range prs {
			t.AppendRow(table.Row{
				branch,
				pr.URL,
				pr.Author,
				pr.Milestone,
				strings.Join(pr.Backports, " "),
			})
		}
	}

	return escapeMarkdownTableColumns(t, 1, 2, 3, 4, 5).RenderMarkdown()
}

// ToGithubOutput marshals the branch→PR map for writing to $GITHUB_OUTPUT.
func (r *ListOpenPRsRes) ToGithubOutput() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}

	out := &ListOpenPRsGithubOutput{BranchPRs: r.BranchPRs}

	b, err := json.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("marshaling list open pull requests GITHUB_OUTPUT to JSON: %w", err)
	}

	return b, nil
}

// sortedBranches returns the branch keys in alphabetical order for stable
// table output.
func (r *ListOpenPRsRes) sortedBranches() []string {
	branches := make([]string, 0, len(r.BranchPRs))
	for b := range r.BranchPRs {
		branches = append(branches, b)
	}
	sort.Strings(branches)
	return branches
}

// listOpenPRs fetches all open (optionally including draft) pull requests
// targeting the given base branch, paginating as needed.
func listOpenPRs(
	ctx context.Context,
	client *gh.Client,
	owner string,
	repo string,
	base string,
	includeDrafts bool,
) ([]*OpenPullRequest, error) {
	ctx = slogctx.Append(
		ctx,
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("base", base),
	)
	slog.Default().DebugContext(ctx, "fetching open pull requests")

	opts := &gh.PullRequestListOptions{
		State:       "open",
		Base:        base,
		ListOptions: gh.ListOptions{PerPage: PerPageMax},
	}

	var results []*OpenPullRequest
	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("fetching pull requests: %w", err)
		}

		for _, pr := range prs {
			if pr.GetDraft() && !includeDrafts {
				continue
			}
			results = append(results, toOpenPullRequest(pr, owner, repo))
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	slog.Default().DebugContext(
		ctx, "fetched open pull requests",
		slog.Int("count", len(results)),
	)

	return results, nil
}

// toOpenPullRequest converts a GitHub API pull request to our OpenPullRequest
// type, extracting backport branches from labels.
func toOpenPullRequest(pr *gh.PullRequest, owner, repo string) *OpenPullRequest {
	number := pr.GetNumber()

	result := &OpenPullRequest{
		Number:    number,
		URL:       fmt.Sprintf("https://github.com/%s/%s/pull/%d", owner, repo, number),
		Author:    pr.GetUser().GetLogin(),
		Milestone: pr.GetMilestone().GetTitle(),
		Backports: extractBackportBranches(pr.Labels),
	}

	return result
}

// extractBackportBranches returns the branch names encoded in any labels of
// the form "backport/<branch>".
func extractBackportBranches(labels []*gh.Label) []string {
	const prefix = "backport/"
	var branches []string
	for _, label := range labels {
		name := label.GetName()
		if strings.HasPrefix(name, prefix) {
			branches = append(branches, strings.TrimPrefix(name, prefix))
		}
	}
	return branches
}
