package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/awsutil"
	"github.com/hashicorp/vault/logical"
)

func getRootConfig(s logical.Storage, clientType string) (*aws.Config, error) {
	credsConfig := &awsutil.CredentialsConfig{}
	var endpoint string

	entry, err := s.Get("config/root")
	if err != nil {
		return nil, err
	}
	if entry != nil {
		var config rootConfig
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, fmt.Errorf("error reading root configuration: %s", err)
		}

		credsConfig.AccessKey = config.AccessKey
		credsConfig.SecretKey = config.SecretKey
		credsConfig.Region = config.Region
		switch {
		case clientType == "iam" && config.IAMEndpoint != "":
			endpoint = *aws.String(config.IAMEndpoint)
		case clientType == "sts" && config.STSEndpoint != "":
			endpoint = *aws.String(config.STSEndpoint)
		}
	}

	if credsConfig.Region == "" {
		credsConfig.Region = os.Getenv("AWS_REGION")
		if credsConfig.Region == "" {
			credsConfig.Region = os.Getenv("AWS_DEFAULT_REGION")
			if credsConfig.Region == "" {
				credsConfig.Region = "us-east-1"
			}
		}
	}

	credsConfig.HTTPClient = cleanhttp.DefaultClient()

	creds, err := credsConfig.GenerateCredentialChain()
	if err != nil {
		return nil, err
	}

	return &aws.Config{
		Credentials: creds,
		Region:      aws.String(credsConfig.Region),
		Endpoint:    &endpoint,
		HTTPClient:  cleanhttp.DefaultClient(),
	}, nil
}

func clientIAM(s logical.Storage) (*iam.IAM, error) {
	awsConfig, err := getRootConfig(s, "iam")
	if err != nil {
		return nil, err
	}

	client := iam.New(session.New(awsConfig))

	if client == nil {
		return nil, fmt.Errorf("could not obtain iam client")
	}
	return client, nil
}

func clientSTS(s logical.Storage) (*sts.STS, error) {
	awsConfig, err := getRootConfig(s, "sts")
	if err != nil {
		return nil, err
	}
	client := sts.New(session.New(awsConfig))

	if client == nil {
		return nil, fmt.Errorf("could not obtain sts client")
	}
	return client, nil
}
