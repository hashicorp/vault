// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Key Management Service API
//
// API for managing and performing operations with keys and vaults.
//

package keymanagement

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//KmsVaultClient a client for KmsVault
type KmsVaultClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewKmsVaultClientWithConfigurationProvider Creates a new default KmsVault client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewKmsVaultClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client KmsVaultClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = KmsVaultClient{BaseClient: baseClient}
	client.BasePath = ""
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *KmsVaultClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).Endpoint("kms")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *KmsVaultClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
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
func (client *KmsVaultClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// CancelVaultDeletion Cancels the scheduled deletion of the specified vault. Canceling a scheduled deletion
// restores the vault and all keys in it to the respective states they were in before
// the deletion was scheduled. All the keys that have already been scheduled deletion before the
// scheduled deletion of the vault will also remain in their state and timeOfDeletion.
func (client KmsVaultClient) CancelVaultDeletion(ctx context.Context, request CancelVaultDeletionRequest) (response CancelVaultDeletionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.cancelVaultDeletion, policy)
	if err != nil {
		if ociResponse != nil {
			response = CancelVaultDeletionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CancelVaultDeletionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CancelVaultDeletionResponse")
	}
	return
}

// cancelVaultDeletion implements the OCIOperation interface (enables retrying operations)
func (client KmsVaultClient) cancelVaultDeletion(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/vaults/{vaultId}/actions/cancelDeletion")
	if err != nil {
		return nil, err
	}

	var response CancelVaultDeletionResponse
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

// ChangeVaultCompartment Moves a vault into a different compartment. When provided, If-Match is checked against ETag values of the resource.
func (client KmsVaultClient) ChangeVaultCompartment(ctx context.Context, request ChangeVaultCompartmentRequest) (response ChangeVaultCompartmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.changeVaultCompartment, policy)
	if err != nil {
		if ociResponse != nil {
			response = ChangeVaultCompartmentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ChangeVaultCompartmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ChangeVaultCompartmentResponse")
	}
	return
}

// changeVaultCompartment implements the OCIOperation interface (enables retrying operations)
func (client KmsVaultClient) changeVaultCompartment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/vaults/{vaultId}/actions/changeCompartment")
	if err != nil {
		return nil, err
	}

	var response ChangeVaultCompartmentResponse
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

// CreateVault Creates a new vault. The type of vault you create determines key
// placement, pricing, and available options. Options include storage
// isolation, a dedicated service endpoint instead of a shared service
// endpoint for API calls, and a dedicated hardware security module (HSM) or a multitenant HSM.
func (client KmsVaultClient) CreateVault(ctx context.Context, request CreateVaultRequest) (response CreateVaultResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createVault, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateVaultResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateVaultResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateVaultResponse")
	}
	return
}

// createVault implements the OCIOperation interface (enables retrying operations)
func (client KmsVaultClient) createVault(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/vaults")
	if err != nil {
		return nil, err
	}

	var response CreateVaultResponse
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

// GetVault Gets the specified vault's configuration information.
func (client KmsVaultClient) GetVault(ctx context.Context, request GetVaultRequest) (response GetVaultResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getVault, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetVaultResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetVaultResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetVaultResponse")
	}
	return
}

// getVault implements the OCIOperation interface (enables retrying operations)
func (client KmsVaultClient) getVault(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/20180608/vaults/{vaultId}")
	if err != nil {
		return nil, err
	}

	var response GetVaultResponse
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

// ListVaults Lists the vaults in the specified compartment.
func (client KmsVaultClient) ListVaults(ctx context.Context, request ListVaultsRequest) (response ListVaultsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listVaults, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListVaultsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListVaultsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListVaultsResponse")
	}
	return
}

// listVaults implements the OCIOperation interface (enables retrying operations)
func (client KmsVaultClient) listVaults(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/20180608/vaults")
	if err != nil {
		return nil, err
	}

	var response ListVaultsResponse
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

// ScheduleVaultDeletion Schedules the deletion of the specified vault. This sets the state of the vault and
// keys that are not scheduled deletion in it to `PENDING_DELETION` and then deletes them
// after the retention period ends.
// The state and the timeOfDeletion of the keys that have already been scheduled for deletion
// will not change. If any keys in it are scheduled for deletion after the specified timeOfDeletion
// for the vault, the call will be rejected with status code 409.
func (client KmsVaultClient) ScheduleVaultDeletion(ctx context.Context, request ScheduleVaultDeletionRequest) (response ScheduleVaultDeletionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.scheduleVaultDeletion, policy)
	if err != nil {
		if ociResponse != nil {
			response = ScheduleVaultDeletionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ScheduleVaultDeletionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ScheduleVaultDeletionResponse")
	}
	return
}

// scheduleVaultDeletion implements the OCIOperation interface (enables retrying operations)
func (client KmsVaultClient) scheduleVaultDeletion(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/vaults/{vaultId}/actions/scheduleDeletion")
	if err != nil {
		return nil, err
	}

	var response ScheduleVaultDeletionResponse
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

// UpdateVault Updates the properties of a vault. Specifically, you can update the
// `displayName`, `freeformTags`, and `definedTags` properties. Furthermore,
// the vault must be in an `ACTIVE` or `CREATING` state to be updated.
func (client KmsVaultClient) UpdateVault(ctx context.Context, request UpdateVaultRequest) (response UpdateVaultResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateVault, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateVaultResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateVaultResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateVaultResponse")
	}
	return
}

// updateVault implements the OCIOperation interface (enables retrying operations)
func (client KmsVaultClient) updateVault(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/20180608/vaults/{vaultId}")
	if err != nil {
		return nil, err
	}

	var response UpdateVaultResponse
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
