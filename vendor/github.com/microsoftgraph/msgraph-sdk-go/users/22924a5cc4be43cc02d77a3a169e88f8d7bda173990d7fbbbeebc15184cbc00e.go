package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder provides operations to manage the detectedMalwareState property of the microsoft.graph.windowsProtectionState entity.
type ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetQueryParameters device malware list
type ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetQueryParameters
}
// ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderInternal instantiates a new ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder and sets the default values.
func NewItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) {
    m := &ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/managedDevices/{managedDevice%2Did}/windowsProtectionState/detectedMalwareState/{windowsDeviceMalwareState%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder instantiates a new ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder and sets the default values.
func NewItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property detectedMalwareState for users
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get device malware list
// returns a WindowsDeviceMalwareStateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsDeviceMalwareStateable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWindowsDeviceMalwareStateFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsDeviceMalwareStateable), nil
}
// Patch update the navigation property detectedMalwareState in users
// returns a WindowsDeviceMalwareStateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsDeviceMalwareStateable, requestConfiguration *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsDeviceMalwareStateable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWindowsDeviceMalwareStateFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsDeviceMalwareStateable), nil
}
// ToDeleteRequestInformation delete navigation property detectedMalwareState for users
// returns a *RequestInformation when successful
func (m *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation device malware list
// returns a *RequestInformation when successful
func (m *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property detectedMalwareState in users
// returns a *RequestInformation when successful
func (m *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsDeviceMalwareStateable, requestConfiguration *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder when successful
func (m *ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) WithUrl(rawUrl string)(*ItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) {
    return NewItemManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
