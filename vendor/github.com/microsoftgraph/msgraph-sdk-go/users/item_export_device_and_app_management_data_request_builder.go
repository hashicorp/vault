package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemExportDeviceAndAppManagementDataRequestBuilder provides operations to call the exportDeviceAndAppManagementData method.
type ItemExportDeviceAndAppManagementDataRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemExportDeviceAndAppManagementDataRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemExportDeviceAndAppManagementDataRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemExportDeviceAndAppManagementDataRequestBuilderInternal instantiates a new ItemExportDeviceAndAppManagementDataRequestBuilder and sets the default values.
func NewItemExportDeviceAndAppManagementDataRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemExportDeviceAndAppManagementDataRequestBuilder) {
    m := &ItemExportDeviceAndAppManagementDataRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/exportDeviceAndAppManagementData()", pathParameters),
    }
    return m
}
// NewItemExportDeviceAndAppManagementDataRequestBuilder instantiates a new ItemExportDeviceAndAppManagementDataRequestBuilder and sets the default values.
func NewItemExportDeviceAndAppManagementDataRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemExportDeviceAndAppManagementDataRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemExportDeviceAndAppManagementDataRequestBuilderInternal(urlParams, requestAdapter)
}
// Get invoke function exportDeviceAndAppManagementData
// returns a DeviceAndAppManagementDataable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemExportDeviceAndAppManagementDataRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemExportDeviceAndAppManagementDataRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceAndAppManagementDataable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceAndAppManagementDataFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceAndAppManagementDataable), nil
}
// ToGetRequestInformation invoke function exportDeviceAndAppManagementData
// returns a *RequestInformation when successful
func (m *ItemExportDeviceAndAppManagementDataRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemExportDeviceAndAppManagementDataRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemExportDeviceAndAppManagementDataRequestBuilder when successful
func (m *ItemExportDeviceAndAppManagementDataRequestBuilder) WithUrl(rawUrl string)(*ItemExportDeviceAndAppManagementDataRequestBuilder) {
    return NewItemExportDeviceAndAppManagementDataRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
