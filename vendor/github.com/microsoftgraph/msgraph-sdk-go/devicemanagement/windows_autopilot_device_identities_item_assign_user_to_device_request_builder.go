package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder provides operations to call the assignUserToDevice method.
type WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewWindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilderInternal instantiates a new WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder and sets the default values.
func NewWindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder) {
    m := &WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/windowsAutopilotDeviceIdentities/{windowsAutopilotDeviceIdentity%2Did}/assignUserToDevice", pathParameters),
    }
    return m
}
// NewWindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder instantiates a new WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder and sets the default values.
func NewWindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewWindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilderInternal(urlParams, requestAdapter)
}
// Post assigns user to Autopilot devices.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-enrollment-windowsautopilotdeviceidentity-assignusertodevice?view=graph-rest-1.0
func (m *WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder) Post(ctx context.Context, body WindowsAutopilotDeviceIdentitiesItemAssignUserToDevicePostRequestBodyable, requestConfiguration *WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation assigns user to Autopilot devices.
// returns a *RequestInformation when successful
func (m *WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder) ToPostRequestInformation(ctx context.Context, body WindowsAutopilotDeviceIdentitiesItemAssignUserToDevicePostRequestBodyable, requestConfiguration *WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder when successful
func (m *WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder) WithUrl(rawUrl string)(*WindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder) {
    return NewWindowsAutopilotDeviceIdentitiesItemAssignUserToDeviceRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
