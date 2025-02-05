// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfigRotateRoot(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/rotate-root",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
			OperationSuffix: "root-iam-credentials",
			OperationVerb:   "rotate",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.pathConfigRotateRootUpdate,
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathConfigRotateRootHelpSyn,
		HelpDescription: pathConfigRotateRootHelpDesc,
	}
}

func (b *backend) pathConfigRotateRootUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := getConfigFromStorage(ctx, req)
	if err != nil {
		return nil, err
	}

	if config.AccessKey == "" || config.SecretKey == "" {
		return logical.ErrorResponse("Cannot call config/rotate-root when either access_key or secret_key is empty"), nil
	}

	if err := b.rotateRootCredential(ctx, req); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"access_key": config.AccessKey,
		},
	}, nil
}

func (b *backend) rotateRootCredential(ctx context.Context, req *logical.Request) error {
	// have to get the client config first because that takes out a read lock
	client, err := b.clientIAM(ctx, req.Storage)
	if err != nil {
		return err
	}
	if client == nil {
		return fmt.Errorf("nil IAM client")
	}

	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	config, err := getConfigFromStorage(ctx, req)
	if err != nil {
		return err
	}

	if config.AccessKey == "" || config.SecretKey == "" {
		return fmt.Errorf("cannot call config/rotate-root when either access_key or secret_key is empty")
	}

	var getUserInput iam.GetUserInput // empty input means get current user
	getUserRes, err := client.GetUserWithContext(ctx, &getUserInput)
	if err != nil {
		return fmt.Errorf("error calling GetUser: %w", err)
	}
	if getUserRes == nil {
		return fmt.Errorf("nil response from GetUser")
	}
	if getUserRes.User == nil {
		return fmt.Errorf("nil user returned from GetUser")
	}
	if getUserRes.User.UserName == nil {
		return fmt.Errorf("nil UserName returned from GetUser")
	}

	createAccessKeyInput := iam.CreateAccessKeyInput{
		UserName: getUserRes.User.UserName,
	}
	createAccessKeyRes, err := client.CreateAccessKeyWithContext(ctx, &createAccessKeyInput)
	if err != nil {
		return fmt.Errorf("error calling CreateAccessKey: %w", err)
	}
	if createAccessKeyRes.AccessKey == nil {
		return fmt.Errorf("nil response from CreateAccessKey")
	}
	if createAccessKeyRes.AccessKey.AccessKeyId == nil || createAccessKeyRes.AccessKey.SecretAccessKey == nil {
		return fmt.Errorf("nil AccessKeyId or SecretAccessKey returned from CreateAccessKey")
	}

	oldAccessKey := config.AccessKey

	config.AccessKey = *createAccessKeyRes.AccessKey.AccessKeyId
	config.SecretKey = *createAccessKeyRes.AccessKey.SecretAccessKey

	putConfigToStorage(ctx, req, config)

	b.iamClient = nil
	b.stsClient = nil

	deleteAccessKeyInput := iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(oldAccessKey),
		UserName:    getUserRes.User.UserName,
	}
	_, err = client.DeleteAccessKeyWithContext(ctx, &deleteAccessKeyInput)
	if err != nil {
		return fmt.Errorf("error deleting old access key: %w", err)
	}

	return nil
}

const pathConfigRotateRootHelpSyn = `
Request to rotate the AWS credentials used by Vault
`

const pathConfigRotateRootHelpDesc = `
This path attempts to rotate the AWS credentials used by Vault for this mount.
It is only valid if Vault has been configured to use AWS IAM credentials via the
config/root endpoint.
`
