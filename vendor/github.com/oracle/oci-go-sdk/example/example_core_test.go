// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
//
// Example code for Core Services API
//

package example

import (
	"context"
	"fmt"
	"log"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/oracle/oci-go-sdk/example/helpers"
)

const (
	vcnDisplayName     = "OCI-GOSDK-Sample-VCN"
	subnetDisplayName1 = "OCI-GOSDK-Sample-Subnet1"
	subnetDisplayName2 = "OCI-GOSDK-Sample-Subnet2"
	subnetDisplayName3 = "OCI-GOSDK-Sample-Subnet3"

	// replace following variables with your instance info
	// this is used by ExampleCreateImageDetails_Polymorphic
	objectStorageURIWtihImage = "[The Object Storage URL for the image which will be used to create an image.]"
)

// ExampleLaunchInstance does create an instance
// NOTE: launch instance will create a new instance and VCN. please make sure delete the instance
// after execute this sample code, otherwise, you will be charged for the running instance
func ExampleLaunchInstance() {
	c, err := core.NewComputeClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)
	ctx := context.Background()

	// create the launch instance request
	request := core.LaunchInstanceRequest{}
	request.CompartmentId = helpers.CompartmentID()
	request.DisplayName = common.String("OCI-Sample-Instance")
	request.AvailabilityDomain = helpers.AvailabilityDomain()

	// create a subnet or get the one already created
	subnet := CreateOrGetSubnet()
	fmt.Println("subnet created")
	request.SubnetId = subnet.Id

	// get a image
	image := listImages(ctx, c)[30]
	fmt.Println("list images")
	request.ImageId = image.Id

	// get all the shapes and filter the list by compatibility with the image
	shapes := listShapes(ctx, c, request.ImageId)
	fmt.Println("list shapes")
	request.Shape = shapes[1].Shape

	// default retry policy will retry on non-200 response
	request.RequestMetadata = helpers.GetRequestMetadataWithDefaultRetryPolicy()

	createResp, err := c.LaunchInstance(ctx, request)
	helpers.FatalIfError(err)

	fmt.Println("launching instance")

	// should retry condition check which returns a bool value indicating whether to do retry or not
	// it checks the lifecycle status equals to Running or not for this case
	shouldRetryFunc := func(r common.OCIOperationResponse) bool {
		if converted, ok := r.Response.(core.GetInstanceResponse); ok {
			return converted.LifecycleState != core.InstanceLifecycleStateRunning
		}
		return true
	}

	// create get instance request with a retry policy which takes a function
	// to determine shouldRetry or not
	pollingGetRequest := core.GetInstanceRequest{
		InstanceId:      createResp.Instance.Id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(shouldRetryFunc),
	}

	instance, pollError := c.GetInstance(ctx, pollingGetRequest)
	helpers.FatalIfError(pollError)

	fmt.Println("instance launched")

	attachVnicResponse, err := c.AttachVnic(context.Background(), core.AttachVnicRequest{
		AttachVnicDetails: core.AttachVnicDetails{
			CreateVnicDetails: &core.CreateVnicDetails{
				SubnetId:       subnet.Id,
				AssignPublicIp: common.Bool(true),
			},
			InstanceId: instance.Id,
		},
	})

	helpers.FatalIfError(err)
	fmt.Println("vnic attached")

	_, err = c.DetachVnic(context.Background(), core.DetachVnicRequest{
		VnicAttachmentId: attachVnicResponse.Id,
	})

	helpers.FatalIfError(err)
	fmt.Println("vnic dettached")

	defer func() {
		terminateInstance(ctx, c, createResp.Id)

		client, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
		helpers.FatalIfError(clerr)

		vcnID := subnet.VcnId
		deleteSubnet(ctx, client, subnet.Id)
		deleteVcn(ctx, client, vcnID)
	}()

	// Output:
	// subnet created
	// list images
	// list shapes
	// launching instance
	// instance launched
	// vnic attached
	// vnic dettached
	// terminating instance
	// instance terminated
	// deleteing subnet
	// subnet deleted
	// deleteing VCN
	// VCN deleted
}

// ExampleCreateImageDetails_Polymorphic creates a boot disk image for the specified instance or
// imports an exported image from the Oracle Cloud Infrastructure Object Storage service.
func ExampleCreateImageDetails_Polymorphic() {
	request := core.CreateImageRequest{}
	request.CompartmentId = helpers.CompartmentID()

	// you can import an image based on the Object Storage URL 'core.ImageSourceViaObjectStorageUriDetails'
	// or based on the namespace, bucket name and object name 'core.ImageSourceViaObjectStorageTupleDetails'
	// following example shows how to import image from object storage uri, you can use another one:
	// request.ImageSourceDetails = core.ImageSourceViaObjectStorageTupleDetails
	sourceDetails := core.ImageSourceViaObjectStorageUriDetails{}
	sourceDetails.SourceUri = common.String(objectStorageURIWtihImage)

	request.ImageSourceDetails = sourceDetails

	c, err := core.NewComputeClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	_, err = c.CreateImage(context.Background(), request)
	helpers.FatalIfError(err)
	fmt.Println("image created")
}

// CreateOrGetVcn either creates a new Virtual Cloud Network (VCN) or get the one already exist
func CreateOrGetVcn() core.Vcn {
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	ctx := context.Background()

	vcnItems := listVcns(ctx, c)

	for _, element := range vcnItems {
		if *element.DisplayName == vcnDisplayName {
			// VCN already created, return it
			return element
		}
	}

	// create a new VCN
	request := core.CreateVcnRequest{}
	request.CidrBlock = common.String("10.0.0.0/16")
	request.CompartmentId = helpers.CompartmentID()
	request.DisplayName = common.String(vcnDisplayName)
	request.DnsLabel = common.String("vcndns")

	r, err := c.CreateVcn(ctx, request)
	helpers.FatalIfError(err)
	return r.Vcn
}

// CreateSubnet creates a new subnet or get the one already exist
func CreateOrGetSubnet() core.Subnet {
	return CreateOrGetSubnetWithDetails(
		common.String(subnetDisplayName1),
		common.String("10.0.0.0/24"),
		common.String("subnetdns1"),
		helpers.AvailabilityDomain())
}

// CreateOrGetSubnetWithDetails either creates a new Virtual Cloud Network (VCN) or get the one already exist
// with detail info
func CreateOrGetSubnetWithDetails(displayName *string, cidrBlock *string, dnsLabel *string, availableDomain *string) core.Subnet {
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	ctx := context.Background()

	subnets := listSubnets(ctx, c)

	if displayName == nil {
		displayName = common.String(subnetDisplayName1)
	}

	// check if the subnet has already been created
	for _, element := range subnets {
		if *element.DisplayName == *displayName {
			// find the subnet, return it
			return element
		}
	}

	// create a new subnet
	request := core.CreateSubnetRequest{}
	request.AvailabilityDomain = availableDomain
	request.CompartmentId = helpers.CompartmentID()
	request.CidrBlock = cidrBlock
	request.DisplayName = displayName
	request.DnsLabel = dnsLabel
	request.RequestMetadata = helpers.GetRequestMetadataWithDefaultRetryPolicy()

	vcn := CreateOrGetVcn()
	request.VcnId = vcn.Id

	r, err := c.CreateSubnet(ctx, request)
	helpers.FatalIfError(err)

	// retry condition check, stop unitl return true
	pollUntilAvailable := func(r common.OCIOperationResponse) bool {
		if converted, ok := r.Response.(core.GetSubnetResponse); ok {
			return converted.LifecycleState != core.SubnetLifecycleStateAvailable
		}
		return true
	}

	pollGetRequest := core.GetSubnetRequest{
		SubnetId:        r.Id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(pollUntilAvailable),
	}

	// wait for lifecyle become running
	_, pollErr := c.GetSubnet(ctx, pollGetRequest)
	helpers.FatalIfError(pollErr)

	// update the security rules
	getReq := core.GetSecurityListRequest{
		SecurityListId: common.String(r.SecurityListIds[0]),
	}

	getResp, err := c.GetSecurityList(ctx, getReq)
	helpers.FatalIfError(err)

	// this security rule allows remote control the instance
	portRange := core.PortRange{
		Max: common.Int(1521),
		Min: common.Int(1521),
	}

	newRules := append(getResp.IngressSecurityRules, core.IngressSecurityRule{
		Protocol: common.String("6"), // TCP
		Source:   common.String("0.0.0.0/0"),
		TcpOptions: &core.TcpOptions{
			DestinationPortRange: &portRange,
		},
	})

	updateReq := core.UpdateSecurityListRequest{
		SecurityListId: common.String(r.SecurityListIds[0]),
	}

	updateReq.IngressSecurityRules = newRules

	_, err = c.UpdateSecurityList(ctx, updateReq)
	helpers.FatalIfError(err)

	return r.Subnet
}

func listVcns(ctx context.Context, c core.VirtualNetworkClient) []core.Vcn {
	request := core.ListVcnsRequest{
		CompartmentId: helpers.CompartmentID(),
	}

	r, err := c.ListVcns(ctx, request)
	helpers.FatalIfError(err)
	return r.Items
}

func listSubnets(ctx context.Context, c core.VirtualNetworkClient) []core.Subnet {
	vcn := CreateOrGetVcn()

	request := core.ListSubnetsRequest{
		CompartmentId: helpers.CompartmentID(),
		VcnId:         vcn.Id,
	}

	r, err := c.ListSubnets(ctx, request)
	helpers.FatalIfError(err)
	return r.Items
}

// ListImages lists the available images in the specified compartment.
func listImages(ctx context.Context, c core.ComputeClient) []core.Image {
	request := core.ListImagesRequest{
		CompartmentId: helpers.CompartmentID(),
	}

	r, err := c.ListImages(ctx, request)
	helpers.FatalIfError(err)

	return r.Items
}

// ListShapes Lists the shapes that can be used to launch an instance within the specified compartment.
func listShapes(ctx context.Context, c core.ComputeClient, imageID *string) []core.Shape {
	request := core.ListShapesRequest{
		CompartmentId: helpers.CompartmentID(),
		ImageId:       imageID,
	}

	r, err := c.ListShapes(ctx, request)
	helpers.FatalIfError(err)

	if r.Items == nil || len(r.Items) == 0 {
		log.Fatalln("Invalid response from ListShapes")
	}

	return r.Items
}

func terminateInstance(ctx context.Context, c core.ComputeClient, id *string) {
	request := core.TerminateInstanceRequest{
		InstanceId:      id,
		RequestMetadata: helpers.GetRequestMetadataWithDefaultRetryPolicy(),
	}

	_, err := c.TerminateInstance(ctx, request)
	helpers.FatalIfError(err)

	fmt.Println("terminating instance")

	// should retry condition check which returns a bool value indicating whether to do retry or not
	// it checks the lifecycle status equals to Terminated or not for this case
	shouldRetryFunc := func(r common.OCIOperationResponse) bool {
		if converted, ok := r.Response.(core.GetInstanceResponse); ok {
			return converted.LifecycleState != core.InstanceLifecycleStateTerminated
		}
		return true
	}

	pollGetRequest := core.GetInstanceRequest{
		InstanceId:      id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(shouldRetryFunc),
	}

	_, pollErr := c.GetInstance(ctx, pollGetRequest)
	helpers.FatalIfError(pollErr)
	fmt.Println("instance terminated")
}

func deleteVcn(ctx context.Context, c core.VirtualNetworkClient, id *string) {
	request := core.DeleteVcnRequest{
		VcnId:           id,
		RequestMetadata: helpers.GetRequestMetadataWithDefaultRetryPolicy(),
	}

	fmt.Println("deleteing VCN")
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
		VcnId:           id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(shouldRetryFunc),
	}

	_, pollErr := c.GetVcn(ctx, pollGetRequest)
	if serviceError, ok := common.IsServiceError(pollErr); !ok ||
		(ok && serviceError.GetHTTPStatusCode() != 404) {
		// fail if the error is not service error or
		// if the error is service error and status code not equals to 404
		helpers.FatalIfError(pollErr)
	}
	fmt.Println("VCN deleted")
}

func deleteSubnet(ctx context.Context, c core.VirtualNetworkClient, id *string) {
	request := core.DeleteSubnetRequest{
		SubnetId:        id,
		RequestMetadata: helpers.GetRequestMetadataWithDefaultRetryPolicy(),
	}

	_, err := c.DeleteSubnet(context.Background(), request)
	helpers.FatalIfError(err)

	fmt.Println("deleteing subnet")

	// should retry condition check which returns a bool value indicating whether to do retry or not
	// it checks the lifecycle status equals to Terminated or not for this case
	shouldRetryFunc := func(r common.OCIOperationResponse) bool {
		if serviceError, ok := common.IsServiceError(r.Error); ok && serviceError.GetHTTPStatusCode() == 404 {
			// resource been deleted
			return false
		}

		if converted, ok := r.Response.(core.GetSubnetResponse); ok {
			return converted.LifecycleState != core.SubnetLifecycleStateTerminated
		}
		return true
	}

	pollGetRequest := core.GetSubnetRequest{
		SubnetId:        id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(shouldRetryFunc),
	}

	_, pollErr := c.GetSubnet(ctx, pollGetRequest)
	if serviceError, ok := common.IsServiceError(pollErr); !ok ||
		(ok && serviceError.GetHTTPStatusCode() != 404) {
		// fail if the error is not service error or
		// if the error is service error and status code not equals to 404
		helpers.FatalIfError(pollErr)
	}

	fmt.Println("subnet deleted")
}
