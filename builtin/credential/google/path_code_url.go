package google

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const PATH_CODE_URL = "code_url"

func pathCodeUrl(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: PATH_CODE_URL,
		Fields: map[string]*framework.FieldSchema{},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathCodeUrl,
		},
	}
}

func (b *backend) pathCodeUrl(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	config, err := b.Config(req.Storage)

	if err != nil {
		return nil, err
	}

	if config.ApplicationId == "" {
		return logical.ErrorResponse(
			"configure the google credential backend with applicationId first"), nil
	}

	if config.ApplicationSecret == "" {
		return logical.ErrorResponse(
			"configure the google credential backend with applicationSecret first"), nil
	}

	googleConfig := &oauth2.Config{
		ClientID:     config.ApplicationId,
		ClientSecret: config.ApplicationSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{ "email" },
	}

	authUrl := googleConfig.AuthCodeURL("state", oauth2.AccessTypeOnline, oauth2.ApprovalForce)

	return &logical.Response{
		Data: map[string]interface{}{
			"url": authUrl,
		},
	}, nil
}

