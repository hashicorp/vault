// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathConfigRotateRoot() *framework.Path {
	return &framework.Path{
		Pattern: "config/rotate-root",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
			OperationVerb:   "rotate",
			OperationSuffix: "root-credentials",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigRotateRootUpdate,
			},
		},

		HelpSynopsis:    pathConfigRotateRootHelpSyn,
		HelpDescription: pathConfigRotateRootHelpDesc,
	}
}

func (b *backend) pathConfigRotateRootUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.rotateRoot(ctx, req)
}

func (b *backend) rotateRoot(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	// First get the AWS key and secret and validate that we _can_ rotate them.
	// We need the read lock here to prevent anything else from mutating it while we're using it.
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	clientConf, err := b.nonLockedClientConfigEntry(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if clientConf == nil {
		return logical.ErrorResponse(`can't update client config because it's unset`), nil
	}
	if clientConf.AccessKey == "" {
		return logical.ErrorResponse("can't update access_key because it's unset"), nil
	}
	if clientConf.SecretKey == "" {
		return logical.ErrorResponse("can't update secret_key because it's unset"), nil
	}

	// Getting our client through the b.clientIAM method requires values retrieved through
	// the user providing an ARN, which we don't have here, so let's just directly
	// make what we need.
	staticCreds := &credentials.StaticProvider{
		Value: credentials.Value{
			AccessKeyID:     clientConf.AccessKey,
			SecretAccessKey: clientConf.SecretKey,
		},
	}
	// By default, leave the iamEndpoint nil to tell AWS it's unset. However, if it is
	// configured, populate the pointer.
	var iamEndpoint *string
	if clientConf.IAMEndpoint != "" {
		iamEndpoint = aws.String(clientConf.IAMEndpoint)
	}

	// Attempt to retrieve the region, error out if no region is provided.
	region, err := awsutil.GetRegion("")
	if err != nil {
		return nil, fmt.Errorf("error retrieving region: %w", err)
	}

	awsConfig := &aws.Config{
		Credentials: credentials.NewCredentials(staticCreds),
		Endpoint:    iamEndpoint,

		// Generally speaking, GetRegion will use the Vault server's region. However, if this
		// needs to be overridden, an easy way would be to set the AWS_DEFAULT_REGION on the Vault server
		// to the desired region. If that's still insufficient for someone's use case, in the future we
		// could add the ability to specify the region either on the client config or as part of the
		// inbound rotation call.
		Region: aws.String(region),

		// Prevents races.
		HTTPClient: cleanhttp.DefaultClient(),
	}
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	iamClient := getIAMClient(sess)

	// Get the current user's name since it's required to create an access key.
	// Empty input means get the current user.
	var getUserInput iam.GetUserInput
	getUserRes, err := iamClient.GetUserWithContext(ctx, &getUserInput)
	if err != nil {
		return nil, fmt.Errorf("error calling GetUser: %w", err)
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

	// Create the new access key and secret.
	createAccessKeyInput := iam.CreateAccessKeyInput{
		UserName: getUserRes.User.UserName,
	}
	createAccessKeyRes, err := iamClient.CreateAccessKeyWithContext(ctx, &createAccessKeyInput)
	if err != nil {
		return nil, fmt.Errorf("error calling CreateAccessKey: %w", err)
	}
	if createAccessKeyRes.AccessKey == nil {
		return nil, fmt.Errorf("nil response from CreateAccessKey")
	}
	if createAccessKeyRes.AccessKey.AccessKeyId == nil || createAccessKeyRes.AccessKey.SecretAccessKey == nil {
		return nil, fmt.Errorf("nil AccessKeyId or SecretAccessKey returned from CreateAccessKey")
	}

	// We're about to attempt to store the newly created key and secret, but just in case we can't,
	// let's clean up after ourselves.
	storedNewConf := false
	var errs error
	defer func() {
		if storedNewConf {
			return
		}
		// Attempt to delete the access key and secret we created but couldn't store and use.
		deleteAccessKeyInput := iam.DeleteAccessKeyInput{
			AccessKeyId: createAccessKeyRes.AccessKey.AccessKeyId,
			UserName:    getUserRes.User.UserName,
		}
		if _, err := iamClient.DeleteAccessKeyWithContext(ctx, &deleteAccessKeyInput); err != nil {
			// Include this error in the errs returned by this method.
			errs = multierror.Append(errs, fmt.Errorf("error deleting newly created but unstored access key ID %s: %s", *createAccessKeyRes.AccessKey.AccessKeyId, err))
		}
	}()

	oldAccessKey := clientConf.AccessKey
	clientConf.AccessKey = *createAccessKeyRes.AccessKey.AccessKeyId
	clientConf.SecretKey = *createAccessKeyRes.AccessKey.SecretAccessKey

	// Now get ready to update storage, doing everything beforehand so we can minimize how long
	// we need to hold onto the lock.
	newEntry, err := b.configClientToEntry(clientConf)
	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error generating new client config JSON: %w", err))
		return nil, errs
	}

	// Someday we may want to allow the user to send a number of seconds to wait here
	// before deleting the previous access key to allow work to complete. That would allow
	// AWS, which is eventually consistent, to finish populating the new key in all places.
	if err := req.Storage.Put(ctx, newEntry); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error saving new client config: %w", err))
		return nil, errs
	}
	storedNewConf = true

	// Previous cached clients need to be cleared because they may have been made using
	// the soon-to-be-obsolete credentials.
	b.IAMClientsMap = make(map[string]map[string]*iam.IAM)
	b.EC2ClientsMap = make(map[string]map[string]*ec2.EC2)

	// Now to clean up the old key.
	deleteAccessKeyInput := iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(oldAccessKey),
		UserName:    getUserRes.User.UserName,
	}
	if _, err = iamClient.DeleteAccessKeyWithContext(ctx, &deleteAccessKeyInput); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error deleting old access key ID %s: %w", oldAccessKey, err))
		return nil, errs
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"access_key": clientConf.AccessKey,
		},
	}, nil
}

// getIAMClient allows us to change how an IAM client is created
// during testing. The AWS SDK doesn't easily lend itself to testing
// using a Go httptest server because if you inject a test URL into
// the config, the client strips important information about which
// endpoint it's hitting. Per
// https://aws.amazon.com/blogs/developer/mocking-out-then-aws-sdk-for-go-for-unit-testing/,
// this is the recommended approach.
var getIAMClient = func(sess *session.Session) iamiface.IAMAPI {
	return iam.New(sess)
}

const pathConfigRotateRootHelpSyn = `
Request to rotate the AWS credentials used by Vault
`

const pathConfigRotateRootHelpDesc = `
This path attempts to rotate the AWS credentials used by Vault for this mount.
It is only valid if Vault has been configured to use AWS IAM credentials via the
config/client endpoint.
`
