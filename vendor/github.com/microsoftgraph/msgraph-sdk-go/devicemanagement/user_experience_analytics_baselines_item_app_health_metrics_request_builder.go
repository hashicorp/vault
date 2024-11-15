package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder provides operations to manage the appHealthMetrics property of the microsoft.graph.userExperienceAnalyticsBaseline entity.
type UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilderGetQueryParameters the scores and insights for the application health metrics.
type UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilderGetQueryParameters
}
// NewUserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilderInternal instantiates a new UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder and sets the default values.
func NewUserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder) {
    m := &UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/userExperienceAnalyticsBaselines/{userExperienceAnalyticsBaseline%2Did}/appHealthMetrics{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewUserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder instantiates a new UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder and sets the default values.
func NewUserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewUserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the scores and insights for the application health metrics.
// returns a UserExperienceAnalyticsCategoryable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder) Get(ctx context.Context, requestConfiguration *UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserExperienceAnalyticsCategoryable, error) {
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
// ToGetRequestInformation the scores and insights for the application health metrics.
// returns a *RequestInformation when successful
func (m *UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder when successful
func (m *UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder) WithUrl(rawUrl string)(*UserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder) {
    return NewUserExperienceAnalyticsBaselinesItemAppHealthMetricsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
