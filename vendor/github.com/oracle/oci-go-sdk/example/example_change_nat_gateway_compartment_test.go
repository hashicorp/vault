// Copyright (c) 2016, 2019, Oracle and/or its affiliates. All rights reserved.
//
//
// This script provides a basic example of how to move a NAT Gateway from one compartment to another using Go SDK.
// This script will:
//
//    * Read user configuration
//    * Construct VirtualNetworkClient using user configuration
//    * Create VCN and NAT Gateway
//    * Construct ChangeNatGatewayCompartmentDetails()
//    * Call ChangeNatGatewayCompartment() in core.VirtualNetworkClient()
//    * List NAT Gateway before and after compartment move operation
//    * Delete VCN and NAT Gateway
//
//  This script takes the following values from environment variables
//
//    * SOURCE_COMPARTMENT_ID - The OCID of the compartment where the NAT gateway and related resources will be created
//    * DESTINATION_COMPARTMENT_ID - The OCID of the compartment where the NAT gateway will be moved to
//
//

package example

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"log"
	"os"
)

const (
	displayName = "oci-go-sdk-example-ngw"
)

var (
	sourceCompartmentId, destinationCompartmentId string
)

func ExampleChangeNatGatewayCompartment() {

	// Parse environment variables to get sourceCompartmentId and destinationCompartmentId
	parseArgs()
	log.Printf("Performing operations to change NAT Gateway compartment from %s to %s", sourceCompartmentId, destinationCompartmentId)

	// Create VirtualNetworkClient with default configuration
	vcnClient, err := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)
	ctx := context.Background()

	// A VCN is required to create a NAT Gateway
	vcn := createVcnforNatGateway(ctx, vcnClient)
	log.Printf("Created VCN: %s", *vcn.Id)
	log.Printf("")

	// Create NAT Gateway
	natGateway := createNatGateway(ctx, vcnClient, vcn)

	// Change NAT Gateway's compartment
	changeNatGatewayCompartment(ctx, vcnClient, natGateway)

	fmt.Printf("Change NAT Gateway Compartment Completed")
	// Clean up resources
	defer func() {
		deleteNatGateway(ctx, vcnClient, natGateway)
		log.Printf("Deleted NAT Gateway")

		deleteVcnforNatGateway(ctx, vcnClient, vcn)
		log.Printf("Deleted VCN")
	}()

	// Output:
	// Change NAT Gateway Compartment Completed

}

func createVcnforNatGateway(ctx context.Context, c core.VirtualNetworkClient) core.Vcn {
	// create a new VCN
	request := core.CreateVcnRequest{}
	request.CidrBlock = common.String("10.0.0.0/16")
	request.CompartmentId = common.String(sourceCompartmentId)
	request.DisplayName = common.String(displayName)

	r, err := c.CreateVcn(ctx, request)
	helpers.FatalIfError(err)

	// below logic is to wait until VCN is in Available state
	pollUntilAvailable := func(r common.OCIOperationResponse) bool {
		if converted, ok := r.Response.(core.GetVcnResponse); ok {
			return converted.LifecycleState != core.VcnLifecycleStateAvailable
		}
		return true
	}

	pollGetRequest := core.GetVcnRequest{
		VcnId:           r.Id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(pollUntilAvailable),
	}

	// wait for VCN to become Available
	rsp, pollErr := c.GetVcn(ctx, pollGetRequest)
	helpers.FatalIfError(pollErr)
	return rsp.Vcn
}

func deleteVcnforNatGateway(ctx context.Context, c core.VirtualNetworkClient, vcn core.Vcn) {
	request := core.DeleteVcnRequest{
		VcnId:           vcn.Id,
		RequestMetadata: helpers.GetRequestMetadataWithDefaultRetryPolicy(),
	}

	_, err := c.DeleteVcn(ctx, request)
	helpers.FatalIfError(err)

	// should retry condition check which returns a bool value indicating whether to do retry or not
	// it checks if the lifecycle state equals Terminated
	shouldRetryFunc := func(r common.OCIOperationResponse) bool {
		if serviceError, ok := common.IsServiceError(r.Error); ok && serviceError.GetHTTPStatusCode() == 404 {
			// resource been deleted, stop retry
			return false
		}

		if converted, ok := r.Response.(core.GetVcnResponse); ok {
			return converted.LifecycleState != core.VcnLifecycleStateTerminated
		}
		return true
	}

	pollGetRequest := core.GetVcnRequest{
		VcnId:           vcn.Id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(shouldRetryFunc),
	}

	_, pollErr := c.GetVcn(ctx, pollGetRequest)
	if serviceError, ok := common.IsServiceError(pollErr); !ok ||
		(ok && serviceError.GetHTTPStatusCode() != 404) {
		// fail if the error is not service error or
		// if the error is service error and status code not equals to 404
		helpers.FatalIfError(pollErr)
	}
}

func createNatGateway(ctx context.Context, c core.VirtualNetworkClient, vcn core.Vcn) core.NatGateway {

	log.Printf("Creating NAT Gateway")
	log.Printf("=======================================")
	createNatGatewayDetails := core.CreateNatGatewayDetails{
		CompartmentId: common.String(sourceCompartmentId),
		VcnId:         vcn.Id,
		DisplayName:   common.String(displayName),
	}

	request := core.CreateNatGatewayRequest{}
	request.CreateNatGatewayDetails = createNatGatewayDetails

	r, err := c.CreateNatGateway(ctx, request)
	helpers.FatalIfError(err)

	// below logic is to wait until NAT Gateway is in Available state
	pollUntilAvailable := func(r common.OCIOperationResponse) bool {
		if converted, ok := r.Response.(core.GetNatGatewayResponse); ok {
			return converted.LifecycleState != core.NatGatewayLifecycleStateAvailable
		}
		return true
	}

	pollGetRequest := core.GetNatGatewayRequest{
		NatGatewayId:    r.Id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(pollUntilAvailable),
	}

	// wait for lifecyle become Available
	rsp, pollErr := c.GetNatGateway(ctx, pollGetRequest)
	helpers.FatalIfError(pollErr)

	log.Printf("Created NAT Gateway and waited for it to become available %v\n", rsp.NatGateway)
	log.Printf("")
	log.Printf("")

	return rsp.NatGateway
}

func getNatGateway(ctx context.Context, c core.VirtualNetworkClient, natGateway core.NatGateway) core.NatGateway {
	request := core.GetNatGatewayRequest{
		NatGatewayId: natGateway.Id,
	}
	r, err := c.GetNatGateway(ctx, request)
	helpers.FatalIfError(err)
	return r.NatGateway
}

func deleteNatGateway(ctx context.Context, c core.VirtualNetworkClient, natGateway core.NatGateway) {
	request := core.DeleteNatGatewayRequest{
		NatGatewayId:    natGateway.Id,
		RequestMetadata: helpers.GetRequestMetadataWithDefaultRetryPolicy(),
	}

	_, err := c.DeleteNatGateway(ctx, request)
	helpers.FatalIfError(err)

	// should retry condition check which returns a bool value indicating whether to do retry or not
	// it checks the lifecycle status equals to Terminated or not for this case
	shouldRetryFunc := func(r common.OCIOperationResponse) bool {
		if serviceError, ok := common.IsServiceError(r.Error); ok && serviceError.GetHTTPStatusCode() == 404 {
			// resource been deleted, stop retry
			return false
		}

		if converted, ok := r.Response.(core.GetNatGatewayResponse); ok {
			return converted.LifecycleState != core.NatGatewayLifecycleStateTerminated
		}
		return true
	}

	pollGetRequest := core.GetNatGatewayRequest{
		NatGatewayId:    natGateway.Id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(shouldRetryFunc),
	}

	_, pollErr := c.GetNatGateway(ctx, pollGetRequest)
	if serviceError, ok := common.IsServiceError(pollErr); !ok ||
		(ok && serviceError.GetHTTPStatusCode() != 404) {
		// fail if the error is not service error or
		// if the error is service error and status code not equals to 404
		helpers.FatalIfError(pollErr)
	}
}

func changeNatGatewayCompartment(ctx context.Context, c core.VirtualNetworkClient, natGateway core.NatGateway) {
	log.Printf("Changing NAT Gateway's compartment")
	log.Printf("=======================================")
	changeNatGatewayCompartmentDetails := core.ChangeNatGatewayCompartmentDetails{
		CompartmentId: common.String(destinationCompartmentId),
	}

	request := core.ChangeNatGatewayCompartmentRequest{}
	request.NatGatewayId = natGateway.Id
	request.ChangeNatGatewayCompartmentDetails = changeNatGatewayCompartmentDetails

	_, err := c.ChangeNatGatewayCompartment(ctx, request)
	helpers.FatalIfError(err)
	updatedNatGateway := getNatGateway(ctx, c, natGateway)
	log.Printf("NAT Gateway's compartment has been changed  : %v\n", updatedNatGateway)
	log.Printf("")
	log.Printf("")
}

func envUsage() {
	log.Printf("Please set the following environment variables to use ChangeInstanceCompartment()")
	log.Printf(" ")
	log.Printf("   SOURCE_COMPARTMENT_ID    # Required: Source Compartment Id")
	log.Printf("   DESTINATION_COMPARTMENT_ID    # Required: Destination Compartment Id")
	log.Printf(" ")
	os.Exit(1)
}

func parseArgs() {

	sourceCompartmentId = os.Getenv("SOURCE_COMPARTMENT_ID")
	destinationCompartmentId = os.Getenv("DESTINATION_COMPARTMENT_ID")

	if sourceCompartmentId == "" || destinationCompartmentId == "" {
		envUsage()
	}

	log.Printf("SOURCE_COMPARTMENT_ID     : %s", sourceCompartmentId)
	log.Printf("DESTINATION_COMPARTMENT_ID  : %s", destinationCompartmentId)

}
