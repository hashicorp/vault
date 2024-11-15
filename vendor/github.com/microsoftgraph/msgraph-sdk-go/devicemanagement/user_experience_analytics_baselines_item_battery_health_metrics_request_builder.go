package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder provides operations to manage the batteryHealthMetrics property of the microsoft.graph.userExperienceAnalyticsBaseline entity.
type UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilderGetQueryParameters the scores and insights for the battery health metrics.
type UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilderGetQueryParameters
}
// NewUserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilderInternal instantiates a new UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder and sets the default values.
func NewUserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder) {
    m := &UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/userExperienceAnalyticsBaselines/{userExperienceAnalyticsBaseline%2Did}/batteryHealthMetrics{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewUserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder instantiates a new UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder and sets the default values.
func NewUserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewUserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the scores and insights for the battery health metrics.
// returns a UserExperienceAnalyticsCategoryable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder) Get(ctx context.Context, requestConfiguration *UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserExperienceAnalyticsCategoryable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserExperienceAnalyticsCategoryFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserExperienceAnalyticsCategoryable), nil
}
// ToGetRequestInformation the scores and insights for the battery health metrics.
// returns a *RequestInformation when successful
func (m *UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder when successful
func (m *UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder) WithUrl(rawUrl string)(*UserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder) {
    return NewUserExperienceAnalyticsBaselinesItemBatteryHealthMetricsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
