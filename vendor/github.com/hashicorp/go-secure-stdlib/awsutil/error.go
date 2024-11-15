// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsutil

import (
	"errors"

	awsRequest "github.com/aws/aws-sdk-go/aws/request"
	multierror "github.com/hashicorp/go-multierror"
)

var ErrUpstreamRateLimited = errors.New("upstream rate limited")

// CheckAWSError will examine an error and convert to a logical error if
// appropriate. If no appropriate error is found, return nil
func CheckAWSError(err error) error {
	// IsErrorThrottle will check if the error returned is one that matches
	// known request limiting errors:
	// https://github.com/aws/aws-sdk-go/blob/488d634b5a699b9118ac2befb5135922b4a77210/aws/request/retryer.go#L35
	if awsRequest.IsErrorThrottle(err) {
		return ErrUpstreamRateLimited
	}
	return nil
}

// AppendAWSError checks if the given error is a known AWS error we modify,
// and if so then returns a go-multierror, appending the original and the
// AWS error.
// If the error is not an AWS error, or not an error we wish to modify, then
// return the original error.
func AppendAWSError(err error) error {
	if awserr := CheckAWSError(err); awserr != nil {
		err = multierror.Append(err, awserr)
	}
	return err
}
