// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Search Service API
//
// Search for resources in your cloud network.
//

package resourcesearch

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//ResourceSearchClient a client for ResourceSearch
type ResourceSearchClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewResourceSearchClientWithConfigurationProvider Creates a new default ResourceSearch client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewResourceSearchClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client ResourceSearchClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = ResourceSearchClient{BaseClient: baseClient}
	client.BasePath = "20180409"
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *ResourceSearchClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).EndpointForTemplate("query", "https://query.{region}.{secondLevelDomain}")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *ResourceSearchClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
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
func (client *ResourceSearchClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// GetResourceType Gets detailed information about a resource type by using the resource type name.
func (client ResourceSearchClient) GetResourceType(ctx context.Context, request GetResourceTypeRequest) (response GetResourceTypeResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getResourceType, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetResourceTypeResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetResourceTypeResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetResourceTypeResponse")
	}
	return
}

// getResourceType implements the OCIOperation interface (enables retrying operations)
func (client ResourceSearchClient) getResourceType(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/resourceTypes/{name}")
	if err != nil {
		return nil, err
	}

	var response GetResourceTypeResponse
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

// ListResourceTypes Lists all resource types that you can search or query for.
func (client ResourceSearchClient) ListResourceTypes(ctx context.Context, request ListResourceTypesRequest) (response ListResourceTypesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listResourceTypes, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListResourceTypesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListResourceTypesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListResourceTypesResponse")
	}
	return
}

// listResourceTypes implements the OCIOperation interface (enables retrying operations)
func (client ResourceSearchClient) listResourceTypes(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/resourceTypes")
	if err != nil {
		return nil, err
	}

	var response ListResourceTypesResponse
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

// SearchResources Queries any and all compartments in the tenancy to find resources that match the specified criteria.
// Results include resources that you have permission to view and can span different resource types.
// You can also sort results based on a specified resource attribute.
func (client ResourceSearchClient) SearchResources(ctx context.Context, request SearchResourcesRequest) (response SearchResourcesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.searchResources, policy)
	if err != nil {
		if ociResponse != nil {
			response = SearchResourcesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(SearchResourcesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into SearchResourcesResponse")
	}
	return
}

// searchResources implements the OCIOperation interface (enables retrying operations)
func (client ResourceSearchClient) searchResources(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/resources")
	if err != nil {
		return nil, err
	}

	var response SearchResourcesResponse
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
