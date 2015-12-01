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
	awsConfig := &aws.Config{
		Region:     aws.String("us-east-1"),
		HTTPClient: cleanhttp.DefaultClient(),
	}
	entry, err := s.Get("config/user")
	if err != nil {
		return nil, err
	}

	// If the entry is not found, AWS SDK will try to find the credentials
	// from other locations, so it's fine if user hasn't configured this.
	if entry != nil {
		var config rootConfig
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, fmt.Errorf("error reading user configuration: %s", err)
		}
		creds := credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, "")
		awsConfig.Credentials = creds
	}

	return kms.New(session.New(awsConfig)), nil
}
