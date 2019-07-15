// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Identity and Access Management Service API
//
// APIs for managing users, groups, compartments, and policies.
//

package identity

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

//IdentityClient a client for Identity
type IdentityClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// NewIdentityClientWithConfigurationProvider Creates a new default Identity client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewIdentityClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client IdentityClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	client = IdentityClient{BaseClient: baseClient}
	client.BasePath = "20160918"
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *IdentityClient) SetRegion(region string) {
	client.Host = common.StringToRegion(region).Endpoint("identity")
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *IdentityClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
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
func (client *IdentityClient) ConfigurationProvider() *common.ConfigurationProvider {
	return client.config
}

// ActivateMfaTotpDevice Activates the specified MFA TOTP device for the user. Activation requires manual interaction with the Console.
func (client IdentityClient) ActivateMfaTotpDevice(ctx context.Context, request ActivateMfaTotpDeviceRequest) (response ActivateMfaTotpDeviceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.activateMfaTotpDevice, policy)
	if err != nil {
		if ociResponse != nil {
			response = ActivateMfaTotpDeviceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ActivateMfaTotpDeviceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ActivateMfaTotpDeviceResponse")
	}
	return
}

// activateMfaTotpDevice implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) activateMfaTotpDevice(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/users/{userId}/mfaTotpDevices/{mfaTotpDeviceId}/actions/activate")
	if err != nil {
		return nil, err
	}

	var response ActivateMfaTotpDeviceResponse
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

// AddUserToGroup Adds the specified user to the specified group and returns a `UserGroupMembership` object with its own OCID.
// After you send your request, the new object's `lifecycleState` will temporarily be CREATING. Before using the
// object, first make sure its `lifecycleState` has changed to ACTIVE.
func (client IdentityClient) AddUserToGroup(ctx context.Context, request AddUserToGroupRequest) (response AddUserToGroupResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.addUserToGroup, policy)
	if err != nil {
		if ociResponse != nil {
			response = AddUserToGroupResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(AddUserToGroupResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into AddUserToGroupResponse")
	}
	return
}

// addUserToGroup implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) addUserToGroup(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/userGroupMemberships/")
	if err != nil {
		return nil, err
	}

	var response AddUserToGroupResponse
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

// ChangeTagNamespaceCompartment Moves the specified tag namespace to the specified compartment within the same tenancy.
// To move the tag namespace, you must have the manage tag-namespaces permission on both compartments.
// For more information about IAM policies, see Details for IAM (https://docs.cloud.oracle.com/Content/Identity/Reference/iampolicyreference.htm).
// Moving a tag namespace moves all the tag key definitions contained in the tag namespace.
func (client IdentityClient) ChangeTagNamespaceCompartment(ctx context.Context, request ChangeTagNamespaceCompartmentRequest) (response ChangeTagNamespaceCompartmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.changeTagNamespaceCompartment, policy)
	if err != nil {
		if ociResponse != nil {
			response = ChangeTagNamespaceCompartmentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ChangeTagNamespaceCompartmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ChangeTagNamespaceCompartmentResponse")
	}
	return
}

// changeTagNamespaceCompartment implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) changeTagNamespaceCompartment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/tagNamespaces/{tagNamespaceId}/actions/changeCompartment")
	if err != nil {
		return nil, err
	}

	var response ChangeTagNamespaceCompartmentResponse
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

// CreateAuthToken Creates a new auth token for the specified user. For information about what auth tokens are for, see
// Managing User Credentials (https://docs.cloud.oracle.com/Content/Identity/Tasks/managingcredentials.htm).
// You must specify a *description* for the auth token (although it can be an empty string). It does not
// have to be unique, and you can change it anytime with
// UpdateAuthToken.
// Every user has permission to create an auth token for *their own user ID*. An administrator in your organization
// does not need to write a policy to give users this ability. To compare, administrators who have permission to the
// tenancy can use this operation to create an auth token for any user, including themselves.
func (client IdentityClient) CreateAuthToken(ctx context.Context, request CreateAuthTokenRequest) (response CreateAuthTokenResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createAuthToken, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateAuthTokenResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateAuthTokenResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateAuthTokenResponse")
	}
	return
}

// createAuthToken implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createAuthToken(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/users/{userId}/authTokens/")
	if err != nil {
		return nil, err
	}

	var response CreateAuthTokenResponse
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

// CreateCompartment Creates a new compartment in the specified compartment.
// **Important:** Compartments cannot be deleted.
// Specify the parent compartment's OCID as the compartment ID in the request object. Remember that the tenancy
// is simply the root compartment. For information about OCIDs, see
// Resource Identifiers (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
// You must also specify a *name* for the compartment, which must be unique across all compartments in
// your tenancy. You can use this name or the OCID when writing policies that apply
// to the compartment. For more information about policies, see
// How Policies Work (https://docs.cloud.oracle.com/Content/Identity/Concepts/policies.htm).
// You must also specify a *description* for the compartment (although it can be an empty string). It does
// not have to be unique, and you can change it anytime with
// UpdateCompartment.
// After you send your request, the new object's `lifecycleState` will temporarily be CREATING. Before using the
// object, first make sure its `lifecycleState` has changed to ACTIVE.
func (client IdentityClient) CreateCompartment(ctx context.Context, request CreateCompartmentRequest) (response CreateCompartmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createCompartment, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateCompartmentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateCompartmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateCompartmentResponse")
	}
	return
}

// createCompartment implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createCompartment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/compartments/")
	if err != nil {
		return nil, err
	}

	var response CreateCompartmentResponse
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

// CreateCustomerSecretKey Creates a new secret key for the specified user. Secret keys are used for authentication with the Object Storage Service's Amazon S3
// compatible API. For information, see
// Managing User Credentials (https://docs.cloud.oracle.com/Content/Identity/Tasks/managingcredentials.htm).
// You must specify a *description* for the secret key (although it can be an empty string). It does not
// have to be unique, and you can change it anytime with
// UpdateCustomerSecretKey.
// Every user has permission to create a secret key for *their own user ID*. An administrator in your organization
// does not need to write a policy to give users this ability. To compare, administrators who have permission to the
// tenancy can use this operation to create a secret key for any user, including themselves.
func (client IdentityClient) CreateCustomerSecretKey(ctx context.Context, request CreateCustomerSecretKeyRequest) (response CreateCustomerSecretKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createCustomerSecretKey, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateCustomerSecretKeyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateCustomerSecretKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateCustomerSecretKeyResponse")
	}
	return
}

// createCustomerSecretKey implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createCustomerSecretKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/users/{userId}/customerSecretKeys/")
	if err != nil {
		return nil, err
	}

	var response CreateCustomerSecretKeyResponse
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

// CreateDynamicGroup Creates a new dynamic group in your tenancy.
// You must specify your tenancy's OCID as the compartment ID in the request object (remember that the tenancy
// is simply the root compartment). Notice that IAM resources (users, groups, compartments, and some policies)
// reside within the tenancy itself, unlike cloud resources such as compute instances, which typically
// reside within compartments inside the tenancy. For information about OCIDs, see
// Resource Identifiers (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
// You must also specify a *name* for the dynamic group, which must be unique across all dynamic groups in your
// tenancy, and cannot be changed. Note that this name has to be also unique across all groups in your tenancy.
// You can use this name or the OCID when writing policies that apply to the dynamic group. For more information
// about policies, see How Policies Work (https://docs.cloud.oracle.com/Content/Identity/Concepts/policies.htm).
// You must also specify a *description* for the dynamic group (although it can be an empty string). It does not
// have to be unique, and you can change it anytime with UpdateDynamicGroup.
// After you send your request, the new object's `lifecycleState` will temporarily be CREATING. Before using the
// object, first make sure its `lifecycleState` has changed to ACTIVE.
func (client IdentityClient) CreateDynamicGroup(ctx context.Context, request CreateDynamicGroupRequest) (response CreateDynamicGroupResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createDynamicGroup, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateDynamicGroupResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateDynamicGroupResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateDynamicGroupResponse")
	}
	return
}

// createDynamicGroup implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createDynamicGroup(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/dynamicGroups/")
	if err != nil {
		return nil, err
	}

	var response CreateDynamicGroupResponse
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

// CreateGroup Creates a new group in your tenancy.
// You must specify your tenancy's OCID as the compartment ID in the request object (remember that the tenancy
// is simply the root compartment). Notice that IAM resources (users, groups, compartments, and some policies)
// reside within the tenancy itself, unlike cloud resources such as compute instances, which typically
// reside within compartments inside the tenancy. For information about OCIDs, see
// Resource Identifiers (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
// You must also specify a *name* for the group, which must be unique across all groups in your tenancy and
// cannot be changed. You can use this name or the OCID when writing policies that apply to the group. For more
// information about policies, see How Policies Work (https://docs.cloud.oracle.com/Content/Identity/Concepts/policies.htm).
// You must also specify a *description* for the group (although it can be an empty string). It does not
// have to be unique, and you can change it anytime with UpdateGroup.
// After you send your request, the new object's `lifecycleState` will temporarily be CREATING. Before using the
// object, first make sure its `lifecycleState` has changed to ACTIVE.
// After creating the group, you need to put users in it and write policies for it.
// See AddUserToGroup and
// CreatePolicy.
func (client IdentityClient) CreateGroup(ctx context.Context, request CreateGroupRequest) (response CreateGroupResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createGroup, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateGroupResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateGroupResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateGroupResponse")
	}
	return
}

// createGroup implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createGroup(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/groups/")
	if err != nil {
		return nil, err
	}

	var response CreateGroupResponse
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

// CreateIdentityProvider Creates a new identity provider in your tenancy. For more information, see
// Identity Providers and Federation (https://docs.cloud.oracle.com/Content/Identity/Concepts/federation.htm).
// You must specify your tenancy's OCID as the compartment ID in the request object.
// Remember that the tenancy is simply the root compartment. For information about
// OCIDs, see Resource Identifiers (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
// You must also specify a *name* for the `IdentityProvider`, which must be unique
// across all `IdentityProvider` objects in your tenancy and cannot be changed.
// You must also specify a *description* for the `IdentityProvider` (although
// it can be an empty string). It does not have to be unique, and you can change
// it anytime with
// UpdateIdentityProvider.
// After you send your request, the new object's `lifecycleState` will temporarily
// be CREATING. Before using the object, first make sure its `lifecycleState` has
// changed to ACTIVE.
func (client IdentityClient) CreateIdentityProvider(ctx context.Context, request CreateIdentityProviderRequest) (response CreateIdentityProviderResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createIdentityProvider, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateIdentityProviderResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateIdentityProviderResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateIdentityProviderResponse")
	}
	return
}

// createIdentityProvider implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createIdentityProvider(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/identityProviders/")
	if err != nil {
		return nil, err
	}

	var response CreateIdentityProviderResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponseWithPolymorphicBody(httpResponse, &response, &identityprovider{})
	return response, err
}

// CreateIdpGroupMapping Creates a single mapping between an IdP group and an IAM Service
// Group.
func (client IdentityClient) CreateIdpGroupMapping(ctx context.Context, request CreateIdpGroupMappingRequest) (response CreateIdpGroupMappingResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createIdpGroupMapping, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateIdpGroupMappingResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateIdpGroupMappingResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateIdpGroupMappingResponse")
	}
	return
}

// createIdpGroupMapping implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createIdpGroupMapping(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/identityProviders/{identityProviderId}/groupMappings/")
	if err != nil {
		return nil, err
	}

	var response CreateIdpGroupMappingResponse
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

// CreateMfaTotpDevice Creates a new MFA TOTP device for the user. A user can have one MFA TOTP device.
func (client IdentityClient) CreateMfaTotpDevice(ctx context.Context, request CreateMfaTotpDeviceRequest) (response CreateMfaTotpDeviceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createMfaTotpDevice, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateMfaTotpDeviceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateMfaTotpDeviceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateMfaTotpDeviceResponse")
	}
	return
}

// createMfaTotpDevice implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createMfaTotpDevice(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/users/{userId}/mfaTotpDevices")
	if err != nil {
		return nil, err
	}

	var response CreateMfaTotpDeviceResponse
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

// CreateOrResetUIPassword Creates a new Console one-time password for the specified user. For more information about user
// credentials, see User Credentials (https://docs.cloud.oracle.com/Content/Identity/Concepts/usercredentials.htm).
// Use this operation after creating a new user, or if a user forgets their password. The new one-time
// password is returned to you in the response, and you must securely deliver it to the user. They'll
// be prompted to change this password the next time they sign in to the Console. If they don't change
// it within 7 days, the password will expire and you'll need to create a new one-time password for the
// user.
// **Note:** The user's Console login is the unique name you specified when you created the user
// (see CreateUser).
func (client IdentityClient) CreateOrResetUIPassword(ctx context.Context, request CreateOrResetUIPasswordRequest) (response CreateOrResetUIPasswordResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createOrResetUIPassword, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateOrResetUIPasswordResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateOrResetUIPasswordResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateOrResetUIPasswordResponse")
	}
	return
}

// createOrResetUIPassword implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createOrResetUIPassword(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/users/{userId}/uiPassword")
	if err != nil {
		return nil, err
	}

	var response CreateOrResetUIPasswordResponse
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

// CreatePolicy Creates a new policy in the specified compartment (either the tenancy or another of your compartments).
// If you're new to policies, see Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// You must specify a *name* for the policy, which must be unique across all policies in your tenancy
// and cannot be changed.
// You must also specify a *description* for the policy (although it can be an empty string). It does not
// have to be unique, and you can change it anytime with UpdatePolicy.
// You must specify one or more policy statements in the statements array. For information about writing
// policies, see How Policies Work (https://docs.cloud.oracle.com/Content/Identity/Concepts/policies.htm) and
// Common Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/commonpolicies.htm).
// After you send your request, the new object's `lifecycleState` will temporarily be CREATING. Before using the
// object, first make sure its `lifecycleState` has changed to ACTIVE.
// New policies take effect typically within 10 seconds.
func (client IdentityClient) CreatePolicy(ctx context.Context, request CreatePolicyRequest) (response CreatePolicyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createPolicy, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreatePolicyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreatePolicyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreatePolicyResponse")
	}
	return
}

// createPolicy implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createPolicy(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/policies/")
	if err != nil {
		return nil, err
	}

	var response CreatePolicyResponse
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

// CreateRegionSubscription Creates a subscription to a region for a tenancy.
func (client IdentityClient) CreateRegionSubscription(ctx context.Context, request CreateRegionSubscriptionRequest) (response CreateRegionSubscriptionResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createRegionSubscription, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateRegionSubscriptionResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateRegionSubscriptionResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateRegionSubscriptionResponse")
	}
	return
}

// createRegionSubscription implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createRegionSubscription(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/tenancies/{tenancyId}/regionSubscriptions")
	if err != nil {
		return nil, err
	}

	var response CreateRegionSubscriptionResponse
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

// CreateSmtpCredential Creates a new SMTP credential for the specified user. An SMTP credential has an SMTP user name and an SMTP password.
// You must specify a *description* for the SMTP credential (although it can be an empty string). It does not
// have to be unique, and you can change it anytime with
// UpdateSmtpCredential.
func (client IdentityClient) CreateSmtpCredential(ctx context.Context, request CreateSmtpCredentialRequest) (response CreateSmtpCredentialResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createSmtpCredential, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateSmtpCredentialResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateSmtpCredentialResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateSmtpCredentialResponse")
	}
	return
}

// createSmtpCredential implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createSmtpCredential(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/users/{userId}/smtpCredentials/")
	if err != nil {
		return nil, err
	}

	var response CreateSmtpCredentialResponse
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

// CreateSwiftPassword **Deprecated. Use CreateAuthToken instead.**
// Creates a new Swift password for the specified user. For information about what Swift passwords are for, see
// Managing User Credentials (https://docs.cloud.oracle.com/Content/Identity/Tasks/managingcredentials.htm).
// You must specify a *description* for the Swift password (although it can be an empty string). It does not
// have to be unique, and you can change it anytime with
// UpdateSwiftPassword.
// Every user has permission to create a Swift password for *their own user ID*. An administrator in your organization
// does not need to write a policy to give users this ability. To compare, administrators who have permission to the
// tenancy can use this operation to create a Swift password for any user, including themselves.
func (client IdentityClient) CreateSwiftPassword(ctx context.Context, request CreateSwiftPasswordRequest) (response CreateSwiftPasswordResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createSwiftPassword, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateSwiftPasswordResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateSwiftPasswordResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateSwiftPasswordResponse")
	}
	return
}

// createSwiftPassword implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createSwiftPassword(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/users/{userId}/swiftPasswords/")
	if err != nil {
		return nil, err
	}

	var response CreateSwiftPasswordResponse
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

// CreateTag Creates a new tag in the specified tag namespace.
// You must specify either the OCID or the name of the tag namespace that will contain this tag definition.
// You must also specify a *name* for the tag, which must be unique across all tags in the tag namespace
// and cannot be changed. The name can contain any ASCII character except the space (_) or period (.) characters.
// Names are case insensitive. That means, for example, "myTag" and "mytag" are not allowed in the same namespace.
// If you specify a name that's already in use in the tag namespace, a 409 error is returned.
// You must also specify a *description* for the tag.
// It does not have to be unique, and you can change it with
// UpdateTag.
func (client IdentityClient) CreateTag(ctx context.Context, request CreateTagRequest) (response CreateTagResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createTag, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateTagResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateTagResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateTagResponse")
	}
	return
}

// createTag implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createTag(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/tagNamespaces/{tagNamespaceId}/tags")
	if err != nil {
		return nil, err
	}

	var response CreateTagResponse
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

// CreateTagDefault Creates a new tag default in the specified compartment for the specified tag definition.
func (client IdentityClient) CreateTagDefault(ctx context.Context, request CreateTagDefaultRequest) (response CreateTagDefaultResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createTagDefault, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateTagDefaultResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateTagDefaultResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateTagDefaultResponse")
	}
	return
}

// createTagDefault implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createTagDefault(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/tagDefaults")
	if err != nil {
		return nil, err
	}

	var response CreateTagDefaultResponse
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

// CreateTagNamespace Creates a new tag namespace in the specified compartment.
// You must specify the compartment ID in the request object (remember that the tenancy is simply the root
// compartment).
// You must also specify a *name* for the namespace, which must be unique across all namespaces in your tenancy
// and cannot be changed. The name can contain any ASCII character except the space (_) or period (.).
// Names are case insensitive. That means, for example, "myNamespace" and "mynamespace" are not allowed
// in the same tenancy. Once you created a namespace, you cannot change the name.
// If you specify a name that's already in use in the tenancy, a 409 error is returned.
// You must also specify a *description* for the namespace.
// It does not have to be unique, and you can change it with
// UpdateTagNamespace.
// Tag namespaces cannot be deleted, but they can be retired.
// See Retiring Key Definitions and Namespace Definitions (https://docs.cloud.oracle.com/Content/Identity/Concepts/taggingoverview.htm#Retiring) for more information.
func (client IdentityClient) CreateTagNamespace(ctx context.Context, request CreateTagNamespaceRequest) (response CreateTagNamespaceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createTagNamespace, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateTagNamespaceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateTagNamespaceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateTagNamespaceResponse")
	}
	return
}

// createTagNamespace implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createTagNamespace(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/tagNamespaces")
	if err != nil {
		return nil, err
	}

	var response CreateTagNamespaceResponse
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

// CreateUser Creates a new user in your tenancy. For conceptual information about users, your tenancy, and other
// IAM Service components, see Overview of the IAM Service (https://docs.cloud.oracle.com/Content/Identity/Concepts/overview.htm).
// You must specify your tenancy's OCID as the compartment ID in the request object (remember that the
// tenancy is simply the root compartment). Notice that IAM resources (users, groups, compartments, and
// some policies) reside within the tenancy itself, unlike cloud resources such as compute instances,
// which typically reside within compartments inside the tenancy. For information about OCIDs, see
// Resource Identifiers (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
// You must also specify a *name* for the user, which must be unique across all users in your tenancy
// and cannot be changed. Allowed characters: No spaces. Only letters, numerals, hyphens, periods,
// underscores, +, and @. If you specify a name that's already in use, you'll get a 409 error.
// This name will be the user's login to the Console. You might want to pick a
// name that your company's own identity system (e.g., Active Directory, LDAP, etc.) already uses.
// If you delete a user and then create a new user with the same name, they'll be considered different
// users because they have different OCIDs.
// You must also specify a *description* for the user (although it can be an empty string).
// It does not have to be unique, and you can change it anytime with
// UpdateUser. You can use the field to provide the user's
// full name, a description, a nickname, or other information to generally identify the user.
// After you send your request, the new object's `lifecycleState` will temporarily be CREATING. Before
// using the object, first make sure its `lifecycleState` has changed to ACTIVE.
// A new user has no permissions until you place the user in one or more groups (see
// AddUserToGroup). If the user needs to
// access the Console, you need to provide the user a password (see
// CreateOrResetUIPassword).
// If the user needs to access the Oracle Cloud Infrastructure REST API, you need to upload a
// public API signing key for that user (see
// Required Keys and OCIDs (https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm) and also
// UploadApiKey).
// **Important:** Make sure to inform the new user which compartment(s) they have access to.
func (client IdentityClient) CreateUser(ctx context.Context, request CreateUserRequest) (response CreateUserResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.createUser, policy)
	if err != nil {
		if ociResponse != nil {
			response = CreateUserResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(CreateUserResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into CreateUserResponse")
	}
	return
}

// createUser implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) createUser(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/users/")
	if err != nil {
		return nil, err
	}

	var response CreateUserResponse
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

// DeleteApiKey Deletes the specified API signing key for the specified user.
// Every user has permission to use this operation to delete a key for *their own user ID*. An
// administrator in your organization does not need to write a policy to give users this ability.
// To compare, administrators who have permission to the tenancy can use this operation to delete
// a key for any user, including themselves.
func (client IdentityClient) DeleteApiKey(ctx context.Context, request DeleteApiKeyRequest) (response DeleteApiKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteApiKey, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteApiKeyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteApiKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteApiKeyResponse")
	}
	return
}

// deleteApiKey implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteApiKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/users/{userId}/apiKeys/{fingerprint}")
	if err != nil {
		return nil, err
	}

	var response DeleteApiKeyResponse
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

// DeleteAuthToken Deletes the specified auth token for the specified user.
func (client IdentityClient) DeleteAuthToken(ctx context.Context, request DeleteAuthTokenRequest) (response DeleteAuthTokenResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteAuthToken, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteAuthTokenResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteAuthTokenResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteAuthTokenResponse")
	}
	return
}

// deleteAuthToken implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteAuthToken(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/users/{userId}/authTokens/{authTokenId}")
	if err != nil {
		return nil, err
	}

	var response DeleteAuthTokenResponse
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

// DeleteCompartment Deletes the specified compartment. The compartment must be empty.
func (client IdentityClient) DeleteCompartment(ctx context.Context, request DeleteCompartmentRequest) (response DeleteCompartmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteCompartment, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteCompartmentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteCompartmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteCompartmentResponse")
	}
	return
}

// deleteCompartment implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteCompartment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/compartments/{compartmentId}")
	if err != nil {
		return nil, err
	}

	var response DeleteCompartmentResponse
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

// DeleteCustomerSecretKey Deletes the specified secret key for the specified user.
func (client IdentityClient) DeleteCustomerSecretKey(ctx context.Context, request DeleteCustomerSecretKeyRequest) (response DeleteCustomerSecretKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteCustomerSecretKey, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteCustomerSecretKeyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteCustomerSecretKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteCustomerSecretKeyResponse")
	}
	return
}

// deleteCustomerSecretKey implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteCustomerSecretKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/users/{userId}/customerSecretKeys/{customerSecretKeyId}")
	if err != nil {
		return nil, err
	}

	var response DeleteCustomerSecretKeyResponse
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

// DeleteDynamicGroup Deletes the specified dynamic group.
func (client IdentityClient) DeleteDynamicGroup(ctx context.Context, request DeleteDynamicGroupRequest) (response DeleteDynamicGroupResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteDynamicGroup, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteDynamicGroupResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteDynamicGroupResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteDynamicGroupResponse")
	}
	return
}

// deleteDynamicGroup implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteDynamicGroup(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/dynamicGroups/{dynamicGroupId}")
	if err != nil {
		return nil, err
	}

	var response DeleteDynamicGroupResponse
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

// DeleteGroup Deletes the specified group. The group must be empty.
func (client IdentityClient) DeleteGroup(ctx context.Context, request DeleteGroupRequest) (response DeleteGroupResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteGroup, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteGroupResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteGroupResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteGroupResponse")
	}
	return
}

// deleteGroup implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteGroup(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/groups/{groupId}")
	if err != nil {
		return nil, err
	}

	var response DeleteGroupResponse
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

// DeleteIdentityProvider Deletes the specified identity provider. The identity provider must not have
// any group mappings (see IdpGroupMapping).
func (client IdentityClient) DeleteIdentityProvider(ctx context.Context, request DeleteIdentityProviderRequest) (response DeleteIdentityProviderResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteIdentityProvider, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteIdentityProviderResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteIdentityProviderResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteIdentityProviderResponse")
	}
	return
}

// deleteIdentityProvider implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteIdentityProvider(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/identityProviders/{identityProviderId}")
	if err != nil {
		return nil, err
	}

	var response DeleteIdentityProviderResponse
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

// DeleteIdpGroupMapping Deletes the specified group mapping.
func (client IdentityClient) DeleteIdpGroupMapping(ctx context.Context, request DeleteIdpGroupMappingRequest) (response DeleteIdpGroupMappingResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteIdpGroupMapping, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteIdpGroupMappingResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteIdpGroupMappingResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteIdpGroupMappingResponse")
	}
	return
}

// deleteIdpGroupMapping implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteIdpGroupMapping(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/identityProviders/{identityProviderId}/groupMappings/{mappingId}")
	if err != nil {
		return nil, err
	}

	var response DeleteIdpGroupMappingResponse
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

// DeleteMfaTotpDevice Deletes the specified MFA TOTP device for the specified user.
func (client IdentityClient) DeleteMfaTotpDevice(ctx context.Context, request DeleteMfaTotpDeviceRequest) (response DeleteMfaTotpDeviceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteMfaTotpDevice, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteMfaTotpDeviceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteMfaTotpDeviceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteMfaTotpDeviceResponse")
	}
	return
}

// deleteMfaTotpDevice implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteMfaTotpDevice(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/users/{userId}/mfaTotpDevices/{mfaTotpDeviceId}")
	if err != nil {
		return nil, err
	}

	var response DeleteMfaTotpDeviceResponse
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

// DeletePolicy Deletes the specified policy. The deletion takes effect typically within 10 seconds.
func (client IdentityClient) DeletePolicy(ctx context.Context, request DeletePolicyRequest) (response DeletePolicyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deletePolicy, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeletePolicyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeletePolicyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeletePolicyResponse")
	}
	return
}

// deletePolicy implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deletePolicy(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/policies/{policyId}")
	if err != nil {
		return nil, err
	}

	var response DeletePolicyResponse
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

// DeleteSmtpCredential Deletes the specified SMTP credential for the specified user.
func (client IdentityClient) DeleteSmtpCredential(ctx context.Context, request DeleteSmtpCredentialRequest) (response DeleteSmtpCredentialResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteSmtpCredential, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteSmtpCredentialResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteSmtpCredentialResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteSmtpCredentialResponse")
	}
	return
}

// deleteSmtpCredential implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteSmtpCredential(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/users/{userId}/smtpCredentials/{smtpCredentialId}")
	if err != nil {
		return nil, err
	}

	var response DeleteSmtpCredentialResponse
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

// DeleteSwiftPassword **Deprecated. Use DeleteAuthToken instead.**
// Deletes the specified Swift password for the specified user.
func (client IdentityClient) DeleteSwiftPassword(ctx context.Context, request DeleteSwiftPasswordRequest) (response DeleteSwiftPasswordResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteSwiftPassword, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteSwiftPasswordResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteSwiftPasswordResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteSwiftPasswordResponse")
	}
	return
}

// deleteSwiftPassword implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteSwiftPassword(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/users/{userId}/swiftPasswords/{swiftPasswordId}")
	if err != nil {
		return nil, err
	}

	var response DeleteSwiftPasswordResponse
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

// DeleteTag Deletes the the specified tag definition.
func (client IdentityClient) DeleteTag(ctx context.Context, request DeleteTagRequest) (response DeleteTagResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteTag, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteTagResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteTagResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteTagResponse")
	}
	return
}

// deleteTag implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteTag(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/tagNamespaces/{tagNamespaceId}/tags/{tagName}")
	if err != nil {
		return nil, err
	}

	var response DeleteTagResponse
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

// DeleteTagDefault Deletes the the specified tag default.
func (client IdentityClient) DeleteTagDefault(ctx context.Context, request DeleteTagDefaultRequest) (response DeleteTagDefaultResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteTagDefault, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteTagDefaultResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteTagDefaultResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteTagDefaultResponse")
	}
	return
}

// deleteTagDefault implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteTagDefault(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/tagDefaults/{tagDefaultId}")
	if err != nil {
		return nil, err
	}

	var response DeleteTagDefaultResponse
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

// DeleteTagNamespace Delete the specified tag namespace. Only an empty tagnamespace can be deleted.
// If the tag namespace you are trying to delete is not empty, please remove tag definitions from it first.
func (client IdentityClient) DeleteTagNamespace(ctx context.Context, request DeleteTagNamespaceRequest) (response DeleteTagNamespaceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteTagNamespace, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteTagNamespaceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteTagNamespaceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteTagNamespaceResponse")
	}
	return
}

// deleteTagNamespace implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteTagNamespace(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/tagNamespaces/{tagNamespaceId}")
	if err != nil {
		return nil, err
	}

	var response DeleteTagNamespaceResponse
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

// DeleteUser Deletes the specified user. The user must not be in any groups.
func (client IdentityClient) DeleteUser(ctx context.Context, request DeleteUserRequest) (response DeleteUserResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.deleteUser, policy)
	if err != nil {
		if ociResponse != nil {
			response = DeleteUserResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(DeleteUserResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into DeleteUserResponse")
	}
	return
}

// deleteUser implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) deleteUser(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/users/{userId}")
	if err != nil {
		return nil, err
	}

	var response DeleteUserResponse
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

// GenerateTotpSeed Generate seed for the MFA TOTP device.
func (client IdentityClient) GenerateTotpSeed(ctx context.Context, request GenerateTotpSeedRequest) (response GenerateTotpSeedResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.generateTotpSeed, policy)
	if err != nil {
		if ociResponse != nil {
			response = GenerateTotpSeedResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GenerateTotpSeedResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GenerateTotpSeedResponse")
	}
	return
}

// generateTotpSeed implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) generateTotpSeed(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/users/{userId}/mfaTotpDevices/{mfaTotpDeviceId}/actions/generateSeed")
	if err != nil {
		return nil, err
	}

	var response GenerateTotpSeedResponse
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

// GetAuthenticationPolicy Gets the authentication policy for the given tenancy. You must specify your tenant’s OCID as the value for
// the compartment ID (remember that the tenancy is simply the root compartment).
func (client IdentityClient) GetAuthenticationPolicy(ctx context.Context, request GetAuthenticationPolicyRequest) (response GetAuthenticationPolicyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getAuthenticationPolicy, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetAuthenticationPolicyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetAuthenticationPolicyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetAuthenticationPolicyResponse")
	}
	return
}

// getAuthenticationPolicy implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getAuthenticationPolicy(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/authenticationPolicies/{compartmentId}")
	if err != nil {
		return nil, err
	}

	var response GetAuthenticationPolicyResponse
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

// GetCompartment Gets the specified compartment's information.
// This operation does not return a list of all the resources inside the compartment. There is no single
// API operation that does that. Compartments can contain multiple types of resources (instances, block
// storage volumes, etc.). To find out what's in a compartment, you must call the "List" operation for
// each resource type and specify the compartment's OCID as a query parameter in the request. For example,
// call the ListInstances operation in the Cloud Compute
// Service or the ListVolumes operation in Cloud Block Storage.
func (client IdentityClient) GetCompartment(ctx context.Context, request GetCompartmentRequest) (response GetCompartmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getCompartment, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetCompartmentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetCompartmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetCompartmentResponse")
	}
	return
}

// getCompartment implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getCompartment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/compartments/{compartmentId}")
	if err != nil {
		return nil, err
	}

	var response GetCompartmentResponse
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

// GetDynamicGroup Gets the specified dynamic group's information.
func (client IdentityClient) GetDynamicGroup(ctx context.Context, request GetDynamicGroupRequest) (response GetDynamicGroupResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getDynamicGroup, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetDynamicGroupResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetDynamicGroupResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetDynamicGroupResponse")
	}
	return
}

// getDynamicGroup implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getDynamicGroup(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/dynamicGroups/{dynamicGroupId}")
	if err != nil {
		return nil, err
	}

	var response GetDynamicGroupResponse
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

// GetGroup Gets the specified group's information.
// This operation does not return a list of all the users in the group. To do that, use
// ListUserGroupMemberships and
// provide the group's OCID as a query parameter in the request.
func (client IdentityClient) GetGroup(ctx context.Context, request GetGroupRequest) (response GetGroupResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getGroup, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetGroupResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetGroupResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetGroupResponse")
	}
	return
}

// getGroup implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getGroup(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/groups/{groupId}")
	if err != nil {
		return nil, err
	}

	var response GetGroupResponse
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

// GetIdentityProvider Gets the specified identity provider's information.
func (client IdentityClient) GetIdentityProvider(ctx context.Context, request GetIdentityProviderRequest) (response GetIdentityProviderResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getIdentityProvider, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetIdentityProviderResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetIdentityProviderResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetIdentityProviderResponse")
	}
	return
}

// getIdentityProvider implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getIdentityProvider(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/identityProviders/{identityProviderId}")
	if err != nil {
		return nil, err
	}

	var response GetIdentityProviderResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponseWithPolymorphicBody(httpResponse, &response, &identityprovider{})
	return response, err
}

// GetIdpGroupMapping Gets the specified group mapping.
func (client IdentityClient) GetIdpGroupMapping(ctx context.Context, request GetIdpGroupMappingRequest) (response GetIdpGroupMappingResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getIdpGroupMapping, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetIdpGroupMappingResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetIdpGroupMappingResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetIdpGroupMappingResponse")
	}
	return
}

// getIdpGroupMapping implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getIdpGroupMapping(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/identityProviders/{identityProviderId}/groupMappings/{mappingId}")
	if err != nil {
		return nil, err
	}

	var response GetIdpGroupMappingResponse
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

// GetMfaTotpDevice Get the specified MFA TOTP device for the specified user.
func (client IdentityClient) GetMfaTotpDevice(ctx context.Context, request GetMfaTotpDeviceRequest) (response GetMfaTotpDeviceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getMfaTotpDevice, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetMfaTotpDeviceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetMfaTotpDeviceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetMfaTotpDeviceResponse")
	}
	return
}

// getMfaTotpDevice implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getMfaTotpDevice(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/users/{userId}/mfaTotpDevices/{mfaTotpDeviceId}")
	if err != nil {
		return nil, err
	}

	var response GetMfaTotpDeviceResponse
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

// GetPolicy Gets the specified policy's information.
func (client IdentityClient) GetPolicy(ctx context.Context, request GetPolicyRequest) (response GetPolicyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getPolicy, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetPolicyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetPolicyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetPolicyResponse")
	}
	return
}

// getPolicy implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getPolicy(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/policies/{policyId}")
	if err != nil {
		return nil, err
	}

	var response GetPolicyResponse
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

// GetTag Gets the specified tag's information.
func (client IdentityClient) GetTag(ctx context.Context, request GetTagRequest) (response GetTagResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getTag, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetTagResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetTagResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetTagResponse")
	}
	return
}

// getTag implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getTag(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/tagNamespaces/{tagNamespaceId}/tags/{tagName}")
	if err != nil {
		return nil, err
	}

	var response GetTagResponse
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

// GetTagDefault Retrieves the specified tag default.
func (client IdentityClient) GetTagDefault(ctx context.Context, request GetTagDefaultRequest) (response GetTagDefaultResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getTagDefault, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetTagDefaultResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetTagDefaultResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetTagDefaultResponse")
	}
	return
}

// getTagDefault implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getTagDefault(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/tagDefaults/{tagDefaultId}")
	if err != nil {
		return nil, err
	}

	var response GetTagDefaultResponse
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

// GetTagNamespace Gets the specified tag namespace's information.
func (client IdentityClient) GetTagNamespace(ctx context.Context, request GetTagNamespaceRequest) (response GetTagNamespaceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getTagNamespace, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetTagNamespaceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetTagNamespaceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetTagNamespaceResponse")
	}
	return
}

// getTagNamespace implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getTagNamespace(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/tagNamespaces/{tagNamespaceId}")
	if err != nil {
		return nil, err
	}

	var response GetTagNamespaceResponse
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

// GetTenancy Get the specified tenancy's information.
func (client IdentityClient) GetTenancy(ctx context.Context, request GetTenancyRequest) (response GetTenancyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getTenancy, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetTenancyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetTenancyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetTenancyResponse")
	}
	return
}

// getTenancy implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getTenancy(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/tenancies/{tenancyId}")
	if err != nil {
		return nil, err
	}

	var response GetTenancyResponse
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

// GetUser Gets the specified user's information.
func (client IdentityClient) GetUser(ctx context.Context, request GetUserRequest) (response GetUserResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getUser, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetUserResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetUserResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetUserResponse")
	}
	return
}

// getUser implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getUser(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/users/{userId}")
	if err != nil {
		return nil, err
	}

	var response GetUserResponse
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

// GetUserGroupMembership Gets the specified UserGroupMembership's information.
func (client IdentityClient) GetUserGroupMembership(ctx context.Context, request GetUserGroupMembershipRequest) (response GetUserGroupMembershipResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getUserGroupMembership, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetUserGroupMembershipResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetUserGroupMembershipResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetUserGroupMembershipResponse")
	}
	return
}

// getUserGroupMembership implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getUserGroupMembership(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/userGroupMemberships/{userGroupMembershipId}")
	if err != nil {
		return nil, err
	}

	var response GetUserGroupMembershipResponse
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

// GetUserUIPasswordInformation Gets the specified user's console password information. The returned object contains the user's OCID,
// but not the password itself. The actual password is returned only when created or reset.
func (client IdentityClient) GetUserUIPasswordInformation(ctx context.Context, request GetUserUIPasswordInformationRequest) (response GetUserUIPasswordInformationResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getUserUIPasswordInformation, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetUserUIPasswordInformationResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetUserUIPasswordInformationResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetUserUIPasswordInformationResponse")
	}
	return
}

// getUserUIPasswordInformation implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getUserUIPasswordInformation(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/users/{userId}/uiPassword")
	if err != nil {
		return nil, err
	}

	var response GetUserUIPasswordInformationResponse
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

// GetWorkRequest Gets details on a specified work request. The workRequestID is returned in the opc-workrequest-id header
// for any asynchronous operation in the Identity and Access Management service.
func (client IdentityClient) GetWorkRequest(ctx context.Context, request GetWorkRequestRequest) (response GetWorkRequestResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.getWorkRequest, policy)
	if err != nil {
		if ociResponse != nil {
			response = GetWorkRequestResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(GetWorkRequestResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into GetWorkRequestResponse")
	}
	return
}

// getWorkRequest implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) getWorkRequest(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/workRequests/{workRequestId}")
	if err != nil {
		return nil, err
	}

	var response GetWorkRequestResponse
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

// ListApiKeys Lists the API signing keys for the specified user. A user can have a maximum of three keys.
// Every user has permission to use this API call for *their own user ID*.  An administrator in your
// organization does not need to write a policy to give users this ability.
func (client IdentityClient) ListApiKeys(ctx context.Context, request ListApiKeysRequest) (response ListApiKeysResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listApiKeys, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListApiKeysResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListApiKeysResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListApiKeysResponse")
	}
	return
}

// listApiKeys implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listApiKeys(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/users/{userId}/apiKeys/")
	if err != nil {
		return nil, err
	}

	var response ListApiKeysResponse
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

// ListAuthTokens Lists the auth tokens for the specified user. The returned object contains the token's OCID, but not
// the token itself. The actual token is returned only upon creation.
func (client IdentityClient) ListAuthTokens(ctx context.Context, request ListAuthTokensRequest) (response ListAuthTokensResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listAuthTokens, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListAuthTokensResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListAuthTokensResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListAuthTokensResponse")
	}
	return
}

// listAuthTokens implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listAuthTokens(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/users/{userId}/authTokens/")
	if err != nil {
		return nil, err
	}

	var response ListAuthTokensResponse
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

// ListAvailabilityDomains Lists the availability domains in your tenancy. Specify the OCID of either the tenancy or another
// of your compartments as the value for the compartment ID (remember that the tenancy is simply the root compartment).
// See Where to Get the Tenancy's OCID and User's OCID (https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#five).
// Note that the order of the results returned can change if availability domains are added or removed; therefore, do not
// create a dependency on the list order.
func (client IdentityClient) ListAvailabilityDomains(ctx context.Context, request ListAvailabilityDomainsRequest) (response ListAvailabilityDomainsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listAvailabilityDomains, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListAvailabilityDomainsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListAvailabilityDomainsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListAvailabilityDomainsResponse")
	}
	return
}

// listAvailabilityDomains implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listAvailabilityDomains(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/availabilityDomains/")
	if err != nil {
		return nil, err
	}

	var response ListAvailabilityDomainsResponse
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

// ListCompartments Lists the compartments in a specified compartment. The members of the list
// returned depends on the values set for several parameters.
// With the exception of the tenancy (root compartment), the ListCompartments operation
// returns only the first-level child compartments in the parent compartment specified in
// `compartmentId`. The list does not include any subcompartments of the child
// compartments (grandchildren).
// The parameter `accessLevel` specifies whether to return only those compartments for which the
// requestor has INSPECT permissions on at least one resource directly
// or indirectly (the resource can be in a subcompartment).
// The parameter `compartmentIdInSubtree` applies only when you perform ListCompartments on the
// tenancy (root compartment). When set to true, the entire hierarchy of compartments can be returned.
// To get a full list of all compartments and subcompartments in the tenancy (root compartment),
// set the parameter `compartmentIdInSubtree` to true and `accessLevel` to ANY.
// See Where to Get the Tenancy's OCID and User's OCID (https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#five).
func (client IdentityClient) ListCompartments(ctx context.Context, request ListCompartmentsRequest) (response ListCompartmentsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listCompartments, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListCompartmentsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListCompartmentsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListCompartmentsResponse")
	}
	return
}

// listCompartments implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listCompartments(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/compartments/")
	if err != nil {
		return nil, err
	}

	var response ListCompartmentsResponse
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

// ListCostTrackingTags Lists all the tags enabled for cost-tracking in the specified tenancy. For information about
// cost-tracking tags, see Using Cost-tracking Tags (https://docs.cloud.oracle.com/Content/Identity/Concepts/taggingoverview.htm#costs).
func (client IdentityClient) ListCostTrackingTags(ctx context.Context, request ListCostTrackingTagsRequest) (response ListCostTrackingTagsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listCostTrackingTags, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListCostTrackingTagsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListCostTrackingTagsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListCostTrackingTagsResponse")
	}
	return
}

// listCostTrackingTags implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listCostTrackingTags(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/tagNamespaces/actions/listCostTrackingTags")
	if err != nil {
		return nil, err
	}

	var response ListCostTrackingTagsResponse
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

// ListCustomerSecretKeys Lists the secret keys for the specified user. The returned object contains the secret key's OCID, but not
// the secret key itself. The actual secret key is returned only upon creation.
func (client IdentityClient) ListCustomerSecretKeys(ctx context.Context, request ListCustomerSecretKeysRequest) (response ListCustomerSecretKeysResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listCustomerSecretKeys, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListCustomerSecretKeysResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListCustomerSecretKeysResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListCustomerSecretKeysResponse")
	}
	return
}

// listCustomerSecretKeys implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listCustomerSecretKeys(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/users/{userId}/customerSecretKeys/")
	if err != nil {
		return nil, err
	}

	var response ListCustomerSecretKeysResponse
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

// ListDynamicGroups Lists the dynamic groups in your tenancy. You must specify your tenancy's OCID as the value for
// the compartment ID (remember that the tenancy is simply the root compartment).
// See Where to Get the Tenancy's OCID and User's OCID (https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#five).
func (client IdentityClient) ListDynamicGroups(ctx context.Context, request ListDynamicGroupsRequest) (response ListDynamicGroupsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listDynamicGroups, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListDynamicGroupsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListDynamicGroupsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListDynamicGroupsResponse")
	}
	return
}

// listDynamicGroups implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listDynamicGroups(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/dynamicGroups/")
	if err != nil {
		return nil, err
	}

	var response ListDynamicGroupsResponse
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

// ListFaultDomains Lists the Fault Domains in your tenancy. Specify the OCID of either the tenancy or another
// of your compartments as the value for the compartment ID (remember that the tenancy is simply the root compartment).
// See Where to Get the Tenancy's OCID and User's OCID (https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#five).
func (client IdentityClient) ListFaultDomains(ctx context.Context, request ListFaultDomainsRequest) (response ListFaultDomainsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listFaultDomains, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListFaultDomainsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListFaultDomainsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListFaultDomainsResponse")
	}
	return
}

// listFaultDomains implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listFaultDomains(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/faultDomains/")
	if err != nil {
		return nil, err
	}

	var response ListFaultDomainsResponse
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

// ListGroups Lists the groups in your tenancy. You must specify your tenancy's OCID as the value for
// the compartment ID (remember that the tenancy is simply the root compartment).
// See Where to Get the Tenancy's OCID and User's OCID (https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#five).
func (client IdentityClient) ListGroups(ctx context.Context, request ListGroupsRequest) (response ListGroupsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listGroups, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListGroupsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListGroupsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListGroupsResponse")
	}
	return
}

// listGroups implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listGroups(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/groups/")
	if err != nil {
		return nil, err
	}

	var response ListGroupsResponse
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

// ListIdentityProviderGroups Lists the identity provider groups.
func (client IdentityClient) ListIdentityProviderGroups(ctx context.Context, request ListIdentityProviderGroupsRequest) (response ListIdentityProviderGroupsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listIdentityProviderGroups, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListIdentityProviderGroupsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListIdentityProviderGroupsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListIdentityProviderGroupsResponse")
	}
	return
}

// listIdentityProviderGroups implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listIdentityProviderGroups(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/identityProviders/{identityProviderId}/groups/")
	if err != nil {
		return nil, err
	}

	var response ListIdentityProviderGroupsResponse
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

//listidentityprovider allows to unmarshal list of polymorphic IdentityProvider
type listidentityprovider []identityprovider

//UnmarshalPolymorphicJSON unmarshals polymorphic json list of items
func (m *listidentityprovider) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {
	res := make([]IdentityProvider, len(*m))
	for i, v := range *m {
		nn, err := v.UnmarshalPolymorphicJSON(v.JsonData)
		if err != nil {
			return nil, err
		}
		res[i] = nn.(IdentityProvider)
	}
	return res, nil
}

// ListIdentityProviders Lists all the identity providers in your tenancy. You must specify the identity provider type (e.g., `SAML2` for
// identity providers using the SAML2.0 protocol). You must specify your tenancy's OCID as the value for the
// compartment ID (remember that the tenancy is simply the root compartment).
// See Where to Get the Tenancy's OCID and User's OCID (https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#five).
func (client IdentityClient) ListIdentityProviders(ctx context.Context, request ListIdentityProvidersRequest) (response ListIdentityProvidersResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listIdentityProviders, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListIdentityProvidersResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListIdentityProvidersResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListIdentityProvidersResponse")
	}
	return
}

// listIdentityProviders implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listIdentityProviders(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/identityProviders/")
	if err != nil {
		return nil, err
	}

	var response ListIdentityProvidersResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponseWithPolymorphicBody(httpResponse, &response, &listidentityprovider{})
	return response, err
}

// ListIdpGroupMappings Lists the group mappings for the specified identity provider.
func (client IdentityClient) ListIdpGroupMappings(ctx context.Context, request ListIdpGroupMappingsRequest) (response ListIdpGroupMappingsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listIdpGroupMappings, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListIdpGroupMappingsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListIdpGroupMappingsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListIdpGroupMappingsResponse")
	}
	return
}

// listIdpGroupMappings implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listIdpGroupMappings(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/identityProviders/{identityProviderId}/groupMappings/")
	if err != nil {
		return nil, err
	}

	var response ListIdpGroupMappingsResponse
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

// ListMfaTotpDevices Lists the MFA TOTP devices for the specified user. The returned object contains the device's OCID, but not
// the seed. The seed is returned only upon creation or when the IAM service regenerates the MFA seed for the device.
func (client IdentityClient) ListMfaTotpDevices(ctx context.Context, request ListMfaTotpDevicesRequest) (response ListMfaTotpDevicesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listMfaTotpDevices, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListMfaTotpDevicesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListMfaTotpDevicesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListMfaTotpDevicesResponse")
	}
	return
}

// listMfaTotpDevices implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listMfaTotpDevices(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/users/{userId}/mfaTotpDevices")
	if err != nil {
		return nil, err
	}

	var response ListMfaTotpDevicesResponse
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

// ListPolicies Lists the policies in the specified compartment (either the tenancy or another of your compartments).
// See Where to Get the Tenancy's OCID and User's OCID (https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#five).
// To determine which policies apply to a particular group or compartment, you must view the individual
// statements inside all your policies. There isn't a way to automatically obtain that information via the API.
func (client IdentityClient) ListPolicies(ctx context.Context, request ListPoliciesRequest) (response ListPoliciesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listPolicies, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListPoliciesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListPoliciesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListPoliciesResponse")
	}
	return
}

// listPolicies implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listPolicies(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/policies/")
	if err != nil {
		return nil, err
	}

	var response ListPoliciesResponse
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

// ListRegionSubscriptions Lists the region subscriptions for the specified tenancy.
func (client IdentityClient) ListRegionSubscriptions(ctx context.Context, request ListRegionSubscriptionsRequest) (response ListRegionSubscriptionsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listRegionSubscriptions, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListRegionSubscriptionsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListRegionSubscriptionsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListRegionSubscriptionsResponse")
	}
	return
}

// listRegionSubscriptions implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listRegionSubscriptions(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/tenancies/{tenancyId}/regionSubscriptions")
	if err != nil {
		return nil, err
	}

	var response ListRegionSubscriptionsResponse
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

// ListRegions Lists all the regions offered by Oracle Cloud Infrastructure.
func (client IdentityClient) ListRegions(ctx context.Context) (response ListRegionsResponse, err error) {
	var ociResponse common.OCIResponse
	ociResponse, err = client.listRegions(ctx)
	if err != nil {
		if ociResponse != nil {
			response = ListRegionsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListRegionsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListRegionsResponse")
	}
	return
}

// listRegions performs the request (retry policy is not enabled without a request object)
func (client IdentityClient) listRegions(ctx context.Context) (common.OCIResponse, error) {
	httpRequest := common.MakeDefaultHTTPRequest(http.MethodGet, "/regions")
	var err error

	var response ListRegionsResponse
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

// ListSmtpCredentials Lists the SMTP credentials for the specified user. The returned object contains the credential's OCID,
// the SMTP user name but not the SMTP password. The SMTP password is returned only upon creation.
func (client IdentityClient) ListSmtpCredentials(ctx context.Context, request ListSmtpCredentialsRequest) (response ListSmtpCredentialsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listSmtpCredentials, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListSmtpCredentialsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListSmtpCredentialsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListSmtpCredentialsResponse")
	}
	return
}

// listSmtpCredentials implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listSmtpCredentials(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/users/{userId}/smtpCredentials/")
	if err != nil {
		return nil, err
	}

	var response ListSmtpCredentialsResponse
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

// ListSwiftPasswords **Deprecated. Use ListAuthTokens instead.**
// Lists the Swift passwords for the specified user. The returned object contains the password's OCID, but not
// the password itself. The actual password is returned only upon creation.
func (client IdentityClient) ListSwiftPasswords(ctx context.Context, request ListSwiftPasswordsRequest) (response ListSwiftPasswordsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listSwiftPasswords, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListSwiftPasswordsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListSwiftPasswordsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListSwiftPasswordsResponse")
	}
	return
}

// listSwiftPasswords implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listSwiftPasswords(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/users/{userId}/swiftPasswords/")
	if err != nil {
		return nil, err
	}

	var response ListSwiftPasswordsResponse
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

// ListTagDefaults Lists the tag defaults for tag definitions in the specified compartment.
func (client IdentityClient) ListTagDefaults(ctx context.Context, request ListTagDefaultsRequest) (response ListTagDefaultsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listTagDefaults, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListTagDefaultsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListTagDefaultsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListTagDefaultsResponse")
	}
	return
}

// listTagDefaults implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listTagDefaults(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/tagDefaults")
	if err != nil {
		return nil, err
	}

	var response ListTagDefaultsResponse
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

// ListTagNamespaces Lists the tag namespaces in the specified compartment.
func (client IdentityClient) ListTagNamespaces(ctx context.Context, request ListTagNamespacesRequest) (response ListTagNamespacesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listTagNamespaces, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListTagNamespacesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListTagNamespacesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListTagNamespacesResponse")
	}
	return
}

// listTagNamespaces implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listTagNamespaces(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/tagNamespaces")
	if err != nil {
		return nil, err
	}

	var response ListTagNamespacesResponse
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

// ListTags Lists the tag definitions in the specified tag namespace.
func (client IdentityClient) ListTags(ctx context.Context, request ListTagsRequest) (response ListTagsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listTags, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListTagsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListTagsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListTagsResponse")
	}
	return
}

// listTags implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listTags(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/tagNamespaces/{tagNamespaceId}/tags")
	if err != nil {
		return nil, err
	}

	var response ListTagsResponse
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

// ListUserGroupMemberships Lists the `UserGroupMembership` objects in your tenancy. You must specify your tenancy's OCID
// as the value for the compartment ID
// (see Where to Get the Tenancy's OCID and User's OCID (https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#five)).
// You must also then filter the list in one of these ways:
// - You can limit the results to just the memberships for a given user by specifying a `userId`.
// - Similarly, you can limit the results to just the memberships for a given group by specifying a `groupId`.
// - You can set both the `userId` and `groupId` to determine if the specified user is in the specified group.
// If the answer is no, the response is an empty list.
// - Although`userId` and `groupId` are not indvidually required, you must set one of them.
func (client IdentityClient) ListUserGroupMemberships(ctx context.Context, request ListUserGroupMembershipsRequest) (response ListUserGroupMembershipsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listUserGroupMemberships, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListUserGroupMembershipsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListUserGroupMembershipsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListUserGroupMembershipsResponse")
	}
	return
}

// listUserGroupMemberships implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listUserGroupMemberships(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/userGroupMemberships/")
	if err != nil {
		return nil, err
	}

	var response ListUserGroupMembershipsResponse
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

// ListUsers Lists the users in your tenancy. You must specify your tenancy's OCID as the value for the
// compartment ID (remember that the tenancy is simply the root compartment).
// See Where to Get the Tenancy's OCID and User's OCID (https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#five).
func (client IdentityClient) ListUsers(ctx context.Context, request ListUsersRequest) (response ListUsersResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listUsers, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListUsersResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListUsersResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListUsersResponse")
	}
	return
}

// listUsers implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listUsers(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/users/")
	if err != nil {
		return nil, err
	}

	var response ListUsersResponse
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

// ListWorkRequests Lists the work requests in compartment.
func (client IdentityClient) ListWorkRequests(ctx context.Context, request ListWorkRequestsRequest) (response ListWorkRequestsResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.listWorkRequests, policy)
	if err != nil {
		if ociResponse != nil {
			response = ListWorkRequestsResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ListWorkRequestsResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ListWorkRequestsResponse")
	}
	return
}

// listWorkRequests implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) listWorkRequests(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodGet, "/workRequests/")
	if err != nil {
		return nil, err
	}

	var response ListWorkRequestsResponse
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

// RemoveUserFromGroup Removes a user from a group by deleting the corresponding `UserGroupMembership`.
func (client IdentityClient) RemoveUserFromGroup(ctx context.Context, request RemoveUserFromGroupRequest) (response RemoveUserFromGroupResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.removeUserFromGroup, policy)
	if err != nil {
		if ociResponse != nil {
			response = RemoveUserFromGroupResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(RemoveUserFromGroupResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into RemoveUserFromGroupResponse")
	}
	return
}

// removeUserFromGroup implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) removeUserFromGroup(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodDelete, "/userGroupMemberships/{userGroupMembershipId}")
	if err != nil {
		return nil, err
	}

	var response RemoveUserFromGroupResponse
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

// ResetIdpScimClient Resets the OAuth2 client credentials for the SCIM client associated with this identity provider.
func (client IdentityClient) ResetIdpScimClient(ctx context.Context, request ResetIdpScimClientRequest) (response ResetIdpScimClientResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.resetIdpScimClient, policy)
	if err != nil {
		if ociResponse != nil {
			response = ResetIdpScimClientResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(ResetIdpScimClientResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into ResetIdpScimClientResponse")
	}
	return
}

// resetIdpScimClient implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) resetIdpScimClient(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/identityProviders/{identityProviderId}/actions/resetScimClient/")
	if err != nil {
		return nil, err
	}

	var response ResetIdpScimClientResponse
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

// UpdateAuthToken Updates the specified auth token's description.
func (client IdentityClient) UpdateAuthToken(ctx context.Context, request UpdateAuthTokenRequest) (response UpdateAuthTokenResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateAuthToken, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateAuthTokenResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateAuthTokenResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateAuthTokenResponse")
	}
	return
}

// updateAuthToken implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateAuthToken(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/users/{userId}/authTokens/{authTokenId}")
	if err != nil {
		return nil, err
	}

	var response UpdateAuthTokenResponse
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

// UpdateAuthenticationPolicy Updates authentication policy for the specified tenancy
func (client IdentityClient) UpdateAuthenticationPolicy(ctx context.Context, request UpdateAuthenticationPolicyRequest) (response UpdateAuthenticationPolicyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateAuthenticationPolicy, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateAuthenticationPolicyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateAuthenticationPolicyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateAuthenticationPolicyResponse")
	}
	return
}

// updateAuthenticationPolicy implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateAuthenticationPolicy(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/authenticationPolicies/{compartmentId}")
	if err != nil {
		return nil, err
	}

	var response UpdateAuthenticationPolicyResponse
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

// UpdateCompartment Updates the specified compartment's description or name. You can't update the root compartment.
func (client IdentityClient) UpdateCompartment(ctx context.Context, request UpdateCompartmentRequest) (response UpdateCompartmentResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateCompartment, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateCompartmentResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateCompartmentResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateCompartmentResponse")
	}
	return
}

// updateCompartment implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateCompartment(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/compartments/{compartmentId}")
	if err != nil {
		return nil, err
	}

	var response UpdateCompartmentResponse
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

// UpdateCustomerSecretKey Updates the specified secret key's description.
func (client IdentityClient) UpdateCustomerSecretKey(ctx context.Context, request UpdateCustomerSecretKeyRequest) (response UpdateCustomerSecretKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateCustomerSecretKey, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateCustomerSecretKeyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateCustomerSecretKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateCustomerSecretKeyResponse")
	}
	return
}

// updateCustomerSecretKey implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateCustomerSecretKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/users/{userId}/customerSecretKeys/{customerSecretKeyId}")
	if err != nil {
		return nil, err
	}

	var response UpdateCustomerSecretKeyResponse
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

// UpdateDynamicGroup Updates the specified dynamic group.
func (client IdentityClient) UpdateDynamicGroup(ctx context.Context, request UpdateDynamicGroupRequest) (response UpdateDynamicGroupResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateDynamicGroup, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateDynamicGroupResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateDynamicGroupResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateDynamicGroupResponse")
	}
	return
}

// updateDynamicGroup implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateDynamicGroup(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/dynamicGroups/{dynamicGroupId}")
	if err != nil {
		return nil, err
	}

	var response UpdateDynamicGroupResponse
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

// UpdateGroup Updates the specified group.
func (client IdentityClient) UpdateGroup(ctx context.Context, request UpdateGroupRequest) (response UpdateGroupResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateGroup, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateGroupResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateGroupResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateGroupResponse")
	}
	return
}

// updateGroup implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateGroup(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/groups/{groupId}")
	if err != nil {
		return nil, err
	}

	var response UpdateGroupResponse
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

// UpdateIdentityProvider Updates the specified identity provider.
func (client IdentityClient) UpdateIdentityProvider(ctx context.Context, request UpdateIdentityProviderRequest) (response UpdateIdentityProviderResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateIdentityProvider, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateIdentityProviderResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateIdentityProviderResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateIdentityProviderResponse")
	}
	return
}

// updateIdentityProvider implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateIdentityProvider(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/identityProviders/{identityProviderId}")
	if err != nil {
		return nil, err
	}

	var response UpdateIdentityProviderResponse
	var httpResponse *http.Response
	httpResponse, err = client.Call(ctx, &httpRequest)
	defer common.CloseBodyIfValid(httpResponse)
	response.RawResponse = httpResponse
	if err != nil {
		return response, err
	}

	err = common.UnmarshalResponseWithPolymorphicBody(httpResponse, &response, &identityprovider{})
	return response, err
}

// UpdateIdpGroupMapping Updates the specified group mapping.
func (client IdentityClient) UpdateIdpGroupMapping(ctx context.Context, request UpdateIdpGroupMappingRequest) (response UpdateIdpGroupMappingResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateIdpGroupMapping, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateIdpGroupMappingResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateIdpGroupMappingResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateIdpGroupMappingResponse")
	}
	return
}

// updateIdpGroupMapping implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateIdpGroupMapping(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/identityProviders/{identityProviderId}/groupMappings/{mappingId}")
	if err != nil {
		return nil, err
	}

	var response UpdateIdpGroupMappingResponse
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

// UpdatePolicy Updates the specified policy. You can update the description or the policy statements themselves.
// Policy changes take effect typically within 10 seconds.
func (client IdentityClient) UpdatePolicy(ctx context.Context, request UpdatePolicyRequest) (response UpdatePolicyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updatePolicy, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdatePolicyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdatePolicyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdatePolicyResponse")
	}
	return
}

// updatePolicy implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updatePolicy(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/policies/{policyId}")
	if err != nil {
		return nil, err
	}

	var response UpdatePolicyResponse
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

// UpdateSmtpCredential Updates the specified SMTP credential's description.
func (client IdentityClient) UpdateSmtpCredential(ctx context.Context, request UpdateSmtpCredentialRequest) (response UpdateSmtpCredentialResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateSmtpCredential, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateSmtpCredentialResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateSmtpCredentialResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateSmtpCredentialResponse")
	}
	return
}

// updateSmtpCredential implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateSmtpCredential(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/users/{userId}/smtpCredentials/{smtpCredentialId}")
	if err != nil {
		return nil, err
	}

	var response UpdateSmtpCredentialResponse
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

// UpdateSwiftPassword **Deprecated. Use UpdateAuthToken instead.**
// Updates the specified Swift password's description.
func (client IdentityClient) UpdateSwiftPassword(ctx context.Context, request UpdateSwiftPasswordRequest) (response UpdateSwiftPasswordResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateSwiftPassword, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateSwiftPasswordResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateSwiftPasswordResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateSwiftPasswordResponse")
	}
	return
}

// updateSwiftPassword implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateSwiftPassword(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/users/{userId}/swiftPasswords/{swiftPasswordId}")
	if err != nil {
		return nil, err
	}

	var response UpdateSwiftPasswordResponse
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

// UpdateTag Updates the the specified tag definition. You can update `description`, and `isRetired`.
func (client IdentityClient) UpdateTag(ctx context.Context, request UpdateTagRequest) (response UpdateTagResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateTag, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateTagResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateTagResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateTagResponse")
	}
	return
}

// updateTag implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateTag(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/tagNamespaces/{tagNamespaceId}/tags/{tagName}")
	if err != nil {
		return nil, err
	}

	var response UpdateTagResponse
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

// UpdateTagDefault Updates the specified tag default. You can update the following field: `value`.
func (client IdentityClient) UpdateTagDefault(ctx context.Context, request UpdateTagDefaultRequest) (response UpdateTagDefaultResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateTagDefault, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateTagDefaultResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateTagDefaultResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateTagDefaultResponse")
	}
	return
}

// updateTagDefault implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateTagDefault(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/tagDefaults/{tagDefaultId}")
	if err != nil {
		return nil, err
	}

	var response UpdateTagDefaultResponse
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

// UpdateTagNamespace Updates the the specified tag namespace. You can't update the namespace name.
// Updating `isRetired` to 'true' retires the namespace and all the tag definitions in the namespace. Reactivating a
// namespace (changing `isRetired` from 'true' to 'false') does not reactivate tag definitions.
// To reactivate the tag definitions, you must reactivate each one indvidually *after* you reactivate the namespace,
// using UpdateTag. For more information about retiring tag namespaces, see
// Retiring Key Definitions and Namespace Definitions (https://docs.cloud.oracle.com/Content/Identity/Concepts/taggingoverview.htm#Retiring).
// You can't add a namespace with the same name as a retired namespace in the same tenancy.
func (client IdentityClient) UpdateTagNamespace(ctx context.Context, request UpdateTagNamespaceRequest) (response UpdateTagNamespaceResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateTagNamespace, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateTagNamespaceResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateTagNamespaceResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateTagNamespaceResponse")
	}
	return
}

// updateTagNamespace implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateTagNamespace(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/tagNamespaces/{tagNamespaceId}")
	if err != nil {
		return nil, err
	}

	var response UpdateTagNamespaceResponse
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

// UpdateUser Updates the description of the specified user.
func (client IdentityClient) UpdateUser(ctx context.Context, request UpdateUserRequest) (response UpdateUserResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateUser, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateUserResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateUserResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateUserResponse")
	}
	return
}

// updateUser implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateUser(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/users/{userId}")
	if err != nil {
		return nil, err
	}

	var response UpdateUserResponse
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

// UpdateUserCapabilities Updates the capabilities of the specified user.
func (client IdentityClient) UpdateUserCapabilities(ctx context.Context, request UpdateUserCapabilitiesRequest) (response UpdateUserCapabilitiesResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateUserCapabilities, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateUserCapabilitiesResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateUserCapabilitiesResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateUserCapabilitiesResponse")
	}
	return
}

// updateUserCapabilities implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateUserCapabilities(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/users/{userId}/capabilities/")
	if err != nil {
		return nil, err
	}

	var response UpdateUserCapabilitiesResponse
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

// UpdateUserState Updates the state of the specified user.
func (client IdentityClient) UpdateUserState(ctx context.Context, request UpdateUserStateRequest) (response UpdateUserStateResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}
	ociResponse, err = common.Retry(ctx, request, client.updateUserState, policy)
	if err != nil {
		if ociResponse != nil {
			response = UpdateUserStateResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UpdateUserStateResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UpdateUserStateResponse")
	}
	return
}

// updateUserState implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) updateUserState(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPut, "/users/{userId}/state/")
	if err != nil {
		return nil, err
	}

	var response UpdateUserStateResponse
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

// UploadApiKey Uploads an API signing key for the specified user.
// Every user has permission to use this operation to upload a key for *their own user ID*. An
// administrator in your organization does not need to write a policy to give users this ability.
// To compare, administrators who have permission to the tenancy can use this operation to upload a
// key for any user, including themselves.
// **Important:** Even though you have permission to upload an API key, you might not yet
// have permission to do much else. If you try calling an operation unrelated to your own credential
// management (e.g., `ListUsers`, `LaunchInstance`) and receive an "unauthorized" error,
// check with an administrator to confirm which IAM Service group(s) you're in and what access
// you have. Also confirm you're working in the correct compartment.
// After you send your request, the new object's `lifecycleState` will temporarily be CREATING. Before using
// the object, first make sure its `lifecycleState` has changed to ACTIVE.
func (client IdentityClient) UploadApiKey(ctx context.Context, request UploadApiKeyRequest) (response UploadApiKeyResponse, err error) {
	var ociResponse common.OCIResponse
	policy := common.NoRetryPolicy()
	if request.RetryPolicy() != nil {
		policy = *request.RetryPolicy()
	}

	if !(request.OpcRetryToken != nil && *request.OpcRetryToken != "") {
		request.OpcRetryToken = common.String(common.RetryToken())
	}

	ociResponse, err = common.Retry(ctx, request, client.uploadApiKey, policy)
	if err != nil {
		if ociResponse != nil {
			response = UploadApiKeyResponse{RawResponse: ociResponse.HTTPResponse()}
		}
		return
	}
	if convertedResponse, ok := ociResponse.(UploadApiKeyResponse); ok {
		response = convertedResponse
	} else {
		err = fmt.Errorf("failed to convert OCIResponse into UploadApiKeyResponse")
	}
	return
}

// uploadApiKey implements the OCIOperation interface (enables retrying operations)
func (client IdentityClient) uploadApiKey(ctx context.Context, request common.OCIRequest) (common.OCIResponse, error) {
	httpRequest, err := request.HTTPRequest(http.MethodPost, "/users/{userId}/apiKeys/")
	if err != nil {
		return nil, err
	}

	var response UploadApiKeyResponse
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
