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
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/jedib0t/go-pretty/v6/table"
	slogctx "github.com/veqryn/slog-context"
)

// CreateBackportReq is a request to create a backport pull request from another
// pull request. The request has been designed to work when triggered in a
// Github Actions workflow where the only required values are present in the
// github event context. That assumes a pull request event:
//
//	pull_request_target:
//	  types: closed
//
// The request ought to be guarded so as to nominally trigger only on merges:
//
//	if: github.even.pull_request.merged"
//
// See Run() for more details around how the request determines which branches
// to backport to, whether or not the backport commits need to be amended for
// excluded CE files, or whether or not the backport can be skipped entirely.
//
// NOTE: At this time the request only supports a single squashed merge commit.
type CreateBackportReq struct {
	// The Github Owner. E.g. "hashicorp"
	Owner string
	// The Github Repo. E.g. "vault-enterprise"
	Repo string
	// The Pull Request ID Number of the PR that we wish to backport.
	PullNumber uint
	// BaseOrigin is the name of the remote for the base ref of the pull request.
	// E.g. "origin".
	BaseOrigin string

	// The local directory where to clone the repository:
	//    https://github.com/<Owner>/<Repo>.git.
	// If the directory is configured it either must exist. When unset, a
	// temporary directory will be created and used automatically.
	RepoDir string

	// ReleaseVersionConfigPath is the path to .release/versions.hcl. We use this
	// file to determine which branches are active so that we can automatically
	// determine which origins to backport depending on the given tags.
	ReleaseVersionConfigPath string
	// ReleaseRecurseDepth defined how many directories back we're allowed to
	// scan to search for .release/versions.hcl. This is incompatible with
	// ReleaseVersionConfigPath.
	ReleaseRecurseDepth uint

	// CEExclude are changed files groups for files that ought to be excluded
	// when creating CE backports. E.g. ["enterprise"]
	CEExclude changed.FileGroups
	// CEBranchPrefix is the prefix used for CE branches. E.g. "ce"
	CEBranchPrefix string
	// CEAllowInactiveGroups are changed file groups for files that ought to be
	// allowed to be backported to inactive CE branches. Eg. ["docs", "pipeline"]
	CEAllowInactiveGroups changed.FileGroups

	// NOTE: The following fields are for testing purposes only and might be
	// removed after the cutover to the new workflow.

	// EntBranchPrefix is an ent branch prefix. This is only used for testing
	// before we migrate to the tool full time.
	EntBranchPrefix string

	// BackportLabelPrefix is the backport label prefix. E.g. "backport". This
	// should only be used for testing before the new workflow is active.
	BackportLabelPrefix string
}

// NewCreateBackportReqOpt is a functional option to set fields when calling
// NewCreateBackportPRReq()
type NewCreateBackportReqOpt func(*CreateBackportReq)

// CreateBackportRes is a respose of creating a backport pull request
type CreateBackportRes struct {
	OriginPullRequest *libgithub.PullRequest            `json:"origin_pull_request,omitempty"`
	Branch            string                            `json:"branch,omitempty"`
	Attempts          map[string]*CreateBackportAttempt `json:"attempts,omitempty"`
	Comment           *libgithub.IssueComment           `json:"comment,omitempty"`
	Error             error                             `json:"-"`
	// Use a separate field so we marshal the error message to a string value
	ErrorMessage string `json:"error,omitempty"`
}

// Labels are just a collection of github labels that we have created various
// helper functions for.
type Labels []*libgithub.Label

// CreateBackportAttempt is an attempt at creating a backport for target
// branch reference.
type CreateBackportAttempt struct {
	BaseRef       string                 `json:"base_ref,omitempty"`
	TargetRef     string                 `json:"target_ref,omitempty"`
	Error         error                  `json:"error,omitempty"`
	Skipped       bool                   `json:"skipped,omitempty"`
	SkippedReason string                 `json:"skipped_reason,omitempty"`
	PullRequest   *libgithub.PullRequest `json:"pull_request,omitempty"`
}

// NewCreateBackportReq takes variable options and returns a new
// CreateBackportPRReq.
func NewCreateBackportReq(opts ...NewCreateBackportReqOpt) *CreateBackportReq {
	req := &CreateBackportReq{
		Owner:               "hashicorp",
		Repo:                "vault-enterprise",
		ReleaseRecurseDepth: 3,
		CEExclude:           changed.FileGroups{changed.FileGroupEnterprise},
		CEBranchPrefix:      "ce",
		CEAllowInactiveGroups: changed.FileGroups{
			changed.FileGroupChangelog,
		},
		BaseOrigin:          "origin",
		BackportLabelPrefix: "backport",
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}

// WithCreateBackportReqOwner sets the Owner
func WithCreateBackportReqOwner(owner string) NewCreateBackportReqOpt {
	return func(req *CreateBackportReq) {
		req.Owner = owner
	}
}

// WithCreateBrackportReqRepo sets the Repo
func WithCreateBrackportReqRepo(repo string) NewCreateBackportReqOpt {
	return func(req *CreateBackportReq) {
		req.Repo = repo
	}
}

// WithCreateBrackportReqRepoDir sets the RepoDir
func WithCreateBrackportReqRepoDir(dir string) NewCreateBackportReqOpt {
	return func(req *CreateBackportReq) {
		req.RepoDir = dir
	}
}

// WithCreateBrackportReqPullNumber sets the PullNumber
func WithCreateBrackportReqPullNumber(number uint) NewCreateBackportReqOpt {
	return func(req *CreateBackportReq) {
		req.PullNumber = number
	}
}

// WithCreateBrackportReqBaseOrigin sets the BaseOrigin
func WithCreateBrackportReqBaseOrigin(origin string) NewCreateBackportReqOpt {
	return func(req *CreateBackportReq) {
		req.BaseOrigin = origin
	}
}

// WithCreateBrackportReqReleaseRecurseDepth sets the ReleaseRecurseDepth
func WithCreateBrackportReqReleaseRecurseDepth(depth uint) NewCreateBackportReqOpt {
	return func(req *CreateBackportReq) {
		req.ReleaseRecurseDepth = depth
	}
}

// WithCreateBrackportReqCEExclude sets the CEExclude
func WithCreateBrackportReqCEExclude(exclude changed.FileGroups) NewCreateBackportReqOpt {
	return func(req *CreateBackportReq) {
		req.CEExclude = exclude
	}
}

// WithCreateBrackportReqCEBranchPrefix sets the CEBranchPrefix
func WithCreateBrackportReqCEBranchPrefix(prefix string) NewCreateBackportReqOpt {
	return func(req *CreateBackportReq) {
		req.CEBranchPrefix = prefix
	}
}

// WithCreateBrackportReqAllowInactiveGroups sets the CEAllowInactiveGroups
func WithCreateBrackportReqAllowInactiveGroups(groups changed.FileGroups) NewCreateBackportReqOpt {
	return func(req *CreateBackportReq) {
		req.CEAllowInactiveGroups = groups
	}
}

// WithCreateBrackportReqEntBranchPrefix sets the EntBranchPrefix
func WithCreateBrackportReqEntBranchPrefix(prefix string) NewCreateBackportReqOpt {
	return func(req *CreateBackportReq) {
		req.EntBranchPrefix = prefix
	}
}

// WithCreateBrackportReqBackportLabelPrefix sets the BackportLabelPrefix
func WithCreateBrackportReqBackportLabelPrefix(prefix string) NewCreateBackportReqOpt {
	return func(req *CreateBackportReq) {
		req.BackportLabelPrefix = prefix
	}
}

// Run runs the backport request to create backports for every target branch
// as needed.
//
// If the base references is to an enteprise branch, that is, the base reference
// branch does not contain the CEBranchPrefix, then a backport to the
// corresponding CE branch is assumed and will be created.
//
// If the base reference is to a CE branch then backports are only created if
// there are backport labels present.
//
// Backport labels should be listed in the same schema as .release/versions.hcl:
// E.g. "release/1.19.x". The correct backport branches will be used depending
// on whether or not base branch of the PR is enteprise or CE.
//
// Enterprise branches will only ever backport to the corresponding ce branch
// and to other enterprise branches. When those enterprise branches are merged
// we'll create the CE backports.
//
// There are many factors to conside when backporting to a CE branch. The
// request will automatically inspect the changed files of a PR to determine
// if the PR contains non-enterprise files that need to be backported. In the
// event we've only changed enterprise files we'll skip the CE backport.
// If we've changed both enterprise and non-enterprise files the backport will
// automatically remove the enterprise files.
//
// We also factor in whether or not a CE branch is "active". If the branch is
// inactive we'll skip backporting unless the change includes docs, pipeline
// changes, or README changes. This allows docs authors to write docs against
// enteprise branches and have them backported without having to do it manually.
//
// We also do our best to update the source pull request with a comment that
// outlines each backport and its status.
//
// This request designed to always return a response, even if things go wrong.
// We will always attempt to run all backport references even if some fail.
// As such we don't return an error here but do embed them in the response for
// more control and precise handling. Callers should use Err() on the response
// to get a singular error, or they can inspect the Error field for each
// backport attempt.
func (r *CreateBackportReq) Run(
	ctx context.Context,
	github *libgithub.Client,
	git *libgit.Client,
) (res *CreateBackportRes) {
	res = &CreateBackportRes{Attempts: map[string]*CreateBackportAttempt{}}

	slog.Default().DebugContext(slogctx.Append(ctx,
		slog.String("owner", r.Owner),
		slog.String("repo", r.Repo),
		slog.String("repo-dir", r.RepoDir),
		slog.Uint64("pull-number", uint64(r.PullNumber)),
		slog.String("base-origin", r.BaseOrigin),
		slog.String("config-path", r.ReleaseVersionConfigPath),
		slog.Uint64("config-path-recurse-depth", uint64(r.ReleaseRecurseDepth)),
		slog.String("ce-branch-prefix", r.CEBranchPrefix),
		slog.String("ce-allow-inactive", strings.Join(r.CEAllowInactiveGroups.Groups(), ",")),
		slog.String("ce-exclude", strings.Join(r.CEExclude.Groups(), ",")),
		slog.String("ent-branch-prefix", r.EntBranchPrefix),
		slog.String("backport-label-prefix", r.BackportLabelPrefix),
	), "running create backport pr request")

	initialDir, err := os.Getwd()
	if err != nil {
		res.Error = fmt.Errorf("getting current working directory: %w", err)
		return res
	}

	// Whenever possible we try to update base pull request with a status update
	// on how the backporting has gone.
	defer func() {
		// Make sure we return a response even if we fail
		if res == nil {
			res = &CreateBackportRes{}
		}

		// Set any known errors on the response before we create a comment, as the
		// error will be used in the comment body if present.
		res.Error = errors.Join(res.Error, os.Chdir(initialDir))
		var err1 error
		res.Comment, err1 = createPullRequestComment(
			ctx, github, r.Owner, r.Repo, int(r.PullNumber), res.CommentBody(),
		)

		// Set our finalized error on our response and also update our returned error
		res.Error = errors.Join(res.Error, err1)
	}()

	// Make sure we have required and valid fields
	res.Error = r.Validate(ctx)
	if res.Error != nil {
		return res
	}

	// Make sure we've been given a valid location for a repo and/or create a
	// temporary one
	var tmpDir bool
	r.RepoDir, res.Error, tmpDir = ensureGitRepoDir(ctx, r.RepoDir)
	if res.Error != nil {
		return res
	}
	if tmpDir {
		defer os.RemoveAll(r.RepoDir)
	}

	// Get our pull request details
	res.OriginPullRequest, res.Error = getPullRequest(
		ctx, github, r.Owner, r.Repo, int(r.PullNumber),
	)
	if res.Error != nil {
		return res
	}

	// Make sure our PR is merged and has a merge SHA
	if !res.OriginPullRequest.GetMerged() {
		res.Error = errors.New("cannot backport unmerged PR")
		return res
	}
	if res.OriginPullRequest.GetMergeCommitSHA() == "" {
		res.Error = errors.New("no merge commit SHA is associated with the PR")
		return res
	}

	// Determine which CE branches are active. Do this before we change our
	// working directory since the path given could be relative to the original
	// path.
	var activeVersions map[string]*releases.Version
	activeVersions, res.Error = r.getActiveVersions(ctx)
	if res.Error != nil {
		return res
	}

	// Clone the remote repository and fetch the base ref, which is the branch our
	// pull request was created against. These will change our working directory
	// into RepoDir
	baseRef := res.OriginPullRequest.GetBase().GetRef()
	_, err = os.Stat(filepath.Join(r.RepoDir, ".git"))
	if err == nil {
		res.Error = initializeExistingRepo(
			ctx, git, r.RepoDir, r.BaseOrigin, baseRef,
		)
	} else {
		res.Error = initializeNewRepo(
			ctx, git, r.RepoDir, r.Owner, r.Repo, r.BaseOrigin, baseRef,
		)
	}
	if res.Error != nil {
		return res
	}

	// Get the list of changed files and determine if our PR modified any files
	// in CEExclude.
	var changedFiles *ListChangedFilesRes
	changedFiles, res.Error = r.getChangedFiles(ctx, github)
	if res.Error != nil {
		return res
	}

	// Determine base references we want to backport and create backports for each
	// reference. In cases where the reference starts with the CEBranchPrefix then
	// we'll remove any files that are in exclude groups.
	for _, ref := range r.determineBackportRefs(ctx, baseRef, res.OriginPullRequest.Labels) {
		res.Attempts[ref] = r.backportRef(
			ctx, git, github, res.OriginPullRequest, activeVersions, changedFiles, ref,
		)

		if attempt := res.Attempts[ref]; attempt != nil && attempt.Error != nil {
			// Something went wrong attempting to backport the reference. Reset our
			// repository to ensure that our next attempt does not start in a nasty
			// state.
			resetRes, err := git.Reset(ctx, &libgit.ResetOpts{
				Mode:    libgit.ResetModeHard,
				Treeish: fmt.Sprintf("%s/%s", r.BaseOrigin, baseRef),
			})
			if err != nil {
				res.Error = errors.Join(res.Error, fmt.Errorf(
					"resetting repository after failed attempt: %s: %w", resetRes.String(), err),
				)
				// If we can't reset the repository there's no point in trying further
				// attempts as we must assume something has gone horribly wrong.
				break
			}
		}
	}

	return res
}

// Validate validates the request to ensure that all required fields are present
func (r *CreateBackportReq) Validate(ctx context.Context) error {
	if r == nil {
		return fmt.Errorf("unitialized")
	}

	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("validating create backport pr requests: %w", err)
		}
	}()

	slog.Default().DebugContext(ctx, "validating create backport pr request")

	if r.Owner == "" {
		return errors.New("no github organization has been provided")
	}

	if r.Repo == "" {
		return errors.New("no github repository has been provided")
	}

	if r.BaseOrigin == "" {
		return errors.New("no base origin has been configued")
	}

	if r.PullNumber == 0 {
		return errors.New("no pull request number or commit SHA has been provided")
	}

	if r.CEBranchPrefix == "" {
		return errors.New("no ce branch prefix has been configured")
	}

	if r.CEExclude == nil {
		return errors.New("ce-exclude has not been initialized")
	}

	if r.CEAllowInactiveGroups == nil {
		return errors.New("ce inactive-allowed has not been initialized")
	}

	if r.BackportLabelPrefix == "" {
		return errors.New("no backport label prefix has been configured")
	}

	return nil
}

// AttemptErrors are any potential errors encountered during our backport attempts
func (r *CreateBackportRes) AttemptErrors() []error {
	if r == nil || len(r.Attempts) < 1 {
		return nil
	}

	errs := []error{}
	for _, k := range slices.Sorted(maps.Keys(r.Attempts)) {
		a := r.Attempts[k]
		if a.Error == nil {
			continue
		}
		errs = append(errs, a.Error)
	}

	return errs
}

// CommentBody is the markdown comment body that we'll attempt to set on the
// pull request
func (r *CreateBackportRes) CommentBody() string {
	if r == nil {
		return "no backport response has been initialized"
	}

	t := r.ToTable()
	err := r.Err()
	if err == nil {
		t.SetTitle("Backport workflow completed!")
		return t.RenderMarkdown()
	}

	if t.Length() == 0 {
		// If we don't have any rows in our table then we never made it far enough
		// to have attempts. As such, there's no need to render a table so we'll
		// just return an error
		return "## Backport workflow failed!\n\nError: " + err.Error()
	}

	// Render out our table but put the error message in the caption
	t.SetTitle("Backport workflow failed!")
	if r.Error != nil {
		// Set the caption to the top-level error only as any attempt errors are
		// nested in the table.
		t.SetCaption("Error: " + r.Error.Error())
	}

	return t.RenderMarkdown()
}

// Err returns a single combined error comprised of any issues that might have
// arisen during Run() but also that of any individual backport attempt.
func (r *CreateBackportRes) Err() error {
	if r == nil {
		return fmt.Errorf("uninitialized")
	}

	return errors.Join(r.Error, errors.Join(r.AttemptErrors()...))
}

// ToJSON marshals the response to JSON.
func (r *CreateBackportRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling create backport pr response to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CreateBackportRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{
		"Base Branch", "Target Branch", "URL", "Skipped Reason", "Error",
	})

	for _, version := range slices.Sorted(maps.Keys(r.Attempts)) {
		values := r.Attempts[version]
		row := table.Row{values.BaseRef, values.TargetRef}
		if values.PullRequest != nil {
			row = append(row, values.PullRequest.GetHTMLURL())
		} else {
			row = append(row, nil)
		}
		valErr := ""
		if values.Error != nil {
			valErr = values.Error.Error()
		}
		row = append(row, values.SkippedReason, valErr)

		t.AppendRow(row)
	}

	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t
}

// backportBranchNameForRef returns then branch name to use for our backport,
// e.g. ce/backport/1.19.x/my-feature-branch
func (r CreateBackportReq) backportBranchNameForRef(
	ref string,
	prBranch string,
) string {
	name := fmt.Sprintf("backport/%s/%s", ref, prBranch)
	if len(name) > 250 {
		// Handle Githubs branch name max length
		name = name[:250]
	}

	return name
}

func (r *CreateBackportReq) backportRef(
	ctx context.Context,
	git *libgit.Client,
	github *libgithub.Client,
	pr *libgithub.PullRequest,
	activeVersions map[string]*releases.Version,
	changedFiles *ListChangedFilesRes,
	ref string, // the full base ref of the branch we're backporting to
) *CreateBackportAttempt {
	res := &CreateBackportAttempt{BaseRef: ref}

	baseRefVersion := r.baseRefVersion(ref)
	// Get the name of our PR branch. We'll use this in our backport branch names
	// to make it easier to find the source.
	prBranch := pr.GetHead().GetRef()
	// The branch name for our backport, e.g. ce/backport/1.19.x/my-feature-branch
	branchName := r.backportBranchNameForRef(ref, prBranch)
	res.TargetRef = branchName
	commitSHA := pr.GetMergeCommitSHA()
	bigCtx := slogctx.Append(ctx,
		slog.String("target-base-ref", ref),
		slog.String("target-ref-version", baseRefVersion),
		slog.String("target-branch", branchName),
		slog.String("pr-branch", prBranch),
		slog.String("commit-sha", commitSHA),
	)

	if reason, shouldSkip := r.shouldSkipRef(
		ctx, baseRefVersion, ref, activeVersions, changedFiles,
	); shouldSkip {
		slog.Default().InfoContext(slogctx.Append(bigCtx,
			slog.String("base-ref-version", baseRefVersion),
			slog.String("target-ref", ref),
			slog.String("reason", reason),
		), "skipping backport")

		res.Skipped = true
		res.SkippedReason = reason

		return res
	}

	slog.Default().DebugContext(bigCtx, "creating backport pull request")
	slog.Default().DebugContext(ctx, "fetching backport target branch base ref")
	fetchRes, err := git.Fetch(ctx, &libgit.FetchOpts{
		// Fetch the ref but also provide a local tracking branch of the same name
		// e.g. "git fetch origin main:main"
		Refspec:     []string{r.BaseOrigin, fmt.Sprintf("%s:%s", ref, ref)},
		SetUpstream: true,
		Porcelain:   true,
	})
	if err != nil {
		res.Error = fmt.Errorf("fetching target branch base ref: %s, %w", fetchRes.String(), err)
		return res
	}

	slog.Default().DebugContext(ctx, "checking out new backport branch")
	checkoutRes, err := git.Checkout(ctx, &libgit.CheckoutOpts{
		NewBranchForceCheckout: branchName, // -B
		Branch:                 ref,
	})
	if err != nil {
		res.Error = fmt.Errorf("checking out new backport branch: %s: %w", checkoutRes.String(), err)
		return res
	}

	// Try and backport the commit
	if r.hasCEPrefix(ref) && changedFiles.Groups.Any(r.CEExclude) {
		// We're backporting enterprise to CE but the commit has files we don't
		// want to include. If we try and cherry-pick the commit it will almost
		// certainly fail unless the enterprise only file is new.
		res.Error = r.backportCECommitWithPatch(ctx, git, pr, changedFiles, commitSHA)
	} else {
		// We're backporting everything else. Simply cherry-pick the commit.
		slog.Default().DebugContext(ctx, "cherry-picking")
		cherryPickRes, err := git.CherryPick(ctx, &libgit.CherryPickOpts{
			FF:       true,
			Empty:    libgit.EmptyCommitKeep,
			Commit:   commitSHA,
			Strategy: libgit.MergeStrategyORT,
			StrategyOptions: []libgit.MergeStrategyOption{
				libgit.MergeStrategyOptionTheirs,
				libgit.MergeStrategyOptionIgnoreSpaceChange,
			},
		})
		if err != nil {
			res.Error = fmt.Errorf("cherry-picking backport merge commit: %s: %w", cherryPickRes.String(), err)
		}
	}

	// If our backport failed we still want to create a pull request for our
	// failed backport. There's still some debate and the validity of this approach
	// but our current process for ensuring backports have been merged is auditing
	// the open pull requests for a branch. Until that changes we'll need to do
	// this.
	if res.Error != nil {
		err = resetAndCreateNOOPCommit(ctx, git, ref)
		if err != nil {
			res.Error = errors.Join(res.Error, err)
		}
	}

	pushRes, err := git.Push(ctx, &libgit.PushOpts{
		Repository: r.BaseOrigin,
		Refspec:    []string{branchName},
	})
	if err != nil {
		res.Error = errors.Join(res.Error, fmt.Errorf("pushing backport branch: %s: %w", pushRes.String(), err))

		// If we didn't successfully push the branch we can't open a PR so it's time
		// to return.
		return res
	}

	// Generate title, removing existing "Backport" prefix to avoid stuttering
	cleanTitle := pr.GetTitle()
	if strings.HasPrefix(strings.ToLower(cleanTitle), "backport ") {
		cleanTitle = strings.TrimSpace(cleanTitle[9:]) // Remove "Backport " prefix
	}
	prTitle := fmt.Sprintf("Backport %s into %s", cleanTitle, ref)

	// Choose template based on whether this is a CE backport
	templateName := "backport-pr-message.tmpl"
	if r.hasCEPrefix(ref) {
		templateName = "backport-ce-pr-message.tmpl"
	}

	prBody, err := renderEmbeddedTemplate(templateName, struct {
		OriginPullRequest *libgithub.PullRequest
		Attempt           *CreateBackportAttempt
	}{pr, res})
	if err != nil {
		res.Error = fmt.Errorf("creating backport pull request body %w", err)
		return res
	}
	limitedPRBody := limitCharacters(prBody)
	res.PullRequest, _, err = github.PullRequests.Create(
		ctx, r.Owner, r.Repo, &libgithub.NewPullRequest{
			Title:    &prTitle,
			Head:     &branchName,
			HeadRepo: &r.Repo,
			Base:     &ref,
			Body:     &limitedPRBody,
		},
	)
	if err != nil {
		res.Error = fmt.Errorf("creating backport pull request %w", err)
		return res
	}

	// Assign the pull request to the actor that merged the pull request and/or the
	// person(s) that it was assigned to.
	err = addAssignees(
		ctx,
		github,
		r.Owner,
		r.Repo,
		int(res.PullRequest.GetNumber()),
		[]string{pr.GetAssignee().GetLogin(), pr.GetMergedBy().GetLogin()},
	)
	if err != nil {
		res.Error = fmt.Errorf("assigning ownership to backport pull request %w", err)
		return res
	}

	return res
}

// backportCECommitWithPatch backports a commit to the currently checked out
// branch and will omit and excluded files for CE backports. This commit
// backport strategy involves creating a new diff patch and applying it rather
// than a cherry-pick. We do this so as to not require fixing bad cherry-picks
// when modifying enterprise only files that don't exist on the CE branch.
func (r *CreateBackportReq) backportCECommitWithPatch(
	ctx context.Context,
	git *libgit.Client,
	pr *libgithub.PullRequest,
	changedFiles *ListChangedFilesRes,
	commitSHA string,
) error {
	var err error
	// Get a list of files that do not include excluded groups.
	files := changed.Files{}
	for _, file := range changedFiles.Files {
		if file.Groups.Any(r.CEExclude) {
			slog.Default().DebugContext(slogctx.Append(ctx,
				slog.String("file", file.Name()),
			), "skipping file as it is in one-or-more excluded groups")
		} else {
			slog.Default().DebugContext(slogctx.Append(ctx,
				slog.String("file", file.Name()),
			), "including changed file")
			files = append(files, file)
		}
	}

	// Create a unified patch of just the files we want to backport.
	tmpDir, err := os.MkdirTemp("", "ce-backport-patch")
	if err != nil {
		return fmt.Errorf("creating temporary directory for CE patches: %w", err)
	}
	patchFile := filepath.Join(tmpDir, pr.GetBase().GetSHA()+".patch")

	patchRes, err := git.Show(ctx, &libgit.ShowOpts{
		DiffAlgorithm: libgit.DiffAlgorithmMyers,
		// Use mboxrd so that we can we use 'git am' to apply and commit the patch
		// and inherit all metadata from the source commit.
		Format:   "mboxrd",
		NoColor:  true,
		Output:   patchFile,
		Object:   commitSHA,
		Patch:    true,
		PathSpec: files.Names(),
	})
	if err != nil {
		return fmt.Errorf("creating CE backport patch %s: %w", patchRes.String(), err)
	}

	// Apply the patch and commit it with the original details
	amRes, err := git.Am(ctx, &libgit.AmOpts{
		CommitterDateIsAuthorDate: true,
		Empty:                     libgit.EmptyCommitKeep,
		KeepNonPatch:              true,
		ThreeWayMerge:             true,
		Whitespace:                libgit.WhitespaceActionFix,
		Mbox:                      []string{patchFile},
	})
	if err != nil {
		return fmt.Errorf("apply CE backport patch: %s: %w", amRes.String(), err)
	}

	return nil
}

// baseRefVersion represents the baseRef as an active branch version. Active
// branch versions are defined in .release/versions.hcl and ought to be
// considered the source of truth for which CE branches are active. The output
// also maps 1:1 to with backport labels. e.g.
//
//	ce/main                => main
//	ent/main               => main
//	main                   => main
//	ce/release/1.19.x      => release/1.19.x
//	release/1.19.x+ent     => release/1.19.x
//	ent/release/1.19.x+ent => release/1.19.x
func (r *CreateBackportReq) baseRefVersion(ref string) string {
	switch {
	case r.hasCEPrefix(ref):
		return strings.TrimSuffix(strings.TrimPrefix(ref, r.CEBranchPrefix+"/"), "+ent")
	case r.hasEntPrefix(ref):
		return strings.TrimSuffix(strings.TrimPrefix(ref, r.EntBranchPrefix+"/"), "+ent")
	default:
		return strings.TrimSuffix(ref, "+ent")
	}
}

// determineBackportRefs determines which backport target branches are candidates
// to backport to depending on a combination of our source pull requests base
// reference and the labels that are present on the pull request.
//
// If the base reference of the original PR is main, we assume we ought to
// backport to ce/main.
//
// Any non-main backport references are derived from the original pull requests
// labels. The valid labels are translated to the corresponding references
// that match the source pull requests base reference type: enterprise or
// community
func (r *CreateBackportReq) determineBackportRefs(
	ctx context.Context,
	baseRef string,
	labels Labels,
) (res []string) {
	baseRefVersion := r.baseRefVersion(baseRef)
	slog.Default().DebugContext(slogctx.Append(ctx,
		slog.String("labels", strings.Join(labels.Names(), " ")),
		slog.String("base-ref", baseRef),
		slog.String("base-ref-version", baseRefVersion),
	), "determining backport base references from pull request labels")

	defer func() {
		if len(res) < 1 {
			res = nil
		}
	}()

	if r.isEntRef(baseRef) {
		// We're dealing an enterprise PR. Always backport to the corresponding
		// CE branch if it's active.
		if baseRefVersion == "main" {
			res = append(res, fmt.Sprintf("%s/main", r.CEBranchPrefix))
		} else {
			res = append(res, fmt.Sprintf("%s/%s", r.CEBranchPrefix, baseRefVersion))
		}

		// Backport to all enterprise release branches that match our backport labels
		for _, label := range labels.Names() {
			parts := strings.SplitN(label, "/", 2)
			if len(parts) != 2 || parts[0] != r.BackportLabelPrefix {
				slog.Default().DebugContext(slogctx.Append(ctx,
					slog.String("label", label),
					slog.String("backport-label-prefix", r.BackportLabelPrefix),
				), "skipping label because it does not match the backport label prefix")
				continue
			}

			if parts[1] == baseRefVersion {
				slog.Default().WarnContext(slogctx.Append(ctx,
					slog.String("label", label),
					slog.String("base-ref-version", baseRefVersion),
				), "skipping label because we cannot backport to the same reference")
				continue
			}

			if r.EntBranchPrefix == "" {
				res = append(res, fmt.Sprintf("release/%s+ent", parts[1]))
			} else {
				res = append(res, fmt.Sprintf("%s/release/%s+ent", r.EntBranchPrefix, parts[1]))
			}
		}
	} else {
		// We're dealing with a CE PR. Backport to all CE release branches that match
		// our backport labels
		for _, label := range labels.Names() {
			parts := strings.SplitN(label, "/", 2)
			if len(parts) != 2 || parts[0] != r.BackportLabelPrefix {
				slog.Default().DebugContext(slogctx.Append(ctx,
					slog.String("label", label),
					slog.String("backport-label-prefix", r.BackportLabelPrefix),
				), "skipping label because it does not match the backport label prefix")
				continue
			}

			if parts[1] == baseRefVersion {
				slog.Default().WarnContext(slogctx.Append(ctx,
					slog.String("label", label),
					slog.String("base-ref-version", baseRefVersion),
				), "skipping label because we cannot backport to the same reference")

				continue
			}

			res = append(res, fmt.Sprintf("%s/release/%s", r.CEBranchPrefix, parts[1]))
		}
	}

	slog.Default().DebugContext(slogctx.Append(ctx,
		slog.String("refs", strings.Join(res, ",")),
	), "determined target backport references")

	return res
}

// getActiveVersions gets the active versions from .release/versions.hcl
func (r *CreateBackportReq) getActiveVersions(
	ctx context.Context,
) (map[string]*releases.Version, error) {
	req := &releases.ListActiveVersionsReq{
		Recurse:                  r.ReleaseRecurseDepth,
		ReleaseVersionConfigPath: r.ReleaseVersionConfigPath,
	}
	res, err := req.Run(ctx)
	if err != nil {
		return nil, err
	}

	return res.VersionsConfig.ActiveVersion.Versions, nil
}

// getChangedFiles gets a list of files that changed in the PR and determines
// whether or not we need to worry about excluding some or all of them for CE
// backports.
func (r *CreateBackportReq) getChangedFiles(
	ctx context.Context,
	github *libgithub.Client,
) (*ListChangedFilesRes, error) {
	req := ListChangedFilesReq{
		Owner:      r.Owner,
		Repo:       r.Repo,
		PullNumber: int(r.PullNumber),
		GroupFiles: true,
	}
	res, err := req.Run(ctx, github)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Names returns the label names as slice of strings
func (l Labels) Names() []string {
	if len(l) < 1 {
		return nil
	}

	res := []string{}
	for label := range slices.Values(l) {
		if label != nil {
			res = append(res, label.GetName())
		}
	}

	return res
}

// hasCEPrefix takes a branch reference and determines whether or not it starts
// with the CEBranchPrefix.
func (r *CreateBackportReq) hasCEPrefix(ref string) bool {
	return strings.HasPrefix(ref, r.CEBranchPrefix+"/")
}

// hasEntPrefix takes a branch reference and determines whether or not it starts
// with the EntBranchPrefix.
func (r *CreateBackportReq) hasEntPrefix(ref string) bool {
	if r.EntBranchPrefix == "" {
		return false
	}

	return strings.HasPrefix(ref, r.EntBranchPrefix+"/")
}

// isEntRef takes a branch reference and determines whether or not it refers to
// an enterprise branch.
func (r *CreateBackportReq) isEntRef(ref string) bool {
	return !r.hasCEPrefix(ref)
}

// shouldSkipRef determines whether or we ought to backport to a given branch
// reference. It considers whether or not the base ref is for enterprise or
// CE, which files have changed and which CE branches are active.
func (r *CreateBackportReq) shouldSkipRef(
	ctx context.Context,
	baseRefVersion string,
	ref string,
	activeVersions map[string]*releases.Version,
	changedFiles *ListChangedFilesRes,
) (string, bool) {
	slog.Default().DebugContext(slogctx.Append(ctx,
		slog.String("base-ref-version", baseRefVersion),
		slog.String("target-ref", ref),
	), "determining whether to skip backport")

	if changedFiles == nil || len(changedFiles.Files) < 1 {
		return "no files were changed", true
	}

	if baseRefVersion == "" {
		return "missing base ref", true
	}

	if ref == "" {
		return "missing fef", true
	}

	if r.isEntRef(ref) {
		// It's an enterprise backport so we'll always do it.
		return "references to enterprise branches always backported", false
	}

	// Check if all of our files belong to excluded groups, i.e. they're all
	// files in the "enterprise" group.
	if changedFiles.Files.EachHasAnyGroup(r.CEExclude) {
		return fmt.Sprintf(
			"all changed files are in excluded groups: %s", r.CEExclude.String(),
		), true
	}

	if ref == r.CEBranchPrefix+"/main" {
		return "ce/main is always active and there are CE allowed files", false
	}

	// Check if there are inactive-allowed changed files, i.e. docs or pipeline
	// files are included so we'll always backport to the CE branch.
	if r.CEAllowInactiveGroups.Any(changedFiles.Groups) {
		return fmt.Sprintf(
			"one or more changed file groups [%s] are included in allowed inactive changed file groups [%s]",
			changedFiles.Groups.String(), r.CEAllowInactiveGroups.String(),
		), false
	}

	// Check if ce branch is active or not
	prefix := "release/"
	if r.EntBranchPrefix != "" {
		prefix = r.EntBranchPrefix + "/"
	}
	version, ok := strings.CutPrefix(baseRefVersion, prefix)
	if ok {
		if ver, ok := activeVersions[version]; ok {
			if ver.CEActive {
				return "CE branch is active", false
			}
			return "CE branch is inactive", true
		}
	}

	return fmt.Sprintf(
		"could not find branch in active branches configuration: %s", baseRefVersion,
	), true
}
