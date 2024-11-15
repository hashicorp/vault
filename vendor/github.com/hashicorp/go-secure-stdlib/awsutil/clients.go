// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsutil

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// IAMAPIFunc is a factory function for returning an IAM interface,
// useful for supplying mock interfaces for testing IAM. The session
// is passed into the function in the same way as done with the
// standard iam.New() constructor.
type IAMAPIFunc func(sess *session.Session) (iamiface.IAMAPI, error)

// STSAPIFunc is a factory function for returning a STS interface,
// useful for supplying mock interfaces for testing STS. The session
// is passed into the function in the same way as done with the
// standard sts.New() constructor.
type STSAPIFunc func(sess *session.Session) (stsiface.STSAPI, error)

// IAMClient returns an IAM client.
//
// Supported options: WithSession, WithIAMAPIFunc.
//
// If WithIAMAPIFunc is supplied, the included function is used as
// the IAM client constructor instead. This can be used for Mocking
// the IAM API.
func (c *CredentialsConfig) IAMClient(opt ...Option) (iamiface.IAMAPI, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("error reading options: %w", err)
	}

	sess := opts.withAwsSession
	if sess == nil {
		sess, err = c.GetSession(opt...)
		if err != nil {
			return nil, fmt.Errorf("error calling GetSession: %w", err)
		}
	}

	if opts.withIAMAPIFunc != nil {
		return opts.withIAMAPIFunc(sess)
	}

	client := iam.New(sess)
	if client == nil {
		return nil, errors.New("could not obtain iam client from session")
	}

	return client, nil
}

// STSClient returns a STS client.
//
// Supported options: WithSession, WithSTSAPIFunc.
//
// If WithSTSAPIFunc is supplied, the included function is used as
// the STS client constructor instead. This can be used for Mocking
// the STS API.
func (c *CredentialsConfig) STSClient(opt ...Option) (stsiface.STSAPI, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("error reading options: %w", err)
	}

	sess := opts.withAwsSession
	if sess == nil {
		sess, err = c.GetSession(opt...)
		if err != nil {
			return nil, fmt.Errorf("error calling GetSession: %w", err)
		}
	}

	if opts.withSTSAPIFunc != nil {
		return opts.withSTSAPIFunc(sess)
	}

	client := sts.New(sess)
	if client == nil {
		return nil, errors.New("could not obtain sts client from session")
	}

	return client, nil
}
