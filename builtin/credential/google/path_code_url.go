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

	authURL := codeUrl(config.ApplicationID, config.ApplicationSecret)

	return &logical.Response{
		Data: map[string]interface{}{
			"url": authURL,
		},
	}, nil
}

func codeUrl(applicationId string, applicationSecret string) string {

	googleConfig := &oauth2.Config{
		ClientID:     applicationId,
		ClientSecret: applicationSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{ "email" },
	}

	authURL := googleConfig.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	return authURL
}

