package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder provides operations to manage the metricDevices property of the microsoft.graph.userExperienceAnalyticsWorkFromAnywhereMetric entity.
type UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderGetQueryParameters the work from anywhere metric devices. Read-only.
type UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderGetQueryParameters
}
// UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByUserExperienceAnalyticsWorkFromAnywhereDeviceId provides operations to manage the metricDevices property of the microsoft.graph.userExperienceAnalyticsWorkFromAnywhereMetric entity.
// returns a *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesUserExperienceAnalyticsWorkFromAnywhereDeviceItemRequestBuilder when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder) ByUserExperienceAnalyticsWorkFromAnywhereDeviceId(userExperienceAnalyticsWorkFromAnywhereDeviceId string)(*UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesUserExperienceAnalyticsWorkFromAnywhereDeviceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if userExperienceAnalyticsWorkFromAnywhereDeviceId != "" {
        urlTplParams["userExperienceAnalyticsWorkFromAnywhereDevice%2Did"] = userExperienceAnalyticsWorkFromAnywhereDeviceId
    }
    return NewUserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesUserExperienceAnalyticsWorkFromAnywhereDeviceItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewUserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderInternal instantiates a new UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder and sets the default values.
func NewUserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder) {
    m := &UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/userExperienceAnalyticsWorkFromAnywhereMetrics/{userExperienceAnalyticsWorkFromAnywhereMetric%2Did}/metricDevices{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewUserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder instantiates a new UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder and sets the default values.
func NewUserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewUserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesCountRequestBuilder when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder) Count()(*UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesCountRequestBuilder) {
    return NewUserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get the work from anywhere metric devices. Read-only.
// returns a UserExperienceAnalyticsWorkFromAnywhereDeviceCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder) Get(ctx context.Context, requestConfiguration *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserExperienceAnalyticsWorkFromAnywhereDeviceCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserExperienceAnalyticsWorkFromAnywhereDeviceCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserExperienceAnalyticsWorkFromAnywhereDeviceCollectionResponseable), nil
}
// Post create new navigation property to metricDevices for deviceManagement
// returns a UserExperienceAnalyticsWorkFromAnywhereDeviceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserExperienceAnalyticsWorkFromAnywhereDeviceable, requestConfiguration *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserExperienceAnalyticsWorkFromAnywhereDeviceable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserExperienceAnalyticsWorkFromAnywhereDeviceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserExperienceAnalyticsWorkFromAnywhereDeviceable), nil
}
// ToGetRequestInformation the work from anywhere metric devices. Read-only.
// returns a *RequestInformation when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create new navigation property to metricDevices for deviceManagement
// returns a *RequestInformation when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserExperienceAnalyticsWorkFromAnywhereDeviceable, requestConfiguration *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder when successful
func (m *UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder) WithUrl(rawUrl string)(*UserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder) {
    return NewUserExperienceAnalyticsWorkFromAnywhereMetricsItemMetricDevicesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
