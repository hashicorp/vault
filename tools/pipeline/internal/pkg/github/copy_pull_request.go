// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	libgithub "github.com/google/go-github/v68/github"
	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	"github.com/jedib0t/go-pretty/v6/table"
	slogctx "github.com/veqryn/slog-context"
)

// CopyPullRequestReq is a request to copy a pull request from the CE repo to
// the Ent repo.
type CopyPullRequestReq struct {
	FromOwner       string
	FromRepo        string
	FromOrigin      string
	ToOwner         string
	ToRepo          string
	ToOrigin        string
	PullNumber      uint
	RepoDir         string
	EntBranchSuffix string // add +ent to release/* branches
}

// CopyPullRequestRes is a copy pull request response.
type CopyPullRequestRes struct {
	Error             error                   `json:"error,omitempty"`
	Request           *CopyPullRequestReq     `json:"request,omitempty"`
	OriginPullRequest *libgithub.PullRequest  `json:"origin_pull_request,omitempty"`
	PullRequest       *libgithub.PullRequest  `json:"pull_request,omitempty"`
	Comment           *libgithub.IssueComment `json:"comment,omitempty"`
}

// Run runs the request to copy a pull request from the CE repo to the Ent repo.
func (r *CopyPullRequestReq) Run(
	ctx context.Context,
	github *libgithub.Client,
	git *libgit.Client,
) (*CopyPullRequestRes, error) {
	var err error
	res := &CopyPullRequestRes{Request: r}

	slog.Default().DebugContext(slogctx.Append(ctx,
		slog.String("from-owner", r.FromOwner),
		slog.String("from-repo", r.FromRepo),
		slog.String("from-origin", r.FromOrigin),
		slog.String("to-owner", r.ToOwner),
		slog.String("to-repo", r.ToRepo),
		slog.String("to-origin", r.ToOrigin),
		slog.String("repo-dir", r.RepoDir),
		slog.Uint64("pull-number", uint64(r.PullNumber)),
		slog.String("ent-branch-suffix", r.EntBranchSuffix),
	), "copying pull request")

	initialDir, err := os.Getwd()
	if err != nil {
		return res, fmt.Errorf("getting current working directory: %w", err)
	}

	// Whenever possible we try to update base pull request with a status update
	// on how the copying has gone.
	createComment := func() {
		// Make sure we return a response even if we fail
		if res == nil {
			res = &CopyPullRequestRes{Request: r}
		}

		// Figure out the comment body. Worst case it ought to be whatever error
		// we've returned.
		var body string
		if err != nil {
			body = err.Error()
		}

		// Set any known errors on the response before we create a comment, as the
		// error will be used in the comment body if present.
		err = errors.Join(err, os.Chdir(initialDir))
		body = res.CommentBody(err)
		var err1 error
		res.Comment, err1 = createPullRequestComment(
			ctx,
			github,
			r.FromOwner,
			r.FromRepo,
			int(r.PullNumber),
			body,
		)

		// Set our finalized error on our response and also update our returned error
		err = errors.Join(err, err1)
	}
	defer createComment()

	// Make sure we have required and valid fields
	err = r.Validate(ctx)
	if err != nil {
		return res, err
	}

	// Make sure we've been given a valid location for a repo and/or create a
	// temporary one
	var tmpDir bool
	r.RepoDir, err, tmpDir = ensureGitRepoDir(ctx, r.RepoDir)
	if err != nil {
		return res, err
	}
	if tmpDir {
		// defer os.RemoveAll(r.RepoDir)
	}

	// Get our pull request details
	res.OriginPullRequest, err = getPullRequest(
		ctx, github, r.FromOwner, r.FromRepo, int(r.PullNumber),
	)
	if err != nil {
		return res, err
	}

	// Determine our pull request base ref. Handle the fact that enterprise
	// release branches contain the +ent suffix.
	baseRef := res.OriginPullRequest.GetBase().GetRef()
	if strings.HasPrefix(baseRef, "release/") {
		baseRef = baseRef + r.EntBranchSuffix
	}

	// Clone the remote repository and fetch the base ref, which is the branch our
	// pull request was created against. These will change our working directory
	// into RepoDir
	_, err = os.Stat(filepath.Join(r.RepoDir, ".git"))
	if err == nil {
		err = initializeExistingRepo(
			ctx, git, r.RepoDir, r.ToOrigin, baseRef,
		)
	} else {
		err = initializeNewRepo(
			ctx, git, r.RepoDir, r.ToOwner, r.ToRepo, r.ToOrigin, baseRef,
		)
	}
	if err != nil {
		return res, err
	}

	prBranch := res.OriginPullRequest.GetHead().GetRef()
	prBranchRef := "remotes/" + r.FromOrigin + "/" + prBranch

	// Add our from upstream as a remote and fetch our PR branch
	slog.Default().DebugContext(ctx, "adding CE upstream and fetching PR branch")
	remoteRes, err := git.Remote(ctx, &libgit.RemoteOpts{
		Command: libgit.RemoteCommandAdd,
		Track:   []string{prBranch},
		Fetch:   true,
		Name:    r.FromOrigin,
		URL:     fmt.Sprintf("https://github.com/%s/%s.git", r.FromOwner, r.FromRepo),
	})
	if err != nil {
		err = fmt.Errorf("fetching target branch base ref: %s, %w", remoteRes.String(), err)
		return res, err
	}

	// Create a new branch for our copied changes.
	branchName := r.copyBranchNameForRef(baseRef, prBranch)
	// We don't have local references so create a new branch from our tracking branch
	baseBranch := "remotes/" + r.ToOrigin + "/" + baseRef
	slog.Default().DebugContext(ctx, "checking out new copy branch")
	checkoutRes, err := git.Checkout(ctx, &libgit.CheckoutOpts{
		NewBranchForceCheckout: branchName, // -B
		Branch:                 baseBranch,
	})
	if err != nil {
		return res, fmt.Errorf("checking out new copy branch: %s: %w", checkoutRes.String(), err)
	}

	// Generate a merge commit message. While git is able to generate a nice merge
	// commit with a summary of all commit headers, we create our own that
	// includes 'Co-Authored-By:' trailers in the commit message. As we always
	// squash all commits into a single merge commit this helps to retain
	// attribution for our source author.
	commits, err := listPullRequestCommits(ctx, github, r.FromOwner, r.FromRepo, int(r.PullNumber))
	if err != nil {
		return res, err
	}

	commitMessageFile, err := renderEmbeddedTemplateToTmpFile("copy-pr-commit-message.tmpl", struct {
		CoAuthoredByTrailers []string
		OriginPullRequest    *libgithub.PullRequest
		TargetRef            string
	}{
		r.getCoAuthoredByTrailers(commits),
		res.OriginPullRequest,
		baseRef,
	})
	if err != nil {
		return res, fmt.Errorf("creating merge commit message: %w", err)
	}
	defer func() {
		commitMessageFile.Close()
		_ = os.Remove(commitMessageFile.Name())
	}()

	slog.Default().DebugContext(ctx, "merging CE PR branch into new copy branch")
	mergeRes, mergeErr := git.Merge(ctx, &libgit.MergeOpts{
		File:     commitMessageFile.Name(),
		NoVerify: true,
		Strategy: libgit.MergeStrategyORT,
		StrategyOptions: []libgit.MergeStrategyOption{
			libgit.MergeStrategyOptionTheirs,
			libgit.MergeStrategyOptionIgnoreSpaceChange,
		},
		IntoName: baseRef,
		Commit:   prBranchRef,
	})
	if mergeErr != nil {
		mergeErr = fmt.Errorf("merging CE PR branch into new copy branch: %s: %w", mergeRes.String(), mergeErr)
	}

	// If our merge failed we still want to create a pull request for our
	// failed copy so that a manual fix can be performed.
	if mergeErr != nil {
		err := resetAndCreateNOOPCommit(ctx, git, baseBranch)
		if err != nil {
			err = errors.Join(mergeErr, err)

			// Something wen't wrong trying to create our no-op commit. There's
			// nothing more we can do but return our error at this point.
			return res, err
		}
	}

	slog.Default().DebugContext(ctx, "pushing new branch to enterprise")
	pushRes, err := git.Push(ctx, &libgit.PushOpts{
		Repository: r.ToOrigin,
		Refspec:    []string{branchName},
	})
	if err != nil {
		err = fmt.Errorf("pushing copied branch: %s: %w", pushRes.String(), err)
		return res, errors.Join(mergeErr, err)
	}

	prTitle := fmt.Sprintf("Copy %s into %s", res.OriginPullRequest.GetTitle(), baseRef)
	prBody, err := renderEmbeddedTemplate("copy-pr-message.tmpl", struct {
		Error             error
		OriginPullRequest *libgithub.PullRequest
		TargetRef         string
	}{
		mergeErr,
		res.OriginPullRequest,
		baseRef,
	})
	if err != nil {
		err = fmt.Errorf("creating copy pull request body %w", err)
		return res, errors.Join(mergeErr, err)
	}

	res.PullRequest, _, err = github.PullRequests.Create(
		ctx, r.ToOwner, r.ToRepo, &libgithub.NewPullRequest{
			Title:    &prTitle,
			Head:     &branchName,
			HeadRepo: &r.ToRepo,
			Base:     &baseRef,
			Body:     &prBody,
		},
	)
	if err != nil {
		err = fmt.Errorf("creating copy pull request %w", err)
		return res, errors.Join(mergeErr, err)
	}

	// Assign the pull request to the actor that was assigned the original
	// pull request and anybody that approved it.
	reviews, err := listPullRequestReviews(ctx, github, r.FromOwner, r.FromRepo, int(r.PullNumber))
	if err != nil {
		return res, err
	}
	err = addAssignees(
		ctx,
		github,
		r.ToOwner,
		r.ToRepo,
		int(res.PullRequest.GetNumber()),
		append(r.getApproverLogins(reviews), res.OriginPullRequest.GetAssignee().GetLogin()),
	)
	if err != nil {
		err = fmt.Errorf("assigning ownership to copy pull request %w", err)
		return res, errors.Join(mergeErr, err)
	}

	return res, nil
}

// copyBranchNameForRef returns then branch name to use for our PR copy operation.
// e.g. copy/release/1.19.x+ent/my-feature-branch
func (r CopyPullRequestReq) copyBranchNameForRef(
	ref string,
	prBranch string,
) string {
	name := fmt.Sprintf("copy/%s/%s", ref, prBranch)
	if len(name) > 250 {
		// Handle Githubs branch name max length
		name = name[:250]
	}

	return name
}

// validate ensures that we've been given the minimum filter arguments necessary to complete a
// request. It is always recommended that additional fitlers be given to reduce the response size
// and not exhaust API limits.
func (r *CopyPullRequestReq) Validate(ctx context.Context) error {
	if r == nil {
		return errors.New("failed to initialize request")
	}

	if r.FromOrigin == "" {
		return errors.New("no github from origin has been provided")
	}

	if r.FromOwner == "" {
		return errors.New("no github from owner has been provided")
	}

	if r.FromRepo == "" {
		return errors.New("no github repository has been provided")
	}

	if r.ToOrigin == "" {
		return errors.New("no github to origin has been provided")
	}

	if r.ToOwner == "" {
		return errors.New("no github to owner has been provided")
	}

	if r.ToRepo == "" {
		return errors.New("no github to repository has been provided")
	}

	if r.PullNumber == 0 {
		return errors.New("no github pull request number has been provided")
	}

	return nil
}

// CommentBody is the markdown comment body that we'll attempt to set on the
// pull request
func (r *CopyPullRequestRes) CommentBody(err error) string {
	if r == nil {
		return "no copy pull request response has been initialized"
	}

	t := r.ToTable(err)
	if err == nil {
		t.SetTitle("Copy workflow completed!")
		return t.RenderMarkdown()
	}

	if t.Length() == 0 {
		// If we don't have any rows in our table then there's no need to render a
		// table so we'll just return an error
		return "## Copying pull request failed!\n\nError: " + err.Error()
	}

	// Render out our table but put the error message in the caption
	t.SetTitle("Copy pull request failed!")
	// Set the caption to the top-level error only as any attempt errors are
	// nested in the table.
	t.SetCaption("Error: " + err.Error())

	return t.RenderMarkdown()
}

// ToJSON marshals the response to JSON.
func (r *CopyPullRequestRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling list changed files to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CopyPullRequestRes) ToTable(err error) table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{
		"From", "To", "Error",
	})

	row := table.Row{nil, nil}
	if r.Request != nil {
		from := r.Request.FromOwner + "/" + r.Request.FromRepo
		if pr := r.OriginPullRequest; pr != nil {
			from = fmt.Sprintf("[%s#%d](%s)", from, pr.GetID(), pr.GetHTMLURL())
		}
		to := r.Request.ToOwner + "/" + r.Request.ToRepo
		if pr := r.PullRequest; pr != nil {
			to = fmt.Sprintf("[%s#%d](%s)", to, pr.GetID(), pr.GetHTMLURL())
		}
		row = table.Row{from, to}
	}
	if err != nil {
		row = append(row, err.Error())
	} else {
		row = append(row, nil)
	}
	t.AppendRow(row)

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t
}

func (r *CopyPullRequestReq) getCoAuthoredByTrailers(commits []*libgithub.RepositoryCommit) []string {
	if len(commits) < 1 {
		return nil
	}

	seen := map[string]struct{}{}
	trailers := []string{}

	for _, repoCommit := range commits {
		commit := repoCommit.GetCommit()
		if commit == nil {
			continue
		}
		author := commit.GetAuthor()
		if author == nil {
			continue
		}
		email := author.GetEmail()
		if email == "" {
			continue
		}
		if _, ok := seen[email]; ok {
			continue
		}
		seen[email] = struct{}{}
		trailers = append(trailers, fmt.Sprintf("Co-Authored-By: %s <%s>", author.GetName(), email))
	}

	return trailers
}

func (r *CopyPullRequestReq) getApproverLogins(reviews []*libgithub.PullRequestReview) []string {
	if len(reviews) < 1 {
		return nil
	}

	logins := map[string]struct{}{}
	for _, review := range reviews {
		if review == nil || review.State == nil || *review.State != "APPROVED" {
			continue
		}
		if login := review.GetUser().GetLogin(); login != "" {
			logins[login] = struct{}{}
		}
	}

	return slices.Sorted(maps.Keys(logins))
}
