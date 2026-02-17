// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	libgithub "github.com/google/go-github/v81/github"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shurcooL/githubv4"
	slogctx "github.com/veqryn/slog-context"
)

// CloseCopiedOriginPullRequestReq is a request to copy a pull request from the CE repo to
// the Ent repo.
type CloseCopiedOriginPullRequestReq struct {
	Owner      string
	Repo       string
	PullNumber uint
}

// CloseCopiedOriginPullRequestRes is a copy pull request response.
type CloseCopiedOriginPullRequestRes struct {
	CopiedClosingIssues []*ClosingIssueRef      `json:"copied_closing_issues,omitempty"`
	CopiedComment       *libgithub.IssueComment `json:"copy_comment,omitempty"`
	CopiedPullRequest   *libgithub.PullRequest  `json:"copy_pull_request,omitempty"`
	OriginClosingIssues []*ClosingIssueRef      `json:"origin_associated_issues,omitempty"`
	OriginComment       *libgithub.IssueComment `json:"origin_comment,omitempty"`
	OriginPullRequest   *libgithub.PullRequest  `json:"origin_pull_request,omitempty"`
}

// ClosingIssueRefs represents our Github GraphQL query for finding issues
// associated with our Pull Request that ought to be automatically closed.
//
// The raw query looks something like:
//
//		{
//		  repository(owner: $owner, name: $repo) {
//		    pullRequest(number: $number) {
//		      repository {
//		        nameWithOwner
//		      }
//		      number
//		      closingIssuesReferences(first: 100) {
//		        edges {
//		          node {
//		            url
//	              number
//		            title
//		            closed
//		            repository {
//		              nameWithOwner
//		            }
//		          }
//		        }
//		      }
//		    }
//		  }
//		}
type ClosingIssueRefs struct {
	Repository struct {
		PullRequest struct {
			Repository struct {
				NameWithOwner string `json:"name_with_owner,omitempty"`
			} `json:"repository"`
			Number                  int `json:"number,omitempty"`
			ClosingIssuesReferences struct {
				Edges []struct {
					Node *ClosingIssueRef `json:"node"`
				} `json:"edges"`
			} `json:"closing_issues_references" graphql:"closingIssuesReferences(first: 100)"`
		} `json:"pull_request" graphql:"pullRequest(number: $number)"`
	} `json:"repository" graphql:"repository(owner: $owner, name: $repo)"`
}

// ClosingIssueRef is an issue that is associated with a pull request.
type ClosingIssueRef struct {
	URL        string `json:"url,omitempty" graphql:"url"`
	Number     int    `json:"number,omitempty"`
	Title      string `json:"title,omitempty"`
	Closed     bool   `json:"closed,omitempty"`
	Repository struct {
		NameWithOwner string `json:"name_with_owner,omitempty"`
	} `json:"repository"`
}

// Run runs the request to copy a pull request from the CE repo to the Ent repo.
func (r *CloseCopiedOriginPullRequestReq) Run(
	ctx context.Context,
	githubV3 *libgithub.Client,
	githubV4 *githubv4.Client,
) (*CloseCopiedOriginPullRequestRes, error) {
	slog.Default().DebugContext(slogctx.Append(ctx,
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.Uint64("pull-number", uint64(r.PullNumber)),
	), "closing copied pull request")

	res := &CloseCopiedOriginPullRequestRes{
		OriginClosingIssues: []*ClosingIssueRef{},
		CopiedClosingIssues: []*ClosingIssueRef{},
	}
	var err error
	originOwner, originRepo := "", ""
	var originNumber uint = 0

	// Whenever possible we try to update base pull request with a status update
	// on how the copying has gone.
	createComments := func() {
		// Make sure we return a response even if we fail
		if res == nil {
			res = &CloseCopiedOriginPullRequestRes{
				OriginClosingIssues: []*ClosingIssueRef{},
				CopiedClosingIssues: []*ClosingIssueRef{},
			}
		}

		var err1, err2 error
		res.CopiedComment, err1 = createPullRequestComment(
			ctx,
			githubV3,
			r.Owner,
			r.Repo,
			int(r.PullNumber),
			res.copiedCommentBody(err),
		)

		res.OriginComment, err2 = createPullRequestComment(
			ctx,
			githubV3,
			originOwner,
			originRepo,
			int(originNumber),
			res.originCommentBody(err),
		)

		// Set our finalized error on our response and also update our returned error
		err = errors.Join(err, err1, err2)
	}
	defer createComments()

	err = r.validate()
	if err != nil {
		return res, err
	}

	// Get the pull details of the copied PR
	res.CopiedPullRequest, err = getPullRequest(ctx, githubV3, r.Owner, r.Repo, int(r.PullNumber))
	if err != nil {
		return res, err
	}

	// Determine the origin PR from the copied PR branch name
	originOwner, originRepo, originNumber, _, err = decodeCopyPullRequestBranch(res.CopiedPullRequest.GetHead().GetRef())
	if err != nil {
		return res, err
	}
	slog.Default().DebugContext(slogctx.Append(ctx,
		slog.String("origin-owner", originOwner),
		slog.String("origin-repo", originRepo),
		slog.Uint64("origin-pull-number", uint64(originNumber)),
	), "decoded origin pull request information from copied pull request")

	// Get the pull details of the origin PR
	res.OriginPullRequest, err = getPullRequest(ctx, githubV3, originOwner, originRepo, int(originNumber))
	if err != nil {
		return res, err
	}

	// Close the origin PR if it's not closed already
	if res.OriginPullRequest.GetState() != "closed" {
		err = closePullRequest(ctx, githubV3, originOwner, originRepo, int(originNumber))
		if err != nil {
			return res, fmt.Errorf("unable to close origin pull request: %w", err)
		}
	} else {
		slog.Default().DebugContext(ctx, "origin pull request has already been closed")
	}

	// Close any issues associated with either the origin or copied PRs
	res.OriginClosingIssues, err = listPullRequestClosingIssues(
		ctx, githubV4, originOwner, originRepo, int(originNumber),
	)
	if err != nil {
		return res, err
	}

	res.CopiedClosingIssues, err = listPullRequestClosingIssues(
		ctx, githubV4, r.Owner, r.Repo, int(r.PullNumber),
	)
	if err != nil {
		return res, err
	}

	slog.Default().DebugContext(ctx, "closing any open associated issues")
	for _, issue := range slices.CompactFunc(
		append(res.OriginClosingIssues, res.CopiedClosingIssues...),
		closingIssueRefEqual,
	) {
		if issue.Closed {
			slog.Default().DebugContext(slogctx.Append(ctx,
				slog.String("repository", issue.Repository.NameWithOwner),
				slog.Int("number", issue.Number),
			), "associated issue is already closed")
			continue
		}

		nwo := issue.Repository.NameWithOwner
		parts := strings.SplitN(nwo, "/", 2)
		if len(parts) != 2 {
			return res, fmt.Errorf("could not determine repo and owner from associated issue %d, got: %s",
				issue.Number,
				nwo,
			)
		}

		err = closeIssue(ctx, githubV3, parts[0], parts[1], issue.Number)
		if err != nil {
			return res, fmt.Errorf("unable to close associated issue: %w", err)
		}
	}

	return res, nil
}

// validate ensures that we've been given all required fields necessary to
// perform the request.
func (r *CloseCopiedOriginPullRequestReq) validate() error {
	if r == nil {
		return errors.New("failed to initialize request")
	}

	if r.Owner == "" {
		return errors.New("no github owner has been provided")
	}

	if r.Repo == "" {
		return errors.New("no github repository has been provided")
	}

	if r.PullNumber == 0 {
		return errors.New("no github pull request number has been provided")
	}

	return nil
}

// copiedCommentBody is the markdown comment body that we'll attempt to set on
// the copied pull request.
func (r *CloseCopiedOriginPullRequestRes) copiedCommentBody(err error) string {
	if r == nil {
		return "no close copied origin pull request response has been initialized"
	}

	t := r.ToTable(err)
	if err == nil {
		t.SetTitle("Origin pull request and all associated issues have been closed!")
		return t.RenderMarkdown()
	}

	if t.Length() == 0 {
		// If we don't have any rows in our table then there's no need to render a
		// table so we'll just return an error
		return "## Closing origin pull request failed!\n\nError: " + err.Error()
	}

	// Render out our table but put the error message in the caption
	t.SetTitle("Closing origin pull request failed!")
	// Set the caption to the top-level error only as any attempt errors are
	// nested in the table.
	t.SetCaption("Error: " + err.Error())

	return t.RenderMarkdown()
}

// originCommentBody is the markdown comment body that we'll attempt to set on
// the origin pull request.
func (r *CloseCopiedOriginPullRequestRes) originCommentBody(err error) string {
	if r == nil {
		return "no close copied origin pull request response has been initialized"
	}

	t := r.ToTable(err)
	if err == nil {
		t.SetTitle("Copied pull request has been merged!")
		return t.RenderMarkdown()
	}

	if t.Length() == 0 {
		// If we don't have any rows in our table then there's no need to render a
		// table so we'll just return an error
		return "## Copied pull request has been merged!\n\nError: " + err.Error()
	}

	// Render out our table but put the error message in the caption
	t.SetTitle("Copy pull request failed!")
	// Set the caption to the top-level error only as any attempt errors are
	// nested in the table.
	t.SetCaption("Error: " + err.Error())

	return t.RenderMarkdown()
}

// ToJSON marshals the response to JSON.
func (r *CloseCopiedOriginPullRequestRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling list changed files to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CloseCopiedOriginPullRequestRes) ToTable(err error) table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{
		"Origin Pull Request", "Copied Pull Request", "Commit SHA", "Error",
	})
	row := table.Row{
		r.OriginPullRequest.GetHTMLURL(),
		r.CopiedPullRequest.GetHTMLURL(),
		fmt.Sprintf(
			"https://github.com/%s/commit/%s",
			r.CopiedPullRequest.GetHead().GetRepo().GetFullName(),
			r.CopiedPullRequest.GetMergeCommitSHA(),
		),
	}
	if err != nil {
		row = append(row, err.Error())
	}

	t.AppendRow(row)

	closedIssues := slices.CompactFunc(
		append(r.OriginClosingIssues, r.CopiedClosingIssues...),
		closingIssueRefEqual,
	)
	if len(closedIssues) > 0 {
		urls := []string{}
		for _, ci := range closedIssues {
			urls = append(urls, ci.URL)
		}
		t.AppendFooter(table.Row{
			"", "Closed Issues", strings.Join(urls, "\n"),
		})
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t
}
