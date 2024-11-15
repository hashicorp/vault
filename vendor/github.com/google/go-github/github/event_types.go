// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// These event types are shared between the Events API and used as Webhook payloads.

package github

// CheckRunEvent is triggered when a check run is "created", "updated", or "re-requested".
// The Webhook event name is "check_run".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#checkrunevent
type CheckRunEvent struct {
	CheckRun *CheckRun `json:"check_run,omitempty"`
	// The action performed. Can be "created", "updated" or "re-requested".
	Action *string `json:"action,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Org          *Organization `json:"organization,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// CheckSuiteEvent is triggered when a check suite is "completed", "requested", or "re-requested".
// The Webhook event name is "check_suite".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#checksuiteevent
type CheckSuiteEvent struct {
	CheckSuite *CheckSuite `json:"check_suite,omitempty"`
	// The action performed. Can be "completed", "requested" or "re-requested".
	Action *string `json:"action,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Org          *Organization `json:"organization,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// CommitCommentEvent is triggered when a commit comment is created.
// The Webhook event name is "commit_comment".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#commitcommentevent
type CommitCommentEvent struct {
	Comment *RepositoryComment `json:"comment,omitempty"`

	// The following fields are only populated by Webhook events.
	Action       *string       `json:"action,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// CreateEvent represents a created repository, branch, or tag.
// The Webhook event name is "create".
//
// Note: webhooks will not receive this event for created repositories.
// Additionally, webhooks will not receive this event for tags if more
// than three tags are pushed at once.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#createevent
type CreateEvent struct {
	Ref *string `json:"ref,omitempty"`
	// RefType is the object that was created. Possible values are: "repository", "branch", "tag".
	RefType      *string `json:"ref_type,omitempty"`
	MasterBranch *string `json:"master_branch,omitempty"`
	Description  *string `json:"description,omitempty"`

	// The following fields are only populated by Webhook events.
	PusherType   *string       `json:"pusher_type,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// DeleteEvent represents a deleted branch or tag.
// The Webhook event name is "delete".
//
// Note: webhooks will not receive this event for tags if more than three tags
// are deleted at once.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#deleteevent
type DeleteEvent struct {
	Ref *string `json:"ref,omitempty"`
	// RefType is the object that was deleted. Possible values are: "branch", "tag".
	RefType *string `json:"ref_type,omitempty"`

	// The following fields are only populated by Webhook events.
	PusherType   *string       `json:"pusher_type,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// DeploymentEvent represents a deployment.
// The Webhook event name is "deployment".
//
// Events of this type are not visible in timelines, they are only used to trigger hooks.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#deploymentevent
type DeploymentEvent struct {
	Deployment *Deployment `json:"deployment,omitempty"`
	Repo       *Repository `json:"repository,omitempty"`

	// The following fields are only populated by Webhook events.
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// DeploymentStatusEvent represents a deployment status.
// The Webhook event name is "deployment_status".
//
// Events of this type are not visible in timelines, they are only used to trigger hooks.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#deploymentstatusevent
type DeploymentStatusEvent struct {
	Deployment       *Deployment       `json:"deployment,omitempty"`
	DeploymentStatus *DeploymentStatus `json:"deployment_status,omitempty"`
	Repo             *Repository       `json:"repository,omitempty"`

	// The following fields are only populated by Webhook events.
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// ForkEvent is triggered when a user forks a repository.
// The Webhook event name is "fork".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#forkevent
type ForkEvent struct {
	// Forkee is the created repository.
	Forkee *Repository `json:"forkee,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// Page represents a single Wiki page.
type Page struct {
	PageName *string `json:"page_name,omitempty"`
	Title    *string `json:"title,omitempty"`
	Summary  *string `json:"summary,omitempty"`
	Action   *string `json:"action,omitempty"`
	SHA      *string `json:"sha,omitempty"`
	HTMLURL  *string `json:"html_url,omitempty"`
}

// GollumEvent is triggered when a Wiki page is created or updated.
// The Webhook event name is "gollum".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#gollumevent
type GollumEvent struct {
	Pages []*Page `json:"pages,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// EditChange represents the changes when an issue, pull request, or comment has
// been edited.
type EditChange struct {
	Title *struct {
		From *string `json:"from,omitempty"`
	} `json:"title,omitempty"`
	Body *struct {
		From *string `json:"from,omitempty"`
	} `json:"body,omitempty"`
}

// ProjectChange represents the changes when a project has been edited.
type ProjectChange struct {
	Name *struct {
		From *string `json:"from,omitempty"`
	} `json:"name,omitempty"`
	Body *struct {
		From *string `json:"from,omitempty"`
	} `json:"body,omitempty"`
}

// ProjectCardChange represents the changes when a project card has been edited.
type ProjectCardChange struct {
	Note *struct {
		From *string `json:"from,omitempty"`
	} `json:"note,omitempty"`
}

// ProjectColumnChange represents the changes when a project column has been edited.
type ProjectColumnChange struct {
	Name *struct {
		From *string `json:"from,omitempty"`
	} `json:"name,omitempty"`
}

// TeamChange represents the changes when a team has been edited.
type TeamChange struct {
	Description *struct {
		From *string `json:"from,omitempty"`
	} `json:"description,omitempty"`
	Name *struct {
		From *string `json:"from,omitempty"`
	} `json:"name,omitempty"`
	Privacy *struct {
		From *string `json:"from,omitempty"`
	} `json:"privacy,omitempty"`
	Repository *struct {
		Permissions *struct {
			From *struct {
				Admin *bool `json:"admin,omitempty"`
				Pull  *bool `json:"pull,omitempty"`
				Push  *bool `json:"push,omitempty"`
			} `json:"from,omitempty"`
		} `json:"permissions,omitempty"`
	} `json:"repository,omitempty"`
}

// InstallationEvent is triggered when a GitHub App has been installed or uninstalled.
// The Webhook event name is "installation".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#installationevent
type InstallationEvent struct {
	// The action that was performed. Can be either "created" or "deleted".
	Action       *string       `json:"action,omitempty"`
	Repositories []*Repository `json:"repositories,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// InstallationRepositoriesEvent is triggered when a repository is added or
// removed from an installation. The Webhook event name is "installation_repositories".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#installationrepositoriesevent
type InstallationRepositoriesEvent struct {
	// The action that was performed. Can be either "added" or "removed".
	Action              *string       `json:"action,omitempty"`
	RepositoriesAdded   []*Repository `json:"repositories_added,omitempty"`
	RepositoriesRemoved []*Repository `json:"repositories_removed,omitempty"`
	RepositorySelection *string       `json:"repository_selection,omitempty"`
	Sender              *User         `json:"sender,omitempty"`
	Installation        *Installation `json:"installation,omitempty"`
}

// IssueCommentEvent is triggered when an issue comment is created on an issue
// or pull request.
// The Webhook event name is "issue_comment".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#issuecommentevent
type IssueCommentEvent struct {
	// Action is the action that was performed on the comment.
	// Possible values are: "created", "edited", "deleted".
	Action  *string       `json:"action,omitempty"`
	Issue   *Issue        `json:"issue,omitempty"`
	Comment *IssueComment `json:"comment,omitempty"`

	// The following fields are only populated by Webhook events.
	Changes      *EditChange   `json:"changes,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// IssuesEvent is triggered when an issue is assigned, unassigned, labeled,
// unlabeled, opened, closed, or reopened.
// The Webhook event name is "issues".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#issuesevent
type IssuesEvent struct {
	// Action is the action that was performed. Possible values are: "assigned",
	// "unassigned", "labeled", "unlabeled", "opened", "closed", "reopened", "edited".
	Action   *string `json:"action,omitempty"`
	Issue    *Issue  `json:"issue,omitempty"`
	Assignee *User   `json:"assignee,omitempty"`
	Label    *Label  `json:"label,omitempty"`

	// The following fields are only populated by Webhook events.
	Changes      *EditChange   `json:"changes,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// LabelEvent is triggered when a repository's label is created, edited, or deleted.
// The Webhook event name is "label"
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#labelevent
type LabelEvent struct {
	// Action is the action that was performed. Possible values are:
	// "created", "edited", "deleted"
	Action *string `json:"action,omitempty"`
	Label  *Label  `json:"label,omitempty"`

	// The following fields are only populated by Webhook events.
	Changes      *EditChange   `json:"changes,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Org          *Organization `json:"organization,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// MarketplacePurchaseEvent is triggered when a user purchases, cancels, or changes
// their GitHub Marketplace plan.
// Webhook event name "marketplace_purchase".
//
// Github API docs: https://developer.github.com/v3/activity/events/types/#marketplacepurchaseevent
type MarketplacePurchaseEvent struct {
	// Action is the action that was performed. Possible values are:
	// "purchased", "cancelled", "changed".
	Action *string `json:"action,omitempty"`

	// The following fields are only populated by Webhook events.
	EffectiveDate               *Timestamp           `json:"effective_date,omitempty"`
	MarketplacePurchase         *MarketplacePurchase `json:"marketplace_purchase,omitempty"`
	PreviousMarketplacePurchase *MarketplacePurchase `json:"previous_marketplace_purchase,omitempty"`
	Sender                      *User                `json:"sender,omitempty"`
	Installation                *Installation        `json:"installation,omitempty"`
}

// MemberEvent is triggered when a user is added as a collaborator to a repository.
// The Webhook event name is "member".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#memberevent
type MemberEvent struct {
	// Action is the action that was performed. Possible value is: "added".
	Action *string `json:"action,omitempty"`
	Member *User   `json:"member,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// MembershipEvent is triggered when a user is added or removed from a team.
// The Webhook event name is "membership".
//
// Events of this type are not visible in timelines, they are only used to
// trigger organization webhooks.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#membershipevent
type MembershipEvent struct {
	// Action is the action that was performed. Possible values are: "added", "removed".
	Action *string `json:"action,omitempty"`
	// Scope is the scope of the membership. Possible value is: "team".
	Scope  *string `json:"scope,omitempty"`
	Member *User   `json:"member,omitempty"`
	Team   *Team   `json:"team,omitempty"`

	// The following fields are only populated by Webhook events.
	Org          *Organization `json:"organization,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// MilestoneEvent is triggered when a milestone is created, closed, opened, edited, or deleted.
// The Webhook event name is "milestone".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#milestoneevent
type MilestoneEvent struct {
	// Action is the action that was performed. Possible values are:
	// "created", "closed", "opened", "edited", "deleted"
	Action    *string    `json:"action,omitempty"`
	Milestone *Milestone `json:"milestone,omitempty"`

	// The following fields are only populated by Webhook events.
	Changes      *EditChange   `json:"changes,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Org          *Organization `json:"organization,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// OrganizationEvent is triggered when a user is added, removed, or invited to an organization.
// Events of this type are not visible in timelines. These events are only used to trigger organization hooks.
// Webhook event name is "organization".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#organizationevent
type OrganizationEvent struct {
	// Action is the action that was performed.
	// Can be one of "member_added", "member_removed", or "member_invited".
	Action *string `json:"action,omitempty"`

	// Invitaion is the invitation for the user or email if the action is "member_invited".
	Invitation *Invitation `json:"invitation,omitempty"`

	// Membership is the membership between the user and the organization.
	// Not present when the action is "member_invited".
	Membership *Membership `json:"membership,omitempty"`

	Organization *Organization `json:"organization,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// OrgBlockEvent is triggered when an organization blocks or unblocks a user.
// The Webhook event name is "org_block".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#orgblockevent
type OrgBlockEvent struct {
	// Action is the action that was performed.
	// Can be "blocked" or "unblocked".
	Action       *string       `json:"action,omitempty"`
	BlockedUser  *User         `json:"blocked_user,omitempty"`
	Organization *Organization `json:"organization,omitempty"`
	Sender       *User         `json:"sender,omitempty"`

	// The following fields are only populated by Webhook events.
	Installation *Installation `json:"installation,omitempty"`
}

// PageBuildEvent represents an attempted build of a GitHub Pages site, whether
// successful or not.
// The Webhook event name is "page_build".
//
// This event is triggered on push to a GitHub Pages enabled branch (gh-pages
// for project pages, master for user and organization pages).
//
// Events of this type are not visible in timelines, they are only used to trigger hooks.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pagebuildevent
type PageBuildEvent struct {
	Build *PagesBuild `json:"build,omitempty"`

	// The following fields are only populated by Webhook events.
	ID           *int64        `json:"id,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// PingEvent is triggered when a Webhook is added to GitHub.
//
// GitHub API docs: https://developer.github.com/webhooks/#ping-event
type PingEvent struct {
	// Random string of GitHub zen.
	Zen *string `json:"zen,omitempty"`
	// The ID of the webhook that triggered the ping.
	HookID *int64 `json:"hook_id,omitempty"`
	// The webhook configuration.
	Hook         *Hook         `json:"hook,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// ProjectEvent is triggered when project is created, modified or deleted.
// The webhook event name is "project".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#projectevent
type ProjectEvent struct {
	Action  *string        `json:"action,omitempty"`
	Changes *ProjectChange `json:"changes,omitempty"`
	Project *Project       `json:"project,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Org          *Organization `json:"organization,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// ProjectCardEvent is triggered when a project card is created, updated, moved, converted to an issue, or deleted.
// The webhook event name is "project_card".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#projectcardevent
type ProjectCardEvent struct {
	Action      *string            `json:"action,omitempty"`
	Changes     *ProjectCardChange `json:"changes,omitempty"`
	AfterID     *int64             `json:"after_id,omitempty"`
	ProjectCard *ProjectCard       `json:"project_card,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Org          *Organization `json:"organization,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// ProjectColumnEvent is triggered when a project column is created, updated, moved, or deleted.
// The webhook event name is "project_column".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#projectcolumnevent
type ProjectColumnEvent struct {
	Action        *string              `json:"action,omitempty"`
	Changes       *ProjectColumnChange `json:"changes,omitempty"`
	AfterID       *int64               `json:"after_id,omitempty"`
	ProjectColumn *ProjectColumn       `json:"project_column,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Org          *Organization `json:"organization,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// PublicEvent is triggered when a private repository is open sourced.
// According to GitHub: "Without a doubt: the best GitHub event."
// The Webhook event name is "public".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#publicevent
type PublicEvent struct {
	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// PullRequestEvent is triggered when a pull request is assigned, unassigned,
// labeled, unlabeled, opened, closed, reopened, or synchronized.
// The Webhook event name is "pull_request".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pullrequestevent
type PullRequestEvent struct {
	// Action is the action that was performed. Possible values are:
	// "assigned", "unassigned", "review_requested", "review_request_removed", "labeled", "unlabeled",
	// "opened", "closed", "reopened", "synchronize", "edited".
	// If the action is "closed" and the merged key is false,
	// the pull request was closed with unmerged commits. If the action is "closed"
	// and the merged key is true, the pull request was merged.
	Action      *string      `json:"action,omitempty"`
	Number      *int         `json:"number,omitempty"`
	PullRequest *PullRequest `json:"pull_request,omitempty"`

	// The following fields are only populated by Webhook events.
	Changes *EditChange `json:"changes,omitempty"`
	// RequestedReviewer is populated in "review_requested", "review_request_removed" event deliveries.
	// A request affecting multiple reviewers at once is split into multiple
	// such event deliveries, each with a single, different RequestedReviewer.
	RequestedReviewer *User         `json:"requested_reviewer,omitempty"`
	Repo              *Repository   `json:"repository,omitempty"`
	Sender            *User         `json:"sender,omitempty"`
	Installation      *Installation `json:"installation,omitempty"`
	Label             *Label        `json:"label,omitempty"` // Populated in "labeled" event deliveries.
}

// PullRequestReviewEvent is triggered when a review is submitted on a pull
// request.
// The Webhook event name is "pull_request_review".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pullrequestreviewevent
type PullRequestReviewEvent struct {
	// Action is always "submitted".
	Action      *string            `json:"action,omitempty"`
	Review      *PullRequestReview `json:"review,omitempty"`
	PullRequest *PullRequest       `json:"pull_request,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`

	// The following field is only present when the webhook is triggered on
	// a repository belonging to an organization.
	Organization *Organization `json:"organization,omitempty"`
}

// PullRequestReviewCommentEvent is triggered when a comment is created on a
// portion of the unified diff of a pull request.
// The Webhook event name is "pull_request_review_comment".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pullrequestreviewcommentevent
type PullRequestReviewCommentEvent struct {
	// Action is the action that was performed on the comment.
	// Possible values are: "created", "edited", "deleted".
	Action      *string             `json:"action,omitempty"`
	PullRequest *PullRequest        `json:"pull_request,omitempty"`
	Comment     *PullRequestComment `json:"comment,omitempty"`

	// The following fields are only populated by Webhook events.
	Changes      *EditChange   `json:"changes,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// PushEvent represents a git push to a GitHub repository.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pushevent
type PushEvent struct {
	PushID       *int64            `json:"push_id,omitempty"`
	Head         *string           `json:"head,omitempty"`
	Ref          *string           `json:"ref,omitempty"`
	Size         *int              `json:"size,omitempty"`
	Commits      []PushEventCommit `json:"commits,omitempty"`
	Before       *string           `json:"before,omitempty"`
	DistinctSize *int              `json:"distinct_size,omitempty"`

	// The following fields are only populated by Webhook events.
	After        *string              `json:"after,omitempty"`
	Created      *bool                `json:"created,omitempty"`
	Deleted      *bool                `json:"deleted,omitempty"`
	Forced       *bool                `json:"forced,omitempty"`
	BaseRef      *string              `json:"base_ref,omitempty"`
	Compare      *string              `json:"compare,omitempty"`
	Repo         *PushEventRepository `json:"repository,omitempty"`
	HeadCommit   *PushEventCommit     `json:"head_commit,omitempty"`
	Pusher       *User                `json:"pusher,omitempty"`
	Sender       *User                `json:"sender,omitempty"`
	Installation *Installation        `json:"installation,omitempty"`
}

func (p PushEvent) String() string {
	return Stringify(p)
}

// PushEventCommit represents a git commit in a GitHub PushEvent.
type PushEventCommit struct {
	Message  *string       `json:"message,omitempty"`
	Author   *CommitAuthor `json:"author,omitempty"`
	URL      *string       `json:"url,omitempty"`
	Distinct *bool         `json:"distinct,omitempty"`

	// The following fields are only populated by Events API.
	SHA *string `json:"sha,omitempty"`

	// The following fields are only populated by Webhook events.
	ID        *string       `json:"id,omitempty"`
	TreeID    *string       `json:"tree_id,omitempty"`
	Timestamp *Timestamp    `json:"timestamp,omitempty"`
	Committer *CommitAuthor `json:"committer,omitempty"`
	Added     []string      `json:"added,omitempty"`
	Removed   []string      `json:"removed,omitempty"`
	Modified  []string      `json:"modified,omitempty"`
}

func (p PushEventCommit) String() string {
	return Stringify(p)
}

// PushEventRepository represents the repo object in a PushEvent payload.
type PushEventRepository struct {
	ID              *int64              `json:"id,omitempty"`
	NodeID          *string             `json:"node_id,omitempty"`
	Name            *string             `json:"name,omitempty"`
	FullName        *string             `json:"full_name,omitempty"`
	Owner           *PushEventRepoOwner `json:"owner,omitempty"`
	Private         *bool               `json:"private,omitempty"`
	Description     *string             `json:"description,omitempty"`
	Fork            *bool               `json:"fork,omitempty"`
	CreatedAt       *Timestamp          `json:"created_at,omitempty"`
	PushedAt        *Timestamp          `json:"pushed_at,omitempty"`
	UpdatedAt       *Timestamp          `json:"updated_at,omitempty"`
	Homepage        *string             `json:"homepage,omitempty"`
	Size            *int                `json:"size,omitempty"`
	StargazersCount *int                `json:"stargazers_count,omitempty"`
	WatchersCount   *int                `json:"watchers_count,omitempty"`
	Language        *string             `json:"language,omitempty"`
	HasIssues       *bool               `json:"has_issues,omitempty"`
	HasDownloads    *bool               `json:"has_downloads,omitempty"`
	HasWiki         *bool               `json:"has_wiki,omitempty"`
	HasPages        *bool               `json:"has_pages,omitempty"`
	ForksCount      *int                `json:"forks_count,omitempty"`
	OpenIssuesCount *int                `json:"open_issues_count,omitempty"`
	DefaultBranch   *string             `json:"default_branch,omitempty"`
	MasterBranch    *string             `json:"master_branch,omitempty"`
	Organization    *string             `json:"organization,omitempty"`
	URL             *string             `json:"url,omitempty"`
	ArchiveURL      *string             `json:"archive_url,omitempty"`
	HTMLURL         *string             `json:"html_url,omitempty"`
	StatusesURL     *string             `json:"statuses_url,omitempty"`
	GitURL          *string             `json:"git_url,omitempty"`
	SSHURL          *string             `json:"ssh_url,omitempty"`
	CloneURL        *string             `json:"clone_url,omitempty"`
	SVNURL          *string             `json:"svn_url,omitempty"`
}

// PushEventRepoOwner is a basic representation of user/org in a PushEvent payload.
type PushEventRepoOwner struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// ReleaseEvent is triggered when a release is published.
// The Webhook event name is "release".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#releaseevent
type ReleaseEvent struct {
	// Action is the action that was performed. Possible value is: "published".
	Action  *string            `json:"action,omitempty"`
	Release *RepositoryRelease `json:"release,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// RepositoryEvent is triggered when a repository is created.
// The Webhook event name is "repository".
//
// Events of this type are not visible in timelines, they are only used to
// trigger organization webhooks.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#repositoryevent
type RepositoryEvent struct {
	// Action is the action that was performed. Possible values are: "created", "deleted",
	// "publicized", "privatized".
	Action *string     `json:"action,omitempty"`
	Repo   *Repository `json:"repository,omitempty"`

	// The following fields are only populated by Webhook events.
	Org          *Organization `json:"organization,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// StatusEvent is triggered when the status of a Git commit changes.
// The Webhook event name is "status".
//
// Events of this type are not visible in timelines, they are only used to
// trigger hooks.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#statusevent
type StatusEvent struct {
	SHA *string `json:"sha,omitempty"`
	// State is the new state. Possible values are: "pending", "success", "failure", "error".
	State       *string   `json:"state,omitempty"`
	Description *string   `json:"description,omitempty"`
	TargetURL   *string   `json:"target_url,omitempty"`
	Branches    []*Branch `json:"branches,omitempty"`

	// The following fields are only populated by Webhook events.
	ID           *int64            `json:"id,omitempty"`
	Name         *string           `json:"name,omitempty"`
	Context      *string           `json:"context,omitempty"`
	Commit       *RepositoryCommit `json:"commit,omitempty"`
	CreatedAt    *Timestamp        `json:"created_at,omitempty"`
	UpdatedAt    *Timestamp        `json:"updated_at,omitempty"`
	Repo         *Repository       `json:"repository,omitempty"`
	Sender       *User             `json:"sender,omitempty"`
	Installation *Installation     `json:"installation,omitempty"`
}

// TeamEvent is triggered when an organization's team is created, modified or deleted.
// The Webhook event name is "team".
//
// Events of this type are not visible in timelines. These events are only used
// to trigger hooks.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#teamevent
type TeamEvent struct {
	Action  *string     `json:"action,omitempty"`
	Team    *Team       `json:"team,omitempty"`
	Changes *TeamChange `json:"changes,omitempty"`
	Repo    *Repository `json:"repository,omitempty"`

	// The following fields are only populated by Webhook events.
	Org          *Organization `json:"organization,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// TeamAddEvent is triggered when a repository is added to a team.
// The Webhook event name is "team_add".
//
// Events of this type are not visible in timelines. These events are only used
// to trigger hooks.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#teamaddevent
type TeamAddEvent struct {
	Team *Team       `json:"team,omitempty"`
	Repo *Repository `json:"repository,omitempty"`

	// The following fields are only populated by Webhook events.
	Org          *Organization `json:"organization,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}

// WatchEvent is related to starring a repository, not watching. See this API
// blog post for an explanation: https://developer.github.com/changes/2012-09-05-watcher-api/
//
// The event’s actor is the user who starred a repository, and the event’s
// repository is the repository that was starred.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#watchevent
type WatchEvent struct {
	// Action is the action that was performed. Possible value is: "started".
	Action *string `json:"action,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
	Installation *Installation `json:"installation,omitempty"`
}
