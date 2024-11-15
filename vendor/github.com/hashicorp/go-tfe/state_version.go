// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

// Compile-time proof of interface implementation.
var _ StateVersions = (*stateVersions)(nil)

// StateVersionStatus are available state version status values
type StateVersionStatus string

// Available state version statuses.
const (
	StateVersionPending   StateVersionStatus = "pending"
	StateVersionFinalized StateVersionStatus = "finalized"
	StateVersionDiscarded StateVersionStatus = "discarded"
)

// StateVersions describes all the state version related methods that
// the Terraform Enterprise API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/state-versions
type StateVersions interface {
	// List all the state versions for a given workspace.
	List(ctx context.Context, options *StateVersionListOptions) (*StateVersionList, error)

	// Create a new state version for the given workspace.
	Create(ctx context.Context, workspaceID string, options StateVersionCreateOptions) (*StateVersion, error)

	// Upload creates a new state version but uploads the state content directly to the object store.
	// This is a more resilient form of Create and is the recommended approach to creating state versions.
	Upload(ctx context.Context, workspaceID string, options StateVersionUploadOptions) (*StateVersion, error)

	// Read a state version by its ID.
	Read(ctx context.Context, svID string) (*StateVersion, error)

	// ReadWithOptions reads a state version by its ID using the options supplied
	ReadWithOptions(ctx context.Context, svID string, options *StateVersionReadOptions) (*StateVersion, error)

	// ReadCurrent reads the latest available state from the given workspace.
	ReadCurrent(ctx context.Context, workspaceID string) (*StateVersion, error)

	// ReadCurrentWithOptions reads the latest available state from the given workspace using the options supplied
	ReadCurrentWithOptions(ctx context.Context, workspaceID string, options *StateVersionCurrentOptions) (*StateVersion, error)

	// Download retrieves the actual stored state of a state version
	Download(ctx context.Context, url string) ([]byte, error)

	// ListOutputs retrieves all the outputs of a state version by its ID. IMPORTANT: HCP Terraform might
	// process outputs asynchronously. When consuming outputs or other async StateVersion fields, be sure to
	// wait for ResourcesProcessed to become `true` before assuming they are empty.
	ListOutputs(ctx context.Context, svID string, options *StateVersionOutputsListOptions) (*StateVersionOutputsList, error)

	// SoftDeleteBackingData soft deletes the state version's backing data
	// **Note: This functionality is only available in Terraform Enterprise.**
	SoftDeleteBackingData(ctx context.Context, svID string) error

	// RestoreBackingData restores a soft deleted state version's backing data
	// **Note: This functionality is only available in Terraform Enterprise.**
	RestoreBackingData(ctx context.Context, svID string) error

	// PermanentlyDeleteBackingData permanently deletes a soft deleted state version's backing data
	// **Note: This functionality is only available in Terraform Enterprise.**
	PermanentlyDeleteBackingData(ctx context.Context, svID string) error
}

// stateVersions implements StateVersions.
type stateVersions struct {
	client *Client
}

// StateVersionList represents a list of state versions.
type StateVersionList struct {
	*Pagination
	Items []*StateVersion
}

// StateVersion represents a Terraform Enterprise state version.
type StateVersion struct {
	ID               string             `jsonapi:"primary,state-versions"`
	CreatedAt        time.Time          `jsonapi:"attr,created-at,iso8601"`
	DownloadURL      string             `jsonapi:"attr,hosted-state-download-url"`
	UploadURL        string             `jsonapi:"attr,hosted-state-upload-url"`
	Status           StateVersionStatus `jsonapi:"attr,status"`
	JSONUploadURL    string             `jsonapi:"attr,hosted-json-state-upload-url"`
	JSONDownloadURL  string             `jsonapi:"attr,hosted-json-state-download-url"`
	Serial           int64              `jsonapi:"attr,serial"`
	VCSCommitSHA     string             `jsonapi:"attr,vcs-commit-sha"`
	VCSCommitURL     string             `jsonapi:"attr,vcs-commit-url"`
	BillableRUMCount *uint32            `jsonapi:"attr,billable-rum-count"`
	// Whether HCP Terraform has finished populating any StateVersion fields that required async processing.
	// If `false`, some fields may appear empty even if they should actually contain data; see comments on
	// individual fields for details.
	ResourcesProcessed bool `jsonapi:"attr,resources-processed"`
	StateVersion       int  `jsonapi:"attr,state-version"`
	// Populated asynchronously.
	TerraformVersion string `jsonapi:"attr,terraform-version"`
	// Populated asynchronously.
	Modules *StateVersionModules `jsonapi:"attr,modules"`
	// Populated asynchronously.
	Providers *StateVersionProviders `jsonapi:"attr,providers"`
	// Populated asynchronously.
	Resources []*StateVersionResources `jsonapi:"attr,resources"`

	// Relations
	Run     *Run                  `jsonapi:"relation,run"`
	Outputs []*StateVersionOutput `jsonapi:"relation,outputs"`
}

// StateVersionOutputsList represents a list of StateVersionOutput items.
type StateVersionOutputsList struct {
	*Pagination
	Items []*StateVersionOutput
}

// StateVersionListOptions represents the options for listing state versions.
type StateVersionListOptions struct {
	ListOptions
	Organization string `url:"filter[organization][name]"`
	Workspace    string `url:"filter[workspace][name]"`
}

// StateVersionIncludeOpt represents the available options for include query params.
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/state-versions#available-related-resources
type StateVersionIncludeOpt string

const (
	SVcreatedby               StateVersionIncludeOpt = "created_by"
	SVrun                     StateVersionIncludeOpt = "run"
	SVrunCreatedBy            StateVersionIncludeOpt = "run.created_by"
	SVrunConfigurationVersion StateVersionIncludeOpt = "run.configuration_version"
	SVoutputs                 StateVersionIncludeOpt = "outputs"
)

// StateVersionReadOptions represents the options for reading state version.
type StateVersionReadOptions struct {
	// Optional: A list of relations to include. See available resources:
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/state-versions#available-related-resources
	Include []StateVersionIncludeOpt `url:"include,omitempty"`
}

// StateVersionOutputsListOptions represents the options for listing state
// version outputs.
type StateVersionOutputsListOptions struct {
	ListOptions
}

// StateVersionCurrentOptions represents the options for reading the current state version.
type StateVersionCurrentOptions struct {
	// Optional: A list of relations to include. See available resources:
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/state-versions#available-related-resources
	Include []StateVersionIncludeOpt `url:"include,omitempty"`
}

// StateVersionCreateOptions represents the options for creating a state version.
type StateVersionCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,state-versions"`

	// Optional: The lineage of the state.
	Lineage *string `jsonapi:"attr,lineage,omitempty"`

	// Required: The MD5 hash of the state version.
	MD5 *string `jsonapi:"attr,md5"`

	// Required: The serial of the state.
	Serial *int64 `jsonapi:"attr,serial"`

	// Optional: The base64 encoded state.
	State *string `jsonapi:"attr,state,omitempty"`

	// Optional: Force can be set to skip certain validations. Wrong use
	// of this flag can cause data loss, so USE WITH CAUTION!
	Force *bool `jsonapi:"attr,force,omitempty"`

	// Optional: Specifies the run to associate the state with.
	Run *Run `jsonapi:"relation,run,omitempty"`

	// Optional: The external, json representation of state data, base64 encoded.
	// https://developer.hashicorp.com/terraform/internals/json-format#state-representation
	// Supplying this state representation can provide more details to the platform
	// about the current terraform state.
	JSONState *string `jsonapi:"attr,json-state,omitempty"`
	// Optional: The external, json representation of state outputs, base64 encoded. Supplying this field
	// will provide more detailed output type information to TFE.
	// For more information on the contents of this field: https://developer.hashicorp.com/terraform/internals/json-format#values-representation
	// about the current terraform state.
	JSONStateOutputs *string `jsonapi:"attr,json-state-outputs,omitempty"`
}

type StateVersionUploadOptions struct {
	StateVersionCreateOptions

	RawState     []byte
	RawJSONState []byte
}

type StateVersionModules struct {
	Root StateVersionModuleRoot `jsonapi:"attr,root"`
}

type StateVersionModuleRoot struct {
	NullResource         int `jsonapi:"attr,null-resource"`
	TerraformRemoteState int `jsonapi:"attr,data.terraform-remote-state"`
}

type StateVersionProviders struct {
	Data ProviderData `jsonapi:"attr,provider[map]string"`
}

type ProviderData struct {
	NullResource         int `json:"null-resource"`
	TerraformRemoteState int `json:"data.terraform-remote-state"`
}

type StateVersionResources struct {
	Name     string `jsonapi:"attr,name"`
	Count    int    `jsonapi:"attr,count"`
	Type     string `jsonapi:"attr,type"`
	Module   string `jsonapi:"attr,module"`
	Provider string `jsonapi:"attr,provider"`
}

// List all the state versions for a given workspace.
func (s *stateVersions) List(ctx context.Context, options *StateVersionListOptions) (*StateVersionList, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", "state-versions", options)
	if err != nil {
		return nil, err
	}

	svl := &StateVersionList{}
	err = req.Do(ctx, svl)
	if err != nil {
		return nil, err
	}

	return svl, nil
}

// Create a new state version for the given workspace.
func (s *stateVersions) Create(ctx context.Context, workspaceID string, options StateVersionCreateOptions) (*StateVersion, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("workspaces/%s/state-versions", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	sv := &StateVersion{}
	err = req.Do(ctx, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

// Upload creates a new state version but uploads the state content directly to the object store.
// This is a more resilient form of Create and is the recommended approach to creating state versions.
func (s *stateVersions) Upload(ctx context.Context, workspaceID string, options StateVersionUploadOptions) (*StateVersion, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	sv, err := s.Create(ctx, workspaceID, options.StateVersionCreateOptions)
	if err != nil {
		if strings.Contains(err.Error(), "param is missing or the value is empty: state") {
			return nil, ErrStateVersionUploadNotSupported
		}
		return nil, err
	}

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.client.doForeignPUTRequest(ctx, sv.UploadURL, bytes.NewReader(options.RawState))
	})
	if options.RawJSONState != nil {
		g.Go(func() error {
			return s.client.doForeignPUTRequest(ctx, sv.JSONUploadURL, bytes.NewReader(options.RawJSONState))
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Re-read the state version to get the updated status, if available
	return s.Read(ctx, sv.ID)
}

// Read a state version by its ID.
func (s *stateVersions) ReadWithOptions(ctx context.Context, svID string, options *StateVersionReadOptions) (*StateVersion, error) {
	if !validStringID(&svID) {
		return nil, ErrInvalidStateVerID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("state-versions/%s", url.PathEscape(svID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	sv := &StateVersion{}
	err = req.Do(ctx, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

// Read a state version by its ID.
func (s *stateVersions) Read(ctx context.Context, svID string) (*StateVersion, error) {
	return s.ReadWithOptions(ctx, svID, nil)
}

// ReadCurrentWithOptions reads the latest available state from the given workspace using the options supplied.
func (s *stateVersions) ReadCurrentWithOptions(ctx context.Context, workspaceID string, options *StateVersionCurrentOptions) (*StateVersion, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("workspaces/%s/current-state-version", url.PathEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	sv := &StateVersion{}
	err = req.Do(ctx, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

// ReadCurrent reads the latest available state from the given workspace.
func (s *stateVersions) ReadCurrent(ctx context.Context, workspaceID string) (*StateVersion, error) {
	return s.ReadCurrentWithOptions(ctx, workspaceID, nil)
}

// Download retrieves the actual stored state of a state version
func (s *stateVersions) Download(ctx context.Context, u string) ([]byte, error) {
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	var buf bytes.Buffer
	err = req.Do(ctx, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ListOutputs retrieves all the outputs of a state version by its ID. IMPORTANT: HCP Terraform might
// process outputs asynchronously. When consuming outputs or other async StateVersion fields, be sure to
// wait for ResourcesProcessed to become `true` before assuming they are empty.
func (s *stateVersions) ListOutputs(ctx context.Context, svID string, options *StateVersionOutputsListOptions) (*StateVersionOutputsList, error) {
	if !validStringID(&svID) {
		return nil, ErrInvalidStateVerID
	}

	u := fmt.Sprintf("state-versions/%s/outputs", url.PathEscape(svID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	sv := &StateVersionOutputsList{}
	err = req.Do(ctx, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

func (s *stateVersions) SoftDeleteBackingData(ctx context.Context, svID string) error {
	return s.manageBackingData(ctx, svID, "soft_delete_backing_data")
}

func (s *stateVersions) RestoreBackingData(ctx context.Context, svID string) error {
	return s.manageBackingData(ctx, svID, "restore_backing_data")
}

func (s *stateVersions) PermanentlyDeleteBackingData(ctx context.Context, svID string) error {
	return s.manageBackingData(ctx, svID, "permanently_delete_backing_data")
}

func (s *stateVersions) manageBackingData(ctx context.Context, svID, action string) error {
	if !validStringID(&svID) {
		return ErrInvalidStateVerID
	}

	u := fmt.Sprintf("state-versions/%s/actions/%s", svID, action)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// check that StateVersionListOptions fields had valid values
func (o *StateVersionListOptions) valid() error {
	if o == nil {
		return ErrRequiredStateVerListOps
	}
	if !validString(&o.Organization) {
		return ErrRequiredOrg
	}
	if !validString(&o.Workspace) {
		return ErrRequiredWorkspace
	}
	return nil
}

func (o StateVersionCreateOptions) valid() error {
	if !validString(o.MD5) {
		return ErrRequiredM5
	}
	if o.Serial == nil {
		return ErrRequiredSerial
	}
	return nil
}

func (o StateVersionUploadOptions) valid() error {
	if err := o.StateVersionCreateOptions.valid(); err != nil {
		return err
	}
	if o.State != nil || o.JSONState != nil {
		return ErrStateMustBeOmitted
	}
	if o.RawState == nil {
		return ErrRequiredRawState
	}
	return nil
}

func (o *StateVersionReadOptions) valid() error {
	return nil
}
func (o *StateVersionCurrentOptions) valid() error {
	return nil
}
