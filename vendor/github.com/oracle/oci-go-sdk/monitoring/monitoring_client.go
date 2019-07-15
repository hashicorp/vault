// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Monitoring API
//
// Use the Monitoring API to manage metric queries and alarms for assessing the health, capacity, and performance of your cloud resources.
// For information about monitoring, see Monitoring Overview (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm).
//

package monitoring

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//MonitoringClient a client for Monitoring
type MonitoringClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewMonitoringClientWithConfigurationProvider Creates a new default Monitoring client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewMonitoringClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client MonitoringClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = MonitoringClient{BaseClient: baseClient}
	client.BasePath = "20180401"
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *MonitoringClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).Endpoint("telemetry")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *MonitoringClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
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
func (client *MonitoringClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// CreateAlarm Creates a new alarm in the specified compartment.
func (client MonitoringClient) CreateAlarm(ctx context.Context, request CreateAlarmRequest) (response CreateAlarmResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createAlarm, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateAlarmResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateAlarmResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateAlarmResponse")
	}
	return
}

// createAlarm implements the OCIOperation interface (enables retrying operations)
func (client MonitoringClient) createAlarm(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/alarms")
	if err != nil {
		return nil, err
	}

	var response CreateAlarmResponse
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

// DeleteAlarm Deletes the specified alarm.
func (client MonitoringClient) DeleteAlarm(ctx context.Context, request DeleteAlarmRequest) (response DeleteAlarmResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteAlarm, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteAlarmResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteAlarmResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteAlarmResponse")
	}
	return
}

// deleteAlarm implements the OCIOperation interface (enables retrying operations)
func (client MonitoringClient) deleteAlarm(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/alarms/{alarmId}")
	if err != nil {
		return nil, err
	}

	var response DeleteAlarmResponse
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

// GetAlarm Gets the specified alarm.
func (client MonitoringClient) GetAlarm(ctx context.Context, request GetAlarmRequest) (response GetAlarmResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getAlarm, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetAlarmResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetAlarmResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetAlarmResponse")
	}
	return
}

// getAlarm implements the OCIOperation interface (enables retrying operations)
func (client MonitoringClient) getAlarm(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/alarms/{alarmId}")
	if err != nil {
		return nil, err
	}

	var response GetAlarmResponse
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

// GetAlarmHistory Get the history of the specified alarm.
func (client MonitoringClient) GetAlarmHistory(ctx context.Context, request GetAlarmHistoryRequest) (response GetAlarmHistoryResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getAlarmHistory, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetAlarmHistoryResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetAlarmHistoryResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetAlarmHistoryResponse")
	}
	return
}

// getAlarmHistory implements the OCIOperation interface (enables retrying operations)
func (client MonitoringClient) getAlarmHistory(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/alarms/{alarmId}/history")
	if err != nil {
		return nil, err
	}

	var response GetAlarmHistoryResponse
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

// ListAlarms Lists the alarms for the specified compartment.
func (client MonitoringClient) ListAlarms(ctx context.Context, request ListAlarmsRequest) (response ListAlarmsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listAlarms, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListAlarmsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListAlarmsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListAlarmsResponse")
	}
	return
}

// listAlarms implements the OCIOperation interface (enables retrying operations)
func (client MonitoringClient) listAlarms(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/alarms")
	if err != nil {
		return nil, err
	}

	var response ListAlarmsResponse
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

// ListAlarmsStatus List the status of each alarm in the specified compartment.
func (client MonitoringClient) ListAlarmsStatus(ctx context.Context, request ListAlarmsStatusRequest) (response ListAlarmsStatusResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listAlarmsStatus, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListAlarmsStatusResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListAlarmsStatusResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListAlarmsStatusResponse")
	}
	return
}

// listAlarmsStatus implements the OCIOperation interface (enables retrying operations)
func (client MonitoringClient) listAlarmsStatus(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/alarms/status")
	if err != nil {
		return nil, err
	}

	var response ListAlarmsStatusResponse
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

// ListMetrics Returns metric definitions that match the criteria specified in the request. Compartment OCID required.
// For information about metrics, see Metrics Overview (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm#MetricsOverview).
func (client MonitoringClient) ListMetrics(ctx context.Context, request ListMetricsRequest) (response ListMetricsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listMetrics, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListMetricsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListMetricsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListMetricsResponse")
	}
	return
}

// listMetrics implements the OCIOperation interface (enables retrying operations)
func (client MonitoringClient) listMetrics(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/metrics/actions/listMetrics")
	if err != nil {
		return nil, err
	}

	var response ListMetricsResponse
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

// PostMetricData Publishes raw metric data points to the Monitoring service.
// For more information about publishing metrics, see Publishing Custom Metrics (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Tasks/publishingcustommetrics.htm).
// The endpoints for this operation differ from other Monitoring operations. Replace the string `telemetry` with `telemetry-ingestion` in the endpoint, as in the following example:
// https://telemetry-ingestion.eu-frankfurt-1.oraclecloud.com
func (client MonitoringClient) PostMetricData(ctx context.Context, request PostMetricDataRequest) (response PostMetricDataResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.postMetricData, policy)
	if err != nil {
		if ociResponse != nil {
			response = PostMetricDataResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(PostMetricDataResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into PostMetricDataResponse")
	}
	return
}

// postMetricData implements the OCIOperation interface (enables retrying operations)
func (client MonitoringClient) postMetricData(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/metrics")
	if err != nil {
		return nil, err
	}

	var response PostMetricDataResponse
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

// RemoveAlarmSuppression Removes any existing suppression for the specified alarm.
func (client MonitoringClient) RemoveAlarmSuppression(ctx context.Context, request RemoveAlarmSuppressionRequest) (response RemoveAlarmSuppressionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.removeAlarmSuppression, policy)
	if err != nil {
		if ociResponse != nil {
			response = RemoveAlarmSuppressionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(RemoveAlarmSuppressionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into RemoveAlarmSuppressionResponse")
	}
	return
}

// removeAlarmSuppression implements the OCIOperation interface (enables retrying operations)
func (client MonitoringClient) removeAlarmSuppression(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/alarms/{alarmId}/actions/removeSuppression")
	if err != nil {
		return nil, err
	}

	var response RemoveAlarmSuppressionResponse
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

// SummarizeMetricsData Returns aggregated data that match the criteria specified in the request. Compartment OCID required.
// For information on metric queries, see Building Metric Queries (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Tasks/buildingqueries.htm).
func (client MonitoringClient) SummarizeMetricsData(ctx context.Context, request SummarizeMetricsDataRequest) (response SummarizeMetricsDataResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.summarizeMetricsData, policy)
	if err != nil {
		if ociResponse != nil {
			response = SummarizeMetricsDataResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(SummarizeMetricsDataResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into SummarizeMetricsDataResponse")
	}
	return
}

// summarizeMetricsData implements the OCIOperation interface (enables retrying operations)
func (client MonitoringClient) summarizeMetricsData(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/metrics/actions/summarizeMetricsData")
	if err != nil {
		return nil, err
	}

	var response SummarizeMetricsDataResponse
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

// UpdateAlarm Updates the specified alarm.
func (client MonitoringClient) UpdateAlarm(ctx context.Context, request UpdateAlarmRequest) (response UpdateAlarmResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateAlarm, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateAlarmResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateAlarmResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateAlarmResponse")
	}
	return
}

// updateAlarm implements the OCIOperation interface (enables retrying operations)
func (client MonitoringClient) updateAlarm(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/alarms/{alarmId}")
	if err != nil {
		return nil, err
	}

	var response UpdateAlarmResponse
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
