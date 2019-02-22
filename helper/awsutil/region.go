package awsutil

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetOrDefaultRegion(configuredRegion string) string {
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
		// This is odd, let's do what we can to pluckily carry on.
		return defaultRegion
	}

	// This will hang for ~10 seconds if we aren't running on an EC2 instance.
	region, err := ec2metadata.New(sess).Region()
	if err != nil {
		// This means we're not on an EC2 instance.
		return defaultRegion
	}
	return region
}
