package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/vault/logical"
)

func clientIAM(s logical.Storage) (*iam.IAM, error) {
	entry, err := s.Get("config/root")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf(
			"root credentials haven't been configured. Please configure\n" +
				"them at the 'config/root' endpoint")
	}

	var config rootConfig
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %s", err)
	}

	creds := credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, "")
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String(config.Region),
	}
	return iam.New(awsConfig), nil
}
