package tfe

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ AdminTerraformVersions = (*adminTerraformVersions)(nil)

// AdminTerraformVersions describes all the admin terraform versions related methods that
// the Terraform Enterprise API supports.
// Note that admin terraform versions are only available in Terraform Enterprise.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/admin/terraform-versions.html
type AdminTerraformVersions interface {
	// List all the terraform versions.
	List(ctx context.Context, options AdminTerraformVersionsListOptions) (*AdminTerraformVersionsList, error)

	// Read a terraform version by its ID.
	Read(ctx context.Context, id string) (*AdminTerraformVersion, error)

	// Create a terraform version.
	Create(ctx context.Context, options AdminTerraformVersionCreateOptions) (*AdminTerraformVersion, error)

	// Update a terraform version.
	Update(ctx context.Context, id string, options AdminTerraformVersionUpdateOptions) (*AdminTerraformVersion, error)

	// Delete a terraform version
	Delete(ctx context.Context, id string) error
}

// adminTerraformVersions implements AdminTerraformVersions.
type adminTerraformVersions struct {
	client *Client
}

// AdminTerraformVersion represents a Terraform Version
type AdminTerraformVersion struct {
	ID        string    `jsonapi:"primary,terraform-versions"`
	Version   string    `jsonapi:"attr,version"`
	URL       string    `jsonapi:"attr,url"`
	Sha       string    `jsonapi:"attr,sha"`
	Official  bool      `jsonapi:"attr,official"`
	Enabled   bool      `jsonapi:"attr,enabled"`
	Beta      bool      `jsonapi:"attr,beta"`
	Usage     int       `jsonapi:"attr,usage"`
	CreatedAt time.Time `jsonapi:"attr,created-at,iso8601"`
}

// AdminTerraformVersionsListOptions represents the options for listing
// terraform versions.
type AdminTerraformVersionsListOptions struct {
	ListOptions
}

// AdminTerraformVersionsList represents a list of terraform versions.
type AdminTerraformVersionsList struct {
	*Pagination
	Items []*AdminTerraformVersion
}

// List all the terraform versions.
func (a *adminTerraformVersions) List(ctx context.Context, options AdminTerraformVersionsListOptions) (*AdminTerraformVersionsList, error) {
	req, err := a.client.newRequest("GET", "admin/terraform-versions", &options)
	if err != nil {
		return nil, err
	}

	tvl := &AdminTerraformVersionsList{}
	err = a.client.do(ctx, req, tvl)
	if err != nil {
		return nil, err
	}

	return tvl, nil
}

// Read a terraform version by its ID.
func (a *adminTerraformVersions) Read(ctx context.Context, id string) (*AdminTerraformVersion, error) {
	if !validStringID(&id) {
		return nil, ErrInvalidTerraformVersionID
	}

	u := fmt.Sprintf("admin/terraform-versions/%s", url.QueryEscape(id))
	req, err := a.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	tfv := &AdminTerraformVersion{}
	err = a.client.do(ctx, req, tfv)
	if err != nil {
		return nil, err
	}

	return tfv, nil
}

// AdminTerraformVersionCreateOptions for creating a terraform version.
// https://www.terraform.io/docs/cloud/api/admin/terraform-versions.html#request-body
type AdminTerraformVersionCreateOptions struct {
	Type     string  `jsonapi:"primary,terraform-versions"`
	Version  *string `jsonapi:"attr,version"`
	URL      *string `jsonapi:"attr,url"`
	Sha      *string `jsonapi:"attr,sha"`
	Official *bool   `jsonapi:"attr,official,omitempty"`
	Enabled  *bool   `jsonapi:"attr,enabled,omitempty"`
	Beta     *bool   `jsonapi:"attr,beta,omitempty"`
}

// Create a new terraform version.
func (a *adminTerraformVersions) Create(ctx context.Context, options AdminTerraformVersionCreateOptions) (*AdminTerraformVersion, error) {
	req, err := a.client.newRequest("POST", "admin/terraform-versions", &options)
	if err != nil {
		return nil, err
	}

	tfv := &AdminTerraformVersion{}
	err = a.client.do(ctx, req, tfv)
	if err != nil {
		return nil, err
	}

	return tfv, nil
}

// AdminTerraformVersionUpdateOptions for updating terraform version.
// https://www.terraform.io/docs/cloud/api/admin/terraform-versions.html#request-body
type AdminTerraformVersionUpdateOptions struct {
	Type     string  `jsonapi:"primary,terraform-versions"`
	Version  *string `jsonapi:"attr,version,omitempty"`
	URL      *string `jsonapi:"attr,url,omitempty"`
	Sha      *string `jsonapi:"attr,sha,omitempty"`
	Official *bool   `jsonapi:"attr,official,omitempty"`
	Enabled  *bool   `jsonapi:"attr,enabled,omitempty"`
	Beta     *bool   `jsonapi:"attr,beta,omitempty"`
}

// Update an existing terraform version.
func (a *adminTerraformVersions) Update(ctx context.Context, id string, options AdminTerraformVersionUpdateOptions) (*AdminTerraformVersion, error) {
	if !validStringID(&id) {
		return nil, ErrInvalidTerraformVersionID
	}

	u := fmt.Sprintf("admin/terraform-versions/%s", url.QueryEscape(id))
	req, err := a.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	tfv := &AdminTerraformVersion{}
	err = a.client.do(ctx, req, tfv)
	if err != nil {
		return nil, err
	}

	return tfv, nil
}

// Delete a terraform version.
func (a *adminTerraformVersions) Delete(ctx context.Context, id string) error {
	if !validStringID(&id) {
		return ErrInvalidTerraformVersionID
	}

	u := fmt.Sprintf("admin/terraform-versions/%s", url.QueryEscape(id))
	req, err := a.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return a.client.do(ctx, req, nil)
}
