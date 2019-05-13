package awsutil

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	hclog "github.com/hashicorp/go-hclog"
)

// "us-east-1 is used because it's where AWS first provides support for new features,
// is a widely used region, and is the most common one for some services like STS.
const DefaultRegion = "us-east-1"

var ec2MetadataBaseURL = "http://169.254.169.254"

/*
It's impossible to mimic "normal" AWS behavior here because it's not consistent
or well-defined. For example, boto3, the Python SDK (which the aws cli uses),
loads `~/.aws/config` by default and only reads the `AWS_DEFAULT_REGION` environment
variable (and not `AWS_REGION`, while the golang SDK does _mostly_ the opposite -- it
reads the region **only** from `AWS_REGION` and not at all `~/.aws/config`, **unless**
the `AWS_SDK_LOAD_CONFIG` environment variable is set. So, we must define our own
approach to walking AWS config and deciding what to use.

Our chosen approach is:

	"More specific takes precedence over less specific."

1. User-provided configuration is the most explicit.
2. Environment variables are potentially shared across many invocations and so they have less precedence.
3. Configuration in `~/.aws/config` is shared across all invocations of a given user and so this has even less precedence.
4. Configuration retrieved from the EC2 instance metadata service is shared by all invocations on a given machine, and so it has the lowest precedence.

This approach should be used in future updates to this logic.
*/
func GetOrDefaultRegion(logger hclog.Logger, configuredRegion string) string {
	if configuredRegion != "" {
		return configuredRegion
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
