package deviceappmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder provides operations to manage the deviceStates property of the microsoft.graph.userInstallStateSummary entity.
type ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderGetQueryParameters the install state of the eBook.
type ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderGetQueryParameters
}
// ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderInternal instantiates a new ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder and sets the default values.
func NewManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder) {
    m := &ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceAppManagement/managedEBooks/{managedEBook%2Did}/userStateSummary/{userInstallStateSummary%2Did}/deviceStates/{deviceInstallState%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder instantiates a new ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder and sets the default values.
func NewManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property deviceStates for deviceAppManagement
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the install state of the eBook.
// returns a DeviceInstallStateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceInstallStateable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceInstallStateFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceInstallStateable), nil
}
// Patch update the navigation property deviceStates in deviceAppManagement
// returns a DeviceInstallStateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceInstallStateable, requestConfiguration *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceInstallStateable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceInstallStateFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceInstallStateable), nil
}
// ToDeleteRequestInformation delete navigation property deviceStates for deviceAppManagement
// returns a *RequestInformation when successful
func (m *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation the install state of the eBook.
// returns a *RequestInformation when successful
func (m *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property deviceStates in deviceAppManagement
// returns a *RequestInformation when successful
func (m *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceInstallStateable, requestConfiguration *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder when successful
func (m *ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder) WithUrl(rawUrl string)(*ManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder) {
    return NewManagedEBooksItemUserStateSummaryItemDeviceStatesDeviceInstallStateItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
