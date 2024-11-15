// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Compile-time proof of interface implementation.
var _ RegistryModules = (*registryModules)(nil)

// RegistryModules describes all the registry module related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/private-registry/modules
type RegistryModules interface {
	// List all the registory modules within an organization
	List(ctx context.Context, organization string, options *RegistryModuleListOptions) (*RegistryModuleList, error)

	// ListCommits List the commits for the registry module
	// This returns the latest 20 commits for the connected VCS repo.
	// Pagination is not applicable due to inconsistent support from the VCS providers.
	ListCommits(ctx context.Context, moduleID RegistryModuleID) (*CommitList, error)

	// Create a registry module without a VCS repo
	Create(ctx context.Context, organization string, options RegistryModuleCreateOptions) (*RegistryModule, error)

	// Create a registry module version
	CreateVersion(ctx context.Context, moduleID RegistryModuleID, options RegistryModuleCreateVersionOptions) (*RegistryModuleVersion, error)

	// Create and publish a registry module with a VCS repo
	CreateWithVCSConnection(ctx context.Context, options RegistryModuleCreateWithVCSConnectionOptions) (*RegistryModule, error)

	// Read a registry module
	Read(ctx context.Context, moduleID RegistryModuleID) (*RegistryModule, error)

	// ReadVersion Read a registry module version
	ReadVersion(ctx context.Context, moduleID RegistryModuleID, version string) (*RegistryModuleVersion, error)

	// Delete a registry module
	// Warning: This method is deprecated and will be removed from a future version of go-tfe. Use DeleteByName instead.
	Delete(ctx context.Context, organization string, name string) error

	// Delete a registry module by name
	DeleteByName(ctx context.Context, module RegistryModuleID) error

	// Delete a specified provider for the given module along with all its versions
	DeleteProvider(ctx context.Context, moduleID RegistryModuleID) error

	// Delete a specified version for the given provider of the module
	DeleteVersion(ctx context.Context, moduleID RegistryModuleID, version string) error

	// Update properties of a registry module
	Update(ctx context.Context, moduleID RegistryModuleID, options RegistryModuleUpdateOptions) (*RegistryModule, error)

	// Upload Terraform configuration files for the provided registry module version. It
	// requires a path to the configuration files on disk, which will be packaged by
	// hashicorp/go-slug before being uploaded.
	Upload(ctx context.Context, rmv RegistryModuleVersion, path string) error

	// Upload a tar gzip archive to the specified configuration version upload URL.
	UploadTarGzip(ctx context.Context, url string, r io.Reader) error
}

// registryModules implements RegistryModules.
type registryModules struct {
	client *Client
}

// RegistryModuleStatus represents the status of the registry module
type RegistryModuleStatus string

// List of available registry module statuses
const (
	RegistryModuleStatusPending       RegistryModuleStatus = "pending"
	RegistryModuleStatusNoVersionTags RegistryModuleStatus = "no_version_tags"
	RegistryModuleStatusSetupFailed   RegistryModuleStatus = "setup_failed"
	RegistryModuleStatusSetupComplete RegistryModuleStatus = "setup_complete"
)

// RegistryModuleVersionStatus represents the status of a specific version of a registry module
type RegistryModuleVersionStatus string

// List of available registry module version statuses
const (
	RegistryModuleVersionStatusPending             RegistryModuleVersionStatus = "pending"
	RegistryModuleVersionStatusCloning             RegistryModuleVersionStatus = "cloning"
	RegistryModuleVersionStatusCloneFailed         RegistryModuleVersionStatus = "clone_failed"
	RegistryModuleVersionStatusRegIngressReqFailed RegistryModuleVersionStatus = "reg_ingress_req_failed"
	RegistryModuleVersionStatusRegIngressing       RegistryModuleVersionStatus = "reg_ingressing"
	RegistryModuleVersionStatusRegIngressFailed    RegistryModuleVersionStatus = "reg_ingress_failed"
	RegistryModuleVersionStatusOk                  RegistryModuleVersionStatus = "ok"
)

type PublishingMechanism string

const (
	PublishingMechanismBranch PublishingMechanism = "branch"
	PublishingMechanismTag    PublishingMechanism = "git_tag"
)

// RegistryModuleID represents the set of IDs that identify a RegistryModule
// Use NewPublicRegistryModuleID or NewPrivateRegistryModuleID to build one

type RegistryModuleID struct {
	// The unique ID of the module. If given, the other fields are ignored.
	ID string
	// The organization the module belongs to, see RegistryModule.Organization.Name
	Organization string
	// The name of the module, see RegistryModule.Name
	Name string
	// The module's provider, see RegistryModule.Provider
	Provider string
	// The namespace of the module. For private modules this is the name of the organization that owns the module
	// Required for public modules
	Namespace string
	// Either public or private. If not provided, defaults to private
	RegistryName RegistryName
}

// RegistryModuleList represents a list of registry modules.
type RegistryModuleList struct {
	*Pagination
	Items []*RegistryModule
}

// CommitList represents a list of the latest commits from the registry module
type CommitList struct {
	*Pagination
	Items []*Commit
}

// RegistryModule represents a registry module
type RegistryModule struct {
	ID                  string                          `jsonapi:"primary,registry-modules"`
	Name                string                          `jsonapi:"attr,name"`
	Provider            string                          `jsonapi:"attr,provider"`
	RegistryName        RegistryName                    `jsonapi:"attr,registry-name"`
	Namespace           string                          `jsonapi:"attr,namespace"`
	NoCode              bool                            `jsonapi:"attr,no-code"`
	Permissions         *RegistryModulePermissions      `jsonapi:"attr,permissions"`
	PublishingMechanism PublishingMechanism             `jsonapi:"attr,publishing-mechanism"`
	Status              RegistryModuleStatus            `jsonapi:"attr,status"`
	TestConfig          *TestConfig                     `jsonapi:"attr,test-config"`
	VCSRepo             *VCSRepo                        `jsonapi:"attr,vcs-repo"`
	VersionStatuses     []RegistryModuleVersionStatuses `jsonapi:"attr,version-statuses"`
	CreatedAt           string                          `jsonapi:"attr,created-at"`
	UpdatedAt           string                          `jsonapi:"attr,updated-at"`

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`
}

// Commit represents a commit
type Commit struct {
	ID              string `jsonapi:"primary,commit"`
	Sha             string `jsonapi:"attr,sha"`
	Date            string `jsonapi:"attr,date"`
	URL             string `jsonapi:"attr,url"`
	Author          string `jsonapi:"attr,author"`
	AuthorAvatarURL string `jsonapi:"attr,author-avatar-url"`
	AuthorHTMLURL   string `jsonapi:"attr,author-html-url"`
	Message         string `jsonapi:"attr,message"`
}

// RegistryModuleVersion represents a registry module version
type RegistryModuleVersion struct {
	ID        string                      `jsonapi:"primary,registry-module-versions"`
	Source    string                      `jsonapi:"attr,source"`
	Status    RegistryModuleVersionStatus `jsonapi:"attr,status"`
	Version   string                      `jsonapi:"attr,version"`
	CreatedAt string                      `jsonapi:"attr,created-at"`
	UpdatedAt string                      `jsonapi:"attr,updated-at"`

	// Relations
	RegistryModule *RegistryModule `jsonapi:"relation,registry-module"`

	// Links
	Links map[string]interface{} `jsonapi:"links,omitempty"`
}

type RegistryModulePermissions struct {
	CanDelete bool `jsonapi:"attr,can-delete"`
	CanResync bool `jsonapi:"attr,can-resync"`
	CanRetry  bool `jsonapi:"attr,can-retry"`
}

type RegistryModuleVersionStatuses struct {
	Version string                      `jsonapi:"attr,version"`
	Status  RegistryModuleVersionStatus `jsonapi:"attr,status"`
	Error   string                      `jsonapi:"attr,error"`
}

// RegistryModuleListOptions represents the options for listing registry modules.
type RegistryModuleListOptions struct {
	ListOptions
}

// RegistryModuleCreateOptions is used when creating a registry module without a VCS repo
type RegistryModuleCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,registry-modules"`
	// Required:
	Name *string `jsonapi:"attr,name"`
	// Required:
	Provider *string `jsonapi:"attr,provider"`
	// Optional: Whether this is a publicly maintained module or private. Must be either public or private.
	// Defaults to private if not specified
	RegistryName RegistryName `jsonapi:"attr,registry-name,omitempty"`
	// Optional: The namespace of this module. Required for public modules only.
	Namespace string `jsonapi:"attr,namespace,omitempty"`
	// Optional: If set to true the module is enabled for no-code provisioning.
	// **Note: This field is still in BETA and subject to change.**
	NoCode *bool `jsonapi:"attr,no-code,omitempty"`
}

// RegistryModuleCreateVersionOptions is used when creating a registry module version
type RegistryModuleCreateVersionOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,registry-module-versions"`

	Version *string `jsonapi:"attr,version"`

	CommitSHA *string `jsonapi:"attr,commit-sha"`
}

// RegistryModuleCreateWithVCSConnectionOptions is used when creating a registry module with a VCS repo
type RegistryModuleCreateWithVCSConnectionOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,registry-modules"`

	// Required: VCS repository information
	VCSRepo *RegistryModuleVCSRepoOptions `jsonapi:"attr,vcs-repo"`

	// Optional: If Branch is set within VCSRepo then InitialVersion sets the
	// initial version of the newly created branch-based registry module. If
	// Branch is not set within VCSRepo then InitialVersion is ignored.
	//
	// Defaults to "0.0.0".
	//
	// **Note: This field is still in BETA and subject to change.**
	InitialVersion *string `jsonapi:"attr,initial-version,omitempty"`

	// Optional: Flag to enable tests for the module
	// **Note: This field is still in BETA and subject to change.**
	TestConfig *RegistryModuleTestConfigOptions `jsonapi:"attr,test-config,omitempty"`
}

// RegistryModuleCreateVersionOptions is used when updating a registry module
type RegistryModuleUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-updating
	Type string `jsonapi:"primary,registry-modules"`

	// Optional: Flag to enable no-code provisioning for the whole module.
	// **Note: This field is still in BETA and subject to change.**
	NoCode *bool `jsonapi:"attr,no-code,omitempty"`

	// Optional: Flag to enable tests for the module
	// **Note: This field is still in BETA and subject to change.**
	TestConfig *RegistryModuleTestConfigOptions `jsonapi:"attr,test-config,omitempty"`

	VCSRepo *RegistryModuleVCSRepoUpdateOptions `jsonapi:"attr,vcs-repo,omitempty"`
}

type RegistryModuleTestConfigOptions struct {
	TestsEnabled *bool `jsonapi:"attr,tests-enabled,omitempty"`
}

type RegistryModuleVCSRepoOptions struct {
	Identifier        *string `json:"identifier"` // Required
	OAuthTokenID      *string `json:"oauth-token-id,omitempty"`
	DisplayIdentifier *string `json:"display-identifier,omitempty"` // Required
	GHAInstallationID *string `json:"github-app-installation-id,omitempty"`
	OrganizationName  *string `json:"organization-name,omitempty"`

	// Optional: If set, the newly created registry module will be branch-based
	// with the starting branch set to Branch.
	//
	// **Note: This field is still in BETA and subject to change.**
	Branch *string `json:"branch,omitempty"`
	Tags   *bool   `json:"tags,omitempty"`
}

type RegistryModuleVCSRepoUpdateOptions struct {
	// The Branch and Tag fields are used to determine
	// the PublishingMechanism for a RegistryModule that has a VCS a connection.
	// When a value for Branch is provided, the Tags field is removed on the server
	// When a value for Tags is provided, the Branch field is removed on the server
	// **Note: This field is still in BETA and subject to change.**
	Branch *string `json:"branch,omitempty"`
	Tags   *bool   `json:"tags,omitempty"`
}

// List all the registory modules within an organization.
func (r *registryModules) List(ctx context.Context, organization string, options *RegistryModuleListOptions) (*RegistryModuleList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/registry-modules", url.PathEscape(organization))
	req, err := r.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	ml := &RegistryModuleList{}
	err = req.Do(ctx, ml)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

// List the last 20 commits for the registry modules within an organization.
func (r *registryModules) ListCommits(ctx context.Context, moduleID RegistryModuleID) (*CommitList, error) {
	if !validStringID(&moduleID.Organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf(
		"organizations/%s/registry-modules/private/%s/%s/%s/commits",
		url.PathEscape(moduleID.Organization),
		url.PathEscape(moduleID.Organization),
		url.PathEscape(moduleID.Name),
		url.PathEscape(moduleID.Provider),
	)
	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	cl := &CommitList{}
	err = req.Do(ctx, cl)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

// Upload uploads Terraform configuration files for the provided registry module version. It
// requires a path to the configuration files on disk, which will be packaged by
// hashicorp/go-slug before being uploaded.
func (r *registryModules) Upload(ctx context.Context, rmv RegistryModuleVersion, path string) error {
	uploadURL, ok := rmv.Links["upload"].(string)
	if !ok {
		return fmt.Errorf("provided RegistryModuleVersion does not contain an upload link")
	}

	body, err := packContents(path)
	if err != nil {
		return err
	}

	return r.UploadTarGzip(ctx, uploadURL, body)
}

// UploadTarGzip is used to upload Terraform configuration files contained a tar gzip archive.
// Any stream implementing io.Reader can be passed into this method. This method is also
// particularly useful for tar streams created by non-default go-slug configurations.
//
// **Note**: This method does not validate the content being uploaded and is therefore the caller's
// responsibility to ensure the raw content is a valid Terraform configuration.
func (r *registryModules) UploadTarGzip(ctx context.Context, uploadURL string, archive io.Reader) error {
	return r.client.doForeignPUTRequest(ctx, uploadURL, archive)
}

// Create a new registry module without a VCS repo
func (r *registryModules) Create(ctx context.Context, organization string, options RegistryModuleCreateOptions) (*RegistryModule, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	if options.NoCode != nil {
		log.Println("[WARN] Support for using the NoCode field is deprecated as of release 1.22.0 and may be removed in a future version. The preferred way to create a no-code module is with the registryNoCodeModules.Create method.")
	}

	u := fmt.Sprintf(
		"organizations/%s/registry-modules",
		url.PathEscape(organization),
	)
	req, err := r.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	rm := &RegistryModule{}
	err = req.Do(ctx, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

func (r *registryModules) Update(ctx context.Context, moduleID RegistryModuleID, options RegistryModuleUpdateOptions) (*RegistryModule, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}

	if moduleID.RegistryName == "" {
		log.Println("[WARN] Support for using the RegistryModuleID without RegistryName is deprecated as of release 1.5.0 and may be removed in a future version. The preferred method is to include the RegistryName in RegistryModuleID.")
		moduleID.RegistryName = PrivateRegistry
	}

	if moduleID.RegistryName == PrivateRegistry && strings.TrimSpace(moduleID.Namespace) == "" {
		log.Println("[WARN] Support for using the RegistryModuleID without Namespace is deprecated as of release 1.5.0 and may be removed in a future version. The preferred method is to include the Namespace in RegistryModuleID.")
		moduleID.Namespace = moduleID.Organization
	}

	if options.NoCode != nil {
		log.Println("[WARN] Support for using the NoCode field is deprecated as of release 1.22.0 and may be removed in a future version. The preferred way to update a no-code module is with the registryNoCodeModules.Update method.")
	}

	if options.VCSRepo != nil {
		if options.VCSRepo.Tags != nil && *options.VCSRepo.Tags && validString(options.VCSRepo.Branch) {
			return nil, ErrBranchMustBeEmptyWhenTagsEnabled
		}
	}

	org := url.PathEscape(moduleID.Organization)
	registryName := url.PathEscape(string(moduleID.RegistryName))
	namespace := url.PathEscape(moduleID.Namespace)
	name := url.PathEscape(moduleID.Name)
	provider := url.PathEscape(moduleID.Provider)
	registryModuleURL := fmt.Sprintf("organizations/%s/registry-modules/%s/%s/%s/%s", org, registryName, namespace, name, provider)

	req, err := r.client.NewRequest(http.MethodPatch, registryModuleURL, &options)
	if err != nil {
		return nil, err
	}

	rm := &RegistryModule{}
	if err := req.Do(ctx, rm); err != nil {
		return nil, err
	}

	return rm, nil
}

// CreateVersion creates a new registry module version
func (r *registryModules) CreateVersion(ctx context.Context, moduleID RegistryModuleID, options RegistryModuleCreateVersionOptions) (*RegistryModuleVersion, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf(
		"registry-modules/%s/%s/%s/versions",
		url.PathEscape(moduleID.Organization),
		url.PathEscape(moduleID.Name),
		url.PathEscape(moduleID.Provider),
	)
	req, err := r.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	rmv := &RegistryModuleVersion{}
	err = req.Do(ctx, rmv)
	if err != nil {
		return nil, err
	}

	return rmv, nil
}

// CreateWithVCSConnection is used to create and publish a new registry module with a VCS repo
func (r *registryModules) CreateWithVCSConnection(ctx context.Context, options RegistryModuleCreateWithVCSConnectionOptions) (*RegistryModule, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}
	var u string
	if options.VCSRepo.OAuthTokenID != nil && options.VCSRepo.Branch == nil {
		u = "registry-modules"
	} else {
		u = fmt.Sprintf(
			"organizations/%s/registry-modules/vcs",
			url.PathEscape(*options.VCSRepo.OrganizationName),
		)
	}
	req, err := r.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	rm := &RegistryModule{}
	err = req.Do(ctx, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

// Read a specific registry module
func (r *registryModules) Read(ctx context.Context, moduleID RegistryModuleID) (*RegistryModule, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}

	var u string
	if moduleID.ID == "" {
		if moduleID.RegistryName == "" {
			log.Println("[WARN] Support for using the RegistryModuleID without RegistryName is deprecated as of release 1.5.0 and may be removed in a future version. The preferred method is to include the RegistryName in RegistryModuleID.")
			moduleID.RegistryName = PrivateRegistry
		}

		if moduleID.RegistryName == PrivateRegistry && strings.TrimSpace(moduleID.Namespace) == "" {
			log.Println("[WARN] Support for using the RegistryModuleID without Namespace is deprecated as of release 1.5.0 and may be removed in a future version. The preferred method is to include the Namespace in RegistryModuleID.")
			moduleID.Namespace = moduleID.Organization
		}

		u = fmt.Sprintf(
			"organizations/%s/registry-modules/%s/%s/%s/%s",
			url.PathEscape(moduleID.Organization),
			url.PathEscape(string(moduleID.RegistryName)),
			url.PathEscape(moduleID.Namespace),
			url.PathEscape(moduleID.Name),
			url.PathEscape(moduleID.Provider),
		)
	} else {
		u = fmt.Sprintf("registry-modules/%s", url.PathEscape(moduleID.ID))
	}

	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	rm := &RegistryModule{}
	err = req.Do(ctx, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}
func (r *registryModules) ReadVersion(ctx context.Context, moduleID RegistryModuleID, version string) (*RegistryModuleVersion, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}
	if !validString(&version) {
		return nil, ErrRequiredVersion
	}
	if !validStringID(&version) {
		return nil, ErrInvalidVersion
	}
	u := fmt.Sprintf(
		"organizations/%s/registry-modules/private/%s/%s/%s/version?module_version=%s",
		url.PathEscape(moduleID.Organization),
		url.PathEscape(moduleID.Organization),
		url.PathEscape(moduleID.Name),
		url.PathEscape(moduleID.Provider),
		url.PathEscape(version),
	)
	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	rmv := &RegistryModuleVersion{}
	err = req.Do(ctx, rmv)
	if err != nil {
		return nil, err
	}

	return rmv, nil
}

// Delete is used to delete the entire registry module
// Warning: This method is deprecated and will be removed from a future version of go-tfe. Use DeleteByName instead.
// See API Docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/private-registry/modules#delete-a-module
func (r *registryModules) Delete(ctx context.Context, organization, name string) error {
	if !validStringID(&organization) {
		return ErrInvalidOrg
	}
	if !validString(&name) {
		return ErrRequiredName
	}
	if !validStringID(&name) {
		return ErrInvalidName
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s",
		url.PathEscape(organization),
		url.PathEscape(name),
	)
	req, err := r.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// DeleteByName is used to delete the entire registry module
func (r *registryModules) DeleteByName(ctx context.Context, module RegistryModuleID) error {
	if err := module.validWhenDeleteByName(); err != nil {
		return err
	}

	u := fmt.Sprintf(
		"organizations/%s/registry-modules/%s/%s/%s",
		url.PathEscape(module.Organization),
		url.PathEscape(string(module.RegistryName)),
		url.PathEscape(module.Namespace),
		url.PathEscape(module.Name),
	)

	req, err := r.client.NewRequest("DELETE", u, nil)
	if err != nil && errors.Is(err, ErrResourceNotFound) {
		return r.Delete(ctx, module.Organization, module.Name)
	}

	return req.Do(ctx, nil)
}

// Delete a specified provider for the given module along with all its versions
func (r *registryModules) DeleteProvider(ctx context.Context, moduleID RegistryModuleID) error {
	if err := moduleID.validWhenDeleteByProvider(); err != nil {
		return err
	}

	u := fmt.Sprintf(
		"organizations/%s/registry-modules/%s/%s/%s/%s",
		url.PathEscape(moduleID.Organization),
		url.PathEscape(string(moduleID.RegistryName)),
		url.PathEscape(moduleID.Namespace),
		url.PathEscape(moduleID.Name),
		url.PathEscape(moduleID.Provider),
	)

	req, err := r.client.NewRequest("DELETE", u, nil)

	if err != nil && errors.Is(err, ErrResourceNotFound) {
		return r.deprecatedDeleteProvider(ctx, moduleID)
	}

	return req.Do(ctx, nil)
}

// Delete a specified version for the given provider of the module
func (r *registryModules) DeleteVersion(ctx context.Context, moduleID RegistryModuleID, version string) error {
	if err := moduleID.valid(); err != nil {
		return err
	}
	if !validString(&version) {
		return ErrRequiredVersion
	}
	if !validVersion(version) {
		return ErrInvalidVersion
	}

	u := fmt.Sprintf(
		"organizations/%s/registry-modules/%s/%s/%s/%s/%s",
		url.PathEscape(moduleID.Organization),
		url.PathEscape(string(moduleID.RegistryName)),
		url.PathEscape(moduleID.Namespace),
		url.PathEscape(moduleID.Name),
		url.PathEscape(moduleID.Provider),
		url.PathEscape(version),
	)
	req, err := r.client.NewRequest("DELETE", u, nil)
	if err != nil && errors.Is(err, ErrResourceNotFound) {
		return r.deprecatedDeleteVersion(ctx, moduleID, version)
	}

	return req.Do(ctx, nil)
}

func (o RegistryModuleID) valid() error {
	if validString(&o.ID) && validStringID(&o.ID) {
		return nil
	}

	if !validStringID(&o.Organization) {
		return ErrInvalidOrg
	}

	if !validString(&o.Name) {
		return ErrRequiredName
	}

	if !validStringID(&o.Name) {
		return ErrInvalidName
	}

	if !validString(&o.Provider) {
		return ErrRequiredProvider
	}

	if !validStringID(&o.Provider) {
		return ErrInvalidProvider
	}

	switch o.RegistryName {
	case PublicRegistry:
		if !validString(&o.Namespace) {
			return ErrRequiredNamespace
		}
	case PrivateRegistry:
	case "":
		// no-op:  RegistryName is optional
	// for all other string
	default:
		return ErrInvalidRegistryName
	}

	return nil
}

func (o RegistryModuleID) validWhenDeleteByProvider() error {
	if !validStringID(&o.Organization) {
		return ErrInvalidOrg
	}

	if !validString(&o.Name) {
		return ErrRequiredName
	}

	if !validStringID(&o.Name) {
		return ErrInvalidName
	}

	if !validString(&o.Provider) {
		return ErrRequiredProvider
	}

	if !validStringID(&o.Provider) {
		return ErrInvalidProvider
	}
	// RegistryName is required in this DELETE call
	switch o.RegistryName {
	case PublicRegistry:
		if !validString(&o.Namespace) {
			return ErrRequiredNamespace
		}
	case PrivateRegistry:
	case "":
		return ErrInvalidRegistryName
	default:
		return ErrInvalidRegistryName
	}

	return nil
}

func (o RegistryModuleID) validWhenDeleteByName() error {
	if !validStringID(&o.Organization) {
		return ErrInvalidOrg
	}

	if !validString(&o.Name) {
		return ErrRequiredName
	}

	if !validStringID(&o.Name) {
		return ErrInvalidName
	}

	// RegistryName is required in this DELETE call
	switch o.RegistryName {
	case PublicRegistry:
		if !validString(&o.Namespace) {
			return ErrRequiredNamespace
		}
	case PrivateRegistry:
	case "":
		return ErrInvalidRegistryName
	default:
		return ErrInvalidRegistryName
	}

	return nil
}

func (o RegistryModuleCreateOptions) valid() error {
	if !validString(o.Name) {
		return ErrRequiredName
	}
	if !validStringID(o.Name) {
		return ErrInvalidName
	}
	if !validString(o.Provider) {
		return ErrRequiredProvider
	}
	if !validStringID(o.Provider) {
		return ErrInvalidProvider
	}

	switch o.RegistryName {
	case PublicRegistry:
		if !validString(&o.Namespace) {
			return ErrRequiredNamespace
		}
	case PrivateRegistry:
		if validString(&o.Namespace) {
			return ErrUnsupportedBothNamespaceAndPrivateRegistryName
		}
	case "":
		// no-op:  RegistryName is optional
	// for all other string
	default:
		return ErrInvalidRegistryName
	}
	return nil
}

func (o RegistryModuleCreateVersionOptions) valid() error {
	if !validString(o.Version) {
		return ErrRequiredVersion
	}
	if !validVersion(*o.Version) {
		return ErrInvalidVersion
	}
	return nil
}

func (o RegistryModuleCreateWithVCSConnectionOptions) valid() error {
	if o.VCSRepo == nil {
		return ErrRequiredVCSRepo
	}

	if o.TestConfig != nil && o.TestConfig.TestsEnabled != nil {
		if *o.TestConfig.TestsEnabled {
			if !validString(o.VCSRepo.Branch) {
				return ErrRequiredBranchWhenTestsEnabled
			}
		}
	}

	if o.VCSRepo.Tags != nil && *o.VCSRepo.Tags {
		if validString(o.VCSRepo.Branch) {
			return ErrBranchMustBeEmptyWhenTagsEnabled
		}
	}

	return o.VCSRepo.valid()
}

func (o RegistryModuleVCSRepoOptions) valid() error {
	if !validString(o.Identifier) {
		return ErrRequiredIdentifier
	}
	if !validString(o.OAuthTokenID) && !validString(o.GHAInstallationID) {
		return ErrRequiredOauthTokenOrGithubAppInstallationID
	}
	if (!validString(o.OAuthTokenID) && validString(o.GHAInstallationID)) || validString(o.Branch) {
		if !validString(o.OrganizationName) {
			return ErrInvalidOrg
		}
	}
	if !validString(o.DisplayIdentifier) {
		return ErrRequiredDisplayIdentifier
	}
	return nil
}

func (r *registryModules) deprecatedDeleteProvider(ctx context.Context, moduleID RegistryModuleID) error {
	if err := moduleID.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s/%s",
		url.PathEscape(moduleID.Organization),
		url.PathEscape(moduleID.Name),
		url.PathEscape(moduleID.Provider),
	)
	req, err := r.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (r *registryModules) deprecatedDeleteVersion(ctx context.Context, moduleID RegistryModuleID, version string) error {
	if err := moduleID.valid(); err != nil {
		return err
	}
	if !validString(&version) {
		return ErrRequiredVersion
	}
	if !validVersion(version) {
		return ErrInvalidVersion
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s/%s/%s",
		url.PathEscape(moduleID.Organization),
		url.PathEscape(moduleID.Name),
		url.PathEscape(moduleID.Provider),
		url.PathEscape(version),
	)
	req, err := r.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func NewPublicRegistryModuleID(organization, namespace, name, provider string) RegistryModuleID {
	return RegistryModuleID{
		Organization: organization,
		Namespace:    namespace,
		Name:         name,
		RegistryName: PublicRegistry,
		Provider:     provider,
	}
}

func NewPrivateRegistryModuleID(organization, name, provider string) RegistryModuleID {
	return RegistryModuleID{
		Organization: organization,
		Namespace:    organization,
		Name:         name,
		RegistryName: PrivateRegistry,
		Provider:     provider,
	}
}
