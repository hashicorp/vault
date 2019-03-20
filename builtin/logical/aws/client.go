package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/awsutil"
	"github.com/hashicorp/vault/logical"
)

// NOTE: The caller is required to ensure that b.clientMutex is at least read locked
func getRootConfig(ctx context.Context, s logical.Storage, clientType string) (*aws.Config, error) {
	credsConfig := &awsutil.CredentialsConfig{}
	var endpoint string
	var maxRetries int = aws.UseServiceDefaultRetries

	entry, err := s.Get(ctx, "config/root")
	if err != nil {
		return nil, err
	}
	if entry != nil {
		var config rootConfig
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, errwrap.Wrapf("error reading root configuration: {{err}}", err)
		}

		credsConfig.AccessKey = config.AccessKey
		credsConfig.SecretKey = config.SecretKey
		credsConfig.Region = config.Region
		maxRetries = config.MaxRetries
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
		MaxRetries:  aws.Int(maxRetries),
	}, nil
}

func nonCachedClientIAM(ctx context.Context, s logical.Storage) (*iam.IAM, error) {
	awsConfig, err := getRootConfig(ctx, s, "iam")
	if err != nil {
		return nil, err
	}

	client := iam.New(session.New(awsConfig))

	if client == nil {
		return nil, fmt.Errorf("could not obtain iam client")
	}
	return client, nil
}

func nonCachedClientSTS(ctx context.Context, s logical.Storage) (*sts.STS, error) {
	awsConfig, err := getRootConfig(ctx, s, "sts")
	if err != nil {
		return nil, err
	}
	client := sts.New(session.New(awsConfig))

	if client == nil {
		return nil, fmt.Errorf("could not obtain sts client")
	}
	return client, nil
}
