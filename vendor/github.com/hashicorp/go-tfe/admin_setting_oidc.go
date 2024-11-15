// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
)

// Compile-time proof of interface implementation.
var _ OIDCSettings = (*adminOIDCSettings)(nil)

// OidcSettings describes all the OIDC admin settings for the Admin Setting API.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings
type OIDCSettings interface {
	// Rotate the key used for signing OIDC tokens for workload identity
	RotateKey(ctx context.Context) error

	// Trim old version of the key used for signing OIDC tokens for workload identity
	TrimKey(ctx context.Context) error
}

type adminOIDCSettings struct {
	client *Client
}

// Rotate the key used for signing OIDC tokens for workload identity
func (a *adminOIDCSettings) RotateKey(ctx context.Context) error {
	req, err := a.client.NewRequest("POST", "admin/oidc-settings/actions/rotate-key", nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Trim old version of the key used for signing OIDC tokens for workload identity
func (a *adminOIDCSettings) TrimKey(ctx context.Context) error {
	req, err := a.client.NewRequest("POST", "admin/oidc-settings/actions/trim-key", nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}
