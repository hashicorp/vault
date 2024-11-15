package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder provides operations to manage the resourcePerformanceMetrics property of the microsoft.graph.userExperienceAnalyticsBaseline entity.
type UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilderGetQueryParameters the scores and insights for the resource performance metrics.
type UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilderGetQueryParameters
}
// NewUserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilderInternal instantiates a new UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder and sets the default values.
func NewUserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder) {
    m := &UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/userExperienceAnalyticsBaselines/{userExperienceAnalyticsBaseline%2Did}/resourcePerformanceMetrics{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewUserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder instantiates a new UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder and sets the default values.
func NewUserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewUserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the scores and insights for the resource performance metrics.
// returns a UserExperienceAnalyticsCategoryable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder) Get(ctx context.Context, requestConfiguration *UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserExperienceAnalyticsCategoryable, error) {
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
// ToGetRequestInformation the scores and insights for the resource performance metrics.
// returns a *RequestInformation when successful
func (m *UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder when successful
func (m *UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder) WithUrl(rawUrl string)(*UserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder) {
    return NewUserExperienceAnalyticsBaselinesItemResourcePerformanceMetricsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
