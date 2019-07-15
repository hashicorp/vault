// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Streaming Service API
//
// The API for the Streaming Service.
//

package streaming

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//StreamAdminClient a client for StreamAdmin
type StreamAdminClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewStreamAdminClientWithConfigurationProvider Creates a new default StreamAdmin client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewStreamAdminClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client StreamAdminClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = StreamAdminClient{BaseClient: baseClient}
	client.BasePath = "20180418"
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *StreamAdminClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).EndpointForTemplate("streams", "https://streaming.{region}.oci.{secondLevelDomain}")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *StreamAdminClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
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
func (client *StreamAdminClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// CreateStream Starts the provisioning of a new stream.
// To track the progress of the provisioning, you can periodically call GetStream.
// In the response, the `lifecycleState` parameter of the Stream object tells you its current state.
func (client StreamAdminClient) CreateStream(ctx context.Context, request CreateStreamRequest) (response CreateStreamResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.createStream, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateStreamResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateStreamResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateStreamResponse")
	}
	return
}

// createStream implements the OCIOperation interface (enables retrying operations)
func (client StreamAdminClient) createStream(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/streams")
	if err != nil {
		return nil, err
	}

	var response CreateStreamResponse
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

// DeleteStream Deletes a stream and its content. Stream contents are deleted immediately. The service retains records of the stream itself for 90 days after deletion.
// The `lifecycleState` parameter of the `Stream` object changes to `DELETING` and the stream becomes inaccessible for read or write operations.
// To verify that a stream has been deleted, make a GetStream request. If the call returns the stream's
// lifecycle state as `DELETED`, then the stream has been deleted. If the call returns a "404 Not Found" error, that means all records of the
// stream have been deleted.
func (client StreamAdminClient) DeleteStream(ctx context.Context, request DeleteStreamRequest) (response DeleteStreamResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteStream, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteStreamResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteStreamResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteStreamResponse")
	}
	return
}

// deleteStream implements the OCIOperation interface (enables retrying operations)
func (client StreamAdminClient) deleteStream(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/streams/{streamId}")
	if err != nil {
		return nil, err
	}

	var response DeleteStreamResponse
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

// GetStream Gets detailed information about a stream, including the number of partitions.
func (client StreamAdminClient) GetStream(ctx context.Context, request GetStreamRequest) (response GetStreamResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getStream, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetStreamResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetStreamResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetStreamResponse")
	}
	return
}

// getStream implements the OCIOperation interface (enables retrying operations)
func (client StreamAdminClient) getStream(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/streams/{streamId}")
	if err != nil {
		return nil, err
	}

	var response GetStreamResponse
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

// ListStreams Lists the streams.
func (client StreamAdminClient) ListStreams(ctx context.Context, request ListStreamsRequest) (response ListStreamsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listStreams, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListStreamsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListStreamsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListStreamsResponse")
	}
	return
}

// listStreams implements the OCIOperation interface (enables retrying operations)
func (client StreamAdminClient) listStreams(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/streams")
	if err != nil {
		return nil, err
	}

	var response ListStreamsResponse
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

// UpdateStream Updates the tags applied to the stream.
func (client StreamAdminClient) UpdateStream(ctx context.Context, request UpdateStreamRequest) (response UpdateStreamResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateStream, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateStreamResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateStreamResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateStreamResponse")
	}
	return
}

// updateStream implements the OCIOperation interface (enables retrying operations)
func (client StreamAdminClient) updateStream(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/streams/{streamId}")
	if err != nil {
		return nil, err
	}

	var response UpdateStreamResponse
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
