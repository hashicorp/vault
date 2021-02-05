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

//KmsCryptoClient a client for KmsCrypto
type KmsCryptoClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewKmsCryptoClientWithConfigurationProvider Creates a new default KmsCrypto client with the given configuration provider.
// the configuration provider will be used for the default signer
func NewKmsCryptoClientWithConfigurationProvider(configProvider common.ConfigurationProvider, endpoint string) (client KmsCryptoClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	return newKmsCryptoClientFromBaseClient(baseClient, configProvider, endpoint)
}

// NewKmsCryptoClientWithOboToken Creates a new default KmsCrypto client with the given configuration provider.
// The obotoken will be added to default headers and signed; the configuration provider will be used for the signer
//
func NewKmsCryptoClientWithOboToken(configProvider common.ConfigurationProvider, oboToken string, endpoint string) (client KmsCryptoClient, err error) {
	baseClient, err := common.NewClientWithOboToken(configProvider, oboToken)
	if err != nil {
		return
	}

	return newKmsCryptoClientFromBaseClient(baseClient, configProvider, endpoint)
}

func newKmsCryptoClientFromBaseClient(baseClient common.BaseClient, configProvider common.ConfigurationProvider, endpoint string) (client KmsCryptoClient, err error) {
	client = KmsCryptoClient{BaseClient: baseClient}
	client.BasePath = ""
	client.Host = endpoint
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *KmsCryptoClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
	if ok, err := common.IsConfigurationProviderValid(configProvider); !ok {
		return err
	}

	client.config = &configProvider
	return nil
}

// ConfigurationProvider the ConfigurationProvider used in this client, or null if none set
func (client *KmsCryptoClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// Decrypt Decrypts data using the given DecryptDataDetails (https://docs.cloud.oracle.com/api/#/en/key/release/datatypes/DecryptDataDetails) resource.
func (client KmsCryptoClient) Decrypt(ctx context.Context, request DecryptRequest) (response DecryptResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.decrypt, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = DecryptResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = DecryptResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DecryptResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DecryptResponse")
	}
	return
}

// decrypt implements the OCIOperation interface (enables retrying operations)
func (client KmsCryptoClient) decrypt(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/decrypt")
	if err != nil {
		return nil, err
	}

	var response DecryptResponse
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

// Encrypt Encrypts data using the given EncryptDataDetails (https://docs.cloud.oracle.com/api/#/en/key/release/datatypes/EncryptDataDetails) resource.
// Plaintext included in the example request is a base64-encoded value of a UTF-8 string.
func (client KmsCryptoClient) Encrypt(ctx context.Context, request EncryptRequest) (response EncryptResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.encrypt, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = EncryptResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = EncryptResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(EncryptResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into EncryptResponse")
	}
	return
}

// encrypt implements the OCIOperation interface (enables retrying operations)
func (client KmsCryptoClient) encrypt(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/encrypt")
	if err != nil {
		return nil, err
	}

	var response EncryptResponse
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

// GenerateDataEncryptionKey Generates a key that you can use to encrypt or decrypt data.
func (client KmsCryptoClient) GenerateDataEncryptionKey(ctx context.Context, request GenerateDataEncryptionKeyRequest) (response GenerateDataEncryptionKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.generateDataEncryptionKey, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = GenerateDataEncryptionKeyResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = GenerateDataEncryptionKeyResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GenerateDataEncryptionKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GenerateDataEncryptionKeyResponse")
	}
	return
}

// generateDataEncryptionKey implements the OCIOperation interface (enables retrying operations)
func (client KmsCryptoClient) generateDataEncryptionKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/generateDataEncryptionKey")
	if err != nil {
		return nil, err
	}

	var response GenerateDataEncryptionKeyResponse
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
