package appId

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
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
			logical.UpdateOperation: b.pathLogin,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLogin(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	appId := data.Get("app_id").(string)
	userId := data.Get("user_id").(string)

	// Ensure both appId and userId are provided
	if appId == "" || userId == "" {
		return logical.ErrorResponse("missing 'app_id' or 'user_id'"), nil
	}

	user, err := b.User(req.Storage, userId)
	if err != nil {
		return nil, err
	}

	// If there is a CIDR block restriction, check that
	if user.CidrBlock != nil {
		var addr string
		if req.Connection != nil {
			addr = req.Connection.RemoteAddr
		}
		if addr == "" || !user.CidrBlock.Contains(net.ParseIP(addr)) {
			return logical.ErrorResponse("unauthorized source address"), nil
		}
	}

	// Verify that the app is in the list
	found := false
	appIdBytes := []byte(strings.Join(user.AppIds, ","))
	for _, app := range user.AppIds {
		match := []byte(strings.TrimSpace(app))
		// Protect against a timing attack with the app_id comparison
		if subtle.ConstantTimeCompare(match, appIdBytes) == 1 {
			found = true
		}
	}
	if !found {
		return logical.ErrorResponse("invalid user ID or app ID"), nil
	}

	app, err := b.App(req.Storage, appId)
	if err != nil {
		return nil, err
	}

	// Get the policies associated with the app
	policies, err := b.MapAppId.Policies(req.Storage, appId)
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
			DisplayName: app.DisplayName,
			Policies:    policies,
			Metadata:    metadata,
			LeaseOptions: logical.LeaseOptions{
				TTL:         app.TTL,
				GracePeriod: app.TTL / 10,
				Renewable:   app.Renewable,
			},
			InternalData: map[string]interface{}{
				"user-id": userId,
				"app-id": appId,
			},
		},
	}, nil
}

func (b *backend) pathLoginRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	appId, ok := req.Auth.InternalData["app-id"].(string)
	if ! ok {
		return logical.ErrorResponse("could not find the app-id for the token provided"), nil
	}

	app, err := b.App(req.Storage, appId)
	if err != nil {
		return nil, err
	}

	if ! app.Renewable {
		return logical.ErrorResponse("lease is not renewable"), nil
	}

	return framework.LeaseExtend(app.MaxTTL, 0, false)(req, d)
}

const pathLoginSyn = `
Log in with an App ID and User ID.
`

const pathLoginDesc = `
This endpoint authenticates using an application ID, user ID and potential the IP address of the connecting client.
`
