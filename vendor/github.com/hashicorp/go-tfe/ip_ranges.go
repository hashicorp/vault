// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
)

// Compile-time proof of interface implementation.
var _ IPRanges = (*ipRanges)(nil)

// IP Ranges provides a list of HCP Terraform or Terraform Enterprise's IP ranges.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/ip-ranges
type IPRanges interface {
	// Retrieve HCP Terraform IP ranges. If `modifiedSince` is not an empty string
	// then it will only return the IP ranges changes since that date.
	// The format for `modifiedSince` can be found here:
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/If-Modified-Since
	Read(ctx context.Context, modifiedSince string) (*IPRange, error)
}

// ipRanges implements IPRanges interface.
type ipRanges struct {
	client *Client
}

// IPRange represents a list of HCP Terraform's IP ranges
type IPRange struct {
	// List of IP ranges in CIDR notation used for connections from user site to HCP Terraform APIs
	API []string `json:"api"`
	// List of IP ranges in CIDR notation used for notifications
	Notifications []string `json:"notifications"`
	// List of IP ranges in CIDR notation used for outbound requests from Sentinel policies
	Sentinel []string `json:"sentinel"`
	// List of IP ranges in CIDR notation used for connecting to VCS providers
	VCS []string `json:"vcs"`
}

// Read an IPRange that was not modified since the specified date.
func (i *ipRanges) Read(ctx context.Context, modifiedSince string) (*IPRange, error) {
	req, err := i.client.NewRequest("GET", "/api/meta/ip-ranges", nil)
	if err != nil {
		return nil, err
	}

	if modifiedSince != "" {
		req.Header.Add("If-Modified-Since", modifiedSince)
	}

	ir := &IPRange{}
	err = req.DoJSON(ctx, ir)
	if err != nil {
		return nil, err
	}

	return ir, nil
}
