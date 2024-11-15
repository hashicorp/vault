package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder provides operations to manage the webinars property of the microsoft.graph.virtualEventsRoot entity.
type VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderGetQueryParameters read the properties and relationships of a virtualEventWebinar object.
type VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderGetQueryParameters
}
// VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewVirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderInternal instantiates a new VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder and sets the default values.
func NewVirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) {
    m := &VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/virtualEvents/webinars/{virtualEventWebinar%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewVirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder instantiates a new VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder and sets the default values.
func NewVirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewVirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property webinars for solutions
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get read the properties and relationships of a virtualEventWebinar object.
// returns a VirtualEventWebinarable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/virtualeventwebinar-get?view=graph-rest-1.0
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) Get(ctx context.Context, requestConfiguration *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventWebinarable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateVirtualEventWebinarFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventWebinarable), nil
}
// Patch update the properties of a virtualEventWebinar object.
// returns a VirtualEventWebinarable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/virtualeventwebinar-update?view=graph-rest-1.0
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventWebinarable, requestConfiguration *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventWebinarable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateVirtualEventWebinarFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventWebinarable), nil
}
// Presenters provides operations to manage the presenters property of the microsoft.graph.virtualEvent entity.
// returns a *VirtualEventsWebinarsItemPresentersRequestBuilder when successful
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) Presenters()(*VirtualEventsWebinarsItemPresentersRequestBuilder) {
    return NewVirtualEventsWebinarsItemPresentersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RegistrationConfiguration provides operations to manage the registrationConfiguration property of the microsoft.graph.virtualEventWebinar entity.
// returns a *VirtualEventsWebinarsItemRegistrationConfigurationRequestBuilder when successful
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) RegistrationConfiguration()(*VirtualEventsWebinarsItemRegistrationConfigurationRequestBuilder) {
    return NewVirtualEventsWebinarsItemRegistrationConfigurationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Registrations provides operations to manage the registrations property of the microsoft.graph.virtualEventWebinar entity.
// returns a *VirtualEventsWebinarsItemRegistrationsRequestBuilder when successful
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) Registrations()(*VirtualEventsWebinarsItemRegistrationsRequestBuilder) {
    return NewVirtualEventsWebinarsItemRegistrationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RegistrationsWithEmail provides operations to manage the registrations property of the microsoft.graph.virtualEventWebinar entity.
// returns a *VirtualEventsWebinarsItemRegistrationsWithEmailRequestBuilder when successful
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) RegistrationsWithEmail(email *string)(*VirtualEventsWebinarsItemRegistrationsWithEmailRequestBuilder) {
    return NewVirtualEventsWebinarsItemRegistrationsWithEmailRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, email)
}
// RegistrationsWithUserId provides operations to manage the registrations property of the microsoft.graph.virtualEventWebinar entity.
// returns a *VirtualEventsWebinarsItemRegistrationsWithUserIdRequestBuilder when successful
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) RegistrationsWithUserId(userId *string)(*VirtualEventsWebinarsItemRegistrationsWithUserIdRequestBuilder) {
    return NewVirtualEventsWebinarsItemRegistrationsWithUserIdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, userId)
}
// Sessions provides operations to manage the sessions property of the microsoft.graph.virtualEvent entity.
// returns a *VirtualEventsWebinarsItemSessionsRequestBuilder when successful
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) Sessions()(*VirtualEventsWebinarsItemSessionsRequestBuilder) {
    return NewVirtualEventsWebinarsItemSessionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property webinars for solutions
// returns a *RequestInformation when successful
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read the properties and relationships of a virtualEventWebinar object.
// returns a *RequestInformation when successful
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of a virtualEventWebinar object.
// returns a *RequestInformation when successful
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventWebinarable, requestConfiguration *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder when successful
func (m *VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) WithUrl(rawUrl string)(*VirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder) {
    return NewVirtualEventsWebinarsVirtualEventWebinarItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
