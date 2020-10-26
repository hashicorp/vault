package awsutil

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

type MockIAM struct {
	iamiface.IAMAPI

	CreateAccessKeyOutput *iam.CreateAccessKeyOutput
	DeleteAccessKeyOutput *iam.DeleteAccessKeyOutput
	GetUserOutput         *iam.GetUserOutput
}

func (m *MockIAM) CreateAccessKey(*iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error) {
	return m.CreateAccessKeyOutput, nil
}

func (m *MockIAM) DeleteAccessKey(*iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
	return m.DeleteAccessKeyOutput, nil
}

func (m *MockIAM) GetUser(*iam.GetUserInput) (*iam.GetUserOutput, error) {
	return m.GetUserOutput, nil
}
