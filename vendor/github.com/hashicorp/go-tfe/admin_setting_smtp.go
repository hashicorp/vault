// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
)

// Compile-time proof of interface implementation.
var _ SMTPSettings = (*adminSMTPSettings)(nil)

// SMTPSettings describes all the SMTP admin settings for the Admin Setting API https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings
type SMTPSettings interface {
	// Read returns the SMTP settings.
	Read(ctx context.Context) (*AdminSMTPSetting, error)

	// Update updates SMTP settings.
	Update(ctx context.Context, options AdminSMTPSettingsUpdateOptions) (*AdminSMTPSetting, error)
}

type adminSMTPSettings struct {
	client *Client
}

// SMTPAuthType represents valid SMTP Auth types.
type SMTPAuthType string

// List of all SMTP auth types.
const (
	SMTPAuthNone  SMTPAuthType = "none"
	SMTPAuthPlain SMTPAuthType = "plain"
	SMTPAuthLogin SMTPAuthType = "login"
)

// AdminSMTPSetting represents a the SMTP settings in Terraform Enterprise.
type AdminSMTPSetting struct {
	ID       string       `jsonapi:"primary,smtp-settings"`
	Enabled  bool         `jsonapi:"attr,enabled"`
	Host     string       `jsonapi:"attr,host"`
	Port     int          `jsonapi:"attr,port"`
	Sender   string       `jsonapi:"attr,sender"`
	Auth     SMTPAuthType `jsonapi:"attr,auth"`
	Username string       `jsonapi:"attr,username"`
}

// Read returns the SMTP settings.
func (a *adminSMTPSettings) Read(ctx context.Context) (*AdminSMTPSetting, error) {
	req, err := a.client.NewRequest("GET", "admin/smtp-settings", nil)
	if err != nil {
		return nil, err
	}

	smtp := &AdminSMTPSetting{}
	err = req.Do(ctx, smtp)
	if err != nil {
		return nil, err
	}

	return smtp, nil
}

// AdminSMTPSettingsUpdateOptions represents the admin options for updating
// SMTP settings.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings#request-body-3
type AdminSMTPSettingsUpdateOptions struct {
	Enabled          *bool         `jsonapi:"attr,enabled,omitempty"`
	Host             *string       `jsonapi:"attr,host,omitempty"`
	Port             *int          `jsonapi:"attr,port,omitempty"`
	Sender           *string       `jsonapi:"attr,sender,omitempty"`
	Auth             *SMTPAuthType `jsonapi:"attr,auth,omitempty"`
	Username         *string       `jsonapi:"attr,username,omitempty"`
	Password         *string       `jsonapi:"attr,password,omitempty"`
	TestEmailAddress *string       `jsonapi:"attr,test-email-address,omitempty"`
}

// Update updates the SMTP settings.
func (a *adminSMTPSettings) Update(ctx context.Context, options AdminSMTPSettingsUpdateOptions) (*AdminSMTPSetting, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := a.client.NewRequest("PATCH", "admin/smtp-settings", &options)
	if err != nil {
		return nil, err
	}

	smtp := &AdminSMTPSetting{}
	err = req.Do(ctx, smtp)
	if err != nil {
		return nil, err
	}

	return smtp, nil
}

func (o AdminSMTPSettingsUpdateOptions) valid() error {
	if validString((*string)(o.Auth)) {
		if err := validateAdminSettingSMTPAuth(*o.Auth); err != nil {
			return err
		}
	}

	return nil
}

func validateAdminSettingSMTPAuth(authVal SMTPAuthType) error {
	switch authVal {
	case SMTPAuthNone, SMTPAuthPlain, SMTPAuthLogin:
		// do nothing
	default:
		return ErrInvalidSMTPAuth
	}

	return nil
}
