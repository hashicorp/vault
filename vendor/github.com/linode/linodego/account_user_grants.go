package linodego

import (
	"context"
)

type GrantPermissionLevel string

const (
	AccessLevelReadOnly  GrantPermissionLevel = "read_only"
	AccessLevelReadWrite GrantPermissionLevel = "read_write"
)

type GlobalUserGrants struct {
	AccountAccess        *GrantPermissionLevel `json:"account_access"`
	AddDatabases         bool                  `json:"add_databases"`
	AddDomains           bool                  `json:"add_domains"`
	AddFirewalls         bool                  `json:"add_firewalls"`
	AddImages            bool                  `json:"add_images"`
	AddLinodes           bool                  `json:"add_linodes"`
	AddLongview          bool                  `json:"add_longview"`
	AddNodeBalancers     bool                  `json:"add_nodebalancers"`
	AddStackScripts      bool                  `json:"add_stackscripts"`
	AddVolumes           bool                  `json:"add_volumes"`
	CancelAccount        bool                  `json:"cancel_account"`
	LongviewSubscription bool                  `json:"longview_subscription"`
}

type EntityUserGrant struct {
	ID          int                   `json:"id"`
	Permissions *GrantPermissionLevel `json:"permissions"`
}

type GrantedEntity struct {
	ID          int                  `json:"id"`
	Label       string               `json:"label"`
	Permissions GrantPermissionLevel `json:"permissions"`
}

type UserGrants struct {
	Database     []GrantedEntity `json:"database"`
	Domain       []GrantedEntity `json:"domain"`
	Firewall     []GrantedEntity `json:"firewall"`
	Image        []GrantedEntity `json:"image"`
	Linode       []GrantedEntity `json:"linode"`
	Longview     []GrantedEntity `json:"longview"`
	NodeBalancer []GrantedEntity `json:"nodebalancer"`
	StackScript  []GrantedEntity `json:"stackscript"`
	Volume       []GrantedEntity `json:"volume"`

	Global GlobalUserGrants `json:"global"`
}

type UserGrantsUpdateOptions struct {
	Database     []GrantedEntity   `json:"database,omitempty"`
	Domain       []EntityUserGrant `json:"domain,omitempty"`
	Firewall     []EntityUserGrant `json:"firewall,omitempty"`
	Image        []EntityUserGrant `json:"image,omitempty"`
	Linode       []EntityUserGrant `json:"linode,omitempty"`
	Longview     []EntityUserGrant `json:"longview,omitempty"`
	NodeBalancer []EntityUserGrant `json:"nodebalancer,omitempty"`
	StackScript  []EntityUserGrant `json:"stackscript,omitempty"`
	Volume       []EntityUserGrant `json:"volume,omitempty"`

	Global GlobalUserGrants `json:"global"`
}

func (c *Client) GetUserGrants(ctx context.Context, username string) (*UserGrants, error) {
	e := formatAPIPath("account/users/%s/grants", username)
	response, err := doGETRequest[UserGrants](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) UpdateUserGrants(ctx context.Context, username string, opts UserGrantsUpdateOptions) (*UserGrants, error) {
	e := formatAPIPath("account/users/%s/grants", username)
	response, err := doPUTRequest[UserGrants](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}
