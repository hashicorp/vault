// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Announcements Service API
//
// Manage Oracle Cloud Infrastructure console announcements.
//

package announcementsservice

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//AnnouncementClient a client for Announcement
type AnnouncementClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewAnnouncementClientWithConfigurationProvider Creates a new default Announcement client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewAnnouncementClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client AnnouncementClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = AnnouncementClient{BaseClient: baseClient}
	client.BasePath = "20180904"
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *AnnouncementClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).Endpoint("announcements")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *AnnouncementClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
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
func (client *AnnouncementClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// GetAnnouncement Gets the details of a specific announcement.
func (client AnnouncementClient) GetAnnouncement(ctx context.Context, request GetAnnouncementRequest) (response GetAnnouncementResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getAnnouncement, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetAnnouncementResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetAnnouncementResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetAnnouncementResponse")
	}
	return
}

// getAnnouncement implements the OCIOperation interface (enables retrying operations)
func (client AnnouncementClient) getAnnouncement(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/announcements/{announcementId}")
	if err != nil {
		return nil, err
	}

	var response GetAnnouncementResponse
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

// GetAnnouncementUserStatus Gets information about whether a specific announcement was acknowledged by a user.
func (client AnnouncementClient) GetAnnouncementUserStatus(ctx context.Context, request GetAnnouncementUserStatusRequest) (response GetAnnouncementUserStatusResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getAnnouncementUserStatus, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetAnnouncementUserStatusResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetAnnouncementUserStatusResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetAnnouncementUserStatusResponse")
	}
	return
}

// getAnnouncementUserStatus implements the OCIOperation interface (enables retrying operations)
func (client AnnouncementClient) getAnnouncementUserStatus(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/announcements/{announcementId}/userStatus")
	if err != nil {
		return nil, err
	}

	var response GetAnnouncementUserStatusResponse
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

// ListAnnouncements Gets a list of announcements for the current tenancy.
func (client AnnouncementClient) ListAnnouncements(ctx context.Context, request ListAnnouncementsRequest) (response ListAnnouncementsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listAnnouncements, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListAnnouncementsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListAnnouncementsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListAnnouncementsResponse")
	}
	return
}

// listAnnouncements implements the OCIOperation interface (enables retrying operations)
func (client AnnouncementClient) listAnnouncements(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/announcements")
	if err != nil {
		return nil, err
	}

	var response ListAnnouncementsResponse
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

// UpdateAnnouncementUserStatus Updates the status of the specified announcement with regard to whether it has been marked as read.
func (client AnnouncementClient) UpdateAnnouncementUserStatus(ctx context.Context, request UpdateAnnouncementUserStatusRequest) (response UpdateAnnouncementUserStatusResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateAnnouncementUserStatus, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateAnnouncementUserStatusResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateAnnouncementUserStatusResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateAnnouncementUserStatusResponse")
	}
	return
}

// updateAnnouncementUserStatus implements the OCIOperation interface (enables retrying operations)
func (client AnnouncementClient) updateAnnouncementUserStatus(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/announcements/{announcementId}/userStatus")
	if err != nil {
		return nil, err
	}

	var response UpdateAnnouncementUserStatusResponse
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
