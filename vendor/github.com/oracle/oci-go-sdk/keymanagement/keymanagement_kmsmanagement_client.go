// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
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

//KmsManagementClient a client for KmsManagement
type KmsManagementClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewKmsManagementClientWithConfigurationProvider Creates a new default KmsManagement client with the given configuration provider.
// the configuration provider will be used for the default signer
func NewKmsManagementClientWithConfigurationProvider(configProvider common.ConfigurationProvider, endpoint string) (client KmsManagementClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	return newKmsManagementClientFromBaseClient(baseClient, configProvider, endpoint)
}

// NewKmsManagementClientWithOboToken Creates a new default KmsManagement client with the given configuration provider.
// The obotoken will be added to default headers and signed; the configuration provider will be used for the signer
//
func NewKmsManagementClientWithOboToken(configProvider common.ConfigurationProvider, oboToken string, endpoint string) (client KmsManagementClient, err error) {
	baseClient, err := common.NewClientWithOboToken(configProvider, oboToken)
	if err != nil {
		return
	}

	return newKmsManagementClientFromBaseClient(baseClient, configProvider, endpoint)
}

func newKmsManagementClientFromBaseClient(baseClient common.BaseClient, configProvider common.ConfigurationProvider, endpoint string) (client KmsManagementClient, err error) {
	client = KmsManagementClient{BaseClient: baseClient}
	client.BasePath = ""
	client.Host = endpoint
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *KmsManagementClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
	if ok, err := common.IsConfigurationProviderValid(configProvider); !ok {
		return err
	}

	client.config = &configProvider
	return nil
}

// ConfigurationProvider the ConfigurationProvider used in this client, or null if none set
func (client *KmsManagementClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// BackupKey Backs up an encrypted file that contains all key versions and metadata of the specified key so that you can restore
// the key later. The file also contains the metadata of the vault that the key belonged to.
func (client KmsManagementClient) BackupKey(ctx context.Context, request BackupKeyRequest) (response BackupKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.backupKey, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = BackupKeyResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = BackupKeyResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(BackupKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into BackupKeyResponse")
	}
	return
}

// backupKey implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) backupKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/{keyId}/actions/backup")
	if err != nil {
		return nil, err
	}

	var response BackupKeyResponse
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

// CancelKeyDeletion Cancels the scheduled deletion of the specified key. Canceling
// a scheduled deletion restores the key's lifecycle state to what
// it was before its scheduled deletion.
// As a provisioning operation, this call is subject to a Key Management limit that applies to
// the total number of requests across all provisioning write operations. Key Management might
// throttle this call to reject an otherwise valid request when the total rate of provisioning
// write operations exceeds 10 requests per second for a given tenancy.
func (client KmsManagementClient) CancelKeyDeletion(ctx context.Context, request CancelKeyDeletionRequest) (response CancelKeyDeletionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.cancelKeyDeletion, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = CancelKeyDeletionResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = CancelKeyDeletionResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CancelKeyDeletionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CancelKeyDeletionResponse")
	}
	return
}

// cancelKeyDeletion implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) cancelKeyDeletion(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/{keyId}/actions/cancelDeletion")
	if err != nil {
		return nil, err
	}

	var response CancelKeyDeletionResponse
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

// CancelKeyVersionDeletion Cancels the scheduled deletion of the specified key version. Canceling
// a scheduled deletion restores the key version to its lifecycle state from
// before its scheduled deletion.
// As a provisioning operation, this call is subject to a Key Management limit that applies to
// the total number of requests across all provisioning write operations. Key Management might
// throttle this call to reject an otherwise valid request when the total rate of provisioning
// write operations exceeds 10 requests per second for a given tenancy.
func (client KmsManagementClient) CancelKeyVersionDeletion(ctx context.Context, request CancelKeyVersionDeletionRequest) (response CancelKeyVersionDeletionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.cancelKeyVersionDeletion, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = CancelKeyVersionDeletionResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = CancelKeyVersionDeletionResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CancelKeyVersionDeletionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CancelKeyVersionDeletionResponse")
	}
	return
}

// cancelKeyVersionDeletion implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) cancelKeyVersionDeletion(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/{keyId}/keyVersions/{keyVersionId}/actions/cancelDeletion")
	if err != nil {
		return nil, err
	}

	var response CancelKeyVersionDeletionResponse
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

// ChangeKeyCompartment Moves a key into a different compartment within the same tenancy. For information about
// moving resources between compartments, see Moving Resources to a Different Compartment (https://docs.cloud.oracle.com/iaas/Content/Identity/Tasks/managingcompartments.htm#moveRes).
// When provided, if-match is checked against the ETag values of the key.
// As a provisioning operation, this call is subject to a Key Management limit that applies to
// the total number of requests across all provisioning write operations. Key Management might
// throttle this call to reject an otherwise valid request when the total rate of provisioning
// write operations exceeds 10 requests per second for a given tenancy.
func (client KmsManagementClient) ChangeKeyCompartment(ctx context.Context, request ChangeKeyCompartmentRequest) (response ChangeKeyCompartmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.changeKeyCompartment, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = ChangeKeyCompartmentResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = ChangeKeyCompartmentResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ChangeKeyCompartmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ChangeKeyCompartmentResponse")
	}
	return
}

// changeKeyCompartment implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) changeKeyCompartment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/{keyId}/actions/changeCompartment")
	if err != nil {
		return nil, err
	}

	var response ChangeKeyCompartmentResponse
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

// CreateKey Creates a new master encryption key.
// As a management operation, this call is subject to a Key Management limit that applies to the total
// number of requests across all management write operations. Key Management might throttle this call
// to reject an otherwise valid request when the total rate of management write operations exceeds 10
// requests per second for a given tenancy.
func (client KmsManagementClient) CreateKey(ctx context.Context, request CreateKeyRequest) (response CreateKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createKey, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = CreateKeyResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = CreateKeyResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateKeyResponse")
	}
	return
}

// createKey implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) createKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys")
	if err != nil {
		return nil, err
	}

	var response CreateKeyResponse
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

// CreateKeyVersion Generates a new KeyVersion (https://docs.cloud.oracle.com/api/#/en/key/release/KeyVersion/) resource that provides new cryptographic
// material for a master encryption key. The key must be in an `ENABLED` state to be rotated.
// As a management operation, this call is subject to a Key Management limit that applies to the total number
// of requests across all  management write operations. Key Management might throttle this call to reject an
// otherwise valid request when the total rate of management write operations exceeds 10 requests per second
// for a given tenancy.
func (client KmsManagementClient) CreateKeyVersion(ctx context.Context, request CreateKeyVersionRequest) (response CreateKeyVersionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createKeyVersion, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = CreateKeyVersionResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = CreateKeyVersionResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateKeyVersionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateKeyVersionResponse")
	}
	return
}

// createKeyVersion implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) createKeyVersion(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/{keyId}/keyVersions")
	if err != nil {
		return nil, err
	}

	var response CreateKeyVersionResponse
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

// DisableKey Disables a master encryption key so it can no longer be used for encryption, decryption, or
// generating new data encryption keys.
// As a management operation, this call is subject to a Key Management limit that applies to the total number
// of requests across all management write operations. Key Management might throttle this call to reject an
// otherwise valid request when the total rate of management write operations exceeds 10 requests per second
// for a given tenancy.
func (client KmsManagementClient) DisableKey(ctx context.Context, request DisableKeyRequest) (response DisableKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.disableKey, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = DisableKeyResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = DisableKeyResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DisableKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DisableKeyResponse")
	}
	return
}

// disableKey implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) disableKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/{keyId}/actions/disable")
	if err != nil {
		return nil, err
	}

	var response DisableKeyResponse
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

// EnableKey Enables a master encryption key so it can be used for encryption, decryption, or
// generating new data encryption keys.
// As a management operation, this call is subject to a Key Management limit that applies to the total number
// of requests across all management write operations. Key Management might throttle this call to reject an
// otherwise valid request when the total rate of management write operations exceeds 10 requests per second
// for a given tenancy.
func (client KmsManagementClient) EnableKey(ctx context.Context, request EnableKeyRequest) (response EnableKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.enableKey, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = EnableKeyResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = EnableKeyResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(EnableKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into EnableKeyResponse")
	}
	return
}

// enableKey implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) enableKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/{keyId}/actions/enable")
	if err != nil {
		return nil, err
	}

	var response EnableKeyResponse
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

// GetKey Gets information about the specified master encryption key.
// As a management operation, this call is subject to a Key Management limit that applies to the total number
// of requests across all management read operations. Key Management might throttle this call to reject an
// otherwise valid request when the total rate of management read operations exceeds 10 requests per second for
// a given tenancy.
func (client KmsManagementClient) GetKey(ctx context.Context, request GetKeyRequest) (response GetKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getKey, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = GetKeyResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = GetKeyResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetKeyResponse")
	}
	return
}

// getKey implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) getKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/20180608/keys/{keyId}")
	if err != nil {
		return nil, err
	}

	var response GetKeyResponse
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

// GetKeyVersion Gets information about the specified key version.
// As a management operation, this call is subject to a Key Management limit that applies to the total number
// of requests across all management read operations. Key Management might throttle this call to reject an
// otherwise valid request when the total rate of management read operations exceeds 10 requests per second
// for a given tenancy.
func (client KmsManagementClient) GetKeyVersion(ctx context.Context, request GetKeyVersionRequest) (response GetKeyVersionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getKeyVersion, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = GetKeyVersionResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = GetKeyVersionResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetKeyVersionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetKeyVersionResponse")
	}
	return
}

// getKeyVersion implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) getKeyVersion(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/20180608/keys/{keyId}/keyVersions/{keyVersionId}")
	if err != nil {
		return nil, err
	}

	var response GetKeyVersionResponse
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

// GetWrappingKey Gets details about the public RSA wrapping key associated with the vault in the endpoint. Each vault has an RSA key-pair that wraps and
// unwraps AES key material for import into Key Management.
func (client KmsManagementClient) GetWrappingKey(ctx context.Context, request GetWrappingKeyRequest) (response GetWrappingKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getWrappingKey, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = GetWrappingKeyResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = GetWrappingKeyResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetWrappingKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetWrappingKeyResponse")
	}
	return
}

// getWrappingKey implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) getWrappingKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/20180608/wrappingKeys")
	if err != nil {
		return nil, err
	}

	var response GetWrappingKeyResponse
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

// ImportKey Imports AES key material to create a new key with. The key material must be base64-encoded and
// wrapped by the vault's public RSA wrapping key before you can import it. Key Management supports AES symmetric keys
// that are exactly 16, 24, or 32 bytes. Furthermore, the key length must match what you specify at the time of import.
func (client KmsManagementClient) ImportKey(ctx context.Context, request ImportKeyRequest) (response ImportKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.importKey, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = ImportKeyResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = ImportKeyResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ImportKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ImportKeyResponse")
	}
	return
}

// importKey implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) importKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/import")
	if err != nil {
		return nil, err
	}

	var response ImportKeyResponse
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

// ImportKeyVersion Imports AES key material to create a new key version with, and then rotates the key to begin using the new
// key version. The key material must be base64-encoded and wrapped by the vault's public RSA wrapping key
// before you can import it. Key Management supports AES symmetric keys that are exactly 16, 24, or 32 bytes.
// Furthermore, the key length must match the length of the specified key and what you specify as the length
// at the time of import.
func (client KmsManagementClient) ImportKeyVersion(ctx context.Context, request ImportKeyVersionRequest) (response ImportKeyVersionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.importKeyVersion, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = ImportKeyVersionResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = ImportKeyVersionResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ImportKeyVersionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ImportKeyVersionResponse")
	}
	return
}

// importKeyVersion implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) importKeyVersion(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/{keyId}/keyVersions/import")
	if err != nil {
		return nil, err
	}

	var response ImportKeyVersionResponse
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

// ListKeyVersions Lists all KeyVersion (https://docs.cloud.oracle.com/api/#/en/key/release/KeyVersion/) resources for the specified
// master encryption key.
// As a management operation, this call is subject to a Key Management limit that applies to the total number
// of requests across all management read operations. Key Management might throttle this call to reject an
// otherwise valid request when the total rate of management read operations exceeds 10 requests per second
// for a given tenancy.
func (client KmsManagementClient) ListKeyVersions(ctx context.Context, request ListKeyVersionsRequest) (response ListKeyVersionsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listKeyVersions, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = ListKeyVersionsResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = ListKeyVersionsResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListKeyVersionsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListKeyVersionsResponse")
	}
	return
}

// listKeyVersions implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) listKeyVersions(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/20180608/keys/{keyId}/keyVersions")
	if err != nil {
		return nil, err
	}

	var response ListKeyVersionsResponse
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

// ListKeys Lists the master encryption keys in the specified vault and compartment.
// As a management operation, this call is subject to a Key Management limit that applies to the total number
// of requests across all management read operations. Key Management might throttle this call to reject an
// otherwise valid request when the total rate of management read operations exceeds 10 requests per second
// for a given tenancy.
func (client KmsManagementClient) ListKeys(ctx context.Context, request ListKeysRequest) (response ListKeysResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listKeys, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = ListKeysResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = ListKeysResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListKeysResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListKeysResponse")
	}
	return
}

// listKeys implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) listKeys(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/20180608/keys")
	if err != nil {
		return nil, err
	}

	var response ListKeysResponse
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

// RestoreKeyFromFile Restores the specified key to the specified vault, based on information in the backup file provided.
// If the vault doesn't exist, the operation returns a response with a 404 HTTP status error code. You
// need to first restore the vault associated with the key.
func (client KmsManagementClient) RestoreKeyFromFile(ctx context.Context, request RestoreKeyFromFileRequest) (response RestoreKeyFromFileResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.restoreKeyFromFile, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = RestoreKeyFromFileResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = RestoreKeyFromFileResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(RestoreKeyFromFileResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into RestoreKeyFromFileResponse")
	}
	return
}

// restoreKeyFromFile implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) restoreKeyFromFile(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/actions/restoreFromFile")
	if err != nil {
		return nil, err
	}

	var response RestoreKeyFromFileResponse
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

// RestoreKeyFromObjectStore Restores the specified key to the specified vault from an Oracle Cloud Infrastructure
// Object Storage location. If the vault doesn't exist, the operation returns a response with a
// 404 HTTP status error code. You need to first restore the vault associated with the key.
func (client KmsManagementClient) RestoreKeyFromObjectStore(ctx context.Context, request RestoreKeyFromObjectStoreRequest) (response RestoreKeyFromObjectStoreResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.restoreKeyFromObjectStore, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = RestoreKeyFromObjectStoreResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = RestoreKeyFromObjectStoreResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(RestoreKeyFromObjectStoreResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into RestoreKeyFromObjectStoreResponse")
	}
	return
}

// restoreKeyFromObjectStore implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) restoreKeyFromObjectStore(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/actions/restoreFromObjectStore")
	if err != nil {
		return nil, err
	}

	var response RestoreKeyFromObjectStoreResponse
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

// ScheduleKeyDeletion Schedules the deletion of the specified key. This sets the lifecycle state of the key
// to `PENDING_DELETION` and then deletes it after the specified retention period ends.
// As a provisioning operation, this call is subject to a Key Management limit that applies to
// the total number of requests across all provisioning write operations. Key Management might
// throttle this call to reject an otherwise valid request when the total rate of provisioning
// write operations exceeds 10 requests per second for a given tenancy.
func (client KmsManagementClient) ScheduleKeyDeletion(ctx context.Context, request ScheduleKeyDeletionRequest) (response ScheduleKeyDeletionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.scheduleKeyDeletion, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = ScheduleKeyDeletionResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = ScheduleKeyDeletionResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ScheduleKeyDeletionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ScheduleKeyDeletionResponse")
	}
	return
}

// scheduleKeyDeletion implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) scheduleKeyDeletion(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/{keyId}/actions/scheduleDeletion")
	if err != nil {
		return nil, err
	}

	var response ScheduleKeyDeletionResponse
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

// ScheduleKeyVersionDeletion Schedules the deletion of the specified key version. This sets the lifecycle state of the key version
// to `PENDING_DELETION` and then deletes it after the specified retention period ends.
// As a provisioning operation, this call is subject to a Key Management limit that applies to
// the total number of requests across all provisioning write operations. Key Management might
// throttle this call to reject an otherwise valid request when the total rate of provisioning
// write operations exceeds 10 requests per second for a given tenancy.
func (client KmsManagementClient) ScheduleKeyVersionDeletion(ctx context.Context, request ScheduleKeyVersionDeletionRequest) (response ScheduleKeyVersionDeletionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.scheduleKeyVersionDeletion, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = ScheduleKeyVersionDeletionResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = ScheduleKeyVersionDeletionResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ScheduleKeyVersionDeletionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ScheduleKeyVersionDeletionResponse")
	}
	return
}

// scheduleKeyVersionDeletion implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) scheduleKeyVersionDeletion(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/keys/{keyId}/keyVersions/{keyVersionId}/actions/scheduleDeletion")
	if err != nil {
		return nil, err
	}

	var response ScheduleKeyVersionDeletionResponse
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

// UpdateKey Updates the properties of a master encryption key. Specifically, you can update the
// `displayName`, `freeformTags`, and `definedTags` properties. Furthermore,
// the key must in an ENABLED or CREATING state to be updated.
// As a management operation, this call is subject to a Key Management limit that applies to the total number
// of requests across all management write operations. Key Management might throttle this call to reject an
// otherwise valid request when the total rate of management write operations exceeds 10 requests per second
// for a given tenancy.
func (client KmsManagementClient) UpdateKey(ctx context.Context, request UpdateKeyRequest) (response UpdateKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateKey, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = UpdateKeyResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = UpdateKeyResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateKeyResponse")
	}
	return
}

// updateKey implements the OCIOperation interface (enables retrying operations)
func (client KmsManagementClient) updateKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/20180608/keys/{keyId}")
	if err != nil {
		return nil, err
	}

	var response UpdateKeyResponse
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
