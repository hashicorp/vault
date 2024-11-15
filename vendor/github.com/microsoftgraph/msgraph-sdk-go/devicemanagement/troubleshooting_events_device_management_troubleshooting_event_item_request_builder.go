package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder provides operations to manage the troubleshootingEvents property of the microsoft.graph.deviceManagement entity.
type TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderGetQueryParameters the list of troubleshooting events for the tenant.
type TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderGetQueryParameters
}
// TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewTroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderInternal instantiates a new TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder and sets the default values.
func NewTroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder) {
    m := &TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/troubleshootingEvents/{deviceManagementTroubleshootingEvent%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewTroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder instantiates a new TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder and sets the default values.
func NewTroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property troubleshootingEvents for deviceManagement
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the list of troubleshooting events for the tenant.
// returns a DeviceManagementTroubleshootingEventable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder) Get(ctx context.Context, requestConfiguration *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementTroubleshootingEventable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceManagementTroubleshootingEventFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementTroubleshootingEventable), nil
}
// Patch update the navigation property troubleshootingEvents in deviceManagement
// returns a DeviceManagementTroubleshootingEventable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementTroubleshootingEventable, requestConfiguration *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementTroubleshootingEventable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceManagementTroubleshootingEventFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementTroubleshootingEventable), nil
}
// ToDeleteRequestInformation delete navigation property troubleshootingEvents for deviceManagement
// returns a *RequestInformation when successful
func (m *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation the list of troubleshooting events for the tenant.
// returns a *RequestInformation when successful
func (m *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property troubleshootingEvents in deviceManagement
// returns a *RequestInformation when successful
func (m *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementTroubleshootingEventable, requestConfiguration *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder when successful
func (m *TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder) WithUrl(rawUrl string)(*TroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder) {
    return NewTroubleshootingEventsDeviceManagementTroubleshootingEventItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
