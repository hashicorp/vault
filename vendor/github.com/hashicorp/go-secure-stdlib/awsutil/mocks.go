// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsutil

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// MockOptionErr provides a mock option error for use with testing.
func MockOptionErr(withErr error) Option {
	return func(_ *options) error {
		return withErr
	}
}

// MockIAM provides a way to mock the AWS IAM API.
type MockIAM struct {
	iamiface.IAMAPI

	CreateAccessKeyOutput *iam.CreateAccessKeyOutput
	CreateAccessKeyError  error
	DeleteAccessKeyError  error
	ListAccessKeysOutput  *iam.ListAccessKeysOutput
	ListAccessKeysError   error
	GetUserOutput         *iam.GetUserOutput
	GetUserError          error
}

// MockIAMOption is a function for setting the various fields on a MockIAM
// object.
type MockIAMOption func(m *MockIAM) error

// WithCreateAccessKeyOutput sets the output for the CreateAccessKey method.
func WithCreateAccessKeyOutput(o *iam.CreateAccessKeyOutput) MockIAMOption {
	return func(m *MockIAM) error {
		m.CreateAccessKeyOutput = o
		return nil
	}
}

// WithCreateAccessKeyError sets the error output for the CreateAccessKey
// method.
func WithCreateAccessKeyError(e error) MockIAMOption {
	return func(m *MockIAM) error {
		m.CreateAccessKeyError = e
		return nil
	}
}

// WithDeleteAccessKeyError sets the error output for the DeleteAccessKey
// method.
func WithDeleteAccessKeyError(e error) MockIAMOption {
	return func(m *MockIAM) error {
		m.DeleteAccessKeyError = e
		return nil
	}
}

// WithListAccessKeysOutput sets the output for the ListAccessKeys method.
func WithListAccessKeysOutput(o *iam.ListAccessKeysOutput) MockIAMOption {
	return func(m *MockIAM) error {
		m.ListAccessKeysOutput = o
		return nil
	}
}

// WithListAccessKeysError sets the error output for the ListAccessKeys method.
func WithListAccessKeysError(e error) MockIAMOption {
	return func(m *MockIAM) error {
		m.ListAccessKeysError = e
		return nil
	}
}

// WithGetUserOutput sets the output for the GetUser method.
func WithGetUserOutput(o *iam.GetUserOutput) MockIAMOption {
	return func(m *MockIAM) error {
		m.GetUserOutput = o
		return nil
	}
}

// WithGetUserError sets the error output for the GetUser method.
func WithGetUserError(e error) MockIAMOption {
	return func(m *MockIAM) error {
		m.GetUserError = e
		return nil
	}
}

// NewMockIAM provides a factory function to use with the WithIAMAPIFunc
// option.
func NewMockIAM(opts ...MockIAMOption) IAMAPIFunc {
	return func(_ *session.Session) (iamiface.IAMAPI, error) {
		m := new(MockIAM)
		for _, opt := range opts {
			if err := opt(m); err != nil {
				return nil, err
			}
		}

		return m, nil
	}
}

func (m *MockIAM) CreateAccessKey(*iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error) {
	if m.CreateAccessKeyError != nil {
		return nil, m.CreateAccessKeyError
	}

	return m.CreateAccessKeyOutput, nil
}

func (m *MockIAM) DeleteAccessKey(*iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
	return &iam.DeleteAccessKeyOutput{}, m.DeleteAccessKeyError
}

func (m *MockIAM) ListAccessKeys(*iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
	if m.ListAccessKeysError != nil {
		return nil, m.ListAccessKeysError
	}

	return m.ListAccessKeysOutput, nil
}

func (m *MockIAM) GetUser(*iam.GetUserInput) (*iam.GetUserOutput, error) {
	if m.GetUserError != nil {
		return nil, m.GetUserError
	}

	return m.GetUserOutput, nil
}

// MockSTS provides a way to mock the AWS STS API.
type MockSTS struct {
	stsiface.STSAPI

	GetCallerIdentityOutput *sts.GetCallerIdentityOutput
	GetCallerIdentityError  error
}

// MockSTSOption is a function for setting the various fields on a MockSTS
// object.
type MockSTSOption func(m *MockSTS) error

// WithGetCallerIdentityOutput sets the output for the GetCallerIdentity
// method.
func WithGetCallerIdentityOutput(o *sts.GetCallerIdentityOutput) MockSTSOption {
	return func(m *MockSTS) error {
		m.GetCallerIdentityOutput = o
		return nil
	}
}

// WithGetCallerIdentityError sets the error output for the GetCallerIdentity
// method.
func WithGetCallerIdentityError(e error) MockSTSOption {
	return func(m *MockSTS) error {
		m.GetCallerIdentityError = e
		return nil
	}
}

// NewMockSTS provides a factory function to use with the WithSTSAPIFunc
// option.
//
// If withGetCallerIdentityError is supplied, calls to GetCallerIdentity will
// return the supplied error. Otherwise, a basic mock API output is returned.
func NewMockSTS(opts ...MockSTSOption) STSAPIFunc {
	return func(_ *session.Session) (stsiface.STSAPI, error) {
		m := new(MockSTS)
		for _, opt := range opts {
			if err := opt(m); err != nil {
				return nil, err
			}
		}

		return m, nil
	}
}

func (m *MockSTS) GetCallerIdentity(_ *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	if m.GetCallerIdentityError != nil {
		return nil, m.GetCallerIdentityError
	}

	return m.GetCallerIdentityOutput, nil
}
