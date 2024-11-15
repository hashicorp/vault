package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder provides operations to manage the detectedMalwareState property of the microsoft.graph.windowsProtectionState entity.
type ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetQueryParameters device malware list
type ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetQueryParameters
}
// ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderInternal instantiates a new ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder and sets the default values.
func NewManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) {
    m := &ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/managedDevices/{managedDevice%2Did}/windowsProtectionState/detectedMalwareState/{windowsDeviceMalwareState%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder instantiates a new ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder and sets the default values.
func NewManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property detectedMalwareState for deviceManagement
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderDeleteRequestConfiguration)(error) {
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
func (m *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsDeviceMalwareStateable, error) {
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
// Patch update the navigation property detectedMalwareState in deviceManagement
// returns a WindowsDeviceMalwareStateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsDeviceMalwareStateable, requestConfiguration *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsDeviceMalwareStateable, error) {
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
// ToDeleteRequestInformation delete navigation property detectedMalwareState for deviceManagement
// returns a *RequestInformation when successful
func (m *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property detectedMalwareState in deviceManagement
// returns a *RequestInformation when successful
func (m *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsDeviceMalwareStateable, requestConfiguration *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder when successful
func (m *ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) WithUrl(rawUrl string)(*ManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder) {
    return NewManagedDevicesItemWindowsProtectionStateDetectedMalwareStateWindowsDeviceMalwareStateItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
