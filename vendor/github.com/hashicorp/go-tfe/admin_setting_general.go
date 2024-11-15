// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
)

// Compile-time proof of interface implementation.
var _ GeneralSettings = (*adminGeneralSettings)(nil)

// GeneralSettings describes the general admin settings for the Admin Setting API.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings
type GeneralSettings interface {
	// Read returns the general settings
	Read(ctx context.Context) (*AdminGeneralSetting, error)

	// Update updates general settings.
	Update(ctx context.Context, options AdminGeneralSettingsUpdateOptions) (*AdminGeneralSetting, error)
}

type adminGeneralSettings struct {
	client *Client
}

// AdminGeneralSetting represents a the general settings in Terraform Enterprise.
type AdminGeneralSetting struct {
	ID                               string `jsonapi:"primary,general-settings"`
	LimitUserOrganizationCreation    bool   `jsonapi:"attr,limit-user-organization-creation"`
	APIRateLimitingEnabled           bool   `jsonapi:"attr,api-rate-limiting-enabled"`
	APIRateLimit                     int    `jsonapi:"attr,api-rate-limit"`
	SendPassingStatusesEnabled       bool   `jsonapi:"attr,send-passing-statuses-for-untriggered-speculative-plans"`
	AllowSpeculativePlansOnPR        bool   `jsonapi:"attr,allow-speculative-plans-on-pull-requests-from-forks"`
	RequireTwoFactorForAdmin         bool   `jsonapi:"attr,require-two-factor-for-admins"`
	FairRunQueuingEnabled            bool   `jsonapi:"attr,fair-run-queuing-enabled"`
	LimitOrgsPerUser                 bool   `jsonapi:"attr,limit-organizations-per-user"`
	DefaultOrgsPerUserCeiling        int    `jsonapi:"attr,default-organizations-per-user-ceiling"`
	LimitWorkspacesPerOrg            bool   `jsonapi:"attr,limit-workspaces-per-organization"`
	DefaultWorkspacesPerOrgCeiling   int    `jsonapi:"attr,default-workspaces-per-organization-ceiling"`
	TerraformBuildWorkerApplyTimeout string `jsonapi:"attr,terraform-build-worker-apply-timeout"`
	TerraformBuildWorkerPlanTimeout  string `jsonapi:"attr,terraform-build-worker-plan-timeout"`
	ApplyTimeout                     string `jsonapi:"attr,apply-timeout"`
	PlanTimeout                      string `jsonapi:"attr,plan-timeout"`
	DefaultRemoteStateAccess         bool   `jsonapi:"attr,default-remote-state-access"`
}

// AdminGeneralSettingsUpdateOptions represents the admin options for updating
// general settings.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings#request-body
type AdminGeneralSettingsUpdateOptions struct {
	LimitUserOrgCreation              *bool   `jsonapi:"attr,limit-user-organization-creation,omitempty"`
	APIRateLimitingEnabled            *bool   `jsonapi:"attr,api-rate-limiting-enabled,omitempty"`
	APIRateLimit                      *int    `jsonapi:"attr,api-rate-limit,omitempty"`
	SendPassingStatusUntriggeredPlans *bool   `jsonapi:"attr,send-passing-statuses-for-untriggered-speculative-plans,omitempty"`
	AllowSpeculativePlansOnPR         *bool   `jsonapi:"attr,allow-speculative-plans-on-pull-requests-from-forks,omitempty"`
	DefaultRemoteStateAccess          *bool   `jsonapi:"attr,default-remote-state-access,omitempty"`
	ApplyTimeout                      *string `jsonapi:"attr,apply-timeout"`
	PlanTimeout                       *string `jsonapi:"attr,plan-timeout"`
}

// Read returns the general settings.
func (a *adminGeneralSettings) Read(ctx context.Context) (*AdminGeneralSetting, error) {
	req, err := a.client.NewRequest("GET", "admin/general-settings", nil)
	if err != nil {
		return nil, err
	}

	ags := &AdminGeneralSetting{}
	err = req.Do(ctx, ags)
	if err != nil {
		return nil, err
	}

	return ags, nil
}

// Update updates the general settings.
func (a *adminGeneralSettings) Update(ctx context.Context, options AdminGeneralSettingsUpdateOptions) (*AdminGeneralSetting, error) {
	req, err := a.client.NewRequest("PATCH", "admin/general-settings", &options)
	if err != nil {
		return nil, err
	}

	ags := &AdminGeneralSetting{}
	err = req.Do(ctx, ags)
	if err != nil {
		return nil, err
	}

	return ags, nil
}
