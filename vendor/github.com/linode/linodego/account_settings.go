package linodego

import (
	"context"
)

// AccountSettings are the account wide flags or plans that effect new resources
type AccountSettings struct {
	// The default backups enrollment status for all new Linodes for all users on the account.  When enabled, backups are mandatory per instance.
	BackupsEnabled bool `json:"backups_enabled"`

	// Wether or not Linode Managed service is enabled for the account.
	Managed bool `json:"managed"`

	// Wether or not the Network Helper is enabled for all new Linode Instance Configs on the account.
	NetworkHelper bool `json:"network_helper"`

	// A plan name like "longview-3"..."longview-100", or a nil value for to cancel any existing subscription plan.
	LongviewSubscription *string `json:"longview_subscription"`

	// A string like "disabled", "suspended", or "active" describing the status of this account’s Object Storage service enrollment.
	ObjectStorage *string `json:"object_storage"`
}

// AccountSettingsUpdateOptions are the updateable account wide flags or plans that effect new resources.
type AccountSettingsUpdateOptions struct {
	// The default backups enrollment status for all new Linodes for all users on the account.  When enabled, backups are mandatory per instance.
	BackupsEnabled *bool `json:"backups_enabled,omitempty"`

	// A plan name like "longview-3"..."longview-100", or a nil value for to cancel any existing subscription plan.
	// Deprecated: Use PUT /longview/plan instead to update the LongviewSubscription
	LongviewSubscription *string `json:"longview_subscription,omitempty"`

	// The default network helper setting for all new Linodes and Linode Configs for all users on the account.
	NetworkHelper *bool `json:"network_helper,omitempty"`
}

// GetAccountSettings gets the account wide flags or plans that effect new resources
func (c *Client) GetAccountSettings(ctx context.Context) (*AccountSettings, error) {
	e := "account/settings"

	response, err := doGETRequest[AccountSettings](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateAccountSettings updates the settings associated with the account
func (c *Client) UpdateAccountSettings(ctx context.Context, opts AccountSettingsUpdateOptions) (*AccountSettings, error) {
	e := "account/settings"

	response, err := doPUTRequest[AccountSettings](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}
