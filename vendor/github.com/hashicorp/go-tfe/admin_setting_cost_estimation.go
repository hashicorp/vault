// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
)

// Compile-time proof of interface implementation.
var _ CostEstimationSettings = (*adminCostEstimationSettings)(nil)

// CostEstimationSettings describes all the cost estimation admin settings for the Admin Setting API.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings
type CostEstimationSettings interface {
	// Read returns the cost estimation settings.
	Read(ctx context.Context) (*AdminCostEstimationSetting, error)

	// Update updates the cost estimation settings.
	Update(ctx context.Context, options AdminCostEstimationSettingOptions) (*AdminCostEstimationSetting, error)
}

type adminCostEstimationSettings struct {
	client *Client
}

// AdminCostEstimationSetting represents the admin cost estimation settings.
type AdminCostEstimationSetting struct {
	ID                        string `jsonapi:"primary,cost-estimation-settings"`
	Enabled                   bool   `jsonapi:"attr,enabled"`
	AWSAccessKeyID            string `jsonapi:"attr,aws-access-key-id"`
	AWSAccessKey              string `jsonapi:"attr,aws-secret-key"`
	AWSEnabled                bool   `jsonapi:"attr,aws-enabled"`
	AWSInstanceProfileEnabled bool   `jsonapi:"attr,aws-instance-profile-enabled"`
	GCPCredentials            string `jsonapi:"attr,gcp-credentials"`
	GCPEnabled                bool   `jsonapi:"attr,gcp-enabled"`
	AzureEnabled              bool   `jsonapi:"attr,azure-enabled"`
	AzureClientID             string `jsonapi:"attr,azure-client-id"`
	AzureClientSecret         string `jsonapi:"attr,azure-client-secret"`
	AzureSubscriptionID       string `jsonapi:"attr,azure-subscription-id"`
	AzureTenantID             string `jsonapi:"attr,azure-tenant-id"`
}

// AdminCostEstimationSettingOptions represents the admin options for updating
// the cost estimation settings.
// https://developer.hashicorp.com/terraform/enterprise/api-docs/admin/settings#request-body-1
type AdminCostEstimationSettingOptions struct {
	Enabled             *bool   `jsonapi:"attr,enabled,omitempty"`
	AWSAccessKeyID      *string `jsonapi:"attr,aws-access-key-id,omitempty"`
	AWSAccessKey        *string `jsonapi:"attr,aws-secret-key,omitempty"`
	GCPCredentials      *string `jsonapi:"attr,gcp-credentials,omitempty"`
	AzureClientID       *string `jsonapi:"attr,azure-client-id,omitempty"`
	AzureClientSecret   *string `jsonapi:"attr,azure-client-secret,omitempty"`
	AzureSubscriptionID *string `jsonapi:"attr,azure-subscription-id,omitempty"`
	AzureTenantID       *string `jsonapi:"attr,azure-tenant-id,omitempty"`
}

// Read returns the cost estimation settings.
func (a *adminCostEstimationSettings) Read(ctx context.Context) (*AdminCostEstimationSetting, error) {
	req, err := a.client.NewRequest("GET", "admin/cost-estimation-settings", nil)
	if err != nil {
		return nil, err
	}

	ace := &AdminCostEstimationSetting{}
	err = req.Do(ctx, ace)
	if err != nil {
		return nil, err
	}

	return ace, nil
}

// Update updates the cost-estimation settings.
func (a *adminCostEstimationSettings) Update(ctx context.Context, options AdminCostEstimationSettingOptions) (*AdminCostEstimationSetting, error) {
	req, err := a.client.NewRequest("PATCH", "admin/cost-estimation-settings", &options)
	if err != nil {
		return nil, err
	}

	ace := &AdminCostEstimationSetting{}
	err = req.Do(ctx, ace)
	if err != nil {
		return nil, err
	}

	return ace, nil
}
