package aws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathStaticCreds = "static-creds"

	paramAccessKeyID      = "access_key_id"
	paramSecretsAccessKey = "secret_access_key"
)

type awsCredentials struct {
	AccessKeyID     string `json:"access_key_id" structs:"access_key_id" mapstructure:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key" structs:"secret_access_key" mapstructure:"secret_access_key"`
}

func pathStaticCredentials(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("%s/%s", pathStaticCreds, framework.GenericNameWithAtRegex(paramRoleName)),
		Fields: map[string]*framework.FieldSchema{
			paramRoleName: {
				Type:        framework.TypeString,
				Description: descRoleName,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathStaticCredsRead,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: http.StatusText(http.StatusOK),
						Fields: map[string]*framework.FieldSchema{
							paramAccessKeyID: {
								Type:        framework.TypeString,
								Description: descAccessKeyID,
							},
							paramSecretsAccessKey: {
								Type:        framework.TypeString,
								Description: descSecretAccessKey,
							},
						},
					}},
				},
			},
		},

		HelpSynopsis:    pathStaticCredsHelpSyn,
		HelpDescription: pathStaticCredsHelpDesc,
	}
}

func (b *backend) pathStaticCredsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName, ok := data.GetOk(paramRoleName)
	if !ok {
		return nil, fmt.Errorf("missing '%s' parameter", paramRoleName)
	}

	entry, err := req.Storage.Get(ctx, formatCredsStoragePath(roleName.(string)))
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials for role '%q': %w", roleName, err)
	}
	if entry == nil {
		return nil, nil
	}

	var credentials awsCredentials
	if err := entry.DecodeJSON(&credentials); err != nil {
		return nil, fmt.Errorf("failed to decode credentials: %w", err)
	}

	return &logical.Response{
		Data: structs.New(credentials).Map(),
	}, nil
}

func formatCredsStoragePath(roleName string) string {
	return fmt.Sprintf("%s/%s", pathStaticCreds, roleName)
}

const pathStaticCredsHelpSyn = `
`

const pathStaticCredsHelpDesc = `
`

const (
	descAccessKeyID     = ""
	descSecretAccessKey = ""
)
