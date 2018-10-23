package awskms

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
)

const awsTestKeyID = "foo"

func NewAWSKMSTestSeal() *AWSKMSSeal {
	s := NewSeal(logging.NewVaultLogger(log.Trace))
	s.client = &mockAWSKMSSealClient{
		keyID: aws.String(awsTestKeyID),
	}
	return s
}

type mockAWSKMSSealClient struct {
	kmsiface.KMSAPI
	keyID *string
}

// Encrypt is a mocked call that returns a base64 encoded string.
func (m *mockAWSKMSSealClient) Encrypt(input *kms.EncryptInput) (*kms.EncryptOutput, error) {
	m.keyID = input.KeyId

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(input.Plaintext)))
	base64.StdEncoding.Encode(encoded, input.Plaintext)

	return &kms.EncryptOutput{
		CiphertextBlob: encoded,
		KeyId:          input.KeyId,
	}, nil
}

// Decrypt is a mocked call that returns a decoded base64 string.
func (m *mockAWSKMSSealClient) Decrypt(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
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
		KeyId:     m.keyID,
		Plaintext: decoded,
	}, nil
}

// DescribeKey is a mocked call that returns the keyID.
func (m *mockAWSKMSSealClient) DescribeKey(input *kms.DescribeKeyInput) (*kms.DescribeKeyOutput, error) {
	if m.keyID == nil {
		return nil, awserr.New(kms.ErrCodeNotFoundException, "key not found", nil)
	}

	return &kms.DescribeKeyOutput{
		KeyMetadata: &kms.KeyMetadata{
			KeyId: m.keyID,
		},
	}, nil
}
