// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"strings"
)

// RepositoriesService handles communication with the repository related
// methods of the GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/repos/
type RepositoriesService service

// Repository represents a GitHub repository.
type Repository struct {
	ID               *int64           `json:"id,omitempty"`
	NodeID           *string          `json:"node_id,omitempty"`
	Owner            *User            `json:"owner,omitempty"`
	Name             *string          `json:"name,omitempty"`
	FullName         *string          `json:"full_name,omitempty"`
	Description      *string          `json:"description,omitempty"`
	Homepage         *string          `json:"homepage,omitempty"`
	CodeOfConduct    *CodeOfConduct   `json:"code_of_conduct,omitempty"`
	DefaultBranch    *string          `json:"default_branch,omitempty"`
	MasterBranch     *string          `json:"master_branch,omitempty"`
	CreatedAt        *Timestamp       `json:"created_at,omitempty"`
	PushedAt         *Timestamp       `json:"pushed_at,omitempty"`
	UpdatedAt        *Timestamp       `json:"updated_at,omitempty"`
	HTMLURL          *string          `json:"html_url,omitempty"`
	CloneURL         *string          `json:"clone_url,omitempty"`
	GitURL           *string          `json:"git_url,omitempty"`
	MirrorURL        *string          `json:"mirror_url,omitempty"`
	SSHURL           *string          `json:"ssh_url,omitempty"`
	SVNURL           *string          `json:"svn_url,omitempty"`
	Language         *string          `json:"language,omitempty"`
	Fork             *bool            `json:"fork,omitempty"`
	ForksCount       *int             `json:"forks_count,omitempty"`
	NetworkCount     *int             `json:"network_count,omitempty"`
	OpenIssuesCount  *int             `json:"open_issues_count,omitempty"`
	StargazersCount  *int             `json:"stargazers_count,omitempty"`
	SubscribersCount *int             `json:"subscribers_count,omitempty"`
	WatchersCount    *int             `json:"watchers_count,omitempty"`
	Size             *int             `json:"size,omitempty"`
	AutoInit         *bool            `json:"auto_init,omitempty"`
	Parent           *Repository      `json:"parent,omitempty"`
	Source           *Repository      `json:"source,omitempty"`
	Organization     *Organization    `json:"organization,omitempty"`
	Permissions      *map[string]bool `json:"permissions,omitempty"`
	AllowRebaseMerge *bool            `json:"allow_rebase_merge,omitempty"`
	AllowSquashMerge *bool            `json:"allow_squash_merge,omitempty"`
	AllowMergeCommit *bool            `json:"allow_merge_commit,omitempty"`
	Topics           []string         `json:"topics,omitempty"`

	// Only provided when using RepositoriesService.Get while in preview
	License *License `json:"license,omitempty"`

	// Additional mutable fields when creating and editing a repository
	Private           *bool   `json:"private,omitempty"`
	HasIssues         *bool   `json:"has_issues,omitempty"`
	HasWiki           *bool   `json:"has_wiki,omitempty"`
	HasPages          *bool   `json:"has_pages,omitempty"`
	HasProjects       *bool   `json:"has_projects,omitempty"`
	HasDownloads      *bool   `json:"has_downloads,omitempty"`
	LicenseTemplate   *string `json:"license_template,omitempty"`
	GitignoreTemplate *string `json:"gitignore_template,omitempty"`
	Archived          *bool   `json:"archived,omitempty"`

	// Creating an organization repository. Required for non-owners.
	TeamID *int64 `json:"team_id,omitempty"`

	// API URLs
	URL              *string `json:"url,omitempty"`
	ArchiveURL       *string `json:"archive_url,omitempty"`
	AssigneesURL     *string `json:"assignees_url,omitempty"`
	BlobsURL         *string `json:"blobs_url,omitempty"`
	BranchesURL      *string `json:"branches_url,omitempty"`
	CollaboratorsURL *string `json:"collaborators_url,omitempty"`
	CommentsURL      *string `json:"comments_url,omitempty"`
	CommitsURL       *string `json:"commits_url,omitempty"`
	CompareURL       *string `json:"compare_url,omitempty"`
	ContentsURL      *string `json:"contents_url,omitempty"`
	ContributorsURL  *string `json:"contributors_url,omitempty"`
	DeploymentsURL   *string `json:"deployments_url,omitempty"`
	DownloadsURL     *string `json:"downloads_url,omitempty"`
	EventsURL        *string `json:"events_url,omitempty"`
	ForksURL         *string `json:"forks_url,omitempty"`
	GitCommitsURL    *string `json:"git_commits_url,omitempty"`
	GitRefsURL       *string `json:"git_refs_url,omitempty"`
	GitTagsURL       *string `json:"git_tags_url,omitempty"`
	HooksURL         *string `json:"hooks_url,omitempty"`
	IssueCommentURL  *string `json:"issue_comment_url,omitempty"`
	IssueEventsURL   *string `json:"issue_events_url,omitempty"`
	IssuesURL        *string `json:"issues_url,omitempty"`
	KeysURL          *string `json:"keys_url,omitempty"`
	LabelsURL        *string `json:"labels_url,omitempty"`
	LanguagesURL     *string `json:"languages_url,omitempty"`
	MergesURL        *string `json:"merges_url,omitempty"`
	MilestonesURL    *string `json:"milestones_url,omitempty"`
	NotificationsURL *string `json:"notifications_url,omitempty"`
	PullsURL         *string `json:"pulls_url,omitempty"`
	ReleasesURL      *string `json:"releases_url,omitempty"`
	StargazersURL    *string `json:"stargazers_url,omitempty"`
	StatusesURL      *string `json:"statuses_url,omitempty"`
	SubscribersURL   *string `json:"subscribers_url,omitempty"`
	SubscriptionURL  *string `json:"subscription_url,omitempty"`
	TagsURL          *string `json:"tags_url,omitempty"`
	TreesURL         *string `json:"trees_url,omitempty"`
	TeamsURL         *string `json:"teams_url,omitempty"`

	// TextMatches is only populated from search results that request text matches
	// See: search.go and https://developer.github.com/v3/search/#text-match-metadata
	TextMatches []TextMatch `json:"text_matches,omitempty"`
}

func (r Repository) String() string {
	return Stringify(r)
}

// RepositoryListOptions specifies the optional parameters to the
// RepositoriesService.List method.
type RepositoryListOptions struct {
	// Visibility of repositories to list. Can be one of all, public, or private.
	// Default: all
	Visibility string `url:"visibility,omitempty"`

	// List repos of given affiliation[s].
	// Comma-separated list of values. Can include:
	// * owner: Repositories that are owned by the authenticated user.
	// * collaborator: Repositories that the user has been added to as a
	//   collaborator.
	// * organization_member: Repositories that the user has access to through
	//   being a member of an organization. This includes every repository on
	//   every team that the user is on.
	// Default: owner,collaborator,organization_member
	Affiliation string `url:"affiliation,omitempty"`

	// Type of repositories to list.
	// Can be one of all, owner, public, private, member. Default: all
	// Will cause a 422 error if used in the same request as visibility or
	// affiliation.
	Type string `url:"type,omitempty"`

	// How to sort the repository list. Can be one of created, updated, pushed,
	// full_name. Default: full_name
	Sort string `url:"sort,omitempty"`

	// Direction in which to sort repositories. Can be one of asc or desc.
	// Default: when using full_name: asc; otherwise desc
	Direction string `url:"direction,omitempty"`

	ListOptions
}

// List the repositories for a user. Passing the empty string will list
// repositories for the authenticated user.
//
// GitHub API docs: https://developer.github.com/v3/repos/#list-user-repositories
func (s *RepositoriesService) List(ctx context.Context, user string, opt *RepositoryListOptions) ([]*Repository, *Response, error) {
	var u string
	if user != "" {
		u = fmt.Sprintf("users/%v/repos", user)
	} else {
		u = "user/repos"
	}
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	acceptHeaders := []string{mediaTypeCodesOfConductPreview, mediaTypeTopicsPreview}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))

	var repos []*Repository
	resp, err := s.client.Do(ctx, req, &repos)
	if err != nil {
		return nil, resp, err
	}

	return repos, resp, nil
}

// RepositoryListByOrgOptions specifies the optional parameters to the
// RepositoriesService.ListByOrg method.
type RepositoryListByOrgOptions struct {
	// Type of repositories to list. Possible values are: all, public, private,
	// forks, sources, member. Default is "all".
	Type string `url:"type,omitempty"`

	ListOptions
}

// ListByOrg lists the repositories for an organization.
//
// GitHub API docs: https://developer.github.com/v3/repos/#list-organization-repositories
func (s *RepositoriesService) ListByOrg(ctx context.Context, org string, opt *RepositoryListByOrgOptions) ([]*Repository, *Response, error) {
	u := fmt.Sprintf("orgs/%v/repos", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept headers when APIs fully launch.
	acceptHeaders := []string{mediaTypeCodesOfConductPreview, mediaTypeTopicsPreview}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))

	var repos []*Repository
	resp, err := s.client.Do(ctx, req, &repos)
	if err != nil {
		return nil, resp, err
	}

	return repos, resp, nil
}

// RepositoryListAllOptions specifies the optional parameters to the
// RepositoriesService.ListAll method.
type RepositoryListAllOptions struct {
	// ID of the last repository seen
	Since int64 `url:"since,omitempty"`
}

// ListAll lists all GitHub repositories in the order that they were created.
//
// GitHub API docs: https://developer.github.com/v3/repos/#list-all-public-repositories
func (s *RepositoriesService) ListAll(ctx context.Context, opt *RepositoryListAllOptions) ([]*Repository, *Response, error) {
	u, err := addOptions("repositories", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var repos []*Repository
	resp, err := s.client.Do(ctx, req, &repos)
	if err != nil {
		return nil, resp, err
	}

	return repos, resp, nil
}

// Create a new repository. If an organization is specified, the new
// repository will be created under that org. If the empty string is
// specified, it will be created for the authenticated user.
//
// GitHub API docs: https://developer.github.com/v3/repos/#create
func (s *RepositoriesService) Create(ctx context.Context, org string, repo *Repository) (*Repository, *Response, error) {
	var u string
	if org != "" {
		u = fmt.Sprintf("orgs/%v/repos", org)
	} else {
		u = "user/repos"
	}

	req, err := s.client.NewRequest("POST", u, repo)
	if err != nil {
		return nil, nil, err
	}

	r := new(Repository)
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// Get fetches a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/#get
func (s *RepositoriesService) Get(ctx context.Context, owner, repo string) (*Repository, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when the license support fully launches
	// https://developer.github.com/v3/licenses/#get-a-repositorys-license
	acceptHeaders := []string{mediaTypeCodesOfConductPreview, mediaTypeTopicsPreview}
	req.Header.Set("Accept", strings.Join(acceptHeaders, ", "))

	repository := new(Repository)
	resp, err := s.client.Do(ctx, req, repository)
	if err != nil {
		return nil, resp, err
	}

	return repository, resp, nil
}

// GetCodeOfConduct gets the contents of a repository's code of conduct.
//
// GitHub API docs: https://developer.github.com/v3/codes_of_conduct/#get-the-contents-of-a-repositorys-code-of-conduct
func (s *RepositoriesService) GetCodeOfConduct(ctx context.Context, owner, repo string) (*CodeOfConduct, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/community/code_of_conduct", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeCodesOfConductPreview)

	coc := new(CodeOfConduct)
	resp, err := s.client.Do(ctx, req, coc)
	if err != nil {
		return nil, resp, err
	}

	return coc, resp, nil
}

// GetByID fetches a repository.
//
// Note: GetByID uses the undocumented GitHub API endpoint /repositories/:id.
func (s *RepositoriesService) GetByID(ctx context.Context, id int64) (*Repository, *Response, error) {
	u := fmt.Sprintf("repositories/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	repository := new(Repository)
	resp, err := s.client.Do(ctx, req, repository)
	if err != nil {
		return nil, resp, err
	}

	return repository, resp, nil
}

// Edit updates a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/#edit
func (s *RepositoriesService) Edit(ctx context.Context, owner, repo string, repository *Repository) (*Repository, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v", owner, repo)
	req, err := s.client.NewRequest("PATCH", u, repository)
	if err != nil {
		return nil, nil, err
	}

	r := new(Repository)
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// Delete a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/#delete-a-repository
func (s *RepositoriesService) Delete(ctx context.Context, owner, repo string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v", owner, repo)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Contributor represents a repository contributor
type Contributor struct {
	Login             *string `json:"login,omitempty"`
	ID                *int64  `json:"id,omitempty"`
	AvatarURL         *string `json:"avatar_url,omitempty"`
	GravatarID        *string `json:"gravatar_id,omitempty"`
	URL               *string `json:"url,omitempty"`
	HTMLURL           *string `json:"html_url,omitempty"`
	FollowersURL      *string `json:"followers_url,omitempty"`
	FollowingURL      *string `json:"following_url,omitempty"`
	GistsURL          *string `json:"gists_url,omitempty"`
	StarredURL        *string `json:"starred_url,omitempty"`
	SubscriptionsURL  *string `json:"subscriptions_url,omitempty"`
	OrganizationsURL  *string `json:"organizations_url,omitempty"`
	ReposURL          *string `json:"repos_url,omitempty"`
	EventsURL         *string `json:"events_url,omitempty"`
	ReceivedEventsURL *string `json:"received_events_url,omitempty"`
	Type              *string `json:"type,omitempty"`
	SiteAdmin         *bool   `json:"site_admin,omitempty"`
	Contributions     *int    `json:"contributions,omitempty"`
}

// ListContributorsOptions specifies the optional parameters to the
// RepositoriesService.ListContributors method.
type ListContributorsOptions struct {
	// Include anonymous contributors in results or not
	Anon string `url:"anon,omitempty"`

	ListOptions
}

// ListContributors lists contributors for a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/#list-contributors
func (s *RepositoriesService) ListContributors(ctx context.Context, owner string, repository string, opt *ListContributorsOptions) ([]*Contributor, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/contributors", owner, repository)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var contributor []*Contributor
	resp, err := s.client.Do(ctx, req, &contributor)
	if err != nil {
		return nil, nil, err
	}

	return contributor, resp, nil
}

// ListLanguages lists languages for the specified repository. The returned map
// specifies the languages and the number of bytes of code written in that
// language. For example:
//
//     {
//       "C": 78769,
//       "Python": 7769
//     }
//
// GitHub API docs: https://developer.github.com/v3/repos/#list-languages
func (s *RepositoriesService) ListLanguages(ctx context.Context, owner string, repo string) (map[string]int, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/languages", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	languages := make(map[string]int)
	resp, err := s.client.Do(ctx, req, &languages)
	if err != nil {
		return nil, resp, err
	}

	return languages, resp, nil
}

// ListTeams lists the teams for the specified repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/#list-teams
func (s *RepositoriesService) ListTeams(ctx context.Context, owner string, repo string, opt *ListOptions) ([]*Team, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/teams", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeNestedTeamsPreview)

	var teams []*Team
	resp, err := s.client.Do(ctx, req, &teams)
	if err != nil {
		return nil, resp, err
	}

	return teams, resp, nil
}

// RepositoryTag represents a repository tag.
type RepositoryTag struct {
	Name       *string `json:"name,omitempty"`
	Commit     *Commit `json:"commit,omitempty"`
	ZipballURL *string `json:"zipball_url,omitempty"`
	TarballURL *string `json:"tarball_url,omitempty"`
}

// ListTags lists tags for the specified repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/#list-tags
func (s *RepositoriesService) ListTags(ctx context.Context, owner string, repo string, opt *ListOptions) ([]*RepositoryTag, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/tags", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tags []*RepositoryTag
	resp, err := s.client.Do(ctx, req, &tags)
	if err != nil {
		return nil, resp, err
	}

	return tags, resp, nil
}

// Branch represents a repository branch
type Branch struct {
	Name      *string           `json:"name,omitempty"`
	Commit    *RepositoryCommit `json:"commit,omitempty"`
	Protected *bool             `json:"protected,omitempty"`
}

// Protection represents a repository branch's protection.
type Protection struct {
	RequiredStatusChecks       *RequiredStatusChecks          `json:"required_status_checks"`
	RequiredPullRequestReviews *PullRequestReviewsEnforcement `json:"required_pull_request_reviews"`
	EnforceAdmins              *AdminEnforcement              `json:"enforce_admins"`
	Restrictions               *BranchRestrictions            `json:"restrictions"`
}

// ProtectionRequest represents a request to create/edit a branch's protection.
type ProtectionRequest struct {
	RequiredStatusChecks       *RequiredStatusChecks                 `json:"required_status_checks"`
	RequiredPullRequestReviews *PullRequestReviewsEnforcementRequest `json:"required_pull_request_reviews"`
	EnforceAdmins              bool                                  `json:"enforce_admins"`
	Restrictions               *BranchRestrictionsRequest            `json:"restrictions"`
}

// RequiredStatusChecks represents the protection status of a individual branch.
type RequiredStatusChecks struct {
	// Require branches to be up to date before merging. (Required.)
	Strict bool `json:"strict"`
	// The list of status checks to require in order to merge into this
	// branch. (Required; use []string{} instead of nil for empty list.)
	Contexts []string `json:"contexts"`
}

// RequiredStatusChecksRequest represents a request to edit a protected branch's status checks.
type RequiredStatusChecksRequest struct {
	Strict   *bool    `json:"strict,omitempty"`
	Contexts []string `json:"contexts,omitempty"`
}

// PullRequestReviewsEnforcement represents the pull request reviews enforcement of a protected branch.
type PullRequestReviewsEnforcement struct {
	// Specifies which users and teams can dismiss pull request reviews.
	DismissalRestrictions DismissalRestrictions `json:"dismissal_restrictions"`
	// Specifies if approved reviews are dismissed automatically, when a new commit is pushed.
	DismissStaleReviews bool `json:"dismiss_stale_reviews"`
	// RequireCodeOwnerReviews specifies if an approved review is required in pull requests including files with a designated code owner.
	RequireCodeOwnerReviews bool `json:"require_code_owner_reviews"`
	// RequiredApprovingReviewCount specifies the number of approvals required before the pull request can be merged.
	// Valid values are 1-6.
	RequiredApprovingReviewCount int `json:"required_approving_review_count"`
}

// PullRequestReviewsEnforcementRequest represents request to set the pull request review
// enforcement of a protected branch. It is separate from PullRequestReviewsEnforcement above
// because the request structure is different from the response structure.
type PullRequestReviewsEnforcementRequest struct {
	// Specifies which users and teams should be allowed to dismiss pull request reviews.
	// User and team dismissal restrictions are only available for
	// organization-owned repositories. Must be nil for personal repositories.
	DismissalRestrictionsRequest *DismissalRestrictionsRequest `json:"dismissal_restrictions,omitempty"`
	// Specifies if approved reviews can be dismissed automatically, when a new commit is pushed. (Required)
	DismissStaleReviews bool `json:"dismiss_stale_reviews"`
	// RequireCodeOwnerReviews specifies if an approved review is required in pull requests including files with a designated code owner.
	RequireCodeOwnerReviews bool `json:"require_code_owner_reviews"`
	// RequiredApprovingReviewCount specifies the number of approvals required before the pull request can be merged.
	// Valid values are 1-6.
	RequiredApprovingReviewCount int `json:"required_approving_review_count"`
}

// PullRequestReviewsEnforcementUpdate represents request to patch the pull request review
// enforcement of a protected branch. It is separate from PullRequestReviewsEnforcementRequest above
// because the patch request does not require all fields to be initialized.
type PullRequestReviewsEnforcementUpdate struct {
	// Specifies which users and teams can dismiss pull request reviews. Can be omitted.
	DismissalRestrictionsRequest *DismissalRestrictionsRequest `json:"dismissal_restrictions,omitempty"`
	// Specifies if approved reviews can be dismissed automatically, when a new commit is pushed. Can be omitted.
	DismissStaleReviews *bool `json:"dismiss_stale_reviews,omitempty"`
	// RequireCodeOwnerReviews specifies if an approved review is required in pull requests including files with a designated code owner.
	RequireCodeOwnerReviews bool `json:"require_code_owner_reviews,omitempty"`
	// RequiredApprovingReviewCount specifies the number of approvals required before the pull request can be merged.
	// Valid values are 1 - 6.
	RequiredApprovingReviewCount int `json:"required_approving_review_count"`
}

// AdminEnforcement represents the configuration to enforce required status checks for repository administrators.
type AdminEnforcement struct {
	URL     *string `json:"url,omitempty"`
	Enabled bool    `json:"enabled"`
}

// BranchRestrictions represents the restriction that only certain users or
// teams may push to a branch.
type BranchRestrictions struct {
	// The list of user logins with push access.
	Users []*User `json:"users"`
	// The list of team slugs with push access.
	Teams []*Team `json:"teams"`
}

// BranchRestrictionsRequest represents the request to create/edit the
// restriction that only certain users or teams may push to a branch. It is
// separate from BranchRestrictions above because the request structure is
// different from the response structure.
type BranchRestrictionsRequest struct {
	// The list of user logins with push access. (Required; use []string{} instead of nil for empty list.)
	Users []string `json:"users"`
	// The list of team slugs with push access. (Required; use []string{} instead of nil for empty list.)
	Teams []string `json:"teams"`
}

// DismissalRestrictions specifies which users and teams can dismiss pull request reviews.
type DismissalRestrictions struct {
	// The list of users who can dimiss pull request reviews.
	Users []*User `json:"users"`
	// The list of teams which can dismiss pull request reviews.
	Teams []*Team `json:"teams"`
}

// DismissalRestrictionsRequest represents the request to create/edit the
// restriction to allows only specific users or teams to dimiss pull request reviews. It is
// separate from DismissalRestrictions above because the request structure is
// different from the response structure.
// Note: Both Users and Teams must be nil, or both must be non-nil.
type DismissalRestrictionsRequest struct {
	// The list of user logins who can dismiss pull request reviews. (Required; use nil to disable dismissal_restrictions or &[]string{} otherwise.)
	Users *[]string `json:"users,omitempty"`
	// The list of team slugs which can dismiss pull request reviews. (Required; use nil to disable dismissal_restrictions or &[]string{} otherwise.)
	Teams *[]string `json:"teams,omitempty"`
}

// ListBranches lists branches for the specified repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/#list-branches
func (s *RepositoriesService) ListBranches(ctx context.Context, owner string, repo string, opt *ListOptions) ([]*Branch, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	var branches []*Branch
	resp, err := s.client.Do(ctx, req, &branches)
	if err != nil {
		return nil, resp, err
	}

	return branches, resp, nil
}

// GetBranch gets the specified branch for a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/#get-branch
func (s *RepositoriesService) GetBranch(ctx context.Context, owner, repo, branch string) (*Branch, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v", owner, repo, branch)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	b := new(Branch)
	resp, err := s.client.Do(ctx, req, b)
	if err != nil {
		return nil, resp, err
	}

	return b, resp, nil
}

// GetBranchProtection gets the protection of a given branch.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#get-branch-protection
func (s *RepositoriesService) GetBranchProtection(ctx context.Context, owner, repo, branch string) (*Protection, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection", owner, repo, branch)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	p := new(Protection)
	resp, err := s.client.Do(ctx, req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// GetRequiredStatusChecks gets the required status checks for a given protected branch.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#get-required-status-checks-of-protected-branch
func (s *RepositoriesService) GetRequiredStatusChecks(ctx context.Context, owner, repo, branch string) (*RequiredStatusChecks, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection/required_status_checks", owner, repo, branch)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	p := new(RequiredStatusChecks)
	resp, err := s.client.Do(ctx, req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// ListRequiredStatusChecksContexts lists the required status checks contexts for a given protected branch.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#list-required-status-checks-contexts-of-protected-branch
func (s *RepositoriesService) ListRequiredStatusChecksContexts(ctx context.Context, owner, repo, branch string) (contexts []string, resp *Response, err error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection/required_status_checks/contexts", owner, repo, branch)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	resp, err = s.client.Do(ctx, req, &contexts)
	if err != nil {
		return nil, resp, err
	}

	return contexts, resp, nil
}

// UpdateBranchProtection updates the protection of a given branch.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#update-branch-protection
func (s *RepositoriesService) UpdateBranchProtection(ctx context.Context, owner, repo, branch string, preq *ProtectionRequest) (*Protection, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection", owner, repo, branch)
	req, err := s.client.NewRequest("PUT", u, preq)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	p := new(Protection)
	resp, err := s.client.Do(ctx, req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// RemoveBranchProtection removes the protection of a given branch.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#remove-branch-protection
func (s *RepositoriesService) RemoveBranchProtection(ctx context.Context, owner, repo, branch string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection", owner, repo, branch)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	return s.client.Do(ctx, req, nil)
}

// UpdateRequiredStatusChecks updates the required status checks for a given protected branch.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#update-required-status-checks-of-protected-branch
func (s *RepositoriesService) UpdateRequiredStatusChecks(ctx context.Context, owner, repo, branch string, sreq *RequiredStatusChecksRequest) (*RequiredStatusChecks, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection/required_status_checks", owner, repo, branch)
	req, err := s.client.NewRequest("PATCH", u, sreq)
	if err != nil {
		return nil, nil, err
	}

	sc := new(RequiredStatusChecks)
	resp, err := s.client.Do(ctx, req, sc)
	if err != nil {
		return nil, resp, err
	}

	return sc, resp, nil
}

// License gets the contents of a repository's license if one is detected.
//
// GitHub API docs: https://developer.github.com/v3/licenses/#get-the-contents-of-a-repositorys-license
func (s *RepositoriesService) License(ctx context.Context, owner, repo string) (*RepositoryLicense, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/license", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	r := &RepositoryLicense{}
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// GetPullRequestReviewEnforcement gets pull request review enforcement of a protected branch.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#get-pull-request-review-enforcement-of-protected-branch
func (s *RepositoriesService) GetPullRequestReviewEnforcement(ctx context.Context, owner, repo, branch string) (*PullRequestReviewsEnforcement, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection/required_pull_request_reviews", owner, repo, branch)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	r := new(PullRequestReviewsEnforcement)
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// UpdatePullRequestReviewEnforcement patches pull request review enforcement of a protected branch.
// It requires admin access and branch protection to be enabled.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#update-pull-request-review-enforcement-of-protected-branch
func (s *RepositoriesService) UpdatePullRequestReviewEnforcement(ctx context.Context, owner, repo, branch string, patch *PullRequestReviewsEnforcementUpdate) (*PullRequestReviewsEnforcement, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection/required_pull_request_reviews", owner, repo, branch)
	req, err := s.client.NewRequest("PATCH", u, patch)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	r := new(PullRequestReviewsEnforcement)
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}

// DisableDismissalRestrictions disables dismissal restrictions of a protected branch.
// It requires admin access and branch protection to be enabled.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#update-pull-request-review-enforcement-of-protected-branch
func (s *RepositoriesService) DisableDismissalRestrictions(ctx context.Context, owner, repo, branch string) (*PullRequestReviewsEnforcement, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection/required_pull_request_reviews", owner, repo, branch)

	data := struct {
		R []interface{} `json:"dismissal_restrictions"`
	}{[]interface{}{}}

	req, err := s.client.NewRequest("PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	r := new(PullRequestReviewsEnforcement)
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}

// RemovePullRequestReviewEnforcement removes pull request enforcement of a protected branch.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#remove-pull-request-review-enforcement-of-protected-branch
func (s *RepositoriesService) RemovePullRequestReviewEnforcement(ctx context.Context, owner, repo, branch string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection/required_pull_request_reviews", owner, repo, branch)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	return s.client.Do(ctx, req, nil)
}

// GetAdminEnforcement gets admin enforcement information of a protected branch.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#get-admin-enforcement-of-protected-branch
func (s *RepositoriesService) GetAdminEnforcement(ctx context.Context, owner, repo, branch string) (*AdminEnforcement, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection/enforce_admins", owner, repo, branch)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	r := new(AdminEnforcement)
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}

// AddAdminEnforcement adds admin enforcement to a protected branch.
// It requires admin access and branch protection to be enabled.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#add-admin-enforcement-of-protected-branch
func (s *RepositoriesService) AddAdminEnforcement(ctx context.Context, owner, repo, branch string) (*AdminEnforcement, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection/enforce_admins", owner, repo, branch)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	r := new(AdminEnforcement)
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}

// RemoveAdminEnforcement removes admin enforcement from a protected branch.
//
// GitHub API docs: https://developer.github.com/v3/repos/branches/#remove-admin-enforcement-of-protected-branch
func (s *RepositoriesService) RemoveAdminEnforcement(ctx context.Context, owner, repo, branch string) (*Response, error) {
	u := fmt.Sprintf("repos/%v/%v/branches/%v/protection/enforce_admins", owner, repo, branch)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	// TODO: remove custom Accept header when this API fully launches
	req.Header.Set("Accept", mediaTypeRequiredApprovingReviewsPreview)

	return s.client.Do(ctx, req, nil)
}

// repositoryTopics represents a collection of repository topics.
type repositoryTopics struct {
	Names []string `json:"names"`
}

// ListAllTopics lists topics for a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/#list-all-topics-for-a-repository
func (s *RepositoriesService) ListAllTopics(ctx context.Context, owner, repo string) ([]string, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/topics", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeTopicsPreview)

	topics := new(repositoryTopics)
	resp, err := s.client.Do(ctx, req, topics)
	if err != nil {
		return nil, resp, err
	}

	return topics.Names, resp, nil
}

// ReplaceAllTopics replaces topics for a repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/#replace-all-topics-for-a-repository
func (s *RepositoriesService) ReplaceAllTopics(ctx context.Context, owner, repo string, topics []string) ([]string, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/topics", owner, repo)
	t := &repositoryTopics{
		Names: topics,
	}
	if t.Names == nil {
		t.Names = []string{}
	}
	req, err := s.client.NewRequest("PUT", u, t)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeTopicsPreview)

	t = new(repositoryTopics)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t.Names, resp, nil
}

// TransferRequest represents a request to transfer a repository.
type TransferRequest struct {
	NewOwner string  `json:"new_owner"`
	TeamID   []int64 `json:"team_ids,omitempty"`
}

// Transfer transfers a repository from one account or organization to another.
//
// This method might return an *AcceptedError and a status code of
// 202. This is because this is the status that GitHub returns to signify that
// it has now scheduled the transfer of the repository in a background task.
// A follow up request, after a delay of a second or so, should result
// in a successful request.
//
// GitHub API docs: https://developer.github.com/v3/repos/#transfer-a-repository
func (s *RepositoriesService) Transfer(ctx context.Context, owner, repo string, transfer TransferRequest) (*Repository, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/transfer", owner, repo)

	req, err := s.client.NewRequest("POST", u, &transfer)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeRepositoryTransferPreview)

	r := new(Repository)
	resp, err := s.client.Do(ctx, req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, nil
}
