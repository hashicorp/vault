package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder provides operations to manage the deviceMalwareStates property of the microsoft.graph.windowsMalwareInformation entity.
type WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderGetQueryParameters read properties and relationships of the malwareStateForWindowsDevice object.
type WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderGetQueryParameters
}
// WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewWindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderInternal instantiates a new WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder and sets the default values.
func NewWindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder) {
    m := &WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/windowsMalwareInformation/{windowsMalwareInformation%2Did}/deviceMalwareStates/{malwareStateForWindowsDevice%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewWindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder instantiates a new WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder and sets the default values.
func NewWindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewWindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete deletes a malwareStateForWindowsDevice.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-devices-malwarestateforwindowsdevice-delete?view=graph-rest-1.0
func (m *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
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
// Get read properties and relationships of the malwareStateForWindowsDevice object.
// returns a MalwareStateForWindowsDeviceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-devices-malwarestateforwindowsdevice-get?view=graph-rest-1.0
func (m *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder) Get(ctx context.Context, requestConfiguration *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MalwareStateForWindowsDeviceable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMalwareStateForWindowsDeviceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MalwareStateForWindowsDeviceable), nil
}
// Patch update the properties of a malwareStateForWindowsDevice object.
// returns a MalwareStateForWindowsDeviceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-devices-malwarestateforwindowsdevice-update?view=graph-rest-1.0
func (m *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MalwareStateForWindowsDeviceable, requestConfiguration *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MalwareStateForWindowsDeviceable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMalwareStateForWindowsDeviceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MalwareStateForWindowsDeviceable), nil
}
// ToDeleteRequestInformation deletes a malwareStateForWindowsDevice.
// returns a *RequestInformation when successful
func (m *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read properties and relationships of the malwareStateForWindowsDevice object.
// returns a *RequestInformation when successful
func (m *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToPatchRequestInformation update the properties of a malwareStateForWindowsDevice object.
// returns a *RequestInformation when successful
func (m *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MalwareStateForWindowsDeviceable, requestConfiguration *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// returns a *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder when successful
func (m *WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder) WithUrl(rawUrl string)(*WindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder) {
    return NewWindowsMalwareInformationItemDeviceMalwareStatesMalwareStateForWindowsDeviceItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
