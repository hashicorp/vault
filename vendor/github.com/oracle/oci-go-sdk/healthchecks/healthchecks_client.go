// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Health Checks API
//
// API for the Health Checks service. Use this API to manage endpoint probes and monitors.
// For more information, see
// Overview of the Health Checks Service (https://docs.cloud.oracle.com/iaas/Content/HealthChecks/Concepts/healthchecks.htm).
//

package healthchecks

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//HealthChecksClient a client for HealthChecks
type HealthChecksClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewHealthChecksClientWithConfigurationProvider Creates a new default HealthChecks client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewHealthChecksClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client HealthChecksClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = HealthChecksClient{BaseClient: baseClient}
	client.BasePath = "20180501"
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *HealthChecksClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).Endpoint("healthchecks")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *HealthChecksClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
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
func (client *HealthChecksClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// CreateHttpMonitor Creates an HTTP monitor. Vantage points will be automatically selected if not specified,
// and probes will be initiated from each vantage point to each of the targets at the frequency
// specified by `intervalInSeconds`.
func (client HealthChecksClient) CreateHttpMonitor(ctx context.Context, request CreateHttpMonitorRequest) (response CreateHttpMonitorResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createHttpMonitor, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateHttpMonitorResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateHttpMonitorResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateHttpMonitorResponse")
	}
	return
}

// createHttpMonitor implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) createHttpMonitor(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/httpMonitors")
	if err != nil {
		return nil, err
	}

	var response CreateHttpMonitorResponse
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

// CreateOnDemandHttpProbe Creates an on-demand HTTP probe. The location response header contains the URL for
// fetching the probe results.
// *Note:* On-demand probe configurations are not saved.
func (client HealthChecksClient) CreateOnDemandHttpProbe(ctx context.Context, request CreateOnDemandHttpProbeRequest) (response CreateOnDemandHttpProbeResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.createOnDemandHttpProbe, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateOnDemandHttpProbeResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateOnDemandHttpProbeResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateOnDemandHttpProbeResponse")
	}
	return
}

// createOnDemandHttpProbe implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) createOnDemandHttpProbe(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/httpProbeResults")
	if err != nil {
		return nil, err
	}

	var response CreateOnDemandHttpProbeResponse
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

// CreateOnDemandPingProbe Creates an on-demand ping probe. The location response header contains the URL for
// fetching probe results.
// *Note:* The on-demand probe configuration is not saved.
func (client HealthChecksClient) CreateOnDemandPingProbe(ctx context.Context, request CreateOnDemandPingProbeRequest) (response CreateOnDemandPingProbeResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.createOnDemandPingProbe, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateOnDemandPingProbeResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateOnDemandPingProbeResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateOnDemandPingProbeResponse")
	}
	return
}

// createOnDemandPingProbe implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) createOnDemandPingProbe(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/pingProbeResults")
	if err != nil {
		return nil, err
	}

	var response CreateOnDemandPingProbeResponse
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

// CreatePingMonitor Creates a ping monitor. Vantage points will be automatically selected if not specified,
// and probes will be initiated from each vantage point to each of the targets at the frequency
// specified by `intervalInSeconds`.
func (client HealthChecksClient) CreatePingMonitor(ctx context.Context, request CreatePingMonitorRequest) (response CreatePingMonitorResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createPingMonitor, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreatePingMonitorResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreatePingMonitorResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreatePingMonitorResponse")
	}
	return
}

// createPingMonitor implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) createPingMonitor(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/pingMonitors")
	if err != nil {
		return nil, err
	}

	var response CreatePingMonitorResponse
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

// DeleteHttpMonitor Deletes the HTTP monitor and its configuration. All future probes of this
// monitor are stopped. Results associated with the monitor are not deleted.
func (client HealthChecksClient) DeleteHttpMonitor(ctx context.Context, request DeleteHttpMonitorRequest) (response DeleteHttpMonitorResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteHttpMonitor, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteHttpMonitorResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteHttpMonitorResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteHttpMonitorResponse")
	}
	return
}

// deleteHttpMonitor implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) deleteHttpMonitor(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/httpMonitors/{monitorId}")
	if err != nil {
		return nil, err
	}

	var response DeleteHttpMonitorResponse
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

// DeletePingMonitor Deletes the ping monitor and its configuration. All future probes of this
// monitor are stopped. Results associated with the monitor are not deleted.
func (client HealthChecksClient) DeletePingMonitor(ctx context.Context, request DeletePingMonitorRequest) (response DeletePingMonitorResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deletePingMonitor, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeletePingMonitorResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeletePingMonitorResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeletePingMonitorResponse")
	}
	return
}

// deletePingMonitor implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) deletePingMonitor(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/pingMonitors/{monitorId}")
	if err != nil {
		return nil, err
	}

	var response DeletePingMonitorResponse
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

// GetHttpMonitor Gets the configuration for the specified monitor.
func (client HealthChecksClient) GetHttpMonitor(ctx context.Context, request GetHttpMonitorRequest) (response GetHttpMonitorResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getHttpMonitor, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetHttpMonitorResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetHttpMonitorResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetHttpMonitorResponse")
	}
	return
}

// getHttpMonitor implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) getHttpMonitor(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/httpMonitors/{monitorId}")
	if err != nil {
		return nil, err
	}

	var response GetHttpMonitorResponse
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

// GetPingMonitor Gets the configuration for the specified ping monitor.
func (client HealthChecksClient) GetPingMonitor(ctx context.Context, request GetPingMonitorRequest) (response GetPingMonitorResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getPingMonitor, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetPingMonitorResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetPingMonitorResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetPingMonitorResponse")
	}
	return
}

// getPingMonitor implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) getPingMonitor(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/pingMonitors/{monitorId}")
	if err != nil {
		return nil, err
	}

	var response GetPingMonitorResponse
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

// ListHealthChecksVantagePoints Gets information about all vantage points available to the user.
func (client HealthChecksClient) ListHealthChecksVantagePoints(ctx context.Context, request ListHealthChecksVantagePointsRequest) (response ListHealthChecksVantagePointsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listHealthChecksVantagePoints, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListHealthChecksVantagePointsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListHealthChecksVantagePointsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListHealthChecksVantagePointsResponse")
	}
	return
}

// listHealthChecksVantagePoints implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) listHealthChecksVantagePoints(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/vantagePoints")
	if err != nil {
		return nil, err
	}

	var response ListHealthChecksVantagePointsResponse
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

// ListHttpMonitors Gets a list of HTTP monitors.
func (client HealthChecksClient) ListHttpMonitors(ctx context.Context, request ListHttpMonitorsRequest) (response ListHttpMonitorsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listHttpMonitors, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListHttpMonitorsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListHttpMonitorsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListHttpMonitorsResponse")
	}
	return
}

// listHttpMonitors implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) listHttpMonitors(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/httpMonitors")
	if err != nil {
		return nil, err
	}

	var response ListHttpMonitorsResponse
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

// ListHttpProbeResults Gets the HTTP probe results for the specified probe or monitor, where
// the `probeConfigurationId` is the OCID of either a monitor or an
// on-demand probe.
func (client HealthChecksClient) ListHttpProbeResults(ctx context.Context, request ListHttpProbeResultsRequest) (response ListHttpProbeResultsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listHttpProbeResults, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListHttpProbeResultsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListHttpProbeResultsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListHttpProbeResultsResponse")
	}
	return
}

// listHttpProbeResults implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) listHttpProbeResults(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/httpProbeResults/{probeConfigurationId}")
	if err != nil {
		return nil, err
	}

	var response ListHttpProbeResultsResponse
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

// ListPingMonitors Gets a list of configured ping monitors.
// Results are paginated based on `page` and `limit`.  The `opc-next-page` header provides
// a URL for fetching the next page.
func (client HealthChecksClient) ListPingMonitors(ctx context.Context, request ListPingMonitorsRequest) (response ListPingMonitorsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listPingMonitors, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListPingMonitorsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListPingMonitorsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListPingMonitorsResponse")
	}
	return
}

// listPingMonitors implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) listPingMonitors(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/pingMonitors")
	if err != nil {
		return nil, err
	}

	var response ListPingMonitorsResponse
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

// ListPingProbeResults Returns the results for the specified probe, where the `probeConfigurationId`
// is the OCID of either a monitor or an on-demand probe.
// Results are paginated based on `page` and `limit`.  The `opc-next-page` header provides
// a URL for fetching the next page.  Use `sortOrder` to set the order of the
// results.  If `sortOrder` is unspecified, results are sorted in ascending order by
// `startTime`.
func (client HealthChecksClient) ListPingProbeResults(ctx context.Context, request ListPingProbeResultsRequest) (response ListPingProbeResultsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listPingProbeResults, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListPingProbeResultsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListPingProbeResultsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListPingProbeResultsResponse")
	}
	return
}

// listPingProbeResults implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) listPingProbeResults(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/pingProbeResults/{probeConfigurationId}")
	if err != nil {
		return nil, err
	}

	var response ListPingProbeResultsResponse
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

// UpdateHttpMonitor Updates the configuration of the specified HTTP monitor. Only the fields
// specified in the request body will be updated; all other configuration
// properties will remain unchanged.
func (client HealthChecksClient) UpdateHttpMonitor(ctx context.Context, request UpdateHttpMonitorRequest) (response UpdateHttpMonitorResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateHttpMonitor, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateHttpMonitorResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateHttpMonitorResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateHttpMonitorResponse")
	}
	return
}

// updateHttpMonitor implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) updateHttpMonitor(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/httpMonitors/{monitorId}")
	if err != nil {
		return nil, err
	}

	var response UpdateHttpMonitorResponse
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

// UpdatePingMonitor Updates the configuration of the specified ping monitor. Only the fields
// specified in the request body will be updated; all other configuration properties
// will remain unchanged.
func (client HealthChecksClient) UpdatePingMonitor(ctx context.Context, request UpdatePingMonitorRequest) (response UpdatePingMonitorResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updatePingMonitor, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdatePingMonitorResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdatePingMonitorResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdatePingMonitorResponse")
	}
	return
}

// updatePingMonitor implements the OCIOperation interface (enables retrying operations)
func (client HealthChecksClient) updatePingMonitor(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/pingMonitors/{monitorId}")
	if err != nil {
		return nil, err
	}

	var response UpdatePingMonitorResponse
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
