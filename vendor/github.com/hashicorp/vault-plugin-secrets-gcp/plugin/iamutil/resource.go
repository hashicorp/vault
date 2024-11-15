// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iamutil

import (
	"context"

	"github.com/hashicorp/go-gcp-common/gcputil"
)

// Resource handles constructing HTTP requests for getting and
// setting IAM policies.
type Resource interface {
	GetIamPolicy(context.Context, *ApiHandle) (*Policy, error)
	SetIamPolicy(context.Context, *ApiHandle, *Policy) (*Policy, error)
	GetConfig() *RestResource
	GetRelativeId() *gcputil.RelativeResourceName
}

type RestResource struct {
	// Name is the base name of the resource
	// i.e. for a GCE instance: "instance"
	Name string

	// TypeKey is the identifying path for the resource, or
	// the RESTful resource identifier without resource IDs
	// i.e. For a GCE instance: "projects/zones/instances"
	TypeKey string

	// Service is the name of the service this resource belongs to.
	Service string

	// IsPreferredVersion is true if this version of the API/resource is preferred.
	IsPreferredVersion bool

	// HTTP metadata for getting Policy data in GCP
	GetMethod RestMethod

	// HTTP metadata for setting Policy data in GCP
	SetMethod RestMethod

	// Ordered parameters to be replaced in method paths
	Parameters []string

	// Mapping of collection ids onto the parameter to be replaced
	CollectionReplacementKeys map[string]string
}

type RestMethod struct {
	HttpMethod    string
	BaseURL       string
	Path          string
	RequestFormat string
}
