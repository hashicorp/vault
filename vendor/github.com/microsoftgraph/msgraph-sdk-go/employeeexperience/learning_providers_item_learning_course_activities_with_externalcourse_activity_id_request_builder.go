package employeeexperience

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder provides operations to manage the learningCourseActivities property of the microsoft.graph.learningProvider entity.
type LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderGetQueryParameters get learningCourseActivities from employeeExperience
type LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderGetQueryParameters
}
// LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewLearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderInternal instantiates a new LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder and sets the default values.
func NewLearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, externalcourseActivityId *string)(*LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder) {
    m := &LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/employeeExperience/learningProviders/{learningProvider%2Did}/learningCourseActivities(externalcourseActivityId='{externalcourseActivityId}'){?%24expand,%24select}", pathParameters),
    }
    if externalcourseActivityId != nil {
        m.BaseRequestBuilder.PathParameters["externalcourseActivityId"] = *externalcourseActivityId
    }
    return m
}
// NewLearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder instantiates a new LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder and sets the default values.
func NewLearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewLearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// Delete delete a learningCourseActivity object using the course activity ID of either an assignment or a self-initiated activity.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/learningcourseactivity-delete?view=graph-rest-1.0
func (m *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder) Delete(ctx context.Context, requestConfiguration *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get learningCourseActivities from employeeExperience
// returns a LearningCourseActivityable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder) Get(ctx context.Context, requestConfiguration *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LearningCourseActivityable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateLearningCourseActivityFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LearningCourseActivityable), nil
}
// Patch update the properties of a learningCourseActivity object. 
// returns a LearningCourseActivityable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/learningcourseactivity-update?view=graph-rest-1.0
func (m *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LearningCourseActivityable, requestConfiguration *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LearningCourseActivityable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateLearningCourseActivityFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LearningCourseActivityable), nil
}
// ToDeleteRequestInformation delete a learningCourseActivity object using the course activity ID of either an assignment or a self-initiated activity.
// returns a *RequestInformation when successful
func (m *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get learningCourseActivities from employeeExperience
// returns a *RequestInformation when successful
func (m *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of a learningCourseActivity object. 
// returns a *RequestInformation when successful
func (m *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LearningCourseActivityable, requestConfiguration *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder when successful
func (m *LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder) WithUrl(rawUrl string)(*LearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder) {
    return NewLearningProvidersItemLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
