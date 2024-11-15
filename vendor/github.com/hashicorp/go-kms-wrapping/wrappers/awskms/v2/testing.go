// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awskms

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

const awsTestKeyId = "foo"

func NewAwsKmsTestWrapper() *Wrapper {
	s := NewWrapper()
	s.client = &mockClient{
		keyId: aws.String(awsTestKeyId),
	}
	return s
}

type mockClient struct {
	kmsiface.KMSAPI
	keyId *string
}

// Encrypt is a mocked call that returns a base64 encoded string.
func (m *mockClient) Encrypt(input *kms.EncryptInput) (*kms.EncryptOutput, error) {
	m.keyId = input.KeyId

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(input.Plaintext)))
	base64.StdEncoding.Encode(encoded, input.Plaintext)

	return &kms.EncryptOutput{
		CiphertextBlob: encoded,
		KeyId:          input.KeyId,
	}, nil
}

// Decrypt is a mocked call that returns a decoded base64 string.
func (m *mockClient) Decrypt(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	decLen := base64.StdEncoding.DecodedLen(len(input.CiphertextBlob))
	decoded := make([]byte, decLen)
	len, err := base64.StdEncoding.Decode(decoded, input.CiphertextBlob)
	if err != nil {
		return nil, err
	}

	if len < decLen {
		decoded = decoded[:len]
	}

	return &kms.DecryptOutput{
		KeyId:     m.keyId,
		Plaintext: decoded,
	}, nil
}

// DescribeKey is a mocked call that returns the keyId.
func (m *mockClient) DescribeKey(input *kms.DescribeKeyInput) (*kms.DescribeKeyOutput, error) {
	if m.keyId == nil {
		return nil, awserr.New(kms.ErrCodeNotFoundException, "key not found", nil)
	}

	return &kms.DescribeKeyOutput{
		KeyMetadata: &kms.KeyMetadata{
			KeyId: m.keyId,
		},
	}, nil
}
