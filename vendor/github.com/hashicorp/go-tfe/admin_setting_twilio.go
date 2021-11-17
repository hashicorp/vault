package tfe

import (
	"context"
)

// Compile-time proof of interface implementation.
var _ TwilioSettings = (*adminTwilioSettings)(nil)

// TwilioSettings describes all the Twilio admin settings.
type TwilioSettings interface {
	// Read returns the Twilio settings.
	Read(ctx context.Context) (*AdminTwilioSetting, error)

	// Update updates Twilio settings.
	Update(ctx context.Context, options AdminTwilioSettingsUpdateOptions) (*AdminTwilioSetting, error)

	// Verify verifies Twilio settings.
	Verify(ctx context.Context, options AdminTwilioSettingsVerifyOptions) error
}

type adminTwilioSettings struct {
	client *Client
}

// AdminTwilioSetting represents the Twilio settings in Terraform Enterprise.
type AdminTwilioSetting struct {
	ID         string `jsonapi:"primary,twilio-settings"`
	Enabled    bool   `jsonapi:"attr,enabled"`
	AccountSid string `jsonapi:"attr,account-sid"`
	FromNumber string `jsonapi:"attr,from-number"`
}

// Read returns the Twilio settings.
func (a *adminTwilioSettings) Read(ctx context.Context) (*AdminTwilioSetting, error) {
	req, err := a.client.newRequest("GET", "admin/twilio-settings", nil)
	if err != nil {
		return nil, err
	}

	twilio := &AdminTwilioSetting{}
	err = a.client.do(ctx, req, twilio)
	if err != nil {
		return nil, err
	}

	return twilio, nil
}

// AdminTwilioSettingsUpdateOptions represents the admin options for updating
// Twilio settings.
// https://www.terraform.io/docs/cloud/api/admin/settings.html#request-body-4
type AdminTwilioSettingsUpdateOptions struct {
	Enabled    *bool   `jsonapi:"attr,enabled,omitempty"`
	AccountSid *string `jsonapi:"attr,account-sid,omitempty"`
	AuthToken  *string `jsonapi:"attr,auth-token,omitempty"`
	FromNumber *string `jsonapi:"attr,from-number,omitempty"`
}

// Update updates the Twilio settings.
func (a *adminTwilioSettings) Update(ctx context.Context, options AdminTwilioSettingsUpdateOptions) (*AdminTwilioSetting, error) {
	req, err := a.client.newRequest("PATCH", "admin/twilio-settings", &options)
	if err != nil {
		return nil, err
	}

	twilio := &AdminTwilioSetting{}
	err = a.client.do(ctx, req, twilio)
	if err != nil {
		return nil, err
	}

	return twilio, nil
}

// AdminTwilioSettingsVerifyOptions represents the test number to verify Twilio.
// https://www.terraform.io/docs/cloud/api/admin/settings.html#verify-twilio-settings
type AdminTwilioSettingsVerifyOptions struct {
	TestNumber *string `jsonapi:"attr,test-number"`
}

// Verify verifies Twilio settings.
func (a *adminTwilioSettings) Verify(ctx context.Context, options AdminTwilioSettingsVerifyOptions) error {
	req, err := a.client.newRequest("PATCH", "admin/twilio-settings/verify", &options)
	if err != nil {
		return err
	}

	return a.client.do(ctx, req, nil)
}
