// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Notification API
//
// Use the Notification API to broadcast messages to distributed components by topic, using a publish-subscribe pattern.
// For information about managing topics, subscriptions, and messages, see Notification Overview (https://docs.cloud.oracle.com/iaas/Content/Notification/Concepts/notificationoverview.htm).
//

package ons

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//NotificationDataPlaneClient a client for NotificationDataPlane
type NotificationDataPlaneClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewNotificationDataPlaneClientWithConfigurationProvider Creates a new default NotificationDataPlane client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewNotificationDataPlaneClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client NotificationDataPlaneClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = NotificationDataPlaneClient{BaseClient: baseClient}
	client.BasePath = "20181201"
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *NotificationDataPlaneClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).Endpoint("notification")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *NotificationDataPlaneClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
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
func (client *NotificationDataPlaneClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// CreateSubscription Creates a subscription for the specified topic.
func (client NotificationDataPlaneClient) CreateSubscription(ctx context.Context, request CreateSubscriptionRequest) (response CreateSubscriptionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createSubscription, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateSubscriptionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateSubscriptionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateSubscriptionResponse")
	}
	return
}

// createSubscription implements the OCIOperation interface (enables retrying operations)
func (client NotificationDataPlaneClient) createSubscription(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/subscriptions")
	if err != nil {
		return nil, err
	}

	var response CreateSubscriptionResponse
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

// DeleteSubscription Deletes the specified subscription.
func (client NotificationDataPlaneClient) DeleteSubscription(ctx context.Context, request DeleteSubscriptionRequest) (response DeleteSubscriptionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteSubscription, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteSubscriptionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteSubscriptionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteSubscriptionResponse")
	}
	return
}

// deleteSubscription implements the OCIOperation interface (enables retrying operations)
func (client NotificationDataPlaneClient) deleteSubscription(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/subscriptions/{subscriptionId}")
	if err != nil {
		return nil, err
	}

	var response DeleteSubscriptionResponse
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

// GetConfirmSubscription Gets the confirmation details for the specified subscription.
func (client NotificationDataPlaneClient) GetConfirmSubscription(ctx context.Context, request GetConfirmSubscriptionRequest) (response GetConfirmSubscriptionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getConfirmSubscription, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetConfirmSubscriptionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetConfirmSubscriptionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetConfirmSubscriptionResponse")
	}
	return
}

// getConfirmSubscription implements the OCIOperation interface (enables retrying operations)
func (client NotificationDataPlaneClient) getConfirmSubscription(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/subscriptions/{id}/confirmation")
	if err != nil {
		return nil, err
	}

	var response GetConfirmSubscriptionResponse
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

// GetSubscription Gets the specified subscription's configuration information.
func (client NotificationDataPlaneClient) GetSubscription(ctx context.Context, request GetSubscriptionRequest) (response GetSubscriptionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getSubscription, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetSubscriptionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetSubscriptionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetSubscriptionResponse")
	}
	return
}

// getSubscription implements the OCIOperation interface (enables retrying operations)
func (client NotificationDataPlaneClient) getSubscription(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/subscriptions/{subscriptionId}")
	if err != nil {
		return nil, err
	}

	var response GetSubscriptionResponse
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

// GetUnsubscription Gets the unsubscription details for the specified subscription.
func (client NotificationDataPlaneClient) GetUnsubscription(ctx context.Context, request GetUnsubscriptionRequest) (response GetUnsubscriptionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getUnsubscription, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetUnsubscriptionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetUnsubscriptionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetUnsubscriptionResponse")
	}
	return
}

// getUnsubscription implements the OCIOperation interface (enables retrying operations)
func (client NotificationDataPlaneClient) getUnsubscription(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/subscriptions/{id}/unsubscription")
	if err != nil {
		return nil, err
	}

	var response GetUnsubscriptionResponse
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

// ListSubscriptions Lists the subscriptions in the specified compartment or topic.
func (client NotificationDataPlaneClient) ListSubscriptions(ctx context.Context, request ListSubscriptionsRequest) (response ListSubscriptionsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listSubscriptions, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListSubscriptionsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListSubscriptionsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListSubscriptionsResponse")
	}
	return
}

// listSubscriptions implements the OCIOperation interface (enables retrying operations)
func (client NotificationDataPlaneClient) listSubscriptions(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/subscriptions")
	if err != nil {
		return nil, err
	}

	var response ListSubscriptionsResponse
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

// PublishMessage Publishes a message to the specified topic. For more information about publishing messages, see Publishing Messages (https://docs.cloud.oracle.com/iaas/Content/Notification/Tasks/publishingmessages.htm).
func (client NotificationDataPlaneClient) PublishMessage(ctx context.Context, request PublishMessageRequest) (response PublishMessageResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.publishMessage, policy)
	if err != nil {
		if ociResponse != nil {
			response = PublishMessageResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(PublishMessageResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into PublishMessageResponse")
	}
	return
}

// publishMessage implements the OCIOperation interface (enables retrying operations)
func (client NotificationDataPlaneClient) publishMessage(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/topics/{topicId}/messages")
	if err != nil {
		return nil, err
	}

	var response PublishMessageResponse
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

// ResendSubscriptionConfirmation Resends the confirmation details for the specified subscription.
func (client NotificationDataPlaneClient) ResendSubscriptionConfirmation(ctx context.Context, request ResendSubscriptionConfirmationRequest) (response ResendSubscriptionConfirmationResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.resendSubscriptionConfirmation, policy)
	if err != nil {
		if ociResponse != nil {
			response = ResendSubscriptionConfirmationResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ResendSubscriptionConfirmationResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ResendSubscriptionConfirmationResponse")
	}
	return
}

// resendSubscriptionConfirmation implements the OCIOperation interface (enables retrying operations)
func (client NotificationDataPlaneClient) resendSubscriptionConfirmation(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/subscriptions/{id}/resendConfirmation")
	if err != nil {
		return nil, err
	}

	var response ResendSubscriptionConfirmationResponse
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

// UpdateSubscription Updates the specified subscription's configuration.
func (client NotificationDataPlaneClient) UpdateSubscription(ctx context.Context, request UpdateSubscriptionRequest) (response UpdateSubscriptionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateSubscription, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateSubscriptionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateSubscriptionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateSubscriptionResponse")
	}
	return
}

// updateSubscription implements the OCIOperation interface (enables retrying operations)
func (client NotificationDataPlaneClient) updateSubscription(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/subscriptions/{subscriptionId}")
	if err != nil {
		return nil, err
	}

	var response UpdateSubscriptionResponse
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
