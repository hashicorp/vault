// Copyright IBM Corp. 2016, 2025
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

	libgithub "github.com/google/go-github/v81/github"
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
	Error             error                         `json:"error,omitempty"`
	Request           *CopyPullRequestReq           `json:"request,omitempty"`
	OriginPullRequest *libgithub.PullRequest        `json:"origin_pull_request,omitempty"`
	PullRequest       *libgithub.PullRequest        `json:"pull_request,omitempty"`
	Comment           *libgithub.IssueComment       `json:"comment,omitempty"`
	SkippedCommits    []*libgithub.RepositoryCommit `json:"skipped_commits,omitempty"`
}

// Run runs the request to copy a pull request from the CE repo to the Ent repo.
func (r *CopyPullRequestReq) Run(
	ctx context.Context,
	github *libgithub.Client,
	git *libgit.Client,
) (*CopyPullRequestRes, error) {
	var err error
	res := &CopyPullRequestRes{
		Request:        r,
		SkippedCommits: []*libgithub.RepositoryCommit{},
	}

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

		// Set any known errors on the response before we create a comment, as the
		// error will be used in the comment body if present.
		err = errors.Join(err, os.Chdir(initialDir))
		var err1 error
		res.Comment, err1 = createPullRequestComment(
			ctx,
			github,
			r.FromOwner,
			r.FromRepo,
			int(r.PullNumber),
			res.CommentBody(err),
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
		defer os.RemoveAll(r.RepoDir)
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
	// pull request was created against. This will change our working directory
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

	// Add our from upstream as a remote and fetch our PR branch
	slog.Default().DebugContext(ctx, "adding CE upstream and fetching PR branch")
	remoteRes, err := git.Remote(ctx, &libgit.RemoteOpts{
		Command: libgit.RemoteCommandAdd,
		Track:   []string{prBranch},
		Fetch:   true,
		Name:    r.FromOrigin,
		URL:     res.OriginPullRequest.GetHead().GetRepo().GetCloneURL(),
	})
	if err != nil {
		err = fmt.Errorf("fetching target branch base ref: %s, %w", remoteRes.String(), err)
		return res, err
	}

	// Create a new branch for our copied changes. Encode the details of our origin
	// pull request into the branch name so that future post-merge operations can
	// determine the origin PR using only the branch name.
	branchName := encodeCopyPullRequestBranch(r.FromOwner, r.FromRepo, r.PullNumber, prBranch)
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

	// Get a list of commits we're going to cherry-pick into our new branch.
	commits, err := listPullRequestCommits(ctx, github, r.FromOwner, r.FromRepo, int(r.PullNumber))
	if err != nil {
		return res, err
	}

	// Generate an empty commit message that includes 'Co-Authored-By:' trailers
	// in the commit message. As we always squash all commits into a single merge
	// commit this helps to retain attribution for our source author(s).
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
		return res, fmt.Errorf("creating copy attribution commit message: %w", err)
	}
	defer func() {
		commitMessageFile.Close()
		_ = os.Remove(commitMessageFile.Name())
	}()

	attrCommitRes, err := git.Commit(ctx, &libgit.CommitOpts{
		AllowEmpty: true,
		File:       commitMessageFile.Name(),
		NoVerify:   true,
		NoEdit:     true,
	})
	if err != nil {
		return res, fmt.Errorf("committing attribution: %s: %w", attrCommitRes.String(), err)
	}

	slog.Default().DebugContext(ctx, "cherry-picking CE PR branch commits into new copy branch")
	var cherryPickErr error
	var cherryPickRes *libgit.ExecResponse
	for _, commit := range commits {
		// We only want to cherry-pick non-merge commits. To determine that we'll
		// see if the commit has more than one parent and skip it if it does.
		cherryPickRes, cherryPickErr = git.Show(ctx, &libgit.ShowOpts{
			Format: "%ph",
			Quiet:  true,
			Object: commit.GetSHA(),
		})
		if cherryPickErr != nil {
			break
		}

		parents := strings.TrimSpace(string(cherryPickRes.Stdout))
		if len(strings.Split(parents, " ")) > 1 {
			slog.Default().DebugContext(slogctx.Append(ctx,
				slog.String("sha", commit.GetSHA()),
				slog.String("parents", parents),
			), "skipping merge commit")

			res.SkippedCommits = append(res.SkippedCommits, commit)

			continue
		}

		cherryPickRes, cherryPickErr = git.CherryPick(ctx, &libgit.CherryPickOpts{
			FF:       true,
			Empty:    libgit.EmptyCommitKeep,
			Commit:   commit.GetSHA(),
			Strategy: libgit.MergeStrategyORT,
			StrategyOptions: []libgit.MergeStrategyOption{
				libgit.MergeStrategyOptionTheirs,
				libgit.MergeStrategyOptionIgnoreSpaceChange,
			},
		})
		if cherryPickErr != nil {
			break
		}
	}

	if cherryPickErr != nil {
		cherryPickErr = fmt.Errorf(
			"cherry-picking CE PR branch commits into new copy branch: %s: %w",
			cherryPickRes.String(),
			cherryPickErr,
		)
	}

	// If our merge failed we still want to create a pull request for our
	// failed copy so that a manual fix can be performed.
	if cherryPickErr != nil {
		err := resetAndCreateNOOPCommit(ctx, git, baseBranch)
		if err != nil {
			err = errors.Join(cherryPickErr, err)

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
		return res, errors.Join(cherryPickErr, err)
	}

	prTitle := fmt.Sprintf("Copy %s into %s", res.OriginPullRequest.GetTitle(), baseRef)
	prBody, err := renderEmbeddedTemplate("copy-pr-message.tmpl", struct {
		Error             error
		OriginPullRequest *libgithub.PullRequest
		TargetRef         string
	}{
		cherryPickErr,
		res.OriginPullRequest,
		baseRef,
	})
	if err != nil {
		err = fmt.Errorf("creating copy pull request body %w", err)
		return res, errors.Join(cherryPickErr, err)
	}
	limitedPRBody := limitCharacters(prBody)

	res.PullRequest, _, err = github.PullRequests.Create(
		ctx, r.ToOwner, r.ToRepo, &libgithub.NewPullRequest{
			Title:    &prTitle,
			Head:     &branchName,
			HeadRepo: &r.ToRepo,
			Base:     &baseRef,
			Body:     &limitedPRBody,
		},
	)
	if err != nil {
		err = fmt.Errorf("creating copy pull request %w", err)
		return res, errors.Join(cherryPickErr, err)
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
		return res, errors.Join(cherryPickErr, err)
	}

	return res, nil
}

// Validate ensures that we've been given the minimum filter arguments necessary to complete a
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
		return errors.New("no github from repository has been provided")
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
		"From", "To", "Skipped Merge Commits", "Error",
	})

	row := table.Row{nil, nil}
	if r.Request != nil {
		from := r.Request.FromOwner + "/" + r.Request.FromRepo
		if pr := r.OriginPullRequest; pr != nil {
			from = fmt.Sprintf("[%s#%d](%s)", from, pr.GetNumber(), pr.GetHTMLURL())
		}
		to := r.Request.ToOwner + "/" + r.Request.ToRepo
		if pr := r.PullRequest; pr != nil {
			to = fmt.Sprintf("[%s#%d](%s)", to, pr.GetNumber(), pr.GetHTMLURL())
		}
		row = table.Row{from, to, len(r.SkippedCommits)}
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
