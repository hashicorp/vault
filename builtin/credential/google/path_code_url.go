package google

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"

	"golang.org/x/oauth2"
)

const codeURLPath = "code_url"
const codeURLResponsePropertyName = "url"
const readCodeUrlPathHelp = "run 'vault read auth/" + BackendName + "/" + codeURLPath + "' for a link to obtain an auth code from google"

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

	if config.ApplicationID == "" || config.ApplicationSecret == "" {
		return logical.ErrorResponse(writeConfigPathHelp), nil
	}

	authURL := codeUrl(config.ApplicationID, config.ApplicationSecret)

	return &logical.Response{
		Data: map[string]interface{}{
			codeURLResponsePropertyName: authURL,
		},
	}, nil
}

func codeUrl(applicationId string, applicationSecret string) string {

	googleConfig := applicationOauth2Config(applicationId, applicationSecret)

	authURL := googleConfig.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	return authURL
}

