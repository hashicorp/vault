// Copyright (c) 2016, 2019, Oracle and/or its affiliates. All rights reserved.
//
//
// This script provides a basic example of how to move a service gateway from one compartment to another using Go SDK.
// This script will:
//
//    * Read user configuration
//    * Construct VirtualNetworkClient using user configuration
//    * Create VCN and Service Gateway
//    * Call ChangeServiceGatewayCompartment() in core.VirtualNetworkClient()
//    * Get Service Gateway to see the updated compartment ID
//    * Delete Service Gateway and VCN
//    * List Instance and its attached resources before and after move operation
//
//  This script takes the following values from environment variables
//
//    * SRC_COMPARTMENT_ID    - Source Compartment ID where the service gateway and VCN should be created
//    * DEST_COMPARTMENT_ID   - Destination Compartment ID where the service gateway should be moved to
//
// Additionally this script assumes that the Default OCI config is setup

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
	serviceGatewayVcnDisplayName = "OCI-GOSDK-Sample"
)

var (
	srcCompartmentId, destCompartmentId string
)

func ExampleChangeServiceGatewayCompartment() {

	// Parse environment variables to get srcCompartmentId, destCompartmentId
	parseEnvVariables()

	// Create VirtualNetworkClient with default configuration
	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)
	ctx := context.Background()

	log.Printf("Create Vcn ")
	vcn := createSgwVcn(ctx, client)
	log.Printf("VCN OCID : %s", *vcn.Id)

	log.Printf("Create Service Gateway")
	sgw := createServiceGateway(ctx, client, vcn)
	log.Printf("Service Gateway OCID : %s", *sgw.Id)

	log.Printf("Change Service Gateway Compartment")
	changeServiceGatewayCompartment(ctx, client, sgw)
	updatedsgw := getServiceGateway(ctx, client, sgw)
	log.Printf("Updated Service Gateway Compartment : %s", *updatedsgw.CompartmentId)

	fmt.Printf("change compartment completed")

	// clean up resources
	defer func() {
		log.Printf("Delete Service Gateway")
		deleteServiceGateway(ctx, client, sgw)
		log.Printf("Deleted Service Gateway")

		log.Printf("Delete VCN")
		deleteSgwVcn(ctx, client, vcn)
		log.Printf("Deleted VCN")
	}()

	// Output:
	// change compartment completed

}

func createSgwVcn(ctx context.Context, c core.VirtualNetworkClient) core.Vcn {
	// create a new VCN
	request := core.CreateVcnRequest{}
	request.CidrBlock = common.String("10.0.0.0/16")
	request.CompartmentId = common.String(srcCompartmentId)
	request.DisplayName = common.String(serviceGatewayVcnDisplayName)

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

	// wait for lifecyle become Available
	rsp, pollErr := c.GetVcn(ctx, pollGetRequest)
	helpers.FatalIfError(pollErr)
	return rsp.Vcn
}

func deleteSgwVcn(ctx context.Context, c core.VirtualNetworkClient, vcn core.Vcn) {
	request := core.DeleteVcnRequest{
		VcnId:           vcn.Id,
		RequestMetadata: helpers.GetRequestMetadataWithDefaultRetryPolicy(),
	}

	_, err := c.DeleteVcn(ctx, request)
	helpers.FatalIfError(err)

	// should retry condition check which returns a bool value indicating whether to do retry or not
	// it checks the lifecycle status equals to Terminated or not for this case
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

func createServiceGateway(ctx context.Context, c core.VirtualNetworkClient, vcn core.Vcn) core.ServiceGateway {

	// Update the services field to required Oracle Services
	var services = []core.ServiceIdRequestDetails{}
	createServiceGatewayDetails := core.CreateServiceGatewayDetails{
		CompartmentId: common.String(srcCompartmentId),
		VcnId:         vcn.Id,
		DisplayName:   common.String(serviceGatewayVcnDisplayName),
		Services:      services,
	}

	// create a new VCN
	request := core.CreateServiceGatewayRequest{}
	request.CreateServiceGatewayDetails = createServiceGatewayDetails

	r, err := c.CreateServiceGateway(ctx, request)
	helpers.FatalIfError(err)

	// below logic is to wait until VCN is in Available state
	pollUntilAvailable := func(r common.OCIOperationResponse) bool {
		if converted, ok := r.Response.(core.GetServiceGatewayResponse); ok {
			return converted.LifecycleState != core.ServiceGatewayLifecycleStateAvailable
		}
		return true
	}

	pollGetRequest := core.GetServiceGatewayRequest{
		ServiceGatewayId: r.Id,
		RequestMetadata:  helpers.GetRequestMetadataWithCustomizedRetryPolicy(pollUntilAvailable),
	}

	// wait for lifecyle become Available
	rsp, pollErr := c.GetServiceGateway(ctx, pollGetRequest)
	helpers.FatalIfError(pollErr)
	return rsp.ServiceGateway
}

func deleteServiceGateway(ctx context.Context, c core.VirtualNetworkClient, serviceGateway core.ServiceGateway) {
	request := core.DeleteServiceGatewayRequest{
		ServiceGatewayId: serviceGateway.Id,
		RequestMetadata:  helpers.GetRequestMetadataWithDefaultRetryPolicy(),
	}

	_, err := c.DeleteServiceGateway(ctx, request)
	helpers.FatalIfError(err)

	// should retry condition check which returns a bool value indicating whether to do retry or not
	// it checks the lifecycle status equals to Terminated or not for this case
	shouldRetryFunc := func(r common.OCIOperationResponse) bool {
		if serviceError, ok := common.IsServiceError(r.Error); ok && serviceError.GetHTTPStatusCode() == 404 {
			// resource been deleted, stop retry
			return false
		}

		if converted, ok := r.Response.(core.GetServiceGatewayResponse); ok {
			return converted.LifecycleState != core.ServiceGatewayLifecycleStateTerminated
		}
		return true
	}

	pollGetRequest := core.GetServiceGatewayRequest{
		ServiceGatewayId: serviceGateway.Id,
		RequestMetadata:  helpers.GetRequestMetadataWithCustomizedRetryPolicy(shouldRetryFunc),
	}

	_, pollErr := c.GetServiceGateway(ctx, pollGetRequest)
	if serviceError, ok := common.IsServiceError(pollErr); !ok ||
		(ok && serviceError.GetHTTPStatusCode() != 404) {
		// fail if the error is not service error or
		// if the error is service error and status code not equals to 404
		helpers.FatalIfError(pollErr)
	}
}

func changeServiceGatewayCompartment(ctx context.Context, c core.VirtualNetworkClient, serviceGateway core.ServiceGateway) {
	changeCompartmentDetails := core.ChangeServiceGatewayCompartmentDetails{
		CompartmentId: common.String(destCompartmentId),
	}

	request := core.ChangeServiceGatewayCompartmentRequest{}
	request.ServiceGatewayId = serviceGateway.Id
	request.ChangeServiceGatewayCompartmentDetails = changeCompartmentDetails

	_, err := c.ChangeServiceGatewayCompartment(ctx, request)
	helpers.FatalIfError(err)

}

func getServiceGateway(ctx context.Context, c core.VirtualNetworkClient, serviceGateway core.ServiceGateway) core.ServiceGateway {
	request := core.GetServiceGatewayRequest{
		ServiceGatewayId: serviceGateway.Id,
	}

	r, err := c.GetServiceGateway(ctx, request)
	helpers.FatalIfError(err)

	return r.ServiceGateway
}

func printUsage() {
	fmt.Printf("Please set the following environment variables to use ChangeServiceGatewayCompartment()")
	fmt.Printf(" ")
	fmt.Printf("   SRC_COMPARTMENT_ID       # Required: Source Compartment Id")
	fmt.Printf("   DEST_COMPARTMENT_ID	    # Required: Destination Compartment Id")
	fmt.Printf(" ")
	os.Exit(1)
}

func parseEnvVariables() {

	srcCompartmentId = os.Getenv("SRC_COMPARTMENT_ID")
	destCompartmentId = os.Getenv("DEST_COMPARTMENT_ID")

	if srcCompartmentId == "" || destCompartmentId == "" {
		printUsage()
	}

	log.Printf("SRC_COMPARTMENT_ID     : %s", srcCompartmentId)
	log.Printf("DEST_COMPARTMENT_ID  : %s", destCompartmentId)
}
