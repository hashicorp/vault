package appId

import (
	"context"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLoginWithAppIDPath(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login/(?P<app_id>.+)",
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
			logical.UpdateOperation: b.pathLogin,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
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
			logical.UpdateOperation:         b.pathLogin,
			logical.AliasLookaheadOperation: b.pathLoginAliasLookahead,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLoginAliasLookahead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	appId := data.Get("app_id").(string)

	if appId == "" {
		return nil, fmt.Errorf("missing app_id")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: appId,
			},
		},
	}, nil
}

func (b *backend) pathLogin(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	appId := data.Get("app_id").(string)
	userId := data.Get("user_id").(string)

	var displayName string
	if dispName, resp, err := b.verifyCredentials(ctx, req, appId, userId); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		displayName = dispName
	}

	// Get the policies associated with the app
	policies, err := b.MapAppId.Policies(ctx, req.Storage, appId)
	if err != nil {
		return nil, err
	}

	// Store hashes of the app ID and user ID for the metadata
	appIdHash := sha1.Sum([]byte(appId))
	userIdHash := sha1.Sum([]byte(userId))
	metadata := map[string]string{
		"app-id":  "sha1:" + hex.EncodeToString(appIdHash[:]),
		"user-id": "sha1:" + hex.EncodeToString(userIdHash[:]),
	}

	return &logical.Response{
		Auth: &logical.Auth{
			InternalData: map[string]interface{}{
				"app-id":  appId,
				"user-id": userId,
			},
			DisplayName: displayName,
			Policies:    policies,
			Metadata:    metadata,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
			},
			Alias: &logical.Alias{
				Name: appId,
			},
		},
	}, nil
}

func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	appId := req.Auth.InternalData["app-id"].(string)
	userId := req.Auth.InternalData["user-id"].(string)

	// Skipping CIDR verification to enable renewal from machines other than
	// the ones encompassed by CIDR block.
	if _, resp, err := b.verifyCredentials(ctx, req, appId, userId); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	}

	// Get the policies associated with the app
	mapPolicies, err := b.MapAppId.Policies(ctx, req.Storage, appId)
	if err != nil {
		return nil, err
	}
	if !policyutil.EquivalentPolicies(mapPolicies, req.Auth.TokenPolicies) {
		return nil, fmt.Errorf("policies do not match")
	}

	return &logical.Response{Auth: req.Auth}, nil
}

func (b *backend) verifyCredentials(ctx context.Context, req *logical.Request, appId, userId string) (string, *logical.Response, error) {
	// Ensure both appId and userId are provided
	if appId == "" || userId == "" {
		return "", logical.ErrorResponse("missing 'app_id' or 'user_id'"), nil
	}

	// Look up the apps that this user is allowed to access
	appsMap, err := b.MapUserId.Get(ctx, req.Storage, userId)
	if err != nil {
		return "", nil, err
	}
	if appsMap == nil {
		return "", logical.ErrorResponse("invalid user ID or app ID"), nil
	}

	// If there is a CIDR block restriction, check that
	if raw, ok := appsMap["cidr_block"]; ok {
		_, cidr, err := net.ParseCIDR(raw.(string))
		if err != nil {
			return "", nil, errwrap.Wrapf("invalid restriction cidr: {{err}}", err)
		}

		var addr string
		if req.Connection != nil {
			addr = req.Connection.RemoteAddr
		}
		if addr == "" || !cidr.Contains(net.ParseIP(addr)) {
			return "", logical.ErrorResponse("unauthorized source address"), nil
		}
	}

	appsRaw, ok := appsMap["value"]
	if !ok {
		appsRaw = ""
	}

	apps, ok := appsRaw.(string)
	if !ok {
		return "", nil, fmt.Errorf("mapping is not a string")
	}

	// Verify that the app is in the list
	found := false
	appIdBytes := []byte(appId)
	for _, app := range strings.Split(apps, ",") {
		match := []byte(strings.TrimSpace(app))
		// Protect against a timing attack with the app_id comparison
		if subtle.ConstantTimeCompare(match, appIdBytes) == 1 {
			found = true
		}
	}
	if !found {
		return "", logical.ErrorResponse("invalid user ID or app ID"), nil
	}

	// Get the raw data associated with the app
	appRaw, err := b.MapAppId.Get(ctx, req.Storage, appId)
	if err != nil {
		return "", nil, err
	}
	if appRaw == nil {
		return "", logical.ErrorResponse("invalid user ID or app ID"), nil
	}
	var displayName string
	if raw, ok := appRaw["display_name"]; ok {
		displayName = raw.(string)
	}

	return displayName, nil, nil
}

const pathLoginSyn = `
Log in with an App ID and User ID.
`

const pathLoginDesc = `
This endpoint authenticates using an application ID, user ID and potential the IP address of the connecting client.
`
