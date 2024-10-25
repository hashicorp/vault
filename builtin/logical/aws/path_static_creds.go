// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathStaticCreds = "static-creds"

	paramAccessKeyID      = "access_key"
	paramSecretsAccessKey = "secret_key"
)

type awsCredentials struct {
	AccessKeyID     string     `json:"access_key" structs:"access_key" mapstructure:"access_key"`
	Expiration      *time.Time `json:"expiration,omitempty" structs:"expiration" mapstructure:"expiration"`
	SecretAccessKey string     `json:"secret_key" structs:"secret_key" mapstructure:"secret_key"`
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

func (a *awsCredentials) priority(role staticRoleEntry) int64 {
	if a.Expiration != nil {
		return a.Expiration.Unix()
	}
	return time.Now().Add(role.RotationPeriod).Unix()
}

const pathStaticCredsHelpSyn = `Retrieve static credentials from the named role.`

const pathStaticCredsHelpDesc = `
This path reads AWS credentials for a certain static role. The keys are rotated
periodically according to their configuration, and will return the same password
until they are rotated.`

const (
	descAccessKeyID     = "The access key of the AWS Credential"
	descSecretAccessKey = "The secret key of the AWS Credential"
)
