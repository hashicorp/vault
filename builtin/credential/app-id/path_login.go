package appId

import (
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
	apps, err := b.MapUserId.Get(req.Storage, userId)
	if err != nil {
		return nil, err
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

	return &logical.Response{
		Auth: &logical.Auth{
			Policies: policies,
		},
	}, nil
}
