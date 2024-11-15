package employeeexperience

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// EmployeeExperienceRequestBuilder provides operations to manage the employeeExperience singleton.
type EmployeeExperienceRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// EmployeeExperienceRequestBuilderGetQueryParameters get employeeExperience
type EmployeeExperienceRequestBuilderGetQueryParameters struct {
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EmployeeExperienceRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EmployeeExperienceRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EmployeeExperienceRequestBuilderGetQueryParameters
}
// EmployeeExperienceRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EmployeeExperienceRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Communities provides operations to manage the communities property of the microsoft.graph.employeeExperience entity.
// returns a *CommunitiesRequestBuilder when successful
func (m *EmployeeExperienceRequestBuilder) Communities()(*CommunitiesRequestBuilder) {
    return NewCommunitiesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewEmployeeExperienceRequestBuilderInternal instantiates a new EmployeeExperienceRequestBuilder and sets the default values.
func NewEmployeeExperienceRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EmployeeExperienceRequestBuilder) {
    m := &EmployeeExperienceRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/employeeExperience{?%24select}", pathParameters),
    }
    return m
}
// NewEmployeeExperienceRequestBuilder instantiates a new EmployeeExperienceRequestBuilder and sets the default values.
func NewEmployeeExperienceRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EmployeeExperienceRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEmployeeExperienceRequestBuilderInternal(urlParams, requestAdapter)
}
// EngagementAsyncOperations provides operations to manage the engagementAsyncOperations property of the microsoft.graph.employeeExperience entity.
// returns a *EngagementAsyncOperationsRequestBuilder when successful
func (m *EmployeeExperienceRequestBuilder) EngagementAsyncOperations()(*EngagementAsyncOperationsRequestBuilder) {
    return NewEngagementAsyncOperationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get employeeExperience
// returns a EmployeeExperienceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EmployeeExperienceRequestBuilder) Get(ctx context.Context, requestConfiguration *EmployeeExperienceRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EmployeeExperienceable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEmployeeExperienceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EmployeeExperienceable), nil
}
// LearningCourseActivities provides operations to manage the learningCourseActivities property of the microsoft.graph.employeeExperience entity.
// returns a *LearningCourseActivitiesRequestBuilder when successful
func (m *EmployeeExperienceRequestBuilder) LearningCourseActivities()(*LearningCourseActivitiesRequestBuilder) {
    return NewLearningCourseActivitiesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LearningCourseActivitiesWithExternalcourseActivityId provides operations to manage the learningCourseActivities property of the microsoft.graph.employeeExperience entity.
// returns a *LearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder when successful
func (m *EmployeeExperienceRequestBuilder) LearningCourseActivitiesWithExternalcourseActivityId(externalcourseActivityId *string)(*LearningCourseActivitiesWithExternalcourseActivityIdRequestBuilder) {
    return NewLearningCourseActivitiesWithExternalcourseActivityIdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, externalcourseActivityId)
}
// LearningProviders provides operations to manage the learningProviders property of the microsoft.graph.employeeExperience entity.
// returns a *LearningProvidersRequestBuilder when successful
func (m *EmployeeExperienceRequestBuilder) LearningProviders()(*LearningProvidersRequestBuilder) {
    return NewLearningProvidersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update employeeExperience
// returns a EmployeeExperienceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EmployeeExperienceRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EmployeeExperienceable, requestConfiguration *EmployeeExperienceRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EmployeeExperienceable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEmployeeExperienceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EmployeeExperienceable), nil
}
// ToGetRequestInformation get employeeExperience
// returns a *RequestInformation when successful
func (m *EmployeeExperienceRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *EmployeeExperienceRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update employeeExperience
// returns a *RequestInformation when successful
func (m *EmployeeExperienceRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EmployeeExperienceable, requestConfiguration *EmployeeExperienceRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *EmployeeExperienceRequestBuilder when successful
func (m *EmployeeExperienceRequestBuilder) WithUrl(rawUrl string)(*EmployeeExperienceRequestBuilder) {
    return NewEmployeeExperienceRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
