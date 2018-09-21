package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigRotateRoot(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/rotate-root",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigRotateRootUpdate,
		},

		HelpSynopsis:    pathConfigRotateRootHelpSyn,
		HelpDescription: pathConfigRotateRootHelpDesc,
	}
}

func (b *backend) pathConfigRotateRootUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// have to get the client config first because that takes out a read lock
	client, err := b.clientIAM(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("nil IAM client")
	}

	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	rawRootConfig, err := req.Storage.Get(ctx, "config/root")
	if err != nil {
		return nil, err
	}
	if rawRootConfig == nil {
		return nil, fmt.Errorf("no configuration found for config/root")
	}
	var config rootConfig
	if err := rawRootConfig.DecodeJSON(&config); err != nil {
		return nil, errwrap.Wrapf("error reading root configuration: {{err}}", err)
	}

	if config.AccessKey == "" || config.SecretKey == "" {
		return logical.ErrorResponse("Cannot call config/rotate-root when either access_key or secret_key is empty"), nil
	}

	var getUserInput iam.GetUserInput // empty input means get current user
	getUserRes, err := client.GetUser(&getUserInput)
	if err != nil {
		return nil, errwrap.Wrapf("error calling GetUser: {{err}}", err)
	}
	if getUserRes == nil {
		return nil, fmt.Errorf("nil response from GetUser")
	}
	if getUserRes.User == nil {
		return nil, fmt.Errorf("nil user returned from GetUser")
	}
	if getUserRes.User.UserName == nil {
		return nil, fmt.Errorf("nil UserName returned from GetUser")
	}

	createAccessKeyInput := iam.CreateAccessKeyInput{
		UserName: getUserRes.User.UserName,
	}
	createAccessKeyRes, err := client.CreateAccessKey(&createAccessKeyInput)
	if err != nil {
		return nil, errwrap.Wrapf("error calling CreateAccessKey: {{err}}", err)
	}
	if createAccessKeyRes.AccessKey == nil {
		return nil, fmt.Errorf("nil response from CreateAccessKey")
	}
	if createAccessKeyRes.AccessKey.AccessKeyId == nil || createAccessKeyRes.AccessKey.SecretAccessKey == nil {
		return nil, fmt.Errorf("nil AccessKeyId or SecretAccessKey returned from CreateAccessKey")
	}

	oldAccessKey := config.AccessKey

	config.AccessKey = *createAccessKeyRes.AccessKey.AccessKeyId
	config.SecretKey = *createAccessKeyRes.AccessKey.SecretAccessKey

	newEntry, err := logical.StorageEntryJSON("config/root", config)
	if err != nil {
		return nil, errwrap.Wrapf("error generating new config/root JSON: {{err}}", err)
	}
	if err := req.Storage.Put(ctx, newEntry); err != nil {
		return nil, errwrap.Wrapf("error saving new config/root: {{err}}", err)
	}

	b.iamClient = nil
	b.stsClient = nil

	deleteAccessKeyInput := iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(oldAccessKey),
		UserName:    getUserRes.User.UserName,
	}
	_, err = client.DeleteAccessKey(&deleteAccessKeyInput)
	if err != nil {
		return nil, errwrap.Wrapf("error deleting old access key: {{err}}", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"access_key": config.AccessKey,
		},
	}, nil
}

const pathConfigRotateRootHelpSyn = `
Request to rotate the AWS credentials used by Vault
`

const pathConfigRotateRootHelpDesc = `
This path attempts to rotate the AWS credentials used by Vault for this mount.
It is only valid if Vault has been configured to use AWS IAM credentials via the
config/root endpoint.
`
