package awsutil

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	hclog "github.com/hashicorp/go-hclog"
)

func GetOrDefaultRegion(logger hclog.Logger, configuredRegion string) string {
	// We default to us-east-1 because it's a widely used region
	// and is also where AWS first rolls out new features.
	defaultRegion := "us-east-1"

	// We prefer env variables to configured ones because they
	// serve as a way to change the application's configuration
	// on the fly rather than via restarting some process.
	if os.Getenv("AWS_REGION") != "" {
		return os.Getenv("AWS_REGION")
	}
	if os.Getenv("AWS_DEFAULT_REGION") != "" {
		return os.Getenv("AWS_DEFAULT_REGION")
	}

	// If a region was configured, it's time to use it.
	if configuredRegion != "" {
		return configuredRegion
	}

	// Nothing was configured, let's try to get the region from EC2 instance metadata.
	sess, err := session.NewSession(nil)
	if err != nil {
		logger.Warn(fmt.Sprintf("unable to start session, defaulting region to %s", defaultRegion))
		return defaultRegion
	}

	metadata := ec2metadata.New(sess)
	if !metadata.Available() {
		return defaultRegion
	}

	region, err := metadata.Region()
	if err != nil {
		logger.Warn("unable to retrieve region from instance metadata, defaulting region to %s", defaultRegion)
		return defaultRegion
	}
	return region
}
