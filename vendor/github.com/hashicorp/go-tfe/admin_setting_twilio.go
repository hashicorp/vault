// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
)

// Compile-time proof of interface implementation.
var _ TwilioSettings = (*adminTwilioSettings)(nil)

// TwilioSettings describes all the Twilio admin settings for the Admin Setting API.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings
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
	req, err := a.client.NewRequest("GET", "admin/twilio-settings", nil)
	if err != nil {
		return nil, err
	}

	twilio := &AdminTwilioSetting{}
	err = req.Do(ctx, twilio)
	if err != nil {
		return nil, err
	}

	return twilio, nil
}

// AdminTwilioSettingsUpdateOptions represents the admin options for updating
// Twilio settings.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings#request-body-4
type AdminTwilioSettingsUpdateOptions struct {
	Enabled    *bool   `jsonapi:"attr,enabled,omitempty"`
	AccountSid *string `jsonapi:"attr,account-sid,omitempty"`
	AuthToken  *string `jsonapi:"attr,auth-token,omitempty"`
	FromNumber *string `jsonapi:"attr,from-number,omitempty"`
}

// AdminTwilioSettingsVerifyOptions represents the test number to verify Twilio.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings#verify-twilio-settings
type AdminTwilioSettingsVerifyOptions struct {
	TestNumber *string `jsonapi:"attr,test-number"` // Required
}

// Update updates the Twilio settings.
func (a *adminTwilioSettings) Update(ctx context.Context, options AdminTwilioSettingsUpdateOptions) (*AdminTwilioSetting, error) {
	req, err := a.client.NewRequest("PATCH", "admin/twilio-settings", &options)
	if err != nil {
		return nil, err
	}

	twilio := &AdminTwilioSetting{}
	err = req.Do(ctx, twilio)
	if err != nil {
		return nil, err
	}

	return twilio, nil
}

// Verify verifies Twilio settings.
func (a *adminTwilioSettings) Verify(ctx context.Context, options AdminTwilioSettingsVerifyOptions) error {
	if err := options.valid(); err != nil {
		return err
	}
	req, err := a.client.NewRequest("PATCH", "admin/twilio-settings/verify", &options)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o AdminTwilioSettingsVerifyOptions) valid() error {
	if !validString(o.TestNumber) {
		return ErrRequiredTestNumber
	}

	return nil
}
