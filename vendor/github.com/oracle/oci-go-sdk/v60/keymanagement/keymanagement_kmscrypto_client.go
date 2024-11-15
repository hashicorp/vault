// Copyright (c) 2016, 2018, 2022, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Vault Service Key Management API
//
// API for managing and performing operations with keys and vaults. (For the API for managing secrets, see the Vault Service
// Secret Management API. For the API for retrieving secrets, see the Vault Service Secret Retrieval API.)
//

package keymanagement

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/v60/common"
	"github.com/oracle/oci-go-sdk/v60/common/auth"
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
	provider, err := auth.GetGenericConfigurationProvider(configProvider)
	if err != nil {
		return client, err
	}
	baseClient, e := common.NewClientWithConfig(provider)
	if e != nil {
		return client, e
	}
	return newKmsCryptoClientFromBaseClient(baseClient, provider, endpoint)
}

// NewKmsCryptoClientWithOboToken Creates a new default KmsCrypto client with the given configuration provider.
// The obotoken will be added to default headers and signed; the configuration provider will be used for the signer
//
func NewKmsCryptoClientWithOboToken(configProvider common.ConfigurationProvider, oboToken string, endpoint string) (client KmsCryptoClient, err error) {
	baseClient, err := common.NewClientWithOboToken(configProvider, oboToken)
	if err != nil {
		return client, err
	}

	return newKmsCryptoClientFromBaseClient(baseClient, configProvider, endpoint)
}

func newKmsCryptoClientFromBaseClient(baseClient common.BaseClient, configProvider common.ConfigurationProvider, endpoint string) (client KmsCryptoClient, err error) {
	// KmsCrypto service default circuit breaker is enabled
	baseClient.Configuration.CircuitBreaker = common.NewCircuitBreaker(common.DefaultCircuitBreakerSettingWithServiceName())
	common.ConfigCircuitBreakerFromEnvVar(&baseClient)
	common.ConfigCircuitBreakerFromGlobalVar(&baseClient)

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

// Decrypt Decrypts data using the given DecryptDataDetails (https://docs.cloud.oracle.com/api/#/en/key/latest/datatypes/DecryptDataDetails) resource.
//
// See also
//
// Click https://docs.cloud.oracle.com/en-us/iaas/tools/go-sdk-examples/latest/keymanagement/Decrypt.go.html to see an example of how to use Decrypt API.
func (client KmsCryptoClient) Decrypt(ctx context.Context, request DecryptRequest) (response DecryptResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if client.RetryPolicy() != nil {
		policy = *client.RetryPolicy()
	}
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
func (client KmsCryptoClient) decrypt(ctx context.Context, request common.OCIRequest, binaryReqBody *common.OCIReadSeekCloser, extraHeaders map[string]string) (common.OCIResponse, error) {

	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/decrypt", binaryReqBody, extraHeaders)
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

// Encrypt Encrypts data using the given EncryptDataDetails (https://docs.cloud.oracle.com/api/#/en/key/latest/datatypes/EncryptDataDetails) resource.
// Plaintext included in the example request is a base64-encoded value of a UTF-8 string.
//
// See also
//
// Click https://docs.cloud.oracle.com/en-us/iaas/tools/go-sdk-examples/latest/keymanagement/Encrypt.go.html to see an example of how to use Encrypt API.
func (client KmsCryptoClient) Encrypt(ctx context.Context, request EncryptRequest) (response EncryptResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if client.RetryPolicy() != nil {
		policy = *client.RetryPolicy()
	}
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
func (client KmsCryptoClient) encrypt(ctx context.Context, request common.OCIRequest, binaryReqBody *common.OCIReadSeekCloser, extraHeaders map[string]string) (common.OCIResponse, error) {

	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/encrypt", binaryReqBody, extraHeaders)
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

// ExportKey Exports a specific version of a master encryption key according to the details of the request. For their protection,
// keys that you create and store on a hardware security module (HSM) can never leave the HSM. You can only export keys
// stored on the server. For export, the key version is encrypted by an RSA public key that you provide.
//
// See also
//
// Click https://docs.cloud.oracle.com/en-us/iaas/tools/go-sdk-examples/latest/keymanagement/ExportKey.go.html to see an example of how to use ExportKey API.
func (client KmsCryptoClient) ExportKey(ctx context.Context, request ExportKeyRequest) (response ExportKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if client.RetryPolicy() != nil {
		policy = *client.RetryPolicy()
	}
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.exportKey, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = ExportKeyResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = ExportKeyResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ExportKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ExportKeyResponse")
	}
	return
}

// exportKey implements the OCIOperation interface (enables retrying operations)
func (client KmsCryptoClient) exportKey(ctx context.Context, request common.OCIRequest, binaryReqBody *common.OCIReadSeekCloser, extraHeaders map[string]string) (common.OCIResponse, error) {

	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/exportKey", binaryReqBody, extraHeaders)
	if err != nil {
		return nil, err
	}

	var response ExportKeyResponse
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
//
// See also
//
// Click https://docs.cloud.oracle.com/en-us/iaas/tools/go-sdk-examples/latest/keymanagement/GenerateDataEncryptionKey.go.html to see an example of how to use GenerateDataEncryptionKey API.
func (client KmsCryptoClient) GenerateDataEncryptionKey(ctx context.Context, request GenerateDataEncryptionKeyRequest) (response GenerateDataEncryptionKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if client.RetryPolicy() != nil {
		policy = *client.RetryPolicy()
	}
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
func (client KmsCryptoClient) generateDataEncryptionKey(ctx context.Context, request common.OCIRequest, binaryReqBody *common.OCIReadSeekCloser, extraHeaders map[string]string) (common.OCIResponse, error) {

	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/generateDataEncryptionKey", binaryReqBody, extraHeaders)
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

// Sign Creates a digital signature for a message or message digest by using the private key of a public-private key pair,
// also known as an asymmetric key. To verify the generated signature, you can use the Verify (https://docs.cloud.oracle.com/api/#/en/key/latest/VerifiedData/Verify)
// operation. Or, if you want to validate the signature outside of the service, you can do so by using the public key of the same asymmetric key.
//
// See also
//
// Click https://docs.cloud.oracle.com/en-us/iaas/tools/go-sdk-examples/latest/keymanagement/Sign.go.html to see an example of how to use Sign API.
func (client KmsCryptoClient) Sign(ctx context.Context, request SignRequest) (response SignResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if client.RetryPolicy() != nil {
		policy = *client.RetryPolicy()
	}
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.sign, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = SignResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = SignResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(SignResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into SignResponse")
	}
	return
}

// sign implements the OCIOperation interface (enables retrying operations)
func (client KmsCryptoClient) sign(ctx context.Context, request common.OCIRequest, binaryReqBody *common.OCIReadSeekCloser, extraHeaders map[string]string) (common.OCIResponse, error) {

	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/sign", binaryReqBody, extraHeaders)
	if err != nil {
		return nil, err
	}

	var response SignResponse
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

// Verify Verifies a digital signature that was generated by the Sign (https://docs.cloud.oracle.com/api/#/en/key/latest/SignedData/Sign) operation
// by using the public key of the same asymmetric key that was used to sign the data. If you want to validate the
// digital signature outside of the service, you can do so by using the public key of the asymmetric key.
//
// See also
//
// Click https://docs.cloud.oracle.com/en-us/iaas/tools/go-sdk-examples/latest/keymanagement/Verify.go.html to see an example of how to use Verify API.
func (client KmsCryptoClient) Verify(ctx context.Context, request VerifyRequest) (response VerifyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if client.RetryPolicy() != nil {
		policy = *client.RetryPolicy()
	}
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.verify, policy)
	if err != nil {
		if ociResponse != nil {
			if httpResponse := ociResponse.HTTPResponse(); httpResponse != nil {
				opcRequestId := httpResponse.Header.Get("opc-request-id")
				response = VerifyResponse{RawResponse: httpResponse, OpcRequestId: &opcRequestId}
			} else {
				response = VerifyResponse{}
			}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(VerifyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into VerifyResponse")
	}
	return
}

// verify implements the OCIOperation interface (enables retrying operations)
func (client KmsCryptoClient) verify(ctx context.Context, request common.OCIRequest, binaryReqBody *common.OCIReadSeekCloser, extraHeaders map[string]string) (common.OCIResponse, error) {

	httpRequest, err := request.HTTPRequest(http.MethodPost, "/20180608/verify", binaryReqBody, extraHeaders)
	if err != nil {
		return nil, err
	}

	var response VerifyResponse
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
