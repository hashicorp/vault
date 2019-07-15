// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Email Delivery API
//
// API for the Email Delivery service. Use this API to send high-volume, application-generated
// emails. For more information, see Overview of the Email Delivery Service (https://docs.cloud.oracle.com/iaas/Content/Email/Concepts/overview.htm).
//
// **Note:** Write actions (POST, UPDATE, DELETE) may take several minutes to propagate and be reflected by the API. If a subsequent read request fails to reflect your changes, wait a few minutes and try again.
//

package email

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//EmailClient a client for Email
type EmailClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewEmailClientWithConfigurationProvider Creates a new default Email client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewEmailClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client EmailClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = EmailClient{BaseClient: baseClient}
	client.BasePath = "20170907"
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *EmailClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).EndpointForTemplate("email", "https://email.{region}.{secondLevelDomain}")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *EmailClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
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
func (client *EmailClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// ChangeSenderCompartment Moves a sender into a different compartment. When provided, If-Match is checked against ETag values of the resource.
func (client EmailClient) ChangeSenderCompartment(ctx context.Context, request ChangeSenderCompartmentRequest) (response ChangeSenderCompartmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.changeSenderCompartment, policy)
	if err != nil {
		if ociResponse != nil {
			response = ChangeSenderCompartmentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ChangeSenderCompartmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ChangeSenderCompartmentResponse")
	}
	return
}

// changeSenderCompartment implements the OCIOperation interface (enables retrying operations)
func (client EmailClient) changeSenderCompartment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/senders/{senderId}/actions/changeCompartment")
	if err != nil {
		return nil, err
	}

	var response ChangeSenderCompartmentResponse
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

// CreateSender Creates a sender for a tenancy in a given compartment.
func (client EmailClient) CreateSender(ctx context.Context, request CreateSenderRequest) (response CreateSenderResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.createSender, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateSenderResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateSenderResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateSenderResponse")
	}
	return
}

// createSender implements the OCIOperation interface (enables retrying operations)
func (client EmailClient) createSender(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/senders")
	if err != nil {
		return nil, err
	}

	var response CreateSenderResponse
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

// CreateSuppression Adds recipient email addresses to the suppression list for a tenancy.
// Addresses added to the suppression list via the API are denoted as
// "MANUAL" in the `reason` field. *Note:* All email addresses added to the
// suppression list are normalized to include only lowercase letters.
func (client EmailClient) CreateSuppression(ctx context.Context, request CreateSuppressionRequest) (response CreateSuppressionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.createSuppression, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateSuppressionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateSuppressionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateSuppressionResponse")
	}
	return
}

// createSuppression implements the OCIOperation interface (enables retrying operations)
func (client EmailClient) createSuppression(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/suppressions")
	if err != nil {
		return nil, err
	}

	var response CreateSuppressionResponse
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

// DeleteSender Deletes an approved sender for a tenancy in a given compartment for a
// provided `senderId`.
func (client EmailClient) DeleteSender(ctx context.Context, request DeleteSenderRequest) (response DeleteSenderResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteSender, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteSenderResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteSenderResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteSenderResponse")
	}
	return
}

// deleteSender implements the OCIOperation interface (enables retrying operations)
func (client EmailClient) deleteSender(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/senders/{senderId}")
	if err != nil {
		return nil, err
	}

	var response DeleteSenderResponse
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

// DeleteSuppression Removes a suppressed recipient email address from the suppression list
// for a tenancy in a given compartment for a provided `suppressionId`.
func (client EmailClient) DeleteSuppression(ctx context.Context, request DeleteSuppressionRequest) (response DeleteSuppressionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteSuppression, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteSuppressionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteSuppressionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteSuppressionResponse")
	}
	return
}

// deleteSuppression implements the OCIOperation interface (enables retrying operations)
func (client EmailClient) deleteSuppression(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/suppressions/{suppressionId}")
	if err != nil {
		return nil, err
	}

	var response DeleteSuppressionResponse
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

// GetSender Gets an approved sender for a given `senderId`.
func (client EmailClient) GetSender(ctx context.Context, request GetSenderRequest) (response GetSenderResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getSender, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetSenderResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetSenderResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetSenderResponse")
	}
	return
}

// getSender implements the OCIOperation interface (enables retrying operations)
func (client EmailClient) getSender(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/senders/{senderId}")
	if err != nil {
		return nil, err
	}

	var response GetSenderResponse
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

// GetSuppression Gets the details of a suppressed recipient email address for a given
// `suppressionId`. Each suppression is given a unique OCID.
func (client EmailClient) GetSuppression(ctx context.Context, request GetSuppressionRequest) (response GetSuppressionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getSuppression, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetSuppressionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetSuppressionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetSuppressionResponse")
	}
	return
}

// getSuppression implements the OCIOperation interface (enables retrying operations)
func (client EmailClient) getSuppression(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/suppressions/{suppressionId}")
	if err != nil {
		return nil, err
	}

	var response GetSuppressionResponse
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

// ListSenders Gets a collection of approved sender email addresses and sender IDs.
func (client EmailClient) ListSenders(ctx context.Context, request ListSendersRequest) (response ListSendersResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listSenders, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListSendersResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListSendersResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListSendersResponse")
	}
	return
}

// listSenders implements the OCIOperation interface (enables retrying operations)
func (client EmailClient) listSenders(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/senders")
	if err != nil {
		return nil, err
	}

	var response ListSendersResponse
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

// ListSuppressions Gets a list of suppressed recipient email addresses for a user. The
// `compartmentId` for suppressions must be a tenancy OCID. The returned list
// is sorted by creation time in descending order.
func (client EmailClient) ListSuppressions(ctx context.Context, request ListSuppressionsRequest) (response ListSuppressionsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listSuppressions, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListSuppressionsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListSuppressionsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListSuppressionsResponse")
	}
	return
}

// listSuppressions implements the OCIOperation interface (enables retrying operations)
func (client EmailClient) listSuppressions(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/suppressions")
	if err != nil {
		return nil, err
	}

	var response ListSuppressionsResponse
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

// UpdateSender Replaces the set of tags for a sender with the tags provided. If either freeform
// or defined tags are omitted, the tags for that set remain the same. Each set must
// include the full set of tags for the sender, partial updates are not permitted.
// For more information about tagging, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
func (client EmailClient) UpdateSender(ctx context.Context, request UpdateSenderRequest) (response UpdateSenderResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateSender, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateSenderResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateSenderResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateSenderResponse")
	}
	return
}

// updateSender implements the OCIOperation interface (enables retrying operations)
func (client EmailClient) updateSender(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/senders/{senderId}")
	if err != nil {
		return nil, err
	}

	var response UpdateSenderResponse
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
