// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//ComputeClient a client for Compute
type ComputeClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewComputeClientWithConfigurationProvider Creates a new default Compute client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewComputeClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client ComputeClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = ComputeClient{BaseClient: baseClient}
	client.BasePath = "20160918"
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *ComputeClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).EndpointForTemplate("iaas", "https://iaas.{region}.{secondLevelDomain}")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *ComputeClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
	if ok, err := common.IsConfigurationProviderValid(configProvider); !ok {
		return err
	}

	// Error has been checked already
	region, _ := configProvider.Region()
	client.SetRegion(region)
	client.config = &configProvider
	return nil
}

// ConfigurationProvider the ConfigurationProvider used in this client, or null if none set
func (client *ComputeClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// AttachBootVolume Attaches the specified boot volume to the specified instance.
func (client ComputeClient) AttachBootVolume(ctx context.Context, request AttachBootVolumeRequest) (response AttachBootVolumeResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.attachBootVolume, policy)
	if err != nil {
		if ociResponse != nil {
			response = AttachBootVolumeResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(AttachBootVolumeResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into AttachBootVolumeResponse")
	}
	return
}

// attachBootVolume implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) attachBootVolume(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/bootVolumeAttachments/")
	if err != nil {
		return nil, err
	}

	var response AttachBootVolumeResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// AttachVnic Creates a secondary VNIC and attaches it to the specified instance.
// For more information about secondary VNICs, see
// Virtual Network Interface Cards (VNICs) (https://docs.cloud.oracle.com/Content/Network/Tasks/managingVNICs.htm).
func (client ComputeClient) AttachVnic(ctx context.Context, request AttachVnicRequest) (response AttachVnicResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.attachVnic, policy)
	if err != nil {
		if ociResponse != nil {
			response = AttachVnicResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(AttachVnicResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into AttachVnicResponse")
	}
	return
}

// attachVnic implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) attachVnic(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/vnicAttachments/")
	if err != nil {
		return nil, err
	}

	var response AttachVnicResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// AttachVolume Attaches the specified storage volume to the specified instance.
func (client ComputeClient) AttachVolume(ctx context.Context, request AttachVolumeRequest) (response AttachVolumeResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.attachVolume, policy)
	if err != nil {
		if ociResponse != nil {
			response = AttachVolumeResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(AttachVolumeResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into AttachVolumeResponse")
	}
	return
}

// attachVolume implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) attachVolume(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/volumeAttachments/")
	if err != nil {
		return nil, err
	}

	var response AttachVolumeResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponseWithPolymorphicBody(httpResponse, &response, &volumeattachment{})
	return response, err
}

// CaptureConsoleHistory Captures the most recent serial console data (up to a megabyte) for the
// specified instance.
// The `CaptureConsoleHistory` operation works with the other console history operations
// as described below.
// 1. Use `CaptureConsoleHistory` to request the capture of up to a megabyte of the
// most recent console history. This call returns a `ConsoleHistory`
// object. The object will have a state of REQUESTED.
// 2. Wait for the capture operation to succeed by polling `GetConsoleHistory` with
// the identifier of the console history metadata. The state of the
// `ConsoleHistory` object will go from REQUESTED to GETTING-HISTORY and
// then SUCCEEDED (or FAILED).
// 3. Use `GetConsoleHistoryContent` to get the actual console history data (not the
// metadata).
// 4. Optionally, use `DeleteConsoleHistory` to delete the console history metadata
// and the console history data.
func (client ComputeClient) CaptureConsoleHistory(ctx context.Context, request CaptureConsoleHistoryRequest) (response CaptureConsoleHistoryResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.captureConsoleHistory, policy)
	if err != nil {
		if ociResponse != nil {
			response = CaptureConsoleHistoryResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CaptureConsoleHistoryResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CaptureConsoleHistoryResponse")
	}
	return
}

// captureConsoleHistory implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) captureConsoleHistory(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/instanceConsoleHistories/")
	if err != nil {
		return nil, err
	}

	var response CaptureConsoleHistoryResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ChangeImageCompartment Moves an image into a different compartment within the same tenancy. For information about moving
// resources between compartments, see
// Moving Resources to a Different Compartment (https://docs.cloud.oracle.com/iaas/Content/Identity/Tasks/managingcompartments.htm#moveRes).
func (client ComputeClient) ChangeImageCompartment(ctx context.Context, request ChangeImageCompartmentRequest) (response ChangeImageCompartmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.changeImageCompartment, policy)
	if err != nil {
		if ociResponse != nil {
			response = ChangeImageCompartmentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ChangeImageCompartmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ChangeImageCompartmentResponse")
	}
	return
}

// changeImageCompartment implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) changeImageCompartment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/images/{imageId}/actions/changeCompartment")
	if err != nil {
		return nil, err
	}

	var response ChangeImageCompartmentResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// CreateAppCatalogSubscription Create a subscription for listing resource version for a compartment. It will take some time to propagate to all regions.
func (client ComputeClient) CreateAppCatalogSubscription(ctx context.Context, request CreateAppCatalogSubscriptionRequest) (response CreateAppCatalogSubscriptionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createAppCatalogSubscription, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateAppCatalogSubscriptionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateAppCatalogSubscriptionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateAppCatalogSubscriptionResponse")
	}
	return
}

// createAppCatalogSubscription implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) createAppCatalogSubscription(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/appCatalogSubscriptions")
	if err != nil {
		return nil, err
	}

	var response CreateAppCatalogSubscriptionResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// CreateImage Creates a boot disk image for the specified instance or imports an exported image from the Oracle Cloud Infrastructure Object Storage service.
// When creating a new image, you must provide the OCID of the instance you want to use as the basis for the image, and
// the OCID of the compartment containing that instance. For more information about images,
// see Managing Custom Images (https://docs.cloud.oracle.com/Content/Compute/Tasks/managingcustomimages.htm).
// When importing an exported image from Object Storage, you specify the source information
// in ImageSourceDetails.
// When importing an image based on the namespace, bucket name, and object name,
// use ImageSourceViaObjectStorageTupleDetails.
// When importing an image based on the Object Storage URL, use
// ImageSourceViaObjectStorageUriDetails.
// See Object Storage URLs (https://docs.cloud.oracle.com/Content/Compute/Tasks/imageimportexport.htm#URLs) and Using Pre-Authenticated Requests (https://docs.cloud.oracle.com/Content/Object/Tasks/usingpreauthenticatedrequests.htm)
// for constructing URLs for image import/export.
// For more information about importing exported images, see
// Image Import/Export (https://docs.cloud.oracle.com/Content/Compute/Tasks/imageimportexport.htm).
// You may optionally specify a *display name* for the image, which is simply a friendly name or description.
// It does not have to be unique, and you can change it. See UpdateImage.
// Avoid entering confidential information.
func (client ComputeClient) CreateImage(ctx context.Context, request CreateImageRequest) (response CreateImageResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createImage, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateImageResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateImageResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateImageResponse")
	}
	return
}

// createImage implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) createImage(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/images")
	if err != nil {
		return nil, err
	}

	var response CreateImageResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// CreateInstanceConsoleConnection Creates a new console connection to the specified instance.
// Once the console connection has been created and is available,
// you connect to the console using SSH.
// For more information about console access, see Accessing the Console (https://docs.cloud.oracle.com/Content/Compute/References/serialconsole.htm).
func (client ComputeClient) CreateInstanceConsoleConnection(ctx context.Context, request CreateInstanceConsoleConnectionRequest) (response CreateInstanceConsoleConnectionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createInstanceConsoleConnection, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateInstanceConsoleConnectionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateInstanceConsoleConnectionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateInstanceConsoleConnectionResponse")
	}
	return
}

// createInstanceConsoleConnection implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) createInstanceConsoleConnection(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/instanceConsoleConnections")
	if err != nil {
		return nil, err
	}

	var response CreateInstanceConsoleConnectionResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// DeleteAppCatalogSubscription Delete a subscription for a listing resource version for a compartment.
func (client ComputeClient) DeleteAppCatalogSubscription(ctx context.Context, request DeleteAppCatalogSubscriptionRequest) (response DeleteAppCatalogSubscriptionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteAppCatalogSubscription, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteAppCatalogSubscriptionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteAppCatalogSubscriptionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteAppCatalogSubscriptionResponse")
	}
	return
}

// deleteAppCatalogSubscription implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) deleteAppCatalogSubscription(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/appCatalogSubscriptions")
	if err != nil {
		return nil, err
	}

	var response DeleteAppCatalogSubscriptionResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// DeleteConsoleHistory Deletes the specified console history metadata and the console history data.
func (client ComputeClient) DeleteConsoleHistory(ctx context.Context, request DeleteConsoleHistoryRequest) (response DeleteConsoleHistoryResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteConsoleHistory, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteConsoleHistoryResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteConsoleHistoryResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteConsoleHistoryResponse")
	}
	return
}

// deleteConsoleHistory implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) deleteConsoleHistory(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/instanceConsoleHistories/{instanceConsoleHistoryId}")
	if err != nil {
		return nil, err
	}

	var response DeleteConsoleHistoryResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// DeleteImage Deletes an image.
func (client ComputeClient) DeleteImage(ctx context.Context, request DeleteImageRequest) (response DeleteImageResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteImage, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteImageResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteImageResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteImageResponse")
	}
	return
}

// deleteImage implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) deleteImage(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/images/{imageId}")
	if err != nil {
		return nil, err
	}

	var response DeleteImageResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// DeleteInstanceConsoleConnection Deletes the specified instance console connection.
func (client ComputeClient) DeleteInstanceConsoleConnection(ctx context.Context, request DeleteInstanceConsoleConnectionRequest) (response DeleteInstanceConsoleConnectionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteInstanceConsoleConnection, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteInstanceConsoleConnectionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteInstanceConsoleConnectionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteInstanceConsoleConnectionResponse")
	}
	return
}

// deleteInstanceConsoleConnection implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) deleteInstanceConsoleConnection(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/instanceConsoleConnections/{instanceConsoleConnectionId}")
	if err != nil {
		return nil, err
	}

	var response DeleteInstanceConsoleConnectionResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// DetachBootVolume Detaches a boot volume from an instance. You must specify the OCID of the boot volume attachment.
// This is an asynchronous operation. The attachment's `lifecycleState` will change to DETACHING temporarily
// until the attachment is completely removed.
func (client ComputeClient) DetachBootVolume(ctx context.Context, request DetachBootVolumeRequest) (response DetachBootVolumeResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.detachBootVolume, policy)
	if err != nil {
		if ociResponse != nil {
			response = DetachBootVolumeResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DetachBootVolumeResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DetachBootVolumeResponse")
	}
	return
}

// detachBootVolume implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) detachBootVolume(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/bootVolumeAttachments/{bootVolumeAttachmentId}")
	if err != nil {
		return nil, err
	}

	var response DetachBootVolumeResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// DetachVnic Detaches and deletes the specified secondary VNIC.
// This operation cannot be used on the instance's primary VNIC.
// When you terminate an instance, all attached VNICs (primary
// and secondary) are automatically detached and deleted.
// **Important:** If the VNIC has a
// PrivateIp that is the
// target of a route rule (https://docs.cloud.oracle.com/Content/Network/Tasks/managingroutetables.htm#privateip),
// deleting the VNIC causes that route rule to blackhole and the traffic
// will be dropped.
func (client ComputeClient) DetachVnic(ctx context.Context, request DetachVnicRequest) (response DetachVnicResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.detachVnic, policy)
	if err != nil {
		if ociResponse != nil {
			response = DetachVnicResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DetachVnicResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DetachVnicResponse")
	}
	return
}

// detachVnic implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) detachVnic(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/vnicAttachments/{vnicAttachmentId}")
	if err != nil {
		return nil, err
	}

	var response DetachVnicResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// DetachVolume Detaches a storage volume from an instance. You must specify the OCID of the volume attachment.
// This is an asynchronous operation. The attachment's `lifecycleState` will change to DETACHING temporarily
// until the attachment is completely removed.
func (client ComputeClient) DetachVolume(ctx context.Context, request DetachVolumeRequest) (response DetachVolumeResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.detachVolume, policy)
	if err != nil {
		if ociResponse != nil {
			response = DetachVolumeResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DetachVolumeResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DetachVolumeResponse")
	}
	return
}

// detachVolume implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) detachVolume(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/volumeAttachments/{volumeAttachmentId}")
	if err != nil {
		return nil, err
	}

	var response DetachVolumeResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ExportImage Exports the specified image to the Oracle Cloud Infrastructure Object Storage service. You can use the Object Storage URL,
// or the namespace, bucket name, and object name when specifying the location to export to.
// For more information about exporting images, see Image Import/Export (https://docs.cloud.oracle.com/Content/Compute/Tasks/imageimportexport.htm).
// To perform an image export, you need write access to the Object Storage bucket for the image,
// see Let Users Write Objects to Object Storage Buckets (https://docs.cloud.oracle.com/Content/Identity/Concepts/commonpolicies.htm#Let4).
// See Object Storage URLs (https://docs.cloud.oracle.com/Content/Compute/Tasks/imageimportexport.htm#URLs) and Using Pre-Authenticated Requests (https://docs.cloud.oracle.com/Content/Object/Tasks/usingpreauthenticatedrequests.htm)
// for constructing URLs for image import/export.
func (client ComputeClient) ExportImage(ctx context.Context, request ExportImageRequest) (response ExportImageResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.exportImage, policy)
	if err != nil {
		if ociResponse != nil {
			response = ExportImageResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ExportImageResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ExportImageResponse")
	}
	return
}

// exportImage implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) exportImage(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/images/{imageId}/actions/export")
	if err != nil {
		return nil, err
	}

	var response ExportImageResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// GetAppCatalogListing Gets the specified listing.
func (client ComputeClient) GetAppCatalogListing(ctx context.Context, request GetAppCatalogListingRequest) (response GetAppCatalogListingResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getAppCatalogListing, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetAppCatalogListingResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetAppCatalogListingResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetAppCatalogListingResponse")
	}
	return
}

// getAppCatalogListing implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getAppCatalogListing(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/appCatalogListings/{listingId}")
	if err != nil {
		return nil, err
	}

	var response GetAppCatalogListingResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// GetAppCatalogListingAgreements Retrieves the agreements for a particular resource version of a listing.
func (client ComputeClient) GetAppCatalogListingAgreements(ctx context.Context, request GetAppCatalogListingAgreementsRequest) (response GetAppCatalogListingAgreementsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getAppCatalogListingAgreements, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetAppCatalogListingAgreementsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetAppCatalogListingAgreementsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetAppCatalogListingAgreementsResponse")
	}
	return
}

// getAppCatalogListingAgreements implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getAppCatalogListingAgreements(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/appCatalogListings/{listingId}/resourceVersions/{resourceVersion}/agreements")
	if err != nil {
		return nil, err
	}

	var response GetAppCatalogListingAgreementsResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// GetAppCatalogListingResourceVersion Gets the specified listing resource version.
func (client ComputeClient) GetAppCatalogListingResourceVersion(ctx context.Context, request GetAppCatalogListingResourceVersionRequest) (response GetAppCatalogListingResourceVersionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getAppCatalogListingResourceVersion, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetAppCatalogListingResourceVersionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetAppCatalogListingResourceVersionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetAppCatalogListingResourceVersionResponse")
	}
	return
}

// getAppCatalogListingResourceVersion implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getAppCatalogListingResourceVersion(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/appCatalogListings/{listingId}/resourceVersions/{resourceVersion}")
	if err != nil {
		return nil, err
	}

	var response GetAppCatalogListingResourceVersionResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// GetBootVolumeAttachment Gets information about the specified boot volume attachment.
func (client ComputeClient) GetBootVolumeAttachment(ctx context.Context, request GetBootVolumeAttachmentRequest) (response GetBootVolumeAttachmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getBootVolumeAttachment, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetBootVolumeAttachmentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetBootVolumeAttachmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetBootVolumeAttachmentResponse")
	}
	return
}

// getBootVolumeAttachment implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getBootVolumeAttachment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/bootVolumeAttachments/{bootVolumeAttachmentId}")
	if err != nil {
		return nil, err
	}

	var response GetBootVolumeAttachmentResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// GetConsoleHistory Shows the metadata for the specified console history.
// See CaptureConsoleHistory
// for details about using the console history operations.
func (client ComputeClient) GetConsoleHistory(ctx context.Context, request GetConsoleHistoryRequest) (response GetConsoleHistoryResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getConsoleHistory, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetConsoleHistoryResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetConsoleHistoryResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetConsoleHistoryResponse")
	}
	return
}

// getConsoleHistory implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getConsoleHistory(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/instanceConsoleHistories/{instanceConsoleHistoryId}")
	if err != nil {
		return nil, err
	}

	var response GetConsoleHistoryResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// GetConsoleHistoryContent Gets the actual console history data (not the metadata).
// See CaptureConsoleHistory
// for details about using the console history operations.
func (client ComputeClient) GetConsoleHistoryContent(ctx context.Context, request GetConsoleHistoryContentRequest) (response GetConsoleHistoryContentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getConsoleHistoryContent, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetConsoleHistoryContentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetConsoleHistoryContentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetConsoleHistoryContentResponse")
	}
	return
}

// getConsoleHistoryContent implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getConsoleHistoryContent(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/instanceConsoleHistories/{instanceConsoleHistoryId}/data")
	if err != nil {
		return nil, err
	}

	var response GetConsoleHistoryContentResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// GetImage Gets the specified image.
func (client ComputeClient) GetImage(ctx context.Context, request GetImageRequest) (response GetImageResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getImage, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetImageResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetImageResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetImageResponse")
	}
	return
}

// getImage implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getImage(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/images/{imageId}")
	if err != nil {
		return nil, err
	}

	var response GetImageResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// GetInstance Gets information about the specified instance.
func (client ComputeClient) GetInstance(ctx context.Context, request GetInstanceRequest) (response GetInstanceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getInstance, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetInstanceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetInstanceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetInstanceResponse")
	}
	return
}

// getInstance implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getInstance(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/instances/{instanceId}")
	if err != nil {
		return nil, err
	}

	var response GetInstanceResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// GetInstanceConsoleConnection Gets the specified instance console connection's information.
func (client ComputeClient) GetInstanceConsoleConnection(ctx context.Context, request GetInstanceConsoleConnectionRequest) (response GetInstanceConsoleConnectionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getInstanceConsoleConnection, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetInstanceConsoleConnectionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetInstanceConsoleConnectionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetInstanceConsoleConnectionResponse")
	}
	return
}

// getInstanceConsoleConnection implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getInstanceConsoleConnection(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/instanceConsoleConnections/{instanceConsoleConnectionId}")
	if err != nil {
		return nil, err
	}

	var response GetInstanceConsoleConnectionResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// GetVnicAttachment Gets the information for the specified VNIC attachment.
func (client ComputeClient) GetVnicAttachment(ctx context.Context, request GetVnicAttachmentRequest) (response GetVnicAttachmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getVnicAttachment, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetVnicAttachmentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetVnicAttachmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetVnicAttachmentResponse")
	}
	return
}

// getVnicAttachment implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getVnicAttachment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/vnicAttachments/{vnicAttachmentId}")
	if err != nil {
		return nil, err
	}

	var response GetVnicAttachmentResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// GetVolumeAttachment Gets information about the specified volume attachment.
func (client ComputeClient) GetVolumeAttachment(ctx context.Context, request GetVolumeAttachmentRequest) (response GetVolumeAttachmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getVolumeAttachment, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetVolumeAttachmentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetVolumeAttachmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetVolumeAttachmentResponse")
	}
	return
}

// getVolumeAttachment implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getVolumeAttachment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/volumeAttachments/{volumeAttachmentId}")
	if err != nil {
		return nil, err
	}

	var response GetVolumeAttachmentResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponseWithPolymorphicBody(httpResponse, &response, &volumeattachment{})
	return response, err
}

// GetWindowsInstanceInitialCredentials Gets the generated credentials for the instance. Only works for instances that require password to log in (E.g. Windows).
// For certain OS'es, users will be forced to change the initial credentials.
func (client ComputeClient) GetWindowsInstanceInitialCredentials(ctx context.Context, request GetWindowsInstanceInitialCredentialsRequest) (response GetWindowsInstanceInitialCredentialsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getWindowsInstanceInitialCredentials, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetWindowsInstanceInitialCredentialsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetWindowsInstanceInitialCredentialsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetWindowsInstanceInitialCredentialsResponse")
	}
	return
}

// getWindowsInstanceInitialCredentials implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) getWindowsInstanceInitialCredentials(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/instances/{instanceId}/initialCredentials")
	if err != nil {
		return nil, err
	}

	var response GetWindowsInstanceInitialCredentialsResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// InstanceAction Performs one of the following power actions on the specified instance:
// - **START** - Powers on the instance.
// - **STOP** - Powers off the instance.
// - **SOFTRESET** - Gracefully reboots instance by sending a shutdown command to the operating system and then powers the instance back on.
// - **SOFTSTOP** - Gracefully shuts down instance by sending a shutdown command to the operating system.
// - **RESET** - Powers off the instance and then powers it back on.
// For more information see Stopping and Starting an Instance (https://docs.cloud.oracle.com/Content/Compute/Tasks/restartinginstance.htm).
func (client ComputeClient) InstanceAction(ctx context.Context, request InstanceActionRequest) (response InstanceActionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.instanceAction, policy)
	if err != nil {
		if ociResponse != nil {
			response = InstanceActionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(InstanceActionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into InstanceActionResponse")
	}
	return
}

// instanceAction implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) instanceAction(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/instances/{instanceId}")
	if err != nil {
		return nil, err
	}

	var response InstanceActionResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// LaunchInstance Creates a new instance in the specified compartment and the specified availability domain.
// For general information about instances, see
// Overview of the Compute Service (https://docs.cloud.oracle.com/Content/Compute/Concepts/computeoverview.htm).
// For information about access control and compartments, see
// Overview of the IAM Service (https://docs.cloud.oracle.com/Content/Identity/Concepts/overview.htm).
// For information about availability domains, see
// Regions and Availability Domains (https://docs.cloud.oracle.com/Content/General/Concepts/regions.htm).
// To get a list of availability domains, use the `ListAvailabilityDomains` operation
// in the Identity and Access Management Service API.
// All Oracle Cloud Infrastructure resources, including instances, get an Oracle-assigned,
// unique ID called an Oracle Cloud Identifier (OCID).
// When you create a resource, you can find its OCID in the response. You can
// also retrieve a resource's OCID by using a List API operation
// on that resource type, or by viewing the resource in the Console.
// To launch an instance using an image or a boot volume use the `sourceDetails` parameter in LaunchInstanceDetails.
// When you launch an instance, it is automatically attached to a virtual
// network interface card (VNIC), called the *primary VNIC*. The VNIC
// has a private IP address from the subnet's CIDR. You can either assign a
// private IP address of your choice or let Oracle automatically assign one.
// You can choose whether the instance has a public IP address. To retrieve the
// addresses, use the ListVnicAttachments
// operation to get the VNIC ID for the instance, and then call
// GetVnic with the VNIC ID.
// You can later add secondary VNICs to an instance. For more information, see
// Virtual Network Interface Cards (VNICs) (https://docs.cloud.oracle.com/Content/Network/Tasks/managingVNICs.htm).
func (client ComputeClient) LaunchInstance(ctx context.Context, request LaunchInstanceRequest) (response LaunchInstanceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.launchInstance, policy)
	if err != nil {
		if ociResponse != nil {
			response = LaunchInstanceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(LaunchInstanceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into LaunchInstanceResponse")
	}
	return
}

// launchInstance implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) launchInstance(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/instances/")
	if err != nil {
		return nil, err
	}

	var response LaunchInstanceResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ListAppCatalogListingResourceVersions Gets all resource versions for a particular listing.
func (client ComputeClient) ListAppCatalogListingResourceVersions(ctx context.Context, request ListAppCatalogListingResourceVersionsRequest) (response ListAppCatalogListingResourceVersionsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listAppCatalogListingResourceVersions, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListAppCatalogListingResourceVersionsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListAppCatalogListingResourceVersionsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListAppCatalogListingResourceVersionsResponse")
	}
	return
}

// listAppCatalogListingResourceVersions implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listAppCatalogListingResourceVersions(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/appCatalogListings/{listingId}/resourceVersions")
	if err != nil {
		return nil, err
	}

	var response ListAppCatalogListingResourceVersionsResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ListAppCatalogListings Lists the published listings.
func (client ComputeClient) ListAppCatalogListings(ctx context.Context, request ListAppCatalogListingsRequest) (response ListAppCatalogListingsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listAppCatalogListings, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListAppCatalogListingsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListAppCatalogListingsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListAppCatalogListingsResponse")
	}
	return
}

// listAppCatalogListings implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listAppCatalogListings(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/appCatalogListings")
	if err != nil {
		return nil, err
	}

	var response ListAppCatalogListingsResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ListAppCatalogSubscriptions Lists subscriptions for a compartment.
func (client ComputeClient) ListAppCatalogSubscriptions(ctx context.Context, request ListAppCatalogSubscriptionsRequest) (response ListAppCatalogSubscriptionsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listAppCatalogSubscriptions, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListAppCatalogSubscriptionsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListAppCatalogSubscriptionsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListAppCatalogSubscriptionsResponse")
	}
	return
}

// listAppCatalogSubscriptions implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listAppCatalogSubscriptions(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/appCatalogSubscriptions")
	if err != nil {
		return nil, err
	}

	var response ListAppCatalogSubscriptionsResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ListBootVolumeAttachments Lists the boot volume attachments in the specified compartment. You can filter the
// list by specifying an instance OCID, boot volume OCID, or both.
func (client ComputeClient) ListBootVolumeAttachments(ctx context.Context, request ListBootVolumeAttachmentsRequest) (response ListBootVolumeAttachmentsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listBootVolumeAttachments, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListBootVolumeAttachmentsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListBootVolumeAttachmentsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListBootVolumeAttachmentsResponse")
	}
	return
}

// listBootVolumeAttachments implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listBootVolumeAttachments(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/bootVolumeAttachments/")
	if err != nil {
		return nil, err
	}

	var response ListBootVolumeAttachmentsResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ListConsoleHistories Lists the console history metadata for the specified compartment or instance.
func (client ComputeClient) ListConsoleHistories(ctx context.Context, request ListConsoleHistoriesRequest) (response ListConsoleHistoriesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listConsoleHistories, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListConsoleHistoriesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListConsoleHistoriesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListConsoleHistoriesResponse")
	}
	return
}

// listConsoleHistories implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listConsoleHistories(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/instanceConsoleHistories/")
	if err != nil {
		return nil, err
	}

	var response ListConsoleHistoriesResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ListImages Lists the available images in the specified compartment, including both
// Oracle-provided images (https://docs.cloud.oracle.com/Content/Compute/References/images.htm) and
// custom images (https://docs.cloud.oracle.com/Content/Compute/Tasks/managingcustomimages.htm) that have
// been created. The list of images returned is ordered to first show all
// Oracle-provided images, then all custom images.
// The order of images returned may change when new images are released.
func (client ComputeClient) ListImages(ctx context.Context, request ListImagesRequest) (response ListImagesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listImages, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListImagesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListImagesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListImagesResponse")
	}
	return
}

// listImages implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listImages(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/images")
	if err != nil {
		return nil, err
	}

	var response ListImagesResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ListInstanceConsoleConnections Lists the console connections for the specified compartment or instance.
// For more information about console access, see Accessing the Console (https://docs.cloud.oracle.com/Content/Compute/References/serialconsole.htm).
func (client ComputeClient) ListInstanceConsoleConnections(ctx context.Context, request ListInstanceConsoleConnectionsRequest) (response ListInstanceConsoleConnectionsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listInstanceConsoleConnections, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListInstanceConsoleConnectionsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListInstanceConsoleConnectionsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListInstanceConsoleConnectionsResponse")
	}
	return
}

// listInstanceConsoleConnections implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listInstanceConsoleConnections(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/instanceConsoleConnections")
	if err != nil {
		return nil, err
	}

	var response ListInstanceConsoleConnectionsResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ListInstanceDevices Gets a list of all the devices for given instance. You can optionally filter results by device availability.
func (client ComputeClient) ListInstanceDevices(ctx context.Context, request ListInstanceDevicesRequest) (response ListInstanceDevicesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listInstanceDevices, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListInstanceDevicesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListInstanceDevicesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListInstanceDevicesResponse")
	}
	return
}

// listInstanceDevices implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listInstanceDevices(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/instances/{instanceId}/devices")
	if err != nil {
		return nil, err
	}

	var response ListInstanceDevicesResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ListInstances Lists the instances in the specified compartment and the specified availability domain.
// You can filter the results by specifying an instance name (the list will include all the identically-named
// instances in the compartment).
func (client ComputeClient) ListInstances(ctx context.Context, request ListInstancesRequest) (response ListInstancesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listInstances, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListInstancesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListInstancesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListInstancesResponse")
	}
	return
}

// listInstances implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listInstances(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/instances/")
	if err != nil {
		return nil, err
	}

	var response ListInstancesResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ListShapes Lists the shapes that can be used to launch an instance within the specified compartment. You can
// filter the list by compatibility with a specific image.
func (client ComputeClient) ListShapes(ctx context.Context, request ListShapesRequest) (response ListShapesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listShapes, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListShapesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListShapesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListShapesResponse")
	}
	return
}

// listShapes implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listShapes(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/shapes")
	if err != nil {
		return nil, err
	}

	var response ListShapesResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// ListVnicAttachments Lists the VNIC attachments in the specified compartment. A VNIC attachment
// resides in the same compartment as the attached instance. The list can be
// filtered by instance, VNIC, or availability domain.
func (client ComputeClient) ListVnicAttachments(ctx context.Context, request ListVnicAttachmentsRequest) (response ListVnicAttachmentsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listVnicAttachments, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListVnicAttachmentsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListVnicAttachmentsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListVnicAttachmentsResponse")
	}
	return
}

// listVnicAttachments implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listVnicAttachments(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/vnicAttachments/")
	if err != nil {
		return nil, err
	}

	var response ListVnicAttachmentsResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

//listvolumeattachment allows to unmarshal list of polymorphic VolumeAttachment
type listvolumeattachment []volumeattachment

//UnmarshalPolymorphicJSON unmarshals polymorphic json list of items
func (m *listvolumeattachment) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {
	res := make([]VolumeAttachment, len(*m))
	for i, v := range *m {
		nn, err := v.UnmarshalPolymorphicJSON(v.JsonData)
		if err != nil {
			return nil, err
		}
		res[i] = nn.(VolumeAttachment)
	}
	return res, nil
}

// ListVolumeAttachments Lists the volume attachments in the specified compartment. You can filter the
// list by specifying an instance OCID, volume OCID, or both.
// Currently, the only supported volume attachment type are IScsiVolumeAttachment and
// ParavirtualizedVolumeAttachment.
func (client ComputeClient) ListVolumeAttachments(ctx context.Context, request ListVolumeAttachmentsRequest) (response ListVolumeAttachmentsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listVolumeAttachments, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListVolumeAttachmentsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListVolumeAttachmentsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListVolumeAttachmentsResponse")
	}
	return
}

// listVolumeAttachments implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) listVolumeAttachments(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/volumeAttachments/")
	if err != nil {
		return nil, err
	}

	var response ListVolumeAttachmentsResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponseWithPolymorphicBody(httpResponse, &response, &listvolumeattachment{})
	return response, err
}

// TerminateInstance Terminates the specified instance. Any attached VNICs and volumes are automatically detached
// when the instance terminates.
// To preserve the boot volume associated with the instance, specify `true` for `PreserveBootVolumeQueryParam`.
// To delete the boot volume when the instance is deleted, specify `false` or do not specify a value for `PreserveBootVolumeQueryParam`.
// This is an asynchronous operation. The instance's `lifecycleState` will change to TERMINATING temporarily
// until the instance is completely removed.
func (client ComputeClient) TerminateInstance(ctx context.Context, request TerminateInstanceRequest) (response TerminateInstanceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.terminateInstance, policy)
	if err != nil {
		if ociResponse != nil {
			response = TerminateInstanceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(TerminateInstanceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into TerminateInstanceResponse")
	}
	return
}

// terminateInstance implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) terminateInstance(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/instances/{instanceId}")
	if err != nil {
		return nil, err
	}

	var response TerminateInstanceResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// UpdateConsoleHistory Updates the specified console history metadata.
func (client ComputeClient) UpdateConsoleHistory(ctx context.Context, request UpdateConsoleHistoryRequest) (response UpdateConsoleHistoryResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateConsoleHistory, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateConsoleHistoryResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateConsoleHistoryResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateConsoleHistoryResponse")
	}
	return
}

// updateConsoleHistory implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) updateConsoleHistory(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/instanceConsoleHistories/{instanceConsoleHistoryId}")
	if err != nil {
		return nil, err
	}

	var response UpdateConsoleHistoryResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// UpdateImage Updates the display name of the image. Avoid entering confidential information.
func (client ComputeClient) UpdateImage(ctx context.Context, request UpdateImageRequest) (response UpdateImageResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.updateImage, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateImageResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateImageResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateImageResponse")
	}
	return
}

// updateImage implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) updateImage(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/images/{imageId}")
	if err != nil {
		return nil, err
	}

	var response UpdateImageResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}

// UpdateInstance Updates certain fields on the specified instance. Fields that are not provided in the
// request will not be updated. Avoid entering confidential information.
// Changes to metadata fields will be reflected in the instance metadata service (this may take
// up to a minute).
// The OCID of the instance remains the same.
func (client ComputeClient) UpdateInstance(ctx context.Context, request UpdateInstanceRequest) (response UpdateInstanceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.updateInstance, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateInstanceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateInstanceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateInstanceResponse")
	}
	return
}

// updateInstance implements the OCIOperation interface (enables retrying operations)
func (client ComputeClient) updateInstance(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/instances/{instanceId}")
	if err != nil {
		return nil, err
	}

	var response UpdateInstanceResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponse(httpResponse, &response)
	return response, err
}
