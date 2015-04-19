package aws

import (
	"fmt"

	"github.com/hashicorp/aws-sdk-go/aws"
	"github.com/hashicorp/aws-sdk-go/gen/iam"
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

	creds := aws.Creds(config.AccessKey, config.SecretKey, "")
	return iam.New(creds, config.Region, nil), nil
}
