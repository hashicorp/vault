// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
)

// Compile-time proof of interface implementation.
var _ SAMLSettings = (*adminSAMLSettings)(nil)

// SAMLSettings describes all the SAML admin settings for the Admin Setting API.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings
type SAMLSettings interface {
	// Read returns the SAML settings.
	Read(ctx context.Context) (*AdminSAMLSetting, error)

	// Update updates the SAML settings.
	Update(ctx context.Context, options AdminSAMLSettingsUpdateOptions) (*AdminSAMLSetting, error)

	// RevokeIdpCert revokes the older IdP certificate when the new IdP
	// certificate is known to be functioning correctly.
	RevokeIdpCert(ctx context.Context) (*AdminSAMLSetting, error)
}

type adminSAMLSettings struct {
	client *Client
}

// AdminSAMLSetting represents the SAML settings in Terraform Enterprise.
type AdminSAMLSetting struct {
	ID                        string `jsonapi:"primary,saml-settings"`
	Enabled                   bool   `jsonapi:"attr,enabled"`
	Debug                     bool   `jsonapi:"attr,debug"`
	AuthnRequestsSigned       bool   `jsonapi:"attr,authn-requests-signed"`
	WantAssertionsSigned      bool   `jsonapi:"attr,want-assertions-signed"`
	TeamManagementEnabled     bool   `jsonapi:"attr,team-management-enabled"`
	OldIDPCert                string `jsonapi:"attr,old-idp-cert"`
	IDPCert                   string `jsonapi:"attr,idp-cert"`
	SLOEndpointURL            string `jsonapi:"attr,slo-endpoint-url"`
	SSOEndpointURL            string `jsonapi:"attr,sso-endpoint-url"`
	AttrUsername              string `jsonapi:"attr,attr-username"`
	AttrGroups                string `jsonapi:"attr,attr-groups"`
	AttrSiteAdmin             string `jsonapi:"attr,attr-site-admin"`
	SiteAdminRole             string `jsonapi:"attr,site-admin-role"`
	SSOAPITokenSessionTimeout int    `jsonapi:"attr,sso-api-token-session-timeout"`
	ACSConsumerURL            string `jsonapi:"attr,acs-consumer-url"`
	MetadataURL               string `jsonapi:"attr,metadata-url"`
	Certificate               string `jsonapi:"attr,certificate"`
	PrivateKey                string `jsonapi:"attr,private-key"`
	SignatureSigningMethod    string `jsonapi:"attr,signature-signing-method"`
	SignatureDigestMethod     string `jsonapi:"attr,signature-digest-method"`
}

// Read returns the SAML settings.
func (a *adminSAMLSettings) Read(ctx context.Context) (*AdminSAMLSetting, error) {
	req, err := a.client.NewRequest("GET", "admin/saml-settings", nil)
	if err != nil {
		return nil, err
	}

	saml := &AdminSAMLSetting{}
	err = req.Do(ctx, saml)
	if err != nil {
		return nil, err
	}

	return saml, nil
}

// AdminSAMLSettingsUpdateOptions represents the admin options for updating
// SAML settings.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings#request-body-2
type AdminSAMLSettingsUpdateOptions struct {
	Enabled                   *bool   `jsonapi:"attr,enabled,omitempty"`
	Debug                     *bool   `jsonapi:"attr,debug,omitempty"`
	IDPCert                   *string `jsonapi:"attr,idp-cert,omitempty"`
	Certificate               *string `jsonapi:"attr,certificate,omitempty"`
	PrivateKey                *string `jsonapi:"attr,private-key,omitempty"`
	SLOEndpointURL            *string `jsonapi:"attr,slo-endpoint-url,omitempty"`
	SSOEndpointURL            *string `jsonapi:"attr,sso-endpoint-url,omitempty"`
	AttrUsername              *string `jsonapi:"attr,attr-username,omitempty"`
	AttrGroups                *string `jsonapi:"attr,attr-groups,omitempty"`
	AttrSiteAdmin             *string `jsonapi:"attr,attr-site-admin,omitempty"`
	SiteAdminRole             *string `jsonapi:"attr,site-admin-role,omitempty"`
	SSOAPITokenSessionTimeout *int    `jsonapi:"attr,sso-api-token-session-timeout,omitempty"`
	TeamManagementEnabled     *bool   `jsonapi:"attr,team-management-enabled,omitempty"`
	AuthnRequestsSigned       *bool   `jsonapi:"attr,authn-requests-signed,omitempty"`
	WantAssertionsSigned      *bool   `jsonapi:"attr,want-assertions-signed,omitempty"`
	SignatureSigningMethod    *string `jsonapi:"attr,signature-signing-method,omitempty"`
	SignatureDigestMethod     *string `jsonapi:"attr,signature-digest-method,omitempty"`
}

// Update updates the SAML settings.
func (a *adminSAMLSettings) Update(ctx context.Context, options AdminSAMLSettingsUpdateOptions) (*AdminSAMLSetting, error) {
	req, err := a.client.NewRequest("PATCH", "admin/saml-settings", &options)
	if err != nil {
		return nil, err
	}

	saml := &AdminSAMLSetting{}
	err = req.Do(ctx, saml)
	if err != nil {
		return nil, err
	}

	return saml, nil
}

// RevokeIdpCert revokes the older IdP certificate when the new IdP
// certificate is known to be functioning correctly.
func (a *adminSAMLSettings) RevokeIdpCert(ctx context.Context) (*AdminSAMLSetting, error) {
	req, err := a.client.NewRequest("POST", "admin/saml-settings/actions/revoke-old-certificate", nil)
	if err != nil {
		return nil, err
	}

	saml := &AdminSAMLSetting{}
	err = req.Do(ctx, saml)
	if err != nil {
		return nil, err
	}

	return saml, nil
}
