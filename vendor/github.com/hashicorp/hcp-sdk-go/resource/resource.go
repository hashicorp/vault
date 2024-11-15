// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-shared/v1/models"
)

const (
	tokenOrganization = "organization"
	tokenProject      = "project"
	tokenSep          = "/"
)

// Resource is a representation of a HCP resource identifier
type Resource struct {
	// ID uniquely identifies a resource within an HCP project
	ID string
	// Type is the name of the kind of resource identified
	Type string
	// Organization is the UUID of the HCP organization the resource belongs to
	Organization string
	// Project is the UUID of the HCP project the resource belongs to
	Project string
}

// String encodes the resource identifier in the following canonical format:
//
//	"organization/<Organization UUID>/project/<Project UUID>/<Type>/<ID>"
//
// Example:
//
//	"organization/ccbdd191-5dc3-4a73-9e05-6ac30ca67992/project/36019e0d-ed59-4df6-9990-05bb7fc793b6/hashicorp.consul.linked-cluster/prod-on-prem"
func (r Resource) String() string {
	return strings.Join([]string{
		tokenOrganization, r.Organization,
		tokenProject, r.Project,
		r.Type, r.ID,
	}, tokenSep)
}

// Location returns a *models.HashicorpCloudLocationLocation initialized with the Resource's organization and project IDs.
func (r Resource) Location() *models.HashicorpCloudLocationLocation {
	return &models.HashicorpCloudLocationLocation{OrganizationID: r.Organization, ProjectID: r.Project}
}

// Link returns a *models.HashicorpCloudLocationLink initialized with values from the Resource
func (r Resource) Link() *models.HashicorpCloudLocationLink {
	return &models.HashicorpCloudLocationLink{
		ID:       r.ID,
		Type:     r.Type,
		Location: r.Location(),
	}
}

// FromLink converts a models.HashicorpCloudLocationLink to a Resource.
func FromLink(l *models.HashicorpCloudLocationLink) (r Resource, err error) {
	if l == nil || l.Location == nil {
		return r, parseErr(errors.New("link and link.Location must not be nil"))
	}
	return Resource{
		ID:           l.ID,
		Type:         l.Type,
		Organization: l.Location.OrganizationID,
		Project:      l.Location.ProjectID,
	}, nil
}

// FromString parses the string encoding of a resource identifier.
func FromString(str string) (r Resource, err error) {
	err = r.UnmarshalText([]byte(str))
	return r, err
}

// UnmarshalText implements the encoding.TextUnmarshaler interface and parses an encoded Resource.
func (r *Resource) UnmarshalText(text []byte) error {
	if r == nil {
		return fmt.Errorf("resource cannot be nil")
	}

	parts := strings.SplitN(string(text), tokenSep, 6)
	if len(parts) != 6 {
		return parseErr(fmt.Errorf("unexpected number of tokens %d", len(parts)))
	}

	if parts[0] != tokenOrganization {
		return parseErr(fmt.Errorf("unexpected token %q", parts[0]))
	}
	if err := validation.Validate(parts[1], is.UUID); err != nil {
		return parseErr(err)
	}
	if parts[2] != tokenProject {
		return parseErr(fmt.Errorf("unexpected token %q", parts[2]))
	}
	if err := validation.Validate(parts[3], is.UUID); err != nil {
		return parseErr(err)
	}

	r.Organization = parts[1]
	r.Project = parts[3]
	r.Type = parts[4]
	r.ID = parts[5]
	return nil
}

func parseErr(err error) error {
	return fmt.Errorf("could not parse resource: %w", err)
}

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String.
func (r *Resource) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

// MarshalJSON implements the json.Marshaler interface
func (r *Resource) MarshalJSON() ([]byte, error) {
	return []byte("\"" + r.String() + "\""), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (r *Resource) UnmarshalJSON(bytes []byte) error {
	return r.UnmarshalText(bytes[1 : len(bytes)-1])
}

// Must is a helper function that wraps a call to a function returning (Resource, error) such as FromLink or FromString
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//
//	var packageResource = resource.Must(resource.FromString("..."))
func Must(r Resource, err error) Resource {
	if err != nil {
		panic(err)
	}
	return r
}
