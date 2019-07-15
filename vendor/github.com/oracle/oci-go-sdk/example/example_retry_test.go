// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
//
// Example code for using retry for SDK APIs
//

package example

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/identity"
)

// ExampleRetry shows how to use retry for Create and Delete groups, please
// refer to example_core_test.go->ExampleLaunchInstance for more examples
func ExampleRetry() {
	// create and delete group with retry
	client, clerr := identity.NewIdentityClientWithConfigurationProvider(common.DefaultConfigProvider())
	ctx := context.Background()
	helpers.FatalIfError(clerr)

	request := identity.CreateGroupRequest{}
	request.CompartmentId = helpers.RootCompartmentID()
	request.Name = common.String("GoSDK_Sample_Group")
	request.Description = common.String("GoSDK Sample Group Description")

	// maximum times of retry
	attempts := uint(10)

	// retry for all non-200 status code
	retryOnAllNon200ResponseCodes := func(r common.OCIOperationResponse) bool {
		return !(r.Error == nil && 199 < r.Response.HTTPResponse().StatusCode && r.Response.HTTPResponse().StatusCode < 300)
	}

	nextDuration := func(r common.OCIOperationResponse) time.Duration {
		// you might want wait longer for next retry when your previous one failed
		// this function will return the duration as:
		// 1s, 2s, 4s, 8s, 16s, 32s, 64s etc...
		return time.Duration(math.Pow(float64(2), float64(r.AttemptNumber-1))) * time.Second
	}

	defaultRetryPolicy := common.NewRetryPolicy(attempts, retryOnAllNon200ResponseCodes, nextDuration)

	// create request metadata for retry
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: &defaultRetryPolicy,
	}

	resp, err := client.CreateGroup(ctx, request)
	helpers.FatalIfError(err)
	fmt.Println("Creating Group")

	// Get with polling
	shouldRetry := func(r common.OCIOperationResponse) bool {
		if _, isServiceError := common.IsServiceError(r.Error); isServiceError {
			// not service error, could be network error or other errors which prevents
			// request send to server, will do retry here
			return true
		}

		if converted, ok := r.Response.(identity.GetGroupResponse); ok {
			// do the retry until lifecycle state become active
			return converted.LifecycleState != identity.GroupLifecycleStateActive
		}

		return true
	}

	lifecycleStateCheckRetryPolicy := common.NewRetryPolicy(attempts, shouldRetry, nextDuration)

	getRequest := identity.GetGroupRequest{
		GroupId: resp.Id,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: &lifecycleStateCheckRetryPolicy,
		},
	}

	_, errAfterPolling := client.GetGroup(ctx, getRequest)
	helpers.FatalIfError(errAfterPolling)
	fmt.Println("Group Created")

	defer func() {
		// if we've successfully created a group, make sure that we delete it
		rDel := identity.DeleteGroupRequest{
			GroupId: resp.Id,
			RequestMetadata: common.RequestMetadata{
				RetryPolicy: &defaultRetryPolicy,
			},
		}

		_, err = client.DeleteGroup(ctx, rDel)
		helpers.FatalIfError(err)
		fmt.Println("Group Deleted")
	}()

	// Output:
	// Creating Group
	// Group Created
	// Group Deleted
}
