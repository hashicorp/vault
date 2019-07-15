// Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
//
// Example code for Compute Management Services API
//
//

/**
 * This class provides an example of how you can create and manage an Instance Pool. It will:
 * <ul>
 * <li>Create the InstanceConfiguration</li>
 * <li>Create a pool of size 1 based off of that configuration.</li>
 * <li>Wait for the pool to go to Running state.</li>
 * <li>Update the pool to a size of 2.</li>
 * <li>Wait for the InstancePool to scale up.</li>
 * <li>Attached a load balancer to the pool.</li>
 * <li>Wait for the load balancer to become ATTACHED.</li>
 * <li>Clean everything up.</li>
 * </ul>
 */

package example

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/oracle/oci-go-sdk/example/helpers"
)

var (
	imageId, ad, subnetId, loadBalancerId, loadBalancerBackendSetName, compartmentId string
)

// Example to showcase instance pool create and operations, and eventual teardown
func ExampleCreateAndWaitForRunningInstancePool() {
	InstancePoolsParseEnvironmentVariables()

	ctx := context.Background()

	computeMgmtClient, err := core.NewComputeManagementClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	createInstanceConfigurationResponse, _ := createInstanceConfiguration(ctx, computeMgmtClient, imageId, compartmentId)
	fmt.Println("Instance configuration created")

	instanceConfiguration := createInstanceConfigurationResponse.InstanceConfiguration

	instancePool, _ := createInstancePool(ctx, computeMgmtClient, *instanceConfiguration.Id, subnetId, ad, compartmentId)
	fmt.Println("Instance pool created")

	// waiting until the poll reaches running state
	pollUntilDesiredState(ctx, computeMgmtClient, instancePool, core.InstancePoolLifecycleStateRunning)

	setInstancePoolSize(ctx, computeMgmtClient, *instancePool.Id, 2)

	// updating the pool size will make it go to scaling state first and then to running
	pollUntilDesiredState(ctx, computeMgmtClient, instancePool, core.InstancePoolLifecycleStateScaling)
	pollUntilDesiredState(ctx, computeMgmtClient, instancePool, core.InstancePoolLifecycleStateRunning)

	// attach load balancer to the created pool
	attachLBtoInstancePool(ctx, computeMgmtClient, *instancePool.Id, loadBalancerId, loadBalancerBackendSetName)

	// poll for instance pool until the lb becomes attached
	pollUntilDesiredLBAttachmentState(ctx, computeMgmtClient, instancePool, core.InstancePoolLoadBalancerAttachmentLifecycleStateAttached)

	// gets the targeted load balancer information
	getLbAttachmentForPool(ctx, computeMgmtClient, *instancePool.Id)

	// clean up resources
	defer func() {
		terminateInstancePool(ctx, computeMgmtClient, *instancePool.Id)
		fmt.Println("Terminated Instance Pool")

		deleteInstanceConfiguration(ctx, computeMgmtClient, *instanceConfiguration.Id)
		fmt.Println("Deleted Instance Configuration")
	}()

	// Output:
	// Instance configuration created
	// Instance pool created
	// Instance pool is RUNNING
	// Instance pool is SCALING
	// Instance pool is RUNNING
	// Instance pool attachment is ATTACHED
	// Instance pool attachment has vnic PrimaryVnic
	// Terminated Instance Pool
	// Deleted Instance Configuration
}

// Usage printing
func InstancePoolsUsage() {
	log.Printf("Please set the following environment variables to run Instance Pool sample")
	log.Printf(" ")
	log.Printf("   IMAGE_ID       # Required: Image Id to use")
	log.Printf("   COMPARTMENT_ID    # Required: Compartment Id to use")
	log.Printf("   AD          # Required: AD to use")
	log.Printf("   SUBNET_ID   # Required: Subnet to use")
	log.Printf("   LB_ID   # Required: Load balancer to use")
	log.Printf("   LB_BACKEND_SET_NAME   # Required: Load balancer backend set name to use")
	log.Printf(" ")
	os.Exit(1)
}

// Args parser
func InstancePoolsParseEnvironmentVariables() {

	imageId = os.Getenv("IMAGE_ID")
	compartmentId = os.Getenv("COMPARTMENT_ID")
	ad = os.Getenv("AD")
	subnetId = os.Getenv("SUBNET_ID")
	loadBalancerId = os.Getenv("LB_ID")
	loadBalancerBackendSetName = os.Getenv("LB_BACKEND_SET_NAME")

	if imageId == "" ||
		compartmentId == "" ||
		ad == "" ||
		subnetId == "" ||
		loadBalancerId == "" ||
		loadBalancerBackendSetName == "" {
		InstancePoolsUsage()
	}

	log.Printf("IMAGE_ID     : %s", imageId)
	log.Printf("COMPARTMENT_ID  : %s", compartmentId)
	log.Printf("AD     : %s", ad)
	log.Printf("SUBNET_ID  : %s", subnetId)
	log.Printf("LB_ID  : %s", loadBalancerId)
	log.Printf("LB_BACKEND_SET_NAME  : %s", loadBalancerBackendSetName)
}

// helper method to create an instance configuration
func createInstanceConfiguration(ctx context.Context, client core.ComputeManagementClient, imageId string, compartmentId string) (response core.CreateInstanceConfigurationResponse, err error) {
	vnicDetails := core.InstanceConfigurationCreateVnicDetails{}

	sourceDetails := core.InstanceConfigurationInstanceSourceViaImageDetails{
		ImageId: &imageId,
	}

	displayName := "Instance Configuration Example"
	shape := "VM.Standard2.1"

	launchDetails := core.InstanceConfigurationLaunchInstanceDetails{
		CompartmentId:     &compartmentId,
		DisplayName:       &displayName,
		CreateVnicDetails: &vnicDetails,
		Shape:             &shape,
		SourceDetails:     &sourceDetails,
	}

	instanceDetails := core.ComputeInstanceDetails{
		LaunchDetails: &launchDetails,
	}

	configurationDetails := core.CreateInstanceConfigurationDetails{
		DisplayName:     &displayName,
		CompartmentId:   &compartmentId,
		InstanceDetails: &instanceDetails,
	}

	req := core.CreateInstanceConfigurationRequest{
		CreateInstanceConfiguration: configurationDetails,
	}

	response, err = client.CreateInstanceConfiguration(ctx, req)
	helpers.FatalIfError(err)

	return
}

// helper method to create an instance pool
func createInstancePool(ctx context.Context, client core.ComputeManagementClient, instanceConfigurationId string,
	subnetId string, availabilityDomain string, compartmentId string) (response core.CreateInstancePoolResponse, err error) {

	displayName := "Instance Pool Example"
	size := 1

	req := core.CreateInstancePoolRequest{
		CreateInstancePoolDetails: core.CreateInstancePoolDetails{
			CompartmentId:           &compartmentId,
			InstanceConfigurationId: &instanceConfigurationId,
			PlacementConfigurations: []core.CreateInstancePoolPlacementConfigurationDetails{
				{
					PrimarySubnetId:    &subnetId,
					AvailabilityDomain: &availabilityDomain,
				},
			},
			Size:        &size,
			DisplayName: &displayName,
		},
	}

	response, err = client.CreateInstancePool(ctx, req)
	return
}

// helper method to terminate an instance configuration
func terminateInstancePool(ctx context.Context, client core.ComputeManagementClient,
	poolId string) (response core.TerminateInstancePoolResponse, err error) {

	req := core.TerminateInstancePoolRequest{
		InstancePoolId: &poolId,
	}

	response, err = client.TerminateInstancePool(ctx, req)
	helpers.FatalIfError(err)

	return
}

// helper method to delete an instance configuration
func deleteInstanceConfiguration(ctx context.Context, client core.ComputeManagementClient,
	instanceConfigurationId string) (response core.DeleteInstanceConfigurationResponse, err error) {

	req := core.DeleteInstanceConfigurationRequest{
		InstanceConfigurationId: &instanceConfigurationId,
	}

	response, err = client.DeleteInstanceConfiguration(ctx, req)
	helpers.FatalIfError(err)

	return
}

// helper method to update an instance pool size
func setInstancePoolSize(ctx context.Context, client core.ComputeManagementClient,
	poolId string, newSize int) (response core.UpdateInstancePoolResponse, err error) {

	updateDetails := core.UpdateInstancePoolDetails{
		Size: &newSize,
	}

	req := core.UpdateInstancePoolRequest{
		InstancePoolId:            &poolId,
		UpdateInstancePoolDetails: updateDetails,
	}

	response, err = client.UpdateInstancePool(ctx, req)
	helpers.FatalIfError(err)

	return
}

// helper method to update an instance pool size
func attachLBtoInstancePool(ctx context.Context, client core.ComputeManagementClient,
	poolId string, loadBalancerId string, lbBackendSetName string) (response core.AttachLoadBalancerResponse, err error) {

	port := 80
	vnic := "PrimaryVnic"

	attachDetails := core.AttachLoadBalancerDetails{
		LoadBalancerId: &loadBalancerId,
		BackendSetName: &loadBalancerBackendSetName,
		Port:           &port,
		VnicSelection:  &vnic,
	}

	req := core.AttachLoadBalancerRequest{
		InstancePoolId:            &poolId,
		AttachLoadBalancerDetails: attachDetails,
	}

	response, err = client.AttachLoadBalancer(ctx, req)
	helpers.FatalIfError(err)

	return
}

// helper method to pool until an instance pool reaches the specified desired state
func pollUntilDesiredState(ctx context.Context, computeMgmtClient core.ComputeManagementClient,
	instancePool core.CreateInstancePoolResponse, desiredState core.InstancePoolLifecycleStateEnum) {
	// should retry condition check which returns a bool value indicating whether to do retry or not
	// it checks the lifecycle status equals to Running or not for this case
	shouldRetryFunc := func(r common.OCIOperationResponse) bool {
		if converted, ok := r.Response.(core.GetInstancePoolResponse); ok {
			return converted.LifecycleState != desiredState
		}
		return true
	}
	// create get instance pool request with a retry policy which takes a function
	// to determine shouldRetry or not
	pollingGetRequest := core.GetInstancePoolRequest{
		InstancePoolId:  instancePool.Id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(shouldRetryFunc),
	}
	_, pollError := computeMgmtClient.GetInstancePool(ctx, pollingGetRequest)
	helpers.FatalIfError(pollError)
	fmt.Println("Instance pool is", desiredState)
}

// helper method to pool until an instance pool lb attachment reaches the specified desired state
func pollUntilDesiredLBAttachmentState(ctx context.Context, computeMgmtClient core.ComputeManagementClient,
	instancePool core.CreateInstancePoolResponse, desiredState core.InstancePoolLoadBalancerAttachmentLifecycleStateEnum) {
	// should retry condition check which returns a bool value indicating whether to do retry or not
	// it checks the lifecycle status equals to Running or not for this case
	shouldRetryFunc := func(r common.OCIOperationResponse) bool {
		if converted, ok := r.Response.(core.GetInstancePoolResponse); ok {
			attachments := converted.LoadBalancers

			for i := range attachments {
				if attachments[i].LifecycleState != desiredState {
					return true
				}
			}

			return false
		}
		return true
	}
	// create get instance pool request with a retry policy which takes a function
	// to determine shouldRetry or not
	pollingGetRequest := core.GetInstancePoolRequest{
		InstancePoolId:  instancePool.Id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(shouldRetryFunc),
	}
	_, pollError := computeMgmtClient.GetInstancePool(ctx, pollingGetRequest)
	helpers.FatalIfError(pollError)
	fmt.Println("Instance pool attachment is", desiredState)
}

// function showing how to get lb attachment info for a pool
func getLbAttachmentForPool(ctx context.Context, computeMgmtClient core.ComputeManagementClient,
	instancePoolId string) {

	// gets the fresh instance pool info which, after lb attaching, now has lb attachment information
	getReq := core.GetInstancePoolRequest{
		InstancePoolId: &instancePoolId,
	}

	instancePoolResp, _ := computeMgmtClient.GetInstancePool(ctx, getReq)

	// takes the 1st load balancer attachment id from the pool
	lbAttachmentId := instancePoolResp.LoadBalancers[0].Id

	req := core.GetInstancePoolLoadBalancerAttachmentRequest{
		InstancePoolId:                       &instancePoolId,
		InstancePoolLoadBalancerAttachmentId: lbAttachmentId,
	}

	response, _ := computeMgmtClient.GetInstancePoolLoadBalancerAttachment(ctx, req)
	fmt.Println("Instance pool attachment has vnic", *response.VnicSelection)
}
