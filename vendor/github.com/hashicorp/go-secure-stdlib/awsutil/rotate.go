// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

// RotateKeys takes the access key and secret key from this credentials config
// and first creates a new access/secret key, then deletes the old access key.
// If deletion of the old access key is successful, the new access key/secret
// key are written into the credentials config and nil is returned. On any
// error, the old credentials are not overwritten. This ensures that any
// generated new secret key never leaves this function in case of an error, even
// though it will still result in an extraneous access key existing; we do also
// try to delete the new one to clean up, although it's unlikely that will work
// if the old one could not be deleted.
//
// Supported options: WithEnvironmentCredentials, WithSharedCredentials,
// WithAwsSession, WithUsername, WithValidityCheckTimeout, WithIAMAPIFunc,
// WithSTSAPIFunc
//
// Note that WithValidityCheckTimeout here, when non-zero, controls the
// WithValidityCheckTimeout option on access key creation. See CreateAccessKey
// for more details.
func (c *CredentialsConfig) RotateKeys(opt ...Option) error {
	if c.AccessKey == "" || c.SecretKey == "" {
		return errors.New("cannot rotate credentials when either access_key or secret_key is empty")
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return fmt.Errorf("error reading options in RotateKeys: %w", err)
	}

	sess := opts.withAwsSession
	if sess == nil {
		sess, err = c.GetSession(opt...)
		if err != nil {
			return fmt.Errorf("error calling GetSession: %w", err)
		}
	}

	sessOpt := append(opt, WithAwsSession(sess))
	createAccessKeyRes, err := c.CreateAccessKey(sessOpt...)
	if err != nil {
		return fmt.Errorf("error calling CreateAccessKey: %w", err)
	}

	err = c.DeleteAccessKey(c.AccessKey, append(sessOpt, WithUsername(*createAccessKeyRes.AccessKey.UserName))...)
	if err != nil {
		return fmt.Errorf("error deleting old access key: %w", err)
	}

	c.AccessKey = *createAccessKeyRes.AccessKey.AccessKeyId
	c.SecretKey = *createAccessKeyRes.AccessKey.SecretAccessKey

	return nil
}

// CreateAccessKey creates a new access/secret key pair.
//
// Supported options: WithEnvironmentCredentials, WithSharedCredentials,
// WithAwsSession, WithUsername, WithValidityCheckTimeout, WithIAMAPIFunc,
// WithSTSAPIFunc
//
// When WithValidityCheckTimeout is non-zero, it specifies a timeout to wait on
// the created credentials to be valid and ready for use.
func (c *CredentialsConfig) CreateAccessKey(opt ...Option) (*iam.CreateAccessKeyOutput, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("error reading options in CreateAccessKey: %w", err)
	}

	client, err := c.IAMClient(opt...)
	if err != nil {
		return nil, fmt.Errorf("error loading IAM client: %w", err)
	}

	var getUserInput iam.GetUserInput
	if opts.withUsername != "" {
		getUserInput.SetUserName(opts.withUsername)
	} // otherwise, empty input means get current user
	getUserRes, err := client.GetUser(&getUserInput)
	if err != nil {
		return nil, fmt.Errorf("error calling aws.GetUser: %w", err)
	}
	if getUserRes == nil {
		return nil, fmt.Errorf("nil response from aws.GetUser")
	}
	if getUserRes.User == nil {
		return nil, fmt.Errorf("nil user returned from aws.GetUser")
	}
	if getUserRes.User.UserName == nil {
		return nil, fmt.Errorf("nil UserName returned from aws.GetUser")
	}

	createAccessKeyInput := iam.CreateAccessKeyInput{
		UserName: getUserRes.User.UserName,
	}
	createAccessKeyRes, err := client.CreateAccessKey(&createAccessKeyInput)
	if err != nil {
		return nil, fmt.Errorf("error calling aws.CreateAccessKey: %w", err)
	}
	if createAccessKeyRes == nil {
		return nil, fmt.Errorf("nil response from aws.CreateAccessKey")
	}
	if createAccessKeyRes.AccessKey == nil {
		return nil, fmt.Errorf("nil access key in response from aws.CreateAccessKey")
	}
	if createAccessKeyRes.AccessKey.AccessKeyId == nil || createAccessKeyRes.AccessKey.SecretAccessKey == nil {
		return nil, fmt.Errorf("nil AccessKeyId or SecretAccessKey returned from aws.CreateAccessKey")
	}

	// Check the credentials to make sure they are usable. We only do
	// this if withValidityCheckTimeout is non-zero to ensue that we don't
	// immediately fail due to eventual consistency.
	if opts.withValidityCheckTimeout != 0 {
		newC := &CredentialsConfig{
			AccessKey: *createAccessKeyRes.AccessKey.AccessKeyId,
			SecretKey: *createAccessKeyRes.AccessKey.SecretAccessKey,
		}

		if _, err := newC.GetCallerIdentity(
			WithValidityCheckTimeout(opts.withValidityCheckTimeout),
			WithSTSAPIFunc(opts.withSTSAPIFunc),
		); err != nil {
			return nil, fmt.Errorf("error verifying new credentials: %w", err)
		}
	}

	return createAccessKeyRes, nil
}

// DeleteAccessKey deletes an access key.
//
// Supported options: WithEnvironmentCredentials, WithSharedCredentials,
// WithAwsSession, WithUserName, WithIAMAPIFunc
func (c *CredentialsConfig) DeleteAccessKey(accessKeyId string, opt ...Option) error {
	opts, err := getOpts(opt...)
	if err != nil {
		return fmt.Errorf("error reading options in RotateKeys: %w", err)
	}

	client, err := c.IAMClient(opt...)
	if err != nil {
		return fmt.Errorf("error loading IAM client: %w", err)
	}

	deleteAccessKeyInput := iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(accessKeyId),
	}
	if opts.withUsername != "" {
		deleteAccessKeyInput.SetUserName(opts.withUsername)
	}

	_, err = client.DeleteAccessKey(&deleteAccessKeyInput)
	if err != nil {
		return fmt.Errorf("error deleting old access key: %w", err)
	}

	return nil
}

// GetSession returns an AWS session configured according to the various values
// in the CredentialsConfig object. This can be passed into iam.New or sts.New
// as appropriate.
//
// Supported options: WithEnvironmentCredentials, WithSharedCredentials,
// WithAwsSession, WithClientType
func (c *CredentialsConfig) GetSession(opt ...Option) (*session.Session, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("error reading options in GetSession: %w", err)
	}

	creds, err := c.GenerateCredentialChain(opt...)
	if err != nil {
		return nil, err
	}

	var endpoint string
	switch opts.withClientType {
	case "sts":
		endpoint = c.STSEndpoint
	case "iam":
		endpoint = c.IAMEndpoint
	default:
		return nil, fmt.Errorf("unknown client type %q in GetSession", opts.withClientType)
	}

	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String(c.Region),
		Endpoint:    aws.String(endpoint),
		HTTPClient:  c.HTTPClient,
		MaxRetries:  c.MaxRetries,
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("error getting new session: %w", err)
	}

	return sess, nil
}

// GetCallerIdentity runs sts.GetCallerIdentity for the current set
// credentials. This can be used to check that credentials are valid,
// in addition to checking details about the effective logged in
// account and user ID.
//
// Supported options: WithEnvironmentCredentials,
// WithSharedCredentials, WithAwsSession, WithValidityCheckTimeout
func (c *CredentialsConfig) GetCallerIdentity(opt ...Option) (*sts.GetCallerIdentityOutput, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("error reading options in GetCallerIdentity: %w", err)
	}

	client, err := c.STSClient(opt...)
	if err != nil {
		return nil, fmt.Errorf("error loading STS client: %w", err)
	}

	delay := time.Second
	timeoutCtx, cancel := context.WithTimeout(context.Background(), opts.withValidityCheckTimeout)
	defer cancel()
	for {
		cid, err := client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
		if err == nil {
			return cid, nil
		}

		// TODO: can add a context here for external cancellation in the future
		select {
		case <-time.After(delay):
			// pass

		case <-timeoutCtx.Done():
			// Format our error based on how we were called.
			if opts.withValidityCheckTimeout == 0 {
				// There was no timeout, just return the error unwrapped.
				return nil, err
			}

			// Otherwise, return the error wrapped in a timeout error.
			return nil, fmt.Errorf("timeout after %s waiting for success: %w", opts.withValidityCheckTimeout, err)
		}
	}
}
