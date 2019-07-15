// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Budgets API
//
// Use the Budgets API to manage budgets and budget alerts.
//

package budget

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//BudgetClient a client for Budget
type BudgetClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewBudgetClientWithConfigurationProvider Creates a new default Budget client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewBudgetClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client BudgetClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = BudgetClient{BaseClient: baseClient}
	client.BasePath = "20190111"
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *BudgetClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).EndpointForTemplate("budget", "https://usage.{region}.oci.{secondLevelDomain}")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *BudgetClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
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
func (client *BudgetClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// CreateAlertRule Creates a new Alert Rule.
func (client BudgetClient) CreateAlertRule(ctx context.Context, request CreateAlertRuleRequest) (response CreateAlertRuleResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createAlertRule, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateAlertRuleResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateAlertRuleResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateAlertRuleResponse")
	}
	return
}

// createAlertRule implements the OCIOperation interface (enables retrying operations)
func (client BudgetClient) createAlertRule(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/budgets/{budgetId}/alertRules")
	if err != nil {
		return nil, err
	}

	var response CreateAlertRuleResponse
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

// CreateBudget Creates a new Budget.
func (client BudgetClient) CreateBudget(ctx context.Context, request CreateBudgetRequest) (response CreateBudgetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createBudget, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateBudgetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateBudgetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateBudgetResponse")
	}
	return
}

// createBudget implements the OCIOperation interface (enables retrying operations)
func (client BudgetClient) createBudget(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/budgets")
	if err != nil {
		return nil, err
	}

	var response CreateBudgetResponse
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

// DeleteAlertRule Deletes a specified Alert Rule resource.
func (client BudgetClient) DeleteAlertRule(ctx context.Context, request DeleteAlertRuleRequest) (response DeleteAlertRuleResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteAlertRule, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteAlertRuleResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteAlertRuleResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteAlertRuleResponse")
	}
	return
}

// deleteAlertRule implements the OCIOperation interface (enables retrying operations)
func (client BudgetClient) deleteAlertRule(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/budgets/{budgetId}/alertRules/{alertRuleId}")
	if err != nil {
		return nil, err
	}

	var response DeleteAlertRuleResponse
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

// DeleteBudget Deletes a specified Budget resource
func (client BudgetClient) DeleteBudget(ctx context.Context, request DeleteBudgetRequest) (response DeleteBudgetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteBudget, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteBudgetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteBudgetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteBudgetResponse")
	}
	return
}

// deleteBudget implements the OCIOperation interface (enables retrying operations)
func (client BudgetClient) deleteBudget(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/budgets/{budgetId}")
	if err != nil {
		return nil, err
	}

	var response DeleteBudgetResponse
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

// GetAlertRule Gets an Alert Rule for a specified Budget.
func (client BudgetClient) GetAlertRule(ctx context.Context, request GetAlertRuleRequest) (response GetAlertRuleResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getAlertRule, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetAlertRuleResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetAlertRuleResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetAlertRuleResponse")
	}
	return
}

// getAlertRule implements the OCIOperation interface (enables retrying operations)
func (client BudgetClient) getAlertRule(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/budgets/{budgetId}/alertRules/{alertRuleId}")
	if err != nil {
		return nil, err
	}

	var response GetAlertRuleResponse
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

// GetBudget Gets a Budget by identifier
func (client BudgetClient) GetBudget(ctx context.Context, request GetBudgetRequest) (response GetBudgetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getBudget, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetBudgetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetBudgetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetBudgetResponse")
	}
	return
}

// getBudget implements the OCIOperation interface (enables retrying operations)
func (client BudgetClient) getBudget(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/budgets/{budgetId}")
	if err != nil {
		return nil, err
	}

	var response GetBudgetResponse
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

// ListAlertRules Returns a list of Alert Rules for a specified Budget.
func (client BudgetClient) ListAlertRules(ctx context.Context, request ListAlertRulesRequest) (response ListAlertRulesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listAlertRules, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListAlertRulesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListAlertRulesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListAlertRulesResponse")
	}
	return
}

// listAlertRules implements the OCIOperation interface (enables retrying operations)
func (client BudgetClient) listAlertRules(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/budgets/{budgetId}/alertRules")
	if err != nil {
		return nil, err
	}

	var response ListAlertRulesResponse
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

// ListBudgets Gets a list of all Budgets in a compartment.
func (client BudgetClient) ListBudgets(ctx context.Context, request ListBudgetsRequest) (response ListBudgetsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listBudgets, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListBudgetsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListBudgetsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListBudgetsResponse")
	}
	return
}

// listBudgets implements the OCIOperation interface (enables retrying operations)
func (client BudgetClient) listBudgets(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/budgets")
	if err != nil {
		return nil, err
	}

	var response ListBudgetsResponse
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

// UpdateAlertRule Update an Alert Rule for the budget identified by the OCID.
func (client BudgetClient) UpdateAlertRule(ctx context.Context, request UpdateAlertRuleRequest) (response UpdateAlertRuleResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateAlertRule, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateAlertRuleResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateAlertRuleResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateAlertRuleResponse")
	}
	return
}

// updateAlertRule implements the OCIOperation interface (enables retrying operations)
func (client BudgetClient) updateAlertRule(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/budgets/{budgetId}/alertRules/{alertRuleId}")
	if err != nil {
		return nil, err
	}

	var response UpdateAlertRuleResponse
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

// UpdateBudget Update a Budget identified by the OCID
func (client BudgetClient) UpdateBudget(ctx context.Context, request UpdateBudgetRequest) (response UpdateBudgetResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateBudget, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateBudgetResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateBudgetResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateBudgetResponse")
	}
	return
}

// updateBudget implements the OCIOperation interface (enables retrying operations)
func (client BudgetClient) updateBudget(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/budgets/{budgetId}")
	if err != nil {
		return nil, err
	}

	var response UpdateBudgetResponse
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
