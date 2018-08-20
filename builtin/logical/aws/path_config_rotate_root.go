package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigRotateRoot() *framework.Path {
	return &framework.Path{
		Pattern: "config/rotate-root",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: pathConfigRotateRootUpdate,
		},

		HelpSynopsis:    pathConfigRotateRootHelpSyn,
		HelpDescription: pathConfigRotateRootHelpDesc,
	}
}

func pathConfigRotateRootUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// TODO: Add locking around reading/writing the config path/rootConfig
	rawRootConfig, err := req.Storage.Get(ctx, "config/root")
	if err != nil {
		return nil, err
	}
	if rawRootConfig == nil {
		return logical.ErrorResponse("no configuration found for config/root"), nil
	}
	var config rootConfig
	if err := rawRootConfig.DecodeJSON(&config); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error reading root configuration: %v", err)), nil
	}

	if config.AccessKey == "" || config.SecretKey == "" {
		return logical.ErrorResponse("cannot call config/rotate-root when either access_key or secret_key is empty"), nil
	}

	client, err := clientIAM(ctx, req.Storage)
	if err == nil {
		return logical.ErrorResponse(fmt.Sprintf("error retrieving IAM client: %v", err)), nil
	}
	if client == nil {
		return logical.ErrorResponse("nil IAM client"), nil
	}
	var getUserInput iam.GetUserInput // empty input means get current user
	getUserRes, err := client.GetUser(&getUserInput)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error calling GetUser: %v", err)), nil
	}
	if getUserRes == nil {
		return logical.ErrorResponse("nil response from GetUser"), nil
	}
	if getUserRes.User == nil {
		return logical.ErrorResponse("nil user returned from GetUser"), nil
	}
	if getUserRes.User.UserName == nil {
		return logical.ErrorResponse("nil UserName returnjd from GetUser"), nil
	}

	createAccessKeyInput := iam.CreateAccessKeyInput{
		UserName: getUserRes.User.UserName,
	}
	createAccessKeyRes, err := client.CreateAccessKey(&createAccessKeyInput)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error calling CreateAccessKey: %v", err)), nil
	}
	if createAccessKeyRes.AccessKey == nil {
		return logical.ErrorResponse("nil response from CreateAccessKey"), nil
	}
	if createAccessKeyRes.AccessKey.AccessKeyId == nil || createAccessKeyRes.AccessKey.SecretAccessKey == nil {
		return logical.ErrorResponse("nil AccessKeyId or SecretAccessKey returned from CreateAccessKey"), nil
	}

	oldAccessKey := config.AccessKey

	config.AccessKey = *createAccessKeyRes.AccessKey.AccessKeyId
	config.SecretKey = *createAccessKeyRes.AccessKey.SecretAccessKey

	newEntry, err := logical.StorageEntryJSON("config/root", config)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error generating new root config: %v", err)), nil
	}
	// TODO: Should we return the full error message? Are there any scenarios in which it could expose
	// the underlying creds? (The idea of rotate-root is that it should never expose credentials.)
	if err := req.Storage.Put(ctx, newEntry); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error saving new root config: %v", err)), nil
	}

	deleteAccessKeyInput := iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(oldAccessKey),
		UserName:    getUserRes.User.UserName,
	}
	_, err = client.DeleteAccessKey(&deleteAccessKeyInput)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error deleting old access key: %v", err)), nil
	}

	return nil, nil
}

const pathConfigRotateRootHelpSyn = `
Request to rotate the AWS credentials used by Vault
`

const pathConfigRotateRootHelpDesc = `
This path attempts to rotate the AWS credentials used by Vault for this mount.
It is only valid if Vault has been configured to use AWS IAM credentials via the
config/root endpoint.
`
