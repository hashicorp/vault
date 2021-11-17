package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-01-01-preview/authorization"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/go-multierror"
)

const (
	// DefaultGraphMicrosoftComURI is the default URI used for the service MS Graph API
	DefaultGraphMicrosoftComURI = "https://graph.microsoft.com"
)

var _ ApplicationsClient = (*AppClient)(nil)
var _ GroupsClient = (*AppClient)(nil)
var _ ServicePrincipalClient = (*AppClient)(nil)

type AppClient struct {
	client authorization.BaseClient
}

func NewMSGraphApplicationClient(subscriptionId string, userAgentExtension string, auth autorest.Authorizer) (*AppClient, error) {
	client := authorization.NewWithBaseURI(DefaultGraphMicrosoftComURI, subscriptionId)
	client.Authorizer = auth

	if userAgentExtension != "" {
		err := client.AddToUserAgent(userAgentExtension)
		if err != nil {
			return nil, fmt.Errorf("failed to add extension to user agent")
		}
	}

	ac := &AppClient{
		client: client,
	}
	return ac, nil
}

func (c *AppClient) AddToUserAgent(extension string) error {
	return c.client.AddToUserAgent(extension)
}

func (c *AppClient) GetApplication(ctx context.Context, applicationObjectID string) (ApplicationResult, error) {
	var result ApplicationResult
	req, err := c.getApplicationPreparer(ctx, applicationObjectID)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "provider", "GetApplication", nil, "Failure preparing request")
	}

	resp, err := c.getApplicationSender(req)
	if err != nil {
		result = ApplicationResult{
			Response: autorest.Response{Response: resp},
		}
		return result, autorest.NewErrorWithError(err, "provider", "GetApplication", resp, "Failure sending request")
	}

	result, err = c.getApplicationResponder(resp)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "provider", "GetApplication", resp, "Failure responding to request")
	}

	return result, nil
}

type listApplicationsResponse struct {
	Value []ApplicationResult `json:"value"`
}

func (c *AppClient) ListApplications(ctx context.Context, filter string) ([]ApplicationResult, error) {
	filterArgs := url.Values{}
	if filter != "" {
		filterArgs.Set("$filter", filter)
	}
	preparer := c.GetPreparer(
		autorest.AsGet(),
		autorest.WithPath(fmt.Sprintf("/v1.0/applications?%s", filterArgs.Encode())),
	)
	listAppResp := listApplicationsResponse{}
	err := c.SendRequest(ctx, preparer,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&listAppResp),
	)
	if err != nil {
		return nil, err
	}

	return listAppResp.Value, nil
}

// CreateApplication create a new Azure application object.
func (c *AppClient) CreateApplication(ctx context.Context, displayName string) (ApplicationResult, error) {
	var result ApplicationResult

	req, err := c.createApplicationPreparer(ctx, displayName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "provider", "CreateApplication", nil, "Failure preparing request")
	}

	resp, err := c.createApplicationSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "provider", "CreateApplication", resp, "Failure sending request")
	}

	result, err = c.createApplicationResponder(resp)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "provider", "CreateApplication", resp, "Failure responding to request")

	}

	return result, nil
}

// DeleteApplication deletes an Azure application object.
// This will in turn remove the service principal (but not the role assignments).
func (c *AppClient) DeleteApplication(ctx context.Context, applicationObjectID string) error {
	req, err := c.deleteApplicationPreparer(ctx, applicationObjectID)
	if err != nil {
		return autorest.NewErrorWithError(err, "provider", "DeleteApplication", nil, "Failure preparing request")
	}

	resp, err := c.deleteApplicationSender(req)
	if err != nil {
		return autorest.NewErrorWithError(err, "provider", "DeleteApplication", resp, "Failure sending request")
	}

	err = autorest.Respond(
		resp,
		c.client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusNotFound),
		autorest.ByClosing())

	if err != nil {
		return autorest.NewErrorWithError(err, "provider", "DeleteApplication", resp, "Failure responding to request")
	}
	return nil
}

func (c *AppClient) AddApplicationPassword(ctx context.Context, applicationObjectID string, displayName string, endDateTime time.Time) (PasswordCredentialResult, error) {
	req, err := c.addPasswordPreparer(ctx, applicationObjectID, displayName, date.Time{endDateTime})
	if err != nil {
		return PasswordCredentialResult{}, autorest.NewErrorWithError(err, "provider", "AddApplicationPassword", nil, "Failure preparing request")
	}

	resp, err := c.addPasswordSender(req)
	if err != nil {
		result := PasswordCredentialResult{
			Response: autorest.Response{Response: resp},
		}
		return result, autorest.NewErrorWithError(err, "provider", "AddApplicationPassword", resp, "Failure sending request")
	}

	result, err := c.addPasswordResponder(resp)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "provider", "AddApplicationPassword", resp, "Failure responding to request")
	}

	return result, nil
}

func (c *AppClient) RemoveApplicationPassword(ctx context.Context, applicationObjectID string, keyID string) error {
	req, err := c.removePasswordPreparer(ctx, applicationObjectID, keyID)
	if err != nil {
		return autorest.NewErrorWithError(err, "provider", "RemoveApplicationPassword", nil, "Failure preparing request")
	}

	resp, err := c.removePasswordSender(req)
	if err != nil {
		return autorest.NewErrorWithError(err, "provider", "RemoveApplicationPassword", resp, "Failure sending request")
	}

	_, err = c.removePasswordResponder(resp)
	if err != nil {
		return autorest.NewErrorWithError(err, "provider", "RemoveApplicationPassword", resp, "Failure responding to request")
	}

	return nil
}

func (c AppClient) getApplicationPreparer(ctx context.Context, applicationObjectID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"applicationObjectId": autorest.Encode("path", applicationObjectID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.client.BaseURI),
		autorest.WithPathParameters("/v1.0/applications/{applicationObjectId}", pathParameters),
		c.client.WithAuthorization())
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (c AppClient) getApplicationSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(c.client.RetryAttempts, c.client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(c.client, req, sd...)
}

func (c AppClient) getApplicationResponder(resp *http.Response) (ApplicationResult, error) {
	var result ApplicationResult
	err := autorest.Respond(
		resp,
		c.client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return result, err
}

func (c AppClient) addPasswordPreparer(ctx context.Context, applicationObjectID string, displayName string, endDateTime date.Time) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"applicationObjectId": autorest.Encode("path", applicationObjectID),
	}

	parameters := struct {
		PasswordCredential *PasswordCredential `json:"passwordCredential"`
	}{
		PasswordCredential: &PasswordCredential{
			DisplayName: to.StringPtr(displayName),
			EndDate:     &endDateTime,
		},
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.client.BaseURI),
		autorest.WithPathParameters("/v1.0/applications/{applicationObjectId}/addPassword", pathParameters),
		autorest.WithJSON(parameters),
		c.client.WithAuthorization())
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (c AppClient) addPasswordSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(c.client.RetryAttempts, c.client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(c.client, req, sd...)
}

func (c AppClient) addPasswordResponder(resp *http.Response) (PasswordCredentialResult, error) {
	var result PasswordCredentialResult
	err := autorest.Respond(
		resp,
		c.client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return result, err
}

func (c AppClient) removePasswordPreparer(ctx context.Context, applicationObjectID string, keyID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"applicationObjectId": autorest.Encode("path", applicationObjectID),
	}

	parameters := struct {
		KeyID string `json:"keyId"`
	}{
		KeyID: keyID,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.client.BaseURI),
		autorest.WithPathParameters("/v1.0/applications/{applicationObjectId}/removePassword", pathParameters),
		autorest.WithJSON(parameters),
		c.client.WithAuthorization())
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (c AppClient) removePasswordSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(c.client.RetryAttempts, c.client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(c.client, req, sd...)
}

func (c AppClient) removePasswordResponder(resp *http.Response) (autorest.Response, error) {
	var result autorest.Response
	err := autorest.Respond(
		resp,
		c.client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = resp
	return result, err
}

func (c AppClient) createApplicationPreparer(ctx context.Context, displayName string) (*http.Request, error) {
	parameters := struct {
		DisplayName *string `json:"displayName"`
	}{
		DisplayName: to.StringPtr(displayName),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.client.BaseURI),
		autorest.WithPath("/v1.0/applications"),
		autorest.WithJSON(parameters),
		c.client.WithAuthorization())
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (c AppClient) createApplicationSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(c.client.RetryAttempts, c.client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(c.client, req, sd...)
}

func (c AppClient) createApplicationResponder(resp *http.Response) (ApplicationResult, error) {
	var result ApplicationResult
	err := autorest.Respond(
		resp,
		c.client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return result, err
}

func (c AppClient) deleteApplicationPreparer(ctx context.Context, applicationObjectID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"applicationObjectId": autorest.Encode("path", applicationObjectID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.client.BaseURI),
		autorest.WithPathParameters("/v1.0/applications/{applicationObjectId}", pathParameters),
		c.client.WithAuthorization())
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

func (c AppClient) deleteApplicationSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(c.client.RetryAttempts, c.client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(c.client, req, sd...)
}

func (c AppClient) AddGroupMember(ctx context.Context, groupObjectID string, memberObjectID string) error {
	if groupObjectID == "" {
		return fmt.Errorf("missing groupObjectID")
	}
	pathParams := map[string]interface{}{
		"groupObjectID": groupObjectID,
	}
	body := map[string]interface{}{
		"@odata.id": fmt.Sprintf("%s/v1.0/directoryObjects/%s", DefaultGraphMicrosoftComURI, memberObjectID),
	}
	preparer := c.GetPreparer(
		autorest.AsPost(),
		autorest.WithPathParameters("/v1.0/groups/{groupObjectID}/members/$ref", pathParams),
		autorest.WithJSON(body),
	)
	return c.SendRequest(ctx, preparer, azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent))
}

func (c AppClient) RemoveGroupMember(ctx context.Context, groupObjectID, memberObjectID string) error {
	if groupObjectID == "" {
		return fmt.Errorf("missing groupObjectID")
	}
	if memberObjectID == "" {
		return fmt.Errorf("missing memberObjectID")
	}
	pathParams := map[string]interface{}{
		"groupObjectID":  groupObjectID,
		"memberObjectID": memberObjectID,
	}

	preparer := c.GetPreparer(
		autorest.AsDelete(),
		autorest.WithPathParameters("/v1.0/groups/{groupObjectID}/members/{memberObjectID}/$ref", pathParams),
	)
	return c.SendRequest(ctx, preparer, azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent))
}

// groupResponse is a struct representation of the data we care about coming back from
// the ms-graph API. This is not the same as `Group` because this information is
// slightly different from the AAD implementation and there should be an abstraction
// between the ms-graph API itself and the API this package presents.
type groupResponse struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

func (c AppClient) GetGroup(ctx context.Context, groupID string) (Group, error) {
	if groupID == "" {
		return Group{}, fmt.Errorf("missing groupID")
	}
	pathParams := map[string]interface{}{
		"groupID": groupID,
	}

	preparer := c.GetPreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("/v1.0/groups/{groupID}", pathParams),
	)

	groupResp := groupResponse{}
	err := c.SendRequest(ctx, preparer,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByUnmarshallingJSON(&groupResp),
	)
	if err != nil {
		return Group{}, err
	}

	group := Group{
		ID:          groupResp.ID,
		DisplayName: groupResp.DisplayName,
	}

	return group, nil
}

// listGroupsResponse is a struct representation of the data we care about
// coming back from the ms-graph API
type listGroupsResponse struct {
	Groups []groupResponse `json:"value"`
}

func (c AppClient) ListGroups(ctx context.Context, filter string) ([]Group, error) {
	filterArgs := url.Values{}
	if filter != "" {
		filterArgs.Set("$filter", filter)
	}

	preparer := c.GetPreparer(
		autorest.AsGet(),
		autorest.WithPath(fmt.Sprintf("/v1.0/groups?%s", filterArgs.Encode())),
	)

	respBody := listGroupsResponse{}
	err := c.SendRequest(ctx, preparer,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByUnmarshallingJSON(&respBody),
	)
	if err != nil {
		return nil, err
	}

	groups := []Group{}
	for _, rawGroup := range respBody.Groups {
		if rawGroup.ID == "" {
			return nil, fmt.Errorf("missing group ID from response")
		}

		group := Group{
			ID:          rawGroup.ID,
			DisplayName: rawGroup.DisplayName,
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (c *AppClient) CreateServicePrincipal(ctx context.Context, appID string, startDate time.Time, endDate time.Time) (string, string, error) {
	spID, err := c.createServicePrincipal(ctx, appID)
	if err != nil {
		return "", "", err
	}
	password, err := c.setPasswordForServicePrincipal(ctx, spID, startDate, endDate)
	if err != nil {
		dErr := c.deleteServicePrincipal(ctx, spID)
		merr := multierror.Append(err, dErr)
		return "", "", merr.ErrorOrNil()
	}
	return spID, password, nil
}

func (c *AppClient) createServicePrincipal(ctx context.Context, appID string) (string, error) {
	body := map[string]interface{}{
		"appId":          appID,
		"accountEnabled": true,
	}
	preparer := c.GetPreparer(
		autorest.AsPost(),
		autorest.WithPath("/v1.0/servicePrincipals"),
		autorest.WithJSON(body),
	)

	respBody := createServicePrincipalResponse{}
	err := c.SendRequest(ctx, preparer,
		autorest.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&respBody),
	)
	if err != nil {
		return "", err
	}

	return respBody.ID, nil
}

func (c *AppClient) setPasswordForServicePrincipal(ctx context.Context, spID string, startDate time.Time, endDate time.Time) (string, error) {
	pathParams := map[string]interface{}{
		"id": spID,
	}
	reqBody := map[string]interface{}{
		"startDateTime": startDate.UTC().Format("2006-01-02T15:04:05Z"),
		"endDateTime":   endDate.UTC().Format("2006-01-02T15:04:05Z"),
	}

	preparer := c.GetPreparer(
		autorest.AsPost(),
		autorest.WithPathParameters("/v1.0/servicePrincipals/{id}/addPassword", pathParams),
		autorest.WithJSON(reqBody),
	)

	respBody := PasswordCredential{}
	err := c.SendRequest(ctx, preparer,
		autorest.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByUnmarshallingJSON(&respBody),
	)
	if err != nil {
		return "", err
	}
	return *respBody.SecretText, nil
}

type createServicePrincipalResponse struct {
	ID string `json:"id"`
}

func (c *AppClient) deleteServicePrincipal(ctx context.Context, spID string) error {
	pathParams := map[string]interface{}{
		"id": spID,
	}

	preparer := c.GetPreparer(
		autorest.AsDelete(),
		autorest.WithPathParameters("/v1.0/servicePrincipals/{id}", pathParams),
	)

	return c.SendRequest(ctx, preparer, autorest.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent))
}

func (c *AppClient) GetPreparer(prepareDecorators ...autorest.PrepareDecorator) autorest.Preparer {
	decs := []autorest.PrepareDecorator{
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.WithBaseURL(c.client.BaseURI),
		c.client.WithAuthorization(),
	}
	decs = append(decs, prepareDecorators...)
	preparer := autorest.CreatePreparer(decs...)
	return preparer
}

func (c *AppClient) SendRequest(ctx context.Context, preparer autorest.Preparer, respDecs ...autorest.RespondDecorator) error {
	req, err := preparer.Prepare((&http.Request{}).WithContext(ctx))
	if err != nil {
		return err
	}

	sender := autorest.GetSendDecorators(req.Context(),
		autorest.DoRetryForStatusCodes(c.client.RetryAttempts, c.client.RetryDuration, autorest.StatusCodesForRetry...),
	)
	resp, err := autorest.SendWithSender(c.client, req, sender...)
	if err != nil {
		return err
	}

	// Put ByInspecting() before any provided decorators
	respDecs = append([]autorest.RespondDecorator{c.client.ByInspecting()}, respDecs...)
	respDecs = append(respDecs, autorest.ByClosing())

	return autorest.Respond(resp, respDecs...)
}
