// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"io"
	"net/url"
)

// StackSources describes all the stack-sources related methods that the HCP Terraform API supports.
// NOTE WELL: This is a beta feature and is subject to change until noted otherwise in the
// release notes.
type StackSources interface {
	// Read retrieves a stack source by its ID.
	Read(ctx context.Context, stackSourceID string) (*StackSource, error)

	// CreateAndUpload packages and uploads the specified Terraform Stacks
	// configuration files in association with a Stack.
	CreateAndUpload(ctx context.Context, stackID string, path string, opts *CreateStackSourceOptions) (*StackSource, error)

	// UploadTarGzip is used to upload Terraform configuration files contained a tar gzip archive.
	// Any stream implementing io.Reader can be passed into this method. This method is also
	// particularly useful for tar streams created by non-default go-slug configurations.
	//
	// **Note**: This method does not validate the content being uploaded and is therefore the caller's
	// responsibility to ensure the raw content is a valid Terraform configuration.
	UploadTarGzip(ctx context.Context, uploadURL string, archive io.Reader) error
}

type CreateStackSourceOptions struct {
	SelectedDeployments []string `jsonapi:"attr,selected-deployments,omitempty"`
}

var _ StackSources = (*stackSources)(nil)

type stackSources struct {
	client *Client
}

// StackSource represents a source of Terraform Stacks configuration files.
type StackSource struct {
	ID                 string              `jsonapi:"primary,stack-sources"`
	UploadURL          *string             `jsonapi:"attr,upload-url"`
	StackConfiguration *StackConfiguration `jsonapi:"relation,stack-configuration"`
	Stack              *Stack              `jsonapi:"relation,stack"`
}

// Read retrieves a stack source by its ID.
func (s *stackSources) Read(ctx context.Context, stackSourceID string) (*StackSource, error) {
	u := fmt.Sprintf("stack-sources/%s", url.PathEscape(stackSourceID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	ss := &StackSource{}
	err = req.Do(ctx, ss)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

// CreateAndUpload packages and uploads the specified Terraform Stacks
// configuration files in association with a Stack.
func (s *stackSources) CreateAndUpload(ctx context.Context, stackID, path string, opts *CreateStackSourceOptions) (*StackSource, error) {
	if opts == nil {
		opts = &CreateStackSourceOptions{}
	}
	u := fmt.Sprintf("stacks/%s/stack-sources", url.PathEscape(stackID))
	req, err := s.client.NewRequest("POST", u, opts)
	if err != nil {
		return nil, err
	}

	ss := &StackSource{}
	err = req.Do(ctx, ss)
	if err != nil {
		return nil, err
	}

	body, err := packContents(path)
	if err != nil {
		return nil, err
	}

	return ss, s.UploadTarGzip(ctx, *ss.UploadURL, body)
}

// UploadTarGzip is used to upload Terraform configuration files contained a tar gzip archive.
// Any stream implementing io.Reader can be passed into this method. This method is also
// particularly useful for tar streams created by non-default go-slug configurations.
//
// **Note**: This method does not validate the content being uploaded and is therefore the caller's
// responsibility to ensure the raw content is a valid Terraform configuration.
func (s *stackSources) UploadTarGzip(ctx context.Context, uploadURL string, archive io.Reader) error {
	return s.client.doForeignPUTRequest(ctx, uploadURL, archive)
}
