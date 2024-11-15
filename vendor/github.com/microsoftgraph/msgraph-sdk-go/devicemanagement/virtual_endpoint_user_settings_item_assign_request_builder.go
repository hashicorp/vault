package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// VirtualEndpointUserSettingsItemAssignRequestBuilder provides operations to call the assign method.
type VirtualEndpointUserSettingsItemAssignRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// VirtualEndpointUserSettingsItemAssignRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEndpointUserSettingsItemAssignRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewVirtualEndpointUserSettingsItemAssignRequestBuilderInternal instantiates a new VirtualEndpointUserSettingsItemAssignRequestBuilder and sets the default values.
func NewVirtualEndpointUserSettingsItemAssignRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEndpointUserSettingsItemAssignRequestBuilder) {
    m := &VirtualEndpointUserSettingsItemAssignRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/virtualEndpoint/userSettings/{cloudPcUserSetting%2Did}/assign", pathParameters),
    }
    return m
}
// NewVirtualEndpointUserSettingsItemAssignRequestBuilder instantiates a new VirtualEndpointUserSettingsItemAssignRequestBuilder and sets the default values.
func NewVirtualEndpointUserSettingsItemAssignRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEndpointUserSettingsItemAssignRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewVirtualEndpointUserSettingsItemAssignRequestBuilderInternal(urlParams, requestAdapter)
}
// Post assign a cloudPcUserSetting to user groups.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/cloudpcusersetting-assign?view=graph-rest-1.0
func (m *VirtualEndpointUserSettingsItemAssignRequestBuilder) Post(ctx context.Context, body VirtualEndpointUserSettingsItemAssignPostRequestBodyable, requestConfiguration *VirtualEndpointUserSettingsItemAssignRequestBuilderPostRequestConfiguration)(error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// ToPostRequestInformation assign a cloudPcUserSetting to user groups.
// returns a *RequestInformation when successful
func (m *VirtualEndpointUserSettingsItemAssignRequestBuilder) ToPostRequestInformation(ctx context.Context, body VirtualEndpointUserSettingsItemAssignPostRequestBodyable, requestConfiguration *VirtualEndpointUserSettingsItemAssignRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    err := requestInfo.SetContentFromParsable(ctx, m.BaseRequestBuilder.RequestAdapter, "application/json", body)
    if err != nil {
        return nil, err
    }
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *VirtualEndpointUserSettingsItemAssignRequestBuilder when successful
func (m *VirtualEndpointUserSettingsItemAssignRequestBuilder) WithUrl(rawUrl string)(*VirtualEndpointUserSettingsItemAssignRequestBuilder) {
    return NewVirtualEndpointUserSettingsItemAssignRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
