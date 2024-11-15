// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Compile-time proof of interface implementation
var _ GPGKeys = (*gpgKeys)(nil)

// GPGKeys describes all the GPG key related methods that the Terraform Private Registry API supports.
//
// TFE API Docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/private-registry/gpg-keys
type GPGKeys interface {
	// Lists GPG keys in a private registry.
	ListPrivate(ctx context.Context, options GPGKeyListOptions) (*GPGKeyList, error)

	// Uploads a GPG Key to a private registry scoped with a namespace.
	Create(ctx context.Context, registryName RegistryName, options GPGKeyCreateOptions) (*GPGKey, error)

	// Read a GPG key.
	Read(ctx context.Context, keyID GPGKeyID) (*GPGKey, error)

	// Update a GPG key.
	Update(ctx context.Context, keyID GPGKeyID, options GPGKeyUpdateOptions) (*GPGKey, error)

	// Delete a GPG key.
	Delete(ctx context.Context, keyID GPGKeyID) error
}

// gpgKeys implements GPGKeys
type gpgKeys struct {
	client *Client
}

// GPGKeyList represents a list of GPG keys.
type GPGKeyList struct {
	*Pagination
	Items []*GPGKey
}

// GPGKey represents a signed GPG key for a HCP Terraform or Terraform Enterprise private provider.
type GPGKey struct {
	ID             string    `jsonapi:"primary,gpg-keys"`
	AsciiArmor     string    `jsonapi:"attr,ascii-armor"`
	CreatedAt      time.Time `jsonapi:"attr,created-at,iso8601"`
	KeyID          string    `jsonapi:"attr,key-id"`
	Namespace      string    `jsonapi:"attr,namespace"`
	Source         string    `jsonapi:"attr,source"`
	SourceURL      *string   `jsonapi:"attr,source-url"`
	TrustSignature string    `jsonapi:"attr,trust-signature"`
	UpdatedAt      time.Time `jsonapi:"attr,updated-at,iso8601"`
}

// GPGKeyID represents the set of identifiers used to fetch a GPG key.
type GPGKeyID struct {
	RegistryName RegistryName
	Namespace    string
	KeyID        string
}

// GPGKeyListOptions represents all the available options to list keys in a registry.
type GPGKeyListOptions struct {
	ListOptions

	// Required: A list of one or more namespaces. Must be authorized HCP Terraform or Terraform Enterprise organization names.
	Namespaces []string `url:"filter[namespace]"`
}

// GPGKeyCreateOptions represents all the available options used to create a GPG key.
type GPGKeyCreateOptions struct {
	Type       string `jsonapi:"primary,gpg-keys"`
	Namespace  string `jsonapi:"attr,namespace"`
	AsciiArmor string `jsonapi:"attr,ascii-armor"`
}

// GPGKeyCreateOptions represents all the available options used to update a GPG key.
type GPGKeyUpdateOptions struct {
	Type      string `jsonapi:"primary,gpg-keys"`
	Namespace string `jsonapi:"attr,namespace"`
}

// ListPrivate lists the private registry GPG keys for specified namespaces.
func (s *gpgKeys) ListPrivate(ctx context.Context, options GPGKeyListOptions) (*GPGKeyList, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("/api/registry/%s/v2/gpg-keys", url.PathEscape(string(PrivateRegistry)))
	req, err := s.client.NewRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	keyl := &GPGKeyList{}
	err = req.Do(ctx, keyl)
	if err != nil {
		return nil, err
	}

	return keyl, nil
}

func (s *gpgKeys) Create(ctx context.Context, registryName RegistryName, options GPGKeyCreateOptions) (*GPGKey, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	if registryName != PrivateRegistry {
		return nil, ErrInvalidRegistryName
	}

	u := fmt.Sprintf("/api/registry/%s/v2/gpg-keys", url.PathEscape(string(registryName)))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	g := &GPGKey{}
	err = req.Do(ctx, g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (s *gpgKeys) Read(ctx context.Context, keyID GPGKeyID) (*GPGKey, error) {
	if err := keyID.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("/api/registry/%s/v2/gpg-keys/%s/%s",
		url.PathEscape(string(keyID.RegistryName)),
		url.PathEscape(keyID.Namespace),
		url.PathEscape(keyID.KeyID),
	)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	g := &GPGKey{}
	err = req.Do(ctx, g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (s *gpgKeys) Update(ctx context.Context, keyID GPGKeyID, options GPGKeyUpdateOptions) (*GPGKey, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	if err := keyID.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("/api/registry/%s/v2/gpg-keys/%s/%s",
		url.PathEscape(string(keyID.RegistryName)),
		url.PathEscape(keyID.Namespace),
		url.PathEscape(keyID.KeyID),
	)
	req, err := s.client.NewRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	g := &GPGKey{}
	err = req.Do(ctx, g)
	if err != nil {
		if strings.Contains(err.Error(), "namespace not authorized") {
			return nil, ErrNamespaceNotAuthorized
		}
		return nil, err
	}

	return g, nil
}

func (s *gpgKeys) Delete(ctx context.Context, keyID GPGKeyID) error {
	if err := keyID.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("/api/registry/%s/v2/gpg-keys/%s/%s",
		url.PathEscape(string(keyID.RegistryName)),
		url.PathEscape(keyID.Namespace),
		url.PathEscape(keyID.KeyID),
	)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o GPGKeyID) valid() error {
	if o.RegistryName != PrivateRegistry {
		return ErrInvalidRegistryName
	}

	if !validString(&o.Namespace) {
		return ErrInvalidNamespace
	}

	if !validString(&o.KeyID) {
		return ErrInvalidKeyID
	}

	return nil
}

func (o *GPGKeyListOptions) valid() error {
	if len(o.Namespaces) == 0 {
		return ErrInvalidNamespace
	}

	for _, namespace := range o.Namespaces {
		if namespace == "" || !validString(&namespace) {
			return ErrInvalidNamespace
		}
	}

	return nil
}

func (o GPGKeyCreateOptions) valid() error {
	if !validString(&o.Namespace) {
		return ErrInvalidNamespace
	}

	if !validString(&o.AsciiArmor) {
		return ErrInvalidAsciiArmor
	}

	return nil
}

func (o GPGKeyUpdateOptions) valid() error {
	if !validString(&o.Namespace) {
		return ErrInvalidNamespace
	}

	return nil
}
