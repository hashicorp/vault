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
	"github.com/jedib0t/go-pretty/v6/table"
)

// ChangelogStrategy is how we make a tag's changelog entry.
type ChangelogStrategy string

// ChangelogStrategyGithubReleaseNotes prepends a "## <tag>" section with the
// tag's GitHub generated release notes.
const ChangelogStrategyGithubReleaseNotes ChangelogStrategy = "github-release-notes"

// ChangelogStrategies returns the supported changelog strategies.
func ChangelogStrategies() []string {
	return []string{string(ChangelogStrategyGithubReleaseNotes)}
}

// CreateChangelogPRReq is a request to open a pull request that adds an entry
// for a tag to a changelog file.
type CreateChangelogPRReq struct {
	Owner               string
	Repo                string
	Tag                 string
	Branch              string
	BaseBranch          string
	Path                string
	Strategy            ChangelogStrategy
	Reviewers           []string
	WriteToGithubOutput bool
}

// CreateChangelogPRRes is a CreateChangelogPRReq response.
type CreateChangelogPRRes struct {
	Tag     string `json:"tag"`
	URL     string `json:"url"`
	Created bool   `json:"created"`
}

// changelogEditor makes the changelog entry for a ChangelogStrategy.
type changelogEditor interface {
	// hasEntry determines if the changelog content already has the entry.
	hasEntry(content string) bool
	// newEntry makes the entry.
	newEntry(ctx context.Context) (string, error)
	// addEntry returns the changelog content with the entry added.
	addEntry(content, entry string) string
}

// Run opens the changelog pull request unless the base branch already has an
// entry for the tag or the pull request is already open.
func (r *CreateChangelogPRReq) Run(ctx context.Context, client *gh.Client) (*CreateChangelogPRRes, error) {
	if err := r.validate(client); err != nil {
		return nil, err
	}
	editor, err := r.newChangelogEditor(client)
	if err != nil {
		return nil, err
	}
	res := &CreateChangelogPRRes{Tag: r.Tag}

	baseSHA, err := getBranchHeadSHA(ctx, client, r.Owner, r.Repo, r.BaseBranch)
	if err != nil {
		return nil, err
	}
	_, content, err := getFileContent(ctx, client, r.Owner, r.Repo, r.Path, baseSHA)
	if err != nil {
		return nil, err
	}
	if editor.hasEntry(content) {
		slog.Default().DebugContext(ctx, "base branch already has a changelog entry for the tag",
			slog.String("base_branch", r.BaseBranch),
			slog.String("tag", r.Tag),
		)
		return res, nil
	}

	// Make the entry before we create the branch so that a failure, like a
	// missing tag, doesn't leave an empty branch behind.
	entry, err := editor.newEntry(ctx)
	if err != nil {
		return nil, err
	}
	if err := ensureBranch(ctx, client, r.Owner, r.Repo, r.Branch, baseSHA); err != nil {
		return nil, err
	}
	if err := r.updateChangelog(ctx, client, editor, entry); err != nil {
		return nil, err
	}

	return r.ensurePullRequest(ctx, client, res)
}

func (r *CreateChangelogPRReq) validate(client *gh.Client) error {
	switch {
	case r == nil:
		return errors.New("uninitialized request")
	case client == nil:
		return errors.New("GitHub client is required")
	case strings.TrimSpace(r.Owner) == "":
		return errors.New("owner is required")
	case strings.TrimSpace(r.Repo) == "":
		return errors.New("repo is required")
	case strings.TrimSpace(r.Tag) == "":
		return errors.New("tag is required")
	case strings.TrimSpace(r.Branch) == "":
		return errors.New("branch is required")
	case strings.TrimSpace(r.BaseBranch) == "":
		return errors.New("base branch is required")
	case r.Branch == r.BaseBranch:
		return errors.New("branch and base branch must be different")
	case strings.TrimSpace(r.Path) == "":
		return errors.New("path is required")
	}

	return nil
}

func (r *CreateChangelogPRReq) newChangelogEditor(client *gh.Client) (changelogEditor, error) {
	switch r.Strategy {
	case ChangelogStrategyGithubReleaseNotes:
		return &githubReleaseNotesEditor{client: client, owner: r.Owner, repo: r.Repo, tag: r.Tag}, nil
	default:
		return nil, fmt.Errorf("unsupported changelog strategy %q, expected one of: %s",
			r.Strategy, strings.Join(ChangelogStrategies(), ", "),
		)
	}
}

// updateChangelog adds the entry to the changelog on our branch unless an
// earlier run already did.
func (r *CreateChangelogPRReq) updateChangelog(
	ctx context.Context,
	client *gh.Client,
	editor changelogEditor,
	entry string,
) error {
	file, content, err := getFileContent(ctx, client, r.Owner, r.Repo, r.Path, r.Branch)
	if err != nil {
		return err
	}
	if editor.hasEntry(content) {
		return nil
	}

	slog.Default().DebugContext(ctx, "updating changelog",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.String("branch", r.Branch),
		slog.String("path", r.Path),
	)
	_, _, err = client.Repositories.UpdateFile(ctx, r.Owner, r.Repo, r.Path, &gh.RepositoryContentFileOptions{
		Message: gh.Ptr("docs: update changelog for " + r.Tag),
		Content: []byte(editor.addEntry(content, entry)),
		SHA:     file.SHA,
		Branch:  gh.Ptr(r.Branch),
	})
	if err != nil {
		return fmt.Errorf("update %s on %s: %w", r.Path, r.Branch, err)
	}

	return nil
}

// ensurePullRequest returns the open pull request for our branch, or opens one.
func (r *CreateChangelogPRReq) ensurePullRequest(
	ctx context.Context,
	client *gh.Client,
	res *CreateChangelogPRRes,
) (*CreateChangelogPRRes, error) {
	slog.Default().DebugContext(ctx, "checking for changelog pull request",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.String("branch", r.Branch),
		slog.String("base_branch", r.BaseBranch),
	)
	pulls, _, err := client.PullRequests.List(ctx, r.Owner, r.Repo, &gh.PullRequestListOptions{
		Head:  r.Owner + ":" + r.Branch,
		Base:  r.BaseBranch,
		State: "open",
	})
	if err != nil {
		return nil, fmt.Errorf("list pull requests for %s: %w", r.Branch, err)
	}
	if len(pulls) > 0 {
		res.URL = pulls[0].GetHTMLURL()
		return res, nil
	}

	slog.Default().DebugContext(ctx, "creating changelog pull request",
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.String("branch", r.Branch),
		slog.String("base_branch", r.BaseBranch),
	)
	pull, _, err := client.PullRequests.Create(ctx, r.Owner, r.Repo, &gh.NewPullRequest{
		Title: gh.Ptr("[CHANGELOG] Prepare for " + r.Tag + " release"),
		Head:  gh.Ptr(r.Branch),
		Base:  gh.Ptr(r.BaseBranch),
		Body: gh.Ptr(
			"Changelog for " + r.Tag + ".\n\nGenerated by the pipeline github create changelog-pr command.",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("create changelog pull request: %w", err)
	}
	res.URL = pull.GetHTMLURL()
	res.Created = true

	if err := addReviewers(ctx, client, r.Owner, r.Repo, pull.GetNumber(), r.Reviewers); err != nil {
		return nil, fmt.Errorf("request reviews on %s: %w", res.URL, err)
	}

	return res, nil
}

// ToJSON marshals the response to JSON.
func (r *CreateChangelogPRRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling create changelog pull request response to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CreateChangelogPRRes) ToTable() string {
	if r == nil {
		return ""
	}
	t := table.NewWriter()
	t.Style().Options = table.Options{}
	t.AppendHeader(table.Row{"Tag", "URL", "Status"})
	t.AppendRow(table.Row{r.Tag, r.URL, createStatus(r.Created)})
	t.SuppressTrailingSpaces()

	return t.Render()
}

// githubReleaseNotesEditor implements ChangelogStrategyGithubReleaseNotes.
type githubReleaseNotesEditor struct {
	client *gh.Client
	owner  string
	repo   string
	tag    string
}

// hasEntry determines if a changelog has a "## <tag>" heading.
func (e *githubReleaseNotesEditor) hasEntry(content string) bool {
	return strings.Contains("\n"+content, "\n## "+e.tag+"\n")
}

func (e *githubReleaseNotesEditor) newEntry(ctx context.Context) (string, error) {
	notes, err := generateReleaseNotes(ctx, e.client, e.owner, e.repo, e.tag)
	if err != nil {
		return "", err
	}

	return "## " + e.tag + "\n\n" + notes, nil
}

func (e *githubReleaseNotesEditor) addEntry(content, entry string) string {
	return entry + "\n\n" + content
}

// getFileContent returns a file and its decoded content at a ref.
func getFileContent(
	ctx context.Context,
	client *gh.Client,
	owner string,
	repo string,
	path string,
	ref string,
) (*gh.RepositoryContent, string, error) {
	slog.Default().DebugContext(ctx, "reading file",
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("path", path),
		slog.String("ref", ref),
	)
	file, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, &gh.RepositoryContentGetOptions{
		Ref: ref,
	})
	if err != nil {
		return nil, "", fmt.Errorf("get %s on %s: %w", path, ref, err)
	}
	if file == nil {
		return nil, "", fmt.Errorf("%s on %s is not a file", path, ref)
	}
	content, err := file.GetContent()
	if err != nil {
		return nil, "", fmt.Errorf("decode %s on %s: %w", path, ref, err)
	}

	return file, content, nil
}

// ensureBranch creates the branch at sha unless it already exists.
func ensureBranch(ctx context.Context, client *gh.Client, owner, repo, branch, sha string) error {
	slog.Default().DebugContext(ctx, "checking for branch",
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("branch", branch),
	)
	_, _, err := client.Git.GetRef(ctx, owner, repo, "heads/"+branch)
	if err == nil {
		return nil
	}
	if !isNotFound(err) {
		return fmt.Errorf("get branch %s: %w", branch, err)
	}

	slog.Default().DebugContext(ctx, "creating branch",
		slog.String("owner", owner),
		slog.String("repo", repo),
		slog.String("branch", branch),
		slog.String("sha", sha),
	)
	if _, _, err := client.Git.CreateRef(ctx, owner, repo, gh.CreateRef{
		Ref: "refs/heads/" + branch,
		SHA: sha,
	}); err != nil {
		return fmt.Errorf("create branch %s: %w", branch, err)
	}

	return nil
}
