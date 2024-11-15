// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsutil

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/errwrap"
)

// "us-east-1 is used because it's where AWS first provides support for new features,
// is a widely used region, and is the most common one for some services like STS.
const DefaultRegion = "us-east-1"

// This is nil by default, but is exposed in case it needs to be changed for tests.
var ec2Endpoint *string

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
func GetRegion(configuredRegion string) (string, error) {
	if configuredRegion != "" {
		return configuredRegion, nil
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return "", errwrap.Wrapf("got error when starting session: {{err}}", err)
	}

	region := aws.StringValue(sess.Config.Region)
	if region != "" {
		return region, nil
	}

	metadata := ec2metadata.New(sess, &aws.Config{
		Endpoint:                          ec2Endpoint,
		EC2MetadataDisableTimeoutOverride: aws.Bool(true),
		HTTPClient: &http.Client{
			Timeout: time.Second,
		},
	})
	if !metadata.Available() {
		return DefaultRegion, nil
	}

	region, err = metadata.Region()
	if err != nil {
		return "", errwrap.Wrapf("unable to retrieve region from instance metadata: {{err}}", err)
	}

	return region, nil
}
