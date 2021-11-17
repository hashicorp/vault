package tfe

// AdminSettings describes all the admin settings related methods that the Terraform Enterprise API supports.
// Note that admin settings are only available in Terraform Enterprise.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/admin/settings.html
// AdminSettings todo
type AdminSettings struct {
	General        GeneralSettings
	SAML           SAMLSettings
	CostEstimation CostEstimationSettings
	SMTP           SMTPSettings
	Twilio         TwilioSettings
	Customization  CustomizationSettings
}

func newAdminSettings(client *Client) *AdminSettings {
	return &AdminSettings{
		General:        &adminGeneralSettings{client: client},
		SAML:           &adminSAMLSettings{client: client},
		CostEstimation: &adminCostEstimationSettings{client: client},
		SMTP:           &adminSMTPSettings{client: client},
		Twilio:         &adminTwilioSettings{client: client},
		Customization:  &adminCustomizationSettings{client: client},
	}
}
