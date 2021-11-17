package tfe

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	slug "github.com/hashicorp/go-slug"
)

// Compile-time proof of interface implementation.
var _ ConfigurationVersions = (*configurationVersions)(nil)

// ConfigurationVersions describes all the configuration version related
// methods that the Terraform Enterprise API supports.
//
// TFE API docs:
// https://www.terraform.io/docs/enterprise/api/configuration-versions.html
type ConfigurationVersions interface {
	// List returns all configuration versions of a workspace.
	List(ctx context.Context, workspaceID string, options ConfigurationVersionListOptions) (*ConfigurationVersionList, error)

	// Create is used to create a new configuration version. The created
	// configuration version will be usable once data is uploaded to it.
	Create(ctx context.Context, workspaceID string, options ConfigurationVersionCreateOptions) (*ConfigurationVersion, error)

	// Read a configuration version by its ID.
	Read(ctx context.Context, cvID string) (*ConfigurationVersion, error)

	// ReadWithOptions reads a configuration version by its ID using the options supplied
	ReadWithOptions(ctx context.Context, cvID string, options *ConfigurationVersionReadOptions) (*ConfigurationVersion, error)

	// Upload packages and uploads Terraform configuration files. It requires
	// the upload URL from a configuration version and the full path to the
	// configuration files on disk.
	Upload(ctx context.Context, url string, path string) error
}

// configurationVersions implements ConfigurationVersions.
type configurationVersions struct {
	client *Client
}

// ConfigurationStatus represents a configuration version status.
type ConfigurationStatus string

//List all available configuration version statuses.
const (
	ConfigurationErrored  ConfigurationStatus = "errored"
	ConfigurationPending  ConfigurationStatus = "pending"
	ConfigurationUploaded ConfigurationStatus = "uploaded"
)

// ConfigurationSource represents a source of a configuration version.
type ConfigurationSource string

// List all available configuration version sources.
const (
	ConfigurationSourceAPI       ConfigurationSource = "tfe-api"
	ConfigurationSourceBitbucket ConfigurationSource = "bitbucket"
	ConfigurationSourceGithub    ConfigurationSource = "github"
	ConfigurationSourceGitlab    ConfigurationSource = "gitlab"
	ConfigurationSourceTerraform ConfigurationSource = "terraform"
)

// ConfigurationVersionList represents a list of configuration versions.
type ConfigurationVersionList struct {
	*Pagination
	Items []*ConfigurationVersion
}

// ConfigurationVersion is a representation of an uploaded or ingressed
// Terraform configuration in TFE. A workspace must have at least one
// configuration version before any runs may be queued on it.
type ConfigurationVersion struct {
	ID               string              `jsonapi:"primary,configuration-versions"`
	AutoQueueRuns    bool                `jsonapi:"attr,auto-queue-runs"`
	Error            string              `jsonapi:"attr,error"`
	ErrorMessage     string              `jsonapi:"attr,error-message"`
	Source           ConfigurationSource `jsonapi:"attr,source"`
	Speculative      bool                `jsonapi:"attr,speculative "`
	Status           ConfigurationStatus `jsonapi:"attr,status"`
	StatusTimestamps *CVStatusTimestamps `jsonapi:"attr,status-timestamps"`
	UploadURL        string              `jsonapi:"attr,upload-url"`

	// Relations
	IngressAttributes *IngressAttributes `jsonapi:"relation,ingress-attributes"`
}

// CVStatusTimestamps holds the timestamps for individual configuration version
// statuses.
type CVStatusTimestamps struct {
	FinishedAt time.Time `jsonapi:"attr,finished-at,rfc3339"`
	QueuedAt   time.Time `jsonapi:"attr,queued-at,rfc3339"`
	StartedAt  time.Time `jsonapi:"attr,started-at,rfc3339"`
}

// ConfigurationVersionReadOptions represents the options for reading a configuration version.
type ConfigurationVersionReadOptions struct {
	Include string `url:"include"`
}

// ConfigurationVersionListOptions represents the options for listing
// configuration versions.
type ConfigurationVersionListOptions struct {
	ListOptions

	// A list of relations to include. See available resources:
	// https://www.terraform.io/docs/cloud/api/configuration-versions.html#available-related-resources
	Include *string `url:"include"`
}

// IngressAttributes include commit information associated with configuration versions sourced from VCS.
type IngressAttributes struct {
	ID                string `jsonapi:"primary,ingress-attributes"`
	Branch            string `jsonapi:"attr,branch"`
	CloneURL          string `jsonapi:"attr,clone-url"`
	CommitMessage     string `jsonapi:"attr,commit-message"`
	CommitSHA         string `jsonapi:"attr,commit-sha"`
	CommitURL         string `jsonapi:"attr,commit-url"`
	CompareURL        string `jsonapi:"attr,compare-url"`
	Identifier        string `jsonapi:"attr,identifier"`
	IsPullRequest     bool   `jsonapi:"attr,is-pull-request"`
	OnDefaultBranch   bool   `jsonapi:"attr,on-default-branch"`
	PullRequestNumber int    `jsonapi:"attr,pull-request-number"`
	PullRequestURL    string `jsonapi:"attr,pull-request-url"`
	PullRequestTitle  string `jsonapi:"attr,pull-request-title"`
	PullRequestBody   string `jsonapi:"attr,pull-request-body"`
	Tag               string `jsonapi:"attr,tag"`
	SenderUsername    string `jsonapi:"attr,sender-username"`
	SenderAvatarURL   string `jsonapi:"attr,sender-avatar-url"`
	SenderHTMLURL     string `jsonapi:"attr,sender-html-url"`

	// Links
	Links map[string]interface{} `jsonapi:"links,omitempty"`
}

// List returns all configuration versions of a workspace.
func (s *configurationVersions) List(ctx context.Context, workspaceID string, options ConfigurationVersionListOptions) (*ConfigurationVersionList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/configuration-versions", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	cvl := &ConfigurationVersionList{}
	err = s.client.do(ctx, req, cvl)
	if err != nil {
		return nil, err
	}

	return cvl, nil
}

// ConfigurationVersionCreateOptions represents the options for creating a
// configuration version.
type ConfigurationVersionCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,configuration-versions"`

	// When true, runs are queued automatically when the configuration version
	// is uploaded.
	AutoQueueRuns *bool `jsonapi:"attr,auto-queue-runs,omitempty"`

	// When true, this configuration version can only be used for planning.
	Speculative *bool `jsonapi:"attr,speculative,omitempty"`
}

// Create is used to create a new configuration version. The created
// configuration version will be usable once data is uploaded to it.
func (s *configurationVersions) Create(ctx context.Context, workspaceID string, options ConfigurationVersionCreateOptions) (*ConfigurationVersion, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/configuration-versions", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	cv := &ConfigurationVersion{}
	err = s.client.do(ctx, req, cv)
	if err != nil {
		return nil, err
	}

	return cv, nil
}

// Read a configuration version by its ID.
func (s *configurationVersions) Read(ctx context.Context, cvID string) (*ConfigurationVersion, error) {
	return s.ReadWithOptions(ctx, cvID, nil)
}

// Read a configuration version by its ID with the given options.
func (s *configurationVersions) ReadWithOptions(ctx context.Context, cvID string, options *ConfigurationVersionReadOptions) (*ConfigurationVersion, error) {
	if !validStringID(&cvID) {
		return nil, ErrInvalidConfigVersionID
	}

	u := fmt.Sprintf("configuration-versions/%s", url.QueryEscape(cvID))
	req, err := s.client.newRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	cv := &ConfigurationVersion{}
	err = s.client.do(ctx, req, cv)
	if err != nil {
		return nil, err
	}

	return cv, nil
}

// Upload packages and uploads Terraform configuration files. It requires the
// upload URL from a configuration version and the path to the configuration
// files on disk.
func (s *configurationVersions) Upload(ctx context.Context, url, path string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !file.Mode().IsDir() {
		return ErrMissingDirectory
	}

	body := bytes.NewBuffer(nil)

	_, err = slug.Pack(path, body, true)
	if err != nil {
		return err
	}

	req, err := s.client.newRequest("PUT", url, body)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
