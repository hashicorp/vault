package google

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const codeURLPath = "code_url"

func pathCodeURL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: codeURLPath,
		Fields: map[string]*framework.FieldSchema{},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathCodeURL,
		},
	}
}

func (b *backend) pathCodeURL(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	config, err := b.Config(req.Storage)

	if err != nil {
		return nil, err
	}

	if config.ApplicationID == "" {
		return logical.ErrorResponse(configErrorMsg), nil
	}

	if config.ApplicationSecret == "" {
		return logical.ErrorResponse(configErrorMsg), nil
	}

	googleConfig := &oauth2.Config{
		ClientID:     config.ApplicationID,
		ClientSecret: config.ApplicationSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{ "email" },
	}

	authURL := googleConfig.AuthCodeURL("state", oauth2.AccessTypeOnline, oauth2.ApprovalForce)

	return &logical.Response{
		Data: map[string]interface{}{
			"url": authURL,
		},
	}, nil
}

