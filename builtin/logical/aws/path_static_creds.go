// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

	paramAccessKeyID      = "access_key"
	paramSecretsAccessKey = "secret_key"
	paramSESSMTPPassword  = "ses_smtp_password"
)

type awsCredentials struct {
	AccessKeyID     string `json:"access_key" structs:"access_key" mapstructure:"access_key"`
	SecretAccessKey string `json:"secret_key" structs:"secret_key" mapstructure:"secret_key"`
	SESSMTPPassword string `json:"ses_smtp_password" structs:"ses_smtp_password" mapstructures:"ses_smtp_password"`
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
							paramSESSMTPPassword: {
								Type:        framework.TypeString,
								Description: descSESSMTPPassword,
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
		return nil, fmt.Errorf("missing %q parameter", paramRoleName)
	}

	entry, err := req.Storage.Get(ctx, formatCredsStoragePath(roleName.(string)))
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials for role %q: %w", roleName, err)
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

const pathStaticCredsHelpSyn = `Retrieve static credentials from the named role.`

const pathStaticCredsHelpDesc = `
This path reads AWS credentials for a certain static role. The keys are rotated
periodically according to their configuration, and will return the same password
until they are rotated.`

const (
	descAccessKeyID     = "The access key of the AWS Credential"
	descSecretAccessKey = "The secret key of the AWS Credential"
	descSESSMTPPassword = "Secret access key converted into an SES SMTP password"
)
