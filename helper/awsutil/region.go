package awsutil

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	hclog "github.com/hashicorp/go-hclog"
)

const DefaultRegion = "us-east-1"

var (
	RegionEnvKeys = []string{"AWS_REGION", "AWS_DEFAULT_REGION"}

	ec2MetadataBaseURL = "http://169.254.169.254"
)

func GetOrDefaultRegion(logger hclog.Logger, configuredRegion string) string {
	if configuredRegion != "" {
		return configuredRegion
	}

	for _, envKey := range RegionEnvKeys {
		envVal := os.Getenv(envKey)
		if envVal != "" {
			return envVal
		}
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("unable to start session, defaulting region to %s", DefaultRegion))
		return DefaultRegion
	}

	region := aws.StringValue(sess.Config.Region)
	if region != "" {
		return region
	}

	metadata := ec2metadata.New(sess, &aws.Config{
		Endpoint:                          aws.String(ec2MetadataBaseURL + "/latest"),
		EC2MetadataDisableTimeoutOverride: aws.Bool(true),
		HTTPClient: &http.Client{
			Timeout: time.Second,
		},
	})
	if !metadata.Available() {
		return DefaultRegion
	}

	region, err = metadata.Region()
	if err != nil {
		logger.Warn("unable to retrieve region from instance metadata, defaulting region to %s", DefaultRegion)
		return DefaultRegion
	}
	return region
}
