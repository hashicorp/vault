// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Load Balancing API
//
// API for the Load Balancing service. Use this API to manage load balancers, backend sets, and related items. For more
// information, see Overview of Load Balancing (https://docs.cloud.oracle.com/iaas/Content/Balance/Concepts/balanceoverview.htm).
//

package loadbalancer

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//LoadBalancerClient a client for LoadBalancer
type LoadBalancerClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewLoadBalancerClientWithConfigurationProvider Creates a new default LoadBalancer client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewLoadBalancerClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client LoadBalancerClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = LoadBalancerClient{BaseClient: baseClient}
	client.BasePath = "20170115"
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *LoadBalancerClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).EndpointForTemplate("iaas", "https://iaas.{region}.{secondLevelDomain}")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *LoadBalancerClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
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
func (client *LoadBalancerClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// CreateBackend Adds a backend server to a backend set.
func (client LoadBalancerClient) CreateBackend(ctx context.Context, request CreateBackendRequest) (response CreateBackendResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createBackend, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateBackendResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateBackendResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateBackendResponse")
	}
	return
}

// createBackend implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) createBackend(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}/backends")
	if err != nil {
		return nil, err
	}

	var response CreateBackendResponse
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

// CreateBackendSet Adds a backend set to a load balancer.
func (client LoadBalancerClient) CreateBackendSet(ctx context.Context, request CreateBackendSetRequest) (response CreateBackendSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createBackendSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateBackendSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateBackendSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateBackendSetResponse")
	}
	return
}

// createBackendSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) createBackendSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/loadBalancers/{loadBalancerId}/backendSets")
	if err != nil {
		return nil, err
	}

	var response CreateBackendSetResponse
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

// CreateCertificate Creates an asynchronous request to add an SSL certificate bundle.
func (client LoadBalancerClient) CreateCertificate(ctx context.Context, request CreateCertificateRequest) (response CreateCertificateResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createCertificate, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateCertificateResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateCertificateResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateCertificateResponse")
	}
	return
}

// createCertificate implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) createCertificate(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/loadBalancers/{loadBalancerId}/certificates")
	if err != nil {
		return nil, err
	}

	var response CreateCertificateResponse
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

// CreateHostname Adds a hostname resource to the specified load balancer. For more information, see
// Managing Request Routing (https://docs.cloud.oracle.com/Content/Balance/Tasks/managingrequest.htm).
func (client LoadBalancerClient) CreateHostname(ctx context.Context, request CreateHostnameRequest) (response CreateHostnameResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createHostname, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateHostnameResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateHostnameResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateHostnameResponse")
	}
	return
}

// createHostname implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) createHostname(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/loadBalancers/{loadBalancerId}/hostnames")
	if err != nil {
		return nil, err
	}

	var response CreateHostnameResponse
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

// CreateListener Adds a listener to a load balancer.
func (client LoadBalancerClient) CreateListener(ctx context.Context, request CreateListenerRequest) (response CreateListenerResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createListener, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateListenerResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateListenerResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateListenerResponse")
	}
	return
}

// createListener implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) createListener(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/loadBalancers/{loadBalancerId}/listeners")
	if err != nil {
		return nil, err
	}

	var response CreateListenerResponse
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

// CreateLoadBalancer Creates a new load balancer in the specified compartment. For general information about load balancers,
// see Overview of the Load Balancing Service (https://docs.cloud.oracle.com/Content/Balance/Concepts/balanceoverview.htm).
// For the purposes of access control, you must provide the OCID of the compartment where you want
// the load balancer to reside. Notice that the load balancer doesn't have to be in the same compartment as the VCN
// or backend set. If you're not sure which compartment to use, put the load balancer in the same compartment as the VCN.
// For information about access control and compartments, see
// Overview of the IAM Service (https://docs.cloud.oracle.com/Content/Identity/Concepts/overview.htm).
// You must specify a display name for the load balancer. It does not have to be unique, and you can change it.
// For information about Availability Domains, see
// Regions and Availability Domains (https://docs.cloud.oracle.com/Content/General/Concepts/regions.htm).
// To get a list of Availability Domains, use the `ListAvailabilityDomains` operation
// in the Identity and Access Management Service API.
// All Oracle Cloud Infrastructure resources, including load balancers, get an Oracle-assigned,
// unique ID called an Oracle Cloud Identifier (OCID). When you create a resource, you can find its OCID
// in the response. You can also retrieve a resource's OCID by using a List API operation on that resource type,
// or by viewing the resource in the Console. Fore more information, see
// Resource Identifiers (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
// After you send your request, the new object's state will temporarily be PROVISIONING. Before using the
// object, first make sure its state has changed to RUNNING.
// When you create a load balancer, the system assigns an IP address.
// To get the IP address, use the GetLoadBalancer operation.
func (client LoadBalancerClient) CreateLoadBalancer(ctx context.Context, request CreateLoadBalancerRequest) (response CreateLoadBalancerResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createLoadBalancer, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateLoadBalancerResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateLoadBalancerResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateLoadBalancerResponse")
	}
	return
}

// createLoadBalancer implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) createLoadBalancer(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/loadBalancers")
	if err != nil {
		return nil, err
	}

	var response CreateLoadBalancerResponse
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

// CreatePathRouteSet Adds a path route set to a load balancer. For more information, see
// Managing Request Routing (https://docs.cloud.oracle.com/Content/Balance/Tasks/managingrequest.htm).
func (client LoadBalancerClient) CreatePathRouteSet(ctx context.Context, request CreatePathRouteSetRequest) (response CreatePathRouteSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createPathRouteSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreatePathRouteSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreatePathRouteSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreatePathRouteSetResponse")
	}
	return
}

// createPathRouteSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) createPathRouteSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/loadBalancers/{loadBalancerId}/pathRouteSets")
	if err != nil {
		return nil, err
	}

	var response CreatePathRouteSetResponse
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

// CreateRuleSet Creates a new rule set associated with the specified load balancer. For more information, see
// Managing Rule Sets (https://docs.cloud.oracle.com/Content/Balance/Tasks/managingrulesets.htm).
func (client LoadBalancerClient) CreateRuleSet(ctx context.Context, request CreateRuleSetRequest) (response CreateRuleSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.createRuleSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateRuleSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateRuleSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateRuleSetResponse")
	}
	return
}

// createRuleSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) createRuleSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/loadBalancers/{loadBalancerId}/ruleSets")
	if err != nil {
		return nil, err
	}

	var response CreateRuleSetResponse
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

// DeleteBackend Removes a backend server from a given load balancer and backend set.
func (client LoadBalancerClient) DeleteBackend(ctx context.Context, request DeleteBackendRequest) (response DeleteBackendResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteBackend, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteBackendResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteBackendResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteBackendResponse")
	}
	return
}

// deleteBackend implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) deleteBackend(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}/backends/{backendName}")
	if err != nil {
		return nil, err
	}

	var response DeleteBackendResponse
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

// DeleteBackendSet Deletes the specified backend set. Note that deleting a backend set removes its backend servers from the load balancer.
// Before you can delete a backend set, you must remove it from any active listeners.
func (client LoadBalancerClient) DeleteBackendSet(ctx context.Context, request DeleteBackendSetRequest) (response DeleteBackendSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteBackendSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteBackendSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteBackendSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteBackendSetResponse")
	}
	return
}

// deleteBackendSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) deleteBackendSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}")
	if err != nil {
		return nil, err
	}

	var response DeleteBackendSetResponse
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

// DeleteCertificate Deletes an SSL certificate bundle from a load balancer.
func (client LoadBalancerClient) DeleteCertificate(ctx context.Context, request DeleteCertificateRequest) (response DeleteCertificateResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteCertificate, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteCertificateResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteCertificateResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteCertificateResponse")
	}
	return
}

// deleteCertificate implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) deleteCertificate(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/loadBalancers/{loadBalancerId}/certificates/{certificateName}")
	if err != nil {
		return nil, err
	}

	var response DeleteCertificateResponse
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

// DeleteHostname Deletes a hostname resource from the specified load balancer.
func (client LoadBalancerClient) DeleteHostname(ctx context.Context, request DeleteHostnameRequest) (response DeleteHostnameResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteHostname, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteHostnameResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteHostnameResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteHostnameResponse")
	}
	return
}

// deleteHostname implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) deleteHostname(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/loadBalancers/{loadBalancerId}/hostnames/{name}")
	if err != nil {
		return nil, err
	}

	var response DeleteHostnameResponse
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

// DeleteListener Deletes a listener from a load balancer.
func (client LoadBalancerClient) DeleteListener(ctx context.Context, request DeleteListenerRequest) (response DeleteListenerResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteListener, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteListenerResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteListenerResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteListenerResponse")
	}
	return
}

// deleteListener implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) deleteListener(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/loadBalancers/{loadBalancerId}/listeners/{listenerName}")
	if err != nil {
		return nil, err
	}

	var response DeleteListenerResponse
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

// DeleteLoadBalancer Stops a load balancer and removes it from service.
func (client LoadBalancerClient) DeleteLoadBalancer(ctx context.Context, request DeleteLoadBalancerRequest) (response DeleteLoadBalancerResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteLoadBalancer, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteLoadBalancerResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteLoadBalancerResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteLoadBalancerResponse")
	}
	return
}

// deleteLoadBalancer implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) deleteLoadBalancer(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/loadBalancers/{loadBalancerId}")
	if err != nil {
		return nil, err
	}

	var response DeleteLoadBalancerResponse
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

// DeletePathRouteSet Deletes a path route set from the specified load balancer.
// To delete a path route rule from a path route set, use the
// UpdatePathRouteSet operation.
func (client LoadBalancerClient) DeletePathRouteSet(ctx context.Context, request DeletePathRouteSetRequest) (response DeletePathRouteSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deletePathRouteSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeletePathRouteSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeletePathRouteSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeletePathRouteSetResponse")
	}
	return
}

// deletePathRouteSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) deletePathRouteSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/loadBalancers/{loadBalancerId}/pathRouteSets/{pathRouteSetName}")
	if err != nil {
		return nil, err
	}

	var response DeletePathRouteSetResponse
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

// DeleteRuleSet Deletes a rule set from the specified load balancer.
// To delete a rule from a rule set, use the
// UpdateRuleSet operation.
func (client LoadBalancerClient) DeleteRuleSet(ctx context.Context, request DeleteRuleSetRequest) (response DeleteRuleSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteRuleSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteRuleSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteRuleSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteRuleSetResponse")
	}
	return
}

// deleteRuleSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) deleteRuleSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/loadBalancers/{loadBalancerId}/ruleSets/{ruleSetName}")
	if err != nil {
		return nil, err
	}

	var response DeleteRuleSetResponse
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

// GetBackend Gets the specified backend server's configuration information.
func (client LoadBalancerClient) GetBackend(ctx context.Context, request GetBackendRequest) (response GetBackendResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getBackend, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetBackendResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetBackendResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetBackendResponse")
	}
	return
}

// getBackend implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) getBackend(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}/backends/{backendName}")
	if err != nil {
		return nil, err
	}

	var response GetBackendResponse
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

// GetBackendHealth Gets the current health status of the specified backend server.
func (client LoadBalancerClient) GetBackendHealth(ctx context.Context, request GetBackendHealthRequest) (response GetBackendHealthResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getBackendHealth, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetBackendHealthResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetBackendHealthResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetBackendHealthResponse")
	}
	return
}

// getBackendHealth implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) getBackendHealth(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}/backends/{backendName}/health")
	if err != nil {
		return nil, err
	}

	var response GetBackendHealthResponse
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

// GetBackendSet Gets the specified backend set's configuration information.
func (client LoadBalancerClient) GetBackendSet(ctx context.Context, request GetBackendSetRequest) (response GetBackendSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getBackendSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetBackendSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetBackendSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetBackendSetResponse")
	}
	return
}

// getBackendSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) getBackendSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}")
	if err != nil {
		return nil, err
	}

	var response GetBackendSetResponse
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

// GetBackendSetHealth Gets the health status for the specified backend set.
func (client LoadBalancerClient) GetBackendSetHealth(ctx context.Context, request GetBackendSetHealthRequest) (response GetBackendSetHealthResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getBackendSetHealth, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetBackendSetHealthResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetBackendSetHealthResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetBackendSetHealthResponse")
	}
	return
}

// getBackendSetHealth implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) getBackendSetHealth(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}/health")
	if err != nil {
		return nil, err
	}

	var response GetBackendSetHealthResponse
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

// GetHealthChecker Gets the health check policy information for a given load balancer and backend set.
func (client LoadBalancerClient) GetHealthChecker(ctx context.Context, request GetHealthCheckerRequest) (response GetHealthCheckerResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getHealthChecker, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetHealthCheckerResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetHealthCheckerResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetHealthCheckerResponse")
	}
	return
}

// getHealthChecker implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) getHealthChecker(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}/healthChecker")
	if err != nil {
		return nil, err
	}

	var response GetHealthCheckerResponse
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

// GetHostname Gets the specified hostname resource's configuration information.
func (client LoadBalancerClient) GetHostname(ctx context.Context, request GetHostnameRequest) (response GetHostnameResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getHostname, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetHostnameResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetHostnameResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetHostnameResponse")
	}
	return
}

// getHostname implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) getHostname(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/hostnames/{name}")
	if err != nil {
		return nil, err
	}

	var response GetHostnameResponse
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

// GetLoadBalancer Gets the specified load balancer's configuration information.
func (client LoadBalancerClient) GetLoadBalancer(ctx context.Context, request GetLoadBalancerRequest) (response GetLoadBalancerResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getLoadBalancer, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetLoadBalancerResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetLoadBalancerResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetLoadBalancerResponse")
	}
	return
}

// getLoadBalancer implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) getLoadBalancer(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}")
	if err != nil {
		return nil, err
	}

	var response GetLoadBalancerResponse
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

// GetLoadBalancerHealth Gets the health status for the specified load balancer.
func (client LoadBalancerClient) GetLoadBalancerHealth(ctx context.Context, request GetLoadBalancerHealthRequest) (response GetLoadBalancerHealthResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getLoadBalancerHealth, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetLoadBalancerHealthResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetLoadBalancerHealthResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetLoadBalancerHealthResponse")
	}
	return
}

// getLoadBalancerHealth implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) getLoadBalancerHealth(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/health")
	if err != nil {
		return nil, err
	}

	var response GetLoadBalancerHealthResponse
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

// GetPathRouteSet Gets the specified path route set's configuration information.
func (client LoadBalancerClient) GetPathRouteSet(ctx context.Context, request GetPathRouteSetRequest) (response GetPathRouteSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getPathRouteSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetPathRouteSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetPathRouteSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetPathRouteSetResponse")
	}
	return
}

// getPathRouteSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) getPathRouteSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/pathRouteSets/{pathRouteSetName}")
	if err != nil {
		return nil, err
	}

	var response GetPathRouteSetResponse
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

// GetRuleSet Gets the specified set of rules.
func (client LoadBalancerClient) GetRuleSet(ctx context.Context, request GetRuleSetRequest) (response GetRuleSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getRuleSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetRuleSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetRuleSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetRuleSetResponse")
	}
	return
}

// getRuleSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) getRuleSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/ruleSets/{ruleSetName}")
	if err != nil {
		return nil, err
	}

	var response GetRuleSetResponse
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

// GetWorkRequest Gets the details of a work request.
func (client LoadBalancerClient) GetWorkRequest(ctx context.Context, request GetWorkRequestRequest) (response GetWorkRequestResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getWorkRequest, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetWorkRequestResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetWorkRequestResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetWorkRequestResponse")
	}
	return
}

// getWorkRequest implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) getWorkRequest(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancerWorkRequests/{workRequestId}")
	if err != nil {
		return nil, err
	}

	var response GetWorkRequestResponse
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

// ListBackendSets Lists all backend sets associated with a given load balancer.
func (client LoadBalancerClient) ListBackendSets(ctx context.Context, request ListBackendSetsRequest) (response ListBackendSetsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listBackendSets, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListBackendSetsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListBackendSetsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListBackendSetsResponse")
	}
	return
}

// listBackendSets implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) listBackendSets(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/backendSets")
	if err != nil {
		return nil, err
	}

	var response ListBackendSetsResponse
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

// ListBackends Lists the backend servers for a given load balancer and backend set.
func (client LoadBalancerClient) ListBackends(ctx context.Context, request ListBackendsRequest) (response ListBackendsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listBackends, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListBackendsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListBackendsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListBackendsResponse")
	}
	return
}

// listBackends implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) listBackends(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}/backends")
	if err != nil {
		return nil, err
	}

	var response ListBackendsResponse
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

// ListCertificates Lists all SSL certificates bundles associated with a given load balancer.
func (client LoadBalancerClient) ListCertificates(ctx context.Context, request ListCertificatesRequest) (response ListCertificatesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listCertificates, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListCertificatesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListCertificatesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListCertificatesResponse")
	}
	return
}

// listCertificates implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) listCertificates(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/certificates")
	if err != nil {
		return nil, err
	}

	var response ListCertificatesResponse
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

// ListHostnames Lists all hostname resources associated with the specified load balancer.
func (client LoadBalancerClient) ListHostnames(ctx context.Context, request ListHostnamesRequest) (response ListHostnamesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listHostnames, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListHostnamesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListHostnamesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListHostnamesResponse")
	}
	return
}

// listHostnames implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) listHostnames(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/hostnames")
	if err != nil {
		return nil, err
	}

	var response ListHostnamesResponse
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

// ListLoadBalancerHealths Lists the summary health statuses for all load balancers in the specified compartment.
func (client LoadBalancerClient) ListLoadBalancerHealths(ctx context.Context, request ListLoadBalancerHealthsRequest) (response ListLoadBalancerHealthsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listLoadBalancerHealths, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListLoadBalancerHealthsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListLoadBalancerHealthsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListLoadBalancerHealthsResponse")
	}
	return
}

// listLoadBalancerHealths implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) listLoadBalancerHealths(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancerHealths")
	if err != nil {
		return nil, err
	}

	var response ListLoadBalancerHealthsResponse
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

// ListLoadBalancers Lists all load balancers in the specified compartment.
func (client LoadBalancerClient) ListLoadBalancers(ctx context.Context, request ListLoadBalancersRequest) (response ListLoadBalancersResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listLoadBalancers, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListLoadBalancersResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListLoadBalancersResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListLoadBalancersResponse")
	}
	return
}

// listLoadBalancers implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) listLoadBalancers(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers")
	if err != nil {
		return nil, err
	}

	var response ListLoadBalancersResponse
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

// ListPathRouteSets Lists all path route sets associated with the specified load balancer.
func (client LoadBalancerClient) ListPathRouteSets(ctx context.Context, request ListPathRouteSetsRequest) (response ListPathRouteSetsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listPathRouteSets, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListPathRouteSetsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListPathRouteSetsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListPathRouteSetsResponse")
	}
	return
}

// listPathRouteSets implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) listPathRouteSets(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/pathRouteSets")
	if err != nil {
		return nil, err
	}

	var response ListPathRouteSetsResponse
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

// ListPolicies Lists the available load balancer policies.
func (client LoadBalancerClient) ListPolicies(ctx context.Context, request ListPoliciesRequest) (response ListPoliciesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listPolicies, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListPoliciesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListPoliciesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListPoliciesResponse")
	}
	return
}

// listPolicies implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) listPolicies(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancerPolicies")
	if err != nil {
		return nil, err
	}

	var response ListPoliciesResponse
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

// ListProtocols Lists all supported traffic protocols.
func (client LoadBalancerClient) ListProtocols(ctx context.Context, request ListProtocolsRequest) (response ListProtocolsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listProtocols, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListProtocolsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListProtocolsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListProtocolsResponse")
	}
	return
}

// listProtocols implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) listProtocols(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancerProtocols")
	if err != nil {
		return nil, err
	}

	var response ListProtocolsResponse
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

// ListRuleSets Lists all rule sets associated with the specified load balancer.
func (client LoadBalancerClient) ListRuleSets(ctx context.Context, request ListRuleSetsRequest) (response ListRuleSetsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listRuleSets, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListRuleSetsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListRuleSetsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListRuleSetsResponse")
	}
	return
}

// listRuleSets implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) listRuleSets(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/ruleSets")
	if err != nil {
		return nil, err
	}

	var response ListRuleSetsResponse
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

// ListShapes Lists the valid load balancer shapes.
func (client LoadBalancerClient) ListShapes(ctx context.Context, request ListShapesRequest) (response ListShapesResponse, err error) {
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
func (client LoadBalancerClient) listShapes(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancerShapes")
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

// ListWorkRequests Lists the work requests for a given load balancer.
func (client LoadBalancerClient) ListWorkRequests(ctx context.Context, request ListWorkRequestsRequest) (response ListWorkRequestsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listWorkRequests, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListWorkRequestsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListWorkRequestsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListWorkRequestsResponse")
	}
	return
}

// listWorkRequests implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) listWorkRequests(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/loadBalancers/{loadBalancerId}/workRequests")
	if err != nil {
		return nil, err
	}

	var response ListWorkRequestsResponse
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

// UpdateBackend Updates the configuration of a backend server within the specified backend set.
func (client LoadBalancerClient) UpdateBackend(ctx context.Context, request UpdateBackendRequest) (response UpdateBackendResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.updateBackend, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateBackendResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateBackendResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateBackendResponse")
	}
	return
}

// updateBackend implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) updateBackend(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}/backends/{backendName}")
	if err != nil {
		return nil, err
	}

	var response UpdateBackendResponse
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

// UpdateBackendSet Updates a backend set.
func (client LoadBalancerClient) UpdateBackendSet(ctx context.Context, request UpdateBackendSetRequest) (response UpdateBackendSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.updateBackendSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateBackendSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateBackendSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateBackendSetResponse")
	}
	return
}

// updateBackendSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) updateBackendSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}")
	if err != nil {
		return nil, err
	}

	var response UpdateBackendSetResponse
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

// UpdateHealthChecker Updates the health check policy for a given load balancer and backend set.
func (client LoadBalancerClient) UpdateHealthChecker(ctx context.Context, request UpdateHealthCheckerRequest) (response UpdateHealthCheckerResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.updateHealthChecker, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateHealthCheckerResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateHealthCheckerResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateHealthCheckerResponse")
	}
	return
}

// updateHealthChecker implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) updateHealthChecker(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/loadBalancers/{loadBalancerId}/backendSets/{backendSetName}/healthChecker")
	if err != nil {
		return nil, err
	}

	var response UpdateHealthCheckerResponse
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

// UpdateHostname Overwrites an existing hostname resource on the specified load balancer. Use this operation to change a
// virtual hostname.
func (client LoadBalancerClient) UpdateHostname(ctx context.Context, request UpdateHostnameRequest) (response UpdateHostnameResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateHostname, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateHostnameResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateHostnameResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateHostnameResponse")
	}
	return
}

// updateHostname implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) updateHostname(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/loadBalancers/{loadBalancerId}/hostnames/{name}")
	if err != nil {
		return nil, err
	}

	var response UpdateHostnameResponse
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

// UpdateListener Updates a listener for a given load balancer.
func (client LoadBalancerClient) UpdateListener(ctx context.Context, request UpdateListenerRequest) (response UpdateListenerResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.updateListener, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateListenerResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateListenerResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateListenerResponse")
	}
	return
}

// updateListener implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) updateListener(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/loadBalancers/{loadBalancerId}/listeners/{listenerName}")
	if err != nil {
		return nil, err
	}

	var response UpdateListenerResponse
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

// UpdateLoadBalancer Updates a load balancer's configuration.
func (client LoadBalancerClient) UpdateLoadBalancer(ctx context.Context, request UpdateLoadBalancerRequest) (response UpdateLoadBalancerResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.updateLoadBalancer, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateLoadBalancerResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateLoadBalancerResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateLoadBalancerResponse")
	}
	return
}

// updateLoadBalancer implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) updateLoadBalancer(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/loadBalancers/{loadBalancerId}")
	if err != nil {
		return nil, err
	}

	var response UpdateLoadBalancerResponse
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

// UpdateNetworkSecurityGroups Updates the network security groups to be used by a load balancer.
func (client LoadBalancerClient) UpdateNetworkSecurityGroups(ctx context.Context, request UpdateNetworkSecurityGroupsRequest) (response UpdateNetworkSecurityGroupsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.updateNetworkSecurityGroups, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateNetworkSecurityGroupsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateNetworkSecurityGroupsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateNetworkSecurityGroupsResponse")
	}
	return
}

// updateNetworkSecurityGroups implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) updateNetworkSecurityGroups(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/loadBalancers/{loadBalancerId}/networkSecurityGroups")
	if err != nil {
		return nil, err
	}

	var response UpdateNetworkSecurityGroupsResponse
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

// UpdatePathRouteSet Overwrites an existing path route set on the specified load balancer. Use this operation to add, delete, or alter
// path route rules in a path route set.
// To add a new path route rule to a path route set, the `pathRoutes` in the
// UpdatePathRouteSetDetails object must include
// both the new path route rule to add and the existing path route rules to retain.
func (client LoadBalancerClient) UpdatePathRouteSet(ctx context.Context, request UpdatePathRouteSetRequest) (response UpdatePathRouteSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.updatePathRouteSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdatePathRouteSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdatePathRouteSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdatePathRouteSetResponse")
	}
	return
}

// updatePathRouteSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) updatePathRouteSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/loadBalancers/{loadBalancerId}/pathRouteSets/{pathRouteSetName}")
	if err != nil {
		return nil, err
	}

	var response UpdatePathRouteSetResponse
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

// UpdateRuleSet Overwrites an existing set of rules on the specified load balancer. Use this operation to add or alter
// the rules in a rule set.
// To add a new rule to a set, the body must include both the new rule to add and the existing rules to retain.
func (client LoadBalancerClient) UpdateRuleSet(ctx context.Context, request UpdateRuleSetRequest) (response UpdateRuleSetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateRuleSet, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateRuleSetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateRuleSetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateRuleSetResponse")
	}
	return
}

// updateRuleSet implements the OCIOperation interface (enables retrying operations)
func (client LoadBalancerClient) updateRuleSet(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/loadBalancers/{loadBalancerId}/ruleSets/{ruleSetName}")
	if err != nil {
		return nil, err
	}

	var response UpdateRuleSetResponse
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
