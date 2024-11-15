package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder provides operations to call the import method.
type ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilderInternal instantiates a new ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder and sets the default values.
func NewImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder) {
    m := &ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/importedWindowsAutopilotDeviceIdentities/import", pathParameters),
    }
    return m
}
// NewImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder instantiates a new ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder and sets the default values.
func NewImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilderInternal(urlParams, requestAdapter)
}
// Post not yet documented
// Deprecated: This method is obsolete. Use PostAsImportPostResponse instead.
// returns a ImportedWindowsAutopilotDeviceIdentitiesImportResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-enrollment-importedwindowsautopilotdeviceidentity-import?view=graph-rest-1.0
func (m *ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder) Post(ctx context.Context, body ImportedWindowsAutopilotDeviceIdentitiesImportPostRequestBodyable, requestConfiguration *ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilderPostRequestConfiguration)(ImportedWindowsAutopilotDeviceIdentitiesImportResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateImportedWindowsAutopilotDeviceIdentitiesImportResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ImportedWindowsAutopilotDeviceIdentitiesImportResponseable), nil
}
// PostAsImportPostResponse not yet documented
// returns a ImportedWindowsAutopilotDeviceIdentitiesImportPostResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-enrollment-importedwindowsautopilotdeviceidentity-import?view=graph-rest-1.0
func (m *ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder) PostAsImportPostResponse(ctx context.Context, body ImportedWindowsAutopilotDeviceIdentitiesImportPostRequestBodyable, requestConfiguration *ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilderPostRequestConfiguration)(ImportedWindowsAutopilotDeviceIdentitiesImportPostResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateImportedWindowsAutopilotDeviceIdentitiesImportPostResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ImportedWindowsAutopilotDeviceIdentitiesImportPostResponseable), nil
}
// ToPostRequestInformation not yet documented
// returns a *RequestInformation when successful
func (m *ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder) ToPostRequestInformation(ctx context.Context, body ImportedWindowsAutopilotDeviceIdentitiesImportPostRequestBodyable, requestConfiguration *ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder when successful
func (m *ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder) WithUrl(rawUrl string)(*ImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder) {
    return NewImportedWindowsAutopilotDeviceIdentitiesImportRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
