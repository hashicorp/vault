package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder provides operations to manage the questions property of the microsoft.graph.virtualEventRegistrationConfiguration entity.
type VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderGetQueryParameters registration questions.
type VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderGetQueryParameters
}
// VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewVirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderInternal instantiates a new VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder and sets the default values.
func NewVirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder) {
    m := &VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/virtualEvents/webinars/{virtualEventWebinar%2Did}/registrationConfiguration/questions/{virtualEventRegistrationQuestionBase%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewVirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder instantiates a new VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder and sets the default values.
func NewVirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewVirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete a registration question from a webinar. The question can either be a predefined registration question or a custom registration question. 
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/virtualeventregistrationquestionbase-delete?view=graph-rest-1.0
func (m *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get registration questions.
// returns a VirtualEventRegistrationQuestionBaseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder) Get(ctx context.Context, requestConfiguration *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventRegistrationQuestionBaseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateVirtualEventRegistrationQuestionBaseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventRegistrationQuestionBaseable), nil
}
// Patch update the navigation property questions in solutions
// returns a VirtualEventRegistrationQuestionBaseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventRegistrationQuestionBaseable, requestConfiguration *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventRegistrationQuestionBaseable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateVirtualEventRegistrationQuestionBaseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventRegistrationQuestionBaseable), nil
}
// ToDeleteRequestInformation delete a registration question from a webinar. The question can either be a predefined registration question or a custom registration question. 
// returns a *RequestInformation when successful
func (m *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation registration questions.
// returns a *RequestInformation when successful
func (m *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property questions in solutions
// returns a *RequestInformation when successful
func (m *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventRegistrationQuestionBaseable, requestConfiguration *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder when successful
func (m *VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder) WithUrl(rawUrl string)(*VirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder) {
    return NewVirtualEventsWebinarsItemRegistrationConfigurationQuestionsVirtualEventRegistrationQuestionBaseItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
