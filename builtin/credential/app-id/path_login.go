package appId

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields: map[string]*framework.FieldSchema{
			"app_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The unique app ID",
			},

			"user_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The unique user ID",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathLogin,
		},
	}
}

func (b *backend) pathLogin(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	appId := data.Get("app_id").(string)
	userId := data.Get("user_id").(string)

	// Look up the apps that this user is allowed to access
	appsMap, err := b.MapUserId.Get(req.Storage, userId)
	if err != nil {
		return nil, err
	}

	// If there is a CIDR block restriction, check that
	if raw, ok := appsMap["cidr_block"]; ok {
		_, cidr, err := net.ParseCIDR(raw.(string))
		if err != nil {
			return nil, fmt.Errorf("invalid restriction cidr: %s", err)
		}

		var addr string
		if req.Connection != nil {
			addr = req.Connection.RemoteAddr
		}
		if addr == "" || !cidr.Contains(net.ParseIP(addr)) {
			return logical.ErrorResponse("unauthorized source address"), nil
		}
	}

	appsRaw, ok := appsMap["value"]
	if !ok {
		appsRaw = ""
	}

	apps, ok := appsRaw.(string)
	if !ok {
		return nil, fmt.Errorf("internal error: mapping is not a string")
	}

	// Verify that the app is in the list
	found := false
	for _, app := range strings.Split(apps, ",") {
		if strings.TrimSpace(app) == appId {
			found = true
		}
	}
	if !found {
		return logical.ErrorResponse("invalid user ID or app ID"), nil
	}

	// Get the policies associated with the app
	policies, err := b.MapAppId.Policies(req.Storage, appId)
	if err != nil {
		return nil, err
	}

	// Get the raw data associated with the app
	appRaw, err := b.MapAppId.Get(req.Storage, appId)
	if err != nil {
		return nil, err
	}

	// Check if we have a display name
	var displayName string
	if raw, ok := appRaw["display_name"]; ok {
		displayName = raw.(string)
	}

	return &logical.Response{
		Auth: &logical.Auth{
			DisplayName: displayName,
			Policies:    policies,
		},
	}, nil
}
