// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/jsonapi"
)

// Compile-time proof of interface implementation.
var _ Workspaces = (*workspaces)(nil)

// Workspaces describes all the workspace related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspaces
type Workspaces interface {
	// List all the workspaces within an organization.
	List(ctx context.Context, organization string, options *WorkspaceListOptions) (*WorkspaceList, error)

	// Create is used to create a new workspace.
	Create(ctx context.Context, organization string, options WorkspaceCreateOptions) (*Workspace, error)

	// Read a workspace by its name and organization name.
	Read(ctx context.Context, organization string, workspace string) (*Workspace, error)

	// ReadWithOptions reads a workspace by name and organization name with given options.
	ReadWithOptions(ctx context.Context, organization string, workspace string, options *WorkspaceReadOptions) (*Workspace, error)

	// Readme gets the readme of a workspace by its ID.
	Readme(ctx context.Context, workspaceID string) (io.Reader, error)

	// ReadByID reads a workspace by its ID.
	ReadByID(ctx context.Context, workspaceID string) (*Workspace, error)

	// ReadByIDWithOptions reads a workspace by its ID with the given options.
	ReadByIDWithOptions(ctx context.Context, workspaceID string, options *WorkspaceReadOptions) (*Workspace, error)

	// Update settings of an existing workspace.
	Update(ctx context.Context, organization string, workspace string, options WorkspaceUpdateOptions) (*Workspace, error)

	// UpdateByID updates the settings of an existing workspace.
	UpdateByID(ctx context.Context, workspaceID string, options WorkspaceUpdateOptions) (*Workspace, error)

	// Delete a workspace by its name.
	Delete(ctx context.Context, organization string, workspace string) error

	// DeleteByID deletes a workspace by its ID.
	DeleteByID(ctx context.Context, workspaceID string) error

	// SafeDelete a workspace by its name.
	SafeDelete(ctx context.Context, organization string, workspace string) error

	// SafeDeleteByID deletes a workspace by its ID.
	SafeDeleteByID(ctx context.Context, workspaceID string) error

	// RemoveVCSConnection from a workspace.
	RemoveVCSConnection(ctx context.Context, organization, workspace string) (*Workspace, error)

	// RemoveVCSConnectionByID removes a VCS connection from a workspace.
	RemoveVCSConnectionByID(ctx context.Context, workspaceID string) (*Workspace, error)

	// Lock a workspace by its ID.
	Lock(ctx context.Context, workspaceID string, options WorkspaceLockOptions) (*Workspace, error)

	// Unlock a workspace by its ID.
	Unlock(ctx context.Context, workspaceID string) (*Workspace, error)

	// ForceUnlock a workspace by its ID.
	ForceUnlock(ctx context.Context, workspaceID string) (*Workspace, error)

	// AssignSSHKey to a workspace.
	AssignSSHKey(ctx context.Context, workspaceID string, options WorkspaceAssignSSHKeyOptions) (*Workspace, error)

	// UnassignSSHKey from a workspace.
	UnassignSSHKey(ctx context.Context, workspaceID string) (*Workspace, error)

	// ListRemoteStateConsumers reads the remote state consumers for a workspace.
	ListRemoteStateConsumers(ctx context.Context, workspaceID string, options *RemoteStateConsumersListOptions) (*WorkspaceList, error)

	// AddRemoteStateConsumers adds remote state consumers to a workspace.
	AddRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceAddRemoteStateConsumersOptions) error

	// RemoveRemoteStateConsumers removes remote state consumers from a workspace.
	RemoveRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceRemoveRemoteStateConsumersOptions) error

	// UpdateRemoteStateConsumers updates all the remote state consumers for a workspace
	// to match the workspaces in the update options.
	UpdateRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceUpdateRemoteStateConsumersOptions) error

	// ListTags reads the tags for a workspace.
	ListTags(ctx context.Context, workspaceID string, options *WorkspaceTagListOptions) (*TagList, error)

	// AddTags appends tags to a workspace
	AddTags(ctx context.Context, workspaceID string, options WorkspaceAddTagsOptions) error

	// RemoveTags removes tags from a workspace
	RemoveTags(ctx context.Context, workspaceID string, options WorkspaceRemoveTagsOptions) error

	// ReadDataRetentionPolicy reads a workspace's data retention policy
	//
	// Deprecated: Use ReadDataRetentionPolicyChoice instead.
	// **Note: This functionality is only available in Terraform Enterprise versions v202311-1 and v202312-1.**
	ReadDataRetentionPolicy(ctx context.Context, workspaceID string) (*DataRetentionPolicy, error)

	// ReadDataRetentionPolicyChoice reads a workspace's data retention policy
	// **Note: This functionality is only available in Terraform Enterprise.**
	ReadDataRetentionPolicyChoice(ctx context.Context, workspaceID string) (*DataRetentionPolicyChoice, error)

	// SetDataRetentionPolicy sets a workspace's data retention policy to delete data older than a certain number of days
	//
	// Deprecated: Use SetDataRetentionPolicyDeleteOlder instead
	// **Note: This functionality is only available in Terraform Enterprise versions v202311-1 and v202312-1.**
	SetDataRetentionPolicy(ctx context.Context, workspaceID string, options DataRetentionPolicySetOptions) (*DataRetentionPolicy, error)

	// SetDataRetentionPolicyDeleteOlder sets a workspace's data retention policy to delete data older than a certain number of days
	// **Note: This functionality is only available in Terraform Enterprise.**
	SetDataRetentionPolicyDeleteOlder(ctx context.Context, workspaceID string, options DataRetentionPolicyDeleteOlderSetOptions) (*DataRetentionPolicyDeleteOlder, error)

	// SetDataRetentionPolicyDontDelete sets a workspace's data retention policy to explicitly not delete data
	// **Note: This functionality is only available in Terraform Enterprise.**
	SetDataRetentionPolicyDontDelete(ctx context.Context, workspaceID string, options DataRetentionPolicyDontDeleteSetOptions) (*DataRetentionPolicyDontDelete, error)

	// DeleteDataRetentionPolicy deletes a workspace's data retention policy
	// **Note: This functionality is only available in Terraform Enterprise.**
	DeleteDataRetentionPolicy(ctx context.Context, workspaceID string) error

	// ListTagBindings lists all tag bindings associated with the workspace.
	ListTagBindings(ctx context.Context, workspaceID string) ([]*TagBinding, error)

	// AddTagBindings adds or modifies the value of existing tag binding keys for a workspace.
	AddTagBindings(ctx context.Context, workspaceID string, options WorkspaceAddTagBindingsOptions) ([]*TagBinding, error)
}

// workspaces implements Workspaces.
type workspaces struct {
	client *Client
}

// WorkspaceList represents a list of workspaces.
type WorkspaceList struct {
	*Pagination
	Items []*Workspace
}

// WorkspaceAddTagBindingsOptions represents the options for adding tag bindings
// to a workspace.
type WorkspaceAddTagBindingsOptions struct {
	TagBindings []*TagBinding
}

// LockedByChoice is a choice type struct that represents the possible values
// within a polymorphic relation. If a value is available, exactly one field
// will be non-nil.
type LockedByChoice struct {
	Run  *Run
	User *User
	Team *Team
}

// Workspace represents a Terraform Enterprise workspace.
type Workspace struct {
	ID                          string                          `jsonapi:"primary,workspaces"`
	Actions                     *WorkspaceActions               `jsonapi:"attr,actions"`
	AllowDestroyPlan            bool                            `jsonapi:"attr,allow-destroy-plan"`
	AssessmentsEnabled          bool                            `jsonapi:"attr,assessments-enabled"`
	AutoApply                   bool                            `jsonapi:"attr,auto-apply"`
	AutoApplyRunTrigger         bool                            `jsonapi:"attr,auto-apply-run-trigger"`
	AutoDestroyAt               jsonapi.NullableAttr[time.Time] `jsonapi:"attr,auto-destroy-at,iso8601,omitempty"`
	AutoDestroyActivityDuration jsonapi.NullableAttr[string]    `jsonapi:"attr,auto-destroy-activity-duration,omitempty"`
	CanQueueDestroyPlan         bool                            `jsonapi:"attr,can-queue-destroy-plan"`
	CreatedAt                   time.Time                       `jsonapi:"attr,created-at,iso8601"`
	Description                 string                          `jsonapi:"attr,description"`
	Environment                 string                          `jsonapi:"attr,environment"`
	ExecutionMode               string                          `jsonapi:"attr,execution-mode"`
	FileTriggersEnabled         bool                            `jsonapi:"attr,file-triggers-enabled"`
	GlobalRemoteState           bool                            `jsonapi:"attr,global-remote-state"`
	Locked                      bool                            `jsonapi:"attr,locked"`
	MigrationEnvironment        string                          `jsonapi:"attr,migration-environment"`
	Name                        string                          `jsonapi:"attr,name"`
	NoCodeUpgradeAvailable      bool                            `jsonapi:"attr,no-code-upgrade-available"`
	Operations                  bool                            `jsonapi:"attr,operations"`
	Permissions                 *WorkspacePermissions           `jsonapi:"attr,permissions"`
	QueueAllRuns                bool                            `jsonapi:"attr,queue-all-runs"`
	SpeculativeEnabled          bool                            `jsonapi:"attr,speculative-enabled"`
	SourceName                  string                          `jsonapi:"attr,source-name"`
	SourceURL                   string                          `jsonapi:"attr,source-url"`
	StructuredRunOutputEnabled  bool                            `jsonapi:"attr,structured-run-output-enabled"`
	TerraformVersion            string                          `jsonapi:"attr,terraform-version"`
	TriggerPrefixes             []string                        `jsonapi:"attr,trigger-prefixes"`
	TriggerPatterns             []string                        `jsonapi:"attr,trigger-patterns"`
	VCSRepo                     *VCSRepo                        `jsonapi:"attr,vcs-repo"`
	WorkingDirectory            string                          `jsonapi:"attr,working-directory"`
	UpdatedAt                   time.Time                       `jsonapi:"attr,updated-at,iso8601"`
	ResourceCount               int                             `jsonapi:"attr,resource-count"`
	ApplyDurationAverage        time.Duration                   `jsonapi:"attr,apply-duration-average"`
	PlanDurationAverage         time.Duration                   `jsonapi:"attr,plan-duration-average"`
	PolicyCheckFailures         int                             `jsonapi:"attr,policy-check-failures"`
	RunFailures                 int                             `jsonapi:"attr,run-failures"`
	RunsCount                   int                             `jsonapi:"attr,workspace-kpis-runs-count"`
	TagNames                    []string                        `jsonapi:"attr,tag-names"`
	SettingOverwrites           *WorkspaceSettingOverwrites     `jsonapi:"attr,setting-overwrites"`

	// Relations
	AgentPool                   *AgentPool            `jsonapi:"relation,agent-pool"`
	CurrentRun                  *Run                  `jsonapi:"relation,current-run"`
	CurrentStateVersion         *StateVersion         `jsonapi:"relation,current-state-version"`
	Organization                *Organization         `jsonapi:"relation,organization"`
	SSHKey                      *SSHKey               `jsonapi:"relation,ssh-key"`
	Outputs                     []*WorkspaceOutputs   `jsonapi:"relation,outputs"`
	Project                     *Project              `jsonapi:"relation,project"`
	Tags                        []*Tag                `jsonapi:"relation,tags"`
	CurrentConfigurationVersion *ConfigurationVersion `jsonapi:"relation,current-configuration-version,omitempty"`
	LockedBy                    *LockedByChoice       `jsonapi:"polyrelation,locked-by"`
	Variables                   []*Variable           `jsonapi:"relation,vars"`
	TagBindings                 []*TagBinding         `jsonapi:"relation,tag-bindings"`

	// Deprecated: Use DataRetentionPolicyChoice instead.
	DataRetentionPolicy *DataRetentionPolicy
	// **Note: This functionality is only available in Terraform Enterprise.**
	DataRetentionPolicyChoice *DataRetentionPolicyChoice `jsonapi:"polyrelation,data-retention-policy"`

	// Links
	Links map[string]interface{} `jsonapi:"links,omitempty"`
}

type WorkspaceOutputs struct {
	ID        string      `jsonapi:"primary,workspace-outputs"`
	Name      string      `jsonapi:"attr,name"`
	Sensitive bool        `jsonapi:"attr,sensitive"`
	Type      string      `jsonapi:"attr,output-type"`
	Value     interface{} `jsonapi:"attr,value"`
}

// workspaceWithReadme is the same as a workspace but it has a readme.
type workspaceWithReadme struct {
	ID     string           `jsonapi:"primary,workspaces"`
	Readme *workspaceReadme `jsonapi:"relation,readme"`
}

// workspaceReadme contains the readme of the workspace.
type workspaceReadme struct {
	ID          string `jsonapi:"primary,workspace-readme"`
	RawMarkdown string `jsonapi:"attr,raw-markdown"`
}

// VCSRepo contains the configuration of a VCS integration.
type VCSRepo struct {
	Branch            string `jsonapi:"attr,branch"`
	DisplayIdentifier string `jsonapi:"attr,display-identifier"`
	Identifier        string `jsonapi:"attr,identifier"`
	IngressSubmodules bool   `jsonapi:"attr,ingress-submodules"`
	OAuthTokenID      string `jsonapi:"attr,oauth-token-id"`
	GHAInstallationID string `jsonapi:"attr,github-app-installation-id"`
	RepositoryHTTPURL string `jsonapi:"attr,repository-http-url"`
	ServiceProvider   string `jsonapi:"attr,service-provider"`
	Tags              bool   `jsonapi:"attr,tags"`
	TagsRegex         string `jsonapi:"attr,tags-regex"`
	WebhookURL        string `jsonapi:"attr,webhook-url"`
}

// Note: the fields of this struct are bool pointers instead of bool values, in order to simplify support for
// future TFE versions that support *some but not all* of the inherited defaults that go-tfe knows about.
type WorkspaceSettingOverwrites struct {
	ExecutionMode *bool `jsonapi:"attr,execution-mode"`
	AgentPool     *bool `jsonapi:"attr,agent-pool"`
}

// WorkspaceActions represents the workspace actions.
type WorkspaceActions struct {
	IsDestroyable bool `jsonapi:"attr,is-destroyable"`
}

// WorkspacePermissions represents the workspace permissions.
type WorkspacePermissions struct {
	CanDestroy        bool  `jsonapi:"attr,can-destroy"`
	CanForceUnlock    bool  `jsonapi:"attr,can-force-unlock"`
	CanLock           bool  `jsonapi:"attr,can-lock"`
	CanManageRunTasks bool  `jsonapi:"attr,can-manage-run-tasks"`
	CanQueueApply     bool  `jsonapi:"attr,can-queue-apply"`
	CanQueueDestroy   bool  `jsonapi:"attr,can-queue-destroy"`
	CanQueueRun       bool  `jsonapi:"attr,can-queue-run"`
	CanReadSettings   bool  `jsonapi:"attr,can-read-settings"`
	CanUnlock         bool  `jsonapi:"attr,can-unlock"`
	CanUpdate         bool  `jsonapi:"attr,can-update"`
	CanUpdateVariable bool  `jsonapi:"attr,can-update-variable"`
	CanForceDelete    *bool `jsonapi:"attr,can-force-delete"` // pointer b/c it will be useful to check if this property exists, as opposed to having it default to false
}

// WSIncludeOpt represents the available options for include query params.
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspaces#available-related-resources
type WSIncludeOpt string

const (
	WSOrganization               WSIncludeOpt = "organization"
	WSCurrentConfigVer           WSIncludeOpt = "current_configuration_version"
	WSCurrentConfigVerIngress    WSIncludeOpt = "current_configuration_version.ingress_attributes"
	WSCurrentRun                 WSIncludeOpt = "current_run"
	WSCurrentRunPlan             WSIncludeOpt = "current_run.plan"
	WSCurrentRunConfigVer        WSIncludeOpt = "current_run.configuration_version"
	WSCurrentrunConfigVerIngress WSIncludeOpt = "current_run.configuration_version.ingress_attributes"
	WSLockedBy                   WSIncludeOpt = "locked_by"
	WSReadme                     WSIncludeOpt = "readme"
	WSOutputs                    WSIncludeOpt = "outputs"
	WSCurrentStateVer            WSIncludeOpt = "current-state-version"
	WSProject                    WSIncludeOpt = "project"
)

// WorkspaceReadOptions represents the options for reading a workspace.
type WorkspaceReadOptions struct {
	// Optional: A list of relations to include.
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspaces#available-related-resources
	Include []WSIncludeOpt `url:"include,omitempty"`
}

// WorkspaceListOptions represents the options for listing workspaces.
type WorkspaceListOptions struct {
	ListOptions

	// Optional: A search string (partial workspace name) used to filter the results.
	Search string `url:"search[name],omitempty"`

	// Optional: A search string (comma-separated tag names) used to filter the results.
	Tags string `url:"search[tags],omitempty"`

	// Optional: A search string (comma-separated tag names to exclude) used to filter the results.
	ExcludeTags string `url:"search[exclude-tags],omitempty"`

	// Optional: A search on substring matching to filter the results.
	WildcardName string `url:"search[wildcard-name],omitempty"`

	// Optional: A filter string to list all the workspaces linked to a given project id in the organization.
	ProjectID string `url:"filter[project][id],omitempty"`

	// Optional: A filter string to list all the workspaces filtered by current run status.
	CurrentRunStatus string `url:"filter[current-run][status],omitempty"`

	// Optional: A filter string to list workspaces filtered by key/value tags.
	// These are not annotated and therefore not encoded by go-querystring
	TagBindings []*TagBinding

	// Optional: A list of relations to include. See available resources https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspaces#available-related-resources
	Include []WSIncludeOpt `url:"include,omitempty"`

	// Optional: May sort on "name" (the default) and "current-run.created-at" (which sorts by the time of the current run)
	// Prepending a hyphen to the sort parameter will reverse the order (e.g. "-name" to reverse the default order)
	Sort string `url:"sort,omitempty"`
}

// WorkspaceCreateOptions represents the options for creating a new workspace.
type WorkspaceCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,workspaces"`

	// Required when: execution-mode is set to agent. The ID of the agent pool
	// belonging to the workspace's organization. This value must not be specified
	// if execution-mode is set to remote or local or if operations is set to true.
	AgentPoolID *string `jsonapi:"attr,agent-pool-id,omitempty"`

	// Optional: Whether destroy plans can be queued on the workspace.
	AllowDestroyPlan *bool `jsonapi:"attr,allow-destroy-plan,omitempty"`

	// Optional: Whether to enable health assessments (drift detection etc.) for the workspace.
	// Reference: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspaces#create-a-workspace
	// Requires remote execution mode, HCP Terraform Business entitlement, and a valid agent pool to work
	AssessmentsEnabled *bool `jsonapi:"attr,assessments-enabled,omitempty"`

	// Optional: Whether to automatically apply changes when a Terraform plan is successful.
	AutoApply *bool `jsonapi:"attr,auto-apply,omitempty"`

	// Optional: Whether to automatically apply changes for runs that are created by run triggers
	// from another workspace.
	AutoApplyRunTrigger *bool `jsonapi:"attr,auto-apply-run-trigger,omitempty"`

	// Optional: The time after which an automatic destroy run will be queued
	AutoDestroyAt jsonapi.NullableAttr[time.Time] `jsonapi:"attr,auto-destroy-at,iso8601,omitempty"`

	// Optional: The period of time to wait after workspace activity to trigger a destroy run. The format
	// should roughly match a Go duration string limited to days and hours, e.g. "24h" or "1d".
	AutoDestroyActivityDuration jsonapi.NullableAttr[string] `jsonapi:"attr,auto-destroy-activity-duration,omitempty"`

	// Optional: A description for the workspace.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Optional: Which execution mode to use. Valid values are remote, local, and agent.
	// When set to local, the workspace will be used for state storage only.
	// This value must not be specified if operations is specified.
	// 'agent' execution mode is not available in Terraform Enterprise.
	ExecutionMode *string `jsonapi:"attr,execution-mode,omitempty"`

	// Optional: Whether to filter runs based on the changed files in a VCS push. If
	// enabled, the working directory and trigger prefixes describe a set of
	// paths which must contain changes for a VCS push to trigger a run. If
	// disabled, any push will trigger a run.
	FileTriggersEnabled *bool `jsonapi:"attr,file-triggers-enabled,omitempty"`

	GlobalRemoteState *bool `jsonapi:"attr,global-remote-state,omitempty"`

	// Optional: The legacy TFE environment to use as the source of the migration, in the
	// form organization/environment. Omit this unless you are migrating a legacy
	// environment.
	MigrationEnvironment *string `jsonapi:"attr,migration-environment,omitempty"`

	// The name of the workspace, which can only include letters, numbers, -,
	// and _. This will be used as an identifier and must be unique in the
	// organization.
	Name *string `jsonapi:"attr,name"`

	// DEPRECATED. Whether the workspace will use remote or local execution mode.
	// Use ExecutionMode instead.
	Operations *bool `jsonapi:"attr,operations,omitempty"`

	// Whether to queue all runs. Unless this is set to true, runs triggered by
	// a webhook will not be queued until at least one run is manually queued.
	QueueAllRuns *bool `jsonapi:"attr,queue-all-runs,omitempty"`

	// Whether this workspace allows speculative plans. Setting this to false
	// prevents HCP Terraform or the Terraform Enterprise instance from
	// running plans on pull requests, which can improve security if the VCS
	// repository is public or includes untrusted contributors.
	SpeculativeEnabled *bool `jsonapi:"attr,speculative-enabled,omitempty"`

	// BETA. A friendly name for the application or client creating this
	// workspace. If set, this will be displayed on the workspace as
	// "Created via <SOURCE NAME>".
	SourceName *string `jsonapi:"attr,source-name,omitempty"`

	// BETA. A URL for the application or client creating this workspace. This
	// can be the URL of a related resource in another app, or a link to
	// documentation or other info about the client.
	SourceURL *string `jsonapi:"attr,source-url,omitempty"`

	// BETA. Enable the experimental advanced run user interface.
	// This only applies to runs using Terraform version 0.15.2 or newer,
	// and runs executed using older versions will see the classic experience
	// regardless of this setting.
	StructuredRunOutputEnabled *bool `jsonapi:"attr,structured-run-output-enabled,omitempty"`

	// The version of Terraform to use for this workspace. Upon creating a
	// workspace, the latest version is selected unless otherwise specified.
	TerraformVersion *string `jsonapi:"attr,terraform-version,omitempty"`

	// List of repository-root-relative paths which list all locations to be
	// tracked for changes. See FileTriggersEnabled above for more details.
	TriggerPrefixes []string `jsonapi:"attr,trigger-prefixes,omitempty"`

	// Optional: List of patterns used to match against changed files in order
	// to decide whether to trigger a run or not.
	TriggerPatterns []string `jsonapi:"attr,trigger-patterns,omitempty"`

	// Settings for the workspace's VCS repository. If omitted, the workspace is
	// created without a VCS repo. If included, you must specify at least the
	// oauth-token-id and identifier keys below.
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo,omitempty"`

	// A relative path that Terraform will execute within. This defaults to the
	// root of your repository and is typically set to a subdirectory matching the
	// environment when multiple environments exist within the same repository.
	WorkingDirectory *string `jsonapi:"attr,working-directory,omitempty"`

	// A list of tags to attach to the workspace. If the tag does not already
	// exist, it is created and added to the workspace.
	Tags []*Tag `jsonapi:"relation,tags,omitempty"`

	// Optional: Struct of booleans, which indicate whether the workspace
	// specifies its own values for various settings. If you mark a setting as
	// `false` in this struct, it will clear the workspace's existing value for
	// that setting and defer to the default value that its project or
	// organization provides.
	//
	// In general, it's not necessary to mark a setting as `true` in this
	// struct; if you provide a literal value for a setting, HCP Terraform will
	// automatically update its overwrites field to `true`. If you do choose to
	// manually mark a setting as overwritten, you must provide a value for that
	// setting at the same time.
	SettingOverwrites *WorkspaceSettingOverwritesOptions `jsonapi:"attr,setting-overwrites,omitempty"`

	// Associated Project with the workspace. If not provided, default project
	// of the organization will be assigned to the workspace.
	Project *Project `jsonapi:"relation,project,omitempty"`

	// Associated TagBindings of the workspace.
	TagBindings []*TagBinding `jsonapi:"relation,tag-bindings,omitempty"`
}

// TODO: move this struct out. VCSRepoOptions is used by workspaces, policy sets, and registry modules
// VCSRepoOptions represents the configuration options of a VCS integration.
type VCSRepoOptions struct {
	Branch            *string `json:"branch,omitempty"`
	Identifier        *string `json:"identifier,omitempty"`
	IngressSubmodules *bool   `json:"ingress-submodules,omitempty"`
	OAuthTokenID      *string `json:"oauth-token-id,omitempty"`
	TagsRegex         *string `json:"tags-regex,omitempty"`
	GHAInstallationID *string `json:"github-app-installation-id,omitempty"`
}

type WorkspaceSettingOverwritesOptions struct {
	// If false, the workspace will defer to its organization or project's DefaultExecutionMode value.
	ExecutionMode *bool `json:"execution-mode,omitempty"`
	// If false, the workspace will defer to its organization or project's DefaultAgentPool value.
	AgentPool *bool `json:"agent-pool,omitempty"`
}

// WorkspaceUpdateOptions represents the options for updating a workspace.
type WorkspaceUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,workspaces"`

	// Required when: execution-mode is set to agent. The ID of the agent pool
	// belonging to the workspace's organization. This value must not be specified
	// if execution-mode is set to remote or local or if operations is set to true.
	AgentPoolID *string `jsonapi:"attr,agent-pool-id,omitempty"`

	// Optional: Whether destroy plans can be queued on the workspace.
	AllowDestroyPlan *bool `jsonapi:"attr,allow-destroy-plan,omitempty"`

	// Optional: Whether to enable health assessments (drift detection etc.) for the workspace.
	// Reference: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/workspaces#update-a-workspace
	// Requires remote execution mode, HCP Terraform Business entitlement, and a valid agent pool to work
	AssessmentsEnabled *bool `jsonapi:"attr,assessments-enabled,omitempty"`

	// Optional: Whether to automatically apply changes when a Terraform plan is successful.
	AutoApply *bool `jsonapi:"attr,auto-apply,omitempty"`

	// Optional: Whether to automatically apply changes for runs that are created by run triggers
	// from another workspace.
	AutoApplyRunTrigger *bool `jsonapi:"attr,auto-apply-run-trigger,omitempty"`

	// Optional: The time after which an automatic destroy run will be queued
	AutoDestroyAt jsonapi.NullableAttr[time.Time] `jsonapi:"attr,auto-destroy-at,iso8601,omitempty"`

	// Optional: The period of time to wait after workspace activity to trigger a destroy run. The format
	// should roughly match a Go duration string limited to days and hours, e.g. "24h" or "1d".
	AutoDestroyActivityDuration jsonapi.NullableAttr[string] `jsonapi:"attr,auto-destroy-activity-duration,omitempty"`

	// Optional: A new name for the workspace, which can only include letters, numbers, -,
	// and _. This will be used as an identifier and must be unique in the
	// organization. Warning: Changing a workspace's name changes its URL in the
	// API and UI.
	Name *string `jsonapi:"attr,name,omitempty"`

	// Optional: A description for the workspace.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Optional: Which execution mode to use. Valid values are remote, local, and agent.
	// When set to local, the workspace will be used for state storage only.
	// This value must not be specified if operations is specified.
	// 'agent' execution mode is not available in Terraform Enterprise.
	ExecutionMode *string `jsonapi:"attr,execution-mode,omitempty"`

	// Optional: Whether to filter runs based on the changed files in a VCS push. If
	// enabled, the working directory and trigger prefixes describe a set of
	// paths which must contain changes for a VCS push to trigger a run. If
	// disabled, any push will trigger a run.
	FileTriggersEnabled *bool `jsonapi:"attr,file-triggers-enabled,omitempty"`

	// Optional:
	GlobalRemoteState *bool `jsonapi:"attr,global-remote-state,omitempty"`

	// DEPRECATED. Whether the workspace will use remote or local execution mode.
	// Use ExecutionMode instead.
	Operations *bool `jsonapi:"attr,operations,omitempty"`

	// Optional: Whether to queue all runs. Unless this is set to true, runs triggered by
	// a webhook will not be queued until at least one run is manually queued.
	QueueAllRuns *bool `jsonapi:"attr,queue-all-runs,omitempty"`

	// Optional: Whether this workspace allows speculative plans. Setting this to false
	// prevents HCP Terraform or the Terraform Enterprise instance from
	// running plans on pull requests, which can improve security if the VCS
	// repository is public or includes untrusted contributors.
	SpeculativeEnabled *bool `jsonapi:"attr,speculative-enabled,omitempty"`

	// BETA. Enable the experimental advanced run user interface.
	// This only applies to runs using Terraform version 0.15.2 or newer,
	// and runs executed using older versions will see the classic experience
	// regardless of this setting.
	StructuredRunOutputEnabled *bool `jsonapi:"attr,structured-run-output-enabled,omitempty"`

	// Optional: The version of Terraform to use for this workspace.
	TerraformVersion *string `jsonapi:"attr,terraform-version,omitempty"`

	// Optional: List of repository-root-relative paths which list all locations to be
	// tracked for changes. See FileTriggersEnabled above for more details.
	TriggerPrefixes []string `jsonapi:"attr,trigger-prefixes,omitempty"`

	// Optional: List of patterns used to match against changed files in order
	// to decide whether to trigger a run or not.
	TriggerPatterns []string `jsonapi:"attr,trigger-patterns,omitempty"`

	// Optional: To delete a workspace's existing VCS repo, specify null instead of an
	// object. To modify a workspace's existing VCS repo, include whichever of
	// the keys below you wish to modify. To add a new VCS repo to a workspace
	// that didn't previously have one, include at least the oauth-token-id and
	// identifier keys.
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo,omitempty"`

	// Optional: A relative path that Terraform will execute within. This defaults to the
	// root of your repository and is typically set to a subdirectory matching
	// the environment when multiple environments exist within the same
	// repository.
	WorkingDirectory *string `jsonapi:"attr,working-directory,omitempty"`

	// Optional: Struct of booleans, which indicate whether the workspace
	// specifies its own values for various settings. If you mark a setting as
	// `false` in this struct, it will clear the workspace's existing value for
	// that setting and defer to the default value that its project or
	// organization provides.
	//
	// In general, it's not necessary to mark a setting as `true` in this
	// struct; if you provide a literal value for a setting, HCP Terraform will
	// automatically update its overwrites field to `true`. If you do choose to
	// manually mark a setting as overwritten, you must provide a value for that
	// setting at the same time.
	SettingOverwrites *WorkspaceSettingOverwritesOptions `jsonapi:"attr,setting-overwrites,omitempty"`

	// Associated Project with the workspace. If not provided, default project
	// of the organization will be assigned to the workspace
	Project *Project `jsonapi:"relation,project,omitempty"`

	// Associated TagBindings of the project. Note that this will replace
	// all existing tag bindings.
	TagBindings []*TagBinding `jsonapi:"relation,tag-bindings,omitempty"`
}

// WorkspaceLockOptions represents the options for locking a workspace.
type WorkspaceLockOptions struct {
	// Specifies the reason for locking the workspace.
	Reason *string `jsonapi:"attr,reason,omitempty"`
}

// workspaceRemoveVCSConnectionOptions
type workspaceRemoveVCSConnectionOptions struct {
	ID      string          `jsonapi:"primary,workspaces"`
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo"`
}

// WorkspaceAssignSSHKeyOptions represents the options to assign an SSH key to
// a workspace.
type WorkspaceAssignSSHKeyOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,workspaces"`

	// The SSH key ID to assign.
	SSHKeyID *string `jsonapi:"attr,id"`
}

// workspaceUnassignSSHKeyOptions represents the options to unassign an SSH key
// to a workspace.
type workspaceUnassignSSHKeyOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,workspaces"`

	// Must be nil to unset the currently assigned SSH key.
	SSHKeyID *string `jsonapi:"attr,id"`
}

type RemoteStateConsumersListOptions struct {
	ListOptions
}

// WorkspaceAddRemoteStateConsumersOptions represents the options for adding remote state consumers
// to a workspace.
type WorkspaceAddRemoteStateConsumersOptions struct {
	// The workspaces to add as remote state consumers to the workspace.
	Workspaces []*Workspace
}

// WorkspaceRemoveRemoteStateConsumersOptions represents the options for removing remote state
// consumers from a workspace.
type WorkspaceRemoveRemoteStateConsumersOptions struct {
	// The workspaces to remove as remote state consumers from the workspace.
	Workspaces []*Workspace
}

// WorkspaceUpdateRemoteStateConsumersOptions represents the options for
// updatintg remote state consumers from a workspace.
type WorkspaceUpdateRemoteStateConsumersOptions struct {
	// The workspaces to update remote state consumers for the workspace.
	Workspaces []*Workspace
}

type WorkspaceTagListOptions struct {
	ListOptions

	// A query string used to filter workspace tags.
	// Any workspace tag with a name partially matching this value will be returned.
	Query *string `url:"name,omitempty"`
}

type WorkspaceAddTagsOptions struct {
	Tags []*Tag
}

type WorkspaceRemoveTagsOptions struct {
	Tags []*Tag
}

// List all the workspaces within an organization.
func (s *workspaces) List(ctx context.Context, organization string, options *WorkspaceListOptions) (*WorkspaceList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	var tagFilters map[string][]string
	if options != nil {
		tagFilters = encodeTagFiltersAsParams(options.TagBindings)
	}

	// Encode parameters that cannot be encoded by go-querystring
	u := fmt.Sprintf("organizations/%s/workspaces", url.PathEscape(organization))
	req, err := s.client.NewRequestWithAdditionalQueryParams("GET", u, options, tagFilters)
	if err != nil {
		return nil, err
	}

	wl := &WorkspaceList{}
	err = req.Do(ctx, wl)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

func (s *workspaces) ListTagBindings(ctx context.Context, workspaceID string) ([]*TagBinding, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/tag-bindings", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var list struct {
		*Pagination
		Items []*TagBinding
	}

	err = req.Do(ctx, &list)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

// AddTagBindings adds or modifies the value of existing tag binding keys for a workspace.
func (s *workspaces) AddTagBindings(ctx context.Context, workspaceID string, options WorkspaceAddTagBindingsOptions) ([]*TagBinding, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("workspaces/%s/tag-bindings", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("PATCH", u, options.TagBindings)
	if err != nil {
		return nil, err
	}

	var response = struct {
		*Pagination
		Items []*TagBinding
	}{}
	err = req.Do(ctx, &response)

	return response.Items, err
}

// Create is used to create a new workspace.
func (s *workspaces) Create(ctx context.Context, organization string, options WorkspaceCreateOptions) (*Workspace, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/workspaces", url.PathEscape(organization))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Read a workspace by its name and organization name.
func (s *workspaces) Read(ctx context.Context, organization, workspace string) (*Workspace, error) {
	return s.ReadWithOptions(ctx, organization, workspace, nil)
}

// ReadWithOptions reads a workspace by name and organization name with given options.
func (s *workspaces) ReadWithOptions(ctx context.Context, organization, workspace string, options *WorkspaceReadOptions) (*Workspace, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if !validStringID(&workspace) {
		return nil, ErrInvalidWorkspaceValue
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf(
		"organizations/%s/workspaces/%s",
		url.PathEscape(organization),
		url.PathEscape(workspace),
	)
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	// Manually populate the deprecated DataRetentionPolicy field
	w.DataRetentionPolicy = w.DataRetentionPolicyChoice.ConvertToLegacyStruct()

	// durations come over in ms
	w.ApplyDurationAverage *= time.Millisecond
	w.PlanDurationAverage *= time.Millisecond

	return w, nil
}

// ReadByID reads a workspace by its ID.
func (s *workspaces) ReadByID(ctx context.Context, workspaceID string) (*Workspace, error) {
	return s.ReadByIDWithOptions(ctx, workspaceID, nil)
}

// ReadByIDWithOptions reads a workspace by its ID with the given options.
func (s *workspaces) ReadByIDWithOptions(ctx context.Context, workspaceID string, options *WorkspaceReadOptions) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	// Manually populate the deprecated DataRetentionPolicy field
	if w.DataRetentionPolicyChoice != nil {
		w.DataRetentionPolicy = w.DataRetentionPolicyChoice.ConvertToLegacyStruct()
	}

	// durations come over in ms
	w.ApplyDurationAverage *= time.Millisecond
	w.PlanDurationAverage *= time.Millisecond

	return w, nil
}

// Readme gets the readme of a workspace by its ID.
func (s *workspaces) Readme(ctx context.Context, workspaceID string) (io.Reader, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s?include=readme", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	r := &workspaceWithReadme{}
	err = req.Do(ctx, r)
	if err != nil {
		return nil, err
	}
	if r.Readme == nil {
		return nil, nil
	}

	return strings.NewReader(r.Readme.RawMarkdown), nil
}

// Update settings of an existing workspace.
func (s *workspaces) Update(ctx context.Context, organization, workspace string, options WorkspaceUpdateOptions) (*Workspace, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if !validStringID(&workspace) {
		return nil, ErrInvalidWorkspaceValue
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf(
		"organizations/%s/workspaces/%s",
		url.PathEscape(organization),
		url.PathEscape(workspace),
	)
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// UpdateByID updates the settings of an existing workspace.
func (s *workspaces) UpdateByID(ctx context.Context, workspaceID string, options WorkspaceUpdateOptions) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Delete a workspace by its name.
func (s *workspaces) Delete(ctx context.Context, organization, workspace string) error {
	if !validStringID(&organization) {
		return ErrInvalidOrg
	}
	if !validStringID(&workspace) {
		return ErrInvalidWorkspaceValue
	}

	u := fmt.Sprintf(
		"organizations/%s/workspaces/%s",
		url.PathEscape(organization),
		url.PathEscape(workspace),
	)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// DeleteByID deletes a workspace by its ID.
func (s *workspaces) DeleteByID(ctx context.Context, workspaceID string) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// SafeDelete a workspace by its name.
func (s *workspaces) SafeDelete(ctx context.Context, organization, workspace string) error {
	if !validStringID(&organization) {
		return ErrInvalidOrg
	}
	if !validStringID(&workspace) {
		return ErrInvalidWorkspaceValue
	}

	u := fmt.Sprintf(
		"organizations/%s/workspaces/%s/actions/safe-delete",
		url.PathEscape(organization),
		url.PathEscape(workspace),
	)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// SafeDeleteByID safely deletes a workspace by its ID.
func (s *workspaces) SafeDeleteByID(ctx context.Context, workspaceID string) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/actions/safe-delete", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// RemoveVCSConnection from a workspace.
func (s *workspaces) RemoveVCSConnection(ctx context.Context, organization, workspace string) (*Workspace, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if !validStringID(&workspace) {
		return nil, ErrInvalidWorkspaceValue
	}

	u := fmt.Sprintf(
		"organizations/%s/workspaces/%s",
		url.PathEscape(organization),
		url.PathEscape(workspace),
	)

	req, err := s.client.NewRequest("PATCH", u, &workspaceRemoveVCSConnectionOptions{})
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// RemoveVCSConnectionByID removes a VCS connection from a workspace.
func (s *workspaces) RemoveVCSConnectionByID(ctx context.Context, workspaceID string) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s", url.PathEscape(workspaceID))

	req, err := s.client.NewRequest("PATCH", u, &workspaceRemoveVCSConnectionOptions{})
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Lock a workspace by its ID.
func (s *workspaces) Lock(ctx context.Context, workspaceID string, options WorkspaceLockOptions) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/actions/lock", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Unlock a workspace by its ID.
func (s *workspaces) Unlock(ctx context.Context, workspaceID string) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/actions/unlock", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		if strings.Contains(err.Error(), "latest state version is still pending") {
			return nil, ErrWorkspaceLockedStateVersionStillPending
		}
		return nil, err
	}

	return w, nil
}

// ForceUnlock a workspace by its ID.
func (s *workspaces) ForceUnlock(ctx context.Context, workspaceID string) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/actions/force-unlock", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// AssignSSHKey to a workspace.
func (s *workspaces) AssignSSHKey(ctx context.Context, workspaceID string, options WorkspaceAssignSSHKeyOptions) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("workspaces/%s/relationships/ssh-key", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// UnassignSSHKey from a workspace.
func (s *workspaces) UnassignSSHKey(ctx context.Context, workspaceID string) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/relationships/ssh-key", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("PATCH", u, &workspaceUnassignSSHKeyOptions{})
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = req.Do(ctx, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// RemoteStateConsumers returns the remote state consumers for a given workspace.
func (s *workspaces) ListRemoteStateConsumers(ctx context.Context, workspaceID string, options *RemoteStateConsumersListOptions) (*WorkspaceList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/relationships/remote-state-consumers", url.PathEscape(workspaceID))

	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	wl := &WorkspaceList{}
	err = req.Do(ctx, wl)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

// AddRemoteStateConsumere adds the remote state consumers to a given workspace.
func (s *workspaces) AddRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceAddRemoteStateConsumersOptions) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("workspaces/%s/relationships/remote-state-consumers", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("POST", u, options.Workspaces)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// RemoveRemoteStateConsumers removes the remote state consumers for a given workspace.
func (s *workspaces) RemoveRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceRemoveRemoteStateConsumersOptions) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("workspaces/%s/relationships/remote-state-consumers", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("DELETE", u, options.Workspaces)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// UpdateRemoteStateConsumers removes the remote state consumers for a given workspace.
func (s *workspaces) UpdateRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceUpdateRemoteStateConsumersOptions) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("workspaces/%s/relationships/remote-state-consumers", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("PATCH", u, options.Workspaces)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// ListTags returns the tags for a given workspace.
func (s *workspaces) ListTags(ctx context.Context, workspaceID string, options *WorkspaceTagListOptions) (*TagList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/relationships/tags", url.PathEscape(workspaceID))

	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	tl := &TagList{}
	err = req.Do(ctx, tl)
	if err != nil {
		return nil, err
	}

	return tl, nil
}

// AddTags adds a list of tags to a workspace.
func (s *workspaces) AddTags(ctx context.Context, workspaceID string, options WorkspaceAddTagsOptions) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("workspaces/%s/relationships/tags", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("POST", u, options.Tags)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// RemoveTags removes a list of tags from a workspace.
func (s *workspaces) RemoveTags(ctx context.Context, workspaceID string, options WorkspaceRemoveTagsOptions) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("workspaces/%s/relationships/tags", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("DELETE", u, options.Tags)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (s *workspaces) ReadDataRetentionPolicy(ctx context.Context, workspaceID string) (*DataRetentionPolicy, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/relationships/data-retention-policy", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	dataRetentionPolicy := &DataRetentionPolicy{}
	err = req.Do(ctx, dataRetentionPolicy)

	if err != nil {
		// try to detect known issue where this function is used with TFE >= 202401,
		// and direct user towards the V2 function
		if drpUnmarshalEr.MatchString(err.Error()) {
			return nil, fmt.Errorf("error reading deprecated DataRetentionPolicy, use ReadDataRetentionPolicyChoice instead")
		}
		return nil, err
	}

	return dataRetentionPolicy, nil
}

func (s *workspaces) ReadDataRetentionPolicyChoice(ctx context.Context, workspaceID string) (*DataRetentionPolicyChoice, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	// The API to read the drp is workspaces/<id>/relationships/data-retention-policy
	// However, this API can return multiple "types" (e.g. data-retention-policy-delete-olders, or data-retention-policy-dont-deletes)
	// Ideally we would deserialize this directly into the choice type (DataRetentionPolicyChoice)...however, there isn't a way to
	// tell the current jsonapi implementation that the direct result of an endpoint could be different types. Relationships can be polymorphic,
	// but the direct result of an endpoint can't be (as far as the jsonapi implementation is concerned)

	// Instead, we need to figure out the type of the data retention policy first, and deserialize it into the matching model. We
	// can then create a choice type manually
	ws, err := s.ReadByID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	// there is no drp (of a known type)
	if ws.DataRetentionPolicyChoice == nil || !ws.DataRetentionPolicyChoice.IsPopulated() {
		return ws.DataRetentionPolicyChoice, nil
	}

	u := s.dataRetentionPolicyLink(workspaceID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	dataRetentionPolicy := &DataRetentionPolicyChoice{}
	// if reading the workspace told us it was a "delete older policy" deserialize into the DeleteOlder portion of the choice model
	if ws.DataRetentionPolicyChoice.DataRetentionPolicyDeleteOlder != nil {
		deleteOlder := &DataRetentionPolicyDeleteOlder{}
		err = req.Do(ctx, deleteOlder)
		dataRetentionPolicy.DataRetentionPolicyDeleteOlder = deleteOlder

		// if reading the workspace told us it was a "delete older policy" deserialize into the DeleteOlder portion of the choice model
	} else if ws.DataRetentionPolicyChoice.DataRetentionPolicyDontDelete != nil {
		dontDelete := &DataRetentionPolicyDontDelete{}
		err = req.Do(ctx, dontDelete)
		dataRetentionPolicy.DataRetentionPolicyDontDelete = dontDelete
	} else if ws.DataRetentionPolicyChoice != nil {
		legacyDrp := &DataRetentionPolicy{}
		err = req.Do(ctx, legacyDrp)
		dataRetentionPolicy.DataRetentionPolicy = legacyDrp
	}

	if err != nil {
		return nil, err
	}

	return dataRetentionPolicy, nil
}

func (s *workspaces) SetDataRetentionPolicy(ctx context.Context, workspaceID string, options DataRetentionPolicySetOptions) (*DataRetentionPolicy, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := s.dataRetentionPolicyLink(workspaceID)
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	dataRetentionPolicy := &DataRetentionPolicy{}
	err = req.Do(ctx, dataRetentionPolicy)

	if err != nil {
		return nil, err
	}

	return dataRetentionPolicy, nil
}

func (s *workspaces) SetDataRetentionPolicyDeleteOlder(ctx context.Context, workspaceID string, options DataRetentionPolicyDeleteOlderSetOptions) (*DataRetentionPolicyDeleteOlder, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := s.dataRetentionPolicyLink(workspaceID)
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	dataRetentionPolicy := &DataRetentionPolicyDeleteOlder{}
	err = req.Do(ctx, dataRetentionPolicy)

	if err != nil {
		return nil, err
	}

	return dataRetentionPolicy, nil
}

func (s *workspaces) SetDataRetentionPolicyDontDelete(ctx context.Context, workspaceID string, options DataRetentionPolicyDontDeleteSetOptions) (*DataRetentionPolicyDontDelete, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := s.dataRetentionPolicyLink(workspaceID)
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	dataRetentionPolicy := &DataRetentionPolicyDontDelete{}
	err = req.Do(ctx, dataRetentionPolicy)

	if err != nil {
		return nil, err
	}

	return dataRetentionPolicy, nil
}

func (s *workspaces) DeleteDataRetentionPolicy(ctx context.Context, workspaceID string) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}

	u := s.dataRetentionPolicyLink(workspaceID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o WorkspaceAddTagBindingsOptions) valid() error {
	if len(o.TagBindings) == 0 {
		return ErrRequiredTagBindings
	}

	return nil
}

func (o WorkspaceCreateOptions) valid() error {
	if !validString(o.Name) {
		return ErrRequiredName
	}
	if !validStringID(o.Name) {
		return ErrInvalidName
	}
	if o.Operations != nil && o.ExecutionMode != nil {
		return ErrUnsupportedOperations
	}
	if o.AgentPoolID != nil && (o.ExecutionMode == nil || *o.ExecutionMode != "agent") {
		return ErrRequiredAgentMode
	}
	if o.AgentPoolID == nil && (o.ExecutionMode != nil && *o.ExecutionMode == "agent") {
		return ErrRequiredAgentPoolID
	}
	if len(o.TriggerPrefixes) > 0 &&
		o.TriggerPatterns != nil && len(o.TriggerPatterns) > 0 {
		return ErrUnsupportedBothTriggerPatternsAndPrefixes
	}
	if tagRegexDefined(o.VCSRepo) &&
		o.TriggerPatterns != nil && len(o.TriggerPatterns) > 0 {
		return ErrUnsupportedBothTagsRegexAndTriggerPatterns
	}
	if tagRegexDefined(o.VCSRepo) &&
		o.TriggerPrefixes != nil && len(o.TriggerPrefixes) > 0 {
		return ErrUnsupportedBothTagsRegexAndTriggerPrefixes
	}
	if tagRegexDefined(o.VCSRepo) &&
		o.FileTriggersEnabled != nil && *o.FileTriggersEnabled {
		return ErrUnsupportedBothTagsRegexAndFileTriggersEnabled
	}

	return nil
}

func (o WorkspaceUpdateOptions) valid() error {
	if o.Name != nil && !validStringID(o.Name) {
		return ErrInvalidName
	}
	if o.Operations != nil && o.ExecutionMode != nil {
		return ErrUnsupportedOperations
	}
	if o.AgentPoolID == nil && (o.ExecutionMode != nil && *o.ExecutionMode == "agent") {
		return ErrRequiredAgentPoolID
	}
	if len(o.TriggerPrefixes) > 0 &&
		o.TriggerPatterns != nil && len(o.TriggerPatterns) > 0 {
		return ErrUnsupportedBothTriggerPatternsAndPrefixes
	}

	if tagRegexDefined(o.VCSRepo) &&
		o.TriggerPatterns != nil && len(o.TriggerPatterns) > 0 {
		return ErrUnsupportedBothTagsRegexAndTriggerPatterns
	}
	if tagRegexDefined(o.VCSRepo) &&
		o.TriggerPrefixes != nil && len(o.TriggerPrefixes) > 0 {
		return ErrUnsupportedBothTagsRegexAndTriggerPrefixes
	}
	if tagRegexDefined(o.VCSRepo) &&
		o.FileTriggersEnabled != nil && *o.FileTriggersEnabled {
		return ErrUnsupportedBothTagsRegexAndFileTriggersEnabled
	}

	return nil
}

func (o WorkspaceAssignSSHKeyOptions) valid() error {
	if !validString(o.SSHKeyID) {
		return ErrRequiredSHHKeyID
	}
	if !validStringID(o.SSHKeyID) {
		return ErrInvalidSHHKeyID
	}
	return nil
}

func (o WorkspaceAddRemoteStateConsumersOptions) valid() error {
	if o.Workspaces == nil {
		return ErrWorkspacesRequired
	}
	if len(o.Workspaces) == 0 {
		return ErrWorkspaceMinLimit
	}
	return nil
}

func (o WorkspaceRemoveRemoteStateConsumersOptions) valid() error {
	if o.Workspaces == nil {
		return ErrWorkspacesRequired
	}
	if len(o.Workspaces) == 0 {
		return ErrWorkspaceMinLimit
	}
	return nil
}

func (o WorkspaceUpdateRemoteStateConsumersOptions) valid() error {
	if o.Workspaces == nil {
		return ErrWorkspacesRequired
	}
	if len(o.Workspaces) == 0 {
		return ErrWorkspaceMinLimit
	}
	return nil
}

func (o WorkspaceAddTagsOptions) valid() error {
	if len(o.Tags) == 0 {
		return ErrMissingTagIdentifier
	}
	for _, s := range o.Tags {
		if s.Name == "" && s.ID == "" {
			return ErrMissingTagIdentifier
		}
	}

	return nil
}

func (o WorkspaceRemoveTagsOptions) valid() error {
	if len(o.Tags) == 0 {
		return ErrMissingTagIdentifier
	}
	for _, s := range o.Tags {
		if s.Name == "" && s.ID == "" {
			return ErrMissingTagIdentifier
		}
	}

	return nil
}

func (o *WorkspaceListOptions) valid() error {
	return nil
}

func (o *WorkspaceReadOptions) valid() error {
	return nil
}

func tagRegexDefined(options *VCSRepoOptions) bool {
	if options == nil {
		return false
	}
	if options.TagsRegex != nil && *options.TagsRegex != "" {
		return true
	}
	return false
}

func (s *workspaces) dataRetentionPolicyLink(wsID string) string {
	return fmt.Sprintf("workspaces/%s/relationships/data-retention-policy", url.PathEscape(wsID))
}
