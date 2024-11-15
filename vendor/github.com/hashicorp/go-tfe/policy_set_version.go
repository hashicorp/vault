// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ PolicySetVersions = (*policySetVersions)(nil)

// PolicySetVersions describes all the Policy Set Version related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/policy-sets#create-a-policy-set-version
type PolicySetVersions interface {
	// Create is used to create a new Policy Set Version.
	Create(ctx context.Context, policySetID string) (*PolicySetVersion, error)

	// Read is used to read a Policy Set Version by its ID.
	Read(ctx context.Context, policySetVersionID string) (*PolicySetVersion, error)

	// Upload uploads policy files. It takes a Policy Set Version and a path
	// to the set of sentinel files, which will be packaged by hashicorp/go-slug
	// before being uploaded.
	Upload(ctx context.Context, psv PolicySetVersion, path string) error
}

// policySetVersions implements PolicySetVersions.
type policySetVersions struct {
	client *Client
}

// PolicySetVersionSource represents a source type of a policy set version.
type PolicySetVersionSource string

// List all available sources for a Policy Set Version.
const (
	PolicySetVersionSourceAPI       PolicySetVersionSource = "tfe-api"
	PolicySetVersionSourceADO       PolicySetVersionSource = "ado"
	PolicySetVersionSourceBitBucket PolicySetVersionSource = "bitbucket"
	PolicySetVersionSourceGitHub    PolicySetVersionSource = "github"
	PolicySetVersionSourceGitLab    PolicySetVersionSource = "gitlab"
)

// PolicySetVersionStatus represents a policy set version status.
type PolicySetVersionStatus string

// List all available policy set version statuses.
const (
	PolicySetVersionErrored    PolicySetVersionStatus = "errored"
	PolicySetVersionIngressing PolicySetVersionStatus = "ingressing"
	PolicySetVersionPending    PolicySetVersionStatus = "pending"
	PolicySetVersionReady      PolicySetVersionStatus = "ready"
)

// PolicySetVersionStatusTimestamps holds the timestamps for individual policy
// set version statuses.
type PolicySetVersionStatusTimestamps struct {
	PendingAt    time.Time `jsonapi:"attr,pending-at,rfc3339"`
	IngressingAt time.Time `jsonapi:"attr,ingressing-at,rfc3339"`
	ReadyAt      time.Time `jsonapi:"attr,ready-at,rfc3339"`
	ErroredAt    time.Time `jsonapi:"attr,errored-at,rfc3339"`
}

// PolicySetVersion represents a Terraform Enterprise Policy Set Version
type PolicySetVersion struct {
	ID               string                           `jsonapi:"primary,policy-set-versions"`
	Source           PolicySetVersionSource           `jsonapi:"attr,source"`
	Status           PolicySetVersionStatus           `jsonapi:"attr,status"`
	StatusTimestamps PolicySetVersionStatusTimestamps `jsonapi:"attr,status-timestamps"`
	Error            string                           `jsonapi:"attr,error"`
	ErrorMessage     string                           `jsonapi:"attr,error-message"`
	CreatedAt        time.Time                        `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt        time.Time                        `jsonapi:"attr,updated-at,iso8601"`

	// Relations
	PolicySet *PolicySet `jsonapi:"relation,policy-set"`

	// Links
	Links map[string]interface{} `jsonapi:"links,omitempty"`
}

func (p PolicySetVersion) uploadURL() (string, error) {
	uploadURL, ok := p.Links["upload"].(string)
	if !ok {
		return uploadURL, fmt.Errorf("the Policy Set Version does not contain an upload link")
	}

	if uploadURL == "" {
		return uploadURL, fmt.Errorf("the Policy Set Version upload URL is empty")
	}

	return uploadURL, nil
}

// Create is used to create a new Policy Set Version.
func (p *policySetVersions) Create(ctx context.Context, policySetID string) (*PolicySetVersion, error) {
	if !validStringID(&policySetID) {
		return nil, ErrInvalidPolicySetID
	}

	u := fmt.Sprintf("policy-sets/%s/versions", url.PathEscape(policySetID))
	req, err := p.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	psv := &PolicySetVersion{}
	err = req.Do(ctx, psv)
	if err != nil {
		return nil, err
	}

	return psv, nil
}

// Read is used to read a Policy Set Version by its ID.
func (p *policySetVersions) Read(ctx context.Context, policySetVersionID string) (*PolicySetVersion, error) {
	if !validStringID(&policySetVersionID) {
		return nil, ErrInvalidPolicySetID
	}

	u := fmt.Sprintf("policy-set-versions/%s", url.PathEscape(policySetVersionID))
	req, err := p.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	psv := &PolicySetVersion{}
	err = req.Do(ctx, psv)
	if err != nil {
		return nil, err
	}

	return psv, nil
}

// Upload uploads policy files. It takes a Policy Set Version and a path
// to the set of sentinel files, which will be packaged by hashicorp/go-slug
// before being uploaded.
func (p *policySetVersions) Upload(ctx context.Context, psv PolicySetVersion, path string) error {
	uploadURL, err := psv.uploadURL()
	if err != nil {
		return err
	}

	body, err := packContents(path)
	if err != nil {
		return err
	}

	return p.client.doForeignPUTRequest(ctx, uploadURL, body)
}
