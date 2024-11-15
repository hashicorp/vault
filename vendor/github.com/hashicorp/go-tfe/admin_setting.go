// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

// AdminSettings describes all the admin settings related methods that the Terraform Enterprise API supports.
// Note that admin settings are only available in Terraform Enterprise.
//
// TFE API docs: https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings
type AdminSettings struct {
	General        GeneralSettings
	SAML           SAMLSettings
	CostEstimation CostEstimationSettings
	SMTP           SMTPSettings
	Twilio         TwilioSettings
	Customization  CustomizationSettings
	OIDC           OIDCSettings
}

func newAdminSettings(client *Client) *AdminSettings {
	return &AdminSettings{
		General:        &adminGeneralSettings{client: client},
		SAML:           &adminSAMLSettings{client: client},
		CostEstimation: &adminCostEstimationSettings{client: client},
		SMTP:           &adminSMTPSettings{client: client},
		Twilio:         &adminTwilioSettings{client: client},
		Customization:  &adminCustomizationSettings{client: client},
		OIDC:           &adminOIDCSettings{client: client},
	}
}
