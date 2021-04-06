package tfe

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ RegistryModules = (*registryModules)(nil)

// RegistryModules describes all the registry module related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/modules.html
type RegistryModules interface {
	// Create a registry module without a VCS repo
	Create(ctx context.Context, organization string, options RegistryModuleCreateOptions) (*RegistryModule, error)

	// Create a registry module version
	CreateVersion(ctx context.Context, organization string, name string, provider string, options RegistryModuleCreateVersionOptions) (*RegistryModuleVersion, error)

	// Create and publish a registry module with a VCS repo
	CreateWithVCSConnection(ctx context.Context, options RegistryModuleCreateWithVCSConnectionOptions) (*RegistryModule, error)

	// Read a registry module
	Read(ctx context.Context, organization string, name string, provider string) (*RegistryModule, error)

	// Delete a registry module
	Delete(ctx context.Context, organization string, name string) error

	// Delete a specific registry module provider
	DeleteProvider(ctx context.Context, organization string, name string, provider string) error

	// Delete a specific registry module version
	DeleteVersion(ctx context.Context, organization string, name string, provider string, version string) error
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

// RegistryModule represents a registry module
type RegistryModule struct {
	ID              string                          `jsonapi:"primary,registry-modules"`
	Name            string                          `jsonapi:"attr,name"`
	Provider        string                          `jsonapi:"attr,provider"`
	Permissions     *RegistryModulePermissions      `jsonapi:"attr,permissions"`
	Status          RegistryModuleStatus            `jsonapi:"attr,status"`
	VCSRepo         *VCSRepo                        `jsonapi:"attr,vcs-repo"`
	VersionStatuses []RegistryModuleVersionStatuses `jsonapi:"attr,version-statuses"`
	CreatedAt       string                          `jsonapi:"attr,created-at"`
	UpdatedAt       string                          `jsonapi:"attr,updated-at"`

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`
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
}

type RegistryModulePermissions struct {
	CanDelete bool `json:"can-delete"`
	CanResync bool `json:"can-resync"`
	CanRetry  bool `json:"can-retry"`
}

type RegistryModuleVersionStatuses struct {
	Version string                      `json:"version"`
	Status  RegistryModuleVersionStatus `json:"status"`
	Error   string                      `json:"error"`
}

// RegistryModuleCreateOptions is used when creating a registry module without a VCS repo
type RegistryModuleCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,registry-modules"`

	Name     *string `jsonapi:"attr,name"`
	Provider *string `jsonapi:"attr,provider"`
}

func (o RegistryModuleCreateOptions) valid() error {
	if !validString(o.Name) {
		return errors.New("name is required")
	}
	if !validStringID(o.Name) {
		return errors.New("invalid value for name")
	}
	if !validString(o.Provider) {
		return errors.New("provider is required")
	}
	if !validStringID(o.Provider) {
		return errors.New("invalid value for provider")
	}
	return nil
}

// Create a new registry module without a VCS repo
func (r *registryModules) Create(ctx context.Context, organization string, options RegistryModuleCreateOptions) (*RegistryModule, error) {
	if !validStringID(&organization) {
		return nil, errors.New("invalid value for organization")
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf(
		"organizations/%s/registry-modules",
		url.QueryEscape(organization),
	)
	req, err := r.client.newRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	rm := &RegistryModule{}
	err = r.client.do(ctx, req, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

// RegistryModuleCreateVersionOptions is used when creating a registry module version
type RegistryModuleCreateVersionOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,registry-module-versions"`

	Version *string `jsonapi:"attr,version"`
}

func (o RegistryModuleCreateVersionOptions) valid() error {
	if !validString(o.Version) {
		return errors.New("version is required")
	}
	if !validStringID(o.Version) {
		return errors.New("invalid value for version")
	}
	return nil
}

// Create a new registry module version
func (r *registryModules) CreateVersion(ctx context.Context, organization string, name string, provider string, options RegistryModuleCreateVersionOptions) (*RegistryModuleVersion, error) {
	if !validStringID(&organization) {
		return nil, errors.New("invalid value for organization")
	}
	if !validString(&name) {
		return nil, errors.New("name is required")
	}
	if !validStringID(&name) {
		return nil, errors.New("invalid value for name")
	}
	if !validString(&provider) {
		return nil, errors.New("provider is required")
	}
	if !validStringID(&provider) {
		return nil, errors.New("invalid value for provider")
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf(
		"registry-modules/%s/%s/%s/versions",
		url.QueryEscape(organization),
		url.QueryEscape(name),
		url.QueryEscape(provider),
	)
	req, err := r.client.newRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	rmv := &RegistryModuleVersion{}
	err = r.client.do(ctx, req, rmv)
	if err != nil {
		return nil, err
	}

	return rmv, nil
}

// RegistryModuleCreateWithVCSConnectionOptions is used when creating a registry module with a VCS repo
type RegistryModuleCreateWithVCSConnectionOptions struct {
        ID string `jsonapi:"primary,registry-modules"`

	// VCS repository information
	VCSRepo *RegistryModuleVCSRepoOptions `jsonapi:"attr,vcs-repo"`
}

func (o RegistryModuleCreateWithVCSConnectionOptions) valid() error {
	if o.VCSRepo == nil {
		return errors.New("vcs repo is required")
	}
	return o.VCSRepo.valid()
}

type RegistryModuleVCSRepoOptions struct {
	Identifier        *string `json:"identifier"`
	OAuthTokenID      *string `json:"oauth-token-id"`
	DisplayIdentifier *string `json:"display-identifier"`
}

func (o RegistryModuleVCSRepoOptions) valid() error {
	if !validString(o.Identifier) {
		return errors.New("identifier is required")
	}
	if !validString(o.OAuthTokenID) {
		return errors.New("oauth token ID is required")
	}
	if !validString(o.DisplayIdentifier) {
		return errors.New("display identifier is required")
	}
	return nil
}

// CreateWithVCSConnection is used to create and publish a new registry module with a VCS repo
func (r *registryModules) CreateWithVCSConnection(ctx context.Context, options RegistryModuleCreateWithVCSConnectionOptions) (*RegistryModule, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

        // Make sure we don't send a user provided ID.
        options.ID = ""

	req, err := r.client.newRequest("POST", "registry-modules", &options)
	if err != nil {
		return nil, err
	}

	rm := &RegistryModule{}
	err = r.client.do(ctx, req, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

// Read a specific registry module
func (r *registryModules) Read(ctx context.Context, organization string, name string, provider string) (*RegistryModule, error) {
	if !validStringID(&organization) {
		return nil, errors.New("invalid value for organization")
	}
	if !validString(&name) {
		return nil, errors.New("name is required")
	}
	if !validStringID(&name) {
		return nil, errors.New("invalid value for name")
	}
	if !validString(&provider) {
		return nil, errors.New("provider is required")
	}
	if !validStringID(&provider) {
		return nil, errors.New("invalid value for provider")
	}

	u := fmt.Sprintf(
		"registry-modules/show/%s/%s/%s",
		url.QueryEscape(organization),
		url.QueryEscape(name),
		url.QueryEscape(provider),
	)
	req, err := r.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	rm := &RegistryModule{}
	err = r.client.do(ctx, req, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

// Delete is used to delete the entire registry module
func (r *registryModules) Delete(ctx context.Context, organization string, name string) error {
	if !validStringID(&organization) {
		return errors.New("invalid value for organization")
	}
	if !validString(&name) {
		return errors.New("name is required")
	}
	if !validStringID(&name) {
		return errors.New("invalid value for name")
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s",
		url.QueryEscape(organization),
		url.QueryEscape(name),
	)
	req, err := r.client.newRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return r.client.do(ctx, req, nil)
}

// DeleteProvider is used to delete the specific registry module provider
func (r *registryModules) DeleteProvider(ctx context.Context, organization string, name string, provider string) error {
	if !validStringID(&organization) {
		return errors.New("invalid value for organization")
	}
	if !validString(&name) {
		return errors.New("name is required")
	}
	if !validStringID(&name) {
		return errors.New("invalid value for name")
	}
	if !validString(&provider) {
		return errors.New("provider is required")
	}
	if !validStringID(&provider) {
		return errors.New("invalid value for provider")
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s/%s",
		url.QueryEscape(organization),
		url.QueryEscape(name),
		url.QueryEscape(provider),
	)
	req, err := r.client.newRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return r.client.do(ctx, req, nil)
}

// DeleteVersion is used to delete the specific registry module version
func (r *registryModules) DeleteVersion(ctx context.Context, organization string, name string, provider string, version string) error {
	if !validStringID(&organization) {
		return errors.New("invalid value for organization")
	}
	if !validString(&name) {
		return errors.New("name is required")
	}
	if !validStringID(&name) {
		return errors.New("invalid value for name")
	}
	if !validString(&provider) {
		return errors.New("provider is required")
	}
	if !validStringID(&provider) {
		return errors.New("invalid value for provider")
	}
	if !validString(&version) {
		return errors.New("version is required")
	}
	if !validStringID(&version) {
		return errors.New("invalid value for version")
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s/%s/%s",
		url.QueryEscape(organization),
		url.QueryEscape(name),
		url.QueryEscape(provider),
		url.QueryEscape(version),
	)
	req, err := r.client.newRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return r.client.do(ctx, req, nil)
}
