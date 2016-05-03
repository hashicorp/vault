package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/awsutil"
	"github.com/hashicorp/vault/logical"
)

func getRootConfig(s logical.Storage) (*aws.Config, error) {
	credConfig := &awsutil.AWSCredentialsConfig{}

	entry, err := s.Get("config/root")
	if err != nil {
		return nil, err
	}
	if entry != nil {
		var config rootConfig
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, fmt.Errorf("error reading root configuration: %s", err)
		}

		credConfig.AccessKey = config.AccessKey
		credConfig.SecretKey = config.SecretKey
		credConfig.Region = config.Region
	}

	creds, err := awsutil.GenerateCredentialChain(credConfig)
	if err != nil {
		return nil, err
	}

	return &aws.Config{
		Credentials: creds,
		Region:      aws.String(credConfig.Region),
		HTTPClient:  cleanhttp.DefaultClient(),
	}, nil
}

func clientIAM(s logical.Storage) (*iam.IAM, error) {
	awsConfig, _ := getRootConfig(s)
	return iam.New(session.New(awsConfig)), nil
}

func clientSTS(s logical.Storage) (*sts.STS, error) {
	awsConfig, _ := getRootConfig(s)
	return sts.New(session.New(awsConfig)), nil
}
