package awsKms

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/logical"
)

func clientKMS(s logical.Storage) (*kms.KMS, error) {
	entry, err := s.Get("config/user")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf(
			"user credentials haven't been configured. Please configure\n" +
				"them at the 'config/user' endpoint")
	}

	var config rootConfig
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("error reading user configuration: %s", err)
	}

	creds := credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, "")
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"),
		HTTPClient:  cleanhttp.DefaultClient(),
	}

	return kms.New(session.New(awsConfig)), nil
}
