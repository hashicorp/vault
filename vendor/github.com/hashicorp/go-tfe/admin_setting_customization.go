package tfe

import (
	"context"
)

// Compile-time proof of interface implementation.
var _ CustomizationSettings = (*adminCustomizationSettings)(nil)

// CustomizationSettings describes all the Customization admin settings.
type CustomizationSettings interface {
	// Read returns the customization settings.
	Read(ctx context.Context) (*AdminCustomizationSetting, error)

	// Update updates the customization settings.
	Update(ctx context.Context, options AdminCustomizationSettingsUpdateOptions) (*AdminCustomizationSetting, error)
}

type adminCustomizationSettings struct {
	client *Client
}

// AdminCustomizationSetting represents the Customization settings in Terraform Enterprise.
type AdminCustomizationSetting struct {
	ID           string `jsonapi:"primary,customization-settings"`
	SupportEmail string `jsonapi:"attr,support-email-address"`
	LoginHelp    string `jsonapi:"attr,login-help"`
	Footer       string `jsonapi:"attr,footer"`
	Error        string `jsonapi:"attr,error"`
	NewUser      string `jsonapi:"attr,new-user"`
}

// Read returns the Customization settings.
func (a *adminCustomizationSettings) Read(ctx context.Context) (*AdminCustomizationSetting, error) {
	req, err := a.client.newRequest("GET", "admin/customization-settings", nil)
	if err != nil {
		return nil, err
	}

	cs := &AdminCustomizationSetting{}
	err = a.client.do(ctx, req, cs)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

// AdminCustomizationSettingsUpdateOptions represents the admin options for updating
// Customization settings.
// https://www.terraform.io/docs/cloud/api/admin/settings.html#request-body-6
type AdminCustomizationSettingsUpdateOptions struct {
	SupportEmail *string `jsonapi:"attr,support-email-address,omitempty"`
	LoginHelp    *string `jsonapi:"attr,login-help,omitempty"`
	Footer       *string `jsonapi:"attr,footer,omitempty"`
	Error        *string `jsonapi:"attr,error,omitempty"`
	NewUser      *string `jsonapi:"attr,new-user,omitempty"`
}

// Update updates the customization settings.
func (a *adminCustomizationSettings) Update(ctx context.Context, options AdminCustomizationSettingsUpdateOptions) (*AdminCustomizationSetting, error) {
	req, err := a.client.newRequest("PATCH", "admin/customization-settings", &options)
	if err != nil {
		return nil, err
	}

	cs := &AdminCustomizationSetting{}
	err = a.client.do(ctx, req, cs)
	if err != nil {
		return nil, err
	}

	return cs, nil
}
